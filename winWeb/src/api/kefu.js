import { request } from '@/utils/request.js'

/**
 * 获取客服列表（公开接口，无需鉴权）
 */
export const getKefuList = () => {
  return request({
    url: '/kefu/getKefuPublic',
    method: 'get'
  })
}

/**
 * 获取客服系统配置（公开接口）
 */
export const getCsConfig = () => {
  return request({
    url: '/cs/config/get',
    method: 'get'
  })
}

/**
 * 按 key 获取系统配置（公开接口）
 */
export const getSysConfigByKey = (configKey) => {
  return request({
    url: '/sysConfig/getSysConfigByKey',
    method: 'get',
    params: { configKey }
  })
}
