package api

type ApiGroup struct {
	UserLearningAssetApi
	EnglishWordApi
	VideoSubtitleApi
	CheckinApi
	ContentApi
	UserDataApi
	AdminApi
	DiaryApi
}

var ApiGroupApp = new(ApiGroup)
