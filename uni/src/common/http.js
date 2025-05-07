/**
 * ===========
 * 统一请求发送
 * ===========
 * token 和 请求拦截需要根据自己业务修改
 */
import config from './config.js'
import { login } from './login.js';
export default function http(opts, data = {}) {
	let token = uni.getStorageSync('token') || '';
	let requestUrl = ''
	for (let key in data) {
		if (data[key] === '' || data[key] === null) {
			delete data[key]
		}
	}
	if (opts.url.substr(0, 7).toLowerCase() == "http://" || opts.url.substr(0, 8).toLowerCase() == "https://") {
		requestUrl = opts.url
	} else {
		requestUrl = config.baseUrl + opts.url
	}
	return new Promise((resolve, reject) => {
		uni.request({
			url: requestUrl,
			data: data,
			method: opts.method,
			header: {
				"Authorization":token,
				...opts.header
			},
			dataType: 'json',
			success: ({
				data: res
			}) => {
				console.log('uni.request--1--',{
					url: requestUrl,
					data: data,
					method: opts.method,
				})
				console.log('uni.request--2--',res)
				if(res.code ==401) {
					reject(res)
					uni.removeStorageSync("userInfo") // 清除本地用户信息
					uni.removeStorageSync("token")
					login(true,(token,userInfo) =>{
						
					})
					return;
				}
				if (res.code !== 200) {
					reject(res)
					uni.showToast({
						title:res.msg || '服务器错误',
						icon: 'none'
					})
					return;
				}
				// 正常
				resolve(res)
			},
			fail: err => {
				console.log('error-----',err)
				reject(err)
				uni.showToast({
					title: '网络连接异常',
					icon: 'none'
				})
			}
		})
	})
}
