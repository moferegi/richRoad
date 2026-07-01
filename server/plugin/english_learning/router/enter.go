package router

type RouterGroup struct {
	UserLearningAssetRouter
	EnglishWordRouter
	VideoSubtitleRouter
	CheckinRouter
	ContentRouter
	UserDataRouter
}

var RouterGroupApp = new(RouterGroup)
