package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SkuSpecRouter struct{}

// InitSkuSpecRouter 初始化 SKU规格字典 路由信息
func (s *SkuSpecRouter) InitSkuSpecRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	skuSpecRouter := Router.Group("skuSpec").Use(middleware.OperationRecord())
	skuSpecRouterWithoutRecord := Router.Group("skuSpec")
	{
		skuSpecRouter.POST("createSkuSpec", skuSpecApi.CreateSkuSpec)             // 新建SKU规格字典
		skuSpecRouter.DELETE("deleteSkuSpec", skuSpecApi.DeleteSkuSpec)           // 删除SKU规格字典
		skuSpecRouter.DELETE("deleteSkuSpecByIds", skuSpecApi.DeleteSkuSpecByIds) // 批量删除
		skuSpecRouter.PUT("updateSkuSpec", skuSpecApi.UpdateSkuSpec)              // 更新SKU规格字典
	}
	{
		skuSpecRouterWithoutRecord.GET("findSkuSpec", skuSpecApi.FindSkuSpec)       // 根据ID获取
		skuSpecRouterWithoutRecord.GET("getSkuSpecList", skuSpecApi.GetSkuSpecList) // 获取列表
		skuSpecRouterWithoutRecord.GET("getAllSkuSpecs", skuSpecApi.GetAllSkuSpecs) // 获取全部（下拉）
	}
	_ = PublicRouter
}
