package client

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/datatypes"
)

// TryonRechargeOrder 试衣币充值订单
// 状态: 0-待支付 8-待后台确认 1-已支付 4-已取消
type TryonRechargeOrder struct {
	global.GVA_MODEL
	UserID       uint           `json:"userID" form:"userID" gorm:"column:user_id;index;comment:用户ID;"`
	Status       string         `json:"status" form:"status" gorm:"column:status;size:20;default:0;index;comment:订单状态(0待支付/8待后台确认/1已支付/4已取消);"`
	Points       int            `json:"points" form:"points" gorm:"column:points;comment:充值试衣币数量;"`
	Amount       int            `json:"amount" form:"amount" gorm:"column:amount;comment:充值金额(分);"`
	Currency     string         `json:"currency" form:"currency" gorm:"column:currency;size:20;default:CNY;comment:货币单位;"`
	PayMethod    string         `json:"payMethod" form:"payMethod" gorm:"column:pay_method;size:40;comment:支付方式;"`
	PlanSnapshot datatypes.JSON `json:"planSnapshot" form:"planSnapshot" gorm:"column:plan_snapshot;type:json;comment:套餐快照;"`
	CloseTime    time.Time      `json:"closeTime" form:"closeTime" gorm:"column:close_time;index;comment:关闭时间;"`
	PaidAt       *time.Time     `json:"paidAt" form:"paidAt" gorm:"column:paid_at;comment:确认支付时间;"`
	CancelledAt  *time.Time     `json:"cancelledAt" form:"cancelledAt" gorm:"column:cancelled_at;comment:取消时间;"`
	OutTradeNo   string         `json:"outTradeNo" form:"outTradeNo" gorm:"column:out_trade_no;size:128;index;comment:商户订单号;"`
	Remark       string         `json:"remark" form:"remark" gorm:"column:remark;size:255;comment:备注;"`
}

func (TryonRechargeOrder) TableName() string {
	return "client_tryon_recharge_order"
}
