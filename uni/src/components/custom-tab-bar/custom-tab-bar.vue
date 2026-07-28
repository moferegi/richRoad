<template>
  <view class="custom-tab-bar">
    <view class="tab-bar-inner">
      <view
        v-for="(tab, index) in tabs"
        :key="tab.pagePath"
        class="tab-item"
        :class="{ on: currentIndex === index }"
        @tap="switchTab(index)"
      >
        <image
          class="tab-icon"
          :src="currentIndex === index ? tab.iconActive : tab.icon"
          mode="aspectFit"
        />
        <text class="tab-label">{{ tabTexts[index] }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
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
  {
    pagePath: '/pages/learning/home',
    i18nKey: 'englishTabHome',
    fallback: '首页',
    icon: '/static/images/tabBar/learn-home.svg',
    iconActive: '/static/images/tabBar/learn-home-white.svg'
  },
  {
    pagePath: '/pages/learning/diary',
    i18nKey: 'englishTabDiary',
    fallback: '日记',
    icon: '/static/images/tabBar/learn-diary.svg',
    iconActive: '/static/images/tabBar/learn-diary-white.svg'
  },
  {
    pagePath: '/pages/learning/typing',
    i18nKey: 'englishTabTyping',
    fallback: '单词',
    icon: '/static/images/tabBar/learn-word.svg',
    iconActive: '/static/images/tabBar/learn-word-white.svg'
  },
  {
    pagePath: '/pages/learning/profile',
    i18nKey: 'englishTabMy',
    fallback: '我的',
    icon: '/static/images/tabBar/learn-my.svg',
    iconActive: '/static/images/tabBar/learn-my-white.svg'
  }
]

const tabTexts = computed(() => tabs.map(tab => t(tab.i18nKey, tab.fallback)))

const currentIndex = ref(0)

const updateIndex = () => {
  const pages = getCurrentPages()
  if (!pages || pages.length === 0) return
  const route = '/' + pages[pages.length - 1].route
  const idx = tabs.findIndex(t => t.pagePath === route)
  if (idx >= 0 && idx !== currentIndex.value) {
    currentIndex.value = idx
  }
}

const switchTab = (index) => {
  if (index === currentIndex.value) return
  // 点击立即更新，激活响应无延迟
  currentIndex.value = index
  uni.switchTab({ url: tabs[index].pagePath })
}

onShow(() => {
  setTimeout(() => {
    updateIndex()
  }, 200)
  setTimeout(() => {
    updateIndex()
  }, 500)
})

onMounted(() => {
  setTimeout(() => {
    updateIndex()
  }, 300)
})
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
  height: 120rpx;
  align-items: center;
  justify-content: space-around;
  padding: 0 12rpx;
}

.tab-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4rpx;
  padding: 10rpx 32rpx;
  border-radius: 36rpx;
  background: transparent;
  transition: background 0.2s ease, box-shadow 0.2s ease;
}

.tab-item.on {
  background: #6D5BFF;
  box-shadow: 0 6rpx 16rpx rgba(109, 91, 255, 0.35);
}

.tab-icon {
  width: 44rpx;
  height: 44rpx;
  line-height: 1;
  transition: transform 0.2s ease;
}

.tab-item.on .tab-icon {
  transform: scale(1.05);
}

.tab-label {
  font-size: 22rpx;
  color: #A9AECB;
  font-weight: 500;
  line-height: 1;
  transition: color 0.2s ease, font-weight 0.2s ease;
}

.tab-item.on .tab-label {
  color: #FFFFFF;
  font-weight: 600;
}
</style>
