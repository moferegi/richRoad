import {request} from '@/utils/request.js'

// 直接下单
export const placeOrder  = (data) => {
    return request({
		url:'/order/placeOrder',
        method: 'post',
		data: data
    })
}

// 购物车下单
export const placeOrderByCart  = (data) => {
    return request({
		url:'/order/placeOrderByCart',
        method: 'post',
		data: data
    })
}

// 变更订单积分抵扣
export const changeOrderPoints = (params) => {
    return request({
        url: '/order/changeOrderPoints',
        method: 'post',
        params: params
    })
}

// 改变订单购物券
export const changeOrderCoupon  = (params) => {
  return request({
  url:`/order/changeOrderCoupon?orderID=${params.orderID}&couponNum=${params.couponNum}`,
  method: 'post'
  })
}

// @Tags Order
// @Summary 创建订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Order true "创建订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /order/createOrder [post]
export const createOrder = (data) => {
  return request({
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
  return request({
    url: '/order/deleteOrder',
    method: 'delete',
    params
  })
}

export const updateOrderStatus = (data) => {
  return request({
    url: `/order/updateOrderStatus?ID=${data.ID}&status=${data.status}`,
    method: 'post',
    data
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
  return request({
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
  return request({
    url: '/order/updateOrder',
    method: 'post',
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
  return request({
    url: `/order/findOrder?ID=${params.id}`,
    method: 'get',
  })
}



export const selfOrderComment = (orderID,goodID,SKUID) => {
  return request({
    url: `/order/selfOrderComment?ID=${orderID}&goodID=${goodID}&SKUID=${SKUID}`,
    method: 'get'
  })
}

export const selfOrder = (params) => {
  return request({
    url: `/order/selfOrder?ID=${params}`,
    method: 'get'
  })
}

export const selfOrderList = (params) => {
  return request({
    url: `/order/selfOrderList?ID=${params}`,
    method: 'get'
  })
}

export const SelfOrderList = (params) => {
  const query = typeof params === 'string'
    ? `status=${params}`
    : Object.entries(params).filter(([,v]) => v !== undefined && v !== '').map(([k,v]) => `${k}=${v}`).join('&')
  return request({
    url: `/order/selfOrderList?${query}`,
    method: 'get',
  })
}

// @Tags Order
// @Summary 查看物流
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "查看物流"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /order/CheckRouters [get]
export const checkRouters = (params) => {
  return request({
    url: `/order/checkRouters?express=${params}`,
    method: 'get',
  })
}

// 申请退款
export const applyRefund = (data) => {
  return request({
    url: '/order/applyRefund',
    method: 'post',
    data
  })
}

// 获取浏览信息
export const getGoodHistory = (params) => {
  return request({
    url: '/good/getGoodHistory',
    method: 'get',
  })
}

// 清空浏览历史
export const clearGoodHistory = () => {
  return request({
    url: '/good/clearGoodHistory',
    method: 'delete',
  })
}
