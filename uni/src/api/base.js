import {request} from '@/utils/request.js'

export const getCaptcha  = () => {
    return request({
		url:'/base/captcha',
        method: 'post'
    })
}

export const register  = (data) => {
    return request({
		url:'/clientUser/register',
        method: 'post',
        data:data
    })
}

export const login  = (data) => {
    return request({
		url:'/clientUser/login',
        method: 'post',
        data:data
    })
}


export const getUserInfo  = () => {
    return request({
		url:'/clientUser/getUserInfo',
        method: 'get'
    })
}

export const getOpenID  = (params) => {
    return request({
		url: `/clientUser/getOpenID?code=${params}`,
        method: 'get'
    })
}

// 获取唤醒微信支付的参数
export const getPayParams  = (data) => {
    return request({
		url:'/wxpay/getPayParams',
        method: 'post',
        data:data
    })
}

// 获取支付结果
export const getOrderById  = (params) => {
    return request({
		url:`/wxpay/getOrderById?orderID=${params}`,
        method: 'get'
    })
}

// 修改用户信息
export const setClientUserInfo  = (data) => {
    return request({
		url:`/clientUser/setClientUserInfo`,
        method: 'post',
		data,
    })
}


export const checkNeedPay  = (data) => {
    return request({
		url:'/wxpay/checkNeedPay',
        method: 'post',
        data:data
    })
}

// 获取我的邀请信息
export const getMyInviteInfo = () => {
    return request({
        url: '/clientUser/getMyInviteInfo',
        method: 'get'
    })
}

// 获取我的下级列表
export const getMySubordinates = (params) => {
    return request({
        url: '/clientUser/getMySubordinates',
        method: 'get',
        params
    })
}

// 手机号登录
export const phoneLogin = (data) => {
    return request({
        url: '/clientUser/phoneLogin',
        method: 'post',
        data
    })
}

// 手机号注册
export const phoneRegister = (data) => {
    return request({
        url: '/clientUser/phoneRegister',
        method: 'post',
        data
    })
}

// 修改密码
export const changePassword = (data) => {
    return request({
        url: '/clientUser/changePassword',
        method: 'post',
        data
    })
}

// 验证后设置手机号
export const setPhoneVerified = (data) => {
    return request({
        url: '/clientUser/setPhoneVerified',
        method: 'post',
        data
    })
}