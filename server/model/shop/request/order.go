package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type OrderSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	Status         string     `json:"status" form:"status"`
	UserID         *int       `json:"userID" form:"userID"`
	GoodID         *int       `json:"goodID" form:"goodID"`
	IsPresale      *bool      `json:"isPresale" form:"isPresale"`
	PayMethod      string     `json:"payMethod" form:"payMethod"`
	Sort           string     `json:"sort" form:"sort"`
	Order          string     `json:"order" form:"order"`
	request.PageInfo
}

// BatchUpdateOrderStatusReq 批量更新订单状态请求
type BatchUpdateOrderStatusReq struct {
	IDs    []string `json:"IDs" binding:"required"`    // 订单ID列表
	Status string   `json:"status" binding:"required"` // 目标状态
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
	CartIDs   []uint `json:"cartIDs" form:"cartIDs"`     // 购物车项ID(为空则全部)
	AddressID uint   `json:"addressID" form:"addressID"` // 收货地址ID
	UsePoints bool   `json:"usePoints" form:"usePoints"` // 是否使用积分抵扣
	PayMethod string `json:"payMethod" form:"payMethod"` // 支付方式(contact/qrcode/wechat/alipay/bank_card_cn/bank_card_us/bank_card_mn/paypal)
	CouponNum string `json:"couponNum" form:"couponNum"` // 优惠券号
}

// PresaleListRequest 预售商品列表请求
type PresaleListRequest struct {
	Limit          int    `json:"limit" form:"limit"`                   // 限制数量(首页用)
	PresaleEnabled *bool  `json:"presaleEnabled" form:"presaleEnabled"` // 预售开关筛选
	OrderBy        string `json:"orderBy" form:"orderBy"`
	OrderDir       string `json:"orderDir" form:"orderDir"`
	request.PageInfo
}
