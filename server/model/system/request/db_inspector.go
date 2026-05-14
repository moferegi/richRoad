package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

// DBInspectorOverviewSearch 数据库巡检查询参数
type DBInspectorOverviewSearch struct {
	request.PageInfo
	Keyword string `json:"keyword" form:"keyword"`
}

// DBInspectorAutoFixReq 自动修复请求
type DBInspectorAutoFixReq struct {
	Target    string `json:"target" form:"target" binding:"required"`
	TableName string `json:"tableName" form:"tableName"`
}

// DBInspectorDeleteByRangeReq 按日期范围真删除请求
type DBInspectorDeleteByRangeReq struct {
	TableName          string `json:"tableName" form:"tableName" binding:"required"`
	TimeColumn         string `json:"timeColumn" form:"timeColumn"`
	StartTime          string `json:"startTime" form:"startTime" binding:"required"`
	EndTime            string `json:"endTime" form:"endTime" binding:"required"`
	ConfirmText        string `json:"confirmText" form:"confirmText" binding:"required"`
	DeleteRelatedFiles bool   `json:"deleteRelatedFiles" form:"deleteRelatedFiles"`
	MaxDeleteRows      int64  `json:"maxDeleteRows" form:"maxDeleteRows"`
}
