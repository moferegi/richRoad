import {request} from '@/utils/request.js'

// 获取启用的国际区号列表
export const getEnabledPhoneAreaCodes = () => {
  return request({
    url: '/phoneAreaCode/getEnabledPhoneAreaCodes',
    method: 'get'
  })
}
