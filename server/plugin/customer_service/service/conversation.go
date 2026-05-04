package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model/request"
	"gorm.io/gorm"
)

type ConversationService struct{}

const pendingDispatchInterval = 15 * time.Second

type lastMessageRow struct {
	ConversationID uint   `gorm:"column:conversation_id"`
	Content        string `gorm:"column:content"`
	MsgType        string `gorm:"column:msg_type"`
	Revoked        bool   `gorm:"column:revoked"`
}

type unreadCountRow struct {
	ConversationID uint  `gorm:"column:conversation_id"`
	Cnt            int64 `gorm:"column:cnt"`
}

func summarizeLastMessage(content, msgType string, revoked bool) string {
	if revoked {
		return "[已撤回]"
	}
	switch msgType {
	case model.MsgTypeImage:
		return "[图片]"
	case model.MsgTypeEmoji:
		return "[表情]"
	default:
		if content == "" {
			return "暂无消息"
		}
		return content
	}
}

// fillConversationSummaries 批量填充会话摘要信息（最后一条消息 + 未读数）
func fillConversationSummaries(list []model.CsConversation, unreadSenderType string) {
	if len(list) == 0 {
		return
	}
	ids := make([]uint, 0, len(list))
	for _, c := range list {
		ids = append(ids, c.ID)
	}

	// 最后一条消息
	var rows []lastMessageRow
	global.GVA_DB.Raw(`
		SELECT m.conversation_id, m.content, m.msg_type, m.revoked
		FROM cs_messages m
		INNER JOIN (
			SELECT conversation_id, MAX(created_at) AS max_created_at
			FROM cs_messages
			WHERE conversation_id IN ?
			GROUP BY conversation_id
		) latest
		ON latest.conversation_id = m.conversation_id AND latest.max_created_at = m.created_at
		WHERE m.conversation_id IN ?
		ORDER BY m.created_at DESC
	`, ids, ids).Scan(&rows)

	lastMap := make(map[uint]string, len(rows))
	for _, r := range rows {
		if _, ok := lastMap[r.ConversationID]; ok {
			continue
		}
		lastMap[r.ConversationID] = summarizeLastMessage(r.Content, r.MsgType, r.Revoked)
	}

	// 未读计数
	var unreadRows []unreadCountRow
	global.GVA_DB.Model(&model.CsMessage{}).
		Select("conversation_id, COUNT(1) AS cnt").
		Where("conversation_id IN ? AND sender_type = ? AND read_at IS NULL", ids, unreadSenderType).
		Group("conversation_id").
		Scan(&unreadRows)

	unreadMap := make(map[uint]int, len(unreadRows))
	for _, r := range unreadRows {
		unreadMap[r.ConversationID] = int(r.Cnt)
	}

	for i := range list {
		if msg, ok := lastMap[list[i].ID]; ok {
			list[i].LastMsg = msg
		}
		list[i].UnreadCount = unreadMap[list[i].ID]
	}
}

// CreateOrGetActive 获取或创建用户的活跃/排队会话
func (s *ConversationService) CreateOrGetActive(clientUserID uint) (*model.CsConversation, bool, error) {
	// 平台客服开关：关闭时不允许创建/进入会话
	cfg, _ := Service.ConfigService.GetConfig()
	if !cfg.PlatEnabled {
		return nil, false, errors.New("平台客服暂未开启")
	}

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
		// 找到已有会话：如果仍在排队，趁用户重新上线再试一次分配
		if conv.Status == model.ConvStatusPending {
			_ = s.tryAssignAgent(&conv)
		}
		tmp := []model.CsConversation{conv}
		fillConversationSummaries(tmp, model.SenderTypeAgent)
		conv = tmp[0]
		return &conv, false, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, err
	}

	// 无活跃会话时，优先复用最近一次已关闭会话，保证用户重新进入可看到历史记录
	if err := global.GVA_DB.Where("client_user_id = ? AND status = ?", clientUserID, model.ConvStatusClosed).
		Order("updated_at DESC").First(&conv).Error; err == nil {
		if err := global.GVA_DB.Model(&conv).Updates(map[string]interface{}{
			"agent_user_id":  nil,
			"status":         model.ConvStatusPending,
			"queue_position": 0,
			"closed_at":      nil,
			"closed_by":      "",
			"rating":         nil,
			"rated_at":       nil,
		}).Error; err != nil {
			return nil, false, err
		}
		conv.AgentUserID = nil
		conv.Status = model.ConvStatusPending
		conv.QueuePosition = 0
		conv.ClosedAt = nil
		conv.ClosedBy = ""
		conv.Rating = nil
		conv.RatedAt = nil

		_ = s.tryAssignAgent(&conv)
		tmp := []model.CsConversation{conv}
		fillConversationSummaries(tmp, model.SenderTypeAgent)
		conv = tmp[0]
		return &conv, false, nil
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
	tmp := []model.CsConversation{conv}
	fillConversationSummaries(tmp, model.SenderTypeAgent)
	conv = tmp[0]

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
	// 更新剩余排队会话位置
	go s.updateQueuePositions()
	return nil
}

// ensureAgentCapacity 校验坐席是否还有可接待容量
func (s *ConversationService) ensureAgentCapacity(agentUserID uint) error {
	agent, err := Service.AgentService.GetEnabledByUserID(agentUserID)
	if err != nil {
		return err
	}
	var activeCnt int64
	if err := global.GVA_DB.Model(&model.CsConversation{}).
		Where("agent_user_id = ? AND status = ?", agentUserID, model.ConvStatusActive).
		Count(&activeCnt).Error; err != nil {
		return err
	}
	if activeCnt >= int64(agent.MaxSessions) {
		return errors.New("当前坐席会话已满，请稍后再试")
	}
	return nil
}

// AcceptByAgent 坐席手动接入排队会话
func (s *ConversationService) AcceptByAgent(convID, agentUserID uint) (*model.CsConversation, error) {
	var conv model.CsConversation
	if err := global.GVA_DB.First(&conv, convID).Error; err != nil {
		return nil, errors.New("会话不存在")
	}
	if conv.Status == model.ConvStatusClosed {
		return nil, errors.New("会话已关闭，无法接入")
	}
	if conv.AgentUserID != nil {
		if *conv.AgentUserID == agentUserID {
			return &conv, nil
		}
		return nil, errors.New("会话已被其他坐席接入")
	}
	if conv.Status != model.ConvStatusPending {
		return nil, errors.New("当前会话状态不可接入")
	}
	if err := s.ensureAgentCapacity(agentUserID); err != nil {
		return nil, err
	}
	result := global.GVA_DB.Model(&model.CsConversation{}).
		Where("id = ? AND status = ? AND agent_user_id IS NULL", convID, model.ConvStatusPending).
		Updates(map[string]interface{}{
			"agent_user_id":  agentUserID,
			"status":         model.ConvStatusActive,
			"queue_position": 0,
		})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		var latest model.CsConversation
		if err := global.GVA_DB.First(&latest, convID).Error; err != nil {
			return nil, errors.New("会话不存在")
		}
		if latest.AgentUserID != nil {
			if *latest.AgentUserID == agentUserID {
				return &latest, nil
			}
			return nil, errors.New("会话已被其他坐席接入")
		}
		return nil, errors.New("当前会话状态不可接入")
	}
	conv.AgentUserID = &agentUserID
	conv.Status = model.ConvStatusActive
	conv.QueuePosition = 0

	CSHub.SendToUser(conv.ClientUserID, MakeFrame(WSEventAssigned, map[string]interface{}{
		"conversationId": conv.ID,
		"agentUserId":    agentUserID,
	}))
	CSHub.SendToAgent(agentUserID, MakeFrame(WSEventNewConv, conv))
	go s.updateQueuePositions()
	return &conv, nil
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
	if _, err := Service.AgentService.GetEnabledByUserID(agentUserID); err != nil {
		return err
	}
	var conv model.CsConversation
	if err := global.GVA_DB.First(&conv, convID).Error; err != nil {
		return err
	}
	if conv.Status == model.ConvStatusClosed {
		return errors.New("会话已关闭，无法分配")
	}
	if conv.AgentUserID != nil && *conv.AgentUserID == agentUserID {
		return nil
	}
	if err := s.ensureAgentCapacity(agentUserID); err != nil {
		return err
	}
	return s.assignToAgent(&conv, agentUserID)
}

// Transfer 转接会话给另一个坐席
func (s *ConversationService) Transfer(convID, targetAgentUserID uint) error {
	if _, err := Service.AgentService.GetEnabledByUserID(targetAgentUserID); err != nil {
		return err
	}
	var conv model.CsConversation
	if err := global.GVA_DB.First(&conv, convID).Error; err != nil {
		return err
	}
	if conv.Status == model.ConvStatusClosed {
		return errors.New("会话已关闭，无法转接")
	}
	oldAgentID := conv.AgentUserID
	if oldAgentID != nil && *oldAgentID == targetAgentUserID {
		return nil
	}
	if err := s.ensureAgentCapacity(targetAgentUserID); err != nil {
		return err
	}
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

// tryAssignPending 关闭会话后尝试分配等待中的会话（取最早一条）
func (s *ConversationService) tryAssignPending() {
	var pending model.CsConversation
	if err := global.GVA_DB.Where("status = ?", model.ConvStatusPending).
		Order("created_at ASC").First(&pending).Error; err != nil {
		return
	}
	_ = s.tryAssignAgent(&pending)
}

// TryAssignAllPending 坐席上线时触发：将所有排队中的会话逐一尝试分配
func (s *ConversationService) TryAssignAllPending() {
	var pendings []model.CsConversation
	if err := global.GVA_DB.Where("status = ?", model.ConvStatusPending).
		Order("created_at ASC").Find(&pendings).Error; err != nil {
		return
	}
	for i := range pendings {
		_ = s.tryAssignAgent(&pendings[i])
	}
}

// StartPendingDispatchWorker 启动排队会话自动分配重试任务
func (s *ConversationService) StartPendingDispatchWorker() {
	go func() {
		ticker := time.NewTicker(pendingDispatchInterval)
		defer ticker.Stop()
		for range ticker.C {
			if len(CSHub.OnlineAgentIDs()) == 0 {
				continue
			}
			s.TryAssignAllPending()
		}
	}()
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
	// 批量填充 ClientNickname
	fillClientNicknames(list)
	// 坐席/管理列表默认统计“用户消息未读”
	fillConversationSummaries(list, model.SenderTypeUser)
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
	tmp := []model.CsConversation{conv}
	fillConversationSummaries(tmp, model.SenderTypeAgent)
	conv = tmp[0]
	return &conv, nil
}

// GetByID 根据会话ID查询会话
func (s *ConversationService) GetByID(convID uint) (*model.CsConversation, error) {
	var conv model.CsConversation
	if err := global.GVA_DB.First(&conv, convID).Error; err != nil {
		return nil, errors.New("会话不存在")
	}
	return &conv, nil
}

// CheckClientOwnership 校验会话是否属于当前客户端用户
func (s *ConversationService) CheckClientOwnership(clientUserID, convID uint) error {
	conv, err := s.GetByID(convID)
	if err != nil {
		return err
	}
	if conv.ClientUserID != clientUserID {
		return errors.New("无权访问该会话")
	}
	return nil
}

// CheckAgentAccess 校验坐席是否可访问会话
// allowPending=true 时允许访问全局待分配会话（pending 且 agent_user_id 为空）
func (s *ConversationService) CheckAgentAccess(agentUserID, convID uint, allowPending bool) error {
	conv, err := s.GetByID(convID)
	if err != nil {
		return err
	}
	if allowPending && conv.Status == model.ConvStatusPending && conv.AgentUserID == nil {
		return nil
	}
	if conv.AgentUserID == nil || *conv.AgentUserID != agentUserID {
		return errors.New("无权操作该会话")
	}
	return nil
}

// fillClientNicknames 批量查询并填充会话列表中的客户端用户昵称
func fillClientNicknames(list []model.CsConversation) {
	if len(list) == 0 {
		return
	}
	ids := make([]uint, 0, len(list))
	for _, c := range list {
		ids = append(ids, c.ClientUserID)
	}
	var users []struct {
		ID       uint
		Nickname string
	}
	global.GVA_DB.Model(&clientModel.ClientUser{}).
		Select("id, nickname").
		Where("id IN ?", ids).
		Find(&users)
	nickMap := make(map[uint]string, len(users))
	for _, u := range users {
		if u.Nickname != "" {
			nickMap[u.ID] = u.Nickname
		}
	}
	for i := range list {
		if nick, ok := nickMap[list[i].ClientUserID]; ok {
			list[i].ClientNickname = nick
		}
	}
}
