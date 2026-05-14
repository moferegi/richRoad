import service from '@/utils/request'

/**
 * 获取数据库巡检总览
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const getDBInspectorOverview = (params) => {
  return service({
    url: '/dbInspector/getOverview',
    method: 'get',
    params
  })
}

/**
 * 自动修复数据库或Redis连接
 * @param {Object} data 修复参数
 * @returns {Promise}
 */
export const autoFixDBInspector = (data) => {
  return service({
    url: '/dbInspector/autoFix',
    method: 'post',
    data
  })
}

/**
 * 按日期范围真删除记录并联动文件清理
 * @param {Object} data 删除参数
 * @returns {Promise}
 */
export const deleteRecordsByRange = (data) => {
  return service({
    url: '/dbInspector/deleteRecordsByRange',
    method: 'post',
    data
  })
}
