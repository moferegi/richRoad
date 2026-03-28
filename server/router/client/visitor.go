package client

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type VisitorRouter struct {
}

// InitVisitorRouter 初始化 访客 路由信息
func (s *VisitorRouter) InitVisitorRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	visitorRouter := Router.Group("visitor").Use(middleware.OperationRecord())
	visitorRouterWithoutRecord := Router.Group("visitor")
	visitorRouterWithoutAuth := PublicRouter.Group("visitor")

	var visitorApi = v1.ApiGroupApp.ClientApiGroup.VisitorApi
	{
		visitorRouter.POST("aggregateDailySummary", visitorApi.AggregateDailySummary) // 手动触发日汇总
	}
	{
		visitorRouterWithoutRecord.GET("getVisitorLogList", visitorApi.GetVisitorLogList)         // 获取访客日志列表
		visitorRouterWithoutRecord.GET("getVisitorSummaryList", visitorApi.GetVisitorSummaryList) // 获取访客汇总列表
		visitorRouterWithoutRecord.GET("getTodayStats", visitorApi.GetTodayStats)                 // 获取今日统计
	}
	{
		visitorRouterWithoutAuth.POST("heartbeat", visitorApi.Heartbeat) // 访客心跳上报（无需认证）
	}
}
