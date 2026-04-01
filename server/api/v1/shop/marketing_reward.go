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

type MarketingRewardApi struct{}

var marketingRewardService = service.ServiceGroupApp.ShopServiceGroup.MarketingRewardService

// CreateMarketingReward 创建营销奖励规则
// @Tags MarketingReward
// @Summary 创建营销奖励规则
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.MarketingReward true "奖励规则信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /marketingReward/createMarketingReward [post]
func (api *MarketingRewardApi) CreateMarketingReward(c *gin.Context) {
	var info shop.MarketingReward
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := marketingRewardService.CreateMarketingReward(&info); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteMarketingReward 删除营销奖励规则
// @Tags MarketingReward
// @Summary 删除营销奖励规则
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.MarketingReward true "删除奖励规则"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /marketingReward/deleteMarketingReward [delete]
func (api *MarketingRewardApi) DeleteMarketingReward(c *gin.Context) {
	ID := c.Query("ID")
	if err := marketingRewardService.DeleteMarketingReward(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// UpdateMarketingReward 更新营销奖励规则
// @Tags MarketingReward
// @Summary 更新营销奖励规则
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.MarketingReward true "更新奖励规则"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /marketingReward/updateMarketingReward [put]
func (api *MarketingRewardApi) UpdateMarketingReward(c *gin.Context) {
	var info shop.MarketingReward
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := marketingRewardService.UpdateMarketingReward(info); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// GetMarketingRewardList 获取营销奖励规则列表
// @Tags MarketingReward
// @Summary 获取营销奖励规则列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shopReq.MarketingRewardSearch true "分页获取奖励规则列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /marketingReward/getMarketingRewardList [get]
func (api *MarketingRewardApi) GetMarketingRewardList(c *gin.Context) {
	var pageInfo shopReq.MarketingRewardSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := marketingRewardService.GetMarketingRewardList(pageInfo); err != nil {
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

// GetMarketingRewardByType 根据类型获取奖励规则
// @Tags MarketingReward
// @Summary 根据类型获取奖励规则
// @accept application/json
// @Produce application/json
// @Param triggerType query string true "触发类型(register/sub_register/sign_in/order)"
// @Success 200 {object} response.Response{data=shop.MarketingReward,msg=string} "获取成功"
// @Router /marketingReward/getMarketingRewardByType [get]
func (api *MarketingRewardApi) GetMarketingRewardByType(c *gin.Context) {
	triggerType := c.Query("triggerType")
	if triggerType == "" {
		response.FailWithMessage("触发类型不能为空", c)
		return
	}
	if data, err := marketingRewardService.GetMarketingRewardByType(triggerType); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(data, "获取成功", c)
	}
}
