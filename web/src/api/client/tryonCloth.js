import service from '@/utils/request'

/**
 * 管理端分页获取我的衣橱列表
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const getTryonClothList = (params) => {
  return service({
    url: '/tryonCloth/getTryonClothList',
    method: 'get',
    params,
  })
}

/**
 * 管理端删除我的衣橱
 * @param {Object} params 查询参数
 * @param {number|string} params.ID 衣橱ID
 * @returns {Promise}
 */
export const deleteTryonCloth = (params) => {
  return service({
    url: '/tryonCloth/deleteTryonCloth',
    method: 'delete',
    params,
  })
}

/**
 * 管理端批量删除我的衣橱
 * @param {Object} params 查询参数
 * @param {Array<number|string>} params['IDs[]'] 衣橱ID列表
 * @returns {Promise}
 */
export const deleteTryonClothByIds = (params) => {
  return service({
    url: '/tryonCloth/deleteTryonClothByIds',
    method: 'delete',
    params,
  })
}
