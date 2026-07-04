package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// VideoCategory 视频分类
type VideoCategory struct {
	global.GVA_MODEL
	Name       string `json:"name" gorm:"type:json;comment:视频分类名称(支持多语言)"`
	StorageKey string `json:"storageKey" gorm:"size:64;default:'';comment:视频资源存储目录键(如liblib)"`
	Sort       int    `json:"sort" gorm:"default:0;comment:排序索引"`
}

// VideoSeries 视频集簇 (类似电视剧系列)
type VideoSeries struct {
	global.GVA_MODEL
	CategoryID uint    `json:"categoryId" gorm:"index;comment:所属视频分类ID"`
	Name       string  `json:"name" gorm:"type:json;comment:剧集名称(支持多语言)"`
	CoverID    uint    `json:"coverId" gorm:"comment:视频封面图(关联附件库)"`
	CoverUrl   string  `json:"coverUrl" gorm:"comment:视频封面图外链"`
	Price      float64 `json:"price" gorm:"comment:单买价格(0代表免费)"`
	NeedVip    *bool   `json:"needVip" gorm:"default:false;comment:是否需要开通VIP才可观看"`
	ViewCount  int     `json:"viewCount" gorm:"default:0;comment:外部展示总浏览量(次)"`
	UserCount  int     `json:"userCount" gorm:"default:0;comment:外部展示总观看人数(人)"`
}

// VideoEpisode 视频具体单集 (播放实体)
type VideoEpisode struct {
	global.GVA_MODEL
	SeriesID     uint   `json:"seriesId" gorm:"index;comment:所属剧集ID"`
	Name         string `json:"name" gorm:"type:json;comment:单集名称(支持多语言)"`
	VideoUrl     string `json:"videoUrl" gorm:"comment:外部视频链接"`
	TrialPercent int    `json:"trialPercent" gorm:"default:8;comment:无权限时最大可试看百分比(默认8%)"`
	HasFullAuth  bool   `json:"hasFullAuth" gorm:"-"`
	Sort         int    `json:"sort" gorm:"default:0;comment:集数排序(第几集)"`
}

// VideoSubtitle 视频外部字幕挂载配置
type VideoSubtitle struct {
	global.GVA_MODEL
	EpisodeID   uint   `json:"episodeId" gorm:"index;comment:所属单集ID"`
	Language    string `json:"language" gorm:"comment:字幕所属语种(如 en, zh)"`
	Format      string `json:"format" gorm:"comment:字幕格式(srt/vtt等)"`
	SubtitleUrl string `json:"subtitleUrl" gorm:"comment:外部字幕链接"`
}

// VideoSentence 视频字幕文本入库 (用于词汇碰对、高亮和允许收藏)
type VideoSentence struct {
	global.GVA_MODEL
	EpisodeID uint    `json:"episodeId" gorm:"index;comment:所属单集ID"`
	StartTime float64 `json:"startTime" gorm:"comment:字幕开始秒数"`
	EndTime   float64 `json:"endTime" gorm:"comment:字幕结束秒数"`
	English   string  `json:"english" gorm:"type:text;comment:原文字幕(包含<w>变色等预先匹配标记)"`
	Translate string  `json:"translate" gorm:"type:json;comment:对应系统拥有的全语种翻译JSON存储"`
}
