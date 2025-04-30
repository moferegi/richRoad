<template>
  <view class="content">
    <!-- 顶部搜索栏 -->
    <view class="search-bar">
      <view class="search-box" @click="goToSearch">
        <uni-icons type="search" size="18" color="#999"></uni-icons>
        <text class="placeholder">搜索商品</text>
      </view>
    </view>

    <scroll-view scroll-y="true" class="scroll-Y" @scrolltolower="debouncedLower" refresher-enabled @refresherrefresh="onRefresh" :refresher-triggered="isRefreshing">
      <!-- 轮播图区域 -->
      <view class="swiper-section">
        <home-swiper :list="list"></home-swiper>
      </view>

      <!-- 分类导航 -->
      <view class="category-section">
        <v-tabs
            :tabs="tabs"
            @change="changeTabs"
            v-model="selectTab"
            activeColor="#FF6A6A"
            lineHeight="4"
            lineWidth="20"
            itemWidth="120rpx"
            bold
        >
        </v-tabs>
      </view>

      <!-- 商品展示区 -->
      <view class="goods-section">
        <home-flow :flowData="flowData"></home-flow>
      </view>

      <!-- 底部加载状态 -->
      <view class="loading-status" v-if="flowData.length > 0">
        <view v-if="isBottom" class="no-more">—— 已经到底啦 ——</view>
        <view v-else class="loading">
          <view class="loading-icon"></view>
          <text>加载中...</text>
        </view>
      </view>

      <!-- 空状态 -->
      <view class="empty-state" v-if="flowData.length === 0 && !isLoading">
        <uni-icons type="shop" size="60" color="#ddd"></uni-icons>
        <text>暂无商品</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref } from 'vue';
import homeSwiper from './components/home-swiper.vue'
import homeFlow from './components/home-flow.vue'
import { getCategoryMobile, getGoodList } from '@/api/homePage.js'
import { getBannerList } from '@/api/homePage.js'
import VTabs from "@/components/v-tabs/v-tabs.vue";

// 轮播图相关业务逻辑
const list = ref([])
const initBanner = async () => {
  const res = await getBannerList()
  list.value = res.data.list
}
initBanner()

// 商品相关属性
const params = ref({
  page: 1,
  pageSize: 10,
  categoryID: 0
})

const flowData = ref([])
const isBottom = ref(false)
const isLoading = ref(true)
const isRefreshing = ref(false)

// 跳转到搜索页面
const goToSearch = () => {
  uni.navigateTo({
    url: '/pages/search/index'
  })
}

// 下拉刷新处理
const onRefresh = async () => {
  isRefreshing.value = true
  params.value.page = 1
  flowData.value = []
  isBottom.value = false
  await lower(true)
  setTimeout(() => {
    isRefreshing.value = false
  }, 800)
}

// 获取商品相关业务逻辑
const lower = async (isRefresh = false) => {
  if(isBottom.value && !isRefresh) {
    return
  } else {
    isLoading.value = true
    if(!isRefresh) {
      params.value.page += 1
    }

    try {
      const res = await getGoodList(params.value)
      if (res.code === 0 && res.data.list.length) {
        flowData.value.push(...res.data.list)
        isBottom.value = false
      } else {
        isBottom.value = true
      }
    } catch (error) {
      console.error('获取商品列表失败', error)
    } finally {
      isLoading.value = false
    }
  }
}

const selectTab = ref("全部")

// 切换tabs
const changeTabs = async (index) => {
  params.value.categoryID = tabsMap.value[index]
  params.value.page = 0
  isBottom.value = false
  flowData.value = []
  // 拿到index.id作为categoryID去调用/good/getGoodList接口
  lower()
}

// 分类tabs相关业务逻辑
const gridList = ref([])
const tabsMap = ref({'全部': 0})
const tabs = ref([])
const initCategory = async () => {
  const res = await getCategoryMobile()
  if (res.code === 0 && res.data.length) {
    gridList.value = res.data
    // 赋值给新数组，新数组清洗数据改变desc为name{ name: '电影' },{ name: '科技' }的格式
    res.data.forEach(item => {
      tabsMap.value[item.title] = item.id
      tabs.value.push(item.title)
    })
    tabs.value.unshift('全部')
  }
}
initCategory()
changeTabs("全部")


// 防抖函数
const debounce = (func, delay) => {
  let debounceTimer;
  return function(...args) {
    if (debounceTimer) clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      func.apply(this, args);
    }, delay);
  };
};


// 防抖包装的 lower 方法
const debouncedLower = debounce(()=>lower(false), 300);


</script>

<style lang="scss">
.content {
  width: 100%;
  background-color: #f8f8f8;
  position: relative;

  .search-bar {
    position: sticky;
    top: 0;
    z-index: 100;
    padding: 20rpx 30rpx;
    background-color: #ffffff;
    box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.05);

    .search-box {
      display: flex;
      align-items: center;
      height: 70rpx;
      background-color: #f5f5f5;
      border-radius: 35rpx;
      padding: 0 30rpx;

      .placeholder {
        margin-left: 10rpx;
        font-size: 28rpx;
        color: #999;
      }
    }
  }

  .swiper-section {
    margin: 20rpx 0;
    border-radius: 20rpx;
    overflow: hidden;
    box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);
  }

  .category-section {
    background-color: #ffffff;
    padding: 20rpx 0;
    margin-bottom: 20rpx;
    border-radius: 16rpx;
    box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.03);
  }

  .goods-section {
    padding: 10rpx 0;
  }

  .loading-status {
    text-align: center;
    padding: 30rpx 0;
    color: #999;
    font-size: 24rpx;

    .no-more {
      color: #999;
    }

    .loading {
      display: flex;
      align-items: center;
      justify-content: center;

      .loading-icon {
        width: 30rpx;
        height: 30rpx;
        border: 4rpx solid #f3f3f3;
        border-top: 4rpx solid #FF6A6A;
        border-radius: 50%;
        margin-right: 10rpx;
        animation: spin 1s linear infinite;
      }
    }
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 100rpx 0;
    color: #999;

    text {
      margin-top: 20rpx;
      font-size: 28rpx;
    }
  }
}

.scroll-Y {
  height: 100vh;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>
