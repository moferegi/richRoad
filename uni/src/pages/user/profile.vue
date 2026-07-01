<template>
  <view class="profile-page">
    <view class="profile-bg"></view>

    <view class="profile-navbar">
      <view class="profile-navbar-status"></view>
      <view class="profile-navbar-content">
        <view class="profile-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#0f172a" />
        </view>
        <text class="profile-navbar-title">{{ $t('personalInfo') }}</text>
        <view class="profile-nav-space"></view>
      </view>
    </view>

    <view class="profile-content">
      <view class="profile-head-card">
        <text class="profile-head-label">{{ $t('nickname') }}</text>
        <text class="profile-head-name">{{ displayName }}</text>
      </view>

      <view class="profile-card">
        <view class="profile-item" @tap="editNickName">
          <text class="profile-item-label">{{ $t('nickname') }}</text>
          <view class="profile-item-right">
            <text class="profile-item-value">{{ displayName }}</text>
            <uni-icons type="right" size="14" color="rgba(15,23,42,0.34)" />
          </view>
        </view>

        <view class="profile-item">
          <text class="profile-item-label">{{ $t('phone') }}</text>
          <view class="profile-item-right">
            <template v-if="hasPhone">
              <text class="profile-item-value">{{ phoneDisplay }}</text>
              <view class="eye-btn" @tap="togglePhoneVisible">
                <uni-icons :type="showPhoneRaw ? 'eye-slash' : 'eye'" size="18" color="#2563eb" />
              </view>
            </template>
            <template v-else>
              <text class="profile-item-empty">{{ $t('notSet') }}</text>
              <view class="set-btn" @tap="goSetPhone">{{ $t('setPhone') }}</view>
            </template>
          </view>
        </view>

        <view class="profile-item">
          <text class="profile-item-label">{{ $t('email') }}</text>
          <view class="profile-item-right">
            <template v-if="hasEmail">
              <text class="profile-item-value">{{ emailDisplay }}</text>
              <view class="eye-btn" @tap="toggleEmailVisible">
                <uni-icons :type="showEmailRaw ? 'eye-slash' : 'eye'" size="18" color="#2563eb" />
              </view>
            </template>
            <template v-else>
              <text class="profile-item-empty">{{ $t('notSet') }}</text>
              <view class="set-btn" @tap="openEmailModal">{{ $t('confirmSet') }}</view>
            </template>
          </view>
        </view>

        <view class="profile-item no-border">
          <text class="profile-item-label">{{ $t('registerTime') }}</text>
          <view class="profile-item-right">
            <text class="profile-item-value">{{ formatTime(userInfo.createdAt || userInfo.CreatedAt) }}</text>
          </view>
        </view>
      </view>

      <view class="profile-action-card">
        <view class="profile-action-item" @tap="goChangePassword">
          <uni-icons type="locked" size="18" color="#2563eb" />
          <text class="profile-action-text">{{ $t('changePassword') }}</text>
          <uni-icons type="right" size="14" color="rgba(15,23,42,0.34)" />
        </view>
      </view>

      <view class="profile-logout" @tap="logout">
        <text>{{ $t('logout') }}</text>
      </view>
    </view>

    <view v-if="showNickNameModal" class="modal-mask" @tap.stop="showNickNameModal = false">
      <view class="modal-box" @tap.stop>
        <text class="modal-title">{{ $t('editNickname') }}</text>
        <input class="modal-input" v-model="newNickName" :placeholder="$t('enterNickname')" maxlength="20" />
        <view class="modal-btns">
          <view class="modal-btn modal-cancel" @tap="showNickNameModal = false">
            <text>{{ $t('cancel') }}</text>
          </view>
          <view class="modal-btn modal-confirm" @tap="saveNickName">
            <text>{{ $t('confirm') }}</text>
          </view>
        </view>
      </view>
    </view>

    <view v-if="showEmailModal" class="modal-mask" @tap.stop="showEmailModal = false">
      <view class="modal-box" @tap.stop>
        <text class="modal-title">{{ $t('email') }}</text>
        <input class="modal-input" v-model="newEmail" placeholder="example@mail.com" />
        <view class="modal-btns">
          <view class="modal-btn modal-cancel" @tap="showEmailModal = false">
            <text>{{ $t('cancel') }}</text>
          </view>
          <view class="modal-btn modal-confirm" @tap="saveEmail">
            <text>{{ $t('confirm') }}</text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getUserInfo, setClientUserInfo } from '@/api/base'
import { useUserStore } from '@/pinia/modules/user'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const userStore = useUserStore()
const userInfo = ref({})
const showNickNameModal = ref(false)
const showEmailModal = ref(false)
const newNickName = ref('')
const newEmail = ref('')
const showPhoneRaw = ref(false)
const showEmailRaw = ref(false)

const displayName = computed(() => {
  return userInfo.value.nickname || userInfo.value.username || '-'
})

const hasPhone = computed(() => {
  return !!String(userInfo.value.phone || '').trim()
})

const hasEmail = computed(() => {
  return !!String(userInfo.value.email || '').trim()
})

const maskedPhone = computed(() => {
  const phone = userInfo.value.phone || ''
  if (phone.length >= 7) {
    return phone.substring(0, 3) + '****' + phone.substring(phone.length - 4)
  }
  return phone || $t.value('notSet')
})

const phoneDisplay = computed(() => {
  if (!hasPhone.value) return $t.value('notSet')
  return showPhoneRaw.value ? (userInfo.value.phone || '') : maskedPhone.value
})

const maskedEmail = computed(() => {
  const value = String(userInfo.value.email || '').trim()
  if (!value) return $t.value('notSet')
  const parts = value.split('@')
  if (parts.length !== 2) return value
  const local = parts[0]
  const host = parts[1]
  if (local.length <= 2) {
    return `${local[0] || '*'}***@${host}`
  }
  return `${local.slice(0, 2)}***${local.slice(-1)}@${host}`
})

const emailDisplay = computed(() => {
  if (!hasEmail.value) return $t.value('notSet')
  return showEmailRaw.value ? (userInfo.value.email || '') : maskedEmail.value
})

const formatTime = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const loadUserInfo = async () => {
  try {
    const res = await getUserInfo()
    if (res.code === 0 && res.data) {
      userInfo.value = res.data
    }
  } catch (e) {
    console.error('加载用户信息失败', e)
  }
}

onShow(() => {
  const token = uni.getStorageSync('x-token')
  if (!token) {
    uni.redirectTo({ url: '/pages/user/login' })
    return
  }
  loadUserInfo()
})

const togglePhoneVisible = () => {
  showPhoneRaw.value = !showPhoneRaw.value
}

const toggleEmailVisible = () => {
  showEmailRaw.value = !showEmailRaw.value
}

const editNickName = () => {
  newNickName.value = displayName.value === '-' ? '' : displayName.value
  showNickNameModal.value = true
}

const saveNickName = async () => {
  if (!newNickName.value.trim()) {
    uni.showToast({ title: $t.value('enterNickname'), icon: 'none' })
    return
  }
  try {
    const value = newNickName.value.trim()
    const res = await setClientUserInfo({ key: 'nickname', value })
    if (res.code === 0) {
      userInfo.value.nickname = value
      uni.setStorageSync('userInfo', { ...(uni.getStorageSync('userInfo') || {}), nickname: value })
      showNickNameModal.value = false
      uni.showToast({ title: $t.value('updateSuccess'), icon: 'none' })
    }
  } catch (e) {
    uni.showToast({ title: $t.value('operationFailed'), icon: 'none' })
  }
}

const goSetPhone = () => {
  const phone = encodeURIComponent(String(userInfo.value.phone || ''))
  uni.navigateTo({ url: `/pages/user/setPhone?phone=${phone}` })
}

const openEmailModal = () => {
  newEmail.value = String(userInfo.value.email || '').trim()
  showEmailModal.value = true
}

const saveEmail = async () => {
  const email = String(newEmail.value || '').trim()
  if (!email) {
    uni.showToast({ title: $t.value('notSet'), icon: 'none' })
    return
  }

  const emailRule = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRule.test(email)) {
    uni.showToast({ title: $t.value('operationFailed'), icon: 'none' })
    return
  }

  try {
    const res = await setClientUserInfo({ key: 'email', value: email })
    if (res.code === 0) {
      userInfo.value.email = email
      uni.setStorageSync('userInfo', { ...(uni.getStorageSync('userInfo') || {}), email })
      showEmailModal.value = false
      uni.showToast({ title: $t.value('setSuccess'), icon: 'none' })
    }
  } catch (e) {
    uni.showToast({ title: $t.value('operationFailed'), icon: 'none' })
  }
}

const goChangePassword = () => {
  uni.navigateTo({ url: '/pages/user/changePassword' })
}

const logout = () => {
  uni.showModal({
    title: $t.value('confirmLogout'),
    confirmColor: '#e50914',
    cancelText: $t.value('cancel'),
    confirmText: $t.value('confirm'),
    success: (res) => {
      if (res.confirm) {
        uni.removeStorageSync('x-token')
        uni.removeStorageSync('userInfo')
        userStore.setToken('')
        uni.reLaunch({ url: '/pages/user/login' })
      }
    }
  })
}

const goBack = () => {
  uni.navigateBack({ fail: () => uni.switchTab({ url: '/pages/learning/profile' }) })
}
</script>

<style lang="scss" scoped>
.profile-page {
  min-height: 100vh;
  background: #f4f7fb;
}

.profile-bg {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 420rpx;
  pointer-events: none;
  background:
    radial-gradient(120% 80% at 100% -10%, rgba(14, 165, 233, 0.2) 0%, transparent 60%),
    radial-gradient(120% 80% at 0% 0%, rgba(37, 99, 235, 0.2) 0%, transparent 60%);
}

.profile-navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  background: rgba(244, 247, 251, 0.88);
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
  border-bottom: 1rpx solid rgba(148, 163, 184, 0.2);
}

.profile-navbar-status {
  height: var(--status-bar-height, 44rpx);
}

.profile-navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 88rpx;
  padding: 0 24rpx;
}

.profile-back,
.profile-nav-space {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.profile-back {
  border-radius: 50%;
  background: #ffffff;
  box-shadow: 0 8rpx 20rpx rgba(15, 23, 42, 0.08);
}

.profile-navbar-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #0f172a;
}

.profile-content {
  padding: calc(var(--status-bar-height, 44rpx) + 88rpx + 40rpx) 24rpx 60rpx;
}

.profile-head-card {
  border-radius: 20rpx;
  padding: 24rpx;
  margin-bottom: 18rpx;
  background: linear-gradient(135deg, #2563eb, #0ea5e9);
  color: #fff;
  box-shadow: 0 16rpx 28rpx rgba(37, 99, 235, 0.24);
}

.profile-head-label {
  display: block;
  font-size: 24rpx;
  opacity: 0.85;
}

.profile-head-name {
  display: block;
  margin-top: 8rpx;
  font-size: 34rpx;
  font-weight: 700;
}

.profile-card,
.profile-action-card {
  background: #ffffff;
  border-radius: 18rpx;
  box-shadow: 0 12rpx 28rpx rgba(15, 23, 42, 0.08);
}

.profile-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 30rpx 24rpx;
  border-bottom: 1rpx solid #eef2f7;
}

.profile-item.no-border {
  border-bottom: none;
}

.profile-item-label {
  color: rgba(15, 23, 42, 0.62);
  font-size: 27rpx;
}

.profile-item-right {
  display: flex;
  align-items: center;
  gap: 10rpx;
}

.profile-item-value {
  font-size: 27rpx;
  color: #0f172a;
}

.profile-item-empty {
  font-size: 27rpx;
  color: rgba(15, 23, 42, 0.42);
}

.set-btn {
  padding: 8rpx 16rpx;
  border-radius: 999rpx;
  font-size: 22rpx;
  color: #2563eb;
  border: 1rpx solid rgba(37, 99, 235, 0.35);
  background: rgba(37, 99, 235, 0.08);
}

.eye-btn {
  width: 44rpx;
  height: 44rpx;
  border-radius: 22rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(37, 99, 235, 0.1);
}

.profile-action-card {
  margin-top: 14rpx;
}

.profile-action-item {
  display: flex;
  align-items: center;
  gap: 14rpx;
  padding: 28rpx 24rpx;
}

.profile-action-text {
  flex: 1;
  color: #0f172a;
  font-size: 28rpx;
  font-weight: 600;
}

.profile-logout {
  margin-top: 32rpx;
  height: 88rpx;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #e11d48;
  font-size: 30rpx;
  font-weight: 600;
  background: rgba(225, 29, 72, 0.1);
  border: 1rpx solid rgba(225, 29, 72, 0.26);
}

.modal-mask {
  position: fixed;
  inset: 0;
  z-index: 200;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(15, 23, 42, 0.42);
}

.modal-box {
  width: 600rpx;
  border-radius: 20rpx;
  padding: 36rpx;
  background: #ffffff;
  box-shadow: 0 20rpx 40rpx rgba(15, 23, 42, 0.22);
}

.modal-title {
  text-align: center;
  font-size: 32rpx;
  color: #0f172a;
  font-weight: 700;
}

.modal-input {
  margin-top: 24rpx;
  height: 84rpx;
  border-radius: 12rpx;
  background: #f8fafc;
  border: 1rpx solid #dbe3ee;
  padding: 0 24rpx;
  color: #0f172a;
  font-size: 28rpx;
}

.modal-btns {
  display: flex;
  gap: 16rpx;
  margin-top: 28rpx;
}

.modal-btn {
  flex: 1;
  height: 78rpx;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28rpx;
  font-weight: 600;
}

.modal-cancel {
  color: #475569;
  background: #f1f5f9;
}

.modal-confirm {
  background: linear-gradient(135deg, #2563eb, #0ea5e9);
  color: #fff;
}
</style>
