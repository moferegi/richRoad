import {getMinProLoginInfo} from '@/api/loginApi.js'
/**
 * 微信登陆方法
 * @param {Function} callback  登陆信息回调
 * 页面中引入次js 导出 login
 * import {login} from '@/common/login.js'
 * 在需要登陆的地方调用
 * token 登陆鉴权
 * userInfo 用户信息
 * login((token,userInfo) => {
 *  console.log(token,userInfo)
 *  })
 */
export const login = (callback) => {
	let token = uni.getStorageSync('token')
	if (token) {
		callback && callback(token, uni.getStorageSync('userInfo'))
		return
	} else {
	uni.login({
		provider: 'weixin',
		success: (res) => {
			callback(res)
			// //接口根据项目情况而定 callback返回 token和用户信息
			// getMinProLoginInfo({code:res.code}).then((res) => {
			// 	console.log(res)
			// 	uni.setStorageSync('userInfo',res.data)
			// 	uni.setStorageSync('token',res.data.token || '')
			// 	if(isBinDing && !res.data.bindStatus && !res.data.token) {
			// 		uni.navigateTo({
			// 			url:"/pages/binding/binding"
			// 		})
			// 	} else {
			// 		callback && callback(res.data.token,res.data)
			// 	}
			// })
		},
		fail: () => {
			uni.showToast({
				icon: "none",
				title: '登录失败'
			})
		}
	})
	}
}