<script>
	import {useUserStore} from "@/pinia/modules/user.js"
	import { useLangStore } from "@/pinia/modules/lang.js"
	import { getOpenID } from '@/api/base.js'
	import {myRouter} from '@/utils/permission.js'
	import { visitorHeartbeat } from '@/api/visitor.js'
	import { generateFingerprint, getSessionId, getPlatform } from '@/utils/fingerprint.js'
	export default {
		onLaunch: function() {
			const userStore = useUserStore()
			const langStore = useLangStore()
			userStore.getInfo()
			langStore.initLangs()
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
