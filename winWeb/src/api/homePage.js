import { request } from '@/utils/request.js'

// 获取轮播图列表
export const getBannerList = () => {
  return request({
    url: '/banner/getBannerList',
    method: 'get',
    data: { isEnabled: true }
  })
}
