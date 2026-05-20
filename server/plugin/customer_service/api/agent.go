package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AgentApi struct{}

// GetAgentList 获取坐席列表
// @Tags CustomerServiceAdmin
// @Summary 获取坐席列表（含在线状态）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.CsAgent} "获取成功"
// @Router /cs/admin/agent/list [get]
func (a *AgentApi) GetAgentList(c *gin.Context) {
	list, err := service.Service.AgentService.GetList()
	if err != nil {
		global.GVA_LOG.Error("获取坐席列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// CreateAgent 创建坐席
// @Tags CustomerServiceAdmin
// @Summary 添加坐席
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.CsAgent true "坐席信息"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /cs/admin/agent/create [post]
func (a *AgentApi) CreateAgent(c *gin.Context) {
	var agent model.CsAgent
	if err := c.ShouldBindJSON(&agent); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if agent.UserID == 0 {
		response.FailWithMessage("用户ID不能为空", c)
		return
	}
	// 先尝试获取已有记录
	existing, err := service.Service.AgentService.GetOrCreate(agent.UserID)
	if err != nil {
		response.FailWithMessage("创建失败", c)
		return
	}
	// 应用表单中的配置字段
	updates := map[string]interface{}{
		"max_sessions": agent.MaxSessions,
		"is_enabled":   agent.IsEnabled,
	}
	if agent.Nickname != "" {
		updates["nickname"] = agent.Nickname
	}
	if agent.MaxSessions < 1 {
		updates["max_sessions"] = 5 // 默认值
	}
	global.GVA_DB.Model(&model.CsAgent{}).Where("id = ?", existing.ID).Updates(updates)
	response.OkWithMessage("创建成功", c)
}

// UpdateAgent 更新坐席配置
// @Tags CustomerServiceAdmin
// @Summary 更新坐席配置（最大会话数、是否启用、昵称）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body csReq.UpdateAgentReq true "坐席配置"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /cs/admin/agent/update [put]
func (a *AgentApi) UpdateAgent(c *gin.Context) {
	var req csReq.UpdateAgentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	// 使用 Updates(map) 只更新指定列，避免 Save() 将零值字段（UserID、CreatedAt 等）覆盖写入
	updates := map[string]interface{}{
		"max_sessions": req.MaxSessions,
		"nickname":     req.Nickname,
	}
	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
	}
	if err := global.GVA_DB.Model(&model.CsAgent{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		global.GVA_LOG.Error("更新坐席失败", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteAgent 删除坐席
// @Tags CustomerServiceAdmin
// @Summary 删除坐席
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query int true "坐席ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /cs/admin/agent/delete [delete]
func (a *AgentApi) DeleteAgent(c *gin.Context) {
	var agent model.CsAgent
	if err := c.ShouldBindQuery(&agent); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := service.Service.AgentService.Delete(agent.ID); err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
