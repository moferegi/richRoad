package model

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// CsConfig 客服系统全局配置（单例，ID=1）
type CsConfig struct {
	global.GVA_MODEL
	PlatEnabled bool `json:"platEnabled" gorm:"column:plat_enabled;default:false;comment:是否启用平台内置客服系统"`
}

func (CsConfig) TableName() string {
	return "cs_configs"
}
