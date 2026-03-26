package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type KefuRouter struct {}

// InitKefuRouter 初始化 客服 路由信息
func (s *KefuRouter) InitKefuRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	kefuRouter := Router.Group("kefu").Use(middleware.OperationRecord())
	kefuRouterWithoutRecord := Router.Group("kefu")
	kefuRouterWithoutAuth := PublicRouter.Group("kefu")
	{
		kefuRouter.POST("createKefu", kefuApi.CreateKefu)   // 新建客服
		kefuRouter.DELETE("deleteKefu", kefuApi.DeleteKefu) // 删除客服
		kefuRouter.DELETE("deleteKefuByIds", kefuApi.DeleteKefuByIds) // 批量删除客服
		kefuRouter.PUT("updateKefu", kefuApi.UpdateKefu)    // 更新客服
	}
	{
		kefuRouterWithoutRecord.GET("findKefu", kefuApi.FindKefu)        // 根据ID获取客服
		kefuRouterWithoutRecord.GET("getKefuList", kefuApi.GetKefuList)  // 获取客服列表
	}
	{
	    kefuRouterWithoutAuth.GET("getKefuPublic", kefuApi.GetKefuPublic)  // 客服开放接口
	}
}
