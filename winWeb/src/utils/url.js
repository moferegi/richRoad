import { baseUrl } from '@/utils/request.js'

// 缓存外部域名
let cachedDomain = ''
let domainLoading = false
let domainPromise = null

// CDN域名缓存
let cdnDomain = ''

/**
 * 初始化外部域名（从后端获取默认外部资源域名）
 */
export const initExternalDomain = async () => {
  if (cachedDomain || domainLoading) {
    if (domainPromise) return domainPromise
    return
  }
  domainLoading = true
  domainPromise = (async () => {
    try {
      const { getDefaultDomain } = await import('@/api/sysConfig.js')
      const res = await getDefaultDomain()
      if (res.code === 0 && res.data) {
        cachedDomain = res.data.replace(/\/+$/, '')
      }
    } catch (e) {
      console.warn('获取外部域名失败', e)
    } finally {
      domainLoading = false
      domainPromise = null
    }
  })()
  return domainPromise
}

/**
 * 获取外部链接完整URL：如果是相对路径，自动拼接外部默认域名
 */
export const getExternalUrl = (url) => {
  if (!url) return ''
  const str = String(url)
  if (str.slice(0, 4) === 'http') return str
  if (str.slice(0, 4) === 'data') return str
  if (cachedDomain) {
    const sep = str.slice(0, 1) === '/' ? '' : '/'
    return cachedDomain + sep + str
  }
  // 兜底使用 getUrl
  return getUrl(str)
}

export const getUrl = (url) => {
  if (!url) return ''
  const path = baseUrl
  if (url.slice(0, 4) === 'data') return url
  if (url && url.slice(0, 4) !== 'http') {
    if (path === '/') return url
    if (url.slice(0, 1) === '/') return path + url
    return path + '/' + url
  }
  return url
}

/**
 * 获取CDN资源URL
 */
export const getCdnUrl = (url) => {
  if (!url) return ''
  const str = String(url)
  if (str.slice(0, 4) === 'http') return str
  if (str.slice(0, 4) === 'data') return str
  if (cdnDomain) {
    const sep = str.slice(0, 1) === '/' ? '' : '/'
    return cdnDomain + sep + str
  }
  return getUrl(str)
}

export const setCdnDomain = (domain) => {
  cdnDomain = domain ? domain.replace(/\/+$/, '') : ''
}
