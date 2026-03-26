
// 自动生成模板Kefu
package shop
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 客服 结构体  Kefu
type Kefu struct {
    global.GVA_MODEL
  Name  *string `json:"name" form:"name" gorm:"comment:客服姓名;column:name;" binding:"required"`  //姓名
  Avatar  string `json:"avatar" form:"avatar" gorm:"comment:客服头像;column:avatar;"`  //头像
  Status  string `json:"status" form:"status" gorm:"comment:客服状态;column:status;type:enum('在线','离线','忙碌');" binding:"required"`  //状态
  Link  *string `json:"link" form:"link" gorm:"column:link;"`  //链接
}


// TableName 客服 Kefu自定义表名 kefu
func (Kefu) TableName() string {
    return "kefu"
}





