package client

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// SignIn 签到记录 结构体
type SignIn struct {
	global.GVA_MODEL
	UserID   uint      `json:"userID" form:"userID" gorm:"column:user_id;index;comment:用户ID;"`                       //用户ID
	SignDate time.Time `json:"signDate" form:"signDate" gorm:"column:sign_date;type:date;index;comment:签到日期;"`        //签到日期
}

// TableName 签到记录自定义表名
func (SignIn) TableName() string {
	return "client_sign_in"
}
