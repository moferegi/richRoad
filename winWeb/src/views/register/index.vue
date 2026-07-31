<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getCaptcha, register, phoneRegister } from '@/api/base'
import { getEnabledPhoneAreaCodes } from '@/api/phoneAreaCode'
import { getLoginConfig } from '@/api/sysConfig'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()

// 注册模式
const registerMode = ref('username') // username | phone
const usernameRegisterEnabled = ref(true)
const phoneRegisterEnabled = ref(false)

// 表单数据
const form = reactive({
  username: '',
  phone: '',
  password: '',
  confirmPassword: '',
  captcha: '',
  inviteCode: ''
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
      usernameRegisterEnabled.value = cfg.usernameLoginEnabled !== false
      phoneRegisterEnabled.value = cfg.phoneRegisterEnabled === true
      if (!usernameRegisterEnabled.value && phoneRegisterEnabled.value) {
        registerMode.value = 'phone'
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

// 注册
async function handleRegister() {
  if (loading.value) return

  if (registerMode.value === 'username' && !form.username.trim()) {
    alert('请输入用户名')
    return
  }
  if (registerMode.value === 'phone' && !form.phone.trim()) {
    alert('请输入手机号')
    return
  }
  if (!form.password.trim()) {
    alert('请输入密码')
    return
  }
  if (form.password !== form.confirmPassword) {
    alert(t('register.passwordNotMatch'))
    return
  }
  if (!form.captcha.trim()) {
    alert('请输入验证码')
    return
  }

  loading.value = true
  try {
    const data = {
      password: form.password,
      captcha: form.captcha.trim(),
      inviteCode: form.inviteCode.trim() || undefined
    }

    let res
    if (registerMode.value === 'username') {
      res = await register({
        ...data,
        username: form.username.trim()
      })
    } else {
      res = await phoneRegister({
        ...data,
        areaCode: selectedAreaCode.value,
        phone: form.phone.trim()
      })
    }

    if (res.code === 0) {
      alert(t('register.registerSuccess'))
      router.push('/login')
    } else {
      alert(res.msg || '注册失败')
      getCaptchaFunc()
    }
  } catch (e) {
    alert(e.message || '注册失败')
    getCaptchaFunc()
  } finally {
    loading.value = false
  }
}

function goLogin() {
  router.push('/login')
}

// 从 URL 获取邀请码
function getInviteCodeFromUrl() {
  const code = route.query.inviteCode || route.query.code || ''
  if (code) {
    form.inviteCode = String(code)
  }
}

onMounted(() => {
  getCaptchaFunc()
  loadLoginConfig()
  loadAreaCodes()
  getInviteCodeFromUrl()
})
</script>

<template>
  <div class="register-page">
    <!-- 左侧品牌区 -->
    <div class="register-left">
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

      <div class="benefit-list">
        <div class="benefit-item">
          <div class="benefit-number">01</div>
          <div class="benefit-content">
            <h4>海量学习资源</h4>
            <p>视频、日记、单词全方位覆盖</p>
          </div>
        </div>
        <div class="benefit-item">
          <div class="benefit-number">02</div>
          <div class="benefit-content">
            <h4>打卡奖励机制</h4>
            <p>每日学习打卡，赢取积分好礼</p>
          </div>
        </div>
        <div class="benefit-item">
          <div class="benefit-number">03</div>
          <div class="benefit-content">
            <h4>科学学习方法</h4>
            <p>跟读、跟打、复习加深记忆</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧表单区 -->
    <div class="register-right">
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
          <h2>{{ t('register.title') }}</h2>
          <p>{{ t('register.subtitle') }}</p>
        </div>

        <!-- 模式切换 -->
        <div v-if="usernameRegisterEnabled && phoneRegisterEnabled" class="mode-tabs">
          <div
            class="mode-tab"
            :class="{ active: registerMode === 'username' }"
            @click="registerMode = 'username'"
          >
            {{ t('register.accountRegister') }}
          </div>
          <div
            class="mode-tab"
            :class="{ active: registerMode === 'phone' }"
            @click="registerMode = 'phone'"
          >
            {{ t('register.phoneRegister') }}
          </div>
          <div class="mode-tab-indicator" :style="{ left: registerMode === 'username' ? '0' : '50%' }"></div>
        </div>

        <!-- 表单卡片 -->
        <div class="form-card">
          <!-- 用户名 -->
          <div v-if="registerMode === 'username'" class="form-group">
            <label class="form-label">{{ t('register.username') }}</label>
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
                :placeholder="t('register.username')"
                maxlength="30"
                @keyup.enter="handleRegister"
              />
            </div>
          </div>

          <!-- 手机号 -->
          <div v-if="registerMode === 'phone'" class="form-group">
            <label class="form-label">{{ t('register.phone') }}</label>
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
                :placeholder="t('register.phone')"
                maxlength="15"
                @keyup.enter="handleRegister"
              />
            </div>
          </div>

          <!-- 密码 -->
          <div class="form-group">
            <label class="form-label">{{ t('register.password') }}</label>
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
                :placeholder="t('register.password')"
                maxlength="18"
                @keyup.enter="handleRegister"
              />
            </div>
          </div>

          <!-- 确认密码 -->
          <div class="form-group">
            <label class="form-label">{{ t('register.confirmPassword') }}</label>
            <div class="input-wrapper">
              <span class="input-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M9 12l2 2 4-4"></path>
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
                </svg>
              </span>
              <input
                v-model="form.confirmPassword"
                type="password"
                class="form-input"
                :placeholder="t('register.confirmPassword')"
                maxlength="18"
                @keyup.enter="handleRegister"
              />
            </div>
          </div>

          <!-- 验证码 -->
          <div class="form-group">
            <label class="form-label">{{ t('register.captcha') }}</label>
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
                :placeholder="t('register.captcha')"
                maxlength="6"
                @keyup.enter="handleRegister"
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

          <!-- 邀请码 -->
          <div class="form-group">
            <label class="form-label">{{ t('register.inviteCode') }}</label>
            <div class="input-wrapper">
              <span class="input-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path>
                  <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path>
                </svg>
              </span>
              <input
                v-model="form.inviteCode"
                type="text"
                class="form-input"
                :placeholder="t('register.inviteCode') + ' (选填)'"
                maxlength="30"
                @keyup.enter="handleRegister"
              />
            </div>
          </div>

          <!-- 注册按钮 -->
          <button class="register-btn" :disabled="loading" @click="handleRegister">
            <span v-if="!loading">{{ t('register.registerBtn') }}</span>
            <span v-else class="btn-loading">...</span>
          </button>

          <!-- 登录入口 -->
          <div class="login-link">
            <span>{{ t('login.haveAccount') }}</span>
            <a class="link-text" @click="goLogin">{{ t('login.goLogin') }}</a>
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
.register-page {
  display: flex;
  min-height: 100vh;
  background: $bg-page;
}

// ========== 左侧品牌区 ==========
.register-left {
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
      width: 280px;
      height: 280px;
      top: -60px;
      right: -60px;
    }

    &.decor-2 {
      width: 180px;
      height: 180px;
      bottom: 25%;
      left: 10%;
      background: rgba(255, 255, 255, 0.08);
    }

    &.decor-3 {
      width: 100px;
      height: 100px;
      top: 35%;
      left: 25%;
      background: rgba(255, 255, 255, 0.12);
    }
  }
}

.brand-section {
  text-align: center;
  color: $text-white;
  margin-bottom: 48px;
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

.benefit-list {
  max-width: 380px;
  position: relative;
  z-index: 1;

  .benefit-item {
    display: flex;
    align-items: flex-start;
    gap: 20px;
    padding: 20px 24px;
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

    .benefit-number {
      font-size: 28px;
      font-weight: $font-weight-bold;
      opacity: 0.4;
      line-height: 1;
      font-family: monospace;
      flex-shrink: 0;
    }

    .benefit-content {
      h4 {
        font-size: 16px;
        font-weight: $font-weight-semibold;
        margin-bottom: 6px;
      }

      p {
        font-size: 13px;
        opacity: 0.8;
        line-height: 1.5;
      }
    }
  }
}

// ========== 右侧表单区 ==========
.register-right {
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
  padding: 40px 56px;
  overflow-y: auto;
}

.form-header {
  margin-bottom: 24px;

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
  margin-bottom: 24px;

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
    margin-bottom: 16px;

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
  height: 46px;
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

.register-btn {
  width: 100%;
  height: 48px;
  margin-top: 24px;
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

.login-link {
  text-align: center;
  margin-top: 16px;
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
  .register-left {
    display: none;
  }

  .register-right {
    width: 100%;
  }
}
</style>
