import { myRouter }from  '@/utils/permission.js'
 export const baseUrl = 'http://localhost:8888'
export const request = ({url, data, header, method}) => {
    return new Promise((resolve, reject) => {
        uni.request({
            url: baseUrl + url, //仅为示例，并非真实接口地址。
            data: data || '',
            method,
            header: {
                'x-token': uni.getStorageSync('x-token') ,//自定义请求头信息
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
                    myRouter("/pages/tabBar/index")
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