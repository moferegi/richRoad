package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PointRecordRouter struct {}

// InitPointRecordRouter 初始化 积分记录管理 路由信息
func (s *PointRecordRouter) InitPointRecordRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	cprRouter := Router.Group("cpr").Use(middleware.OperationRecord())
	cprRouterWithoutRecord := Router.Group("cpr")
	cprRouterWithoutAuth := PublicRouter.Group("cpr")
	{
		cprRouter.POST("createPointRecord", cprApi.CreatePointRecord)   // 新建积分记录管理
		cprRouter.DELETE("deletePointRecord", cprApi.DeletePointRecord) // 删除积分记录管理
		cprRouter.DELETE("deletePointRecordByIds", cprApi.DeletePointRecordByIds) // 批量删除积分记录管理
		cprRouter.PUT("updatePointRecord", cprApi.UpdatePointRecord)    // 更新积分记录管理
	}
	{
		cprRouterWithoutRecord.GET("findPointRecord", cprApi.FindPointRecord)        // 根据ID获取积分记录管理
		cprRouterWithoutRecord.GET("getPointRecordList", cprApi.GetPointRecordList)  // 获取积分记录管理列表
	}
	{
	    cprRouterWithoutAuth.GET("getPointRecordPublic", cprApi.GetPointRecordPublic)  // 积分记录管理开放接口
	}
}
