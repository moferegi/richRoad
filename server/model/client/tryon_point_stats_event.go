package client

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TryonPointStatsEvent 试衣币统计事实事件（独立于流水表，供统计页查询）
type TryonPointStatsEvent struct {
	global.GVA_MODEL
	SourcePointRecordID uint      `json:"sourcePointRecordID" form:"sourcePointRecordID" gorm:"column:source_point_record_id;default:0;uniqueIndex;comment:来源积分流水ID(client_point_records.id);"`
	EventAt             time.Time `json:"eventAt" form:"eventAt" gorm:"column:event_at;index;comment:事件发生时间;"`
	UserID              uint      `json:"userID" form:"userID" gorm:"column:user_id;default:0;index;comment:用户ID;"`
	ChangeType          string    `json:"changeType" form:"changeType" gorm:"column:change_type;size:20;index;comment:增减类型(increase/decrease);"`
	OperationType       string    `json:"operationType" form:"operationType" gorm:"column:operation_type;size:50;index;comment:操作类型;"`
	PointChange         int       `json:"pointChange" form:"pointChange" gorm:"column:point_change;default:0;comment:变化值(带符号);"`
	PointAmount         int       `json:"pointAmount" form:"pointAmount" gorm:"column:point_amount;default:0;comment:变化绝对值;"`
	Reason              string    `json:"reason" form:"reason" gorm:"column:reason;size:200;comment:变化原因;"`
}

func (TryonPointStatsEvent) TableName() string {
	return "client_tryon_point_stats_event"
}
