package client

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"gorm.io/gorm"
)

type ExternalLinkDomainService struct{}

// CreateExternalLinkDomain 创建外部链接域名
func (s *ExternalLinkDomainService) CreateExternalLinkDomain(domain client.ExternalLinkDomain) error {
	return global.GVA_DB.Create(&domain).Error
}

// DeleteExternalLinkDomain 删除外部链接域名
func (s *ExternalLinkDomainService) DeleteExternalLinkDomain(id uint) error {
	return global.GVA_DB.Delete(&client.ExternalLinkDomain{}, id).Error
}

// UpdateExternalLinkDomain 更新外部链接域名
func (s *ExternalLinkDomainService) UpdateExternalLinkDomain(domain client.ExternalLinkDomain) error {
	// 如果设为默认，先将其他的取消默认
	if domain.IsDefault != nil && *domain.IsDefault {
		global.GVA_DB.Model(&client.ExternalLinkDomain{}).Where("id != ?", domain.ID).Update("is_default", false)
	}
	return global.GVA_DB.Model(&client.ExternalLinkDomain{}).Where("id = ?", domain.ID).Updates(&domain).Error
}

// GetExternalLinkDomainList 分页获取外部链接域名列表
func (s *ExternalLinkDomainService) GetExternalLinkDomainList(info clientReq.ExternalLinkDomainSearch) (list []client.ExternalLinkDomain, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&client.ExternalLinkDomain{})

	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.IsEnabled != nil {
		db = db.Where("is_enabled = ?", *info.IsEnabled)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("sort DESC, id ASC").Limit(limit).Offset(offset).Find(&list).Error
	return
}

// GetDefaultDomain 获取默认域名
func (s *ExternalLinkDomainService) GetDefaultDomain() (string, error) {
	var domain client.ExternalLinkDomain
	err := global.GVA_DB.Where("is_default = ? AND is_enabled = ?", true, true).First(&domain).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有默认的，取第一个启用的
			err = global.GVA_DB.Where("is_enabled = ?", true).Order("sort DESC, id ASC").First(&domain).Error
			if err != nil {
				return "", nil // 没有任何域名配置，返回空
			}
		} else {
			return "", err
		}
	}
	return domain.Domain, nil
}

// SetDefaultDomain 设置默认域名
func (s *ExternalLinkDomainService) SetDefaultDomain(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 先取消所有默认
		if err := tx.Model(&client.ExternalLinkDomain{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			return err
		}
		// 设置新默认并确保启用
		return tx.Model(&client.ExternalLinkDomain{}).Where("id = ?", id).Updates(map[string]interface{}{
			"is_default": true,
			"is_enabled": true,
		}).Error
	})
}
