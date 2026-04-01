import {request} from '@/utils/request.js'

// 获取系统配置（按分组）
export const getSysConfigByGroup = (group) => {
  return request({
    url: '/sysConfig/getSysConfigByGroup',
    method: 'get',
    params: { configGroup: group }
  })
}

// 获取公告配置
export const getAnnouncementConfig = () => {
  return request({
    url: '/sysConfig/getAnnouncementConfig',
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
