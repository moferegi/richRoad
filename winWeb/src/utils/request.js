import axios from 'axios'
import storage from './storage'
import router from '@/router'

// 导出 baseUrl 供其他模块使用（如 url.js、player 的 m3u8 重写）
export const baseURL = import.meta.env.VITE_API_BASE_URL
export let baseUrl = baseURL

const service = axios.create({
  baseURL,
  timeout: 30000
})

// 成功 code（与 uni 端、后端保持一致：0 表示成功）
const SUCCESS_CODE = 0

// 请求拦截
service.interceptors.request.use(
  (config) => {
    const token = storage.get('x-token')
    if (token) {
      config.headers['x-token'] = token
    }
    // 多语言 header
    const lang = storage.get('app-lang') || 'zh'
    config.headers['Accept-Language'] = lang
    // 客户端标识
    config.headers['X-Client-Platform'] = 'winweb'

    // 过滤空值参数
    if (config.params) {
      Object.keys(config.params).forEach((key) => {
        if (config.params[key] === '' || config.params[key] === null || config.params[key] === undefined) {
          delete config.params[key]
        }
      })
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截
service.interceptors.response.use(
  (response) => {
    // 处理 new-token（与 uni 端逻辑一致）
    const newToken = response.headers['new-token']
    if (newToken) {
      storage.set('x-token', newToken)
    }

    // HTTP 401：未授权
    if (response.status === 401) {
      storage.remove('x-token')
      storage.remove('userInfo')
      router.push('/login')
      return Promise.reject(new Error('Unauthorized'))
    }

    const res = response.data

    // 业务 code 为 0 表示成功，直接返回
    if (res.code === SUCCESS_CODE) {
      return res
    }

    // 业务错误：统一提示
    const msg = res.msg || 'Request failed'
    console.warn('[API Error]', msg, res)

    // 业务层面的 401 也跳转登录
    if (res.code === 401) {
      storage.remove('x-token')
      storage.remove('userInfo')
      router.push('/login')
    }

    return Promise.reject(new Error(msg))
  },
  (error) => {
    console.error('request error', error)

    // HTTP 401
    if (error.response?.status === 401) {
      storage.remove('x-token')
      storage.remove('userInfo')
      router.push('/login')
    }

    const msg = error.response?.data?.msg || error.message || 'Network Error'
    return Promise.reject(new Error(msg))
  }
)

export default service

/**
 * 请求方法封装
 */
export function request(options = {}) {
  const { url, method = 'GET', data, params, headers } = options
  return service({
    url,
    method,
    data,
    params,
    headers
  })
}
