package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type CategorySearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	ParentID       uint       `json:"parentId" form:"parentId"`
	Title          string     `json:"title" form:"title"`
	ShowInUni      *bool      `json:"showInUni" form:"showInUni"` // uni端筛选：仅显示showInUni=true的分类
	request.PageInfo
}
