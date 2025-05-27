package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type OrderRouter struct {
}

// InitOrderRouter 初始化 订单 路由信息
func (s *OrderRouter) InitOrderRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	orderRouter := Router.Group("order").Use(middleware.OperationRecord())
	orderRouterWithoutRecord := Router.Group("order")
	orderRouterWithoutAuth := PublicRouter.Group("order")

	var orderApi = v1.ApiGroupApp.ShopApiGroup.OrderApi
	{
		orderRouter.POST("createOrder", orderApi.CreateOrder)             // 新建订单
		orderRouter.DELETE("deleteOrder", orderApi.DeleteOrder)           // 删除订单
		orderRouter.DELETE("deleteOrderByIds", orderApi.DeleteOrderByIds) // 批量删除订单
		orderRouter.PUT("updateOrder", orderApi.UpdateOrder)              // 更新订单

		orderRouter.POST("placeOrder", orderApi.PlaceOrder)               // 下单
		orderRouter.POST("placeOrderByCart", orderApi.PlaceOrderByCart)   // 购物车下单
		orderRouter.POST("changeOrderCoupon", orderApi.ChangeOrderCoupon) // 变更订单优惠券
		orderRouter.POST("updateOrderStatus", orderApi.UpdateOrderStatus) // 更新订单状态
		orderRouter.GET("selfOrder", orderApi.SelfOrder)                  // 获取单一订单
		orderRouter.GET("selfOrderList", orderApi.SelfOrderList)          // 获取订单列表
		orderRouter.GET("selfOrderComment", orderApi.SelfOrderComment)    // 获取评价详情

		orderRouter.POST("updateOrder", orderApi.UpdateOrder) // 更新订单

	}
	{
		orderRouterWithoutRecord.GET("findOrder", orderApi.FindOrder)       // 根据ID获取订单
		orderRouterWithoutRecord.GET("getOrderList", orderApi.GetOrderList) // 获取订单列表
	}
	{
		orderRouterWithoutAuth.GET("getOrderPublic", orderApi.GetOrderPublic) // 获取订单列表
		orderRouterWithoutAuth.GET("checkRouters", orderApi.CheckRouters)     // 获取订单列表
	}
}
