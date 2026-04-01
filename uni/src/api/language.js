import {request} from '@/utils/request.js'

// 获取启用的语言列表
export const getEnabledLanguages = () => {
  return request({
    url: '/language/getEnabledLanguages',
    method: 'get'
  })
}
