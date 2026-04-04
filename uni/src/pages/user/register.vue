<template>
  <view class="nf-page">
    <view class="nf-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <text class="nf-back-icon">&#xe603;</text>
        </view>
        <text class="nf-navbar-title">{{ $t('navRegister') }}</text>
        <view class="nf-lang-btn" @tap="showLangPicker = true">
          <text class="nf-lang-label">{{ langLabel }}</text>
        </view>
      </view>
    </view>

    <view class="nf-container">
      <!-- 标题 -->
      <view class="nf-header">
        <text class="nf-title">{{ $t('registerBtn') }}</text>
        <view class="nf-title-line"></view>
      </view>

      <!-- 模式切换（仅在两种注册方式都启用时显示） -->
      <view class="nf-mode-switch" v-if="phoneLoginEnabled && usernameLoginEnabled">
        <view class="nf-mode-tab" :class="{active: registerMode === 'username'}" @tap="registerMode = 'username'">
          <text>{{ $t('usernameRegisterBtn') }}</text>
        </view>
        <view class="nf-mode-tab" :class="{active: registerMode === 'phone'}" @tap="registerMode = 'phone'">
          <text>{{ $t('phoneRegisterBtn') }}</text>
        </view>
      </view>

      <!-- 表单卡片 -->
      <view class="nf-card">
        <!-- 用户名模式 -->
        <view class="nf-field" v-if="registerMode === 'username'">
          <text class="nf-label">{{ $t('account') }}</text>
          <input class="nf-input" :placeholder="$t('accountPlaceholder')" maxlength="12" v-model="form.username" />
        </view>

        <!-- 手机号模式 -->
        <view class="nf-field" v-if="registerMode === 'phone'">
          <text class="nf-label">{{ $t('phoneNumber') }}</text>
          <view class="nf-phone-row">
            <view class="nf-area-code-btn" @tap="showAreaCodePicker = true">
              <text class="nf-area-code-text">{{ selectedAreaCode || '+86' }}</text>
              <text class="nf-area-code-arrow">▼</text>
            </view>
            <input class="nf-input nf-phone-input" :placeholder="$t('phonePlaceholder')" maxlength="15" v-model="form.phone" type="number" />
          </view>
        </view>

        <view class="nf-field">
          <text class="nf-label">{{ $t('password') }}</text>
          <input class="nf-input" type="password" maxlength="18" :placeholder="$t('passwordPlaceholder')" v-model="form.password" />
        </view>
        <view class="nf-field">
          <text class="nf-label">{{ $t('repeatPassword') }}</text>
          <input class="nf-input" type="password" maxlength="18" :placeholder="$t('passwordPlaceholder')" v-model="form.rePassword" />
        </view>

        <!-- 验证码（所有模式都显示） -->
        <view class="nf-field nf-captcha-field">
          <text class="nf-label">{{ $t('captcha') }}</text>
          <view class="nf-captcha-row">
            <input class="nf-input nf-captcha-input" :placeholder="$t('captchaPlaceholder')" v-model="form.captcha" />
            <image class="nf-captcha-img" @tap="getCaptchaFunc()" :src="captchaImg" mode="aspectFit"></image>
          </view>
        </view>

        <!-- 邀请码 -->
        <view class="nf-field">
          <text class="nf-label">{{ $t('inviteCodeLabel') }}</text>
          <input class="nf-input" :placeholder="$t('inviteCodePlaceholder')" maxlength="20" v-model="form.inviteCode" />
        </view>
      </view>

      <!-- 按钮 -->
      <view class="nf-actions">
        <button class="nf-btn nf-btn-primary" @tap="registerFunc()">{{ $t('registerBtn') }}</button>
        <button class="nf-btn nf-btn-ghost" @tap="toLogin()">{{ $t('goToLogin') }}</button>
      </view>
    </view>

    <!-- 语言弹窗 -->
    <lang-switch v-model="showLangPicker" />

    <!-- 区号选择弹窗 -->
    <view class="nf-popup-mask" v-if="showAreaCodePicker" @tap="showAreaCodePicker = false">
      <view class="nf-popup-content" @tap.stop>
        <view class="nf-popup-title">{{ $t('selectAreaCode') }}</view>
        <scroll-view scroll-y class="nf-popup-scroll">
          <view class="nf-area-item" v-for="item in areaCodes" :key="item.ID" @tap="selectArea(item)">
            <text class="nf-area-name">{{ $lt(item.countryName) }}</text>
            <text class="nf-area-code-val">{{ item.areaCode }}</text>
          </view>
        </scroll-view>
      </view>
    </view>
  </view>
</template>
<script setup>
	import {
		onLoad,
	} from '@dcloudio/uni-app';

	import {
		reactive,
		ref,
		computed,
		onMounted
	} from 'vue';

	import {
		register,
		getCaptcha,
		phoneRegister
	} from "@/api/base.js"
	import { getEnabledPhoneAreaCodes } from "@/api/phoneAreaCode.js"
	import { getLoginConfig } from "@/api/sysConfig.js"

	import {useUserStore} from "@/pinia/modules/user.js"
	import { useLangStore } from '@/pinia/modules/lang.js'
	import langSwitch from '@/components/lang-switch/lang-switch.vue'

	const langStore = useLangStore()
	const $t = computed(() => langStore.$t)
	const $lt = computed(() => langStore.$lt)
	const langLabel = computed(() => {
	  const map = { zh: '中', en: 'EN', mn: 'MN', 'zh-TW': '繁', th: 'TH', hi: 'HI', id: 'ID' }
	  return map[langStore.locale] || langStore.locale.slice(0, 2).toUpperCase()
	})
	const showLangPicker = ref(false)

	const goBack = () => {
	  uni.navigateBack({ fail: () => uni.switchTab({ url: '/pages/tabBar/index' }) })
	}

	const userStore = useUserStore()
	const token = userStore.token || ''

	if(token){
		uni.switchTab({ url: '/pages/tabBar/index' })
	}

	// 注册模式
	const registerMode = ref('username') // 'username' | 'phone'
	const phoneLoginEnabled = ref(false)
	const usernameLoginEnabled = ref(true)
	const showAreaCodePicker = ref(false)
	const areaCodes = ref([])
	const selectedAreaCode = ref('+86')

	const form = reactive({
		username: "",
		password: "",
		rePassword: "",
		inviteCode: "",
		phone: "",
		areaCode: "+86",
		captcha: "",
		captchaId: ""
	})

	const captchaImg = ref('')

	const getCaptchaFunc = async () => {
		const res = await getCaptcha()
		if (res.code === 0) {
			captchaImg.value = res.data.picPath
			form.captchaId = res.data.captchaId
		}
	}

	// 加载配置
	const loadConfig = async () => {
		try {
			const res = await getLoginConfig()
			if (res.code === 0 && res.data) {
				phoneLoginEnabled.value = res.data.phone_login_enabled === 'true'
				usernameLoginEnabled.value = res.data.username_login_enabled !== 'false'
				// 互斥：两者不能同时关闭
				if (!phoneLoginEnabled.value && !usernameLoginEnabled.value) {
					usernameLoginEnabled.value = true
				}
				// 根据配置决定默认注册模式
				if (!usernameLoginEnabled.value && phoneLoginEnabled.value) {
					registerMode.value = 'phone'
				}
			}
		} catch(e) {}
	}

	// 加载区号列表
	const loadAreaCodes = async () => {
		try {
			const res = await getEnabledPhoneAreaCodes()
			if (res.code === 0 && res.data && res.data.list) {
				areaCodes.value = res.data.list
				if (areaCodes.value.length > 0) {
					selectedAreaCode.value = areaCodes.value[0].areaCode
					form.areaCode = areaCodes.value[0].areaCode
				}
			}
		} catch(e) {}
	}

	const selectArea = (item) => {
		selectedAreaCode.value = item.areaCode
		form.areaCode = item.areaCode
		showAreaCodePicker.value = false
	}

	// 从URL参数获取邀请码
	onLoad((options) => {
		if (options && options.inviteCode) {
			form.inviteCode = options.inviteCode
		}
	})

	onMounted(() => {
		getCaptchaFunc()
		loadConfig()
		loadAreaCodes()
	})

	const toLogin = () => {
		uni.navigateTo({ url: '/pages/user/login' })
	}

	const registerFunc = async () => {
		if (registerMode.value === 'username') {
			// 用户名注册
			if (!form.username) {
				uni.showToast({ title: $t.value('enterUsername'), icon: 'none' })
				return
			}
			if (!form.password) {
				uni.showToast({ title: $t.value('enterPassword'), icon: 'none' })
				return
			}
			if (!form.rePassword) {
				uni.showToast({ title: $t.value('enterRePassword'), icon: 'none' })
				return
			}
			if (form.password !== form.rePassword) {
				uni.showToast({ title: $t.value('passwordMismatch'), icon: 'none' })
				return
			}
			if (!form.captcha) {
				uni.showToast({ title: $t.value('enterCaptcha'), icon: 'none' })
				return
			}
			const res = await register({
				username: form.username,
				password: form.password,
				rePassword: form.rePassword,
				captcha: form.captcha,
				captchaId: form.captchaId,
				inviteCode: form.inviteCode
			})
			if (res.code === 0) {
				uni.showToast({ title: $t.value('registerSuccess'), icon: 'none' })
				toLogin()
			} else {
				getCaptchaFunc()
			}
		} else {
			// 手机号注册
			if (!form.phone) {
				uni.showToast({ title: $t.value('phonePlaceholder'), icon: 'none' })
				return
			}
			if (!form.password) {
				uni.showToast({ title: $t.value('enterPassword'), icon: 'none' })
				return
			}
			if (!form.rePassword) {
				uni.showToast({ title: $t.value('enterRePassword'), icon: 'none' })
				return
			}
			if (form.password !== form.rePassword) {
				uni.showToast({ title: $t.value('passwordMismatch'), icon: 'none' })
				return
			}
			if (!form.captcha) {
				uni.showToast({ title: $t.value('enterCaptcha'), icon: 'none' })
				return
			}
			const res = await phoneRegister({
				areaCode: form.areaCode,
				phone: form.phone,
				password: form.password,
				captcha: form.captcha,
				captchaId: form.captchaId,
				inviteCode: form.inviteCode
			})
			if (res.code === 0) {
				uni.showToast({ title: $t.value('registerSuccess'), icon: 'none' })
				toLogin()
			} else {
				getCaptchaFunc()
			}
		}
	}
</script>
<style lang="scss" scoped>
page { background-color: #000; }

.nf-page {
  min-height: 100vh;
  background: #000;
  position: relative;
  overflow: hidden;
}

.nf-bg {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  z-index: 0;
  background:
    radial-gradient(ellipse at 85% 25%, rgba(229, 9, 20, 0.25) 0%, transparent 50%),
    radial-gradient(ellipse at 15% 65%, rgba(229, 9, 20, 0.15) 0%, transparent 50%),
    radial-gradient(ellipse at 50% 90%, rgba(229, 9, 20, 0.1) 0%, transparent 40%),
    #000;
}

/* 自定义导航栏 */
.nf-navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(20px);
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
}

.nf-navbar-status {
  height: var(--status-bar-height, 44px);
}

.nf-navbar-content {
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24rpx;
}

.nf-navbar-back {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background 0.3s;

  &:active {
    background: rgba(255, 255, 255, 0.1);
  }
}

.nf-back-icon {
  font-size: 36rpx;
  color: #fff;
  font-family: 'iconfont';
  &::before { content: '←'; font-family: inherit; }
}

.nf-navbar-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #fff;
  letter-spacing: 2rpx;
}

.nf-lang-btn {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);
  border: 2rpx solid rgba(255, 255, 255, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  transition: all 0.3s;

  &:active {
    background: rgba(255, 255, 255, 0.15);
    transform: scale(0.92);
  }
}

.nf-lang-label {
  font-size: 22rpx;
  font-weight: 700;
  color: #fff;
  line-height: 1;
  text-align: center;
}

.nf-container {
  position: relative;
  z-index: 1;
  padding: 0 60rpx;
  padding-top: calc(var(--status-bar-height, 44px) + 88rpx + 40rpx);
}

.nf-header {
  margin-bottom: 64rpx;
}

.nf-title {
  font-size: 52rpx;
  font-weight: 800;
  color: #fff;
  letter-spacing: 4rpx;
}

.nf-title-line {
  width: 80rpx;
  height: 6rpx;
  border-radius: 3rpx;
  background: linear-gradient(90deg, #e50914, #ff6b6b);
  margin-top: 20rpx;
}

.nf-card {
  background: rgba(255, 255, 255, 0.05);
  border: 2rpx solid rgba(255, 255, 255, 0.08);
  border-radius: 28rpx;
  padding: 40rpx 36rpx;
  backdrop-filter: blur(20px);
  box-shadow: 0 8rpx 40rpx rgba(0, 0, 0, 0.4);
}

.nf-field {
  margin-bottom: 36rpx;

  &:last-child { margin-bottom: 0; }
}

.nf-label {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.5);
  font-weight: 600;
  letter-spacing: 2rpx;
  text-transform: uppercase;
  margin-bottom: 12rpx;
  display: block;
}

.nf-input {
  width: 100%;
  height: 88rpx;
  background: rgba(255, 255, 255, 0.06);
  border: 2rpx solid rgba(255, 255, 255, 0.1);
  border-radius: 16rpx;
  padding: 0 24rpx;
  color: #fff;
  font-size: 30rpx;
  box-sizing: border-box;
  transition: border-color 0.3s;

  &:focus {
    border-color: rgba(229, 9, 20, 0.6);
  }
}

.nf-actions {
  margin-top: 56rpx;
  display: flex;
  flex-direction: column;
  gap: 24rpx;
}

.nf-btn {
  width: 100%;
  height: 96rpx;
  border-radius: 48rpx;
  font-size: 32rpx;
  font-weight: 700;
  letter-spacing: 4rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  margin: 0;
  line-height: 1;
  overflow: hidden;
  white-space: nowrap;
  transition: all 0.3s;

  &:active { transform: scale(0.97); }
  &::after { border: none; }
}

.nf-btn-primary {
  background: linear-gradient(135deg, #e50914 0%, #ff4d4d 100%);
  color: #fff;
  box-shadow: 0 8rpx 32rpx rgba(229, 9, 20, 0.4);

  &:active {
    box-shadow: 0 4rpx 16rpx rgba(229, 9, 20, 0.5);
  }
}

.nf-btn-ghost {
  background: transparent;
  color: rgba(255, 255, 255, 0.6);
  border: 2rpx solid rgba(255, 255, 255, 0.12);

  &:active {
    background: rgba(255, 255, 255, 0.05);
    color: #fff;
  }
}

/* 模式切换 */
.nf-mode-switch {
  display: flex;
  margin-bottom: 32rpx;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 16rpx;
  padding: 6rpx;
}

.nf-mode-tab {
  flex: 1;
  text-align: center;
  padding: 16rpx 0;
  border-radius: 12rpx;
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.5);
  line-height: 1.2;
  overflow: hidden;
  white-space: nowrap;
  transition: all 0.3s;

  &.active {
    background: rgba(229, 9, 20, 0.3);
    color: #fff;
    font-weight: 600;
  }
}

/* 手机号输入行 */
.nf-phone-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.nf-area-code-btn {
  height: 88rpx;
  padding: 0 20rpx;
  background: rgba(255, 255, 255, 0.06);
  border: 2rpx solid rgba(255, 255, 255, 0.1);
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  gap: 8rpx;
  flex-shrink: 0;
}

.nf-area-code-text {
  color: #fff;
  font-size: 28rpx;
  font-weight: 600;
}

.nf-area-code-arrow {
  color: rgba(255, 255, 255, 0.4);
  font-size: 20rpx;
}

.nf-phone-input {
  flex: 1;
}

.nf-captcha-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.nf-captcha-input {
  flex: 1;
}

.nf-captcha-img {
  width: 200rpx;
  height: 88rpx;
  border-radius: 16rpx;
  background: rgba(255, 255, 255, 0.08);
  flex-shrink: 0;
}

/* 区号弹窗 */
.nf-popup-mask {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  z-index: 200;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.nf-popup-content {
  width: 100%;
  max-height: 60vh;
  background: #1a1a1a;
  border-radius: 28rpx 28rpx 0 0;
  padding: 32rpx 0;
}

.nf-popup-title {
  text-align: center;
  font-size: 32rpx;
  font-weight: 700;
  color: #fff;
  padding-bottom: 24rpx;
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.08);
}

.nf-popup-scroll {
  max-height: 50vh;
}

.nf-area-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 40rpx;
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.05);

  &:active {
    background: rgba(255, 255, 255, 0.05);
  }
}

.nf-area-name {
  color: #fff;
  font-size: 28rpx;
}

.nf-area-code-val {
  color: rgba(229, 9, 20, 0.8);
  font-size: 28rpx;
  font-weight: 600;
}
</style>
