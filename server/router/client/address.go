package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AddressRouter struct {
}

// InitAddressRouter 初始化 用户地址 路由信息
func (s *AddressRouter) InitAddressRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	addressRouter := Router.Group("address").Use(middleware.OperationRecord())
	addressRouterWithoutRecord := Router.Group("address")
	addressRouterWithoutAuth := PublicRouter.Group("address")

	var addressApi = v1.ApiGroupApp.ClientApiGroup.AddressApi
	{
		addressRouter.POST("createAddress", addressApi.CreateAddress)             // 新建用户地址
		addressRouter.DELETE("deleteAddress", addressApi.DeleteAddress)           // 删除用户地址
		addressRouter.DELETE("deleteAddressByIds", addressApi.DeleteAddressByIds) // 批量删除用户地址
		addressRouter.PUT("updateAddress", addressApi.UpdateAddress)              // 更新用户地址
	}
	{
		addressRouterWithoutRecord.GET("findAddress", addressApi.FindAddress)             // 根据ID获取用户地址
		addressRouterWithoutRecord.GET("getDefaultAddress", addressApi.GetDefaultAddress) // 获取用户默认地址
		addressRouterWithoutRecord.GET("getAddressList", addressApi.GetAddressList)       // 获取用户地址列表
	}
	{
		addressRouterWithoutAuth.GET("getAddressDataSource", addressApi.GetAddressDataSource) // 获取用户地址数据源
		addressRouterWithoutAuth.GET("getAddressPublic", addressApi.GetAddressPublic)         // 获取用户地址列表
	}
}
