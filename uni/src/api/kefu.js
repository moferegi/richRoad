import { request } from '@/utils/request.js'

const withI18nFallback = (params = {}, options = {}) => {
  if (options && options.includeI18n) {
    return { ...params, includeI18n: 1 }
  }
  return params
}

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
 * 获取客服系统配置（公开接口，用于判断是否启用平台客服）
 * @returns {Promise} { platEnabled: bool }
 */
export const getCsConfig = () => {
  return request({
    url: '/cs/config/get',
    method: 'get'
  })
}

/**
 * 按 key 获取系统配置（公开接口）
 * @param {string} configKey 配置键
 * @param {{ includeI18n?: boolean }} [options] 可选全量多语言回退
 * @returns {Promise} 配置值字符串
 */
export const getSysConfigByKey = (configKey, options = {}) => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: withI18nFallback({ configKey }, options)
  })
}

/**
 * 获取消息历史记录（分页，降序，page=1 为最新一页）
 * @param {{ conversationId: number, page: number, pageSize: number }} params
 * @returns {Promise} PageResult: { list, total, page, pageSize }
 */
export const getMessageHistory = (params) => {
  return request({
    url: '/cs/message/history',
    method: 'get',
    params
  })
}

/**
 * 用户评价会话满意度
 * @param {{ conversationId: number, rating: number }} data rating 1-5
 * @returns {Promise}
 */
export const rateConversation = (data) => {
  return request({
    url: '/cs/conversation/rate',
    method: 'post',
    data
  })
}
