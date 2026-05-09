package client

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// TryonCloth 我的衣橱
// Category 取值: upper/lower/onepiece/shoes
type TryonCloth struct {
	global.GVA_MODEL
	UserID   uint   `json:"userID" form:"userID" gorm:"column:user_id;comment:用户ID;index"`
	Name     string `json:"name" form:"name" gorm:"column:name;size:128;comment:衣橱名称"`
	Category string `json:"category" form:"category" gorm:"column:category;size:32;comment:分类(upper/lower/onepiece/shoes);index"`
	Image    string `json:"image" form:"image" gorm:"column:image;type:text;comment:衣橱图片"`
}

// TableName 我的衣橱 TryonCloth自定义表名 client_tryon_cloth
func (TryonCloth) TableName() string {
	return "client_tryon_cloth"
}
