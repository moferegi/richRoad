package shop

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	clientService "github.com/flipped-aurora/gin-vue-admin/server/service/client"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderService struct {
}

const (
	shopOrderStatusPendingConfirm = "8"

	orderCreateRateLimitConfigKey = "security_order_create_rate_limit_per_minute"
	orderCreateRateLimitEnvKey    = "CS_ORDER_CREATE_RATE_LIMIT_PER_MINUTE"

	orderPendingLimitConfigKey = "security_order_pending_limit_per_user"
	orderPendingLimitEnvKey    = "CS_ORDER_PENDING_LIMIT_PER_USER"

	orderPendingConfirmCloseMinutesConfigKey = "order_pending_confirm_close_minutes"
	orderPendingConfirmCloseMinutesEnvKey    = "CS_ORDER_PENDING_CONFIRM_CLOSE_MINUTES"

	defaultOrderCreateRateLimitPerMinute   int64 = 30
	defaultOrderPendingLimitPerUser        int64 = 10
	defaultOrderPendingConfirmCloseMinutes       = 180
)

func (orderService *OrderService) enforceOrderCreateRateLimit(userID uint) error {
	limit := utils.GetInt64Setting(orderCreateRateLimitConfigKey, orderCreateRateLimitEnvKey, defaultOrderCreateRateLimitPerMinute)
	if limit <= 0 {
		return nil
	}

	if global.GVA_REDIS != nil {
		ctx := context.Background()
		key := fmt.Sprintf("shop:order:create:rate:user:%d", userID)
		count, err := global.GVA_REDIS.Incr(ctx, key).Result()
		if err == nil {
			if count == 1 {
				_ = global.GVA_REDIS.Expire(ctx, key, time.Minute).Err()
			}
			if count > limit {
				return errors.New("orderCreateTooFrequent")
			}
			return nil
		}
		global.GVA_LOG.Warn("商城下单频率校验降级为DB", zap.Error(err), zap.Uint("userID", userID))
	}

	var recentCount int64
	if err := global.GVA_DB.Model(&shop.Order{}).
		Where("user_id = ? AND created_at >= ?", userID, time.Now().Add(-time.Minute)).
		Count(&recentCount).Error; err != nil {
		global.GVA_LOG.Warn("商城下单频率DB校验失败，已降级放行", zap.Error(err), zap.Uint("userID", userID))
		return nil
	}

	if recentCount >= limit {
		return errors.New("orderCreateTooFrequent")
	}

	return nil
}

func (orderService *OrderService) enforceOrderPendingLimit(tx *gorm.DB, userID uint) error {
	limit := utils.GetInt64Setting(orderPendingLimitConfigKey, orderPendingLimitEnvKey, defaultOrderPendingLimitPerUser)
	if limit <= 0 {
		return nil
	}

	var pendingCount int64
	err := tx.Model(&shop.Order{}).
		Where("user_id = ? AND status IN ? AND (close_time IS NULL OR close_time > ?)", userID, []string{"0", shopOrderStatusPendingConfirm}, time.Now()).
		Count(&pendingCount).Error
	if err != nil {
		return err
	}

	if pendingCount >= limit {
		return errors.New("orderPendingLimitExceeded")
	}

	return nil
}

func (orderService *OrderService) ensureOrderPendingAndNotExpired(order *shop.Order) error {
	if order == nil {
		return errors.New("orderNotFound")
	}
	if order.Status != "0" {
		return errors.New("orderStateInvalidForUpdate")
	}
	if !order.CloseTime.IsZero() && time.Now().After(order.CloseTime) {
		return errors.New("orderExpiredRecreate")
	}
	return nil
}

func (orderService *OrderService) getOrderPendingConfirmCloseMinutes(tx *gorm.DB) int {
	_ = tx
	minutes := utils.GetInt64Setting(orderPendingConfirmCloseMinutesConfigKey, orderPendingConfirmCloseMinutesEnvKey, int64(defaultOrderPendingConfirmCloseMinutes))
	if minutes <= 0 {
		return defaultOrderPendingConfirmCloseMinutes
	}
	return int(minutes)
}

func (orderService *OrderService) validateOrderStatusTransition(order *shop.Order, status string) error {
	if order == nil {
		return errors.New("orderNotFound")
	}
	if strings.TrimSpace(status) == "" {
		return errors.New("orderStatusInvalid")
	}
	if order.Status == status {
		return nil
	}

	switch status {
	case "1":
		if order.Status != "0" && order.Status != shopOrderStatusPendingConfirm {
			return errors.New("orderStateInvalidForConfirmPayment")
		}
	case "3":
		if order.Status != "2" {
			return errors.New("orderStateInvalidForConfirmReceive")
		}
	case "4":
		if order.Status != "0" && order.Status != shopOrderStatusPendingConfirm {
			return errors.New("orderStateInvalidForCancel")
		}
	case "5":
		if order.Status != "6" {
			return errors.New("orderStateInvalidForRefund")
		}
	case shopOrderStatusPendingConfirm:
		if order.Status != "0" {
			return errors.New("orderStateInvalidForSubmitPayment")
		}
		if strings.ToLower(strings.TrimSpace(order.PayMethod)) != "qrcode" {
			return errors.New("orderPayMethodNotQrcodeForSubmit")
		}
		if !order.CloseTime.IsZero() && time.Now().After(order.CloseTime) {
			return errors.New("orderExpiredRecreate")
		}
	default:
		return errors.New("orderStatusInvalid")
	}

	return nil
}

func normalizeSettlementCurrency(value string) string {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return ""
	}
	return strings.ReplaceAll(normalized, "_", "-")
}

func resolveCurrencySymbolByConfigValue(rawValue, settlementCurrency string) string {
	trimmed := strings.TrimSpace(rawValue)
	if trimmed == "" {
		return ""
	}

	if strings.HasPrefix(trimmed, "{") {
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(trimmed), &parsed); err == nil && len(parsed) > 0 {
			candidates := make([]string, 0, 8)
			currency := normalizeSettlementCurrency(settlementCurrency)
			if currency != "" {
				candidates = append(candidates, currency)
				if strings.Contains(currency, "-") {
					candidates = append(candidates, strings.SplitN(currency, "-", 2)[0])
				}
			}
			candidates = append(candidates, "zh", "en", "mn", "zh-TW")

			for _, key := range candidates {
				if key == "" {
					continue
				}
				if value, ok := parsed[key]; ok {
					resolved := strings.TrimSpace(fmt.Sprintf("%v", value))
					if resolved != "" {
						return resolved
					}
				}
			}

			for _, value := range parsed {
				resolved := strings.TrimSpace(fmt.Sprintf("%v", value))
				if resolved != "" {
					return resolved
				}
			}
		}
	}

	return trimmed
}

func (orderService *OrderService) fillOrderSettlementCurrency(tx *gorm.DB, order *shop.Order) {
	if order == nil {
		return
	}

	settlementCurrency := normalizeSettlementCurrency(order.SettlementCurrency)
	settlementSymbol := strings.TrimSpace(order.SettlementCurrencySymbol)

	if settlementSymbol == "" {
		var sysConf client.SysConfig
		if err := tx.Where("config_key = ?", "currency_symbol").First(&sysConf).Error; err == nil {
			settlementSymbol = resolveCurrencySymbolByConfigValue(sysConf.ConfigValue, settlementCurrency)
		}
	}

	if settlementCurrency == "" {
		settlementCurrency = "default"
	}
	if settlementSymbol == "" {
		settlementSymbol = "¥"
	}

	order.SettlementCurrency = settlementCurrency
	order.SettlementCurrencySymbol = settlementSymbol
}

// CreateOrder 创建订单记录
// Author [piexlmax](https://github.com/piexlmax)
func (orderService *OrderService) CreateOrder(order *shop.Order) (err error) {
	err = global.GVA_DB.Create(order).Error
	return err
}

func (orderService *OrderService) generateUniqueShopOutTradeNo(tx *gorm.DB) (string, error) {
	for i := 0; i < 10; i++ {
		orderNo, err := utils.GenerateBusinessOrderNo("sp")
		if err != nil {
			return "", err
		}

		var count int64
		if err := tx.Model(&shop.Order{}).Where("out_trade_no = ?", orderNo).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return orderNo, nil
		}
	}

	return "", errors.New("orderNoGenFail")
}

func (orderService *OrderService) ChangeOrderCoupon(userID uint, orderID string, couponNum string) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var order shop.Order
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? and user_id = ?", orderID, userID).
			Preload("Detail").
			First(&order).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("orderNotFound")
			}
			return err
		}
		if err = orderService.ensureOrderPendingAndNotExpired(&order); err != nil {
			return err
		}

		// 把原来的券解锁
		if order.CouponNum != "" {
			// 解锁原来的券
			err = tx.Model(&shop.CouponOrderUser{}).Where("coupon_num = ?", order.CouponNum).Update("order_id", nil).Error
			if err != nil {
				return err
			}
		}

		// 如果 couponNum 为空，表示取消优惠券
		if couponNum == "" {
			order.TotalPrice = int(order.OriginPrice)
			order.CouponNum = ""
			order.Discount = 0
			err = tx.Model(&order).
				Update("total_price", order.TotalPrice).
				Update("coupon_num", "").
				Update("discount", 0).
				Error
			return err
		}

		var couponOrderUser shop.CouponOrderUser
		shopUserID := int(userID)
		err = tx.Where("coupon_num = ? AND order_id IS NULL AND shop_user_id = ?", couponNum, shopUserID).First(&couponOrderUser).Error
		if err != nil {
			return errors.New("orderCouponUnavailable")
		}
		var coupon shop.Coupon
		err = tx.Where("id = ?", couponOrderUser.CouponID).First(&coupon).Error
		if err != nil {
			return err
		}
		// 校验优惠券是否过期
		if coupon.EndTime != nil && time.Now().After(*coupon.EndTime) {
			return errors.New("orderCouponExpired")
		}
		if coupon.ProductID != nil || coupon.ProductIDs != "" {
			allowedIDs := couponAllowedProductIDs(coupon)
			if len(allowedIDs) > 0 {
				hasGoods := false
				for _, detail := range order.Detail {
					if allowedIDs[int(detail.GoodID)] {
						hasGoods = true
						break
					}
				}
				if !hasGoods {
					return errors.New("orderCouponNotApplicable")
				}
			}
		}
		order.TotalPrice = 0
		if order.OriginPrice > coupon.Discount {
			order.TotalPrice = int(order.OriginPrice - coupon.Discount)
		}
		order.CouponNum = couponNum
		order.Discount = coupon.Discount
		err = tx.Model(&order).
			Update("total_price", order.TotalPrice).
			Update("coupon_num", couponNum).
			Update("discount", coupon.Discount).
			Error
		if err != nil {
			return err
		}
		// 使用新的券
		used := true
		now := time.Now()
		err = tx.Model(&shop.CouponOrderUser{}).Where("coupon_num = ?", couponNum).Updates(map[string]interface{}{
			"order_id": order.ID,
			"status":   &used,
			"used_at":  &now,
		}).Error
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// ChangeOrderPoints 变更订单积分抵扣
func (orderService *OrderService) ChangeOrderPoints(userID uint, orderID string, usePoints bool) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var order shop.Order
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? and user_id = ?", orderID, userID).
			Preload("Detail").
			First(&order).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("orderNotFound")
			}
			return err
		}
		if err = orderService.ensureOrderPendingAndNotExpired(&order); err != nil {
			return err
		}

		// 1. 如果之前使用了积分，需要先恢复积分
		if order.UsePoints && order.PointsUsed > 0 {
			pointRecordService := clientService.PointRecordService{}
			userID := int(order.UserID)
			assetType := client.AssetTypePoint
			pointChange := int(order.PointsUsed) // 恢复积分，正数
			changeType := "increase"
			operationType := "point_exchange"
			reason := "取消积分抵扣，恢复积分"
			relatedOrderId := int(order.ID)

			pointRecord := client.PointRecord{
				AssetType:      &assetType,
				UserId:         &userID,
				PointChange:    &pointChange,
				ChangeType:     &changeType,
				OperationType:  &operationType,
				Reason:         &reason,
				RelatedOrderId: &relatedOrderId,
			}
			ctxWithTx := context.WithValue(context.Background(), "tx", tx)
			err = pointRecordService.CreatePointRecord(ctxWithTx, &pointRecord)
			if err != nil {
				return err
			}
		}

		// 2. 重新计算订单价格
		order.TotalPrice = int(order.OriginPrice - order.Discount)
		order.UsePoints = usePoints
		order.PointsUsed = 0

		// 3. 如果现在要使用积分，进行积分扣除
		if usePoints {
			// 3a. 校验每个商品是否允许积分抵扣，并计算总的最大可用积分
			var goodMaxPoints int // 商品层面允许的最大积分总和
			allGoodsAllow := true
			for _, detail := range order.Detail {
				var good shop.Good
				if e := tx.Where("id = ?", detail.GoodID).First(&good).Error; e != nil {
					return errors.New("orderGoodInfoQueryFailed")
				}
				if good.PointsEnabled == nil || !*good.PointsEnabled {
					allGoodsAllow = false
					break
				}
				// 检查同商品使用积分次数限制
				if good.PointsUseTimes != nil && *good.PointsUseTimes > 0 {
					var usedCount int64
					tx.Model(&shop.Order{}).
						Where("user_id = ? AND id != ? AND use_points = ? AND status NOT IN ?", order.UserID, order.ID, true, []string{"4", "5"}).
						Joins("JOIN shop_order_detail ON shop_order_detail.order_id = shop_order.id").
						Where("shop_order_detail.good_id = ?", detail.GoodID).
						Count(&usedCount)
					if int(usedCount) >= *good.PointsUseTimes {
						return errors.New("orderPointsUseTimesLimitReached")
					}
				}
				// 累加每个商品的最大可用积分（按数量乘）
				if good.PointsMaxUse != nil && *good.PointsMaxUse > 0 {
					goodMaxPoints += *good.PointsMaxUse * int(detail.Quantity)
				}
			}
			if !allGoodsAllow {
				return errors.New("orderContainsGoodsNotSupportPoints")
			}

			var user client.ClientUser
			err = tx.Where("id = ?", order.UserID).First(&user).Error
			if err != nil {
				return err
			}
			// 3b. 计算可抵扣的积分数：min(用户积分, 订单金额, 商品最大可用积分)
			maxPointsCanUse := order.TotalPrice
			if goodMaxPoints > 0 && goodMaxPoints < maxPointsCanUse {
				maxPointsCanUse = goodMaxPoints
			}
			pointsToUse := user.Point
			if pointsToUse > maxPointsCanUse {
				pointsToUse = maxPointsCanUse
			}
			if pointsToUse > 0 {
				pointRecordService := clientService.PointRecordService{}
				userID := int(order.UserID)
				assetType := client.AssetTypePoint
				pointChange := -int(pointsToUse)
				changeType := "decrease"
				operationType := "point_exchange"
				reason := "订单积分抵扣"
				relatedOrderId := int(order.ID)

				pointRecord := client.PointRecord{
					AssetType:      &assetType,
					UserId:         &userID,
					PointChange:    &pointChange,
					ChangeType:     &changeType,
					OperationType:  &operationType,
					Reason:         &reason,
					RelatedOrderId: &relatedOrderId,
				}
				ctxWithTx := context.WithValue(context.Background(), "tx", tx)
				err = pointRecordService.CreatePointRecord(ctxWithTx, &pointRecord)
				if err != nil {
					return err
				}
				order.TotalPrice -= pointsToUse
				order.PointsUsed = uint(pointsToUse)
			}
		}

		// 4. 确保订单金额不为负数
		if order.TotalPrice < 0 {
			order.TotalPrice = 0
		}

		// 5. 更新订单信息
		err = tx.Model(&order).Updates(map[string]interface{}{
			"total_price": order.TotalPrice,
			"use_points":  order.UsePoints,
			"points_used": order.PointsUsed,
		}).Error
		if err != nil {
			return err
		}
		return nil
	})
	return
}

func (orderService *OrderService) PlaceOrder(order *shop.Order) (OrderID uint, orderNo string, err error) {
	if order == nil {
		return 0, "", errors.New("invalidOrder")
	}
	if order.UserID == 0 {
		return 0, "", errors.New("loginRequired")
	}
	if err = orderService.enforceOrderCreateRateLimit(order.UserID); err != nil {
		return 0, "", err
	}

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err = orderService.enforceOrderPendingLimit(tx, order.UserID); err != nil {
			return err
		}

		order.OriginPrice = 0
		presaleService := PresaleService{}
		for i := range order.Detail {
			// 判断库存
			var sku shop.Sku
			err = tx.Where("id = ?", order.Detail[i].SKUID).First(&sku).Error
			if err != nil {
				return err
			}
			if sku.Inventory < order.Detail[i].Quantity {
				return errors.New("orderInventoryInsufficient")
			}

			// 预售商品检查与防超卖
			var good shop.Good
			err = tx.Where("id = ?", sku.GoodID).First(&good).Error
			if err != nil {
				return err
			}
			if good.IsPresale != nil && *good.IsPresale {
				available, message, checkErr := presaleService.CheckPresaleAvailable(sku.GoodID)
				if checkErr != nil {
					return checkErr
				}
				if !available {
					return errors.New(message)
				}
				if err = presaleService.IncrementPresaleSold(tx, sku.GoodID, int(order.Detail[i].Quantity)); err != nil {
					return err
				}
				order.IsPresale = true
			}
			order.Detail[i].GoodID = sku.GoodID

			// 乐观锁扣减库存：WHERE inventory >= quantity 防止并发超卖
			result := tx.Model(&shop.Sku{}).Where("id = ? AND inventory >= ?", order.Detail[i].SKUID, order.Detail[i].Quantity).
				Update("inventory", gorm.Expr("inventory - ?", order.Detail[i].Quantity))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New("orderInventoryInsufficientRetry")
			}
			// sale_num 在付款确认时增加，而非下单时
			order.Detail[i].Price = sku.Price
			order.Detail[i].PriceI18n = sku.PriceI18n
			order.OriginPrice += order.Detail[i].Quantity * order.Detail[i].Price
			// 减扣库存
		}
		order.TotalPrice = int(order.OriginPrice)
		if order.CouponNum != "" {
			var couponOrderUser shop.CouponOrderUser
			shopUserID := int(order.UserID)
			err = tx.Where("coupon_num = ? AND order_id IS NULL AND shop_user_id = ?", order.CouponNum, shopUserID).First(&couponOrderUser).Error
			if err != nil {
				return errors.New("orderCouponUnavailable")
			}
			var coupon shop.Coupon
			err = tx.Where("id = ?", couponOrderUser.CouponID).First(&coupon).Error
			if err != nil {
				return err
			}
			// 校验优惠券是否过期
			if coupon.EndTime != nil && time.Now().After(*coupon.EndTime) {
				return errors.New("orderCouponExpired")
			}
			if coupon.ProductID != nil {
				hasGoods := false
				for _, detail := range order.Detail {
					if int(detail.ID) == *coupon.ProductID {
						hasGoods = true
						break
					}
				}
				if !hasGoods {
					return errors.New("orderCouponNotApplicable")
				}
			}
			order.TotalPrice = 0
			if order.OriginPrice > coupon.Discount {
				order.TotalPrice = int(order.OriginPrice - coupon.Discount)
			}
			order.Discount = coupon.Discount
		}

		// 如果订单金额为0或负数，确保设置为0（0元购）
		if order.TotalPrice <= 0 {
			order.TotalPrice = 0
		}
		order.Status = "0"

		// 只有在前端未提供地址信息时，才使用默认地址
		if order.Name == "" && order.Phone == "" {
			if addr, err := orderService.GetDefaultAddress(order.UserID); err == nil {
				order.Name = addr.Name
				order.Phone = addr.Phone
				order.Province = addr.ProvinceStr
				order.City = addr.CityStr
				order.Area = addr.AreaStr
				order.Street = addr.Street
			}
		}
		// 从系统配置读取订单自动关闭时间（分钟），默认15分钟
		closeMinutes := 15
		var sysConf client.SysConfig
		if e := tx.Where("config_key = ?", "order_close_minutes").First(&sysConf).Error; e == nil {
			if v, pe := strconv.Atoi(sysConf.ConfigValue); pe == nil && v > 0 {
				closeMinutes = v
			}
		}
		order.CloseTime = time.Now().Add(time.Duration(closeMinutes) * time.Minute)
		orderService.fillOrderSettlementCurrency(tx, order)

		generatedOrderNo, noErr := orderService.generateUniqueShopOutTradeNo(tx)
		if noErr != nil {
			return noErr
		}
		order.OutTradeNo = generatedOrderNo

		err = tx.Create(order).Error
		if err != nil {
			return err
		}
		OrderID = order.ID
		orderNo = order.OutTradeNo
		used := true
		now := time.Now()
		err = tx.Model(&shop.CouponOrderUser{}).Where("coupon_num = ?", order.CouponNum).Updates(map[string]interface{}{
			"order_id": order.ID,
			"status":   &used,
			"used_at":  &now,
		}).Error
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// 购物车下单
func (orderService *OrderService) PlaceOrderByCart(userID uint, req shopReq.PlaceOrderByCartRequest) (OrderID uint, orderNo string, err error) {
	if userID == 0 {
		return 0, "", errors.New("loginRequired")
	}
	if err = orderService.enforceOrderCreateRateLimit(userID); err != nil {
		return 0, "", err
	}

	// 1. 获取购物车中的商品
	var carts []shop.Cart
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err = orderService.enforceOrderPendingLimit(tx, userID); err != nil {
			return err
		}

		query := tx.Where("user_id = ?", userID).Preload("Good").Preload("SKU")
		if len(req.CartIDs) > 0 {
			query = query.Where("id IN ?", req.CartIDs)
		}
		query.Find(&carts)
		if len(carts) == 0 {
			return errors.New("orderCartEmpty")
		}
		// 2. 创建订单
		order := shop.Order{
			UserID:                   userID,
			Status:                   "0",
			PayMethod:                req.PayMethod,
			SettlementCurrency:       req.SettlementCurrency,
			SettlementCurrencySymbol: req.SettlementCurrencySymbol,
		}
		order.OriginPrice = 0
		// 从购物车设置订单详情
		order.Detail = make([]shop.OrderDetail, 0)
		presaleService := PresaleService{}
		for i := range carts {
			// 判断当前库存是否充足
			if carts[i].SKU.Inventory < carts[i].Quantity {
				return errors.New("orderInventoryInsufficient")
			}

			// 预售商品检查与防超卖
			if carts[i].Good.IsPresale != nil && *carts[i].Good.IsPresale {
				available, message, checkErr := presaleService.CheckPresaleAvailable(carts[i].GoodID)
				if checkErr != nil {
					return checkErr
				}
				if !available {
					return errors.New(message)
				}
				if err = presaleService.IncrementPresaleSold(tx, carts[i].GoodID, int(carts[i].Quantity)); err != nil {
					return err
				}
				order.IsPresale = true
			}

			order.Detail = append(order.Detail, shop.OrderDetail{
				GoodID:    carts[i].GoodID,
				SKUID:     carts[i].SKUID,
				Quantity:  carts[i].Quantity,
				Price:     carts[i].SKU.Price,
				PriceI18n: carts[i].SKU.PriceI18n,
			})
			order.OriginPrice += carts[i].Quantity * carts[i].SKU.Price
			// 乐观锁扣减库存：WHERE inventory >= quantity 防止并发超卖
			result := tx.Model(&shop.Sku{}).Where("id = ? AND inventory >= ?", carts[i].SKUID, carts[i].Quantity).
				Update("inventory", gorm.Expr("inventory - ?", carts[i].Quantity))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New("orderInventoryInsufficientRetry")
			}
			// sale_num 在付款确认时增加，而非下单时
		}

		// 只有在前端未提供地址信息时，才使用默认地址
		if order.Name == "" && order.Phone == "" {
			if addr, err := orderService.GetDefaultAddress(order.UserID); err == nil {
				order.Name = addr.Name
				order.Phone = addr.Phone
				order.Province = addr.ProvinceStr
				order.City = addr.CityStr
				order.Area = addr.AreaStr
				order.Street = addr.Street
			}
		}
		// 从系统配置读取订单自动关闭时间（分钟），默认15分钟
		closeMinutes := 15
		var sysConf client.SysConfig
		if e := tx.Where("config_key = ?", "order_close_minutes").First(&sysConf).Error; e == nil {
			if v, pe := strconv.Atoi(sysConf.ConfigValue); pe == nil && v > 0 {
				closeMinutes = v
			}
		}
		order.CloseTime = time.Now().Add(time.Duration(closeMinutes) * time.Minute)

		// 3. 创建订单详情

		order.TotalPrice = int(order.OriginPrice)

		// 如果订单金额为0或负数，确保设置为0（0元购）
		if order.TotalPrice <= 0 {
			order.TotalPrice = 0
		}
		orderService.fillOrderSettlementCurrency(tx, &order)

		generatedOrderNo, noErr := orderService.generateUniqueShopOutTradeNo(tx)
		if noErr != nil {
			return noErr
		}
		order.OutTradeNo = generatedOrderNo

		err = tx.Create(&order).Error
		if err != nil {
			return err
		}

		// 4. 删除已下单的购物车商品
		cartIDs := make([]uint, len(carts))
		for i, c := range carts {
			cartIDs[i] = c.ID
		}
		err = tx.Where("id IN ?", cartIDs).Delete(&shop.Cart{}).Error
		if err != nil {
			return err
		}
		OrderID = order.ID
		orderNo = order.OutTradeNo
		return nil
	})
	return
}

func (orderService *OrderService) UpdateOrderStatus(db *gorm.DB, orderID string, status string) (err error) {
	if db == nil {
		db = global.GVA_DB
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		var order shop.Order
		var user client.ClientUser
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, "id = ?", orderID).Error
		if err != nil {
			return err
		}

		if err = orderService.validateOrderStatusTransition(&order, status); err != nil {
			return err
		}
		if order.Status == status {
			return nil
		}

		updates := map[string]interface{}{"status": status}
		if status == shopOrderStatusPendingConfirm {
			pendingConfirmCloseMinutes := orderService.getOrderPendingConfirmCloseMinutes(tx)
			if pendingConfirmCloseMinutes > 0 {
				updates["close_time"] = time.Now().Add(time.Duration(pendingConfirmCloseMinutes) * time.Minute)
			}
		}
		err = tx.Model(&order).Updates(updates).Error
		if err != nil {
			return err
		}
		err = tx.Where("id = ?", order.UserID).First(&user).Error
		if err != nil {
			return err
		}
		totalPrice := order.TotalPrice
		// 订单取消（status="4"用户取消 或 status="5"系统/退款取消）时，恢复库存
		if status == "4" || status == "5" {
			err = tx.Preload("Detail").First(&order, "id = ?", orderID).Error
			if err != nil {
				return err
			}
			// 恢复SKU库存
			for _, detail := range order.Detail {
				err = tx.Model(&shop.Sku{}).Where("id = ?", detail.SKUID).Update("inventory", gorm.Expr("inventory + ?", detail.Quantity)).Error
				if err != nil {
					return err
				}
			}
			// 退款（status="5"）时恢复销量（销量在付款确认时增加的）
			if status == "5" {
				for _, detail := range order.Detail {
					tx.Model(&shop.Sku{}).Where("id = ? AND sale_num >= ?", detail.SKUID, detail.Quantity).
						Update("sale_num", gorm.Expr("sale_num - ?", detail.Quantity))
					tx.Model(&shop.Good{}).Where("id = ? AND sale_num >= ?", detail.GoodID, detail.Quantity).
						Update("sale_num", gorm.Expr("sale_num - ?", detail.Quantity))
				}
			}
			// 如果是预售商品，恢复预售已售数量
			if order.IsPresale {
				for _, detail := range order.Detail {
					tx.Model(&shop.Good{}).Where("id = ? AND presale_sold > 0", detail.GoodID).
						Update("presale_sold", gorm.Expr("presale_sold - ?", detail.Quantity))
				}
			}
			// 设置取消时间
			now := time.Now()
			tx.Model(&order).Update("cancelled_at", now)

			// 恢复优惠券：清除 order_id 绑定，如果订单使用了优惠券
			if order.CouponNum != "" {
				// 查到该券记录
				var couponRecord shop.CouponOrderUser
				if e := tx.Where("coupon_num = ?", order.CouponNum).First(&couponRecord).Error; e == nil {
					// 清除订单绑定，重置状态为未使用
					unused := false
					tx.Model(&couponRecord).Updates(map[string]interface{}{
						"order_id": nil,
						"status":   &unused,
					})
				}
			}

			// 如果订单未付款但使用了积分，也需要返还积分
			if !order.Pointed && order.UsePoints && order.PointsUsed > 0 {
				pointRecordService := &clientService.PointRecordService{}
				userIdInt := int(user.ID)
				assetType := client.AssetTypePoint
				changeType := "increase"
				pointChange := int(order.PointsUsed)
				operationType := "point_refund"
				reason := "订单取消，返还已使用积分"
				orderIdStr := strconv.Itoa(int(order.ID))
				orderIdInt, _ := strconv.Atoi(orderIdStr)
				remark := "订单ID: " + orderIdStr

				pointRecord := &client.PointRecord{
					AssetType:      &assetType,
					UserId:         &userIdInt,
					ChangeType:     &changeType,
					PointChange:    &pointChange,
					OperationType:  &operationType,
					Reason:         &reason,
					RelatedOrderId: &orderIdInt,
					Remark:         &remark,
				}

				ctxWithTx := context.WithValue(context.Background(), "tx", tx)
				err = pointRecordService.CreatePointRecord(ctxWithTx, pointRecord)
				if err != nil {
					return err
				}
				// 清除订单积分使用标记
				err = tx.Model(&order).Updates(map[string]interface{}{
					"use_points":  false,
					"points_used": 0,
				}).Error
				if err != nil {
					return err
				}
			}

			// 如果订单已积分（Pointed=true），还需要处理积分返还和扣除
			if order.Pointed {
				// 1. 如果订单使用了积分，先返还已使用的积分
				if order.UsePoints && order.PointsUsed > 0 {
					pointRecordService := &clientService.PointRecordService{}
					userIdInt := int(user.ID)
					assetType := client.AssetTypePoint
					changeType := "increase"
					pointChange := int(order.PointsUsed) // 返还积分，正数
					operationType := "point_refund"
					reason := "订单取消，返还已使用积分"
					orderIdStr := strconv.Itoa(int(order.ID))
					orderIdInt, _ := strconv.Atoi(orderIdStr)
					remark := "订单ID: " + orderIdStr

					pointRecord := &client.PointRecord{
						AssetType:      &assetType,
						UserId:         &userIdInt,
						ChangeType:     &changeType,
						PointChange:    &pointChange,
						OperationType:  &operationType,
						Reason:         &reason,
						RelatedOrderId: &orderIdInt,
						Remark:         &remark,
					}

					ctxWithTx := context.WithValue(context.Background(), "tx", tx)
					err = pointRecordService.CreatePointRecord(ctxWithTx, pointRecord)
					if err != nil {
						return err
					}
				}

				// 2. 扣除已获得的积分奖励
				pointRecordService := &clientService.PointRecordService{}
				userIdInt := int(user.ID)
				assetType := client.AssetTypePoint
				changeType := "decrease"
				pointChange := int(totalPrice / 100)
				operationType := "refund_return"
				reason := "订单取消，扣除已获得积分"
				orderIdStr := strconv.Itoa(int(order.ID))
				orderIdInt, _ := strconv.Atoi(orderIdStr)
				remark := "订单ID: " + orderIdStr

				pointRecord := &client.PointRecord{
					AssetType:      &assetType,
					UserId:         &userIdInt,
					ChangeType:     &changeType,
					PointChange:    &pointChange,
					OperationType:  &operationType,
					Reason:         &reason,
					RelatedOrderId: &orderIdInt,
					Remark:         &remark,
				}

				ctxWithTx := context.WithValue(context.Background(), "tx", tx)
				err = pointRecordService.CreatePointRecord(ctxWithTx, pointRecord)
				if err != nil {
					return err
				}
				err = tx.Model(&order).Update("pointed", false).Error
				if err != nil {
					return err
				}
			}
			return nil
		}

		if status == "1" && !order.Pointed {
			// 付款确认后增加销量
			err = tx.Preload("Detail").First(&order, "id = ?", orderID).Error
			if err != nil {
				return err
			}
			for _, detail := range order.Detail {
				tx.Model(&shop.Sku{}).Where("id = ?", detail.SKUID).
					Update("sale_num", gorm.Expr("sale_num + ?", detail.Quantity))
				tx.Model(&shop.Good{}).Where("id = ?", detail.GoodID).
					Update("sale_num", gorm.Expr("sale_num + ?", detail.Quantity))
			}

			// 幂等检查：查看该订单是否已经发放过奖励（防止取消后重新支付重复发放）
			var rewardCount int64
			tx.Model(&client.PointRecord{}).
				Where("related_order_id = ? AND operation_type = ? AND (asset_type = ? OR asset_type IS NULL)", order.ID, "order_complete", client.AssetTypePoint).
				Count(&rewardCount)
			if rewardCount == 0 {
				// 使用营销奖励配置发放订单奖励（积分+优惠券）
				mrService := &MarketingRewardService{}
				if err := mrService.TriggerReward(order.UserID, "order", "order_complete", "订单完成，获得奖励", order.ID); err != nil {
					global.GVA_LOG.Error("订单营销奖励发放失败", zap.Error(err))
				}

				// 检查下级首单奖励逻辑
				mrService.TriggerSubOrderReward(tx, order.UserID, order.ID)
			}

			err = tx.Model(&order).Update("pointed", true).Error
			if err != nil {
				return err
			}
		}
		return err
	})
	return err
}

// UpdateOrderStatusForUser 普通用户更新订单状态（仅本人订单，且严格限制状态迁移）
func (orderService *OrderService) UpdateOrderStatusForUser(userID uint, orderID string, status string) error {
	if userID == 0 {
		return errors.New("loginRequired")
	}
	if strings.TrimSpace(orderID) == "" {
		return errors.New("orderIDRequired")
	}
	if status != "3" && status != "4" && status != shopOrderStatusPendingConfirm {
		return errors.New("orderStatusInvalid")
	}

	var order shop.Order
	err := global.GVA_DB.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("orderNotFound")
		}
		return err
	}

	if status == "4" && order.Status != "0" {
		return errors.New("orderStateInvalidForCancel")
	}
	if status == "3" && order.Status != "2" {
		return errors.New("orderStateInvalidForConfirmReceive")
	}
	if status == shopOrderStatusPendingConfirm {
		if order.Status == shopOrderStatusPendingConfirm {
			return nil
		}
		if order.Status != "0" {
			return errors.New("orderStateInvalidForSubmitPayment")
		}
		if !order.CloseTime.IsZero() && time.Now().After(order.CloseTime) {
			return errors.New("orderExpiredRecreate")
		}
		if strings.ToLower(strings.TrimSpace(order.PayMethod)) != "qrcode" {
			return errors.New("orderPayMethodNotQrcodeForSubmit")
		}
	}

	return orderService.UpdateOrderStatus(nil, orderID, status)
}

// DeleteOrder 删除订单记录
// Author [piexlmax](https://github.com/piexlmax)
func (orderService *OrderService) DeleteOrder(ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.Order{}, "id = ?", ID).Error
	return err
}

// ConfirmPayment 管理员确认收款
func (orderService *OrderService) ConfirmPayment(orderID string) (err error) {
	var order shop.Order
	if err = global.GVA_DB.First(&order, "id = ?", orderID).Error; err != nil {
		return err
	}
	if order.Status != "0" && order.Status != shopOrderStatusPendingConfirm {
		return errors.New("orderConfirmPaymentStateInvalid")
	}
	now := time.Now()
	// 先设置 paidAt
	if err = global.GVA_DB.Model(&order).Update("paid_at", now).Error; err != nil {
		return err
	}
	// 调用 UpdateOrderStatus 处理状态变更和积分逻辑
	return orderService.UpdateOrderStatus(nil, orderID, "1")
}

// DeleteOrderByIds 批量删除订单记录
// Author [piexlmax](https://github.com/piexlmax)
func (orderService *OrderService) DeleteOrderByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]shop.Order{}, "id in ?", IDs).Error
	return err
}

// UpdateOrder 更新订单记录
// Author [piexlmax](https://github.com/piexlmax)
func (orderService *OrderService) UpdateOrder(order shop.Order) (err error) {
	err = global.GVA_DB.Model(&shop.Order{}).Where("id = ?", order.ID).Updates(&order).Error
	return err
}

// UpdateOrderForUser 普通用户更新订单（仅允许本人待支付订单修改地址和支付方式）
func (orderService *OrderService) UpdateOrderForUser(userID uint, order shop.Order) (err error) {
	if userID == 0 {
		return errors.New("loginRequired")
	}
	if order.ID == 0 {
		return errors.New("invalidOrder")
	}

	var existing shop.Order
	err = global.GVA_DB.Where("id = ? AND user_id = ?", order.ID, userID).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("orderNotFound")
		}
		return err
	}

	if existing.Status != "0" {
		return errors.New("orderStateInvalidForUpdate")
	}

	updates := map[string]interface{}{}
	if trimmed := strings.TrimSpace(order.PayMethod); trimmed != "" {
		updates["pay_method"] = trimmed
	}
	if trimmed := strings.TrimSpace(order.Name); trimmed != "" {
		updates["name"] = trimmed
	}
	if trimmed := strings.TrimSpace(order.Phone); trimmed != "" {
		updates["phone"] = trimmed
	}
	if trimmed := strings.TrimSpace(order.Province); trimmed != "" {
		updates["province"] = trimmed
	}
	if trimmed := strings.TrimSpace(order.City); trimmed != "" {
		updates["city"] = trimmed
	}
	if trimmed := strings.TrimSpace(order.Area); trimmed != "" {
		updates["area"] = trimmed
	}
	if trimmed := strings.TrimSpace(order.Street); trimmed != "" {
		updates["street"] = trimmed
	}

	if len(updates) == 0 {
		return nil
	}

	err = global.GVA_DB.Model(&shop.Order{}).Where("id = ? AND user_id = ?", order.ID, userID).Updates(updates).Error
	return err
}

// GetOrder 根据ID获取订单记录
// Author [piexlmax](https://github.com/piexlmax)
func (orderService *OrderService) GetOrder(ID string, userID uint) (order shop.OrderRes, err error) {
	db := global.GVA_DB.Where("id = ?", ID).Preload("Detail").Preload("Detail.Good").Preload("Detail.SKU").Preload("Comment")
	if userID != 0 {
		db = db.Where("user_id = ?", userID)
	}
	err = db.First(&order).Error
	return
}

func (orderService *OrderService) SelfOrderComment(ID, goodID, SKUID string, userID uint) (order shop.OrderCommentRes, err error) {
	var comment shop.Comment
	var detail shop.OrderDetailRes
	_ = global.GVA_DB.First(&comment, "order_id = ? && good_id = ? && sku_id = ?", ID, goodID, SKUID).Error

	err = global.GVA_DB.Preload("Good").Preload("SKU").First(&detail, "order_id = ? && good_id = ? && sku_id = ?", ID, goodID, SKUID).Error
	if err != nil {
		return order, err
	}
	db := global.GVA_DB.Where("id = ?", ID)
	if userID != 0 {
		db = db.Where("user_id = ?", userID)
	}
	err = db.First(&order).Error
	order.Comment = comment
	order.Detail = detail
	return
}

// GetOrderInfoList 分页获取订单记录
// Author [piexlmax](https://github.com/piexlmax)
func (orderService *OrderService) GetOrderInfoList(info shopReq.OrderSearch) (list []shop.OrderRes, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	// 创建db
	db := global.GVA_DB.Model(&shop.OrderRes{}).Preload("Detail").Preload("Detail.Good").Preload("Detail.SKU")
	var orders []shop.OrderRes

	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.UserID != nil {
		db = db.Where("user_id = ?", info.UserID)
	}

	// 根据info.status查询
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}

	cErr := db.Count(&total).Error
	if cErr != nil {
		return
	}
	// 应用分页
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}

	// 按排序字段排序（白名单限制）
	orderClause := "created_at desc"
	orderMap := map[string]bool{
		"id":          true,
		"created_at":  true,
		"total_price": true,
		"status":      true,
	}
	if orderMap[info.Sort] {
		orderClause = info.Sort
		if info.Order == "descending" {
			orderClause += " desc"
		}
	}
	db = db.Order(orderClause)

	// 执行查询
	err = db.Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetDefaultAddress 获取用户默认地址 如果用户没有设置则拉取表内最后一条
func (orderService *OrderService) GetDefaultAddress(UserID uint) (address client.Address, err error) {
	// 尝试获取默认地址
	err = global.GVA_DB.Where("user_id = ? and active = ?", UserID, true).First(&address).Error
	if err != nil {
		// 如果所有active都为0，则使用Order()按表内所有时间排序后返回最后一条
		err = global.GVA_DB.Where("user_id = ?", UserID).Order("created_at desc").First(&address).Error
	}
	return address, nil
}

// ApplyRefund 用户申请退款
func (orderService *OrderService) ApplyRefund(userID uint, req shopReq.RefundApplyReq) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var order shop.Order
		if err := tx.Where("id = ? and user_id = ?", req.OrderID, userID).First(&order).Error; err != nil {
			return errors.New("orderNotFound")
		}
		switch order.Status {
		case "0":
			return errors.New("orderRefundNotPaid")
		case shopOrderStatusPendingConfirm:
			return errors.New("orderRefundPendingConfirmUnsupported")
		case "4":
			return errors.New("orderRefundCanceled")
		case "5":
			return errors.New("orderAlreadyRefunded")
		case "6":
			return errors.New("orderRefundProcessing")
		}
		allowStatus := map[string]bool{"1": true, "2": true, "3": true, "7": true}
		if !allowStatus[order.Status] {
			return errors.New("orderStateInvalidForApplyRefund")
		}

		imagesJSON := []byte("[]")
		if len(req.Images) > 0 {
			if bytes, marshalErr := json.Marshal(req.Images); marshalErr == nil {
				imagesJSON = bytes
			}
		}
		now := time.Now()
		updates := map[string]interface{}{
			"refund_reason":     req.Reason,
			"refund_images":     datatypes.JSON(imagesJSON),
			"refund_applied_at": &now,
			"status":            "6",
		}
		if err := tx.Model(&shop.Order{}).Where("id = ?", order.ID).Updates(updates).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// RefundOrder 后台处理退款
func (orderService *OrderService) RefundOrder(orderID uint, remark string) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var order shop.Order
		if err := tx.First(&order, "id = ?", orderID).Error; err != nil {
			return errors.New("orderNotFound")
		}
		if order.Status != "6" {
			return errors.New("orderRefundStateNotApplying")
		}
		if err := orderService.UpdateOrderStatus(tx, strconv.Itoa(int(orderID)), "5"); err != nil {
			return err
		}
		now := time.Now()
		updates := map[string]interface{}{
			"refund_remark":     remark,
			"refund_handled_at": &now,
		}
		if err := tx.Model(&shop.Order{}).Where("id = ?", orderID).Updates(updates).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// BatchUpdateOrderStatus 批量更新订单状态
func (orderService *OrderService) BatchUpdateOrderStatus(IDs []string, status string) error {
	for _, id := range IDs {
		orderID := strings.TrimSpace(id)
		if orderID == "" {
			continue
		}
		if err := orderService.UpdateOrderStatus(nil, orderID, status); err != nil {
			return err
		}
	}
	return nil
}
