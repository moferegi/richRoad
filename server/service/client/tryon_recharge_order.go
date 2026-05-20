package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
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

	tryonRechargeCreateRateLimitConfigKey = "security_tryon_recharge_create_rate_limit_per_minute"
	tryonRechargeCreateRateLimitEnvKey    = "CS_TRYON_RECHARGE_CREATE_RATE_LIMIT_PER_MINUTE"

	tryonRechargePendingLimitConfigKey = "security_tryon_recharge_pending_limit_per_user"
	tryonRechargePendingLimitEnvKey    = "CS_TRYON_RECHARGE_PENDING_LIMIT_PER_USER"

	tryonRechargePendingConfirmCloseMinutesConfigKey = "tryon_recharge_pending_confirm_close_minutes"
	tryonRechargePendingConfirmCloseMinutesEnvKey    = "CS_TRYON_RECHARGE_PENDING_CONFIRM_CLOSE_MINUTES"

	defaultTryonRechargeCreateRateLimitPerMinute   int64 = 20
	defaultTryonRechargePendingLimitPerUser        int64 = 8
	defaultTryonRechargePendingConfirmCloseMinutes       = 180
)

var paymentMethodPriorityKeys = []string{"zh", "en", "mn", "zh-TW", "th", "hi", "id", "vi", "ar", "ja", "ko", "ms"}

// TryonRechargeOrderService 试衣币充值订单服务
type TryonRechargeOrderService struct{}

func normalizePayMethod(payMethod string) string {
	method := strings.TrimSpace(strings.ToLower(payMethod))
	if method == "" {
		return "contact"
	}
	return method
}

func normalizeSettlementCurrency(value string) string {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return ""
	}
	return strings.ReplaceAll(normalized, "_", "-")
}

func buildSettlementLocaleCandidates(settlementCurrency string) []string {
	normalized := normalizeSettlementCurrency(settlementCurrency)
	candidates := make([]string, 0, len(paymentMethodPriorityKeys)+4)
	if normalized != "" {
		candidates = append(candidates, normalized)
		if strings.Contains(normalized, "-") {
			candidates = append(candidates, strings.SplitN(normalized, "-", 2)[0])
		}
	}
	candidates = append(candidates, "default")
	candidates = append(candidates, paymentMethodPriorityKeys...)

	seen := make(map[string]struct{}, len(candidates))
	result := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" {
			continue
		}
		if _, ok := seen[candidate]; ok {
			continue
		}
		seen[candidate] = struct{}{}
		result = append(result, candidate)
	}
	return result
}

func resolveCurrencySymbolByConfigValue(rawValue, settlementCurrency string) string {
	trimmed := strings.TrimSpace(rawValue)
	if trimmed == "" {
		return ""
	}

	if strings.HasPrefix(trimmed, "{") {
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(trimmed), &parsed); err == nil && len(parsed) > 0 {
			for _, key := range buildSettlementLocaleCandidates(settlementCurrency) {
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

func (s *TryonRechargeOrderService) resolveSettlementCurrencyAndSymbol(tx *gorm.DB, req clientReq.CreateTryonRechargeOrderReq) (string, string) {
	settlementCurrency := normalizeSettlementCurrency(req.SettlementCurrency)
	settlementSymbol := strings.TrimSpace(req.SettlementCurrencySymbol)

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

	return settlementCurrency, settlementSymbol
}

func parsePriceToCents(raw string) (int, error) {
	text := sanitizeNumericText(raw, true)
	if text == "" {
		return 0, errors.New("tryonRechargeAmountFormatError")
	}
	value, err := strconv.ParseFloat(text, 64)
	if err != nil || value < 0 {
		return 0, errors.New("tryonRechargeAmountFormatError")
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

func parseFenFromAny(value interface{}) (int, bool) {
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
	}
	return 0, false
}

func parseLocalizedIntFromAny(value interface{}, settlementCurrency string) (int, bool) {
	if value == nil {
		return 0, false
	}
	if typed, ok := value.(map[string]interface{}); ok {
		for _, key := range buildSettlementLocaleCandidates(settlementCurrency) {
			if field, exists := typed[key]; exists {
				if result, parsed := parseIntFromAny(field); parsed {
					return result, true
				}
			}
		}
	}
	return parseIntFromAny(value)
}

func parseLocalizedCentsFromAny(value interface{}, settlementCurrency string) (int, bool) {
	if value == nil {
		return 0, false
	}
	if typed, ok := value.(map[string]interface{}); ok {
		for _, key := range buildSettlementLocaleCandidates(settlementCurrency) {
			if field, exists := typed[key]; exists {
				if cents, parsed := parseCentsFromAny(field); parsed {
					return cents, true
				}
			}
		}
	}
	return parseCentsFromAny(value)
}

func parseLocalizedFenFromAny(value interface{}, settlementCurrency string) (int, bool) {
	if value == nil {
		return 0, false
	}
	if typed, ok := value.(map[string]interface{}); ok {
		for _, key := range buildSettlementLocaleCandidates(settlementCurrency) {
			if field, exists := typed[key]; exists {
				if fen, parsed := parseFenFromAny(field); parsed {
					return fen, true
				}
			}
		}
	}
	return parseFenFromAny(value)
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

func (s *TryonRechargeOrderService) getPendingConfirmCloseMinutes() int {
	minutes := utils.GetInt64Setting(tryonRechargePendingConfirmCloseMinutesConfigKey, tryonRechargePendingConfirmCloseMinutesEnvKey, int64(defaultTryonRechargePendingConfirmCloseMinutes))
	if minutes <= 0 {
		return defaultTryonRechargePendingConfirmCloseMinutes
	}
	return int(minutes)
}

func (s *TryonRechargeOrderService) enforceTryonRechargeCreateRateLimit(userID uint) error {
	limit := utils.GetInt64Setting(tryonRechargeCreateRateLimitConfigKey, tryonRechargeCreateRateLimitEnvKey, defaultTryonRechargeCreateRateLimitPerMinute)
	if limit <= 0 {
		return nil
	}

	if global.GVA_REDIS != nil {
		ctx := context.Background()
		key := fmt.Sprintf("tryon:recharge:create:rate:user:%d", userID)
		count, err := global.GVA_REDIS.Incr(ctx, key).Result()
		if err == nil {
			if count == 1 {
				_ = global.GVA_REDIS.Expire(ctx, key, time.Minute).Err()
			}
			if count > limit {
				return errors.New("requestTooFrequent")
			}
			return nil
		}
		global.GVA_LOG.Warn("试衣币充值建单频率校验降级为DB", zap.Error(err), zap.Uint("userID", userID))
	}

	var recentCount int64
	if err := global.GVA_DB.Model(&client.TryonRechargeOrder{}).
		Where("user_id = ? AND created_at >= ?", userID, time.Now().Add(-time.Minute)).
		Count(&recentCount).Error; err != nil {
		global.GVA_LOG.Warn("试衣币充值建单频率DB校验失败，已降级放行", zap.Error(err), zap.Uint("userID", userID))
		return nil
	}

	if recentCount >= limit {
		return errors.New("requestTooFrequent")
	}

	return nil
}

func (s *TryonRechargeOrderService) enforceTryonRechargePendingLimit(tx *gorm.DB, userID uint) error {
	limit := utils.GetInt64Setting(tryonRechargePendingLimitConfigKey, tryonRechargePendingLimitEnvKey, defaultTryonRechargePendingLimitPerUser)
	if limit <= 0 {
		return nil
	}

	var pendingCount int64
	err := tx.Model(&client.TryonRechargeOrder{}).
		Where("user_id = ? AND status IN ? AND (close_time IS NULL OR close_time > ?)", userID, []string{tryonRechargeOrderStatusPending, tryonRechargeOrderStatusReview}, time.Now()).
		Count(&pendingCount).Error
	if err != nil {
		return err
	}

	if pendingCount >= limit {
		return errors.New("tryonRechargePendingLimitExceeded")
	}

	return nil
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

	return "", errors.New("tryonRechargeOrderNoGenFail")
}

func (s *TryonRechargeOrderService) findMatchedRechargePlan(tx *gorm.DB, points int, amountCents int, settlementCurrency string) (map[string]interface{}, int, error) {
	var config client.SysConfig
	err := tx.Where("config_key = ?", "tryon_recharge_plans").First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errors.New("tryonRechargePlanNotConfigured")
		}
		return nil, 0, err
	}

	var plans []map[string]interface{}
	if unmarshalErr := json.Unmarshal([]byte(config.ConfigValue), &plans); unmarshalErr != nil {
		return nil, 0, errors.New("tryonRechargePlanConfigInvalid")
	}

	for _, plan := range plans {
		planPoints, pointsParsed := parseLocalizedIntFromAny(plan["points"], settlementCurrency)
		if !pointsParsed || planPoints != points {
			continue
		}

		planAmount := 0
		amountParsed := false

		if localizedFen, parsed := parseLocalizedFenFromAny(plan["priceI18n"], settlementCurrency); parsed {
			planAmount = localizedFen
			amountParsed = true
		}

		if !amountParsed {
			switch plan["price"].(type) {
			case int, int32, int64, float32, float64, json.Number:
				if baseFen, parsed := parseFenFromAny(plan["price"]); parsed {
					planAmount = baseFen
					amountParsed = true
				}
			}
		}

		if !amountParsed {
			if localizedCents, parsed := parseLocalizedCentsFromAny(plan["price"], settlementCurrency); parsed {
				planAmount = localizedCents
				amountParsed = true
			}
		}

		if amountParsed && planAmount == amountCents {
			return plan, planAmount, nil
		}
	}
	return nil, 0, errors.New("tryonRechargePlanChanged")
}

// CreateTryonRechargeOrder 创建充值订单
func (s *TryonRechargeOrderService) CreateTryonRechargeOrder(ctx context.Context, userID uint, req clientReq.CreateTryonRechargeOrderReq) (order client.TryonRechargeOrder, err error) {
	if userID == 0 {
		return order, errors.New("loginRequired")
	}
	if req.Points <= 0 {
		return order, errors.New("tryonRechargePointsMustPositive")
	}
	amountCents, parseErr := parsePriceToCents(req.Price)
	if parseErr != nil {
		return order, parseErr
	}
	if amountCents <= 0 {
		return order, errors.New("tryonRechargeAmountMustPositive")
	}
	if err = s.enforceTryonRechargeCreateRateLimit(userID); err != nil {
		return order, err
	}

	payMethod := normalizePayMethod(req.PayMethod)
	err = global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = s.enforceTryonRechargePendingLimit(tx, userID); err != nil {
			return err
		}

		settlementCurrency, settlementCurrencySymbol := s.resolveSettlementCurrencyAndSymbol(tx, req)

		matchedPlan, matchedPlanAmount, matchErr := s.findMatchedRechargePlan(tx, req.Points, amountCents, settlementCurrency)
		if matchErr != nil {
			return matchErr
		}

		snapshotPlan := make(map[string]interface{}, len(matchedPlan)+4)
		for key, value := range matchedPlan {
			snapshotPlan[key] = value
		}
		snapshotPlan["selectedPriceFen"] = matchedPlanAmount
		snapshotPlan["settlementCurrency"] = settlementCurrency
		snapshotPlan["settlementCurrencySymbol"] = settlementCurrencySymbol
		snapshotPlan["priceSnapshotVersion"] = 1

		snapshotBytes, marshalErr := json.Marshal(snapshotPlan)
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
			UserID:                   userID,
			Status:                   tryonRechargeOrderStatusPending,
			Points:                   req.Points,
			Amount:                   amountCents,
			Currency:                 "CNY",
			SettlementCurrency:       settlementCurrency,
			SettlementCurrencySymbol: settlementCurrencySymbol,
			PayMethod:                payMethod,
			PlanSnapshot:             datatypes.JSON(snapshotBytes),
			CloseTime:                now.Add(time.Duration(closeMinutes) * time.Minute),
			OutTradeNo:               outTradeNo,
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
		return errors.New("loginRequired")
	}
	if orderID == 0 {
		return errors.New("invalidOrder")
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var order client.TryonRechargeOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("orderNotFound")
			}
			return err
		}

		if order.Status == tryonRechargeOrderStatusReview {
			return nil
		}
		if order.Status != tryonRechargeOrderStatusPending {
			return errors.New("tryonRechargeOrderStateInvalidForSubmitPayment")
		}
		if !order.CloseTime.IsZero() && time.Now().After(order.CloseTime) {
			return errors.New("tryonRechargeOrderExpired")
		}
		if normalizePayMethod(order.PayMethod) != "qrcode" {
			return errors.New("tryonRechargeOrderPayMethodNotQrcode")
		}

		updates := map[string]interface{}{"status": tryonRechargeOrderStatusReview}
		pendingConfirmCloseMinutes := s.getPendingConfirmCloseMinutes()
		if pendingConfirmCloseMinutes > 0 {
			updates["close_time"] = time.Now().Add(time.Duration(pendingConfirmCloseMinutes) * time.Minute)
		}
		return tx.Model(&client.TryonRechargeOrder{}).Where("id = ?", order.ID).Updates(updates).Error
	})
}

// UpdateTryonRechargeOrderPayMethod 更新待支付充值订单的支付方式
func (s *TryonRechargeOrderService) UpdateTryonRechargeOrderPayMethod(userID uint, req clientReq.UpdateTryonRechargeOrderPayMethodReq) error {
	if userID == 0 {
		return errors.New("loginRequired")
	}
	payMethod := normalizePayMethod(req.PayMethod)
	if payMethod == "" {
		return errors.New("payMethodRequired")
	}
	now := time.Now()

	result := global.GVA_DB.Model(&client.TryonRechargeOrder{}).
		Where("id = ? AND user_id = ? AND status = ?", req.ID, userID, tryonRechargeOrderStatusPending).
		Where("(close_time IS NULL OR close_time > ?)", now).
		Updates(map[string]interface{}{"pay_method": payMethod})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("tryonRechargeOrderNotEditable")
	}
	return nil
}

// CancelTryonRechargeOrder 取消充值订单
func (s *TryonRechargeOrderService) CancelTryonRechargeOrder(userID uint, orderID uint) error {
	if userID == 0 {
		return errors.New("loginRequired")
	}
	if orderID == 0 {
		return errors.New("invalidOrder")
	}
	now := time.Now()
	result := global.GVA_DB.Model(&client.TryonRechargeOrder{}).
		Where("id = ? AND user_id = ? AND status = ?", orderID, userID, tryonRechargeOrderStatusPending).
		Updates(map[string]interface{}{"status": tryonRechargeOrderStatusCanceled, "cancelled_at": &now})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("tryonRechargeOrderNotCancelable")
	}
	return nil
}

// ConfirmTryonRechargeOrderPayment 确认充值订单支付（管理端）
func (s *TryonRechargeOrderService) ConfirmTryonRechargeOrderPayment(ctx context.Context, orderID uint, remark string) error {
	if orderID == 0 {
		return errors.New("invalidOrder")
	}

	pointRecordService := PointRecordService{}
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order client.TryonRechargeOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", orderID).First(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("orderNotFound")
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
			return errors.New("tryonRechargeOrderStateInvalidForConfirm")
		}
		if !order.CloseTime.IsZero() && time.Now().After(order.CloseTime) {
			return errors.New("tryonRechargeOrderExpired")
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
		reason := `{"zh":"试衣币充值到账","en":"Try-on coins recharge credited","mn":"Туршилтын зоос дансанд орлоо","ja":"試着コインのチャージが入金されました","ko":"피팅 코인 충전이 반영되었습니다","ms":"Tambah nilai syiling cuba pakaian telah dikreditkan"}`
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
		return order, errors.New("invalidOrder")
	}
	err = global.GVA_DB.Where("id = ?", orderID).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return order, errors.New("orderNotFound")
	}
	return order, err
}

// GetMyTryonRechargeOrderByID 获取我的充值订单
func (s *TryonRechargeOrderService) GetMyTryonRechargeOrderByID(userID uint, orderID uint) (order client.TryonRechargeOrder, err error) {
	if userID == 0 {
		return order, errors.New("loginRequired")
	}
	if orderID == 0 {
		return order, errors.New("invalidOrder")
	}
	err = global.GVA_DB.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return order, errors.New("orderNotFound")
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
		return list, 0, errors.New("loginRequired")
	}
	search.UserID = userID
	return s.GetTryonRechargeOrderList(search)
}
