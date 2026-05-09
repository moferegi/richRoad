package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// CreateTryonClothReq 创建我的衣橱请求
type CreateTryonClothReq struct {
	Name     string `json:"name" form:"name"`
	Category string `json:"category" form:"category" binding:"required,oneof=upper lower onepiece shoes"`
	Image    string `json:"image" form:"image" binding:"required"`
}

// UpdateTryonClothReq 更新我的衣橱请求
type UpdateTryonClothReq struct {
	ID       uint   `json:"ID" form:"ID" binding:"required"`
	Name     string `json:"name" form:"name" binding:"required"`
	Category string `json:"category" form:"category" binding:"omitempty,oneof=upper lower onepiece shoes"`
}

// TryonClothSearch 我的衣橱查询条件
type TryonClothSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	UserID         uint       `json:"userID" form:"userID"`
	Name           string     `json:"name" form:"name"`
	Category       string     `json:"category" form:"category"`
	request.PageInfo
}
