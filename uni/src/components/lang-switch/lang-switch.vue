<template>
  <view v-if="visible" class="lang-mask" @tap.self="close">
    <view class="lang-dialog" :class="{ 'lang-dialog-show': animShow }">
      <view class="lang-header">
        <text class="lang-header-title">{{ $t('selectLanguage') }}</text>
        <view class="lang-header-line"></view>
      </view>

      <view class="lang-options">
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
      </view>

      <view class="lang-footer" @tap="close">
        <text class="lang-cancel">{{ $t('cancel') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, watch, nextTick, computed, onMounted } from 'vue'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const props = defineProps({
  modelValue: { type: Boolean, default: false }
})
const emit = defineEmits(['update:modelValue', 'change'])

const visible = ref(false)
const animShow = ref(false)
const currentLang = computed(() => langStore.locale)

// 回退硬编码列表（API加载前或加载失败时使用）
const fallbackLangs = [
  { value: 'zh', label: 'ZH', native: 'Simplified Chinese', flag: '🇨🇳' },
  { value: 'zh-TW', label: 'ZH-TW', native: 'Traditional Chinese', flag: '🇹🇼' },
  { value: 'en', label: 'English', native: 'English', flag: '🇬🇧' },
  { value: 'mn', label: 'Монгол', native: 'Монгол хэл', flag: '🇲🇳' },
  { value: 'th', label: 'ไทย', native: 'ภาษาไทย', flag: '🇹🇭' },
  { value: 'hi', label: 'हिन्दी', native: 'हिन्दी', flag: '🇮🇳' },
  { value: 'id', label: 'Bahasa', native: 'Bahasa Indonesia', flag: '🇮🇩' },
]

// 优先使用后端启用的语言列表
const langs = computed(() => {
  return langStore.enabledLangs.length > 0 ? langStore.enabledLangs : fallbackLangs
})

onMounted(() => {
  langStore.initLangs()
})

watch(() => props.modelValue, (val) => {
  if (val) {
    visible.value = true
    nextTick(() => { animShow.value = true })
  } else {
    animShow.value = false
    setTimeout(() => { visible.value = false }, 300)
  }
}, { immediate: true })

const selectLang = (lang) => {
  langStore.setLocale(lang)
  emit('change', lang)
  close()
}

const close = () => {
  emit('update:modelValue', false)
}
</script>

<style scoped lang="scss">
.lang-mask {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  z-index: 9999;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.lang-dialog {
  width: 100%;
  max-width: 750rpx;
  background: linear-gradient(180deg, #1e1e1e 0%, #141414 100%);
  border-radius: 32rpx 32rpx 0 0;
  padding: 0 0 env(safe-area-inset-bottom, 0);
  transform: translateY(100%);
  transition: transform 0.3s cubic-bezier(0.22, 1, 0.36, 1);
  box-shadow: 0 -8rpx 40rpx rgba(229, 9, 20, 0.15);
}

.lang-dialog-show {
  transform: translateY(0);
}

.lang-header {
  padding: 40rpx 40rpx 0;
  text-align: center;
}

.lang-header-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #fff;
  letter-spacing: 2rpx;
}

.lang-header-line {
  width: 60rpx;
  height: 6rpx;
  border-radius: 3rpx;
  background: linear-gradient(90deg, #e50914, #ff6b6b);
  margin: 20rpx auto 0;
}

.lang-options {
  padding: 32rpx 32rpx 0;
}

.lang-option {
  display: flex;
  align-items: center;
  padding: 28rpx 24rpx;
  margin-bottom: 16rpx;
  border-radius: 20rpx;
  background: rgba(255, 255, 255, 0.04);
  border: 2rpx solid rgba(255, 255, 255, 0.06);
  transition: all 0.25s ease;

  &:active {
    transform: scale(0.97);
    background: rgba(255, 255, 255, 0.08);
  }
}

.lang-option-active {
  background: rgba(229, 9, 20, 0.12);
  border-color: rgba(229, 9, 20, 0.4);
  box-shadow: 0 4rpx 20rpx rgba(229, 9, 20, 0.15);
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
  color: #fff;
}

.lang-option-native {
  font-size: 22rpx;
  color: rgba(255, 255, 255, 0.45);
  margin-top: 4rpx;
}

.lang-check {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  background: linear-gradient(135deg, #e50914, #ff4d4d);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4rpx 12rpx rgba(229, 9, 20, 0.4);
}

.lang-check-icon {
  color: #fff;
  font-size: 28rpx;
  font-weight: bold;
}

.lang-footer {
  padding: 24rpx 32rpx 40rpx;
  text-align: center;
}

.lang-cancel {
  display: inline-block;
  padding: 18rpx 120rpx;
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.5);
  border-radius: 40rpx;
  border: 2rpx solid rgba(255, 255, 255, 0.1);
  transition: all 0.25s ease;

  &:active {
    background: rgba(255, 255, 255, 0.05);
    color: #fff;
  }
}
</style>
