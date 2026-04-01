package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
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
