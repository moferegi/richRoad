package client

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
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
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("签到成功", c)
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
	}, "获取成功", c)
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
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     records,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		}, "获取成功", c)
	}
}
