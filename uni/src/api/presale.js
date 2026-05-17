import {request} from '@/utils/request.js'

const withI18nFallback = (params = {}, options = {}) => {
  if (options && options.includeI18n) {
    return { ...params, includeI18n: 1 }
  }
  return params
}

// 获取预售商品列表
export const getPresaleGoodList = (params = {}, options = {}) => {
  return request({
    url: '/presale/getPresaleGoodList',
    method: 'get',
    params: withI18nFallback(params, options)
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
