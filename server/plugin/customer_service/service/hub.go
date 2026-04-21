package service

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = 54 * time.Second
	maxMsgSize = 1024 * 1024 // 1MB
)

// WS 事件类型（服务端 → 客户端）
const (
	WSEventMessage  = "message"  // 新消息
	WSEventRevoke   = "revoke"   // 消息已撤回
	WSEventRead     = "read"     // 已读回执
	WSEventAssigned = "assigned" // 会话已分配
	WSEventClosed   = "closed"   // 会话已关闭
	WSEventQueue    = "queue"    // 排队状态更新
	WSEventTyping   = "typing"   // 正在输入
	WSEventPong     = "pong"     // 心跳响应
	WSEventError    = "error"    // 错误通知
	WSEventTransfer = "transfer" // 会话被转接
	WSEventNewConv  = "new_conv" // 坐席收到新会话
)

// WSFrame WebSocket 通信帧
type WSFrame struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

// MakeFrame 构造 WebSocket 帧
func MakeFrame(event string, data interface{}) []byte {
	raw, _ := json.Marshal(data)
	frame := WSFrame{Event: event, Data: raw}
	b, _ := json.Marshal(frame)
	return b
}

// Connection 一个 WebSocket 连接
type Connection struct {
	ID   uint   // userID (client) 或 sysUserID (agent)
	Kind string // "user" | "agent"
	Conn *websocket.Conn
	Send chan []byte
	hub  *Hub
}

// Hub 管理所有 WebSocket 连接
type Hub struct {
	mu         sync.RWMutex
	userConns  map[uint]*Connection // clientUserID → conn
	agentConns map[uint]*Connection // sysUserID    → conn
}

// CSHub 全局 Hub 单例
var CSHub = newHub()

func newHub() *Hub {
	return &Hub{
		userConns:  make(map[uint]*Connection),
		agentConns: make(map[uint]*Connection),
	}
}

// RegisterUser 注册用户连接（踢掉旧连接）
func (h *Hub) RegisterUser(userID uint, conn *websocket.Conn) *Connection {
	c := &Connection{ID: userID, Kind: "user", Conn: conn, Send: make(chan []byte, 256), hub: h}
	h.mu.Lock()
	if old, ok := h.userConns[userID]; ok {
		close(old.Send)
	}
	h.userConns[userID] = c
	h.mu.Unlock()
	return c
}

// RegisterAgent 注册坐席连接（踢掉旧连接）
func (h *Hub) RegisterAgent(agentUserID uint, conn *websocket.Conn) *Connection {
	c := &Connection{ID: agentUserID, Kind: "agent", Conn: conn, Send: make(chan []byte, 512), hub: h}
	h.mu.Lock()
	if old, ok := h.agentConns[agentUserID]; ok {
		close(old.Send)
	}
	h.agentConns[agentUserID] = c
	h.mu.Unlock()
	return c
}

// UnregisterUser 注销用户连接
func (h *Hub) UnregisterUser(userID uint) {
	h.mu.Lock()
	if c, ok := h.userConns[userID]; ok && c.hub != nil {
		delete(h.userConns, userID)
	}
	h.mu.Unlock()
}

// UnregisterAgent 注销坐席连接
func (h *Hub) UnregisterAgent(agentUserID uint) {
	h.mu.Lock()
	if c, ok := h.agentConns[agentUserID]; ok && c.hub != nil {
		delete(h.agentConns, agentUserID)
	}
	h.mu.Unlock()
}

// SendToUser 向用户推送消息，返回是否在线
func (h *Hub) SendToUser(userID uint, data []byte) bool {
	h.mu.RLock()
	c, ok := h.userConns[userID]
	h.mu.RUnlock()
	if !ok {
		return false
	}
	select {
	case c.Send <- data:
		return true
	default:
		return false
	}
}

// SendToAgent 向坐席推送消息，返回是否在线
func (h *Hub) SendToAgent(agentUserID uint, data []byte) bool {
	h.mu.RLock()
	c, ok := h.agentConns[agentUserID]
	h.mu.RUnlock()
	if !ok {
		return false
	}
	select {
	case c.Send <- data:
		return true
	default:
		return false
	}
}

// BroadcastToAgents 向所有在线坐席广播
func (h *Hub) BroadcastToAgents(data []byte) {
	h.mu.RLock()
	conns := make([]*Connection, 0, len(h.agentConns))
	for _, c := range h.agentConns {
		conns = append(conns, c)
	}
	h.mu.RUnlock()
	for _, c := range conns {
		select {
		case c.Send <- data:
		default:
		}
	}
}

// IsUserOnline 判断用户是否在线
func (h *Hub) IsUserOnline(userID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.userConns[userID]
	return ok
}

// IsAgentOnline 判断坐席是否在线
func (h *Hub) IsAgentOnline(agentUserID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.agentConns[agentUserID]
	return ok
}

// OnlineAgentIDs 获取所有在线坐席的 sysUserID 列表
func (h *Hub) OnlineAgentIDs() []uint {
	h.mu.RLock()
	defer h.mu.RUnlock()
	ids := make([]uint, 0, len(h.agentConns))
	for id := range h.agentConns {
		ids = append(ids, id)
	}
	return ids
}

// ---- 读写 pump ----

// ReadPump 持续读 WebSocket 帧，断线后清理连接
func (c *Connection) ReadPump(onMessage func([]byte)) {
	defer func() {
		if c.Kind == "user" {
			c.hub.UnregisterUser(c.ID)
		} else {
			c.hub.UnregisterAgent(c.ID)
		}
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMsgSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		return c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		onMessage(msg)
	}
}

// WritePump 持续向 WebSocket 写帧，支持 ping
func (c *Connection) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// channel 已关闭
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
