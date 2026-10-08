<template>
  <view class="ad-player-overlay" v-if="visible" @touchmove.stop.prevent="() => {}">
    <view class="ad-player-modal">
      <!-- 顶部信息 -->
      <view class="ad-player-header">
        <text class="ad-player-title">{{ currentAd.title || '广告' }}</text>
        <view class="ad-player-progress">
          <text class="ad-progress-text">{{ currentWatchSeconds }}s / {{ currentAd.minWatchSeconds }}s</text>
        </view>
        <!-- 关闭按钮：仅当达到最低观看时长后显示 -->
        <view v-if="canClose" class="ad-close-btn" @click="closePlayer">
          <text class="ad-close-icon">✕</text>
        </view>
      </view>

      <!-- 媒体内容 -->
      <view class="ad-player-content">
        <!-- 视频播放 -->
        <video
          v-if="currentAd.mediaType === 'video'"
          :src="currentAd.mediaUrl"
          class="ad-video"
          :autoplay="true"
          :controls="false"
          :show-center-play-btn="false"
          :show-play-btn="false"
          :enable-progress-gesture="false"
          object-fit="contain"
          @timeupdate="onVideoTimeUpdate"
          @ended="onVideoEnded"
          @error="onVideoError"
        />
        <!-- 图片展示 -->
        <image
          v-else-if="currentAd.mediaType === 'image'"
          :src="currentAd.mediaUrl"
          class="ad-image"
          mode="aspectFit"
        />
      </view>

      <!-- 底部操作 -->
      <view class="ad-player-footer">
        <!-- 进度条 -->
        <view class="ad-progress-bar">
          <view class="ad-progress-fill" :style="{ width: progressPercent + '%' }" />
        </view>
        <!-- 按钮区域 -->
        <view class="ad-footer-btns">
          <view v-if="!canClose" class="ad-countdown-tip">
            <text class="ad-tip-text">请观看 {{ currentAd.minWatchSeconds - currentWatchSeconds }} 秒后可关闭</text>
          </view>
          <view v-else class="ad-footer-actions">
            <!-- 跳转链接 -->
            <view v-if="currentAd.linkUrl" class="ad-link-btn" @click="openLink">
              <text class="ad-link-text">查看详情</text>
            </view>
            <!-- 显示关闭按钮 -->
            <view class="ad-close-action" @click="closePlayer">
              <text class="ad-close-action-text">关闭</text>
            </view>
          </view>
        </view>
        <!-- 下一个广告按钮 -->
        <view v-if="canClose && adIndex < ads.length - 1" class="ad-next-btn" @click="nextAd">
          <text class="ad-next-text">下一个广告 →</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, watch, onBeforeUnmount } from 'vue'
import { reportWatch } from '@/api/ad.js'

const props = defineProps({
  ads: { type: Array, default: () => [] },
  positionId: { type: Number, default: 0 },
  visible: { type: Boolean, default: false }
})

const emit = defineEmits(['close', 'update:visible'])

const adIndex = ref(0)
const currentWatchSeconds = ref(0)
const watchTimer = ref(null)
const reportedSet = ref(new Set())

const currentAd = computed(() => props.ads[adIndex.value] || {})
const canClose = computed(() => currentWatchSeconds.value >= (currentAd.value.minWatchSeconds || 5))
const progressPercent = computed(() => {
  const min = currentAd.value.minWatchSeconds || 5
  return Math.min(100, Math.round((currentWatchSeconds.value / min) * 100))
})

const startWatchTimer = () => {
  stopWatchTimer()
  watchTimer.value = setInterval(() => {
    currentWatchSeconds.value++
  }, 1000)
}

const stopWatchTimer = () => {
  if (watchTimer.value) {
    clearInterval(watchTimer.value)
    watchTimer.value = null
  }
}

const reportCurrentAd = async () => {
  const adId = currentAd.value.id || currentAd.value.ID
  if (!adId || reportedSet.value.has(adId)) return
  reportedSet.value.add(adId)
  try {
    await reportWatch({
      adVideoId: adId,
      positionId: props.positionId,
      watchedSeconds: currentWatchSeconds.value
    })
  } catch (e) { /* ignore */ }
}

const onVideoTimeUpdate = (e) => {
  const sec = Math.floor(e.detail?.currentTime || 0)
  currentWatchSeconds.value = sec
}

const onVideoEnded = () => {
  // 视频播放完毕，自动下一个
  reportCurrentAd()
  stopWatchTimer()
  if (adIndex.value < props.ads.length - 1) {
    nextAd()
  } else {
    closePlayer()
  }
}

const onVideoError = () => {
  // 视频加载失败，跳过
  stopWatchTimer()
  if (adIndex.value < props.ads.length - 1) {
    nextAd()
  } else {
    closePlayer()
  }
}

const nextAd = () => {
  reportCurrentAd()
  stopWatchTimer()
  currentWatchSeconds.value = 0
  adIndex.value++
}

const openLink = () => {
  const link = currentAd.value.linkUrl
  if (link) {
    // #ifdef H5
    window.open(link, '_blank')
    // #endif
    // #ifdef APP-PLUS
    plus.runtime.openURL(link)
    // #endif
    // #ifdef MP
    uni.setClipboardData({ data: link })
    // #endif
  }
}

const closePlayer = () => {
  if (!canClose.value) return
  reportCurrentAd()
  stopWatchTimer()
  adIndex.value = 0
  currentWatchSeconds.value = 0
  reportedSet.value = new Set()
  emit('update:visible', false)
  emit('close')
}

// 拦截返回键
const handleBackPress = () => {
  if (props.visible) {
    if (canClose.value) {
      closePlayer()
    }
    return true // 阻止返回
  }
  return false
}

watch(() => props.visible, (val) => {
  if (val) {
    adIndex.value = 0
    currentWatchSeconds.value = 0
    reportedSet.value = new Set()
    // 图片类型用计时器，视频类型用 timeupdate
    if (props.ads[0]?.mediaType === 'image') {
      startWatchTimer()
    }
    // 拦截返回键
    // #ifdef APP-PLUS
    plus.key.addEventListener('backbutton', handleBackPress)
    // #endif
  } else {
    stopWatchTimer()
    // #ifdef APP-PLUS
    plus.key.removeEventListener('backbutton', handleBackPress)
    // #endif
  }
})

// 当切换到图片广告时，使用计时器
watch(adIndex, (newIdx) => {
  if (props.visible && props.ads[newIdx]?.mediaType === 'image') {
    startWatchTimer()
  } else {
    stopWatchTimer()
  }
})

onBeforeUnmount(() => {
  stopWatchTimer()
  // #ifdef APP-PLUS
  plus.key.removeEventListener('backbutton', handleBackPress)
  // #endif
})
</script>

<style scoped>
.ad-player-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.85);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
}

.ad-player-modal {
  width: 90%;
  max-width: 700rpx;
  background: #1a1a2e;
  border-radius: 24rpx;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.ad-player-header {
  display: flex;
  align-items: center;
  padding: 20rpx 24rpx;
  background: rgba(0, 0, 0, 0.4);
  position: relative;
}

.ad-player-title {
  flex: 1;
  font-size: 28rpx;
  color: #fff;
  font-weight: 600;
}

.ad-player-progress {
  margin-right: 16rpx;
}

.ad-progress-text {
  font-size: 22rpx;
  color: #ff9500;
  font-weight: 500;
}

.ad-close-btn {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
}

.ad-close-icon {
  font-size: 28rpx;
  color: #fff;
}

.ad-player-content {
  width: 100%;
  aspect-ratio: 16/9;
  background: #000;
  display: flex;
  align-items: center;
  justify-content: center;
}

.ad-video {
  width: 100%;
  height: 100%;
}

.ad-image {
  width: 100%;
  height: 100%;
}

.ad-player-footer {
  padding: 20rpx 24rpx;
}

.ad-progress-bar {
  width: 100%;
  height: 6rpx;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3rpx;
  margin-bottom: 16rpx;
  overflow: hidden;
}

.ad-progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #ff9500, #ff5e3a);
  border-radius: 3rpx;
  transition: width 0.3s ease;
}

.ad-footer-btns {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 60rpx;
}

.ad-countdown-tip {
  text-align: center;
}

.ad-tip-text {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.5);
}

.ad-footer-actions {
  display: flex;
  align-items: center;
  gap: 24rpx;
}

.ad-link-btn {
  padding: 12rpx 28rpx;
  background: rgba(255, 149, 0, 0.15);
  border: 1rpx solid rgba(255, 149, 0, 0.3);
  border-radius: 30rpx;
}

.ad-link-text {
  font-size: 24rpx;
  color: #ff9500;
}

.ad-close-action {
  padding: 12rpx 32rpx;
  background: linear-gradient(135deg, #ff9500, #ff5e3a);
  border-radius: 30rpx;
}

.ad-close-action-text {
  font-size: 26rpx;
  color: #fff;
  font-weight: 600;
}

.ad-next-btn {
  margin-top: 12rpx;
  text-align: center;
  padding: 12rpx;
}

.ad-next-text {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.6);
}
</style>