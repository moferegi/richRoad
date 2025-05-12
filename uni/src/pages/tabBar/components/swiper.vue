<script setup>

</script>

<template>
  <view class="swiper-section">
    <swiper class="swiper" circular autoplay :indicator-dots="true" indicator-color="rgba(255,255,255,0.6)" indicator-active-color="#fff">
      <swiper-item v-for="(item, index) in 3" :key="index">
        <view class="swiper-item">
          <view class="swiper-content">
            <view class="swiper-tag" v-if="index === 2">新人惊喜福利</view>
            <view class="swiper-title">全球好货当地价</view>
            <view class="swiper-title">美妆新人会场</view>
            <view class="swiper-subtitle">雅诗兰黛420豆沙色89元</view>
            <view class="swiper-indicator">{{ index + 1 }} / 3</view>
            <view class="swiper-product-image"></view>
          </view>
        </view>
      </swiper-item>
    </swiper>
  </view>
</template>

<style scoped lang="scss">
// 轮播图区域
.scroll-Y {
  flex: 1;
  width: 100%;
}

.swiper-section {
  padding: 24rpx;
  background-color: #fff;
}

.swiper {
  height: 320rpx;
  border-radius: 16rpx;
  overflow: hidden;
}

.swiper-item {
  width: 100%;
  height: 100%;
  background-color: #FF7E7E; // 桃红色背景，像原型图
}

.swiper-content {
  position: relative;
  width: 100%;
  height: 100%;
  padding: 40rpx;
  box-sizing: border-box;
  color: white;
}

.swiper-tag {
  background-color: #ff4a4a;
  color: white;
  font-size: 20rpx;
  padding: 4rpx 12rpx;
  border-radius: 20rpx;
  display: inline-block;
  margin-bottom: 10rpx;
}

.swiper-title {
  font-size: 36rpx;
  font-weight: bold;
  margin-bottom: 8rpx;
  line-height: 1.2;
}

.swiper-subtitle {
  font-size: 28rpx;
  background-color: rgba(255, 255, 255, 0.2);
  display: inline-block;
  padding: 8rpx 16rpx;
  border-radius: 24rpx;
  margin-top: 8rpx;
}

.swiper-indicator {
  position: absolute;
  bottom: 20rpx;
  right: 20rpx;
  background-color: rgba(0, 0, 0, 0.3);
  color: white;
  font-size: 24rpx;
  padding: 4rpx 16rpx;
  border-radius: 20rpx;
}

.swiper-product-image {
  position: absolute;
  right: 40rpx;
  top: 80rpx;
  width: 240rpx;
  height: 200rpx;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 16rpx;
}
</style>

