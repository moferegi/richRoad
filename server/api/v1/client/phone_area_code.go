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

type PhoneAreaCodeApi struct{}

var phoneAreaCodeService = service.ServiceGroupApp.ClientServiceGroup.PhoneAreaCodeService

// CreatePhoneAreaCode 创建国际区号
// @Tags PhoneAreaCode
// @Summary 创建国际区号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.PhoneAreaCode true "区号信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /phoneAreaCode/createPhoneAreaCode [post]
func (api *PhoneAreaCodeApi) CreatePhoneAreaCode(c *gin.Context) {
	var info client.PhoneAreaCode
	err := c.ShouldBindJSON(&info)
	if err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	if err := phoneAreaCodeService.CreatePhoneAreaCode(&info); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		failClientWithKey(c, "createFail")
	} else {
		response.OkWithMessage(i18n.T(c, "createSuccess"), c)
	}
}

// DeletePhoneAreaCode 删除国际区号
// @Tags PhoneAreaCode
// @Summary 删除国际区号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.PhoneAreaCode true "删除区号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /phoneAreaCode/deletePhoneAreaCode [delete]
func (api *PhoneAreaCodeApi) DeletePhoneAreaCode(c *gin.Context) {
	ID := c.Query("ID")
	if err := phoneAreaCodeService.DeletePhoneAreaCode(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		failClientWithKey(c, "deleteFail")
	} else {
		response.OkWithMessage(i18n.T(c, "deleteSuccess"), c)
	}
}

// UpdatePhoneAreaCode 更新国际区号
// @Tags PhoneAreaCode
// @Summary 更新国际区号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.PhoneAreaCode true "更新区号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /phoneAreaCode/updatePhoneAreaCode [put]
func (api *PhoneAreaCodeApi) UpdatePhoneAreaCode(c *gin.Context) {
	var info client.PhoneAreaCode
	err := c.ShouldBindJSON(&info)
	if err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	if err := phoneAreaCodeService.UpdatePhoneAreaCode(info); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		failClientWithKey(c, "updateFail")
	} else {
		response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
	}
}

// GetPhoneAreaCodeList 获取区号列表（管理端）
// @Tags PhoneAreaCode
// @Summary 获取区号列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query clientReq.PhoneAreaCodeSearch true "分页获取区号列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /phoneAreaCode/getPhoneAreaCodeList [get]
func (api *PhoneAreaCodeApi) GetPhoneAreaCodeList(c *gin.Context) {
	var pageInfo clientReq.PhoneAreaCodeSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	if list, total, err := phoneAreaCodeService.GetPhoneAreaCodeList(pageInfo); err != nil {
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

// GetEnabledPhoneAreaCodes 获取启用的区号列表（客户端）
// @Tags PhoneAreaCode
// @Summary 获取启用的区号列表
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]client.PhoneAreaCode,msg=string} "获取成功"
// @Router /phoneAreaCode/getEnabledPhoneAreaCodes [get]
func (api *PhoneAreaCodeApi) GetEnabledPhoneAreaCodes(c *gin.Context) {
	if list, err := phoneAreaCodeService.GetEnabledPhoneAreaCodes(); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		failClientWithKey(c, "getFail")
	} else {
		response.OkWithDetailed(i18n.LocalizeResponseData(c, list), i18n.T(c, "getSuccess"), c)
	}
}
