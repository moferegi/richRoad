import { defineStore } from 'pinia'
import {ref, watch} from 'vue'
import { login, getUserInfo, phoneLogin } from '@/api/base.js'
import { myRouter } from '../../utils/permission'

export const useUserStore = defineStore('user', () => {
    const userInfo = ref({
        nickname: '',
        username: '',
		ID: ''
    })
    const token = ref(uni.getStorageSync('x-token') || '')

    const setToken = (val) => {
        token.value = val
    }
    const setUserInfo = (val) => {
        userInfo.value = val
		uni.setStorageSync('ID',val.ID)
    }


    // 登录
    const loginIn = async (loginInfo) => {
        const res = await login(loginInfo)
        if(res.code === 0) {
           await setToken(res.data.token)
			uni.switchTab({
				url: '/pages/tabBar/index'
			})
            return true
        }
        uni.showToast({	// 提示错误信息
            icon:'none',
            title: res.msg,
            duration: 3000
        })
        return false
    }

    // 手机号登录
    const phoneLoginIn = async (loginInfo) => {
        const res = await phoneLogin(loginInfo)
        if(res.code === 0) {
           await setToken(res.data.token)
			uni.switchTab({
				url: '/pages/tabBar/index'
			})
            return true
        }
        uni.showToast({
            icon:'none',
            title: res.msg,
            duration: 3000
        })
        return false
    }

    const loginOut = async() =>{
        uni.removeStorageSync("x-token")
        uni.removeStorageSync("ID")
        clearUserInfo()
		uni.switchTab({
			url: '/pages/tabBar/index'
		})
		// myRouter('/pages/tabBar/index',true)
    }

	const getInfo = async ()=> {
		const token = uni.getStorageSync('x-token')
		if(token){
			const infoRes = await getUserInfo()
			if(infoRes.code === 0){
				await setUserInfo(infoRes.data)
                uni.setStorageSync('userInfo',infoRes.data)
			}
		}
	}

    // 清空用户信息
    const clearUserInfo = () => {
        setUserInfo(
            {
                nickname: '',
                username: '',
				ID:''
            }
        )
        uni.removeStorageSync("userInfo")
        setToken("")
    }

    watch(()=>token.value,()=>{
        uni.setStorageSync('x-token',token.value)
		if(token.value){
			getInfo()
		}else{
			clearUserInfo()
		}
    })

    return {
        userInfo,
        token,
        loginIn,
        phoneLoginIn,
        setToken,
        setUserInfo,
        getInfo,
        clearUserInfo,
        loginOut
    }
})
