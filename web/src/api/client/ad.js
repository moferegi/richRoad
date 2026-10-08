import service from '@/utils/request'

// ========== 广告位置 ==========

export const createAdPosition = (data) => {
  return service({ url: '/ad/createAdPosition', method: 'post', data })
}

export const deleteAdPosition = (data) => {
  return service({ url: '/ad/deleteAdPosition', method: 'delete', data })
}

export const updateAdPosition = (data) => {
  return service({ url: '/ad/updateAdPosition', method: 'put', data })
}

export const getAdPositionList = (params) => {
  return service({ url: '/ad/getAdPositionList', method: 'get', params })
}

// ========== 广告视频 ==========

export const createAdVideo = (data) => {
  return service({ url: '/ad/createAdVideo', method: 'post', data })
}

export const deleteAdVideo = (data) => {
  return service({ url: '/ad/deleteAdVideo', method: 'delete', data })
}

export const updateAdVideo = (data) => {
  return service({ url: '/ad/updateAdVideo', method: 'put', data })
}

export const getAdVideoList = (params) => {
  return service({ url: '/ad/getAdVideoList', method: 'get', params })
}

export const sliceAdVideo = (formData) => {
  return service({
    url: '/ad/sliceAdVideo',
    method: 'post',
    headers: { 'Content-Type': 'multipart/form-data' },
    data: formData
  })
}

export const checkFfmpeg = () => {
  return service({ url: '/ad/checkFfmpeg', method: 'get' })
}

// ========== 观看记录 ==========

export const getWatchRecordList = (params) => {
  return service({ url: '/ad/getWatchRecordList', method: 'get', params })
}

export const getWatchStats = (params) => {
  return service({ url: '/ad/getWatchStats', method: 'get', params })
}