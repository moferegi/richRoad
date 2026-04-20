<template>
  <view class="nf-page">
    <view class="nf-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff"></uni-icons>
        </view>
        <text class="nf-navbar-title">{{ $t('changePassword') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <view class="nf-container">
      <view class="nf-header">
        <text class="nf-title">{{ $t('changePassword') }}</text>
        <view class="nf-title-line"></view>
      </view>

      <view class="nf-card">
        <view class="nf-field">
          <text class="nf-label">{{ $t('oldPassword') }}</text>
          <input class="nf-input" type="password" :placeholder="$t('oldPasswordPlaceholder')" v-model="form.oldPassword" />
        </view>
        <view class="nf-field">
          <text class="nf-label">{{ $t('newPassword') }}</text>
          <input class="nf-input" type="password" :placeholder="$t('newPasswordPlaceholder')" v-model="form.newPassword" />
        </view>
        <view class="nf-field">
          <text class="nf-label">{{ $t('confirmNewPassword') }}</text>
          <input class="nf-input" type="password" :placeholder="$t('confirmNewPasswordPlaceholder')" v-model="form.confirmPassword" />
        </view>
        <view class="nf-field nf-captcha-field">
          <text class="nf-label">{{ $t('captcha') }}</text>
          <view class="nf-captcha-row">
            <input class="nf-input nf-captcha-input" :placeholder="$t('captchaPlaceholder')" v-model="form.captcha" />
            <image class="nf-captcha-img" @tap="getCaptchaFunc()" :src="captchaImg" mode="aspectFit"></image>
          </view>
        </view>
      </view>

      <view class="nf-actions">
        <button class="nf-btn nf-btn-primary" @tap="onSubmit">{{ $t('confirmChange') }}</button>
        <button class="nf-btn nf-btn-kefu" @tap="goKefu">{{ $t('contactCustomerService') }}</button>
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
.nf-page { min-height: 100vh; background: #141414; position: relative; }
.nf-bg { position: absolute; top: 0; left: 0; right: 0; height: 500rpx; background: linear-gradient(180deg, rgba(229,9,20,0.15), transparent); }
.nf-navbar { position: relative; z-index: 10; }
.nf-navbar-status { height: var(--status-bar-height, 44px); }
.nf-navbar-content { display: flex; align-items: center; height: 88rpx; padding: 0 24rpx; }
.nf-navbar-back { width: 64rpx; height: 64rpx; border-radius: 50%; background: rgba(255,255,255,0.06); border: 1rpx solid rgba(255,255,255,0.1); display: flex; align-items: center; justify-content: center; }
.nf-navbar-back:active { background: rgba(255,255,255,0.12); transform: scale(0.93); }
.nf-navbar-title { flex: 1; text-align: center; font-size: 32rpx; font-weight: 600; color: #fff; }
.nf-container { position: relative; z-index: 5; padding: 40rpx; }
.nf-header { margin-bottom: 60rpx; text-align: center; }
.nf-title { font-size: 48rpx; font-weight: bold; color: #fff; }
.nf-title-line { width: 80rpx; height: 6rpx; background: #e50914; margin: 16rpx auto 0; border-radius: 3rpx; }
.nf-card { background: rgba(255,255,255,0.05); border-radius: 20rpx; padding: 40rpx 32rpx; margin-bottom: 40rpx; }
.nf-field { margin-bottom: 32rpx; }
.nf-label { font-size: 24rpx; color: rgba(255,255,255,0.5); margin-bottom: 12rpx; display: block; }
.nf-input { height: 80rpx; background: rgba(255,255,255,0.08); border-radius: 12rpx; padding: 0 24rpx; font-size: 28rpx; color: #fff; }
.nf-captcha-field {}
.nf-captcha-row { display: flex; align-items: center; gap: 16rpx; }
.nf-captcha-input { flex: 1; }
.nf-captcha-img { width: 200rpx; height: 80rpx; border-radius: 12rpx; }
.nf-actions { margin-top: 20rpx; }
.nf-btn { width: 100%; height: 88rpx; line-height: 88rpx; border-radius: 44rpx; font-size: 30rpx; font-weight: 600; text-align: center; margin-bottom: 24rpx; border: none; }
.nf-btn-primary { background: linear-gradient(135deg, #e50914, #b20710); color: #fff; }
.nf-btn-kefu { background: rgba(255,255,255,0.08); color: rgba(255,255,255,0.7); border: 1rpx solid rgba(255,255,255,0.1); }
.nf-btn-kefu:active { background: rgba(255,255,255,0.12); }
</style>
