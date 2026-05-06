<template>
  <view class="pwd-page">
    <view class="pwd-bg"></view>

    <view class="pwd-navbar">
      <view class="pwd-navbar-status"></view>
      <view class="pwd-navbar-content">
        <view class="pwd-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#0f172a"></uni-icons>
        </view>
        <text class="pwd-navbar-title">{{ $t('changePassword') }}</text>
        <view class="pwd-nav-space"></view>
      </view>
    </view>

    <view class="pwd-container">
      <view class="pwd-card">
        <view class="pwd-field">
          <text class="pwd-label">{{ $t('oldPassword') }}</text>
          <input class="pwd-input" type="password" :placeholder="$t('oldPasswordPlaceholder')" v-model="form.oldPassword" />
        </view>
        <view class="pwd-field">
          <text class="pwd-label">{{ $t('newPassword') }}</text>
          <input class="pwd-input" type="password" :placeholder="$t('newPasswordPlaceholder')" v-model="form.newPassword" />
        </view>
        <view class="pwd-field">
          <text class="pwd-label">{{ $t('confirmNewPassword') }}</text>
          <input class="pwd-input" type="password" :placeholder="$t('confirmNewPasswordPlaceholder')" v-model="form.confirmPassword" />
        </view>
        <view class="pwd-field">
          <text class="pwd-label">{{ $t('captcha') }}</text>
          <view class="pwd-captcha-row">
            <input class="pwd-input pwd-captcha-input" :placeholder="$t('captchaPlaceholder')" v-model="form.captcha" />
            <image class="pwd-captcha-img" @tap="getCaptchaFunc()" :src="captchaImg" mode="aspectFit"></image>
          </view>
        </view>
      </view>

      <view class="pwd-actions">
        <button class="pwd-btn pwd-btn-primary" @tap="onSubmit">{{ $t('confirmChange') }}</button>
        <button class="pwd-btn pwd-btn-secondary" @tap="goKefu">{{ $t('contactCustomerService') }}</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { reactive, ref, computed } from 'vue'
import { getCaptcha, changePassword } from '@/api/base.js'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const goBack = () => {
  uni.navigateBack({ fail: () => uni.switchTab({ url: '/pages/tabBar/my/index' }) })
}

const goKefu = () => {
  uni.navigateTo({ url: '/pages/kefu/index' })
}

// Backend error message i18n mapping
const errorMsgMap = {
  '旧密码错误': 'oldPasswordWrong',
  '用户不存在': 'userNotFound',
  '验证码请求过于频繁，请稍后再试': 'captchaRateLimit',
  '密码错误': 'passwordWrong',
}

const form = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
  captcha: '',
  captchaId: '',
})

const captchaImg = ref('')

const getCaptchaFunc = async () => {
  const res = await getCaptcha()
  if (res.code === 0) {
    captchaImg.value = res.data.picPath
    form.captchaId = res.data.captchaId
  }
}
getCaptchaFunc()

const onSubmit = async () => {
  if (!form.oldPassword) {
    uni.showToast({ title: $t.value('oldPasswordPlaceholder'), icon: 'none' })
    return
  }
  if (!form.newPassword) {
    uni.showToast({ title: $t.value('newPasswordPlaceholder'), icon: 'none' })
    return
  }
  if (form.newPassword !== form.confirmPassword) {
    uni.showToast({ title: $t.value('passwordMismatch'), icon: 'none' })
    return
  }
  if (!form.captcha) {
    uni.showToast({ title: $t.value('enterCaptcha'), icon: 'none' })
    return
  }

  const res = await changePassword({
    oldPassword: form.oldPassword,
    newPassword: form.newPassword,
    method: 'old_password',
    captcha: form.captcha,
    captchaId: form.captchaId,
  })
  if (res.code === 0) {
    uni.showToast({ title: $t.value('changeSuccess'), icon: 'success' })
    setTimeout(() => goBack(), 1500)
  } else {
    const key = errorMsgMap[res.msg]
    if (key) {
      uni.showToast({ title: $t.value(key), icon: 'none' })
    }
    getCaptchaFunc()
  }
}
</script>

<style scoped>
.pwd-page {
  min-height: 100vh;
  background: #f4f7fb;
}

.pwd-bg {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 360rpx;
  pointer-events: none;
  background:
    radial-gradient(120% 80% at 100% -10%, rgba(14, 165, 233, 0.2) 0%, transparent 60%),
    radial-gradient(120% 80% at 0% 0%, rgba(37, 99, 235, 0.2) 0%, transparent 60%);
}

.pwd-navbar {
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

.pwd-navbar-status {
  height: var(--status-bar-height, 44px);
}

.pwd-navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 88rpx;
  padding: 0 24rpx;
}

.pwd-back,
.pwd-nav-space {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pwd-back {
  border-radius: 50%;
  background: #ffffff;
  box-shadow: 0 8rpx 20rpx rgba(15, 23, 42, 0.08);
}

.pwd-navbar-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #0f172a;
}

.pwd-container {
  position: relative;
  padding: calc(var(--status-bar-height, 44px) + 88rpx + 30rpx) 24rpx 30rpx;
}

.pwd-card {
  border-radius: 18rpx;
  background: #ffffff;
  box-shadow: 0 12rpx 28rpx rgba(15, 23, 42, 0.08);
  padding: 26rpx 24rpx 12rpx;
}

.pwd-field {
  margin-bottom: 22rpx;
}

.pwd-label {
  display: block;
  margin-bottom: 10rpx;
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.58);
}

.pwd-input {
  height: 82rpx;
  border-radius: 12rpx;
  background: #f8fafc;
  border: 1rpx solid #dbe3ee;
  padding: 0 22rpx;
  font-size: 28rpx;
  color: #0f172a;
}

.pwd-captcha-row {
  display: flex;
  align-items: center;
  gap: 14rpx;
}

.pwd-captcha-input {
  flex: 1;
}

.pwd-captcha-img {
  width: 220rpx;
  height: 82rpx;
  border-radius: 12rpx;
  border: 1rpx solid #dbe3ee;
  background: #ffffff;
}

.pwd-actions {
  margin-top: 26rpx;
}

.pwd-btn {
  width: 100%;
  height: 88rpx;
  line-height: 88rpx;
  border-radius: 44rpx;
  font-size: 30rpx;
  font-weight: 600;
  text-align: center;
  margin-bottom: 18rpx;
  border: none;
}

.pwd-btn-primary {
  color: #ffffff;
  background: linear-gradient(135deg, #2563eb, #0ea5e9);
}

.pwd-btn-secondary {
  color: #334155;
  background: #ffffff;
  border: 1rpx solid #dbe3ee;
}
</style>
