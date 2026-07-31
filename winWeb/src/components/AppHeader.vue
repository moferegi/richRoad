<script setup>
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/pinia'
import { computed } from 'vue'

const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()

const isLoggedIn = computed(() => userStore.isLoggedIn)
const username = computed(() => userStore.username)

function goHome() {
  router.push('/home')
}

function goProfile() {
  router.push('/learning/profile')
}

function handleLogout() {
  userStore.loginOut()
  router.push('/login')
}

function goLogin() {
  router.push('/login')
}
</script>

<template>
  <header class="app-header">
    <div class="header-inner">
      <div class="header-left" @click="goHome">
        <div class="logo">
          <div class="logo-icon">R</div>
          <span class="logo-text">RichRoad</span>
        </div>
      </div>

      <div class="header-center">
        <div class="search-box">
          <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"></circle>
            <path d="m21 21-4.35-4.35"></path>
          </svg>
          <input type="text" :placeholder="t('common.search')" />
        </div>
      </div>

      <div class="header-right">
        <div v-if="isLoggedIn" class="user-section">
          <div class="avatar" @click="goProfile">
            <span>{{ username.charAt(0).toUpperCase() }}</span>
          </div>
          <div class="user-info">
            <div class="user-name">{{ username }}</div>
            <button class="logout-btn" @click="handleLogout">
              {{ t('common.logout') }}
            </button>
          </div>
        </div>
        <button v-else class="btn btn-primary btn-sm" @click="goLogin">
          {{ t('common.login') }}
        </button>
      </div>
    </div>
  </header>
</template>

<style lang="scss" scoped>
.app-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: $header-height;
  background: $bg-header;
  border-bottom: 1px solid $border-light;
  box-shadow: $shadow-sm;
  z-index: $z-sticky;
}

.header-inner {
  display: flex;
  align-items: center;
  height: 100%;
  padding: 0 $spacing-xl;
  max-width: 100%;
}

.header-left {
  display: flex;
  align-items: center;
  width: $sidebar-width - $spacing-xl;
  cursor: pointer;

  .logo {
    display: flex;
    align-items: center;
    gap: $spacing-sm;

    .logo-icon {
      width: 36px;
      height: 36px;
      background: $gradient-primary;
      border-radius: $radius-md;
      display: flex;
      align-items: center;
      justify-content: center;
      color: $text-white;
      font-weight: $font-weight-bold;
      font-size: $font-size-lg;
    }

    .logo-text {
      font-size: $font-size-lg;
      font-weight: $font-weight-bold;
      background: $gradient-primary;
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }
  }
}

.header-center {
  flex: 1;
  display: flex;
  justify-content: center;
  padding: 0 $spacing-xxl;

  .search-box {
    position: relative;
    width: 100%;
    max-width: 480px;

    .search-icon {
      position: absolute;
      left: 14px;
      top: 50%;
      transform: translateY(-50%);
      width: 18px;
      height: 18px;
      color: $text-tertiary;
    }

    input {
      width: 100%;
      height: 40px;
      padding: 0 16px 0 40px;
      background: $bg-input;
      border: 1px solid transparent;
      border-radius: $radius-pill;
      font-size: $font-size-sm;
      color: $text-primary;
      transition: all $transition-fast;

      &:focus {
        border-color: $primary-color;
        background: $bg-card;
        box-shadow: 0 0 0 3px rgba(109, 91, 255, 0.1);
      }

      &::placeholder {
        color: $text-placeholder;
      }
    }
  }
}

.header-right {
  display: flex;
  align-items: center;
  gap: $spacing-md;

  .user-section {
    display: flex;
    align-items: center;
    gap: $spacing-md;
    cursor: pointer;

    .avatar {
      width: 38px;
      height: 38px;
      background: $gradient-primary;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      color: $text-white;
      font-weight: $font-weight-semibold;
      font-size: $font-size-md;
      transition: transform $transition-fast;

      &:hover {
        transform: scale(1.05);
      }
    }

    .user-info {
      .user-name {
        font-size: $font-size-sm;
        font-weight: $font-weight-medium;
        color: $text-primary;
      }

      .logout-btn {
        font-size: $font-size-xs;
        color: $text-tertiary;
        padding: 0;
        margin-top: 2px;
        transition: color $transition-fast;

        &:hover {
          color: $error-color;
        }
      }
    }
  }
}
</style>
