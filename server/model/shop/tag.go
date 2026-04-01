
// 自动生成模板Tag
package shop
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 标签 结构体  Tag
type Tag struct {
	global.GVA_MODEL
	Name        *string `json:"name" form:"name" gorm:"column:name;" binding:"required"`                                    //标签名
	NameI18n    string  `json:"nameI18n" form:"nameI18n" gorm:"column:name_i18n;type:text;comment:标签名多语言JSON;"`          //标签名多语言
	Description *string `json:"description" form:"description" gorm:"column:description;"`                                //描述
	Color       *string `json:"color" form:"color" gorm:"column:color;"`                                                  //颜色
}


// TableName 标签 Tag自定义表名 Tags
func (Tag) TableName() string {
    return "Tags"
}





