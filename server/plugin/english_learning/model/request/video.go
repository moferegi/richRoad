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
