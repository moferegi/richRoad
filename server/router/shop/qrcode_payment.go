package shop

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type QrcodePaymentRouter struct{}

// InitQrcodePaymentRouter 初始化 收款码 路由信息
func (s *QrcodePaymentRouter) InitQrcodePaymentRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	qrcodePaymentRouter := Router.Group("qrcodePayment").Use(middleware.OperationRecord())
	qrcodePaymentRouterWithoutAuth := PublicRouter.Group("qrcodePayment")

	var qrcodePaymentApi = v1.ApiGroupApp.ShopApiGroup.QrcodePaymentApi
	{
		qrcodePaymentRouter.POST("createQrcodePayment", qrcodePaymentApi.CreateQrcodePayment)       // 创建收款码
		qrcodePaymentRouter.DELETE("deleteQrcodePayment", qrcodePaymentApi.DeleteQrcodePayment)     // 删除收款码
		qrcodePaymentRouter.PUT("updateQrcodePayment", qrcodePaymentApi.UpdateQrcodePayment)        // 更新收款码
		qrcodePaymentRouter.GET("getQrcodePaymentList", qrcodePaymentApi.GetQrcodePaymentList)      // 获取收款码列表
	}
	{
		qrcodePaymentRouterWithoutAuth.GET("getEnabledQrcodePayments", qrcodePaymentApi.GetEnabledQrcodePayments) // 客户端获取启用的收款码
	}
}
