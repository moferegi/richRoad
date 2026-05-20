package shop

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type GoodRouter struct {
}

// InitGoodRouter 初始化 商品 路由信息
func (s *GoodRouter) InitGoodRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	goodRouter := Router.Group("good").Use(middleware.OperationRecord())
	goodRouterWithoutRecord := Router.Group("good")
	goodRouterWithoutAuth := PublicRouter.Group("good")

	var goodApi = v1.ApiGroupApp.ShopApiGroup.GoodApi
	{
		goodRouter.POST("createGood", goodApi.CreateGood)             // 新建商品
		goodRouter.DELETE("deleteGood", goodApi.DeleteGood)           // 删除商品
		goodRouter.DELETE("deleteGoodByIds", goodApi.DeleteGoodByIds) // 批量删除商品
		goodRouter.PUT("updateGood", goodApi.UpdateGood)              // 更新商品
		goodRouter.DELETE("clearGoodHistory", goodApi.ClearGoodHistory)
	}
	{
		goodRouterWithoutAuth.GET("findGood", goodApi.FindGood)             // 根据ID获取商品
		goodRouterWithoutAuth.GET("getGoodList", goodApi.GetGoodList)       // 获取商品列表
		goodRouterWithoutAuth.GET("getGoodHistory", goodApi.GetGoodHistory) // 获取商品列表
	}
	{
		goodRouterWithoutRecord.GET("getGoodPublic", goodApi.GetGoodPublic) // 获取商品列表
	}
}
