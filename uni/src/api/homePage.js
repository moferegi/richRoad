import {request} from '@/utils/request.js'

const withI18nFallback = (params = {}, options = {}) => {
    if (options && options.includeI18n) {
        return { ...params, includeI18n: 1 }
    }
    return params
}

// 获取轮播图列表
export const getBannerList  = (options = {}) => {
    return request({
		url:'/banner/getBannerList',
        method: 'get',
        params: withI18nFallback({ isEnabled: true }, options)
    })
}

// 获取分类
export const getCategoryMobile  = (params = {}, options = {}) => {
    return request({
		url:'/category/getCategoryMobile',
        method: 'get',
        params: withI18nFallback(params, options)
    })
}


// 获取商品推荐
export const getGoodList = (params = {}, options = {}) => {
    return request({
		url:`/good/getGoodList`,
        method: 'get',
        params: withI18nFallback({ ...params, status: true, excludeHiddenCategories: true, isPresale: false }, options)
    })
}

// 获取促销信息
export const getPromotionPublic = () => {
    return request({
        url:'/promo/getPromotionPublic',
        method: 'get'
    })
}


export const getChildrenCategoryAndProduct  = (params = {}, options = {}) => {
    return request({
		url:'/category/getChildrenCategoryAndProduct',
        method: 'get',
        params: withI18nFallback(params, options)
    })
}
