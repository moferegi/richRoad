package model

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// CsBlacklist 客服黑名单（禁止联系客服）
type CsBlacklist struct {
	global.GVA_MODEL
	ClientUserID uint   `json:"clientUserId" gorm:"column:client_user_id;uniqueIndex;not null;comment:客户端用户ID"`
	Reason       string `json:"reason" gorm:"column:reason;comment:封禁原因"`
	CreatedBy    uint   `json:"createdBy" gorm:"column:created_by;comment:操作人系统用户ID"`
	// 非持久化字段
	ClientNickname string `json:"clientNickname" gorm:"-"`
}

func (CsBlacklist) TableName() string { return "cs_blacklists" }
