package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CommentApi struct{}

// CreateComment 创建用户评论
// @Tags Comment
// @Summary 创建用户评论
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Comment true "创建用户评论"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /comment/createComment [post]
func (commentApi *CommentApi) CreateComment(c *gin.Context) {
	var comment shop.Comment
	err := c.ShouldBindJSON(&comment)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	comment.UserID = utils.GetUserID(c)
	err = commentService.CreateComment(&comment)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteComment 删除用户评论
// @Tags Comment
// @Summary 删除用户评论
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Comment true "删除用户评论"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /comment/deleteComment [delete]
func (commentApi *CommentApi) DeleteComment(c *gin.Context) {
	ID := c.Query("ID")
	err := commentService.DeleteComment(ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteCommentByIds 批量删除用户评论
// @Tags Comment
// @Summary 批量删除用户评论
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /comment/deleteCommentByIds [delete]
func (commentApi *CommentApi) DeleteCommentByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	err := commentService.DeleteCommentByIds(IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateComment 更新用户评论
// @Tags Comment
// @Summary 更新用户评论
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Comment true "更新用户评论"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /comment/updateComment [put]
func (commentApi *CommentApi) UpdateComment(c *gin.Context) {
	var comment shop.Comment
	err := c.ShouldBindJSON(&comment)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	err = commentService.UpdateComment(comment)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindComment 用id查询用户评论
// @Tags Comment
// @Summary 用id查询用户评论
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shop.Comment true "用id查询用户评论"
// @Success 200 {object} response.Response{data=shop.Comment,msg=string} "查询成功"
// @Router /comment/findComment [get]
func (commentApi *CommentApi) FindComment(c *gin.Context) {
	ID := c.Query("ID")
	recomment, err := commentService.GetComment(ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(i18n.LocalizeResponseData(c, recomment), c)
}

// GetComment 用id查询用户评论
// @Tags Comment
// @Summary 用id查询用户评论
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shop.Comment true "用id查询用户评论"
// @Success 200 {object} response.Response{data=shop.Comment,msg=string} "查询成功"
// @Router /comment/getComment [get]
func (commentApi *CommentApi) GetComment(c *gin.Context) {
	ID := c.Query("ID")
	recomment, err := commentService.GetCommentBk(ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(recomment, c)
}

// GetCommentList 分页获取用户评论列表
// @Tags Comment
// @Summary 分页获取用户评论列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shopReq.CommentSearch true "分页获取用户评论列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /comment/getCommentList [get]
func (commentApi *CommentApi) GetCommentList(c *gin.Context) {
	var pageInfo shopReq.CommentSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := commentService.GetCommentInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(i18n.LocalizeResponseData(c, response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}), "获取成功", c)
}

// GetCommentPublic 不需要鉴权的用户评论接口
// @Tags Comment
// @Summary 不需要鉴权的用户评论接口
// @accept application/json
// @Produce application/json
// @Param data query shopReq.CommentSearch true "分页获取用户评论列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /comment/getCommentPublic [get]
func (commentApi *CommentApi) GetCommentPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的用户评论接口信息",
	}, "获取成功", c)
}
