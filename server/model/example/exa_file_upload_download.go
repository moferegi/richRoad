package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type ExaFileUploadAndDownload struct {
	global.GVA_MODEL
	CreatedBy    uint   `json:"createdBy" form:"createdBy" gorm:"column:created_by;default:0;index;comment:创建者ID"`
	Name         string `json:"name" form:"name" gorm:"column:name;comment:文件名"`                                // 文件名
	ClassId      int    `json:"classId" form:"classId" gorm:"default:0;type:int;column:class_id;comment:分类id;"` // 分类id
	Url          string `json:"url" form:"url" gorm:"column:url;comment:文件地址"`                                  // 文件地址
	ThumbnailURL string `json:"thumbnailUrl" form:"thumbnailUrl" gorm:"column:thumbnail_url;type:text;comment:缩略图地址"`
	Categories   string `json:"categories" form:"categories" gorm:"column:categories;size:64;comment:分类标签"`
	Position     string `json:"position" form:"position" gorm:"column:position;size:32;default:other;comment:上传位置"`
	Keywords     string `json:"keywords" form:"keywords" gorm:"column:keywords;type:text;comment:关键词"`
	Size         int64  `json:"size" form:"size" gorm:"column:size;type:bigint;default:0;comment:文件真实大小(字节)"`
	Tag          string `json:"tag" form:"tag" gorm:"column:tag;comment:文件标签"` // 文件标签
	Key          string `json:"key" form:"key" gorm:"column:key;comment:编号"`   // 编号
}

func (ExaFileUploadAndDownload) TableName() string {
	return "exa_file_upload_and_downloads"
}
