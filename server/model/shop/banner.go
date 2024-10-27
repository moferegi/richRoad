// 自动生成模板Banner
package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 轮播图 结构体  Banner
type Banner struct {
	global.GVA_MODEL
	Title string `json:"title" form:"title" gorm:"column:title;comment:轮播标题;" binding:"required"` //轮播标题
	Src   string `json:"src" form:"src" gorm:"column:src;comment:图片连接;" binding:"required"`       //图片连接
	Href  string `json:"href" form:"href" gorm:"column:href;comment:图片跳转链接;"`                     //跳转链接
}

// TableName 轮播图 Banner自定义表名 shop_banner
func (Banner) TableName() string {
	return "shop_banner"
}
