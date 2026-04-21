package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// ConversationSearch 会话搜索参数
type ConversationSearch struct {
	request.PageInfo
	Status       string     `json:"status" form:"status"`
	AgentUserID  *uint      `json:"agentUserId" form:"agentUserId"`
	ClientUserID *uint      `json:"clientUserId" form:"clientUserId"`
	StartTime    *time.Time `json:"startTime" form:"startTime"`
	EndTime      *time.Time `json:"endTime" form:"endTime"`
}

// MessageSearch 消息历史搜索参数
type MessageSearch struct {
	request.PageInfo
	ConversationID uint `json:"conversationId" form:"conversationId" binding:"required"`
}

// SendMessageReq REST发送消息请求
type SendMessageReq struct {
	ConversationID uint   `json:"conversationId" binding:"required"`
	MsgType        string `json:"msgType" binding:"required"`
	Content        string `json:"content" binding:"required"`
	ClientMsgID    string `json:"clientMsgId"`
}

// RevokeMessageReq 撤回消息请求
type RevokeMessageReq struct {
	MessageID uint `json:"messageId" binding:"required"`
}

// RateConversationReq 评价会话请求
type RateConversationReq struct {
	ConversationID uint `json:"conversationId" binding:"required"`
	Rating         int  `json:"rating" binding:"required,min=1,max=5"`
}

// TransferConversationReq 转接会话请求
type TransferConversationReq struct {
	ConversationID    uint `json:"conversationId" binding:"required"`
	TargetAgentUserID uint `json:"targetAgentUserId" binding:"required"`
}

// CloseConversationReq 关闭会话请求
type CloseConversationReq struct {
	ConversationID uint   `json:"conversationId" binding:"required"`
	ClosedBy       string `json:"closedBy"` // user/agent/system
}

// UpdateAgentReq 更新坐席请求
type UpdateAgentReq struct {
	ID          uint   `json:"ID" binding:"required"`
	MaxSessions int    `json:"maxSessions" binding:"min=1,max=50"`
	IsEnabled   *bool  `json:"isEnabled"`
	Nickname    string `json:"nickname"`
}

// QuickReplySearch 快捷回复搜索
type QuickReplySearch struct {
	request.PageInfo
	Title string `json:"title" form:"title"`
}

// BlacklistSearch 黑名单搜索
type BlacklistSearch struct {
	request.PageInfo
}

// WSMessageFrame WebSocket 消息帧(客户端发往服务端)
type WSMessageFrame struct {
	Event string `json:"event"` // send_message / revoke / read / typing / ping
	// send_message 字段
	MsgType     string `json:"msgType"`
	Content     string `json:"content"`
	ClientMsgID string `json:"clientMsgId"`
	// revoke 字段
	MessageID uint `json:"messageId"`
}
