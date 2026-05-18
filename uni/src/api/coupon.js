import {request} from '@/utils/request.js'

const withI18nFallback = (params = {}, options = {}) => {
  if (options && options.includeI18n) {
    return { ...params, includeI18n: 1 }
  }
  return params
}

// 添加购物车
export const getAllClaimCoupon  = (data, options = {}) => {
    return request({ 
		url:'/cou/getAllClaimCoupon',
        method: 'post',
    data,
    params: withI18nFallback({}, options)
    })
}

// 领取优惠券
export const claimCouponByUser  = (data) => {
  return request({ 
  url:'/cou/claimCouponByUser',
      method: 'post',
  data: data
  })
}

