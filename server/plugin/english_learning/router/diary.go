package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/api"
	"github.com/gin-gonic/gin"
)

type DiaryRouter struct{}

func (s *DiaryRouter) InitDiaryRouter(Router *gin.RouterGroup) {
	diaryRouter := Router.Group("diary")
	var diaryApi = api.ApiGroupApp.DiaryApi
	{
		diaryRouter.POST("createDiaryCategory", diaryApi.CreateDiaryCategory)
		diaryRouter.PUT("updateDiaryCategory", diaryApi.UpdateDiaryCategory)
		diaryRouter.DELETE("deleteDiaryCategory", diaryApi.DeleteDiaryCategory)
		diaryRouter.GET("findDiaryCategory", diaryApi.FindDiaryCategory)
		diaryRouter.GET("getDiaryCategoryList", diaryApi.GetDiaryCategoryList)

		diaryRouter.POST("createDiary", diaryApi.CreateDiary)
		diaryRouter.PUT("updateDiary", diaryApi.UpdateDiary)
		diaryRouter.DELETE("deleteDiary", diaryApi.DeleteDiary)
		diaryRouter.GET("findDiary", diaryApi.FindDiary)
		diaryRouter.GET("getDiaryList", diaryApi.GetDiaryList)
		diaryRouter.GET("getDiaryListByTag", diaryApi.GetDiaryListByTag)

		diaryRouter.GET("getSentenceList", diaryApi.GetDiarySentenceList)
		diaryRouter.GET("getSentenceListAll", diaryApi.GetDiarySentenceListAll)
		diaryRouter.POST("parseDiarySubtitle", diaryApi.ParseDiarySubtitle)
		diaryRouter.POST("scanKeywords", diaryApi.ScanDiaryKeywords)
		diaryRouter.POST("parseDiarySubtitleFiles", diaryApi.ParseDiarySubtitleFiles)
		diaryRouter.GET("getKeywords", diaryApi.GetDiaryKeywords)
		diaryRouter.PUT("rehighlightSentences", diaryApi.RehighlightDiarySentences)
		diaryRouter.PUT("updateSentenceList", diaryApi.UpdateDiarySentenceList)
	}
}
