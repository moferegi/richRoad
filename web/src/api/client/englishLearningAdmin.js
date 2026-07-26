import service from '@/utils/request'

// 管理端：签到记录列表
export const getAdminCheckinRecordList = (params) => {
  return service({
    url: '/englishLearning/admin/getCheckinRecordList',
    method: 'get',
    params
  })
}

// 管理端：积分记录列表
export const getAdminPointRecordList = (params) => {
  return service({
    url: '/englishLearning/admin/getPointRecordList',
    method: 'get',
    params
  })
}

// 管理端：时长明细列表
export const getAdminFreeTimeRecordList = (params) => {
  return service({
    url: '/englishLearning/admin/getFreeTimeRecordList',
    method: 'get',
    params
  })
}

// 管理端：观看历史列表
export const getAdminWatchHistoryList = (params) => {
  return service({
    url: '/englishLearning/admin/getWatchHistoryList',
    method: 'get',
    params
  })
}

// 管理端：收藏列表
export const getAdminCollectionList = (params) => {
  return service({
    url: '/englishLearning/admin/getCollectionList',
    method: 'get',
    params
  })
}

// 管理端：错词本列表
export const getAdminWordErrorLogList = (params) => {
  return service({
    url: '/englishLearning/admin/getWordErrorLogList',
    method: 'get',
    params
  })
}

// 管理端：用户搜索（用于下拉选择）
export const getAdminUserList = (params) => {
  return service({
    url: '/englishLearning/admin/getUserList',
    method: 'get',
    params
  })
}
