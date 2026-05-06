package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// CreateTryonModelReq 创建我的模特请求
type CreateTryonModelReq struct {
	Name  string `json:"name" form:"name"`
	Image string `json:"image" form:"image" binding:"required"`
}

// UpdateTryonModelReq 重命名我的模特请求
type UpdateTryonModelReq struct {
	ID   uint   `json:"ID" form:"ID" binding:"required"`
	Name string `json:"name" form:"name" binding:"required"`
}

// TryonModelSearch 我的模特查询条件
type TryonModelSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	UserID         uint       `json:"userID" form:"userID"`
	Name           string     `json:"name" form:"name"`
	request.PageInfo
}
