package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model"
)

type ConfigService struct{}

// GetConfig 获取客服全局配置（不存在时返回默认值）
func (s *ConfigService) GetConfig() (model.CsConfig, error) {
	var cfg model.CsConfig
	err := global.GVA_DB.First(&cfg).Error
	if err != nil {
		// 未初始化时返回默认配置
		return model.CsConfig{PlatEnabled: false}, nil
	}
	return cfg, nil
}

// UpdatePlatEnabled 更新平台客服开关
func (s *ConfigService) UpdatePlatEnabled(enabled bool) error {
	var cfg model.CsConfig
	result := global.GVA_DB.First(&cfg)
	if result.Error != nil {
		// 首次创建
		cfg = model.CsConfig{PlatEnabled: enabled}
		return global.GVA_DB.Create(&cfg).Error
	}
	return global.GVA_DB.Model(&cfg).Update("plat_enabled", enabled).Error
}
