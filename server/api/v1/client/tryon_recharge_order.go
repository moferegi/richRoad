package client

import (
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TryonRechargeOrderApi 试衣币充值订单API
type TryonRechargeOrderApi struct{}

var tryonRechargeOrderService = service.ServiceGroupApp.ClientServiceGroup.TryonRechargeOrderService

// CreateTryonRechargeOrder 创建试衣币充值订单
// @Tags TryonRechargeOrder
// @Summary 创建试衣币充值订单（客户端）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body clientReq.CreateTryonRechargeOrderReq true "充值参数"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "创建成功"
// @Router /tryonRechargeOrder/createTryonRechargeOrder [post]
func (api *TryonRechargeOrderApi) CreateTryonRechargeOrder(c *gin.Context) {
	var req clientReq.CreateTryonRechargeOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	userID := utils.GetUserID(c)
	order, err := tryonRechargeOrderService.CreateTryonRechargeOrder(c.Request.Context(), userID, req)
	if err != nil {
		global.GVA_LOG.Error("创建充值订单失败", zap.Error(err), zap.Uint("userID", userID))
		response.FailWithMessage(i18n.T(c, err.Error()), c)
		return
	}

	response.OkWithDetailed(i18n.LocalizeResponseData(c, gin.H{"order": order}), i18n.T(c, "createSuccess"), c)
}

// UpdateTryonRechargeOrderPayMethod 更新充值订单支付方式
// @Tags TryonRechargeOrder
// @Summary 更新充值订单支付方式（仅待支付订单）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body clientReq.UpdateTryonRechargeOrderPayMethodReq true "订单ID与支付方式"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /tryonRechargeOrder/updateTryonRechargeOrderPayMethod [post]
func (api *TryonRechargeOrderApi) UpdateTryonRechargeOrderPayMethod(c *gin.Context) {
	var req clientReq.UpdateTryonRechargeOrderPayMethodReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	userID := utils.GetUserID(c)
	if err := tryonRechargeOrderService.UpdateTryonRechargeOrderPayMethod(userID, req); err != nil {
		global.GVA_LOG.Error("更新充值订单支付方式失败", zap.Error(err), zap.Uint("userID", userID), zap.Uint("orderID", req.ID))
		response.FailWithMessage(i18n.T(c, err.Error()), c)
		return
	}

	response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
}

// SubmitTryonRechargeOrderPayment 提交充值订单付款确认
// @Tags TryonRechargeOrder
// @Summary 提交充值订单付款确认（仅扫码支付且待支付）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body clientReq.SubmitTryonRechargeOrderPaymentReq true "订单ID"
// @Success 200 {object} response.Response{msg=string} "提交成功"
// @Router /tryonRechargeOrder/submitTryonRechargeOrderPayment [post]
func (api *TryonRechargeOrderApi) SubmitTryonRechargeOrderPayment(c *gin.Context) {
	var req clientReq.SubmitTryonRechargeOrderPaymentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	userID := utils.GetUserID(c)
	if err := tryonRechargeOrderService.SubmitTryonRechargeOrderPayment(userID, req.ID); err != nil {
		global.GVA_LOG.Error("提交充值订单付款确认失败", zap.Error(err), zap.Uint("userID", userID), zap.Uint("orderID", req.ID))
		response.FailWithMessage(i18n.T(c, err.Error()), c)
		return
	}

	response.OkWithMessage(i18n.T(c, "submitSuccess"), c)
}

// CancelTryonRechargeOrder 取消充值订单
// @Tags TryonRechargeOrder
// @Summary 取消充值订单（仅待支付）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query int true "订单ID"
// @Success 200 {object} response.Response{msg=string} "取消成功"
// @Router /tryonRechargeOrder/cancelTryonRechargeOrder [post]
func (api *TryonRechargeOrderApi) CancelTryonRechargeOrder(c *gin.Context) {
	id, err := strconv.ParseUint(strings.TrimSpace(c.Query("ID")), 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage(i18n.T(c, "invalidOrder"), c)
		return
	}

	userID := utils.GetUserID(c)
	if err = tryonRechargeOrderService.CancelTryonRechargeOrder(userID, uint(id)); err != nil {
		global.GVA_LOG.Error("取消充值订单失败", zap.Error(err), zap.Uint("userID", userID), zap.Uint64("orderID", id))
		response.FailWithMessage(i18n.T(c, err.Error()), c)
		return
	}

	response.OkWithMessage(i18n.T(c, "cancelSuccess"), c)
}

// ConfirmTryonRechargeOrderPayment 确认充值订单支付（管理端）
// @Tags TryonRechargeOrder
// @Summary 确认充值订单支付（管理端）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query int true "订单ID"
// @Param remark query string false "备注"
// @Success 200 {object} response.Response{msg=string} "确认成功"
// @Router /tryonRechargeOrder/confirmTryonRechargeOrderPayment [post]
func (api *TryonRechargeOrderApi) ConfirmTryonRechargeOrderPayment(c *gin.Context) {
	if authorityID := utils.GetUserAuthorityId(c); authorityID != 888 && authorityID != 8881 {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return
	}

	id, err := strconv.ParseUint(strings.TrimSpace(c.Query("ID")), 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage(i18n.T(c, "invalidOrder"), c)
		return
	}

	remark := strings.TrimSpace(c.Query("remark"))
	if err = tryonRechargeOrderService.ConfirmTryonRechargeOrderPayment(c.Request.Context(), uint(id), remark); err != nil {
		global.GVA_LOG.Error("确认充值订单支付失败", zap.Error(err), zap.Uint64("orderID", id))
		response.FailWithMessage(i18n.T(c, err.Error()), c)
		return
	}

	response.OkWithMessage(i18n.T(c, "confirmSuccess"), c)
}

// SelfTryonRechargeOrder 获取我的充值订单详情
// @Tags TryonRechargeOrder
// @Summary 获取我的充值订单详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query int true "订单ID"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "查询成功"
// @Router /tryonRechargeOrder/selfTryonRechargeOrder [get]
func (api *TryonRechargeOrderApi) SelfTryonRechargeOrder(c *gin.Context) {
	id, err := strconv.ParseUint(strings.TrimSpace(c.Query("ID")), 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage(i18n.T(c, "invalidOrder"), c)
		return
	}

	userID := utils.GetUserID(c)
	order, err := tryonRechargeOrderService.GetMyTryonRechargeOrderByID(userID, uint(id))
	if err != nil {
		global.GVA_LOG.Error("查询我的充值订单失败", zap.Error(err), zap.Uint("userID", userID), zap.Uint64("orderID", id))
		response.FailWithMessage(i18n.T(c, err.Error()), c)
		return
	}

	response.OkWithDetailed(i18n.LocalizeResponseData(c, gin.H{"order": order}), i18n.T(c, "querySuccess"), c)
}

// FindTryonRechargeOrder 获取充值订单详情（管理端）
// @Tags TryonRechargeOrder
// @Summary 获取充值订单详情（管理端）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query int true "订单ID"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "查询成功"
// @Router /tryonRechargeOrder/findTryonRechargeOrder [get]
func (api *TryonRechargeOrderApi) FindTryonRechargeOrder(c *gin.Context) {
	if authorityID := utils.GetUserAuthorityId(c); authorityID != 888 && authorityID != 8881 {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return
	}

	id, err := strconv.ParseUint(strings.TrimSpace(c.Query("ID")), 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage(i18n.T(c, "invalidOrder"), c)
		return
	}

	order, err := tryonRechargeOrderService.GetTryonRechargeOrderByID(uint(id))
	if err != nil {
		global.GVA_LOG.Error("查询充值订单失败", zap.Error(err), zap.Uint64("orderID", id))
		response.FailWithMessage(i18n.T(c, err.Error()), c)
		return
	}

	response.OkWithDetailed(i18n.LocalizeResponseData(c, gin.H{"order": order}), i18n.T(c, "querySuccess"), c)
}

// GetMyTryonRechargeOrderList 获取我的充值订单列表
// @Tags TryonRechargeOrder
// @Summary 获取我的充值订单列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query clientReq.TryonRechargeOrderSearch true "分页查询条件"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "查询成功"
// @Router /tryonRechargeOrder/getMyTryonRechargeOrderList [get]
func (api *TryonRechargeOrderApi) GetMyTryonRechargeOrderList(c *gin.Context) {
	var search clientReq.TryonRechargeOrderSearch
	if err := c.ShouldBindQuery(&search); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	userID := utils.GetUserID(c)
	list, total, err := tryonRechargeOrderService.GetMyTryonRechargeOrderList(userID, search)
	if err != nil {
		global.GVA_LOG.Error("获取我的充值订单列表失败", zap.Error(err), zap.Uint("userID", userID))
		response.FailWithMessage(i18n.T(c, err.Error()), c)
		return
	}

	response.OkWithDetailed(i18n.LocalizeResponseData(c, response.PageResult{List: list, Total: total, Page: search.Page, PageSize: search.PageSize}), i18n.T(c, "querySuccess"), c)
}

// GetTryonRechargeOrderList 获取充值订单列表（管理端）
// @Tags TryonRechargeOrder
// @Summary 获取充值订单列表（管理端）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query clientReq.TryonRechargeOrderSearch true "分页查询条件"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "查询成功"
// @Router /tryonRechargeOrder/getTryonRechargeOrderList [get]
func (api *TryonRechargeOrderApi) GetTryonRechargeOrderList(c *gin.Context) {
	if authorityID := utils.GetUserAuthorityId(c); authorityID != 888 && authorityID != 8881 {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return
	}

	var search clientReq.TryonRechargeOrderSearch
	if err := c.ShouldBindQuery(&search); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	list, total, err := tryonRechargeOrderService.GetTryonRechargeOrderList(search)
	if err != nil {
		global.GVA_LOG.Error("获取充值订单列表失败", zap.Error(err))
		response.FailWithMessage(i18n.T(c, err.Error()), c)
		return
	}

	response.OkWithDetailed(i18n.LocalizeResponseData(c, response.PageResult{List: list, Total: total, Page: search.Page, PageSize: search.PageSize}), i18n.T(c, "querySuccess"), c)
}
