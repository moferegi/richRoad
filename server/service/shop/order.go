package shop

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	clientService "github.com/flipped-aurora/gin-vue-admin/server/service/client"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OrderService struct {
}

// CreateOrder 创建订单记录
// Author [piexlmax](https://github.com/piexlmax)
func (orderService *OrderService) CreateOrder(order *shop.Order) (err error) {
	err = global.GVA_DB.Create(order).Error
	return err
}

func (orderService *OrderService) ChangeOrderCoupon(userID uint, orderID string, couponNum string) (err error) {
	var order shop.Order
	err = global.GVA_DB.Where("id = ? and user_id = ?", orderID, userID).Preload("Detail").First(&order).Error
	if err != nil {
		return err
	}

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
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
		err = tx.Where("coupon_num = ? and order_id IS NULL", couponNum).First(&couponOrderUser).Error
		if err != nil {
			return err
		}
		var coupon shop.Coupon
		err = tx.Where("id = ?", couponOrderUser.CouponID).First(&coupon).Error
		if err != nil {
			return err
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
					return errors.New("当前订单不可使用此券")
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
		err = tx.Model(&shop.CouponOrderUser{}).Where("coupon_num = ?", couponNum).Update("order_id", order.ID).Error
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// ChangeOrderPoints 变更订单积分抵扣
func (orderService *OrderService) ChangeOrderPoints(userID uint, orderID string, usePoints bool) (err error) {
	var order shop.Order
	err = global.GVA_DB.Where("id = ? and user_id = ?", orderID, userID).Preload("Detail").First(&order).Error
	if err != nil {
		return err
	}

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 1. 如果之前使用了积分，需要先恢复积分
		if order.UsePoints && order.PointsUsed > 0 {
			pointRecordService := clientService.PointRecordService{}
			userID := int(order.UserID)
			pointChange := int(order.PointsUsed) // 恢复积分，正数
			changeType := "increase"
			operationType := "point_exchange"
			reason := "取消积分抵扣，恢复积分"
			relatedOrderId := int(order.ID)

			pointRecord := client.PointRecord{
				UserId:         &userID,
				PointChange:    &pointChange,
				ChangeType:     &changeType,
				OperationType:  &operationType,
				Reason:         &reason,
				RelatedOrderId: &relatedOrderId,
			}
			// 将事务挂在context上传递
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
			var user client.ClientUser
			err = tx.Where("id = ?", order.UserID).First(&user).Error
			if err != nil {
				return err
			}
			// 计算可抵扣的积分数（不能超过订单金额）
			maxPointsCanUse := order.TotalPrice
			pointsToUse := user.Point
			if pointsToUse > maxPointsCanUse {
				pointsToUse = maxPointsCanUse
			}
			if pointsToUse > 0 {
				// 扣除积分
				pointRecordService := clientService.PointRecordService{}
				userID := int(order.UserID)
				pointChange := -int(pointsToUse) // 扣除积分，负数
				changeType := "decrease"
				operationType := "point_exchange"
				reason := "订单积分抵扣"
				relatedOrderId := int(order.ID)

				pointRecord := client.PointRecord{
					UserId:         &userID,
					PointChange:    &pointChange,
					ChangeType:     &changeType,
					OperationType:  &operationType,
					Reason:         &reason,
					RelatedOrderId: &relatedOrderId,
				}
				// 将事务挂在context上传递
				ctxWithTx := context.WithValue(context.Background(), "tx", tx)
				err = pointRecordService.CreatePointRecord(ctxWithTx, &pointRecord)
				if err != nil {
					return err
				}
				// 更新订单价格和积分使用信息
				order.TotalPrice -= pointsToUse
				order.PointsUsed = uint(pointsToUse)
			}
		}

		// 4. 确保订单金额不为负数
		if order.TotalPrice < 0 {
			order.TotalPrice = 0
		}

		// 5. 更新订单信息
		err = tx.Model(&order).
			Update("total_price", order.TotalPrice).
			Update("use_points", order.UsePoints).
			Update("points_used", order.PointsUsed).
			Error
		if err != nil {
			return err
		}
		return nil
	})
	return
}

func (orderService *OrderService) PlaceOrder(order *shop.Order) (OrderID uint, err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
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
				return errors.New("库存不足")
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
				return errors.New("库存不足，请重试")
			}
			err = tx.Model(&shop.Sku{}).Where("id = ?", order.Detail[i].SKUID).
				Update("sale_num", gorm.Expr("sale_num + ?", order.Detail[i].Quantity)).Error
			if err != nil {
				return err
			}
			err = tx.Model(&shop.Good{}).Where("id = ?", sku.GoodID).
				Update("sale_num", gorm.Expr("sale_num + ?", order.Detail[i].Quantity)).Error
			if err != nil {
				return err
			}
			order.Detail[i].Price = sku.Price
			order.OriginPrice += order.Detail[i].Quantity * order.Detail[i].Price
			// 减扣库存
		}
		order.TotalPrice = int(order.OriginPrice)
		if order.CouponNum != "" {
			var couponOrderUser shop.CouponOrderUser
			err = tx.Where("coupon_num = ? and order_id IS NULL", order.CouponNum).First(&couponOrderUser).Error
			if err != nil {
				return err
			}
			var coupon shop.Coupon
			err = tx.Where("id = ?", couponOrderUser.CouponID).First(&coupon).Error
			if err != nil {
				return err
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
					return errors.New("当前订单不可使用此券")
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

		if addr, err := orderService.GetDefaultAddress(order.UserID); err == nil {
			order.Name = addr.Name
			order.Phone = addr.Phone
			order.Province = addr.ProvinceStr
			order.City = addr.CityStr
			order.Area = addr.AreaStr
			order.Street = addr.Street
		}
		// 从系统配置读取订单自动关闭时间（分钟），默认15分钟
		closeMinutes := 15
		var sysConf client.SysConfig
		if e := tx.Where("config_key = ?", "order_auto_close_minutes").First(&sysConf).Error; e == nil {
			if v, pe := strconv.Atoi(sysConf.ConfigValue); pe == nil && v > 0 {
				closeMinutes = v
			}
		}
		order.CloseTime = time.Now().Add(time.Duration(closeMinutes) * time.Minute)

		err = tx.Create(order).Error
		if err != nil {
			return err
		}
		OrderID = order.ID
		err = tx.Model(&shop.CouponOrderUser{}).Where("coupon_num = ?", order.CouponNum).Update("order_id", order.ID).Error
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// 购物车下单
func (orderService *OrderService) PlaceOrderByCart(userID uint, req shopReq.PlaceOrderByCartRequest) (OrderID uint, err error) {
	// 1. 获取购物车中的商品
	var carts []shop.Cart
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		query := tx.Where("user_id = ?", userID).Preload("Good").Preload("SKU")
		if len(req.CartIDs) > 0 {
			query = query.Where("id IN ?", req.CartIDs)
		}
		query.Find(&carts)
		if len(carts) == 0 {
			return errors.New("购物车为空")
		}
		// 2. 创建订单
		order := shop.Order{
			UserID:    userID,
			Status:    "0",
			PayMethod: req.PayMethod,
		}
		order.OriginPrice = 0
		// 从购物车设置订单详情
		order.Detail = make([]shop.OrderDetail, 0)
		presaleService := PresaleService{}
		for i := range carts {
			// 判断当前库存是否充足
			if carts[i].SKU.Inventory < carts[i].Quantity {
				return errors.New("库存不足")
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
				GoodID:   carts[i].GoodID,
				SKUID:    carts[i].SKUID,
				Quantity: carts[i].Quantity,
				Price:    carts[i].SKU.Price,
			})
			order.OriginPrice += carts[i].Quantity * carts[i].SKU.Price
			// 乐观锁扣减库存：WHERE inventory >= quantity 防止并发超卖
			result := tx.Model(&shop.Sku{}).Where("id = ? AND inventory >= ?", carts[i].SKUID, carts[i].Quantity).
				Update("inventory", gorm.Expr("inventory - ?", carts[i].Quantity))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New("库存不足，请重试")
			}
			err = tx.Model(&shop.Sku{}).Where("id = ?", carts[i].SKUID).
				Update("sale_num", gorm.Expr("sale_num + ?", carts[i].Quantity)).Error
			if err != nil {
				return err
			}
			err = tx.Model(&shop.Good{}).Where("id = ?", carts[i].GoodID).Update("sale_num", gorm.Expr("sale_num + ?", carts[i].Quantity)).Error
			if err != nil {
				return err
			}
		}

		if addr, err := orderService.GetDefaultAddress(order.UserID); err == nil {
			order.Name = addr.Name
			order.Phone = addr.Phone
			order.Province = addr.ProvinceStr
			order.City = addr.CityStr
			order.Area = addr.AreaStr
			order.Street = addr.Street
		}
		// 从系统配置读取订单自动关闭时间（分钟），默认15分钟
		closeMinutes := 15
		var sysConf client.SysConfig
		if e := tx.Where("config_key = ?", "order_auto_close_minutes").First(&sysConf).Error; e == nil {
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
		err = tx.First(&order, "id = ?", orderID).Update("status", status).Error
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

			// 如果订单已积分（Pointed=true），还需要处理积分返还和扣除
			if order.Pointed {
				// 1. 如果订单使用了积分，先返还已使用的积分
				if order.UsePoints && order.PointsUsed > 0 {
					pointRecordService := &clientService.PointRecordService{}
					userIdInt := int(user.ID)
					changeType := "increase"
					pointChange := int(order.PointsUsed) // 返还积分，正数
					operationType := "point_refund"
					reason := "订单取消，返还已使用积分"
					orderIdStr := strconv.Itoa(int(order.ID))
					orderIdInt, _ := strconv.Atoi(orderIdStr)
					remark := "订单ID: " + orderIdStr

					pointRecord := &client.PointRecord{
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
				changeType := "decrease"
				pointChange := int(totalPrice / 100)
				operationType := "refund_return"
				reason := "订单取消，扣除已获得积分"
				orderIdStr := strconv.Itoa(int(order.ID))
				orderIdInt, _ := strconv.Atoi(orderIdStr)
				remark := "订单ID: " + orderIdStr

				pointRecord := &client.PointRecord{
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
			// 幂等检查：查看该订单是否已经发放过奖励（防止取消后重新支付重复发放）
			var rewardCount int64
			tx.Model(&client.PointRecord{}).Where("related_order_id = ? AND operation_type = ?", order.ID, "order_complete").Count(&rewardCount)
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
	if order.Status != "0" {
		return fmt.Errorf("只能对待付款订单确认收款，当前状态: %s", order.Status)
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
			return errors.New("订单不存在")
		}
		switch order.Status {
		case "0":
			return errors.New("订单未支付，无法申请退款")
		case "4":
			return errors.New("订单已取消，无法申请退款")
		case "5":
			return errors.New("订单已退款")
		case "6":
			return errors.New("退款申请处理中")
		}
		allowStatus := map[string]bool{"1": true, "2": true, "3": true, "7": true}
		if !allowStatus[order.Status] {
			return errors.New("当前订单状态不允许退款")
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
			return errors.New("订单不存在")
		}
		if order.Status != "6" {
			return errors.New("订单未处于退款申请中")
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
