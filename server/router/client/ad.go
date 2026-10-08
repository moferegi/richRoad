package client

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AdRouter struct{}

func (s *AdRouter) InitAdRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	adRouter := Router.Group("ad").Use(middleware.OperationRecord())
	adRouterWithoutRecord := Router.Group("ad")
	adPublicRouter := PublicRouter.Group("ad")

	var adApi = v1.ApiGroupApp.ClientApiGroup.AdApi
	{
		// 广告位置
		adRouter.POST("createAdPosition", adApi.CreateAdPosition)
		adRouter.DELETE("deleteAdPosition", adApi.DeleteAdPosition)
		adRouter.PUT("updateAdPosition", adApi.UpdateAdPosition)
		// 广告视频
		adRouter.POST("createAdVideo", adApi.CreateAdVideo)
		adRouter.DELETE("deleteAdVideo", adApi.DeleteAdVideo)
		adRouter.PUT("updateAdVideo", adApi.UpdateAdVideo)
		// 切片上传
		adRouter.POST("sliceAdVideo", adApi.SliceAdVideo)
	}
	{
		// 查询类
		adRouterWithoutRecord.GET("getAdPositionList", adApi.GetAdPositionList)
		adRouterWithoutRecord.GET("getAdVideoList", adApi.GetAdVideoList)
		adRouterWithoutRecord.GET("checkFfmpeg", adApi.CheckFfmpeg)
		adRouterWithoutRecord.GET("getWatchRecordList", adApi.GetWatchRecordList)
		adRouterWithoutRecord.GET("getWatchStats", adApi.GetWatchStats)
	}
	{
		// Uni端公开接口
		adPublicRouter.GET("getAdByPosition", adApi.GetAdByPosition)
		adPublicRouter.POST("reportWatch", adApi.ReportWatch)
	}
}