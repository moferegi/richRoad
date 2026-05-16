package shop

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Popup 弹窗管理 结构体
type Popup struct {
	global.GVA_MODEL
	Title        string     `json:"title" form:"title" gorm:"column:title;type:text;comment:弹窗标题(JSON多语言);"`                                //弹窗标题
	Image        string     `json:"image" form:"image" gorm:"column:image;size:500;comment:弹窗图片;"`                                          //弹窗图片
	ExternalPath string     `json:"externalPath" form:"externalPath" gorm:"column:external_path;size:500;comment:外部图片路径;"`                  //外部图片路径
	Link         string     `json:"link" form:"link" gorm:"column:link;size:500;comment:点击跳转链接;"`                                           //点击跳转链接
	Content      string     `json:"content" form:"content" gorm:"column:content;type:text;comment:富文本内容(JSON多语言);"`                         //富文本内容
	PopupType    string     `json:"popupType" form:"popupType" gorm:"column:popup_type;size:20;default:image;comment:弹窗类型(image/content);"` //弹窗类型
	Pages        string     `json:"pages" form:"pages" gorm:"column:pages;size:1000;comment:展示页面路径(逗号分隔);"`                                 //展示页面路径
	StartTime    *time.Time `json:"startTime" form:"startTime" gorm:"column:start_time;comment:生效时间;"`                                      //生效时间
	EndTime      *time.Time `json:"endTime" form:"endTime" gorm:"column:end_time;comment:失效时间;"`                                            //失效时间
	ClientType   string     `json:"clientType" form:"clientType" gorm:"column:client_type;size:20;default:all;comment:客户端类型(uni/web/all);"` //客户端类型
	OnceOnly     *bool      `json:"onceOnly" form:"onceOnly" gorm:"column:once_only;default:false;comment:是否只弹一次;"`                         //是否只弹一次
	Sort         *int       `json:"sort" form:"sort" gorm:"column:sort;default:0;comment:排序;"`                                              //排序
	Closeable    *bool      `json:"closeable" form:"closeable" gorm:"column:closeable;default:true;comment:是否可关闭;"`                         //是否可关闭
	Position     string     `json:"position" form:"position" gorm:"column:position;size:30;default:home;comment:展示位置(兼容旧版);"`               //展示位置(兼容)
	IsEnabled    *bool      `json:"isEnabled" form:"isEnabled" gorm:"column:is_enabled;default:true;comment:是否启用;"`                         //是否启用
}

// TableName 弹窗管理自定义表名
func (Popup) TableName() string {
	return "shop_popup"
}
