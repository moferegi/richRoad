package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type OrderSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	Status         string     `json:"status" form:"status"`
	UserID         *int       `json:"userID" form:"userID"`
	IsPresale      *bool      `json:"isPresale" form:"isPresale"`
	PayMethod      string     `json:"payMethod" form:"payMethod"`
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
	CartIDs   []uint `json:"cartIDs" form:"cartIDs"`       // 购物车项ID(为空则全部)
	AddressID uint   `json:"addressID" form:"addressID"`   // 收货地址ID
	UsePoints bool   `json:"usePoints" form:"usePoints"`   // 是否使用积分抵扣
	PayMethod string `json:"payMethod" form:"payMethod"`   // 支付方式(contact/qrcode)
	CouponNum string `json:"couponNum" form:"couponNum"`   // 优惠券号
}

// PresaleListRequest 预售商品列表请求
type PresaleListRequest struct {
	Limit int `json:"limit" form:"limit"` // 限制数量(首页用)
	request.PageInfo
}
