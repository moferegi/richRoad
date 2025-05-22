import {request} from '@/utils/request.js'

// 添加购物车
export const getAllClaimCoupon  = (data) => {
    return request({ 
		url:'/cou/getAllClaimCoupon',
        method: 'post',
		data: data
    })
}
