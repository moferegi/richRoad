<template>
  <view class="gva-lazy-image">
    <image
      v-if="normalizedSrc"
      class="gva-lazy-image__img"
      :src="normalizedSrc"
      :mode="mode"
      :style="imageStyle"
      :lazy-load="lazyLoad"
      @load="handleLoad"
      @error="handleError"
    />

    <view
      v-if="showStatusMask"
      class="gva-lazy-image__status"
      :class="{ 'gva-lazy-image__status--error': loadStatus === 'error' }"
    >
      <text class="gva-lazy-image__status-text">{{ statusText }}</text>
    </view>
  </view>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useLangStore } from '@/pinia/modules/lang.js'

const props = defineProps({
  src: {
    type: String,
    default: ''
  },
  mode: {
    type: String,
    default: 'aspectFill'
  },
  lazyLoad: {
    type: Boolean,
    default: true
  },
  loadingText: {
    type: String,
    default: ''
  },
  failedText: {
    type: String,
    default: ''
  }
})
const emit = defineEmits(['load', 'error'])

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const normalizedSrc = computed(() => String(props.src || '').trim())
const loadStatus = ref(normalizedSrc.value ? 'loading' : 'error')

watch(normalizedSrc, (val) => {
  loadStatus.value = val ? 'loading' : 'error'
}, { immediate: true })

const normalizedLoadingText = computed(() => {
  if (props.loadingText) return props.loadingText
  const raw = String($t.value('loading') || '').trim()
  if (!raw) return 'Loading'
  return raw.replace(/(\.\.\.|…+)$/g, '').trim() || raw
})

const normalizedFailedText = computed(() => {
  if (props.failedText) return props.failedText
  const raw = String($t.value('loadFailed') || '').trim()
  if (!raw) return 'Load failed'
  const compact = raw.split(/[，,]/)[0].trim()
  return compact || raw
})

const statusText = computed(() => {
  return loadStatus.value === 'error' ? normalizedFailedText.value : normalizedLoadingText.value
})

const showStatusMask = computed(() => {
  return loadStatus.value !== 'loaded'
})

const imageStyle = computed(() => {
  const currentMode = String(props.mode || '').trim()
  if (currentMode === 'widthFix') {
    return {
      width: '100%',
      height: 'auto'
    }
  }
  if (currentMode === 'heightFix') {
    return {
      width: 'auto',
      height: '100%'
    }
  }
  return {
    width: '100%',
    height: '100%'
  }
})

const handleLoad = (event) => {
  loadStatus.value = 'loaded'
  emit('load', event)
}

const handleError = (event) => {
  loadStatus.value = 'error'
  emit('error', event)
}
</script>

<style scoped>
.gva-lazy-image {
  position: relative;
  overflow: hidden;
  display: block;
  background: rgba(148, 163, 184, 0.12);
}

.gva-lazy-image__img {
  display: block;
}

.gva-lazy-image__status {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(226, 232, 240, 0.66);
  padding: 8rpx;
  box-sizing: border-box;
}

.gva-lazy-image__status--error {
  background: rgba(254, 226, 226, 0.72);
}

.gva-lazy-image__status-text {
  font-size: 22rpx;
  line-height: 1.3;
  color: rgba(15, 23, 42, 0.72);
  text-align: center;
}
</style>
