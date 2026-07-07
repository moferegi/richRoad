import { myRouter }from  '@/utils/permission.js'
import { localText, t } from '@/utils/i18n.js'
import CryptoJS from 'crypto-js'

const BACKEND_ERROR_KEYS = [
    'captchaRateLimit',
    'oldPasswordWrong',
    'userNotFound',
    'passwordWrong',
    'phoneAlreadyRegistered',
    'phoneFormatInvalid',
    'phoneOrPasswordWrong',
    'usernameOrPasswordWrong',
    'loginLocked',
    'registerIPLimit',
    'alreadySigned',
    'aliyunBeautifyAccessKeyMissing',
    'modelAliyunRetouchNeedStandardOSS',
    'reason_tryonBeautifyDeduct',
    'reason_tryonBeautifyRefund',
    'tryonBeautifyAlreadyUsed',
    'tryonBeautifyEndpointInvalid',
    'tryonBeautifyRefundFailedContactAdmin',
    'tryonBeautifyRequestFailed',
    'tryonBeautifyResponseInvalid',
    'tryonBeautifyResultEmpty',
    'tryonBeautifyStatusConflict',
    'tryonBeautifyUnsupported',
    'tryonClothCategoryInvalid',
    'tryonClothNameRequired',
    'tryonClothNotFoundOrNoPermission',
    'aliyunParsingModelUnavailable',
    'dbNotInitialized',
    'modelTokenQuotaExhausted',
    'pointRecordIDMissing',
    'tryonParsingRefundFailedContactAdmin',
    'tryonRefinerModelUnavailable',
    'tryonResultImageRequired',
]

const BACKEND_ERROR_ALIAS_LANGS = ['zh', 'zh-TW', 'en', 'mn', 'th', 'hi', 'id', 'vi', 'ar', 'ja', 'ko', 'ms']

const buildBackendErrorAliasMap = () => {
    const aliasMap = {}

    BACKEND_ERROR_KEYS.forEach((key) => {
        aliasMap[key] = key

        BACKEND_ERROR_ALIAS_LANGS.forEach((lang) => {
            const text = String(t(key, lang) || '').trim()
            if (!text) return
            aliasMap[text] = key
            aliasMap[text.replace(/\s+/g, '')] = key
        })
    })

    return aliasMap
}

const backendErrorAliasMap = buildBackendErrorAliasMap()

const normalizeBackendMessage = (value) => {
    if (value === null || value === undefined) {
        return ''
    }

    if (typeof value === 'object') {
        return String(localText(value) || '').trim()
    }

    const rawMsg = String(value).trim()
    if (!rawMsg) {
        return ''
    }

    if (rawMsg.charAt(0) === '{') {
        const localized = String(localText(rawMsg) || '').trim()
        if (localized && localized !== rawMsg) {
            return localized
        }
    }

    const normalizedMsg = rawMsg.replace(/\s+/g, '')
    const i18nKey = backendErrorAliasMap[rawMsg] || backendErrorAliasMap[normalizedMsg]
    if (i18nKey) {
        return t(i18nKey)
    }

    const translated = t(rawMsg)
    return translated !== rawMsg ? translated : rawMsg
}

const hasOwn = Object.prototype.hasOwnProperty
const AUTH_ROUTES = ['/pages/user/login', '/pages/user/register']

const safeDecode = (value) => {
    if (value === null || value === undefined) return ''
    const text = String(value).trim()
    if (!text) return ''
    try {
        return decodeURIComponent(text)
    } catch (e) {
        return text
    }
}

const getRouteFromH5Location = () => {
    // #ifdef H5
    try {
        const hash = String(window.location.hash || '')
        if (hash.startsWith('#/')) {
            const route = hash.slice(1).split('?')[0]
            return route || ''
        }
        const pathname = String(window.location.pathname || '').trim()
        if (pathname.startsWith('/pages/')) {
            return pathname.split('?')[0]
        }
    } catch (e) {
        return ''
    }
    // #endif
    return ''
}

const getRouteFromLaunchOptions = () => {
    try {
        if (typeof uni.getLaunchOptionsSync !== 'function') {
            return ''
        }
        const launch = uni.getLaunchOptionsSync() || {}
        const path = String(launch.path || '').trim()
        if (!path) {
            return ''
        }
        return path.startsWith('/') ? path : `/${path}`
    } catch (e) {
        return ''
    }
}

const getCurrentRoutePath = () => {
    const pages = getCurrentPages()
    if (pages && pages.length > 0) {
        const route = pages[pages.length - 1]?.route || ''
        return route ? `/${route}` : ''
    }
    return getRouteFromH5Location() || getRouteFromLaunchOptions() || ''
}

const isAuthRoute = (routePath = '') => {
    return AUTH_ROUTES.some(item => routePath.startsWith(item))
}

const getInviteCodeFromCurrentPage = () => {
    try {
        const pages = getCurrentPages()
        if (!pages || pages.length === 0) return ''
        const options = pages[pages.length - 1]?.options || {}
        return safeDecode(options.inviteCode || options.invite_code || options.code || '')
    } catch (e) {
        return ''
    }
}

const getInviteCodeFromH5Location = () => {
    // #ifdef H5
    try {
        const fromSearch = new URLSearchParams(window.location.search || '').get('inviteCode')
        if (fromSearch) return safeDecode(fromSearch)
        const hash = window.location.hash || ''
        const hashQuery = hash.includes('?') ? hash.slice(hash.indexOf('?') + 1) : ''
        const fromHashQuery = new URLSearchParams(hashQuery).get('inviteCode')
        if (fromHashQuery) return safeDecode(fromHashQuery)
    } catch (e) {
        return ''
    }
    // #endif
    return ''
}

const resolveInviteCode = () => {
    return (
        getInviteCodeFromCurrentPage() ||
        getInviteCodeFromH5Location() ||
        safeDecode(uni.getStorageSync('pendingInviteCode') || '')
    )
}

const buildLoginRedirectUrl = () => {
    const inviteCode = resolveInviteCode()
    if (!inviteCode) return '/pages/user/login'
    uni.setStorageSync('pendingInviteCode', inviteCode)
    return `/pages/user/login?inviteCode=${encodeURIComponent(inviteCode)}`
}

const redirectToLogin = () => {
    uni.removeStorageSync('x-token')
    uni.removeStorageSync('userInfo')
    const currentRoute = getCurrentRoutePath()
    if (isAuthRoute(currentRoute)) {
        return false
    }
    uni.reLaunch({ url: buildLoginRedirectUrl() })
    return true
}

const isApiResponseObject = (payload) => {
    return !!payload && typeof payload === 'object' && !Array.isArray(payload) && hasOwn.call(payload, 'code')
}

const normalizeApiResponse = (payload) => {
    if (isApiResponseObject(payload)) {
        return payload
    }
    return {
        code: -1,
        data: null,
        msg: t('apiResponseInvalid'),
    }
}

// 定义并导出 baseUrl 变量
export let baseUrl = '/api'
if (process.env.NODE_ENV === 'development') {
    // 开发环境：继续使用 Vite 代理，解决本地跨域
    // #ifdef H5
    baseUrl = '/api'
    // #endif
    
    // #ifndef H5
    baseUrl = 'http://localhost:8888' // 如果本地手机调试，记得改成本机内网IP
    // #endif
} else {
    // 生产环境：直接指向你的公网后端域名
    // 注意：由于你之前配置的 Nginx 已经处理了转发，这里不需要加端口号
    baseUrl = 'https://clothapi.235235.vip'
}

const inflightGetRequests = new Map()
const ENABLE_CLIENT_I18N_FALLBACK = uni.getStorageSync('client-i18n-fallback') === '1'
const i18nNormalizeEndpoints = [
    '/popup/getActivePopups',
    '/phoneAreaCode/getEnabledPhoneAreaCodes',
    '/kefu/',
    '/englishLearning/content/getCategoryList',
    '/englishLearning/content/getChapterList',
    '/englishLearning/content/getVideoCategoryList',
    '/englishLearning/content/getVideoSeriesList',
    '/englishLearning/content/findVideoSeries',
    '/englishLearning/content/findVideoEpisode',
    '/englishLearning/content/getVideoEpisodeList',
    '/englishLearning/word/getWordList',
    '/englishLearning/word/findWord',
    '/englishLearning/word/getErrorLogList',
    '/englishLearning/video/getSentenceList',
    '/englishLearning/userData/getWatchHistoryList',
    '/englishLearning/userData/getCollectionList',
    '/englishLearning/userData/getCollectionDetailList',
    '/englishLearning/checkin/getCheckinRecordList',
    '/englishLearning/checkin/getPointRecordList',
    '/englishLearning/asset/getFreeTimeRecordList'
]
const i18nFieldKeys = new Set([
    'name', 'title', 'content', 'countryName',
    'seriesName', 'episodeName', 'description',
    'explanation', 'translate', 'announcementContent'
])
const i18nLangKeyReg = /^[a-z]{2}(?:-[A-Za-z]{2})?$/

const isObjectLike = (value) => !!value && typeof value === 'object' && !Array.isArray(value)

const isLikelyI18nMapObject = (value) => {
    if (!isObjectLike(value)) {
        return false
    }
    const keys = Object.keys(value)
    if (keys.length === 0 || keys.length > 16) {
        return false
    }
    let langKeyCount = 0
    for (const key of keys) {
        if (i18nLangKeyReg.test(key)) {
            langKeyCount += 1
            continue
        }
        return false
    }
    return langKeyCount > 0
}

const normalizeI18nFieldValue = (value, lang) => {
    if (value === null || value === undefined) {
        return ''
    }
    if (isLikelyI18nMapObject(value)) {
        return localText(value, lang)
    }
    if (typeof value === 'string' && value.charAt(0) === '{') {
        try {
            const parsed = JSON.parse(value)
            if (isLikelyI18nMapObject(parsed)) {
                return localText(parsed, lang)
            }
        } catch (e) {
            return value
        }
    }
    return value
}

const normalizeI18nPayload = (payload, lang, parentKey = '') => {
    if (Array.isArray(payload)) {
        return payload.map((item) => normalizeI18nPayload(item, lang, parentKey))
    }
    if (!isObjectLike(payload)) {
        return payload
    }

    if (isLikelyI18nMapObject(payload)) {
        return localText(payload, lang)
    }

    const out = {}
    Object.keys(payload).forEach((key) => {
        const rawValue = payload[key]
        if (i18nFieldKeys.has(key)) {
            out[key] = normalizeI18nFieldValue(rawValue, lang)
            return
        }
        out[key] = normalizeI18nPayload(rawValue, lang, key)
    })
    return out
}

const shouldNormalizeI18nByUrl = (rawUrl) => {
    const purePath = String(rawUrl || '').split('?')[0]
    return i18nNormalizeEndpoints.some((item) => {
        if (item.endsWith('/')) {
            return purePath.startsWith(item)
        }
        return purePath === item
    })
}

const isUniEncryptEnabled = () => uni.getStorageSync('learning-api-encrypt-enabled') !== '0'

const deriveEncryptKey = (token) => {
    const base = `${String(token || '').trim()}|uni-api-v1`
    return CryptoJS.SHA256(base)
}

const deriveSignKey = (token) => CryptoJS.SHA256(`${String(token || '').trim()}|uni-api-sign-v1`).toString()

const buildReqSign = ({ method, path, rawQuery, ts, nonce, token }) => {
    const material = [
        String(method || '').toUpperCase().trim(),
        String(path || '').trim(),
        String(rawQuery || '').trim(),
        String(ts || '').trim(),
        String(nonce || '').trim(),
    ].join('|')
    return CryptoJS.HmacSHA256(material, CryptoJS.enc.Hex.parse(deriveSignKey(token))).toString(CryptoJS.enc.Hex)
}

const buildSignHeaders = ({ finalUrl, method, token }) => {
    const nowTs = String(Math.floor(Date.now() / 1000))
    const nonce = `${Date.now()}-${Math.random().toString(36).slice(2, 10)}`
    const rawPath = String(finalUrl || '').replace(String(baseUrl || ''), '')
    const parts = rawPath.split('?')
    const path = parts[0] || ''
    const rawQuery = parts[1] || ''
    const sign = buildReqSign({ method, path, rawQuery, ts: nowTs, nonce, token })
    return {
        'X-Req-Ts': nowTs,
        'X-Req-Nonce': nonce,
        'X-Req-Sign': sign,
    }
}

const tryDecryptPayloadData = (rawData, token) => {
    if (!rawData || typeof rawData !== 'object' || Array.isArray(rawData)) {
        return rawData
    }
    const payload = String(rawData.payload || '')
    const iv = String(rawData.iv || '')
    if (!payload || !iv) {
        return rawData
    }
    try {
        const key = deriveEncryptKey(token)
        const ivWord = CryptoJS.enc.Base64.parse(iv)
        const cipherText = CryptoJS.enc.Base64.parse(payload)
        const decrypted = CryptoJS.AES.decrypt(
            { ciphertext: cipherText },
            key,
            {
                iv: ivWord,
                mode: CryptoJS.mode.CBC,
                padding: CryptoJS.pad.Pkcs7,
            }
        )
        const utf8 = decrypted.toString(CryptoJS.enc.Utf8)
        if (!utf8) {
            return rawData
        }
        return JSON.parse(utf8)
    } catch (e) {
        return rawData
    }
}

export const request = ({url, data, header, method, params}) => {
    // 处理 params 参数拼接到 url
    let finalUrl = baseUrl + url;
    if (params && Object.keys(params).length > 0) {
        const queryString = Object.keys(params)
            .map(key => `${encodeURIComponent(key)}=${encodeURIComponent(params[key])}`)
            .join('&');
        finalUrl += (url.includes('?') ? '&' : '?') + queryString;
    }

    const requestMethod = String(method || 'get').toUpperCase()
    const currentLang = uni.getStorageSync('app-lang') || 'zh'
    const dedupeKey = `${requestMethod}::${finalUrl}::${currentLang}`
    if (requestMethod === 'GET' && inflightGetRequests.has(dedupeKey)) {
        return inflightGetRequests.get(dedupeKey)
    }

    const requestPromise = new Promise((resolve, reject) => {
        const token = uni.getStorageSync('x-token')
        const signHeaders = buildSignHeaders({ finalUrl, method: requestMethod, token })
        uni.request({
            url: finalUrl, // 使用拼接后的 URL
            data: data || '',
            method,
            header: {
                'x-token': token,
                'Accept-Language': uni.getStorageSync('app-lang') || 'zh',
                'X-Client-Platform': 'uni',
                'X-Resp-Encrypt': isUniEncryptEnabled() ? '1' : '0',
                ...signHeaders,
                ...header
            },
            success: (res) => {
				if (res.header['new-token']) {
                    uni.setStorageSync('x-token',res.header['new-token'])
                }
                if (res.statusCode === 401) {
                    const redirected = redirectToLogin()
                    if (!redirected) {
                        resolve(normalizeApiResponse(res.data))
                    }
                    return
                }
                // 封禁检测：后端返回 banned=true 时强制退出登录并跳转登录页
                if (res.data && res.data.data && res.data.data.banned) {
                    redirectToLogin()
                    uni.showModal({
                        title: '',
                        content: normalizeBackendMessage(res.data.msg) || t('operationFailed'),
                        showCancel: false,
                        confirmText: t('confirm'),
                        success: () => {
                            const redirected = redirectToLogin()
                            if (!redirected) {
                                uni.reLaunch({ url: buildLoginRedirectUrl() })
                            }
                        }
                    })
                    resolve(res.data)
                    return
                }
                // 维护模式：503 + maintenance=true → 缓存配置并跳转维护页（不弹toast）
                if (res.statusCode === 503 && res.data && res.data.data && res.data.data.maintenance) {
                    uni.setStorageSync('maintenance_config', JSON.stringify(res.data.data))
                    const pages = getCurrentPages()
                    const currentRoute = pages.length > 0 ? '/' + pages[pages.length - 1].route : ''
                    if (currentRoute !== '/pages/maintenance/index') {
                        uni.redirectTo({ url: '/pages/maintenance/index' })
                    }
                    resolve(res.data)
                    return
                }
                const payload = normalizeApiResponse(res.data)
                if (payload.code === 0 && payload.data && typeof payload.data === 'object' && payload.data.alg === 'aes-cbc') {
                    payload.data = tryDecryptPayloadData(payload.data, token)
                }
                if (!isApiResponseObject(res.data)) {
                    uni.showToast({
                        title: payload.msg,
                        icon: 'none'
                    })
                    resolve(payload)
                    return
                }
                // 通用错误提示：排除上面已处理的特殊状态
                if(payload.code != 0){
                    const msg = normalizeBackendMessage(payload.msg)
                    if (msg) {
						uni.showToast({
							title: msg,
							icon: 'none'
						});
                    }
				}
                const shouldNormalizeI18n = ENABLE_CLIENT_I18N_FALLBACK && payload.code === 0 && shouldNormalizeI18nByUrl(url)
                if (shouldNormalizeI18n) {
                    resolve({
                        ...payload,
                        data: normalizeI18nPayload(payload.data, currentLang)
                    })
                    return
                }
                resolve(payload)

            },
            fail: (err) => {
                reject(err)
            },
            timeout: 30000
        });
    })

    if (requestMethod === 'GET') {
        inflightGetRequests.set(dedupeKey, requestPromise)
        requestPromise.finally(() => {
            inflightGetRequests.delete(dedupeKey)
        })
    }

    return requestPromise
}
