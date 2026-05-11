import service from '@/utils/request'

/**
 * 获取系统参数列表
 * @param {Object} data 查询参数
 * @returns {Promise}
 */
export const getSysConfigList = (data) => {
  return service({
    url: '/sysConfig/getSysConfigList',
    method: 'get',
    params: data
  })
}

/**
 * 更新系统参数
 * @param {Object} data { id, configValue, remark }
 * @returns {Promise}
 */
export const updateSysConfig = (data) => {
  return service({
    url: '/sysConfig/updateSysConfig',
    method: 'put',
    data: data
  })
}

/**
 * 获取阿里试衣模型剩余额度估算（管理端）
 * @param {Object} params 查询参数
 * @param {string} params.modelKey 模型key，不传返回全部
 * @returns {Promise}
 */
export const getAliyunTryonQuotaEstimate = (params = {}) => {
  return service({
    url: '/sysConfig/getAliyunTryonQuotaEstimate',
    method: 'get',
    params
  })
}

/**
 * 分页获取模型调用日志
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const getModelCallLogList = (params = {}) => {
  return service({
    url: '/sysConfig/getModelCallLogList',
    method: 'get',
    params
  })
}
