<template>
  <view class="flash-sale">
    <view class="flash-header">
      <text class="flash-title">限时秒杀</text>
      <text class="flash-countdown">8点场 07:13:55</text>
      <uni-icons type="right" size="14" color="#999"></uni-icons>
    </view>

    <view class="product-list">
      <view class="product-item" v-for="(item, index) in productData" :key="index">
        <view class="product-image"></view>
        <text class="product-title">{{ item.title }}</text>
        <text class="product-price">¥ {{ item.price }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
const props = defineProps({
  productData: {
    type: Array,
    default: () => [] // 修改为函数返回空数组
  }
})
</script>

<style scoped lang="scss">
// 限时秒杀区域
.flash-sale {
  margin: 24rpx;
  background-color: #fff;
  border-radius: 24rpx;
  padding: 24rpx;
}

.flash-header {
  display: flex;
  align-items: center;
  margin-bottom: 24rpx;
}

.flash-title {
  font-size: 32rpx;
  font-weight: bold;
  color: #FF4500;
  margin-right: auto;
}

.flash-countdown {
  font-size: 24rpx;
  color: #666;
  margin-right: 8rpx;
}

.product-list {
  display: flex;
  overflow-x: auto;
  padding-bottom: 16rpx;
  &::-webkit-scrollbar {
    display: none;
  }
}

.product-item {
  width: 240rpx;
  margin-right: 24rpx;
  flex-shrink: 0;
}

.product-image {
  width: 240rpx;
  height: 300rpx;
  background-color: #f5f5f5;
  border-radius: 16rpx;
  margin-bottom: 12rpx;
}

.product-title {
  font-size: 24rpx;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 8rpx;
}

.product-price {
  font-size: 28rpx;
  color: #FF4500;
  font-weight: bold;
}
</style>

