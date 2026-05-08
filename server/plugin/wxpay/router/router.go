package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/api"
	"github.com/gin-gonic/gin"
)

type WxpayRouter struct {
}

func (s *WxpayRouter) InitWxpayRouter(Router *gin.RouterGroup) {
	authRouter := Router.Group("").Use(middleware.JWTAuth())
	publicRouter := Router.Group("")
	plugApi := api.ApiGroupApp.WxpayApi
	{
		authRouter.POST("getPayCode", plugApi.GetPayCode)
		authRouter.POST("getPayParams", plugApi.GetPayParams)
		authRouter.POST("checkNeedPay", plugApi.CheckNeedPay)

		authRouter.GET("getOrderById", plugApi.GetOrderById)
	}

	{
		publicRouter.POST("payAction", plugApi.PayAction)
	}
}
