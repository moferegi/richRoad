package api

import (
	"errors"
	"io"
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EnglishWordApi struct{}

type EnglishWordDetailResponse struct {
	model.EnglishWord
	CategoryIDs []uint `json:"categoryIds"`
	ChapterIDs  []uint `json:"chapterIds"`
}

func normalizeWordListPage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return page, pageSize
}

// 获取单例Service
var englishWordService = service.ServiceGroupApp.EnglishWordService

func parseWordIDFromQuery(c *gin.Context) (uint, bool) {
	idStr := c.Query("ID")
	if idStr == "" {
		idStr = c.Query("id")
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage("ID参数错误", c)
		return 0, false
	}
	return uint(id), true
}

// CreateEnglishWord 创建单词库记录
// @Tags     EnglishWord
// @Summary  后台录入英文词汇，无音频时唤起 TTS 自动填色
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.CreateEnglishWordReq true "单词属性与章节挂载映射"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /englishLearning/word/create [post]
func (a *EnglishWordApi) CreateEnglishWord(c *gin.Context) {
	var req request.CreateEnglishWordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	wordEntity := model.EnglishWord{
		Word:         req.Word,
		PhoneticUS:   req.PhoneticUS, // 后续也可写脚本用 JSON 词典进行查漏补缺
		PhoneticUK:   req.PhoneticUK,
		PartOfSpeech: req.PartOfSpeech,
		AudioUS:      req.AudioUS,
		AudioUK:      req.AudioUK,
		Explanation:  req.Explanation,
	}

	if err := englishWordService.CreateWord(wordEntity, req.CategoryIDs, req.ChapterIDs, req.Sentences); err != nil {
		global.GVA_LOG.Error("录入英语单词失败!", zap.Error(err))
		response.FailWithMessage("录入失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("录入并绑定成功", c)
}

// FindEnglishWord 获取单词详情
// @Tags     EnglishWord
// @Summary  根据ID查询单词详情
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=model.EnglishWord,msg=string} "查询成功"
// @Router   /englishLearning/word/findWord [get]
func (a *EnglishWordApi) FindEnglishWord(c *gin.Context) {
	id, ok := parseWordIDFromQuery(c)
	if !ok {
		return
	}

	data, err := englishWordService.GetWord(uint(id))
	if err != nil {
		global.GVA_LOG.Error("查询单词详情失败", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	categoryIDs, chapterIDs, err := englishWordService.GetWordBindingIDs(uint(id))
	if err != nil {
		global.GVA_LOG.Error("查询单词绑定关系失败", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	sentences, err := englishWordService.GetWordSentences(uint(id))
	if err != nil {
		global.GVA_LOG.Error("查询单词造句失败", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	data.Sentences = sentences

	response.OkWithData(utils.LocalizeI18nPayloadByContext(c, EnglishWordDetailResponse{
		EnglishWord: data,
		CategoryIDs: categoryIDs,
		ChapterIDs:  chapterIDs,
	}), c)
}

// UpdateEnglishWord 更新单词信息
// @Tags     EnglishWord
// @Summary  更新英语单词信息（可选重建章节关联）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.UpdateEnglishWordReq true "更新参数"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /englishLearning/word/update [put]
func (a *EnglishWordApi) UpdateEnglishWord(c *gin.Context) {
	var req request.UpdateEnglishWordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := englishWordService.UpdateWord(req); err != nil {
		global.GVA_LOG.Error("更新英语单词失败", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// DeleteEnglishWord 删除单词
// @Tags     EnglishWord
// @Summary  删除英语单词及其关联数据
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /englishLearning/word/delete [delete]
func (a *EnglishWordApi) DeleteEnglishWord(c *gin.Context) {
	id, ok := parseWordIDFromQuery(c)
	if !ok {
		return
	}

	if err := englishWordService.DeleteWord(id); err != nil {
		global.GVA_LOG.Error("删除英语单词失败", zap.Error(err))
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// RegenerateWordAudio 重生成单词发音
// @Tags     EnglishWord
// @Summary  按配置的TTS服务重生成美式/英式发音链接
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.RegenerateWordAudioReq true "重生成参数"
// @Success  200  {object} response.Response{msg=string} "重生成成功"
// @Router   /englishLearning/word/regenerateAudio [post]
func (a *EnglishWordApi) RegenerateWordAudio(c *gin.Context) {
	var req request.RegenerateWordAudioReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := englishWordService.RegenerateWordAudio(req); err != nil {
		global.GVA_LOG.Error("重生成单词发音失败", zap.Error(err))
		response.FailWithMessage("重生成失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("重生成成功", c)
}

// PreflightTTS 预检TTS可用性
// @Tags     EnglishWord
// @Summary  使用当前配置对TTS服务进行联调预检
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.TTSPreflightReq false "预检参数，可选指定单词与美英发音开关"
// @Success  200  {object} response.Response{data=service.TTSPreflightResult,msg=string} "预检通过"
// @Router   /englishLearning/word/preflightTTS [post]
func (a *EnglishWordApi) PreflightTTS(c *gin.Context) {
	var req request.TTSPreflightReq
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		response.FailWithMessage("参数错误", c)
		return
	}

	data, err := englishWordService.PreflightTTS(req)
	if err != nil {
		global.GVA_LOG.Error("TTS预检失败", zap.Error(err))
		response.FailWithMessage("预检失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(data, "预检通过", c)
}

// GetWordList 分页获取某章节单词列表
// @Tags     EnglishWord
// @Summary  分页获取某章节单词列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.WordListSearch true "分页参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/word/getWordList [get]
func (a *EnglishWordApi) GetWordList(c *gin.Context) {
	var pageInfo request.WordListSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	page, pageSize := normalizeWordListPage(pageInfo.Page, pageInfo.PageSize)

	list, total, err := englishWordService.GetWordListByChapter(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取单词列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(response.PageResult{List: utils.LocalizeI18nPayloadByContext(c, list), Total: total, Page: page, PageSize: pageSize}, "获取成功", c)
}

// ReportWordError 上报错词
// @Tags     EnglishWord
// @Summary  跟打错误时上报错词记录
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.ReportWordErrorReq true "错词参数"
// @Success  200  {object} response.Response{msg=string} "上报成功"
// @Router   /englishLearning/word/reportError [post]
func (a *EnglishWordApi) ReportWordError(c *gin.Context) {
	var req request.ReportWordErrorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	if err = englishWordService.ReportWordError(claims.BaseClaims.ID, req); err != nil {
		global.GVA_LOG.Error("上报错词失败", zap.Error(err))
		response.FailWithMessage("上报失败", c)
		return
	}

	response.OkWithMessage("上报成功", c)
}

// GetWordErrorLogList 获取错题本列表
// @Tags     EnglishWord
// @Summary  获取当前用户错题本列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.WordErrorLogSearch true "分页参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/word/getErrorLogList [get]
func (a *EnglishWordApi) GetWordErrorLogList(c *gin.Context) {
	var pageInfo request.WordErrorLogSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	page, pageSize := normalizeWordListPage(pageInfo.Page, pageInfo.PageSize)

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	list, total, err := englishWordService.GetWordErrorLogList(claims.BaseClaims.ID, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取错题本列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(response.PageResult{List: utils.LocalizeI18nPayloadByContext(c, list), Total: total, Page: page, PageSize: pageSize}, "获取成功", c)
}

// DeleteWordErrorLog 删除错题本记录
// @Tags     EnglishWord
// @Summary  删除当前用户错题本中的一条记录
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.DeleteWordErrorLogReq true "要删除的单词ID"
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /englishLearning/word/deleteErrorLog [delete]
func (a *EnglishWordApi) DeleteWordErrorLog(c *gin.Context) {
	var req request.DeleteWordErrorLogReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	if err = englishWordService.DeleteWordErrorLog(claims.BaseClaims.ID, req.WordID); err != nil {
		global.GVA_LOG.Error("删除错题本记录失败", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// UpsertWordFromSQL 按SQL行批量导入时的单词新增/更新
// @Tags     EnglishWord
// @Summary  按单词唯一键执行新增或更新，并补充分类/章节绑定
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.UpsertWordFromSQLReq true "导入行参数"
// @Success  200  {object} response.Response{data=service.UpsertWordFromSQLResult,msg=string} "处理成功"
// @Router   /englishLearning/word/upsertSqlWord [post]
func (a *EnglishWordApi) UpsertWordFromSQL(c *gin.Context) {
	var req request.UpsertWordFromSQLReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	data, err := englishWordService.UpsertWordFromSQL(req)
	if err != nil {
		global.GVA_LOG.Error("SQL单词导入处理失败", zap.Error(err))
		response.FailWithMessage("处理失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(data, "处理成功", c)
}

// BatchFillWordFromDictionary 按分类批量补全词典信息
// @Tags     EnglishWord
// @Summary  按分类批量补全音标/词性/释义/例句，并可选翻译释义与例句
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.BatchFillWordFromDictionaryReq true "批量补全参数"
// @Success  200  {object} response.Response{data=service.BatchFillWordFromDictionaryResult,msg=string} "处理成功"
// @Router   /englishLearning/word/batchFillFromDictionary [post]
func (a *EnglishWordApi) BatchFillWordFromDictionary(c *gin.Context) {
	var req request.BatchFillWordFromDictionaryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	data, err := englishWordService.BatchFillWordFromDictionary(req)
	if err != nil {
		global.GVA_LOG.Error("词典批量补全失败", zap.Error(err))
		response.FailWithMessage("处理失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(data, "处理成功", c)
}
