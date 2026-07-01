<template>
  <view class="phone-page">
    <view class="phone-bg"></view>

    <view class="phone-navbar">
      <view class="phone-navbar-status"></view>
      <view class="phone-navbar-content">
        <view class="phone-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#0f172a" />
        </view>
        <text class="phone-navbar-title">{{ $t('setPhone') }}</text>
        <view class="phone-nav-space"></view>
      </view>
    </view>

    <view class="phone-container">
      <view class="phone-desc">{{ $t('setPhoneDesc') }}</view>

      <view class="phone-card">
        <view class="phone-field" v-if="currentPhone">
          <text class="phone-label">{{ $t('currentPhone') }}</text>
          <view class="phone-current-box">
            <text class="phone-current-text">{{ maskedCurrentPhone }}</text>
          </view>
        </view>

        <view class="phone-field">
          <text class="phone-label">{{ $t('newPhone') }}</text>
          <input class="phone-input" type="number" :placeholder="$t('enterNewPhone')" v-model="form.phone" maxlength="20" />
        </view>

        <view class="phone-field">
          <text class="phone-label">{{ $t('password') }}</text>
          <input class="phone-input" type="password" :placeholder="$t('enterPassword')" v-model="form.password" />
        </view>

        <view class="phone-field">
          <text class="phone-label">{{ $t('captcha') }}</text>
          <view class="phone-captcha-row">
            <input class="phone-input phone-captcha-input" :placeholder="$t('captchaPlaceholder')" v-model="form.captcha" />
            <image class="phone-captcha-img" @tap="getCaptchaFunc()" :src="captchaImg" mode="aspectFit"></image>
          </view>
        </view>
      </view>

      <button class="phone-submit" @tap="onSubmit" :disabled="submitting">
        {{ submitting ? $t('submitting') : $t('confirmSet') }}
      </button>
    </view>
  </view>
</template>

<script setup>
import { reactive, ref, computed } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { getCaptcha, setPhoneVerified } from '@/api/base.js'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)
const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')

const currentPhone = ref('')
const submitting = ref(false)
const lastLoadedLocale = ref('')

onLoad((options) => {
  if (options.phone) {
    currentPhone.value = options.phone
  }
  lastLoadedLocale.value = locale.value
})

onShow(() => {
  const localeChanged = !!lastLoadedLocale.value && lastLoadedLocale.value !== locale.value
  if (localeChanged) {
    getCaptchaFunc()
  }
  lastLoadedLocale.value = locale.value
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
  uni.navigateBack({ fail: () => uni.switchTab({ url: '/pages/learning/profile' }) })
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
.phone-page {
  min-height: 100vh;
  background: #f4f7fb;
}

.phone-bg {
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

.phone-navbar {
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

.phone-navbar-status {
  height: var(--status-bar-height, 44px);
}

.phone-navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 88rpx;
  padding: 0 24rpx;
}

.phone-back,
.phone-nav-space {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.phone-back {
  border-radius: 50%;
  background: #ffffff;
  box-shadow: 0 8rpx 20rpx rgba(15, 23, 42, 0.08);
}

.phone-navbar-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #0f172a;
}

.phone-container {
  position: relative;
  padding: calc(var(--status-bar-height, 44px) + 88rpx + 30rpx) 24rpx 30rpx;
}

.phone-desc {
  margin-bottom: 16rpx;
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.56);
}

.phone-card {
  border-radius: 18rpx;
  background: #ffffff;
  box-shadow: 0 12rpx 28rpx rgba(15, 23, 42, 0.08);
  padding: 26rpx 24rpx 12rpx;
}

.phone-field {
  margin-bottom: 22rpx;
}

.phone-label {
  display: block;
  margin-bottom: 10rpx;
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.58);
}

.phone-input {
  height: 82rpx;
  border-radius: 12rpx;
  background: #f8fafc;
  border: 1rpx solid #dbe3ee;
  padding: 0 22rpx;
  font-size: 28rpx;
  color: #0f172a;
}

.phone-current-box {
  height: 82rpx;
  border-radius: 12rpx;
  border: 1rpx solid #dbe3ee;
  background: #f8fafc;
  padding: 0 22rpx;
  display: flex;
  align-items: center;
}

.phone-current-text {
  font-size: 28rpx;
  color: #0f172a;
  letter-spacing: 2rpx;
}

.phone-captcha-row {
  display: flex;
  align-items: center;
  gap: 14rpx;
}

.phone-captcha-input {
  flex: 1;
}

.phone-captcha-img {
  width: 220rpx;
  height: 82rpx;
  border-radius: 12rpx;
  border: 1rpx solid #dbe3ee;
  background: #ffffff;
}

.phone-submit {
  margin-top: 24rpx;
  width: 100%;
  height: 88rpx;
  line-height: 88rpx;
  border-radius: 44rpx;
  border: none;
  color: #ffffff;
  font-size: 30rpx;
  font-weight: 600;
  background: linear-gradient(135deg, #2563eb, #0ea5e9);
}

.phone-submit[disabled] {
  opacity: 0.6;
}
</style>
