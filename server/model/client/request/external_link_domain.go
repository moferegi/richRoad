package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type ExternalLinkDomainSearch struct {
	Name      string `json:"name" form:"name"`
	IsEnabled *bool  `json:"isEnabled" form:"isEnabled"`
	request.PageInfo
}
