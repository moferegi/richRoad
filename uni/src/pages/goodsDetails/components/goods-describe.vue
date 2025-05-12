<script setup>

</script>

<template>
  <view class="detail-section">
    <view class="detail-title">图文详情</view>
    <view class="detail-content">
      <view class="detail-image" :style="{ backgroundColor: '#f5f5f5', height: '200px' }"></view>
      <view class="detail-image" :style="{ backgroundColor: '#f0f0f0', height: '300px' }"></view>
    </view>
  </view>
</template>

<style scoped lang="scss">
/* 图文详情 */
.detail-section {
  background-color: #fff;
  margin-top: 20rpx;
  padding-bottom: 40rpx;
}

.detail-title {
  text-align: center;
  font-size: 32rpx;
  color: #333;
  padding: 30rpx 0;
  position: relative;
}

.detail-title::before,
.detail-title::after {
  content: '';
  position: absolute;
  top: 50%;
  width: 120rpx;
  height: 1px; // 保留 1px 线条，细腻、轻盈
  background-color: #ddd;
}

.detail-title::before {
  left: 100rpx;
}

.detail-title::after {
  right: 100rpx;
}

.detail-content {
  padding: 0 30rpx;
}

.detail-image {
  width: 100%;
  margin-bottom: 20rpx;
  border-radius: 8rpx;
}
</style>

