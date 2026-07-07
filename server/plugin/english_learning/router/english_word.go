package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/api"
	"github.com/gin-gonic/gin"
)

type EnglishWordRouter struct{}

func (s *EnglishWordRouter) InitEnglishWordRouter(Router *gin.RouterGroup) {
	wordRouter := Router.Group("word")
	var wordApi = api.ApiGroupApp.EnglishWordApi
	{
		wordRouter.POST("create", wordApi.CreateEnglishWord) // 创建单词
		wordRouter.PUT("update", wordApi.UpdateEnglishWord)
		wordRouter.DELETE("delete", wordApi.DeleteEnglishWord)
		wordRouter.POST("regenerateAudio", wordApi.RegenerateWordAudio)
		wordRouter.POST("preflightTTS", wordApi.PreflightTTS)
		wordRouter.GET("findWord", wordApi.FindEnglishWord)
		wordRouter.GET("getWordList", wordApi.GetWordList)
		wordRouter.POST("reportError", wordApi.ReportWordError)
		wordRouter.GET("getErrorLogList", wordApi.GetWordErrorLogList)
		wordRouter.DELETE("deleteErrorLog", wordApi.DeleteWordErrorLog)
		wordRouter.POST("upsertSqlWord", wordApi.UpsertWordFromSQL)
		wordRouter.POST("batchFillFromDictionary", wordApi.BatchFillWordFromDictionary)
	}
}
