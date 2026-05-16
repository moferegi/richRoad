package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// QrcodePayment 二维码收款 结构体
type QrcodePayment struct {
	global.GVA_MODEL
	Name         string `json:"name" form:"name" gorm:"column:name;size:100;comment:收款码名称;" binding:"required"`      //收款码名称
	NameI18n     string `json:"nameI18n" form:"nameI18n" gorm:"column:name_i18n;type:text;comment:收款码名称多语言JSON;"`    //收款码名称多语言JSON
	Image        string `json:"image" form:"image" gorm:"column:image;size:500;comment:二维码图片路径;"`                    //二维码图片路径
	ExternalPath string `json:"externalPath" form:"externalPath" gorm:"column:external_path;size:500;comment:外部路径;"` //外部路径
	IsEnabled    *bool  `json:"isEnabled" form:"isEnabled" gorm:"column:is_enabled;default:true;comment:是否启用;"`      //是否启用
	Sort         *int   `json:"sort" form:"sort" gorm:"column:sort;default:0;comment:排序;"`                           //排序
}

// TableName 二维码收款自定义表名
func (QrcodePayment) TableName() string {
	return "shop_qrcode_payment"
}
