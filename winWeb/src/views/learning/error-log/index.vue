<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getWordErrorLogList, deleteWordErrorLog } from '@/api/learning'

const { t } = useI18n()
const router = useRouter()

const errorList = ref([])
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const finished = ref(false)

// 确认删除
const confirmDialog = ref({ visible: false, item: null })

// 错误总数
const totalErrors = computed(() => {
  let total = 0
  errorList.value.forEach(item => {
    total += item.errorCount || item.count || 0
  })
  return total
})

// 加载错误列表
async function loadErrors(reset = false) {
  if (loading.value) return
  if (reset) {
    page.value = 1
    finished.value = false
    errorList.value = []
  }
  if (finished.value) return

  loading.value = true
  try {
    const res = await getWordErrorLogList({
      page: page.value,
      pageSize: pageSize.value
    })
    if (res.code === 0) {
      const list = res.data?.list || res.data || []
      errorList.value = reset ? list : [...errorList.value, ...list]
      if (list.length < pageSize.value) {
        finished.value = true
      } else {
        page.value++
      }
    }
  } catch (e) {
    console.warn('load error log failed', e)
  } finally {
    loading.value = false
  }
}

// 获取单词
function getWord(item) {
  return item.word || item.wordText || item.english || item.name || ''
}

// 获取释义
function getMeaning(item) {
  return item.meaning || item.explanation || item.translate || item.zh || ''
}

// 获取错误次数
function getErrorCount(item) {
  return item.errorCount || item.count || item.times || 0
}

// 获取期望字符
function getExpected(item) {
  return item.expectedChar || item.expected || item.correctChar || ''
}

// 获取输入字符
function getInput(item) {
  return item.inputChar || item.input || item.wrongChar || ''
}

// 获取错误位置
function getPosition(item) {
  return item.position ?? item.errorIndex ?? ''
}

// 格式化时间
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

// 打开删除确认
function openDeleteConfirm(item, e) {
  e.stopPropagation()
  confirmDialog.value = { visible: true, item }
}

// 删除错误记录
async function doDelete() {
  const item = confirmDialog.value.item
  if (!item) return

  const wordId = item.wordId || item.id || item.ID
  try {
    await deleteWordErrorLog(wordId)
    // 从列表移除
    const idx = errorList.value.findIndex(e => (e.wordId || e.id || e.ID) === wordId)
    if (idx >= 0) {
      errorList.value.splice(idx, 1)
    }
  } catch (e) {
    console.warn('delete error log failed', e)
  } finally {
    confirmDialog.value = { visible: false, item: null }
  }
}

function cancelDelete() {
  confirmDialog.value = { visible: false, item: null }
}

// 去练习
function goPractice() {
  router.push('/learning/typing')
}

onMounted(() => {
  loadErrors(true)
})
</script>

<template>
  <div class="error-log-page">
    <div class="error-log-content">
      <!-- 顶部标题 -->
      <div class="page-header">
        <h1 class="page-title">{{ t('errorLog.title') }}</h1>
      </div>

      <!-- 摘要卡 -->
      <div class="summary-card">
        <div class="summary-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
            <line x1="12" y1="9" x2="12" y2="13"></line>
            <line x1="12" y1="17" x2="12.01" y2="17"></line>
          </svg>
        </div>
        <div class="summary-info">
          <div class="summary-label">{{ t('errorLog.totalErrors') }}</div>
          <div class="summary-value">{{ errorList.length }} <span class="unit">个单词</span></div>
        </div>
        <button class="btn btn-primary practice-btn" @click="goPractice">
          <svg viewBox="0 0 24 24" fill="currentColor">
            <polygon points="5 3 19 12 5 21 5 3"></polygon>
          </svg>
          {{ t('errorLog.practice') }}
        </button>
      </div>

      <!-- 错误记录列表 -->
      <div class="error-list">
        <div
          v-for="item in errorList"
          :key="item.id || item.ID || item.wordId"
          class="error-card"
        >
          <div class="card-left">
            <h3 class="word-text">{{ getWord(item) }}</h3>
            <p class="word-meaning">{{ getMeaning(item) }}</p>
            <div class="error-meta">
              <span v-if="getExpected(item)" class="meta-tag">
                {{ t('errorLog.expected') }}: <strong>{{ getExpected(item) }}</strong>
              </span>
              <span v-if="getInput(item)" class="meta-tag wrong">
                {{ t('errorLog.input') }}: <strong>{{ getInput(item) }}</strong>
              </span>
              <span v-if="getPosition(item) !== ''" class="meta-tag">
                {{ t('errorLog.position') }}: <strong>第 {{ Number(getPosition(item)) + 1 }} 位</strong>
              </span>
            </div>
          </div>
          <div class="card-right">
            <div class="error-count">
              <span class="count-num">{{ getErrorCount(item) }}</span>
              <span class="count-label">{{ t('errorLog.errorTimes') }}</span>
            </div>
            <div class="update-time">
              {{ t('errorLog.updatedAt') }}: {{ formatDate(item.updatedAt || item.updateTime || item.createTime) }}
            </div>
            <button class="delete-btn" @click="openDeleteConfirm(item, $event)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3 6 5 6 21 6"></polyline>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
              </svg>
              {{ t('errorLog.delete') }}
            </button>
          </div>
        </div>

        <!-- 骨架屏 -->
        <div v-for="n in (loading && errorList.length === 0 ? 4 : 0)" :key="'skeleton-' + n" class="error-card skeleton">
          <div class="card-left">
            <div class="skeleton-line title-line"></div>
            <div class="skeleton-line"></div>
            <div class="skeleton-line short-line"></div>
          </div>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading && errorList.length > 0" class="load-more">
        <div class="loading-spinner"></div>
        <span>{{ t('common.loading') }}</span>
      </div>

      <!-- 加载更多 -->
      <div v-else-if="!finished && errorList.length > 0" class="load-more">
        <button class="btn btn-outline" @click="loadErrors()">
          加载更多
        </button>
      </div>

      <!-- 空状态 -->
      <div v-else-if="!loading && errorList.length === 0" class="empty-state">
        <div class="empty-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
            <polyline points="22 4 12 14.01 9 11.01"></polyline>
          </svg>
        </div>
        <p>太棒了！暂无错误记录</p>
        <button class="btn btn-primary empty-btn" @click="goPractice">{{ t('errorLog.practice') }}</button>
      </div>
    </div>

    <!-- 确认弹窗 -->
    <div v-if="confirmDialog.visible" class="modal-mask" @click="cancelDelete">
      <div class="modal-content confirm-modal" @click.stop>
        <div class="confirm-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="3 6 5 6 21 6"></polyline>
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
          </svg>
        </div>
        <h3 class="confirm-title">确认删除？</h3>
        <p class="confirm-desc">删除后错误记录将无法恢复</p>
        <div class="confirm-actions">
          <button class="btn btn-outline" @click="cancelDelete">{{ t('common.cancel') }}</button>
          <button class="btn btn-danger" @click="doDelete">{{ t('common.delete') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use "@/assets/styles/variables.scss" as *;

.error-log-page {
  min-height: 100%;
}

.error-log-content {
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

// 摘要卡
.summary-card {
  display: flex;
  align-items: center;
  gap: $spacing-md;
  padding: $spacing-lg;
  background: linear-gradient(135deg, rgba(255, 59, 48, 0.1) 0%, rgba(255, 107, 107, 0.05) 100%);
  border-radius: $radius-lg;
  margin-bottom: $spacing-lg;
  border: 1px solid rgba(255, 59, 48, 0.15);
}

.summary-icon {
  width: 48px;
  height: 48px;
  border-radius: $radius-md;
  background: linear-gradient(135deg, #FF3B30 0%, #FF6B6B 100%);
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

.summary-info {
  flex: 1;

  .summary-label {
    font-size: 13px;
    color: $text-secondary;
    margin-bottom: 4px;
  }

  .summary-value {
    font-size: 24px;
    font-weight: $font-weight-bold;
    color: $error-color;

    .unit {
      font-size: 13px;
      font-weight: $font-weight-regular;
      color: $text-secondary;
      margin-left: 4px;
    }
  }
}

.practice-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;

  svg {
    width: 14px;
    height: 14px;
  }
}

// 错误列表
.error-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.error-card {
  display: flex;
  justify-content: space-between;
  gap: $spacing-md;
  padding: 18px 20px;
  background: $bg-card;
  border-radius: $radius-lg;
  box-shadow: $shadow-card;
  border: 1px solid $border-light;
  transition: all 0.2s ease;

  &:hover {
    box-shadow: $shadow-md;
    border-color: rgba(255, 59, 48, 0.2);
  }

  &.skeleton {
    pointer-events: none;
  }
}

.card-left {
  flex: 1;
  min-width: 0;
}

.word-text {
  font-size: 18px;
  font-weight: $font-weight-bold;
  color: $text-primary;
  margin: 0 0 6px 0;
}

.word-meaning {
  font-size: 13px;
  color: $text-secondary;
  margin: 0 0 10px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.5;
}

.error-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;

  .meta-tag {
    font-size: 12px;
    color: $text-secondary;
    background: $bg-input;
    padding: 3px 8px;
    border-radius: $radius-sm;

    strong {
      color: $success-color;
      font-weight: $font-weight-semibold;
    }

    &.wrong strong {
      color: $error-color;
    }
  }
}

.card-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8px;
  flex-shrink: 0;
  text-align: right;
}

.error-count {
  .count-num {
    font-size: 24px;
    font-weight: $font-weight-bold;
    color: $error-color;
    line-height: 1;
  }

  .count-label {
    font-size: 11px;
    color: $text-tertiary;
    display: block;
    margin-top: 4px;
  }
}

.update-time {
  font-size: 11px;
  color: $text-tertiary;
}

.delete-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 5px 12px;
  background: transparent;
  border: 1px solid rgba(255, 59, 48, 0.3);
  border-radius: $radius-pill;
  font-size: 12px;
  color: $error-color;
  cursor: pointer;
  transition: all 0.2s ease;

  svg {
    width: 12px;
    height: 12px;
  }

  &:hover {
    background: rgba(255, 59, 48, 0.1);
    border-color: $error-color;
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
    height: 20px;
    width: 40%;
  }

  &.short-line {
    width: 30%;
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

  &.btn-primary {
    background: $gradient-primary;
    color: $text-white;
    box-shadow: 0 4px 12px rgba(109, 91, 255, 0.3);

    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 6px 16px rgba(109, 91, 255, 0.4);
    }
  }

  &.btn-outline {
    background: $bg-card;
    border: 1px solid $border-color;
    color: $text-secondary;

    &:hover {
      border-color: $primary-light;
      color: $primary-color;
    }
  }

  &.btn-danger {
    background: $error-color;
    color: $text-white;

    &:hover {
      background: darken($error-color, 5%);
    }
  }
}

// 空状态
.empty-state {
  text-align: center;
  padding: 60px 0;
  color: $text-tertiary;

  .empty-icon {
    width: 72px;
    height: 72px;
    margin: 0 auto 16px;
    color: $success-color;

    svg {
      width: 100%;
      height: 100%;
    }
  }

  p {
    font-size: 14px;
    margin: 0 0 20px 0;
  }

  .empty-btn {
    display: inline-flex;
  }
}

// 确认弹窗
.modal-mask {
  position: fixed;
  inset: 0;
  background: $bg-mask;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: $z-modal;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal-content {
  background: $bg-card;
  border-radius: $radius-lg;
  box-shadow: $shadow-xl;
  animation: scaleIn 0.2s ease;
}

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.confirm-modal {
  width: 320px;
  padding: 32px 24px 24px;
  text-align: center;

  .confirm-icon {
    width: 56px;
    height: 56px;
    margin: 0 auto 16px;
    border-radius: 50%;
    background: rgba(255, 59, 48, 0.1);
    color: $error-color;
    display: flex;
    align-items: center;
    justify-content: center;

    svg {
      width: 28px;
      height: 28px;
    }
  }

  .confirm-title {
    font-size: 18px;
    font-weight: $font-weight-bold;
    color: $text-primary;
    margin: 0 0 8px 0;
  }

  .confirm-desc {
    font-size: 13px;
    color: $text-secondary;
    margin: 0 0 24px 0;
  }

  .confirm-actions {
    display: flex;
    gap: 12px;

    .btn {
      flex: 1;
    }
  }
}
</style>
