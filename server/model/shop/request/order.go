package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type OrderSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	Status         string     `json:"status" form:"status"`
	UserID         *int       `json:"userID" form:"userID" `
	request.PageInfo
}

// PlaceOrderByCartRequest 购物车下单请求
type PlaceOrderByCartRequest struct {
	UsePoints bool `json:"usePoints" form:"usePoints"` // 是否使用积分抵扣
}
