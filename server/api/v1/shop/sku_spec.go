package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SkuSpecApi struct{}

// CreateSkuSpec 创建SKU规格字典
// @Tags SkuSpec
// @Summary 创建SKU规格字典
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shop.SkuSpec true "创建SKU规格字典"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /skuSpec/createSkuSpec [post]
func (api *SkuSpecApi) CreateSkuSpec(c *gin.Context) {
	ctx := c.Request.Context()
	var spec shop.SkuSpec
	err := c.ShouldBindJSON(&spec)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	err = skuSpecService.CreateSkuSpec(ctx, &spec)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteSkuSpec 删除SKU规格字典
// @Tags SkuSpec
// @Summary 删除SKU规格字典
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shop.SkuSpec true "删除SKU规格字典"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /skuSpec/deleteSkuSpec [delete]
func (api *SkuSpecApi) DeleteSkuSpec(c *gin.Context) {
	ctx := c.Request.Context()
	ID := c.Query("ID")
	err := skuSpecService.DeleteSkuSpec(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteSkuSpecByIds 批量删除SKU规格字典
// @Tags SkuSpec
// @Summary 批量删除SKU规格字典
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除SKU规格字典"
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /skuSpec/deleteSkuSpecByIds [delete]
func (api *SkuSpecApi) DeleteSkuSpecByIds(c *gin.Context) {
	ctx := c.Request.Context()
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	err = skuSpecService.DeleteSkuSpecByIds(ctx, IDS.Ids)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateSkuSpec 更新SKU规格字典
// @Tags SkuSpec
// @Summary 更新SKU规格字典
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shop.SkuSpec true "更新SKU规格字典"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /skuSpec/updateSkuSpec [put]
func (api *SkuSpecApi) UpdateSkuSpec(c *gin.Context) {
	ctx := c.Request.Context()
	var spec shop.SkuSpec
	err := c.ShouldBindJSON(&spec)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	err = skuSpecService.UpdateSkuSpec(ctx, spec)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindSkuSpec 根据ID获取SKU规格字典
// @Tags SkuSpec
// @Summary 根据ID获取SKU规格字典
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query shop.SkuSpec true "根据ID获取SKU规格字典"
// @Success 200 {object} response.Response{data=shop.SkuSpec} "获取成功"
// @Router /skuSpec/findSkuSpec [get]
func (api *SkuSpecApi) FindSkuSpec(c *gin.Context) {
	ctx := c.Request.Context()
	ID := c.Query("ID")
	spec, err := skuSpecService.GetSkuSpec(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(spec, "获取成功", c)
}

// GetSkuSpecList 分页获取SKU规格字典列表
// @Tags SkuSpec
// @Summary 分页获取SKU规格字典列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query shopReq.SkuSpecSearch true "分页获取SKU规格字典列表"
// @Success 200 {object} response.Response{data=response.PageResult} "获取成功"
// @Router /skuSpec/getSkuSpecList [get]
func (api *SkuSpecApi) GetSkuSpecList(c *gin.Context) {
	ctx := c.Request.Context()
	var pageInfo shopReq.SkuSpecSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := skuSpecService.GetSkuSpecInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetAllSkuSpecs 获取所有SKU规格字典（下拉选择用）
// @Tags SkuSpec
// @Summary 获取所有SKU规格字典
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param type query string false "类型：spec或attr"
// @Success 200 {object} response.Response{data=[]shop.SkuSpec} "获取成功"
// @Router /skuSpec/getAllSkuSpecs [get]
func (api *SkuSpecApi) GetAllSkuSpecs(c *gin.Context) {
	ctx := c.Request.Context()
	specType := c.Query("type")
	list, err := skuSpecService.GetAllSkuSpecs(ctx, specType)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}
