package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// CreateTryonTaskReq 创建试衣任务请求
type CreateTryonTaskReq struct {
	RequestID     string `json:"requestID" form:"requestID"`
	ModelKey      string `json:"modelKey" form:"modelKey"`
	EnableRefiner bool   `json:"enableRefiner" form:"enableRefiner"`
	SceneType     string `json:"sceneType" form:"sceneType" binding:"required,oneof=clothes shoes takeoff"`
	TemplatePart  string `json:"templatePart" form:"templatePart" binding:"omitempty,oneof=upper lower"`
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

// DeleteMyTryonTaskReq 删除我的试衣任务请求
type DeleteMyTryonTaskReq struct {
	TaskNo string `json:"taskNo" form:"taskNo" binding:"required"`
}

// ApplyTryonBeautifyReq 试衣任务智能美肤请求
type ApplyTryonBeautifyReq struct {
	TaskID           uint    `json:"taskID" form:"taskID" binding:"required"`
	BeautifyModelKey string  `json:"beautifyModelKey" form:"beautifyModelKey"`
	RetouchDegree    float64 `json:"retouchDegree" form:"retouchDegree"`
	WhiteningDegree  float64 `json:"whiteningDegree" form:"whiteningDegree"`
}
