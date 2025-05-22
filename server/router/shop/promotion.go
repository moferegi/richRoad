package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PromotionRouter struct{}

// InitPromotionRouter 初始化 促销信息 路由信息
func (s *PromotionRouter) InitPromotionRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	promoRouter := Router.Group("promo").Use(middleware.OperationRecord())
	promoRouterWithoutRecord := Router.Group("promo")
	promoRouterWithoutAuth := PublicRouter.Group("promo")
	{
		promoRouter.POST("createPromotion", PromoApi.CreatePromotion)             // 新建促销信息
		promoRouter.DELETE("deletePromotion", PromoApi.DeletePromotion)           // 删除促销信息
		promoRouter.DELETE("deletePromotionByIds", PromoApi.DeletePromotionByIds) // 批量删除促销信息
		promoRouter.PUT("updatePromotion", PromoApi.UpdatePromotion)              // 更新促销信息
	}
	{
		promoRouterWithoutRecord.GET("findPromotion", PromoApi.FindPromotion)       // 根据ID获取促销信息
		promoRouterWithoutRecord.GET("getPromotionList", PromoApi.GetPromotionList) // 获取促销信息列表
	}
	{
		promoRouterWithoutAuth.GET("getPromotionPublic", PromoApi.GetPromotionPublic) // 促销信息开放接口
	}
}
