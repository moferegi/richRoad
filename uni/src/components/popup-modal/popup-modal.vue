<template>
  <view class="nf-popup-mask" v-if="visible" @tap.self="onClose" @touchmove.prevent>
    <view class="nf-popup-wrap" :key="'popup_' + currentIndex" @touchmove.stop>
      <view class="nf-popup-glow nf-popup-glow--top"></view>
      <view class="nf-popup-glow nf-popup-glow--bottom"></view>
      <view class="nf-popup-header" v-if="parsedTitle || popup.closeable !== false">
        <view class="nf-popup-title-wrap">
          <text class="nf-popup-title-text" v-if="parsedTitle">{{ parsedTitle }}</text>
        </view>
        <view class="nf-popup-close" v-if="popup.closeable !== false" @tap="onClose">
          <text class="nf-popup-close-icon">×</text>
        </view>
      </view>

      <!-- 图片类型 -->
      <image
        v-if="popup.popupType !== 'content' && (popup.image || popup.externalPath)"
        class="nf-popup-img"
        :src="getUrl(popup.externalPath || popup.image)"
        mode="widthFix"
        @tap="onImageTap"
      />
      <!-- 内容类型 -->
      <view
        v-if="popup.popupType === 'content' && parsedContent"
        class="nf-popup-content"
        @touchmove.stop
      >
        <rich-text :nodes="parsedContent" />
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, nextTick } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getActivePopups } from '@/api/popup.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { getUrl } from '@/utils/url.js'

const props = defineProps({
  position: { type: String, default: '' },
  pagePath: { type: String, default: '' },
  clientType: { type: String, default: 'uni' }
})

const langStore = useLangStore()
const $lt = computed(() => langStore.$lt)

const visible = ref(false)
const popup = ref({})
const loading = ref(false)

// 弹窗队列，支持多个连续弹窗
const popupQueue = ref([])
const currentIndex = ref(0)

const parsedTitle = computed(() => {
  if (!popup.value.title) return ''
  return $lt.value(popup.value.title)
})

const parsedContent = computed(() => {
  if (!popup.value.content) return ''
  return $lt.value(popup.value.content)
})

const getCurrentPagePath = () => {
  if (props.pagePath) return props.pagePath
  const pages = getCurrentPages()
  if (pages.length > 0) {
    return '/' + pages[pages.length - 1].route
  }
  return ''
}

const showNext = () => {
  if (currentIndex.value >= popupQueue.value.length) {
    visible.value = false
    return
  }
  popup.value = popupQueue.value[currentIndex.value]
  visible.value = true
}

const loadPopup = async () => {
  if (loading.value) return
  loading.value = true
  const currentPage = getCurrentPagePath()
  const params = { clientType: props.clientType }
  if (currentPage) {
    params.page = currentPage
  }
  if (props.position) {
    params.position = props.position
  }

  try {
    const res = await getActivePopups(params)
    if (res.code === 0 && res.data && res.data.length > 0) {
      const storageKey = 'popup_shown_' + (currentPage || props.position || 'default')
      // 过滤掉已弹过的 onceOnly 弹窗
      const filtered = res.data.filter(p => {
        if (p.onceOnly) {
          const shown = uni.getStorageSync(storageKey + '_' + p.ID)
          if (shown) return false
        }
        return true
      })
      if (filtered.length === 0) {
        visible.value = false
        popupQueue.value = []
        currentIndex.value = 0
        return
      }
      popupQueue.value = filtered
      currentIndex.value = 0
      // 标记 onceOnly
      filtered.forEach(p => {
        if (p.onceOnly) {
          uni.setStorageSync(storageKey + '_' + p.ID, '1')
        }
      })
      showNext()
    }
  } finally {
    loading.value = false
  }
}

const onClose = () => {
  visible.value = false
  currentIndex.value++
  // 确保 DOM 先销毁再创建下一个弹窗
  nextTick(() => {
    setTimeout(() => {
      showNext()
    }, 200)
  })
}

const onImageTap = () => {
  if (popup.value.link) {
    uni.navigateTo({ url: popup.value.link, fail: () => {} })
  }
}

onShow(() => {
  loadPopup()
})
</script>

<style scoped>
.nf-popup-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background:
    radial-gradient(120% 120% at 50% 18%, rgba(148, 163, 184, 0.2) 0%, rgba(15, 23, 42, 0) 55%),
    rgba(15, 23, 42, 0.62);
  backdrop-filter: blur(4rpx);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 28rpx;
  box-sizing: border-box;
}

.nf-popup-wrap {
  width: 100%;
  max-width: 640rpx;
  max-height: 84vh;
  border-radius: 32rpx;
  overflow: hidden;
  position: relative;
  background: linear-gradient(165deg, #ffffff 0%, #f8fbff 58%, #f3f6ff 100%);
  border: 1rpx solid rgba(255, 255, 255, 0.84);
  box-shadow:
    0 28rpx 84rpx rgba(15, 23, 42, 0.26),
    0 6rpx 18rpx rgba(30, 41, 59, 0.12);
  transform: translateY(14rpx) scale(0.97);
  opacity: 0;
  animation: nf-popup-enter 220ms cubic-bezier(0.22, 1, 0.36, 1) forwards;
}

.nf-popup-glow {
  position: absolute;
  border-radius: 999rpx;
  pointer-events: none;
}

.nf-popup-glow--top {
  width: 320rpx;
  height: 320rpx;
  right: -120rpx;
  top: -188rpx;
  background: radial-gradient(circle, rgba(56, 189, 248, 0.22) 0%, rgba(56, 189, 248, 0) 72%);
}

.nf-popup-glow--bottom {
  width: 280rpx;
  height: 280rpx;
  left: -118rpx;
  bottom: -176rpx;
  background: radial-gradient(circle, rgba(59, 130, 246, 0.14) 0%, rgba(59, 130, 246, 0) 76%);
}

.nf-popup-header {
  display: flex;
  align-items: center;
  gap: 16rpx;
  padding: 24rpx 24rpx 18rpx 30rpx;
  background: linear-gradient(180deg, rgba(248, 250, 252, 0.92) 0%, rgba(255, 255, 255, 0.98) 100%);
  border-bottom: 1rpx solid rgba(148, 163, 184, 0.16);
  position: relative;
  z-index: 2;
}

.nf-popup-title-wrap {
  flex: 1;
  min-width: 0;
}

.nf-popup-title-text {
  font-size: 30rpx;
  font-weight: 700;
  color: #0f172a;
  line-height: 1.4;
  text-shadow: 0 1rpx 0 rgba(255, 255, 255, 0.6);
}

.nf-popup-img {
  width: 100%;
  display: block;
}

.nf-popup-content {
  padding: 28rpx 30rpx 34rpx;
  max-height: 62vh;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  color: #1e293b;
  font-size: 26rpx;
  line-height: 1.72;
  position: relative;
  z-index: 1;
}

.nf-popup-close {
  width: 56rpx;
  height: 56rpx;
  border-radius: 50%;
  background: linear-gradient(180deg, #ffffff 0%, #eef2f7 100%);
  border: 1rpx solid rgba(148, 163, 184, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow:
    inset 0 1rpx 0 rgba(255, 255, 255, 0.9),
    0 8rpx 18rpx rgba(30, 41, 59, 0.12);
  transition: transform 0.16s ease, box-shadow 0.16s ease;
}

.nf-popup-close:active {
  transform: scale(0.94);
  box-shadow:
    inset 0 1rpx 0 rgba(255, 255, 255, 0.86),
    0 4rpx 10rpx rgba(30, 41, 59, 0.12);
}

.nf-popup-close-icon {
  font-size: 34rpx;
  color: #334155;
  line-height: 1;
}

:deep(.nf-popup-content p) {
  margin: 0 0 14rpx;
}

:deep(.nf-popup-content p:last-child) {
  margin-bottom: 0;
}

@keyframes nf-popup-enter {
  from {
    transform: translateY(14rpx) scale(0.97);
    opacity: 0;
  }
  to {
    transform: translateY(0) scale(1);
    opacity: 1;
  }
}
</style>
