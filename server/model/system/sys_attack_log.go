package system

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// SysAttackLog 攻击行为日志
// attack_type 枚举：login_fail | captcha_fail | register_limit
type SysAttackLog struct {
	global.GVA_MODEL
	IP         string `json:"ip" gorm:"column:ip;index;comment:来源IP"`
	AttackType string `json:"attackType" gorm:"column:attack_type;comment:攻击类型"`
	Username   string `json:"username" gorm:"column:username;comment:尝试的用户名/手机号"`
	Detail     string `json:"detail" gorm:"column:detail;comment:详情"`
}

func (SysAttackLog) TableName() string {
	return "sys_attack_logs"
}
