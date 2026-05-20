package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WxpayApi struct{}

// @Tags Wxpay
// @Summary 获取微信支付二维码和ID
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /wxpay/getPayCode[post]
func (p *WxpayApi) GetPayCode(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	order.CustomerID = utils.GetUserID(c)
	if order.CustomerID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}
	if err, codeUrl, codeId := service.ServiceGroupApp.GetPayCode(order); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage("获取付款码失败", c)
	} else {
		response.OkWithData(gin.H{
			"codeUrl": codeUrl,
			"codeId":  codeId,
		}, c)
	}
}

func (p *WxpayApi) CheckNeedPay(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	order.CustomerID = utils.GetUserID(c)
	if order.CustomerID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}

	err, needPay := service.ServiceGroupApp.CheckNeedPay(order)

	if err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage("检查订单失败", c)
		return
	}
	response.OkWithData(needPay, c)
}

// @Tags Wxpay
// @Summary 获取微信支付二维码和ID
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /wxpay/GetPayParams[post]
func (p *WxpayApi) GetPayParams(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	order.CustomerID = utils.GetUserID(c)
	if order.CustomerID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}
	if err, patConf := service.ServiceGroupApp.GetPayConf(order); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage("获取付款码失败", c)
	} else {
		response.OkWithData(patConf, c)
	}
}

// @Tags Wxpay
// @Summary 获取支付结果
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /wxpay/getOrderById[get]
func (p *WxpayApi) GetOrderById(c *gin.Context) {
	id := c.Query("orderID")
	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}
	if err, data := service.ServiceGroupApp.GetOrderById(id, userID); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage("获取付款码失败", c)
	} else {
		response.OkWithData(data, c)
	}
}

// @Tags Wxpay
// @Summary 回调支付结果
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
