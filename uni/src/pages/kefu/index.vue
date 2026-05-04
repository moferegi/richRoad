<template>
  <view class="nf-kefu">
    <view class="nf-kefu-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff"></uni-icons>
        </view>
        <text class="nf-navbar-title">{{ $t('kefuTitle') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <scroll-view scroll-y :show-scrollbar="false" class="nf-kefu-scroll">
      <!-- 外部客服列表（shop/kefu 开关控制） -->
      <view v-if="showExternalList && kefuList.length" class="nf-kefu-list">
        <view
          class="nf-kefu-card"
          v-for="(item, index) in kefuList"
          :key="index"
          @tap="contactKefu(item)"
        >
          <!-- 头像区 -->
          <view class="nf-kefu-avatar-wrap">
            <image v-if="resolveAvatar(item)" class="nf-kefu-avatar" :src="resolveAvatar(item)" mode="aspectFill" />
            <view v-else class="nf-kefu-avatar nf-kefu-avatar-fallback" :style="{ background: avatarColor(item.name || String(index)) }">
              <text class="nf-kefu-avatar-fallback-text">{{ avatarInitial(item.name) }}</text>
            </view>
            <view class="nf-kefu-status-dot" :class="'nf-dot-' + normalizeStatus(item.status)"></view>
          </view>
          <!-- 信息区 -->
          <view class="nf-kefu-info">
            <text class="nf-kefu-name">{{ item.name }}</text>
            <view class="nf-kefu-status-row">
              <text class="nf-kefu-status-text" :class="'nf-status-' + normalizeStatus(item.status)">
                {{ getStatusText(item.status) }}
              </text>
            </view>
          </view>
          <!-- 联系按钮 -->
          <view class="nf-kefu-action" :class="{ 'nf-action-disabled': normalizeStatus(item.status) === 'offline' }">
            <text class="nf-kefu-action-text">{{ $t('kefuContact') }}</text>
          </view>
        </view>
      </view>

      <!-- 平台内置客服入口（customerService/config 开关控制；若都开，排在最后） -->
      <view v-if="showPlatEntry" :class="['nf-plat-cs-wrap', showExternalList ? 'nf-plat-cs-wrap--tail' : '']">
        <view v-if="showExternalList" class="nf-section-title">{{ $t('kefuPlatformSection') }}</view>
        <view class="nf-plat-cs-card" @tap="enterPlatChat">
          <view class="nf-plat-cs-avatar">
            <uni-icons type="chat-filled" size="40" color="#e50914" />
          </view>
          <view class="nf-plat-cs-info">
            <text class="nf-plat-cs-name">{{ $t('kefuPlatformName') }}</text>
            <text class="nf-plat-cs-desc">{{ $t('kefuPlatformDesc') }}</text>
          </view>
          <view class="nf-kefu-action">
            <text class="nf-kefu-action-text">{{ $t('kefuConsultNow') }}</text>
          </view>
        </view>
      </view>

      <!-- 空状态 -->
      <view class="nf-kefu-empty" v-if="showNothing || showExternalEmpty">
        <view class="nf-kefu-empty-icon">
          <uni-icons type="chat" size="48" color="rgba(229,9,20,0.4)" />
        </view>
        <text class="nf-kefu-empty-text">{{ $t('kefuEmpty') }}</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getKefuList, getCsConfig, getSysConfigByKey } from '@/api/kefu.js'
import { getUrl } from '@/utils/url.js'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const kefuList = ref([])
const isLoading = ref(true)
const platEnabled = ref(false)
const externalEnabled = ref(true)
const csDefaultAvatarUrl = ref('')

const avatarBgPalette = ['#E50914', '#0EA5E9', '#10B981', '#F59E0B', '#6366F1', '#EC4899', '#14B8A6', '#F97316']

const showExternalList = computed(() => externalEnabled.value)
const showPlatEntry = computed(() => platEnabled.value)
const showNothing = computed(() => !showExternalList.value && !showPlatEntry.value)
const showExternalEmpty = computed(() => showExternalList.value && !kefuList.value.length && !isLoading.value && !showPlatEntry.value)

const normalizeStatus = (status) => {
  const map = { '在线': 'online', 'online': 'online', '离线': 'offline', 'offline': 'offline', '忙碌': 'busy', 'busy': 'busy' }
  return map[status] || 'offline'
}

const avatarInitial = (name) => {
  const text = String(name || '').trim()
  if (!text) return 'K'
  return text.charAt(0).toUpperCase()
}

const avatarColor = (seed) => {
  const text = String(seed || 'kefu')
  let hash = 0
  for (let i = 0; i < text.length; i++) {
    hash = (hash * 31 + text.charCodeAt(i)) >>> 0
  }
  return avatarBgPalette[hash % avatarBgPalette.length]
}

const resolveAvatar = (item) => {
  const raw = item?.avatar || item?.externalAvatar || csDefaultAvatarUrl.value
  if (!raw) return ''
  if (/^(https?:)?\/\//i.test(raw) || /^data:/i.test(raw)) return raw
  return getUrl(raw)
}

const getStatusText = (status) => {
  const key = normalizeStatus(status)
  const map = {
    'online': $t.value('kefuOnline'),
    'offline': $t.value('kefuOffline'),
    'busy': $t.value('kefuBusy')
  }
  return map[key] || status
}

const init = async () => {
  isLoading.value = true
  try {
    const [listRes, cfgRes, switchRes] = await Promise.allSettled([
      getKefuList(),
      getCsConfig(),
      getSysConfigByKey('shop_kefu_enabled')
    ])
    if (listRes.status === 'fulfilled' && listRes.value.code === 0) {
      kefuList.value = Array.isArray(listRes.value.data) ? listRes.value.data : (listRes.value.data?.list || [])
    }
    if (cfgRes.status === 'fulfilled' && cfgRes.value.code === 0) {
      platEnabled.value = !!cfgRes.value.data?.platEnabled
      csDefaultAvatarUrl.value = String(cfgRes.value.data?.defaultAvatarUrl || '')
    }
    if (switchRes.status === 'fulfilled' && switchRes.value.code === 0) {
      externalEnabled.value = String(switchRes.value.data).toLowerCase() === 'true'
    }
  } catch (e) {
    console.error('初始化客服页失败', e)
  } finally {
    isLoading.value = false
  }
}

onShow(() => { init() })

// 进入平台客服聊天室
const enterPlatChat = () => {
  uni.navigateTo({ url: '/pages/kefu/chat' })
}

const contactKefu = (item) => {
  if (normalizeStatus(item.status) === 'offline') {
    uni.showToast({ title: $t.value('kefuOffline'), icon: 'none' })
    return
  }
  if (item.link) {
    // #ifdef H5
    window.open(item.link)
    // #endif
    // #ifndef H5
    uni.navigateTo({
      url: `/pages/webview/index?url=${encodeURIComponent(item.link)}`
    })
    // #endif
  } else {
    uni.showToast({ title: $t.value('kefuContact'), icon: 'none' })
  }
}

const goBack = () => {
  uni.navigateBack()
}
</script>

<style lang="scss">
page {
  background-color: #000;
}

.nf-kefu {
  width: 100%;
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #000;
  position: relative;
}

.nf-kefu-bg {
  position: fixed;
  top: 0; left: 0; right: 0;
  height: 500rpx;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(ellipse at 30% 0%, rgba(229, 9, 20, 0.12) 0%, transparent 60%),
    radial-gradient(ellipse at 70% 10%, rgba(229, 9, 20, 0.08) 0%, transparent 50%);
}

/* ===== 导航栏 ===== */
.nf-navbar {
  background: rgba(0, 0, 0, 0.85);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
  padding: 0 28rpx 16rpx;
  position: sticky;
  top: 0;
  z-index: 99;
}

.nf-navbar-status {
  width: 100%;
  height: var(--status-bar-height, 0px);
}

/* #ifdef MP-WEIXIN */
.nf-navbar-status {
  height: var(--status-bar-height, 44px);
}
/* #endif */

.nf-navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 88rpx;
}

.nf-navbar-back {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.06);
  border: 1rpx solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;

  &:active {
    background: rgba(255, 255, 255, 0.12);
    transform: scale(0.93);
  }
}

.nf-navbar-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #fff;
  letter-spacing: 2rpx;
}

/* ===== 滚动区域 ===== */
.nf-kefu-scroll {
  flex: 1;
  width: 100%;
  scrollbar-width: none;
  -ms-overflow-style: none;

  &::-webkit-scrollbar {
    display: none;
    width: 0;
    height: 0;
  }
}

/* ===== 客服列表 ===== */
.nf-kefu-list {
  padding: 24rpx 28rpx;
}

.nf-section-title {
  color: rgba(255, 255, 255, 0.62);
  font-size: 24rpx;
  letter-spacing: 2rpx;
  margin-bottom: 16rpx;
}

/* ===== 平台专属客服卡片 ===== */
.nf-plat-cs-wrap {
  padding: 40rpx 28rpx;
}

.nf-plat-cs-wrap--tail {
  padding-top: 8rpx;
}

.nf-plat-cs-card {
  display: flex;
  align-items: center;
  background: rgba(229, 9, 20, 0.08);
  border: 1rpx solid rgba(229, 9, 20, 0.3);
  border-radius: 24rpx;
  padding: 40rpx 32rpx;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  transition: transform 0.3s, background 0.3s;

  &:active {
    transform: scale(0.985);
    background: rgba(229, 9, 20, 0.14);
  }
}

.nf-plat-cs-avatar {
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  background: rgba(229, 9, 20, 0.12);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-right: 28rpx;
}

.nf-plat-cs-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.nf-plat-cs-name {
  font-size: 34rpx;
  font-weight: 600;
  color: #fff;
}

.nf-plat-cs-desc {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.55);
}

.nf-kefu-card {
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.06);
  border-radius: 24rpx;
  padding: 32rpx 28rpx;
  margin-bottom: 20rpx;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  transition: transform 0.3s, background 0.3s;

  &:active {
    transform: scale(0.985);
    background: rgba(255, 255, 255, 0.07);
  }
}

/* ===== 头像 ===== */
.nf-kefu-avatar-wrap {
  position: relative;
  flex-shrink: 0;
  margin-right: 24rpx;
}

.nf-kefu-avatar {
  width: 108rpx;
  height: 108rpx;
  border-radius: 50%;
  border: 3rpx solid rgba(229, 9, 20, 0.4);
  box-shadow: 0 0 20rpx rgba(229, 9, 20, 0.15);
}

.nf-kefu-avatar-fallback {
  display: flex;
  align-items: center;
  justify-content: center;
}

.nf-kefu-avatar-fallback-text {
  font-size: 34rpx;
  font-weight: 700;
  color: #fff;
}

.nf-kefu-status-dot {
  position: absolute;
  bottom: 4rpx;
  right: 4rpx;
  width: 24rpx;
  height: 24rpx;
  border-radius: 50%;
  border: 3rpx solid #000;
}

.nf-dot-online {
  background: #22c55e;
  box-shadow: 0 0 8rpx rgba(34, 197, 94, 0.6);
}

.nf-dot-offline {
  background: #6b7280;
}

.nf-dot-busy {
  background: #f59e0b;
  box-shadow: 0 0 8rpx rgba(245, 158, 11, 0.6);
}

/* ===== 信息区 ===== */
.nf-kefu-info {
  flex: 1;
  min-width: 0;
}

.nf-kefu-name {
  font-size: 32rpx;
  font-weight: 700;
  color: #fff;
  letter-spacing: 1rpx;
  margin-bottom: 8rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.nf-kefu-status-row {
  display: flex;
  align-items: center;
}

.nf-kefu-status-text {
  font-size: 24rpx;
  font-weight: 500;
}

.nf-status-online {
  color: #22c55e;
}

.nf-status-offline {
  color: #6b7280;
}

.nf-status-busy {
  color: #f59e0b;
}

/* ===== 按钮 ===== */
.nf-kefu-action {
  flex-shrink: 0;
  margin-left: 16rpx;
  padding: 14rpx 32rpx;
  background: linear-gradient(135deg, #e50914, #b20710);
  border-radius: 40rpx;
  box-shadow: 0 4rpx 16rpx rgba(229, 9, 20, 0.35);
  transition: all 0.3s;

  &:active {
    transform: scale(0.95);
    box-shadow: 0 2rpx 8rpx rgba(229, 9, 20, 0.5);
  }
}

.nf-action-disabled {
  background: rgba(255, 255, 255, 0.08);
  box-shadow: none;

  .nf-kefu-action-text {
    color: rgba(255, 255, 255, 0.3);
  }
}

.nf-kefu-action-text {
  font-size: 26rpx;
  font-weight: 600;
  color: #fff;
  letter-spacing: 1rpx;
}

/* ===== 空状态 ===== */
.nf-kefu-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 200rpx 0;
}

.nf-kefu-empty-icon {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  background: rgba(229, 9, 20, 0.08);
  border: 1rpx solid rgba(229, 9, 20, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24rpx;
}

.nf-kefu-empty-text {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.3);
  letter-spacing: 2rpx;
}
</style>
