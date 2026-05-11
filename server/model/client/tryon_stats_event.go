package client

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TryonStatsEvent 试衣统计事实事件（独立于业务任务与调用日志，供统计页查询）
type TryonStatsEvent struct {
	global.GVA_MODEL
	SourceLogID  uint      `json:"sourceLogID" form:"sourceLogID" gorm:"column:source_log_id;default:0;uniqueIndex;comment:来源日志ID(client_model_call_log.id);"`
	EventAt      time.Time `json:"eventAt" form:"eventAt" gorm:"column:event_at;index;comment:事件发生时间;"`
	UserID       uint      `json:"userID" form:"userID" gorm:"column:user_id;default:0;index;comment:用户ID;"`
	TaskID       uint      `json:"taskID" form:"taskID" gorm:"column:task_id;default:0;index;comment:试衣任务ID;"`
	TaskNo       string    `json:"taskNo" form:"taskNo" gorm:"column:task_no;size:40;index;comment:试衣任务编号;"`
	RequestID    string    `json:"requestID" form:"requestID" gorm:"column:request_id;size:80;index;comment:请求幂等ID;"`
	SceneType    string    `json:"sceneType" form:"sceneType" gorm:"column:scene_type;size:20;index;comment:场景类型;"`
	ModelUsage   string    `json:"modelUsage" form:"modelUsage" gorm:"column:model_usage;size:20;index;comment:模型用途;"`
	ModelKey     string    `json:"modelKey" form:"modelKey" gorm:"column:model_key;size:120;index;comment:模型key;"`
	Provider     string    `json:"provider" form:"provider" gorm:"column:provider;size:80;index;comment:模型提供商;"`
	Status       string    `json:"status" form:"status" gorm:"column:status;size:32;index;comment:调用状态;"`
	CostPoints   int       `json:"costPoints" form:"costPoints" gorm:"column:cost_points;default:0;comment:该次模型调用扣币;"`
	RefundPoints int       `json:"refundPoints" form:"refundPoints" gorm:"column:refund_points;default:0;comment:该次模型调用退币;"`
}

func (TryonStatsEvent) TableName() string {
	return "client_tryon_stats_event"
}
