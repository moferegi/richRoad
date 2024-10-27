package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type CommentSearch struct {
	StartCreatedAt   *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt     *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	OrderID          *int       `json:"orderID" form:"orderID" `
	GoodID           *int       `json:"goodID" form:"goodID" `
	SKUID            *int       `json:"SKUID" form:"SKUID" `
	StartShopReplyAt *time.Time `json:"startShopReplyAt" form:"startShopReplyAt"`
	EndShopReplyAt   *time.Time `json:"endShopReplyAt" form:"endShopReplyAt"`
	request.PageInfo
}
