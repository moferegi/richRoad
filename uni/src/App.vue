<script>
	import {useUserStore} from "@/pinia/modules/user.js"
	import { useLangStore } from "@/pinia/modules/lang.js"
	import { useAppConfigStore } from "@/pinia/modules/appConfig.js"
	import { getOpenID } from '@/api/base.js'
	import {myRouter} from '@/utils/permission.js'
	import { visitorHeartbeat } from '@/api/visitor.js'
	import { generateFingerprint, getSessionId, getPlatform } from '@/utils/fingerprint.js'
	import { initExternalDomain, initCdnDomain } from '@/utils/url.js'
	export default {
		onLaunch: function() {
			const userStore = useUserStore()
			const langStore = useLangStore()
			const appConfigStore = useAppConfigStore()
			userStore.getInfo()
			langStore.initLangs()
			initExternalDomain()
			initCdnDomain()
			appConfigStore.loadConfig()
      console.log(123)
			wx.login({
			  success: async (res) => {
			    if (res.code) {
			      //发起网络请求
			      const data = await getOpenID(res.code)
				  if(data.code === 0) {
            console.log( data.data.openid)
					  uni.setStorageSync('openid', data.data.openid)
				  }
            console.log( data.data)
			    } else {
			      console.log('登录失败！' + res.errMsg)
			    }
			  }
			})
			myRouter()
			// 初始化访客指纹
			generateFingerprint()
		},
		onShow: function() {
			console.log('App Show')
			const langStore = useLangStore()
			langStore.updateTabBar(langStore.locale || uni.getStorageSync('app-lang') || 'mn')
			this.reportVisitor()
		},
		onHide: function() {
			console.log('App Hide')
		},
		methods: {
			reportVisitor() {
				try {
					const info = uni.getSystemInfoSync()
					visitorHeartbeat({
						visitorId: generateFingerprint(),
						sessionId: getSessionId(),
						platform: getPlatform(),
						pagePath: getCurrentPages().length > 0 ? '/' + getCurrentPages()[getCurrentPages().length - 1].route : '/',
						referer: '',
						screenWidth: info.screenWidth,
						screenHeight: info.screenHeight,
						language: uni.getStorageSync('app-lang') || 'zh'
					}).catch(() => {})
				} catch (e) {
					console.log('visitor report error', e)
				}
			}
		}
	}
</script>

<style lang="scss">
@import '@/common/common.css';
	/*每个页面公共zcss */
</style>
