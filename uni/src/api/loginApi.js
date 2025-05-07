import http from '@/common/http.js'
//获取当前登录信息
export const getMinProLoginInfo = (data) => {
	return http({
		url: '/getMinProLoginInfo',
		method: 'POST'
	},data)
}
//发送验证码
export const sendVerCode = (data) => {
	return http({
		url: '/sms/sendVerCode',
		method: 'POST'
	},data)
}
//小程序注册
export const minProRegisterUser = (data) => {
	return http({
		url: '/minProRegisterUser',
		method: 'POST'
	},data)
}
//小程序绑定
export const minProBindUser = (data) => {
	return http({
		url: '/minProBindUser',
		method: 'POST'
	},data)
}