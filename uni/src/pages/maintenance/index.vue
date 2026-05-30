<template>
  <view class="nf-maintenance">
    <!-- 背景图 -->
    <image
      v-if="bgImage"
      class="nf-maintenance-bg"
      :src="bgImage"
      mode="aspectFill"
    />
    <view class="nf-maintenance-overlay"></view>

    <view class="nf-maintenance-content">
      <!-- Logo -->
      <view class="nf-maintenance-logo">
        <uni-icons type="info" size="60" color="#e50914" />
      </view>

      <text class="nf-maintenance-title">{{ $t('maintenanceTitle') }}</text>
      <text class="nf-maintenance-desc">{{ $t('maintenanceDesc') }}</text>

      <!-- 弹窗信息 -->
      <view v-if="popupEnabled" class="nf-maintenance-popup">
        <text class="nf-maintenance-popup-title">{{ popupDisplay.title }}</text>
        <text class="nf-maintenance-popup-content">{{ popupDisplay.content }}</text>
      </view>

      <!-- 按钮区域 -->
      <view class="nf-maintenance-actions">
        <view class="nf-maintenance-btn nf-maintenance-btn-kefu" @tap="goKefu">
          <text>{{ $t('contactCustomerService') }}</text>
        </view>
        <view
          v-if="homeBtnEnabled"
          class="nf-maintenance-btn nf-maintenance-btn-home"
          @tap="goHome"
        >
          <text>{{ $t('goHome') }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { getUrl } from '@/utils/url.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)
const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')

const bgImage = ref('')
const popupEnabled = ref(false)
const popupTitle = ref('')
const popupContent = ref('')
const homeBtnEnabled = ref(false)
const lastLoadedLocale = ref('')

const popupDisplay = computed(() => ({
  // 维护页弹窗只读展示缓存配置文案；维护模式判断和按钮跳转仍读取原始开关字段。
  title: $lt.value(popupTitle.value),
  content: $lt.value(popupContent.value),
}))

const resolveImageUrl = (url) => {
  const raw = String(url || '').trim()
  if (!raw) return ''
  return getUrl(raw)
}

const loadConfig = () => {
  try {
    // 从缓存中读取维护模式配置（由 request.js 在收到503时写入）
    const cached = uni.getStorageSync('maintenance_config')
    if (cached) {
      const data = JSON.parse(cached)
      bgImage.value = resolveImageUrl(data.maintenance_bg_image || '')
      popupEnabled.value = data.maintenance_popup_enabled === 'true'
      popupTitle.value = data.maintenance_popup_title || ''
      popupContent.value = data.maintenance_popup_content || ''
      homeBtnEnabled.value = data.maintenance_home_btn_enabled === 'true'
    }
  } catch (e) {
    console.error('加载维护配置失败', e)
  }
}

const goKefu = () => {
  uni.navigateTo({ url: '/pages/kefu/index' })
}

const goHome = () => {
  uni.removeStorageSync('maintenance_config')
  uni.switchTab({ url: '/pages/tabBar/index' })
}

onMounted(() => {
  loadConfig()
  lastLoadedLocale.value = locale.value
})

onShow(() => {
  const localeChanged = !!lastLoadedLocale.value && lastLoadedLocale.value !== locale.value
  if (localeChanged) {
    loadConfig()
  }
  lastLoadedLocale.value = locale.value
})
</script>

<style lang="scss" scoped>
.nf-maintenance {
  position: relative;
  min-height: 100vh;
  background: #141414;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.nf-maintenance-bg {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  width: 100%; height: 100%;
}
.nf-maintenance-overlay {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.6);
}
.nf-maintenance-content {
  position: relative;
  z-index: 10;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60rpx 40rpx;
}
.nf-maintenance-logo {
  margin-bottom: 40rpx;
}
.nf-maintenance-title {
  font-size: 40rpx;
  font-weight: bold;
  color: #fff;
  margin-bottom: 20rpx;
}
.nf-maintenance-desc {
  font-size: 28rpx;
  color: rgba(255,255,255,0.7);
  text-align: center;
  margin-bottom: 40rpx;
}
.nf-maintenance-popup {
  background: rgba(255,255,255,0.1);
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 40rpx;
  width: 100%;
}
.nf-maintenance-popup-title {
  font-size: 32rpx;
  font-weight: bold;
  color: #e50914;
  margin-bottom: 16rpx;
  display: block;
}
.nf-maintenance-popup-content {
  font-size: 26rpx;
  color: rgba(255,255,255,0.8);
  line-height: 1.6;
}
.nf-maintenance-actions {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
  width: 100%;
}
.nf-maintenance-btn {
  height: 88rpx;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 30rpx;
  font-weight: bold;
}
.nf-maintenance-btn-kefu {
  background: #e50914;
  color: #fff;
}
.nf-maintenance-btn-home {
  background: rgba(255,255,255,0.15);
  color: #fff;
}
</style>
