<template>
  <view class="nf-profile">
    <view class="nf-profile-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff" />
        </view>
        <text class="nf-navbar-title">{{ $t('personalInfo') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <view class="nf-profile-content">
      <!-- 头像区域 -->
      <view class="nf-avatar-section">
        <image class="nf-avatar" :src="avatarUrl" mode="aspectFill" @tap="changeAvatar" />
        <text class="nf-avatar-tip">{{ $t('tapToChangeAvatar') }}</text>
      </view>

      <!-- 信息卡片 -->
      <view class="nf-info-card">
        <!-- 用户名 -->
        <view class="nf-info-item" @tap="editNickName">
          <text class="nf-info-label">{{ $t('nickname') }}</text>
          <view class="nf-info-value-row">
            <text class="nf-info-value">{{ userInfo.nickName || '-' }}</text>
            <uni-icons type="right" size="14" color="rgba(255,255,255,0.3)" />
          </view>
        </view>

        <!-- 手机号 -->
        <view class="nf-info-item">
          <text class="nf-info-label">{{ $t('phone') }}</text>
          <view class="nf-info-value-row">
            <text class="nf-info-value">{{ maskedPhone }}</text>
          </view>
        </view>

        <!-- 邮箱 -->
        <view class="nf-info-item">
          <text class="nf-info-label">{{ $t('email') }}</text>
          <view class="nf-info-value-row">
            <text class="nf-info-value">{{ userInfo.email || $t('notSet') }}</text>
          </view>
        </view>

        <!-- 注册时间 -->
        <view class="nf-info-item">
          <text class="nf-info-label">{{ $t('registerTime') }}</text>
          <view class="nf-info-value-row">
            <text class="nf-info-value">{{ formatTime(userInfo.CreatedAt) }}</text>
          </view>
        </view>
      </view>

      <!-- 操作区域 -->
      <view class="nf-action-card">
        <view class="nf-action-item" @tap="goChangePassword">
          <uni-icons type="locked" size="20" color="#e50914" />
          <text class="nf-action-text">{{ $t('changePassword') }}</text>
          <uni-icons type="right" size="14" color="rgba(255,255,255,0.3)" />
        </view>
      </view>

      <!-- 退出登录 -->
      <view class="nf-logout-btn" @tap="logout">
        <text class="nf-logout-text">{{ $t('logout') }}</text>
      </view>
    </view>

    <!-- 修改昵称弹窗 -->
    <view v-if="showNickNameModal" class="nf-modal-mask" @tap.stop="showNickNameModal = false">
      <view class="nf-modal-box" @tap.stop>
        <text class="nf-modal-title">{{ $t('editNickname') }}</text>
        <input
          class="nf-modal-input"
          v-model="newNickName"
          :placeholder="$t('enterNickname')"
          maxlength="20"
        />
        <view class="nf-modal-btns">
          <view class="nf-modal-btn nf-modal-cancel" @tap="showNickNameModal = false">
            <text>{{ $t('cancel') }}</text>
          </view>
          <view class="nf-modal-btn nf-modal-confirm" @tap="saveNickName">
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
const avatarUrl = ref('https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0')
const showNickNameModal = ref(false)
const newNickName = ref('')

const maskedPhone = computed(() => {
  const phone = userInfo.value.phone || ''
  if (phone.length >= 7) {
    return phone.substring(0, 3) + '****' + phone.substring(phone.length - 4)
  }
  return phone || '-'
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
      if (res.data.headerImg) {
        avatarUrl.value = res.data.headerImg
      }
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

const changeAvatar = () => {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: async (res) => {
      const tempPath = res.tempFilePaths[0]
      uni.uploadFile({
        url: uni.getStorageSync('baseUrl') + '/fileUploadAndDownload/upload',
        filePath: tempPath,
        name: 'file',
        header: { 'x-token': uni.getStorageSync('x-token') },
        success: async (uploadRes) => {
          const data = JSON.parse(uploadRes.data)
          if (data.code === 0 && data.data && data.data.file) {
            avatarUrl.value = data.data.file.url
            await setClientUserInfo({ headerImg: data.data.file.url })
            uni.showToast({ title: $t.value('updateSuccess'), icon: 'none' })
          }
        }
      })
    }
  })
}

const editNickName = () => {
  newNickName.value = userInfo.value.nickName || ''
  showNickNameModal.value = true
}

const saveNickName = async () => {
  if (!newNickName.value.trim()) {
    uni.showToast({ title: $t.value('enterNickname'), icon: 'none' })
    return
  }
  try {
    const res = await setClientUserInfo({ nickName: newNickName.value.trim() })
    if (res.code === 0) {
      userInfo.value.nickName = newNickName.value.trim()
      showNickNameModal.value = false
      uni.showToast({ title: $t.value('updateSuccess'), icon: 'none' })
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
  uni.navigateBack({ fail: () => uni.switchTab({ url: '/pages/tabBar/my/index' }) })
}
</script>

<style lang="scss" scoped>
.nf-profile {
  min-height: 100vh;
  background: #141414;
  position: relative;
}
.nf-profile-bg {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: linear-gradient(180deg, #1a1a2e 0%, #141414 100%);
  z-index: 0;
}
.nf-navbar {
  position: fixed;
  top: 0; left: 0; right: 0;
  z-index: 100;
  background: rgba(20, 20, 20, 0.95);
  backdrop-filter: blur(20rpx);
}
.nf-navbar-status {
  height: var(--status-bar-height, 44rpx);
}
.nf-navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 88rpx;
  padding: 0 24rpx;
}
.nf-navbar-back {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}
.nf-navbar-title {
  font-size: 34rpx;
  font-weight: bold;
  color: #fff;
}
.nf-profile-content {
  position: relative;
  z-index: 1;
  padding-top: calc(var(--status-bar-height, 44rpx) + 88rpx + 40rpx);
  padding: calc(var(--status-bar-height, 44rpx) + 88rpx + 40rpx) 24rpx 60rpx;
}
.nf-avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 40rpx;
}
.nf-avatar {
  width: 160rpx;
  height: 160rpx;
  border-radius: 50%;
  border: 4rpx solid rgba(229, 9, 20, 0.5);
}
.nf-avatar-tip {
  font-size: 22rpx;
  color: rgba(255,255,255,0.4);
  margin-top: 12rpx;
}
.nf-info-card {
  background: rgba(255,255,255,0.06);
  border-radius: 16rpx;
  padding: 0 30rpx;
  margin-bottom: 20rpx;
}
.nf-info-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 30rpx 0;
  border-bottom: 1rpx solid rgba(255,255,255,0.06);
}
.nf-info-item:last-child {
  border-bottom: none;
}
.nf-info-label {
  font-size: 28rpx;
  color: rgba(255,255,255,0.6);
}
.nf-info-value-row {
  display: flex;
  align-items: center;
  gap: 8rpx;
}
.nf-info-value {
  font-size: 28rpx;
  color: #fff;
}
.nf-action-card {
  background: rgba(255,255,255,0.06);
  border-radius: 16rpx;
  padding: 0 30rpx;
  margin-bottom: 40rpx;
}
.nf-action-item {
  display: flex;
  align-items: center;
  padding: 30rpx 0;
  gap: 16rpx;
}
.nf-action-text {
  flex: 1;
  font-size: 28rpx;
  color: #fff;
}
.nf-logout-btn {
  height: 88rpx;
  border-radius: 12rpx;
  background: rgba(229, 9, 20, 0.15);
  border: 1rpx solid rgba(229, 9, 20, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
}
.nf-logout-text {
  font-size: 30rpx;
  font-weight: bold;
  color: #e50914;
}
/* 弹窗 */
.nf-modal-mask {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.7);
  z-index: 200;
  display: flex;
  align-items: center;
  justify-content: center;
}
.nf-modal-box {
  width: 600rpx;
  background: #222;
  border-radius: 20rpx;
  padding: 40rpx;
}
.nf-modal-title {
  font-size: 32rpx;
  font-weight: bold;
  color: #fff;
  text-align: center;
  margin-bottom: 30rpx;
}
.nf-modal-input {
  height: 80rpx;
  background: rgba(255,255,255,0.08);
  border-radius: 12rpx;
  padding: 0 24rpx;
  color: #fff;
  font-size: 28rpx;
  margin-bottom: 30rpx;
}
.nf-modal-btns {
  display: flex;
  gap: 20rpx;
}
.nf-modal-btn {
  flex: 1;
  height: 80rpx;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28rpx;
  font-weight: bold;
}
.nf-modal-cancel {
  background: rgba(255,255,255,0.1);
  color: rgba(255,255,255,0.7);
}
.nf-modal-confirm {
  background: #e50914;
  color: #fff;
}
</style>
