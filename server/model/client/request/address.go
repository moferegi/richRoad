package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type AddressSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`

	Name     string `json:"name" form:"name" `
	UserID   uint   `json:"userID" form:"userID" `
	OrderBy  string `json:"orderBy" form:"orderBy"`
	OrderDir string `json:"orderDir" form:"orderDir"`
	request.PageInfo
}
