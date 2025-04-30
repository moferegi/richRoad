import {request} from '@/utils/request.js'


// 获取是否收藏
export const findCollect  = (params) => {
    return request({
		url:`/collect/findCollect?goodID=${params.goodID}`,
        method: 'get'
    })
}

// 分页获取收藏列表
export const getCollectList  = (params) => {
    return request({
		url:`/collect/getCollectList?page=${params.page}&pageSize=${params.pageSize}`,
        method: 'get'
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
