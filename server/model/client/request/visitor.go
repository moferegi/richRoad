package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// HeartbeatRequest 心跳上报请求
type HeartbeatRequest struct {
	VisitorID     string `json:"visitorId" form:"visitorId" binding:"max=64"`
	PagePath      string `json:"pagePath" form:"pagePath" binding:"max=255"`
	Referer       string `json:"referer" form:"referer" binding:"max=512"`
	Platform      string `json:"platform" form:"platform" binding:"max=32"`
	EventCategory string `json:"eventCategory" form:"eventCategory" binding:"max=64"`
	EventAction   string `json:"eventAction" form:"eventAction" binding:"max=64"`
	EventLabel    string `json:"eventLabel" form:"eventLabel" binding:"max=128"`
	EventValue    int    `json:"eventValue" form:"eventValue"`
	EventExtra    string `json:"eventExtra" form:"eventExtra" binding:"max=4096"`
	ScreenWidth   int    `json:"screenWidth" form:"screenWidth"`
	ScreenHeight  int    `json:"screenHeight" form:"screenHeight"`
	Language      string `json:"language" form:"language" binding:"max=16"`
	SessionID     string `json:"sessionId" form:"sessionId" binding:"max=64"`
}

// KefuGuideStatsSearch 客服引导统计查询条件
type KefuGuideStatsSearch struct {
	StartDate string `json:"startDate" form:"startDate"` // 开始日期
	EndDate   string `json:"endDate" form:"endDate"`     // 结束日期
	Platform  string `json:"platform" form:"platform"`   // 平台
}

// VisitorLogSearch 访客日志查询条件
type VisitorLogSearch struct {
	request.PageInfo
	StartDate string `json:"startDate" form:"startDate"` // 开始日期
	EndDate   string `json:"endDate" form:"endDate"`     // 结束日期
	VisitorID string `json:"visitorId" form:"visitorId"` // 访客ID
	IP        string `json:"ip" form:"ip"`               // IP地址
	Platform  string `json:"platform" form:"platform"`   // 平台
	OrderBy   string `json:"orderBy" form:"orderBy"`
	OrderDir  string `json:"orderDir" form:"orderDir"`
}

// VisitorSummarySearch 访客汇总查询条件
type VisitorSummarySearch struct {
	request.PageInfo
	StartDate string `json:"startDate" form:"startDate"` // 开始日期
	EndDate   string `json:"endDate" form:"endDate"`     // 结束日期
}
