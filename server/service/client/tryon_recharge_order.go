package client

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	tryonRechargeOrderStatusPending  = "0"
	tryonRechargeOrderStatusReview   = "8"
	tryonRechargeOrderStatusPaid     = "1"
	tryonRechargeOrderStatusCanceled = "4"

	defaultTryonRechargeCloseMinutes = 20
)

var paymentMethodPriorityKeys = []string{"zh", "en", "mn", "zh-TW", "th", "hi", "id"}

// TryonRechargeOrderService 试衣币充值订单服务
type TryonRechargeOrderService struct{}

func normalizePayMethod(payMethod string) string {
	method := strings.TrimSpace(strings.ToLower(payMethod))
	if method == "" {
		return "contact"
	}
	return method
}

func parsePriceToCents(raw string) (int, error) {
	text := sanitizeNumericText(raw, true)
	if text == "" {
		return 0, errors.New("充值金额格式不正确")
	}
	value, err := strconv.ParseFloat(text, 64)
	if err != nil || value < 0 {
		return 0, errors.New("充值金额格式不正确")
	}
	return int(math.Round(value * 100)), nil
}

func sanitizeNumericText(raw string, allowDot bool) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	var builder strings.Builder
	dotUsed := false
	for _, ch := range raw {
		if ch >= '0' && ch <= '9' {
			builder.WriteRune(ch)
			continue
		}
		if allowDot && ch == '.' && !dotUsed {
			dotUsed = true
			if builder.Len() == 0 {
				builder.WriteRune('0')
			}
			builder.WriteRune(ch)
		}
	}
	return builder.String()
}

func parseIntFromAny(value interface{}) (int, bool) {
	switch typed := value.(type) {
	case nil:
		return 0, false
	case int:
		return typed, true
	case int32:
		return int(typed), true
	case int64:
		return int(typed), true
	case float32:
		return int(math.Round(float64(typed))), true
	case float64:
		return int(math.Round(typed)), true
	case json.Number:
		intValue, err := typed.Int64()
		if err == nil {
			return int(intValue), true
		}
		floatValue, ferr := typed.Float64()
		if ferr == nil {
			return int(math.Round(floatValue)), true
		}
		return 0, false
	case string:
		text := sanitizeNumericText(typed, false)
		if text == "" {
			return 0, false
		}
		intValue, err := strconv.Atoi(text)
		if err != nil {
			return 0, false
		}
		return intValue, true
	case map[string]interface{}:
		for _, key := range paymentMethodPriorityKeys {
			if field, ok := typed[key]; ok {
				if result, parsed := parseIntFromAny(field); parsed {
					return result, true
				}
			}
		}
		for _, field := range typed {
			if result, parsed := parseIntFromAny(field); parsed {
				return result, true
			}
		}
	case []interface{}:
		for _, field := range typed {
			if result, parsed := parseIntFromAny(field); parsed {
				return result, true
			}
		}
	}
	return 0, false
}

func parseCentsFromAny(value interface{}) (int, bool) {
	switch typed := value.(type) {
	case nil:
		return 0, false
	case int:
		return typed * 100, true
	case int32:
		return int(typed) * 100, true
	case int64:
		return int(typed) * 100, true
	case float32:
		return int(math.Round(float64(typed) * 100)), true
	case float64:
		return int(math.Round(typed * 100)), true
	case json.Number:
		floatValue, err := typed.Float64()
		if err != nil {
			return 0, false
		}
		return int(math.Round(floatValue * 100)), true
	case string:
		cents, err := parsePriceToCents(typed)
		if err != nil {
			return 0, false
		}
		return cents, true
	case map[string]interface{}:
		for _, key := range paymentMethodPriorityKeys {
			if field, ok := typed[key]; ok {
				if cents, parsed := parseCentsFromAny(field); parsed {
					return cents, true
				}
			}
		}
		for _, field := range typed {
			if cents, parsed := parseCentsFromAny(field); parsed {
				return cents, true
			}
		}
	case []interface{}:
		for _, field := range typed {
			if cents, parsed := parseCentsFromAny(field); parsed {
				return cents, true
			}
		}
	}
	return 0, false
}

func (s *TryonRechargeOrderService) getOrderCloseMinutes(tx *gorm.DB) int {
	var config client.SysConfig
	err := tx.Where("config_key = ?", "order_close_minutes").First(&config).Error
	if err != nil {
		return defaultTryonRechargeCloseMinutes
	}
	minutes, parseErr := strconv.Atoi(strings.TrimSpace(config.ConfigValue))
	if parseErr != nil || minutes <= 0 {
		return defaultTryonRechargeCloseMinutes
	}
	return minutes
}

func (s *TryonRechargeOrderService) generateUniqueTryonRechargeOrderNo(tx *gorm.DB) (string, error) {
	for i := 0; i < 10; i++ {
		orderNo, err := utils.GenerateBusinessOrderNo("sy")
		if err != nil {
			return "", err
		}

		var count int64
		if err := tx.Model(&client.TryonRechargeOrder{}).Where("out_trade_no = ?", orderNo).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return orderNo, nil
		}
	}

	return "", errors.New("订单号生成失败，请稍后重试")
}

func (s *TryonRechargeOrderService) findMatchedRechargePlan(tx *gorm.DB, points int, amountCents int) (map[string]interface{}, error) {
	var config client.SysConfig
	err := tx.Where("config_key = ?", "tryon_recharge_plans").First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("充值套餐未配置")
		}
		return nil, err
	}

	var plans []map[string]interface{}
	if unmarshalErr := json.Unmarshal([]byte(config.ConfigValue), &plans); unmarshalErr != nil {
		return nil, errors.New("充值套餐配置格式错误")
	}

	for _, plan := range plans {
		planPoints, pointsParsed := parseIntFromAny(plan["points"])
		planAmount, amountParsed := parseCentsFromAny(plan["price"])
		if pointsParsed && amountParsed && planPoints == points && planAmount == amountCents {
			return plan, nil
		}
	}
	return nil, errors.New("充值套餐已变更，请刷新页面后重试")
}

// CreateTryonRechargeOrder 创建充值订单
func (s *TryonRechargeOrderService) CreateTryonRechargeOrder(ctx context.Context, userID uint, req clientReq.CreateTryonRechargeOrderReq) (order client.TryonRechargeOrder, err error) {
	if userID == 0 {
		return order, errors.New("请先登录")
	}
	if req.Points <= 0 {
		return order, errors.New("充值点数必须大于0")
	}
	amountCents, parseErr := parsePriceToCents(req.Price)
	if parseErr != nil {
		return order, parseErr
	}
	if amountCents <= 0 {
		return order, errors.New("充值金额必须大于0")
	}

	payMethod := normalizePayMethod(req.PayMethod)
	err = global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		matchedPlan, matchErr := s.findMatchedRechargePlan(tx, req.Points, amountCents)
		if matchErr != nil {
			return matchErr
		}

		snapshotBytes, marshalErr := json.Marshal(matchedPlan)
		if marshalErr != nil {
			return marshalErr
		}

		outTradeNo, noErr := s.generateUniqueTryonRechargeOrderNo(tx)
		if noErr != nil {
			return noErr
		}

		closeMinutes := s.getOrderCloseMinutes(tx)
		now := time.Now()

		order = client.TryonRechargeOrder{
			UserID:       userID,
			Status:       tryonRechargeOrderStatusPending,
			Points:       req.Points,
			Amount:       amountCents,
			Currency:     "CNY",
			PayMethod:    payMethod,
			PlanSnapshot: datatypes.JSON(snapshotBytes),
			CloseTime:    now.Add(time.Duration(closeMinutes) * time.Minute),
			OutTradeNo:   outTradeNo,
		}
		if createErr := tx.Create(&order).Error; createErr != nil {
			return createErr
		}
		return nil
	})
	return order, err
}

// SubmitTryonRechargeOrderPayment 用户提交充值订单付款确认（仅扫码支付且待支付订单）
func (s *TryonRechargeOrderService) SubmitTryonRechargeOrderPayment(userID uint, orderID uint) error {
	if userID == 0 {
		return errors.New("请先登录")
	}
	if orderID == 0 {
		return errors.New("订单参数错误")
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var order client.TryonRechargeOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("订单不存在")
			}
			return err
		}

		if order.Status == tryonRechargeOrderStatusReview {
			return nil
		}
		if order.Status != tryonRechargeOrderStatusPending {
			return errors.New("当前订单状态不支持提交付款确认")
		}
		if normalizePayMethod(order.PayMethod) != "qrcode" {
			return errors.New("仅扫码支付订单可提交付款确认")
		}

		return tx.Model(&client.TryonRechargeOrder{}).Where("id = ?", order.ID).Update("status", tryonRechargeOrderStatusReview).Error
	})
}

// UpdateTryonRechargeOrderPayMethod 更新待支付充值订单的支付方式
func (s *TryonRechargeOrderService) UpdateTryonRechargeOrderPayMethod(userID uint, req clientReq.UpdateTryonRechargeOrderPayMethodReq) error {
	if userID == 0 {
		return errors.New("请先登录")
	}
	payMethod := normalizePayMethod(req.PayMethod)
	if payMethod == "" {
		return errors.New("支付方式不能为空")
	}

	result := global.GVA_DB.Model(&client.TryonRechargeOrder{}).
		Where("id = ? AND user_id = ? AND status = ?", req.ID, userID, tryonRechargeOrderStatusPending).
		Updates(map[string]interface{}{"pay_method": payMethod})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("订单不存在或当前状态不支持修改")
	}
	return nil
}

// CancelTryonRechargeOrder 取消充值订单
func (s *TryonRechargeOrderService) CancelTryonRechargeOrder(userID uint, orderID uint) error {
	if userID == 0 {
		return errors.New("请先登录")
	}
	if orderID == 0 {
		return errors.New("订单参数错误")
	}
	now := time.Now()
	result := global.GVA_DB.Model(&client.TryonRechargeOrder{}).
		Where("id = ? AND user_id = ? AND status = ?", orderID, userID, tryonRechargeOrderStatusPending).
		Updates(map[string]interface{}{"status": tryonRechargeOrderStatusCanceled, "cancelled_at": &now})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("订单不存在或当前状态不支持取消")
	}
	return nil
}

// ConfirmTryonRechargeOrderPayment 确认充值订单支付（管理端）
func (s *TryonRechargeOrderService) ConfirmTryonRechargeOrderPayment(ctx context.Context, orderID uint, remark string) error {
	if orderID == 0 {
		return errors.New("订单参数错误")
	}

	pointRecordService := PointRecordService{}
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order client.TryonRechargeOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", orderID).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("订单不存在")
			}
			return err
		}

		switch order.Status {
		case tryonRechargeOrderStatusPaid:
			return nil
		case tryonRechargeOrderStatusPending:
			// continue
		case tryonRechargeOrderStatusReview:
			// continue
		default:
			return errors.New("当前订单状态不支持确认支付")
		}

		now := time.Now()
		if err := tx.Model(&client.TryonRechargeOrder{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
			"status":  tryonRechargeOrderStatusPaid,
			"paid_at": &now,
			"remark":  strings.TrimSpace(remark),
		}).Error; err != nil {
			return err
		}

		txCtx := context.WithValue(ctx, "tx", tx)
		reason := `{"zh":"试衣币充值到账","en":"Try-on coins recharge credited","mn":"Туршилтын зоос дансанд орлоо"}`
		pointRecord := buildPointRecord(order.UserID, client.AssetTypeTryonPoint, "increase", order.Points, "tryon_recharge", reason, strings.TrimSpace(order.PayMethod))
		relatedOrderID := int(order.ID)
		pointRecord.RelatedOrderId = &relatedOrderID
		if err := pointRecordService.CreatePointRecord(txCtx, pointRecord); err != nil {
			return err
		}
		return nil
	})
}

// GetTryonRechargeOrderByID 获取充值订单（管理端）
func (s *TryonRechargeOrderService) GetTryonRechargeOrderByID(orderID uint) (order client.TryonRechargeOrder, err error) {
	if orderID == 0 {
		return order, errors.New("订单参数错误")
	}
	err = global.GVA_DB.Where("id = ?", orderID).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return order, errors.New("订单不存在")
	}
	return order, err
}

// GetMyTryonRechargeOrderByID 获取我的充值订单
func (s *TryonRechargeOrderService) GetMyTryonRechargeOrderByID(userID uint, orderID uint) (order client.TryonRechargeOrder, err error) {
	if userID == 0 {
		return order, errors.New("请先登录")
	}
	if orderID == 0 {
		return order, errors.New("订单参数错误")
	}
	err = global.GVA_DB.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return order, errors.New("订单不存在")
	}
	return order, err
}

func buildTryonRechargeOrderSearch(db *gorm.DB, search clientReq.TryonRechargeOrderSearch) *gorm.DB {
	if search.UserID > 0 {
		db = db.Where("user_id = ?", search.UserID)
	}
	if status := strings.TrimSpace(search.Status); status != "" {
		db = db.Where("status = ?", status)
	}
	if search.StartCreatedAt != nil {
		db = db.Where("created_at >= ?", *search.StartCreatedAt)
	}
	if search.EndCreatedAt != nil {
		db = db.Where("created_at <= ?", *search.EndCreatedAt)
	}
	return db
}

// GetTryonRechargeOrderList 获取充值订单列表（管理端）
func (s *TryonRechargeOrderService) GetTryonRechargeOrderList(search clientReq.TryonRechargeOrderSearch) (list []client.TryonRechargeOrder, total int64, err error) {
	page := search.Page
	pageSize := search.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := pageSize * (page - 1)

	db := buildTryonRechargeOrderSearch(global.GVA_DB.Model(&client.TryonRechargeOrder{}), search)
	if err = db.Count(&total).Error; err != nil {
		return list, total, err
	}
	err = db.Order("id desc").Limit(pageSize).Offset(offset).Find(&list).Error
	return list, total, err
}

// GetMyTryonRechargeOrderList 获取我的充值订单列表
func (s *TryonRechargeOrderService) GetMyTryonRechargeOrderList(userID uint, search clientReq.TryonRechargeOrderSearch) (list []client.TryonRechargeOrder, total int64, err error) {
	if userID == 0 {
		return list, 0, errors.New("请先登录")
	}
	search.UserID = userID
	return s.GetTryonRechargeOrderList(search)
}
