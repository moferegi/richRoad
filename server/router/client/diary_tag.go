package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DiaryTagRouter struct{}

// InitDiaryTagRouter 初始化 日记标签 路由信息
func (r *DiaryTagRouter) InitDiaryTagRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	diaryTagRouter := Router.Group("diaryTag").Use(middleware.OperationRecord())
	diaryTagRouterWithoutRecord := Router.Group("diaryTag")
	diaryTagRouterWithoutAuth := PublicRouter.Group("diaryTag")
	{
		diaryTagRouter.POST("createDiaryTag", diaryTagApi.CreateDiaryTag)
		diaryTagRouter.DELETE("deleteDiaryTag", diaryTagApi.DeleteDiaryTag)
		diaryTagRouter.DELETE("deleteDiaryTagByIds", diaryTagApi.DeleteDiaryTagByIds)
		diaryTagRouter.PUT("updateDiaryTag", diaryTagApi.UpdateDiaryTag)
	}
	{
		diaryTagRouterWithoutRecord.GET("findDiaryTag", diaryTagApi.FindDiaryTag)
		diaryTagRouterWithoutRecord.GET("getDiaryTagList", diaryTagApi.GetDiaryTagList)
	}
	{
		diaryTagRouterWithoutAuth.GET("getDiaryTagPublic", diaryTagApi.GetDiaryTagPublic)
	}
}