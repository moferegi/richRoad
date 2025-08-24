
// 自动生成模板PointRecord
package client
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 积分记录管理 结构体  PointRecord
type PointRecord struct {
    global.GVA_MODEL
  UserId  *int `json:"userId" form:"userId" gorm:"comment:用户ID;column:user_id;" binding:"required"`  //用户ID
  ChangeType  *string `json:"changeType" form:"changeType" gorm:"comment:增减类型;column:change_type;size:20;" binding:"required"`  //增减类型
  PointChange  *int `json:"pointChange" form:"pointChange" gorm:"comment:积分变化数量，正数为增加，负数为减少;column:point_change;" binding:"required"`  //积分变化
  OperationType  *string `json:"operationType" form:"operationType" gorm:"comment:操作类型;column:operation_type;size:50;" binding:"required"`  //操作类型
  Reason  *string `json:"reason" form:"reason" gorm:"comment:积分变化的具体原因;column:reason;size:200;" binding:"required"`  //变化原因
  CurrentPoints  *int `json:"currentPoints" form:"currentPoints" gorm:"comment:操作后的当前积分总数;column:current_points;"`  //当前积分
  RelatedOrderId  *int `json:"relatedOrderId" form:"relatedOrderId" gorm:"comment:关联的订单ID（如果有）;column:related_order_id;"`  //关联订单ID
  Remark  *string `json:"remark" form:"remark" gorm:"comment:备注信息;column:remark;size:500;"`  //备注
}


// TableName 积分记录管理 PointRecord自定义表名 client_point_records
func (PointRecord) TableName() string {
    return "client_point_records"
}





