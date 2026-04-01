package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type PhoneAreaCodeSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
}

// PhoneLoginRequest 手机号登录请求
type PhoneLoginRequest struct {
	AreaCode string `json:"areaCode" binding:"required"` // 国际区号
	Phone    string `json:"phone" binding:"required"`     // 手机号
	Password string `json:"password" binding:"required"`  // 密码
	Captcha  string `json:"captcha" binding:"required"`   // 验证码
	CaptchaId string `json:"captchaId" binding:"required"` // 验证码ID
}

// PhoneRegisterRequest 手机号注册请求
type PhoneRegisterRequest struct {
	AreaCode   string `json:"areaCode" binding:"required"`   // 国际区号
	Phone      string `json:"phone" binding:"required"`      // 手机号
	Password   string `json:"password" binding:"required"`   // 密码
	Captcha    string `json:"captcha" binding:"required"`    // 验证码
	CaptchaId  string `json:"captchaId" binding:"required"`  // 验证码ID
	InviteCode string `json:"inviteCode"`                    // 邀请码
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`  // 旧密码(方式一)
	NewPassword string `json:"newPassword" binding:"required"`  // 新密码
	Method      string `json:"method" binding:"required"`       // 修改方式: old_password / phone_verify
	Captcha     string `json:"captcha" binding:"required"`      // 系统验证码
	CaptchaId   string `json:"captchaId" binding:"required"`    // 验证码ID
}
