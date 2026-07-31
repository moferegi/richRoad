<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { findVideoSeries, getVideoEpisodeList, getSeriesWatchProgressList } from '@/api/learning'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const seriesId = computed(() => route.params.seriesId)

// 系列信息
const seriesInfo = ref(null)
const loading = ref(true)

// 集数列表
const episodeList = ref([])

// 观看进度
const progressMap = ref({}) // { episodeId: progress }
const lastWatchedEpisodeId = ref(null)

// 加载系列详情
async function loadSeriesInfo() {
  try {
    const res = await findVideoSeries(seriesId.value)
    if (res.code === 0) {
      seriesInfo.value = res.data
    }
  } catch (e) {
    console.warn('load series info failed', e)
  }
}

// 加载集数列表
async function loadEpisodes() {
  try {
    const res = await getVideoEpisodeList({
      seriesId: seriesId.value,
      pageSize: 300
    })
    if (res.code === 0) {
      episodeList.value = res.data?.list || res.data || []
    }
  } catch (e) {
    console.warn('load episodes failed', e)
  }
}

// 加载观看进度
async function loadProgress() {
  try {
    const res = await getSeriesWatchProgressList(seriesId.value)
    if (res.code === 0) {
      const list = res.data?.list || res.data || []
      const map = {}
      let lastId = null
      let lastTime = 0
      list.forEach(item => {
        map[item.episodeId || item.ID] = item.progress || 0
        if (item.updatedAt > lastTime) {
          lastTime = item.updatedAt
          lastId = item.episodeId || item.ID
        }
      })
      progressMap.value = map
      lastWatchedEpisodeId.value = lastId
    }
  } catch (e) {
    console.warn('load progress failed', e)
  }
}

// 跳转播放器
function goPlayer(episodeId) {
  router.push(`/learning/player/${episodeId}`)
}

// 继续观看
function continueWatch() {
  const targetId = lastWatchedEpisodeId.value || (episodeList.value[0]?.ID || episodeList.value[0]?.id)
  if (targetId) {
    goPlayer(targetId)
  }
}

// 返回
function goBack() {
  router.back()
}

// 格式化数字
function formatNumber(num) {
  if (!num) return 0
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w'
  }
  return num
}

// 获取进度百分比
function getProgress(episodeId) {
  return progressMap.value[episodeId] || 0
}

// 集数标题格式化
function formatEpisodeTitle(ep, index) {
  const name = ep.nameI18n || ep.name || ''
  return name || `第 ${index + 1} 集`
}

onMounted(async () => {
  loading.value = true
  await Promise.all([loadSeriesInfo(), loadEpisodes()])
  loadProgress()
  loading.value = false
})
</script>

<template>
  <div class="video-detail-page">
    <!-- 顶部 Hero -->
    <div v-if="seriesInfo" class="detail-hero">
      <div class="hero-bg">
        <img
          v-if="seriesInfo.coverUrl || seriesInfo.cover"
          :src="seriesInfo.coverUrl || seriesInfo.cover"
          alt=""
          class="bg-img"
        />
        <div class="bg-mask"></div>
      </div>
      <div class="hero-content">
        <button class="back-btn" @click="goBack">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="15 18 9 12 15 6"></polyline>
          </svg>
          返回
        </button>

        <div class="series-main">
          <div class="series-cover">
            <img
              v-if="seriesInfo.coverUrl || seriesInfo.cover"
              :src="seriesInfo.coverUrl || seriesInfo.cover"
              :alt="seriesInfo.nameI18n || seriesInfo.name"
              class="cover-img"
            />
            <div v-else class="cover-placeholder">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polygon points="23 7 16 12 23 17 23 7"></polygon>
                <rect x="1" y="5" width="15" height="14" rx="2" ry="2"></rect>
              </svg>
            </div>
          </div>

          <div class="series-info">
            <h1 class="series-title">{{ seriesInfo.nameI18n || seriesInfo.name }}</h1>
            <p v-if="seriesInfo.description || seriesInfo.desc" class="series-desc">
              {{ seriesInfo.description || seriesInfo.desc }}
            </p>

            <div class="series-stats">
              <span class="stat-item">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                  <circle cx="12" cy="12" r="3"></circle>
                </svg>
                {{ formatNumber(seriesInfo.viewCount) }} {{ t('videoDetail.viewCount') }}
              </span>
              <span class="stat-item">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
                  <circle cx="9" cy="7" r="4"></circle>
                  <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
                  <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
                </svg>
                {{ formatNumber(seriesInfo.userCount) }} {{ t('videoDetail.userCount') }}
              </span>
              <span class="stat-item">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                  <polyline points="14 2 14 8 20 8"></polyline>
                </svg>
                共 {{ episodeList.length }} 集
              </span>
            </div>

            <div class="series-actions">
              <button class="btn btn-primary btn-lg" @click="continueWatch">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <polygon points="5 3 19 12 5 21 5 3"></polygon>
                </svg>
                {{ lastWatchedEpisodeId ? t('videoDetail.continueWatch') : '开始观看' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 集数列表 -->
    <div class="episode-section">
      <div class="section-header">
        <h3 class="section-title">{{ t('videoDetail.episodeList') }}</h3>
        <span class="episode-count">共 {{ episodeList.length }} 集</span>
      </div>

      <div v-if="loading" class="loading-wrap">
        <div class="loading-spinner"></div>
        <span>{{ t('common.loading') }}</span>
      </div>

      <div v-else class="episode-grid">
        <div
          v-for="(ep, index) in episodeList"
          :key="ep.ID || ep.id"
          class="episode-card"
          :class="{ watched: getProgress(ep.ID || ep.id) > 0 }"
          @click="goPlayer(ep.ID || ep.id)"
        >
          <div class="ep-index">{{ index + 1 }}</div>
          <div class="ep-info">
            <div class="ep-title ellipsis">{{ formatEpisodeTitle(ep, index) }}</div>
            <div v-if="getProgress(ep.ID || ep.id) > 0" class="ep-progress">
              <div class="progress-bar">
                <div class="progress-fill" :style="{ width: Math.min(getProgress(ep.ID || ep.id), 100) + '%' }"></div>
              </div>
              <span class="progress-text">{{ Math.min(getProgress(ep.ID || ep.id), 100).toFixed(0) }}%</span>
            </div>
          </div>
          <div class="ep-play">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <polygon points="5 3 19 12 5 21 5 3"></polygon>
            </svg>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.video-detail-page {
  min-height: 100%;
  max-width: $content-max-width;
  margin: 0 auto;
}

// ========== Hero 区域 ==========
.detail-hero {
  position: relative;
  border-radius: $radius-xl;
  overflow: hidden;
  margin-bottom: $spacing-xl;
  box-shadow: $shadow-lg;
}

.hero-bg {
  position: absolute;
  inset: 0;

  .bg-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    filter: blur(20px);
    transform: scale(1.1);
  }

  .bg-mask {
    position: absolute;
    inset: 0;
    background: linear-gradient(135deg, rgba(109, 91, 255, 0.85), rgba(62, 49, 179, 0.9));
  }
}

.hero-content {
  position: relative;
  z-index: 1;
  padding: $spacing-xl;
  color: $text-white;
}

.back-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  border-radius: $radius-pill;
  color: $text-white;
  font-size: 14px;
  margin-bottom: $spacing-lg;
  transition: all 0.2s ease;

  &:hover {
    background: rgba(255, 255, 255, 0.25);
  }

  svg {
    width: 16px;
    height: 16px;
  }
}

.series-main {
  display: flex;
  gap: $spacing-xl;
}

.series-cover {
  width: 220px;
  flex-shrink: 0;
  border-radius: $radius-lg;
  overflow: hidden;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);

  .cover-img {
    width: 100%;
    aspect-ratio: 16 / 10;
    object-fit: cover;
    display: block;
  }

  .cover-placeholder {
    width: 100%;
    aspect-ratio: 16 / 10;
    background: rgba(255, 255, 255, 0.1);
    display: flex;
    align-items: center;
    justify-content: center;
    color: rgba(255, 255, 255, 0.5);

    svg {
      width: 60px;
      height: 60px;
    }
  }
}

.series-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.series-title {
  font-size: 28px;
  font-weight: $font-weight-bold;
  margin-bottom: 12px;
  line-height: 1.3;
}

.series-desc {
  font-size: 14px;
  opacity: 0.85;
  line-height: 1.6;
  margin-bottom: 20px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.series-stats {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 24px;
  font-size: 13px;
  opacity: 0.9;

  .stat-item {
    display: flex;
    align-items: center;
    gap: 6px;

    svg {
      width: 16px;
      height: 16px;
    }
  }
}

.series-actions {
  .btn {
    gap: 8px;

    svg {
      width: 18px;
      height: 18px;
    }
  }
}

// ========== 集数列表 ==========
.episode-section {
  background: $bg-card;
  border-radius: $radius-xl;
  padding: $spacing-xl;
  box-shadow: $shadow-card;
  border: 1px solid $border-light;
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

  .episode-count {
    font-size: 13px;
    color: $text-tertiary;
  }
}

.loading-wrap {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  padding: 40px 0;
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

.episode-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.episode-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  background: $bg-input;
  border-radius: $radius-md;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid transparent;

  &:hover {
    background: $bg-hover;
    border-color: $border-primary;
    transform: translateX(4px);

    .ep-play {
      background: $gradient-primary;
      color: $text-white;
    }
  }

  &.watched {
    .ep-index {
      background: $gradient-primary;
      color: $text-white;
    }
  }
}

.ep-index {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
  background: $bg-card;
  border-radius: $radius-sm;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: $font-weight-semibold;
  color: $text-secondary;
}

.ep-info {
  flex: 1;
  min-width: 0;

  .ep-title {
    font-size: 14px;
    font-weight: $font-weight-medium;
    color: $text-primary;
    margin-bottom: 6px;
  }

  .ep-progress {
    display: flex;
    align-items: center;
    gap: 8px;

    .progress-bar {
      flex: 1;
      height: 4px;
      background: $border-color;
      border-radius: 2px;
      overflow: hidden;

      .progress-fill {
        height: 100%;
        background: $gradient-primary;
        border-radius: 2px;
        transition: width 0.3s ease;
      }
    }

    .progress-text {
      font-size: 11px;
      color: $text-tertiary;
      width: 32px;
      text-align: right;
    }
  }
}

.ep-play {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
  background: $bg-card;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: $primary-color;
  transition: all 0.2s ease;

  svg {
    width: 14px;
    height: 14px;
    margin-left: 2px;
  }
}

// 响应式
@media (max-width: 900px) {
  .series-main {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  .series-stats {
    justify-content: center;
  }

  .episode-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
