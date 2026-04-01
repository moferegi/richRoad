package shop

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PresaleRouter struct{}

// InitPresaleRouter 初始化 预售 路由信息
func (s *PresaleRouter) InitPresaleRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	presaleRouter := Router.Group("presale").Use(middleware.OperationRecord())
	presaleRouterWithoutAuth := PublicRouter.Group("presale")

	var presaleApi = v1.ApiGroupApp.ShopApiGroup.PresaleApi
	{
		presaleRouter.GET("getPresaleParticipants", presaleApi.GetPresaleParticipants) // 管理端获取预售参与者
	}
	{
		presaleRouterWithoutAuth.GET("getPresaleGoodList", presaleApi.GetPresaleGoodList)  // 客户端获取预售商品列表
		presaleRouterWithoutAuth.GET("checkAvailable", presaleApi.CheckPresaleAvailable)    // 客户端检查预售可用性
	}
}
