package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/api"
	"github.com/gin-gonic/gin"
)

type UserLearningAssetRouter struct{}

func (s *UserLearningAssetRouter) InitUserLearningAssetRouter(Router *gin.RouterGroup) {
	// authRouter 指的是这里应该经过了 JWT 等验证，如果通过全局挂载就不需要二次挂
	assetRouter := Router.Group("asset")
	var assetApi = api.ApiGroupApp.UserLearningAssetApi
	{
		assetRouter.POST("heartbeat", assetApi.Heartbeat) // 发送观影心跳
		assetRouter.GET("getFreeTimeRecordList", assetApi.GetFreeTimeRecordList)
		assetRouter.POST("grantEntitlement", assetApi.GrantEntitlement)
		assetRouter.DELETE("revokeEntitlement", assetApi.RevokeEntitlement)
		assetRouter.GET("getEntitlementList", assetApi.GetEntitlementList)
	}
}
