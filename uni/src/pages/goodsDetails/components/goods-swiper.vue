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
    <swiper-item
      v-for="item in normalizedList"
      :key="item.key"
      class="swiper-item"
      @tap="handleSwiperItemTap(item._rawIndex, item)">
      <!-- 视频类型 -->
      <video
          v-if="item.type === 'video'"
          class="swiper-item-media"
          :src="item.src"
          :poster="item.src"
          controls
          object-fit="contain"></video>
      <!-- 图片类型 -->
        <LazyImage
          v-else
          mode="aspectFit"
          class="swiper-item-media"
          :src="item.src">
        </LazyImage>
      <!-- 文字叠加层 -->
      <view v-if="item.text" class="swiper-text-overlay" :style="item._overlayStyle">
          <text class="swiper-text" :style="item._textStyle">{{ item.text }}</text>
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
const previewImageUrls = ref([])

// 判断是否是视频 URL
const isVideoUrl = (url) => /\.(mp4|mov|webm|m3u8|ts)(\?|$)/i.test(url || '')

const buildItem = (item, index = 0) => {
  if (typeof item === 'string') {
    const src = getUrl(item)
        return { _rawIndex: index, key: `${src || item || 'string'}-${index}`, src, type: isVideoUrl(item) ? 'video' : 'image', text: '', textColor: '#fff', textSize: 14, textPosition: 'bottom', _overlayStyle: {}, _textStyle: {} }
  }
  if (typeof item === 'object' && item !== null) {
    const rawSrc = item.externalUrl ? getExternalUrl(item.externalUrl) : getUrl(item.url || '')
    const rawText = item.text || ''
    const displayText = localText(rawText) || rawText
        const textPosition = item.textPosition || 'bottom'
        const justifyContent = textPosition === 'top' ? 'flex-start' : textPosition === 'center' ? 'center' : 'flex-end'
        const textColor = item.textColor || '#fff'
        const textSize = item.textSize || 14
    return {
      _rawIndex: index,
      key: `${item.id || item.ID || rawSrc || displayText || 'object'}-${index}`,
      src: rawSrc,
      type: item.type || 'image',
      text: displayText,
          textColor,
          textSize,
          textPosition,
          // 轮播文字层样式只服务当前展示；视频签名、图片预览和原始轮播数据不受影响。
          _overlayStyle: { justifyContent },
          _textStyle: { color: textColor, fontSize: `${textSize}px` },
    }
  }
    return { _rawIndex: index, key: `empty-${index}`, src: '', type: 'image', text: '', _overlayStyle: {}, _textStyle: {} }
}

// 预览只需要非视频图片 URL；在归一化列表更新后缓存，避免每次点击重新 filter/map。
const refreshPreviewImageUrls = (items) => {
  previewImageUrls.value = (Array.isArray(items) ? items : [])
    .filter((entry) => entry && entry.type !== 'video' && entry.src)
    .map((entry) => entry.src)
}

const handleSwiperItemTap = (index, item) => {
  if (!item || item.type === 'video') {
    return
  }
  const imageUrls = previewImageUrls.value
  if (!imageUrls.length) {
    return
  }
  const currentSrc = normalizedList.value[index]?.src
  const current = imageUrls.includes(currentSrc) ? currentSrc : imageUrls[0]
  uni.previewImage({
    current,
    urls: imageUrls,
    indicator: 'number',
    loop: true,
  })
}

watch(() => props.list, async (list) => {
  if (!list || list.length === 0) {
    normalizedList.value = []
    previewImageUrls.value = []
    return
  }
  const items = list.map(buildItem)

  // 视频类型先用空 src 占位，避免未签名请求触发 Worker 403
  normalizedList.value = items.map(item =>
    item.type === 'video' ? { ...item, src: '' } : item
  )
  refreshPreviewImageUrls(normalizedList.value)

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
  refreshPreviewImageUrls(signed)
}, { immediate: true })
</script>

<style lang="scss" scoped>
.goods-swiper {
  width: 100%;
  height: 62vh;
  min-height: 520rpx;
  max-height: 860rpx;
  background: #000;
}

.swiper-item {
  width: 100%;
  height: 100%;
  position: relative;
}

.swiper-item-media {
  width: 100%;
  height: 100%;
  object-fit: contain;
  background: #000;
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
