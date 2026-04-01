import {request} from '@/utils/request.js'

// 获取预售商品列表
export const getPresaleGoodList = (params) => {
  return request({
    url: '/presale/getPresaleGoodList',
    method: 'get',
    params
  })
}

// 检查预售是否可参与
export const checkPresaleAvailable = (params) => {
  return request({
    url: '/presale/checkAvailable',
    method: 'get',
    params
  })
}
