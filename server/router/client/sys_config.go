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
	sysConfigPublicRouter := PublicRouter.Group("sysConfig")

	var sysConfigApi = v1.ApiGroupApp.ClientApiGroup.SysConfigApi
	{
		sysConfigRouter.PUT("updateSysConfig", sysConfigApi.UpdateSysConfig) // 更新系统参数
	}
	{
		sysConfigRouterWithoutRecord.GET("getSysConfigList", sysConfigApi.GetSysConfigList)       // 获取系统参数列表
		sysConfigRouterWithoutRecord.GET("getSysConfigByGroup", sysConfigApi.GetSysConfigByGroup) // 按分组获取参数
	}
	{
		sysConfigPublicRouter.GET("getAnnouncementConfig", sysConfigApi.GetAnnouncementConfig) // 公告配置(公开)
		sysConfigPublicRouter.GET("getLoginConfig", sysConfigApi.GetLoginConfig)               // 登录配置(公开)
		sysConfigPublicRouter.GET("getSysConfigByKey", sysConfigApi.GetSysConfigByKey)         // 按key获取参数(公开)
	}
}
