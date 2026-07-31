<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/pinia'
import { getCaptcha } from '@/api/base'
import { getEnabledPhoneAreaCodes } from '@/api/phoneAreaCode'
import { getLoginConfig } from '@/api/sysConfig'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// 登录模式
const loginMode = ref('username') // username | phone
const usernameLoginEnabled = ref(true)
const phoneLoginEnabled = ref(true)

// 表单数据
const form = reactive({
  username: '',
  phone: '',
  password: '',
  captcha: ''
})

// 验证码
const captchaImg = ref('')

// 区号
const areaCodeList = ref([])
const selectedAreaCode = ref('+86')
const showAreaCodePicker = ref(false)

// 加载状态
const loading = ref(false)

// 语言切换
const langLabel = computed(() => {
  const map = { zh: '中文', en: 'English' }
  return map[locale.value] || '中文'
})

function toggleLang() {
  locale.value = locale.value === 'zh' ? 'en' : 'zh'
}

// 获取验证码
async function getCaptchaFunc() {
  try {
    const res = await getCaptcha()
    if (res.code === 0) {
      captchaImg.value = 'data:image/png;base64,' + res.data
    }
  } catch (e) {
    console.warn('get captcha failed', e)
  }
}

// 获取区号列表
async function loadAreaCodes() {
  try {
    const res = await getEnabledPhoneAreaCodes()
    if (res.code === 0) {
      areaCodeList.value = res.data || []
      if (res.data?.length > 0) {
        selectedAreaCode.value = res.data[0].areaCode || '+86'
      }
    }
  } catch (e) {
    console.warn('load area codes failed', e)
  }
}

// 获取登录配置
async function loadLoginConfig() {
  try {
    const res = await getLoginConfig()
    if (res.code === 0) {
      const cfg = res.data || {}
      usernameLoginEnabled.value = cfg.usernameLoginEnabled !== false
      phoneLoginEnabled.value = cfg.phoneLoginEnabled === true
      if (!usernameLoginEnabled.value && phoneLoginEnabled.value) {
        loginMode.value = 'phone'
      }
    }
  } catch (e) {
    console.warn('load login config failed', e)
  }
}

// 选择区号
function selectArea(item) {
  selectedAreaCode.value = item.areaCode
  showAreaCodePicker.value = false
}

// 登录
async function handleLogin() {
  if (loading.value) return

  if (loginMode.value === 'username' && !form.username.trim()) {
    alert('请输入用户名')
    return
  }
  if (loginMode.value === 'phone' && !form.phone.trim()) {
    alert('请输入手机号')
    return
  }
  if (!form.password.trim()) {
    alert('请输入密码')
    return
  }
  if (!form.captcha.trim()) {
    alert('请输入验证码')
    return
  }

  loading.value = true
  try {
    let res
    if (loginMode.value === 'username') {
      res = await userStore.loginIn({
        username: form.username.trim(),
        password: form.password,
        captcha: form.captcha.trim()
      })
    } else {
      res = await userStore.phoneLoginIn({
        areaCode: selectedAreaCode.value,
        phone: form.phone.trim(),
        password: form.password,
        captcha: form.captcha.trim()
      })
    }

    if (res.code === 0) {
      const redirect = route.query.redirect || '/home'
      router.push(redirect)
    } else {
      alert(res.msg || '登录失败')
      getCaptchaFunc()
    }
  } catch (e) {
    alert(e.message || '登录失败')
    getCaptchaFunc()
  } finally {
    loading.value = false
  }
}

function goRegister() {
  router.push('/register')
}

onMounted(() => {
  getCaptchaFunc()
  loadLoginConfig()
  loadAreaCodes()
})
</script>

<template>
  <div class="login-page">
    <!-- 左侧品牌区 -->
    <div class="login-left">
      <div class="decor-circle decor-1"></div>
      <div class="decor-circle decor-2"></div>
      <div class="decor-circle decor-3"></div>

      <div class="brand-section">
        <div class="brand-logo">
          <span>R</span>
        </div>
        <h1 class="brand-name">RichRoad</h1>
        <p class="brand-slogan">{{ t('home.slogan') }}</p>
      </div>

      <div class="feature-list">
        <div class="feature-item">
          <div class="feature-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="23 7 16 12 23 17 23 7"></polygon>
              <rect x="1" y="5" width="15" height="14" rx="2" ry="2"></rect>
            </svg>
          </div>
          <div class="feature-text">
            <h4>视频学习</h4>
            <p>精选视频课程，沉浸式学习体验</p>
          </div>
        </div>
        <div class="feature-item">
          <div class="feature-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"></path>
              <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"></path>
            </svg>
          </div>
          <div class="feature-text">
            <h4>每日日记</h4>
            <p>每天一篇，培养英语阅读习惯</p>
          </div>
        </div>
        <div class="feature-item">
          <div class="feature-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z"></path>
            </svg>
          </div>
          <div class="feature-text">
            <h4>单词跟打</h4>
            <p>边打字边记单词，高效记忆</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧表单区 -->
    <div class="login-right">
      <div class="lang-switcher" @click="toggleLang">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="2" y1="12" x2="22" y2="12"></line>
          <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path>
        </svg>
        <span>{{ langLabel }}</span>
      </div>

      <div class="form-wrapper">
        <div class="form-header">
          <h2>{{ t('login.title') }}</h2>
          <p>{{ t('login.subtitle') }}</p>
        </div>

        <!-- 模式切换 -->
        <div v-if="usernameLoginEnabled && phoneLoginEnabled" class="mode-tabs">
          <div
            class="mode-tab"
            :class="{ active: loginMode === 'username' }"
            @click="loginMode = 'username'"
          >
            {{ t('login.accountLogin') }}
          </div>
          <div
            class="mode-tab"
            :class="{ active: loginMode === 'phone' }"
            @click="loginMode = 'phone'"
          >
            {{ t('login.phoneLogin') }}
          </div>
          <div class="mode-tab-indicator" :style="{ left: loginMode === 'username' ? '0' : '50%' }"></div>
        </div>

        <!-- 表单卡片 -->
        <div class="form-card">
          <!-- 用户名 -->
          <div v-if="loginMode === 'username'" class="form-group">
            <label class="form-label">{{ t('login.username') }}</label>
            <div class="input-wrapper">
              <span class="input-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                  <circle cx="12" cy="7" r="4"></circle>
                </svg>
              </span>
              <input
                v-model="form.username"
                type="text"
                class="form-input"
                :placeholder="t('login.username')"
                maxlength="30"
                @keyup.enter="handleLogin"
              />
            </div>
          </div>

          <!-- 手机号 -->
          <div v-if="loginMode === 'phone'" class="form-group">
            <label class="form-label">{{ t('login.phone') }}</label>
            <div class="input-wrapper">
              <span class="input-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="5" y="2" width="14" height="20" rx="2" ry="2"></rect>
                  <line x1="12" y1="18" x2="12.01" y2="18"></line>
                </svg>
              </span>
              <div class="area-code-selector" @click="showAreaCodePicker = true">
                <span class="area-code-text">{{ selectedAreaCode }}</span>
                <svg class="area-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="6 9 12 15 18 9"></polyline>
                </svg>
              </div>
              <div class="divider-v"></div>
              <input
                v-model="form.phone"
                type="tel"
                class="form-input phone-input"
                :placeholder="t('login.phone')"
                maxlength="15"
                @keyup.enter="handleLogin"
              />
            </div>
          </div>

          <!-- 密码 -->
          <div class="form-group">
            <label class="form-label">{{ t('login.password') }}</label>
            <div class="input-wrapper">
              <span class="input-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
                </svg>
              </span>
              <input
                v-model="form.password"
                type="password"
                class="form-input"
                :placeholder="t('login.password')"
                maxlength="18"
                @keyup.enter="handleLogin"
              />
            </div>
          </div>

          <!-- 验证码 -->
          <div class="form-group">
            <label class="form-label">{{ t('login.captcha') }}</label>
            <div class="input-wrapper">
              <span class="input-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M9 12l2 2 4-4"></path>
                  <path d="M21 12c-1 0-3-1-3-3s2-3 3-3 3 1 3 3-2 3-3 3z"></path>
                  <circle cx="12" cy="12" r="10"></circle>
                </svg>
              </span>
              <input
                v-model="form.captcha"
                type="text"
                class="form-input captcha-input"
                :placeholder="t('login.captcha')"
                maxlength="6"
                @keyup.enter="handleLogin"
              />
              <img
                v-if="captchaImg"
                class="captcha-image"
                :src="captchaImg"
                @click="getCaptchaFunc"
                alt="captcha"
              />
            </div>
          </div>

          <!-- 登录按钮 -->
          <button class="login-btn" :disabled="loading" @click="handleLogin">
            <span v-if="!loading">{{ t('login.loginBtn') }}</span>
            <span v-else class="btn-loading">...</span>
          </button>

          <!-- 注册入口 -->
          <div class="register-link">
            <span>{{ t('login.noAccount') }}</span>
            <a class="link-text" @click="goRegister">{{ t('login.goRegister') }}</a>
          </div>
        </div>
      </div>
    </div>

    <!-- 区号选择弹窗 -->
    <div v-if="showAreaCodePicker" class="modal-mask" @click.self="showAreaCodePicker = false">
      <div class="modal-box">
        <div class="modal-header">
          <h3>{{ t('login.areaCode') }}</h3>
          <button class="modal-close" @click="showAreaCodePicker = false">×</button>
        </div>
        <div class="modal-body">
          <div
            v-for="item in areaCodeList"
            :key="item.areaCode"
            class="area-item"
            :class="{ active: selectedAreaCode === item.areaCode }"
            @click="selectArea(item)"
          >
            <span class="area-country">{{ item.countryName || item.country }}</span>
            <span class="area-code">{{ item.areaCode }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.login-page {
  display: flex;
  min-height: 100vh;
  background: $bg-page;
}

// ========== 左侧品牌区 ==========
.login-left {
  flex: 1;
  position: relative;
  background: $gradient-hero;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 80px 60px;
  overflow: hidden;

  .decor-circle {
    position: absolute;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.1);

    &.decor-1 {
      width: 300px;
      height: 300px;
      top: -80px;
      left: -80px;
    }

    &.decor-2 {
      width: 200px;
      height: 200px;
      bottom: 20%;
      right: 10%;
      background: rgba(255, 255, 255, 0.08);
    }

    &.decor-3 {
      width: 120px;
      height: 120px;
      top: 30%;
      right: 20%;
      background: rgba(255, 255, 255, 0.12);
    }
  }
}

.brand-section {
  text-align: center;
  color: $text-white;
  margin-bottom: 60px;
  position: relative;
  z-index: 1;

  .brand-logo {
    width: 80px;
    height: 80px;
    margin: 0 auto 24px;
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(10px);
    border-radius: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 40px;
    font-weight: $font-weight-bold;
    border: 2px solid rgba(255, 255, 255, 0.3);
  }

  .brand-name {
    font-size: 42px;
    font-weight: $font-weight-bold;
    margin-bottom: 12px;
    letter-spacing: 1px;
  }

  .brand-slogan {
    font-size: 16px;
    opacity: 0.85;
    font-weight: $font-weight-regular;
  }
}

.feature-list {
  max-width: 360px;
  position: relative;
  z-index: 1;

  .feature-item {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px 20px;
    margin-bottom: 16px;
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    border-radius: 16px;
    border: 1px solid rgba(255, 255, 255, 0.15);
    color: $text-white;
    transition: all 0.3s ease;

    &:hover {
      transform: translateX(8px);
      background: rgba(255, 255, 255, 0.18);
    }

    .feature-icon {
      width: 44px;
      height: 44px;
      background: rgba(255, 255, 255, 0.2);
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;

      svg {
        width: 22px;
        height: 22px;
      }
    }

    .feature-text {
      h4 {
        font-size: 16px;
        font-weight: $font-weight-semibold;
        margin-bottom: 4px;
      }

      p {
        font-size: 13px;
        opacity: 0.8;
      }
    }
  }
}

// ========== 右侧表单区 ==========
.login-right {
  width: 560px;
  background: $bg-card;
  display: flex;
  flex-direction: column;
  position: relative;
}

.lang-switcher {
  position: absolute;
  top: 24px;
  right: 32px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: $bg-input;
  border-radius: $radius-pill;
  cursor: pointer;
  color: $text-secondary;
  font-size: $font-size-sm;
  transition: all 0.2s ease;

  &:hover {
    background: $bg-hover;
    color: $primary-color;
  }

  svg {
    width: 16px;
    height: 16px;
  }
}

.form-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 48px 56px;
}

.form-header {
  margin-bottom: 32px;

  h2 {
    font-size: 28px;
    font-weight: $font-weight-bold;
    color: $text-primary;
    margin-bottom: 8px;
  }

  p {
    font-size: 14px;
    color: $text-secondary;
  }
}

// 模式切换
.mode-tabs {
  position: relative;
  display: flex;
  background: $bg-input;
  border-radius: $radius-pill;
  padding: 4px;
  margin-bottom: 28px;

  .mode-tab {
    flex: 1;
    text-align: center;
    padding: 10px 0;
    font-size: 14px;
    color: $text-secondary;
    cursor: pointer;
    position: relative;
    z-index: 1;
    transition: color 0.25s ease;

    &.active {
      color: $primary-color;
      font-weight: $font-weight-semibold;
    }
  }

  .mode-tab-indicator {
    position: absolute;
    top: 4px;
    bottom: 4px;
    width: 50%;
    background: $bg-card;
    border-radius: $radius-pill;
    box-shadow: 0 2px 8px rgba(109, 91, 255, 0.15);
    transition: left 0.3s ease;
  }
}

// 表单
.form-card {
  .form-group {
    margin-bottom: 20px;

    .form-label {
      display: block;
      font-size: 13px;
      font-weight: $font-weight-medium;
      color: $text-secondary;
      margin-bottom: 8px;
    }
  }
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  height: 48px;
  background: $bg-input;
  border: 1.5px solid transparent;
  border-radius: $radius-md;
  transition: all 0.2s ease;

  &:focus-within {
    background: $bg-card;
    border-color: $primary-color;
    box-shadow: 0 0 0 4px rgba(109, 91, 255, 0.1);
  }

  .input-icon {
    width: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: $text-tertiary;

    svg {
      width: 18px;
      height: 18px;
    }
  }

  .form-input {
    flex: 1;
    height: 100%;
    background: transparent;
    font-size: 15px;
    color: $text-primary;
    padding-right: 16px;

    &::placeholder {
      color: $text-placeholder;
    }
  }

  .phone-input {
    padding-right: 16px;
  }

  .area-code-selector {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 0 12px 0 4px;
    cursor: pointer;
    color: $text-primary;
    font-size: 14px;
    font-weight: $font-weight-medium;

    .area-code-text {
      white-space: nowrap;
    }

    .area-arrow {
      width: 14px;
      height: 14px;
      color: $text-tertiary;
    }
  }

  .divider-v {
    width: 1px;
    height: 20px;
    background: $border-color;
    margin: 0 4px;
  }

  .captcha-input {
    padding-right: 120px;
  }

  .captcha-image {
    position: absolute;
    right: 6px;
    top: 50%;
    transform: translateY(-50%);
    height: 36px;
    border-radius: 6px;
    cursor: pointer;
  }
}

.login-btn {
  width: 100%;
  height: 50px;
  margin-top: 28px;
  background: $gradient-primary;
  color: $text-white;
  border-radius: $radius-md;
  font-size: 16px;
  font-weight: $font-weight-semibold;
  box-shadow: 0 8px 20px rgba(109, 91, 255, 0.35);
  transition: all 0.25s ease;

  &:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 12px 28px rgba(109, 91, 255, 0.45);
  }

  &:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }

  .btn-loading {
    display: inline-block;
    animation: pulse 1s infinite;
  }
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.register-link {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: $text-tertiary;

  .link-text {
    color: $primary-color;
    font-weight: $font-weight-medium;
    cursor: pointer;
    margin-left: 4px;

    &:hover {
      text-decoration: underline;
    }
  }
}

// ========== 区号选择弹窗 ==========
.modal-mask {
  position: fixed;
  inset: 0;
  background: $bg-mask;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: $z-modal;
}

.modal-box {
  width: 420px;
  max-height: 70vh;
  background: $bg-card;
  border-radius: $radius-xl;
  display: flex;
  flex-direction: column;
  box-shadow: $shadow-xl;
  animation: modalIn 0.25s ease;
}

@keyframes modalIn {
  from {
    opacity: 0;
    transform: scale(0.95) translateY(-10px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid $border-light;

  h3 {
    font-size: 18px;
    font-weight: $font-weight-semibold;
    color: $text-primary;
  }

  .modal-close {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
    color: $text-tertiary;
    border-radius: 50%;
    transition: all 0.2s ease;

    &:hover {
      background: $bg-hover;
      color: $text-primary;
    }
  }
}

.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.area-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 24px;
  cursor: pointer;
  transition: background 0.15s ease;

  &:hover {
    background: $bg-hover;
  }

  &.active {
    background: $bg-active;
    color: $primary-color;
    font-weight: $font-weight-medium;
  }

  .area-country {
    font-size: 14px;
  }

  .area-code {
    font-size: 14px;
    font-family: monospace;
  }
}

// 响应式
@media (max-width: 960px) {
  .login-left {
    display: none;
  }

  .login-right {
    width: 100%;
  }
}
</style>
