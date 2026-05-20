package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

type BlacklistApi struct{}

// GetBlacklist 获取黑名单列表
// @Tags CustomerServiceAdmin
// @Summary 获取客服黑名单列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query csReq.BlacklistSearch true "分页参数"
// @Success 200 {object} response.Response{data=response.PageResult} "获取成功"
// @Router /cs/admin/blacklist/list [get]
func (a *BlacklistApi) GetBlacklist(c *gin.Context) {
	var search csReq.BlacklistSearch
	if err := c.ShouldBindQuery(&search); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := service.Service.BlacklistService.GetList(search)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     search.Page,
		PageSize: search.PageSize,
	}, "获取成功", c)
}

// AddToBlacklist 加入黑名单
// @Tags CustomerServiceAdmin
// @Summary 将用户加入客服黑名单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.CsBlacklist true "黑名单信息"
// @Success 200 {object} response.Response{msg=string} "操作成功"
// @Router /cs/admin/blacklist/add [post]
func (a *BlacklistApi) AddToBlacklist(c *gin.Context) {
	var bl model.CsBlacklist
	if err := c.ShouldBindJSON(&bl); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	operatorID := utils.GetUserID(c)
	if err := service.Service.BlacklistService.Add(bl.ClientUserID, operatorID, bl.Reason); err != nil {
		response.FailWithMessage("操作失败", c)
		return
	}
	response.OkWithMessage("操作成功", c)
}

// RemoveFromBlacklist 移出黑名单
// @Tags CustomerServiceAdmin
// @Summary 移出客服黑名单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query int true "黑名单ID"
// @Success 200 {object} response.Response{msg=string} "操作成功"
// @Router /cs/admin/blacklist/remove [delete]
func (a *BlacklistApi) RemoveFromBlacklist(c *gin.Context) {
	var bl model.CsBlacklist
	if err := c.ShouldBindQuery(&bl); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := service.Service.BlacklistService.Remove(bl.ID); err != nil {
		response.FailWithMessage("操作失败", c)
		return
	}
	response.OkWithMessage("操作成功", c)
}
