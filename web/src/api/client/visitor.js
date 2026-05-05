import service from '@/utils/request'

/**
 * 获取访客日志列表
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const getVisitorLogList = (params) => {
  return service({
    url: '/visitor/getVisitorLogList',
    method: 'get',
    params
  })
}

/**
 * 获取访客汇总列表
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const getVisitorSummaryList = (params) => {
  return service({
    url: '/visitor/getVisitorSummaryList',
    method: 'get',
    params
  })
}

/**
 * 获取今日实时统计
 * @returns {Promise}
 */
export const getTodayStats = () => {
  return service({
    url: '/visitor/getTodayStats',
    method: 'get'
  })
}

/**
 * 获取客服引导漏斗统计
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const getKefuGuideStats = (params) => {
  return service({
    url: '/visitor/getKefuGuideStats',
    method: 'get',
    params
  })
}

/**
 * 手动触发日汇总聚合
 * @param {string} date 日期 YYYY-MM-DD
 * @returns {Promise}
 */
export const aggregateDailySummary = (date) => {
  return service({
    url: '/visitor/aggregateDailySummary',
    method: 'post',
    params: { date }
  })
}
