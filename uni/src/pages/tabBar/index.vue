<template>
  <view class="nf-home">
    <!-- 背景光晕 -->
    <view class="nf-home-bg"></view>

    <!-- 顶部搜索栏 -->
    <view class="nf-topbar">
      <view class="nf-topbar-status"></view>
      <view class="nf-topbar-row">
        <!-- 左侧语言切换 -->
        <view class="nf-topbar-btn" @click="showLangPicker = true">
          <text class="nf-topbar-btn-text">{{ langLabel }}</text>
        </view>

        <!-- 搜索框 -->
        <view class="nf-search">
          <uni-icons type="search" size="16" color="rgba(255,255,255,0.4)"></uni-icons>
          <input
            class="nf-search-input"
            :placeholder="$t('searchPlaceholder')"
            v-model="searchKeyword"
            @confirm="goToSearchWithKeyword"
            @click.stop
          />
        </view>

        <!-- 搜索按钮 -->
        <view class="nf-topbar-btn nf-topbar-btn-red" @click="goToSearchWithKeyword">
          <uni-icons type="search" size="18" color="#fff"></uni-icons>
        </view>
      </view>
    </view>

    <scroll-view scroll-y="true" :show-scrollbar="false" class="nf-scroll" @scrolltolower="debouncedLower">
      <!-- 轮播图区域 -->
      <swpiers :lists="list"></swpiers>
      <!-- 限时秒杀区域 -->
      <seckilling v-if="products.length > 0" :productData="products"></seckilling>
      <!-- 分类导航 -->
      <categories :categoriesData="gridList" v-model="activeCategoryID" @change="onCategoryChange"></categories>
      <!-- 商品展示区 -->
      <view class="nf-goods-section" :class="{ 'nf-goods-fade': switching }">
        <noPaginRowGoodList
            :goodsList="flowData"
            :current-page="params.page"
            :page-size="params.pageSize"
            :total="totalCount"
            @load-more="handleAutoLoadMore"
        />
      </view>
      <!-- 底部加载状态 -->
      <view class="nf-loading-status" v-if="flowData.length > 0">
        <gva-divider :text="isBottom ? $t('reachedBottom') : $t('loading')"></gva-divider>
      </view>
      <!-- 空状态 -->
      <view class="nf-empty" v-if="flowData.length === 0 && !isLoading">
        <view class="nf-empty-icon">
          <uni-icons type="shop" size="48" color="rgba(229,9,20,0.4)" />
        </view>
        <text class="nf-empty-text">{{ $t('noGoods') }}</text>
      </view>
    </scroll-view>

    <!-- 语言切换弹窗 -->
    <lang-switch v-model="showLangPicker" />
  </view>
</template>
<script setup>
import  { ref, computed } from 'vue';
import { getCategoryMobile, getGoodList } from '@/api/homePage.js'
import { getBannerList } from '@/api/homePage.js'
import noPaginRowGoodList from '@/components/good-list/no-pagin-row-good-list.vue'
import swpiers from './components/swiper.vue'
import categories from './components/categories.vue'
import seckilling from './components/seckilling.vue';
import langSwitch from '@/components/lang-switch/lang-switch.vue'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)
const langLabel = computed(() => {
  const map = { zh: '中', en: 'EN', mn: 'MN' }
  return map[langStore.locale] || '中'
})
const showLangPicker = ref(false)

// 启动时恢复 tabBar 语言
langStore.updateTabBar(langStore.locale)

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
const activeCategoryID = ref(0)
const params = ref({
  page: 1,
  pageSize: 10,
  categoryID: 0
})

const flowData = ref([])
const isBottom = ref(false)
const isLoading = ref(true)
const switching = ref(false)

// 清空搜索关键词
searchKeyword.value = ''
// 带关键词跳转到搜索页面
const goToSearchWithKeyword = () => {
  const keyword = searchKeyword.value.trim()
  if (!keyword) {
    uni.showToast({
      title: $t.value('searchEmpty'),
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
  if (isRefresh) {
    params.value.page = 1
  } else {
    params.value.page += 1
  }

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


// 分类切换（带淡入淡出）
const onCategoryChange = async (categoryID) => {
  params.value.categoryID = categoryID || 0
  // 淡出（保留旧数据撑高度）
  switching.value = true
  await new Promise(r => setTimeout(r, 220))
  // 请求新数据
  params.value.page = 1
  isBottom.value = false
  try {
    const res = await getGoodList(params.value)
    if (res.code === 0) {
      const listData = res.data.list || []
      flowData.value = listData
      totalCount.value = res.data.total ?? listData.length
      isBottom.value = listData.length < params.value.pageSize
    } else {
      flowData.value = []
      isBottom.value = true
    }
  } catch (e) {
    flowData.value = []
  }
  // 淡入
  switching.value = false
}

// 分类tabs相关业务逻辑
const gridList = ref([])
const initCategory = async () => {
  const res = await getCategoryMobile()
  if (res.code === 0 && res.data.length) {
    gridList.value = res.data
  }
}
initCategory()
// 初始加载全部商品
lower(true)


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

.nf-home {
  width: 100%;
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #000;
  position: relative;
}

/* 背景装饰光晕 */
.nf-home-bg {
  position: fixed;
  top: 0; left: 0; right: 0;
  height: 500rpx;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(ellipse at 20% 0%, rgba(229, 9, 20, 0.12) 0%, transparent 60%),
    radial-gradient(ellipse at 80% 10%, rgba(229, 9, 20, 0.08) 0%, transparent 50%);
}

/* ===== 顶部栏 ===== */
.nf-topbar {
  background: rgba(0, 0, 0, 0.85);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
  padding: 0 28rpx 16rpx;
  position: sticky;
  top: 0;
  z-index: 99;
}

.nf-topbar-status {
  width: 100%;
  height: var(--status-bar-height, 0px);
}

/* #ifdef MP-WEIXIN */
.nf-topbar-status {
  height: var(--status-bar-height, 44px);
}
/* #endif */

.nf-topbar-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-top: 12rpx;
}

.nf-topbar-btn {
  width: 68rpx;
  height: 68rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.06);
  border: 1rpx solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.3s;

  &:active {
    background: rgba(255, 255, 255, 0.12);
    transform: scale(0.93);
  }
}

.nf-topbar-btn-red {
  background: linear-gradient(135deg, #e50914, #b20710);
  border: none;
  box-shadow: 0 4rpx 16rpx rgba(229, 9, 20, 0.35);

  &:active {
    box-shadow: 0 2rpx 8rpx rgba(229, 9, 20, 0.5);
  }
}

.nf-topbar-btn-text {
  font-size: 22rpx;
  font-weight: 800;
  color: #fff;
  letter-spacing: 0;
}

/* ===== 搜索框 ===== */
.nf-search {
  flex: 1;
  height: 72rpx;
  background: rgba(255, 255, 255, 0.06);
  border: 1rpx solid rgba(255, 255, 255, 0.1);
  border-radius: 36rpx;
  display: flex;
  align-items: center;
  padding: 0 24rpx;
  gap: 12rpx;
  transition: border-color 0.3s;
}

.nf-search-input {
  flex: 1;
  color: #fff;
  font-size: 28rpx;
  border: none;
  outline: none;
  background: transparent;

  &::placeholder {
    color: rgba(255, 255, 255, 0.35);
  }
}

/* ===== 滚动区域 ===== */
.nf-scroll {
  flex: 1;
  width: 100%;
  background: #000;
}

/* ===== 商品区块 ===== */
.nf-goods-section {
  background: #000;
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.nf-goods-fade {
  opacity: 0;
  transform: translateY(8rpx);
}

/* ===== 加载状态 ===== */
.nf-loading-status {
  padding: 40rpx 0;
  background: #000;
}

/* ===== 空状态 ===== */
.nf-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0 80rpx;
}

.nf-empty-icon {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  background: rgba(229, 9, 20, 0.08);
  border: 1rpx solid rgba(229, 9, 20, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24rpx;
}

.nf-empty-text {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.3);
  letter-spacing: 2rpx;
}
</style>

