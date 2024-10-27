// 自动生成模板Sku
package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/datatypes"
)

// sku 结构体  Sku
type Sku struct {
	global.GVA_MODEL
	Name        string         `json:"name" form:"name" gorm:"column:name;comment:名称;" binding:"required"`          //名称
	Picture     string         `json:"picture" form:"picture" gorm:"column:picture;comment:图片;"`                    //图片
	Description string         `json:"description" form:"description" gorm:"column:description;comment:介绍;"`        //介绍
	Price       uint           `json:"price" form:"price" gorm:"column:price;comment:价格;"`                          //价格
	Inventory   uint           `json:"inventory" form:"inventory" gorm:"column:inventory;comment:余量;"`              //余量
	Specs       datatypes.JSON `json:"specs" form:"specs" gorm:"column:specs;comment:规格;type:text;"`                //规格
	Attrs       datatypes.JSON `json:"attrs" form:"attrs" gorm:"column:attrs;comment:属性;type:text;"`                //属性
	GoodID      uint           `json:"goodID" form:"goodID" gorm:"column:good_id;comment:商品ID;" binding:"required"` //商品ID
	SaleNum     uint           `json:"saleNum" form:"saleNum" gorm:"column:sale_num;comment:销量;"`                   //销量
}

// TableName sku Sku自定义表名 shop_sku
func (Sku) TableName() string {
	return "shop_sku"
}
