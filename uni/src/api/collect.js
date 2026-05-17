import {request} from '@/utils/request.js'

const withI18nFallback = (params = {}, options = {}) => {
    if (options && options.includeI18n) {
        return { ...params, includeI18n: 1 }
    }
    return params
}


// 获取是否收藏
export const findCollect  = (params, options = {}) => {
    return request({
		url:'/collect/findCollect',
        method: 'get',
        params: withI18nFallback({ goodID: params.goodID }, options)
    })
}

// 分页获取收藏列表
export const getCollectList  = (params, options = {}) => {
    return request({
		url:'/collect/getCollectList',
        method: 'get',
        params: withI18nFallback({
            page: params.page,
            pageSize: params.pageSize,
        }, options)
    })
}


// 收藏和取消收藏
export const createCollect  = (params) => {
    return request({
		url: "/collect/createCollect",
        method: 'post',
		data: params
    })
}
