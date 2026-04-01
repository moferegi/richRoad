<template>
  <view class="nf-home">
    <!-- 背景光晕 -->
    <view class="nf-home-bg"></view>

    <!-- 顶部搜索栏 -->
    <view class="nf-topbar" id="nf-topbar">
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

    <scroll-view
      scroll-y="true"
      :show-scrollbar="false"
      class="nf-scroll"
      :scroll-top="scrollTopVal"
      @scroll="onScroll"
      @scrolltolower="debouncedLower"
    >
      <!-- 轮播图区域 -->
      <swpiers :lists="list"></swpiers>

      <!-- 1. 公告走马灯 -->
      <announcement-marquee
        :enabled="announcementConfig.enabled"
        :text="announcementConfig.text"
        :textColor="announcementConfig.textColor"
        :speed="announcementConfig.speed"
      />

      <!-- 2. 预售商品区域 -->
      <view class="nf-presale-home" v-if="presaleList.length > 0">
        <view class="nf-presale-header">
          <text class="nf-presale-title">🔥 {{ $t('presaleSection') }}</text>
          <view class="nf-presale-more" @tap="goPresaleList">
            <text class="nf-presale-more-text">{{ $t('viewMore') }}</text>
            <uni-icons type="right" size="12" color="rgba(255,255,255,0.5)" />
          </view>
        </view>
        <scroll-view scroll-x class="nf-presale-scroll">
          <view class="nf-presale-items">
            <view
              class="nf-presale-item"
              v-for="(item, idx) in presaleList"
              :key="idx"
              @tap="goGoodsDetail(item)"
            >
              <image class="nf-presale-img" :src="getUrl(item.good && item.good.imageUrl)" mode="aspectFill" />
              <view class="nf-presale-info">
                <text class="nf-presale-name">{{ $lt(item.good && item.good.title) }}</text>
                <text class="nf-presale-price">¥{{ formatPrice(item.presalePrice || (item.good && item.good.price)) }}</text>
              </view>
              <view class="nf-presale-badge-tag">{{ $t('presale') }}</view>
            </view>
          </view>
        </scroll-view>
      </view>

      <!-- 分类导航 (带id用于吸顶检测) -->
      <view id="nf-category-anchor"></view>
      <categories
        ref="categoriesRef"
        :categoriesData="gridList"
        v-model="activeCategoryID"
        @change="onCategoryChange"
      ></categories>

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
      <view class="nf-footer" v-if="flowData.length > 0">
        <text class="nf-footer-text">{{ isBottom ? $t('reachedBottom') : $t('loading') }}</text>
      </view>
      <!-- 空状态 -->
      <view class="nf-empty" v-if="flowData.length === 0 && !isLoading">
        <view class="nf-empty-icon">
          <uni-icons type="shop" size="48" color="rgba(229,9,20,0.4)" />
        </view>
        <text class="nf-empty-text">{{ $t('noGoods') }}</text>
      </view>
    </scroll-view>

    <!-- 5. 吸顶分类栏 (当原始分类滚出视口时显示) -->
    <view class="nf-sticky-category" v-if="showStickyCategory">
      <categories
        :categoriesData="gridList"
        v-model="activeCategoryID"
        @change="onCategoryChange"
      ></categories>
    </view>

    <!-- 5. 回到顶部按钮 -->
    <view class="nf-back-top" v-if="showBackTop" @tap="scrollToTop">
      <uni-icons type="up" size="20" color="#fff" />
      <text class="nf-back-top-text">{{ $t('backToTop') }}</text>
    </view>

    <!-- 6. 签到悬浮按钮 -->
    <view class="nf-sign-float" v-if="signInEnabled" @tap="goSignIn">
      <text class="nf-sign-float-text">{{ $t('signIn') }}</text>
    </view>

    <!-- 语言切换弹窗 -->
    <lang-switch v-model="showLangPicker" />

    <!-- 首页弹窗 -->
    <popup-modal position="home" client-type="uni" />
  </view>
</template>
<script setup>
import  { ref, computed, onUnmounted } from 'vue';
import { getCategoryMobile, getGoodList } from '@/api/homePage.js'
import { getBannerList } from '@/api/homePage.js'
import { getPresaleGoodList } from '@/api/presale.js'
import { getAnnouncementConfig, getSignInEnabled, getPresaleHomeCount } from '@/api/sysConfig.js'
import { doSignIn, getSignInStatus } from '@/api/signIn.js'
import noPaginRowGoodList from '@/components/good-list/no-pagin-row-good-list.vue'
import swpiers from './components/swiper.vue'
import categories from './components/categories.vue'
import announcementMarquee from '@/components/announcement-marquee/announcement-marquee.vue'
import langSwitch from '@/components/lang-switch/lang-switch.vue'
import popupModal from '@/components/popup-modal/popup-modal.vue'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useUserStore } from '@/pinia/modules/user.js'
import { getUrl } from '@/utils/url.js'
import { onShow } from '@dcloudio/uni-app'

const langStore = useLangStore()
const userStore = useUserStore()
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)
const langLabel = computed(() => {
  const map = { zh: '中', en: 'EN', mn: 'MN', 'zh-TW': '繁', th: 'ไทย', hi: 'हि', id: 'ID' }
  return map[langStore.locale] || '中'
})
const showLangPicker = ref(false)

// 启动时恢复 tabBar 语言
langStore.updateTabBar(langStore.locale)

// ========== 轮播图 ==========
const list = ref([])
const searchKeyword = ref('')

const initBanner = async () => {
  const res = await getBannerList()
  list.value = res.data.list
}
initBanner()

// ========== 1. 公告走马灯 ==========
const announcementConfig = ref({
  enabled: false,
  text: '',
  textColor: '#fff',
  speed: 60
})

const initAnnouncement = async () => {
  try {
    const res = await getAnnouncementConfig()
    if (res.code === 0 && res.data) {
      const d = res.data
      announcementConfig.value = {
        enabled: d.enabled === true || d.enabled === 'true',
        text: typeof d.content === 'string' ? (langStore.$lt(d.content) || d.content) : (d.content || ''),
        textColor: d.textColor || '#fff',
        speed: parseInt(d.speed) || 60
      }
    }
  } catch (e) {
    console.error('获取公告失败', e)
  }
}
initAnnouncement()

// ========== 2. 预售商品 ==========
const presaleList = ref([])
const presaleHomeCount = ref(4)

const initPresale = async () => {
  try {
    // 获取后台配置的首页展示数量
    const countRes = await getPresaleHomeCount()
    if (countRes.code === 0 && countRes.data && countRes.data.configValue) {
      presaleHomeCount.value = parseInt(countRes.data.configValue) || 4
    }
  } catch (e) {}

  try {
    const res = await getPresaleGoodList({ page: 1, pageSize: presaleHomeCount.value })
    if (res.code === 0 && res.data.list) {
      presaleList.value = res.data.list
    }
  } catch (e) {
    console.error('获取预售商品失败', e)
  }
}
initPresale()

const formatPrice = (priceInCents) => {
  if (!priceInCents && priceInCents !== 0) return '0.00'
  return (parseInt(priceInCents) / 100).toFixed(2)
}

const goPresaleList = () => {
  uni.navigateTo({ url: '/pages/presale/list' })
}

const goGoodsDetail = (item) => {
  if (item.good && item.good.ID) {
    uni.navigateTo({ url: '/pages/goodsDetails/goodsDetails?id=' + item.good.ID })
  }
}

// ========== 5. 滚动状态：吸顶分类 + 回到顶部 ==========
const showStickyCategory = ref(false)
const showBackTop = ref(false)
const scrollTopVal = ref(0)
let categoryAnchorTop = 0
let oldScrollTop = 0

// 获取分类栏位置 (延迟到渲染完成)
const getCategoryAnchorTop = () => {
  const query = uni.createSelectorQuery()
  query.select('#nf-category-anchor').boundingClientRect((rect) => {
    if (rect) categoryAnchorTop = rect.top + oldScrollTop
  }).exec()
}

const onScroll = (e) => {
  const scrollTop = e.detail.scrollTop
  oldScrollTop = scrollTop
  // 分类区是否滚出视口
  if (categoryAnchorTop > 0) {
    showStickyCategory.value = scrollTop > categoryAnchorTop
  }
  // 回到顶部按钮
  showBackTop.value = scrollTop > 600
}

const scrollToTop = () => {
  scrollTopVal.value = oldScrollTop // 先设一个非0值
  setTimeout(() => {
    scrollTopVal.value = 0
  }, 50)
}

// ========== 6. 签到悬浮按钮 ==========
const signInEnabled = ref(false)

const initSignIn = async () => {
  try {
    const res = await getSignInEnabled()
    if (res.code === 0 && res.data && res.data.configValue) {
      signInEnabled.value = res.data.configValue === 'true'
    }
  } catch (e) {}
}
initSignIn()

const goSignIn = async () => {
  const token = userStore.token || ''
  if (!token) {
    uni.showToast({ title: $t.value('loginFirst'), icon: 'none' })
    uni.navigateTo({ url: '/pages/user/login' })
    return
  }
  try {
    // 检查是否已签到
    const statusRes = await getSignInStatus()
    if (statusRes.code === 0 && statusRes.data && statusRes.data.signedToday) {
      uni.showToast({ title: $t.value('alreadySigned'), icon: 'none' })
      return
    }
    const res = await doSignIn()
    if (res.code === 0) {
      const points = res.data && res.data.points ? res.data.points : 0
      uni.showToast({
        title: $t.value('signInSuccess') + (points > 0 ? ` +${points}` : ''),
        icon: 'success'
      })
    } else {
      uni.showToast({ title: res.msg || $t.value('alreadySigned'), icon: 'none' })
    }
  } catch (e) {
    uni.showToast({ title: $t.value('alreadySigned'), icon: 'none' })
  }
}

// ========== 商品列表 ==========
const products = ref([])
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

searchKeyword.value = ''

const goToSearchWithKeyword = () => {
  const keyword = searchKeyword.value.trim()
  if (!keyword) {
    uni.showToast({ title: $t.value('searchEmpty'), icon: 'none' })
    return
  }
  uni.setStorageSync('searchKeyword', keyword)
  uni.navigateTo({ url: '/pages/search/index' })
}

const getRecommend = async () => {
  const res = await getGoodList({ recommend: true })
  if (res.code === 0 && res.data.list.length) {
    products.value = res.data.list
  }
}
getRecommend()

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
      if (isRefresh) {
        flowData.value = listData
      } else {
        flowData.value.push(...listData)
      }
      totalCount.value = res.data.total ?? flowData.value.length
      isBottom.value = listData.length < params.value.pageSize
    } else {
      isBottom.value = true
    }
  } catch (error) {
    console.error('获取商品列表失败', error)
  } finally {
    isLoading.value = false
  }
}

const totalCount = ref(0)

const handleAutoLoadMore = () => {
  lower(false)
}

// 分类切换（带淡入淡出）
const onCategoryChange = async (categoryID) => {
  params.value.categoryID = categoryID || 0
  switching.value = true
  await new Promise(r => setTimeout(r, 220))
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
  switching.value = false
}

const gridList = ref([])
const initCategory = async () => {
  const res = await getCategoryMobile()
  if (res.code === 0 && res.data.length) {
    gridList.value = res.data
  }
}
initCategory()
lower(true)

// 延迟获取分类锚点位置
onShow(() => {
  setTimeout(getCategoryAnchorTop, 500)
})

const debounce = (func, delay) => {
  let debounceTimer;
  return function(...args) {
    if (debounceTimer) clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      func.apply(this, args);
    }, delay);
  };
};

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
.nf-footer {
  padding: 40rpx;
  text-align: center;
}

.nf-footer-text {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.25);
  letter-spacing: 2rpx;
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

/* ===== 预售区域 ===== */
.nf-presale-home {
  margin: 16rpx 20rpx 0;
}

.nf-presale-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.nf-presale-title {
  font-size: 28rpx;
  color: #fff;
  font-weight: 700;
}

.nf-presale-more {
  display: flex;
  align-items: center;
  gap: 4rpx;
}

.nf-presale-more-text {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.5);
}

.nf-presale-scroll {
  white-space: nowrap;
}

.nf-presale-items {
  display: flex;
  gap: 16rpx;
  padding-bottom: 8rpx;
}

.nf-presale-item {
  flex-shrink: 0;
  width: 240rpx;
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.06);
  border-radius: 16rpx;
  overflow: hidden;
  position: relative;

  &:active {
    transform: scale(0.97);
  }
}

.nf-presale-img {
  width: 240rpx;
  height: 240rpx;
}

.nf-presale-info {
  padding: 10rpx 12rpx;
}

.nf-presale-name {
  font-size: 22rpx;
  color: #fff;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.nf-presale-price {
  font-size: 26rpx;
  color: #e50914;
  font-weight: 700;
  display: block;
  margin-top: 4rpx;
}

.nf-presale-badge-tag {
  position: absolute;
  top: 8rpx;
  left: 8rpx;
  background: linear-gradient(135deg, #e50914, #b20710);
  color: #fff;
  font-size: 18rpx;
  padding: 2rpx 12rpx;
  border-radius: 8rpx;
  font-weight: 600;
}

/* ===== 吸顶分类栏 ===== */
.nf-sticky-category {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 98;
  background: rgba(0, 0, 0, 0.95);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
  padding-top: calc(var(--status-bar-height, 0px) + 96rpx);
}

/* ===== 回到顶部按钮 ===== */
.nf-back-top {
  position: fixed;
  right: 28rpx;
  bottom: 200rpx;
  width: 88rpx;
  height: 88rpx;
  border-radius: 50%;
  background: rgba(229, 9, 20, 0.85);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4rpx 20rpx rgba(229, 9, 20, 0.4);
  z-index: 97;
  transition: all 0.3s;

  &:active {
    transform: scale(0.9);
  }
}

.nf-back-top-text {
  font-size: 18rpx;
  color: #fff;
  margin-top: 2rpx;
}

/* ===== 签到悬浮按钮 ===== */
.nf-sign-float {
  position: fixed;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
  background: linear-gradient(135deg, #e50914, #b20710);
  padding: 16rpx 12rpx;
  border-radius: 16rpx 0 0 16rpx;
  z-index: 96;
  box-shadow: -4rpx 0 16rpx rgba(229, 9, 20, 0.4);
  writing-mode: vertical-rl;

  &:active {
    transform: translateY(-50%) scale(0.95);
  }
}

.nf-sign-float-text {
  font-size: 22rpx;
  color: #fff;
  font-weight: 700;
  letter-spacing: 4rpx;
}
</style>

