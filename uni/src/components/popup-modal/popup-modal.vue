<template>
  <view class="nf-popup-mask" v-if="visible" @tap.self="onClose">
    <view class="nf-popup-wrap">
      <image
        v-if="popup.image || popup.externalPath"
        class="nf-popup-img"
        :src="popup.externalPath || getUrl(popup.image)"
        mode="widthFix"
        @tap="onImageTap"
      />
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
import { ref, computed, onMounted } from 'vue'
import { getActivePopups } from '@/api/popup.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { getUrl } from '@/utils/url.js'

const props = defineProps({
  position: { type: String, default: 'home' },
  clientType: { type: String, default: 'uni' }
})

const langStore = useLangStore()
const $lt = computed(() => langStore.$lt)

const visible = ref(false)
const popup = ref({})

const parsedTitle = computed(() => {
  if (!popup.value.title) return ''
  return $lt.value(popup.value.title)
})

const loadPopup = async () => {
  const storageKey = 'popup_shown_' + props.position
  const res = await getActivePopups({ position: props.position, clientType: props.clientType })
  if (res.code === 0 && res.data && res.data.length > 0) {
    const p = res.data[0]
    // 如果设置为只弹一次，检查是否已展示过
    if (p.onceOnly) {
      const shown = uni.getStorageSync(storageKey)
      if (shown === String(p.ID)) return
    }
    popup.value = p
    visible.value = true
    if (p.onceOnly) {
      uni.setStorageSync(storageKey, String(p.ID))
    }
  }
}

const onClose = () => {
  visible.value = false
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
