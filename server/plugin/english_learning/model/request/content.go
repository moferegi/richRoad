package request

import commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type EnglishCategorySearch struct {
	NeedVip *bool `json:"needVip" form:"needVip"`
	commonReq.PageInfo
}

type EnglishChapterSearch struct {
	CategoryID uint `json:"categoryId" form:"categoryId"`
	commonReq.PageInfo
}

type VideoCategorySearch struct {
	ShowHome *bool `json:"showHome" form:"showHome"`
	commonReq.PageInfo
}

type VideoSeriesSearch struct {
	CategoryID uint   `json:"categoryId" form:"categoryId"`
	ShowHome   *bool  `json:"showHome" form:"showHome"`
	Keyword    string `json:"keyword" form:"keyword"`
	commonReq.PageInfo
}

type VideoEpisodeSearch struct {
	SeriesID uint   `json:"seriesId" form:"seriesId"`
	Keyword  string `json:"keyword" form:"keyword"`
	commonReq.PageInfo
}
