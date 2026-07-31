<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getVideoTagList } from '@/api/learning'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const tagList = ref([])
const loading = ref(false)
const selectedTags = ref([])

const tagGroups = computed(() => {
  const groups = {}
  tagList.value.forEach(tag => {
    const groupKey = tag.groupId || tag.group || tag.categoryId || 'default'
    const groupName = tag.groupName || tag.group || tag.categoryName || t('tagFilter.title')
    if (!groups[groupKey]) {
      groups[groupKey] = {
        key: groupKey,
        name: groupName,
        tags: []
      }
    }
    groups[groupKey].tags.push(tag)
  })
  return Object.values(groups)
})

async function loadTags() {
  loading.value = true
  try {
    const res = await getVideoTagList()
    if (res.code === 0) {
      tagList.value = res.data?.list || res.data || []
    }
    const initialTags = route.query.tags
    if (initialTags) {
      const ids = String(initialTags).split(',').map(Number).filter(Boolean)
      selectedTags.value = ids
    }
  } catch (e) {
    console.warn('load tags failed', e)
  } finally {
    loading.value = false
  }
}

function isSelected(tag) {
  const id = tag.ID || tag.id
  return selectedTags.value.includes(id)
}

function toggleTag(tag) {
  const id = tag.ID || tag.id
  const idx = selectedTags.value.indexOf(id)
  if (idx >= 0) {
    selectedTags.value.splice(idx, 1)
  } else {
    selectedTags.value.push(id)
  }
}

function removeTag(tagId) {
  const idx = selectedTags.value.indexOf(tagId)
  if (idx >= 0) {
    selectedTags.value.splice(idx, 1)
  }
}

function clearAll() {
  selectedTags.value = []
}

function getTagName(tag) {
  return tag.nameI18n || tag.name || tag.tagName || tag.title || '标签'
}

function getSelectedTagInfo(tagId) {
  return tagList.value.find(t => (t.ID || t.id) === tagId) || null
}

function handleConfirm() {
  const tagParam = selectedTags.value.join(',')
  router.push({
    path: '/',
    query: { tags: tagParam || undefined }
  })
}

onMounted(() => {
  loadTags()
})
</script>

<template>
  <div class="tag-filter-page">
    <div class="tag-filter-content">
      <div class="page-header">
        <h1 class="page-title">{{ t('tagFilter.title') }}</h1>
        <p class="page-subtitle">选择感兴趣的标签，精准匹配学习内容</p>
      </div>

      <div v-if="selectedTags.length > 0" class="selected-section">
        <div class="selected-header">
          <span class="selected-label">已选 {{ selectedTags.length }} 个</span>
          <button class="clear-btn" @click="clearAll">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="3 6 5 6 21 6"></polyline>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
            </svg>
            全部清除
          </button>
        </div>
        <div class="selected-chips">
          <span
            v-for="tagId in selectedTags"
            :key="tagId"
            class="selected-chip"
          >
            {{ getTagName(getSelectedTagInfo(tagId) || {}) }}
            <button class="chip-remove" @click="removeTag(tagId)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </span>
        </div>
      </div>

      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <span>{{ t('common.loading') }}</span>
      </div>

      <div v-else-if="tagList.length === 0" class="empty-state">
        <div class="empty-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"></path>
            <line x1="7" y1="7" x2="7.01" y2="7"></line>
          </svg>
        </div>
        <p>{{ t('common.empty') }}</p>
      </div>

      <div v-else class="tag-groups">
        <div
          v-for="group in tagGroups"
          :key="group.key"
          class="tag-group"
        >
          <h3 class="group-title">{{ group.name }}</h3>
          <div class="tag-cloud">
            <button
              v-for="tag in group.tags"
              :key="tag.ID || tag.id"
              class="tag-chip"
              :class="{ active: isSelected(tag) }"
              @click="toggleTag(tag)"
            >
              {{ getTagName(tag) }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="footer-action">
      <button class="btn btn-confirm" @click="handleConfirm">
        {{ t('tagFilter.confirm') }}
      </button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use "@/assets/styles/variables.scss" as *;

.tag-filter-page {
  min-height: 100%;
  padding-bottom: 100px;
}

.tag-filter-content {
  max-width: $content-max-width;
  margin: 0 auto;
  padding: $spacing-xl;
}

.page-header {
  margin-bottom: $spacing-xl;

  .page-title {
    font-size: $font-size-title;
    font-weight: $font-weight-bold;
    color: $text-primary;
    margin: 0 0 6px 0;
  }

  .page-subtitle {
    font-size: $font-size-base;
    color: $text-secondary;
    margin: 0;
  }
}

.selected-section {
  background: $bg-card;
  border-radius: $radius-lg;
  padding: $spacing-md $spacing-lg;
  margin-bottom: $spacing-xl;
  box-shadow: $shadow-card;
  border: 1px solid $border-light;
}

.selected-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: $spacing-sm;

  .selected-label {
    font-size: $font-size-sm;
    color: $text-secondary;
    font-weight: $font-weight-medium;
  }
}

.clear-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 12px;
  background: transparent;
  border: none;
  color: $text-tertiary;
  font-size: $font-size-sm;
  cursor: pointer;
  border-radius: $radius-pill;
  transition: all 0.2s ease;

  svg {
    width: 14px;
    height: 14px;
  }

  &:hover {
    color: $error-color;
    background: rgba(255, 59, 48, 0.08);
  }
}

.selected-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.selected-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: $gradient-primary;
  color: $text-white;
  border-radius: $radius-pill;
  font-size: $font-size-sm;
  font-weight: $font-weight-medium;
  box-shadow: 0 2px 8px rgba(109, 91, 255, 0.25);
}

.chip-remove {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  padding: 0;
  background: rgba(255, 255, 255, 0.2);
  border: none;
  border-radius: 50%;
  color: $text-white;
  cursor: pointer;
  transition: all 0.2s ease;

  svg {
    width: 10px;
    height: 10px;
  }

  &:hover {
    background: rgba(255, 255, 255, 0.35);
  }
}

.tag-groups {
  display: flex;
  flex-direction: column;
  gap: $spacing-xl;
}

.tag-group {
  background: $bg-card;
  border-radius: $radius-lg;
  padding: $spacing-lg;
  box-shadow: $shadow-card;
  border: 1px solid $border-light;
}

.group-title {
  font-size: $font-size-md;
  font-weight: $font-weight-semibold;
  color: $text-primary;
  margin: 0 0 $spacing-md 0;
}

.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.tag-chip {
  display: inline-flex;
  align-items: center;
  padding: 8px 18px;
  background: $bg-input;
  border: 1px solid $border-color;
  border-radius: $radius-pill;
  font-size: $font-size-sm;
  color: $text-secondary;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;

  &:hover {
    border-color: $primary-light;
    color: $primary-color;
    background: rgba(109, 91, 255, 0.06);
  }

  &.active {
    background: $gradient-primary;
    border-color: transparent;
    color: $text-white;
    font-weight: $font-weight-medium;
    box-shadow: 0 2px 8px rgba(109, 91, 255, 0.3);
  }
}

.loading-state,
.empty-state {
  text-align: center;
  padding: $spacing-xxl 0;
  color: $text-tertiary;

  .loading-spinner {
    width: 28px;
    height: 28px;
    margin: 0 auto $spacing-sm;
    border: 3px solid $border-color;
    border-top-color: $primary-color;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  .empty-icon {
    width: 72px;
    height: 72px;
    margin: 0 auto $spacing-sm;
    color: $border-primary;

    svg {
      width: 100%;
      height: 100%;
    }
  }

  p {
    font-size: $font-size-base;
    margin: 0;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.footer-action {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: $spacing-md $spacing-xl;
  background: linear-gradient(to top, $bg-page 60%, transparent);
  display: flex;
  justify-content: center;
  z-index: $z-sticky;
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 12px 48px;
  border-radius: $radius-pill;
  font-size: $font-size-base;
  font-weight: $font-weight-semibold;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;

  &.btn-confirm {
    background: $gradient-primary;
    color: $text-white;
    box-shadow: 0 4px 16px rgba(109, 91, 255, 0.35);
    min-width: 200px;

    &:hover {
      opacity: 0.9;
      transform: translateY(-1px);
      box-shadow: 0 6px 20px rgba(109, 91, 255, 0.45);
    }

    &:active {
      transform: translateY(0);
    }
  }
}
</style>
