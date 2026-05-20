package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PointRecordApi struct{}

func isPointRecordAdmin(authorityId uint) bool {
	return authorityId == 888 || authorityId == 8881
}

// CreatePointRecord 创建积分记录管理
// @Tags PointRecord
// @Summary 创建积分记录管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body client.PointRecord true "创建积分记录管理"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /cpr/createPointRecord [post]
func (cprApi *PointRecordApi) CreatePointRecord(c *gin.Context) {
	if !isPointRecordAdmin(utils.GetUserAuthorityId(c)) {
		failClientWithKey(c, "noPermission")
		return
	}

	// 创建业务用Context
	ctx := c.Request.Context()

	var cpr client.PointRecord
	err := c.ShouldBindJSON(&cpr)
	if err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	err = cprService.CreatePointRecord(ctx, &cpr)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		failClientWithKey(c, "createFail")
		return
	}
	response.OkWithMessage(i18n.T(c, "createSuccess"), c)
}

// DeletePointRecord 删除积分记录管理
// @Tags PointRecord
// @Summary 删除积分记录管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body client.PointRecord true "删除积分记录管理"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /cpr/deletePointRecord [delete]
func (cprApi *PointRecordApi) DeletePointRecord(c *gin.Context) {
	if !isPointRecordAdmin(utils.GetUserAuthorityId(c)) {
		failClientWithKey(c, "noPermission")
		return
	}

	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := cprService.DeletePointRecord(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		failClientWithKey(c, "deleteFail")
		return
	}
	response.OkWithMessage(i18n.T(c, "deleteSuccess"), c)
}

// DeletePointRecordByIds 批量删除积分记录管理
// @Tags PointRecord
// @Summary 批量删除积分记录管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /cpr/deletePointRecordByIds [delete]
func (cprApi *PointRecordApi) DeletePointRecordByIds(c *gin.Context) {
	if !isPointRecordAdmin(utils.GetUserAuthorityId(c)) {
		failClientWithKey(c, "noPermission")
		return
	}

	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := cprService.DeletePointRecordByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		failClientWithKey(c, "batchDeleteFail")
		return
	}
	response.OkWithMessage(i18n.T(c, "batchDeleteSuccess"), c)
}

// UpdatePointRecord 更新积分记录管理
// @Tags PointRecord
// @Summary 更新积分记录管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body client.PointRecord true "更新积分记录管理"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /cpr/updatePointRecord [put]
func (cprApi *PointRecordApi) UpdatePointRecord(c *gin.Context) {
	if !isPointRecordAdmin(utils.GetUserAuthorityId(c)) {
		failClientWithKey(c, "noPermission")
		return
	}

	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var cpr client.PointRecord
	err := c.ShouldBindJSON(&cpr)
	if err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	err = cprService.UpdatePointRecord(ctx, cpr)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		failClientWithKey(c, "updateFail")
		return
	}
	response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
}

// FindPointRecord 用id查询积分记录管理
// @Tags PointRecord
// @Summary 用id查询积分记录管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询积分记录管理"
// @Success 200 {object} response.Response{data=client.PointRecord,msg=string} "查询成功"
// @Router /cpr/findPointRecord [get]
func (cprApi *PointRecordApi) FindPointRecord(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	recpr, err := cprService.GetPointRecord(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		failClientWithKey(c, "queryFail")
		return
	}
	if !isPointRecordAdmin(utils.GetUserAuthorityId(c)) {
		uid := int(utils.GetUserID(c))
		if recpr.UserId == nil || *recpr.UserId != uid {
			failClientWithKey(c, "noPermission")
			return
		}
	}
	response.OkWithData(recpr, c)
}

// GetPointRecordList 分页获取积分记录管理列表
// @Tags PointRecord
// @Summary 分页获取积分记录管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query clientReq.PointRecordSearch true "分页获取积分记录管理列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /cpr/getPointRecordList [get]
func (cprApi *PointRecordApi) GetPointRecordList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo clientReq.PointRecordSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	// 非管理员只能查看自己的积分记录
	authID := utils.GetUserAuthorityId(c)
	if authID != 888 {
		uid := int(utils.GetUserID(c))
		pageInfo.UserId = &uid
	}
	list, total, err := cprService.GetPointRecordInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		failClientWithKey(c, "getFail")
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, i18n.T(c, "getSuccess"), c)
}

// GetPointRecordStats 获取试衣币统计（管理端）
// @Tags PointRecord
// @Summary 获取试衣币统计（管理端）
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query clientReq.TryonPointStatsSearch true "查询试衣币统计"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取成功"
// @Router /cpr/getPointRecordStats [get]
func (cprApi *PointRecordApi) GetPointRecordStats(c *gin.Context) {
	ctx := c.Request.Context()

	var search clientReq.TryonPointStatsSearch
	if err := c.ShouldBindQuery(&search); err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}

	authID := utils.GetUserAuthorityId(c)
	if authID != 888 {
		uid := int(utils.GetUserID(c))
		search.UserId = &uid
	}

	stats, err := cprService.GetTryonPointStats(ctx, search)
	if err != nil {
		global.GVA_LOG.Error("获取试衣币统计失败!", zap.Error(err))
		failClientWithKey(c, "getFail")
		return
	}

	response.OkWithDetailed(stats, i18n.T(c, "getSuccess"), c)
}

// GetPointRecordPublic 不需要鉴权的积分记录管理接口
// @Tags PointRecord
// @Summary 不需要鉴权的积分记录管理接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /cpr/getPointRecordPublic [get]
func (cprApi *PointRecordApi) GetPointRecordPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	cprService.GetPointRecordPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的积分记录管理接口信息",
	}, i18n.T(c, "getSuccess"), c)
}
