import service from '@/utils/request'

// @Tags Collect
// @Summary 创建收藏
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Collect true "创建收藏"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /collect/createCollect [post]
export const createCollect = (data) => {
  return service({
    url: '/collect/createCollect',
    method: 'post',
    data
  })
}

// @Tags Collect
// @Summary 删除收藏
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Collect true "删除收藏"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /collect/deleteCollect [delete]
export const deleteCollect = (params) => {
  return service({
    url: '/collect/deleteCollect',
    method: 'delete',
    params
  })
}

// @Tags Collect
// @Summary 批量删除收藏
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除收藏"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /collect/deleteCollect [delete]
export const deleteCollectByIds = (params) => {
  return service({
    url: '/collect/deleteCollectByIds',
    method: 'delete',
    params
  })
}

// @Tags Collect
// @Summary 更新收藏
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Collect true "更新收藏"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /collect/updateCollect [put]
export const updateCollect = (data) => {
  return service({
    url: '/collect/updateCollect',
    method: 'put',
    data
  })
}

// @Tags Collect
// @Summary 用id查询收藏
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Collect true "用id查询收藏"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /collect/findCollect [get]
export const findCollect = (params) => {
  return service({
    url: '/collect/findCollect',
    method: 'get',
    params
  })
}

// @Tags Collect
// @Summary 分页获取收藏列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取收藏列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /collect/getCollectList [get]
export const getCollectList = (params) => {
  return service({
    url: '/collect/getCollectList',
    method: 'get',
    params
  })
}
