package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type ExaAttachmentCategorySearch struct {
	ClassId  int    `json:"classId" form:"classId"`
	Position string `json:"position" form:"position"`
	request.PageInfo
}

type CreateScanUploadTicketRequest struct {
	ClassId        int    `json:"classId" form:"classId"`
	Folder         string `json:"folder" form:"folder" binding:"omitempty,max=256"`
	UploadType     string `json:"uploadType" form:"uploadType" binding:"omitempty,max=64"`
	UploadPosition string `json:"uploadPosition" form:"uploadPosition" binding:"omitempty,max=64"`
}

type SignURLRequest struct {
	FilePath string `json:"filePath" binding:"required,max=2048"`
}
