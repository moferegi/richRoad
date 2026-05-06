// 自动生成模板ClientUser
package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/google/uuid"
)

// 客户端用户 结构体  ClientUser
type ClientUser struct {
	global.GVA_MODEL
	UUID             uuid.UUID `json:"uuid" form:"uuid" gorm:"column:uuid;comment:UUID;"`                                                          //UUID
	Username         string    `json:"username" form:"username" gorm:"column:username;comment:用户名;"`                                               //用户名
	Password         string    `json:"password" form:"password" gorm:"column:password;comment:密码;"`                                                //密码
	Avatar           string    `json:"avatar" form:"avatar" gorm:"column:avatar;comment:头像;type:text"`                                             //头像
	Nickname         string    `json:"nickname" form:"nickname" gorm:"column:nickname;comment:昵称;"`                                                //昵称
	Gender           string    `json:"gender" form:"gender" gorm:"column:gender;comment:性别;"`                                                      //性别
	Phone            string    `json:"phone" form:"phone" gorm:"column:phone;comment:手机号;"`                                                        //手机号
	Email            string    `json:"email" form:"email" gorm:"column:email;comment:邮箱;"`                                                         //邮箱
	OpenID           string    `json:"openID" form:"openID" gorm:"column:open_id;comment:OpenID;"`                                                 //OpenID
	Point            int       `json:"point" form:"point" gorm:"column:point;comment:积分;default:0"`                                                //积分
	TryonPoint       int       `json:"tryonPoint" form:"tryonPoint" gorm:"column:tryon_point;comment:试衣币;default:0"`                               //试衣币
	InviteCode       string    `json:"inviteCode" form:"inviteCode" gorm:"column:invite_code;uniqueIndex;comment:邀请码;size:16"`                     //邀请码
	InvitedBy        uint      `json:"invitedBy" form:"invitedBy" gorm:"column:invited_by;comment:邀请人ID;default:0;index"`                          //邀请人ID
	SubOrderRewarded bool      `json:"subOrderRewarded" form:"subOrderRewarded" gorm:"column:sub_order_rewarded;default:false;comment:下级首单已奖励上级;"` //下级首单已奖励上级
	AreaCode         string    `json:"areaCode" form:"areaCode" gorm:"column:area_code;size:10;comment:国际区号;"`                                     //国际区号
	Banned           *bool     `json:"banned" form:"banned" gorm:"column:banned;default:false;comment:是否封禁;"`                                      //是否封禁
	CreatedBy        uint      `gorm:"column:created_by;comment:创建者"`
	UpdatedBy        uint      `gorm:"column:updated_by;comment:更新者"`
	DeletedBy        uint      `gorm:"column:deleted_by;comment:删除者"`
}

// TableName 客户端用户 ClientUser自定义表名 client_user
func (ClientUser) TableName() string {
	return "client_user"
}
