import { request } from '@/utils/request.js'

// 获取系统配置（按分组）
export const getSysConfigByGroup = (group) => {
  return request({
    url: '/sysConfig/getSysConfigByGroup',
    method: 'get',
    params: { configGroup: group }
  })
}

// 获取单个配置值
export const getSysConfigByKey = (configKey) => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: { configKey }
  })
}

// 获取登录相关配置（公开接口）
export const getLoginConfig = () => {
  return request({
    url: '/sysConfig/getLoginConfig',
    method: 'get'
  })
}

// 获取积分兑换比率
export const getPointsExchangeRate = () => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: { configKey: 'points_exchange_rate' }
  })
}

// 获取默认外部链接域名
export const getDefaultDomain = () => {
  return request({
    url: '/extDomain/getDefaultDomain',
    method: 'get'
  })
}

// 获取应用名称
export const getAppName = () => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: { configKey: 'app_name' }
  })
}

// 获取应用Logo
export const getAppLogo = () => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: { configKey: 'app_logo' }
  })
}
