import { request } from '@/utils/request.js'

// 创建试衣币充值订单
export const createTryonRechargeOrder = (data) => {
  return request({
    url: '/tryonRechargeOrder/createTryonRechargeOrder',
    method: 'post',
    data,
  })
}

// 更新充值订单支付方式（仅待支付）
export const updateTryonRechargeOrderPayMethod = (data) => {
  return request({
    url: '/tryonRechargeOrder/updateTryonRechargeOrderPayMethod',
    method: 'post',
    data,
  })
}

// 提交充值订单付款确认（仅扫码支付且待支付）
export const submitTryonRechargeOrderPayment = (data) => {
  return request({
    url: '/tryonRechargeOrder/submitTryonRechargeOrderPayment',
    method: 'post',
    data,
  })
}

// 取消充值订单（仅待支付）
export const cancelTryonRechargeOrder = (id) => {
  return request({
    url: '/tryonRechargeOrder/cancelTryonRechargeOrder',
    method: 'post',
    params: { ID: id },
  })
}

// 获取我的充值订单详情
export const selfTryonRechargeOrder = (id) => {
  return request({
    url: '/tryonRechargeOrder/selfTryonRechargeOrder',
    method: 'get',
    params: { ID: id },
  })
}

// 获取我的充值订单列表
export const getMyTryonRechargeOrderList = (params = {}) => {
  return request({
    url: '/tryonRechargeOrder/getMyTryonRechargeOrderList',
    method: 'get',
    params,
  })
}
