import service from '@/utils/request'

const createRequestId = () => {
  return `web_${Date.now()}_${Math.random().toString(16).slice(2, 10)}`
}

/**
 * 创建试衣任务（自动补充requestID，避免重试重复扣币）
 * @param {Object} data 请求参数
 * @param {string} data.sceneType 场景类型 clothes/shoes
 * @param {string} data.sourceImage 用户原图
 * @param {string} data.templateImage 模板图
 * @param {string} [data.requestID] 幂等请求ID，不传则自动生成
 * @returns {Promise}
 */
export const createTryonTask = (data) => {
  const payload = {
    ...data,
  }
  if (!payload.requestID) {
    payload.requestID = createRequestId()
  }
  return service({
    url: '/tryonTask/createTryonTask',
    method: 'post',
    data: payload,
  })
}

/**
 * 按ID查询试衣任务
 * @param {Object} params 查询参数
 * @param {number|string} params.ID 任务ID
 * @returns {Promise}
 */
export const findTryonTask = (params) => {
  return service({
    url: '/tryonTask/findTryonTask',
    method: 'get',
    params,
  })
}

/**
 * 获取我的试衣任务列表
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const getMyTryonTaskList = (params) => {
  return service({
    url: '/tryonTask/getMyTryonTaskList',
    method: 'get',
    params,
  })
}

/**
 * 管理端分页获取试衣任务列表
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const getTryonTaskList = (params) => {
  return service({
    url: '/tryonTask/getTryonTaskList',
    method: 'get',
    params,
  })
}

/**
 * 管理端获取试衣任务统计
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const getTryonTaskStats = (params) => {
  return service({
    url: '/tryonTask/getTryonTaskStats',
    method: 'get',
    params,
  })
}

/**
 * 管理端获取试衣任务趋势
 * @param {Object} params 查询参数
 * @returns {Promise}
 */
export const getTryonTaskTrend = (params) => {
  return service({
    url: '/tryonTask/getTryonTaskTrend',
    method: 'get',
    params,
  })
}

/**
 * 管理端删除试衣任务
 * @param {Object} params 查询参数
 * @param {number|string} params.ID 任务ID
 * @returns {Promise}
 */
export const deleteTryonTask = (params) => {
  return service({
    url: '/tryonTask/deleteTryonTask',
    method: 'delete',
    params,
  })
}

/**
 * 管理端批量删除试衣任务
 * @param {Object} params 查询参数
 * @param {Array<number|string>} params['IDs[]'] 任务ID列表
 * @returns {Promise}
 */
export const deleteTryonTaskByIds = (params) => {
  return service({
    url: '/tryonTask/deleteTryonTaskByIds',
    method: 'delete',
    params,
  })
}

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
