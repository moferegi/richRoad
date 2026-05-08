package client

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

// TryonRechargeOrderRouter 试衣币充值订单路由
type TryonRechargeOrderRouter struct{}

// InitTryonRechargeOrderRouter 初始化试衣币充值订单路由
func (s *TryonRechargeOrderRouter) InitTryonRechargeOrderRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	_ = PublicRouter
	tryonRechargeOrderRouter := Router.Group("tryonRechargeOrder")
	tryonRechargeOrderApi := v1.ApiGroupApp.ClientApiGroup.TryonRechargeOrderApi
	{
		tryonRechargeOrderRouter.POST("createTryonRechargeOrder", tryonRechargeOrderApi.CreateTryonRechargeOrder)
		tryonRechargeOrderRouter.POST("updateTryonRechargeOrderPayMethod", tryonRechargeOrderApi.UpdateTryonRechargeOrderPayMethod)
		tryonRechargeOrderRouter.POST("submitTryonRechargeOrderPayment", tryonRechargeOrderApi.SubmitTryonRechargeOrderPayment)
		tryonRechargeOrderRouter.POST("cancelTryonRechargeOrder", tryonRechargeOrderApi.CancelTryonRechargeOrder)
		tryonRechargeOrderRouter.POST("confirmTryonRechargeOrderPayment", tryonRechargeOrderApi.ConfirmTryonRechargeOrderPayment)
		tryonRechargeOrderRouter.GET("selfTryonRechargeOrder", tryonRechargeOrderApi.SelfTryonRechargeOrder)
		tryonRechargeOrderRouter.GET("findTryonRechargeOrder", tryonRechargeOrderApi.FindTryonRechargeOrder)
		tryonRechargeOrderRouter.GET("getMyTryonRechargeOrderList", tryonRechargeOrderApi.GetMyTryonRechargeOrderList)
		tryonRechargeOrderRouter.GET("getTryonRechargeOrderList", tryonRechargeOrderApi.GetTryonRechargeOrderList)
	}
}
