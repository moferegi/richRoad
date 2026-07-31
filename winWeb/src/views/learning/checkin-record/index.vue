<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { getCheckinStats, doCheckin, getLearningCheckinRecordList } from '@/api/learning'

const { t } = useI18n()

// 统计数据
const stats = ref(null)
const loadingStats = ref(true)

// 记录列表
const recordList = ref([])
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const finished = ref(false)

// 打卡中
const checkingIn = ref(false)

// 加载统计
async function loadStats() {
  loadingStats.value = true
  try {
    const res = await getCheckinStats()
    if (res.code === 0) {
      stats.value = res.data
    }
  } catch (e) {
    console.warn('load checkin stats failed', e)
  } finally {
    loadingStats.value = false
  }
}

// 加载记录列表
async function loadRecords(reset = false) {
  if (loading.value) return
  if (reset) {
    page.value = 1
    finished.value = false
    recordList.value = []
  }
  if (finished.value) return

  loading.value = true
  try {
    const res = await getLearningCheckinRecordList({
      page: page.value,
      pageSize: pageSize.value
    })
    if (res.code === 0) {
      const list = res.data?.list || res.data || []
      recordList.value = reset ? list : [...recordList.value, ...list]
      if (list.length < pageSize.value) {
        finished.value = true
      } else {
        page.value++
      }
    }
  } catch (e) {
    console.warn('load checkin records failed', e)
  } finally {
    loading.value = false
  }
}

// 执行打卡
async function handleCheckin() {
  if (checkingIn.value || stats.value?.todayChecked) return

  checkingIn.value = true
  try {
    const res = await doCheckin()
    if (res.code === 0) {
      // 刷新统计和记录
      await loadStats()
      loadRecords(true)
    }
  } catch (e) {
    console.warn('checkin failed', e)
  } finally {
    checkingIn.value = false
  }
}

// 格式化日期
function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return dateStr
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

// 获取奖励积分
function getRewardPoints(item) {
  return item.points || item.rewardPoints || item.reward || 0
}

onMounted(() => {
  loadStats()
  loadRecords(true)
})
</script>

<template>
  <div class="checkin-page">
    <div class="checkin-content">
      <!-- 顶部标题 -->
      <div class="page-header">
        <h1 class="page-title">{{ t('checkinRecord.title') }}</h1>
      </div>

      <!-- 统计卡片 -->
      <div class="stats-card">
        <div class="stats-decor decor-1"></div>
        <div class="stats-decor decor-2"></div>

        <div class="stats-row">
          <div class="stat-item">
            <div class="stat-value">{{ loadingStats ? '--' : (stats?.continuousDays || stats?.continuous || 0) }}</div>
            <div class="stat-label">{{ t('checkinRecord.continuousDays') }}</div>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <div class="stat-value">{{ loadingStats ? '--' : (stats?.totalDays || stats?.total || 0) }}</div>
            <div class="stat-label">{{ t('checkinRecord.totalDays') }}</div>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <div class="stat-today" :class="{ checked: stats?.todayChecked }">
              <svg v-if="stats?.todayChecked" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                <polyline points="20 6 9 17 4 12"></polyline>
              </svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <polyline points="12 6 12 12 16 14"></polyline>
              </svg>
            </div>
            <div class="stat-label">{{ t('checkinRecord.todayStatus') }}</div>
          </div>
        </div>

        <!-- 打卡按钮 -->
        <button
          class="checkin-btn"
          :class="{ checked: stats?.todayChecked }"
          :disabled="stats?.todayChecked || checkingIn"
          @click="handleCheckin"
        >
          <svg v-if="checkingIn" class="btn-spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10" stroke-dasharray="40"></circle>
          </svg>
          <template v-else-if="stats?.todayChecked">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            {{ t('checkinRecord.checked') }}
          </template>
          <template v-else>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
            </svg>
            {{ t('checkinRecord.checkin') }}
          </template>
        </button>
      </div>

      <!-- 打卡记录 -->
      <div class="records-section">
        <h3 class="section-title">打卡记录</h3>

        <div class="record-list">
          <div
            v-for="(item, index) in recordList"
            :key="item.id || item.ID || item.date || index"
            class="record-item"
          >
            <div class="record-left">
              <div class="record-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                  <line x1="16" y1="2" x2="16" y2="6"></line>
                  <line x1="8" y1="2" x2="8" y2="6"></line>
                  <line x1="3" y1="10" x2="21" y2="10"></line>
                </svg>
              </div>
              <div class="record-info">
                <div class="record-date">{{ formatDate(item.date || item.checkinDate || item.createTime) }}</div>
                <div class="record-desc">每日打卡</div>
              </div>
            </div>
            <div class="record-reward">
              <span class="reward-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"></circle>
                  <path d="M12 6v12M8 10h8M8 14h8"></path>
                </svg>
              </span>
              <span class="reward-points">+{{ getRewardPoints(item) }}</span>
            </div>
          </div>

          <!-- 骨架屏 -->
          <div v-for="n in (loading && recordList.length === 0 ? 5 : 0)" :key="'skeleton-' + n" class="record-item skeleton">
            <div class="skeleton-line title-line"></div>
            <div class="skeleton-line short-line"></div>
          </div>
        </div>

        <!-- 加载状态 -->
        <div v-if="loading && recordList.length > 0" class="load-more">
          <div class="loading-spinner"></div>
          <span>{{ t('common.loading') }}</span>
        </div>

        <!-- 加载更多 -->
        <div v-else-if="!finished && recordList.length > 0" class="load-more">
          <button class="btn btn-outline" @click="loadRecords()">
            加载更多
          </button>
        </div>

        <!-- 空状态 -->
        <div v-else-if="!loading && recordList.length === 0" class="empty-state">
          <div class="empty-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
              <line x1="16" y1="2" x2="16" y2="6"></line>
              <line x1="8" y1="2" x2="8" y2="6"></line>
              <line x1="3" y1="10" x2="21" y2="10"></line>
            </svg>
          </div>
          <p>{{ t('common.empty') }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use "@/assets/styles/variables.scss" as *;

.checkin-page {
  min-height: 100%;
}

.checkin-content {
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

// 统计卡片
.stats-card {
  position: relative;
  background: $gradient-hero;
  border-radius: $radius-xl;
  padding: $spacing-xl $spacing-xxl;
  margin-bottom: $spacing-xl;
  overflow: hidden;
  color: $text-white;

  .stats-decor {
    position: absolute;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.08);

    &.decor-1 {
      width: 160px;
      height: 160px;
      top: -60px;
      right: 10%;
    }

    &.decor-2 {
      width: 100px;
      height: 100px;
      bottom: -30px;
      left: 10%;
    }
  }
}

.stats-row {
  display: flex;
  align-items: center;
  justify-content: space-around;
  margin-bottom: $spacing-xl;
  position: relative;
  z-index: 1;
}

.stat-item {
  text-align: center;
  flex: 1;

  .stat-value {
    font-size: 36px;
    font-weight: $font-weight-bold;
    line-height: 1.2;
    margin-bottom: 6px;
  }

  .stat-label {
    font-size: 13px;
    opacity: 0.85;
  }

  .stat-today {
    width: 44px;
    height: 44px;
    margin: 0 auto 6px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(255, 255, 255, 0.2);
    transition: all 0.3s ease;

    svg {
      width: 24px;
      height: 24px;
    }

    &.checked {
      background: rgba(52, 199, 89, 0.9);
    }
  }
}

.stat-divider {
  width: 1px;
  height: 40px;
  background: rgba(255, 255, 255, 0.2);
}

.checkin-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 14px;
  background: $text-white;
  color: $primary-color;
  border: none;
  border-radius: $radius-pill;
  font-size: 16px;
  font-weight: $font-weight-bold;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
  z-index: 1;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);

  svg {
    width: 20px;
    height: 20px;
  }

  &:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.15);
  }

  &.checked {
    background: rgba(255, 255, 255, 0.3);
    color: $text-white;
    cursor: default;
    box-shadow: none;
  }

  &:disabled {
    cursor: not-allowed;
  }

  .btn-spinner {
    animation: spin 0.8s linear infinite;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

// 记录区域
.records-section {
  .section-title {
    font-size: $font-size-lg;
    font-weight: $font-weight-semibold;
    color: $text-primary;
    margin: 0 0 $spacing-md 0;
  }
}

.record-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.record-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 18px;
  background: $bg-card;
  border-radius: $radius-lg;
  box-shadow: $shadow-card;
  border: 1px solid $border-light;
  transition: all 0.2s ease;

  &:hover {
    background: $bg-hover;
  }

  &.skeleton {
    flex-direction: column;
    align-items: stretch;
    gap: 10px;
    pointer-events: none;
  }
}

.record-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.record-icon {
  width: 40px;
  height: 40px;
  border-radius: $radius-md;
  background: linear-gradient(135deg, rgba(109, 91, 255, 0.12), rgba(155, 143, 255, 0.08));
  display: flex;
  align-items: center;
  justify-content: center;
  color: $primary-color;
  flex-shrink: 0;

  svg {
    width: 20px;
    height: 20px;
  }
}

.record-info {
  .record-date {
    font-size: 14px;
    font-weight: $font-weight-medium;
    color: $text-primary;
    margin-bottom: 2px;
  }

  .record-desc {
    font-size: 12px;
    color: $text-tertiary;
  }
}

.record-reward {
  display: flex;
  align-items: center;
  gap: 6px;

  .reward-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    color: #FFB800;

    svg {
      width: 18px;
      height: 18px;
    }
  }

  .reward-points {
    font-size: 16px;
    font-weight: $font-weight-bold;
    color: #FFB800;
  }
}

.skeleton-line {
  height: 14px;
  background: linear-gradient(90deg, $bg-input 25%, $border-light 50%, $bg-input 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;

  &.title-line {
    width: 60%;
  }

  &.short-line {
    width: 30%;
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
  padding: 24px 0;
  color: $text-secondary;
  font-size: 14px;

  .loading-spinner {
    width: 20px;
    height: 20px;
    border: 2px solid $border-color;
    border-top-color: $primary-color;
    border-radius: 50%;
    animation: spin2 0.8s linear infinite;
  }
}

@keyframes spin2 {
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
  padding: 60px 0;
  color: $text-tertiary;

  .empty-icon {
    width: 64px;
    height: 64px;
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
</style>
