import {request} from '@/utils/request.js'

// 获取单个商品详情
export const findGood = (params) => {
    return request({
		url:`/good/findGood?ID=${params}`,
        method: 'get',
    })
}