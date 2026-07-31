<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getWatchHistoryList } from '@/api/learning'

const { t } = useI18n()
const router = useRouter()

const historyList = ref([])
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const finished = ref(false)

// 加载观看历史
async function loadHistory(reset = false) {
  if (loading.value) return
  if (reset) {
    page.value = 1
    finished.value = false
    historyList.value = []
  }
  if (finished.value) return

  loading.value = true
  try {
    const res = await getWatchHistoryList({
      page: page.value,
      pageSize: pageSize.value
    })
    if (res.code === 0) {
      const list = res.data?.list || res.data || []
      historyList.value = reset ? list : [...historyList.value, ...list]
      if (list.length < pageSize.value) {
        finished.value = true
      } else {
        page.value++
      }
    }
  } catch (e) {
    console.warn('load watch history failed', e)
  } finally {
    loading.value = false
  }
}

// 获取系列名
function getSeriesName(item) {
  return item.seriesName || item.seriesTitle || item.seriesNameI18n || item.nameI18n || item.name || '系列视频'
}

// 获取集名
function getEpisodeName(item) {
  return item.episodeName || item.episodeTitle || item.episodeNameI18n || item.episode || `第 ${item.episodeNo || 1} 集`
}

// 获取封面
function getCoverUrl(item) {
  return item.coverUrl || item.cover || item.episodeCover || item.seriesCover || ''
}

// 进度百分比
function getProgress(item) {
  const current = item.currentTime || item.watchProgress || item.progress || 0
  const total = item.duration || item.totalDuration || 0
  if (total <= 0) return 0
  return Math.min(100, Math.round((current / total) * 100))
}

// 格式化时间
function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return dateStr
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes} 分钟前`
  if (hours < 24) return `${hours} 小时前`
  if (days < 7) return `${days} 天前`

  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

// 继续观看
function continueWatch(item, e) {
  e.stopPropagation()
  const id = item.episodeId || item.targetId || item.id || item.ID
  if (id) {
    router.push(`/learning/player/${id}`)
  }
}

// 点击整条
function handleItemClick(item) {
  const seriesId = item.seriesId || item.id || item.ID
  if (seriesId) {
    router.push(`/learning/video-detail/${seriesId}`)
  }
}

onMounted(() => {
  loadHistory(true)
})
</script>

<template>
  <div class="watch-history-page">
    <div class="history-content">
      <!-- 顶部标题 -->
      <div class="page-header">
        <h1 class="page-title">{{ t('watchHistory.title') }}</h1>
      </div>

      <!-- 历史列表 -->
      <div class="history-list">
        <div
          v-for="item in historyList"
          :key="item.id || item.ID || item.episodeId"
          class="history-item"
          @click="handleItemClick(item)"
        >
          <!-- 左侧封面 -->
          <div class="item-cover">
            <img
              v-if="getCoverUrl(item)"
              :src="getCoverUrl(item)"
              :alt="getSeriesName(item)"
              class="cover-img"
            />
            <div v-else class="cover-placeholder">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polygon points="23 7 16 12 23 17 23 7"></polygon>
                <rect x="1" y="5" width="15" height="14" rx="2" ry="2"></rect>
              </svg>
            </div>
            <div class="progress-bar-bottom">
              <div class="progress-fill" :style="{ width: getProgress(item) + '%' }"></div>
            </div>
            <div class="progress-badge">{{ getProgress(item) }}%</div>
          </div>

          <!-- 右侧信息 -->
          <div class="item-info">
            <h3 class="series-name">{{ getSeriesName(item) }}</h3>
            <p class="episode-name">{{ getEpisodeName(item) }}</p>
            <div class="item-meta">
              <div class="meta-row">
                <span class="meta-label">{{ t('watchHistory.progress') }}</span>
                <span class="meta-value progress-text">
                  {{ getProgress(item) }}%
                </span>
              </div>
              <div class="meta-row">
                <span class="meta-label">{{ t('watchHistory.watchedAt') }}</span>
                <span class="meta-value">
                  {{ formatDate(item.watchTime || item.lastWatchTime || item.updatedAt || item.createTime) }}
                </span>
              </div>
            </div>
            <button class="continue-btn" @click="continueWatch(item, $event)">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <polygon points="5 3 19 12 5 21 5 3"></polygon>
              </svg>
              {{ t('watchHistory.continueWatch') }}
            </button>
          </div>
        </div>

        <!-- 骨架屏 -->
        <div v-for="n in (loading && historyList.length === 0 ? 5 : 0)" :key="'skeleton-' + n" class="history-item skeleton">
          <div class="item-cover skeleton-block"></div>
          <div class="item-info">
            <div class="skeleton-line title-line"></div>
            <div class="skeleton-line"></div>
            <div class="skeleton-line short-line"></div>
          </div>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading && historyList.length > 0" class="load-more">
        <div class="loading-spinner"></div>
        <span>{{ t('common.loading') }}</span>
      </div>

      <!-- 加载更多 -->
      <div v-else-if="!finished && historyList.length > 0" class="load-more">
        <button class="btn btn-outline" @click="loadHistory()">
          加载更多
        </button>
      </div>

      <!-- 空状态 -->
      <div v-else-if="!loading && historyList.length === 0" class="empty-state">
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
</template>

<style lang="scss" scoped>
@use "@/assets/styles/variables.scss" as *;

.watch-history-page {
  min-height: 100%;
}

.history-content {
  max-width: $content-max-width;
  margin: 0 auto;
  padding: $spacing-xl;
}

// 头部
.page-header {
  margin-bottom: $spacing-lg;

  .page-title {
    font-size: $font-size-title;
    font-weight: $font-weight-bold;
    color: $text-primary;
    margin: 0;
  }
}

// 历史列表
.history-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.history-item {
  display: flex;
  gap: 16px;
  padding: 16px;
  background: $bg-card;
  border-radius: $radius-lg;
  cursor: pointer;
  transition: all 0.25s ease;
  box-shadow: $shadow-card;
  border: 1px solid $border-light;

  &:hover {
    transform: translateY(-2px);
    box-shadow: $shadow-md;

    .series-name {
      color: $primary-color;
    }
  }

  &.skeleton {
    cursor: default;
    pointer-events: none;
  }
}

.item-cover {
  width: 200px;
  aspect-ratio: 16 / 9;
  border-radius: $radius-md;
  overflow: hidden;
  flex-shrink: 0;
  background: $bg-input;
  position: relative;

  .cover-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .cover-placeholder {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    color: $border-primary;

    svg {
      width: 40px;
      height: 40px;
    }
  }

  .progress-bar-bottom {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: rgba(255, 255, 255, 0.3);

    .progress-fill {
      height: 100%;
      background: $gradient-primary;
      border-radius: 0 2px 2px 0;
    }
  }

  .progress-badge {
    position: absolute;
    top: 8px;
    right: 8px;
    padding: 2px 8px;
    background: rgba(0, 0, 0, 0.6);
    color: $text-white;
    font-size: 11px;
    border-radius: $radius-pill;
    font-weight: $font-weight-medium;
    backdrop-filter: blur(4px);
  }

  &.skeleton-block {
    background: linear-gradient(90deg, $bg-input 25%, $border-light 50%, $bg-input 75%);
    background-size: 200% 100%;
    animation: shimmer 1.5s infinite;
  }
}

.item-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.series-name {
  font-size: 16px;
  font-weight: $font-weight-semibold;
  color: $text-primary;
  margin: 0 0 4px 0;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.2s ease;
}

.episode-name {
  font-size: 13px;
  color: $text-secondary;
  margin: 0 0 12px 0;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.item-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 12px;
  flex: 1;

  .meta-row {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 13px;
  }

  .meta-label {
    color: $text-tertiary;
    min-width: 60px;
  }

  .meta-value {
    color: $text-secondary;

    &.progress-text {
      color: $primary-color;
      font-weight: $font-weight-medium;
    }
  }
}

.continue-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 18px;
  background: $gradient-primary;
  color: $text-white;
  border: none;
  border-radius: $radius-pill;
  font-size: 13px;
  font-weight: $font-weight-medium;
  cursor: pointer;
  transition: all 0.2s ease;
  align-self: flex-start;
  box-shadow: 0 2px 8px rgba(109, 91, 255, 0.25);

  svg {
    width: 14px;
    height: 14px;
  }

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(109, 91, 255, 0.35);
  }
}

.skeleton-line {
  height: 14px;
  background: linear-gradient(90deg, $bg-input 25%, $border-light 50%, $bg-input 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
  margin-bottom: 8px;

  &.title-line {
    height: 18px;
    width: 50%;
  }

  &.short-line {
    width: 25%;
    margin-bottom: 0;
  }
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

// 加载更多
.load-more {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  padding: 32px 0;
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

.btn {
  padding: 10px 24px;
  border-radius: $radius-pill;
  font-size: 14px;
  font-weight: $font-weight-medium;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;

  &.btn-outline {
    background: $bg-card;
    border: 1px solid $border-color;
    color: $text-secondary;

    &:hover {
      border-color: $primary-light;
      color: $primary-color;
    }
  }
}

// 空状态
.empty-state {
  text-align: center;
  padding: 80px 0;
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
@media (max-width: 700px) {
  .history-item {
    flex-direction: column;
  }

  .item-cover {
    width: 100%;
  }
}
</style>
