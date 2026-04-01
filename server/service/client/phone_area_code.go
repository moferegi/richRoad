package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
)

type PhoneAreaCodeService struct{}

// CreatePhoneAreaCode 创建国际区号
func (s *PhoneAreaCodeService) CreatePhoneAreaCode(code *client.PhoneAreaCode) (err error) {
	err = global.GVA_DB.Create(code).Error
	return err
}

// DeletePhoneAreaCode 删除国际区号
func (s *PhoneAreaCodeService) DeletePhoneAreaCode(ID string) (err error) {
	err = global.GVA_DB.Delete(&client.PhoneAreaCode{}, "id = ?", ID).Error
	return err
}

// UpdatePhoneAreaCode 更新国际区号
func (s *PhoneAreaCodeService) UpdatePhoneAreaCode(code client.PhoneAreaCode) (err error) {
	err = global.GVA_DB.Model(&client.PhoneAreaCode{}).Where("id = ?", code.ID).Updates(&code).Error
	return err
}

// GetPhoneAreaCode 根据ID获取
func (s *PhoneAreaCodeService) GetPhoneAreaCode(ID string) (code client.PhoneAreaCode, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&code).Error
	return
}

// GetPhoneAreaCodeList 分页获取列表
func (s *PhoneAreaCodeService) GetPhoneAreaCodeList(info clientReq.PhoneAreaCodeSearch) (list []client.PhoneAreaCode, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&client.PhoneAreaCode{})

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

// GetEnabledPhoneAreaCodes 获取已启用的区号(客户端用)
func (s *PhoneAreaCodeService) GetEnabledPhoneAreaCodes() (list []client.PhoneAreaCode, err error) {
	err = global.GVA_DB.Where("is_enabled = ?", true).Order("sort ASC").Find(&list).Error
	return
}
