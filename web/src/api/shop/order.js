import service from '@/utils/request'

// @Tags Order
// @Summary 创建订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Order true "创建订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /order/createOrder [post]
export const createOrder = (data) => {
  return service({
    url: '/order/createOrder',
    method: 'post',
    data
  })
}

// @Tags Order
// @Summary 删除订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Order true "删除订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /order/deleteOrder [delete]
export const deleteOrder = (params) => {
  return service({
    url: '/order/deleteOrder',
    method: 'delete',
    params
  })
}

// @Tags Order
// @Summary 批量删除订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /order/deleteOrder [delete]
export const deleteOrderByIds = (params) => {
  return service({
    url: '/order/deleteOrderByIds',
    method: 'delete',
    params
  })
}

// @Tags Order
// @Summary 更新订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Order true "更新订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /order/updateOrder [put]
export const updateOrder = (data) => {
  return service({
    url: '/order/updateOrder',
    method: 'put',
    data
  })
}

// @Tags Order
// @Summary 用id查询订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Order true "用id查询订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /order/findOrder [get]
export const findOrder = (params) => {
  return service({
    url: '/order/findOrder',
    method: 'get',
    params
  })
}

// @Tags Order
// @Summary 分页获取订单列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取订单列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /order/getOrderList [get]
export const getOrderList = (params) => {
  return service({
    url: '/order/getOrderList',
    method: 'get',
    params
  })
}

export const checkRouters = (params) => {
  return service({
    url: '/order/checkRouters',
    method: 'get',
    params
  })
}

export const refundOrder = (data) => {
  return service({
    url: '/order/refundOrder',
    method: 'post',
    data
  })
}

/**
 * 管理员确认收款
 * @param {Object} params
 * @param {string} params.ID 订单ID
 * @returns {Promise}
 */
export const confirmPayment = (params) => {
  return service({
    url: '/order/confirmPayment',
    method: 'post',
    params
  })
}

/**
 * 批量更新订单状态
 * @param {Object} data
 * @param {string[]} data.IDs 订单ID列表
 * @param {string} data.status 目标状态
 * @returns {Promise}
 */
export const batchUpdateOrderStatus = (data) => {
  return service({
    url: '/order/batchUpdateOrderStatus',
    method: 'post',
    data
  })
}

/**
 * 更新订单状态（客户端）
 * @param {Object} params
 * @param {string} params.ID 订单ID
 * @param {string} params.status 订单状态
 * @returns {Promise}
 */
export const updateOrderStatus = (params) => {
  return service({
    url: '/order/updateOrderStatus',
    method: 'post',
    params
  })
}
