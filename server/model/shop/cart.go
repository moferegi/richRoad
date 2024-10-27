// 自动生成模板Cart
package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 购物车 结构体  Cart
type Cart struct {
	global.GVA_MODEL
	UserID   uint `json:"userID" form:"userID" gorm:"column:user_id;comment:用户ID;"`                    //用户ID
	GoodID   uint `json:"goodID" form:"goodID" gorm:"column:good_id;comment:商品ID;" binding:"required"` //商品ID
	Good     Good `json:"good" form:"good" gorm:"->;foreignKey:GoodID;references:ID"`                  //商品
	SKUID    uint `json:"skuID" form:"skuID" gorm:"column:sku_id;comment:skuID;" binding:"required"`   //skuID
	SKU      Sku  `json:"sku" form:"sku" gorm:"->;foreignKey:SKUID;references:ID"`                     //sku
	Quantity uint `json:"quantity" form:"quantity" gorm:"column:quantity;comment:数量;"`                 //数量
}

// TableName 购物车 Cart自定义表名 shop_cart
func (Cart) TableName() string {
	return "shop_cart"
}
