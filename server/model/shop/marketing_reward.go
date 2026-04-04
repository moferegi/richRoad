package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// MarketingReward 营销奖励配置 结构体
type MarketingReward struct {
	global.GVA_MODEL
	TriggerType     string `json:"triggerType" form:"triggerType" gorm:"column:trigger_type;size:30;uniqueIndex;comment:触发类型(register/sub_register/sign_in/order);"` //触发类型
	IsEnabled       *bool  `json:"isEnabled" form:"isEnabled" gorm:"column:is_enabled;default:false;comment:是否启用;"`                                                  //是否启用
	Points          *int   `json:"points" form:"points" gorm:"column:points;default:0;comment:奖励积分;"`                                                                //奖励积分
	CouponIDs       string `json:"couponIDs" form:"couponIDs" gorm:"column:coupon_ids;size:500;comment:奖励优惠券ID(逗号分隔);"`                                              //奖励优惠券ID
	SubRequireOrder *bool  `json:"subRequireOrder" form:"subRequireOrder" gorm:"column:sub_require_order;default:false;comment:是否要求下级首单付款;"`                         //是否要求下级首单付款
	OrderOnce       *bool  `json:"orderOnce" form:"orderOnce" gorm:"column:order_once;default:false;comment:下单奖励仅发放一次;"`                                             //下单奖励仅发放一次
}

// TableName 营销奖励配置自定义表名
func (MarketingReward) TableName() string {
	return "shop_marketing_reward"
}
