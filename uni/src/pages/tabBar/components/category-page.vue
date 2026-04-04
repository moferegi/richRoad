<template>
  <view class="nf-cate">
    <view class="nf-cate-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <text class="nf-back-arrow">←</text>
        </view>
        <text class="nf-navbar-title">{{ $t('categoryDetail') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <!-- 分类选择横栏 -->
    <view class="nf-cate-bar">
      <scroll-view scroll-x class="nf-cate-scroll">
        <view class="nf-cate-tabs">
          <view
            class="nf-cate-tab"
            :class="{ 'nf-cate-tab-active': selectedIndex === index }"
            v-for="(item, index) in gridList"
            :key="index"
            @tap="goto(item, index)"
          >
            <image class="nf-cate-tab-icon" :src="getUrl(item.externalIconPath || item.icons)" mode="aspectFill"></image>
            <text class="nf-cate-tab-text">{{ $lt(item.title) }}</text>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- 商品内容区域 -->
    <scroll-view scroll-y class="nf-cate-content" @scrolltolower="loadMoreGoods">
      <!-- 商品网格 -->
      <view class="nf-grid" v-if="flowData.length">
        <view
          class="nf-grid-item"
          v-for="(item, index) in flowData"
          :key="index"
          @tap="handleGoodsClick(item)"
        >
          <view class="nf-grid-img-wrap">
            <image class="nf-grid-img" :src="getUrl(item.externalImagePath || item.imageUrl)" mode="aspectFill"></image>
            <view class="nf-grid-img-overlay"></view>
            <!-- 折扣标签 -->
            <view class="nf-grid-badge" v-if="item.discount && item.discount < 10">
              <text class="nf-grid-badge-text">{{ item.discount }}折</text>
            </view>
          </view>
          <view class="nf-grid-info">
            <text class="nf-grid-title">{{ $lt(item.title) }}</text>
            <view class="nf-grid-price-row">
              <text class="nf-grid-price">{{ cs }}{{ formatPrice(item.price) }}</text>
              <text class="nf-grid-sales">{{ $t('sold') }} {{ item.saleNum || 0 }}</text>
            </view>
            <view v-if="item.tags && item.tags.length > 0" class="nf-grid-tags">
              <text
                v-for="tag in item.tags"
                :key="tag.ID"
                class="nf-grid-tag"
                :style="{ color: tag.color, borderColor: tag.color + '55', background: tag.color + '15' }"
              >{{ tag.name }}</text>
            </view>
          </view>
        </view>
      </view>

      <!-- 加载状态 -->
      <view class="nf-cate-loading" v-if="loading">
        <text class="nf-cate-loading-text">{{ $t('loading') }}</text>
      </view>

      <!-- 到底提示 -->
      <view class="nf-cate-loading" v-if="noMore && flowData.length">
        <text class="nf-cate-loading-text">{{ $t('reachedBottom') }}</text>
      </view>

      <!-- 空状态 -->
      <view class="nf-cate-empty" v-if="!flowData.length && !loading">
        <view class="nf-cate-empty-icon">
          <uni-icons type="shop" size="48" color="rgba(229,9,20,0.4)" />
        </view>
        <text class="nf-cate-empty-text">{{ $t('noProductsInCategory') }}</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { getCategoryMobile, getGoodList } from '@/api/homePage.js'
import { onLoad } from '@dcloudio/uni-app'
import { getUrl } from '@/utils/url'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'

const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
const cs = computed(() => appConfigStore.currencySymbol)
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)

const selectedIndex = ref(0)
const currentCategoryID = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const noMore = ref(false)
const gridList = ref([])
const flowData = ref([])

const goBack = () => {
  uni.navigateBack({ fail: () => uni.switchTab({ url: '/pages/tabBar/index' }) })
}

const formatPrice = (priceInCents) => {
  if (!priceInCents && priceInCents !== 0) return '0.00'
  const cents = parseInt(priceInCents)
  if (isNaN(cents)) return '0.00'
  const yuan = Math.floor(cents / 100)
  const remainingCents = cents % 100
  return `${yuan}.${remainingCents.toString().padStart(2, '0')}`
}

onLoad((options) => {
  if (options.id) {
    initCategory(options.id)
  }
})

const initCategory = async (parentID) => {
  const res = await getCategoryMobile({ parentID })
  if (res.code === 0 && res.data.length) {
    gridList.value = res.data
    currentCategoryID.value = res.data[0].ID
  }

  const targetIndex = gridList.value.findIndex(item => item.ID == currentCategoryID.value)
  if (targetIndex !== -1) {
    selectedIndex.value = targetIndex
  }

  if (currentCategoryID.value) {
    await loadGoodsList(currentCategoryID.value, true)
  }
}

const loadGoodsList = async (categoryID, isReset = false) => {
  if (loading.value) return
  loading.value = true

  try {
    if (isReset) {
      currentPage.value = 1
      flowData.value = []
      noMore.value = false
    }

    const res = await getGoodList({
      page: currentPage.value,
      pageSize: pageSize.value,
      categoryID: categoryID
    })

    if (res.code === 0) {
      const newList = res.data?.list || []
      if (isReset || currentPage.value === 1) {
        flowData.value = newList
      } else {
        flowData.value = [...flowData.value, ...newList]
      }
      if (newList.length < pageSize.value) {
        noMore.value = true
      } else {
        currentPage.value += 1
      }
    } else {
      uni.showToast({ title: $t.value('loadFailed'), icon: 'none' })
    }
  } catch (error) {
    console.error('加载商品列表失败:', error)
    uni.showToast({ title: $t.value('loadFailed'), icon: 'none' })
  } finally {
    loading.value = false
  }
}

const loadMoreGoods = () => {
  if (!loading.value && !noMore.value && currentCategoryID.value) {
    loadGoodsList(currentCategoryID.value, false)
  }
}

const goto = async (item, index) => {
  selectedIndex.value = index
  currentCategoryID.value = item.ID
  await loadGoodsList(item.ID, true)
}

const handleGoodsClick = (item) => {
  uni.navigateTo({
    url: '/pages/player/index?id=' + item.ID
  })
}
</script>

<style scoped lang="scss">
.nf-cate {
  min-height: 100vh;
  background: #000;
  display: flex;
  flex-direction: column;
}

.nf-cate-bg {
  position: fixed;
  top: 0; left: 0; right: 0;
  height: 400rpx;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(ellipse at 30% 0%, rgba(229, 9, 20, 0.15) 0%, transparent 55%),
    radial-gradient(ellipse at 70% 10%, rgba(229, 9, 20, 0.08) 0%, transparent 50%);
}

/* ===== 导航栏 ===== */
.nf-navbar {
  position: fixed;
  top: 0; left: 0; right: 0;
  z-index: 100;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
}

.nf-navbar-status {
  height: var(--status-bar-height, 44px);
}

.nf-navbar-content {
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24rpx;
}

.nf-navbar-back {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background 0.3s;
  &:active { background: rgba(255, 255, 255, 0.1); }
}

.nf-back-arrow {
  font-size: 36rpx;
  color: #fff;
}

.nf-navbar-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #fff;
  letter-spacing: 2rpx;
}

/* ===== 分类选择栏 ===== */
.nf-cate-bar {
  position: fixed;
  top: calc(var(--status-bar-height, 44px) + 88rpx);
  left: 0; right: 0;
  z-index: 99;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
}

.nf-cate-scroll {
  width: 100%;
  white-space: nowrap;
}

.nf-cate-tabs {
  display: flex;
  padding: 20rpx 12rpx;
  gap: 8rpx;
}

.nf-cate-tab {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex-shrink: 0;
  min-width: 120rpx;
  padding: 12rpx 16rpx;
  border-radius: 16rpx;
  border: 1rpx solid transparent;
  transition: all 0.3s;

  &:active { transform: scale(0.95); }
}

.nf-cate-tab-active {
  background: rgba(229, 9, 20, 0.12);
  border-color: rgba(229, 9, 20, 0.3);

  .nf-cate-tab-icon {
    border-color: rgba(229, 9, 20, 0.6);
    box-shadow: 0 0 16rpx rgba(229, 9, 20, 0.25);
  }

  .nf-cate-tab-text {
    color: #fff;
    font-weight: 700;
  }
}

.nf-cate-tab-icon {
  width: 72rpx;
  height: 72rpx;
  border-radius: 50%;
  margin-bottom: 8rpx;
  border: 2rpx solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.04);
  transition: all 0.3s;
}

.nf-cate-tab-text {
  font-size: 22rpx;
  color: rgba(255, 255, 255, 0.6);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 120rpx;
  text-align: center;
  letter-spacing: 1rpx;
  transition: all 0.3s;
}

/* ===== 商品内容 ===== */
.nf-cate-content {
  flex: 1;
  /* navbar height + cate-bar height */
  margin-top: calc(var(--status-bar-height, 44px) + 88rpx + 130rpx);
  min-height: calc(100vh - var(--status-bar-height, 44px) - 88rpx - 130rpx);
}

/* ===== 商品网格 ===== */
.nf-grid {
  display: flex;
  flex-wrap: wrap;
  padding: 16rpx 20rpx;
  gap: 16rpx;
}

.nf-grid-item {
  width: calc(50% - 8rpx);
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.06);
  border-radius: 20rpx;
  overflow: hidden;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  transition: transform 0.3s;

  &:active { transform: scale(0.97); }
}

.nf-grid-img-wrap {
  position: relative;
  width: 100%;
  height: 340rpx;
  overflow: hidden;
}

.nf-grid-img {
  width: 100%;
  height: 100%;
}

.nf-grid-img-overlay {
  position: absolute;
  bottom: 0; left: 0; right: 0;
  height: 40%;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.5));
  pointer-events: none;
}

.nf-grid-badge {
  position: absolute;
  top: 12rpx;
  right: 12rpx;
  background: linear-gradient(135deg, #e50914, #b20710);
  padding: 4rpx 14rpx;
  border-radius: 10rpx;
  box-shadow: 0 4rpx 12rpx rgba(229, 9, 20, 0.4);
}

.nf-grid-badge-text {
  font-size: 20rpx;
  color: #fff;
  font-weight: 700;
}

.nf-grid-info {
  padding: 16rpx 16rpx 20rpx;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.nf-grid-title {
  font-size: 26rpx;
  color: #fff;
  font-weight: 600;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.nf-grid-price-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nf-grid-price {
  font-size: 32rpx;
  color: #e50914;
  font-weight: 800;
}

.nf-grid-sales {
  font-size: 20rpx;
  color: rgba(255, 255, 255, 0.35);
}

.nf-grid-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6rpx;
}

.nf-grid-tag {
  font-size: 18rpx;
  padding: 2rpx 10rpx;
  border-radius: 8rpx;
  border: 1rpx solid;
  font-weight: 500;
}

/* ===== 加载与空状态 ===== */
.nf-cate-loading {
  padding: 40rpx;
  text-align: center;
}

.nf-cate-loading-text {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.3);
  letter-spacing: 2rpx;
}

.nf-cate-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 160rpx 0;
}

.nf-cate-empty-icon {
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

.nf-cate-empty-text {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.3);
  letter-spacing: 2rpx;
}
</style>
