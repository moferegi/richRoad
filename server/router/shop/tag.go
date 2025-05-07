package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TagRouter struct {}

// InitTagRouter 初始化 标签 路由信息
func (s *TagRouter) InitTagRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	tagRouter := Router.Group("tag").Use(middleware.OperationRecord())
	tagRouterWithoutRecord := Router.Group("tag")
	tagRouterWithoutAuth := PublicRouter.Group("tag")
	{
		tagRouter.POST("createTag", tagApi.CreateTag)   // 新建标签
		tagRouter.DELETE("deleteTag", tagApi.DeleteTag) // 删除标签
		tagRouter.DELETE("deleteTagByIds", tagApi.DeleteTagByIds) // 批量删除标签
		tagRouter.PUT("updateTag", tagApi.UpdateTag)    // 更新标签
	}
	{
		tagRouterWithoutRecord.GET("findTag", tagApi.FindTag)        // 根据ID获取标签
		tagRouterWithoutRecord.GET("getTagList", tagApi.GetTagList)  // 获取标签列表
	}
	{
	    tagRouterWithoutAuth.GET("getTagPublic", tagApi.GetTagPublic)  // 标签开放接口
	}
}
