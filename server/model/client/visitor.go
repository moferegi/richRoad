package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// VisitorLog 访客日志 结构体
type VisitorLog struct {
	global.GVA_MODEL
	VisitorID     string `json:"visitorId" form:"visitorId" gorm:"column:visitor_id;size:64;index;comment:访客指纹ID;"`
	UserID        uint   `json:"userId" form:"userId" gorm:"column:user_id;default:0;index;comment:注册用户ID(未登录为0);"`
	IP            string `json:"ip" form:"ip" gorm:"column:ip;size:45;index;comment:IP地址;"`
	UserAgent     string `json:"userAgent" form:"userAgent" gorm:"column:user_agent;size:512;comment:浏览器UA;"`
	Platform      string `json:"platform" form:"platform" gorm:"column:platform;size:32;index;comment:平台(ios/android/h5);"`
	PagePath      string `json:"pagePath" form:"pagePath" gorm:"column:page_path;size:255;comment:访问页面路径;"`
	Referer       string `json:"referer" form:"referer" gorm:"column:referer;size:512;comment:来源页面;"`
	EventCategory string `json:"eventCategory" form:"eventCategory" gorm:"column:event_category;size:64;index;comment:事件分类;"`
	EventAction   string `json:"eventAction" form:"eventAction" gorm:"column:event_action;size:64;index;comment:事件动作;"`
	EventLabel    string `json:"eventLabel" form:"eventLabel" gorm:"column:event_label;size:128;comment:事件标签;"`
	EventValue    int    `json:"eventValue" form:"eventValue" gorm:"column:event_value;default:0;comment:事件数值;"`
	EventExtra    string `json:"eventExtra" form:"eventExtra" gorm:"column:event_extra;type:text;comment:事件附加信息(JSON);"`
	ScreenWidth   int    `json:"screenWidth" form:"screenWidth" gorm:"column:screen_width;comment:屏幕宽度;"`
	ScreenHeight  int    `json:"screenHeight" form:"screenHeight" gorm:"column:screen_height;comment:屏幕高度;"`
	Language      string `json:"language" form:"language" gorm:"column:language;size:16;comment:浏览器语言;"`
	SessionID     string `json:"sessionId" form:"sessionId" gorm:"column:session_id;size:64;index;comment:会话ID;"`
	Location      string `json:"location" form:"location" gorm:"column:location;size:128;comment:IP归属地;"`
	CreatedDate   string `json:"-" gorm:"column:created_date;size:10;index;comment:创建日期(冗余字段加速查询);"`
}

// TableName 访客日志 自定义表名
func (VisitorLog) TableName() string {
	return "client_visitor_log"
}

// VisitorSummary 访客汇总统计 结构体
type VisitorSummary struct {
	global.GVA_MODEL
	Date           string `json:"date" form:"date" gorm:"column:date;size:10;uniqueIndex;comment:日期(YYYY-MM-DD);"`
	PV             int64  `json:"pv" form:"pv" gorm:"column:pv;default:0;comment:页面浏览量;"`
	UV             int64  `json:"uv" form:"uv" gorm:"column:uv;default:0;comment:独立访客数;"`
	NewVisitor     int64  `json:"newVisitor" form:"newVisitor" gorm:"column:new_visitor;default:0;comment:新访客数;"`
	RegisteredUser int64  `json:"registeredUser" form:"registeredUser" gorm:"column:registered_user;default:0;comment:已登录用户数;"`
}

// TableName 访客汇总 自定义表名
func (VisitorSummary) TableName() string {
	return "client_visitor_summary"
}
