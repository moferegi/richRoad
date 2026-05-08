<template>
  <view class="nf-search">
    <!-- 状态栏占位 -->
    <view class="nf-search-status"></view>

    <!-- 搜索栏 -->
    <view class="nf-search-header">
      <view class="nf-search-row">
        <!-- 返回首页 -->
        <view class="nf-search-back" @tap="goToHome">
          <uni-icons type="home" size="22" color="#fff" />
        </view>
        <!-- 搜索框 -->
        <view class="nf-search-input-wrap">
          <uni-icons type="search" size="16" color="rgba(255,255,255,0.4)" />
          <input
            class="nf-search-input"
            v-model="searchKeyword"
            :placeholder="$t('searchPlaceholder')"
            placeholder-class="nf-search-ph"
            @confirm="handleSearch"
            confirm-type="search"
            focus
          />
          <view v-if="searchKeyword" class="nf-search-clear" @tap="searchKeyword = ''">
            <uni-icons type="closeempty" size="14" color="rgba(255,255,255,0.4)" />
          </view>
        </view>
        <!-- 搜索按钮 -->
        <view class="nf-search-btn" @tap="handleSearch">
          <text class="nf-search-btn-text">{{ $t('searchBtn') }}</text>
        </view>
      </view>
    </view>

    <!-- 内容区 -->
    <scroll-view scroll-y="true" class="nf-search-scroll" @scrolltolower="handleLoadMore">

      <!-- 搜索结果提示条 -->
      <view class="nf-result-tip" v-if="hasSearched && goodsList.length > 0">
        <text class="nf-result-tip-text">{{ searchResultLabel }}</text>
      </view>

      <!-- 商品列表 -->
      <view class="nf-goods-wrap" v-if="goodsList.length > 0">
        <no-pagin-grid-good-list
          :goodsList="goodsList"
          :loading="isLoading"
          :noMore="isBottom"
          @load-more="handleLoadMore"
        />
      </view>

      <!-- 已加载全部 -->
      <view class="nf-load-all" v-if="goodsList.length > 0 && isBottom">
        <text class="nf-load-all-text">{{ $t('searchLoadAll') }}</text>
      </view>

      <!-- 空状态 -->
      <view class="nf-empty" v-if="goodsList.length === 0 && !isLoading && hasSearched">
        <view class="nf-empty-icon-wrap"><uni-icons type="search" size="72" color="rgba(255,255,255,0.12)" /></view>
        <text class="nf-empty-title">{{ $t('searchNoResult') }}</text>
        <text class="nf-empty-desc">{{ $t('searchNoResultTip') }}</text>
      </view>

      <!-- 初始引导 -->
      <view class="nf-init" v-if="!hasSearched && goodsList.length === 0 && !isLoading">
        <view class="nf-init-icon-wrap"><uni-icons type="search" size="72" color="rgba(255,255,255,0.12)" /></view>
        <text class="nf-init-title">{{ $t('searchInitTitle') }}</text>
        <text class="nf-init-desc">{{ $t('searchInitTip') }}</text>
      </view>

      <!-- 底部安全区 -->
      <view style="height: 60rpx;"></view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getGoodList } from '@/api/homePage.js'
import noPaginGridGoodList from '@/components/good-list/no-pagin-grid-good-list.vue'
import { resolveApiMessage } from '@/utils/i18n.js'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)
const searchResultLabel = computed(() => $t.value('searchResultCount').replace('{' + '{n}}', goodsList.value.length))

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
    uni.showToast({ title: $t.value('searchEmpty'), icon: 'none' })
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
      uni.showToast({ title: resolveApiMessage(res.msg, 'searchNoResult'), icon: 'none' })
      isBottom.value = true
    }
  } catch (error) {
    console.error('搜索商品失败:', error)
    uni.showToast({ title: $t.value('playerLoadFail'), icon: 'none' })
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
page {
  background-color: #000;
}

.nf-search {
  width: 100%;
  height: 100vh;
  background: #000;
  display: flex;
  flex-direction: column;
}

/* 状态栏占位 */
.nf-search-status {
  width: 100%;
  height: var(--status-bar-height, 0px);
  flex-shrink: 0;
}

/* 搜索栏 */
.nf-search-header {
  flex-shrink: 0;
  padding: 16rpx 24rpx 20rpx;
  background: rgba(0, 0, 0, 0.92);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
}

.nf-search-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

/* 首页按钮 */
.nf-search-back {
  flex-shrink: 0;
  width: 68rpx;
  height: 68rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);

  &:active {
    background: rgba(255, 255, 255, 0.16);
    transform: scale(0.92);
  }
}

/* 输入框容器 */
.nf-search-input-wrap {
  flex: 1;
  height: 68rpx;
  display: flex;
  align-items: center;
  gap: 12rpx;
  background: rgba(255, 255, 255, 0.08);
  border: 1rpx solid rgba(255, 255, 255, 0.1);
  border-radius: 8rpx;
  padding: 0 16rpx;
  transition: border-color 0.2s;

  &:focus-within {
    border-color: rgba(229, 9, 20, 0.45);
  }
}

.nf-search-input {
  flex: 1;
  font-size: 28rpx;
  color: #fff;
  background: transparent;
}

.nf-search-ph {
  color: rgba(255, 255, 255, 0.3) !important;
}

.nf-search-clear {
  flex-shrink: 0;
  width: 44rpx;
  height: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;

  &:active { opacity: 0.6; transform: scale(0.88); }
}

/* 搜索按钮 */
.nf-search-btn {
  flex-shrink: 0;
  width: 104rpx;
  height: 68rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #e50914, #b20710);
  border-radius: 8rpx;
  box-shadow: 0 4rpx 12rpx rgba(229, 9, 20, 0.3);

  &:active {
    transform: scale(0.95);
    box-shadow: 0 2rpx 6rpx rgba(229, 9, 20, 0.4);
  }
}

.nf-search-btn-text {
  color: #fff;
  font-size: 26rpx;
  font-weight: 600;
  letter-spacing: 2rpx;
}

/* 滚动区 */
.nf-search-scroll {
  flex: 1;
  overflow-y: auto;
  scrollbar-width: none;
  -ms-overflow-style: none;

  &::-webkit-scrollbar { display: none; width: 0; height: 0; }
}

/* 结果数量提示 */
.nf-result-tip {
  padding: 20rpx 32rpx 12rpx;
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.05);

  .nf-result-tip-text {
    font-size: 24rpx;
    color: rgba(255, 255, 255, 0.45);
    letter-spacing: 0.5rpx;
  }
}

/* 商品列表 */
.nf-goods-wrap {
  padding: 20rpx 12rpx;
}

/* 已加载全部 */
.nf-load-all {
  padding: 40rpx 0;
  text-align: center;

  .nf-load-all-text {
    font-size: 24rpx;
    color: rgba(255, 255, 255, 0.3);
    letter-spacing: 1rpx;
  }
}

/* 空状态 & 初始引导 */
.nf-empty,
.nf-init {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 48rpx;
  min-height: 60vh;
}

.nf-empty-icon-wrap,
.nf-init-icon-wrap {
  margin-bottom: 36rpx;
}

.nf-empty-title,
.nf-init-title {
  font-size: 32rpx;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.55);
  letter-spacing: 1rpx;
  margin-bottom: 12rpx;
  text-align: center;
}

.nf-empty-desc,
.nf-init-desc {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.28);
  text-align: center;
  line-height: 1.6;
}
</style>
