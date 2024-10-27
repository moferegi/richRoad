import service from '@/utils/request'

// @Tags Sku
// @Summary 创建sku
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Sku true "创建sku"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /sku/createSku [post]
export const createSku = (data) => {
  return service({
    url: '/sku/createSku',
    method: 'post',
    data
  })
}

// @Tags Sku
// @Summary 删除sku
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Sku true "删除sku"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sku/deleteSku [delete]
export const deleteSku = (params) => {
  return service({
    url: '/sku/deleteSku',
    method: 'delete',
    params
  })
}

// @Tags Sku
// @Summary 批量删除sku
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除sku"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sku/deleteSku [delete]
export const deleteSkuByIds = (params) => {
  return service({
    url: '/sku/deleteSkuByIds',
    method: 'delete',
    params
  })
}

// @Tags Sku
// @Summary 更新sku
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Sku true "更新sku"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sku/updateSku [put]
export const updateSku = (data) => {
  return service({
    url: '/sku/updateSku',
    method: 'put',
    data
  })
}

// @Tags Sku
// @Summary 用id查询sku
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Sku true "用id查询sku"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sku/findSku [get]
export const findSku = (params) => {
  return service({
    url: '/sku/findSku',
    method: 'get',
    params
  })
}

// @Tags Sku
// @Summary 分页获取sku列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取sku列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sku/getSkuList [get]
export const getSkuList = (params) => {
  return service({
    url: '/sku/getSkuList',
    method: 'get',
    params
  })
}
