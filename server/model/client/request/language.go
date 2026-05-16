package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type SysLanguageSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
}

type TranslateI18nRequest struct {
	Text    string   `json:"text" form:"text" binding:"required"`
	Source  string   `json:"source" form:"source"`
	Targets []string `json:"targets" form:"targets" binding:"required,min=1"`
}
