import { useUserStore } from "@/pinia/modules/user";

export const myRouter = (url,isRedirectTo)=>{
    const userStore = useUserStore()
    const token = userStore.token || ''
	const blackList = ['/pages/orderInfo/orderInfo', '/pages/tabBar/shop/shop', '/pages/order/order']

	const loginList = ['/pages/user/login','/pages/user/register']
    const inBlack = blackList.some(b=>{
        return url&&url.indexOf(b) > -1
    })

	const inLogin = loginList.some(b=>{
        return url&&url.indexOf(b) > -1
    })

	if(inLogin && token) {
		uni.switchTab({
			url: '/pages/tabBar/index'
		})
		return
	}
    if(inBlack && !token) {
        uni.switchTab({
            url: '/pages/user/login'
        })
        return
    }

	if (isRedirectTo){
		uni.navigateTo({
			url: url
		})
		return
	}
	// 对url为undefined时做个非空处理
	if(url){
		uni.switchTab({
			url: url
		})
	}
}
