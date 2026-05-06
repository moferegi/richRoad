package client

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type TryonTaskRouter struct{}

// InitTryonTaskRouter 初始化 试衣任务 路由信息
func (s *TryonTaskRouter) InitTryonTaskRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	_ = PublicRouter
	tryonTaskRouter := Router.Group("tryonTask")

	var tryonTaskApi = v1.ApiGroupApp.ClientApiGroup.TryonTaskApi
	{
		tryonTaskRouter.POST("createTryonTask", tryonTaskApi.CreateTryonTask)             // 创建试衣任务
		tryonTaskRouter.DELETE("deleteTryonTask", tryonTaskApi.DeleteTryonTask)           // 删除试衣任务(管理端)
		tryonTaskRouter.DELETE("deleteTryonTaskByIds", tryonTaskApi.DeleteTryonTaskByIds) // 批量删除试衣任务(管理端)
		tryonTaskRouter.GET("findTryonTask", tryonTaskApi.FindTryonTask)                  // 根据ID查询试衣任务
		tryonTaskRouter.GET("getMyTryonTaskList", tryonTaskApi.GetMyTryonTaskList)        // 获取我的试衣任务列表
		tryonTaskRouter.GET("getTryonTaskList", tryonTaskApi.GetTryonTaskList)            // 获取试衣任务列表(管理端)
		tryonTaskRouter.GET("getTryonTaskStats", tryonTaskApi.GetTryonTaskStats)          // 获取试衣任务统计(管理端)
		tryonTaskRouter.GET("getTryonTaskTrend", tryonTaskApi.GetTryonTaskTrend)          // 获取试衣任务趋势(管理端)
	}
}
