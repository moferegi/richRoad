package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// AdPosition 广告位置
type AdPosition struct {
	global.GVA_MODEL
	Name        string `json:"name" form:"name" gorm:"column:name;size:100;comment:位置名称;"`
	PositionKey string `json:"positionKey" form:"positionKey" gorm:"column:position_key;size:200;uniqueIndex;comment:位置标记(如pages/game/category);"`
	Description string `json:"description" form:"description" gorm:"column:description;size:500;comment:位置描述;"`
	IsEnabled   *bool  `json:"isEnabled" form:"isEnabled" gorm:"column:is_enabled;default:true;comment:是否启用;"`
	Sort        int    `json:"sort" form:"sort" gorm:"column:sort;default:0;comment:排序;"`
}

func (AdPosition) TableName() string {
	return "ad_positions"
}

// AdVideo 广告视频/图片
type AdVideo struct {
	global.GVA_MODEL
	PositionID      uint   `json:"positionId" form:"positionId" gorm:"column:position_id;index;comment:关联广告位置ID;"`
	Title           string `json:"title" form:"title" gorm:"column:title;size:200;comment:广告标题;"`
	MediaType       string `json:"mediaType" form:"mediaType" gorm:"column:media_type;size:20;default:'video';comment:媒体类型(video/image);"`
	MediaUrl        string `json:"mediaUrl" form:"mediaUrl" gorm:"column:media_url;size:500;comment:媒体资源相对路径;"`
	ThumbnailUrl    string `json:"thumbnailUrl" form:"thumbnailUrl" gorm:"column:thumbnail_url;size:500;comment:封面缩略图;"`
	Duration        int    `json:"duration" form:"duration" gorm:"column:duration;default:0;comment:总时长(秒);"`
	MinWatchSeconds int    `json:"minWatchSeconds" form:"minWatchSeconds" gorm:"column:min_watch_seconds;default:5;comment:最少观看秒数;"`
	IsEnabled       *bool  `json:"isEnabled" form:"isEnabled" gorm:"column:is_enabled;default:true;comment:是否启用;"`
	Sort            int    `json:"sort" form:"sort" gorm:"column:sort;default:0;comment:排序(同一位置下播放顺序);"`
	StorageKey      string `json:"storageKey" form:"storageKey" gorm:"column:storage_key;size:200;comment:云存储目录标识;"`
	UploadFolder    string `json:"uploadFolder" form:"uploadFolder" gorm:"column:upload_folder;size:200;comment:上传目录;"`
	LinkUrl         string `json:"linkUrl" form:"linkUrl" gorm:"column:link_url;size:500;comment:跳转链接;"`
}

func (AdVideo) TableName() string {
	return "ad_videos"
}

// AdWatchRecord 广告观看记录
type AdWatchRecord struct {
	global.GVA_MODEL
	UserID        uint   `json:"userId" form:"userId" gorm:"column:user_id;index;comment:用户ID;"`
	AdVideoID     uint   `json:"adVideoId" form:"adVideoId" gorm:"column:ad_video_id;index;comment:广告视频ID;"`
	PositionID    uint   `json:"positionId" form:"positionId" gorm:"column:position_id;index;comment:广告位置ID;"`
	WatchedSeconds int   `json:"watchedSeconds" form:"watchedSeconds" gorm:"column:watched_seconds;default:0;comment:实际观看秒数;"`
	IsCompleted   *bool  `json:"isCompleted" form:"isCompleted" gorm:"column:is_completed;default:false;comment:是否达到最低观看时长;"`
	IP            string `json:"ip" form:"ip" gorm:"column:ip;size:50;comment:用户IP;"`
}

func (AdWatchRecord) TableName() string {
	return "ad_watch_records"
}