package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type SkuSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	GoodID         *int       `json:"goodID" form:"goodID"`
	OrderBy        string     `json:"orderBy" form:"orderBy"`
	OrderDir       string     `json:"orderDir" form:"orderDir"`
	request.PageInfo
}
