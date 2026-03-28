package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type SysConfigSearch struct {
	ConfigGroup string `json:"configGroup" form:"configGroup"`
	ConfigKey   string `json:"configKey" form:"configKey"`
	request.PageInfo
}
