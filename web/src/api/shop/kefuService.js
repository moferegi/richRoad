import service from '@/utils/request'
// @Tags Kefu
// @Summary 创建客服
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Kefu true "创建客服"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /kefu/createKefu [post]
export const createKefu = (data) => {
  return service({
    url: '/kefu/createKefu',
    method: 'post',
    data
  })
}

// @Tags Kefu
// @Summary 删除客服
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Kefu true "删除客服"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /kefu/deleteKefu [delete]
export const deleteKefu = (params) => {
  return service({
    url: '/kefu/deleteKefu',
    method: 'delete',
    params
  })
}

// @Tags Kefu
// @Summary 批量删除客服
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除客服"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /kefu/deleteKefu [delete]
export const deleteKefuByIds = (params) => {
  return service({
    url: '/kefu/deleteKefuByIds',
    method: 'delete',
    params
  })
}

// @Tags Kefu
// @Summary 更新客服
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Kefu true "更新客服"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /kefu/updateKefu [put]
export const updateKefu = (data) => {
  return service({
    url: '/kefu/updateKefu',
    method: 'put',
    data
  })
}

// @Tags Kefu
// @Summary 用id查询客服
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.Kefu true "用id查询客服"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /kefu/findKefu [get]
export const findKefu = (params) => {
  return service({
    url: '/kefu/findKefu',
    method: 'get',
    params
  })
}

// @Tags Kefu
// @Summary 分页获取客服列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取客服列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /kefu/getKefuList [get]
export const getKefuList = (params) => {
  return service({
    url: '/kefu/getKefuList',
    method: 'get',
    params
  })
}

// @Tags Kefu
// @Summary 不需要鉴权的客服接口
// @Accept application/json
// @Produce application/json
// @Param data query shopReq.KefuSearch true "分页获取客服列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /kefu/getKefuPublic [get]
export const getKefuPublic = () => {
  return service({
    url: '/kefu/getKefuPublic',
    method: 'get',
  })
}
