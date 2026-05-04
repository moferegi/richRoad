package api

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func checkWSOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		// 非浏览器客户端通常不带 Origin（如 uni-app 原生）
		return true
	}
	u, err := url.Parse(origin)
	if err != nil || u.Host == "" {
		return false
	}
	reqHost := r.Host
	if strings.EqualFold(u.Host, reqHost) {
		return true
	}
	originHost := strings.ToLower(u.Hostname())
	reqHostOnly := extractHostName(reqHost)
	if originHost == reqHostOnly {
		return true
	}
	// 本地联调放行：localhost / 127.0.0.1 / ::1 之间互通
	if isLoopbackHost(originHost) && isLoopbackHost(reqHostOnly) {
		return true
	}
	// 允许同主域下的子域名互通（如 web.xxx 与 back.xxx）
	allowed := sameRootDomain(originHost, reqHostOnly)
	if !allowed {
		global.GVA_LOG.Warn("WebSocket Origin 校验拒绝",
			zap.String("origin", origin),
			zap.String("originHost", originHost),
			zap.String("requestHost", reqHost),
			zap.String("requestHostOnly", reqHostOnly),
			zap.String("path", r.URL.Path),
		)
	}
	return allowed
}

func extractHostName(host string) string {
	host = strings.TrimSpace(host)
	if host == "" {
		return ""
	}
	if h, _, err := net.SplitHostPort(host); err == nil {
		return strings.ToLower(strings.Trim(h, "[]"))
	}
	return strings.ToLower(strings.Trim(host, "[]"))
}

func isLoopbackHost(host string) bool {
	host = strings.ToLower(strings.TrimSpace(strings.Trim(host, "[]")))
	if host == "" {
		return false
	}
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func sameRootDomain(a, b string) bool {
	if net.ParseIP(a) != nil || net.ParseIP(b) != nil {
		return false
	}
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	if len(aParts) < 2 || len(bParts) < 2 {
		return false
	}
	aRoot := aParts[len(aParts)-2] + "." + aParts[len(aParts)-1]
	bRoot := bParts[len(bParts)-2] + "." + bParts[len(bParts)-1]
	return aRoot == bRoot
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	Subprotocols:    []string{"bearer", "token"},
	// 仅允许同源/同主域来源
	CheckOrigin: checkWSOrigin,
}

type WsApi struct{}

const defaultMaxWSConnsPerIP = 8

var (
	wsMaxConnsPerIP  = loadWSMaxConnsPerIP()
	wsAllowQueryAuth = loadWSAllowQueryAuth()
)

func loadWSMaxConnsPerIP() int {
	v := strings.TrimSpace(os.Getenv("CS_WS_MAX_CONNS_PER_IP"))
	if v == "" {
		return defaultMaxWSConnsPerIP
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		return defaultMaxWSConnsPerIP
	}
	return n
}

func loadWSAllowQueryAuth() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("CS_WS_ALLOW_QUERY_TOKEN")))
	if v == "" {
		// 默认关闭 query token，避免 token 经 URL 透传导致泄漏风险。
		return false
	}
	switch v {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return false
	}
}

func tokenFromProtocolHeader(v string) string {
	if strings.TrimSpace(v) == "" {
		return ""
	}
	parts := strings.Split(v, ",")
	tokens := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			tokens = append(tokens, t)
		}
	}
	if len(tokens) == 0 {
		return ""
	}
	if len(tokens) >= 2 && (strings.EqualFold(tokens[0], "bearer") || strings.EqualFold(tokens[0], "token")) {
		return tokens[1]
	}
	for _, t := range tokens {
		if strings.Count(t, ".") >= 2 {
			return t
		}
	}
	return ""
}

func extractWSToken(c *gin.Context) (string, string) {
	if auth := strings.TrimSpace(c.GetHeader("Authorization")); auth != "" {
		const bearer = "Bearer "
		if len(auth) > len(bearer) && strings.EqualFold(auth[:len(bearer)], bearer) {
			return strings.TrimSpace(auth[len(bearer):]), "authorization"
		}
		return auth, "authorization"
	}
	if t := tokenFromProtocolHeader(c.GetHeader("Sec-WebSocket-Protocol")); t != "" {
		return t, "protocol"
	}
	if t := strings.TrimSpace(c.Query("token")); t != "" {
		return t, "query"
	}
	return "", ""
}

// UserWS 客户端用户 WebSocket 连接入口
// @Tags CustomerServiceWS
// @Summary 客户端用户连接客服WebSocket
// @Security ApiKeyAuth
// @Router /cs/ws [get]
func (a *WsApi) UserWS(c *gin.Context) {
	// 兼容多来源：Authorization / Sec-WebSocket-Protocol / query(token)
	tokenStr, tokenSource := extractWSToken(c)
	if tokenStr == "" {
		response.FailWithMessage("token 不能为空", c)
		return
	}
	if tokenSource == "query" {
		if !wsAllowQueryAuth {
			response.FailWithMessage("当前环境禁止通过URL传递token，请升级客户端", c)
			return
		}
		global.GVA_LOG.Warn("WebSocket 使用 query token，建议升级为子协议/Authorization",
			zap.String("ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
		)
	}
	j := utils.NewJWT()
	claims, err := j.ParseToken(tokenStr)
	if err != nil {
		response.FailWithMessage("token 无效", c)
		return
	}
	if claims.AuthorityId != 8080 {
		response.FailWithMessage("仅客户端用户可连接该通道", c)
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

	wsConn, regErr := service.CSHub.RegisterUser(clientUserID, c.ClientIP(), conn, wsMaxConnsPerIP)
	if regErr != nil {
		global.GVA_LOG.Warn("WebSocket 用户连接被IP限流", zap.String("ip", c.ClientIP()), zap.Uint("userID", clientUserID), zap.Int("limit", wsMaxConnsPerIP))
		_ = conn.WriteMessage(websocket.TextMessage, service.MakeFrame(service.WSEventError, map[string]string{
			"message": regErr.Error(),
		}))
		_ = conn.Close()
		return
	}

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

	// 如果会话仍在排队（分配失败或无坐席），广播给所有在线坐席刷新列表
	if conv.Status == model.ConvStatusPending {
		go service.CSHub.BroadcastToAgents(service.MakeFrame("pending_conv", conv))
	}

	// 发送最近一页消息历史（page=1, pageSize=10），与前端 REST 分页保持一致
	// 这样前端 historyPage=1，上拉时请求 page=2 不会产生数据空隙
	initSearch := csReq.MessageSearch{
		PageInfo:       commonReq.PageInfo{Page: 1, PageSize: 10},
		ConversationID: conv.ID,
	}
	initMsgs, _, _ := service.Service.MessageService.GetHistory(initSearch)
	// GetHistory 返回 DESC 顺序，反转为时间正序再发送
	for i, j := 0, len(initMsgs)-1; i < j; i, j = i+1, j-1 {
		initMsgs[i], initMsgs[j] = initMsgs[j], initMsgs[i]
	}
	if len(initMsgs) > 0 {
		_ = conn.WriteMessage(websocket.TextMessage, service.MakeFrame("history", initMsgs))
	}

	go wsConn.WritePump()
	wsConn.ReadPump(func(raw []byte) {
		a.handleUserMessage(clientUserID, conv.ID, raw)
	})

	// 兜底处理：用户断开后若持续离线，则自动关闭会话，避免坐席端卡在“进行中”
	go func(userID, convID uint) {
		time.Sleep(8 * time.Second)
		if service.CSHub.IsUserOnline(userID) {
			return
		}
		_ = service.Service.ConversationService.Close(convID, model.ConvClosedByUser)
	}(clientUserID, conv.ID)
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
	case "close":
		if convID > 0 {
			if err := service.Service.ConversationService.Close(convID, model.ConvClosedByUser); err != nil {
				service.CSHub.SendToUser(clientUserID, service.MakeFrame(service.WSEventError, map[string]string{
					"message": err.Error(),
				}))
			}
		}
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
	tokenStr, tokenSource := extractWSToken(c)
	if tokenStr == "" {
		response.FailWithMessage("token 不能为空", c)
		return
	}
	if tokenSource == "query" {
		if !wsAllowQueryAuth {
			response.FailWithMessage("当前环境禁止通过URL传递token，请升级客户端", c)
			return
		}
		global.GVA_LOG.Warn("WebSocket 使用 query token，建议升级为子协议/Authorization",
			zap.String("ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
		)
	}
	j := utils.NewJWT()
	claims, err := j.ParseToken(tokenStr)
	if err != nil {
		response.FailWithMessage("token 无效", c)
		return
	}
	if claims.AuthorityId == 8080 {
		response.FailWithMessage("客户端用户不能连接坐席通道", c)
		return
	}
	agentSysUserID := claims.BaseClaims.ID
	if agentSysUserID == 0 {
		response.FailWithMessage("无效的用户", c)
		return
	}

	// 仅允许已配置且启用的坐席连接
	if _, err := service.Service.AgentService.GetEnabledByUserID(agentSysUserID); err != nil {
		response.FailWithMessage("您未被配置为客服坐席", c)
		return
	}

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.GVA_LOG.Error("坐席 WebSocket 升级失败", zap.Error(err))
		return
	}

	wsConn, regErr := service.CSHub.RegisterAgent(agentSysUserID, c.ClientIP(), conn, wsMaxConnsPerIP)
	if regErr != nil {
		global.GVA_LOG.Warn("WebSocket 坐席连接被IP限流", zap.String("ip", c.ClientIP()), zap.Uint("userID", agentSysUserID), zap.Int("limit", wsMaxConnsPerIP))
		_ = conn.WriteMessage(websocket.TextMessage, service.MakeFrame(service.WSEventError, map[string]string{
			"message": regErr.Error(),
		}))
		_ = conn.Close()
		return
	}

	// 发送当前进行中的会话列表
	convSearch := csReq.ConversationSearch{}
	convSearch.PageSize = 50
	convSearch.AgentUserID = &agentSysUserID
	activeStatus := model.ConvStatusActive
	convSearch.Status = activeStatus
	convList, _, _ := service.Service.ConversationService.GetList(convSearch)
	_ = conn.WriteMessage(websocket.TextMessage, service.MakeFrame("convList", convList))

	// 坐席上线后，尝试将所有排队中的会话分配出去
	go service.Service.ConversationService.TryAssignAllPending()

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
		if err := service.Service.ConversationService.CheckAgentAccess(agentSysUserID, convID, true); err != nil {
			service.CSHub.SendToAgent(agentSysUserID, service.MakeFrame(service.WSEventError, map[string]string{
				"message": err.Error(),
			}))
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
			if err := service.Service.ConversationService.CheckAgentAccess(agentSysUserID, convID, true); err == nil {
				_ = service.Service.MessageService.MarkRead(convID, model.SenderTypeAgent)
			}
		}
	case "close":
		if convID > 0 {
			if err := service.Service.ConversationService.CheckAgentAccess(agentSysUserID, convID, false); err != nil {
				service.CSHub.SendToAgent(agentSysUserID, service.MakeFrame(service.WSEventError, map[string]string{
					"message": err.Error(),
				}))
				return
			}
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
