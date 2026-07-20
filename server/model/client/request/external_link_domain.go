package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type ExternalLinkDomainSearch struct {
	Name      string `json:"name" form:"name"`
	IsEnabled *bool  `json:"isEnabled" form:"isEnabled"`
	CloudType string `json:"cloudType" form:"cloudType"`
	request.PageInfo
}

// CloudPingReq 云存储连接检测请求
type CloudPingReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

// CloudPingResp 云存储连接检测响应
type CloudPingResp struct {
	Success bool   `json:"success"`
	Latency int64  `json:"latency"` // 毫秒
	Message string `json:"message"`
}

// CloudListFilesReq 云存储文件列表请求
type CloudListFilesReq struct {
	ID       uint   `json:"id" form:"id" binding:"required"`
	Prefix   string `json:"prefix" form:"prefix"`
	Marker   string `json:"marker" form:"marker"`
	MaxKeys  int    `json:"maxKeys" form:"maxKeys"`
}

// CloudFileItem 云存储文件项
type CloudFileItem struct {
	Key          string `json:"key"`
	Size         int64  `json:"size"`
	LastModified string `json:"lastModified"`
	IsDir        bool   `json:"isDir"`
}

// CloudListFilesResp 云存储文件列表响应
type CloudListFilesResp struct {
	Prefix      string          `json:"prefix"`
	Files       []CloudFileItem `json:"files"`
	Marker      string          `json:"marker"`
	IsTruncated bool            `json:"isTruncated"`
}

// CloudCompareReq 多云目录比对请求
type CloudCompareReq struct {
	ReferenceID uint   `json:"referenceId" form:"referenceId" binding:"required"` // 参考基准云ID
	Prefix      string `json:"prefix" form:"prefix"`
}

// CloudCompareTarget 单个对比目标的比对结果
type CloudCompareTarget struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	CloudType   string            `json:"cloudType"`
	Missing     []CloudCompareFile `json:"missing"`      // 参考有、目标无
	Extra       []CloudCompareFile `json:"extra"`        // 目标有、参考无
	TotalRef    int               `json:"totalRef"`      // 参考云该目录下文件总数
	TotalTarget int               `json:"totalTarget"`   // 目标云该目录下文件总数
	Error       string            `json:"error,omitempty"` // 若获取失败
}

// CloudCompareFile 比对文件条目
type CloudCompareFile struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
}

// CloudCompareResp 多云目录比对响应
type CloudCompareResp struct {
	ReferenceName string              `json:"referenceName"`
	Targets       []CloudCompareTarget `json:"targets"`
}
