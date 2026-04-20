package client

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type SignInRouter struct{}

// InitSignInRouter 初始化 签到 路由信息
func (s *SignInRouter) InitSignInRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	signInRouter := Router.Group("signIn")

	var signInApi = v1.ApiGroupApp.ClientApiGroup.SignInApi
	{
		signInRouter.POST("doSignIn", signInApi.DoSignIn)                // 用户签到
		signInRouter.GET("getSignInStatus", signInApi.GetSignInStatus)   // 获取今日签到状态
		signInRouter.GET("getSignInRecords", signInApi.GetSignInRecords) // 获取签到记录
		signInRouter.GET("getSignInList", signInApi.GetSignInList)       // 管理端：获取签到列表
		signInRouter.DELETE("deleteSignIn", signInApi.DeleteSignIn)      // 管理端：删除签到记录
	}
}
