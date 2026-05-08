package client

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SignInApi struct{}

var signInService = service.ServiceGroupApp.ClientServiceGroup.SignInService

// DoSignIn 用户签到
// @Tags SignIn
// @Summary 用户签到
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "签到成功"
// @Router /signIn/doSignIn [post]
func (api *SignInApi) DoSignIn(c *gin.Context) {
	userID := utils.GetUserID(c)
	if err := signInService.DoSignIn(userID); err != nil {
		global.GVA_LOG.Error("签到失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, err.Error()), c)
	} else {
		// 触发签到奖励（积分+优惠券）
		if err := marketingRewardService.TriggerReward(userID, "sign_in", "sign_in_reward", "签到奖励", 0); err != nil {
			global.GVA_LOG.Error("签到奖励发放失败", zap.Error(err))
		}
		response.OkWithMessage(i18n.T(c, "signInSuccess"), c)
	}
}

// GetSignInStatus 获取今日签到状态
// @Tags SignIn
// @Summary 获取今日签到状态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /signIn/getSignInStatus [get]
func (api *SignInApi) GetSignInStatus(c *gin.Context) {
	userID := utils.GetUserID(c)
	signed := signInService.GetSignInStatus(userID)
	continuousDays := signInService.GetContinuousSignInDays(userID)
	response.OkWithDetailed(gin.H{
		"signed":         signed,
		"continuousDays": continuousDays,
	}, i18n.T(c, "getSuccess"), c)
}

// GetSignInRecords 获取签到记录
// @Tags SignIn
// @Summary 获取签到记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=[]client.SignIn,msg=string} "获取成功"
// @Router /signIn/getSignInRecords [get]
func (api *SignInApi) GetSignInRecords(c *gin.Context) {
	userID := utils.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "30"))
	if records, total, err := signInService.GetSignInRecords(userID, page, pageSize); err != nil {
		global.GVA_LOG.Error("获取签到记录失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     records,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		}, i18n.T(c, "getSuccess"), c)
	}
}

// GetSignInList 管理端获取签到记录列表
// @Tags SignIn
// @Summary 管理端获取所有用户签到记录列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param userId query string false "用户ID"
// @Param username query string false "用户名(模糊搜索)"
// @Param orderKey query string false "排序字段(sign_date/total_days)"
// @Param desc query bool false "是否降序"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /signIn/getSignInList [get]
func (api *SignInApi) GetSignInList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	userId := c.Query("userId")
	username := c.Query("username")
	orderKey := c.Query("orderKey")
	desc := c.Query("desc") == "true"
	startDate := ""
	endDate := ""
	dateRange := c.QueryArray("signDateRange[]")
	if len(dateRange) == 2 {
		startDate = dateRange[0]
		endDate = dateRange[1]
	}

	if list, total, err := signInService.GetSignInList(page, pageSize, userId, username, startDate, endDate, orderKey, desc); err != nil {
		global.GVA_LOG.Error("获取签到列表失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		}, i18n.T(c, "getSuccess"), c)
	}
}

// DeleteSignIn 管理端删除签到记录
// @Tags SignIn
// @Summary 管理端删除签到记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query int true "签到记录ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /signIn/deleteSignIn [delete]
func (api *SignInApi) DeleteSignIn(c *gin.Context) {
	idStr := c.Query("ID")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if err := signInService.DeleteSignIn(uint(id)); err != nil {
		global.GVA_LOG.Error("删除签到记录失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "deleteFail"), c)
	} else {
		response.OkWithMessage(i18n.T(c, "deleteSuccess"), c)
	}
}
