package api

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VideoSubtitleApi struct{}

// 获取单例Service
var videoSubtitleService = service.ServiceGroupApp.VideoSubtitleService

// ParseSubtitle
// @Tags     VideoSubtitle
// @Summary  预解析字幕结构，生成高亮与可被客户端直接渲染或收藏的句库
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.ParseSubtitleReq true "包含视频单集和JSON数组化的字幕片段"
// @Success  200  {object} response.Response{msg=string} "解析成功"
// @Router   /englishLearning/video/parseSubtitle [post]
func (a *VideoSubtitleApi) ParseSubtitle(c *gin.Context) {
	var req request.ParseSubtitleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := videoSubtitleService.ParseAndHighlightSubtitle(req.EpisodeID, req.Subtitles); err != nil {
		global.GVA_LOG.Error("解析并高亮字幕失败!", zap.Error(err))
		response.FailWithMessage("解析失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("字幕成功解析入库并应用高亮", c)
}

// ParseSubtitleFiles
// @Tags     VideoSubtitle
// @Summary  解析字幕文件并入库（英文必填，其他语言可选）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.ParseSubtitleFilesReq true "单集ID、英文字幕地址与多语言字幕地址"
// @Success  200  {object} response.Response{msg=string} "解析成功"
// @Router   /englishLearning/video/parseSubtitleFiles [post]
func (a *VideoSubtitleApi) ParseSubtitleFiles(c *gin.Context) {
	var req request.ParseSubtitleFilesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := videoSubtitleService.ParseAndHighlightSubtitleFromFiles(req); err != nil {
		global.GVA_LOG.Error("解析字幕文件失败!", zap.Error(err))
		response.FailWithMessage("解析失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("字幕文件解析入库成功", c)
}

// GetSentenceList
// @Tags     VideoSubtitle
// @Summary  获取某视频单集的字幕句子列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=object,msg=string} "获取成功"
// @Router   /englishLearning/video/getSentenceList [get]
func (a *VideoSubtitleApi) GetSentenceList(c *gin.Context) {
	episodeIDStr := c.Query("episodeId")
	episodeID, err := strconv.ParseUint(episodeIDStr, 10, 64)
	if err != nil || episodeID == 0 {
		response.FailWithMessage("episodeId参数错误", c)
		return
	}

	list, err := videoSubtitleService.GetSentenceList(uint(episodeID))
	if err != nil {
		global.GVA_LOG.Error("获取字幕句子失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithData(list, c)
}
