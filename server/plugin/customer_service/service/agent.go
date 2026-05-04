package service

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model"
)

type AgentService struct{}

// GetOrCreate 获取或创建坐席记录
func (s *AgentService) GetOrCreate(userID uint) (*model.CsAgent, error) {
	var agent model.CsAgent
	err := global.GVA_DB.Where("user_id = ?", userID).First(&agent).Error
	if err == nil {
		return &agent, nil
	}
	agent = model.CsAgent{
		UserID:      userID,
		MaxSessions: 5,
		IsEnabled:   true,
	}
	if err := global.GVA_DB.Create(&agent).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

// GetEnabledByUserID 获取已启用的坐席（用于接口鉴权）
func (s *AgentService) GetEnabledByUserID(userID uint) (*model.CsAgent, error) {
	var agent model.CsAgent
	if err := global.GVA_DB.Where("user_id = ?", userID).First(&agent).Error; err != nil {
		return nil, errors.New("当前账号不是客服坐席")
	}
	if !agent.IsEnabled {
		return nil, errors.New("当前坐席已被禁用")
	}
	return &agent, nil
}

// GetList 获取坐席列表
func (s *AgentService) GetList() ([]model.CsAgent, error) {
	var agents []model.CsAgent
	if err := global.GVA_DB.Order("id ASC").Find(&agents).Error; err != nil {
		return nil, err
	}
	// 批量查询 sys_users 昵称
	if len(agents) > 0 {
		userIDs := make([]uint, 0, len(agents))
		for _, a := range agents {
			userIDs = append(userIDs, a.UserID)
		}
		type nickRow struct {
			ID       uint
			NickName string
		}
		var rows []nickRow
		global.GVA_DB.Table("sys_users").Select("id, nick_name").Where("id IN ?", userIDs).Scan(&rows)
		nickMap := make(map[uint]string, len(rows))
		for _, r := range rows {
			nickMap[r.ID] = r.NickName
		}
		for i := range agents {
			agents[i].SysNickname = nickMap[agents[i].UserID]
		}
	}
	// 填充在线状态和当前会话数
	for i := range agents {
		if CSHub.IsAgentOnline(agents[i].UserID) {
			agents[i].OnlineStatus = "online"
		} else {
			agents[i].OnlineStatus = "offline"
		}
		var cnt int64
		global.GVA_DB.Model(&model.CsConversation{}).
			Where("agent_user_id = ? AND status = ?", agents[i].UserID, model.ConvStatusActive).
			Count(&cnt)
		agents[i].ActiveSessions = int(cnt)
	}
	return agents, nil
}

// Update 更新坐席配置
func (s *AgentService) Update(agent model.CsAgent) error {
	return global.GVA_DB.Save(&agent).Error
}

// Delete 删除坐席
func (s *AgentService) Delete(id uint) error {
	var agent model.CsAgent
	if err := global.GVA_DB.First(&agent, id).Error; err != nil {
		return err
	}
	var cnt int64
	global.GVA_DB.Model(&model.CsConversation{}).
		Where("agent_user_id = ? AND status = ?", agent.UserID, model.ConvStatusActive).
		Count(&cnt)
	if cnt > 0 {
		return errors.New("该坐席还有进行中的会话，无法删除")
	}
	return global.GVA_DB.Delete(&agent).Error
}

// SetEnabled 启用/禁用坐席
func (s *AgentService) SetEnabled(id uint, enabled bool) error {
	return global.GVA_DB.Model(&model.CsAgent{}).Where("id = ?", id).
		Update("is_enabled", enabled).Error
}
