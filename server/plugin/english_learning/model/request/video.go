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
