package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
)

type EnglishWordSearch struct {
	model.EnglishWord
	request.PageInfo
}

type EnglishWordSentenceReq struct {
	Source    string `json:"source"`
	Translate string `json:"translate"`
	AudioUS   string `json:"audioUs"`
	AudioUK   string `json:"audioUk"`
	VideoID   uint   `json:"videoId"`
	Sort      int    `json:"sort"`
}

type CreateEnglishWordReq struct {
	Word         string                   `json:"word" binding:"required"`
	PhoneticUS   string                   `json:"phoneticUs"`
	PhoneticUK   string                   `json:"phoneticUk"`
	PartOfSpeech string                   `json:"partOfSpeech"`
	AudioUS      string                   `json:"audioUs"`
	AudioUK      string                   `json:"audioUk"`
	Explanation  string                   `json:"explanation"` // JSON string for multiple languages
	CategoryIDs  []uint                   `json:"categoryIds"` // 单词可归属多个分类
	ChapterIDs   []uint                   `json:"chapterIds"`  // 创建单词时自动关联的章节ID列表
	Sentences    []EnglishWordSentenceReq `json:"sentences"`
}

type UpdateEnglishWordReq struct {
	ID           uint                     `json:"ID" binding:"required"`
	Word         string                   `json:"word" binding:"required"`
	PhoneticUS   string                   `json:"phoneticUs"`
	PhoneticUK   string                   `json:"phoneticUk"`
	PartOfSpeech string                   `json:"partOfSpeech"`
	AudioUS      string                   `json:"audioUs"`
	AudioUK      string                   `json:"audioUk"`
	Explanation  string                   `json:"explanation"`
	CategoryIDs  []uint                   `json:"categoryIds"` // 可选；传入时会重建分类绑定关系
	ChapterIDs   []uint                   `json:"chapterIds"`  // 可选；传入时会重建章节绑定关系
	Sentences    []EnglishWordSentenceReq `json:"sentences"`
}

type RegenerateWordAudioReq struct {
	ID           uint  `json:"ID" binding:"required"`
	RegenerateUS *bool `json:"regenerateUs"`
	RegenerateUK *bool `json:"regenerateUk"`
}

type TTSPreflightReq struct {
	Word    string `json:"word"`
	CheckUS *bool  `json:"checkUs"`
	CheckUK *bool  `json:"checkUk"`
}

type WordListSearch struct {
	CategoryID uint   `json:"categoryId" form:"categoryId"`
	ChapterID  uint   `json:"chapterId" form:"chapterId"`
	Keyword    string `json:"keyword" form:"keyword"`
	request.PageInfo
}

type ReportWordErrorReq struct {
	WordID       uint   `json:"wordId" binding:"required"`
	CategoryID   uint   `json:"categoryId"`
	ChapterID    uint   `json:"chapterId"`
	WrongIndex   int    `json:"wrongIndex"`
	ExpectedChar string `json:"expectedChar"`
	InputChar    string `json:"inputChar"`
}

type WordErrorLogSearch struct {
	CategoryID uint `json:"categoryId" form:"categoryId"`
	ChapterID  uint `json:"chapterId" form:"chapterId"`
	request.PageInfo
}

type DeleteWordErrorLogReq struct {
	WordID uint `json:"wordId" binding:"required"`
}

type UpsertWordFromSQLReq struct {
	Word          string `json:"word" binding:"required"`
	TranslateZh   string `json:"translateZh"`
	CategoryID    uint   `json:"categoryId" binding:"required"`
	ChapterID     uint   `json:"chapterId"`
	GenerateAudio bool   `json:"generateAudio"`
}

type BatchFillWordFromDictionaryReq struct {
	CategoryID           uint     `json:"categoryId" binding:"required"`
	FillPhonetic         bool     `json:"fillPhonetic"`
	FillPartOfSpeech     bool     `json:"fillPartOfSpeech"`
	FillExplanation      bool     `json:"fillExplanation"`
	FillSentences        bool     `json:"fillSentences"`
	OverwriteExisting    bool     `json:"overwriteExisting"`
	TranslateExplanation bool     `json:"translateExplanation"`
	TranslateSentences   bool     `json:"translateSentences"`
	TargetLangs          []string `json:"targetLangs"`
	TranslateService     string   `json:"translateService"`
	TranslateURL         string   `json:"translateUrl"`
	Limit                int      `json:"limit"`
}

// BatchCreateWordsReq 批量创建单词（仅入库单词本体，不生成TTS，用于字幕全部单词一键入库）
type BatchCreateWordsReq struct {
	Words []string `json:"words" binding:"required"`
}

// BatchCreateWordsItem 批量创建单词的返回项
type BatchCreateWordsItem struct {
	Word   string `json:"word"`
	WordID uint   `json:"wordId"`
	Newly  bool   `json:"newly"` // 是否本次新创建
}
