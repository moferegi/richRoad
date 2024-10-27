package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SkuRouter struct {
}

// InitSkuRouter 初始化 sku 路由信息
func (s *SkuRouter) InitSkuRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	skuRouter := Router.Group("sku").Use(middleware.OperationRecord())
	skuRouterWithoutRecord := Router.Group("sku")
	skuRouterWithoutAuth := PublicRouter.Group("sku")

	var skuApi = v1.ApiGroupApp.ShopApiGroup.SkuApi
	{
		skuRouter.POST("createSku", skuApi.CreateSku)             // 新建sku
		skuRouter.DELETE("deleteSku", skuApi.DeleteSku)           // 删除sku
		skuRouter.DELETE("deleteSkuByIds", skuApi.DeleteSkuByIds) // 批量删除sku
		skuRouter.PUT("updateSku", skuApi.UpdateSku)              // 更新sku
	}
	{
		skuRouterWithoutRecord.GET("findSku", skuApi.FindSku)       // 根据ID获取sku
		skuRouterWithoutRecord.GET("getSkuList", skuApi.GetSkuList) // 获取sku列表
	}
	{
		skuRouterWithoutAuth.GET("getSkuPublic", skuApi.GetSkuPublic) // 获取sku列表
	}
}
