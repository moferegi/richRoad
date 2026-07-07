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
