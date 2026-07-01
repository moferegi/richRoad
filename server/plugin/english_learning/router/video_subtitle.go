package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/api"
	"github.com/gin-gonic/gin"
)

type VideoSubtitleRouter struct{}

func (s *VideoSubtitleRouter) InitVideoSubtitleRouter(Router *gin.RouterGroup) {
	videoRouter := Router.Group("video")
	var videoSubtitleApi = api.ApiGroupApp.VideoSubtitleApi
	{
		videoRouter.POST("parseSubtitle", videoSubtitleApi.ParseSubtitle) // 解析字幕生成高亮
		videoRouter.GET("getSentenceList", videoSubtitleApi.GetSentenceList)
	}
}
