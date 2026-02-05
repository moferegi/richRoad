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

type RefundApplyReq struct {
	OrderID uint     `json:"orderID" binding:"required"`
	Reason  string   `json:"reason" binding:"required"`
	Images  []string `json:"images"`
}

type RefundHandleReq struct {
	OrderID uint   `json:"orderID" binding:"required"`
	Remark  string `json:"remark"`
}

// PlaceOrderByCartRequest 购物车下单请求
type PlaceOrderByCartRequest struct {
	UsePoints bool `json:"usePoints" form:"usePoints"` // 是否使用积分抵扣
}
