import service from '@/utils/request'

/**
 * 创建外部链接域名
 * @param {Object} data 域名信息
 * @returns {Promise}
 */
export const createExternalLinkDomain = (data) => {
  return service({
    url: '/extDomain/createExternalLinkDomain',
    method: 'post',
    data: data
  })
}

/**
 * 删除外部链接域名
 * @param {Object} data { id }
 * @returns {Promise}
 */
export const deleteExternalLinkDomain = (data) => {
  return service({
    url: '/extDomain/deleteExternalLinkDomain',
    method: 'delete',
    data: data
  })
}

/**
 * 更新外部链接域名
 * @param {Object} data 域名信息
 * @returns {Promise}
 */
export const updateExternalLinkDomain = (data) => {
  return service({
    url: '/extDomain/updateExternalLinkDomain',
    method: 'put',
    data: data
  })
}

/**
 * 获取外部链接域名列表
 * @param {Object} data 分页参数
 * @returns {Promise}
 */
export const getExternalLinkDomainList = (data) => {
  return service({
    url: '/extDomain/getExternalLinkDomainList',
    method: 'get',
    params: data
  })
}

/**
 * 设置默认域名
 * @param {Object} data { id }
 * @returns {Promise}
 */
export const setDefaultDomain = (data) => {
  return service({
    url: '/extDomain/setDefaultDomain',
    method: 'post',
    data: data
  })
}

/**
 * 获取默认域名（公开接口）
 * @returns {Promise}
 */
export const getDefaultDomain = () => {
  return service({
    url: '/extDomain/getDefaultDomain',
    method: 'get'
  })
}
