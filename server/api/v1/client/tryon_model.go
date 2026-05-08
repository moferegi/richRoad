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

type TryonModelApi struct{}

var tryonModelService = service.ServiceGroupApp.ClientServiceGroup.TryonModelService

// CreateTryonModel 创建我的模特
// @Tags TryonModel
// @Summary 创建我的模特
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body clientReq.CreateTryonModelReq true "创建我的模特"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "创建成功"
// @Router /tryonModel/createTryonModel [post]
func (api *TryonModelApi) CreateTryonModel(c *gin.Context) {
	var req clientReq.CreateTryonModelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	userID := utils.GetUserID(c)
	model, err := tryonModelService.CreateTryonModel(userID, req)
	if err != nil {
		global.GVA_LOG.Error("创建我的模特失败!", zap.Error(err), zap.Uint("userID", userID))
		response.FailWithMessage(i18n.T(c, err.Error()), c)
		return
	}

	response.OkWithDetailed(gin.H{"model": model}, i18n.T(c, "createSuccess"), c)
}

// UpdateTryonModel 重命名我的模特
// @Tags TryonModel
// @Summary 重命名我的模特
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body clientReq.UpdateTryonModelReq true "重命名我的模特"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /tryonModel/updateTryonModel [put]
func (api *TryonModelApi) UpdateTryonModel(c *gin.Context) {
	var req clientReq.UpdateTryonModelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	userID := utils.GetUserID(c)
	authorityID := utils.GetUserAuthorityId(c)
	err := tryonModelService.UpdateTryonModel(userID, authorityID, req)
	if err != nil {
		global.GVA_LOG.Error("更新我的模特失败!", zap.Error(err), zap.Uint("userID", userID), zap.Uint("ID", req.ID))
		response.FailWithMessage(i18n.T(c, err.Error()), c)
		return
	}

	response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
}

// DeleteTryonModel 删除我的模特
// @Tags TryonModel
// @Summary 删除我的模特
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query int true "模特ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /tryonModel/deleteTryonModel [delete]
func (api *TryonModelApi) DeleteTryonModel(c *gin.Context) {
	idStr := c.Query("ID")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage(i18n.T(c, "invalidID"), c)
		return
	}

	userID := utils.GetUserID(c)
	authorityID := utils.GetUserAuthorityId(c)
	err = tryonModelService.DeleteTryonModel(userID, authorityID, uint(id))
	if err != nil {
		global.GVA_LOG.Error("删除我的模特失败!", zap.Error(err), zap.Uint("userID", userID), zap.Uint64("ID", id))
		response.FailWithMessage(i18n.T(c, err.Error()), c)
		return
	}

	response.OkWithMessage(i18n.T(c, "deleteSuccess"), c)
}

// DeleteTryonModelByIds 批量删除我的模特
// @Tags TryonModel
// @Summary 批量删除我的模特
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param IDs[] query []int true "模特ID列表"
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /tryonModel/deleteTryonModelByIds [delete]
func (api *TryonModelApi) DeleteTryonModelByIds(c *gin.Context) {
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

	userID := utils.GetUserID(c)
	authorityID := utils.GetUserAuthorityId(c)
	err := tryonModelService.DeleteTryonModelByIds(userID, authorityID, ids)
	if err != nil {
		global.GVA_LOG.Error("批量删除我的模特失败!", zap.Error(err), zap.Uint("userID", userID))
		response.FailWithMessage(i18n.T(c, err.Error()), c)
		return
	}

	response.OkWithMessage(i18n.T(c, "batchDeleteSuccess"), c)
}

// GetMyTryonModelList 获取我的模特列表
// @Tags TryonModel
// @Summary 获取我的模特列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取成功"
// @Router /tryonModel/getMyTryonModelList [get]
func (api *TryonModelApi) GetMyTryonModelList(c *gin.Context) {
	userID := utils.GetUserID(c)
	list, err := tryonModelService.GetMyTryonModelList(userID)
	if err != nil {
		global.GVA_LOG.Error("获取我的模特列表失败!", zap.Error(err), zap.Uint("userID", userID))
		response.FailWithMessage(i18n.T(c, err.Error()), c)
		return
	}

	response.OkWithDetailed(gin.H{"list": list}, i18n.T(c, "getSuccess"), c)
}

// GetTryonModelList 获取模特列表（管理端）
// @Tags TryonModel
// @Summary 获取模特列表（管理端）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query clientReq.TryonModelSearch true "分页查询模特列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /tryonModel/getTryonModelList [get]
func (api *TryonModelApi) GetTryonModelList(c *gin.Context) {
	var pageInfo clientReq.TryonModelSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	if utils.GetUserAuthorityId(c) != 888 {
		pageInfo.UserID = utils.GetUserID(c)
	}

	list, total, err := tryonModelService.GetTryonModelList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取模特列表失败!", zap.Error(err))
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
