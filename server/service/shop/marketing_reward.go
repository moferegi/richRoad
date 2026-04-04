package shop

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	clientService "github.com/flipped-aurora/gin-vue-admin/server/service/client"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MarketingRewardService struct{}

// CreateMarketingReward 创建营销奖励配置
func (s *MarketingRewardService) CreateMarketingReward(reward *shop.MarketingReward) (err error) {
	err = global.GVA_DB.Create(reward).Error
	return err
}

// DeleteMarketingReward 删除营销奖励配置
func (s *MarketingRewardService) DeleteMarketingReward(ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.MarketingReward{}, "id = ?", ID).Error
	return err
}

// UpdateMarketingReward 更新营销奖励配置
func (s *MarketingRewardService) UpdateMarketingReward(reward shop.MarketingReward) (err error) {
	err = global.GVA_DB.Model(&shop.MarketingReward{}).Where("id = ?", reward.ID).Updates(&reward).Error
	return err
}

// GetMarketingReward 根据ID获取营销奖励配置
func (s *MarketingRewardService) GetMarketingReward(ID string) (reward shop.MarketingReward, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&reward).Error
	return
}

// GetMarketingRewardByType 按触发类型获取
func (s *MarketingRewardService) GetMarketingRewardByType(triggerType string) (reward shop.MarketingReward, err error) {
	err = global.GVA_DB.Where("trigger_type = ? AND is_enabled = ?", triggerType, true).First(&reward).Error
	return
}

// GetMarketingRewardList 分页获取营销奖励配置列表
func (s *MarketingRewardService) GetMarketingRewardList(info shopReq.MarketingRewardSearch) (list []shop.MarketingReward, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&shop.MarketingReward{})

	if info.TriggerType != "" {
		db = db.Where("trigger_type = ?", info.TriggerType)
	}
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Find(&list).Error
	return
}

// TriggerReward 统一触发营销奖励（积分+优惠券）
// triggerType: register / sub_register / sign_in / order
// userID: 奖励接收者
// operationType: 积分记录的操作类型
// reason: 积分记录的原因
// relatedOrderID: 关联订单ID（可选，0表示无）
func (s *MarketingRewardService) TriggerReward(userID uint, triggerType string, operationType string, reason string, relatedOrderID uint) error {
	reward, err := s.GetMarketingRewardByType(triggerType)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // 未配置该类型奖励，静默跳过
		}
		return err
	}

	// 下单奖励仅发放一次：检查该用户是否已获得过订单奖励
	if triggerType == "order" && reward.OrderOnce != nil && *reward.OrderOnce {
		var count int64
		global.GVA_DB.Model(&client.PointRecord{}).
			Where("user_id = ? AND operation_type = ?", userID, "order_complete").
			Count(&count)
		if count > 0 {
			return nil // 已发放过，跳过
		}
	}

	// 发放积分
	if reward.Points != nil && *reward.Points > 0 {
		pointRecordService := &clientService.PointRecordService{}
		uid := int(userID)
		changeType := "increase"
		points := *reward.Points
		remark := ""
		if relatedOrderID > 0 {
			remark = "订单ID: " + strconv.Itoa(int(relatedOrderID))
		}

		record := &client.PointRecord{
			UserId:        &uid,
			ChangeType:    &changeType,
			PointChange:   &points,
			OperationType: &operationType,
			Reason:        &reason,
			Remark:        &remark,
		}
		if relatedOrderID > 0 {
			oid := int(relatedOrderID)
			record.RelatedOrderId = &oid
		}
		if err := pointRecordService.CreatePointRecord(context.Background(), record); err != nil {
			global.GVA_LOG.Error("营销奖励积分发放失败", zap.String("type", triggerType), zap.Error(err))
		}
	}

	// 发放优惠券
	if reward.CouponIDs != "" {
		s.issueCouponsToUser(userID, reward.CouponIDs)
	}

	return nil
}

// issueCouponsToUser 为用户发放奖励优惠券
func (s *MarketingRewardService) issueCouponsToUser(userID uint, couponIDsStr string) {
	ids := strings.Split(couponIDsStr, ",")
	for _, idStr := range ids {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		couponID, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}

		// 检查优惠券有效
		var coupon shop.Coupon
		if err := global.GVA_DB.Where("id = ? AND status = ?", couponID, true).First(&coupon).Error; err != nil {
			global.GVA_LOG.Warn("奖励优惠券无效", zap.Int("couponID", couponID))
			continue
		}

		// 检查库存
		if coupon.Quantity != nil && coupon.Claimed != nil && *coupon.Claimed >= *coupon.Quantity {
			global.GVA_LOG.Warn("奖励优惠券库存不足", zap.Int("couponID", couponID))
			continue
		}

		// 生成券码
		snowflakeID, err := generateSnowflakeID()
		if err != nil {
			global.GVA_LOG.Error("生成券码失败", zap.Error(err))
			continue
		}

		status := false
		uid := int(userID)
		now := time.Now()
		couponOrderUser := shop.CouponOrderUser{
			CouponNum:  snowflakeID,
			CouponID:   &couponID,
			UserID:     userID,
			ShopUserID: &uid,
			Status:     &status,
			ClaimedAt:  &now,
		}
		if err := global.GVA_DB.Create(&couponOrderUser).Error; err != nil {
			global.GVA_LOG.Error("奖励优惠券发放失败", zap.Error(err))
			continue
		}

		// 更新已领取数
		global.GVA_DB.Model(&shop.Coupon{}).Where("id = ?", couponID).UpdateColumn("claimed", gorm.Expr("claimed + 1"))
	}
}

// TriggerSubOrderReward 下级首单付款，奖励上级
// 当 sub_register 配置了 SubRequireOrder=true 时，下级首次付款完成后给上级发放奖励
func (s *MarketingRewardService) TriggerSubOrderReward(tx *gorm.DB, buyerUserID uint, orderID uint) {
	if tx == nil {
		tx = global.GVA_DB
	}

	// 查询买家用户信息
	var buyer client.ClientUser
	if err := tx.Where("id = ?", buyerUserID).First(&buyer).Error; err != nil {
		return
	}

	// 没有邀请人或已奖励过
	if buyer.InvitedBy == 0 || buyer.SubOrderRewarded {
		return
	}

	// 查询 sub_register 配置
	reward, err := s.GetMarketingRewardByType("sub_register")
	if err != nil {
		return // 未配置
	}

	// 只有 SubRequireOrder=true 时才在这里触发
	if reward.SubRequireOrder == nil || !*reward.SubRequireOrder {
		return // 注册时已发放
	}

	// 给邀请人发放奖励
	if err := s.TriggerReward(buyer.InvitedBy, "sub_register", "invite_reward", "下级用户首单完成奖励", orderID); err != nil {
		global.GVA_LOG.Error("下级首单奖励发放失败", zap.Error(err))
	}

	// 标记为已奖励
	tx.Model(&client.ClientUser{}).Where("id = ?", buyerUserID).Update("sub_order_rewarded", true)
}
