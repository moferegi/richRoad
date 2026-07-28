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
		videoRouter.POST("parseSubtitle", videoSubtitleApi.ParseSubtitle)           // 解析字幕生成高亮
		videoRouter.POST("scanKeywords", videoSubtitleApi.ScanKeywords)             // 扫描字幕提取关键词
		videoRouter.POST("parseSubtitleFiles", videoSubtitleApi.ParseSubtitleFiles) // 上传字幕文件解析（支持keywordIDs）
		videoRouter.GET("getSentenceList", videoSubtitleApi.GetSentenceList)
		videoRouter.GET("getEpisodeKeywords", videoSubtitleApi.GetEpisodeKeywords)  // 获取单集已高亮重点单词
		videoRouter.GET("getEpisodeSubtitles", videoSubtitleApi.GetEpisodeSubtitles) // 获取单集已上传的字幕文件记录
		videoRouter.PUT("rehighlightSentences", videoSubtitleApi.RehighlightSentences) // 重新高亮已有字幕
		videoRouter.PUT("updateSentenceList", videoSubtitleApi.UpdateSentenceList)
	}
}
