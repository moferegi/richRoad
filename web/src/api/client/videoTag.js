import service from '@/utils/request'

// @Tags VideoTag
// @Summary 创建视频标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.VideoTag true "创建视频标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /videoTag/createVideoTag [post]
export const createVideoTag = (data) => {
  return service({
    url: '/videoTag/createVideoTag',
    method: 'post',
    data
  })
}

// @Tags VideoTag
// @Summary 删除视频标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.VideoTag true "删除视频标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /videoTag/deleteVideoTag [delete]
export const deleteVideoTag = (params) => {
  return service({
    url: '/videoTag/deleteVideoTag',
    method: 'delete',
    params
  })
}

// @Tags VideoTag
// @Summary 批量删除视频标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除视频标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /videoTag/deleteVideoTagByIds [delete]
export const deleteVideoTagByIds = (params) => {
  return service({
    url: '/videoTag/deleteVideoTagByIds',
    method: 'delete',
    params
  })
}

// @Tags VideoTag
// @Summary 更新视频标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.VideoTag true "更新视频标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /videoTag/updateVideoTag [put]
export const updateVideoTag = (data) => {
  return service({
    url: '/videoTag/updateVideoTag',
    method: 'put',
    data
  })
}

// @Tags VideoTag
// @Summary 根据ID获取视频标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.GetById true "ID"
// @Success 200 {object} response.Response{data=model.VideoTag,msg=string} "获取成功"
// @Router /videoTag/findVideoTag [get]
export const findVideoTag = (params) => {
  return service({
    url: '/videoTag/findVideoTag',
    method: 'get',
    params
  })
}

// @Tags VideoTag
// @Summary 分页获取视频标签列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页参数"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /videoTag/getVideoTagList [get]
export const getVideoTagList = (params) => {
  return service({
    url: '/videoTag/getVideoTagList',
    method: 'get',
    params
  })
}

// @Tags VideoTag
// @Summary 获取公开视频标签列表
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.VideoTag,msg=string} "获取成功"
// @Router /videoTag/getVideoTagPublic [get]
export const getVideoTagPublic = (params) => {
  return service({
    url: '/videoTag/getVideoTagPublic',
    method: 'get',
    params
  })
}
