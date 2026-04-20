package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type ClientUserSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`

	Username string `json:"username" form:"username" `
	Nickname string `json:"nickname" form:"nickname" `
	OrderBy  string `json:"orderBy" form:"orderBy"`
	OrderDir string `json:"orderDir" form:"orderDir"`
	request.PageInfo
}

// 用户 结构体  ShopUser
type CreateUser struct {
	Username   string `json:"username" form:"username" gorm:"column:username;comment:用户名;" binding:"required"` //用户名
	Password   string `json:"password" form:"password" gorm:"column:password;comment:密码;" binding:"required"`  //密码
	RePassword string `json:"rePassword" form:"rePassword" gorm:"-"`                                           // 确认密码
	InviteCode string `json:"inviteCode" form:"inviteCode"`                                                    // 邀请码（可选）
	Captcha    string `json:"captcha" form:"captcha"`                                                          // 验证码
	CaptchaId  string `json:"captchaId" form:"captchaId"`                                                      // 验证码ID
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

type SetPhoneRequest struct {
	Password  string `json:"password"`  // 当前密码
	Phone     string `json:"phone"`     // 新手机号
	Captcha   string `json:"captcha"`   // 验证码
	CaptchaId string `json:"captchaId"` // 验证码ID
}
