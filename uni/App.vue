<script>
	import {useUserStore} from "@/pinia/modules/user.js"
	import { getOpenID } from '@/api/base.js'
	import {myRouter} from '@/utils/permission.js'
	export default {
		onLaunch: function() {
			const userStore = useUserStore()
			userStore.getInfo()
			wx.login({
			  success: async (res) => {
			    if (res.code) {
			      //发起网络请求
			      const data = await getOpenID(res.code)
				  if(data.code === 0) {
					  uni.setStorageSync('openid', data.data.openid)
				  }
			    } else {
			      console.log('登录失败！' + res.errMsg)
			    }
			  }
			})
			myRouter()
		},
		onShow: function() {
			console.log('App Show')
		},
		onHide: function() {
			console.log('App Hide')
		}
	}
</script>

<style lang="scss">
@import "uview-plus/index.scss";
	@import '@/common/common.css';
	/*每个页面公共zcss */
</style>
