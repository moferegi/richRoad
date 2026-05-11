package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ModelCallLog 模型调用日志
type ModelCallLog struct {
	global.GVA_MODEL
	UserID          uint   `json:"userID" form:"userID" gorm:"column:user_id;default:0;index;comment:用户ID;"`
	TaskID          uint   `json:"taskID" form:"taskID" gorm:"column:task_id;default:0;index;comment:试衣任务ID;"`
	TaskNo          string `json:"taskNo" form:"taskNo" gorm:"column:task_no;size:40;index;comment:试衣任务编号;"`
	RequestID       string `json:"requestID" form:"requestID" gorm:"column:request_id;size:80;index;comment:请求幂等ID;"`
	SceneType       string `json:"sceneType" form:"sceneType" gorm:"column:scene_type;size:20;index;comment:场景类型;"`
	TemplatePart    string `json:"templatePart" form:"templatePart" gorm:"column:template_part;size:20;comment:模板部位;"`
	EntrySource     string `json:"entrySource" form:"entrySource" gorm:"column:entry_source;size:80;index;comment:入口来源;"`
	Behavior        string `json:"behavior" form:"behavior" gorm:"column:behavior;size:80;index;comment:行为标识;"`
	CallStage       string `json:"callStage" form:"callStage" gorm:"column:call_stage;size:80;index;comment:调用阶段;"`
	ModelKey        string `json:"modelKey" form:"modelKey" gorm:"column:model_key;size:120;index;comment:模型key;"`
	ModelUsage      string `json:"modelUsage" form:"modelUsage" gorm:"column:model_usage;size:20;index;comment:模型用途;"`
	ModelName       string `json:"modelName" form:"modelName" gorm:"column:model_name;size:120;comment:模型名称;"`
	Provider        string `json:"provider" form:"provider" gorm:"column:provider;size:80;index;comment:模型提供商;"`
	EndpointURL     string `json:"endpointURL" form:"endpointURL" gorm:"column:endpoint_url;size:500;comment:调用地址;"`
	QueryURL        string `json:"queryURL" form:"queryURL" gorm:"column:query_url;size:500;comment:查询地址;"`
	RelatedModels   string `json:"relatedModels" form:"relatedModels" gorm:"column:related_models;type:text;comment:协同模型(JSON);"`
	SourceImage     string `json:"sourceImage" form:"sourceImage" gorm:"column:source_image;size:1024;comment:输入人物图;"`
	TemplateImage   string `json:"templateImage" form:"templateImage" gorm:"column:template_image;size:1024;comment:输入服饰图;"`
	ResultImage     string `json:"resultImage" form:"resultImage" gorm:"column:result_image;size:1024;comment:输出结果图;"`
	RequestPayload  string `json:"requestPayload" form:"requestPayload" gorm:"column:request_payload;type:longtext;comment:请求参数(JSON);"`
	ResponsePayload string `json:"responsePayload" form:"responsePayload" gorm:"column:response_payload;type:longtext;comment:响应参数(JSON);"`
	Status          string `json:"status" form:"status" gorm:"column:status;size:32;index;comment:调用状态;"`
	ErrorMessage    string `json:"errorMessage" form:"errorMessage" gorm:"column:error_message;size:1000;comment:错误信息;"`
	CostPoints      int    `json:"costPoints" form:"costPoints" gorm:"column:cost_points;default:0;comment:该次模型调用扣币;"`
	RefundPoints    int    `json:"refundPoints" form:"refundPoints" gorm:"column:refund_points;default:0;comment:该次模型调用退币;"`
	DurationMs      int64  `json:"durationMs" form:"durationMs" gorm:"column:duration_ms;default:0;comment:耗时毫秒;"`
}

func (ModelCallLog) TableName() string {
	return "client_model_call_log"
}
