package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type QrcodePaymentApi struct{}

var qrcodePaymentService = service.ServiceGroupApp.ShopServiceGroup.QrcodePaymentService

// CreateQrcodePayment 创建收款码
// @Tags QrcodePayment
// @Summary 创建收款码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.QrcodePayment true "收款码信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /qrcodePayment/createQrcodePayment [post]
func (api *QrcodePaymentApi) CreateQrcodePayment(c *gin.Context) {
	var info shop.QrcodePayment
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := qrcodePaymentService.CreateQrcodePayment(&info); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteQrcodePayment 删除收款码
// @Tags QrcodePayment
// @Summary 删除收款码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.QrcodePayment true "删除收款码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /qrcodePayment/deleteQrcodePayment [delete]
func (api *QrcodePaymentApi) DeleteQrcodePayment(c *gin.Context) {
	ID := c.Query("ID")
	if err := qrcodePaymentService.DeleteQrcodePayment(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// UpdateQrcodePayment 更新收款码
// @Tags QrcodePayment
// @Summary 更新收款码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.QrcodePayment true "更新收款码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /qrcodePayment/updateQrcodePayment [put]
func (api *QrcodePaymentApi) UpdateQrcodePayment(c *gin.Context) {
	var info shop.QrcodePayment
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := qrcodePaymentService.UpdateQrcodePayment(info); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// GetQrcodePaymentList 获取收款码列表
// @Tags QrcodePayment
// @Summary 获取收款码列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shopReq.QrcodePaymentSearch true "分页获取收款码列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /qrcodePayment/getQrcodePaymentList [get]
func (api *QrcodePaymentApi) GetQrcodePaymentList(c *gin.Context) {
	var pageInfo shopReq.QrcodePaymentSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := qrcodePaymentService.GetQrcodePaymentList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetEnabledQrcodePayments 获取启用的收款码列表（客户端）
// @Tags QrcodePayment
// @Summary 获取启用的收款码列表
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]shop.QrcodePayment,msg=string} "获取成功"
// @Router /qrcodePayment/getEnabledQrcodePayments [get]
func (api *QrcodePaymentApi) GetEnabledQrcodePayments(c *gin.Context) {
	if list, err := qrcodePaymentService.GetEnabledQrcodePayments(); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(i18n.LocalizeResponseData(c, list), "获取成功", c)
	}
}
