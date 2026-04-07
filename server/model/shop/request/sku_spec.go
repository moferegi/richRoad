package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type SkuSpecSearch struct {
	Type           string     `json:"type" form:"type"`
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}
