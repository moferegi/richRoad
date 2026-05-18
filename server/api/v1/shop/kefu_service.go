package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type KefuApi struct{}

// CreateKefu 创建客服
// @Tags Kefu
// @Summary 创建客服
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shop.Kefu true "创建客服"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /kefu/createKefu [post]
func (kefuApi *KefuApi) CreateKefu(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var kefu shop.Kefu
	err := c.ShouldBindJSON(&kefu)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = kefuService.CreateKefu(ctx, &kefu)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteKefu 删除客服
// @Tags Kefu
// @Summary 删除客服
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shop.Kefu true "删除客服"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /kefu/deleteKefu [delete]
func (kefuApi *KefuApi) DeleteKefu(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := kefuService.DeleteKefu(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteKefuByIds 批量删除客服
// @Tags Kefu
// @Summary 批量删除客服
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /kefu/deleteKefuByIds [delete]
func (kefuApi *KefuApi) DeleteKefuByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := kefuService.DeleteKefuByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateKefu 更新客服
// @Tags Kefu
// @Summary 更新客服
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shop.Kefu true "更新客服"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /kefu/updateKefu [put]
func (kefuApi *KefuApi) UpdateKefu(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var kefu shop.Kefu
	err := c.ShouldBindJSON(&kefu)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = kefuService.UpdateKefu(ctx, kefu)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindKefu 用id查询客服
// @Tags Kefu
// @Summary 用id查询客服
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询客服"
// @Success 200 {object} response.Response{data=shop.Kefu,msg=string} "查询成功"
// @Router /kefu/findKefu [get]
func (kefuApi *KefuApi) FindKefu(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	rekefu, err := kefuService.GetKefu(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(rekefu, c)
}

// GetKefuList 分页获取客服列表
// @Tags Kefu
// @Summary 分页获取客服列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query shopReq.KefuSearch true "分页获取客服列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /kefu/getKefuList [get]
func (kefuApi *KefuApi) GetKefuList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo shopReq.KefuSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := kefuService.GetKefuInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetKefuPublic 不需要鉴权的客服接口
// @Tags Kefu
// @Summary 获取客服列表（公开）
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]shop.Kefu,msg=string} "获取成功"
// @Router /kefu/getKefuPublic [get]
func (kefuApi *KefuApi) GetKefuPublic(c *gin.Context) {
	ctx := c.Request.Context()
	list, err := kefuService.GetKefuPublic(ctx)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(i18n.LocalizeResponseData(c, list), "获取成功", c)
}
