package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ConversationApi struct{}

// GetOrCreateConversation 用户发起/获取当前会话
// @Tags CustomerService
// @Summary 客户端用户发起或获取会话
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=model.CsConversation} "获取成功"
// @Router /cs/conversation/getOrCreate [post]
func (a *ConversationApi) GetOrCreateConversation(c *gin.Context) {
	clientUserID := utils.GetUserID(c)
	conv, _, err := service.Service.ConversationService.CreateOrGetActive(clientUserID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(conv, "获取成功", c)
}

// GetConversationList 管理员获取会话列表
// @Tags CustomerServiceAdmin
// @Summary 获取会话列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query csReq.ConversationSearch true "搜索参数"
// @Success 200 {object} response.Response{data=response.PageResult} "获取成功"
// @Router /cs/admin/conversation/list [get]
func (a *ConversationApi) GetConversationList(c *gin.Context) {
	var search csReq.ConversationSearch
	if err := c.ShouldBindQuery(&search); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := service.Service.ConversationService.GetList(search)
	if err != nil {
		global.GVA_LOG.Error("获取会话列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     search.Page,
		PageSize: search.PageSize,
	}, "获取成功", c)
}

// GetAgentConversationList 坐席获取自己的会话列表
// @Tags CustomerService
// @Summary 坐席获取自己的会话列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query csReq.ConversationSearch true "搜索参数"
// @Success 200 {object} response.Response{data=response.PageResult} "获取成功"
// @Router /cs/agent/conversation/list [get]
func (a *ConversationApi) GetAgentConversationList(c *gin.Context) {
	agentSysUserID := utils.GetUserID(c)
	var search csReq.ConversationSearch
	if err := c.ShouldBindQuery(&search); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	search.AgentUserID = &agentSysUserID
	list, total, err := service.Service.ConversationService.GetList(search)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     search.Page,
		PageSize: search.PageSize,
	}, "获取成功", c)
}

// CloseConversation 关闭会话
// @Tags CustomerService
// @Summary 关闭会话
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body csReq.CloseConversationReq true "关闭参数"
// @Success 200 {object} response.Response{msg=string} "关闭成功"
// @Router /cs/agent/conversation/close [post]
func (a *ConversationApi) CloseConversation(c *gin.Context) {
	var req csReq.CloseConversationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	closedBy := req.ClosedBy
	if closedBy == "" {
		closedBy = model.ConvClosedByAgent
	}
	if err := service.Service.ConversationService.Close(req.ConversationID, closedBy); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("关闭成功", c)
}

// TransferConversation 转接会话
// @Tags CustomerServiceAdmin
// @Summary 转接会话给其他坐席
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body csReq.TransferConversationReq true "转接参数"
// @Success 200 {object} response.Response{msg=string} "转接成功"
// @Router /cs/agent/conversation/transfer [post]
func (a *ConversationApi) TransferConversation(c *gin.Context) {
	var req csReq.TransferConversationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.Service.ConversationService.Transfer(req.ConversationID, req.TargetAgentUserID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("转接成功", c)
}

// RateConversation 用户评价会话
// @Tags CustomerService
// @Summary 客户端用户评价会话
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body csReq.RateConversationReq true "评价参数"
// @Success 200 {object} response.Response{msg=string} "评价成功"
// @Router /cs/conversation/rate [post]
func (a *ConversationApi) RateConversation(c *gin.Context) {
	var req csReq.RateConversationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	clientUserID := utils.GetUserID(c)
	if err := service.Service.ConversationService.Rate(req.ConversationID, clientUserID, req.Rating); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("评价成功", c)
}

// AssignConversation 管理员手动分配会话
// @Tags CustomerServiceAdmin
// @Summary 手动分配会话给坐席
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body csReq.TransferConversationReq true "分配参数"
// @Success 200 {object} response.Response{msg=string} "分配成功"
// @Router /cs/admin/conversation/assign [post]
func (a *ConversationApi) AssignConversation(c *gin.Context) {
	var req csReq.TransferConversationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.Service.ConversationService.AssignToAgent(req.ConversationID, req.TargetAgentUserID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("分配成功", c)
}
