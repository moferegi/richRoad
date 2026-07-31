<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/pinia/modules/user'
import { useLangStore } from '@/pinia/modules/lang'
import { getAsset, exchangeTime } from '@/api/learning'
import { getPointsExchangeRate } from '@/api/sysConfig'

const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()
const langStore = useLangStore()

// 用户资产
const asset = ref(null)
const loadingAsset = ref(true)

// 兑换
const showExchangeModal = ref(false)
const exchangePoints = ref('')
const exchangeRate = ref(100) // 默认 100积分 = 1分钟
const exchanging = ref(false)

// 菜单配置
const menuItems = [
  {
    key: 'watchHistory',
    icon: 'history',
    label: 'profile.watchHistory',
    path: '/learning/watch-history'
  },
  {
    key: 'checkinRecord',
    icon: 'calendar',
    label: 'profile.checkinRecord',
    path: '/learning/checkin-record'
  },
  {
    key: 'collections',
    icon: 'star',
    label: 'profile.collections',
    path: '/learning/collections'
  },
  {
    key: 'errorLog',
    icon: 'alert',
    label: 'profile.errorLog',
    path: '/learning/error-log'
  },
  {
    key: 'pointHistory',
    icon: 'coin',
    label: 'profile.pointHistory',
    path: '/learning/point-history'
  },
  {
    key: 'freeTimeHistory',
    icon: 'clock',
    label: 'profile.freeTimeHistory',
    path: '/learning/free-time-history'
  },
  {
    key: 'contactService',
    icon: 'message',
    label: 'profile.contactService',
    path: '/kefu'
  },
  {
    key: 'switchLanguage',
    icon: 'globe',
    label: 'profile.switchLanguage',
    action: 'switchLang'
  }
]

// 加载用户资产
async function loadAsset() {
  loadingAsset.value = true
  try {
    const res = await getAsset()
    if (res.code === 0) {
      asset.value = res.data
    }
  } catch (e) {
    console.warn('load asset failed', e)
  } finally {
    loadingAsset.value = false
  }
}

// 加载兑换比例
async function loadExchangeRate() {
  try {
    const res = await getPointsExchangeRate()
    if (res.code === 0) {
      const val = Number(res.data?.configValue || res.data?.value || 10)
      exchangeRate.value = val > 0 ? val : 100
    }
  } catch (e) {
    console.warn('load exchange rate failed', e)
  }
}

// 打开兑换弹窗
function openExchange() {
  exchangePoints.value = ''
  showExchangeModal.value = true
}

// 执行兑换
async function doExchange() {
  const points = Number(exchangePoints.value)
  if (!points || points <= 0) return
  if (points > (asset.value?.points || 0)) return

  exchanging.value = true
  try {
    const res = await exchangeTime(points)
    if (res.code === 0) {
      showExchangeModal.value = false
      loadAsset()
    }
  } catch (e) {
    console.warn('exchange failed', e)
  } finally {
    exchanging.value = false
  }
}

// 菜单点击
function handleMenuClick(item) {
  if (item.action === 'switchLang') {
    switchLanguage()
  } else if (item.path) {
    router.push(item.path)
  }
}

// 切换语言
function switchLanguage() {
  const next = langStore.locale === 'zh' ? 'en' : 'zh'
  langStore.setLocale(next)
}

// 退出登录
function logout() {
  userStore.loginOut()
  router.push('/login')
}

// 格式化时长
function formatDuration(minutes) {
  if (!minutes) return '0 分钟'
  if (minutes < 60) return `${minutes} 分钟`
  const h = Math.floor(minutes / 60)
  const m = minutes % 60
  return m > 0 ? `${h} 小时 ${m} 分钟` : `${h} 小时`
}

onMounted(() => {
  loadAsset()
  loadExchangeRate()
})
</script>

<template>
  <div class="profile-page">
    <div class="profile-content">
      <!-- 顶部用户信息卡 -->
      <div class="user-card">
        <div class="user-avatar">
          <img v-if="userStore.userInfo?.avatar" :src="userStore.userInfo.avatar" alt="avatar" />
          <div v-else class="avatar-placeholder">
            {{ (userStore.username || 'U').charAt(0).toUpperCase() }}
          </div>
        </div>
        <div class="user-info">
          <h2 class="user-name">{{ userStore.username || '用户' }}</h2>
          <div class="user-badge">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
            </svg>
            <span>积分会员</span>
          </div>
        </div>
      </div>

      <!-- 资产概览 -->
      <div class="asset-overview">
        <div class="asset-card points-card">
          <div class="asset-icon points-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <path d="M12 6v12M8 10h8M8 14h8"></path>
            </svg>
          </div>
          <div class="asset-info">
            <span class="asset-label">{{ t('profile.points') }}</span>
            <span class="asset-value">
              {{ loadingAsset ? '--' : (asset?.totalPoints || 0) }}
            </span>
          </div>
        </div>

        <div class="asset-card time-card">
          <div class="asset-icon time-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <polyline points="12 6 12 12 16 14"></polyline>
            </svg>
          </div>
          <div class="asset-info">
            <span class="asset-label">{{ t('profile.freeTime') }}</span>
            <span class="asset-value">
              {{ loadingAsset ? '--' : formatDuration(asset?.freeMinutes || asset?.freeTime || 0) }}
            </span>
          </div>
        </div>
      </div>

      <!-- 积分兑换 -->
      <div class="exchange-section">
        <div class="section-header">
          <h3 class="section-title">{{ t('profile.exchange') }}</h3>
          <span class="exchange-rate">
            {{ t('profile.exchangeRate') }}: {{ exchangeRate }} 积分 = 1 分钟
          </span>
        </div>
        <button class="btn btn-primary exchange-btn" @click="openExchange">
          {{ t('profile.exchangeBtn') }}
        </button>
      </div>

      <!-- 功能菜单 -->
      <div class="menu-grid">
        <div
          v-for="item in menuItems"
          :key="item.key"
          class="menu-item"
          @click="handleMenuClick(item)"
        >
          <div class="menu-icon" :class="'icon-' + item.icon">
            <!-- 不同图标 -->
            <svg v-if="item.icon === 'history'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="1 4 1 10 7 10"></polyline>
              <path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"></path>
            </svg>
            <svg v-else-if="item.icon === 'calendar'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
              <line x1="16" y1="2" x2="16" y2="6"></line>
              <line x1="8" y1="2" x2="8" y2="6"></line>
              <line x1="3" y1="10" x2="21" y2="10"></line>
            </svg>
            <svg v-else-if="item.icon === 'star'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
            </svg>
            <svg v-else-if="item.icon === 'alert'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
              <line x1="12" y1="9" x2="12" y2="13"></line>
              <line x1="12" y1="17" x2="12.01" y2="17"></line>
            </svg>
            <svg v-else-if="item.icon === 'coin'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <path d="M12 6v12M8 10h8M8 14h8"></path>
            </svg>
            <svg v-else-if="item.icon === 'clock'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <polyline points="12 6 12 12 16 14"></polyline>
            </svg>
            <svg v-else-if="item.icon === 'message'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
            </svg>
            <svg v-else-if="item.icon === 'globe'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="2" y1="12" x2="22" y2="12"></line>
              <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path>
            </svg>
          </div>
          <span class="menu-label">{{ t(item.label) }}</span>
          <svg class="menu-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="9 18 15 12 9 6"></polyline>
          </svg>
        </div>
      </div>

      <!-- 退出登录 -->
      <button class="logout-btn" @click="logout">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
          <polyline points="16 17 21 12 16 7"></polyline>
          <line x1="21" y1="12" x2="9" y2="12"></line>
        </svg>
        {{ t('profile.logout') }}
      </button>
    </div>

    <!-- 兑换弹窗 -->
    <div v-if="showExchangeModal" class="modal-mask" @click="showExchangeModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ t('profile.exchange') }}</h3>
          <button class="close-btn" @click="showExchangeModal = false">×</button>
        </div>
        <div class="modal-body">
          <div class="exchange-info">
            <div class="info-row">
              <span>当前积分</span>
              <span class="highlight">{{ asset?.totalPoints || 0 }}</span>
            </div>
            <div class="info-row">
              <span>{{ t('profile.exchangeRate') }}</span>
              <span>{{ exchangeRate }} 积分 = 1 分钟</span>
            </div>
          </div>

          <div class="form-item">
            <label>输入兑换积分数量</label>
            <input
              v-model="exchangePoints"
              type="number"
              placeholder="请输入积分数量"
              min="1"
              :max="asset?.totalPoints || 0"
            />
          </div>

          <div v-if="exchangePoints && Number(exchangePoints) > 0" class="exchange-preview">
            预计兑换时长：<span class="highlight">{{ Math.floor(Number(exchangePoints) / exchangeRate) }} 分钟</span>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" @click="showExchangeModal = false">{{ t('common.cancel') }}</button>
          <button
            class="btn btn-primary"
            :disabled="!exchangePoints || Number(exchangePoints) <= 0 || Number(exchangePoints) > (asset?.totalPoints || 0) || exchanging"
            @click="doExchange"
          >
            {{ exchanging ? '兑换中...' : t('common.confirm') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use "@/assets/styles/variables.scss" as *;

.profile-page {
  min-height: 100%;
}

.profile-content {
  max-width: $content-max-width;
  margin: 0 auto;
  padding: $spacing-xl;
}

// ========== 用户信息卡 ==========
.user-card {
  display: flex;
  align-items: center;
  gap: $spacing-lg;
  padding: $spacing-xl;
  background: $gradient-hero;
  border-radius: $radius-xl;
  margin-bottom: $spacing-xl;
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    width: 200px;
    height: 200px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.1);
    top: -80px;
    right: -40px;
  }

  &::after {
    content: '';
    position: absolute;
    width: 120px;
    height: 120px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.08);
    bottom: -40px;
    right: 20%;
  }
}

.user-avatar {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  overflow: hidden;
  border: 3px solid rgba(255, 255, 255, 0.3);
  flex-shrink: 0;
  position: relative;
  z-index: 1;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .avatar-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(255, 255, 255, 0.2);
    color: $text-white;
    font-size: 28px;
    font-weight: $font-weight-bold;
  }
}

.user-info {
  flex: 1;
  position: relative;
  z-index: 1;

  .user-name {
    font-size: $font-size-xxl;
    font-weight: $font-weight-bold;
    color: $text-white;
    margin: 0 0 8px 0;
  }

  .user-badge {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 4px 12px;
    background: rgba(255, 255, 255, 0.2);
    border-radius: $radius-pill;
    color: $text-white;
    font-size: 12px;

    svg {
      width: 14px;
      height: 14px;
    }
  }
}

// ========== 资产概览 ==========
.asset-overview {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: $spacing-lg;
  margin-bottom: $spacing-xl;
}

.asset-card {
  display: flex;
  align-items: center;
  gap: $spacing-md;
  padding: $spacing-lg;
  background: $bg-card;
  border-radius: $radius-lg;
  box-shadow: $shadow-card;
  border: 1px solid $border-light;
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: $shadow-md;
  }
}

.asset-icon {
  width: 48px;
  height: 48px;
  border-radius: $radius-md;
  display: flex;
  align-items: center;
  justify-content: center;
  color: $text-white;
  flex-shrink: 0;

  svg {
    width: 24px;
    height: 24px;
  }

  &.points-icon {
    background: linear-gradient(135deg, #FFB800 0%, #FF8A00 100%);
  }

  &.time-icon {
    background: $gradient-primary;
  }
}

.asset-info {
  display: flex;
  flex-direction: column;
  gap: 4px;

  .asset-label {
    font-size: 13px;
    color: $text-secondary;
  }

  .asset-value {
    font-size: 22px;
    font-weight: $font-weight-bold;
    color: $text-primary;
  }
}

// ========== 兑换区域 ==========
.exchange-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: $spacing-lg;
  background: $bg-card;
  border-radius: $radius-lg;
  box-shadow: $shadow-card;
  border: 1px solid $border-light;
  margin-bottom: $spacing-xl;
}

.section-header {
  .section-title {
    font-size: $font-size-lg;
    font-weight: $font-weight-semibold;
    color: $text-primary;
    margin: 0 0 4px 0;
  }

  .exchange-rate {
    font-size: 13px;
    color: $text-secondary;
  }
}

.exchange-btn {
  padding: 10px 28px;
  border-radius: $radius-pill;
  background: $gradient-primary;
  color: $text-white;
  border: none;
  font-size: 14px;
  font-weight: $font-weight-medium;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 4px 12px rgba(109, 91, 255, 0.3);

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 6px 16px rgba(109, 91, 255, 0.4);
  }
}

// ========== 菜单网格 ==========
.menu-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: $spacing-xl;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: $spacing-md $spacing-lg;
  background: $bg-card;
  border-radius: $radius-lg;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: $shadow-card;
  border: 1px solid $border-light;

  &:hover {
    background: $bg-hover;
    transform: translateX(4px);

    .menu-arrow {
      opacity: 1;
      transform: translateX(0);
    }
  }

  .menu-icon {
    width: 40px;
    height: 40px;
    border-radius: $radius-md;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    color: $text-white;

    svg {
      width: 20px;
      height: 20px;
    }

    &.icon-history { background: linear-gradient(135deg, #007AFF 0%, #5AC8FA 100%); }
    &.icon-calendar { background: linear-gradient(135deg, #FF6B6B 0%, #FF9F43 100%); }
    &.icon-star { background: linear-gradient(135deg, #FFB800 0%, #FF8A00 100%); }
    &.icon-alert { background: linear-gradient(135deg, #FF3B30 0%, #FF6B6B 100%); }
    &.icon-coin { background: linear-gradient(135deg, #34C759 0%, #30D158 100%); }
    &.icon-clock { background: $gradient-primary; }
    &.icon-message { background: linear-gradient(135deg, #AF52DE 0%, #DA8FFF 100%); }
    &.icon-globe { background: linear-gradient(135deg, #5856D6 0%, #8E8CEC 100%); }
  }

  .menu-label {
    flex: 1;
    font-size: 14px;
    color: $text-primary;
    font-weight: $font-weight-medium;
  }

  .menu-arrow {
    width: 16px;
    height: 16px;
    color: $text-tertiary;
    opacity: 0.5;
    transform: translateX(-4px);
    transition: all 0.2s ease;
  }
}

// ========== 退出按钮 ==========
.logout-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px;
  background: $bg-card;
  border: 1px solid $border-light;
  border-radius: $radius-lg;
  color: $error-color;
  font-size: 15px;
  font-weight: $font-weight-medium;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: $shadow-card;

  svg {
    width: 18px;
    height: 18px;
  }

  &:hover {
    background: rgba(255, 59, 48, 0.06);
    border-color: rgba(255, 59, 48, 0.2);
  }
}

// ========== 弹窗 ==========
.modal-mask {
  position: fixed;
  inset: 0;
  background: $bg-mask;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: $z-modal;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal-content {
  width: 400px;
  max-width: 90vw;
  background: $bg-card;
  border-radius: $radius-lg;
  overflow: hidden;
  box-shadow: $shadow-xl;
  animation: slideUp 0.25s ease;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid $border-light;

  h3 {
    font-size: 16px;
    font-weight: $font-weight-semibold;
    color: $text-primary;
    margin: 0;
  }

  .close-btn {
    width: 28px;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 22px;
    color: $text-tertiary;
    border-radius: 50%;
    cursor: pointer;
    transition: all 0.15s ease;
    border: none;
    background: transparent;

    &:hover {
      background: $bg-hover;
      color: $text-primary;
    }
  }
}

.modal-body {
  padding: 20px;
}

.exchange-info {
  background: $bg-input;
  border-radius: $radius-md;
  padding: 12px 16px;
  margin-bottom: 20px;

  .info-row {
    display: flex;
    justify-content: space-between;
    font-size: 13px;
    color: $text-secondary;
    padding: 4px 0;

    .highlight {
      color: $primary-color;
      font-weight: $font-weight-semibold;
    }
  }
}

.form-item {
  margin-bottom: 16px;

  label {
    display: block;
    font-size: 13px;
    color: $text-secondary;
    margin-bottom: 8px;
  }

  input {
    width: 100%;
    padding: 12px 16px;
    background: $bg-input;
    border: 1px solid $border-color;
    border-radius: $radius-md;
    font-size: 14px;
    color: $text-primary;
    transition: all 0.2s ease;
    box-sizing: border-box;

    &:focus {
      outline: none;
      border-color: $primary-color;
      background: $bg-card;
    }

    &::placeholder {
      color: $text-placeholder;
    }
  }
}

.exchange-preview {
  font-size: 13px;
  color: $text-secondary;

  .highlight {
    color: $success-color;
    font-weight: $font-weight-semibold;
  }
}

.modal-footer {
  display: flex;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid $border-light;
}

.btn {
  flex: 1;
  padding: 10px 24px;
  border-radius: $radius-pill;
  font-size: 14px;
  font-weight: $font-weight-medium;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;

  &.btn-primary {
    background: $gradient-primary;
    color: $text-white;
    box-shadow: 0 4px 12px rgba(109, 91, 255, 0.3);

    &:hover:not(:disabled) {
      transform: translateY(-1px);
      box-shadow: 0 6px 16px rgba(109, 91, 255, 0.4);
    }

    &:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
  }

  &.btn-outline {
    background: $bg-card;
    border: 1px solid $border-color;
    color: $text-secondary;

    &:hover {
      border-color: $primary-light;
      color: $primary-color;
    }
  }
}

// 响应式
@media (max-width: 600px) {
  .asset-overview {
    grid-template-columns: 1fr;
  }

  .menu-grid {
    grid-template-columns: 1fr;
  }
}
</style>
