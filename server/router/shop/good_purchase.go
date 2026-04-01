package shop

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type GoodPurchaseRouter struct{}

// InitGoodPurchaseRouter 初始化 进货记录 路由信息
func (s *GoodPurchaseRouter) InitGoodPurchaseRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	goodPurchaseRouter := Router.Group("goodPurchase").Use(middleware.OperationRecord())
	goodPurchaseRouterWithoutRecord := Router.Group("goodPurchase")

	var goodPurchaseApi = v1.ApiGroupApp.ShopApiGroup.GoodPurchaseApi
	{
		goodPurchaseRouter.POST("createGoodPurchase", goodPurchaseApi.CreateGoodPurchase)       // 创建进货记录
		goodPurchaseRouter.DELETE("deleteGoodPurchase", goodPurchaseApi.DeleteGoodPurchase)     // 删除进货记录
		goodPurchaseRouter.PUT("updateGoodPurchase", goodPurchaseApi.UpdateGoodPurchase)        // 更新进货记录
	}
	{
		goodPurchaseRouterWithoutRecord.GET("getGoodPurchaseList", goodPurchaseApi.GetGoodPurchaseList)         // 获取进货记录列表
		goodPurchaseRouterWithoutRecord.GET("getGoodPurchaseSummary", goodPurchaseApi.GetGoodPurchaseSummary) // 获取进货汇总
	}
}
