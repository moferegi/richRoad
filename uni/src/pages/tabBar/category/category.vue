<template>
  <view class="page-container">
    <!-- 顶部搜索栏 -->
    <view class="search-header">
      <view class="search-box">
      </view>
    </view>

    <view class="cate-box">
      <!-- 左侧分类导航 -->
      <scroll-view class="left-sidebar" scroll-y="true">
        <view
          v-for="(item, index) in catelist"
          :key="index"
          :class="['category-item', { 'active': activeindex === index }]"
          @tap="checkitem(index, item)"
        >
          <view class="category-text">{{ $lt(item.title) }}</view>
          <view v-if="activeindex === index" class="active-indicator"></view>
        </view>
      </scroll-view>

      <!-- 右侧商品列表 -->
      <scroll-view
        :lower-threshold="lowerThresholdInPx"
        class="right-content"
        scroll-y="true"
        @scrolltolower="handleScrollToLower"
      >
        <!-- 分类标题 -->
        <view v-if="currentCategory" class="content-header">
          <view class="category-banner">
            <view class="banner-title">{{ $lt(currentCategory.title) }}</view>
            <view class="banner-subtitle">{{ $t('categorySubtitle') }}</view>
          </view>
        </view>

        <!-- 商品网格 -->
        <template v-for="category in categoriesWithGoods" :key="category.ID">
          <view v-if="category.goods && category.goods.length > 0" class="category-section">
            <view class="section-title">
              <view class="title-line"></view>
              <text class="title-text">{{ $lt(category.title) }}</text>
              <view class="title-line"></view>
            </view>

            <view class="goods-grid">
              <view
                v-for="item in category.goods"
                :key="item.ID"
                class="goods-card"
                @tap="goto(item)"
              >
                <view class="card-image-container">
                  <image
                    :src="getUrl(item.imageUrl)"
                    class="goods-image"
                    mode="aspectFill"
                    :lazy-load="true"
                  ></image>
                  <view v-if="item.tags && item.tags.length > 0" class="image-tags">
                    <text
                      v-for="tag in item.tags.slice(0, 1)"
                      :key="tag.ID"
                      :style="{
                        background: tag.color || '#ff6b6b'
                      }"
                      class="tag-badge"
                    >
                      {{ $lt(tag.nameI18n || tag.name) }}
                    </text>
                  </view>
                </view>

                <view class="card-content">
                  <view class="goods-title">{{ $lt(item.title) }}</view>
                  <view class="goods-desc">{{ $lt(item.description) }}</view>

                  <view class="price-section">
                    <view class="current-price">
                      <text class="price-symbol">{{ cs }}</text>
                      <text class="price-value">{{ formatGoodsPrice(item) }}</text>
                    </view>
                    <view class="add-cart-btn">
                      <uni-icons type="plus" size="16" color="#fff"></uni-icons>
                    </view>
                  </view>
                </view>
              </view>
            </view>
          </view>
        </template>

        <!-- 加载状态 -->
        <view v-if="loading && categoriesWithGoods.length === 0 && currentCategory" class="loading-state">
          <view class="loading-spinner"></view>
          <text class="loading-text">{{ $t('loading') }}</text>
        </view>

        <!-- 空状态 -->
        <view v-if="categoriesWithGoods.length === 0 && !loading && currentCategory" class="empty-state">
          <view class="empty-icon">📦</view>
          <text class="empty-text">{{ $t('noGoods') }}</text>
          <text class="empty-desc">{{ $t('categoryEmpty') }}</text>
        </view>

        <!-- 底部间距 -->
        <view class="bottom-spacing"></view>
      </scroll-view>
    </view>
  </view>
</template>

<script setup>
import {ref, computed, onMounted} from 'vue'
import {getUrl} from '@/utils/url'
import {getCategoryMobile, getChildrenCategoryAndProduct} from '@/api/homePage.js'
import { formatLocalizedPrice } from '@/utils/price-i18n.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'

const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
const cs = computed(() => appConfigStore.currencySymbol)
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)
const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')

const formatGoodsPrice = (item) => formatLocalizedPrice(item?.price, item?.priceI18n, locale.value)

// 分类数据
const catelist = ref([])
const activeindex = ref(0)
const currentCategory = ref(null)

// 分类和商品数据
const categoriesWithGoods = ref([])
const loading = ref(false)

// 新增：用于存储转换后的 lower-threshold 的 px 值
const lowerThresholdInPx = ref(50) // 给一个默认的px值，在onMounted中会被实际rpx转换值覆盖

// 获取分类数据
const getCategoryData = async () => {
  try {
    const res = await getCategoryMobile()
    if (res && res.data) {
      catelist.value = res.data
      if (catelist.value.length > 0) {
        currentCategory.value = catelist.value[0]
        await loadCategoryData(currentCategory.value.ID)
      } else {
        categoriesWithGoods.value = [];
        currentCategory.value = null;
      }
    } else {
      categoriesWithGoods.value = [];
      currentCategory.value = null;
      uni.showToast({title: $t.value('loadFail'), icon: 'none'});
    }
  } catch (error) {
    console.error('获取分类数据失败:', error)
    categoriesWithGoods.value = [];
    currentCategory.value = null;
    uni.showToast({
      title: $t.value('loadFail'),
      icon: 'none'
    })
  }
}

// 加载分类和商品数据
const loadCategoryData = async (parentID) => {
  if (loading.value) {
    return
  }
  loading.value = true
  try {
    const params = {
      parentID: parentID
    }
    const res = await getChildrenCategoryAndProduct(params)
    if (res && res.code === 0 && res.data) {
      categoriesWithGoods.value = res.data || []
    } else {
      categoriesWithGoods.value = []
      console.error('获取分类商品数据API响应异常:', res)
    }
  } catch (error) {
    console.error('获取分类商品数据失败:', error)
    categoriesWithGoods.value = []
    uni.showToast({
      title: $t.value('loadFail'),
      icon: 'none'
    })
  } finally {
    loading.value = false
  }
}

// 切换分类
const checkitem = async (index, item) => {
  if (activeindex.value === index && categoriesWithGoods.value.length > 0) {
    return;
  }
  activeindex.value = index
  currentCategory.value = item
  await loadCategoryData(item.ID)
}

// 滚动到底部时触发（暂时保留，但不再需要分页加载）
const handleScrollToLower = () => {
  // 新的数据结构一次性加载所有分类和商品，不需要分页
}

// 点击商品
const clickitem = (item) => {
  getgoods(item)
}

const getgoods = (item) => {
  // debug removed
}

const goto = (item) => {
  uni.navigateTo({
    url: '/pages/player/index?id=' + item.ID
  })
}
// 页面加载时执行
onMounted(() => {
  // 将 200rpx 转换为 px 值并设置
  // uni.upx2px 在编译到不同平台时，会基于基准宽度（通常是750rpx）进行转换
  lowerThresholdInPx.value = uni.upx2px(200);
  // console.log(`200rpx is approximately ${lowerThresholdInPx.value}px on this device.`); // 用于调试

  getCategoryData()
})
</script>

<style lang="scss" scoped>
.page-container {
  min-height: 100vh;
  background: #000;

  /* #ifdef H5 */
  min-height: calc(100vh - var(--window-top) - var(--window-bottom));
  /* #endif */
}

// 搜索头部
.search-header {
  padding: 20rpx 30rpx;
  backdrop-filter: blur(10rpx);

  .search-box {
    display: flex;
    align-items: center;
    border-radius: 50rpx;
    padding: 20rpx 30rpx;
    gap: 15rpx;
    transition: all 0.3s ease;

    &:active {
      transform: scale(0.98);
      background: rgba(255, 255, 255, 0.06);
    }

    .search-placeholder {
      color: rgba(255, 255, 255, 0.4);
      font-size: 28rpx;
    }
  }
}

.cate-box {
  display: flex;
  height: calc(100vh - 120rpx);
  background: #141414;
  border-radius: 30rpx 30rpx 0 0;
  overflow: hidden;

  /* #ifdef H5 */
  height: calc(100vh - var(--window-top) - var(--window-bottom) - 120rpx);
  /* #endif */
}

// 左侧分类导航
.left-sidebar {
  width: 200rpx;
  background: linear-gradient(180deg, #1a1a1a 0%, #111 100%);
  position: relative;

  .category-item {
    position: relative;
    margin: 8rpx 0;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

    .category-text {
      padding: 30rpx 20rpx;
      text-align: center;
      font-size: 26rpx;
      color: rgba(255, 255, 255, 0.5);
      font-weight: 500;
      transition: all 0.3s ease;
      position: relative;
      z-index: 2;
    }

    .active-indicator {
      position: absolute;
      right: 0;
      top: 50%;
      transform: translateY(-50%);
      width: 6rpx;
      height: 40rpx;
      background: #e50914;
      border-radius: 6rpx 0 0 6rpx;
      animation: slideIn 0.3s ease;
    }

    &.active {
      background: #141414;
      margin-right: 6rpx;
      border-radius: 25rpx 0 0 25rpx;

      .category-text {
        color: #e50914;
        font-weight: 600;
        font-size: 28rpx;
      }
    }

    &:not(.active):active {
      transform: scale(0.95);
      background: rgba(255, 255, 255, 0.06);
    }
  }
}

// 右侧内容区域
.right-content {
  flex: 1;
  background: #141414;

  .content-header {
    position: sticky;
    top: 0;
    z-index: 10;
    background: linear-gradient(135deg, #1a1a1a 0%, #2a2a2a 100%);

    .category-banner {
      padding: 40rpx 30rpx;
      text-align: center;

      .banner-title {
        font-size: 36rpx;
        font-weight: bold;
        margin-bottom: 10rpx;
        color: #fff;
      }

      .banner-subtitle {
        font-size: 24rpx;
        color: rgba(255, 255, 255, 0.5);
      }
    }
  }

  .category-section {
    padding: 0 20rpx;

    .section-title {
      display: flex;
      align-items: center;
      margin: 40rpx 0 30rpx;

      .title-line {
        flex: 1;
        height: 2rpx;
        background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent);
      }

      .title-text {
        padding: 0 30rpx;
        font-size: 28rpx;
        font-weight: 600;
        color: rgba(255, 255, 255, 0.8);
      }
    }

    .goods-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 20rpx;
      padding-bottom: 30rpx;
    }
  }
}

// 商品卡片
.goods-card {
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.06);
  border-radius: 20rpx;
  overflow: hidden;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

  &:active {
    transform: scale(0.97);
  }

  .card-image-container {
    position: relative;
    width: 100%;
    height: 280rpx;
    overflow: hidden;

    .goods-image {
      width: 100%;
      height: 100%;
      transition: transform 0.3s ease;
    }

    .image-tags {
      position: absolute;
      top: 15rpx;
      left: 15rpx;

      .tag-badge {
        display: inline-block;
        padding: 8rpx 16rpx;
        border-radius: 20rpx;
        font-size: 20rpx;
        color: #fff;
        font-weight: 500;
        text-shadow: 0 1rpx 2rpx rgba(0, 0, 0, 0.3);
      }
    }
  }

  .card-content {
    padding: 25rpx;

    .goods-title {
      font-size: 28rpx;
      font-weight: 600;
      color: #fff;
      line-height: 1.4;
      margin-bottom: 12rpx;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
    }

    .goods-desc {
      font-size: 24rpx;
      color: rgba(255, 255, 255, 0.45);
      line-height: 1.3;
      margin-bottom: 20rpx;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .price-section {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .current-price {
        display: flex;
        align-items: baseline;

        .price-symbol {
          font-size: 24rpx;
          color: #e50914;
          font-weight: 500;
        }

        .price-value {
          font-size: 32rpx;
          color: #e50914;
          font-weight: bold;
          margin-left: 2rpx;
        }
      }

      .add-cart-btn {
        width: 60rpx;
        height: 60rpx;
        background: #e50914;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.3s ease;

        &:active {
          transform: scale(0.9);
          opacity: 0.85;
        }
      }
    }
  }
}

// 加载状态
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 80rpx 0;

  .loading-spinner {
    width: 60rpx;
    height: 60rpx;
    border: 4rpx solid rgba(255, 255, 255, 0.1);
    border-top: 4rpx solid #e50914;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 30rpx;
  }

  .loading-text {
    font-size: 28rpx;
    color: rgba(255, 255, 255, 0.5);
  }
}

// 空状态
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 100rpx 0;

  .empty-icon {
    font-size: 120rpx;
    margin-bottom: 30rpx;
    opacity: 0.6;
  }

  .empty-text {
    font-size: 32rpx;
    color: rgba(255, 255, 255, 0.8);
    font-weight: 600;
    margin-bottom: 15rpx;
  }

  .empty-desc {
    font-size: 26rpx;
    color: rgba(255, 255, 255, 0.4);
  }
}

.bottom-spacing {
  height: 40rpx;
}

// 动画
@keyframes slideIn {
  from {
    transform: translateY(-50%) translateX(20rpx);
    opacity: 0;
  }
  to {
    transform: translateY(-50%) translateX(0);
    opacity: 1;
  }
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

// 响应式设计
@media (max-width: 750rpx) {
  .goods-grid {
    grid-template-columns: 1fr !important;
  }

  .left-sidebar {
    width: 160rpx !important;
  }
}
</style>
