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
        <text class="tab-label" :class="{ on: currentIndex === index }">{{ tab.text }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'

const tabs = [
  { pagePath: '/pages/learning/home', text: '首页', icon: '🏠' },
  { pagePath: '/pages/learning/typing', text: '跟打', icon: '⌨' },
  { pagePath: '/pages/learning/profile', text: '我的', icon: '👤' }
]

const currentIndex = ref(0)

const updateIndex = () => {
  const pages = getCurrentPages()
  if (!pages || pages.length === 0) return
  const route = '/' + pages[pages.length - 1].route
  const idx = tabs.findIndex(t => t.pagePath === route)
  if (idx >= 0) currentIndex.value = idx
}

const switchTab = (index) => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  const route = '/' + (currentPage ? currentPage.route : '')
  if (tabs[index].pagePath === route) return
  uni.switchTab({ url: tabs[index].pagePath })
}

onShow(() => {
  updateIndex()
})

updateIndex()
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
