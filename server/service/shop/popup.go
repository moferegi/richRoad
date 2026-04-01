package shop

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
)

type PopupService struct{}

// CreatePopup 创建弹窗
func (s *PopupService) CreatePopup(popup *shop.Popup) (err error) {
	err = global.GVA_DB.Create(popup).Error
	return err
}

// DeletePopup 删除弹窗
func (s *PopupService) DeletePopup(ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.Popup{}, "id = ?", ID).Error
	return err
}

// UpdatePopup 更新弹窗
func (s *PopupService) UpdatePopup(popup shop.Popup) (err error) {
	err = global.GVA_DB.Model(&shop.Popup{}).Where("id = ?", popup.ID).Updates(&popup).Error
	return err
}

// GetPopup 根据ID获取弹窗
func (s *PopupService) GetPopup(ID string) (popup shop.Popup, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&popup).Error
	return
}

// GetPopupList 分页获取弹窗列表(后台管理)
func (s *PopupService) GetPopupList(info shopReq.PopupSearch) (list []shop.Popup, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&shop.Popup{})

	if info.Position != "" {
		db = db.Where("position = ?", info.Position)
	}
	if info.ClientType != "" {
		db = db.Where("client_type = ? OR client_type = 'all'", info.ClientType)
	}
	if info.IsEnabled != nil {
		db = db.Where("is_enabled = ?", *info.IsEnabled)
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
	err = db.Order("sort ASC").Find(&list).Error
	return
}

// GetActivePopups 获取当前有效弹窗(客户端用)
func (s *PopupService) GetActivePopups(position string, clientType string) (list []shop.Popup, err error) {
	now := time.Now()
	db := global.GVA_DB.Where("is_enabled = ?", true)
	if position != "" {
		db = db.Where("position = ? OR position = 'all'", position)
	}
	if clientType != "" {
		db = db.Where("client_type = ? OR client_type = 'all'", clientType)
	}
	db = db.Where("(start_time IS NULL OR start_time <= ?) AND (end_time IS NULL OR end_time >= ?)", now, now)
	err = db.Order("sort ASC").Find(&list).Error
	return
}
