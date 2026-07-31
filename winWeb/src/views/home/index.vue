<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getBannerList } from '@/api/homePage'
import { getVideoCategoryList, getVideoSeriesList } from '@/api/learning'
import { localText } from '@/utils/i18n'

const { t } = useI18n()
const router = useRouter()

// Banner
const bannerList = ref([])
const currentBanner = ref(0)
let bannerTimer = null

// 分类
const categories = ref([])
const currentCategoryId = ref(0) // 0 表示全部

// 视频列表
const videoList = ref([])
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const finished = ref(false)

// 加载 Banner
async function loadBanners() {
  try {
    const res = await getBannerList()
    if (res.code === 0) {
      bannerList.value = res.data?.list || res.data || []
    }
  } catch (e) {
    console.warn('load banners failed', e)
  }
}

// 加载分类
async function loadCategories() {
  try {
    const res = await getVideoCategoryList()
    if (res.code === 0) {
      categories.value = res.data?.list || res.data || []
    }
  } catch (e) {
    console.warn('load categories failed', e)
  }
}

// 加载视频列表
async function loadVideos(reset = false) {
  if (loading.value) return
  if (reset) {
    page.value = 1
    finished.value = false
    videoList.value = []
  }
  if (finished.value) return

  loading.value = true
  try {
    const params = {
      page: page.value,
      pageSize: pageSize.value
    }
    if (currentCategoryId.value > 0) {
      params.categoryId = currentCategoryId.value
    }
    const res = await getVideoSeriesList(params)
    if (res.code === 0) {
      const list = res.data?.list || res.data || []
      videoList.value = reset ? list : [...videoList.value, ...list]
      if (list.length < pageSize.value) {
        finished.value = true
      } else {
        page.value++
      }
    }
  } catch (e) {
    console.warn('load videos failed', e)
  } finally {
    loading.value = false
  }
}

// 切换分类
function switchCategory(id) {
  if (currentCategoryId.value === id) return
  currentCategoryId.value = id
  loadVideos(true)
}

// Banner 自动轮播
function startBannerAutoPlay() {
  if (bannerList.value.length <= 1) return
  stopBannerAutoPlay()
  bannerTimer = setInterval(() => {
    currentBanner.value = (currentBanner.value + 1) % bannerList.value.length
  }, 4000)
}

function stopBannerAutoPlay() {
  if (bannerTimer) {
    clearInterval(bannerTimer)
    bannerTimer = null
  }
}

function goToBanner(index) {
  currentBanner.value = index
}

// 跳转视频详情
function goVideoDetail(id) {
  router.push(`/learning/video-detail/${id}`)
}

// 跳转标签筛选
function goTagFilter() {
  router.push('/learning/tag-filter')
}

// 格式化数字
function formatNumber(num) {
  if (!num) return 0
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w'
  }
  return num
}

const allCategoryLabel = computed(() => t('common.all'))

onMounted(() => {
  loadBanners()
  loadCategories()
  loadVideos(true)
})
</script>

<template>
  <div class="home-page">
    <!-- 顶部 Hero 区域（装饰性） -->
    <div class="hero-banner">
      <div class="hero-decor decor-1"></div>
      <div class="hero-decor decor-2"></div>
      <div class="hero-content">
        <h2 class="hero-title">RichRoad Learning</h2>
        <p class="hero-subtitle">{{ t('home.slogan') }}</p>
      </div>
    </div>

    <div class="home-content">
      <!-- Banner 轮播 -->
      <div v-if="bannerList.length > 0" class="banner-section">
        <div
          class="banner-wrapper"
          @mouseenter="stopBannerAutoPlay"
          @mouseleave="startBannerAutoPlay"
        >
          <div class="banner-track" :style="{ transform: `translateX(-${currentBanner * 100}%)` }">
            <div
              v-for="(banner, index) in bannerList"
              :key="banner.ID || banner.id || index"
              class="banner-slide"
            >
              <img
                v-if="banner.imageUrl || banner.imgUrl || banner.pic"
                :src="banner.imageUrl || banner.imgUrl || banner.pic"
                :alt="banner.title || ''"
                class="banner-img"
              />
              <div v-else class="banner-placeholder">
                <span>{{ banner.title || 'Banner ' + (index + 1) }}</span>
              </div>
            </div>
          </div>
          <div class="banner-dots">
            <span
              v-for="(_, i) in bannerList"
              :key="i"
              class="banner-dot"
              :class="{ active: i === currentBanner }"
              @click="goToBanner(i)"
            ></span>
          </div>
        </div>
      </div>

      <!-- 分类栏 -->
      <div class="category-section">
        <div class="section-header">
          <h3 class="section-title">{{ t('home.categories') }}</h3>
          <a class="more-tags" @click="goTagFilter">
            {{ t('home.moreTags') }}
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
          </a>
        </div>

        <div class="category-list">
          <div
            class="category-item"
            :class="{ active: currentCategoryId === 0 }"
            @click="switchCategory(0)"
          >
            <span class="category-name">{{ allCategoryLabel }}</span>
          </div>
          <div
            v-for="cat in categories"
            :key="cat.ID || cat.id"
            class="category-item"
            :class="{ active: currentCategoryId === (cat.ID || cat.id) }"
            @click="switchCategory(cat.ID || cat.id)"
          >
            <span class="category-name">{{ localText(cat.name) }}</span>
          </div>
        </div>
      </div>

      <!-- 视频网格 -->
      <div class="video-section">
        <div class="video-grid">
          <div
            v-for="video in videoList"
            :key="video.ID || video.id"
            class="video-card"
            @click="goVideoDetail(video.ID || video.id)"
          >
            <div class="video-cover">
              <img
                v-if="video.coverUrl || video.cover"
                :src="video.coverUrl || video.cover"
                :alt="localText(video.name)"
                class="cover-img"
              />
              <div v-else class="cover-placeholder">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polygon points="23 7 16 12 23 17 23 7"></polygon>
                  <rect x="1" y="5" width="15" height="14" rx="2" ry="2"></rect>
                </svg>
              </div>
              <div class="cover-overlay">
                <div class="play-icon">
                  <svg viewBox="0 0 24 24" fill="currentColor">
                    <polygon points="5 3 19 12 5 21 5 3"></polygon>
                  </svg>
                </div>
              </div>
            </div>
            <div class="video-info">
              <h4 class="video-title">{{ localText(video.name) }}</h4>
              <div class="video-meta">
                <span class="meta-item">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                    <circle cx="12" cy="12" r="3"></circle>
                  </svg>
                  {{ formatNumber(video.viewCount) }} {{ t('home.viewCount') }}
                </span>
                <span class="meta-item">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
                    <circle cx="9" cy="7" r="4"></circle>
                    <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
                    <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
                  </svg>
                  {{ formatNumber(video.userCount) }} {{ t('home.userCount') }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="load-more">
          <div class="loading-spinner"></div>
          <span>{{ t('common.loading') }}</span>
        </div>

        <!-- 加载更多 -->
        <div v-else-if="!finished && videoList.length > 0" class="load-more">
          <button class="btn btn-outline" @click="loadVideos()">
            加载更多
          </button>
        </div>

        <!-- 空状态 -->
        <div v-else-if="!loading && videoList.length === 0" class="empty-state">
          <div class="empty-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <rect x="2" y="4" width="20" height="16" rx="2"></rect>
              <path d="M10 9l5 3-5 3V9z" fill="currentColor"></path>
            </svg>
          </div>
          <p>{{ t('common.empty') }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.home-page {
  min-height: 100%;
}

// ========== Hero 区域 ==========
.hero-banner {
  position: relative;
  height: 200px;
  background: $gradient-hero;
  border-radius: 0 0 $radius-xl $radius-xl;
  overflow: hidden;
  margin-bottom: $spacing-xl;

  .hero-decor {
    position: absolute;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.1);

    &.decor-1 {
      width: 200px;
      height: 200px;
      top: -60px;
      right: 10%;
    }

    &.decor-2 {
      width: 120px;
      height: 120px;
      bottom: -30px;
      left: 15%;
      background: rgba(255, 255, 255, 0.15);
    }
  }

  .hero-content {
    position: relative;
    z-index: 1;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    padding: 0 $spacing-xxl;
    color: $text-white;

    .hero-title {
      font-size: 32px;
      font-weight: $font-weight-bold;
      margin-bottom: 8px;
      letter-spacing: 1px;
    }

    .hero-subtitle {
      font-size: 15px;
      opacity: 0.9;
    }
  }
}

.home-content {
  max-width: $content-max-width;
  margin: 0 auto;
  padding: 0 $spacing-xl;
}

// ========== Banner 轮播 ==========
.banner-section {
  margin-bottom: $spacing-xl;
}

.banner-wrapper {
  position: relative;
  border-radius: $radius-xl;
  overflow: hidden;
  box-shadow: $shadow-lg;
  cursor: pointer;
}

.banner-track {
  display: flex;
  transition: transform 0.6s ease;
}

.banner-slide {
  flex-shrink: 0;
  width: 100%;
  height: 280px;
  position: relative;

  .banner-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .banner-placeholder {
    width: 100%;
    height: 100%;
    background: $gradient-card;
    display: flex;
    align-items: center;
    justify-content: center;
    color: $primary-light;
    font-size: 20px;
    font-weight: $font-weight-medium;
  }
}

.banner-dots {
  position: absolute;
  bottom: 16px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 8px;

  .banner-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.5);
    cursor: pointer;
    transition: all 0.3s ease;

    &.active {
      width: 24px;
      border-radius: 4px;
      background: $text-white;
    }
  }
}

// ========== 分类 ==========
.category-section {
  margin-bottom: $spacing-xl;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: $spacing-lg;

  .section-title {
    font-size: $font-size-xl;
    font-weight: $font-weight-bold;
    color: $text-primary;
  }

  .more-tags {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: $font-size-sm;
    color: $primary-color;
    cursor: pointer;
    transition: all 0.2s ease;

    &:hover {
      opacity: 0.8;
    }

    svg {
      width: 16px;
      height: 16px;
    }
  }
}

.category-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.category-item {
  padding: 10px 22px;
  background: $bg-card;
  border: 1px solid $border-color;
  border-radius: $radius-pill;
  font-size: 14px;
  color: $text-secondary;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;

  &:hover {
    border-color: $primary-light;
    color: $primary-color;
  }

  &.active {
    background: $gradient-primary;
    border-color: transparent;
    color: $text-white;
    font-weight: $font-weight-medium;
    box-shadow: 0 4px 12px rgba(109, 91, 255, 0.3);
  }

  .category-name {
    display: block;
  }
}

// ========== 视频网格 ==========
.video-section {
  padding-bottom: $spacing-xxl;
}

.video-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
}

.video-card {
  background: $bg-card;
  border-radius: $radius-lg;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: $shadow-card;
  border: 1px solid $border-light;

  &:hover {
    transform: translateY(-4px);
    box-shadow: $shadow-lg;

    .cover-overlay {
      opacity: 1;
    }

    .video-title {
      color: $primary-color;
    }
  }
}

.video-cover {
  position: relative;
  width: 100%;
  padding-top: 56.25%; // 16:9
  background: $bg-input;
  overflow: hidden;

  .cover-img {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.3s ease;
  }

  .cover-placeholder {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    color: $border-primary;

    svg {
      width: 48px;
      height: 48px;
    }
  }

  .cover-overlay {
    position: absolute;
    inset: 0;
    background: rgba(26, 27, 58, 0.4);
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: opacity 0.3s ease;

    .play-icon {
      width: 56px;
      height: 56px;
      background: $gradient-primary;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      color: $text-white;
      box-shadow: 0 4px 20px rgba(109, 91, 255, 0.5);

      svg {
        width: 24px;
        height: 24px;
        margin-left: 3px;
      }
    }
  }
}

.video-info {
  padding: 16px;

  .video-title {
    font-size: 15px;
    font-weight: $font-weight-medium;
    color: $text-primary;
    line-height: 1.4;
    margin-bottom: 10px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    transition: color 0.2s ease;
    min-height: 42px;
  }

  .video-meta {
    display: flex;
    align-items: center;
    gap: 16px;
    font-size: 12px;
    color: $text-tertiary;

    .meta-item {
      display: flex;
      align-items: center;
      gap: 4px;

      svg {
        width: 14px;
        height: 14px;
      }
    }
  }
}

// ========== 加载更多 ==========
.load-more {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  padding: 24px 0;
  color: $text-secondary;
  font-size: 14px;

  .loading-spinner {
    width: 20px;
    height: 20px;
    border: 2px solid $border-color;
    border-top-color: $primary-color;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

// ========== 空状态 ==========
.empty-state {
  text-align: center;
  padding: 60px 0;
  color: $text-tertiary;

  .empty-icon {
    width: 80px;
    height: 80px;
    margin: 0 auto 16px;
    color: $border-primary;

    svg {
      width: 100%;
      height: 100%;
    }
  }

  p {
    font-size: 14px;
  }
}

// 响应式
@media (max-width: 1200px) {
  .video-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 900px) {
  .video-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
