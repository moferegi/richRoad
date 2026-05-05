package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// SysConfig 系统参数配置
type SysConfig struct {
	global.GVA_MODEL
	ConfigKey   string `json:"configKey" form:"configKey" gorm:"column:config_key;size:100;uniqueIndex;comment:参数键名;"`
	ConfigValue string `json:"configValue" form:"configValue" gorm:"column:config_value;type:text;comment:参数值;"`
	ConfigName  string `json:"configName" form:"configName" gorm:"column:config_name;size:100;comment:参数名称;"`
	ConfigGroup string `json:"configGroup" form:"configGroup" gorm:"column:config_group;size:50;index;comment:参数分组;"`
	Remark      string `json:"remark" form:"remark" gorm:"column:remark;size:500;comment:备注;"`
}

func (SysConfig) TableName() string {
	return "client_sys_config"
}
