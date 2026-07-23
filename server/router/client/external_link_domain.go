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
		extDomainRouterWithoutRecord.GET("listCloudFiles", extDomainApi.ListCloudFiles)                       // 列出云存储文件
		extDomainRouterWithoutRecord.GET("compareDirectories", extDomainApi.CompareDirectories)               // 多云目录比对
		extDomainRouter.POST("pingCloud", extDomainApi.PingCloud)                                              // 检测云存储连接
		extDomainRouter.POST("deleteCloudFiles", extDomainApi.DeleteCloudFiles)                               // 批量删除云存储文件
		extDomainRouter.POST("uploadCloudFile", extDomainApi.UploadCloudFile)                                 // 上传文件到云存储
		extDomainRouterWithoutRecord.GET("searchCloudFiles", extDomainApi.SearchCloudFiles)                     // 全局搜索云存储文件
		extDomainRouterWithoutRecord.GET("getFileDownloadURL", extDomainApi.GetFileDownloadURL)                 // 获取文件下载链接
		extDomainRouterWithoutRecord.GET("downloadCloudFolder", extDomainApi.DownloadCloudFolder)               // 打包下载目录
	}
	{
		extDomainPublicRouter.GET("getDefaultDomain", extDomainApi.GetDefaultDomain) // 获取默认域名（公开）
	}
}
