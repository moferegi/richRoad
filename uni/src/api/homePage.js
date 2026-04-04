import {request} from '@/utils/request.js'

// 获取轮播图列表
export const getBannerList  = () => {
    return request({
		url:'/banner/getBannerList',
        method: 'get',
        data: { isEnabled: true }
    })
}

// 获取分类
export const getCategoryMobile  = (params) => {
    return request({
		url:'/category/getCategoryMobile',
        method: 'get',
        params
    })
}


// 获取商品推荐
export const getGoodList = (params) => {
    return request({
		url:`/good/getGoodList`,
        method: 'get',
        params: { ...params, status: true } 
    })
}

// 获取促销信息
export const getPromotionPublic = () => {
    return request({
        url:'/promo/getPromotionPublic',
        method: 'get'
    })
}


export const getChildrenCategoryAndProduct  = (params) => {
    return request({
		url:'/category/getChildrenCategoryAndProduct',
        method: 'get',
        params
    })
}
