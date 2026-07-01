package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/api"
	"github.com/gin-gonic/gin"
)

type UserDataRouter struct{}

func (s *UserDataRouter) InitUserDataRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("userData")
	var dataApi = api.ApiGroupApp.UserDataApi
	{
		dataRouter.POST("saveWordProgress", dataApi.SaveWordProgress)
		dataRouter.GET("getWordProgress", dataApi.GetWordProgress)
		dataRouter.GET("getWatchProgress", dataApi.GetWatchProgress)
		dataRouter.GET("getSeriesWatchProgressList", dataApi.GetSeriesWatchProgressList)
		dataRouter.GET("getWatchHistoryList", dataApi.GetWatchHistoryList)

		dataRouter.POST("collect", dataApi.Collect)
		dataRouter.DELETE("uncollect", dataApi.Uncollect)
		dataRouter.GET("getCollectionList", dataApi.GetCollectionList)
		dataRouter.GET("getCollectionDetailList", dataApi.GetCollectionDetailList)

		dataRouter.GET("getAsset", dataApi.GetAsset)
	}
}
