package client

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ClientUserRouter struct {
}

// InitClientUserRouter 初始化 客户端用户 路由信息
func (s *ClientUserRouter) InitClientUserRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	clientUserRouter := Router.Group("clientUser").Use(middleware.OperationRecord())
	clientUserRouterWithoutRecord := Router.Group("clientUser")
	clientUserRouterWithoutAuth := PublicRouter.Group("clientUser")

	var clientUserApi = v1.ApiGroupApp.ClientApiGroup.ClientUserApi
	{
		clientUserRouter.POST("createClientUser", clientUserApi.CreateClientUser)             // 新建客户端用户
		clientUserRouter.DELETE("deleteClientUser", clientUserApi.DeleteClientUser)           // 删除客户端用户
		clientUserRouter.DELETE("deleteClientUserByIds", clientUserApi.DeleteClientUserByIds) // 批量删除客户端用户
		clientUserRouter.PUT("updateClientUser", clientUserApi.UpdateClientUser)              // 更新客户端用户
		clientUserRouter.GET("getUserInfo", clientUserApi.GetUserInfo)                        // 客户端获取自身信息
		clientUserRouter.POST("changePassword", clientUserApi.ChangePassword)                 // 修改密码
	}
	{
		clientUserRouterWithoutRecord.GET("findClientUser", clientUserApi.FindClientUser)       // 根据ID获取客户端用户
		clientUserRouterWithoutRecord.GET("getClientUserList", clientUserApi.GetClientUserList) // 获取客户端用户列表
		clientUserRouterWithoutRecord.GET("getSubordinates", clientUserApi.GetSubordinates)     // 获取下级用户列表
		clientUserRouterWithoutRecord.GET("getMyInviteInfo", clientUserApi.GetMyInviteInfo)     // 获取我的邀请信息
		clientUserRouterWithoutRecord.GET("getMySubordinates", clientUserApi.GetMySubordinates) // 获取我的下级列表
	}
	{
		clientUserRouterWithoutAuth.GET("getOpenID", clientUserApi.GetOpenID)                  // 获取小程序openid
		clientUserRouterWithoutAuth.POST("login", clientUserApi.Login)                         // 客户端登录
		clientUserRouterWithoutAuth.POST("register", clientUserApi.Register)                   // 客户端注册
		clientUserRouterWithoutAuth.POST("phoneLogin", clientUserApi.PhoneLogin)               // 手机号登录
		clientUserRouterWithoutAuth.POST("phoneRegister", clientUserApi.PhoneRegister)         // 手机号注册
		clientUserRouterWithoutAuth.POST("setClientUserInfo", clientUserApi.SetClientUserInfo) // 客户端设置个人信息
	}
}
