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
  position: relative;
  background-color: #FF7E7E; // 桃红色背景
}

.swiper-image {
  width: 100%;
  height: 100%;
  display: block;
}
</style>
