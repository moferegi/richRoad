package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// SKU规格字典 结构体 SkuSpec
// 用于管理可复用的规格/属性键值，支持多语言
type SkuSpec struct {
	global.GVA_MODEL
	Type      string `json:"type" form:"type" gorm:"column:type;size:20;default:spec;comment:类型(spec规格/attr属性);" binding:"required"` //类型
	Label     string `json:"label" form:"label" gorm:"column:label;size:191;comment:默认名称;" binding:"required"`                       //默认名称
	LabelI18n string `json:"labelI18n" form:"labelI18n" gorm:"column:label_i18n;type:text;comment:名称多语言JSON;"`                       //名称多语言
	Value     string `json:"value" form:"value" gorm:"column:value;size:191;comment:默认值;" binding:"required"`                        //默认值
	ValueI18n string `json:"valueI18n" form:"valueI18n" gorm:"column:value_i18n;type:text;comment:值多语言JSON;"`                        //值多语言
	Sort      *int   `json:"sort" form:"sort" gorm:"column:sort;default:0;comment:排序;"`                                              //排序
}

// TableName SKU规格字典 SkuSpec自定义表名 shop_sku_spec
func (SkuSpec) TableName() string {
	return "shop_sku_spec"
}
