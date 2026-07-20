package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type VideoTagRouter struct{}

// InitVideoTagRouter 初始化 视频标签 路由信息（Web 后台管理端用）
func (r *VideoTagRouter) InitVideoTagRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	videoTagRouter := Router.Group("videoTag").Use(middleware.OperationRecord())
	videoTagRouterWithoutRecord := Router.Group("videoTag")
	videoTagRouterWithoutAuth := PublicRouter.Group("videoTag")
	{
		videoTagRouter.POST("createVideoTag", videoTagApi.CreateVideoTag)             // 新建视频标签
		videoTagRouter.DELETE("deleteVideoTag", videoTagApi.DeleteVideoTag)           // 删除视频标签
		videoTagRouter.DELETE("deleteVideoTagByIds", videoTagApi.DeleteVideoTagByIds) // 批量删除视频标签
		videoTagRouter.PUT("updateVideoTag", videoTagApi.UpdateVideoTag)              // 更新视频标签
	}
	{
		videoTagRouterWithoutRecord.GET("findVideoTag", videoTagApi.FindVideoTag)          // 根据ID获取视频标签
		videoTagRouterWithoutRecord.GET("getVideoTagList", videoTagApi.GetVideoTagList)    // 获取视频标签列表
	}
	{
		videoTagRouterWithoutAuth.GET("getVideoTagPublic", videoTagApi.GetVideoTagPublic) // 视频标签开放接口
	}
}
