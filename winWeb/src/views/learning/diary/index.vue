<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getDiaryCategoryList, getDiaryTagPublic, getDiaryListByTag } from '@/api/learning'
import { localText } from '@/utils/i18n'

const { t } = useI18n()
const router = useRouter()

// 分类
const categories = ref([])
const selectedCategoryId = ref(0)
const showCategoryPopup = ref(false)

// 标签
const tags = ref([])
const selectedTagIds = ref([])
const showTagPopup = ref(false)

// 日记列表
const diaryList = ref([])
const page = ref(1)
const pageSize = ref(12)
const loading = ref(false)
const finished = ref(false)

// 加载分类
async function loadCategories() {
  try {
    const res = await getDiaryCategoryList()
    if (res.code === 0) {
      categories.value = res.data?.list || res.data || []
    }
  } catch (e) {
    console.warn('load diary categories failed', e)
  }
}

// 加载标签
async function loadTags() {
  try {
    const res = await getDiaryTagPublic()
    if (res.code === 0) {
      tags.value = res.data?.list || res.data || []
    }
  } catch (e) {
    console.warn('load diary tags failed', e)
  }
}

// 加载日记列表
async function loadDiaries(reset = false) {
  if (loading.value) return
  if (reset) {
    page.value = 1
    finished.value = false
    diaryList.value = []
  }
  if (finished.value) return

  loading.value = true
  try {
    const params = {
      page: page.value,
      pageSize: pageSize.value
    }
    if (selectedCategoryId.value > 0) {
      params.categoryId = selectedCategoryId.value
    }
    if (selectedTagIds.value.length > 0) {
      params.tagIds = selectedTagIds.value.join(',')
    }
    const res = await getDiaryListByTag(params)
    if (res.code === 0) {
      const list = res.data?.list || res.data || []
      diaryList.value = reset ? list : [...diaryList.value, ...list]
      if (list.length < pageSize.value) {
        finished.value = true
      } else {
        page.value++
      }
    }
  } catch (e) {
    console.warn('load diaries failed', e)
  } finally {
    loading.value = false
  }
}

// 选择分类
function selectCategory(id) {
  selectedCategoryId.value = id
  showCategoryPopup.value = false
  loadDiaries(true)
}

// 切换标签选择
function toggleTag(tagId) {
  const idx = selectedTagIds.value.indexOf(tagId)
  if (idx >= 0) {
    selectedTagIds.value.splice(idx, 1)
  } else {
    selectedTagIds.value.push(tagId)
  }
}

// 确认标签
function confirmTags() {
  showTagPopup.value = false
  loadDiaries(true)
}

// 重置筛选
function resetFilters() {
  selectedCategoryId.value = 0
  selectedTagIds.value = []
  loadDiaries(true)
}

// 跳转详情
function goDetail(id) {
  router.push(`/learning/diary-detail/${id}`)
}

// 获取选中分类名称
const selectedCategoryName = () => {
  if (selectedCategoryId.value === 0) return t('common.all')
  const cat = categories.value.find(c => (c.ID || c.id) === selectedCategoryId.value)
  return cat?.name ? localText(cat.name) : t('common.all')
}

// 获取标签名称
const getTagName = (tagId) => {
  const tag = tags.value.find(t => (t.ID || t.id) === tagId)
  return tag?.name || tag?.tagName || ''
}

// 移除已选标签
function removeTag(tagId) {
  const idx = selectedTagIds.value.indexOf(tagId)
  if (idx >= 0) {
    selectedTagIds.value.splice(idx, 1)
    loadDiaries(true)
  }
}

// 截取文本摘要
function truncateText(text, maxLen = 80) {
  if (!text) return ''
  if (text.length <= maxLen) return text
  return text.substring(0, maxLen) + '...'
}

onMounted(() => {
  loadCategories()
  loadTags()
  loadDiaries(true)
})
</script>

<template>
  <div class="diary-page">
    <div class="diary-content">
      <!-- 顶部标题栏 -->
      <div class="page-header">
        <h1 class="page-title">{{ t('diary.title') }}</h1>
      </div>

      <!-- 筛选栏 -->
      <div class="filter-bar">
        <div class="filter-left">
          <button class="filter-btn" @click="showCategoryPopup = true">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 3H2l8 9.46V19l4 2v-8.54L22 3z"></path>
            </svg>
            <span>{{ t('diary.category') }}: {{ selectedCategoryName() }}</span>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="arrow">
              <polyline points="6 9 12 15 18 9"></polyline>
            </svg>
          </button>

          <button class="filter-btn" @click="showTagPopup = true">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"></path>
              <line x1="7" y1="7" x2="7.01" y2="7"></line>
            </svg>
            <span>{{ t('diary.tags') }}</span>
            <span v-if="selectedTagIds.length > 0" class="badge">{{ selectedTagIds.length }}</span>
          </button>

          <button v-if="selectedCategoryId > 0 || selectedTagIds.length > 0" class="reset-btn" @click="resetFilters">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="1 4 1 10 7 10"></polyline>
              <path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"></path>
            </svg>
            {{ t('diary.reset') }}
          </button>
        </div>
      </div>

      <!-- 已选标签 chip 展示 -->
      <div v-if="selectedTagIds.length > 0" class="selected-tags">
        <span
          v-for="tagId in selectedTagIds"
          :key="tagId"
          class="tag-chip"
        >
          {{ getTagName(tagId) }}
          <span class="chip-close" @click="removeTag(tagId)">×</span>
        </span>
      </div>

      <!-- 日记网格 -->
      <div class="diary-grid">
        <div
          v-for="diary in diaryList"
          :key="diary.ID || diary.id"
          class="diary-card"
          @click="goDetail(diary.ID || diary.id)"
        >
          <div class="card-cover">
            <img
              v-if="diary.coverUrl || diary.cover || diary.imageUrl"
              :src="diary.coverUrl || diary.cover || diary.imageUrl"
              :alt="diary.title || diary.name"
              class="cover-img"
            />
            <div v-else class="cover-placeholder">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
                <circle cx="8.5" cy="8.5" r="1.5"></circle>
                <polyline points="21 15 16 10 5 21"></polyline>
              </svg>
            </div>
          </div>
          <div class="card-content">
            <h3 class="card-title">{{ localText(diary.name) || diary.title }}</h3>
            <p class="card-summary">
              {{ truncateText(diary.content || diary.text || diary.summary || '') }}
            </p>
            <button class="view-btn">
              {{ t('diary.view') }}
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="9 18 15 12 9 6"></polyline>
              </svg>
            </button>
          </div>
        </div>

        <!-- 骨架屏 -->
        <div v-for="n in (loading && diaryList.length === 0 ? 6 : 0)" :key="'skeleton-' + n" class="diary-card skeleton">
          <div class="card-cover skeleton-block"></div>
          <div class="card-content">
            <div class="skeleton-line title-line"></div>
            <div class="skeleton-line"></div>
            <div class="skeleton-line short-line"></div>
          </div>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading && diaryList.length > 0" class="load-more">
        <div class="loading-spinner"></div>
        <span>{{ t('common.loading') }}</span>
      </div>

      <!-- 加载更多 -->
      <div v-else-if="!finished && diaryList.length > 0" class="load-more">
        <button class="btn btn-outline" @click="loadDiaries()">
          加载更多
        </button>
      </div>

      <!-- 空状态 -->
      <div v-else-if="!loading && diaryList.length === 0" class="empty-state">
        <div class="empty-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
            <polyline points="14 2 14 8 20 8"></polyline>
            <line x1="16" y1="13" x2="8" y2="13"></line>
            <line x1="16" y1="17" x2="8" y2="17"></line>
          </svg>
        </div>
        <p>{{ t('common.empty') }}</p>
      </div>
    </div>

    <!-- 分类选择弹窗 -->
    <div v-if="showCategoryPopup" class="popup-mask" @click="showCategoryPopup = false">
      <div class="popup-content category-popup" @click.stop>
        <div class="popup-header">
          <h3>{{ t('diary.category') }}</h3>
          <button class="close-btn" @click="showCategoryPopup = false">×</button>
        </div>
        <div class="popup-body">
          <div
            class="category-option"
            :class="{ active: selectedCategoryId === 0 }"
            @click="selectCategory(0)"
          >
            <span>{{ t('common.all') }}</span>
            <svg v-if="selectedCategoryId === 0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
          </div>
          <div
            v-for="cat in categories"
            :key="cat.ID || cat.id"
            class="category-option"
            :class="{ active: selectedCategoryId === (cat.ID || cat.id) }"
            @click="selectCategory(cat.ID || cat.id)"
          >
            <span>{{ localText(cat.name) }}</span>
            <svg v-if="selectedCategoryId === (cat.ID || cat.id)" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- 标签选择弹窗 -->
    <div v-if="showTagPopup" class="popup-mask" @click="showTagPopup = false">
      <div class="popup-content tag-popup" @click.stop>
        <div class="popup-header">
          <h3>{{ t('diary.tags') }}</h3>
          <button class="close-btn" @click="showTagPopup = false">×</button>
        </div>
        <div class="popup-body">
          <div class="tag-grid">
            <div
              v-for="tag in tags"
              :key="tag.ID || tag.id"
              class="tag-option"
              :class="{ active: selectedTagIds.includes(tag.ID || tag.id) }"
              @click="toggleTag(tag.ID || tag.id)"
            >
              {{ tag.name || tag.tagName }}
            </div>
          </div>
        </div>
        <div class="popup-footer">
          <button class="btn btn-outline" @click="showTagPopup = false">{{ t('common.cancel') }}</button>
          <button class="btn btn-primary" @click="confirmTags">{{ t('common.confirm') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use "@/assets/styles/variables.scss" as *;

.diary-page {
  min-height: 100%;
}

.diary-content {
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
  }
}

// 筛选栏
.filter-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: $spacing-md;

  .filter-left {
    display: flex;
    align-items: center;
    gap: 12px;
  }
}

.filter-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 18px;
  background: $bg-card;
  border: 1px solid $border-color;
  border-radius: $radius-pill;
  font-size: 14px;
  color: $text-secondary;
  cursor: pointer;
  transition: all 0.2s ease;

  svg {
    width: 16px;
    height: 16px;
  }

  .arrow {
    width: 14px;
    height: 14px;
  }

  &:hover {
    border-color: $primary-light;
    color: $primary-color;
  }

  .badge {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 18px;
    height: 18px;
    padding: 0 5px;
    background: $gradient-primary;
    color: $text-white;
    font-size: 11px;
    border-radius: 9px;
    font-weight: $font-weight-medium;
  }
}

.reset-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 10px 16px;
  background: transparent;
  border: 1px dashed $border-color;
  border-radius: $radius-pill;
  font-size: 13px;
  color: $text-tertiary;
  cursor: pointer;
  transition: all 0.2s ease;

  svg {
    width: 14px;
    height: 14px;
  }

  &:hover {
    border-color: $primary-light;
    color: $primary-color;
  }
}

// 已选标签
.selected-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: $spacing-lg;

  .tag-chip {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 6px 12px;
    background: $bg-active;
    color: $primary-color;
    border-radius: $radius-pill;
    font-size: 13px;
    font-weight: $font-weight-medium;

    .chip-close {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 16px;
      height: 16px;
      font-size: 16px;
      line-height: 1;
      border-radius: 50%;
      cursor: pointer;
      transition: all 0.15s ease;

      &:hover {
        background: rgba(109, 91, 255, 0.2);
      }
    }
  }
}

// 日记网格
.diary-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.diary-card {
  display: flex;
  background: $bg-card;
  border-radius: $radius-lg;
  overflow: hidden;
  box-shadow: $shadow-card;
  border: 1px solid $border-light;
  cursor: pointer;
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-4px);
    box-shadow: $shadow-lg;

    .view-btn {
      color: $primary-color;
    }
  }

  &.skeleton {
    cursor: default;
    pointer-events: none;
  }
}

.card-cover {
  width: 140px;
  flex-shrink: 0;
  background: $bg-input;
  overflow: hidden;
  position: relative;

  .cover-img {
    width: 100%;
    height: 100%;
    min-height: 180px;
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

  &.skeleton-block {
    background: linear-gradient(90deg, $bg-input 25%, $border-light 50%, $bg-input 75%);
    background-size: 200% 100%;
    animation: shimmer 1.5s infinite;
  }
}

.card-content {
  flex: 1;
  padding: 16px;
  display: flex;
  flex-direction: column;
  min-width: 0;

  .card-title {
    font-size: 15px;
    font-weight: $font-weight-semibold;
    color: $text-primary;
    margin-bottom: 8px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    line-height: 1.4;
  }

  .card-summary {
    flex: 1;
    font-size: 13px;
    color: $text-secondary;
    line-height: 1.6;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
    margin-bottom: 12px;
  }

  .view-btn {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    font-size: 13px;
    color: $text-tertiary;
    font-weight: $font-weight-medium;
    transition: color 0.2s ease;
    align-self: flex-start;

    svg {
      width: 14px;
      height: 14px;
    }
  }

  .skeleton-line {
    height: 14px;
    background: linear-gradient(90deg, $bg-input 25%, $border-light 50%, $bg-input 75%);
    background-size: 200% 100%;
    animation: shimmer 1.5s infinite;
    border-radius: 4px;
    margin-bottom: 10px;

    &.title-line {
      height: 18px;
      width: 80%;
    }

    &.short-line {
      width: 50%;
      margin-bottom: 0;
    }
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

// 弹窗
.popup-mask {
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

.popup-content {
  background: $bg-card;
  border-radius: $radius-lg;
  min-width: 320px;
  max-width: 90vw;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: $shadow-xl;
  animation: slideUp 0.25s ease;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.popup-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid $border-light;

  h3 {
    font-size: 16px;
    font-weight: $font-weight-semibold;
    color: $text-primary;
  }

  .close-btn {
    width: 28px;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 22px;
    color: $text-tertiary;
    border-radius: 50%;
    cursor: pointer;
    transition: all 0.15s ease;

    &:hover {
      background: $bg-hover;
      color: $text-primary;
    }
  }
}

.popup-body {
  flex: 1;
  overflow-y: auto;
  padding: 12px 0;
}

.popup-footer {
  display: flex;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid $border-light;

  .btn {
    flex: 1;
  }
}

// 分类弹窗
.category-popup {
  width: 320px;

  .category-option {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 20px;
    font-size: 14px;
    color: $text-secondary;
    cursor: pointer;
    transition: all 0.15s ease;

    svg {
      width: 18px;
      height: 18px;
    }

    &:hover {
      background: $bg-hover;
      color: $text-primary;
    }

    &.active {
      color: $primary-color;
      font-weight: $font-weight-medium;
      background: $bg-active;
    }
  }
}

// 标签弹窗
.tag-popup {
  width: 420px;

  .tag-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    padding: 12px 20px;
  }

  .tag-option {
    padding: 8px 16px;
    background: $bg-input;
    border: 1px solid $border-color;
    border-radius: $radius-pill;
    font-size: 13px;
    color: $text-secondary;
    cursor: pointer;
    transition: all 0.2s ease;

    &:hover {
      border-color: $primary-light;
      color: $primary-color;
    }

    &.active {
      background: $gradient-primary;
      border-color: transparent;
      color: $text-white;
      font-weight: $font-weight-medium;
    }
  }
}

// 响应式
@media (max-width: 1100px) {
  .diary-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 700px) {
  .diary-grid {
    grid-template-columns: 1fr;
  }
}
</style>
