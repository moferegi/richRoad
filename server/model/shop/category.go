// 自动生成模板Category
package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 商品分类 结构体  Category
type Category struct {
	global.GVA_MODEL
	ParentID         uint       `json:"parentId" form:"parentId" gorm:"column:parent_id;comment:父级分类ID;"`                                       //父级分类ID
	Title            string     `json:"title" form:"title" gorm:"column:title;type:text;comment:分类标题(JSON多语言);"`                                //分类标题(JSON多语言)
	Desc             string     `json:"desc" form:"desc" gorm:"column:desc;comment:分类描述;"`                                                      //分类描述
	Goods            []Good     `json:"goods" gorm:"foreignKey:CategoryID;references:ID"`                                                       //分类下的商品
	Children         []Category `json:"children" gorm:"-"`                                                                                      //子分类
	Icons            string     `json:"icons" form:"icons" gorm:"column:icons;comment:分类图标;"`                                                   //分类图标
	ExternalIconPath string     `json:"externalIconPath" form:"externalIconPath" gorm:"column:external_icon_path;size:500;comment:外部图标路径(优先);"` //外部图标路径
	ShowInUni        *bool      `json:"showInUni" form:"showInUni" gorm:"column:show_in_uni;default:true;comment:是否在uni端显示;"`                   //是否在uni端显示
}

// TableName 商品分类 Category自定义表名 shop_category
func (Category) TableName() string {
	return "shop_category"
}
