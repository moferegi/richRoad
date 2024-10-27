package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	wx_global "github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/utils"
	service "github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	utils2 "github.com/wechatpay-apiv3/wechatpay-go/utils"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type WxpayService struct{}

var shopService = service.ServiceGroupApp.ShopServiceGroup.OrderService

// GetPayParams
func (e *WxpayService) GetPayParams(order model.Order) (err error, payParams string) {
	err, ctx, wxClient := utils.CreateClientAndCtx()
	if err != nil {
		return err, ""
	}
	err, payParams = payOneParams(ctx, wxClient, order)
	return err, payParams
}

func payOneParams(ctx context.Context, wxClient *core.Client, order model.Order) (err error, payParams string) {
	// 下单用户ID
	// 下单单号（总ID 用6位（100000）开始记录）
	// 下单产品名(Description)
	// 下单价格(Total)分
	// 得到prepay_id，以及调起支付所需的参数和签名
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var o shop.Order
		fe := tx.First(&o, "id = ?", order.OrderID).Error
		if fe != nil {
			return fe
		}
		if o.Status != "0" {
			return errors.New("订单已支付或已关闭")
		}
		//先创建一个空的订单
		//然后 订单一定要记录好订单上面关联的产品ID 包含产品的介绍 名字 金额等
		preKey := "test_wxshop"
		no := strconv.Itoa(int(order.OrderID))
		b := o.TotalPrice // 产品价格  分
		timestamp := strconv.Itoa(int(time.Now().Unix()))
		Name := "订单支付"
		payOrderID := fmt.Sprintf("%s_%s_%s", preKey, no, timestamp)
		ctx = context.WithValue(ctx, "no", payOrderID)
		svc := jsapi.JsapiApiService{Client: wxClient}
		resp, result, err := svc.Prepay(ctx,
			jsapi.PrepayRequest{
				Appid:         core.String(wx_global.GlobalConfig.AppID),
				Mchid:         core.String(wx_global.GlobalConfig.MchID),
				Payer:         &jsapi.Payer{Openid: core.String(order.Openid)},
				Description:   core.String(Name),       // 产品名字
				OutTradeNo:    core.String(payOrderID), // 预下单的那个ID 必须大于6位 建议 ID自增从 100000开始
				TimeExpire:    core.Time(time.Now()),
				Attach:        core.String("自定义数据说明"),                        //可以不写
				NotifyUrl:     core.String(wx_global.GlobalConfig.NotifyUrl), // 回调url
				SupportFapiao: core.Bool(false),                              // 开不开发票
				Amount: &jsapi.Amount{
					Currency: core.String("CNY"),
					//Total:    core.Int64(commodity.Place), //商品价格(分)
					Total: core.Int64(int64(b)),
				},
			},
		)
		fmt.Println(result)
		if err != nil {
			log.Println(err)
			return err
		}
		var shopOrder shop.Order
		sfe := tx.First(&shopOrder, "id = ?", order.OrderID).Update("pay_order_id", payOrderID).Error
		if sfe != nil {
			return sfe
		}
		payParams = *resp.PrepayId

		var payOrder model.Order
		payOrder.Appid = wx_global.GlobalConfig.AppID
		payOrder.Mchid = wx_global.GlobalConfig.MchID
		payOrder.OutTradeNo = payOrderID
		payOrder.CustomerID = order.CustomerID
		payOrder.OrderID = order.OrderID
		payOrder.TradeState = "NOTPAY"
		payOrder.Openid = order.Openid
		payOrder.Total = int(b)
		payOrder.PayerTotal = int(b)
		payOrder.Currency = "CNY"
		payOrder.PayerCurrency = "CNY"
		perr := tx.Create(&payOrder).Error
		if perr != nil {
			return perr
		}
		return err
	})
	return err, payParams
}

// 关闭订单  用于紧急关闭订单操作 传入订单ID即可
func (e *WxpayService) ClosePayCode(order model.Order) (err error) {
	err, ctx, client := utils.CreateClientAndCtx()
	if err != nil {
		return err
	}
	err = closeOrder(ctx, client, order)
	return err
}

func closeOrder(ctx context.Context, client *core.Client, order model.Order) error {
	svc := native.NativeApiService{Client: client}

	var shopOrder shop.Order
	sfe := global.GVA_DB.First(&shopOrder, "id = ?", order.OrderID).Error
	if sfe != nil {
		return sfe
	}

	result, err := svc.CloseOrder(ctx,
		native.CloseOrderRequest{
			OutTradeNo: core.String(shopOrder.PayOrderID),
			Mchid:      core.String(wx_global.GlobalConfig.MchID),
		},
	)
	if err != nil {
		// 处理错误
		log.Printf("call CloseOrder err:%s", err)
		return err
	} else {
		// 处理返回结果
		log.Printf("status=%d", result.Response.StatusCode)
		return nil
	}
}

// 查询订单 传入订单ID即可
func (e *WxpayService) GetOrderById(orderID string) (error, model.Order) {
	err, ctx, client := utils.CreateClientAndCtx()
	if err != nil {
		return err, model.Order{}
	}
	return queryOrderByOutTradeNo(ctx, client, orderID)
}

func queryOrderByOutTradeNo(ctx context.Context, client *core.Client, orderID string) (err error, order model.Order) {
	svc := native.NativeApiService{Client: client}

	var shopOrder shop.Order
	sfe := global.GVA_DB.First(&shopOrder, "id = ?", orderID).Error
	if sfe != nil {
		return sfe, model.Order{}
	}

	resp, result, err := svc.QueryOrderByOutTradeNo(ctx,
		native.QueryOrderByOutTradeNoRequest{
			OutTradeNo: core.String(shopOrder.PayOrderID),
			Mchid:      core.String(wx_global.GlobalConfig.MchID),
		},
	)
	if err != nil {
		// 错误处理
		log.Printf("call QueryOrderByOutTradeNo err:%s", err)
		return err, order
	}

	// TradeState
	//SUCCESS：支付成功
	//REFUND：转入退款
	//NOTPAY：未支付
	//CLOSED：已关闭
	//REVOKED：已撤销（仅付款码支付会返回）
	//USERPAYING：用户支付中（仅付款码支付会返回）
	//PAYERROR：支付失败（仅付款码支付会返回）
	order.TradeState = *resp.TradeState
	if order.TradeState == "SUCCESS" {
		var payOrder model.Order
		poe := global.GVA_DB.First(&payOrder, "out_trade_no = ?", *resp.OutTradeNo).Update("trade_state", "SUCCESS").Error
		if poe != nil {
			return poe, order
		}
		soe := shopService.UpdateOrderStatus(nil, orderID, "1")
		if soe != nil {
			return soe, order
		}
	}

	// 处理错误
	// 可以根据订单返回的结果做一些业务逻辑

	log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
	return err, order
}

func (e *WxpayService) PayAction(pay model.PayAction) error {
	p, err := utils2.DecryptAES256GCM(wx_global.GlobalConfig.MchAPIv3Key, pay.Resource.AssociatedData, pay.Resource.Nonce, pay.Resource.Ciphertext)
	if err != nil {
		global.GVA_LOG.Info(p)
		global.GVA_LOG.Info(err.Error())
		return err
	}
	var payOrder model.PayOrder
	if err != nil {
		global.GVA_LOG.Info(p)
		global.GVA_LOG.Info(err.Error())
		return err
	}
	err = json.Unmarshal([]byte(p), &payOrder)
	// payOrder 为回调信息 请自行根据回调信息做业务逻辑
	if err != nil {
		global.GVA_LOG.Info(p)
		global.GVA_LOG.Info(err.Error())
		return err
	}

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		err = tx.First(&model.Order{}, "out_trade_no = ?", payOrder.OutTradeNo).
			Update("trade_state", payOrder.TradeState).
			Update("trade_state_desc", payOrder.TradeStateDesc).
			Update("transaction_id", payOrder.TransactionId).
			Update("trade_type", payOrder.TradeType).
			Update("bank_type", payOrder.BankType).
			Update("attach", payOrder.Attach).
			Update("success_time", payOrder.SuccessTime).
			Error
		if err != nil {
			return err
		}
		if payOrder.TradeState == "SUCCESS" {
			var shopOrder shop.Order
			err = tx.First(&shopOrder, "pay_order_id = ?", payOrder.OutTradeNo).Error
			if err != nil {
				return err
			}
			err = shopService.UpdateOrderStatus(tx, strconv.Itoa(int(shopOrder.ID)), "1")
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
