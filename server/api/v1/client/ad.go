package client

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdApi struct{}

// --- 广告位置 ---

func (a *AdApi) CreateAdPosition(c *gin.Context) {
	var p client.AdPosition
	if err := c.ShouldBindJSON(&p); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := adService.CreateAdPosition(p); err != nil {
		global.GVA_LOG.Error("创建广告位置失败", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (a *AdApi) DeleteAdPosition(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := adService.DeleteAdPosition(req.ID); err != nil {
		global.GVA_LOG.Error("删除广告位置失败", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (a *AdApi) UpdateAdPosition(c *gin.Context) {
	var p client.AdPosition
	if err := c.ShouldBindJSON(&p); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := adService.UpdateAdPosition(p); err != nil {
		global.GVA_LOG.Error("更新广告位置失败", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (a *AdApi) GetAdPositionList(c *gin.Context) {
	var info request.AdPositionSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := adService.GetAdPositionList(info)
	if err != nil {
		global.GVA_LOG.Error("获取广告位置列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     info.Page,
		PageSize: info.PageSize,
	}, "获取成功", c)
}

// --- 广告视频 ---

func (a *AdApi) CreateAdVideo(c *gin.Context) {
	var v client.AdVideo
	if err := c.ShouldBindJSON(&v); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := adService.CreateAdVideo(v); err != nil {
		global.GVA_LOG.Error("创建广告视频失败", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (a *AdApi) DeleteAdVideo(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := adService.DeleteAdVideo(req.ID); err != nil {
		global.GVA_LOG.Error("删除广告视频失败", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (a *AdApi) UpdateAdVideo(c *gin.Context) {
	var v client.AdVideo
	if err := c.ShouldBindJSON(&v); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := adService.UpdateAdVideo(v); err != nil {
		global.GVA_LOG.Error("更新广告视频失败", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (a *AdApi) GetAdVideoList(c *gin.Context) {
	var info request.AdVideoSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := adService.GetAdVideoList(info)
	if err != nil {
		global.GVA_LOG.Error("获取广告视频列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     info.Page,
		PageSize: info.PageSize,
	}, "获取成功", c)
}

// --- 广告视频切片上传 ---

func (a *AdApi) SliceAdVideo(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage("请上传视频文件", c)
		return
	}

	folder := c.PostForm("folder")
	if folder == "" {
		folder = "ad/video"
	}

	result, err := adSliceService.SliceAndUploadAdVideo(file, folder)
	if err != nil {
		global.GVA_LOG.Error("广告视频切片上传失败", zap.Error(err))
		response.FailWithMessage("切片上传失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "切片上传成功", c)
}

func (a *AdApi) CheckFfmpeg(c *gin.Context) {
	if err := adSliceService.CheckFfmpeg(); err != nil {
		response.FailWithMessage("FFmpeg 不可用: "+err.Error(), c)
		return
	}
	response.OkWithMessage("FFmpeg 可用", c)
}

// --- Uni端接口 ---

func (a *AdApi) GetAdByPosition(c *gin.Context) {
	var req request.GetAdByPositionReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	videos, err := adService.GetAdByPosition(req.PositionKey)
	if err != nil {
		response.FailWithMessage("获取广告失败", c)
		return
	}

	// 拼接完整URL
	domain, _ := extDomainService.GetDefaultDomain()
	for i := range videos {
		if videos[i].MediaUrl != "" && domain != "" {
			videos[i].MediaUrl = fmt.Sprintf("%s/%s", domain, videos[i].MediaUrl)
		}
		if videos[i].ThumbnailUrl != "" && domain != "" {
			videos[i].ThumbnailUrl = fmt.Sprintf("%s/%s", domain, videos[i].ThumbnailUrl)
		}
	}

	response.OkWithData(videos, c)
}

func (a *AdApi) ReportWatch(c *gin.Context) {
	var req request.ReportWatchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	userID := utils.GetUserID(c)
	isCompleted := req.WatchedSeconds > 0

	record := client.AdWatchRecord{
		UserID:         userID,
		AdVideoID:      req.AdVideoID,
		PositionID:     req.PositionID,
		WatchedSeconds: req.WatchedSeconds,
		IsCompleted:    &isCompleted,
		IP:             c.ClientIP(),
	}

	if err := adService.CreateWatchRecord(record); err != nil {
		global.GVA_LOG.Error("上报观看记录失败", zap.Error(err))
		response.FailWithMessage("上报失败", c)
		return
	}
	response.OkWithMessage("上报成功", c)
}

// --- 观看记录 ---

func (a *AdApi) GetWatchRecordList(c *gin.Context) {
	var info request.AdWatchRecordSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := adService.GetWatchRecordList(info)
	if err != nil {
		global.GVA_LOG.Error("获取观看记录失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     info.Page,
		PageSize: info.PageSize,
	}, "获取成功", c)
}

// GetWatchStats 观看统计
func (a *AdApi) GetWatchStats(c *gin.Context) {
	type statItem struct {
		AdVideoID  uint   `json:"adVideoId"`
		Title      string `json:"title"`
		TotalViews int64  `json:"totalViews"`
	}

	var stats []statItem
	global.GVA_DB.Table("ad_watch_records r").
		Select("r.ad_video_id, v.title, COUNT(*) as total_views").
		Joins("LEFT JOIN ad_videos v ON v.id = r.ad_video_id").
		Group("r.ad_video_id, v.title").
		Order("total_views DESC").
		Scan(&stats)

	response.OkWithData(stats, c)
}
