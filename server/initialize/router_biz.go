package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router"
	"github.com/gin-gonic/gin"
)

func holder(routers ...*gin.RouterGroup) {
	_ = routers
	_ = router.RouterGroupApp
}
func initBizRouter(routers ...*gin.RouterGroup) {
	privateGroup := routers[0]
	publicGroup := routers[1]
	{
		clientRouter := router.RouterGroupApp.Client
		clientRouter.InitClientUserRouter(privateGroup, publicGroup)
		clientRouter.InitAddressRouter(privateGroup, publicGroup)
		clientRouter.InitCollectRouter(privateGroup, publicGroup)
	}
	{
		shopRouter := router.RouterGroupApp.Shop
		shopRouter.InitBannerRouter(privateGroup, publicGroup)
		shopRouter.InitCategoryRouter(privateGroup, publicGroup)
		shopRouter.InitGoodRouter(privateGroup, publicGroup)
		shopRouter.InitSkuRouter(privateGroup, publicGroup)
		shopRouter.InitCartRouter(privateGroup, publicGroup)
		shopRouter.InitOrderRouter(privateGroup, publicGroup)
		shopRouter.InitCommentRouter(privateGroup, publicGroup)
		shopRouter.InitTagRouter(privateGroup, publicGroup)
		shopRouter.InitCouponRouter(privateGroup, publicGroup)
		shopRouter.InitCouponOrderUserRouter(privateGroup, publicGroup) // 占位方法，保证文件可以正确加载，避免go空变量检测报错，请勿删除。
		shopRouter.InitPromotionRouter(privateGroup, publicGroup)
	}
	holder(publicGroup, privateGroup)
}
