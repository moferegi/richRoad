package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
)

type SysLanguageService struct{}

// CreateSysLanguage 创建语言
func (s *SysLanguageService) CreateSysLanguage(lang *client.SysLanguage) (err error) {
	err = global.GVA_DB.Create(lang).Error
	return err
}

// DeleteSysLanguage 删除语言
func (s *SysLanguageService) DeleteSysLanguage(ID string) (err error) {
	err = global.GVA_DB.Delete(&client.SysLanguage{}, "id = ?", ID).Error
	return err
}

// UpdateSysLanguage 更新语言
func (s *SysLanguageService) UpdateSysLanguage(lang client.SysLanguage) (err error) {
	// 如果设置为默认，先清除其他默认
	if lang.IsDefault != nil && *lang.IsDefault {
		global.GVA_DB.Model(&client.SysLanguage{}).Where("id != ?", lang.ID).Update("is_default", false)
	}
	err = global.GVA_DB.Model(&client.SysLanguage{}).Where("id = ?", lang.ID).Updates(&lang).Error
	return err
}

// GetSysLanguage 根据ID获取语言
func (s *SysLanguageService) GetSysLanguage(ID string) (lang client.SysLanguage, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&lang).Error
	return
}

// GetSysLanguageList 分页获取语言列表
func (s *SysLanguageService) GetSysLanguageList(info clientReq.SysLanguageSearch) (list []client.SysLanguage, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&client.SysLanguage{})

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

// GetEnabledLanguages 获取已启用的语言列表(客户端用)
func (s *SysLanguageService) GetEnabledLanguages() (list []client.SysLanguage, err error) {
	err = global.GVA_DB.Where("is_enabled = ?", true).Order("sort ASC").Find(&list).Error
	return
}
