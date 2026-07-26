package router

type RouterGroup struct {
	UserLearningAssetRouter
	EnglishWordRouter
	VideoSubtitleRouter
	CheckinRouter
	ContentRouter
	UserDataRouter
	AdminRouter
}

var RouterGroupApp = new(RouterGroup)
