package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CartRouter struct {
}

// InitCartRouter 初始化 购物车 路由信息
func (s *CartRouter) InitCartRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	cartRouter := Router.Group("cart").Use(middleware.OperationRecord())
	cartRouterWithoutRecord := Router.Group("cart")
	//cartRouterWithoutAuth := PublicRouter.Group("cart")

	var cartApi = v1.ApiGroupApp.ShopApiGroup.CartApi
	{
		cartRouter.POST("createCart", cartApi.CreateCart)             // 新建购物车
		cartRouter.DELETE("deleteCart", cartApi.DeleteCart)           // 删除购物车
		cartRouter.DELETE("deleteCartByIds", cartApi.DeleteCartByIds) // 批量删除购物车
		cartRouter.PUT("updateCart", cartApi.UpdateCart)              // 更新购物车

		cartRouter.POST("addCart", cartApi.AddCart)        // 加入购物车
		cartRouter.POST("cutCart", cartApi.CutCart)        // 删减购物车
		cartRouter.GET("getSelfCart", cartApi.GetSelfCart) // 获取自身购物车
		cartRouter.GET("clearCart", cartApi.ClearCart)     // 全部删除购物车

	}
	{
		cartRouterWithoutRecord.GET("findCart", cartApi.FindCart)       // 根据ID获取购物车
		cartRouterWithoutRecord.GET("getCartList", cartApi.GetCartList) // 获取购物车列表
	}
	{
		//cartRouterWithoutAuth.GET("getCartPublic", cartApi.GetCartPublic)  // 获取购物车列表
	}
}
