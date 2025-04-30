import {request} from '@/utils/request.js'

// @Tags SkuDetail
// @Summary 创建商品详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SkuDetail true "创建商品详情"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /skuDetail/createSkuDetail [post]
export const createSkuDetail = (data) => {
  return request({
    url: '/skuDetail/createSkuDetail',
    method: 'post',
    data
  })
}

// @Tags SkuDetail
// @Summary 删除商品详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SkuDetail true "删除商品详情"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /skuDetail/deleteSkuDetail [delete]
export const deleteSkuDetail = (params) => {
  return request({
    url: '/skuDetail/deleteSkuDetail',
    method: 'delete',
    params
  })
}

// @Tags SkuDetail
// @Summary 批量删除商品详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除商品详情"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /skuDetail/deleteSkuDetail [delete]
export const deleteSkuDetailByIds = (params) => {
  return request({
    url: '/skuDetail/deleteSkuDetailByIds',
    method: 'delete',
    params
  })
}

// @Tags SkuDetail
// @Summary 更新商品详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SkuDetail true "更新商品详情"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /skuDetail/updateSkuDetail [put]
export const updateSkuDetail = (data) => {
  return request({
    url: '/skuDetail/updateSkuDetail',
    method: 'put',
    data
  })
}

// @Tags SkuDetail
// @Summary 用id查询商品详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.SkuDetail true "用id查询商品详情"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /skuDetail/findSkuDetail [get]
export const findSkuDetail = (params) => {
  return request({
    url: `/skuDetail/findSkuDetail?goodID=${params}`,
    method: 'get',
    params
  })
}

// @Tags SkuDetail
// @Summary 分页获取商品详情列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取商品详情列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /skuDetail/getSkuDetailList [get]
export const getSkuDetailList = (params) => {
  return request({
    url: '/skuDetail/getSkuDetailList',
    method: 'get',
    params
  })
}
