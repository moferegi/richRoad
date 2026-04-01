package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type PopupSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	Position       string     `json:"position" form:"position"`
	ClientType     string     `json:"clientType" form:"clientType"`
	IsEnabled      *bool      `json:"isEnabled" form:"isEnabled"`
	request.PageInfo
}
