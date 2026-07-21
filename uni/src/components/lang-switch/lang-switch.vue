<template>
  <view v-if="visible" class="lang-mask" :class="{ 'lang-mask-closing': !animShow }" @tap="close('cancel')">
    <view class="lang-dialog" :class="{ 'lang-dialog-show': animShow }" @tap.stop @touchmove.stop>
      <view class="lang-header">
        <text class="lang-header-title">{{ $t('selectLanguage') }}</text>
        <view class="lang-header-line"></view>
      </view>

      <scroll-view class="lang-options" scroll-y :show-scrollbar="false">
          <view
            v-for="item in langs"
            :key="item.value"
            class="lang-option"
            :class="{ 'lang-option-active': currentLang === item.value }"
            @tap="selectLang(item.value)"
          >
            <text class="lang-flag">{{ item.flag }}</text>
            <view class="lang-option-info">
              <text class="lang-option-name">{{ item.label }}</text>
              <text class="lang-option-native">{{ item.native }}</text>
            </view>
            <view v-if="currentLang === item.value" class="lang-check">
              <text class="lang-check-icon">✓</text>
            </view>
          </view>
          <view style="height: 16rpx;"></view>
      </scroll-view>

      <!-- 取消按钮：flex-shrink:0 永远不被压缩 -->
      <view class="lang-cancel-wrap" @tap.stop="close('cancel')">
        <text class="lang-cancel">{{ $t('cancel') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, watch, nextTick, computed, onMounted, onUnmounted } from 'vue'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const props = defineProps({
  modelValue: { type: Boolean, default: false }
})
const emit = defineEmits(['update:modelValue', 'change', 'dismiss'])

const visible = ref(false)
const animShow = ref(false)
const hideTimer = ref(null)
const currentLang = computed(() => langStore.locale)
const HIDE_DURATION_MS = 300

const clearHideTimer = () => {
  if (hideTimer.value) {
    clearTimeout(hideTimer.value)
    hideTimer.value = null
  }
}

// 回退硬编码列表（API加载前或加载失败时使用）
const fallbackLangs = [
  { value: 'zh', label: 'ZH', native: 'Simplified Chinese', flag: '🇨🇳' },
  { value: 'zh-TW', label: 'ZH-TW', native: 'Traditional Chinese', flag: '🇹🇼' },
  { value: 'en', label: 'English', native: 'English', flag: '🇬🇧' },
  { value: 'mn', label: 'Монгол', native: 'Монгол хэл', flag: '🇲🇳' },
  { value: 'th', label: 'ไทย', native: 'ภาษาไทย', flag: '🇹🇭' },
  { value: 'hi', label: 'हिन्दी', native: 'हिन्दी', flag: '🇮🇳' },
  { value: 'id', label: 'Bahasa', native: 'Bahasa Indonesia', flag: '🇮🇩' },
  { value: 'vi', label: 'VI', native: 'Tiếng Việt', flag: '🇻🇳' },
  { value: 'ar', label: 'AR', native: 'العربية', flag: '🇸🇦' },
  { value: 'ja', label: 'JA', native: '日本語', flag: '🇯🇵' },
  { value: 'ko', label: 'KO', native: '한국어', flag: '🇰🇷' },
  { value: 'ms', label: 'MS', native: 'Bahasa Melayu', flag: '🇲🇾' },
]

// 优先使用后端启用的语言列表
const langs = computed(() => {
  return langStore.enabledLangs.length > 0 ? langStore.enabledLangs : fallbackLangs
})

onMounted(() => {
  langStore.initLangs()
})

const scrollY = ref(0)

const lockBodyScroll = () => {
  // #ifdef H5
  try {
    scrollY.value = window.scrollY || document.documentElement.scrollTop || 0
    document.body.style.position = 'fixed'
    document.body.style.top = `-${scrollY.value}px`
    document.body.style.width = '100%'
  } catch (_) { /* noop */ }
  // #endif
}

const unlockBodyScroll = () => {
  // #ifdef H5
  try {
    document.body.style.position = ''
    document.body.style.top = ''
    document.body.style.width = ''
    if (scrollY.value > 0) {
      window.scrollTo(0, scrollY.value)
      scrollY.value = 0
    }
  } catch (_) { /* noop */ }
  // #endif
}

watch(() => props.modelValue, (val) => {
  clearHideTimer()
  if (val) {
    visible.value = true
    lockBodyScroll()
    nextTick(() => { animShow.value = true })
  } else {
    animShow.value = false
    hideTimer.value = setTimeout(() => {
      visible.value = false
      hideTimer.value = null
    }, HIDE_DURATION_MS)
  }
}, { immediate: true })

watch(visible, (val) => {
  if (!val) {
    unlockBodyScroll()
  }
})

onUnmounted(() => {
  clearHideTimer()
  unlockBodyScroll()
})

const selectLang = (lang) => {
  langStore.setLocale(lang, { emitRefresh: false })
  emit('change', lang)
  close('select')
  forceRefreshCurrentPage()
}

const forceRefreshCurrentPage = () => {
  const pages = getCurrentPages()
  if (!pages || !pages.length) return
  const current = pages[pages.length - 1]
  const route = String(current?.route || '').trim()
  if (!route) return
  const options = current?.options || {}
  const query = Object.keys(options)
    .filter((key) => options[key] !== undefined && options[key] !== null && String(options[key]).trim() !== '')
    .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(String(options[key]))}`)
    .join('&')
  const url = `/${route}${query ? `?${query}` : ''}`
  uni.reLaunch({ url })
}

const close = (reason = 'cancel') => {
  if (reason === 'cancel' || reason === 'select') {
    langStore.markLanguagePickerPrompted()
  }
  emit('dismiss', reason)
  emit('update:modelValue', false)
}
</script>

<style scoped lang="scss">
.lang-mask {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  z-index: 10001;
  background: rgba(15, 23, 42, 0.36);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  touch-action: none;
}

.lang-mask-closing {
  pointer-events: none;
}

/* 对话框：flex列，max-height限制总高度，overflow:hidden裁剪溢出 */
.lang-dialog {
  width: 100%;
  max-width: 750rpx;
  max-height: 82vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-sizing: border-box;
  background:
    radial-gradient(120% 100% at 100% 0%, rgba(37, 99, 235, 0.14) 0%, transparent 56%),
    #f8fbff;
  border-radius: 32rpx 32rpx 0 0;
  border: 1rpx solid rgba(15, 23, 42, 0.08);
  transform: translateY(100%);
  transition: transform 0.3s cubic-bezier(0.22, 1, 0.36, 1);
  box-shadow: 0 -12rpx 36rpx rgba(15, 23, 42, 0.12);
  touch-action: auto;
}

.lang-dialog-show {
  transform: translateY(0);
}

.lang-header {
  padding: 40rpx 40rpx 0;
  text-align: center;
  flex-shrink: 0;
  box-sizing: border-box;
  width: 100%;
}

.lang-header-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #0f172a;
  letter-spacing: 1rpx;
}

.lang-header-line {
  width: 60rpx;
  height: 6rpx;
  border-radius: 3rpx;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  margin: 20rpx auto 0;
}

.lang-options {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  width: 100%;
  padding: 32rpx 32rpx 0;
  box-sizing: border-box;
}

/* 取消按钮容器：flex-shrink:0 永远在底部 */
.lang-cancel-wrap {
  flex-shrink: 0;
  padding: 20rpx 32rpx calc(36rpx + env(safe-area-inset-bottom, 0px));
  background: #f8fbff;
  display: flex;
  justify-content: center;
  box-sizing: border-box;
  width: 100%;
}

.lang-option {
  display: flex;
  align-items: center;
  padding: 28rpx 24rpx;
  margin-bottom: 16rpx;
  border-radius: 20rpx;
  background: rgba(219, 234, 254, 0.52);
  border: 2rpx solid rgba(37, 99, 235, 0.08);
  box-sizing: border-box;
  width: 100%;
  transition: all 0.25s ease;

  &:active {
    transform: scale(0.97);
    background: rgba(191, 219, 254, 0.65);
  }
}

.lang-option-active {
  background: rgba(37, 99, 235, 0.14);
  border-color: rgba(37, 99, 235, 0.36);
  box-shadow: 0 4rpx 18rpx rgba(37, 99, 235, 0.16);
}

.lang-flag {
  font-size: 48rpx;
  margin-right: 24rpx;
}

.lang-option-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.lang-option-name {
  font-size: 30rpx;
  font-weight: 600;
  color: #0f172a;
}

.lang-option-native {
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.55);
  margin-top: 4rpx;
}

.lang-check {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  background: linear-gradient(135deg, #2563eb, #0ea5e9);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4rpx 12rpx rgba(37, 99, 235, 0.35);
}

.lang-check-icon {
  color: #fff;
  font-size: 28rpx;
  font-weight: bold;
}

.lang-cancel {
  display: inline-block;
  padding: 18rpx 120rpx;
  font-size: 28rpx;
  color: rgba(15, 23, 42, 0.65);
  border-radius: 40rpx;
  border: 2rpx solid rgba(15, 23, 42, 0.1);
  background: rgba(255, 255, 255, 0.78);
  transition: all 0.25s ease;

  &:active {
    background: rgba(219, 234, 254, 0.8);
    color: #0f172a;
  }
}
</style>
