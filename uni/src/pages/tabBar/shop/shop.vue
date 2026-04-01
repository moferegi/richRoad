<template>
  <view class="nf-shop">
    <view class="nf-shop-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff"></uni-icons>
        </view>
        <text class="nf-navbar-title">{{ $t('tabCart') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <view class="nf-shop-body">
      <shop-goods-list></shop-goods-list>

      <!-- 推荐商品 -->
      <view class="nf-recommend" v-if="flowData && flowData.length">
        <view class="nf-recommend-header">
          <view class="nf-recommend-line"></view>
          <text class="nf-recommend-title">{{ $t('guessYouLike') }}</text>
          <view class="nf-recommend-line"></view>
        </view>
        <noPaginGridGoodList :goodsList="flowData"></noPaginGridGoodList>
      </view>
      <view style="height: 180rpx;"></view>
    </view>
  </view>
</template>

<script setup>
import noPaginGridGoodList from '@/components/good-list/no-pagin-grid-good-list.vue'
import shopGoodsList from './components/shop-goods-list.vue'
import { ref, computed } from "vue"
import { getGoodList } from '@/api/homePage.js'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const flowData = ref([])
const params = ref({ page: 1, pageSize: 10, categoryID: 0 })

const init = async () => {
  const res = await getGoodList(params.value)
  if (res.code === 0 && res.data.list.length) {
    flowData.value.push(...res.data.list)
  }
}
init()

const goBack = () => {
  uni.switchTab({ url: '/pages/tabBar/index' })
}
</script>

<style lang="scss">
page {
  background-color: #000;
  min-height: 100vh;
}

.nf-shop {
  min-height: 100vh;
  background: #000;
  position: relative;
}

.nf-shop-bg {
  position: fixed;
  top: 0; left: 0; right: 0;
  height: 500rpx;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(ellipse at 30% 0%, rgba(229, 9, 20, 0.10) 0%, transparent 60%),
    radial-gradient(ellipse at 70% 10%, rgba(229, 9, 20, 0.06) 0%, transparent 50%);
}

.nf-navbar {
  background: rgba(0, 0, 0, 0.85);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
  padding: 0 28rpx 16rpx;
  position: sticky;
  top: 0;
  z-index: 99;
}

.nf-navbar-status {
  width: 100%;
  height: var(--status-bar-height, 0px);
}

.nf-navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 88rpx;
}

.nf-navbar-back {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.06);
  border: 1rpx solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.nf-navbar-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #fff;
  letter-spacing: 2rpx;
}

.nf-shop-body {
  position: relative;
  z-index: 1;
  padding: 20rpx 0 0;
}

.nf-recommend {
  padding: 20rpx 24rpx;
}

.nf-recommend-header {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20rpx;
  padding: 30rpx 0 20rpx;
}

.nf-recommend-line {
  width: 80rpx;
  height: 2rpx;
  background: linear-gradient(90deg, transparent, rgba(229, 9, 20, 0.5), transparent);
}

.nf-recommend-title {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.6);
  letter-spacing: 4rpx;
  font-weight: 500;
}
</style>
