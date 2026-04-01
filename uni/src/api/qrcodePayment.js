import {request} from '@/utils/request.js'

// 获取启用的收款码列表（客户端）
export const getEnabledQrcodePayments = () => {
  return request({
    url: '/qrcodePayment/getEnabledQrcodePayments',
    method: 'get'
  })
}
