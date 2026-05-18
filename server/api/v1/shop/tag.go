package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TagApi struct{}

// CreateTag 创建标签
// @Tags Tag
// @Summary 创建标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shop.Tag true "创建标签"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /tag/createTag [post]
func (tagApi *TagApi) CreateTag(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var tag shop.Tag
	err := c.ShouldBindJSON(&tag)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = tagService.CreateTag(ctx, &tag)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteTag 删除标签
// @Tags Tag
// @Summary 删除标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shop.Tag true "删除标签"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /tag/deleteTag [delete]
func (tagApi *TagApi) DeleteTag(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := tagService.DeleteTag(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteTagByIds 批量删除标签
// @Tags Tag
// @Summary 批量删除标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /tag/deleteTagByIds [delete]
func (tagApi *TagApi) DeleteTagByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := tagService.DeleteTagByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateTag 更新标签
// @Tags Tag
// @Summary 更新标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shop.Tag true "更新标签"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /tag/updateTag [put]
func (tagApi *TagApi) UpdateTag(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var tag shop.Tag
	err := c.ShouldBindJSON(&tag)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = tagService.UpdateTag(ctx, tag)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindTag 用id查询标签
// @Tags Tag
// @Summary 用id查询标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询标签"
// @Success 200 {object} response.Response{data=shop.Tag,msg=string} "查询成功"
// @Router /tag/findTag [get]
func (tagApi *TagApi) FindTag(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	retag, err := tagService.GetTag(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(retag, c)
}

// GetTagList 分页获取标签列表
// @Tags Tag
// @Summary 分页获取标签列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query shopReq.TagSearch true "分页获取标签列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /tag/getTagList [get]
func (tagApi *TagApi) GetTagList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo shopReq.TagSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := tagService.GetTagInfoList(ctx, pageInfo)
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

// GetTagPublic 不需要鉴权的标签接口
// @Tags Tag
// @Summary 不需要鉴权的标签接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /tag/getTagPublic [get]
func (tagApi *TagApi) GetTagPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	tagService.GetTagPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的标签接口信息",
	}, "获取成功", c)
}
