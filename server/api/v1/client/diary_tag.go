package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	serverUtils "github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DiaryTagApi struct{}

// CreateDiaryTag 创建日记标签
// @Tags DiaryTag
// @Summary 创建日记标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body client.DiaryTag true "创建日记标签"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /diaryTag/createDiaryTag [post]
func (api *DiaryTagApi) CreateDiaryTag(c *gin.Context) {
	ctx := c.Request.Context()
	var tag client.DiaryTag
	err := c.ShouldBindJSON(&tag)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = diaryTagService.CreateDiaryTag(ctx, &tag)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteDiaryTag 删除日记标签
// @Tags DiaryTag
// @Summary 删除日记标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /diaryTag/deleteDiaryTag [delete]
func (api *DiaryTagApi) DeleteDiaryTag(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Query("id")
	if id == "" {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	err := diaryTagService.DeleteDiaryTag(ctx, id)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteDiaryTagByIds 批量删除日记标签
// @Tags DiaryTag
// @Summary 批量删除日记标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID集合"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /diaryTag/deleteDiaryTagByIds [delete]
func (api *DiaryTagApi) DeleteDiaryTagByIds(c *gin.Context) {
	ctx := c.Request.Context()
	ids := c.QueryArray("ids[]")
	if len(ids) == 0 {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	err := diaryTagService.DeleteDiaryTagByIds(ctx, ids)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateDiaryTag 更新日记标签
// @Tags DiaryTag
// @Summary 更新日记标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body client.DiaryTag true "更新日记标签"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /diaryTag/updateDiaryTag [put]
func (api *DiaryTagApi) UpdateDiaryTag(c *gin.Context) {
	ctx := c.Request.Context()
	var tag client.DiaryTag
	err := c.ShouldBindJSON(&tag)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = diaryTagService.UpdateDiaryTag(ctx, tag)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindDiaryTag 根据ID获取日记标签
// @Tags DiaryTag
// @Summary 根据ID获取日记标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.GetById true "ID"
// @Success 200 {object} response.Response{data=client.DiaryTag,msg=string} "获取成功"
// @Router /diaryTag/findDiaryTag [get]
func (api *DiaryTagApi) FindDiaryTag(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Query("id")
	if id == "" {
		response.FailWithMessage("ID不能为空", c)
		return
	}
	tag, err := diaryTagService.GetDiaryTag(ctx, id)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(tag, c)
}

// GetDiaryTagList 分页获取日记标签列表
// @Tags DiaryTag
// @Summary 分页获取日记标签列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query clientReq.DiaryTagSearch true "分页参数"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /diaryTag/getDiaryTagList [get]
func (api *DiaryTagApi) GetDiaryTagList(c *gin.Context) {
	ctx := c.Request.Context()
	var pageInfo clientReq.DiaryTagSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := diaryTagService.GetDiaryTagInfoList(ctx, pageInfo)
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

// GetDiaryTagPublic 获取公开日记标签列表
// @Tags DiaryTag
// @Summary 获取公开日记标签列表
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]client.DiaryTag,msg=string} "获取成功"
// @Router /diaryTag/getDiaryTagPublic [get]
func (api *DiaryTagApi) GetDiaryTagPublic(c *gin.Context) {
	ctx := c.Request.Context()
	list, err := diaryTagService.GetDiaryTagPublic(ctx)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithData(serverUtils.LocalizeI18nPayloadByContext(c, list), c)
}