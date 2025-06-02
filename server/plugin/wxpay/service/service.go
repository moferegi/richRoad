package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/payment"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/payment/order/request"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	wx_global "github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	utils2 "github.com/wechatpay-apiv3/wechatpay-go/utils"
	"gorm.io/gorm"
	"log"
	"math"
	"strconv"
)

var shopService = service.ServiceGroupApp.ShopServiceGroup.OrderService

type WxpayService struct{}

func (e *WxpayService) GetPayCode(order model.Order) (err error, CodeUrl string, orderID string) {
	err, ctx, client := utils.CreateClientAndCtx()
	if err != nil {
		return err, "", ""
	}
	err, codeUrl, orderID := payOne(ctx, client, order)
	return err, codeUrl, orderID
}

func (e *WxpayService) CheckNeedPay(order model.Order) (err error, ok bool) {
	var shopOrder shop.Order
	err = global.GVA_DB.First(&shopOrder, "id = ? and user_id = ?", order.OrderID, order.CustomerID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订单不存在"), true
		}
		return err, true
	}
	if shopOrder.TotalPrice == 0 {
		global.GVA_DB.Model(&shopOrder).Update("status", "1")
		return nil, false
	}
	return nil, true
}

func (e *WxpayService) GetPayConf(order model.Order) (err error, payConf any) {
	err, ctx, client := utils.CreateClientAndCtx()
	if err != nil {
		return err, ""
	}
	err, payConf = GetPayConf(ctx, client, order)
	return err, payConf
}

func payOne(ctx context.Context, client *payment.Payment, order model.Order) (err error, codeUrl string, orderID string) {
	// 下单用户ID
	// 下单单号（总ID 用6位（100000）开始记录）
	// 下单产品名(Description)
	// 下单价格(Total)分
	// 得到prepay_id，以及调起支付所需的参数和签名

	//rs, err := client.Security.GetCertificates(ctx)
	//fmt.Println(rs.Data)
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {

		//先创建一个空的订单
		//然后 订单一定要记录好订单上面关联的产品ID 包含产品的介绍 名字 金额等
		snowflakeID := utils.GenerateSnowflakeID() // 获取雪花id
		no := fmt.Sprintf("%d", snowflakeID)       // 转为字符串
		b := int(math.Floor(1))                    // 产品价格  分
		Name := "GVA插件"                            // 订单产品产品名称

		options := &request.RequestNativePrepay{
			Amount: &request.NativeAmount{
				Total:    b,
				Currency: "CNY",
			},
			Attach:      "自定义数据说明",
			Description: Name,
			OutTradeNo:  no,
		}

		response, err := client.Order.TransactionNative(ctx, options)

		if err != nil {
			log.Println(err)
			return err
		}
		orderID = no
		codeUrl = response.CodeURL
		return err
	})
	return err, codeUrl, orderID
}

func GetPayConf(ctx context.Context, client *payment.Payment, order model.Order) (err error, payConf any) {
	// 下单用户ID
	// 下单单号（总ID 用6位（100000）开始记录）
	// 下单产品名(Description)
	// 下单价格(Total)分
	// 得到prepay_id，以及调起支付所需的参数和签名
	var prepayID string
	//rs, err := client.Security.GetCertificates(ctx)
	//fmt.Println(rs.Data)
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
		snowflakeID := utils.GenerateSnowflakeID() // 获取雪花id
		b := o.TotalPrice                          // 产品价格  分
		Name := "订单支付"
		payOrderID := fmt.Sprintf("%d", snowflakeID)

		options := &request.RequestJSAPIPrepay{
			Amount: &request.JSAPIAmount{
				Total:    int(b),
				Currency: "CNY",
			},
			Attach:      "自定义数据说明",
			Description: Name,
			OutTradeNo:  payOrderID,
			Payer: &request.JSAPIPayer{
				OpenID: order.Openid, // 用户的openid， 记得也是动态的。
			},
		}

		response, err := client.Order.JSAPITransaction(ctx, options)
		if err != nil {
			global.GVA_LOG.Error("微信支付下单失败：" + err.Error())
			return err
		}
		if response.ResponseBase.Code != "" {
			global.GVA_LOG.Error("微信支付下单失败：" + response.ResponseBase.Message)
			return errors.New(response.ResponseBase.Message)
		}
		var shopOrder shop.Order
		sfe := tx.First(&shopOrder, "id = ?", order.OrderID).Update("pay_order_id", payOrderID).Error
		if sfe != nil {
			return sfe
		}
		prepayID = response.PrepayID

		var payOrder model.Order
		payOrder.Appid = wx_global.GlobalConfig.AppID
		payOrder.Mchid = wx_global.GlobalConfig.MchID
		payOrder.OutTradeNo = payOrderID
		payOrder.CustomerID = order.CustomerID
		payOrder.TradeState = "NOTPAY"
		payOrder.Openid = order.Openid
		payOrder.Total = int(b)
		payOrder.PayerTotal = int(b)
		payOrder.Currency = "CNY"
		payOrder.PayerCurrency = "CNY"
		perr := tx.Create(&payOrder).Error
		return perr
	})

	payConf, err = client.JSSDK.BridgeConfig(prepayID, false)

	return err, payConf
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

func closeOrder(ctx context.Context, client *payment.Payment, order model.Order) error {
	no := strconv.Itoa(int(order.ID))
	result, err := client.Order.Close(ctx, no)
	if err != nil {
		// 处理错误
		log.Printf("call CloseOrder err:%s", err)
		return err
	} else {
		// 处理返回结果
		global.GVA_LOG.Info(fmt.Sprintf("status = %d", result.ResultCode))
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

func queryOrderByOutTradeNo(ctx context.Context, client *payment.Payment, orderID string) (err error, order model.Order) {

	var o shop.Order

	err = global.GVA_DB.First(&o, "id = ?", orderID).Error

	if err != nil {
		return err, model.Order{}
	}

	result, err := client.Order.QueryByOutTradeNumber(ctx, o.PayOrderID)
	if err != nil {
		// 错误处理
		log.Printf("call QueryOrderByOutTradeNo err:%s", err)
		return err, order
	} else {

		// TradeState
		//SUCCESS：支付成功
		//REFUND：转入退款
		//NOTPAY：未支付
		//CLOSED：已关闭
		//REVOKED：已撤销（仅付款码支付会返回）
		//USERPAYING：用户支付中（仅付款码支付会返回）
		//PAYERROR：支付失败（仅付款码支付会返回）

		order.TradeState = result.TradeState
		if order.TradeState == "SUCCESS" {
			var payOrder model.Order
			poe := global.GVA_DB.First(&payOrder, "out_trade_no = ?", result.OutTradeNo).Update("trade_state", "SUCCESS").Error
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

		log.Printf("status=%d resp=%s", result.TradeState, result)
		return err, order
	}
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
