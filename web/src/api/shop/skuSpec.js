import service from '@/utils/request'

/**
 * 创建SKU规格字典
 * @param {Object} data SKU规格字典数据
 * @returns {Promise}
 */
export const createSkuSpec = (data) => {
  return service({
    url: '/skuSpec/createSkuSpec',
    method: 'post',
    data
  })
}

/**
 * 删除SKU规格字典
 * @param {Object} params 删除参数
 * @returns {Promise}
 */
export const deleteSkuSpec = (params) => {
  return service({
    url: '/skuSpec/deleteSkuSpec',
    method: 'delete',
    params
  })
}

/**
 * 批量删除SKU规格字典
 * @param {Object} data 批量删除参数
 * @returns {Promise}
 */
export const deleteSkuSpecByIds = (data) => {
  return service({
    url: '/skuSpec/deleteSkuSpecByIds',
    method: 'delete',
    data
  })
}

/**
 * 更新SKU规格字典
 * @param {Object} data SKU规格字典数据
 * @returns {Promise}
 */
export const updateSkuSpec = (data) => {
  return service({
    url: '/skuSpec/updateSkuSpec',
    method: 'put',
    data
  })
}

/**
 * 根据ID获取SKU规格字典
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const findSkuSpec = (params) => {
  return service({
    url: '/skuSpec/findSkuSpec',
    method: 'get',
    params
  })
}

/**
 * 分页获取SKU规格字典列表
 * @param {Object} params 分页查询参数
 * @returns {Promise}
 */
export const getSkuSpecList = (params) => {
  return service({
    url: '/skuSpec/getSkuSpecList',
    method: 'get',
    params
  })
}

/**
 * 获取所有SKU规格字典（下拉选择用）
 * @param {string} type 类型：spec或attr
 * @returns {Promise}
 */
export const getAllSkuSpecs = (type) => {
  return service({
    url: '/skuSpec/getAllSkuSpecs',
    method: 'get',
    params: { type }
  })
}
