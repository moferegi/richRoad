package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// CsConversation 客服会话
type CsConversation struct {
	global.GVA_MODEL
	ClientUserID  uint       `json:"clientUserId" gorm:"column:client_user_id;index;not null;comment:客户端用户ID"`
	AgentUserID   *uint      `json:"agentUserId" gorm:"column:agent_user_id;index;comment:坐席系统用户ID"`
	Status        string     `json:"status" gorm:"column:status;default:pending;comment:状态:pending/active/closed"`
	QueuePosition int        `json:"queuePosition" gorm:"column:queue_position;default:0;comment:排队位置(0=不在队列)"`
	ClosedAt      *time.Time `json:"closedAt" gorm:"column:closed_at;comment:关闭时间"`
	ClosedBy      string     `json:"closedBy" gorm:"column:closed_by;comment:关闭方:user/agent/system"`
	Rating        *int       `json:"rating" gorm:"column:rating;comment:满意度评分1-5"`
	RatedAt       *time.Time `json:"ratedAt" gorm:"column:rated_at;comment:评分时间"`
	// 非持久化字段
	UnreadCount    int    `json:"unreadCount" gorm:"-"`
	LastMsg        string `json:"lastMsg" gorm:"-"`
	ClientNickname string `json:"clientNickname" gorm:"-"`
	AgentNickname  string `json:"agentNickname" gorm:"-"`
}

func (CsConversation) TableName() string { return "cs_conversations" }

// 会话状态常量
const (
	ConvStatusPending = "pending" // 待分配
	ConvStatusActive  = "active"  // 服务中
	ConvStatusClosed  = "closed"  // 已关闭
)

// 关闭方常量
const (
	ConvClosedByUser   = "user"
	ConvClosedByAgent  = "agent"
	ConvClosedBySystem = "system"
)
