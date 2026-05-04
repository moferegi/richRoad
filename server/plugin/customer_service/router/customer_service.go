package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/api"
	"github.com/gin-gonic/gin"
)

type CustomerServiceRouter struct{}

// InitCustomerServiceRouter 初始化客服路由
func (r *CustomerServiceRouter) InitCustomerServiceRouter(public, private *gin.RouterGroup) {
	a := api.Api

	// ---- 公开接口（无需登录） ----
	public.GET("/cs/config/get", a.GetCsConfig)

	// ---- WebSocket（无 Casbin，仅 JWT 验证）----
	public.GET("/cs/ws", a.UserWS)
	public.GET("/cs/wsAgent", a.AgentWS)

	// ---- 客户端用户接口（需登录） ----
	csGroup := private.Group("/cs")
	{
		csGroup.POST("/conversation/getOrCreate", a.GetOrCreateConversation)
		csGroup.POST("/conversation/rate", a.RateConversation)
		csGroup.GET("/message/history", a.GetMessageHistory)
		csGroup.POST("/message/upload", a.UploadImage)
		csGroup.POST("/message/revoke", a.RevokeMessage)
		csGroup.GET("/quickReply/all", a.GetAllQuickReplies)
	}

	// ---- 坐席接口（需登录） ----
	agentGroup := private.Group("/cs/agent")
	{
		agentGroup.GET("/conversation/list", a.GetAgentConversationList)
		agentGroup.POST("/conversation/accept", a.AcceptConversation)
		agentGroup.POST("/conversation/close", a.CloseConversation)
		agentGroup.POST("/conversation/transfer", a.TransferConversation)
		agentGroup.POST("/message/send", a.SendMessage)
		agentGroup.POST("/message/upload", a.UploadImage)
		agentGroup.GET("/quickReply/list", a.GetQuickReplyList)
	}

	// ---- 管理员接口（需登录） ----
	adminGroup := private.Group("/cs/admin")
	{
		adminGroup.GET("/conversation/list", a.GetConversationList)
		adminGroup.POST("/conversation/assign", a.AssignConversation)
		adminGroup.GET("/agent/list", a.GetAgentList)
		adminGroup.POST("/agent/create", a.CreateAgent)
		adminGroup.PUT("/agent/update", a.UpdateAgent)
		adminGroup.DELETE("/agent/delete", a.DeleteAgent)
		adminGroup.POST("/quickReply/create", a.CreateQuickReply)
		adminGroup.PUT("/quickReply/update", a.UpdateQuickReply)
		adminGroup.DELETE("/quickReply/delete", a.DeleteQuickReply)
		adminGroup.GET("/blacklist/list", a.GetBlacklist)
		adminGroup.POST("/blacklist/add", a.AddToBlacklist)
		adminGroup.DELETE("/blacklist/remove", a.RemoveFromBlacklist)
		adminGroup.GET("/config/get", a.GetCsConfig)
		adminGroup.PUT("/config/update", middleware.OperationRecord(), a.UpdateCsConfig)
	}
}
