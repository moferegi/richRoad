package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
)

type EnglishWordSearch struct {
	model.EnglishWord
	request.PageInfo
}

type CreateEnglishWordReq struct {
	Word        string `json:"word" binding:"required"`
	PhoneticUS  string `json:"phoneticUs"`
	PhoneticUK  string `json:"phoneticUk"`
	Explanation string `json:"explanation"` // JSON string for multiple languages
	ChapterIDs  []uint `json:"chapterIds"`  // 创建单词时自动关联的章节ID列表
}

type UpdateEnglishWordReq struct {
	ID          uint   `json:"ID" binding:"required"`
	Word        string `json:"word" binding:"required"`
	PhoneticUS  string `json:"phoneticUs"`
	PhoneticUK  string `json:"phoneticUk"`
	Explanation string `json:"explanation"`
	ChapterIDs  []uint `json:"chapterIds"` // 可选；传入时会重建章节绑定关系
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
	ChapterID uint `json:"chapterId" form:"chapterId"`
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
