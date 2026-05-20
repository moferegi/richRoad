package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type SysErrorSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	Form           *string     `json:"form" form:"form"`
	Info           *string     `json:"info" form:"info"`
	request.PageInfo
}

// CreateSysErrorRequest 公开错误上报请求体（限制字段，避免客户端直接写入状态/解决方案）
type CreateSysErrorRequest struct {
	Form  string `json:"form" form:"form" binding:"required,max=256"`
	Info  string `json:"info" form:"info" binding:"required,max=4096"`
	Level string `json:"level" form:"level" binding:"omitempty,max=32"`
}
