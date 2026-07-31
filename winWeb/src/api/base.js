import { request } from '@/utils/request.js'

// 获取图形验证码
export const getCaptcha = () => {
  return request({
    url: '/base/captcha',
    method: 'post'
  })
}

// 用户名注册
export const register = (data) => {
  return request({
    url: '/clientUser/register',
    method: 'post',
    data
  })
}

// 用户名登录
export const login = (data) => {
  return request({
    url: '/clientUser/login',
    method: 'post',
    data
  })
}

// 获取用户信息
export const getUserInfo = () => {
  return request({
    url: '/clientUser/getUserInfo',
    method: 'get'
  })
}

// 手机号登录
export const phoneLogin = (data) => {
  return request({
    url: '/clientUser/phoneLogin',
    method: 'post',
    data
  })
}

// 手机号注册
export const phoneRegister = (data) => {
  return request({
    url: '/clientUser/phoneRegister',
    method: 'post',
    data
  })
}

// 修改用户信息
export const setClientUserInfo = (data) => {
  return request({
    url: '/clientUser/setClientUserInfo',
    method: 'post',
    data
  })
}

// 修改密码
export const changePassword = (data) => {
  return request({
    url: '/clientUser/changePassword',
    method: 'post',
    data
  })
}

// 获取我的邀请信息
export const getMyInviteInfo = () => {
  return request({
    url: '/clientUser/getMyInviteInfo',
    method: 'get'
  })
}
