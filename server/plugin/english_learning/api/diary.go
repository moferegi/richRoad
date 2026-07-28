package api

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	learningMiddleware "github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/service"
	serverUtils "github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DiaryApi struct{}

var diaryService = service.ServiceGroupApp.DiaryService

// CreateDiaryCategory 创建日记分类
// @Tags     DiaryContent
// @Summary  创建日记分类
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.DiaryCategory true "日记分类信息"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /englishLearning/diary/createDiaryCategory [post]
func (a *DiaryApi) CreateDiaryCategory(c *gin.Context) {
	var body model.DiaryCategory
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := diaryService.CreateDiaryCategory(&body); err != nil {
		global.GVA_LOG.Error("创建日记分类失败", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateDiaryCategory 更新日记分类
// @Tags     DiaryContent
// @Summary  更新日记分类
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.DiaryCategory true "日记分类信息"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /englishLearning/diary/updateDiaryCategory [put]
func (a *DiaryApi) UpdateDiaryCategory(c *gin.Context) {
	var body model.DiaryCategory
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if body.ID == 0 {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	if err := diaryService.UpdateDiaryCategory(body); err != nil {
		global.GVA_LOG.Error("更新日记分类失败", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteDiaryCategory 删除日记分类
// @Tags     DiaryContent
// @Summary  删除日记分类
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /englishLearning/diary/deleteDiaryCategory [delete]
func (a *DiaryApi) DeleteDiaryCategory(c *gin.Context) {
	id, ok := parseQueryID(c)
	if !ok {
		return
	}
	if err := diaryService.DeleteDiaryCategory(id); err != nil {
		global.GVA_LOG.Error("删除日记分类失败", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// FindDiaryCategory 获取日记分类详情
// @Tags     DiaryContent
// @Summary  获取日记分类详情
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=model.DiaryCategory,msg=string} "获取成功"
// @Router   /englishLearning/diary/findDiaryCategory [get]
func (a *DiaryApi) FindDiaryCategory(c *gin.Context) {
	id, ok := parseQueryID(c)
	if !ok {
		return
	}
	data, err := diaryService.GetDiaryCategory(id)
	if err != nil {
		global.GVA_LOG.Error("获取日记分类失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(data, c)
}

// GetDiaryCategoryList 分页获取日记分类
// @Tags     DiaryContent
// @Summary  分页获取日记分类
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.DiaryCategorySearch true "分页参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/diary/getDiaryCategoryList [get]
func (a *DiaryApi) GetDiaryCategoryList(c *gin.Context) {
	var pageInfo request.DiaryCategorySearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := diaryService.GetDiaryCategoryList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取日记分类列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: serverUtils.LocalizeI18nPayloadByContext(c, list), Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize}, "获取成功", c)
}

// CreateDiary 创建日记（含标签）
// @Tags     DiaryContent
// @Summary  创建日记
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.CreateDiaryReq true "日记信息"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /englishLearning/diary/createDiary [post]
func (a *DiaryApi) CreateDiary(c *gin.Context) {
	var body request.CreateDiaryReq
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	diary := &model.Diary{
		CategoryID:           body.CategoryID,
		Name:                 body.Name,
		AudioUs:              body.AudioUs,
		AudioUk:              body.AudioUk,
		Duration:             body.Duration,
		ImageI18n:            body.ImageI18n,
		ShowTranslate:        body.ShowTranslate,
		ShowEnglish:          body.ShowEnglish,
		TrialPercent:         body.TrialPercent,
		NeedVip:              body.NeedVip,
		EnglishSubtitleUrl:   body.EnglishSubtitleUrl,
		EnglishSubtitleUrlUk: body.EnglishSubtitleUrlUk,
		Sort:                 body.Sort,
	}
	if err := diaryService.CreateDiaryWithTags(diary, body.TagIds); err != nil {
		global.GVA_LOG.Error("创建日记失败", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateDiary 更新日记（含标签）
// @Tags     DiaryContent
// @Summary  更新日记
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.UpdateDiaryReq true "日记信息"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /englishLearning/diary/updateDiary [put]
func (a *DiaryApi) UpdateDiary(c *gin.Context) {
	var body request.UpdateDiaryReq
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if body.ID == 0 {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	diary := model.Diary{
		GVA_MODEL:            global.GVA_MODEL{ID: body.ID},
		CategoryID:           body.CategoryID,
		Name:                 body.Name,
		AudioUs:              body.AudioUs,
		AudioUk:              body.AudioUk,
		Duration:             body.Duration,
		ImageI18n:            body.ImageI18n,
		ShowTranslate:        body.ShowTranslate,
		ShowEnglish:          body.ShowEnglish,
		TrialPercent:         body.TrialPercent,
		NeedVip:              body.NeedVip,
		EnglishSubtitleUrl:   body.EnglishSubtitleUrl,
		EnglishSubtitleUrlUk: body.EnglishSubtitleUrlUk,
		Sort:                 body.Sort,
	}
	if err := diaryService.UpdateDiaryWithTags(diary, body.TagIds); err != nil {
		global.GVA_LOG.Error("更新日记失败", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteDiary 删除日记
// @Tags     DiaryContent
// @Summary  删除日记
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /englishLearning/diary/deleteDiary [delete]
func (a *DiaryApi) DeleteDiary(c *gin.Context) {
	id, ok := parseQueryID(c)
	if !ok {
		return
	}
	if err := diaryService.DeleteDiary(id); err != nil {
		global.GVA_LOG.Error("删除日记失败", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// FindDiary 获取日记详情
// @Tags     DiaryContent
// @Summary  获取日记详情
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=model.Diary,msg=string} "获取成功"
// @Router   /englishLearning/diary/findDiary [get]
func (a *DiaryApi) FindDiary(c *gin.Context) {
	id, ok := parseQueryID(c)
	if !ok {
		return
	}
	data, err := diaryService.GetDiary(id)
	if err != nil {
		global.GVA_LOG.Error("获取日记失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	hasFullAuth := false
	if data.NeedVip == nil || !*data.NeedVip {
		// 不需要会员的日记，全部可见
		hasFullAuth = true
	} else {
		if userIDValue, exists := c.Get(learningMiddleware.CtxEnglishUserID); exists {
			if userID, castOK := userIDValue.(uint); castOK && userID > 0 {
				decision, decisionErr := contentAuthzService.EvaluateDiaryAccess(userID, data)
				if decisionErr != nil {
					global.GVA_LOG.Warn("评估日记权限失败", zap.Error(decisionErr), zap.Uint("diaryID", data.ID), zap.Uint("userID", userID))
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
	}
	data.HasFullAuth = hasFullAuth

	tags, tagErr := diaryService.GetDiaryTags(data.ID)
	if tagErr != nil {
		global.GVA_LOG.Warn("获取日记标签失败", zap.Error(tagErr), zap.Uint("diaryID", data.ID))
	}

	type diaryWithTags struct {
		model.Diary
		Tags []clientModel.DiaryTag `json:"tags"`
	}
	response.OkWithData(serverUtils.LocalizeI18nPayloadByContext(c, diaryWithTags{Diary: data, Tags: tags}), c)
}

// GetDiaryList 分页获取日记
// @Tags     DiaryContent
// @Summary  分页获取日记
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.DiarySearch true "分页参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/diary/getDiaryList [get]
func (a *DiaryApi) GetDiaryList(c *gin.Context) {
	var pageInfo request.DiarySearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := diaryService.GetDiaryList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取日记列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: serverUtils.LocalizeI18nPayloadByContext(c, list), Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize}, "获取成功", c)
}

type diaryTagItem struct {
	model.Diary
	Tags []clientModel.DiaryTag `json:"tags"`
}

// GetDiaryListByTag 按分类+标签筛选日记
// @Tags     DiaryContent
// @Summary  按分类+标签筛选日记
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.DiaryTagFilterSearch true "筛选参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/diary/getDiaryListByTag [get]
func (a *DiaryApi) GetDiaryListByTag(c *gin.Context) {
	var pageInfo request.DiaryTagFilterSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := diaryService.GetDiaryListByTag(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("按标签获取日记列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	diaryIDs := make([]uint, len(list))
	for i, d := range list {
		diaryIDs[i] = d.ID
	}
	tagMap, _ := diaryService.GetDiaryTagsBatch(diaryIDs)

	result := make([]diaryTagItem, len(list))
	for i, d := range list {
		tags := tagMap[d.ID]
		if tags == nil {
			tags = []clientModel.DiaryTag{}
		}
		result[i] = diaryTagItem{Diary: d, Tags: tags}
	}
	response.OkWithDetailed(response.PageResult{List: serverUtils.LocalizeI18nPayloadByContext(c, result), Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize}, "获取成功", c)
}

// GetDiarySentenceList 获取日记字幕句子列表（移动端，含试看过滤）
// @Tags     DiaryContent
// @Summary  获取日记字幕句子列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=[]model.DiarySentence,msg=string} "获取成功"
// @Router   /englishLearning/diary/getSentenceList [get]
func (a *DiaryApi) GetDiarySentenceList(c *gin.Context) {
	diaryIDStr := c.Query("diaryId")
	if diaryIDStr == "" {
		response.FailWithMessage("diaryId不能为空", c)
		return
	}
	diaryID, err := strconv.ParseUint(diaryIDStr, 10, 64)
	if err != nil || diaryID == 0 {
		response.FailWithMessage("diaryId参数错误", c)
		return
	}

	// 获取日记信息以判断 NeedVip 和 TrialPercent
	diary, err := diaryService.GetDiary(uint(diaryID))
	if err != nil {
		global.GVA_LOG.Error("获取日记信息失败", zap.Error(err), zap.Uint("diaryID", uint(diaryID)))
		response.FailWithMessage("获取失败", c)
		return
	}

	// 评估用户权限
	hasFullAuth := false
	if diary.NeedVip == nil || !*diary.NeedVip {
		// 不需要会员的日记，全部可见
		hasFullAuth = true
	} else {
		if userIDValue, exists := c.Get(learningMiddleware.CtxEnglishUserID); exists {
			if userID, castOK := userIDValue.(uint); castOK && userID > 0 {
				decision, decisionErr := contentAuthzService.EvaluateDiaryAccess(userID, diary)
				if decisionErr != nil {
					global.GVA_LOG.Warn("评估日记权限失败", zap.Error(decisionErr), zap.Uint("diaryID", diary.ID), zap.Uint("userID", userID))
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
	}

	list, err := diaryService.GetDiarySentenceListWithAuth(uint(diaryID), hasFullAuth, diary.TrialPercent)
	if err != nil {
		global.GVA_LOG.Error("获取日记字幕句子列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(serverUtils.LocalizeI18nPayloadByContext(c, list), c)
}

// GetDiarySentenceListAll 获取日记字幕句子列表（Web后台，不过滤试看）
// @Tags     DiaryContent
// @Summary  获取日记字幕句子列表（全部）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=[]model.DiarySentence,msg=string} "获取成功"
// @Router   /englishLearning/diary/getSentenceListAll [get]
func (a *DiaryApi) GetDiarySentenceListAll(c *gin.Context) {
	diaryIDStr := c.Query("diaryId")
	if diaryIDStr == "" {
		response.FailWithMessage("diaryId不能为空", c)
		return
	}
	diaryID, err := strconv.ParseUint(diaryIDStr, 10, 64)
	if err != nil || diaryID == 0 {
		response.FailWithMessage("diaryId参数错误", c)
		return
	}

	list, err := diaryService.GetDiarySentenceList(uint(diaryID))
	if err != nil {
		global.GVA_LOG.Error("获取日记字幕句子列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(serverUtils.LocalizeI18nPayloadByContext(c, list), c)
}

// ParseDiarySubtitle 解析日记字幕（JSON手动输入模式）
// @Tags     DiaryContent
// @Summary  解析日记字幕
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.ParseDiarySubtitleReq true "字幕数据"
// @Success  200  {object} response.Response{msg=string} "解析成功"
// @Router   /englishLearning/diary/parseDiarySubtitle [post]
func (a *DiaryApi) ParseDiarySubtitle(c *gin.Context) {
	var body request.ParseDiarySubtitleReq
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if body.DiaryID == 0 {
		response.FailWithMessage("diaryId不能为空", c)
		return
	}
	if err := diaryService.ParseAndSaveDiarySentences(body.DiaryID, body.Subtitles, nil); err != nil {
		global.GVA_LOG.Error("解析日记字幕失败", zap.Error(err))
		response.FailWithMessage("解析失败", c)
		return
	}
	response.OkWithMessage("解析成功", c)
}

// GetDiaryKeywords 获取日记已高亮重点单词
// @Tags     DiaryContent
// @Summary  获取日记已高亮重点单词列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param   diaryId query uint true "日记ID"
// @Success  200  {object} response.Response{data=[]request.DiaryKeywordItem,msg=string} "获取成功"
// @Router   /englishLearning/diary/getKeywords [get]
func (a *DiaryApi) GetDiaryKeywords(c *gin.Context) {
	diaryIDStr := c.Query("diaryId")
	diaryID, err := strconv.ParseUint(diaryIDStr, 10, 64)
	if err != nil || diaryID == 0 {
		response.FailWithMessage("diaryId参数错误", c)
		return
	}
	items, err := diaryService.GetDiaryKeywords(uint(diaryID))
	if err != nil {
		global.GVA_LOG.Error("获取日记关键词失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(items, c)
}

// RehighlightDiarySentences 重新高亮日记字幕句子
// @Tags     DiaryContent
// @Summary  对已有日记字幕句子重新高亮
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.RehighlightDiarySentencesReq true "日记ID与关键词ID列表"
// @Success  200  {object} response.Response{msg=string} "重新高亮成功"
// @Router   /englishLearning/diary/rehighlightSentences [put]
func (a *DiaryApi) RehighlightDiarySentences(c *gin.Context) {
	var req request.RehighlightDiarySentencesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := diaryService.RehighlightDiarySentences(req); err != nil {
		global.GVA_LOG.Error("重新高亮日记字幕失败", zap.Error(err))
		response.FailWithMessage("重新高亮失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("重新高亮成功", c)
}

// UpdateDiarySentenceList 批量更新日记字幕句子
// @Tags     DiaryContent
// @Summary  批量更新日记字幕句子
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.UpdateDiarySentenceListReq true "日记ID与字幕句子列表"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /englishLearning/diary/updateSentenceList [put]
func (a *DiaryApi) UpdateDiarySentenceList(c *gin.Context) {
	var req request.UpdateDiarySentenceListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := diaryService.UpdateDiarySentenceList(req); err != nil {
		global.GVA_LOG.Error("更新日记字幕句子失败", zap.Error(err))
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// ScanDiaryKeywords 扫描日记字幕关键词
// @Tags     DiaryContent
// @Summary  扫描日记字幕文件提取关键词
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.ScanDiaryKeywordsReq true "扫描请求"
// @Success  200  {object} response.Response{data=[]request.DiaryKeywordItem,msg=string} "扫描成功"
// @Router   /englishLearning/diary/scanKeywords [post]
func (a *DiaryApi) ScanDiaryKeywords(c *gin.Context) {
	var body request.ScanDiaryKeywordsReq
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	items, err := diaryService.ScanDiaryKeywords(body)
	if err != nil {
		global.GVA_LOG.Error("扫描日记关键词失败", zap.Error(err))
		response.FailWithMessage("扫描失败: "+err.Error(), c)
		return
	}
	response.OkWithData(items, c)
}

// ParseDiarySubtitleFiles 从字幕文件解析日记字幕
// @Tags     DiaryContent
// @Summary  从字幕文件解析日记字幕
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.ParseDiarySubtitleFilesReq true "解析请求"
// @Success  200  {object} response.Response{msg=string} "解析成功"
// @Router   /englishLearning/diary/parseDiarySubtitleFiles [post]
func (a *DiaryApi) ParseDiarySubtitleFiles(c *gin.Context) {
	var body request.ParseDiarySubtitleFilesReq
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := diaryService.ParseAndSaveDiarySentencesFromFiles(body); err != nil {
		global.GVA_LOG.Error("解析日记字幕文件失败", zap.Error(err))
		response.FailWithMessage("解析失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("解析成功", c)
}
