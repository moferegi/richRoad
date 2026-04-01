package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type GoodSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	Title          string     `json:"title" form:"title"`
	Status         *bool      `json:"status" form:"status"`
	CategoryID     int        `json:"categoryID" form:"categoryID"`
	Recommend      *bool      `json:"recommend" form:"recommend"`
	Keyword        string     `json:"keyword" form:"keyword"`       // 关键词搜索字段
	IsPresale      *bool      `json:"isPresale" form:"isPresale"`   // 预售商品筛选
	OrderBy        string     `json:"orderBy" form:"orderBy"`       // 排序字段: view_num/sale_num/collect_num/price
	OrderDir       string     `json:"orderDir" form:"orderDir"`     // 排序方向: asc/desc
	request.PageInfo
}
