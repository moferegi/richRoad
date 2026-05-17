import {request} from '@/utils/request.js'

const withI18nFallbackData = (data = {}, options = {}) => {
  if (options && options.includeI18n) {
    return { ...data, includeI18n: 1 }
  }
  return data
}

// 添加购物车
export const getAllClaimCoupon  = (data, options = {}) => {
    return request({ 
		url:'/cou/getAllClaimCoupon',
        method: 'post',
		data: withI18nFallbackData(data, options)
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

