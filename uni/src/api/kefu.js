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
