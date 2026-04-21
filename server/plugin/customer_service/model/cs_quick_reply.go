package model

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// CsQuickReply 快捷回复
type CsQuickReply struct {
	global.GVA_MODEL
	Title     string `json:"title" gorm:"column:title;not null;size:100;comment:标题"`
	Content   string `json:"content" gorm:"column:content;type:text;not null;comment:回复内容"`
	Sort      int    `json:"sort" gorm:"column:sort;default:0;comment:排序值(越小越靠前)"`
	CreatedBy uint   `json:"createdBy" gorm:"column:created_by;comment:创建人ID"`
}

func (CsQuickReply) TableName() string { return "cs_quick_replies" }
