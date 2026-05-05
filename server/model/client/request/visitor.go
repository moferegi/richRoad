package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// HeartbeatRequest 心跳上报请求
type HeartbeatRequest struct {
	VisitorID     string `json:"visitorId" form:"visitorId"`
	PagePath      string `json:"pagePath" form:"pagePath"`
	Referer       string `json:"referer" form:"referer"`
	Platform      string `json:"platform" form:"platform"`
	EventCategory string `json:"eventCategory" form:"eventCategory"`
	EventAction   string `json:"eventAction" form:"eventAction"`
	EventLabel    string `json:"eventLabel" form:"eventLabel"`
	EventValue    int    `json:"eventValue" form:"eventValue"`
	EventExtra    string `json:"eventExtra" form:"eventExtra"`
	ScreenWidth   int    `json:"screenWidth" form:"screenWidth"`
	ScreenHeight  int    `json:"screenHeight" form:"screenHeight"`
	Language      string `json:"language" form:"language"`
	SessionID     string `json:"sessionId" form:"sessionId"`
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
