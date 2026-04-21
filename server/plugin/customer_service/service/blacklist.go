package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model/request"
)

type BlacklistService struct{}

// Add 加入黑名单
func (s *BlacklistService) Add(clientUserID, operatorID uint, reason string) error {
	bl := model.CsBlacklist{
		ClientUserID: clientUserID,
		Reason:       reason,
		CreatedBy:    operatorID,
	}
	return global.GVA_DB.Create(&bl).Error
}

// Remove 移出黑名单
func (s *BlacklistService) Remove(id uint) error {
	return global.GVA_DB.Delete(&model.CsBlacklist{}, id).Error
}

// GetList 分页获取黑名单
func (s *BlacklistService) GetList(search request.BlacklistSearch) ([]model.CsBlacklist, int64, error) {
	db := global.GVA_DB.Model(&model.CsBlacklist{})
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.CsBlacklist
	if err := db.Order("created_at DESC").Scopes(search.Paginate()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
