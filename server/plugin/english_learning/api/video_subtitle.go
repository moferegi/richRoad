package api

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
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

// ScanKeywords
// @Tags     VideoSubtitle
// @Summary  扫描英文字幕提取关键词列表（管理员多选后确认解析）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.ScanKeywordsReq true "单集ID和英文字幕URL"
// @Success  200  {object} response.Response{data=[]request.KeywordItem,msg=string} "扫描成功"
// @Router   /englishLearning/video/scanKeywords [post]
func (a *VideoSubtitleApi) ScanKeywords(c *gin.Context) {
	var req request.ScanKeywordsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	list, err := videoSubtitleService.ScanKeywords(req)
	if err != nil {
		global.GVA_LOG.Error("扫描关键词失败!", zap.Error(err))
		response.FailWithMessage("扫描失败: "+err.Error(), c)
		return
	}

	response.OkWithData(list, c)
}

// ParseSubtitleFiles
// @Tags     VideoSubtitle
// @Summary  解析字幕文件并入库（英文必填，其他语言可选；keywordIDs为管理员确认的词ID列表）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.ParseSubtitleFilesReq true "单集ID、英文字幕地址、多语言字幕地址与关键词ID列表"
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

// GetEpisodeKeywords
// @Tags     VideoSubtitle
// @Summary  获取某单集已高亮的重点单词列表（用于预览弹窗关键词管理）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param   episodeId query uint true "单集ID"
// @Success 200  {object} response.Response{data=[]request.KeywordItem,msg=string} "获取成功"
// @Router   /englishLearning/video/getEpisodeKeywords [get]
func (a *VideoSubtitleApi) GetEpisodeKeywords(c *gin.Context) {
	episodeIDStr := c.Query("episodeId")
	episodeID, err := strconv.ParseUint(episodeIDStr, 10, 64)
	if err != nil || episodeID == 0 {
		response.FailWithMessage("episodeId参数错误", c)
		return
	}

	list, err := videoSubtitleService.GetEpisodeKeywords(uint(episodeID))
	if err != nil {
		global.GVA_LOG.Error("获取重点单词失败", zap.Error(err))
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}

	response.OkWithData(list, c)
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

	response.OkWithData(utils.LocalizeI18nPayloadByContext(c, list), c)
}

// RehighlightSentences
// @Tags     VideoSubtitle
// @Summary  对已有字幕句子重新高亮（预览弹窗中调整关键词后即时应用）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.RehighlightSentencesReq true "单集ID与关键词ID列表"
// @Success  200  {object} response.Response{msg=string} "重新高亮成功"
// @Router   /englishLearning/video/rehighlightSentences [put]
func (a *VideoSubtitleApi) RehighlightSentences(c *gin.Context) {
	var req request.RehighlightSentencesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := videoSubtitleService.RehighlightSentences(req); err != nil {
		global.GVA_LOG.Error("重新高亮字幕失败!", zap.Error(err))
		response.FailWithMessage("重新高亮失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("字幕重新高亮成功", c)
}

// UpdateSentenceList
// @Tags     VideoSubtitle
// @Summary  批量更新单集字幕句子
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.UpdateVideoSentenceListReq true "单集ID与字幕句子列表"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /englishLearning/video/updateSentenceList [put]
func (a *VideoSubtitleApi) UpdateSentenceList(c *gin.Context) {
	var req request.UpdateVideoSentenceListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := videoSubtitleService.UpdateSentenceList(req); err != nil {
		global.GVA_LOG.Error("更新字幕句子失败", zap.Error(err))
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("字幕更新成功", c)
}
