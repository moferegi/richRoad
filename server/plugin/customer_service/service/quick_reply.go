package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model/request"
)

type QuickReplyService struct{}

// Create 创建快捷回复
func (s *QuickReplyService) Create(qr *model.CsQuickReply) error {
	return global.GVA_DB.Create(qr).Error
}

// Update 更新快捷回复
func (s *QuickReplyService) Update(qr *model.CsQuickReply) error {
	return global.GVA_DB.Save(qr).Error
}

// Delete 删除快捷回复
func (s *QuickReplyService) Delete(id uint) error {
	return global.GVA_DB.Delete(&model.CsQuickReply{}, id).Error
}

// GetList 分页获取快捷回复列表
func (s *QuickReplyService) GetList(search request.QuickReplySearch) ([]model.CsQuickReply, int64, error) {
	db := global.GVA_DB.Model(&model.CsQuickReply{})
	if search.Title != "" {
		db = db.Where("title LIKE ?", "%"+search.Title+"%")
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.CsQuickReply
	if err := db.Order("sort ASC, id ASC").Scopes(search.Paginate()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// GetAll 获取全部快捷回复（供坐席工作台使用）
func (s *QuickReplyService) GetAll() ([]model.CsQuickReply, error) {
	var list []model.CsQuickReply
	err := global.GVA_DB.Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}
