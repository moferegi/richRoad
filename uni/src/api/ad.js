import { request } from '@/utils/request.js'

/**
 * 根据位置标记获取可用广告列表
 * @param {string} positionKey 位置标记，如 pages/game/category
 */
export const getAdByPosition = (positionKey) => {
  return request({
    url: '/ad/getAdByPosition',
    method: 'GET',
    params: { positionKey }
  })
}

/**
 * 上报观看记录
 * @param {Object} data { adVideoId, positionId, watchedSeconds }
 */
export const reportWatch = (data) => {
  return request({
    url: '/ad/reportWatch',
    method: 'POST',
    data
  })
}