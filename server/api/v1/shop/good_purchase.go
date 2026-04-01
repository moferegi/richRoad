package shop

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GoodPurchaseApi struct{}

var goodPurchaseService = service.ServiceGroupApp.ShopServiceGroup.GoodPurchaseService

// CreateGoodPurchase 创建进货记录
// @Tags GoodPurchase
// @Summary 创建进货记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.GoodPurchase true "进货记录信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /goodPurchase/createGoodPurchase [post]
func (api *GoodPurchaseApi) CreateGoodPurchase(c *gin.Context) {
	var info shop.GoodPurchase
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := goodPurchaseService.CreateGoodPurchase(&info); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteGoodPurchase 删除进货记录
// @Tags GoodPurchase
// @Summary 删除进货记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.GoodPurchase true "删除进货记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /goodPurchase/deleteGoodPurchase [delete]
func (api *GoodPurchaseApi) DeleteGoodPurchase(c *gin.Context) {
	ID := c.Query("ID")
	if err := goodPurchaseService.DeleteGoodPurchase(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// UpdateGoodPurchase 更新进货记录
// @Tags GoodPurchase
// @Summary 更新进货记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.GoodPurchase true "更新进货记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /goodPurchase/updateGoodPurchase [put]
func (api *GoodPurchaseApi) UpdateGoodPurchase(c *gin.Context) {
	var info shop.GoodPurchase
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := goodPurchaseService.UpdateGoodPurchase(info); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// GetGoodPurchaseList 获取进货记录列表
// @Tags GoodPurchase
// @Summary 获取进货记录列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shopReq.GoodPurchaseSearch true "分页获取进货记录列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /goodPurchase/getGoodPurchaseList [get]
func (api *GoodPurchaseApi) GetGoodPurchaseList(c *gin.Context) {
	var pageInfo shopReq.GoodPurchaseSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := goodPurchaseService.GetGoodPurchaseList(pageInfo); err != nil {
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

// GetGoodPurchaseSummary 获取商品进货汇总
// @Tags GoodPurchase
// @Summary 获取商品进货汇总
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param goodId query string true "商品ID"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /goodPurchase/getGoodPurchaseSummary [get]
func (api *GoodPurchaseApi) GetGoodPurchaseSummary(c *gin.Context) {
	goodId := c.Query("goodId")
	if goodId == "" {
		response.FailWithMessage("商品ID不能为空", c)
		return
	}
	var id uint
	if v, err := strconv.ParseUint(goodId, 10, 64); err == nil {
		id = uint(v)
	}
	if totalQty, totalCost, err := goodPurchaseService.GetGoodPurchaseSummary(id); err != nil {
		global.GVA_LOG.Error("获取汇总失败!", zap.Error(err))
		response.FailWithMessage("获取汇总失败", c)
	} else {
		response.OkWithDetailed(gin.H{"totalQty": totalQty, "totalCost": totalCost}, "获取成功", c)
	}
}
