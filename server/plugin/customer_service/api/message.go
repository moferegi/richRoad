package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MessageApi struct{}

// GetMessageHistory 获取消息历史
// @Tags CustomerService
// @Summary 获取会话消息历史（分页）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query csReq.MessageSearch true "分页参数"
// @Success 200 {object} response.Response{data=response.PageResult} "获取成功"
// @Router /cs/message/history [get]
func (a *MessageApi) GetMessageHistory(c *gin.Context) {
	var search csReq.MessageSearch
	if err := c.ShouldBindQuery(&search); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := service.Service.MessageService.GetHistory(search)
	if err != nil {
		global.GVA_LOG.Error("获取消息历史失败", zap.Error(err))
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

// SendMessage 通过 REST 发送消息（坐席后台补发等场景）
// @Tags CustomerService
// @Summary 发送消息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body csReq.SendMessageReq true "消息内容"
// @Success 200 {object} response.Response{msg=string} "发送成功"
// @Router /cs/message/send [post]
func (a *MessageApi) SendMessage(c *gin.Context) {
	var req csReq.SendMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	senderID := utils.GetUserID(c)
	// 默认 agent 类型（REST 接口一般坐席调用）
	msg, err := service.Service.MessageService.Send(
		"agent", senderID, req.ConversationID,
		req.MsgType, req.Content, req.ClientMsgID,
	)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(msg, "发送成功", c)
}

// RevokeMessage 撤回消息
// @Tags CustomerService
// @Summary 撤回消息（2分钟内）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body csReq.RevokeMessageReq true "消息ID"
// @Success 200 {object} response.Response{msg=string} "撤回成功"
// @Router /cs/message/revoke [post]
func (a *MessageApi) RevokeMessage(c *gin.Context) {
	var req csReq.RevokeMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	senderID := utils.GetUserID(c)
	// 从 authority 区分发送方类型：8080 = 客户端用户，其余 = 坐席
	authID := utils.GetUserAuthorityId(c)
	senderType := "agent"
	if authID == 8080 {
		senderType = "user"
	}
	if err := service.Service.MessageService.Revoke(req.MessageID, senderID, senderType); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("撤回成功", c)
}
