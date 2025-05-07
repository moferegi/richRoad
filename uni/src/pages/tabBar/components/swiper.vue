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
  padding: 12px;
  background-color: #fff;
}

.swiper {
  height: 160px;
  border-radius: 8px;
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
  padding: 20px;
  box-sizing: border-box;
  color: white;
}

.swiper-tag {
  background-color: #ff4a4a;
  color: white;
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 10px;
  display: inline-block;
  margin-bottom: 5px;
}

.swiper-title {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 4px;
  line-height: 1.2;
}

.swiper-subtitle {
  font-size: 14px;
  background-color: rgba(255, 255, 255, 0.2);
  display: inline-block;
  padding: 4px 8px;
  border-radius: 12px;
  margin-top: 4px;
}

.swiper-indicator {
  position: absolute;
  bottom: 10px;
  right: 10px;
  background-color: rgba(0, 0, 0, 0.3);
  color: white;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 10px;
}

.swiper-product-image {
  position: absolute;
  right: 20px;
  top: 40px;
  width: 120px;
  height: 100px;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
}

</style>
