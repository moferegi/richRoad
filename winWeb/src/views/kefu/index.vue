<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { getKefuList, getCsConfig } from '@/api/kefu'

const { t } = useI18n()

const kefuList = ref([])
const csConfig = ref(null)
const loading = ref(false)
const modal = ref({ visible: false, type: '', title: '', content: '' })
const copied = ref(false)

async function loadData() {
  loading.value = true
  try {
    const [kefuRes, csRes] = await Promise.allSettled([
      getKefuList(),
      getCsConfig()
    ])
    if (kefuRes.status === 'fulfilled') {
      const res = kefuRes.value
      kefuList.value = res.data?.list || res.data || []
    }
    if (csRes.status === 'fulfilled') {
      const res = csRes.value
      csConfig.value = res.data || null
    }
  } catch (e) {
    console.warn('load kefu failed', e)
  } finally {
    loading.value = false
  }
}

function isOnline(item) {
  return item.online !== false && item.online !== 0 && item.status !== 0
}

function handleContact(item) {
  if (item.link) {
    window.open(item.link, '_blank')
    return
  }
  if (item.qrCode) {
    modal.value = {
      visible: true,
      type: 'qrcode',
      title: item.name || t('kefu.platformService'),
      content: item.qrCode
    }
    return
  }
  if (item.contactId) {
    modal.value = {
      visible: true,
      type: 'contact',
      title: item.name || t('kefu.platformService'),
      content: item.contactId
    }
    return
  }
}

function handlePlatformContact() {
  if (!csConfig.value) return
  const cfg = csConfig.value
  if (cfg.link) {
    window.open(cfg.link, '_blank')
    return
  }
  if (cfg.qrCode) {
    modal.value = {
      visible: true,
      type: 'qrcode',
      title: t('kefu.platformService'),
      content: cfg.qrCode
    }
    return
  }
  if (cfg.contactId) {
    modal.value = {
      visible: true,
      type: 'contact',
      title: t('kefu.platformService'),
      content: cfg.contactId
    }
    return
  }
}

function closeModal() {
  modal.value = { visible: false, type: '', title: '', content: '' }
  copied.value = false
}

async function copyContact() {
  try {
    await navigator.clipboard.writeText(modal.value.content)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (e) {
    const textarea = document.createElement('textarea')
    textarea.value = modal.value.content
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  }
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="kefu-page">
    <div class="kefu-content">
      <div class="page-header">
        <h1 class="page-title">{{ t('kefu.title') }}</h1>
        <p class="page-subtitle">{{ t('kefu.platformService') }}</p>
      </div>

      <div v-if="csConfig && csConfig.enabled !== false" class="platform-card" @click="handlePlatformContact">
        <div class="platform-left">
          <div class="platform-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
            </svg>
          </div>
          <div class="platform-info">
            <h3 class="platform-title">{{ t('kefu.platformService') }}</h3>
            <p class="platform-desc">7×24 小时在线为您服务</p>
          </div>
        </div>
        <div class="platform-right">
          <button class="btn btn-primary">{{ t('kefu.contact') }}</button>
        </div>
      </div>

      <div class="section-header">
        <h3 class="section-title">{{ t('kefu.externalService') }}</h3>
      </div>

      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <span>{{ t('common.loading') }}</span>
      </div>

      <div v-else-if="kefuList.length === 0" class="empty-state">
        <div class="empty-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
          </svg>
        </div>
        <p>{{ t('common.empty') }}</p>
      </div>

      <div v-else class="kefu-grid">
        <div
          v-for="item in kefuList"
          :key="item.ID || item.id || item.contactId"
          class="kefu-card"
        >
          <div class="card-avatar-wrap">
            <img
              v-if="item.avatar || item.avatarUrl || item.headImg"
              :src="item.avatar || item.avatarUrl || item.headImg"
              :alt="item.name"
              class="card-avatar"
            />
            <div v-else class="card-avatar avatar-placeholder">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                <circle cx="12" cy="7" r="4"></circle>
              </svg>
            </div>
            <span class="status-dot" :class="{ online: isOnline(item), offline: !isOnline(item) }"></span>
          </div>
          <h4 class="card-name">{{ item.name || item.nickname || '客服' }}</h4>
          <p class="card-status" :class="{ online: isOnline(item) }">
            {{ isOnline(item) ? t('kefu.online') : t('kefu.offline') }}
          </p>
          <button class="btn btn-primary contact-btn" @click="handleContact(item)">
            {{ t('kefu.contact') }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="modal.visible" class="modal-mask" @click="closeModal">
      <div class="modal-content kefu-modal" @click.stop>
        <button class="modal-close" @click="closeModal">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
        <h3 class="modal-title">{{ modal.title }}</h3>

        <div v-if="modal.type === 'qrcode'" class="qrcode-wrap">
          <img :src="modal.content" :alt="modal.title" class="qrcode-img" />
        </div>

        <div v-else-if="modal.type === 'contact'" class="contact-wrap">
          <p class="contact-label">{{ t('kefu.platformService') }}</p>
          <div class="contact-id-box">
            <span class="contact-id">{{ modal.content }}</span>
            <button class="copy-btn" :class="{ copied }" @click="copyContact">
              <svg v-if="!copied" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
              </svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="20 6 9 17 4 12"></polyline>
              </svg>
              <span>{{ copied ? t('kefu.copied') : t('kefu.copy') }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use "@/assets/styles/variables.scss" as *;

.kefu-page {
  min-height: 100%;
}

.kefu-content {
  max-width: $content-max-width;
  margin: 0 auto;
  padding: $spacing-xl;
}

.page-header {
  margin-bottom: $spacing-xl;

  .page-title {
    font-size: $font-size-title;
    font-weight: $font-weight-bold;
    color: $text-primary;
    margin: 0 0 6px 0;
  }

  .page-subtitle {
    font-size: $font-size-base;
    color: $text-secondary;
    margin: 0;
  }
}

.platform-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: $spacing-lg;
  background: $gradient-card;
  border-radius: $radius-lg;
  box-shadow: $shadow-card;
  border: 1px solid $border-light;
  cursor: pointer;
  transition: all 0.25s ease;
  margin-bottom: $spacing-xl;

  &:hover {
    box-shadow: $shadow-md;
    transform: translateY(-2px);
  }

  .platform-left {
    display: flex;
    align-items: center;
    gap: $spacing-md;
  }

  .platform-icon {
    width: 56px;
    height: 56px;
    border-radius: $radius-md;
    background: $gradient-primary;
    color: $text-white;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 12px rgba(109, 91, 255, 0.3);

    svg {
      width: 28px;
      height: 28px;
    }
  }

  .platform-title {
    font-size: $font-size-lg;
    font-weight: $font-weight-bold;
    color: $text-primary;
    margin: 0 0 4px 0;
  }

  .platform-desc {
    font-size: $font-size-sm;
    color: $text-secondary;
    margin: 0;
  }
}

.section-header {
  margin-bottom: $spacing-md;

  .section-title {
    font-size: $font-size-lg;
    font-weight: $font-weight-semibold;
    color: $text-primary;
    margin: 0;
  }
}

.kefu-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: $spacing-md;
}

.kefu-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: $spacing-lg $spacing-md;
  background: $bg-card;
  border-radius: $radius-lg;
  box-shadow: $shadow-card;
  border: 1px solid $border-light;
  transition: all 0.25s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: $shadow-md;
  }
}

.card-avatar-wrap {
  position: relative;
  margin-bottom: $spacing-sm;
}

.card-avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  object-fit: cover;
  background: $bg-input;

  &.avatar-placeholder {
    display: flex;
    align-items: center;
    justify-content: center;
    color: $primary-light;
    background: rgba(109, 91, 255, 0.1);

    svg {
      width: 32px;
      height: 32px;
    }
  }
}

.status-dot {
  position: absolute;
  bottom: 2px;
  right: 2px;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  border: 2px solid $bg-card;

  &.online {
    background: $success-color;
  }

  &.offline {
    background: $text-tertiary;
  }
}

.card-name {
  font-size: $font-size-md;
  font-weight: $font-weight-semibold;
  color: $text-primary;
  margin: 0 0 4px 0;
}

.card-status {
  font-size: $font-size-sm;
  color: $text-tertiary;
  margin: 0 0 $spacing-md 0;

  &.online {
    color: $success-color;
  }
}

.contact-btn {
  width: 100%;
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 8px 20px;
  border-radius: $radius-pill;
  font-size: $font-size-sm;
  font-weight: $font-weight-medium;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;

  &.btn-primary {
    background: $gradient-primary;
    color: $text-white;
    box-shadow: 0 2px 8px rgba(109, 91, 255, 0.3);

    &:hover {
      opacity: 0.9;
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(109, 91, 255, 0.4);
    }
  }
}

.loading-state,
.empty-state {
  text-align: center;
  padding: $spacing-xxl 0;
  color: $text-tertiary;

  .loading-spinner {
    width: 28px;
    height: 28px;
    margin: 0 auto $spacing-sm;
    border: 3px solid $border-color;
    border-top-color: $primary-color;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  .empty-icon {
    width: 72px;
    height: 72px;
    margin: 0 auto $spacing-sm;
    color: $border-primary;

    svg {
      width: 100%;
      height: 100%;
    }
  }

  p {
    font-size: $font-size-base;
    margin: 0;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

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
  position: relative;
  background: $bg-card;
  border-radius: $radius-lg;
  box-shadow: $shadow-xl;
  padding: $spacing-lg;
  animation: scaleIn 0.2s ease;
}

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.kefu-modal {
  width: 320px;
  text-align: center;
}

.modal-close {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  color: $text-tertiary;
  cursor: pointer;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;

  svg {
    width: 18px;
    height: 18px;
  }

  &:hover {
    background: $bg-hover;
    color: $text-primary;
  }
}

.modal-title {
  font-size: $font-size-lg;
  font-weight: $font-weight-bold;
  color: $text-primary;
  margin: 0 0 $spacing-md 0;
}

.qrcode-wrap {
  padding: $spacing-md 0;

  .qrcode-img {
    width: 200px;
    height: 200px;
    object-fit: contain;
    border-radius: $radius-md;
  }
}

.contact-wrap {
  padding: $spacing-md 0;
}

.contact-label {
  font-size: $font-size-sm;
  color: $text-secondary;
  margin: 0 0 $spacing-sm 0;
}

.contact-id-box {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: $spacing-sm;
  padding: $spacing-sm $spacing-md;
  background: $bg-input;
  border-radius: $radius-md;
  margin-bottom: $spacing-sm;
}

.contact-id {
  flex: 1;
  font-size: $font-size-md;
  font-weight: $font-weight-medium;
  color: $text-primary;
  text-align: left;
  word-break: break-all;
}

.copy-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  background: $gradient-primary;
  color: $text-white;
  border: none;
  border-radius: $radius-pill;
  font-size: $font-size-sm;
  cursor: pointer;
  flex-shrink: 0;
  transition: all 0.2s ease;

  svg {
    width: 14px;
    height: 14px;
  }

  &:hover {
    opacity: 0.9;
  }

  &.copied {
    background: $success-color;
  }
}
</style>
