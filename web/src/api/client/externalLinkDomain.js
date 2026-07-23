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

/**
 * 检测云存储连接
 * @param {Object} data { id }
 * @returns {Promise}
 */
export const pingCloud = (data) => {
  return service({
    url: '/extDomain/pingCloud',
    method: 'post',
    data: data
  })
}

/**
 * 列出云存储文件
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const listCloudFiles = (params) => {
  return service({
    url: '/extDomain/listCloudFiles',
    method: 'get',
    params: params
  })
}

/**
 * 多云目录比对
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const compareDirectories = (params) => {
  return service({
    url: '/extDomain/compareDirectories',
    method: 'get',
    params: params
  })
}

/**
 * 批量删除云存储文件
 * @param {Object} data { id, keys: string[], password: string }
 * @returns {Promise}
 */
export const deleteCloudFiles = (data) => {
  return service({
    url: '/extDomain/deleteCloudFiles',
    method: 'post',
    data: data
  })
}

/**
 * 上传文件到云存储
 * @param {FormData} formData 含 file、folder、id
 * @returns {Promise}
 */
export const uploadCloudFile = (formData) => {
  return service({
    url: '/extDomain/uploadCloudFile',
    method: 'post',
    headers: { 'Content-Type': 'multipart/form-data' },
    data: formData
  })
}

/**
 * 全局搜索云存储文件
 * @param {Object} params { id, keyword, maxKeys }
 * @returns {Promise}
 */
export const searchCloudFiles = (params) => {
  return service({
    url: '/extDomain/searchCloudFiles',
    method: 'get',
    params: params
  })
}

/**
 * 获取文件下载链接
 * @param {Object} params { id, key }
 * @returns {Promise}
 */
export const getFileDownloadURL = (params) => {
  return service({
    url: '/extDomain/getFileDownloadURL',
    method: 'get',
    params: params
  })
}

/**
 * 获取打包下载目录的 URL（直接打开下载）
 * @param {Object} params { id, prefix }
 * @returns {string} 下载链接
 */
export const getCloudFolderDownloadURL = (params) => {
  const baseURL = service.defaults?.baseURL || ''
  const query = new URLSearchParams(params).toString()
  return `${baseURL}/extDomain/downloadCloudFolder?${query}`
}
