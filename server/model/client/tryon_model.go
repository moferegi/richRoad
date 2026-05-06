package client

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// TryonModel 我的模特
type TryonModel struct {
	global.GVA_MODEL
	UserID uint   `json:"userID" form:"userID" gorm:"column:user_id;comment:用户ID;index"`
	Name   string `json:"name" form:"name" gorm:"column:name;size:128;comment:模特名称"`
	Image  string `json:"image" form:"image" gorm:"column:image;type:text;comment:模特图片"`
}

// TableName 我的模特 TryonModel自定义表名 client_tryon_model
func (TryonModel) TableName() string {
	return "client_tryon_model"
}
