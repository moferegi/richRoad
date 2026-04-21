package system

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// SysBannedIP IP封禁记录
type SysBannedIP struct {
	global.GVA_MODEL
	IP        string     `json:"ip" gorm:"column:ip;uniqueIndex;not null;comment:封禁的IP地址"`
	Reason    string     `json:"reason" gorm:"column:reason;comment:封禁原因"`
	BannedBy  string     `json:"bannedBy" gorm:"column:banned_by;comment:封禁操作人"`
	IsAuto    bool       `json:"isAuto" gorm:"column:is_auto;default:false;comment:是否自动封禁"`
	ExpiredAt *time.Time `json:"expiredAt" gorm:"column:expired_at;comment:过期时间，null表示永久封禁"`
}

func (SysBannedIP) TableName() string {
	return "sys_banned_ips"
}
