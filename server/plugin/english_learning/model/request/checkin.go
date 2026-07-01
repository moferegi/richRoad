package request

import commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type ExchangePointsReq struct {
	Points int `json:"points" binding:"required,min=1"`
}

type PointRecordSearch struct {
	commonReq.PageInfo
}

type CheckinRecordSearch struct {
	commonReq.PageInfo
}
