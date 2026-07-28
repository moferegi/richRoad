package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// DiaryCategory 日记分类
type DiaryCategory struct {
	global.GVA_MODEL
	Name       string `json:"name" gorm:"type:json;comment:日记分类名称(支持多语言)"`
	StorageKey string `json:"storageKey" gorm:"size:64;default:'';comment:日记资源存储目录键"`
	ShowHome   *bool  `json:"showHome" gorm:"default:true;comment:是否在首页展示"`
	Sort       int    `json:"sort" gorm:"default:0;comment:排序索引"`
}

// Diary 日记 (播放实体)
type Diary struct {
	global.GVA_MODEL
	CategoryID           uint    `json:"categoryId" gorm:"index;comment:所属日记分类ID"`
	Name                 string  `json:"name" gorm:"type:json;comment:日记名称(支持多语言)"`
	AudioUs              string  `json:"audioUs" gorm:"type:text;comment:美式语音URL"`
	AudioUk              string  `json:"audioUk" gorm:"type:text;comment:英式语音URL"`
	Duration             float64 `json:"duration" gorm:"default:0;comment:音频时长(秒)"`
	ImageI18n            string  `json:"imageI18n" gorm:"type:json;comment:多语言日记图片"`
	ShowTranslate        *bool   `json:"showTranslate" gorm:"default:true;comment:是否显示对应语言翻译字幕"`
	ShowEnglish          *bool   `json:"showEnglish" gorm:"default:true;comment:是否显示英文字幕"`
	TrialPercent         int     `json:"trialPercent" gorm:"default:8;comment:无权限时最大可试看百分比(默认8%)"`
	NeedVip              *bool   `json:"needVip" gorm:"default:false;comment:是否需要开通VIP才可完整观看"`
	EnglishSubtitleUrl   string  `json:"englishSubtitleUrl" gorm:"type:text;comment:美式英文字幕文件URL"`
	EnglishSubtitleUrlUk string  `json:"englishSubtitleUrlUk" gorm:"type:text;comment:英式英文字幕文件URL(仅时间轴不同)"`
	TranslationSubtitleUrls string `json:"translationSubtitleUrls" gorm:"type:json;comment:多语言翻译字幕文件URL映射JSON"`
	HasFullAuth          bool    `json:"hasFullAuth" gorm:"-"`
	Sort                 int     `json:"sort" gorm:"default:0;comment:排序索引"`
}

// DiarySentence 日记字幕句子 (用于收藏和展示)
type DiarySentence struct {
	global.GVA_MODEL
	DiaryID     uint    `json:"diaryId" gorm:"index;comment:所属日记ID"`
	StartTime   float64 `json:"startTime" gorm:"comment:美式字幕开始秒数"`
	EndTime     float64 `json:"endTime" gorm:"comment:美式字幕结束秒数"`
	StartTimeUk float64 `json:"startTimeUk" gorm:"default:0;comment:英式字幕开始秒数(0表示无UK时间轴，回退用US)"`
	EndTimeUk   float64 `json:"endTimeUk" gorm:"default:0;comment:英式字幕结束秒数(0表示无UK时间轴，回退用US)"`
	English     string  `json:"english" gorm:"type:text;comment:原文字幕(包含<w>变色等预先匹配标记)"`
	Translate   string  `json:"translate" gorm:"type:json;comment:对应系统拥有的全语种翻译JSON存储"`
	Locked      bool    `json:"locked" gorm:"-"` // 仅传输用：试看限制下该句是否锁定
}
