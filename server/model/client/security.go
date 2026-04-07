package client

import "time"

// LoginFailRecord 登录失败记录（DB回退方案，Redis不可用时使用）
type LoginFailRecord struct {
	ID          uint      `gorm:"primaryKey"`
	Username    string    `gorm:"column:username;size:100;uniqueIndex;comment:用户名或手机号"`
	FailCount   int       `gorm:"column:fail_count;default:0;comment:连续失败次数"`
	LockedUntil time.Time `gorm:"column:locked_until;comment:锁定截止时间"`
}

func (LoginFailRecord) TableName() string {
	return "client_login_fail_record"
}

// RegisterIPRecord 注册IP记录（DB回退方案，Redis不可用时使用）
type RegisterIPRecord struct {
	ID            uint   `gorm:"primaryKey"`
	IP            string `gorm:"column:ip;size:50;uniqueIndex;comment:IP地址"`
	RegisterCount int    `gorm:"column:register_count;default:0;comment:注册次数"`
}

func (RegisterIPRecord) TableName() string {
	return "client_register_ip_record"
}
