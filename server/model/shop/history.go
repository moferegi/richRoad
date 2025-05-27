// 自动生成模板Good
package shop

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// 商品 结构体  Good
type History struct {
	global.GVA_MODEL
	GoodID int  `json:"good_id"`
	UserID uint `json:"user_id"`
}

// TableName 商品 Good自定义表名 shop_good
func (History) TableName() string {
	return "shop_good_history"
}
