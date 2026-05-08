package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// CreateTryonRechargeOrderReq 创建试衣币充值订单请求
type CreateTryonRechargeOrderReq struct {
	Points    int    `json:"points" form:"points" binding:"required,min=1"`
	Price     string `json:"price" form:"price" binding:"required"`
	PayMethod string `json:"payMethod" form:"payMethod"`
}

// UpdateTryonRechargeOrderPayMethodReq 更新充值订单支付方式请求
type UpdateTryonRechargeOrderPayMethodReq struct {
	ID        uint   `json:"id" form:"id" binding:"required"`
	PayMethod string `json:"payMethod" form:"payMethod" binding:"required"`
}

// SubmitTryonRechargeOrderPaymentReq 提交充值订单付款确认请求
type SubmitTryonRechargeOrderPaymentReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

// TryonRechargeOrderSearch 充值订单查询条件
type TryonRechargeOrderSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	Status         string     `json:"status" form:"status"`
	UserID         uint       `json:"userID" form:"userID"`
	request.PageInfo
}
