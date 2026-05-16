<template>
  <view class="cart-page">
    <view class="cart-navbar">
      <view class="cart-status"></view>
      <view class="cart-navbar-inner">
        <view class="cart-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#0f172a" />
        </view>
        <text class="cart-title">{{ $t('cart') }}</text>
        <view class="cart-placeholder"></view>
      </view>
    </view>

    <view class="cart-body">
      <shop-goods-list />
    </view>
  </view>
</template>

<script setup>
import shopGoodsList from '@/pages/tabBar/shop/components/shop-goods-list.vue'
import { computed } from 'vue'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const goBack = () => {
  const pages = getCurrentPages()
  if (pages.length > 1) {
    uni.navigateBack()
    return
  }
  uni.switchTab({ url: '/pages/tabBar/index' })
}
</script>

<style lang="scss">
page {
  background: #f4f7fb;
}

.cart-page {
  min-height: 100vh;
  background: radial-gradient(120% 80% at 100% -10%, #dbeafe 0%, transparent 60%), #f4f7fb;
}

.cart-navbar {
  position: sticky;
  top: 0;
  left: 0;
  right: 0;
  z-index: 99;
  background: rgba(244, 247, 251, 0.86);
  backdrop-filter: blur(16px);
  border-bottom: 1rpx solid rgba(15, 23, 42, 0.08);
}

.cart-status {
  height: var(--status-bar-height, 0px);
}

.cart-navbar-inner {
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24rpx;
}

.cart-back {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.84);
  border: 1rpx solid rgba(15, 23, 42, 0.1);
  box-shadow: 0 8rpx 18rpx rgba(15, 23, 42, 0.08);
}

.cart-title {
  font-size: 32rpx;
  color: #0f172a;
  font-weight: 700;
}

.cart-placeholder {
  width: 64rpx;
  height: 64rpx;
}

.cart-body {
  padding-top: 16rpx;
  padding-bottom: calc(env(safe-area-inset-bottom) + 132rpx);
  background: transparent;
}
</style>
