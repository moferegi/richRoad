package model

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// CsConfig 客服系统全局配置（单例，ID=1）
type CsConfig struct {
	global.GVA_MODEL
	PlatEnabled      bool   `json:"platEnabled" gorm:"column:plat_enabled;default:false;comment:是否启用平台内置客服系统"`
	AutoCloseMinutes int    `json:"autoCloseMinutes" gorm:"column:auto_close_minutes;default:30;comment:会话无消息多少分钟后自动关闭(0=不自动关闭)"`
	UploadMaxSizeMB  int    `json:"uploadMaxSizeMB" gorm:"column:upload_max_size_mb;default:5;comment:客服图片上传单文件大小上限(MB)"`
	UploadAllowExt   string `json:"uploadAllowExt" gorm:"column:upload_allow_ext;type:varchar(255);default:'jpg,jpeg,png,webp,gif';comment:客服图片允许上传的扩展名(逗号分隔)"`
	DefaultAvatarURL string `json:"defaultAvatarUrl" gorm:"column:default_avatar_url;type:varchar(512);comment:客服默认头像URL(为空时使用首字母随机底色)"`
}

func (CsConfig) TableName() string {
	return "cs_configs"
}
