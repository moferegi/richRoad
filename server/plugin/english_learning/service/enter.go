package service

type ServiceGroup struct {
	LearningAuthzService
	UserLearningAssetService
	EnglishWordService
	VideoSubtitleService
	CheckinService
	ContentService
	UserDataService
}

var ServiceGroupApp = new(ServiceGroup)
