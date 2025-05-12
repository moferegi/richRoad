package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type CouponOrderUserSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	CouponID       *int       `json:"couponID" form:"couponID"`
	OrderID        *int       `json:"orderID" form:"orderID"`
	request.PageInfo
}
