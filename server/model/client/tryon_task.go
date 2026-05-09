package client

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TryonTask 试衣任务 结构体
type TryonTask struct {
	global.GVA_MODEL
	UserID         uint       `json:"userID" form:"userID" gorm:"column:user_id;index;index:idx_tryon_user_request,unique;comment:用户ID;"`
	RequestID      string     `json:"requestID" form:"requestID" gorm:"column:request_id;size:64;index:idx_tryon_user_request,unique;comment:客户端幂等请求ID;"`
	TaskNo         string     `json:"taskNo" form:"taskNo" gorm:"column:task_no;size:64;uniqueIndex;comment:任务编号;"`
	SceneType      string     `json:"sceneType" form:"sceneType" gorm:"column:scene_type;size:20;index;comment:场景类型(clothes/shoes/takeoff);"`
	Status         string     `json:"status" form:"status" gorm:"column:status;size:20;index;comment:任务状态(processing/success/failed);"`
	SourceImage    string     `json:"sourceImage" form:"sourceImage" gorm:"column:source_image;type:text;comment:用户原图;"`
	TemplateImage  string     `json:"templateImage" form:"templateImage" gorm:"column:template_image;type:text;comment:服饰模板图;"`
	ResultImage    string     `json:"resultImage" form:"resultImage" gorm:"column:result_image;type:text;comment:生成结果图;"`
	Provider       string     `json:"provider" form:"provider" gorm:"column:provider;size:50;comment:模型服务提供方;"`
	ProviderTaskID string     `json:"providerTaskID" form:"providerTaskID" gorm:"column:provider_task_id;size:128;index;comment:第三方异步任务ID;"`
	EnableRefiner  bool       `json:"enableRefiner" form:"enableRefiner" gorm:"column:enable_refiner;default:false;comment:是否开启图片精修;"`
	RefinerStatus  string     `json:"refinerStatus" form:"refinerStatus" gorm:"column:refiner_status;size:20;default:'';comment:精修状态(disabled/pending/processing/success/failed);"`
	RefinerTaskID  string     `json:"refinerTaskID" form:"refinerTaskID" gorm:"column:refiner_task_id;size:128;index;comment:精修任务ID;"`
	BeautifyStatus string     `json:"beautifyStatus" form:"beautifyStatus" gorm:"column:beautify_status;size:20;default:'';comment:美肤状态(disabled/processing/success/failed);"`
	BeautifyTaskNo string     `json:"beautifyTaskNo" form:"beautifyTaskNo" gorm:"column:beautify_task_no;size:64;index;comment:美肤任务编号;"`
	BeautifyResult string     `json:"beautifyResult" form:"beautifyResult" gorm:"column:beautify_result;type:text;comment:美肤结果图;"`
	BeautifyCost   int        `json:"beautifyCost" form:"beautifyCost" gorm:"column:beautify_cost;default:0;comment:美肤扣除试衣币;"`
	BeautifyRefund int        `json:"beautifyRefund" form:"beautifyRefund" gorm:"column:beautify_refund;default:0;comment:美肤退还试衣币;"`
	BeautifyError  string     `json:"beautifyError" form:"beautifyError" gorm:"column:beautify_error;size:500;comment:美肤失败原因;"`
	BeautifyAt     *time.Time `json:"beautifyAt" form:"beautifyAt" gorm:"column:beautify_at;comment:美肤完成时间;"`
	CostPoints     int        `json:"costPoints" form:"costPoints" gorm:"column:cost_points;default:0;comment:扣除试衣币;"`
	RefundPoints   int        `json:"refundPoints" form:"refundPoints" gorm:"column:refund_points;default:0;comment:退还试衣币;"`
	ErrorMessage   string     `json:"errorMessage" form:"errorMessage" gorm:"column:error_message;size:500;comment:失败原因;"`
	CompletedAt    *time.Time `json:"completedAt" form:"completedAt" gorm:"column:completed_at;comment:完成时间;"`
	RefundedAt     *time.Time `json:"refundedAt" form:"refundedAt" gorm:"column:refunded_at;comment:退币时间;"`
}

// TableName 试衣任务自定义表名
func (TryonTask) TableName() string {
	return "client_tryon_task"
}
