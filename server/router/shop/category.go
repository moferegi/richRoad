package shop

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CategoryRouter struct {
}

// InitCategoryRouter 初始化 商品分类 路由信息
func (s *CategoryRouter) InitCategoryRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	categoryRouter := Router.Group("category").Use(middleware.OperationRecord())
	categoryRouterWithoutRecord := Router.Group("category")
	categoryRouterWithoutAuth := PublicRouter.Group("category")

	var categoryApi = v1.ApiGroupApp.ShopApiGroup.CategoryApi
	{
		categoryRouter.POST("createCategory", categoryApi.CreateCategory)             // 新建商品分类
		categoryRouter.DELETE("deleteCategory", categoryApi.DeleteCategory)           // 删除商品分类
		categoryRouter.DELETE("deleteCategoryByIds", categoryApi.DeleteCategoryByIds) // 批量删除商品分类
		categoryRouter.PUT("updateCategory", categoryApi.UpdateCategory)              // 更新商品分类
	}
	{
		categoryRouterWithoutRecord.GET("findCategory", categoryApi.FindCategory)                                 // 根据ID获取商品分类
		categoryRouterWithoutRecord.GET("getCategoryList", categoryApi.GetCategoryList)                           // 获取商品分类列表
		categoryRouterWithoutAuth.GET("getCategoryMobile", categoryApi.GetCategoryMobile)                         // 获取商品分类列表
		categoryRouterWithoutAuth.GET("getChildrenCategoryAndProduct", categoryApi.GetChildrenCategoryAndProduct) // 获取商品分类列表
	}
	{
		categoryRouterWithoutAuth.GET("getCategoryPublic", categoryApi.GetCategoryPublic) // 获取商品分类列表
	}
}
