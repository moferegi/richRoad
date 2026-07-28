package request

import commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

// DiaryCategorySearch 日记分类搜索
type DiaryCategorySearch struct {
	ShowHome *bool `json:"showHome" form:"showHome"`
	commonReq.PageInfo
}

// DiarySearch 日记搜索
type DiarySearch struct {
	CategoryID uint   `json:"categoryId" form:"categoryId"`
	Keyword    string `json:"keyword" form:"keyword"`
	commonReq.PageInfo
}

// DiaryTagFilterSearch 按分类+标签筛选日记
type DiaryTagFilterSearch struct {
	CategoryID uint   `json:"categoryId" form:"categoryId"`
	TagIds     []uint `json:"tagIds" form:"tagIds"`
	commonReq.PageInfo
}

// CreateDiaryReq 创建日记请求
type CreateDiaryReq struct {
	CategoryID           uint    `json:"categoryId"`
	Name                 string  `json:"name"`
	AudioUs              string  `json:"audioUs"`
	AudioUk              string  `json:"audioUk"`
	Duration             float64 `json:"duration"`
	ImageI18n            string  `json:"imageI18n"`
	ShowTranslate        *bool   `json:"showTranslate"`
	ShowEnglish          *bool   `json:"showEnglish"`
	TrialPercent         int     `json:"trialPercent"`
	NeedVip              *bool   `json:"needVip"`
	EnglishSubtitleUrl   string  `json:"englishSubtitleUrl"`
	EnglishSubtitleUrlUk string  `json:"englishSubtitleUrlUk"`
	Sort                 int     `json:"sort"`
	TagIds               []uint  `json:"tagIds"`
}

// UpdateDiaryReq 更新日记请求
type UpdateDiaryReq struct {
	ID                   uint    `json:"ID"`
	CategoryID           uint    `json:"categoryId"`
	Name                 string  `json:"name"`
	AudioUs              string  `json:"audioUs"`
	AudioUk              string  `json:"audioUk"`
	Duration             float64 `json:"duration"`
	ImageI18n            string  `json:"imageI18n"`
	ShowTranslate        *bool   `json:"showTranslate"`
	ShowEnglish          *bool   `json:"showEnglish"`
	TrialPercent         int     `json:"trialPercent"`
	NeedVip              *bool   `json:"needVip"`
	EnglishSubtitleUrl   string  `json:"englishSubtitleUrl"`
	EnglishSubtitleUrlUk string  `json:"englishSubtitleUrlUk"`
	Sort                 int     `json:"sort"`
	TagIds               []uint  `json:"tagIds"`
}

// DiarySubtitleItem 日记字幕条目
type DiarySubtitleItem struct {
	StartTime float64 `json:"startTime"`
	EndTime   float64 `json:"endTime"`
	English   string  `json:"english"`
	Translate string  `json:"translate"`
}

// ParseDiarySubtitleReq 解析日记字幕请求（JSON手动输入模式）
type ParseDiarySubtitleReq struct {
	DiaryID   uint                `json:"diaryId"`
	Subtitles []DiarySubtitleItem `json:"subtitles"`
}

// ScanDiaryKeywordsReq 扫描日记字幕关键词请求
type ScanDiaryKeywordsReq struct {
	DiaryID            uint   `json:"diaryId" binding:"required"`
	EnglishSubtitleURL string `json:"englishSubtitleUrl" binding:"required"`
}

// DiarySubtitleFileItem 日记字幕文件条目
type DiarySubtitleFileItem struct {
	Language    string `json:"language" binding:"required"`
	SubtitleURL string `json:"subtitleUrl" binding:"required"`
}

// ParseDiarySubtitleFilesReq 解析日记字幕文件请求（文件上传模式）
type ParseDiarySubtitleFilesReq struct {
	DiaryID              uint                    `json:"diaryId" binding:"required"`
	EnglishSubtitleURL   string                  `json:"englishSubtitleUrl" binding:"required"`
	EnglishSubtitleURLUk string                  `json:"englishSubtitleUrlUk"`
	TranslationSubtitle  []DiarySubtitleFileItem `json:"translationSubtitle"`
	KeywordIDs           []uint                  `json:"keywordIds"`
}

// DiarySentenceSearch 日记字幕句子搜索
type DiarySentenceSearch struct {
	DiaryID uint `json:"diaryId" form:"diaryId"`
	commonReq.PageInfo
}

// DiaryKeywordItem 字幕中提取的单词项
type DiaryKeywordItem struct {
	Word    string `json:"word"`
	WordID  uint   `json:"wordId"`
	Matched bool   `json:"matched"`
	Count   int    `json:"count"`
}

// RehighlightDiarySentencesReq 对已有日记字幕重新高亮请求
type RehighlightDiarySentencesReq struct {
	DiaryID    uint   `json:"diaryId" binding:"required"`
	KeywordIDs []uint `json:"keywordIds"`
}

// UpdateDiarySentenceItem 更新日记字幕单条记录
type UpdateDiarySentenceItem struct {
	ID        uint    `json:"id"`
	StartTime float64 `json:"startTime"`
	EndTime   float64 `json:"endTime"`
	English   string  `json:"english"`
	Translate string  `json:"translate"`
}

// UpdateDiarySentenceListReq 批量更新日记字幕句子请求
type UpdateDiarySentenceListReq struct {
	DiaryID   uint                      `json:"diaryId" binding:"required"`
	Sentences []UpdateDiarySentenceItem `json:"sentences" binding:"required"`
}
