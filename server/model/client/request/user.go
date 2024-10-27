package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type ClientUserSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`

	Username string `json:"username" form:"username" `
	Nickname string `json:"nickname" form:"nickname" `
	request.PageInfo
}

// 用户 结构体  ShopUser
type CreateUser struct {
	Username   string `json:"username" form:"username" gorm:"column:username;comment:用户名;" binding:"required"` //用户名
	Password   string `json:"password" form:"password" gorm:"column:password;comment:密码;" binding:"required"`  //密码
	RePassword string `json:"rePassword" form:"rePassword" gorm:"-"`                                           // 确认密码
}

type Login struct {
	Username  string `json:"username"`  // 用户名
	Password  string `json:"password"`  // 密码
	Captcha   string `json:"captcha"`   // 验证码
	CaptchaId string `json:"captchaId"` // 验证码ID
}

type UpdateKV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
