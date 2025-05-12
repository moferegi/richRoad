// 自动生成模板CouponOrderUser
package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 优惠券 结构体  CouponOrderUser
type CouponOrderUser struct {
	global.GVA_MODEL
	CouponID   *int  `json:"couponID" form:"couponID" gorm:"comment:优惠券ID;column:coupon_id;" binding:"required"` //优惠券ID
	OrderID    *int  `json:"orderID" form:"orderID" gorm:"comment:订单ID;column:order_id;"`                        //订单ID
	CanUse     bool  `json:"canUse" form:"canUse" gorm:"-"`                                                      //可用
	ShopUserID *int  `json:"shopUserID" form:"shopUserID" gorm:"comment:商户用户ID;column:shop_user_id;"`            //商户用户ID
	Status     *bool `json:"status" form:"status" gorm:"comment:是否已使用;column:status;" binding:"required"`        //已使用
}

// TableName 优惠券 CouponOrderUser自定义表名 coupon_order_user
func (CouponOrderUser) TableName() string {
	return "coupon_order_user"
}
