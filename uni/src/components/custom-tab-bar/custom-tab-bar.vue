<template>
  <view class="custom-tab-bar">
    <view class="tab-bar-inner">
      <view
        v-for="(tab, index) in tabs"
        :key="tab.pagePath"
        class="tab-item"
        @tap="switchTab(index)"
      >
        <text class="tab-icon" :class="{ on: currentIndex === index }">{{ tab.icon }}</text>
        <text class="tab-label" :class="{ on: currentIndex === index }">{{ tabTexts[index] }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { t as i18nT } from '@/utils/i18n.js'

const langStore = useLangStore()
const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')

const t = (key, fallback) => {
  const text = i18nT(key, locale.value)
  if (text && text !== key) return text
  return fallback || key
}

const tabs = [
  { pagePath: '/pages/learning/home', i18nKey: 'englishTabHome', fallback: '首页', icon: '🏠' },
  { pagePath: '/pages/learning/diary', i18nKey: 'englishTabDiary', fallback: '日记', icon: '📖' },
  { pagePath: '/pages/learning/typing', i18nKey: 'englishTabTyping', fallback: '单词', icon: '⌨' },
  { pagePath: '/pages/learning/profile', i18nKey: 'englishTabMy', fallback: '我的', icon: '👤' }
]

const tabTexts = computed(() => tabs.map(tab => t(tab.i18nKey, tab.fallback)))

const currentIndex = ref(0)

const updateIndex = () => {
  const pages = getCurrentPages()
  if (!pages || pages.length === 0) {
    // 页面栈为空时延迟重试
    setTimeout(() => {
      const retryPages = getCurrentPages()
      if (retryPages && retryPages.length > 0) {
        const route = '/' + retryPages[retryPages.length - 1].route
        const idx = tabs.findIndex(t => t.pagePath === route)
        if (idx >= 0) currentIndex.value = idx
      }
    }, 300)
    return
  }
  const route = '/' + pages[pages.length - 1].route
  const idx = tabs.findIndex(t => t.pagePath === route)
  if (idx >= 0) currentIndex.value = idx
}

const switchTab = (index) => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  const route = '/' + (currentPage ? currentPage.route : '')
  if (tabs[index].pagePath === route) return
  currentIndex.value = index
  uni.switchTab({ url: tabs[index].pagePath })
}

onShow(() => {
  // 首次延迟确保页面栈已更新
  setTimeout(() => {
    updateIndex()
  }, 200)
  // 二次确认，防止页面栈更新延迟
  setTimeout(() => {
    updateIndex()
  }, 500)
})

// 初始加载时延迟执行
setTimeout(() => {
  updateIndex()
}, 300)
</script>

<style>
.custom-tab-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 9999;
  background: #FFFFFF;
  border-top: 1rpx solid rgba(108, 91, 255, 0.12);
  padding-bottom: env(safe-area-inset-bottom);
  box-shadow: 0 -4rpx 20rpx rgba(108, 91, 255, 0.06);
}

.tab-bar-inner {
  display: flex;
  height: 88rpx;
  align-items: center;
}

.tab-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2rpx;
}

.tab-icon {
  font-size: 32rpx;
  line-height: 1;
  filter: grayscale(0.5);
  opacity: 0.5;
  transition: all 0.2s;
}

.tab-icon.on {
  filter: grayscale(0);
  opacity: 1;
}

.tab-label {
  font-size: 18rpx;
  color: #A9AECB;
  font-weight: 500;
  line-height: 1;
}

.tab-label.on {
  color: #6D5BFF;
  font-weight: 700;
}
</style>
