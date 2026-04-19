<template>
  <view class="nf-page">
    <view class="nf-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff" />
        </view>
        <text class="nf-navbar-title">{{ $t('setPhone') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <view class="nf-container">
      <view class="nf-header">
        <!-- <text class="nf-title">{{ $t('setPhone') }}</text>
        <view class="nf-title-line"></view> -->
        <text class="nf-subtitle">{{ $t('setPhoneDesc') }}</text>
      </view>

      <view class="nf-card">
        <!-- 当前手机号 -->
        <view class="nf-field" v-if="currentPhone">
          <text class="nf-label">{{ $t('currentPhone') }}</text>
          <view class="nf-current-phone">
            <text class="nf-current-phone-text">{{ maskedCurrentPhone }}</text>
          </view>
        </view>

        <!-- 新手机号 -->
        <view class="nf-field">
          <text class="nf-label">{{ $t('newPhone') }}</text>
          <input class="nf-input" type="number" :placeholder="$t('enterNewPhone')" v-model="form.phone" maxlength="20" />
        </view>

        <!-- 密码 -->
        <view class="nf-field">
          <text class="nf-label">{{ $t('password') }}</text>
          <input class="nf-input" type="password" :placeholder="$t('enterPassword')" v-model="form.password" />
        </view>

        <!-- 验证码 -->
        <view class="nf-field nf-captcha-field">
          <text class="nf-label">{{ $t('captcha') }}</text>
          <view class="nf-captcha-row">
            <input class="nf-input nf-captcha-input" :placeholder="$t('captchaPlaceholder')" v-model="form.captcha" />
            <image class="nf-captcha-img" @tap="getCaptchaFunc()" :src="captchaImg" mode="aspectFit"></image>
          </view>
        </view>
      </view>

      <view class="nf-actions">
        <button class="nf-btn nf-btn-primary" @tap="onSubmit" :disabled="submitting">
          {{ submitting ? $t('submitting') : $t('confirmSet') }}
        </button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { reactive, ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getCaptcha, setPhoneVerified } from '@/api/base.js'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const currentPhone = ref('')
const submitting = ref(false)

onLoad((options) => {
  if (options.phone) {
    currentPhone.value = options.phone
  }
})

const maskedCurrentPhone = computed(() => {
  const p = currentPhone.value
  if (!p || p.length < 4) return p
  const visibleStart = Math.ceil(p.length * 0.3)
  const visibleEnd = Math.ceil(p.length * 0.3)
  const maskedLen = p.length - visibleStart - visibleEnd
  return p.slice(0, visibleStart) + '*'.repeat(Math.max(maskedLen, 2)) + p.slice(p.length - visibleEnd)
})

const goBack = () => {
  uni.navigateBack({ fail: () => uni.switchTab({ url: '/pages/tabBar/my/index' }) })
}

const form = reactive({
  phone: '',
  password: '',
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
  if (!form.phone) {
    uni.showToast({ title: $t.value('enterNewPhone'), icon: 'none' })
    return
  }
  if (!form.password) {
    uni.showToast({ title: $t.value('enterPassword'), icon: 'none' })
    return
  }
  if (!form.captcha) {
    uni.showToast({ title: $t.value('enterCaptcha'), icon: 'none' })
    return
  }

  submitting.value = true
  try {
    const res = await setPhoneVerified({
      phone: form.phone,
      password: form.password,
      captcha: form.captcha,
      captchaId: form.captchaId,
    })
    if (res.code === 0) {
      // 更新本地缓存
      const userInfo = uni.getStorageSync('userInfo') || {}
      userInfo.phone = form.phone
      uni.setStorageSync('userInfo', userInfo)
      uni.showToast({ title: $t.value('setSuccess'), icon: 'success' })
      setTimeout(() => goBack(), 1500)
    } else {
      getCaptchaFunc()
    }
  } catch (e) {
    getCaptchaFunc()
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.nf-page { min-height: 100vh; background: #141414; position: relative; }
.nf-bg { position: absolute; top: 0; left: 0; right: 0; height: 500rpx; background: linear-gradient(180deg, rgba(229,9,20,0.15), transparent); }
.nf-navbar { position: relative; z-index: 10; }
.nf-navbar-status { height: var(--status-bar-height, 44px); }
.nf-navbar-content { display: flex; align-items: center; height: 88rpx; padding: 0 24rpx; }
.nf-navbar-back { width: 60rpx; height: 60rpx; display: flex; align-items: center; justify-content: center; }
.nf-navbar-title { flex: 1; text-align: center; font-size: 32rpx; font-weight: 600; color: #fff; }
.nf-container { position: relative; z-index: 5; padding: 40rpx; }
.nf-header { margin-bottom: 48rpx; text-align: center; }
.nf-title { font-size: 48rpx; font-weight: bold; color: #fff; }
.nf-title-line { width: 80rpx; height: 6rpx; background: #e50914; margin: 16rpx auto 0; border-radius: 3rpx; }
.nf-subtitle { display: block; font-size: 24rpx; color: rgba(255,255,255,0.5); margin-top: 16rpx; }
.nf-card { background: rgba(255,255,255,0.05); border-radius: 20rpx; padding: 40rpx 32rpx; margin-bottom: 40rpx; }
.nf-field { margin-bottom: 32rpx; }
.nf-label { font-size: 24rpx; color: rgba(255,255,255,0.5); margin-bottom: 12rpx; display: block; }
.nf-input { height: 80rpx; background: rgba(255,255,255,0.08); border-radius: 12rpx; padding: 0 24rpx; font-size: 28rpx; color: #fff; }
.nf-current-phone { height: 80rpx; background: rgba(255,255,255,0.04); border-radius: 12rpx; padding: 0 24rpx; display: flex; align-items: center; }
.nf-current-phone-text { font-size: 30rpx; color: rgba(255,255,255,0.7); letter-spacing: 4rpx; }
.nf-captcha-row { display: flex; align-items: center; gap: 16rpx; }
.nf-captcha-input { flex: 1; }
.nf-captcha-img { width: 200rpx; height: 80rpx; border-radius: 12rpx; }
.nf-actions { margin-top: 20rpx; }
.nf-btn { width: 100%; height: 88rpx; line-height: 88rpx; border-radius: 44rpx; font-size: 30rpx; font-weight: 600; text-align: center; margin-bottom: 24rpx; border: none; }
.nf-btn-primary { background: linear-gradient(135deg, #e50914, #b20710); color: #fff; }
.nf-btn-primary[disabled] { opacity: 0.6; }
</style>
