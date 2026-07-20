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
	// Uni 移动端路由组（排除 Casbin RBAC，挂载 UniSignVerify）
	uniPrivateGroup := routers[0]
	// Web 管理后台路由组（挂载 Casbin RBAC）
	privateGroup := routers[1]
	publicGroup := routers[2]
	{
		clientRouter := router.RouterGroupApp.Client
		clientRouter.InitClientUserRouter(uniPrivateGroup, publicGroup)
		clientRouter.InitAddressRouter(uniPrivateGroup, publicGroup)
		clientRouter.InitCollectRouter(uniPrivateGroup, publicGroup)
		clientRouter.InitPointRecordRouter(uniPrivateGroup, publicGroup)
		clientRouter.InitTryonRechargeOrderRouter(uniPrivateGroup, publicGroup)
		clientRouter.InitTryonTaskRouter(uniPrivateGroup, publicGroup)
		clientRouter.InitTryonModelRouter(uniPrivateGroup, publicGroup)
		clientRouter.InitTryonClothRouter(uniPrivateGroup, publicGroup)
		clientRouter.InitVisitorRouter(uniPrivateGroup, publicGroup)
		clientRouter.InitSysConfigRouter(uniPrivateGroup, publicGroup)
		clientRouter.InitLanguageRouter(uniPrivateGroup, publicGroup)
		clientRouter.InitPhoneAreaCodeRouter(uniPrivateGroup, publicGroup)
		clientRouter.InitSignInRouter(uniPrivateGroup, publicGroup)
		clientRouter.InitExternalLinkDomainRouter(uniPrivateGroup, publicGroup)
		clientRouter.InitVideoTagRouter(privateGroup, publicGroup)
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
		shopRouter.InitCouponOrderUserRouter(privateGroup, publicGroup)
		shopRouter.InitPromotionRouter(privateGroup, publicGroup)
		shopRouter.InitKefuRouter(privateGroup, publicGroup)
		shopRouter.InitGoodPurchaseRouter(privateGroup, publicGroup)
		shopRouter.InitQrcodePaymentRouter(privateGroup, publicGroup)
		shopRouter.InitPopupRouter(privateGroup, publicGroup)
		shopRouter.InitMarketingRewardRouter(privateGroup, publicGroup)
		shopRouter.InitDashboardRouter(privateGroup, publicGroup)
		shopRouter.InitPresaleRouter(privateGroup, publicGroup)
		shopRouter.InitSkuSpecRouter(privateGroup, publicGroup)
	}
	holder(publicGroup, uniPrivateGroup, privateGroup)
}

// 占位方法，保证文件可以正确加载，避免go空变量检测报错，请勿删除。
