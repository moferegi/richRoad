package client

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ExternalLinkDomainRouter struct{}

func (s *ExternalLinkDomainRouter) InitExternalLinkDomainRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	extDomainRouter := Router.Group("extDomain").Use(middleware.OperationRecord())
	extDomainRouterWithoutRecord := Router.Group("extDomain")
	extDomainPublicRouter := PublicRouter.Group("extDomain")

	var extDomainApi = v1.ApiGroupApp.ClientApiGroup.ExternalLinkDomainApi
	{
		extDomainRouter.POST("createExternalLinkDomain", extDomainApi.CreateExternalLinkDomain)   // 创建域名
		extDomainRouter.DELETE("deleteExternalLinkDomain", extDomainApi.DeleteExternalLinkDomain) // 删除域名
		extDomainRouter.PUT("updateExternalLinkDomain", extDomainApi.UpdateExternalLinkDomain)    // 更新域名
		extDomainRouter.POST("setDefaultDomain", extDomainApi.SetDefaultDomain)                   // 设置默认域名
	}
	{
		extDomainRouterWithoutRecord.GET("getExternalLinkDomainList", extDomainApi.GetExternalLinkDomainList) // 获取域名列表
	}
	{
		extDomainPublicRouter.GET("getDefaultDomain", extDomainApi.GetDefaultDomain) // 获取默认域名（公开）
	}
}
