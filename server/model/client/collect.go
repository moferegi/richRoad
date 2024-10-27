// 自动生成模板Collect
package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 收藏 结构体  Collect
type Collect struct {
	global.GVA_MODEL
	UserID uint `json:"userID" form:"userID" gorm:"column:user_id;comment:用户ID;"` //用户ID
	GoodID uint `json:"goodID" form:"goodID" gorm:"column:good_id;comment:商品ID;"` //商品ID
}

// TableName 收藏 Collect自定义表名 clinet_collect
func (Collect) TableName() string {
	return "clinet_collect"
}
