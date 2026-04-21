package service

import (
	"errors"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model/request"
	"gorm.io/gorm"
)

type MessageService struct{}

// Send 发送消息（REST / WS 共用逻辑）
func (s *MessageService) Send(senderType string, senderID, convID uint, msgType, content, clientMsgID string) (*model.CsMessage, error) {
	// 检查会话状态
	var conv model.CsConversation
	if err := global.GVA_DB.First(&conv, convID).Error; err != nil {
		return nil, errors.New("会话不存在")
	}
	if conv.Status == model.ConvStatusClosed {
		return nil, errors.New("会话已关闭")
	}

	// client_msg_id 去重
	if clientMsgID != "" {
		var exist model.CsMessage
		if err := global.GVA_DB.Where("client_msg_id = ?", clientMsgID).First(&exist).Error; err == nil {
			return &exist, nil // 幂等返回
		}
	}

	msg := model.CsMessage{
		ConversationID: convID,
		SenderType:     senderType,
		SenderID:       senderID,
		MsgType:        msgType,
		Content:        content,
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
