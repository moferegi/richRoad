package request

type SubtitleItem struct {
	StartTime float64 `json:"startTime"`
	EndTime   float64 `json:"endTime"`
	English   string  `json:"english"`
	Translate string  `json:"translate"` // 多语种JSON结构，例如 {"zh":"你好","mn":"..."}
}

type ParseSubtitleReq struct {
	EpisodeID uint           `json:"episodeId" binding:"required"`
	Subtitles []SubtitleItem `json:"subtitles" binding:"required"`
}

type SubtitleFileItem struct {
	Language    string `json:"language" binding:"required"`
	SubtitleURL string `json:"subtitleUrl" binding:"required"`
}

type ParseSubtitleFilesReq struct {
	EpisodeID           uint               `json:"episodeId" binding:"required"`
	EnglishSubtitleURL  string             `json:"englishSubtitleUrl" binding:"required"`
	TranslationSubtitle []SubtitleFileItem `json:"translationSubtitle"`
	KeywordIDs          []uint             `json:"keywordIds"` // 管理员多选确认的重点单词ID列表
}

// ScanKeywordsReq 扫描字幕提取关键词请求
type ScanKeywordsReq struct {
	EpisodeID          uint   `json:"episodeId" binding:"required"`
	EnglishSubtitleURL string `json:"englishSubtitleUrl" binding:"required"`
}

// KeywordItem 字幕中提取的单词项
type KeywordItem struct {
	Word    string `json:"word"`    // 单词原形
	WordID  uint   `json:"wordId"`  // 单词库中的ID，0表示未匹配
	Matched bool   `json:"matched"` // 是否在单词库中存在
	Count   int    `json:"count"`   // 在字幕中出现的次数
}

type UpdateVideoSentenceItem struct {
	ID        uint    `json:"id" binding:"required"`
	English   string  `json:"english"`
	Translate string  `json:"translate"`
	StartTime float64 `json:"startTime"`
	EndTime   float64 `json:"endTime"`
}

type UpdateVideoSentenceListReq struct {
	EpisodeID uint                      `json:"episodeId" binding:"required"`
	Sentences []UpdateVideoSentenceItem `json:"sentences" binding:"required"`
}

// RehighlightSentencesReq 对已有字幕重新高亮请求
type RehighlightSentencesReq struct {
	EpisodeID  uint   `json:"episodeId" binding:"required"`
	KeywordIDs []uint `json:"keywordIds"` // 重新选定的重点单词ID列表
}

// CreateVideoEpisodeReq 创建视频单集请求（含标签）
type CreateVideoEpisodeReq struct {
	SeriesID     uint    `json:"seriesId"`
	Name         string  `json:"name"`
	VideoUrl     string  `json:"videoUrl"`
	Duration     float64 `json:"duration"`
	VideoType    string  `json:"videoType"`
	TrialPercent int     `json:"trialPercent"`
	Sort         int     `json:"sort"`
	TagIds       []uint  `json:"tagIds"`
}

// UpdateVideoEpisodeReq 更新视频单集请求（含标签）
type UpdateVideoEpisodeReq struct {
	ID           uint    `json:"ID" binding:"required"`
	SeriesID     uint    `json:"seriesId"`
	Name         string  `json:"name"`
	VideoUrl     string  `json:"videoUrl"`
	Duration     float64 `json:"duration"`
	VideoType    string  `json:"videoType"`
	TrialPercent int     `json:"trialPercent"`
	Sort         int     `json:"sort"`
	TagIds       []uint  `json:"tagIds"`
}

// VideoEpisodeTagSearch 按分类+标签筛选单集请求
type VideoEpisodeTagSearch struct {
	CategoryID uint   `json:"categoryId" form:"categoryId"`
	TagIds     []uint `json:"tagIds" form:"tagIds"`
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"pageSize" form:"pageSize"`
}
