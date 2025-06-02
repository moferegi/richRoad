<template>
  <view class="content">
    <!-- 顶部搜索栏 - 修改刘海屏适配 -->
    <view class="status-bar-placeholder"></view>
    <view class="status-bar">
      <view class="test">
        <!-- 左侧菜单图标 -->
        <view class="scan-icon">
          <uni-icons type="bars" size="22" color="#fff"></uni-icons>
        </view>

        <!-- 中间搜索框 -->
        <view class="search-bar" @click="goToSearch">
          <view class="search-icon">
            <uni-icons type="search" size="18" color="#999"></uni-icons>
          </view>
          <text class="search-placeholder">请输入地址 如：大钟寺</text>
        </view>

        <!-- 右侧消息 -->
        <view class="message-icon">
          <uni-icons type="chat" size="22" color="#fff"></uni-icons>
        </view>
      </view>
    </view>
    <scroll-view scroll-y="true" class="scroll-Y" @scrolltolower="debouncedLower" @refresherrefresh="onRefresh" >
      <!-- 轮播图区域 -->
      <swpiers :lists="list"></swpiers>

      <!-- 分类导航 -->
      <categories :categoriesData="gridList"></categories>

      <!-- 限时秒杀区域 -->
      <seckilling :productData="products"></seckilling>
      <!-- 商品展示区 -->
      <view class="goods-section">
        <view class="goods-section">
          <noPaginRowGoodList
            :goodsList="flowData"
            :current-page="params.page"
            :page-size="params.pageSize"
            :total="totalCount"
            :loading="isLoading"
            :load-offset="200"
            @load-more="handleAutoLoadMore"
          ></noPaginRowGoodList>
        </view>
      </view>

      <!-- 底部加载状态 -->
      <view class="loading-status" v-if="flowData.length > 0">
        <gva-divider :text="isBottom?'已经到底啦':'加载中...'"></gva-divider>
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
import  { ref } from 'vue';
import { getCategoryMobile, getGoodList } from '@/api/homePage.js'
import { getBannerList } from '@/api/homePage.js'
import noPaginRowGoodList from '@/components/good-list/no-pagin-row-good-list.vue'
import swpiers from './components/swiper.vue'
import categories from './components/categories.vue'
import seckilling from './components/seckilling.vue';

// 轮播图相关业务逻辑
const list = ref([])
const initBanner = async () => {
  const res = await getBannerList()
  list.value = res.data.list
}
initBanner()

 const products = ref ([])
// 商品相关属性
const params = ref({
  page: 1,
  pageSize: 10,
  categoryID: 0
})

const flowData = ref([])
const isBottom = ref(false)
const isLoading = ref(true)

// 跳转到搜索页面
const goToSearch = () => {
  uni.navigateTo({
    url: '/pages/search/index'
  })
}

const getRecommend = async () =>{
  const res = await getGoodList({
    recommend:true
  })
  if (res.code === 0 && res.data.list.length) {
    products.value = res.data.list
  }
}

getRecommend()

// 修改 lower 函数，确保正确处理页码
const lower = async (isRefresh = false) => {
  if(isBottom.value && !isRefresh) {
    return
  } else {
    isLoading.value = true
    if(!isRefresh) {
      params.value.page += 1  // 这里会正确递增页码
    }

    try {
      const res = await getGoodList(params.value)
      if (res.code === 0 && res.data.list.length) {
        if(isRefresh) {
          flowData.value = res.data.list
        } else {
          flowData.value.push(...res.data.list)
        }
        totalCount.value = res.data.total || flowData.value.length
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

const totalCount = ref(0)

// 处理自动加载更多
const handleAutoLoadMore = () => {
  // 直接调用现有的 lower 方法，不需要修改页码
  // lower 方法内部会自动处理页码递增
  lower(false)
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
page {
  background-color: #f8f8f8;
}

.content {
  width: 100%;
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f8f8f8;
}

// 刘海屏适配
.status-bar-placeholder {
  width: 100%;
  height: var(--status-bar-height);
  background-color: #fea94a;
}

.test {
  display: flex;
  margin-top: 64rpx; // 原 2rem ≈ 64rpx
  width: 100%;
}

// 顶部状态栏
.status-bar {
  background-color: #fea94a;
  display: flex;
  align-items: center;
  padding: 16rpx 24rpx;
  box-sizing: border-box;
  position: sticky;
  top: 0;
  z-index: 99;
  width: 100%;
}

.scan-icon,
.message-icon {
  width: 60rpx;
  height: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.search-bar {
  flex: 1;
  background-color: #fff;
  height: 72rpx;
  border-radius: 36rpx;
  margin: 0 20rpx;
  display: flex;
  align-items: center;
  padding: 0 24rpx;
}

.search-icon {
  margin-right: 12rpx;
}

.search-placeholder {
  color: #999;
  font-size: 28rpx;
}

// 促销横幅
.promo-right {
  background-color: #FFD700;
}



// 右侧图片容器
.promo-image-container {
  position: absolute;
  right: 40rpx;
  display: flex;
  align-items: center;
}

.promo-text {
  color: white;
  font-size: 28rpx;
  font-weight: bold;
}

.promo-tag {
  display: inline-block;
  background-color: #ff4a4a;
  color: white;
  font-size: 20rpx;
  padding: 2rpx 8rpx;
  border-radius: 8rpx;
  margin-top: 8rpx;
}

// 底部加载及空状态
.loading-status {
  padding: 40rpx 0;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80rpx 0;
}

.empty-state text {
  font-size: 28rpx;
  color: #999;
  margin-top: 24rpx;
}
</style>

