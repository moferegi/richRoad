package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PopupApi struct{}

var popupService = service.ServiceGroupApp.ShopServiceGroup.PopupService

// CreatePopup 创建弹窗
// @Tags Popup
// @Summary 创建弹窗
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Popup true "弹窗信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /popup/createPopup [post]
func (api *PopupApi) CreatePopup(c *gin.Context) {
	var info shop.Popup
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := popupService.CreatePopup(&info); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeletePopup 删除弹窗
// @Tags Popup
// @Summary 删除弹窗
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Popup true "删除弹窗"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /popup/deletePopup [delete]
func (api *PopupApi) DeletePopup(c *gin.Context) {
	ID := c.Query("ID")
	if err := popupService.DeletePopup(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// UpdatePopup 更新弹窗
// @Tags Popup
// @Summary 更新弹窗
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Popup true "更新弹窗"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /popup/updatePopup [put]
func (api *PopupApi) UpdatePopup(c *gin.Context) {
	var info shop.Popup
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := popupService.UpdatePopup(info); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// GetPopupList 获取弹窗列表
// @Tags Popup
// @Summary 获取弹窗列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shopReq.PopupSearch true "分页获取弹窗列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /popup/getPopupList [get]
func (api *PopupApi) GetPopupList(c *gin.Context) {
	var pageInfo shopReq.PopupSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := popupService.GetPopupList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetActivePopups 获取当前生效的弹窗（客户端）
// @Tags Popup
// @Summary 获取当前生效的弹窗
// @accept application/json
// @Produce application/json
// @Param position query string false "弹窗位置(兼容旧版)"
// @Param clientType query string false "客户端类型(web/uni/all)"
// @Param page query string false "当前页面路径"
// @Success 200 {object} response.Response{data=[]shop.Popup,msg=string} "获取成功"
// @Router /popup/getActivePopups [get]
func (api *PopupApi) GetActivePopups(c *gin.Context) {
	position := c.Query("position")
	clientType := c.Query("clientType")
	pagePath := c.Query("page")
	if list, err := popupService.GetActivePopups(position, clientType, pagePath); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(i18n.LocalizeResponseData(c, list), "获取成功", c)
	}
}

// GetPopupPagePathOptions 获取弹窗可选页面路径
// @Tags Popup
// @Summary 获取弹窗可选页面路径
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param clientType query string false "客户端类型(uni/web/all)"
// @Success 200 {object} response.Response{data=[]string,msg=string} "获取成功"
// @Router /popup/getPopupPagePathOptions [get]
func (api *PopupApi) GetPopupPagePathOptions(c *gin.Context) {
	clientType := c.Query("clientType")
	if list, err := popupService.GetPopupPagePathOptions(clientType); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(list, "获取成功", c)
	}
}
