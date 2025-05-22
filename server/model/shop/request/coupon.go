package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type CouponSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	Type           string     `json:"type" form:"type"`
	request.PageInfo
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}

// AdminIssueCouponRequest 管理员发券请求结构体
type AdminIssueCouponRequest struct {
	CouponID  int    `json:"couponID" form:"couponID" binding:"required"`   // 优惠券ID
	UserCount int    `json:"userCount" form:"userCount" binding:"required"` // 发放用户数量
	UserIds   []int  `json:"userIds" form:"userIds"`                        // 用户ID列表
	UserLimit string `json:"userLimit" form:"userLimit"`                    // 用户筛选条件，可选
}

// ClaimCouponRequest 用户领取优惠券请求结构体
type ClaimCouponRequest struct {
	CouponID int `json:"couponID" form:"couponID" binding:"required"` // 优惠券ID
}

type GetClaimCouponRequest struct {
	GoodIds []int `json:"goodIds" form:"goodIds"`
}
