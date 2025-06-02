import {request} from '@/utils/request.js'

// @Tags Comment
// @Summary 创建用户评论
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Comment true "创建用户评论"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /comment/createComment [post]
export const createComment = (data) => {
  return request({
    url: '/comment/createComment',
    method: 'post',
    data
  })
}

// @Tags Comment
// @Summary 删除用户评论
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Comment true "删除用户评论"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /comment/deleteComment [delete]
export const deleteComment = (params) => {
  return request({
    url: '/comment/deleteComment',
    method: 'delete',
    params
  })
}

// @Tags Comment
// @Summary 批量删除用户评论
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除用户评论"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /comment/deleteComment [delete]
export const deleteCommentByIds = (params) => {
  return request({
    url: '/comment/deleteCommentByIds',
    method: 'delete',
    params
  })
}

// @Tags Comment
// @Summary 更新用户评论
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Comment true "更新用户评论"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /comment/updateComment [put]
export const updateComment = (data) => {
  return request({
    url: '/comment/updateComment',
    method: 'put',
    data
  })
}

// @Tags Comment
// @Summary 用id查询用户评论
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Comment true "用id查询用户评论"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /comment/findComment [get]
export const findComment = (params) => {
  return request({
    url: `/comment/findComment`,
    method: 'get',
    params
  })
}

// @Tags Comment
// @Summary 分页获取用户评论列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取用户评论列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /comment/getCommentList [get]
export const getCommentList = (params) => {
  return request({
    url: '/comment/getCommentList',
    method: 'get',
    params
  })
}
