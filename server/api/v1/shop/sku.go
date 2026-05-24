package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SkuApi struct {
}

var skuService = service.ServiceGroupApp.ShopServiceGroup.SkuService

// CreateSku 创建sku
// @Tags Sku
// @Summary 创建sku
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Sku true "创建sku"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /sku/createSku [post]
func (skuApi *SkuApi) CreateSku(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	var sku shop.Sku
	err := c.ShouldBindJSON(&sku)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := skuService.CreateSku(&sku); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteSku 删除sku
// @Tags Sku
// @Summary 删除sku
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Sku true "删除sku"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sku/deleteSku [delete]
func (skuApi *SkuApi) DeleteSku(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	ID := c.Query("ID")
	if err := skuService.DeleteSku(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteSkuByIds 批量删除sku
// @Tags Sku
// @Summary 批量删除sku
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /sku/deleteSkuByIds [delete]
func (skuApi *SkuApi) DeleteSkuByIds(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	IDs := c.QueryArray("IDs[]")
	if err := skuService.DeleteSkuByIds(IDs); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateSku 更新sku
// @Tags Sku
// @Summary 更新sku
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Sku true "更新sku"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sku/updateSku [put]
func (skuApi *SkuApi) UpdateSku(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	var sku shop.Sku
	err := c.ShouldBindJSON(&sku)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := skuService.UpdateSku(sku); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindSku 用id查询sku
// @Tags Sku
// @Summary 用id查询sku
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shop.Sku true "用id查询sku"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sku/findSku [get]
func (skuApi *SkuApi) FindSku(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	ID := c.Query("ID")
	if resku, err := skuService.GetSku(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"resku": resku}, c)
	}
}

// GetSkuList 分页获取sku列表
// @Tags Sku
// @Summary 分页获取sku列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shopReq.SkuSearch true "分页获取sku列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sku/getSkuList [get]
func (skuApi *SkuApi) GetSkuList(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	var pageInfo shopReq.SkuSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if list, total, err := skuService.GetSkuInfoList(pageInfo); err != nil {
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

// GetSkuPublic 不需要鉴权的sku接口
// @Tags Sku
// @Summary 不需要鉴权的sku接口
// @accept application/json
// @Produce application/json
// @Param data query shopReq.SkuSearch true "分页获取sku列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sku/getSkuList [get]
func (skuApi *SkuApi) GetSkuPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的sku接口信息",
	}, "获取成功", c)
}
