package api

type ApiGroup struct {
	UserLearningAssetApi
	EnglishWordApi
	VideoSubtitleApi
	CheckinApi
	ContentApi
	UserDataApi
	AdminApi
}

var ApiGroupApp = new(ApiGroup)
