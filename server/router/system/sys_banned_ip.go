package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BannedIPRouter struct{}

// InitBannedIPRouter 注册IP封禁管理路由
func (r *BannedIPRouter) InitBannedIPRouter(Router *gin.RouterGroup) {
	bannedIPRouter := Router.Group("sysBannedIP").Use(middleware.OperationRecord())
	bannedIPRouterWithoutRecord := Router.Group("sysBannedIP")
	bannedIPApi := v1.ApiGroupApp.SystemApiGroup.BannedIPApi
	{
		bannedIPRouter.POST("banIP", bannedIPApi.BanIP)     // 封禁IP
		bannedIPRouter.POST("unbanIP", bannedIPApi.UnbanIP) // 解封IP
	}
	{
		bannedIPRouterWithoutRecord.GET("getBannedIPList", bannedIPApi.GetBannedIPList) // 获取封禁IP列表
		bannedIPRouterWithoutRecord.GET("getAttackStats", bannedIPApi.GetAttackStats)   // 获取攻击统计
	}
}
