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
    let url = item.href;
    // 检查是否是包含协议的完整URL
    if (url.startsWith('http://') || url.startsWith('https://')) {
      // 查找'#/'之后的部分作为目标路径
      const hashIndex = url.indexOf('#');
      if (hashIndex !== -1) {
        const path = url.substring(hashIndex + 1);
        uni.navigateTo({
          url: path
        });
      } else {
        // 如果没有'#/'，则视为外部链接，使用webview打开
        uni.navigateTo({
          url: '/pages/webview/webview?url=' + encodeURIComponent(url)
        });
      }
    } else {
      // 如果是内部路由，直接跳转
      uni.navigateTo({
        url: url
      });
    }
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
