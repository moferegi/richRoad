package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ExternalLinkDomain 外部链接域名管理
type ExternalLinkDomain struct {
	global.GVA_MODEL
	Name      string `json:"name" form:"name" gorm:"column:name;size:100;comment:域名名称;"`
	Domain    string `json:"domain" form:"domain" gorm:"column:domain;size:500;comment:域名地址(如https://cdn.example.com);"`
	IsDefault *bool  `json:"isDefault" form:"isDefault" gorm:"column:is_default;default:false;comment:是否默认;"`
	IsEnabled *bool  `json:"isEnabled" form:"isEnabled" gorm:"column:is_enabled;default:true;comment:是否启用;"`
	Sort      int    `json:"sort" form:"sort" gorm:"column:sort;default:0;comment:排序(越大越前);"`
	Remark    string `json:"remark" form:"remark" gorm:"column:remark;size:500;comment:备注;"`
}

func (ExternalLinkDomain) TableName() string {
	return "client_external_link_domain"
}
