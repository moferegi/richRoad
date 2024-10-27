package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/service"
	wxUtils "github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/utils"
	shopService "github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"time"
)

var orderService = shopService.ServiceGroupApp.ShopServiceGroup.OrderService

type WxpayApi struct{}

func (p *WxpayApi) GetPayParams(c *gin.Context) {
	var payOrder model.Order
	c.ShouldBindJSON(&payOrder)
	payOrder.CustomerID = utils.GetUserID(c)

	if err, payParams := service.ServiceGroupApp.GetPayParams(payOrder); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage("获取支付参数失败:"+err.Error(), c)
	} else {
		var resMap map[string]string
		resMap = make(map[string]string)
		//商户申请的小程序对应的appid，由微信支付生成
		resMap["appId"] = global.GVA_CONFIG.Wxpay.AppID
		// 时间戳，标准北京时间，时区为东八区，自1970年1月1日 0点0分0秒以来的秒数。注意：部分系统取到的值为毫秒级，需要转换成秒(10位数字)。
		resMap["timeStamp"] = strconv.Itoa(int(time.Now().Unix()))
		//随机字符串，不长于32位
		resMap["nonceStr"] = utils.RandomString(32)
		//小程序下单接口返回的prepay_id参数值，提交格式如：prepay_id=***
		resMap["package"] = "prepay_id=" + payParams
		// 签名类型，默认为RSA，仅支持RSA。
		resMap["signType"] = "RSA"
		//签名，使用字段appId、timeStamp、nonceStr、package计算得出的签名值
		resMap["paySign"], err = wxUtils.GeneratePaySign(global.GVA_CONFIG.Wxpay.PemPath, resMap["appId"], resMap["timeStamp"], resMap["nonceStr"], resMap["package"])
		if err != nil {
			global.GVA_LOG.Error("失败!", zap.Error(err))
			response.FailWithMessage("获取支付参数失败:"+err.Error(), c)
			return
		}
		response.OkWithData(resMap, c)
	}
}

// @Tags Wxpay
// @Summary 获取支付结果
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /wxpay/getOrderById[get]
func (p *WxpayApi) GetOrderById(c *gin.Context) {
	id := c.Query("orderID")
	if err, data := service.ServiceGroupApp.GetOrderById(id); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage("获取付款码失败:"+err.Error(), c)
	} else {
		response.OkWithData(data, c)
	}
}

// @Tags Wxpay
// @Summary 获取支付结果
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /wxpay/payAction[post]
func (p *WxpayApi) PayAction(c *gin.Context) {
	var pay model.PayAction
	c.ShouldBindJSON(&pay)
	if err := service.ServiceGroupApp.PayAction(pay); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		c.JSON(200, gin.H{
			"code":    "FAIL",
			"message": "失败",
		})
	} else {
		c.JSON(200, gin.H{
			"code":    "SUCCESS",
			"message": "接收成功",
		})
	}
}
