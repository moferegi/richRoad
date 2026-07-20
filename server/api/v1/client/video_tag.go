package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VideoTagApi struct{}

// CreateVideoTag 创建视频标签
// @Tags VideoTag
// @Summary 创建视频标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body client.VideoTag true "创建视频标签"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /videoTag/createVideoTag [post]
func (api *VideoTagApi) CreateVideoTag(c *gin.Context) {
	ctx := c.Request.Context()
	var tag client.VideoTag
	err := c.ShouldBindJSON(&tag)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = videoTagService.CreateVideoTag(ctx, &tag)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteVideoTag 删除视频标签
// @Tags VideoTag
// @Summary 删除视频标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /videoTag/deleteVideoTag [delete]
func (api *VideoTagApi) DeleteVideoTag(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Query("id")
	if id == "" {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	err := videoTagService.DeleteVideoTag(ctx, id)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteVideoTagByIds 批量删除视频标签
// @Tags VideoTag
// @Summary 批量删除视频标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID集合"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /videoTag/deleteVideoTagByIds [delete]
func (api *VideoTagApi) DeleteVideoTagByIds(c *gin.Context) {
	ctx := c.Request.Context()
	ids := c.QueryArray("ids[]")
	if len(ids) == 0 {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	err := videoTagService.DeleteVideoTagByIds(ctx, ids)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateVideoTag 更新视频标签
// @Tags VideoTag
// @Summary 更新视频标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body client.VideoTag true "更新视频标签"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /videoTag/updateVideoTag [put]
func (api *VideoTagApi) UpdateVideoTag(c *gin.Context) {
	ctx := c.Request.Context()
	var tag client.VideoTag
	err := c.ShouldBindJSON(&tag)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = videoTagService.UpdateVideoTag(ctx, tag)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindVideoTag 根据ID获取视频标签
// @Tags VideoTag
// @Summary 根据ID获取视频标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.GetById true "ID"
// @Success 200 {object} response.Response{data=client.VideoTag,msg=string} "获取成功"
// @Router /videoTag/findVideoTag [get]
func (api *VideoTagApi) FindVideoTag(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Query("id")
	if id == "" {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	tag, err := videoTagService.GetVideoTag(ctx, id)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(tag, c)
}

// GetVideoTagList 分页获取视频标签列表
// @Tags VideoTag
// @Summary 分页获取视频标签列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query clientReq.VideoTagSearch true "分页参数"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /videoTag/getVideoTagList [get]
func (api *VideoTagApi) GetVideoTagList(c *gin.Context) {
	ctx := c.Request.Context()
	var pageInfo clientReq.VideoTagSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := videoTagService.GetVideoTagInfoList(ctx, pageInfo)
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

// GetVideoTagPublic 获取公开视频标签列表
// @Tags VideoTag
// @Summary 获取公开视频标签列表
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]client.VideoTag,msg=string} "获取成功"
// @Router /videoTag/getVideoTagPublic [get]
func (api *VideoTagApi) GetVideoTagPublic(c *gin.Context) {
	ctx := c.Request.Context()
	list, err := videoTagService.GetVideoTagPublic(ctx)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}
