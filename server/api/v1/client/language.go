package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LanguageApi struct{}

var languageService = service.ServiceGroupApp.ClientServiceGroup.SysLanguageService

// CreateLanguage 创建语言
// @Tags Language
// @Summary 创建语言
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.SysLanguage true "语言信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /language/createLanguage [post]
func (api *LanguageApi) CreateLanguage(c *gin.Context) {
	var info client.SysLanguage
	err := c.ShouldBindJSON(&info)
	if err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	if err := languageService.CreateSysLanguage(&info); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		failClientWithKey(c, "createFail")
	} else {
		response.OkWithMessage(i18n.T(c, "createSuccess"), c)
	}
}

// DeleteLanguage 删除语言
// @Tags Language
// @Summary 删除语言
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.SysLanguage true "删除语言"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /language/deleteLanguage [delete]
func (api *LanguageApi) DeleteLanguage(c *gin.Context) {
	ID := c.Query("ID")
	if err := languageService.DeleteSysLanguage(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		failClientWithKey(c, "deleteFail")
	} else {
		response.OkWithMessage(i18n.T(c, "deleteSuccess"), c)
	}
}

// UpdateLanguage 更新语言
// @Tags Language
// @Summary 更新语言
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.SysLanguage true "更新语言"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /language/updateLanguage [put]
func (api *LanguageApi) UpdateLanguage(c *gin.Context) {
	var info client.SysLanguage
	err := c.ShouldBindJSON(&info)
	if err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	if err := languageService.UpdateSysLanguage(info); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		failClientWithKey(c, "updateFail")
	} else {
		response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
	}
}

// GetLanguageList 获取语言列表（管理端）
// @Tags Language
// @Summary 获取语言列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query clientReq.LanguageSearch true "分页获取语言列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /language/getLanguageList [get]
func (api *LanguageApi) GetLanguageList(c *gin.Context) {
	var pageInfo clientReq.SysLanguageSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	if list, total, err := languageService.GetSysLanguageList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		failClientWithKey(c, "getFail")
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, i18n.T(c, "getSuccess"), c)
	}
}

// GetEnabledLanguages 获取启用的语言列表（客户端）
// @Tags Language
// @Summary 获取启用的语言列表
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]client.SysLanguage,msg=string} "获取成功"
// @Router /language/getEnabledLanguages [get]
func (api *LanguageApi) GetEnabledLanguages(c *gin.Context) {
	if list, err := languageService.GetEnabledLanguages(); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		failClientWithKey(c, "getFail")
	} else {
		response.OkWithDetailed(i18n.LocalizeResponseData(c, list), i18n.T(c, "getSuccess"), c)
	}
}

// TranslateI18n 翻译多语言文本
// @Tags Language
// @Summary 翻译多语言文本
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body clientReq.TranslateI18nRequest true "翻译参数"
// @Success 200 {object} response.Response{data=object,msg=string} "翻译成功"
// @Router /language/translateI18n [post]
func (api *LanguageApi) TranslateI18n(c *gin.Context) {
	var req clientReq.TranslateI18nRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}

	translations, err := languageService.TranslateI18n(req)
	if err != nil {
		global.GVA_LOG.Error("翻译失败!", zap.Error(err))
		failClientWithKey(c, "getFail")
		return
	}

	response.OkWithDetailed(gin.H{"translations": translations}, i18n.T(c, "getSuccess"), c)
}
