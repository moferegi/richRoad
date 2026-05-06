package client

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TryonModelRouter struct{}

// InitTryonModelRouter 初始化 我的模特 路由信息
func (s *TryonModelRouter) InitTryonModelRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	_ = PublicRouter
	tryonModelRouter := Router.Group("tryonModel").Use(middleware.OperationRecord())
	tryonModelRouterWithoutRecord := Router.Group("tryonModel")

	var tryonModelApi = v1.ApiGroupApp.ClientApiGroup.TryonModelApi
	{
		tryonModelRouter.POST("createTryonModel", tryonModelApi.CreateTryonModel)             // 创建我的模特
		tryonModelRouter.PUT("updateTryonModel", tryonModelApi.UpdateTryonModel)              // 重命名我的模特
		tryonModelRouter.DELETE("deleteTryonModel", tryonModelApi.DeleteTryonModel)           // 删除我的模特
		tryonModelRouter.DELETE("deleteTryonModelByIds", tryonModelApi.DeleteTryonModelByIds) // 批量删除我的模特
	}
	{
		tryonModelRouterWithoutRecord.GET("getMyTryonModelList", tryonModelApi.GetMyTryonModelList) // 获取我的模特列表
		tryonModelRouterWithoutRecord.GET("getTryonModelList", tryonModelApi.GetTryonModelList)     // 获取模特列表(管理端)
	}
}
