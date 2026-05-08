<template>
  <view class="clothes-page">
    <view class="header">
      <view class="search-row">
        <input class="search-input" v-model="searchName" :placeholder="$t('clothesSearchPlaceholder')" confirm-type="search" @confirm="onSearch" />
        <view class="search-btn" @tap="onSearch">{{ $t('searchAction') }}</view>
      </view>
      <scroll-view
        class="category-tabs"
        scroll-x
        :show-scrollbar="false"
        scroll-with-animation
        :scroll-left="categoryScrollLeft"
        v-if="categoryList.length"
      >
        <view class="category-track">
          <view
            class="category-tab"
            v-for="item in categoryList"
            :key="item.ID || item.id"
            :id="categoryTabDomId(item)"
            :class="{ active: Number(item.ID || item.id) === Number(activeCategoryID) }"
            @tap="switchCategory(item)"
          >
            {{ categoryName(item) || $t('categoryDetail') }}
          </view>
        </view>
      </scroll-view>
    </view>

    <scroll-view class="list-wrap" scroll-y @scrolltolower="onScrollToLower">
      <view class="grid">
        <view class="card" v-for="item in displayList" :key="item.ID || item.id || item.name">
          <image class="card-image" :src="mainImage(item)" mode="aspectFill" />
          <view class="card-body">
            <text class="card-name">{{ goodName(item) || $t('unnamedGoods') }}</text>
            <view class="tag-list">
              <text class="tag" v-for="(tag, idx) in extractTags(item)" :key="idx">{{ tag }}</text>
            </view>
            <view class="btn-row">
              <view class="btn tryon" @tap="chooseTryon(item)">{{ $t('tryOnAction') }}</view>
              <view class="btn buy" @tap="goDetail(item)">{{ $t('buyAction') }}</view>
            </view>
          </view>
        </view>
      </view>

      <view class="load-more-wrap" v-if="displayList.length">
        <text class="all-loaded" v-if="loading">{{ $t('loading') }}</text>
        <text class="all-loaded" v-else-if="!hasMore">{{ $t('allLoadedText') }}</text>
      </view>
      <view style="height: 40rpx"></view>
    </scroll-view>
  </view>
</template>

<script setup>
import { nextTick, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getGoodList, getCategoryMobile } from '@/api/homePage.js'
import { getUrl, getExternalUrl } from '@/utils/url.js'
import { localText, t as i18nT } from '@/utils/i18n.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { usePlayHistoryStore } from '@/pinia/modules/playHistory.js'
import { setSelectedClothes } from '@/utils/tryon.js'

const langStore = useLangStore()
const playHistoryStore = usePlayHistoryStore()
const $t = langStore.$t
const searchName = ref('')
const page = ref(1)
const pageSize = 10
const hasMore = ref(true)
const loading = ref(false)
const goodsList = ref([])
const displayList = ref([])
const categoryList = ref([])
const activeCategoryID = ref(0)
const categoryScrollLeft = ref(0)

const TRYON_LOWER_KEYWORD_LANGS = ['zh', 'zh-TW', 'en', 'mn', 'th', 'hi', 'id']
const DEFAULT_TRYON_LOWER_KEYWORDS = ['pants', 'skirt', 'bottom', 'lower', 'trousers', 'jeans']

const normalizeKeyword = (value) => String(value || '').trim().toLowerCase()

const splitKeywordText = (value) => {
  const text = String(value || '').trim()
  if (!text) return []
  return text
    .split(/[\s,，;；|/]+/)
    .map(normalizeKeyword)
    .filter(Boolean)
}

const getTryonLowerKeywords = () => {
  const keywordSet = new Set(DEFAULT_TRYON_LOWER_KEYWORDS.map(normalizeKeyword))

  TRYON_LOWER_KEYWORD_LANGS.forEach((lang) => {
    ;[
      i18nT('tryonLowerKeywords', lang),
      i18nT('uploadLowerImage', lang),
      i18nT('tryonExampleLowerPrefix', lang),
    ].forEach((raw) => {
      splitKeywordText(raw).forEach((keyword) => keywordSet.add(keyword))
    })
  })

  return [...keywordSet]
}

const inferTryonPart = (item) => {
  const explicitPart = resolveLocaleText(
    item?.tryonPart ||
      item?.part ||
      item?.clothesPart ||
      item?.position ||
      item?.bodyPart ||
      item?.tryonType
  )

  const text = `${explicitPart} ${categoryText(item)} ${goodName(item)} ${extractTags(item).join(' ')}`
  const searchable = normalizeKeyword(text)
  const isLower = getTryonLowerKeywords().some((keyword) => searchable.includes(keyword))
  return isLower ? 'lower' : 'upper'
}

const parseJsonSafe = (raw) => {
  try {
    return JSON.parse(raw)
  } catch {
    return null
  }
}

const resolveLocaleText = (value) => {
  if (!value) return ''
  if (typeof value === 'string') {
    const text = value.trim()
    if (!text) return ''
    if (text.startsWith('{') || text.startsWith('[')) {
      const parsed = parseJsonSafe(text)
      if (parsed) {
        if (Array.isArray(parsed)) {
          return resolveLocaleText(parsed[0])
        }
        return localText(parsed, langStore.locale)
      }
    }
    return text
  }
  if (Array.isArray(value)) {
    return resolveLocaleText(value[0])
  }
  if (typeof value === 'object') {
    return localText(value, langStore.locale)
  }
  return String(value)
}

const categoryName = (item) => {
  return resolveLocaleText(item?.title || item?.name || item?.label)
}

const categoryIdentity = (item) => String(item?.ID ?? item?.id ?? '0')

const categoryTabDomId = (item) => {
  return `category-tab-${categoryIdentity(item).replace(/[^a-zA-Z0-9_-]/g, '_')}`
}

const centerCategoryTab = (item) => {
  if (!item) return

  const tabId = categoryTabDomId(item)
  const query = uni.createSelectorQuery()
  query.select('.category-tabs').boundingClientRect()
  query.select('.category-track').boundingClientRect()
  query.select(`#${tabId}`).boundingClientRect()
  query.exec((res) => {
    const wrap = res && res[0]
    const track = res && res[1]
    const tab = res && res[2]
    if (!wrap || !track || !tab) return

    const maxScroll = Math.max(0, Math.round((track.width || 0) - wrap.width))
    if (maxScroll <= 0) {
      categoryScrollLeft.value = 0
      return
    }

    const delta = tab.left - wrap.left
    const target = categoryScrollLeft.value + delta - (wrap.width - tab.width) / 2
    categoryScrollLeft.value = Math.max(0, Math.min(maxScroll, Math.round(target)))
  })
}

const goodName = (item) => {
  return resolveLocaleText(item?.nameI18n || item?.name || item?.titleI18n || item?.title)
}

const categoryText = (item) => {
  return resolveLocaleText(item?.categoryNameI18n || item?.categoryName || item?.categoryTitle)
}

const refreshDisplay = () => {
  displayList.value = [...goodsList.value]
}

const mainImage = (item) => {
  if (item.externalImagePath) return getExternalUrl(item.externalImagePath)
  return getUrl(item.imageUrl || item.picture || item.image || '')
}

const extractTags = (item) => {
  const tags = []

  const parseTagString = (value) => {
    const text = String(value || '').trim()
    if (!text) return []

    if (text.startsWith('[') || text.startsWith('{')) {
      const parsed = parseJsonSafe(text)
      if (Array.isArray(parsed)) return parsed
      if (parsed) return [parsed]
    }

    if (/[\n\r,;|，]/.test(text)) {
      return text
        .split(/[\n\r,;|，]+/)
        .map(v => v.trim())
        .filter(Boolean)
    }

    return [text]
  }

  const tagLabel = (value) => {
    if (!value) return ''

    if (typeof value === 'string') {
      const list = parseTagString(value)
      if (list.length === 1 && list[0] === value) {
        return resolveLocaleText(value)
      }
      return tagLabel(list[0])
    }

    if (typeof value === 'object') {
      return (
        resolveLocaleText(value.nameI18n) ||
        resolveLocaleText(value.name) ||
        resolveLocaleText(value.tagNameI18n) ||
        resolveLocaleText(value.tagName) ||
        resolveLocaleText(value.label) ||
        resolveLocaleText(value.tag) ||
        resolveLocaleText(value.titleI18n) ||
        resolveLocaleText(value.title) ||
        ''
      )
    }

    return resolveLocaleText(value)
  }

  const appendTag = (value) => {
    const text = tagLabel(value)
    if (text) tags.push(text)
  }

  if (Array.isArray(item.tagNames)) {
    item.tagNames.forEach(appendTag)
  } else if (typeof item.tagNames === 'string') {
    parseTagString(item.tagNames).forEach(appendTag)
  }

  if (Array.isArray(item.tags)) {
    item.tags.forEach(appendTag)
  } else if (typeof item.tags === 'string' && item.tags) {
    const parsed = parseJsonSafe(item.tags)
    if (Array.isArray(parsed)) {
      parsed.forEach(appendTag)
    } else if (parsed) {
      appendTag(parsed)
    } else {
      parseTagString(item.tags).forEach(appendTag)
    }
  } else if (item.tags && typeof item.tags === 'object') {
    appendTag(item.tags)
  }

  if (item.tag) appendTag(item.tag)
  if (item.tagName) appendTag(item.tagName)
  if (item.tag_name) appendTag(item.tag_name)
  return [...new Set(tags)].slice(0, 2)
}

const loadCategories = async () => {
  const res = await getCategoryMobile()
  if (res.code === 0 && Array.isArray(res.data)) {
    const allCategory = {
      ID: 0,
      id: 0,
      title: $t('all'),
    }
    categoryList.value = [allCategory, ...res.data]
    if (!categoryList.value.find(item => Number(item?.ID || item?.id || 0) === Number(activeCategoryID.value))) {
      activeCategoryID.value = 0
    }

    await nextTick()
    const activeItem = categoryList.value.find(item => Number(item?.ID || item?.id || 0) === Number(activeCategoryID.value))
    if (activeItem) {
      centerCategoryTab(activeItem)
    }
  }
}

const buildListParams = () => {
  const params = {
    page: page.value,
    pageSize,
    name: searchName.value,
    keyword: searchName.value,
    sort: 'id',
    order: 'desc',
  }
  if (activeCategoryID.value) {
    params.categoryID = activeCategoryID.value
  }
  return params
}

const loadList = async (reset = false) => {
  if (loading.value) return
  loading.value = true
  try {
    if (reset) {
      page.value = 1
      hasMore.value = true
      goodsList.value = []
    }

    const res = await getGoodList(buildListParams())

    if (res.code === 0 && res.data) {
      const list = Array.isArray(res.data.list) ? res.data.list : []
      if (reset) {
        goodsList.value = list
      } else {
        goodsList.value = goodsList.value.concat(list)
      }

      const total = Number(res.data.total || 0)
      if (total > 0) {
        hasMore.value = goodsList.value.length < total
      } else {
        hasMore.value = list.length >= pageSize
      }
      refreshDisplay()
    }
  } finally {
    loading.value = false
  }
}

const onSearch = () => {
  loadList(true)
}

const switchCategory = (category) => {
  const nextID = Number(category?.ID || category?.id || 0)
  if (nextID === Number(activeCategoryID.value)) return
  activeCategoryID.value = nextID
  nextTick(() => centerCategoryTab(category))
  loadList(true)
}

const loadMore = async () => {
  if (!hasMore.value) return
  page.value += 1
  await loadList(false)
}

const onScrollToLower = () => {
  if (loading.value || !hasMore.value) return
  loadMore()
}

const chooseTryon = (item) => {
  const image = item.externalImagePath ? getExternalUrl(item.externalImagePath) : (item.imageUrl || item.picture || '')
  if (!image) {
    uni.showToast({ title: $t('goodsImageMissing'), icon: 'none' })
    return
  }
  const tryonPart = inferTryonPart(item)
  setSelectedClothes(tryonPart === 'lower' ? { lowerImage: image } : { upperImage: image })
  uni.switchTab({ url: '/pages/tabBar/index' })
}

const goDetail = (item) => {
  const goodID = item.ID || item.id
  if (!goodID) {
    uni.showToast({ title: $t('goodsInfoInvalid'), icon: 'none' })
    return
  }

  playHistoryStore.saveBrowse(goodID, {
    imageUrl: item.externalImagePath || item.imageUrl || item.picture || item.image || '',
    title: goodName(item) || '',
  })

  uni.navigateTo({ url: `/pages/goodsDetails/goodsDetails?id=${goodID}` })
}

onShow(() => {
  if (categoryList.value.length === 0) {
    loadCategories().then(() => loadList(true))
    return
  }
  if (goodsList.value.length === 0) {
    loadList(true)
  }
})
</script>

<style lang="scss">
page {
  background: #f4f7fb;
}

.clothes-page {
  min-height: 100vh;
  background: radial-gradient(120% 80% at 100% -10%, #dbeafe 0%, transparent 60%), #f4f7fb;
  color: #0f172a;
  padding-top: calc(var(--status-bar-height, 0px) + 16rpx);
}

.header {
  padding: 0 20rpx 16rpx;
}

.search-row {
  display: flex;
  gap: 12rpx;
}

.search-input {
  flex: 1;
  height: 72rpx;
  border-radius: 999rpx;
  background: rgba(255, 255, 255, 0.95);
  border: 1rpx solid rgba(15, 23, 42, 0.12);
  padding: 0 24rpx;
  color: #0f172a;
  font-size: 24rpx;
}

.search-btn {
  width: 120rpx;
  height: 72rpx;
  border-radius: 999rpx;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24rpx;
  color: #fff;
  box-shadow: 0 10rpx 26rpx rgba(37, 99, 235, 0.28);
}

.category-tabs {
  margin-top: 16rpx;
  white-space: nowrap;
}

.category-track {
  display: inline-flex;
  gap: 12rpx;
  padding-right: 12rpx;
}

.category-tab {
  height: 64rpx;
  min-width: 160rpx;
  padding: 0 26rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.12);
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(15, 23, 42, 0.58);
  background: rgba(255, 255, 255, 0.8);
}

.category-tab.active {
  color: #0f172a;
  border-color: rgba(37, 99, 235, 0.45);
  background: rgba(219, 234, 254, 0.75);
}

.list-wrap {
  height: calc(100vh - var(--status-bar-height, 0px) - 170rpx);
  padding: 0 20rpx;
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16rpx;
}

.card {
  background: rgba(255, 255, 255, 0.95);
  border: 1rpx solid rgba(15, 23, 42, 0.08);
  border-radius: 14rpx;
  overflow: hidden;
  box-shadow: 0 8rpx 22rpx rgba(15, 23, 42, 0.05);
}

.card-image {
  width: 100%;
  height: 260rpx;
}

.card-body {
  padding: 12rpx;
}

.card-name {
  font-size: 24rpx;
  min-height: 68rpx;
  line-height: 34rpx;
}

.tag-list {
  min-height: 42rpx;
}

.tag {
  display: inline-block;
  margin-right: 8rpx;
  margin-bottom: 8rpx;
  padding: 4rpx 10rpx;
  border-radius: 999rpx;
  font-size: 18rpx;
  background: rgba(219, 234, 254, 0.8);
  color: rgba(30, 64, 175, 0.9);
}

.btn-row {
  display: flex;
  gap: 8rpx;
}

.btn {
  flex: 1;
  height: 54rpx;
  border-radius: 999rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22rpx;
}

.btn.tryon {
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  color: #fff;
}

.btn.buy {
  background: rgba(15, 23, 42, 0.06);
  color: #0f172a;
}

.load-more-wrap {
  margin-top: 24rpx;
  text-align: center;
}

.load-more {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 280rpx;
  height: 66rpx;
  border-radius: 999rpx;
  background: rgba(255, 255, 255, 0.92);
  border: 1rpx solid rgba(15, 23, 42, 0.1);
  color: #0f172a;
  font-size: 24rpx;
}

.all-loaded {
  color: rgba(15, 23, 42, 0.45);
  font-size: 22rpx;
}
</style>
