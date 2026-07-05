<template>
  <view class="nf-page">
    <view class="nf-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#0f172a"></uni-icons>
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
        <view class="nf-brand-row">
          <image v-if="appLogoUrl" class="nf-logo" :src="appLogoUrl" mode="aspectFill"></image>
          <view v-else class="nf-logo nf-logo-fallback">{{ (appName || 'R').slice(0, 1).toUpperCase() }}</view>
          <view class="nf-brand-meta">
            <text class="nf-brand-name">{{ appName || 'RichRoad' }}</text>
            <text class="nf-brand-subtitle">{{ $t('wardrobeSlogan') }}</text>
          </view>
        </view>
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
          <input class="nf-input" :placeholder="usernamePlaceholder" maxlength="12" v-model="form.username" />
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
          <input class="nf-input" type="password" maxlength="18" :placeholder="passwordPlaceholder" v-model="form.password" />
        </view>
        <view class="nf-field">
          <text class="nf-label">{{ $t('repeatPassword') }}</text>
          <input class="nf-input" type="password" maxlength="18" :placeholder="passwordPlaceholder" v-model="form.rePassword" />
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
      <view class="nf-actions" :key="`actions-${langStore.locale}`">
        <button :key="`register-btn-${langStore.locale}`" class="nf-btn nf-btn-primary" @tap="registerFunc()">
          <text>{{ $t('registerBtn') }}</text>
        </button>
        <button :key="`login-btn-${langStore.locale}`" class="nf-btn nf-btn-ghost" @tap="toLogin()">
          <text>{{ $t('goToLogin') }}</text>
        </button>
      </view>
    </view>

    <!-- 语言弹窗 -->
    <lang-switch v-model="showLangPicker" />

    <!-- 区号选择弹窗 -->
    <view class="nf-popup-mask" v-if="showAreaCodePicker" @tap="showAreaCodePicker = false">
      <view class="nf-popup-content" @tap.stop>
        <view class="nf-popup-title">{{ $t('selectAreaCode') }}</view>
        <scroll-view scroll-y class="nf-popup-scroll">
          <view class="nf-area-item" v-for="item in areaCodeViews" :key="item._areaKey" @tap="selectArea(item._raw)">
            <text class="nf-area-name">{{ item._countryNameText }}</text>
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
    onShow,
	} from '@dcloudio/uni-app';

	import {
		reactive,
		ref,
		computed,
    onMounted,
    watch
	} from 'vue';

	import {
		register,
		getCaptcha,
		phoneRegister
	} from "@/api/base.js"
	import { getEnabledPhoneAreaCodes } from "@/api/phoneAreaCode.js"
	import { getLoginConfig } from "@/api/sysConfig.js"

	import { useLangStore } from '@/pinia/modules/lang.js'
  import { useAppConfigStore } from '@/pinia/modules/appConfig.js'
	import langSwitch from '@/components/lang-switch/lang-switch.vue'
  import { getExternalUrl } from '@/utils/url.js'

	const langStore = useLangStore()
  const appConfigStore = useAppConfigStore()
	const $t = computed(() => langStore.$t)
	const $lt = computed(() => langStore.$lt)
  const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')
  const appName = computed(() => appConfigStore.appName || 'RichRoad')
  const appLogoUrl = computed(() => getExternalUrl(appConfigStore.appLogo || ''))
  const langLabel = computed(() => {
    const map = { zh: 'ZH', en: 'EN', mn: 'MN', 'zh-TW': 'TW', th: 'TH', hi: 'HI', id: 'ID', vi: 'VI', ar: 'AR', ja: 'JA', ko: 'KO', ms: 'MS' }
    return map[langStore.locale] || langStore.locale.slice(0, 2).toUpperCase()
	})
	const showLangPicker = ref(false)
  const lastLoadedLocale = ref('')

  const tryAutoShowLangPicker = () => {
    if (showLangPicker.value) return
    if (!langStore.shouldAutoShowLanguagePicker()) return
    showLangPicker.value = true
  }

	const goBack = () => {
    uni.switchTab({ url: '/pages/learning/home' })
	}

	// 注册模式
	const registerMode = ref('username') // 'username' | 'phone'
	const phoneLoginEnabled = ref(false)
	const usernameLoginEnabled = ref(true)
	const showAreaCodePicker = ref(false)
	const areaCodes = ref([])
	const selectedAreaCode = ref('+86')
	const areaCodeViews = computed(() => areaCodes.value.map((item, index) => ({
    ...item,
    // 注册区号弹窗只读展示字段；注册提交仍使用 form.areaCode 和 selectArea 原对象。
    _raw: item,
    _areaKey: item.ID || item.id || `${item.areaCode || 'area'}-${index}`,
    _countryNameText: $lt.value(item.countryName) || item.countryName,
  })))

	// 正则和提示配置
	const usernameRegex = ref('')
	const passwordRegex = ref('')
	const usernameRegexTip = ref('')
	const passwordRegexTip = ref('')

	// 解析多语言JSON配置
	const parseLangTip = (jsonStr) => {
		if (!jsonStr) return ''
		try {
			const obj = JSON.parse(jsonStr)
      return obj[langStore.locale] || obj['en'] || obj['zh'] || ''
		} catch(e) {
			return jsonStr
		}
	}

	// 计算 placeholder
	const usernamePlaceholder = computed(() => {
		const tip = parseLangTip(usernameRegexTip.value)
		return tip || $t.value('accountPlaceholder')
	})
	const passwordPlaceholder = computed(() => {
		const tip = parseLangTip(passwordRegexTip.value)
		return tip || $t.value('passwordPlaceholder')
	})

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

  const reloadLocaleSensitiveData = async () => {
    await Promise.all([
      loadAreaCodes(),
      loadConfig(),
      appConfigStore.loadConfig({ force: true, localeOnly: true })
    ])
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
				// 正则和提示
				usernameRegex.value = res.data.username_regex || ''
				passwordRegex.value = res.data.password_regex || ''
				usernameRegexTip.value = res.data.username_regex_tip || ''
				passwordRegexTip.value = res.data.password_regex_tip || ''
			}
    } catch {
      // ignore config fetch failure
    }
	}

	// 加载区号列表
	const loadAreaCodes = async () => {
		try {
			const res = await getEnabledPhoneAreaCodes()
			if (res.code === 0 && res.data) {
				const list = Array.isArray(res.data) ? res.data : (res.data.list || [])
				areaCodes.value = list
        if (!areaCodes.value.length) return
        const currentCode = form.areaCode || selectedAreaCode.value
        const matched = areaCodes.value.find(item => item.areaCode === currentCode) || areaCodes.value[0]
        selectedAreaCode.value = matched.areaCode
        form.areaCode = matched.areaCode
        selectedAreaItem.value = matched
			}
    } catch {
      // ignore area code fetch failure
    }
	}

	const selectArea = (item) => {
		selectedAreaCode.value = item.areaCode
		form.areaCode = item.areaCode
		selectedAreaItem.value = item
		showAreaCodePicker.value = false
	}
	const selectedAreaItem = ref(null)

  const safeDecode = (value) => {
    if (value === undefined || value === null) return ''
    const raw = String(value).trim()
    if (!raw) return ''
    try {
      return decodeURIComponent(raw)
    } catch (e) {
      return raw
    }
  }

  const parseInviteCodeFromH5Location = () => {
    // #ifdef H5
    try {
      const fromSearch = new URLSearchParams(window.location.search || '').get('inviteCode')
      if (fromSearch) return safeDecode(fromSearch)
      const hash = window.location.hash || ''
      const hashQuery = hash.includes('?') ? hash.slice(hash.indexOf('?') + 1) : ''
      const fromHashQuery = new URLSearchParams(hashQuery).get('inviteCode')
      if (fromHashQuery) return safeDecode(fromHashQuery)
    } catch (e) {
      return ''
    }
    // #endif
    return ''
  }

  const applyInviteCode = (options = {}) => {
    const fromOptions = options?.inviteCode || options?.invite_code || options?.code || ''
    const fromStorage = safeDecode(uni.getStorageSync('pendingInviteCode') || '')
    const resolved = safeDecode(fromOptions) || parseInviteCodeFromH5Location() || fromStorage
    if (resolved) {
      form.inviteCode = resolved
      uni.setStorageSync('pendingInviteCode', resolved)
    }
  }

	// 从URL参数获取邀请码
	onLoad((options) => {
    applyInviteCode(options)
    lastLoadedLocale.value = locale.value
	})

  onShow(() => {
    const localeChanged = !!lastLoadedLocale.value && lastLoadedLocale.value !== locale.value
    if (localeChanged) {
      reloadLocaleSensitiveData()
    }
    lastLoadedLocale.value = locale.value
  })

  watch(() => locale.value, (newLocale, oldLocale) => {
    if (!oldLocale || newLocale === oldLocale) return
    reloadLocaleSensitiveData()
    lastLoadedLocale.value = newLocale
  })

	onMounted(() => {
		getCaptchaFunc()
		loadConfig()
		loadAreaCodes()
    appConfigStore.loadConfig()
    if (!form.inviteCode) {
      applyInviteCode()
    }
    tryAutoShowLangPicker()
	})

  const syncTabBarLocale = () => {
    const lang = langStore.locale || uni.getStorageSync('app-lang') || 'mn'
    langStore.updateTabBar(lang)
    ;[120, 300, 620].forEach((delay) => {
      setTimeout(() => {
        langStore.updateTabBar(lang)
      }, delay)
    })
  }

	const toLogin = () => {
    syncTabBarLocale()
		uni.navigateTo({ url: '/pages/user/login' })
	}

	const registerFunc = async () => {
		if (registerMode.value === 'username') {
			// 用户名注册
			if (!form.username) {
				uni.showToast({ title: $t.value('enterUsername'), icon: 'none' })
				return
			}
			// 用户名正则验证
			if (usernameRegex.value) {
				try {
					const regex = new RegExp(usernameRegex.value)
					if (!regex.test(form.username)) {
						const tip = parseLangTip(usernameRegexTip.value)
						uni.showToast({ title: tip || $t.value('enterUsername'), icon: 'none' })
						return
					}
        } catch {
          // ignore invalid username regex
        }
			}
			if (!form.password) {
				uni.showToast({ title: $t.value('enterPassword'), icon: 'none' })
				return
			}
			// 密码正则验证
			if (passwordRegex.value) {
				try {
					const regex = new RegExp(passwordRegex.value)
					if (!regex.test(form.password)) {
						const tip = parseLangTip(passwordRegexTip.value)
						uni.showToast({ title: tip || $t.value('enterPassword'), icon: 'none' })
						return
					}
        } catch {
          // ignore invalid password regex
        }
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
        uni.removeStorageSync('pendingInviteCode')
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
			// Phone regex validation
			if (selectedAreaItem.value && selectedAreaItem.value.phoneRegex) {
				try {
					const regex = new RegExp(selectedAreaItem.value.phoneRegex)
					if (!regex.test(form.phone)) {
						uni.showToast({ title: $t.value('phoneFormatInvalid'), icon: 'none' })
						return
					}
        } catch {
          // ignore invalid phone regex
        }
			}
			if (!form.password) {
				uni.showToast({ title: $t.value('enterPassword'), icon: 'none' })
				return
			}
			// 密码正则验证（手机号注册）
			if (passwordRegex.value) {
				try {
					const regex = new RegExp(passwordRegex.value)
					if (!regex.test(form.password)) {
						const tip = parseLangTip(passwordRegexTip.value)
						uni.showToast({ title: tip || $t.value('enterPassword'), icon: 'none' })
						return
					}
        } catch {
          // ignore invalid password regex
        }
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
        uni.removeStorageSync('pendingInviteCode')
				uni.showToast({ title: $t.value('registerSuccess'), icon: 'none' })
				toLogin()
			} else {
				getCaptchaFunc()
			}
		}
	}
</script>
<style lang="scss" scoped>
page { background-color: #f8f4e7; }

.nf-page {
  min-height: 100vh;
  background: #f8f4e7;
  position: relative;
}

.nf-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 0;
  background:
    radial-gradient(circle at 8% 15%, rgba(20, 184, 166, 0.26) 0%, transparent 46%),
    radial-gradient(circle at 92% 22%, rgba(249, 115, 22, 0.24) 0%, transparent 48%),
    radial-gradient(circle at 54% 96%, rgba(245, 158, 11, 0.22) 0%, transparent 40%),
    linear-gradient(180deg, #fbf7ee 0%, #f8f4e7 100%);
}

.nf-navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  background: rgba(255, 251, 242, 0.88);
  backdrop-filter: blur(16px);
  border-bottom: 1rpx solid rgba(146, 64, 14, 0.16);
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

.nf-navbar-back,
.nf-lang-btn {
  width: 64rpx;
  height: 64rpx;
  border-radius: 20rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2rpx solid rgba(146, 64, 14, 0.14);
  background: rgba(255, 255, 255, 0.76);
  transition: transform 0.2s ease;

  &:active {
    transform: scale(0.93);
  }
}

.nf-navbar-title {
  font-size: 34rpx;
  font-weight: 800;
  color: #7c2d12;
  letter-spacing: 2rpx;
}

.nf-lang-label {
  font-size: 22rpx;
  font-weight: 800;
  color: #7c2d12;
}

.nf-container {
  position: relative;
  z-index: 1;
  padding: 0 34rpx 52rpx;
  padding-top: calc(var(--status-bar-height, 44px) + 88rpx + 34rpx);
}

.nf-header {
  margin-bottom: 30rpx;
  padding: 24rpx;
  border-radius: 30rpx;
  border: 2rpx solid rgba(20, 184, 166, 0.16);
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.86) 0%, rgba(255, 248, 236, 0.92) 100%);
}

.nf-brand-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-bottom: 18rpx;
}

.nf-logo {
  width: 86rpx;
  height: 86rpx;
  border-radius: 26rpx;
  background: #fff;
  border: 2rpx solid rgba(20, 184, 166, 0.3);
  box-shadow: 0 10rpx 24rpx rgba(146, 64, 14, 0.12);
}

.nf-logo-fallback {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ffffff;
  font-size: 36rpx;
  font-weight: 700;
  background: linear-gradient(135deg, #0f766e 0%, #f97316 100%);
}

.nf-brand-meta {
  display: flex;
  flex-direction: column;
}

.nf-brand-name {
  font-size: 30rpx;
  font-weight: 700;
  color: #7c2d12;
}

.nf-brand-subtitle {
  font-size: 22rpx;
  color: rgba(120, 53, 15, 0.64);
}

.nf-title {
  font-size: 62rpx;
  font-weight: 800;
  color: #0f766e;
  letter-spacing: 3rpx;
}

.nf-title-line {
  width: 112rpx;
  height: 8rpx;
  border-radius: 999rpx;
  background: linear-gradient(90deg, #0f766e 0%, #f97316 100%);
  margin-top: 16rpx;
}

.nf-mode-switch,
.nf-card {
  border-radius: 28rpx;
  border: 2rpx solid rgba(20, 184, 166, 0.2);
  background: rgba(255, 253, 248, 0.94);
  box-shadow: 0 18rpx 38rpx rgba(120, 53, 15, 0.12);
}

.nf-mode-switch {
  display: flex;
  margin-bottom: 18rpx;
  padding: 8rpx;
}

.nf-mode-tab {
  flex: 1;
  text-align: center;
  padding: 16rpx 0;
  border-radius: 16rpx;
  font-size: 26rpx;
  color: rgba(120, 53, 15, 0.68);
  font-weight: 600;

  &.active {
    background: rgba(15, 118, 110, 0.16);
    color: #0f766e;
  }
}

.nf-card {
  padding: 30rpx;
  border-radius: 30rpx;
}

.nf-field {
  margin-bottom: 28rpx;
  &:last-child { margin-bottom: 0; }
}

.nf-label {
  display: block;
  font-size: 23rpx;
  color: rgba(120, 53, 15, 0.68);
  font-weight: 700;
  margin-bottom: 10rpx;
}

.nf-input,
.nf-area-code-btn,
.nf-captcha-img {
  height: 88rpx;
  border-radius: 16rpx;
  border: 2rpx solid rgba(146, 64, 14, 0.14);
  background: #fff;
  box-sizing: border-box;
}

.nf-input {
  width: 100%;
  padding: 0 22rpx;
  color: #334155;
  font-size: 30rpx;
}

.nf-phone-row,
.nf-captcha-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.nf-area-code-btn {
  padding: 0 18rpx;
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.nf-area-code-text { color: #7c2d12; font-size: 27rpx; font-weight: 700; }
.nf-area-code-arrow { color: rgba(120, 53, 15, 0.56); font-size: 20rpx; }
.nf-phone-input, .nf-captcha-input { flex: 1; }
.nf-captcha-img { width: 204rpx; flex-shrink: 0; }

.nf-actions {
  margin-top: 26rpx;
  display: flex;
  flex-direction: column;
  gap: 18rpx;
}

.nf-btn {
  width: 100%;
  height: 96rpx;
  border-radius: 20rpx;
  font-size: 30rpx;
  font-weight: 700;
  border: none;
  margin: 0;
  display: flex;
  align-items: center;
  justify-content: center;

  &:active { transform: translateY(2rpx); }
  &::after { border: none; }
}

.nf-btn-primary {
  color: #fff;
  background: linear-gradient(120deg, #0f766e 0%, #f97316 100%);
  box-shadow: 0 12rpx 24rpx rgba(15, 118, 110, 0.28);
}

.nf-btn-ghost {
  color: #7c2d12;
  background: rgba(255, 255, 255, 0.8);
  border: 2rpx solid rgba(146, 64, 14, 0.14);
}

.nf-popup-mask {
  position: fixed;
  inset: 0;
  z-index: 200;
  background: rgba(15, 23, 42, 0.42);
  display: flex;
  align-items: flex-end;
}

.nf-popup-content {
  width: 100%;
  max-height: 60vh;
  background: #fffdf8;
  border-radius: 30rpx 30rpx 0 0;
  padding: 30rpx 0;
}

.nf-popup-title {
  text-align: center;
  font-size: 31rpx;
  font-weight: 700;
  color: #7c2d12;
  padding-bottom: 20rpx;
  border-bottom: 1rpx solid rgba(146, 64, 14, 0.12);
}

.nf-popup-scroll { max-height: 50vh; }

.nf-area-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 26rpx 40rpx;
  border-bottom: 1rpx solid rgba(146, 64, 14, 0.08);

  &:active { background: rgba(15, 118, 110, 0.1); }
}

.nf-area-name { color: #334155; font-size: 28rpx; }
.nf-area-code-val { color: #0f766e; font-size: 28rpx; font-weight: 700; }
</style>
