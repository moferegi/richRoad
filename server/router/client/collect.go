package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CollectRouter struct {
}

// InitCollectRouter 初始化 收藏 路由信息
func (s *CollectRouter) InitCollectRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	collectRouter := Router.Group("collect").Use(middleware.OperationRecord())
	collectRouterWithoutRecord := Router.Group("collect")
	collectRouterWithoutAuth := PublicRouter.Group("collect")

	var collectApi = v1.ApiGroupApp.ClientApiGroup.CollectApi
	{
		collectRouter.POST("createCollect", collectApi.CreateCollect)             // 新建收藏
		collectRouter.DELETE("deleteCollect", collectApi.DeleteCollect)           // 删除收藏
		collectRouter.DELETE("deleteCollectByIds", collectApi.DeleteCollectByIds) // 批量删除收藏
		collectRouter.PUT("updateCollect", collectApi.UpdateCollect)              // 更新收藏
	}
	{
		collectRouterWithoutRecord.GET("findCollect", collectApi.FindCollect)       // 根据ID获取收藏
		collectRouterWithoutRecord.GET("getCollectList", collectApi.GetCollectList) // 获取收藏列表
	}
	{
		collectRouterWithoutAuth.GET("getCollectPublic", collectApi.GetCollectPublic) // 获取收藏列表
	}
}
