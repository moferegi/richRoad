package api

import (
	"encoding/json"
	"net/http"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	// 允许所有来源（生产环境可收严格一些）
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WsApi struct{}

// UserWS 客户端用户 WebSocket 连接入口
// @Tags CustomerServiceWS
// @Summary 客户端用户连接客服WebSocket
// @Security ApiKeyAuth
// @Router /cs/ws [get]
func (a *WsApi) UserWS(c *gin.Context) {
	// 从 query 参数解析 token（浏览器 WebSocket 不方便设置 header）
	tokenStr := c.Query("token")
	if tokenStr == "" {
		response.FailWithMessage("token 不能为空", c)
		return
	}
	j := utils.NewJWT()
	claims, err := j.ParseToken(tokenStr)
	if err != nil {
		response.FailWithMessage("token 无效", c)
		return
	}
	clientUserID := claims.BaseClaims.ID
	if clientUserID == 0 {
		response.FailWithMessage("无效的用户", c)
		return
	}

	// 升级为 WebSocket
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.GVA_LOG.Error("WebSocket 升级失败", zap.Error(err))
		return
	}

	wsConn := service.CSHub.RegisterUser(clientUserID, conn)

	// 获取或创建会话
	conv, _, err := service.Service.ConversationService.CreateOrGetActive(clientUserID)
	if err != nil {
		_ = conn.WriteMessage(websocket.TextMessage, service.MakeFrame(service.WSEventError, map[string]string{
			"message": err.Error(),
		}))
		_ = conn.Close()
		return
	}

	// 发送当前会话信息（让客户端知道会话状态和队列位置）
	_ = conn.WriteMessage(websocket.TextMessage, service.MakeFrame("session", conv))

	// 发送最近5条消息
	msgs, _ := service.Service.MessageService.GetLast5(conv.ID)
	if len(msgs) > 0 {
		_ = conn.WriteMessage(websocket.TextMessage, service.MakeFrame("history", msgs))
	}

	go wsConn.WritePump()
	wsConn.ReadPump(func(raw []byte) {
		a.handleUserMessage(clientUserID, conv.ID, raw)
	})
}

// handleUserMessage 处理用户发来的 WS 消息
func (a *WsApi) handleUserMessage(clientUserID, convID uint, raw []byte) {
	var frame csReq.WSMessageFrame
	if err := json.Unmarshal(raw, &frame); err != nil {
		return
	}
	switch frame.Event {
	case "send_message":
		if frame.Content == "" || frame.MsgType == "" {
			return
		}
		_, err := service.Service.MessageService.Send(
			model.SenderTypeUser, clientUserID, convID,
			frame.MsgType, frame.Content, frame.ClientMsgID,
		)
		if err != nil {
			service.CSHub.SendToUser(clientUserID, service.MakeFrame(service.WSEventError, map[string]string{
				"message": err.Error(),
			}))
		}
	case "revoke":
		if frame.MessageID == 0 {
			return
		}
		if err := service.Service.MessageService.Revoke(frame.MessageID, clientUserID, model.SenderTypeUser); err != nil {
			service.CSHub.SendToUser(clientUserID, service.MakeFrame(service.WSEventError, map[string]string{
				"message": err.Error(),
			}))
		}
	case "read":
		_ = service.Service.MessageService.MarkRead(convID, model.SenderTypeUser)
	case "ping":
		service.CSHub.SendToUser(clientUserID, service.MakeFrame(service.WSEventPong, nil))
	}
}

// AgentWS 坐席 WebSocket 连接入口
// @Tags CustomerServiceWS
// @Summary 坐席连接客服WebSocket
// @Security ApiKeyAuth
// @Router /cs/wsAgent [get]
func (a *WsApi) AgentWS(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		response.FailWithMessage("token 不能为空", c)
		return
	}
	j := utils.NewJWT()
	claims, err := j.ParseToken(tokenStr)
	if err != nil {
		response.FailWithMessage("token 无效", c)
		return
	}
	agentSysUserID := claims.BaseClaims.ID
	if agentSysUserID == 0 {
		response.FailWithMessage("无效的用户", c)
		return
	}

	// 检查坐席是否配置
	agent, err := service.Service.AgentService.GetOrCreate(agentSysUserID)
	if err != nil || !agent.IsEnabled {
		response.FailWithMessage("您未被配置为客服坐席", c)
		return
	}

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.GVA_LOG.Error("坐席 WebSocket 升级失败", zap.Error(err))
		return
	}

	wsConn := service.CSHub.RegisterAgent(agentSysUserID, conn)

	// 发送当前进行中的会话列表
	convSearch := csReq.ConversationSearch{}
	convSearch.PageSize = 50
	convSearch.AgentUserID = &agentSysUserID
	activeStatus := model.ConvStatusActive
	convSearch.Status = activeStatus
	convList, _, _ := service.Service.ConversationService.GetList(convSearch)
	_ = conn.WriteMessage(websocket.TextMessage, service.MakeFrame("convList", convList))

	go wsConn.WritePump()
	wsConn.ReadPump(func(raw []byte) {
		a.handleAgentMessage(agentSysUserID, raw)
	})
}

// handleAgentMessage 处理坐席发来的 WS 消息
func (a *WsApi) handleAgentMessage(agentSysUserID uint, raw []byte) {
	var frame csReq.WSMessageFrame
	if err := json.Unmarshal(raw, &frame); err != nil {
		return
	}
	// 坐席消息需要指定会话 ID
	// frame.MessageID 这里用作 ConversationID（复用字段）
	convID := frame.MessageID

	switch frame.Event {
	case "send_message":
		if frame.Content == "" || frame.MsgType == "" || convID == 0 {
			return
		}
		_, err := service.Service.MessageService.Send(
			model.SenderTypeAgent, agentSysUserID, convID,
			frame.MsgType, frame.Content, frame.ClientMsgID,
		)
		if err != nil {
			service.CSHub.SendToAgent(agentSysUserID, service.MakeFrame(service.WSEventError, map[string]string{
				"message": err.Error(),
			}))
		}
	case "revoke":
		if frame.MessageID == 0 {
			return
		}
		if err := service.Service.MessageService.Revoke(frame.MessageID, agentSysUserID, model.SenderTypeAgent); err != nil {
			service.CSHub.SendToAgent(agentSysUserID, service.MakeFrame(service.WSEventError, map[string]string{
				"message": err.Error(),
			}))
		}
	case "read":
		if convID > 0 {
			_ = service.Service.MessageService.MarkRead(convID, model.SenderTypeAgent)
		}
	case "close":
		if convID > 0 {
			if err := service.Service.ConversationService.Close(convID, model.ConvClosedByAgent); err != nil {
				service.CSHub.SendToAgent(agentSysUserID, service.MakeFrame(service.WSEventError, map[string]string{
					"message": err.Error(),
				}))
			}
		}
	case "ping":
		service.CSHub.SendToAgent(agentSysUserID, service.MakeFrame(service.WSEventPong, nil))
	}
}
