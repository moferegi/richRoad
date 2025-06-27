package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type GoodSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	Title          string     `json:"title" form:"title"`
	CategoryID     int        `json:"categoryID" form:"categoryID"`
	Recommend      *bool      `json:"recommend" form:"recommend"`
	Keyword        string     `json:"keyword" form:"keyword"` // 新增关键词搜索字段
	request.PageInfo
}
