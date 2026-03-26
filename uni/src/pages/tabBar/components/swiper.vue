<template>
  <view class="swiper-section">
    <swiper class="swiper" circular autoplay :indicator-dots="true" indicator-color="rgba(255,255,255,0.6)" indicator-active-color="#fff">
      <swiper-item v-for="(item, index) in props.lists" :key="index" @click="handleSwiperClick(item)">
        <view class="swiper-item">
          <image class="swiper-image" :src="getUrl(item.src)" mode="aspectFill"></image>
        </view>
      </swiper-item>
    </swiper>
  </view>
</template>

<script setup>
import {getUrl} from "@/utils/url.js"

const props = defineProps({
  lists: {
    type: Array,
    default: () => []
  }
})

// 处理轮播图点击事件
const handleSwiperClick = (item) => {
  if (item.href && item.href.trim() !== '') {
    uni.navigateTo({
      url: item.href
    });
  }
};

</script>

<style scoped lang="scss">
.swiper-section {
  padding: 24rpx 28rpx;
  background: #000;
}

.swiper {
  height: 340rpx;
  border-radius: 20rpx;
  overflow: hidden;
  box-shadow: 0 12rpx 48rpx rgba(0, 0, 0, 0.5), 0 0 60rpx rgba(229, 9, 20, 0.06);
}

.swiper-item {
  width: 100%;
  height: 100%;
  position: relative;
  background: #111;

  &::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 40%;
    background: linear-gradient(transparent, rgba(0, 0, 0, 0.5));
    pointer-events: none;
    z-index: 1;
  }
}

.swiper-image {
  width: 100%;
  height: 100%;
  display: block;
  border-radius: 20rpx;
}
</style>
