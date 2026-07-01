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
	commonReq.PageInfo
}

type VideoSeriesSearch struct {
	CategoryID uint `json:"categoryId" form:"categoryId"`
	commonReq.PageInfo
}

type VideoEpisodeSearch struct {
	SeriesID uint `json:"seriesId" form:"seriesId"`
	commonReq.PageInfo
}
