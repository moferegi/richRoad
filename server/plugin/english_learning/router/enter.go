package router

type RouterGroup struct {
	UserLearningAssetRouter
	EnglishWordRouter
	VideoSubtitleRouter
	CheckinRouter
	ContentRouter
	UserDataRouter
	AdminRouter
	DiaryRouter
}

var RouterGroupApp = new(RouterGroup)
