package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// PhoneAreaCode 国际区号 结构体
type PhoneAreaCode struct {
	global.GVA_MODEL
	CountryName string `json:"countryName" form:"countryName" gorm:"column:country_name;type:text;comment:国家名称(JSON多语言);"` //国家名称(JSON多语言)
	AreaCode    string `json:"areaCode" form:"areaCode" gorm:"column:area_code;size:10;comment:区号;" binding:"required"`     //区号(+976)
	PhoneRegex  string `json:"phoneRegex" form:"phoneRegex" gorm:"column:phone_regex;size:200;comment:手机号正则验证;"`            //手机号正则
	IsEnabled   *bool  `json:"isEnabled" form:"isEnabled" gorm:"column:is_enabled;default:true;comment:是否启用;"`              //是否启用
	Sort        *int   `json:"sort" form:"sort" gorm:"column:sort;default:0;comment:排序;"`                                   //排序
	FlagIcon    string `json:"flagIcon" form:"flagIcon" gorm:"column:flag_icon;size:500;comment:国旗图标;"`                      //国旗图标
}

// TableName 国际区号自定义表名
func (PhoneAreaCode) TableName() string {
	return "client_phone_area_code"
}
