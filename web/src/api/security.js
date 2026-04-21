import service from '@/utils/request'

/**
 * 封禁指定IP
 * @param {Object} data
 * @param {string} data.ip - IP地址
 * @param {string} data.reason - 封禁原因
 * @param {number} data.duration - 封禁时长（分钟），0=永久
 */
export const banIP = (data) => {
  return service({
    url: '/sysBannedIP/banIP',
    method: 'post',
    data
  })
}

/**
 * 解封指定IP
 * @param {Object} data
 * @param {string} data.ip - IP地址
 */
export const unbanIP = (data) => {
  return service({
    url: '/sysBannedIP/unbanIP',
    method: 'post',
    data
  })
}

/**
 * 分页获取封禁IP列表
 * @param {Object} params
 * @param {number} params.page - 页码
 * @param {number} params.pageSize - 每页数量
 * @param {string} [params.ip] - 过滤IP
 * @param {boolean} [params.isAuto] - 是否自动封禁
 */
export const getBannedIPList = (params) => {
  return service({
    url: '/sysBannedIP/getBannedIPList',
    method: 'get',
    params
  })
}

/**
 * 获取IP攻击统计
 * @param {Object} params
 * @param {number} [params.hours] - 统计最近N小时，默认24
 */
export const getAttackStats = (params) => {
  return service({
    url: '/sysBannedIP/getAttackStats',
    method: 'get',
    params
  })
}
