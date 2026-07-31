<script setup>
import { useI18n } from 'vue-i18n'
import { useRouter, useRoute } from 'vue-router'
import { computed } from 'vue'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const navItems = [
  { path: '/home', name: 'home', icon: 'home' },
  { path: '/learning/diary', name: 'diary', icon: 'diary' },
  { path: '/learning/typing', name: 'typing', icon: 'typing' },
  { path: '/learning/profile', name: 'profile', icon: 'profile' }
]

const isActive = (path) => {
  if (path === '/home') return route.path === '/home'
  return route.path.startsWith(path.split('/').slice(0, 3).join('/'))
}

function navigate(path) {
  router.push(path)
}

// 二级菜单 - 个人中心子菜单
const profileSubItems = computed(() => {
  if (!route.path.startsWith('/learning/profile') &&
      !route.path.startsWith('/learning/collections') &&
      !route.path.startsWith('/learning/watch-history') &&
      !route.path.startsWith('/learning/checkin-record') &&
      !route.path.startsWith('/learning/error-log') &&
      !route.path.startsWith('/learning/point-history') &&
      !route.path.startsWith('/learning/free-time-history')) {
    return []
  }
  return [
    { path: '/learning/profile', label: 'profile', icon: 'profile' },
    { path: '/learning/collections', label: 'collections', icon: 'star' },
    { path: '/learning/watch-history', label: 'watchHistory', icon: 'history' },
    { path: '/learning/checkin-record', label: 'checkinRecord', icon: 'check' },
    { path: '/learning/error-log', label: 'errorLog', icon: 'error' },
    { path: '/learning/point-history', label: 'pointHistory', icon: 'coin' },
    { path: '/learning/free-time-history', label: 'freeTimeHistory', icon: 'clock' }
  ]
})
</script>

<template>
  <aside class="app-sidebar">
    <nav class="sidebar-nav">
      <div class="nav-section">
        <div
          v-for="item in navItems"
          :key="item.path"
          class="nav-item"
          :class="{ active: isActive(item.path) }"
          @click="navigate(item.path)"
        >
          <div class="nav-icon">
            <svg v-if="item.icon === 'home'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
              <polyline points="9 22 9 12 15 12 15 22"></polyline>
            </svg>
            <svg v-else-if="item.icon === 'diary'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"></path>
              <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"></path>
            </svg>
            <svg v-else-if="item.icon === 'typing'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="2" y="4" width="20" height="16" rx="2"></rect>
              <path d="M6 8h.01M10 8h.01M14 8h.01M18 8h.01M8 12h.01M12 12h.01M16 12h.01M6 16h12"></path>
            </svg>
            <svg v-else-if="item.icon === 'profile'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
              <circle cx="12" cy="7" r="4"></circle>
            </svg>
          </div>
          <span class="nav-label">{{ t('nav.' + item.name) }}</span>
        </div>
      </div>

      <!-- 个人中心子菜单 -->
      <div v-if="profileSubItems.length > 0" class="nav-sub-section">
        <div
          v-for="item in profileSubItems"
          :key="item.path"
          class="nav-sub-item"
          :class="{ active: route.path === item.path }"
          @click="navigate(item.path)"
        >
          <div class="nav-sub-icon">
            <svg v-if="item.icon === 'profile'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
              <circle cx="12" cy="7" r="4"></circle>
            </svg>
            <svg v-else-if="item.icon === 'star'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
            </svg>
            <svg v-else-if="item.icon === 'history'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <polyline points="12 6 12 12 16 14"></polyline>
            </svg>
            <svg v-else-if="item.icon === 'check'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
              <polyline points="22 4 12 14.01 9 11.01"></polyline>
            </svg>
            <svg v-else-if="item.icon === 'error'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="15" y1="9" x2="9" y2="15"></line>
              <line x1="9" y1="9" x2="15" y2="15"></line>
            </svg>
            <svg v-else-if="item.icon === 'coin'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <path d="M12 6v12M9 9h4.5a2.5 2.5 0 0 1 0 5H9h5"></path>
            </svg>
            <svg v-else-if="item.icon === 'clock'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <polyline points="12 6 12 12 16 14"></polyline>
            </svg>
          </div>
          <span class="nav-sub-label">{{ t('profile.' + item.label) }}</span>
        </div>
      </div>
    </nav>

    <!-- 底部装饰 -->
    <div class="sidebar-footer">
      <div class="footer-deco"></div>
    </div>
  </aside>
</template>

<style lang="scss" scoped>
.app-sidebar {
  position: fixed;
  top: $header-height;
  left: 0;
  bottom: 0;
  width: $sidebar-width;
  background: $bg-sidebar;
  border-right: 1px solid $border-light;
  display: flex;
  flex-direction: column;
  z-index: $z-dropdown;
}

.sidebar-nav {
  flex: 1;
  padding: $spacing-lg 0;
  overflow-y: auto;
}

.nav-section {
  margin-bottom: $spacing-md;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: $spacing-md;
  height: 48px;
  padding: 0 $spacing-xl;
  margin: 2px $spacing-md;
  border-radius: $radius-md;
  cursor: pointer;
  transition: all $transition-fast;
  color: $text-secondary;

  &:hover {
    background: $bg-hover;
    color: $text-primary;
  }

  &.active {
    background: $bg-active;
    color: $primary-color;
    font-weight: $font-weight-medium;

    .nav-icon {
      color: $primary-color;
    }
  }

  .nav-icon {
    width: 20px;
    height: 20px;
    color: inherit;
    display: flex;
    align-items: center;
    justify-content: center;

    svg {
      width: 20px;
      height: 20px;
    }
  }

  .nav-label {
    font-size: $font-size-base;
    flex: 1;
  }
}

.nav-sub-section {
  padding: $spacing-sm 0;
  margin: 0 $spacing-md;
  border-left: 2px solid $border-light;
  margin-left: $spacing-lg;
  padding-left: $spacing-sm;
}

.nav-sub-item {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  height: 40px;
  padding: 0 $spacing-md;
  border-radius: $radius-sm;
  cursor: pointer;
  transition: all $transition-fast;
  color: $text-tertiary;

  &:hover {
    background: $bg-hover;
    color: $text-primary;
  }

  &.active {
    color: $primary-color;
    font-weight: $font-weight-medium;
    background: rgba(109, 91, 255, 0.08);
  }

  .nav-sub-icon {
    width: 16px;
    height: 16px;
    display: flex;
    align-items: center;
    justify-content: center;

    svg {
      width: 16px;
      height: 16px;
    }
  }

  .nav-sub-label {
    font-size: $font-size-sm;
  }
}

.sidebar-footer {
  padding: $spacing-md;
  position: relative;

  .footer-deco {
    height: 80px;
    margin: 0 $spacing-sm;
    background: $gradient-card;
    border-radius: $radius-lg;
    position: relative;
    overflow: hidden;

    &::before {
      content: '';
      position: absolute;
      top: -30px;
      right: -20px;
      width: 80px;
      height: 80px;
      background: $gradient-primary;
      border-radius: 50%;
      opacity: 0.15;
    }
  }
}
</style>
