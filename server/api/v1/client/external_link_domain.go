package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ExternalLinkDomainApi struct{}

// CreateExternalLinkDomain 创建外部链接域名
// @Tags ExternalLinkDomain
// @Summary 创建外部链接域名
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.ExternalLinkDomain true "域名信息"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /extDomain/createExternalLinkDomain [post]
func (a *ExternalLinkDomainApi) CreateExternalLinkDomain(c *gin.Context) {
	var domain client.ExternalLinkDomain
	if err := c.ShouldBindJSON(&domain); err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	if err := extDomainService.CreateExternalLinkDomain(domain); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		failClientWithKey(c, "createFail")
		return
	}
	response.OkWithMessage(i18n.T(c, "createSuccess"), c)
}

// DeleteExternalLinkDomain 删除外部链接域名
// @Tags ExternalLinkDomain
// @Summary 删除外部链接域名
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /extDomain/deleteExternalLinkDomain [delete]
func (a *ExternalLinkDomainApi) DeleteExternalLinkDomain(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	if err := extDomainService.DeleteExternalLinkDomain(req.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		failClientWithKey(c, "deleteFail")
		return
	}
	response.OkWithMessage(i18n.T(c, "deleteSuccess"), c)
}

// UpdateExternalLinkDomain 更新外部链接域名
// @Tags ExternalLinkDomain
// @Summary 更新外部链接域名
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.ExternalLinkDomain true "域名信息"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /extDomain/updateExternalLinkDomain [put]
func (a *ExternalLinkDomainApi) UpdateExternalLinkDomain(c *gin.Context) {
	var domain client.ExternalLinkDomain
	if err := c.ShouldBindJSON(&domain); err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	if err := extDomainService.UpdateExternalLinkDomain(domain); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		failClientWithKey(c, "updateFail")
		return
	}
	response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
}

// GetExternalLinkDomainList 分页获取外部链接域名列表
// @Tags ExternalLinkDomain
// @Summary 分页获取外部链接域名列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.ExternalLinkDomainSearch true "分页参数"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /extDomain/getExternalLinkDomainList [get]
func (a *ExternalLinkDomainApi) GetExternalLinkDomainList(c *gin.Context) {
	var pageInfo request.ExternalLinkDomainSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	list, total, err := extDomainService.GetExternalLinkDomainList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		failClientWithKey(c, "getFail")
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, i18n.T(c, "getSuccess"), c)
}

// SetDefaultDomain 设置默认域名
// @Tags ExternalLinkDomain
// @Summary 设置默认域名
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body object true "ID"
// @Success 200 {object} response.Response{msg=string} "设置成功"
// @Router /extDomain/setDefaultDomain [post]
func (a *ExternalLinkDomainApi) SetDefaultDomain(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	if err := extDomainService.SetDefaultDomain(req.ID); err != nil {
		global.GVA_LOG.Error("设置失败!", zap.Error(err))
		failClientWithKey(c, "setFail")
		return
	}
	response.OkWithMessage(i18n.T(c, "setSuccess"), c)
}

// GetDefaultDomain 获取默认域名（公开接口）
// @Tags ExternalLinkDomain
// @Summary 获取默认域名
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=string,msg=string} "获取成功"
// @Router /extDomain/getDefaultDomain [get]
func (a *ExternalLinkDomainApi) GetDefaultDomain(c *gin.Context) {
	domain, err := extDomainService.GetDefaultDomain()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		failClientWithKey(c, "getFail")
		return
	}
	response.OkWithDetailed(i18n.LocalizeResponseData(c, domain), i18n.T(c, "getSuccess"), c)
}
