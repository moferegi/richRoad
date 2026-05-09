package client

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TryonClothRouter struct{}

// InitTryonClothRouter 初始化 我的衣橱 路由信息
func (s *TryonClothRouter) InitTryonClothRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	_ = PublicRouter
	tryonClothRouter := Router.Group("tryonCloth").Use(middleware.OperationRecord())
	tryonClothRouterWithoutRecord := Router.Group("tryonCloth")

	var tryonClothApi = v1.ApiGroupApp.ClientApiGroup.TryonClothApi
	{
		tryonClothRouter.POST("createTryonCloth", tryonClothApi.CreateTryonCloth)             // 创建我的衣橱
		tryonClothRouter.PUT("updateTryonCloth", tryonClothApi.UpdateTryonCloth)              // 更新我的衣橱
		tryonClothRouter.DELETE("deleteTryonCloth", tryonClothApi.DeleteTryonCloth)           // 删除我的衣橱
		tryonClothRouter.DELETE("deleteTryonClothByIds", tryonClothApi.DeleteTryonClothByIds) // 批量删除我的衣橱
	}
	{
		tryonClothRouterWithoutRecord.GET("getMyTryonClothList", tryonClothApi.GetMyTryonClothList) // 获取我的衣橱列表
		tryonClothRouterWithoutRecord.GET("getTryonClothList", tryonClothApi.GetTryonClothList)     // 获取我的衣橱列表(管理端)
	}
}
