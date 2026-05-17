import {request} from '@/utils/request.js'

const withI18nFallback = (params = {}, options = {}) => {
  if (options && options.includeI18n) {
    return { ...params, includeI18n: 1 }
  }
  return params
}

const withConfigKey = (configKey, options = {}) => withI18nFallback({ configKey }, options)

// 获取系统配置（按分组）
export const getSysConfigByGroup = (group, options = {}) => {
  return request({
    url: '/sysConfig/getSysConfigByGroup',
    method: 'get',
    params: withI18nFallback({ configGroup: group }, options)
  })
}

// 获取单个配置值
export const getSysConfigByKey = (configKey, options = {}) => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: withConfigKey(configKey, options)
  })
}

// 获取公告配置
export const getAnnouncementConfig = (options = {}) => {
  return request({
    url: '/sysConfig/getAnnouncementConfig',
    method: 'get',
    params: withI18nFallback({}, options)
  })
}

// 获取登录相关配置（公开接口，无需登录）
export const getLoginConfig = (options = {}) => {
  return request({
    url: '/sysConfig/getLoginConfig',
    method: 'get',
    params: withI18nFallback({}, options)
  })
}

// 获取签到是否开启
export const getSignInEnabled = (options = {}) => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: withConfigKey('sign_in_enabled', options)
  })
}

// 获取预售首页展示数量
export const getPresaleHomeCount = (options = {}) => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: withConfigKey('presale_home_count', options)
  })
}

// 获取货币后缀
export const getCurrencySuffix = (options = {}) => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: withConfigKey('currency_suffix', options)
  })
}

// 获取货币符号
export const getCurrencySymbol = (options = {}) => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: withConfigKey('currency_symbol', options)
  })
}

// 获取应用名称
export const getAppName = (options = {}) => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: withConfigKey('app_name', options)
  })
}

// 获取应用Logo
export const getAppLogo = (options = {}) => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: withConfigKey('app_logo', options)
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
export const getPointsExchangeRate = (options = {}) => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: withConfigKey('points_exchange_rate', options)
  })
}

// 获取游客初始试衣币
export const getTryonGuestInitPoints = (options = {}) => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: withConfigKey('tryon_guest_init_points', options)
  })
}

// 获取注册奖励试衣币
export const getTryonRegisterRewardPoints = (options = {}) => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: withConfigKey('tryon_register_reward_points', options)
  })
}

// 获取单次试衣消耗
export const getTryonCostPoints = (options = {}) => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: withConfigKey('tryon_cost_points', options)
  })
}

// 获取试衣币充值套餐
export const getTryonRechargePlans = (options = {}) => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: withConfigKey('tryon_recharge_plans', options)
  })
}

// 获取试衣配置（游客初始币、注册奖励、单次消耗、失败退币比例）
export const getTryonConfig = (options = {}) => {
  return request({
    url: '/sysConfig/getTryonConfig',
    method: 'get',
    params: withI18nFallback({}, options)
  })
}

// 获取支付方式配置（人工优先 + 渠道开关）
export const getPaymentConfig = (options = {}) => {
  return request({
    url: '/sysConfig/getPaymentConfig',
    method: 'get',
    params: withI18nFallback({}, options)
  })
}

// 获取 uni 联系客服场景的期望支付方式配置（独立于 payment_manual_methods）
export const getUniPreferredPayConfig = (options = {}) => {
  return request({
    url: '/sysConfig/getUniPreferredPayConfig',
    method: 'get',
    params: withI18nFallback({}, options)
  })
}
