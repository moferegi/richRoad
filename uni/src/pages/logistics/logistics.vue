<template>
	<view class="nf-logistics-page">
		<view class="nf-logistics-header">
			<image class="nf-logistics-icon" src="http://www.liwanying.top/applate-icon/kuaidicheche.png" mode=""></image>
			<view class="nf-logistics-info">
				<view class="nf-logistics-carrier">{{ $t('logisticsCarrier') }}</view>
				<view class="nf-logistics-num-row">
					<text class="nf-logistics-num">{{ expressNum }}</text>
					<view class="nf-logistics-copy" @tap="copyBoard">
						<text>{{ $t('copy') }}</text>
					</view>
				</view>
			</view>
		</view>
		<view class="nf-logistics-steps">
      <uni-steps :options="infoList" direction="column" :active="2"></uni-steps>
		</view>
	</view>
</template>

<script setup>
	import {ref, computed} from 'vue'
  import { checkRouters } from '@/api/order'
  import { onLoad } from '@dcloudio/uni-app'
  import { useLangStore } from '@/pinia/modules/lang.js'

  const langStore = useLangStore()
  const $t = computed(() => langStore.$t)

  const infoList = ref([])
  const express = ref('')
  const expressNum = ref('')
  const init = async () => {
    const res = await checkRouters(express.value)
    if(res.code === 0) {
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
          title: $t.value('copySuccess'),
          icon: 'none'
        })
      }
    })
  }

</script>

<style lang="scss" scoped>
	.nf-logistics-page {
		min-height: 100vh;
		background: #000;
	}
	.nf-logistics-header {
		display: flex;
		align-items: center;
		padding: 30rpx 32rpx;
		background: rgba(255, 255, 255, 0.04);
		border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
	}
	.nf-logistics-icon {
		width: 72rpx;
		height: 72rpx;
		margin-right: 24rpx;
	}
	.nf-logistics-info {
		flex: 1;
	}
	.nf-logistics-carrier {
		color: rgba(255, 255, 255, 0.9);
		font-size: 30rpx;
		font-weight: 600;
		margin-bottom: 8rpx;
	}
	.nf-logistics-num-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}
	.nf-logistics-num {
		color: rgba(255, 255, 255, 0.5);
		font-size: 28rpx;
	}
	.nf-logistics-copy {
		color: #e50914;
		font-size: 28rpx;
		padding: 4rpx 16rpx;
	}
	.nf-logistics-steps {
		padding: 24rpx 0 100rpx 56rpx;
		background: rgba(255, 255, 255, 0.02);
	}
</style>
