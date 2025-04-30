<template>
	<view>
		<view class="bgc_fff flex logistics_code_box b_b_2">
			<image class="logistics_icon m_r_24" src="http://www.liwanying.top/applate-icon/kuaidicheche.png" mode=""></image>
			<view class="flex-fitem">
				<view>顺丰速运</view>
				<view class="flex-aic flexr-jsb font_28">
					<text class="color_666">{{ expressNum }}</text>
					<view class="color_fe5572" @tap="copyBoard">
						<text class="m_r_16">复制</text>
					</view>
				</view>
			</view>
		</view>
		<view class="logistics_steps_box bgc_fff">
      <uni-steps :options="infoList" direction="column" :active="2"></uni-steps>
		</view>
	</view>
</template>

<script setup>
	import {ref} from 'vue'
  import { checkRouters } from '@/api/order'
  import { onLoad } from '@dcloudio/uni-app'

  const infoList = ref([])
  const express = ref('')
  const expressNum = ref('')
  const init = async () => {
    const res = await checkRouters(express.value)
    if(res.code === 0) {
      // 把返回值的res.data.msgData.routeResps.routes洗数据，把time改成desc, remark改成title
      infoList.value = res.data.msgData.routeResps[0].routes
      infoList.value.forEach(item => {
        item.desc = item.acceptTime
        item.title = item.remark
      })
      expressNum.value = res.data.msgData.routeResps[0].mailNo
    }
  }

  onLoad( async (options) => {
    express.value = options.express
    init()
  })

  const copyBoard = () => {
    uni.setClipboardData({
      data: expressNum.value,
      success: function () {
        uni.showToast({
          title: '复制成功',
          icon: 'none'
        })
      }
    })
  }

</script>

<style lang="scss">
	page {
		background-color: #f8f8f8;
	}
	.logistics_notify_box {
		height: 68rpx;
	}
	.logistics_icon {
		width: 72rpx;
		height: 72rpx;
	}
	.logistics_code_box {
		padding: 24rpx 32rpx;
	}
	.logistics_steps_box {
		padding: 24rpx 0 100rpx 56rpx;
	}
	.logistics_steps_icon {
		width: 24rpx;
		height: 24rpx;
		background-color: rgba(255, 112, 0, 0.3);
		border-radius: 50%;
		view {
			width: 12rpx;
			height: 12rpx;
			border-radius: 50%;
		}
	}
	.logistics_steps_icon_out {
		background-color: rgba(194, 199, 204, 0.3);
		view {
			background-color: #C2C7CC;
		}
	}
</style>
