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

type TryonClothApi struct{}

var tryonClothService = service.ServiceGroupApp.ClientServiceGroup.TryonClothService

// CreateTryonCloth 创建我的衣橱
// @Tags TryonCloth
// @Summary 创建我的衣橱
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body clientReq.CreateTryonClothReq true "创建我的衣橱"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "创建成功"
// @Router /tryonCloth/createTryonCloth [post]
func (api *TryonClothApi) CreateTryonCloth(c *gin.Context) {
	var req clientReq.CreateTryonClothReq
	if err := c.ShouldBindJSON(&req); err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}

	userID := utils.GetUserID(c)
	cloth, err := tryonClothService.CreateTryonCloth(userID, req)
	if err != nil {
		global.GVA_LOG.Error("创建我的衣橱失败!", zap.Error(err), zap.Uint("userID", userID))
		failClientWithErr(c, err)
		return
	}

	response.OkWithDetailed(i18n.LocalizeResponseData(c, gin.H{"cloth": cloth}), i18n.T(c, "createSuccess"), c)
}

// UpdateTryonCloth 更新我的衣橱
// @Tags TryonCloth
// @Summary 更新我的衣橱
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body clientReq.UpdateTryonClothReq true "更新我的衣橱"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /tryonCloth/updateTryonCloth [put]
func (api *TryonClothApi) UpdateTryonCloth(c *gin.Context) {
	var req clientReq.UpdateTryonClothReq
	if err := c.ShouldBindJSON(&req); err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}

	userID := utils.GetUserID(c)
	authorityID := utils.GetUserAuthorityId(c)
	err := tryonClothService.UpdateTryonCloth(userID, authorityID, req)
	if err != nil {
		global.GVA_LOG.Error("更新我的衣橱失败!", zap.Error(err), zap.Uint("userID", userID), zap.Uint("ID", req.ID))
		failClientWithErr(c, err)
		return
	}

	response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
}

// DeleteTryonCloth 删除我的衣橱
// @Tags TryonCloth
// @Summary 删除我的衣橱
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query int true "衣橱ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /tryonCloth/deleteTryonCloth [delete]
func (api *TryonClothApi) DeleteTryonCloth(c *gin.Context) {
	idStr := c.Query("ID")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		failClientWithKey(c, "invalidID")
		return
	}

	userID := utils.GetUserID(c)
	authorityID := utils.GetUserAuthorityId(c)
	err = tryonClothService.DeleteTryonCloth(userID, authorityID, uint(id))
	if err != nil {
		global.GVA_LOG.Error("删除我的衣橱失败!", zap.Error(err), zap.Uint("userID", userID), zap.Uint64("ID", id))
		failClientWithErr(c, err)
		return
	}

	response.OkWithMessage(i18n.T(c, "deleteSuccess"), c)
}

// DeleteTryonClothByIds 批量删除我的衣橱
// @Tags TryonCloth
// @Summary 批量删除我的衣橱
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param IDs[] query []int true "衣橱ID列表"
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /tryonCloth/deleteTryonClothByIds [delete]
func (api *TryonClothApi) DeleteTryonClothByIds(c *gin.Context) {
	idStrs := c.QueryArray("IDs[]")
	if len(idStrs) == 0 {
		failClientWithKey(c, "invalidIDs")
		return
	}

	ids := make([]uint, 0, len(idStrs))
	for _, idStr := range idStrs {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil || id == 0 {
			failClientWithKey(c, "invalidIDs")
			return
		}
		ids = append(ids, uint(id))
	}

	userID := utils.GetUserID(c)
	authorityID := utils.GetUserAuthorityId(c)
	err := tryonClothService.DeleteTryonClothByIds(userID, authorityID, ids)
	if err != nil {
		global.GVA_LOG.Error("批量删除我的衣橱失败!", zap.Error(err), zap.Uint("userID", userID))
		failClientWithErr(c, err)
		return
	}

	response.OkWithMessage(i18n.T(c, "batchDeleteSuccess"), c)
}

// GetMyTryonClothList 获取我的衣橱列表
// @Tags TryonCloth
// @Summary 获取我的衣橱列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param category query string false "分类(upper/lower/onepiece/shoes)"
// @Param categories[] query []string false "分类列表(upper/lower/onepiece/shoes)"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取成功"
// @Router /tryonCloth/getMyTryonClothList [get]
func (api *TryonClothApi) GetMyTryonClothList(c *gin.Context) {
	categories := c.QueryArray("categories[]")
	if category := c.Query("category"); category != "" {
		categories = append(categories, category)
	}

	userID := utils.GetUserID(c)
	list, err := tryonClothService.GetMyTryonClothList(userID, categories)
	if err != nil {
		global.GVA_LOG.Error("获取我的衣橱列表失败!", zap.Error(err), zap.Uint("userID", userID))
		failClientWithErr(c, err)
		return
	}

	response.OkWithDetailed(i18n.LocalizeResponseData(c, gin.H{"list": list}), i18n.T(c, "getSuccess"), c)
}

// GetTryonClothList 获取我的衣橱列表（管理端）
// @Tags TryonCloth
// @Summary 获取我的衣橱列表（管理端）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query clientReq.TryonClothSearch true "分页查询我的衣橱列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /tryonCloth/getTryonClothList [get]
func (api *TryonClothApi) GetTryonClothList(c *gin.Context) {
	var pageInfo clientReq.TryonClothSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}

	if utils.GetUserAuthorityId(c) != 888 {
		pageInfo.UserID = utils.GetUserID(c)
	}

	list, total, err := tryonClothService.GetTryonClothList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取我的衣橱列表失败!", zap.Error(err))
		failClientWithErr(c, err)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, i18n.T(c, "getSuccess"), c)
}
