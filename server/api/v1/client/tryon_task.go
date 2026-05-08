package client

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TryonTaskApi struct{}

var tryonTaskService = service.ServiceGroupApp.ClientServiceGroup.TryonTaskService

// DeleteTryonTask 删除试衣任务（管理端）
// @Tags TryonTask
// @Summary 删除试衣任务（管理端）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query int true "任务ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /tryonTask/deleteTryonTask [delete]
func (api *TryonTaskApi) DeleteTryonTask(c *gin.Context) {
	idStr := c.Query("ID")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage(i18n.T(c, "invalidID"), c)
		return
	}

	if err := tryonTaskService.DeleteTryonTask(uint(id)); err != nil {
		global.GVA_LOG.Error("删除试衣任务失败!", zap.Error(err), zap.Uint64("taskID", id))
		response.FailWithMessage(i18n.T(c, "deleteFail"), c)
		return
	}

	response.OkWithMessage(i18n.T(c, "deleteSuccess"), c)
}

// DeleteTryonTaskByIds 批量删除试衣任务（管理端）
// @Tags TryonTask
// @Summary 批量删除试衣任务（管理端）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param IDs[] query []int true "任务ID列表"
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /tryonTask/deleteTryonTaskByIds [delete]
func (api *TryonTaskApi) DeleteTryonTaskByIds(c *gin.Context) {
	idStrs := c.QueryArray("IDs[]")
	if len(idStrs) == 0 {
		response.FailWithMessage(i18n.T(c, "invalidIDs"), c)
		return
	}

	ids := make([]uint, 0, len(idStrs))
	for _, idStr := range idStrs {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil || id == 0 {
			response.FailWithMessage(i18n.T(c, "invalidIDs"), c)
			return
		}
		ids = append(ids, uint(id))
	}

	if err := tryonTaskService.DeleteTryonTaskByIds(ids); err != nil {
		global.GVA_LOG.Error("批量删除试衣任务失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "batchDeleteFail"), c)
		return
	}

	response.OkWithMessage(i18n.T(c, "batchDeleteSuccess"), c)
}

// DeleteMyTryonTask 删除我的试衣任务（客户端）
// @Tags TryonTask
// @Summary 删除我的试衣任务（客户端）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param taskNo query string true "任务编号"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /tryonTask/deleteMyTryonTask [delete]
func (api *TryonTaskApi) DeleteMyTryonTask(c *gin.Context) {
	taskNo := c.Query("taskNo")
	if taskNo == "" {
		response.FailWithMessage(i18n.T(c, "invalidTaskNo"), c)
		return
	}

	userID := utils.GetUserID(c)
	if err := tryonTaskService.DeleteMyTryonTaskByTaskNo(userID, taskNo); err != nil {
		global.GVA_LOG.Error("删除我的试衣任务失败!", zap.Error(err), zap.Uint("userID", userID), zap.String("taskNo", taskNo))
		response.FailWithMessage(i18n.T(c, "deleteFail"), c)
		return
	}

	response.OkWithMessage(i18n.T(c, "deleteSuccess"), c)
}

// CreateTryonTask 创建试衣任务
// @Tags TryonTask
// @Summary 创建试衣任务（扣币后调用模型，失败全额退币）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body clientReq.CreateTryonTaskReq true "创建试衣任务"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "处理成功"
// @Router /tryonTask/createTryonTask [post]
func (api *TryonTaskApi) CreateTryonTask(c *gin.Context) {
	var req clientReq.CreateTryonTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	userID := utils.GetUserID(c)
	task, reused, err := tryonTaskService.CreateTryonTask(c.Request.Context(), userID, req)
	if err != nil {
		global.GVA_LOG.Error("创建试衣任务失败!", zap.Error(err), zap.Uint("userID", userID), zap.Uint("taskID", task.ID))
		response.FailWithDetailed(gin.H{"task": task, "reused": reused}, i18n.T(c, err.Error()), c)
		return
	}

	msg := i18n.T(c, "tryonTaskDone")
	if task.Status == "processing" {
		msg = i18n.T(c, "tryonTaskProcessing")
	}
	if reused {
		msg = i18n.T(c, "requestReusedHistory")
	}
	response.OkWithDetailed(gin.H{"task": task, "reused": reused}, msg, c)
}

// FindTryonTask 根据ID获取试衣任务
// @Tags TryonTask
// @Summary 根据ID获取试衣任务
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query int true "任务ID"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "查询成功"
// @Router /tryonTask/findTryonTask [get]
func (api *TryonTaskApi) FindTryonTask(c *gin.Context) {
	idStr := c.Query("ID")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage(i18n.T(c, "invalidID"), c)
		return
	}

	userID := utils.GetUserID(c)
	task, err := tryonTaskService.GetTryonTaskByID(c.Request.Context(), uint(id), userID)
	if err != nil {
		global.GVA_LOG.Error("查询试衣任务失败!", zap.Error(err), zap.Uint("userID", userID), zap.Uint64("taskID", id))
		response.FailWithMessage(i18n.T(c, "queryFail"), c)
		return
	}

	response.OkWithData(gin.H{"reTryonTask": task}, c)
}

// GetMyTryonTaskList 获取我的试衣任务列表
// @Tags TryonTask
// @Summary 获取我的试衣任务列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query clientReq.TryonTaskSearch true "分页查询我的试衣任务"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /tryonTask/getMyTryonTaskList [get]
func (api *TryonTaskApi) GetMyTryonTaskList(c *gin.Context) {
	var pageInfo clientReq.TryonTaskSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	userID := utils.GetUserID(c)
	list, total, err := tryonTaskService.GetMyTryonTaskList(userID, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取试衣任务列表失败!", zap.Error(err), zap.Uint("userID", userID))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, i18n.T(c, "getSuccess"), c)
}

// GetTryonTaskList 获取试衣任务列表（管理端）
// @Tags TryonTask
// @Summary 获取试衣任务列表（管理端）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query clientReq.TryonTaskSearch true "分页查询试衣任务"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /tryonTask/getTryonTaskList [get]
func (api *TryonTaskApi) GetTryonTaskList(c *gin.Context) {
	var pageInfo clientReq.TryonTaskSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	list, total, err := tryonTaskService.GetTryonTaskList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取试衣任务列表失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, i18n.T(c, "getSuccess"), c)
}

// GetTryonTaskStats 获取试衣任务统计（管理端）
// @Tags TryonTask
// @Summary 获取试衣任务统计（管理端）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query clientReq.TryonTaskSearch true "查询试衣任务统计"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取成功"
// @Router /tryonTask/getTryonTaskStats [get]
func (api *TryonTaskApi) GetTryonTaskStats(c *gin.Context) {
	var search clientReq.TryonTaskSearch
	if err := c.ShouldBindQuery(&search); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	stats, err := tryonTaskService.GetTryonTaskStats(search)
	if err != nil {
		global.GVA_LOG.Error("获取试衣任务统计失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}

	response.OkWithDetailed(stats, i18n.T(c, "getSuccess"), c)
}

// GetTryonTaskTrend 获取试衣任务趋势（管理端）
// @Tags TryonTask
// @Summary 获取试衣任务趋势（管理端）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query clientReq.TryonTaskSearch true "查询试衣任务趋势"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取成功"
// @Router /tryonTask/getTryonTaskTrend [get]
func (api *TryonTaskApi) GetTryonTaskTrend(c *gin.Context) {
	var search clientReq.TryonTaskSearch
	if err := c.ShouldBindQuery(&search); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	list, err := tryonTaskService.GetTryonTaskTrend(search)
	if err != nil {
		global.GVA_LOG.Error("获取试衣任务趋势失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}

	response.OkWithDetailed(gin.H{"list": list}, i18n.T(c, "getSuccess"), c)
}
