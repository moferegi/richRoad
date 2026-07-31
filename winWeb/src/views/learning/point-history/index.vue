<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { getLearningPointRecordList, getAsset } from '@/api/learning'

const { t } = useI18n()

const recordList = ref([])
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const finished = ref(false)
const currentPoints = ref(0)
const loadingAsset = ref(true)

// 加载当前积分
async function loadAsset() {
  loadingAsset.value = true
  try {
    const res = await getAsset()
    if (res.code === 0) {
      currentPoints.value = res.data?.points || 0
    }
  } catch (e) {
    console.warn('load asset failed', e)
  } finally {
    loadingAsset.value = false
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
    const res = await getLearningPointRecordList({
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
    console.warn('load point records failed', e)
  } finally {
    loading.value = false
  }
}

// 获取类型
function getType(item) {
  return item.type || item.typeName || item.reason || '积分变动'
}

// 获取变动值
function getChange(item) {
  return item.points || item.change || item.amount || 0
}

// 是否为增加
function isIncrease(item) {
  return getChange(item) > 0
}

// 格式化日期
function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return dateStr
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const h = String(d.getHours()).padStart(2, '0')
  const min = String(d.getMinutes()).padStart(2, '0')
  return `${y}-${m}-${day} ${h}:${min}`
}

onMounted(() => {
  loadAsset()
  loadRecords(true)
})
</script>

<template>
  <div class="record-page">
    <div class="record-content">
      <!-- 顶部标题 -->
      <div class="page-header">
        <h1 class="page-title">{{ t('pointHistory.title') }}</h1>
      </div>

      <!-- 概览卡 -->
      <div class="overview-card points-card">
        <div class="card-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"></circle>
            <path d="M12 6v12M8 10h8M8 14h8"></path>
          </svg>
        </div>
        <div class="card-info">
          <div class="card-label">{{ t('pointHistory.currentPoints') }}</div>
          <div class="card-value">
            <span class="value-num">{{ loadingAsset ? '--' : currentPoints }}</span>
            <span class="value-unit">积分</span>
          </div>
        </div>
      </div>

      <!-- 列表头部 -->
      <div class="list-header">
        <span>{{ t('pointHistory.time') }}</span>
        <span>{{ t('pointHistory.type') }}</span>
        <span>{{ t('pointHistory.change') }}</span>
      </div>

      <!-- 记录列表 -->
      <div class="record-list">
        <div
          v-for="(item, index) in recordList"
          :key="item.id || item.ID || index"
          class="record-item"
        >
          <div class="item-time">
            {{ formatDate(item.createTime || item.time || item.createdAt || item.date) }}
          </div>
          <div class="item-type">{{ getType(item) }}</div>
          <div class="item-change" :class="{ positive: isIncrease(item), negative: !isIncrease(item) }">
            {{ isIncrease(item) ? '+' : '' }}{{ getChange(item) }}
          </div>
        </div>

        <!-- 骨架屏 -->
        <div v-for="n in (loading && recordList.length === 0 ? 6 : 0)" :key="'skeleton-' + n" class="record-item skeleton">
          <div class="skeleton-line"></div>
          <div class="skeleton-line mid-line"></div>
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
            <circle cx="12" cy="12" r="10"></circle>
            <path d="M12 6v12M8 10h8M8 14h8"></path>
          </svg>
        </div>
        <p>{{ t('common.empty') }}</p>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use "@/assets/styles/variables.scss" as *;

.record-page {
  min-height: 100%;
}

.record-content {
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

// 概览卡
.overview-card {
  display: flex;
  align-items: center;
  gap: $spacing-md;
  padding: $spacing-lg;
  border-radius: $radius-lg;
  margin-bottom: $spacing-lg;
  box-shadow: $shadow-card;
  border: 1px solid $border-light;

  &.points-card {
    background: linear-gradient(135deg, rgba(255, 184, 0, 0.1) 0%, rgba(255, 138, 0, 0.05) 100%);
    border-color: rgba(255, 184, 0, 0.2);
  }

  .card-icon {
    width: 48px;
    height: 48px;
    border-radius: $radius-md;
    background: linear-gradient(135deg, #FFB800 0%, #FF8A00 100%);
    color: $text-white;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;

    svg {
      width: 24px;
      height: 24px;
    }
  }

  .card-info {
    flex: 1;

    .card-label {
      font-size: 13px;
      color: $text-secondary;
      margin-bottom: 4px;
    }

    .card-value {
      .value-num {
        font-size: 28px;
        font-weight: $font-weight-bold;
        color: #FF9500;
      }

      .value-unit {
        font-size: 13px;
        color: $text-secondary;
        margin-left: 6px;
      }
    }
  }
}

// 列表头部
.list-header {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  font-size: 12px;
  color: $text-tertiary;
  border-bottom: 1px solid $border-light;

  & > span:nth-child(1) {
    flex: 1.5;
  }

  & > span:nth-child(2) {
    flex: 2;
    text-align: center;
  }

  & > span:nth-child(3) {
    flex: 1;
    text-align: right;
  }
}

// 记录列表
.record-list {
  background: $bg-card;
  border-radius: $radius-lg;
  box-shadow: $shadow-card;
  border: 1px solid $border-light;
  overflow: hidden;
}

.record-item {
  display: flex;
  align-items: center;
  padding: 14px 16px;
  border-bottom: 1px solid $border-light;
  transition: background 0.2s ease;

  &:last-child {
    border-bottom: none;
  }

  &:hover {
    background: $bg-hover;
  }

  &.skeleton {
    gap: 12px;
    pointer-events: none;
  }

  .item-time {
    flex: 1.5;
    font-size: 13px;
    color: $text-tertiary;
  }

  .item-type {
    flex: 2;
    text-align: center;
    font-size: 14px;
    color: $text-primary;
  }

  .item-change {
    flex: 1;
    text-align: right;
    font-size: 16px;
    font-weight: $font-weight-bold;

    &.positive {
      color: $success-color;
    }

    &.negative {
      color: $error-color;
    }
  }
}

.skeleton-line {
  height: 14px;
  background: linear-gradient(90deg, $bg-input 25%, $border-light 50%, $bg-input 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;

  &:nth-child(1) { flex: 1.5; }
  &:nth-child(2) { flex: 2; }
  &:nth-child(3) { flex: 1; }

  &.mid-line { width: 40%; }
  &.short-line { width: 30%; }
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
