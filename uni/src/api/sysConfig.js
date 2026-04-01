import {request} from '@/utils/request.js'

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

// 获取公告配置
export const getAnnouncementConfig = () => {
  return request({
    url: '/sysConfig/getAnnouncementConfig',
    method: 'get'
  })
}

// 获取登录相关配置（公开接口，无需登录）
export const getLoginConfig = () => {
  return request({
    url: '/sysConfig/getLoginConfig',
    method: 'get'
  })
}

// 获取签到是否开启
export const getSignInEnabled = () => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: { configKey: 'sign_in_enabled' }
  })
}

// 获取预售首页展示数量
export const getPresaleHomeCount = () => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: { configKey: 'presale_home_count' }
  })
}

// 获取货币后缀
export const getCurrencySuffix = () => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: { configKey: 'currency_suffix' }
  })
}

// 获取货币符号
export const getCurrencySymbol = () => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: { configKey: 'currency_symbol' }
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

// 获取默认外部链接域名
export const getDefaultDomain = () => {
  return request({
    url: '/extDomain/getDefaultDomain',
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
