package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// CsMessage 客服消息
type CsMessage struct {
	global.GVA_MODEL
	ConversationID uint       `json:"conversationId" gorm:"column:conversation_id;index;not null;comment:会话ID"`
	SenderType     string     `json:"senderType" gorm:"column:sender_type;comment:发送方:user/agent/system"`
	SenderID       uint       `json:"senderId" gorm:"column:sender_id;comment:发送方ID"`
	MsgType        string     `json:"msgType" gorm:"column:msg_type;default:text;comment:消息类型:text/image/emoji/system"`
	Content        string     `json:"content" gorm:"column:content;type:text;comment:消息内容"`
	ClientMsgID    string     `json:"clientMsgId" gorm:"column:client_msg_id;index;size:64;comment:客户端消息ID(去重)"`
	Revoked        bool       `json:"revoked" gorm:"column:revoked;default:false;comment:是否已撤回"`
	RevokedAt      *time.Time `json:"revokedAt" gorm:"column:revoked_at;comment:撤回时间"`
	ReadAt         *time.Time `json:"readAt" gorm:"column:read_at;comment:已读时间"`
	// 非持久化字段
	SenderNickname string `json:"senderNickname" gorm:"-"`
	SenderAvatar   string `json:"senderAvatar" gorm:"-"`
}

func (CsMessage) TableName() string { return "cs_messages" }

// 消息类型常量
const (
	MsgTypeText   = "text"
	MsgTypeImage  = "image"
	MsgTypeEmoji  = "emoji"
	MsgTypeSystem = "system"
)

// 发送方类型常量
const (
	SenderTypeUser   = "user"
	SenderTypeAgent  = "agent"
	SenderTypeSystem = "system"
)
