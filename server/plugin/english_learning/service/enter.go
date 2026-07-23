package service

type ServiceGroup struct {
	LearningAuthzService
	UserLearningAssetService
	EnglishWordService
	VideoSubtitleService
	CheckinService
	ContentService
	UserDataService
	VideoSliceService
}

var ServiceGroupApp = new(ServiceGroup)
