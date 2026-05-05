package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// CreateTryonTaskReq 创建试衣任务请求
type CreateTryonTaskReq struct {
	RequestID     string `json:"requestID" form:"requestID"`
	ModelKey      string `json:"modelKey" form:"modelKey"`
	SceneType     string `json:"sceneType" form:"sceneType" binding:"required,oneof=clothes shoes takeoff"`
	SourceImage   string `json:"sourceImage" form:"sourceImage" binding:"required"`
	TemplateImage string `json:"templateImage" form:"templateImage" binding:"required"`
}

// TryonTaskSearch 试衣任务查询条件
type TryonTaskSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	TrendDays      int        `json:"trendDays" form:"trendDays"`
	UserID         uint       `json:"userID" form:"userID"`
	TaskNo         string     `json:"taskNo" form:"taskNo"`
	RequestID      string     `json:"requestID" form:"requestID"`
	Status         string     `json:"status" form:"status"`
	SceneType      string     `json:"sceneType" form:"sceneType"`
	request.PageInfo
}
