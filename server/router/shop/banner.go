package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BannerRouter struct {
}

// InitBannerRouter 初始化 轮播图 路由信息
func (s *BannerRouter) InitBannerRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	bannerRouter := Router.Group("banner").Use(middleware.OperationRecord())
	bannerRouterWithoutRecord := Router.Group("banner")
	bannerRouterWithoutAuth := PublicRouter.Group("banner")

	var bannerApi = v1.ApiGroupApp.ShopApiGroup.BannerApi
	{
		bannerRouter.POST("createBanner", bannerApi.CreateBanner)             // 新建轮播图
		bannerRouter.DELETE("deleteBanner", bannerApi.DeleteBanner)           // 删除轮播图
		bannerRouter.DELETE("deleteBannerByIds", bannerApi.DeleteBannerByIds) // 批量删除轮播图
		bannerRouter.PUT("updateBanner", bannerApi.UpdateBanner)              // 更新轮播图
	}
	{
		bannerRouterWithoutRecord.GET("findBanner", bannerApi.FindBanner)     // 根据ID获取轮播图
		bannerRouterWithoutAuth.GET("getBannerList", bannerApi.GetBannerList) // 获取轮播图列表
	}
	{
		bannerRouterWithoutAuth.GET("getBannerPublic", bannerApi.GetBannerPublic) // 获取轮播图列表
	}
}
