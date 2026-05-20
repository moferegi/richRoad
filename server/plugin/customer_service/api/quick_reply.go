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

type QuickReplyApi struct{}

// GetQuickReplyList 获取快捷回复列表
// @Tags CustomerService
// @Summary 获取快捷回复列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query csReq.QuickReplySearch true "搜索参数"
// @Success 200 {object} response.Response{data=response.PageResult} "获取成功"
// @Router /cs/agent/quickReply/list [get]
func (a *QuickReplyApi) GetQuickReplyList(c *gin.Context) {
	userID := utils.GetUserID(c)
	authID := utils.GetUserAuthorityId(c)
	if authID == 8080 {
		response.FailWithMessage("客户端用户不能访问坐席快捷回复列表", c)
		return
	}
	if authID != 888 {
		if _, err := service.Service.AgentService.GetEnabledByUserID(userID); err != nil {
			response.FailWithMessage("坐席不可用", c)
			return
		}
	}

	var search csReq.QuickReplySearch
	if err := c.ShouldBindQuery(&search); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := service.Service.QuickReplyService.GetList(search)
	if err != nil {
		global.GVA_LOG.Error("获取快捷回复失败", zap.Error(err))
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

// GetAllQuickReplies 获取全部快捷回复（供坐席工作台使用）
// @Tags CustomerService
// @Summary 获取全部快捷回复
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.CsQuickReply} "获取成功"
// @Router /cs/quickReply/all [get]
func (a *QuickReplyApi) GetAllQuickReplies(c *gin.Context) {
	list, err := service.Service.QuickReplyService.GetAll()
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// CreateQuickReply 创建快捷回复
// @Tags CustomerServiceAdmin
// @Summary 创建快捷回复
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.CsQuickReply true "快捷回复内容"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /cs/admin/quickReply/create [post]
func (a *QuickReplyApi) CreateQuickReply(c *gin.Context) {
	var qr model.CsQuickReply
	if err := c.ShouldBindJSON(&qr); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	qr.CreatedBy = utils.GetUserID(c)
	if err := service.Service.QuickReplyService.Create(&qr); err != nil {
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateQuickReply 更新快捷回复
// @Tags CustomerServiceAdmin
// @Summary 更新快捷回复
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.CsQuickReply true "快捷回复内容"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /cs/admin/quickReply/update [put]
func (a *QuickReplyApi) UpdateQuickReply(c *gin.Context) {
	var qr model.CsQuickReply
	if err := c.ShouldBindJSON(&qr); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := service.Service.QuickReplyService.Update(&qr); err != nil {
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteQuickReply 删除快捷回复
// @Tags CustomerServiceAdmin
// @Summary 删除快捷回复
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query int true "ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /cs/admin/quickReply/delete [delete]
func (a *QuickReplyApi) DeleteQuickReply(c *gin.Context) {
	var qr model.CsQuickReply
	if err := c.ShouldBindQuery(&qr); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := service.Service.QuickReplyService.Delete(qr.ID); err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
