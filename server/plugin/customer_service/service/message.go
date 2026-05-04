package service

import (
	"context"
	"errors"
	"fmt"
	"html"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model/request"
	"gorm.io/gorm"
)

type MessageService struct{}

const (
	maxMessageTextLen = 2000
	maxImageURLLen    = 2048
)

// checkRateLimit 检查消息发送频率（Redis INCR，每秒最多 2 条）
func checkRateLimit(senderType string, senderID uint) error {
	if global.GVA_REDIS == nil {
		return nil // Redis 未启用，跳过限流
	}
	key := fmt.Sprintf("cs_rate:%s:%d", senderType, senderID)
	ctx := context.Background()
	count, err := global.GVA_REDIS.Incr(ctx, key).Result()
	if err != nil {
		return nil // Redis 异常不阻断正常流程
	}
	if count == 1 {
		// 第一次计数，设置 1 秒过期
		global.GVA_REDIS.Expire(ctx, key, time.Second)
	}
	if count > 2 {
		return errors.New("发送过于频繁，请稍后再试")
	}
	return nil
}

// normalizeAndValidateMessage 校验并规范化消息内容
func normalizeAndValidateMessage(msgType, content string) (string, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return "", errors.New("消息内容不能为空")
	}

	switch msgType {
	case model.MsgTypeText, model.MsgTypeEmoji, model.MsgTypeSystem:
		if utf8.RuneCountInString(content) > maxMessageTextLen {
			return "", fmt.Errorf("消息长度不能超过 %d 字", maxMessageTextLen)
		}
		// 统一做 HTML 转义，降低 XSS 风险
		return html.EscapeString(content), nil
	case model.MsgTypeImage:
		if len(content) > maxImageURLLen {
			return "", errors.New("图片地址过长")
		}
		lower := strings.ToLower(content)
		if strings.HasPrefix(lower, "data:") || strings.HasPrefix(lower, "javascript:") {
			return "", errors.New("图片地址不合法")
		}
		if !(strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") || strings.HasPrefix(content, "/")) {
			return "", errors.New("图片地址不合法")
		}
		return content, nil
	default:
		return "", errors.New("不支持的消息类型")
	}
}

// Send 发送消息（REST / WS 共用逻辑）
func (s *MessageService) Send(senderType string, senderID, convID uint, msgType, content, clientMsgID string) (*model.CsMessage, error) {
	// 频率限制：用户每秒最多 2 条
	if senderType == model.SenderTypeUser {
		if err := checkRateLimit(senderType, senderID); err != nil {
			return nil, err
		}
	}
	// 检查会话状态
	var conv model.CsConversation
	if err := global.GVA_DB.First(&conv, convID).Error; err != nil {
		return nil, errors.New("会话不存在")
	}
	if conv.Status == model.ConvStatusClosed {
		return nil, errors.New("会话已关闭")
	}
	if senderType == model.SenderTypeAgent && conv.Status == model.ConvStatusPending && conv.AgentUserID == nil {
		if _, err := Service.ConversationService.AcceptByAgent(convID, senderID); err != nil {
			return nil, err
		}
		if err := global.GVA_DB.First(&conv, convID).Error; err != nil {
			return nil, errors.New("会话不存在")
		}
	}

	// 发送方与会话归属校验，避免越权发送
	if senderType == model.SenderTypeUser && conv.ClientUserID != senderID {
		return nil, errors.New("无权向该会话发送消息")
	}
	if senderType == model.SenderTypeAgent {
		if conv.AgentUserID == nil || *conv.AgentUserID != senderID {
			return nil, errors.New("无权向该会话发送消息")
		}
	}

	if clientMsgID != "" {
		clientMsgID = strings.TrimSpace(clientMsgID)
		if len(clientMsgID) > 64 {
			return nil, errors.New("clientMsgId 过长")
		}
	}

	normalizedContent, err := normalizeAndValidateMessage(msgType, content)
	if err != nil {
		return nil, err
	}

	// client_msg_id 去重
	if clientMsgID != "" {
		var exist model.CsMessage
		if err := global.GVA_DB.Where(
			"conversation_id = ? AND sender_type = ? AND sender_id = ? AND client_msg_id = ?",
			convID, senderType, senderID, clientMsgID,
		).First(&exist).Error; err == nil {
			return &exist, nil // 幂等返回
		}
	}

	msg := model.CsMessage{
		ConversationID: convID,
		SenderType:     senderType,
		SenderID:       senderID,
		MsgType:        msgType,
		Content:        normalizedContent,
		ClientMsgID:    clientMsgID,
	}
	if err := global.GVA_DB.Create(&msg).Error; err != nil {
		return nil, err
	}

	// 推送 WS
	frame := MakeFrame(WSEventMessage, msg)
	if senderType == model.SenderTypeUser {
		// 用户发消息 → 推送给坐席
		if conv.AgentUserID != nil {
			CSHub.SendToAgent(*conv.AgentUserID, frame)
		} else {
			// 无坐席，广播给所有在线坐席
			CSHub.BroadcastToAgents(frame)
		}
	} else {
		// 坐席/系统发消息 → 推送给用户
		CSHub.SendToUser(conv.ClientUserID, frame)
	}
	return &msg, nil
}

// SaveSystemMessage 保存系统消息（会话状态变更提示）
func (s *MessageService) SaveSystemMessage(convID uint, content string) (*model.CsMessage, error) {
	return s.Send(model.SenderTypeSystem, 0, convID, model.MsgTypeSystem, content, "")
}

// Revoke 撤回消息（仅允许 2 分钟内）
func (s *MessageService) Revoke(msgID, senderID uint, senderType string) error {
	var msg model.CsMessage
	if err := global.GVA_DB.First(&msg, msgID).Error; err != nil {
		return errors.New("消息不存在")
	}
	if msg.SenderID != senderID || msg.SenderType != senderType {
		return errors.New("无权撤回他人消息")
	}
	if msg.Revoked {
		return errors.New("消息已撤回")
	}
	if time.Since(msg.CreatedAt) > 2*time.Minute {
		return errors.New("超过2分钟的消息不能撤回")
	}
	now := time.Now()
	if err := global.GVA_DB.Model(&msg).Updates(map[string]interface{}{
		"revoked":    true,
		"revoked_at": now,
		"content":    "[已撤回]",
	}).Error; err != nil {
		return err
	}

	// 获取会话推送给对方
	var conv model.CsConversation
	if err := global.GVA_DB.First(&conv, msg.ConversationID).Error; err == nil {
		frame := MakeFrame(WSEventRevoke, map[string]interface{}{
			"messageId":      msgID,
			"conversationId": msg.ConversationID,
		})
		if senderType == model.SenderTypeUser {
			if conv.AgentUserID != nil {
				CSHub.SendToAgent(*conv.AgentUserID, frame)
			}
		} else {
			CSHub.SendToUser(conv.ClientUserID, frame)
		}
	}
	return nil
}

// MarkRead 标记消息为已读
func (s *MessageService) MarkRead(convID uint, readerType string) error {
	now := time.Now()
	// 标记对方发来的消息为已读
	var otherType string
	if readerType == model.SenderTypeUser {
		otherType = model.SenderTypeAgent
	} else {
		otherType = model.SenderTypeUser
	}
	return global.GVA_DB.Model(&model.CsMessage{}).
		Where("conversation_id = ? AND sender_type = ? AND read_at IS NULL", convID, otherType).
		Update("read_at", now).Error
}

// GetHistory 分页获取消息历史（降序，最新在前）
func (s *MessageService) GetHistory(search request.MessageSearch) ([]model.CsMessage, int64, error) {
	db := global.GVA_DB.Model(&model.CsMessage{}).Where("conversation_id = ?", search.ConversationID)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.CsMessage
	if err := db.Order("created_at DESC").Scopes(search.Paginate()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// GetLast5 获取会话最近5条消息（坐席打开会话时加载）
func (s *MessageService) GetLast5(convID uint) ([]model.CsMessage, error) {
	var msgs []model.CsMessage
	if err := global.GVA_DB.Where("conversation_id = ?", convID).
		Order("created_at DESC").Limit(5).Find(&msgs).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	// 反转为时间正序
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	return msgs, nil
}
