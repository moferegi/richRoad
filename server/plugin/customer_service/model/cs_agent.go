package model

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// CsAgent 客服坐席配置
type CsAgent struct {
	global.GVA_MODEL
	UserID      uint   `json:"userId" gorm:"column:user_id;uniqueIndex;not null;comment:关联系统用户ID"`
	MaxSessions int    `json:"maxSessions" gorm:"column:max_sessions;default:5;comment:最大并发接待会话数"`
	IsEnabled   bool   `json:"isEnabled" gorm:"column:is_enabled;default:true;comment:是否启用"`
	Nickname    string `json:"nickname" gorm:"column:nickname;size:50;comment:客服昵称(覆盖系统昵称)"`
	// 非持久化字段
	OnlineStatus   string `json:"onlineStatus" gorm:"-"` // online/offline (实时判断)
	ActiveSessions int    `json:"activeSessions" gorm:"-"`
	SysNickname    string `json:"sysNickname" gorm:"-"`
	Avatar         string `json:"avatar" gorm:"-"`
}

func (CsAgent) TableName() string { return "cs_agents" }
