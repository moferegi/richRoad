<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getCollectionDetailList, uncollect } from '@/api/learning'

const { t } = useI18n()
const router = useRouter()

// Tab 配置：all=全部, word=单词(1), sentence=句子(2/5), video=视频(3), diary=日记(4)
const tabs = [
  { key: 'all', label: 'collections.all', targetType: 0 },
  { key: 'word', label: 'collections.words', targetType: 1 },
  { key: 'sentence', label: 'collections.sentences', targetType: 2 },
  { key: 'video', label: 'collections.videos', targetType: 3 },
  { key: 'diary', label: 'collections.diaries', targetType: 4 }
]

const activeTab = ref('all')
const collectionList = ref([])
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const finished = ref(false)

// 确认取消收藏
const confirmDialog = ref({ visible: false, item: null })

// 获取当前 targetType
const currentTargetType = () => {
  const tab = tabs.find(t => t.key === activeTab.value)
  return tab?.targetType || 0
}

// 加载收藏列表
async function loadCollections(reset = false) {
  if (loading.value) return
  if (reset) {
    page.value = 1
    finished.value = false
    collectionList.value = []
  }
  if (finished.value) return

  loading.value = true
  try {
    const params = {
      page: page.value,
      pageSize: pageSize.value
    }
    const tt = currentTargetType()
    if (tt > 0) {
      params.targetType = tt
    }
    const res = await getCollectionDetailList(params)
    if (res.code === 0) {
      const list = res.data?.list || res.data || []
      collectionList.value = reset ? list : [...collectionList.value, ...list]
      if (list.length < pageSize.value) {
        finished.value = true
      } else {
        page.value++
      }
    }
  } catch (e) {
    console.warn('load collections failed', e)
  } finally {
    loading.value = false
  }
}

// 切换 Tab
function switchTab(key) {
  if (activeTab.value === key) return
  activeTab.value = key
  loadCollections(true)
}

// 获取类型标签文本
function getTypeLabel(item) {
  const type = item.targetType
  const map = {
    1: '单词',
    2: '句子',
    3: '视频',
    4: '日记',
    5: '日记句子'
  }
  return map[type] || '收藏'
}

// 获取类型对应的颜色
function getTypeClass(item) {
  const type = item.targetType
  const map = {
    1: 'type-word',
    2: 'type-sentence',
    3: 'type-video',
    4: 'type-diary',
    5: 'type-diary-sentence'
  }
  return map[type] || ''
}

// 点击条目跳转
function handleItemClick(item) {
  const type = item.targetType
  const id = item.targetId || item.id || item.ID
  if (!id) return

  switch (type) {
    case 1: // 单词
      router.push('/learning/typing')
      break
    case 2: // 句子
      router.push(`/learning/player/${item.episodeId || id}`)
      break
    case 3: // 视频
      router.push(`/learning/video-detail/${id}`)
      break
    case 4: // 日记
      router.push(`/learning/diary-detail/${id}`)
      break
    case 5: // 日记句子
      router.push(`/learning/diary-detail/${item.diaryId || id}`)
      break
  }
}

// 打开取消收藏确认
function openUncollectConfirm(item, e) {
  e.stopPropagation()
  confirmDialog.value = { visible: true, item }
}

// 取消收藏
async function doUncollect() {
  const item = confirmDialog.value.item
  if (!item) return

  try {
    await uncollect(item.targetType, item.targetId || item.id || item.ID)
    // 从列表中移除
    const idx = collectionList.value.findIndex(
      c => c.targetId === item.targetId && c.targetType === item.targetType
    )
    if (idx >= 0) {
      collectionList.value.splice(idx, 1)
    }
  } catch (e) {
    console.warn('uncollect failed', e)
  } finally {
    confirmDialog.value = { visible: false, item: null }
  }
}

// 取消对话框
function cancelUncollect() {
  confirmDialog.value = { visible: false, item: null }
}

// 获取标题
function getTitle(item) {
  return item.title || item.name || item.nameI18n || item.word || item.english || item.content || '收藏内容'
}

// 获取副标题/描述
function getSubtitle(item) {
  return item.description || item.summary || item.translate || item.meaning || item.explanation || ''
}

// 格式化时间
function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return dateStr
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

onMounted(() => {
  loadCollections(true)
})
</script>

<template>
  <div class="collections-page">
    <div class="collections-content">
      <!-- 顶部标题 -->
      <div class="page-header">
        <h1 class="page-title">{{ t('collections.title') }}</h1>
      </div>

      <!-- Tab 切换 -->
      <div class="tab-bar">
        <div class="tab-wrapper">
          <button
            v-for="tab in tabs"
            :key="tab.key"
            class="tab-item"
            :class="{ active: activeTab === tab.key }"
            @click="switchTab(tab.key)"
          >
            {{ t(tab.label) }}
          </button>
        </div>
      </div>

      <!-- 收藏列表 -->
      <div class="collection-list">
        <div
          v-for="item in collectionList"
          :key="(item.targetId || item.id || item.ID) + '-' + item.targetType"
          class="collection-item"
          @click="handleItemClick(item)"
        >
          <!-- 左侧缩略图/图标 -->
          <div class="item-thumb">
            <img
              v-if="item.coverUrl || item.cover || item.imageUrl"
              :src="item.coverUrl || item.cover || item.imageUrl"
              :alt="getTitle(item)"
              class="thumb-img"
            />
            <div v-else class="thumb-placeholder" :class="getTypeClass(item)">
              <svg v-if="item.targetType === 1" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M4 7V4h16v3"></path>
                <path d="M9 20h6"></path>
                <path d="M12 4v16"></path>
              </svg>
              <svg v-else-if="item.targetType === 3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polygon points="23 7 16 12 23 17 23 7"></polygon>
                <rect x="1" y="5" width="15" height="14" rx="2" ry="2"></rect>
              </svg>
              <svg v-else-if="item.targetType === 4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                <polyline points="14 2 14 8 20 8"></polyline>
              </svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
              </svg>
            </div>
          </div>

          <!-- 右侧信息 -->
          <div class="item-info">
            <div class="item-top">
              <span class="type-tag" :class="getTypeClass(item)">{{ getTypeLabel(item) }}</span>
              <span class="item-date">{{ formatDate(item.createdAt || item.createTime || item.date) }}</span>
            </div>
            <h3 class="item-title">{{ getTitle(item) }}</h3>
            <p v-if="getSubtitle(item)" class="item-subtitle">{{ getSubtitle(item) }}</p>
            <button class="remove-btn" @click="openUncollectConfirm(item, $event)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3 6 5 6 21 6"></polyline>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
              </svg>
              {{ t('collections.remove') }}
            </button>
          </div>
        </div>

        <!-- 骨架屏 -->
        <div v-for="n in (loading && collectionList.length === 0 ? 5 : 0)" :key="'skeleton-' + n" class="collection-item skeleton">
          <div class="item-thumb skeleton-block"></div>
          <div class="item-info">
            <div class="skeleton-line title-line"></div>
            <div class="skeleton-line"></div>
            <div class="skeleton-line short-line"></div>
          </div>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading && collectionList.length > 0" class="load-more">
        <div class="loading-spinner"></div>
        <span>{{ t('common.loading') }}</span>
      </div>

      <!-- 加载更多 -->
      <div v-else-if="!finished && collectionList.length > 0" class="load-more">
        <button class="btn btn-outline" @click="loadCollections()">
          加载更多
        </button>
      </div>

      <!-- 空状态 -->
      <div v-else-if="!loading && collectionList.length === 0" class="empty-state">
        <div class="empty-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
          </svg>
        </div>
        <p>{{ t('common.empty') }}</p>
      </div>
    </div>

    <!-- 确认弹窗 -->
    <div v-if="confirmDialog.visible" class="modal-mask" @click="cancelUncollect">
      <div class="modal-content confirm-modal" @click.stop>
        <div class="confirm-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="3 6 5 6 21 6"></polyline>
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
          </svg>
        </div>
        <h3 class="confirm-title">确认取消收藏？</h3>
        <p class="confirm-desc">取消收藏后，可在详情页重新收藏</p>
        <div class="confirm-actions">
          <button class="btn btn-outline" @click="cancelUncollect">{{ t('common.cancel') }}</button>
          <button class="btn btn-danger" @click="doUncollect">{{ t('common.delete') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use "@/assets/styles/variables.scss" as *;

.collections-page {
  min-height: 100%;
}

.collections-content {
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

// Tab 栏
.tab-bar {
  margin-bottom: $spacing-lg;
}

.tab-wrapper {
  display: inline-flex;
  background: $bg-card;
  border-radius: $radius-pill;
  padding: 4px;
  box-shadow: $shadow-card;
  border: 1px solid $border-light;
}

.tab-item {
  padding: 8px 20px;
  font-size: 14px;
  color: $text-secondary;
  border-radius: $radius-pill;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
  background: transparent;
  font-weight: $font-weight-medium;

  &:hover {
    color: $primary-color;
  }

  &.active {
    background: $gradient-primary;
    color: $text-white;
    box-shadow: 0 2px 8px rgba(109, 91, 255, 0.3);
  }
}

// 收藏列表
.collection-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.collection-item {
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
    transform: translateX(4px);
    box-shadow: $shadow-md;

    .item-title {
      color: $primary-color;
    }
  }

  &.skeleton {
    cursor: default;
    pointer-events: none;
  }
}

.item-thumb {
  width: 100px;
  height: 100px;
  border-radius: $radius-md;
  overflow: hidden;
  flex-shrink: 0;
  background: $bg-input;

  .thumb-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .thumb-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: $text-white;

    svg {
      width: 32px;
      height: 32px;
    }

    &.type-word { background: linear-gradient(135deg, #34C759 0%, #30D158 100%); }
    &.type-sentence { background: linear-gradient(135deg, #007AFF 0%, #5AC8FA 100%); }
    &.type-video { background: linear-gradient(135deg, #AF52DE 0%, #DA8FFF 100%); }
    &.type-diary { background: linear-gradient(135deg, #FF9500 0%, #FFB800 100%); }
    &.type-diary-sentence { background: linear-gradient(135deg, #FF6B6B 0%, #FF9F43 100%); }
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

.item-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;

  .type-tag {
    display: inline-block;
    padding: 2px 10px;
    font-size: 11px;
    border-radius: $radius-pill;
    font-weight: $font-weight-medium;

    &.type-word { background: rgba(52, 199, 89, 0.12); color: #34C759; }
    &.type-sentence { background: rgba(0, 122, 255, 0.12); color: #007AFF; }
    &.type-video { background: rgba(175, 82, 222, 0.12); color: #AF52DE; }
    &.type-diary { background: rgba(255, 149, 0, 0.12); color: #FF9500; }
    &.type-diary-sentence { background: rgba(255, 107, 107, 0.12); color: #FF6B6B; }
  }

  .item-date {
    font-size: 12px;
    color: $text-tertiary;
  }
}

.item-title {
  font-size: 15px;
  font-weight: $font-weight-semibold;
  color: $text-primary;
  margin: 0 0 4px 0;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.2s ease;
}

.item-subtitle {
  font-size: 13px;
  color: $text-secondary;
  line-height: 1.5;
  margin: 0 0 10px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  flex: 1;
}

.remove-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: transparent;
  border: 1px solid $border-color;
  border-radius: $radius-pill;
  font-size: 12px;
  color: $text-tertiary;
  cursor: pointer;
  transition: all 0.2s ease;
  align-self: flex-start;

  svg {
    width: 12px;
    height: 12px;
  }

  &:hover {
    border-color: $error-color;
    color: $error-color;
    background: rgba(255, 59, 48, 0.06);
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
    width: 60%;
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
