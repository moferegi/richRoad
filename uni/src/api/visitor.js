import { request } from '@/utils/request.js'

/**
 * 访客心跳上报
 * @param {Object} data 心跳数据
 * @returns {Promise}
 */
export const visitorHeartbeat = (data) => {
  return request({
    url: '/visitor/heartbeat',
    method: 'POST',
    data
  })
}
