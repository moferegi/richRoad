<template>
  <view class="nf-marquee-wrap" v-if="enabled && text">
    <view class="nf-marquee-icon">
      <uni-icons type="sound-filled" size="16" :color="textColor || '#e50914'" />
    </view>
    <view class="nf-marquee-track">
      <view
        class="nf-marquee-content"
        :animation="animationData"
      >
        <text class="nf-marquee-text" :style="{ color: textColor || '#fff' }">{{ text }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted, nextTick, getCurrentInstance } from 'vue'

const props = defineProps({
  enabled: { type: Boolean, default: false },
  text: { type: String, default: '' },
  textColor: { type: String, default: '#fff' },
  speed: { type: Number, default: 60 } // px/s
})

const animationData = ref({})
let animTimer = null
let resetTimer = null
let animationToken = 0
const instance = getCurrentInstance()

const getFallbackSize = () => {
  const safeTextLength = Math.max(String(props.text || '').length, 1)
  return {
    textWidth: Math.max(safeTextLength * 16, 80),
    trackWidth: Math.max(uni.getSystemInfoSync().windowWidth || 320, 120)
  }
}

const measureMarqueeSize = () => new Promise((resolve) => {
  const fallback = getFallbackSize()
  if (!instance || !instance.proxy) {
    resolve(fallback)
    return
  }

  const query = uni.createSelectorQuery().in(instance.proxy)
  query.select('.nf-marquee-track').boundingClientRect()
  query.select('.nf-marquee-text').boundingClientRect()
  query.exec((res) => {
    const trackWidth = Number(res?.[0]?.width) || fallback.trackWidth
    const textWidth = Number(res?.[1]?.width) || fallback.textWidth
    resolve({
      trackWidth: Math.max(trackWidth, 1),
      textWidth: Math.max(textWidth, 1)
    })
  })
})

const startAnimation = async () => {
  stopAnimation()
  if (!props.text || !props.enabled) return

  const currentToken = ++animationToken
  const { textWidth, trackWidth } = await measureMarqueeSize()
  if (currentToken !== animationToken) return

  // 使用轨道宽度和真实文字宽度计算位移，避免尾部多余空白造成的长停顿
  const totalDistance = textWidth + trackWidth
  const safeSpeed = Math.max(Number(props.speed) || 60, 10)
  const duration = Math.max((totalDistance / safeSpeed) * 1000, 600)

  const runOnce = () => {
    if (currentToken !== animationToken) return

    const animation = uni.createAnimation({
      duration: 0,
      timingFunction: 'linear'
    })
    animation.translateX(trackWidth).step()
    animationData.value = animation.export()

    resetTimer = setTimeout(() => {
      if (currentToken !== animationToken) return

      const animation2 = uni.createAnimation({
        duration: duration,
        timingFunction: 'linear'
      })
      animation2.translateX(-textWidth).step()
      animationData.value = animation2.export()

      animTimer = setTimeout(() => {
        runOnce()
      }, duration)
    }, 16)
  }

  runOnce()
}

const stopAnimation = () => {
  animationToken++

  if (animTimer) {
    clearTimeout(animTimer)
    animTimer = null
  }

  if (resetTimer) {
    clearTimeout(resetTimer)
    resetTimer = null
  }
}

onMounted(() => {
  nextTick(() => {
    startAnimation()
  })
})

onUnmounted(() => {
  stopAnimation()
})

watch(() => [props.text, props.enabled, props.speed], () => {
  nextTick(() => startAnimation())
})
</script>

<style scoped lang="scss">
.nf-marquee-wrap {
  display: flex;
  align-items: center;
  background: rgba(229, 9, 20, 0.08);
  border: 1rpx solid rgba(229, 9, 20, 0.15);
  border-radius: 12rpx;
  margin: 16rpx 20rpx 0;
  padding: 12rpx 16rpx;
  overflow: hidden;
  height: 60rpx;
}

.nf-marquee-icon {
  flex-shrink: 0;
  margin-right: 12rpx;
  display: flex;
  align-items: center;
}

.nf-marquee-track {
  flex: 1;
  overflow: hidden;
  position: relative;
  height: 36rpx;
}

.nf-marquee-content {
  position: absolute;
  white-space: nowrap;
  display: flex;
  align-items: center;
  height: 100%;
}

.nf-marquee-text {
  font-size: 24rpx;
  white-space: nowrap;
  letter-spacing: 1rpx;
}
</style>
