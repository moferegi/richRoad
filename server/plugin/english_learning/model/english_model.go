package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// EnglishCategory 英语分类包
type EnglishCategory struct {
	global.GVA_MODEL
	Name    string  `json:"name" gorm:"type:json;comment:名称(支持多语言)"`
	Logo    string  `json:"logo" gorm:"comment:分类Logo"`
	Price   float64 `json:"price" gorm:"comment:购买价格(0表示免费)"`
	NeedVip *bool   `json:"needVip" gorm:"default:false;comment:是否需要开通会员"`
	Sort    int     `json:"sort" gorm:"default:0;comment:排序索引"`
}

// EnglishChapter 英语分类中的章节
type EnglishChapter struct {
	global.GVA_MODEL
	CategoryID uint   `json:"categoryId" gorm:"index;comment:对应分类ID"`
	Name       string `json:"name" gorm:"type:json;comment:章节名称(支持多语言)"`
	Sort       int    `json:"sort" gorm:"default:0;comment:排序索引"`
}

// EnglishWord 英语单词字典
type EnglishWord struct {
	global.GVA_MODEL
	Word        string                `json:"word" gorm:"uniqueIndex;comment:单词本体"`
	PhoneticUS  string                `json:"phoneticUs" gorm:"comment:美式音标"`
	PhoneticUK  string                `json:"phoneticUk" gorm:"comment:英式音标"`
	AudioUS     string                `json:"audioUs" gorm:"comment:美式发音音频外链"`
	AudioUK     string                `json:"audioUk" gorm:"comment:英式发音音频外链"`
	Explanation string                `json:"explanation" gorm:"type:json;comment:单词解释(多语言JSON)"`
	Sentences   []EnglishWordSentence `json:"sentences,omitempty" gorm:"-"`
}

// EnglishCategoryWord 英语分类/章节与单词的关联表 (多对多)
type EnglishCategoryWord struct {
	global.GVA_MODEL
	CategoryID uint `json:"categoryId" gorm:"index;comment:对应分类ID"`
	ChapterID  uint `json:"chapterId" gorm:"index;comment:对应章节ID"`
	WordID     uint `json:"wordId" gorm:"index;comment:对应单词ID"`
	Sort       int  `json:"sort" gorm:"default:0;comment:该单词在本章节的排序"`
}

// EnglishWordSentence 英语造句
type EnglishWordSentence struct {
	global.GVA_MODEL
	WordID    uint   `json:"wordId" gorm:"index;comment:对应单词ID"`
	Source    string `json:"source" gorm:"type:text;comment:原句(英文)"`
	Translate string `json:"translate" gorm:"type:json;comment:外文翻译(多语言JSON)"`
	AudioUS   string `json:"audioUs" gorm:"comment:美式发音音频外链"`
	AudioUK   string `json:"audioUk" gorm:"comment:英式发音音频外链"`
	VideoID   uint   `json:"videoId" gorm:"comment:绑定的相关视频片段(视频库单集ID)"`
	Sort      int    `json:"sort" gorm:"default:0;comment:造句展示排序"`
}

// EnglishWordErrorLog 错题本记录（用户-单词维度聚合）
type EnglishWordErrorLog struct {
	global.GVA_MODEL
	UserID           uint   `json:"userId" gorm:"index;uniqueIndex:idx_user_word_error;comment:用户ID"`
	WordID           uint   `json:"wordId" gorm:"index;uniqueIndex:idx_user_word_error;comment:单词ID"`
	CategoryID       uint   `json:"categoryId" gorm:"index;comment:分类ID"`
	ChapterID        uint   `json:"chapterId" gorm:"index;comment:章节ID"`
	WrongCount       int    `json:"wrongCount" gorm:"default:1;comment:累计打错次数"`
	LastWrongIndex   int    `json:"lastWrongIndex" gorm:"default:-1;comment:最近错误字母下标"`
	LastExpectedChar string `json:"lastExpectedChar" gorm:"size:16;comment:最近期望字母"`
	LastInputChar    string `json:"lastInputChar" gorm:"size:16;comment:最近实际输入字母"`
}
