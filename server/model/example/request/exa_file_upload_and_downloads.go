package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type ExaAttachmentCategorySearch struct {
	ClassId  int    `json:"classId" form:"classId"`
	Position string `json:"position" form:"position"`
	request.PageInfo
}
