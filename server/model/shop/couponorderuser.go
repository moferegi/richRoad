// 自动生成模板CouponOrderUser
package shop

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 优惠券 结构体  CouponOrderUser
type CouponOrderUser struct {
	global.GVA_MODEL
	CouponNum  string     `json:"couponNum" form:"couponNum" gorm:"comment:优惠券编号;column:coupon_num;"`                 //优惠券编号
	CouponID   *int       `json:"couponID" form:"couponID" gorm:"comment:优惠券ID;column:coupon_id;" binding:"required"` //优惠券ID
	OrderID    *int       `json:"orderID" form:"orderID" gorm:"comment:订单ID;column:order_id;"`                        //订单ID
	ShopUserID *int       `json:"shopUserID" form:"shopUserID" gorm:"comment:商户用户ID;column:shop_user_id;"`            //商户用户ID
	UserID     uint       `json:"userID" form:"userID" gorm:"column:user_id;index;comment:用户ID;"`                     //用户ID
	Status     *bool      `json:"status" form:"status" gorm:"comment:是否已使用;column:status;"`                           //已使用
	ClaimedAt  *time.Time `json:"claimedAt" form:"claimedAt" gorm:"column:claimed_at;comment:领取时间;"`                  //领取时间
	UsedAt     *time.Time `json:"usedAt" form:"usedAt" gorm:"column:used_at;comment:使用时间;"`                           //使用时间
}

// TableName 优惠券 CouponOrderUser自定义表名 coupon_order_user
func (CouponOrderUser) TableName() string {
	return "coupon_order_user"
}
