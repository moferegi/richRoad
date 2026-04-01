package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// SysLanguage 系统语言 结构体
type SysLanguage struct {
	global.GVA_MODEL
	Code       string `json:"code" form:"code" gorm:"column:code;uniqueIndex;size:10;comment:语言代码;" binding:"required"` //语言代码(mn,zh,en...)
	Name       string `json:"name" form:"name" gorm:"column:name;size:50;comment:语言名称;" binding:"required"`             //语言名称(蒙古国语)
	NativeName string `json:"nativeName" form:"nativeName" gorm:"column:native_name;size:50;comment:原生名称;"`             //原生名称(Монгол)
	IsEnabled  *bool  `json:"isEnabled" form:"isEnabled" gorm:"column:is_enabled;default:true;comment:是否启用;"`           //是否启用
	IsDefault  *bool  `json:"isDefault" form:"isDefault" gorm:"column:is_default;default:false;comment:是否默认;"`          //是否默认
	Sort       *int   `json:"sort" form:"sort" gorm:"column:sort;default:0;comment:排序;"`                                //排序
}

// TableName 系统语言自定义表名
func (SysLanguage) TableName() string {
	return "client_sys_language"
}
