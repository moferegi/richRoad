// 自动生成模板Comment
package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"gorm.io/datatypes"
	"time"
)

// 用户评论 结构体  Comment
type Comment struct {
	global.GVA_MODEL
	UserID      uint               `json:"userID" form:"userID" gorm:"column:user_id;comment:用户ID;"`                       //用户ID
	User        *client.ClientUser `json:"user" form:"user" gorm:"foreignKey:UserID;references:ID;"`                       //用户
	OrderID     uint               `json:"orderID" form:"orderID" gorm:"column:order_id;comment:订单ID;" binding:"required"` //订单ID
	Order       *OrderDetailRes    `json:"order" form:"order" gorm:"foreignKey:OrderID;references:ID;"`                    //订单
	GoodID      uint               `json:"goodID" form:"goodID" gorm:"column:good_id;comment:商品ID;"`                       //商品ID
	SKUID       uint               `json:"SKUID" form:"SKUID" gorm:"column:sku_id;comment:SKUID;"`                         //SKUID
	Pics        datatypes.JSON     `json:"pics" form:"pics" gorm:"column:pics;comment:评论图片"`                               //评论图片
	Rating      uint               `json:"rating" form:"rating" gorm:"column:rating;comment:;" binding:"required"`         //用户评分
	Content     string             `json:"content" form:"content" gorm:"column:content;comment:;" binding:"required"`      //评论内容
	ShopReply   string             `json:"shopReply" form:"shopReply" gorm:"column:shop_reply;comment:;"`                  //商家回复
	ShopReplyAt *time.Time         `json:"shopReplyAt" form:"shopReplyAt" gorm:"column:shop_reply_at;comment:;"`           //商家回复创建时间
}

// TableName 用户评论 Comment自定义表名 shop_comment
func (Comment) TableName() string {
	return "shop_comment"
}
