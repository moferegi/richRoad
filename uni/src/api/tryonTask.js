import { request } from '@/utils/request.js'

const createRequestId = () => {
  return `uni_${Date.now()}_${Math.random().toString(16).slice(2, 10)}`
}

// 创建试衣任务（自动补充requestID，避免重试重复扣币）
export const createTryonTask = (data) => {
  const payload = {
    ...data,
  }
  if (!payload.requestID) {
    payload.requestID = createRequestId()
  }
  return request({
    url: '/tryonTask/createTryonTask',
    method: 'post',
    data: payload,
  })
}

// 按ID查询试衣任务
export const findTryonTask = (params) => {
  return request({
    url: '/tryonTask/findTryonTask',
    method: 'get',
    params,
  })
}

// 获取我的试衣任务列表
export const getMyTryonTaskList = (params) => {
  return request({
    url: '/tryonTask/getMyTryonTaskList',
    method: 'get',
    params,
  })
}

// 创建我的模特
export const createTryonModel = (data) => {
  return request({
    url: '/tryonModel/createTryonModel',
    method: 'post',
    data,
  })
}

// 重命名我的模特
export const updateTryonModel = (data) => {
  return request({
    url: '/tryonModel/updateTryonModel',
    method: 'put',
    data,
  })
}

// 删除我的模特
export const deleteTryonModel = (params) => {
  return request({
    url: '/tryonModel/deleteTryonModel',
    method: 'delete',
    params,
  })
}

// 获取我的模特列表
export const getMyTryonModelList = () => {
  return request({
    url: '/tryonModel/getMyTryonModelList',
    method: 'get',
  })
}
