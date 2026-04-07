// 自动生成模板Banner
package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 轮播图 结构体  Banner
type Banner struct {
	global.GVA_MODEL
	Title        string `json:"title" form:"title" gorm:"column:title;comment:轮播标题;" binding:"required"`                   //轮播标题
	Src          string `json:"src" form:"src" gorm:"column:src;comment:图片连接;" binding:"required"`                         //图片连接
	Href         string `json:"href" form:"href" gorm:"column:href;comment:图片跳转链接;"`                                       //跳转链接
	ExternalPath string `json:"externalPath" form:"externalPath" gorm:"column:external_path;size:500;comment:外部图片路径(优先);"` //外部图片路径
	// === 遮罩设置 ===
	MaskEnabled   *bool  `json:"maskEnabled" form:"maskEnabled" gorm:"column:mask_enabled;default:false;comment:遮罩开关;"`                                      //遮罩开关
	MaskHeight    *int   `json:"maskHeight" form:"maskHeight" gorm:"column:mask_height;default:40;comment:遮罩高度px;"`                                          //遮罩高度
	MaskBgColor   string `json:"maskBgColor" form:"maskBgColor" gorm:"column:mask_bg_color;size:50;default:'rgba(0,0,0,0.5)';comment:遮罩背景色;"`                //遮罩背景色
	MaskText      string `json:"maskText" form:"maskText" gorm:"column:mask_text;type:text;comment:遮罩文字(JSON多语言);"`                                          //遮罩文字
	MaskTextColor string `json:"maskTextColor" form:"maskTextColor" gorm:"column:mask_text_color;size:20;default:#FFFFFF;comment:遮罩文字颜色;"`                   //遮罩文字颜色
	MaskTextSize  *int   `json:"maskTextSize" form:"maskTextSize" gorm:"column:mask_text_size;default:14;comment:遮罩文字大小;"`                                   //遮罩文字大小
	MaskTextAlign string `json:"maskTextAlign" form:"maskTextAlign" gorm:"column:mask_text_align;size:20;default:center;comment:遮罩文字对齐(left/center/right);"` //遮罩文字对齐
	Sort          *int   `json:"sort" form:"sort" gorm:"column:sort;default:0;comment:排序;"`                                                                  //排序
	IsEnabled     *bool  `json:"isEnabled" form:"isEnabled" gorm:"column:is_enabled;default:true;comment:是否启用;"`                                             //是否启用
	GoodID        *uint  `json:"goodID" form:"goodID" gorm:"column:good_id;comment:关联商品ID;"`                                                                 //关联商品ID
}

// TableName 轮播图 Banner自定义表名 shop_banner
func (Banner) TableName() string {
	return "shop_banner"
}
