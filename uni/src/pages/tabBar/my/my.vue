<template>
	<view class="my-page" style="padding-bottom: 60rpx;">
		<view class="my_head_box boxs_bb pos_r">
			<view class="flex flex-aic info-box" v-if="isShow">
        <button class="avatar-wrapper" open-type="chooseAvatar" @chooseavatar="onChooseAvatar">
          <image class="avatar"
                 :src="info.avatar || avatarUrl">
          </image>
        </button>
        <view class="nick-name">
          <view>{{ $t('nickname') }}</view>
          <input name="nickName" type="nickname" :placeholder="$t('nicknamePlaceholder')" class="weui-input" @change="onInputNickName"
                 v-model="info.nickname" />
        </view>
			</view>
      <view v-if="!isShow" class="flex flex-aic info-box unLogin">
        <p @tap="logins">{{ $t('loginPrompt') }}</p>
      </view>
		</view>
		<view class="my_order_card boxs_bb">
			<view class="flex-aic flexr-jsb my_order_card_title">
				<text class="color_333 font_28">{{ $t('myOrders') }}</text>
				<view class="flex flex-aic" @tap="toOrder('all')">
					<text class="color_999 font_28 m_r_8">{{ $t('viewAll') }}</text>
					<uni-icons color="#999999" type="forward" size="14"></uni-icons>
				</view>
			</view>
      <view class="flex-aic flexr-jsa color_333 font_24">
        <view class="tac" @tap="toNav(`/pages/order/order?status=${0}`)">
          <image class="my_order_icon m_b_8" src="../../../static/images/Wallet.png" mode=""></image>
          <view>{{ $t('ordersPending') }}</view>
        </view>
        <view class="tac" @tap="toNav(`/pages/order/order?status=${1}`)">
          <image class="my_order_icon m_b_8" src="../../../static/images/delivery.png" mode="">
          </image>
          <view>{{ $t('ordersShipping') }}</view>
        </view>
        <view class="tac" @tap="toNav(`/pages/order/order?status=${2}`)">
          <image class="my_order_icon m_b_8" src="../../../static/images/package.png"
                 mode=""></image>
          <view>{{ $t('ordersReceiving') }}</view>
        </view>
        <view class="tac" @tap="toNav(`/pages/order/order?status=${3}`)">
          <image class="my_order_icon m_b_8" src="../../../static/images/comment.png"
                 mode=""></image>
          <view>{{ $t('ordersToReview') }}</view>
        </view>

        <view class="tac" @tap="toNav(`/pages/order/order?status=${7}`)">
          <image class="my_order_icon m_b_8" src="../../../static/images/commented.png"
                 mode=""></image>
          <view>{{ $t('ordersReviewed') }}</view>
        </view>

      </view>
		</view>
		<view class="my_tools_box boxs_bb" v-if="columns.length >=1">
			<template v-for="(item,index) in columns.filter(c=>!c.hidden)">
				<view @tap="toPages(item.pages)" class="my_tools_item flex-aic flexr-jsb"
					:class="columns.length-1 == index ? '' : 'b_b_2'">
					<image class="m_r_16 my_tools_img" :src="item.icon" mode=""></image>
					<view class="flex-fitem color_333 font_28">{{$t(item.titleKey)}}</view>
					<uni-icons color="#ccc" type="forward" size="32"></uni-icons>
				</view>
			</template>
		</view>
	</view>
</template>


<script setup>
import { myRouter } from "@/utils/permission";
import { setClientUserInfo } from "@/api/base";
import {useUserStore} from "@/pinia/modules/user.js"
import { onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { ref, computed } from 'vue'
const langStore = useLangStore()
const $t = computed(() => langStore.$t)
const isShow = ref(false)
	const columns = ref([
		{
			titleKey: 'addressManage',
			pages: '/pages/address/address',
			icon: '/static/images/address.png'
		},{
			titleKey: 'myCollectionMenu',
			pages: '/pages/collect/collect',
			icon: '/static/images/collectionIcon.png'
		},{
			titleKey: 'logout',
			pages: 'exit',
			icon: '/static/images/exit.png',
			hidden: true
		}
	])

const defaultAvatarUrl = 'https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0'
const avatarUrl = ref('')
const nickName = ref('默认用户')
avatarUrl.value = defaultAvatarUrl

const info = ref({})
const userStore = useUserStore()
onShow(() => {
const token = userStore.token || ''
  if (token) {
    isShow.value = true
    info.value = uni.getStorageSync("userInfo")
  } else {
    isShow.value = false
  }
  const c = columns.value.find(item=>item.pages==='exit')
  if(c){
	  c.hidden = !isShow.value
  }
})

const logins = () => {
  uni.redirectTo({
    url: '/pages/user/login'
  })
}

// 修改头像的方法
const onChooseAvatar = async (e) => {
  avatarUrl.value = e.detail.avatarUrl
  let tmpFilePath = avatarUrl.value
  // 对微信返回的临时图片链接进行base64编码
  var avatarUrl_base64 = 'data:image/jpeg;base64,' + wx.getFileSystemManager().readFileSync(tmpFilePath, 'base64');
  const res = await setClientUserInfo({
    key: "avatar",
    value: avatarUrl_base64
  })
  if(res.code === 0) {
    uni.showToast({
      title: res.msg,
      icon: 'success'
    })
  } else {
    uni.showToast({
      title: res.msg,
      icon: 'none'
    })
  }
}

// 修改昵称的方法
const onInputNickName = async (e) => {
  nickName.value = e.detail.value
  const res = await setClientUserInfo({
    key: "nickname",
    value: nickName.value
  })
  if (res.code === 0) {
    uni.showToast({
      title: res.msg,
      icon: 'success'
    })
  } else {
    uni.showToast({
      title: res.msg,
      icon: 'none'
    })
  }
}


	const toPages = (pages) => {
		if(pages === 'exit'){
			userStore.loginOut()
			uni.showToast({
				icon: 'none',
				title: $t('logoutSuccess')
			})
			
			return
		}
		if (!pages) {
			uni.showToast({
				icon: 'none',
				title: $t('developing')
			})
		} else {
			uni.navigateTo({
				url: pages
			})
		}
	}

	const toOrder = (val) => {
    myRouter("/pages/order/order",true)
	}
	const toNav = (pages) => {
		uni.navigateTo({
			url: pages
		})
	}
	const toRetreatOrder = () => {
		uni.navigateTo({
			url: '/pages/retreat/retreatOrder'
		})
	}
</script>

<style lang="scss">
	page {
		background-color: #f4f7fb;
	}

	.my-page {
		background: radial-gradient(120% 80% at 100% -10%, #dbeafe 0%, transparent 60%), #f4f7fb;
	}

	.my_head_box {
		height: 284rpx;
		padding: 24rpx 32rpx 0;
		margin-bottom: 20rpx;
		background: linear-gradient(135deg, #2563eb, #0ea5e9);
		border-bottom-left-radius: 28rpx;
		border-bottom-right-radius: 28rpx;
	}

	.my_head_img {
		width: 96rpx;
		height: 96rpx;
		border-radius: 100%;
	}

	.my_integral_view {
		width: 686rpx;
		height: 172rpx;
		box-shadow: 0rpx 4rpx 8rpx 0rpx rgba(238, 238, 238, 1);
		border-radius: 20rpx;
		bottom: -104rpx;
		left: 50%;
		transform: translateX(-50%);
		right: 0;
	}

	.my_binding_box {
		height: 26rpx;
		line-height: 26rpx;
		padding: 0 16rpx;
		background-color: rgba(0, 0, 0, 0.1);
		border-radius: 12rpx;
	}

	.my_integral_line {
		width: 4rpx;
		height: 76rpx;
	}

	.my_order_card {
		margin: 0 auto 24rpx;
		width: 686rpx;
		height: 226rpx;
		box-shadow: 0rpx 14rpx 30rpx 0rpx rgba(15, 23, 42, 0.08);
		border-radius: 16rpx;
		background: #fff;
	}

	.my_order_card_title {
		padding: 20rpx 32rpx 32rpx;
	}

	.my_order_icon {
		height: 60rpx;
		width: 60rpx;
	}

	.my_vip_box {
		padding-left: 40rpx;
		height: 136rpx;
		width: 100%;
		background-image: url(http://www.liwanying.top/applate-icon/vip-bgc.png);
		background-repeat: no-repeat;
		background-size: 100%;
	}

	.my_tools_box {
		width: 686rpx;
		margin: 0 auto;
		padding: 0 32rpx;
		border-radius: 16rpx;
		overflow: hidden;
		background: #fff;
		box-shadow: 0rpx 14rpx 30rpx 0rpx rgba(15, 23, 42, 0.08);
	}

	.my_tools_item {
		padding: 24rpx 0;
	}

	.my_tools_img {
		width: 56rpx;
		height: 56rpx;
	}

  .info-box{
    padding: 20rpx;
		background: rgba(255,255,255,0.22);
		backdrop-filter: blur(6px);
    border-radius: 20rpx;
    .nick-name{
      color: #fff;
    }
    .avatar-wrapper{
      padding: 0;
	  margin: 0;
	  margin-right: 20rpx;
      .avatar{
        width: 120rpx;
        height: 120rpx;
      }
    }
  }

.unLogin{
  height: 200rpx;
  text-align: center;
  color: #fff;
  display: flex;
  justify-content: center;
	align-items: center;
	font-weight: 600;
}
</style>
