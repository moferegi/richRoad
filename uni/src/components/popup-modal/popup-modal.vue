<template>
  <view class="nf-popup-mask" v-if="visible" @tap.self="onClose" @touchmove.stop.prevent>
    <view class="nf-popup-wrap" :key="'popup_' + currentIndex">
      <!-- 图片类型 -->
      <image
        v-if="popup.popupType !== 'content' && (popup.image || popup.externalPath)"
        class="nf-popup-img"
        :src="popup.externalPath ? getExternalUrl(popup.externalPath) : getUrl(popup.image)"
        mode="widthFix"
        @tap="onImageTap"
      />
      <!-- 内容类型 -->
      <view v-if="popup.popupType === 'content' && parsedContent" class="nf-popup-content">
        <rich-text :nodes="parsedContent" />
      </view>
      <view class="nf-popup-title" v-if="parsedTitle">
        <text class="nf-popup-title-text">{{ parsedTitle }}</text>
      </view>
      <view class="nf-popup-close" v-if="popup.closeable !== false" @tap="onClose">
        <text class="nf-popup-close-icon">×</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { getActivePopups } from '@/api/popup.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { getUrl, getExternalUrl } from '@/utils/url.js'

const props = defineProps({
  position: { type: String, default: '' },
  pagePath: { type: String, default: '' },
  clientType: { type: String, default: 'uni' }
})

const langStore = useLangStore()
const $lt = computed(() => langStore.$lt)

const visible = ref(false)
const popup = ref({})

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
  const currentPage = getCurrentPagePath()
  const params = { clientType: props.clientType }
  if (currentPage) {
    params.page = currentPage
  }
  if (props.position) {
    params.position = props.position
  }

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
    if (filtered.length === 0) return
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

onMounted(() => {
  loadPopup()
})
</script>

<style scoped>
.nf-popup-mask {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.7);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
}
.nf-popup-wrap {
  position: relative;
  width: 580rpx;
  border-radius: 24rpx;
  overflow: hidden;
  background: #1a1a1a;
}
.nf-popup-img {
  width: 100%;
}
.nf-popup-content {
  padding: 24rpx 32rpx;
  max-height: 600rpx;
  overflow-y: auto;
  color: #fff;
  font-size: 26rpx;
  line-height: 1.6;
}
.nf-popup-title {
  padding: 24rpx 32rpx;
  text-align: center;
}
.nf-popup-title-text {
  font-size: 28rpx;
  color: #fff;
}
.nf-popup-close {
  position: absolute;
  top: 12rpx;
  right: 12rpx;
  width: 56rpx;
  height: 56rpx;
  border-radius: 50%;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}
.nf-popup-close-icon {
  font-size: 36rpx;
  color: #fff;
  line-height: 1;
}
</style>
