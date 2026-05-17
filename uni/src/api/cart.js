import {request} from '@/utils/request.js'

const withI18nFallback = (params = {}, options = {}) => {
    if (options && options.includeI18n) {
        return { ...params, includeI18n: 1 }
    }
    return params
}

// 添加购物车
export const addCart  = (data) => {
    return request({ 
		url:'/cart/addCart',
        method: 'post',
		data: data
    })
}

// 删除购物车
export const cutCart  = (data) => {
    return request({ 
		url:'/cart/cutCart',
        method: 'post',
		data: data
    })
}

// 删除全部购物车
export const clearCart  = () => {
    return request({ 
		url:'/cart/clearCart',
        method: 'get'
    })
}

// 获取自身购物车 
export const getSelfCart  = (options = {}) => {
    return request({
		url:'/cart/getSelfCart',
        method: 'get',
        params: withI18nFallback({}, options)
    })
}