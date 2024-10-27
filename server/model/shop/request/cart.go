package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type CartSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`

	UserID uint `json:"userID" form:"userID" `
	request.PageInfo
}

type CartCreate struct {
	UserID uint `json:"userID" form:"userID"`
	GoodID uint `json:"goodID" form:"goodID" binding:"required"`
	SKUID  uint `json:"skuID" form:"skuID" binding:"required"`
}
