package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
)

type QrcodePaymentService struct{}

// CreateQrcodePayment 创建收款码
func (s *QrcodePaymentService) CreateQrcodePayment(qr *shop.QrcodePayment) (err error) {
	err = global.GVA_DB.Create(qr).Error
	return err
}

// DeleteQrcodePayment 删除收款码
func (s *QrcodePaymentService) DeleteQrcodePayment(ID string) (err error) {
	err = global.GVA_DB.Delete(&shop.QrcodePayment{}, "id = ?", ID).Error
	return err
}

// UpdateQrcodePayment 更新收款码
func (s *QrcodePaymentService) UpdateQrcodePayment(qr shop.QrcodePayment) (err error) {
	err = global.GVA_DB.Model(&shop.QrcodePayment{}).Where("id = ?", qr.ID).Updates(&qr).Error
	return err
}

// GetQrcodePayment 根据ID获取收款码
func (s *QrcodePaymentService) GetQrcodePayment(ID string) (qr shop.QrcodePayment, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&qr).Error
	return
}

// GetQrcodePaymentList 分页获取收款码列表
func (s *QrcodePaymentService) GetQrcodePaymentList(info shopReq.QrcodePaymentSearch) (list []shop.QrcodePayment, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&shop.QrcodePayment{})

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

// GetEnabledQrcodePayments 获取已启用的收款码列表(客户端用)
func (s *QrcodePaymentService) GetEnabledQrcodePayments() (list []shop.QrcodePayment, err error) {
	err = global.GVA_DB.Where("is_enabled = ?", true).Order("sort ASC").Find(&list).Error
	return
}
