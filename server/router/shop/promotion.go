package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PromotionRouter struct {}

// InitPromotionRouter 初始化 促销信息 路由信息
func (s *PromotionRouter) InitPromotionRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	PromoRouter := Router.Group("Promo").Use(middleware.OperationRecord())
	PromoRouterWithoutRecord := Router.Group("Promo")
	PromoRouterWithoutAuth := PublicRouter.Group("Promo")
	{
		PromoRouter.POST("createPromotion", PromoApi.CreatePromotion)   // 新建促销信息
		PromoRouter.DELETE("deletePromotion", PromoApi.DeletePromotion) // 删除促销信息
		PromoRouter.DELETE("deletePromotionByIds", PromoApi.DeletePromotionByIds) // 批量删除促销信息
		PromoRouter.PUT("updatePromotion", PromoApi.UpdatePromotion)    // 更新促销信息
	}
	{
		PromoRouterWithoutRecord.GET("findPromotion", PromoApi.FindPromotion)        // 根据ID获取促销信息
		PromoRouterWithoutRecord.GET("getPromotionList", PromoApi.GetPromotionList)  // 获取促销信息列表
	}
	{
	    PromoRouterWithoutAuth.GET("getPromotionPublic", PromoApi.GetPromotionPublic)  // 促销信息开放接口
	}
}
