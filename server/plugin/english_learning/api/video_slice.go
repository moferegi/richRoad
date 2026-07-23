package api

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var videoSliceService = service.ServiceGroupApp.VideoSliceService

// CheckFfmpeg 检查 FFmpeg 是否可用
// @Tags     EnglishContent
// @Summary  检查 FFmpeg 是否可用
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=map[string]interface{},msg=string} "检查结果"
// @Router   /englishLearning/content/checkFfmpeg [get]
func (a *ContentApi) CheckFfmpeg(c *gin.Context) {
	available, err := videoSliceService.CheckFfmpeg()
	if err != nil {
		response.OkWithDetailed(map[string]interface{}{
			"available": false,
			"message":   err.Error(),
		}, "FFmpeg 不可用", c)
		return
	}
	response.OkWithDetailed(map[string]interface{}{
		"available": available,
	}, "FFmpeg 可用", c)
}

// SliceVideoEpisode 视频切片为 HLS
// @Tags     EnglishContent
// @Summary  视频切片为 HLS (m3u8)，上传至默认云存储
// @Security ApiKeyAuth
// @accept   multipart/form-data
// @Produce  application/json
// @Param    file       formData file   true  "视频文件"
// @Param    folder     formData string false "云存储目录(如 didi)，不传则自动生成"
// @Param    seriesId   formData int    true  "所属剧集ID"
// @Param    name       formData string true  "单集名称(JSON多语言)"
// @Param    trialPercent formData int  false "试看比例(%)"
// @Param    sort       formData int    false "排序"
// @Param    episodeId  formData int    false "已有单集ID(更新模式)"
// @Success  200  {object} response.Response{data=model.VideoEpisode,msg=string} "切片成功"
// @Router   /englishLearning/content/sliceVideoEpisode [post]
func (a *ContentApi) SliceVideoEpisode(c *gin.Context) {
	// 1. 检查 FFmpeg
	available, checkErr := videoSliceService.CheckFfmpeg()
	if !available {
		response.FailWithMessage("FFmpeg 不可用: "+checkErr.Error(), c)
		return
	}

	// 2. 解析表单参数
	seriesIDStr := c.PostForm("seriesId")
	seriesID, err := strconv.ParseUint(seriesIDStr, 10, 64)
	if err != nil || seriesID == 0 {
		response.FailWithMessage("seriesId 参数错误", c)
		return
	}

	nameJSON := strings.TrimSpace(c.PostForm("name"))
	if nameJSON == "" {
		response.FailWithMessage("name 参数不能为空", c)
		return
	}

	trialPercent := 8
	if tpStr := c.PostForm("trialPercent"); tpStr != "" {
		if tp, err := strconv.Atoi(tpStr); err == nil && tp >= 1 && tp <= 100 {
			trialPercent = tp
		}
	}

	sortOrder := 0
	if sortStr := c.PostForm("sort"); sortStr != "" {
		if s, err := strconv.Atoi(sortStr); err == nil {
			sortOrder = s
		}
	}

	var existingEpisodeID uint
	if epIDStr := c.PostForm("episodeId"); epIDStr != "" {
		if epID, err := strconv.ParseUint(epIDStr, 10, 64); err == nil && epID > 0 {
			existingEpisodeID = uint(epID)
		}
	}

	folder := strings.TrimSpace(c.PostForm("folder"))

	// 3. 接收上传文件
	file, header, upErr := c.Request.FormFile("file")
	if upErr != nil {
		response.FailWithMessage("请上传视频文件", c)
		return
	}
	defer file.Close()

	// 4. 保存视频到 uploads/video/ 目录
	uploadDir := filepath.Join("uploads", "video")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.FailWithMessage("创建上传目录失败: "+err.Error(), c)
		return
	}

	// 保留原始文件名（加 seriesId 前缀避免冲突）
	ext := filepath.Ext(header.Filename)
	safeName := fmt.Sprintf("hls_input_%d%s", seriesID, ext)
	localVideoPath := filepath.Join(uploadDir, safeName)

	dst, createErr := os.Create(localVideoPath)
	if createErr != nil {
		response.FailWithMessage("创建视频文件失败: "+createErr.Error(), c)
		return
	}
	defer dst.Close()
	if _, copyErr := io.Copy(dst, file); copyErr != nil {
		response.FailWithMessage("保存视频文件失败: "+copyErr.Error(), c)
		return
	}
	dst.Close()

	global.GVA_LOG.Info("HLS 切片",
		zap.String("video", localVideoPath),
		zap.String("prefix", folder))

	// 5. 构造云存储前缀
	var cloudPrefix string
	if folder != "" {
		cloudPrefix = strings.TrimRight(folder, "/")
	} else if existingEpisodeID > 0 {
		cloudPrefix = fmt.Sprintf("hls/%d", existingEpisodeID)
	} else {
		cloudPrefix = fmt.Sprintf("hls/series_%d_%s", seriesID, strings.ReplaceAll(header.Filename, ".", "_"))
	}

	// 6. 调用切片服务
	result, sliceErr := videoSliceService.SliceAndUpload(localVideoPath, cloudPrefix, global.GVA_CONFIG.Hls.GlobalKey, existingEpisodeID)
	if sliceErr != nil {
		global.GVA_LOG.Error("HLS 切片失败", zap.Error(sliceErr))
		response.FailWithMessage("视频切片失败: "+sliceErr.Error(), c)
		return
	}

	// 7. 创建或更新 VideoEpisode 记录
	if existingEpisodeID > 0 {
		episode := model.VideoEpisode{
			VideoUrl:  result.M3u8URL,
			Duration:  result.Duration,
			VideoType: "m3u8",
		}
		episode.ID = existingEpisodeID
		if err := contentService.UpdateVideoEpisode(episode); err != nil {
			global.GVA_LOG.Error("更新单集视频URL失败", zap.Error(err))
			response.FailWithMessage("切片成功但更新记录失败: "+err.Error(), c)
			return
		}
		response.OkWithMessage(fmt.Sprintf("HLS 切片成功，共 %d 个分段，时长 %.1f 秒", result.SegmentCount, result.Duration), c)
		return
	}

	// 新建单集
	episode := model.VideoEpisode{
		SeriesID:     uint(seriesID),
		Name:         nameJSON,
		VideoUrl:     result.M3u8URL,
		Duration:     result.Duration,
		VideoType:    "m3u8",
		TrialPercent: trialPercent,
		Sort:         sortOrder,
	}

	if err := contentService.CreateVideoEpisode(&episode); err != nil {
		global.GVA_LOG.Error("创建单集记录失败", zap.Error(err))
		response.FailWithMessage("切片成功但创建记录失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(episode, fmt.Sprintf("HLS 切片成功，共 %d 个分段，时长 %.1f 秒", result.SegmentCount, result.Duration), c)
}
