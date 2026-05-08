import service from '@/utils/request'

/**
 * 管理端分页获取模特列表
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const getTryonModelList = (params) => {
  return service({
    url: '/tryonModel/getTryonModelList',
    method: 'get',
    params,
  })
}

/**
 * 管理端删除模特
 * @param {Object} params 查询参数
 * @param {number|string} params.ID 模特ID
 * @returns {Promise}
 */
export const deleteTryonModel = (params) => {
  return service({
    url: '/tryonModel/deleteTryonModel',
    method: 'delete',
    params,
  })
}

/**
 * 管理端批量删除模特
 * @param {Object} params 查询参数
 * @param {Array<number|string>} params['IDs[]'] 模特ID列表
 * @returns {Promise}
 */
export const deleteTryonModelByIds = (params) => {
  return service({
    url: '/tryonModel/deleteTryonModelByIds',
    method: 'delete',
    params,
  })
}
