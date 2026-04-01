package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type CouponOrderUserSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	CouponID       *int       `json:"couponID" form:"couponID"`
	OrderID        *int       `json:"orderID" form:"orderID"`
	Status         *bool      `json:"status" form:"status"`
	UserID         *uint      `json:"userID" form:"userID"`
	request.PageInfo
}
