import { request } from '@/utils/request.js'

/**
 * 获取客服列表（公开接口，无需鉴权）
 * @returns {Promise} 客服列表数据
 */
export const getKefuList = () => {
  return request({
    url: '/kefu/getKefuPublic',
    method: 'get'
  })
}

/**
 * 评价会话（需登录）
 * @param {Object} data - { conversationId: number, rating: number }
 */
export const rateConversation = (data) => {
  return request({
    url: '/cs/conversation/rate',
    method: 'POST',
    data
  })
}

/**
 * 获取平台客服配置（判断是否开启内置客服）
 * @returns {Promise<{code: number, data: {platEnabled: boolean}}>}
 */
export const getCsConfig = () => {
  return request({
    url: '/cs/config/get',
    method: 'GET'
  })
}
