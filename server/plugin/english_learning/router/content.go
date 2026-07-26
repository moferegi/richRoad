package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/api"
	"github.com/gin-gonic/gin"
)

type ContentRouter struct{}

func (s *ContentRouter) InitContentRouter(Router *gin.RouterGroup) {
	contentRouter := Router.Group("content")
	var contentApi = api.ApiGroupApp.ContentApi
	{
		contentRouter.POST("createCategory", contentApi.CreateCategory)
		contentRouter.PUT("updateCategory", contentApi.UpdateCategory)
		contentRouter.DELETE("deleteCategory", contentApi.DeleteCategory)
		contentRouter.GET("findCategory", contentApi.FindCategory)
		contentRouter.GET("getCategoryList", contentApi.GetCategoryList)

		contentRouter.POST("createChapter", contentApi.CreateChapter)
		contentRouter.PUT("updateChapter", contentApi.UpdateChapter)
		contentRouter.DELETE("deleteChapter", contentApi.DeleteChapter)
		contentRouter.GET("findChapter", contentApi.FindChapter)
		contentRouter.GET("getChapterList", contentApi.GetChapterList)

		contentRouter.POST("createVideoCategory", contentApi.CreateVideoCategory)
		contentRouter.PUT("updateVideoCategory", contentApi.UpdateVideoCategory)
		contentRouter.DELETE("deleteVideoCategory", contentApi.DeleteVideoCategory)
		contentRouter.GET("findVideoCategory", contentApi.FindVideoCategory)
		contentRouter.GET("getVideoCategoryList", contentApi.GetVideoCategoryList)

		contentRouter.POST("createVideoSeries", contentApi.CreateVideoSeries)
		contentRouter.PUT("updateVideoSeries", contentApi.UpdateVideoSeries)
		contentRouter.DELETE("deleteVideoSeries", contentApi.DeleteVideoSeries)
		contentRouter.GET("findVideoSeries", contentApi.FindVideoSeries)
		contentRouter.GET("getVideoSeriesList", contentApi.GetVideoSeriesList)

		contentRouter.POST("createVideoEpisode", contentApi.CreateVideoEpisode)
		contentRouter.PUT("updateVideoEpisode", contentApi.UpdateVideoEpisode)
		contentRouter.DELETE("deleteVideoEpisode", contentApi.DeleteVideoEpisode)
		contentRouter.GET("findVideoEpisode", contentApi.FindVideoEpisode)
		contentRouter.GET("getVideoEpisodeList", contentApi.GetVideoEpisodeList)
		contentRouter.GET("getVideoEpisodeListByTag", contentApi.GetVideoEpisodeListByTag)

		contentRouter.GET("checkFfmpeg", contentApi.CheckFfmpeg)
		contentRouter.POST("sliceVideoEpisode", contentApi.SliceVideoEpisode)
	}
}
