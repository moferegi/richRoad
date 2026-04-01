package shop

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DashboardRouter struct{}

// InitDashboardRouter 初始化 数据看板 路由信息
func (s *DashboardRouter) InitDashboardRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	dashboardRouter := Router.Group("dashboard").Use(middleware.OperationRecord())

	var dashboardApi = v1.ApiGroupApp.ShopApiGroup.DashboardApi
	{
		dashboardRouter.GET("getOverview", dashboardApi.GetDashboardOverview) // 获取数据看板概览
	}
}
