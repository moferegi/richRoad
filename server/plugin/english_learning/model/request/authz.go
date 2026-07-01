package request

import (
	"time"

	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type GrantEntitlementReq struct {
	UserID       uint       `json:"userId" binding:"required"`
	ResourceType string     `json:"resourceType" binding:"required"`
	ResourceID   uint       `json:"resourceId" binding:"required"`
	ExpireAt     *time.Time `json:"expireAt"`
	Remark       string     `json:"remark"`
}

type RevokeEntitlementReq struct {
	UserID       uint   `json:"userId" binding:"required"`
	ResourceType string `json:"resourceType" binding:"required"`
	ResourceID   uint   `json:"resourceId" binding:"required"`
}

type EntitlementSearch struct {
	UserID       uint   `json:"userId" form:"userId"`
	ResourceType string `json:"resourceType" form:"resourceType"`
	commonReq.PageInfo
}
