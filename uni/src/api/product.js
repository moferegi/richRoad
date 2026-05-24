import {request} from '@/utils/request.js'

// 获取单个商品详情
export const findGood = (params, options = {}) => {
    return request({
		url:'/good/getGoodPublic',
        method: 'get',
        params: {
            ID: params,
            ...(options && options.includeI18n ? { includeI18n: 1 } : {})
        }
    })
}
