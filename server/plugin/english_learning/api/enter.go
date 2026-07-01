package api

type ApiGroup struct {
	UserLearningAssetApi
	EnglishWordApi
	VideoSubtitleApi
	CheckinApi
	ContentApi
	UserDataApi
}

var ApiGroupApp = new(ApiGroup)
