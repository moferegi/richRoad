package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GoodApi struct {
}

var goodService = service.ServiceGroupApp.ShopServiceGroup.GoodService

// CreateGood 创建商品
// @Tags Good
// @Summary 创建商品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Good true "创建商品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /good/createGood [post]
func (goodApi *GoodApi) CreateGood(c *gin.Context) {
	var good shop.Good
	err := c.ShouldBindJSON(&good)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := goodService.CreateGood(&good); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteGood 删除商品
// @Tags Good
// @Summary 删除商品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Good true "删除商品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /good/deleteGood [delete]
func (goodApi *GoodApi) DeleteGood(c *gin.Context) {
	ID := c.Query("ID")
	if err := goodService.DeleteGood(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteGoodByIds 批量删除商品
// @Tags Good
// @Summary 批量删除商品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /good/deleteGoodByIds [delete]
func (goodApi *GoodApi) DeleteGoodByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	if err := goodService.DeleteGoodByIds(IDs); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateGood 更新商品
// @Tags Good
// @Summary 更新商品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Good true "更新商品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /good/updateGood [put]
func (goodApi *GoodApi) UpdateGood(c *gin.Context) {
	var good shop.Good
	err := c.ShouldBindJSON(&good)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := goodService.UpdateGood(good); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindGood 用id查询商品
// @Tags Good
// @Summary 用id查询商品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shop.Good true "用id查询商品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /good/findGood [get]
func (goodApi *GoodApi) FindGood(c *gin.Context) {
	ID := c.Query("ID")
	if regood, err := goodService.GetGood(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"regood": regood}, c)
	}
}

// GetGoodList 分页获取商品列表
// @Tags Good
// @Summary 分页获取商品列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shopReq.GoodSearch true "分页获取商品列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /good/getGoodList [get]
func (goodApi *GoodApi) GetGoodList(c *gin.Context) {
	var pageInfo shopReq.GoodSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := goodService.GetGoodInfoList(pageInfo); err != nil {
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

// GetGoodPublic 不需要鉴权的商品接口
// @Tags Good
// @Summary 不需要鉴权的商品接口
// @accept application/json
// @Produce application/json
// @Param data query shopReq.GoodSearch true "分页获取商品列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /good/getGoodList [get]
func (goodApi *GoodApi) GetGoodPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的商品接口信息",
	}, "获取成功", c)
}
