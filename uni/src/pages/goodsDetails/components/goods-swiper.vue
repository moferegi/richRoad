<template>
  <swiper
      class="goods-swiper"
      circular
      :indicator-dots="true"
      :autoplay="true"
      :interval="3000"
      :duration="500"
      indicator-active-color="#e50914"
      indicator-color="rgba(255,255,255,0.4)">
    <swiper-item v-for="(item, index) in normalizedList" :key="index" class="swiper-item">
      <!-- 视频类型 -->
      <video
          v-if="item.type === 'video'"
          class="swiper-item-media"
          :src="item.src"
          :poster="item.src"
          controls
          object-fit="cover"></video>
      <!-- 图片类型 -->
      <image
          v-else
          mode="aspectFill"
          class="swiper-item-media"
          :src="item.src">
      </image>
      <!-- 文字叠加层 -->
      <view v-if="item.text" class="swiper-text-overlay"
        :style="{ justifyContent: item.textPosition === 'top' ? 'flex-start' : item.textPosition === 'center' ? 'center' : 'flex-end' }">
        <text class="swiper-text" :style="{ color: item.textColor || '#fff', fontSize: (item.textSize || 14) + 'px' }">{{ item.text }}</text>
      </view>
    </swiper-item>
  </swiper>
</template>

<script setup>
import { computed } from "vue"
import { getUrl, getExternalUrl } from "@/utils/url.js"

const props = defineProps({
  list: {
    type: Array,
    default: () => []
  }
})

// 兼容两种 banner 格式：字符串数组 或 对象数组
const normalizedList = computed(() => {
  if (!props.list || props.list.length === 0) return []
  return props.list.map(item => {
    if (typeof item === 'string') {
      return { src: getUrl(item), type: 'image', text: '', textColor: '#fff', textSize: 14, textPosition: 'bottom' }
    }
    if (typeof item === 'object' && item !== null) {
      const src = item.externalUrl ? getExternalUrl(item.externalUrl) : getUrl(item.url || '')
      return {
        src,
        type: item.type || 'image',
        text: item.text || '',
        textColor: item.textColor || '#fff',
        textSize: item.textSize || 14,
        textPosition: item.textPosition || 'bottom'
      }
    }
    return { src: '', type: 'image', text: '' }
  })
})
</script>

<style lang="scss" scoped>
.goods-swiper {
  width: 100%;
  height: 26rem;
}

.swiper-item {
  width: 100%;
  height: 100%;
  position: relative;
}

.swiper-item-media {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.swiper-text-overlay {
  position: absolute;
  left: 0; right: 0; top: 0; bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40rpx 30rpx;
  background: linear-gradient(180deg, transparent 50%, rgba(0,0,0,0.5) 100%);
  pointer-events: none;
}

.swiper-text {
  text-shadow: 0 2px 8px rgba(0,0,0,0.6);
  font-weight: 600;
}
</style>
