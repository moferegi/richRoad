import {request} from '@/utils/request.js'

// 获取活跃弹窗列表（客户端）
export const getActivePopups = (params) => {
  return request({
    url: '/popup/getActivePopups',
    method: 'get',
    params
  })
}
