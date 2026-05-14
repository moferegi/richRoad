package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DBInspectorRouter struct{}

// InitDBInspectorRouter 初始化数据库巡检路由
func (s *DBInspectorRouter) InitDBInspectorRouter(Router *gin.RouterGroup) {
	dbInspectorRouter := Router.Group("dbInspector").Use(middleware.OperationRecord())
	dbInspectorRouterWithoutRecord := Router.Group("dbInspector")

	{
		dbInspectorRouter.POST("autoFix", dbInspectorApi.AutoFix)                           // 自动修复数据库或Redis
		dbInspectorRouter.POST("deleteRecordsByRange", dbInspectorApi.DeleteRecordsByRange) // 按日期范围真删除
	}
	{
		dbInspectorRouterWithoutRecord.GET("getOverview", dbInspectorApi.GetOverview) // 获取巡检总览
	}
}
