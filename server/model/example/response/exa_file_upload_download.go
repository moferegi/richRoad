package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/example"

type ExaFileResponse struct {
	File example.ExaFileUploadAndDownload `json:"file"`
}

type ScanUploadTicketResponse struct {
	Ticket    string `json:"ticket"`
	ExpiresAt int64  `json:"expiresAt"`
}
