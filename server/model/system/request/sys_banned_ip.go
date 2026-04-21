package request

import common "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

// BanIPRequest 封禁IP请求
type BanIPRequest struct {
	IP       string `json:"ip" binding:"required" example:"192.168.1.1"`
	Reason   string `json:"reason" example:"手动封禁"`
	Duration int    `json:"duration" example:"60"` // 封禁时长（分钟），0=永久
}

// UnbanIPRequest 解封IP请求
type UnbanIPRequest struct {
	IP string `json:"ip" binding:"required" example:"192.168.1.1"`
}

// BannedIPSearch 封禁IP查询请求
type BannedIPSearch struct {
	common.PageInfo
	IP     string `json:"ip" form:"ip"`
	IsAuto *bool  `json:"isAuto" form:"isAuto"`
}

// AttackStatsRequest 攻击统计查询请求
type AttackStatsRequest struct {
	Hours int `json:"hours" form:"hours" example:"24"` // 查询最近N小时，默认24
}
