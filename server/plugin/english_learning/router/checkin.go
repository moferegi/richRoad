package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/api"
	"github.com/gin-gonic/gin"
)

type CheckinRouter struct{}

func (s *CheckinRouter) InitCheckinRouter(Router *gin.RouterGroup) {
	checkinRouter := Router.Group("checkin")
	var checkinApi = api.ApiGroupApp.CheckinApi
	{
		checkinRouter.POST("do", checkinApi.DoCheckin)          // 打卡
		checkinRouter.POST("exchange", checkinApi.ExchangeTime) // 积分兑换时长
		checkinRouter.GET("getStats", checkinApi.GetStats)
		checkinRouter.GET("getPointRecordList", checkinApi.GetPointRecordList)
		checkinRouter.GET("getCheckinRecordList", checkinApi.GetCheckinRecordList)
	}
}
