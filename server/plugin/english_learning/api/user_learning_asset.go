package api

import (
	"math"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserLearningAssetApi struct{}

// 引入全局服务组
var userLearningAssetService = service.ServiceGroupApp.UserLearningAssetService
var learningAuthzService = service.ServiceGroupApp.LearningAuthzService

const heartbeatUsingSecsUpperBound = 120

func normalizeAssetPage(page int, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return page, pageSize
}

// Heartbeat 视频心跳接收器
// @Tags     EnglishLearning
// @Summary  接收前端心跳以扣减免费时长并保存看课进度
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.HeartbeatReq true "心跳参数"
// @Success  200  {object} response.Response{msg=string} "接收成功"
// @Router   /englishLearning/asset/heartbeat [post]
func (a *UserLearningAssetApi) Heartbeat(c *gin.Context) {
	var req request.HeartbeatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if req.EpisodeID == 0 {
		response.FailWithMessage("参数错误", c)
		return
	}
	if math.IsNaN(req.ProgressSecs) || math.IsInf(req.ProgressSecs, 0) || req.ProgressSecs < 0 {
		req.ProgressSecs = 0
	}

	// 保证使用该接口的必须是登录后的系统合法用户
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	// 为防止恶意请求(比如每秒发个心跳带30分钟的数据)，设置使用时间的合法阈值上限
	if req.UsingTimeSecs > heartbeatUsingSecsUpperBound || req.UsingTimeSecs < 0 {
		global.GVA_LOG.Warn("心跳包时间异常，疑似防薅羊毛拦截", zap.Uint("user", claims.BaseClaims.ID), zap.Int("secs", req.UsingTimeSecs))
		if req.UsingTimeSecs < 0 {
			req.UsingTimeSecs = 0
		} else {
			req.UsingTimeSecs = heartbeatUsingSecsUpperBound
		}
	}

	err = userLearningAssetService.HandleHeartbeat(claims.BaseClaims.ID, req.EpisodeID, req.ProgressSecs, req.UsingTimeSecs)
	if err != nil {
		global.GVA_LOG.Error("心跳处理失败!", zap.Error(err))
		response.FailWithMessage("进度保存失败", c)
		return
	}

	response.OkWithMessage("接收心跳成功", c)
}

// GetFreeTimeRecordList
// @Tags     EnglishLearning
// @Summary  获取英语学习免费时长流水列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.FreeTimeRecordSearch true "分页参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/asset/getFreeTimeRecordList [get]
func (a *UserLearningAssetApi) GetFreeTimeRecordList(c *gin.Context) {
	var query request.FreeTimeRecordSearch
	if err := c.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	list, total, page, pageSize, err := userLearningAssetService.GetFreeTimeRecordList(claims.BaseClaims.ID, query)
	if err != nil {
		global.GVA_LOG.Error("获取英语学习免费时长流水失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	page, pageSize = normalizeAssetPage(page, pageSize)
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, "获取成功", c)
}

// GrantEntitlement 授予用户学习资源权限
// @Tags     EnglishLearning
// @Summary  授予用户分类/剧集/单集粒度的学习权限
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.GrantEntitlementReq true "授权参数"
// @Success  200  {object} response.Response{msg=string} "授权成功"
// @Router   /englishLearning/asset/grantEntitlement [post]
func (a *UserLearningAssetApi) GrantEntitlement(c *gin.Context) {
	var req request.GrantEntitlementReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	if err = learningAuthzService.GrantEntitlement(req, claims.BaseClaims.ID); err != nil {
		global.GVA_LOG.Error("授予学习资源权限失败", zap.Error(err), zap.Uint("operator", claims.BaseClaims.ID), zap.Uint("userId", req.UserID), zap.String("resourceType", req.ResourceType), zap.Uint("resourceId", req.ResourceID))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("授权成功", c)
}

// RevokeEntitlement 撤销用户学习资源权限
// @Tags     EnglishLearning
// @Summary  撤销用户分类/剧集/单集粒度的学习权限
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.RevokeEntitlementReq true "撤销参数"
// @Success  200  {object} response.Response{msg=string} "撤销成功"
// @Router   /englishLearning/asset/revokeEntitlement [delete]
func (a *UserLearningAssetApi) RevokeEntitlement(c *gin.Context) {
	var req request.RevokeEntitlementReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := learningAuthzService.RevokeEntitlement(req); err != nil {
		global.GVA_LOG.Error("撤销学习资源权限失败", zap.Error(err), zap.Uint("userId", req.UserID), zap.String("resourceType", req.ResourceType), zap.Uint("resourceId", req.ResourceID))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("撤销成功", c)
}

// GetEntitlementList 获取用户学习资源权限列表
// @Tags     EnglishLearning
// @Summary  分页获取用户分类/剧集/单集授权记录
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.EntitlementSearch true "分页参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/asset/getEntitlementList [get]
func (a *UserLearningAssetApi) GetEntitlementList(c *gin.Context) {
	var query request.EntitlementSearch
	if err := c.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	list, total, page, pageSize, err := learningAuthzService.GetEntitlementList(query)
	if err != nil {
		global.GVA_LOG.Error("获取学习资源权限列表失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	page, pageSize = normalizeAssetPage(page, pageSize)
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, "获取成功", c)
}
