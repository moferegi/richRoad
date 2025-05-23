// 自动生成模板Coupon
package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// 优惠券 结构体  Coupon
type Coupon struct {
	global.GVA_MODEL
	Name        *string    `json:"name" form:"name" gorm:"comment:优惠券名称;column:name;" binding:"required"`                      //名称
	Description *string    `json:"description" form:"description" gorm:"comment:优惠券描述;column:description;" binding:"required"` //描述
	Discount    uint       `json:"discount" form:"discount" gorm:"comment:折扣金额;column:discount;" binding:"required"`           //折扣
	MinSpend    uint       `json:"minSpend" form:"minSpend" gorm:"comment:最低消费金额;column:min_spend;" binding:"required"`        //最低消费
	ProductID   *int       `json:"productID" form:"productID" gorm:"comment:关联商品ID（仅商品券适用）;column:product_id;"`                //商品ID
	Quantity    *int       `json:"quantity" form:"quantity" gorm:"comment:优惠券总数量;column:quantity;" binding:"required"`         //数量
	Claimed     *int       `json:"claimed" form:"claimed" gorm:"comment:已领取数量;column:claimed;default:0"`                       //已领取
	StartTime   *time.Time `json:"startTime" form:"startTime" gorm:"comment:开始时间;column:start_time;" binding:"required"`       //开始时间
	EndTime     *time.Time `json:"endTime" form:"endTime" gorm:"comment:结束时间;column:end_time;" binding:"required"`             //结束时间
	Status      *bool      `json:"status" form:"status" gorm:"comment:是否启用;column:status;" binding:"required"`                 //启用
}

// TableName 优惠券 Coupon自定义表名 coupon
func (Coupon) TableName() string {
	return "shop_coupon"
}
