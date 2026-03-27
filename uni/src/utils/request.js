import { myRouter }from  '@/utils/permission.js'

// 定义并导出 baseUrl 变量
export let baseUrl = '/api'
// 判断是否为 web 环境
// #ifdef H5
baseUrl = '/api'
// #endif

// 非 web 环境（APP、小程序等）使用完整 URL
// #ifndef H5
baseUrl = 'http://localhost:8888'
// #endif

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
                if(res.data.code != 0){
					uni.showToast({
						title: res.data.msg,
						icon: 'none'
					});
				}
				if (res.header['new-token']) {
                    uni.setStorageSync('x-token',res.header['new-token'])
                }
                if (res.statusCode === 401) {
                    uni.removeStorageSync('x-token')
                    myRouter("/pages/home/index")
                    return
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
