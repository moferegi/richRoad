package api

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	learningMiddleware "github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/service"
	serverUtils "github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ContentApi struct{}

var contentService = service.ServiceGroupApp.ContentService
var contentAuthzService = service.ServiceGroupApp.LearningAuthzService

func parseQueryID(c *gin.Context) (uint, bool) {
	idStr := c.Query("ID")
	if idStr == "" {
		response.FailWithMessage("ID不能为空", c)
		return 0, false
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage("ID参数错误", c)
		return 0, false
	}
	return uint(id), true
}

// CreateCategory 创建英语分类
// @Tags     EnglishContent
// @Summary  创建英语分类
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.EnglishCategory true "分类信息"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /englishLearning/content/createCategory [post]
func (a *ContentApi) CreateCategory(c *gin.Context) {
	var body model.EnglishCategory
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := contentService.CreateCategory(&body); err != nil {
		global.GVA_LOG.Error("创建英语分类失败", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateCategory 更新英语分类
// @Tags     EnglishContent
// @Summary  更新英语分类
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.EnglishCategory true "分类信息"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /englishLearning/content/updateCategory [put]
func (a *ContentApi) UpdateCategory(c *gin.Context) {
	var body model.EnglishCategory
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if body.ID == 0 {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	if err := contentService.UpdateCategory(body); err != nil {
		global.GVA_LOG.Error("更新英语分类失败", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteCategory 删除英语分类
// @Tags     EnglishContent
// @Summary  删除英语分类
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /englishLearning/content/deleteCategory [delete]
func (a *ContentApi) DeleteCategory(c *gin.Context) {
	id, ok := parseQueryID(c)
	if !ok {
		return
	}
	if err := contentService.DeleteCategory(id); err != nil {
		global.GVA_LOG.Error("删除英语分类失败", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// FindCategory 获取英语分类详情
// @Tags     EnglishContent
// @Summary  获取英语分类详情
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=model.EnglishCategory,msg=string} "获取成功"
// @Router   /englishLearning/content/findCategory [get]
func (a *ContentApi) FindCategory(c *gin.Context) {
	id, ok := parseQueryID(c)
	if !ok {
		return
	}
	data, err := contentService.GetCategory(id)
	if err != nil {
		global.GVA_LOG.Error("获取英语分类失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(data, c)
}

// GetCategoryList 分页获取英语分类
// @Tags     EnglishContent
// @Summary  分页获取英语分类
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.EnglishCategorySearch true "分页参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/content/getCategoryList [get]
func (a *ContentApi) GetCategoryList(c *gin.Context) {
	var pageInfo request.EnglishCategorySearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := contentService.GetCategoryList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取英语分类列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: serverUtils.LocalizeI18nPayloadByContext(c, list), Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize}, "获取成功", c)
}

// CreateChapter 创建英语章节
// @Tags     EnglishContent
// @Summary  创建英语章节
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.EnglishChapter true "章节信息"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /englishLearning/content/createChapter [post]
func (a *ContentApi) CreateChapter(c *gin.Context) {
	var body model.EnglishChapter
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := contentService.CreateChapter(&body); err != nil {
		global.GVA_LOG.Error("创建英语章节失败", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateChapter 更新英语章节
// @Tags     EnglishContent
// @Summary  更新英语章节
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.EnglishChapter true "章节信息"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /englishLearning/content/updateChapter [put]
func (a *ContentApi) UpdateChapter(c *gin.Context) {
	var body model.EnglishChapter
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if body.ID == 0 {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	if err := contentService.UpdateChapter(body); err != nil {
		global.GVA_LOG.Error("更新英语章节失败", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteChapter 删除英语章节
// @Tags     EnglishContent
// @Summary  删除英语章节
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /englishLearning/content/deleteChapter [delete]
func (a *ContentApi) DeleteChapter(c *gin.Context) {
	id, ok := parseQueryID(c)
	if !ok {
		return
	}
	if err := contentService.DeleteChapter(id); err != nil {
		global.GVA_LOG.Error("删除英语章节失败", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// FindChapter 获取英语章节详情
// @Tags     EnglishContent
// @Summary  获取英语章节详情
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=model.EnglishChapter,msg=string} "获取成功"
// @Router   /englishLearning/content/findChapter [get]
func (a *ContentApi) FindChapter(c *gin.Context) {
	id, ok := parseQueryID(c)
	if !ok {
		return
	}
	data, err := contentService.GetChapter(id)
	if err != nil {
		global.GVA_LOG.Error("获取英语章节失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(data, c)
}

// GetChapterList 分页获取英语章节
// @Tags     EnglishContent
// @Summary  分页获取英语章节
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.EnglishChapterSearch true "分页参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/content/getChapterList [get]
func (a *ContentApi) GetChapterList(c *gin.Context) {
	var pageInfo request.EnglishChapterSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := contentService.GetChapterList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取英语章节列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: serverUtils.LocalizeI18nPayloadByContext(c, list), Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize}, "获取成功", c)
}

// CreateVideoCategory 创建视频分类
// @Tags     EnglishContent
// @Summary  创建视频分类
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.VideoCategory true "视频分类信息"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /englishLearning/content/createVideoCategory [post]
func (a *ContentApi) CreateVideoCategory(c *gin.Context) {
	var body model.VideoCategory
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := contentService.CreateVideoCategory(&body); err != nil {
		global.GVA_LOG.Error("创建视频分类失败", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateVideoCategory 更新视频分类
// @Tags     EnglishContent
// @Summary  更新视频分类
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.VideoCategory true "视频分类信息"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /englishLearning/content/updateVideoCategory [put]
func (a *ContentApi) UpdateVideoCategory(c *gin.Context) {
	var body model.VideoCategory
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if body.ID == 0 {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	if err := contentService.UpdateVideoCategory(body); err != nil {
		global.GVA_LOG.Error("更新视频分类失败", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteVideoCategory 删除视频分类
// @Tags     EnglishContent
// @Summary  删除视频分类
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /englishLearning/content/deleteVideoCategory [delete]
func (a *ContentApi) DeleteVideoCategory(c *gin.Context) {
	id, ok := parseQueryID(c)
	if !ok {
		return
	}
	if err := contentService.DeleteVideoCategory(id); err != nil {
		global.GVA_LOG.Error("删除视频分类失败", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// FindVideoCategory 获取视频分类详情
// @Tags     EnglishContent
// @Summary  获取视频分类详情
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=model.VideoCategory,msg=string} "获取成功"
// @Router   /englishLearning/content/findVideoCategory [get]
func (a *ContentApi) FindVideoCategory(c *gin.Context) {
	id, ok := parseQueryID(c)
	if !ok {
		return
	}
	data, err := contentService.GetVideoCategory(id)
	if err != nil {
		global.GVA_LOG.Error("获取视频分类失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(data, c)
}

// GetVideoCategoryList 分页获取视频分类
// @Tags     EnglishContent
// @Summary  分页获取视频分类
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.VideoCategorySearch true "分页参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/content/getVideoCategoryList [get]
func (a *ContentApi) GetVideoCategoryList(c *gin.Context) {
	var pageInfo request.VideoCategorySearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := contentService.GetVideoCategoryList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取视频分类列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: serverUtils.LocalizeI18nPayloadByContext(c, list), Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize}, "获取成功", c)
}

// CreateVideoSeries 创建视频剧集
// @Tags     EnglishContent
// @Summary  创建视频剧集
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.VideoSeries true "视频剧集信息"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /englishLearning/content/createVideoSeries [post]
func (a *ContentApi) CreateVideoSeries(c *gin.Context) {
	var body model.VideoSeries
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := contentService.CreateVideoSeries(&body); err != nil {
		global.GVA_LOG.Error("创建视频剧集失败", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateVideoSeries 更新视频剧集
// @Tags     EnglishContent
// @Summary  更新视频剧集
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.VideoSeries true "视频剧集信息"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /englishLearning/content/updateVideoSeries [put]
func (a *ContentApi) UpdateVideoSeries(c *gin.Context) {
	var body model.VideoSeries
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if body.ID == 0 {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	if err := contentService.UpdateVideoSeries(body); err != nil {
		global.GVA_LOG.Error("更新视频剧集失败", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteVideoSeries 删除视频剧集
// @Tags     EnglishContent
// @Summary  删除视频剧集
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /englishLearning/content/deleteVideoSeries [delete]
func (a *ContentApi) DeleteVideoSeries(c *gin.Context) {
	id, ok := parseQueryID(c)
	if !ok {
		return
	}
	if err := contentService.DeleteVideoSeries(id); err != nil {
		global.GVA_LOG.Error("删除视频剧集失败", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// FindVideoSeries 获取视频剧集详情
// @Tags     EnglishContent
// @Summary  获取视频剧集详情
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=model.VideoSeries,msg=string} "获取成功"
// @Router   /englishLearning/content/findVideoSeries [get]
func (a *ContentApi) FindVideoSeries(c *gin.Context) {
	id, ok := parseQueryID(c)
	if !ok {
		return
	}
	data, err := contentService.GetVideoSeries(id)
	if err != nil {
		global.GVA_LOG.Error("获取视频剧集失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(serverUtils.LocalizeI18nPayloadByContext(c, data), c)
}

// GetVideoSeriesList 分页获取视频剧集
// @Tags     EnglishContent
// @Summary  分页获取视频剧集
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.VideoSeriesSearch true "分页参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/content/getVideoSeriesList [get]
func (a *ContentApi) GetVideoSeriesList(c *gin.Context) {
	var pageInfo request.VideoSeriesSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := contentService.GetVideoSeriesList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取视频剧集列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: serverUtils.LocalizeI18nPayloadByContext(c, list), Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize}, "获取成功", c)
}

// CreateVideoEpisode 创建视频单集
// @Tags     EnglishContent
// @Summary  创建视频单集
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.VideoEpisode true "视频单集信息"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /englishLearning/content/createVideoEpisode [post]
func (a *ContentApi) CreateVideoEpisode(c *gin.Context) {
	var body model.VideoEpisode
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := contentService.CreateVideoEpisode(&body); err != nil {
		global.GVA_LOG.Error("创建视频单集失败", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateVideoEpisode 更新视频单集
// @Tags     EnglishContent
// @Summary  更新视频单集
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.VideoEpisode true "视频单集信息"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /englishLearning/content/updateVideoEpisode [put]
func (a *ContentApi) UpdateVideoEpisode(c *gin.Context) {
	var body model.VideoEpisode
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if body.ID == 0 {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	if err := contentService.UpdateVideoEpisode(body); err != nil {
		global.GVA_LOG.Error("更新视频单集失败", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteVideoEpisode 删除视频单集
// @Tags     EnglishContent
// @Summary  删除视频单集
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /englishLearning/content/deleteVideoEpisode [delete]
func (a *ContentApi) DeleteVideoEpisode(c *gin.Context) {
	id, ok := parseQueryID(c)
	if !ok {
		return
	}
	if err := contentService.DeleteVideoEpisode(id); err != nil {
		global.GVA_LOG.Error("删除视频单集失败", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// FindVideoEpisode 获取视频单集详情
// @Tags     EnglishContent
// @Summary  获取视频单集详情
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=model.VideoEpisode,msg=string} "获取成功"
// @Router   /englishLearning/content/findVideoEpisode [get]
func (a *ContentApi) FindVideoEpisode(c *gin.Context) {
	id, ok := parseQueryID(c)
	if !ok {
		return
	}
	data, err := contentService.GetVideoEpisode(id)
	if err != nil {
		global.GVA_LOG.Error("获取视频单集失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	data.VideoUrl = service.SignLearningVideoURL(data.VideoUrl)
	hasFullAuth := false
	if userIDValue, exists := c.Get(learningMiddleware.CtxEnglishUserID); exists {
		if userID, castOK := userIDValue.(uint); castOK && userID > 0 {
			decision, decisionErr := contentAuthzService.EvaluateVideoEpisodeAccess(userID, data)
			if decisionErr != nil {
				global.GVA_LOG.Warn("评估视频单集权限失败，回退到全局权限", zap.Error(decisionErr), zap.Uint("episodeID", data.ID), zap.Uint("userID", userID))
			} else {
				hasFullAuth = decision.HasFullAuth
			}
		}
	}

	if !hasFullAuth {
		if value, exists := c.Get(learningMiddleware.CtxEnglishHasFullAuth); exists {
			if allowed, castOK := value.(bool); castOK {
				hasFullAuth = allowed
			}
		}
	}

	data.HasFullAuth = hasFullAuth
	response.OkWithData(serverUtils.LocalizeI18nPayloadByContext(c, data), c)
}

// GetVideoEpisodeList 分页获取视频单集
// @Tags     EnglishContent
// @Summary  分页获取视频单集
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.VideoEpisodeSearch true "分页参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/content/getVideoEpisodeList [get]
func (a *ContentApi) GetVideoEpisodeList(c *gin.Context) {
	var pageInfo request.VideoEpisodeSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := contentService.GetVideoEpisodeList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取视频单集列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	for i := range list {
		list[i].VideoUrl = service.SignLearningVideoURL(list[i].VideoUrl)
	}
	response.OkWithDetailed(response.PageResult{List: serverUtils.LocalizeI18nPayloadByContext(c, list), Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize}, "获取成功", c)
}
