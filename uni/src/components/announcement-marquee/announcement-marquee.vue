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
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue'

const props = defineProps({
  enabled: { type: Boolean, default: false },
  text: { type: String, default: '' },
  textColor: { type: String, default: '#fff' },
  speed: { type: Number, default: 60 } // px/s
})

const animationData = ref({})
let animTimer = null

const startAnimation = () => {
  stopAnimation()
  if (!props.text || !props.enabled) return

  // 估算文字宽度 (每个字符约16px)
  const textWidth = props.text.length * 16 + 200
  const screenWidth = uni.getSystemInfoSync().windowWidth
  const totalDistance = textWidth + screenWidth
  const duration = (totalDistance / props.speed) * 1000

  const runOnce = () => {
    const animation = uni.createAnimation({
      duration: 0,
      timingFunction: 'linear'
    })
    animation.translateX(screenWidth).step()
    animationData.value = animation.export()

    setTimeout(() => {
      const animation2 = uni.createAnimation({
        duration: duration,
        timingFunction: 'linear'
      })
      animation2.translateX(-textWidth).step()
      animationData.value = animation2.export()

      animTimer = setTimeout(() => {
        runOnce()
      }, duration)
    }, 50)
  }

  runOnce()
}

const stopAnimation = () => {
  if (animTimer) {
    clearTimeout(animTimer)
    animTimer = null
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
