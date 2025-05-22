
// 自动生成模板Promotion
package shop
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 促销信息 结构体  Promotion
type Promotion struct {
    global.GVA_MODEL
  PromotionImage  string `json:"promotionImage" form:"promotionImage" gorm:"column:promotion_name;" binding:"required"`  //促销背景图
  Title  *string `json:"title" form:"title" gorm:"column:title;" binding:"required"`  //促销标题
  Description  *string `json:"description" form:"description" gorm:"column:description;"`  //促销文案
  IsAction  *bool `json:"isAction" form:"isAction" gorm:"column:is_action;" binding:"required"`  //是否启用
  Category  *string `json:"category" form:"category" gorm:"column:category;"`  //促销类别
}


// TableName 促销信息 Promotion自定义表名 shop_Promotions
func (Promotion) TableName() string {
    return "shop_Promotions"
}





