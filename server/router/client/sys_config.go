package client

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysConfigRouter struct{}

func (s *SysConfigRouter) InitSysConfigRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	sysConfigRouter := Router.Group("sysConfig").Use(middleware.OperationRecord())
	sysConfigRouterWithoutRecord := Router.Group("sysConfig")
	_ = PublicRouter

	var sysConfigApi = v1.ApiGroupApp.ClientApiGroup.SysConfigApi
	{
		sysConfigRouter.PUT("updateSysConfig", sysConfigApi.UpdateSysConfig) // 更新系统参数
	}
	{
		sysConfigRouterWithoutRecord.GET("getSysConfigList", sysConfigApi.GetSysConfigList) // 获取系统参数列表
	}
}
