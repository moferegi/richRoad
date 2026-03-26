<template>
  <view class="content">
    <!-- 顶部搜索栏  -->
    <view class="status-bar-placeholder"></view>
    <view class="status-bar">
      <view class="test">
        <!-- 左侧客服按钮 -->
        <view class="message-icon">
          <button class="contact-action-button" open-type="contact" :session-from="sessionFrom">
            <uni-icons type="chat" size="20" color="#fff"></uni-icons>
          </button>
        </view>

        <!-- 中间搜索框 -->
        <view class="search-bar">
          <view class="search-icon">
            <uni-icons type="search" size="18" color="#999"></uni-icons>
          </view>
          <input
            class="search-input"
            placeholder="请输入您想搜索的商品"
            v-model="searchKeyword"
            @confirm="goToSearchWithKeyword"
            @click.stop
          />
        </view>

        <!-- 右侧搜索按钮 -->
        <view class="search-button" @click="goToSearchWithKeyword">
          <uni-icons type="search" size="20" color="#fff"></uni-icons>
        </view>
      </view>
    </view>
    <scroll-view scroll-y="true" class="scroll-Y" @scrolltolower="debouncedLower">
      <!-- 轮播图区域 -->
      <swpiers :lists="list"></swpiers>
      <!-- 分类导航 -->
      <categories :categoriesData="gridList"></categories>
      <!-- 限时秒杀区域 -->
      <seckilling v-if="products.length > 0" :productData="products"></seckilling>
      <!-- 商品展示区 -->
      <view class="goods-section">
        <noPaginRowGoodList
            :goodsList="flowData"
            :current-page="params.page"
            :page-size="params.pageSize"
            :total="totalCount"
            @load-more="handleAutoLoadMore"
        />
      </view>
      <!-- 底部加载状态 -->
      <view class="loading-status" v-if="flowData.length > 0">
        <gva-divider :text="isBottom ? '已经到底啦' : '加载中...'"></gva-divider>
      </view>
      <!-- 空状态 -->
      <view class="empty-state" v-if="flowData.length === 0 && !isLoading">
        <uni-icons type="shop" size="60" color="#ddd" />
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
const searchKeyword = ref('')

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

const sessionFrom = ref('');
// 清空搜索关键词
searchKeyword.value = ''
// 带关键词跳转到搜索页面
const goToSearchWithKeyword = () => {
  const keyword = searchKeyword.value.trim()
  if (!keyword) {
    uni.showToast({
      title: '请输入搜索关键词',
      icon: 'none'
    })
    return
  }
  
  // 使用 uni.setStorageSync 临时存储搜索关键词
  uni.setStorageSync('searchKeyword', keyword)
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
/*const lower = async (isRefresh = false) => {
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
}*/
const lower = async (isRefresh = false) => {
  if (isBottom.value && !isRefresh) return
  isLoading.value = true
  if (!isRefresh) params.value.page += 1

  try {
    const res = await getGoodList(params.value)
    if (res.code === 0) {
      const listData = res.data.list || []
      // 合并或重置
      if (isRefresh) {
        flowData.value = listData
      } else {
        flowData.value.push(...listData)
      }
      // 更新总数
      totalCount.value = res.data.total ?? flowData.value.length
      // 如果返回条数小于 pageSize，说明已经是最后一页
      isBottom.value = listData.length < params.value.pageSize
    } else {
      // 接口异常也视为无更多
      isBottom.value = true
    }
  } catch (error) {
    console.error('获取商品列表失败', error)
  } finally {
    isLoading.value = false
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
  background-color: #000;
}

.content {
  width: 100%;
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #000;
}

// ===== Netflix 风格顶部状态栏 =====
.status-bar-placeholder {
  width: 100%;
  height: var(--status-bar-height, 0px);
  background-color: #000;
}

/* #ifdef MP-WEIXIN */
.status-bar-placeholder {
  height: var(--status-bar-height, 44px);
}
/* #endif */

.status-bar {
  background: linear-gradient(180deg, rgba(0, 0, 0, 0.95) 0%, rgba(0, 0, 0, 0.8) 100%);
  backdrop-filter: blur(10px);
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.1);
  padding: 12rpx 24rpx;
  box-sizing: border-box;
  position: sticky;
  top: 0;
  z-index: 99;
  width: 100%;
  display: flex;
  align-items: center;
}

.test {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-top: 20rpx;
}

// ===== Netflix 风格按钮 =====
.message-icon,
.search-button {
  width: 60rpx;
  height: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.1);
  border: 1rpx solid rgba(255, 255, 255, 0.15);
  transition: all 0.3s ease;

  &:active {
    background-color: rgba(255, 255, 255, 0.2);
    transform: scale(0.95);
  }
}

.contact-action-button {
  background: none;
  border: none;
  padding: 0;
  margin: 0;
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

  &::after {
    border: none;
  }
}

// ===== Netflix 风格搜索框 =====
.search-bar {
  flex: 1;
  background-color: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border: 1rpx solid rgba(255, 255, 255, 0.15);
  height: 72rpx;
  border-radius: 36rpx;
  display: flex;
  align-items: center;
  padding: 0 24rpx;
  gap: 12rpx;
}

.search-icon {
  margin-right: 0;
}

.search-input {
  flex: 1;
  color: #fff;
  font-size: 28rpx;
  border: none;
  outline: none;
  background: transparent;

  &::placeholder {
    color: rgba(255, 255, 255, 0.5);
  }
}

// ===== 滚动区域 =====
.scroll-Y {
  flex: 1;
  width: 100%;
  background-color: #000;
  // padding-top: 116rpx; // 为sticky头部留出足够空间 (status-bar高度约116rpx + 额外缓冲)
}

// ===== 商品展示区 =====
.goods-section {
  background-color: #000;
  padding: 0;
}

// ===== Netflix 风格加载和空状态 =====
.loading-status {
  padding: 40rpx 0;
  background-color: #000;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80rpx 0;
  background-color: #000;
}

.empty-state text {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.4);
  margin-top: 24rpx;
}

// ===== 移除不需要的样式 =====
.promo-right,
.promo-image-container,
.promo-text,
.promo-tag {
  // 这些样式在子组件中处理
}
</style>

