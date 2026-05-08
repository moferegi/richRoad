import service from '@/utils/request'

/**
 * 获取试衣币充值订单列表（管理端）
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const getTryonRechargeOrderList = (params) => {
  return service({
    url: '/tryonRechargeOrder/getTryonRechargeOrderList',
    method: 'get',
    params
  })
}

/**
 * 获取试衣币充值订单详情（管理端）
 * @param {Object} params
 * @param {number|string} params.ID 订单ID
 * @returns {Promise}
 */
export const findTryonRechargeOrder = (params) => {
  return service({
    url: '/tryonRechargeOrder/findTryonRechargeOrder',
    method: 'get',
    params
  })
}

/**
 * 确认试衣币充值订单到账（管理端）
 * @param {Object} params
 * @param {number|string} params.ID 订单ID
 * @param {string} [params.remark] 备注
 * @returns {Promise}
 */
export const confirmTryonRechargeOrderPayment = (params) => {
  return service({
    url: '/tryonRechargeOrder/confirmTryonRechargeOrderPayment',
    method: 'post',
    params
  })
}
