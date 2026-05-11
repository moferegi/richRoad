package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

// ModelCallLogSearch 模型调用日志查询条件
type ModelCallLogSearch struct {
	request.PageInfo
	StartCreatedAt string `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   string `json:"endCreatedAt" form:"endCreatedAt"`
	OnlyFailed     bool   `json:"onlyFailed" form:"onlyFailed"`
	UserID         uint   `json:"userID" form:"userID"`
	SceneType      string `json:"sceneType" form:"sceneType"`
	Behavior       string `json:"behavior" form:"behavior"`
	CallStage      string `json:"callStage" form:"callStage"`
	ModelKey       string `json:"modelKey" form:"modelKey"`
	ModelUsage     string `json:"modelUsage" form:"modelUsage"`
	Provider       string `json:"provider" form:"provider"`
	Status         string `json:"status" form:"status"`
	TaskNo         string `json:"taskNo" form:"taskNo"`
	RequestID      string `json:"requestID" form:"requestID"`
}
