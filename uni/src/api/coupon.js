import {request} from '@/utils/request.js'

// 添加购物车
export const getAllClaimCoupon  = (data) => {
    return request({ 
		url:'/cou/getAllClaimCoupon',
        method: 'post',
		data: data
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

