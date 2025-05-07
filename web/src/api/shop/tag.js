import service from '@/utils/request'
// @Tags Tag
// @Summary 创建标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Tag true "创建标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /tag/createTag [post]
export const createTag = (data) => {
  return service({
    url: '/tag/createTag',
    method: 'post',
    data
  })
}

// @Tags Tag
// @Summary 删除标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Tag true "删除标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /tag/deleteTag [delete]
export const deleteTag = (params) => {
  return service({
    url: '/tag/deleteTag',
    method: 'delete',
    params
  })
}

// @Tags Tag
// @Summary 批量删除标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /tag/deleteTag [delete]
export const deleteTagByIds = (params) => {
  return service({
    url: '/tag/deleteTagByIds',
    method: 'delete',
    params
  })
}

// @Tags Tag
// @Summary 更新标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Tag true "更新标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /tag/updateTag [put]
export const updateTag = (data) => {
  return service({
    url: '/tag/updateTag',
    method: 'put',
    data
  })
}

// @Tags Tag
// @Summary 用id查询标签
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.Tag true "用id查询标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /tag/findTag [get]
export const findTag = (params) => {
  return service({
    url: '/tag/findTag',
    method: 'get',
    params
  })
}

// @Tags Tag
// @Summary 分页获取标签列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取标签列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /tag/getTagList [get]
export const getTagList = (params) => {
  return service({
    url: '/tag/getTagList',
    method: 'get',
    params
  })
}

// @Tags Tag
// @Summary 不需要鉴权的标签接口
// @Accept application/json
// @Produce application/json
// @Param data query shopReq.TagSearch true "分页获取标签列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /tag/getTagPublic [get]
export const getTagPublic = () => {
  return service({
    url: '/tag/getTagPublic',
    method: 'get',
  })
}
