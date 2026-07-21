<template>
  <view class="nf-page">
    <!-- 紫色渐变 Hero 区 -->
    <view class="login-hero">
      <view class="hero-decor-circle decor-1"></view>
      <view class="hero-decor-circle decor-2"></view>
      <view class="hero-top-row">
        <view class="hero-spacer"></view>
        <view class="hero-lang-btn" @tap="showLangPicker = true">
          <text class="hero-lang-text">{{ langLabel }}</text>
        </view>
      </view>
      <view class="hero-brand">
        <view class="hero-logo-wrap">
          <image v-if="appLogoUrl" class="hero-logo" :src="appLogoUrl" mode="aspectFit" />
          <view v-else class="hero-logo hero-logo-fallback">{{ (appName || 'R').slice(0, 1).toUpperCase() }}</view>
        </view>
        <text class="hero-brand-name">{{ appName || 'RichRoad' }}</text>
        <text class="hero-brand-subtitle">{{ $t('wardrobeSlogan') }}</text>
      </view>
      <view class="hero-title-section">
        <text class="hero-title">{{ $t('loginBtn') }}</text>
        <view class="hero-title-line"></view>
      </view>
    </view>

    <!-- 内容区 -->
    <view class="login-body">
      <!-- 模式切换胶囊 -->
      <view class="mode-pills" v-if="phoneLoginEnabled && usernameLoginEnabled">
        <view class="mode-pill" :class="{active: loginMode === 'username'}" @tap="loginMode = 'username'">
          <text class="mode-pill-text">{{ $t('switchToAccountLogin') }}</text>
        </view>
        <view class="mode-pill" :class="{active: loginMode === 'phone'}" @tap="loginMode = 'phone'">
          <text class="mode-pill-text">{{ $t('switchToPhoneLogin') }}</text>
        </view>
      </view>

      <!-- 表单卡片 -->
      <view class="form-card">
        <!-- 用户名模式 -->
        <view class="form-field" v-if="loginMode === 'username'">
          <view class="field-icon-wrap">
            <text class="field-icon">👤</text>
          </view>
          <view class="field-content">
            <text class="field-label">{{ $t('account') }}</text>
            <input class="field-input" name="username" :placeholder="usernamePlaceholder" maxlength="30" v-model="form.username" />
          </view>
        </view>

        <!-- 手机号模式 -->
        <view class="form-field" v-if="loginMode === 'phone'">
          <view class="field-icon-wrap">
            <text class="field-icon">📱</text>
          </view>
          <view class="field-content">
            <text class="field-label">{{ $t('phoneNumber') }}</text>
            <view class="phone-row">
              <view class="area-code-btn" @tap="showAreaCodePicker = true">
                <text class="area-code-text">{{ selectedAreaCode || '+86' }}</text>
                <text class="area-code-arrow">▼</text>
              </view>
              <input class="field-input phone-input" name="phone" :placeholder="$t('phonePlaceholder')" maxlength="15" v-model="form.phone" type="number" />
            </view>
          </view>
        </view>

        <view class="form-field">
          <view class="field-icon-wrap">
            <text class="field-icon">🔒</text>
          </view>
          <view class="field-content">
            <text class="field-label">{{ $t('password') }}</text>
            <input class="field-input" type="password" name="password" maxlength="18" :placeholder="passwordPlaceholder" v-model="form.password" />
          </view>
        </view>

        <view class="form-field">
          <view class="field-icon-wrap">
            <text class="field-icon">🛡</text>
          </view>
          <view class="field-content">
            <text class="field-label">{{ $t('captcha') }}</text>
            <view class="captcha-row">
              <input class="field-input captcha-input" name="captcha" maxlength="18" :placeholder="$t('captchaPlaceholder')" v-model="form.captcha" />
              <image class="captcha-img" @tap="getCaptchaFunc()" :src="captchaImg" mode="aspectFit"></image>
            </view>
          </view>
        </view>
      </view>

      <!-- 按钮 -->
      <view class="action-area" :key="`actions-${langStore.locale}`">
        <button :key="`login-btn-${langStore.locale}`" class="action-btn action-btn-primary" @tap="login()">
          <text class="btn-text">{{ $t('loginBtn') }}</text>
          <text class="btn-arrow">→</text>
        </button>
        <button :key="`register-btn-${langStore.locale}`" class="action-btn action-btn-ghost" @tap="toRegister()">
          <text>{{ $t('goToRegister') }}</text>
        </button>
      </view>

    </view>

    <!-- 语言弹窗 -->
    <lang-switch v-model="showLangPicker" />

    <!-- 区号选择弹窗 -->
    <view class="popup-mask" v-if="showAreaCodePicker" @touchmove.prevent @tap="showAreaCodePicker = false">
      <view class="popup-content" @tap.stop>
        <view class="popup-handle"></view>
        <text class="popup-title">{{ $t('selectAreaCode') }}</text>
        <scroll-view scroll-y class="popup-scroll">
          <view class="area-item" v-for="item in areaCodeViews" :key="item._areaKey" @tap="selectArea(item._raw)">
            <text class="area-name">{{ item._countryNameText }}</text>
            <text class="area-code-val">{{ item.areaCode }}</text>
          </view>
        </scroll-view>
        <view class="popup-cancel" @tap.stop="showAreaCodePicker = false">
          <text>{{ $t('cancel') }}</text>
        </view>
      </view>
    </view>
  </view>
</template>
<script setup>
	import {getCaptcha} from "@/api/base.js"
	import { getEnabledPhoneAreaCodes } from "@/api/phoneAreaCode.js"
	import { getLoginConfig } from "@/api/sysConfig.js"
	import {
		reactive,
		ref,
		computed,
    onMounted,
    onUnmounted,
    watch
	} from 'vue';
  import { onLoad, onShow, onHide } from '@dcloudio/uni-app'

	import {useUserStore} from "@/pinia/modules/user.js"
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
  const inviteCodeFromShare = ref('')
  const localeRefreshing = ref(false)

  const tryAutoShowLangPicker = () => {
    if (showLangPicker.value) return
    if (!langStore.shouldAutoShowLanguagePicker()) return
    showLangPicker.value = true
  }

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
      inviteCodeFromShare.value = resolved
      uni.setStorageSync('pendingInviteCode', resolved)
    }
  }

	const goBack = () => {
    uni.switchTab({ url: '/pages/learning/home' })
	}

	const userStore = useUserStore()

	// 登录模式
	const loginMode = ref('username') // 'username' | 'phone'
	const phoneLoginEnabled = ref(false)
	const usernameLoginEnabled = ref(true)
	const showAreaCodePicker = ref(false)
	const areaCodes = ref([])
	const selectedAreaCode = ref('+86')
	const areaCodeViews = computed(() => areaCodes.value.map((item, index) => ({
    ...item,
    // 登录区号弹窗只读展示字段；登录提交仍使用 form.areaCode 和 selectArea 原对象。
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
		username: '',
		password: '',
		captcha: '',
		captchaId: '',
		phone: '',
		areaCode: '+86'
	})

	// 加载配置
	const loadConfig = async () => {
		try {
			const res = await getLoginConfig()
			if (res.code === 0 && res.data) {
				phoneLoginEnabled.value = res.data.phone_login_enabled === 'true'
				usernameLoginEnabled.value = res.data.username_login_enabled !== 'false' // 默认true
				// 互斥：两者不能同时关闭
				if (!phoneLoginEnabled.value && !usernameLoginEnabled.value) {
					usernameLoginEnabled.value = true
				}
				// 根据配置决定默认模式
				if (res.data.default_login_method === 'phone' && phoneLoginEnabled.value) {
					loginMode.value = 'phone'
				} else if (!usernameLoginEnabled.value && phoneLoginEnabled.value) {
					loginMode.value = 'phone'
				} else {
					loginMode.value = 'username'
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

	const captchaImg = ref("")

	const getCaptchaFunc = async () =>{
		const res = await getCaptcha()
		if(res.code === 0){
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

  onLoad((options) => {
    applyInviteCode(options)
    lastLoadedLocale.value = locale.value
  })

  onShow(() => {
    if (localeRefreshing.value) {
      return
    }
    const localeChanged = !!lastLoadedLocale.value && lastLoadedLocale.value !== locale.value
    if (localeChanged) {
      reloadLocaleSensitiveData()
    }
    lastLoadedLocale.value = locale.value
  })

  onHide(() => {
    showLangPicker.value = false
    showAreaCodePicker.value = false
  })

  watch(() => locale.value, (newLocale, oldLocale) => {
    if (localeRefreshing.value) return
    if (!oldLocale || newLocale === oldLocale) return
    reloadLocaleSensitiveData()
    lastLoadedLocale.value = newLocale
  })

	onMounted(() => {
    if (!inviteCodeFromShare.value) {
      applyInviteCode()
    }
    tryAutoShowLangPicker()
		getCaptchaFunc()
		loadConfig()
		loadAreaCodes()
    appConfigStore.loadConfig()
	})

  const forceRefreshByLocale = () => {
    if (localeRefreshing.value) return
    localeRefreshing.value = true
    const inviteCode = safeDecode(inviteCodeFromShare.value || uni.getStorageSync('pendingInviteCode') || '')
    const nextUrl = inviteCode
      ? `/pages/user/login?inviteCode=${encodeURIComponent(inviteCode)}`
      : '/pages/user/login'
    uni.reLaunch({ url: nextUrl })
  }

  onMounted(() => {
    if (typeof uni.$on === 'function') {
      uni.$on('app:locale-force-refresh', forceRefreshByLocale)
    }
  })

  onUnmounted(() => {
    if (typeof uni.$off === 'function') {
      uni.$off('app:locale-force-refresh', forceRefreshByLocale)
    }
  })


	//当前登录按钮操作
	const login = async () => {
		if (loginMode.value === 'username') {
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
		} else {
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
		if (!form.captcha) {
			uni.showToast({ title: $t.value('enterCaptcha'), icon: 'none' })
			return
		}

		let flag = false
		if (loginMode.value === 'username') {
			flag = await userStore.loginIn(form)
		} else {
			flag = await userStore.phoneLoginIn({
				areaCode: form.areaCode,
				phone: form.phone,
				password: form.password,
				captcha: form.captcha,
				captchaId: form.captchaId
			})
		}

		if(flag){
      uni.removeStorageSync('pendingInviteCode')
			uni.showToast({ title: $t.value('loginSuccess') })
      uni.switchTab({ url: '/pages/learning/home' })
			return
		}
		getCaptchaFunc()
	}
	//注册按钮点击
	const toRegister = () => {
    const inviteCode = String(inviteCodeFromShare.value || '').trim()
    const query = inviteCode ? `?inviteCode=${encodeURIComponent(inviteCode)}` : ''
    uni.navigateTo({ url: `/pages/user/register${query}` })
	}
</script>
<style lang="scss" scoped>
page { background-color: #F5F3FF; }

.nf-page {
  min-height: 100vh;
  background: #F5F3FF;
}

/* === Hero 区 === */
.login-hero {
  position: relative;
  overflow: hidden;
  padding: calc(var(--status-bar-height, 0px) + 36rpx) 32rpx 48rpx;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  border-radius: 0 0 36rpx 36rpx;
}

.hero-decor-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.12);
}

.decor-1 { width: 240rpx; height: 240rpx; top: -80rpx; right: -40rpx; }
.decor-2 { width: 160rpx; height: 160rpx; bottom: -60rpx; left: 200rpx; }

.hero-top-row {
  position: relative;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 32rpx;
}

.hero-back-btn,
.hero-lang-btn {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.25);
  display: flex;
  align-items: center;
  justify-content: center;

  &:active {
    transform: scale(0.92);
    background: rgba(255, 255, 255, 0.35);
  }
}

.hero-back-icon {
  color: #fff;
  font-size: 44rpx;
  line-height: 1;
}

.hero-lang-btn {
  width: auto;
  height: auto;
  border-radius: 999rpx;
  padding: 10rpx 24rpx;
}

.hero-lang-text {
  color: #fff;
  font-size: 26rpx;
  font-weight: 700;
  letter-spacing: 1rpx;
}

.hero-spacer { width: 64rpx; height: 64rpx; }

.hero-brand {
  position: relative;
  z-index: 2;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10rpx;
  margin-bottom: 32rpx;
}

.hero-logo-wrap {
  width: 96rpx;
  height: 96rpx;
  border-radius: 28rpx;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.2);
  border: 2rpx solid rgba(255, 255, 255, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
}

.hero-logo { width: 100%; height: 100%; }

.hero-logo-fallback {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 40rpx;
  font-weight: 800;
}

.hero-brand-name {
  font-size: 36rpx;
  font-weight: 800;
  color: #fff;
  letter-spacing: 1rpx;
}

.hero-brand-subtitle {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.8);
}

.hero-title-section {
  position: relative;
  z-index: 2;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14rpx;
}

.hero-title {
  font-size: 52rpx;
  font-weight: 800;
  color: #fff;
  letter-spacing: 2rpx;
}

.hero-title-line {
  width: 80rpx;
  height: 6rpx;
  border-radius: 999rpx;
  background: rgba(255, 255, 255, 0.6);
}

/* === 内容区 === */
.login-body {
  position: relative;
  z-index: 5;
  padding: 24rpx 32rpx 60rpx;
}

/* === 模式切换胶囊 === */
.mode-pills {
  display: flex;
  gap: 8rpx;
  margin-bottom: 24rpx;
  background: rgba(255, 255, 255, 0.6);
  padding: 6rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.16);
}

.mode-pill {
  flex: 1;
  text-align: center;
  padding: 16rpx 0;
  border-radius: 999rpx;
  font-size: 26rpx;
  color: #6B6F8D;
  font-weight: 500;
  transition: all 0.25s ease;

  &.active {
    background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
    color: #fff;
    font-weight: 700;
    box-shadow: 0 4rpx 16rpx rgba(108, 91, 255, 0.3);
  }
}

.mode-pill-text { font-size: inherit; color: inherit; }

/* === 表单卡片 === */
.form-card {
  background: #FFFFFF;
  border-radius: 28rpx;
  padding: 32rpx 28rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  box-shadow: 0 8rpx 24rpx rgba(108, 91, 255, 0.12);
}

.form-field {
  display: flex;
  gap: 16rpx;
  margin-bottom: 24rpx;

  &:last-child { margin-bottom: 0; }
}

.field-icon-wrap {
  width: 56rpx;
  height: 56rpx;
  border-radius: 50%;
  background: rgba(108, 91, 255, 0.08);
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-top: 30rpx;
}

.field-icon { font-size: 28rpx; }

.field-content {
  flex: 1;
  min-width: 0;
}

.field-label {
  display: block;
  font-size: 22rpx;
  color: #6B6F8D;
  font-weight: 600;
  margin-bottom: 8rpx;
}

.field-input {
  width: 100%;
  height: 80rpx;
  border-radius: 16rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.2);
  background: rgba(245, 243, 255, 0.5);
  padding: 0 20rpx;
  color: #1A1B3A;
  font-size: 28rpx;
  box-sizing: border-box;
}

.phone-row,
.captcha-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.area-code-btn {
  height: 80rpx;
  padding: 0 18rpx;
  border-radius: 16rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.2);
  background: rgba(245, 243, 255, 0.5);
  display: flex;
  align-items: center;
  gap: 8rpx;
  flex-shrink: 0;
}

.area-code-text { color: #1A1B3A; font-size: 26rpx; font-weight: 700; }
.area-code-arrow { color: #6B6F8D; font-size: 18rpx; }
.phone-input, .captcha-input { flex: 1; }

.captcha-img {
  width: 180rpx;
  height: 80rpx;
  border-radius: 16rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.2);
  background: rgba(245, 243, 255, 0.5);
  flex-shrink: 0;
}

/* === 按钮 === */
.action-area {
  margin-top: 28rpx;
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.action-btn {
  width: 100%;
  height: 92rpx;
  border-radius: 999rpx;
  font-size: 28rpx;
  font-weight: 700;
  border: none;
  margin: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8rpx;

  &:active { transform: translateY(2rpx); opacity: 0.9; }
  &::after { border: none; }
}

.action-btn-primary {
  color: #fff;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  box-shadow: 0 8rpx 24rpx rgba(108, 91, 255, 0.28);
}

.btn-arrow { font-size: 28rpx; }

.action-btn-ghost {
  color: #6D5BFF;
  background: rgba(255, 255, 255, 0.8);
  border: 1rpx solid rgba(108, 91, 255, 0.3);
}

/* === 底部语言切换 === */
.bottom-lang {
  margin-top: 32rpx;
  text-align: center;
}

.bottom-lang-label {
  font-size: 24rpx;
  color: #6D5BFF;
  font-weight: 700;
  padding: 10rpx 28rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.3);
  border-radius: 999rpx;
  background: rgba(255, 255, 255, 0.6);
}

/* === 弹窗 === */
.popup-mask {
  position: fixed;
  inset: 0;
  z-index: 10001;
  background: rgba(26, 27, 58, 0.5);
  display: flex;
  align-items: flex-end;
  touch-action: none;
}

.popup-content {
  width: 100%;
  max-height: 70vh;
  background: #FFFFFF;
  border-radius: 28rpx 28rpx 0 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  touch-action: auto;
}

.popup-handle {
  width: 64rpx;
  height: 8rpx;
  border-radius: 999rpx;
  background: rgba(108, 91, 255, 0.2);
  margin: 16rpx auto 20rpx;
}

.popup-title {
  display: block;
  text-align: center;
  font-size: 30rpx;
  font-weight: 800;
  color: #1A1B3A;
  padding-bottom: 20rpx;
  border-bottom: 1rpx solid rgba(108, 91, 255, 0.12);
  margin: 0 30rpx;
  flex-shrink: 0;
}

.popup-scroll {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

.popup-cancel {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 84rpx;
  margin: 16rpx 32rpx;
  margin-bottom: calc(16rpx + env(safe-area-inset-bottom));
  border-radius: 999rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.32);
  background: #fff;
  color: #6D5BFF;
  font-size: 30rpx;
  font-weight: 600;

  &:active {
    background: rgba(108, 91, 255, 0.08);
  }
}

.area-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 26rpx 40rpx;
  border-bottom: 1rpx solid rgba(108, 91, 255, 0.08);

  &:active { background: rgba(108, 91, 255, 0.06); }
}

.area-name { color: #1A1B3A; font-size: 28rpx; }
.area-code-val { color: #6D5BFF; font-size: 28rpx; font-weight: 700; }
</style>
