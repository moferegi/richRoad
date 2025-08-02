<template>
  <view class="search-page">
    <!-- 顶部搜索栏 -->
    <view class="status-bar-placeholder"></view>
    <view class="search-header">
      <view class="search-container">
        <!-- 左侧返回首页按钮 -->
        <view class="back-home-icon" @click="goToHome">
          <uni-icons type="home" size="20" color="#fff"></uni-icons>
        </view>

        <!-- 中间搜索框 -->
        <view class="search-input-container">
          <input
            class="search-input"
            v-model="searchKeyword"
            placeholder="请输入您想搜索的商品"
            @confirm="handleSearch"
            confirm-type="search"
            focus
          />
        </view>

        <!-- 右侧搜索按钮 -->
        <view class="search-button" @click="handleSearch">
          <uni-icons type="search" size="20" color="#fff"></uni-icons>
        </view>
      </view>
    </view>

    <!-- 商品列表区域 -->
    <scroll-view scroll-y="true" class="scroll-container" @scrolltolower="handleLoadMore">
      <view class="goods-section">
        <no-pagin-grid-good-list
          :goodsList="goodsList"
          :loading="isLoading"
          :noMore="isBottom"
          @load-more="handleLoadMore"
        />
      </view>

      <!-- 底部加载状态 -->
      <view class="loading-status" v-if="goodsList.length > 0">
        <gva-divider :text="isBottom ? '已经到底啦' : '加载中...'" />
      </view>

      <!-- 空状态 -->
      <view class="empty-state" v-if="goodsList.length === 0 && !isLoading && hasSearched">
        <uni-icons type="search" size="60" color="#ddd" />
        <text>未找到相关商品</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getGoodList } from '@/api/homePage.js'
import noPaginGridGoodList from '@/components/good-list/no-pagin-grid-good-list.vue'

// 搜索相关
const searchKeyword = ref('')
const hasSearched = ref(false)
const sessionFrom = ref('search-page')

// 商品列表相关
const goodsList = ref([])
const isLoading = ref(false)
const isBottom = ref(false)
const params = ref({
  page: 1,
  pageSize: 10,
  categoryID: 0,
  keyword: '',
  recommend: false
})

// 页面加载时获取参数
onLoad((options) => {
  // 优先从本地存储获取搜索关键词
  const storedKeyword = uni.getStorageSync('searchKeyword')
  if (storedKeyword && storedKeyword.trim()) {
    searchKeyword.value = storedKeyword.trim()
    params.value.keyword = storedKeyword.trim()
    // 使用后立即清除存储
    uni.removeStorageSync('searchKeyword')
    handleSearch()
    return
  }

  // 处理URL参数中的搜索关键词
  if (options.keyword && options.keyword.trim()) {
    searchKeyword.value = options.keyword.trim()
    params.value.keyword = options.keyword.trim()
    handleSearch()
  }

  // 处理分类ID参数 - URL参数都是字符串，需要转换为数字类型
  if (options.categoryID) {
    const categoryID = parseInt(options.categoryID, 10)
    // 验证转换后的数字是否有效
    if (!isNaN(categoryID) && categoryID > 0) {
      params.value.categoryID = categoryID
    }
  }
})

// 返回首页
const goToHome = () => {
  uni.switchTab({
    url: '/pages/tabBar/index'
  })
}

// 执行搜索
const handleSearch = async () => {
  if (!searchKeyword.value.trim()) {
    uni.showToast({
      title: '请输入搜索关键词',
      icon: 'none'
    })
    return
  }

  // 重置搜索状态
  params.value.keyword = searchKeyword.value.trim()
  params.value.page = 1
  goodsList.value = []
  isBottom.value = false
  hasSearched.value = true

  await loadGoodsList(true)
}

// 加载商品列表
const loadGoodsList = async (isRefresh = false) => {
  if (isLoading.value) return

  isLoading.value = true

  try {
    // 构建请求参数，确保参数类型正确
    const requestParams = {
      page: params.value.page,
      pageSize: params.value.pageSize,
      keyword: params.value.keyword || ''
    }

    console.log('搜索请求参数:', requestParams)

    const res = await getGoodList(requestParams)

    if (res.code === 0) {
      const listData = res.data.list || []

      if (isRefresh) {
        goodsList.value = listData
      } else {
        goodsList.value.push(...listData)
      }

      // 判断是否还有更多数据
      isBottom.value = listData.length < params.value.pageSize

      // 如果是刷新且没有数据，显示空状态
      if (isRefresh && listData.length === 0) {
        console.log('搜索结果为空')
      }
    } else {
      console.error('搜索API错误:', res)
      uni.showToast({
        title: res.msg || '搜索失败',
        icon: 'none'
      })
      isBottom.value = true
    }
  } catch (error) {
    console.error('搜索商品失败:', error)
    uni.showToast({
      title: '网络错误，请重试',
      icon: 'none'
    })
  } finally {
    isLoading.value = false
  }
}

// 加载更多
const handleLoadMore = async () => {
  if (isBottom.value || isLoading.value || !hasSearched.value) return

  params.value.page += 1
  await loadGoodsList(false)
}
</script>

<style lang="scss" scoped>
.search-page {
  width: 100%;
  height: 100vh;
  background-color: #ffffff;
  display: flex;
  flex-direction: column;
}

// 刘海屏适配
.status-bar-placeholder {
  width: 100%;
  height: var(--status-bar-height, 0px);
  background-color: #ff4c7d;
}

// 微信小程序刘海屏兼容
/* #ifdef MP-WEIXIN */
.status-bar-placeholder {
  height: var(--status-bar-height, 44px);
}
/* #endif */

.search-header {
  background-color: #ff4c7d;
  padding: 16rpx 24rpx;
  position: sticky;
  top: 0;
  z-index: 99;
}

.search-container {
  display: flex;
  align-items: center;
  gap: 20rpx;
  margin-top: 2rem;
}



// 返回首页按钮样式
.back-home-icon {
  width: 60rpx;
  height: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.3s ease;

  &:active {
    background-color: rgba(255, 255, 255, 0.1);
    transform: scale(0.95);
  }
}

.search-input-container {
  flex: 1;
  background-color: #fff;
  height: 72rpx;
  border-radius: 36rpx;
  display: flex;
  align-items: center;
  padding: 0 24rpx;
}

.search-input {
  flex: 1;
  font-size: 28rpx;
  color: #333;

  &::placeholder {
    color: #999;
  }
}

.search-button {
  width: 60rpx;
  height: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.3s ease;

  &:active {
    background-color: rgba(255, 255, 255, 0.1);
    transform: scale(0.95);
  }
}

.scroll-container {
  flex: 1;
  height: 0;
}

.goods-section {
  padding: 10rpx;
}

.loading-status {
  padding: 40rpx 0;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;

  text {
    font-size: 28rpx;
    color: #999;
    margin-top: 24rpx;
  }
}
</style>
