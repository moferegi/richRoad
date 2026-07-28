package request

import commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type SaveWordProgressReq struct {
	CategoryID uint `json:"categoryId" binding:"required"`
	ChapterID  uint `json:"chapterId"`
	WordIndex  int  `json:"wordIndex"`
}

type CollectionReq struct {
	TargetType int  `json:"targetType" binding:"required"` // 1单词 2视频句子 3视频 4日记 5日记句子
	TargetID   uint `json:"targetId" binding:"required"`
}

type CollectionSearch struct {
	TargetType  int   `json:"targetType" form:"targetType"`
	TargetTypes []int `json:"targetTypes" form:"targetTypes"`
	commonReq.PageInfo
}

type WatchHistorySearch struct {
	commonReq.PageInfo
}
