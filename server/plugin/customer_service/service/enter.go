package service

// ServiceGroup 服务组
type ServiceGroup struct {
	ConversationService
	MessageService
	QuickReplyService
	AgentService
	BlacklistService
	ConfigService
}

// Service 全局服务实例
var Service = new(ServiceGroup)
