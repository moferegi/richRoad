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
        <LazyImage
          v-else
          mode="aspectFill"
          class="swiper-item-media"
          :src="item.src">
        </LazyImage>
      <!-- 文字叠加层 -->
      <view v-if="item.text" class="swiper-text-overlay"
        :style="{ justifyContent: item.textPosition === 'top' ? 'flex-start' : item.textPosition === 'center' ? 'center' : 'flex-end' }">
        <text class="swiper-text" :style="{ color: item.textColor || '#fff', fontSize: (item.textSize || 14) + 'px' }">{{ item.text }}</text>
      </view>
    </swiper-item>
  </swiper>
</template>

<script setup>
import { ref, watch } from "vue"
import { getUrl, getExternalUrl } from "@/utils/url.js"
import { localText } from "@/utils/i18n.js"
import { signURL } from "@/api/fileUpload.js"
import LazyImage from '@/components/lazy-image/lazy-image.vue'

const props = defineProps({
  list: {
    type: Array,
    default: () => []
  }
})

const normalizedList = ref([])

// 判断是否是视频 URL
const isVideoUrl = (url) => /\.(mp4|mov|webm|m3u8|ts)(\?|$)/i.test(url || '')

const buildItem = (item) => {
  if (typeof item === 'string') {
    return { src: getUrl(item), type: isVideoUrl(item) ? 'video' : 'image', text: '', textColor: '#fff', textSize: 14, textPosition: 'bottom' }
  }
  if (typeof item === 'object' && item !== null) {
    const rawSrc = item.externalUrl ? getExternalUrl(item.externalUrl) : getUrl(item.url || '')
    const rawText = item.text || ''
    const displayText = localText(rawText) || rawText
    return {
      src: rawSrc,
      type: item.type || 'image',
      text: displayText,
      textColor: item.textColor || '#fff',
      textSize: item.textSize || 14,
      textPosition: item.textPosition || 'bottom'
    }
  }
  return { src: '', type: 'image', text: '' }
}

watch(() => props.list, async (list) => {
  if (!list || list.length === 0) {
    normalizedList.value = []
    return
  }
  const items = list.map(buildItem)

  // 视频类型先用空 src 占位，避免未签名请求触发 Worker 403
  normalizedList.value = items.map(item =>
    item.type === 'video' ? { ...item, src: '' } : item
  )

  // 异步对视频类型签名，完成后整体替换
  const signed = await Promise.all(items.map(async (item) => {
    if (item.type !== 'video' || !item.src) return item
    try {
      const res = await signURL(item.src)
      if (res.code === 0 && res.data && res.data.url) {
        return { ...item, src: res.data.url }
      }
    } catch (e) {
      console.warn('banner视频签名失败', e)
    }
    return item
  }))
  normalizedList.value = signed
}, { immediate: true })
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
