import { myRouter }from  '@/utils/permission.js'
import { t, localText } from '@/utils/i18n.js'

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
    baseUrl = 'https://back.mnmovie.icu'
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

    return new Promise((resolve, reject) => {
        uni.request({
            url: finalUrl, // 使用拼接后的 URL
            data: data || '',
            method,
            header: {
                'x-token': uni.getStorageSync('x-token'),
                'Accept-Language': uni.getStorageSync('app-lang') || 'zh',
                ...header
            },
            success: (res) => {
				if (res.header['new-token']) {
                    uni.setStorageSync('x-token',res.header['new-token'])
                }
                if (res.statusCode === 401) {
                    uni.removeStorageSync('x-token')
                    uni.removeStorageSync('userInfo')
                    uni.reLaunch({ url: '/pages/user/login' })
                    return
                }
                // 封禁检测：后端返回 banned=true 时强制退出登录并跳转登录页
                if (res.data && res.data.data && res.data.data.banned) {
                    uni.removeStorageSync('x-token')
                    uni.removeStorageSync('userInfo')
                    uni.showModal({
                        title: '',
                        content: res.data.msg || 'Account banned',
                        showCancel: false,
                        confirmText: t('confirm'),
                        success: () => {
                            uni.reLaunch({ url: '/pages/user/login' })
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
                // 通用错误提示：排除上面已处理的特殊状态
                if(res.data.code != 0){
					uni.showToast({
						title: res.data.msg,
						icon: 'none'
					});
				}
                resolve(res.data)

            },
            fail: (err) => {
                reject(err)
            },
            timeout: 30000
        });
    })
}
