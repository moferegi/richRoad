package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model/request"
	"gorm.io/gorm"
)

type ConversationService struct{}

// CreateOrGetActive 获取或创建用户的活跃/排队会话
func (s *ConversationService) CreateOrGetActive(clientUserID uint) (*model.CsConversation, bool, error) {
	// 检查黑名单
	var bl model.CsBlacklist
	if err := global.GVA_DB.Where("client_user_id = ?", clientUserID).First(&bl).Error; err == nil {
		return nil, false, errors.New("您已被限制使用客服功能")
	}

	// 检查是否有未关闭的会话
	var conv model.CsConversation
	err := global.GVA_DB.Where("client_user_id = ? AND status != ?", clientUserID, model.ConvStatusClosed).
		Order("created_at DESC").First(&conv).Error
	if err == nil {
		return &conv, false, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, err
	}

	// 创建新会话
	conv = model.CsConversation{
		ClientUserID: clientUserID,
		Status:       model.ConvStatusPending,
	}
	if err := global.GVA_DB.Create(&conv).Error; err != nil {
		return nil, false, err
	}

	// 尝试立即分配坐席
	_ = s.tryAssignAgent(&conv)

	return &conv, true, nil
}

// tryAssignAgent 尝试为会话分配空闲坐席（轮询策略：最少活跃会话数）
func (s *ConversationService) tryAssignAgent(conv *model.CsConversation) error {
	onlineIDs := CSHub.OnlineAgentIDs()
	if len(onlineIDs) == 0 {
		// 无在线坐席，放入排队
		return s.updateQueuePositions()
	}

	// 查询在线且启用的坐席
	type agentLoad struct {
		UserID      uint
		MaxSessions int
		Active      int64
	}
	var agents []agentLoad
	global.GVA_DB.Raw(`
		SELECT a.user_id, a.max_sessions,
		       COUNT(c.id) AS active
		FROM cs_agents a
		LEFT JOIN cs_conversations c ON c.agent_user_id = a.user_id AND c.status = 'active'
		WHERE a.is_enabled = 1 AND a.user_id IN ?
		GROUP BY a.user_id, a.max_sessions
		HAVING COUNT(c.id) < a.max_sessions
		ORDER BY active ASC
		LIMIT 1
	`, onlineIDs).Scan(&agents)

	if len(agents) == 0 {
		return s.updateQueuePositions()
	}

	agentUserID := agents[0].UserID
	return s.assignToAgent(conv, agentUserID)
}

// assignToAgent 将会话分配给指定坐席
func (s *ConversationService) assignToAgent(conv *model.CsConversation, agentUserID uint) error {
	now := time.Now()
	_ = now
	if err := global.GVA_DB.Model(conv).Updates(map[string]interface{}{
		"agent_user_id":  agentUserID,
		"status":         model.ConvStatusActive,
		"queue_position": 0,
	}).Error; err != nil {
		return err
	}
	conv.AgentUserID = &agentUserID
	conv.Status = model.ConvStatusActive
	conv.QueuePosition = 0

	// 通知用户：已分配坐席
	CSHub.SendToUser(conv.ClientUserID, MakeFrame(WSEventAssigned, map[string]interface{}{
		"conversationId": conv.ID,
		"agentUserId":    agentUserID,
	}))
	// 通知坐席：新会话
	CSHub.SendToAgent(agentUserID, MakeFrame(WSEventNewConv, conv))
	return nil
}

// updateQueuePositions 更新排队会话的位置号
func (s *ConversationService) updateQueuePositions() error {
	var pending []model.CsConversation
	if err := global.GVA_DB.Where("status = ?", model.ConvStatusPending).
		Order("created_at ASC").Find(&pending).Error; err != nil {
		return err
	}
	for i, conv := range pending {
		pos := i + 1
		global.GVA_DB.Model(&conv).Update("queue_position", pos)
		CSHub.SendToUser(conv.ClientUserID, MakeFrame(WSEventQueue, map[string]interface{}{
			"conversationId": conv.ID,
			"position":       pos,
		}))
	}
	return nil
}

// AssignToAgent 手动将会话分配给坐席（管理员操作）
func (s *ConversationService) AssignToAgent(convID, agentUserID uint) error {
	var conv model.CsConversation
	if err := global.GVA_DB.First(&conv, convID).Error; err != nil {
		return err
	}
	if conv.Status == model.ConvStatusClosed {
		return errors.New("会话已关闭，无法分配")
	}
	return s.assignToAgent(&conv, agentUserID)
}

// Transfer 转接会话给另一个坐席
func (s *ConversationService) Transfer(convID, targetAgentUserID uint) error {
	var conv model.CsConversation
	if err := global.GVA_DB.First(&conv, convID).Error; err != nil {
		return err
	}
	if conv.Status == model.ConvStatusClosed {
		return errors.New("会话已关闭，无法转接")
	}
	oldAgentID := conv.AgentUserID
	if err := global.GVA_DB.Model(&conv).Updates(map[string]interface{}{
		"agent_user_id": targetAgentUserID,
		"status":        model.ConvStatusActive,
	}).Error; err != nil {
		return err
	}

	// 系统消息：通知转接
	svc := Service.MessageService
	_, _ = svc.SaveSystemMessage(convID, fmt.Sprintf("会话已转接给其他坐席"))

	// 通知旧坐席
	if oldAgentID != nil {
		CSHub.SendToAgent(*oldAgentID, MakeFrame(WSEventTransfer, map[string]interface{}{
			"conversationId":    convID,
			"targetAgentUserId": targetAgentUserID,
		}))
	}
	// 通知新坐席
	conv.AgentUserID = &targetAgentUserID
	CSHub.SendToAgent(targetAgentUserID, MakeFrame(WSEventNewConv, conv))
	// 通知用户
	CSHub.SendToUser(conv.ClientUserID, MakeFrame(WSEventTransfer, map[string]interface{}{
		"conversationId": convID,
	}))
	return nil
}

// Close 关闭会话
func (s *ConversationService) Close(convID uint, closedBy string) error {
	var conv model.CsConversation
	if err := global.GVA_DB.First(&conv, convID).Error; err != nil {
		return err
	}
	if conv.Status == model.ConvStatusClosed {
		return nil
	}
	now := time.Now()
	if err := global.GVA_DB.Model(&conv).Updates(map[string]interface{}{
		"status":    model.ConvStatusClosed,
		"closed_at": now,
		"closed_by": closedBy,
	}).Error; err != nil {
		return err
	}

	// 系统消息
	svc := Service.MessageService
	_, _ = svc.SaveSystemMessage(convID, "会话已结束")

	closeData := map[string]interface{}{
		"conversationId": convID,
		"closedBy":       closedBy,
	}
	CSHub.SendToUser(conv.ClientUserID, MakeFrame(WSEventClosed, closeData))
	if conv.AgentUserID != nil {
		CSHub.SendToAgent(*conv.AgentUserID, MakeFrame(WSEventClosed, closeData))
	}

	// 有新的 pending 会话，检查能否分配
	go s.tryAssignPending()
	return nil
}

// tryAssignPending 关闭会话后尝试分配等待中的会话
func (s *ConversationService) tryAssignPending() {
	var pending model.CsConversation
	if err := global.GVA_DB.Where("status = ?", model.ConvStatusPending).
		Order("created_at ASC").First(&pending).Error; err != nil {
		return
	}
	_ = s.tryAssignAgent(&pending)
}

// Rate 用户评价会话
func (s *ConversationService) Rate(convID, clientUserID uint, rating int) error {
	var conv model.CsConversation
	if err := global.GVA_DB.Where("id = ? AND client_user_id = ?", convID, clientUserID).First(&conv).Error; err != nil {
		return errors.New("会话不存在")
	}
	if conv.Status != model.ConvStatusClosed {
		return errors.New("请在会话结束后评价")
	}
	if conv.Rating != nil {
		return errors.New("您已评价过该会话")
	}
	now := time.Now()
	return global.GVA_DB.Model(&conv).Updates(map[string]interface{}{
		"rating":   rating,
		"rated_at": now,
	}).Error
}

// GetList 分页获取会话列表（管理员）
func (s *ConversationService) GetList(search request.ConversationSearch) ([]model.CsConversation, int64, error) {
	db := global.GVA_DB.Model(&model.CsConversation{})
	if search.Status != "" {
		db = db.Where("status = ?", search.Status)
	}
	if search.AgentUserID != nil {
		db = db.Where("agent_user_id = ?", *search.AgentUserID)
	}
	if search.ClientUserID != nil {
		db = db.Where("client_user_id = ?", *search.ClientUserID)
	}
	if search.StartTime != nil {
		db = db.Where("created_at >= ?", search.StartTime)
	}
	if search.EndTime != nil {
		db = db.Where("created_at <= ?", search.EndTime)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.CsConversation
	if err := db.Order("created_at DESC").Scopes(search.Paginate()).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// GetUserActive 获取客户端用户当前的活跃/排队会话
func (s *ConversationService) GetUserActive(clientUserID uint) (*model.CsConversation, error) {
	var conv model.CsConversation
	err := global.GVA_DB.Where("client_user_id = ? AND status != ?", clientUserID, model.ConvStatusClosed).
		Order("created_at DESC").First(&conv).Error
	if err != nil {
		return nil, err
	}
	return &conv, nil
}
