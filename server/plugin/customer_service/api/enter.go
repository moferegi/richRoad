package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/service"
)

// ApiGroup API 组
type ApiGroup struct {
	WsApi
	ConversationApi
	MessageApi
	AgentApi
	QuickReplyApi
	BlacklistApi
	ConfigApi
}

// Api 全局 API 实例
var Api = new(ApiGroup)

var (
	conversationService = service.Service.ConversationService
	messageService      = service.Service.MessageService
	agentService        = service.Service.AgentService
	quickReplyService   = service.Service.QuickReplyService
	blacklistService    = service.Service.BlacklistService
)
