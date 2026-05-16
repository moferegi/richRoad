// 自动生成模板Sku
package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/datatypes"
)

// sku 结构体  Sku
type Sku struct {
	global.GVA_MODEL
	Name                string         `json:"name" form:"name" gorm:"column:name;comment:名称(JSON多语言);type:text;" binding:"required"`                           //名称(JSON多语言)
	Picture             string         `json:"picture" form:"picture" gorm:"column:picture;comment:图片;"`                                                        //图片
	ExternalPicturePath string         `json:"externalPicturePath" form:"externalPicturePath" gorm:"column:external_picture_path;size:500;comment:外部图片路径(优先);"` //外部图片路径
	UpperImage          string         `json:"upperImage" form:"upperImage" gorm:"column:upper_image;size:500;comment:上装图;"`                                    //上装图
	LowerImage          string         `json:"lowerImage" form:"lowerImage" gorm:"column:lower_image;size:500;comment:下装图;"`                                    //下装图
	Description         string         `json:"description" form:"description" gorm:"column:description;comment:介绍(JSON多语言);type:text;"`                         //介绍(JSON多语言)
	Price               uint           `json:"price" form:"price" gorm:"column:price;comment:价格;"`                                                              //价格
	PriceI18n           string         `json:"priceI18n" form:"priceI18n" gorm:"column:price_i18n;comment:多语言价格JSON;type:text;"`                                //多语言价格(JSON)
	Inventory           uint           `json:"inventory" form:"inventory" gorm:"column:inventory;comment:余量;"`                                                  //余量
	Specs               datatypes.JSON `json:"specs" form:"specs" gorm:"column:specs;comment:规格;type:text;"`                                                    //规格
	Attrs               datatypes.JSON `json:"attrs" form:"attrs" gorm:"column:attrs;comment:属性;type:text;"`                                                    //属性
	GoodID              uint           `json:"goodID" form:"goodID" gorm:"column:good_id;comment:商品ID;" binding:"required"`                                     //商品ID
	SaleNum             uint           `json:"saleNum" form:"saleNum" gorm:"column:sale_num;comment:销量;"`                                                       //销量
}

// TableName sku Sku自定义表名 shop_sku
func (Sku) TableName() string {
	return "shop_sku"
}
