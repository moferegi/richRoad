import {request} from '@/utils/request.js'

// 执行签到
export const doSignIn = () => {
  return request({
    url: '/signIn/doSignIn',
    method: 'post'
  })
}

// 获取签到状态（今日是否已签到、连续天数等）
export const getSignInStatus = () => {
  return request({
    url: '/signIn/getSignInStatus',
    method: 'get'
  })
}
