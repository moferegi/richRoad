package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// AdPositionSearch 广告位置分页查询
type AdPositionSearch struct {
	Name        string `json:"name" form:"name"`
	PositionKey string `json:"positionKey" form:"positionKey"`
	IsEnabled   *bool  `json:"isEnabled" form:"isEnabled"`
	request.PageInfo
}

// AdVideoSearch 广告视频分页查询
type AdVideoSearch struct {
	PositionID uint   `json:"positionId" form:"positionId"`
	Title      string `json:"title" form:"title"`
	MediaType  string `json:"mediaType" form:"mediaType"`
	IsEnabled  *bool  `json:"isEnabled" form:"isEnabled"`
	request.PageInfo
}

// AdWatchRecordSearch 观看记录分页查询
type AdWatchRecordSearch struct {
	UserID     uint `json:"userId" form:"userId"`
	AdVideoID  uint `json:"adVideoId" form:"adVideoId"`
	PositionID uint `json:"positionId" form:"positionId"`
	request.PageInfo
}

// GetAdByPositionReq Uni端根据位置获取广告请求
type GetAdByPositionReq struct {
	PositionKey string `json:"positionKey" form:"positionKey" binding:"required"`
}

// ReportWatchReq Uni端上报观看记录请求
type ReportWatchReq struct {
	AdVideoID      uint `json:"adVideoId" binding:"required"`
	PositionID     uint `json:"positionId" binding:"required"`
	WatchedSeconds int  `json:"watchedSeconds" binding:"required"`
}