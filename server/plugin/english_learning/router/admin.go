package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/api"
	"github.com/gin-gonic/gin"
)

type AdminRouter struct{}

func (s *AdminRouter) InitAdminRouter(Router *gin.RouterGroup) {
	adminRouter := Router.Group("admin")
	var adminApi = api.ApiGroupApp.AdminApi
	{
		adminRouter.GET("getCheckinRecordList", adminApi.GetCheckinRecordList)
		adminRouter.GET("getPointRecordList", adminApi.GetPointRecordList)
		adminRouter.GET("getFreeTimeRecordList", adminApi.GetFreeTimeRecordList)
		adminRouter.GET("getWatchHistoryList", adminApi.GetWatchHistoryList)
		adminRouter.GET("getCollectionList", adminApi.GetCollectionList)
		adminRouter.GET("getWordErrorLogList", adminApi.GetWordErrorLogList)
		adminRouter.GET("getUserList", adminApi.GetUserList)
	}
}
