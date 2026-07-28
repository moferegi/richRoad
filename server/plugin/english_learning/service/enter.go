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
	AdminService
	DiaryService
}

var ServiceGroupApp = new(ServiceGroup)
