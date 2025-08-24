
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type PointRecordSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
      UserId  *int `json:"userId" form:"userId"` 
      ChangeType  *string `json:"changeType" form:"changeType"` 
      PointChange  *int `json:"pointChange" form:"pointChange"` 
      OperationType  *string `json:"operationType" form:"operationType"` 
      Reason  *string `json:"reason" form:"reason"` 
      CurrentPoints  *int `json:"currentPoints" form:"currentPoints"` 
      RelatedOrderId  *int `json:"relatedOrderId" form:"relatedOrderId"` 
      Remark  *string `json:"remark" form:"remark"` 
    request.PageInfo
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
}
