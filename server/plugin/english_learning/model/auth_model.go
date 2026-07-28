package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

const (
	ResourceTypeEnglishCategory = "english_category"
	ResourceTypeVideoSeries     = "video_series"
	ResourceTypeVideoEpisode    = "video_episode"
	ResourceTypeDiary           = "diary"
)

// UserLearningEntitlement 用户学习资源授权表（支持分类/剧集/单集粒度）
type UserLearningEntitlement struct {
	global.GVA_MODEL
	UserID       uint       `json:"userId" gorm:"uniqueIndex:idx_user_resource_entitlement,priority:1;index;comment:用户ID"`
	ResourceType string     `json:"resourceType" gorm:"size:32;uniqueIndex:idx_user_resource_entitlement,priority:2;index;comment:资源类型(english_category/video_series/video_episode)"`
	ResourceID   uint       `json:"resourceId" gorm:"uniqueIndex:idx_user_resource_entitlement,priority:3;index;comment:资源ID"`
	Source       string     `json:"source" gorm:"size:32;default:manual_grant;comment:授权来源"`
	GrantedBy    uint       `json:"grantedBy" gorm:"default:0;comment:授权人ID"`
	ExpireAt     *time.Time `json:"expireAt" gorm:"index;comment:授权过期时间(为空表示永久)"`
	Remark       string     `json:"remark" gorm:"size:500;comment:备注"`
}
