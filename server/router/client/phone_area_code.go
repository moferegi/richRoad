package client

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PhoneAreaCodeRouter struct{}

// InitPhoneAreaCodeRouter 初始化 国际区号 路由信息
func (s *PhoneAreaCodeRouter) InitPhoneAreaCodeRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	phoneAreaCodeRouter := Router.Group("phoneAreaCode").Use(middleware.OperationRecord())
	phoneAreaCodeRouterWithoutAuth := PublicRouter.Group("phoneAreaCode")

	var phoneAreaCodeApi = v1.ApiGroupApp.ClientApiGroup.PhoneAreaCodeApi
	{
		phoneAreaCodeRouter.POST("createPhoneAreaCode", phoneAreaCodeApi.CreatePhoneAreaCode)       // 创建区号
		phoneAreaCodeRouter.DELETE("deletePhoneAreaCode", phoneAreaCodeApi.DeletePhoneAreaCode)     // 删除区号
		phoneAreaCodeRouter.PUT("updatePhoneAreaCode", phoneAreaCodeApi.UpdatePhoneAreaCode)        // 更新区号
		phoneAreaCodeRouter.GET("getPhoneAreaCodeList", phoneAreaCodeApi.GetPhoneAreaCodeList)      // 获取区号列表（管理端）
	}
	{
		phoneAreaCodeRouterWithoutAuth.GET("getEnabledPhoneAreaCodes", phoneAreaCodeApi.GetEnabledPhoneAreaCodes) // 客户端获取启用的区号列表
	}
}
