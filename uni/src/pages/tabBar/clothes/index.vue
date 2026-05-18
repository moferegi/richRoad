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
        :scroll-into-view="categoryScrollIntoView"
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
      <view class="mobile-scroller-wrapper">
        <view class="waterfall-column">
          <view class="card" :class="cardClass(colIndex, 'left')" v-for="(item, colIndex) in leftColumnCards" :key="`${item.ID || item.id || item.name || colIndex}-left`">
            <view class="card-media" :class="cardMediaClass(colIndex, 'left')" @tap="previewImage(item)">
              <LazyImage class="card-image" :src="mainImage(item)" mode="scaleToFill" />
              <view class="card-cinema-shadow"></view>
              <view class="card-cinema-glow"></view>
              <view class="card-body">
                <text class="card-name">{{ goodName(item) || $t('unnamedGoods') }}</text>
                <view class="btn-row">
                  <view class="btn tryon" @tap.stop="chooseTryon(item)">{{ $t('tryOnAction') }}</view>
                  <view class="btn buy" @tap.stop="goDetail(item)">{{ $t('buyAction') }}</view>
                </view>
              </view>
            </view>
          </view>
        </view>

        <view class="waterfall-column">
          <view class="card" :class="cardClass(colIndex, 'right')" v-for="(item, colIndex) in rightColumnCards" :key="`${item.ID || item.id || item.name || colIndex}-right`">
            <view class="card-media" :class="cardMediaClass(colIndex, 'right')" @tap="previewImage(item)">
              <LazyImage class="card-image" :src="mainImage(item)" mode="scaleToFill" />
              <view class="card-cinema-shadow"></view>
              <view class="card-cinema-glow"></view>
              <view class="card-body">
                <text class="card-name">{{ goodName(item) || $t('unnamedGoods') }}</text>
                <view class="btn-row">
                  <view class="btn tryon" @tap.stop="chooseTryon(item)">{{ $t('tryOnAction') }}</view>
                  <view class="btn buy" @tap.stop="goDetail(item)">{{ $t('buyAction') }}</view>
                </view>
              </view>
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
import { computed, nextTick, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getGoodList, getCategoryMobile } from '@/api/homePage.js'
import { getUrl, getExternalUrl } from '@/utils/url.js'
import { t as i18nT } from '@/utils/i18n.js'
import { useI18nDisplay } from '@/composables/useI18nDisplay.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { usePlayHistoryStore } from '@/pinia/modules/playHistory.js'
import { setSelectedClothes } from '@/utils/tryon.js'
import LazyImage from '@/components/lazy-image/lazy-image.vue'

const langStore = useLangStore()
const playHistoryStore = usePlayHistoryStore()
const $t = langStore.$t
const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')
const { resolveDisplayText } = useI18nDisplay(locale)
const lastLoadedLocale = ref('')
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
const categoryScrollIntoView = ref('')

const TRYON_LOWER_KEYWORD_LANGS = ['zh', 'zh-TW', 'en', 'mn', 'th', 'hi', 'id', 'vi', 'ar', 'ja', 'ko', 'ms']
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
  if (Array.isArray(value)) {
    return resolveLocaleText(value[0])
  }

  if (typeof value === 'string') {
    const text = value.trim()
    if (!text) return ''
    if (text.startsWith('[')) {
      const parsed = parseJsonSafe(text)
      if (Array.isArray(parsed) && parsed.length > 0) {
        return resolveLocaleText(parsed[0])
      }
    }
    return resolveDisplayText(text, text)
  }

  return resolveDisplayText(value, '')
}

const categoryName = (item) => {
  if (Number(item?.ID || item?.id || 0) === 0) {
    return $t('all')
  }
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

    const tabOffset = tab.left - track.left
    const target = tabOffset - (wrap.width - tab.width) / 2
    categoryScrollLeft.value = Math.max(0, Math.min(maxScroll, Math.round(target)))
  })
}

const ensureCategoryCentered = (item) => {
  if (!item) return
  categoryScrollIntoView.value = categoryTabDomId(item)
  nextTick(() => {
    centerCategoryTab(item)
    setTimeout(() => centerCategoryTab(item), 80)
    setTimeout(() => centerCategoryTab(item), 180)
  })
}

const goodName = (item) => {
  return resolveLocaleText(item?.nameI18n || item?.name || item?.titleI18n || item?.title)
}

const goodDesc = (item) => {
  return resolveLocaleText(
    item?.descriptionI18n ||
    item?.description ||
    item?.descI18n ||
    item?.desc ||
    item?.summaryI18n ||
    item?.summary
  )
}

const categoryText = (item) => {
  return resolveLocaleText(item?.categoryNameI18n || item?.categoryName || item?.categoryTitle)
}

const refreshDisplay = () => {
  displayList.value = [...goodsList.value]
}

const leftColumnCards = computed(() => {
  return displayList.value.filter((_, index) => index % 2 === 0)
})

const rightColumnCards = computed(() => {
  return displayList.value.filter((_, index) => index % 2 === 1)
})

const cardMediaClass = (columnIndex, side) => {
  if (columnIndex === 0) return side === 'left' ? 'media-tall' : 'media-short'
  const odd = columnIndex % 2 === 1
  if (side === 'left') return odd ? 'media-short' : 'media-tall'
  return odd ? 'media-tall' : 'media-short'
}

const cardClass = (columnIndex, side) => {
  if (columnIndex === 0) return 'card-flat'
  const odd = columnIndex % 2 === 1
  if (side === 'left') return odd ? 'card-shift-a' : 'card-shift-b'
  return odd ? 'card-shift-c' : 'card-shift-d'
}

const mainImage = (item) => {
  if (item.externalImagePath) return getExternalUrl(item.externalImagePath)
  return getUrl(item.imageUrl || item.picture || item.image || '')
}

const getFileSizeAsync = (filePath) => {
  return new Promise((resolve) => {
    uni.getFileInfo({
      filePath,
      success: (res) => resolve(Number(res.size || 0)),
      fail: () => resolve(0),
    })
  })
}

const downloadFileAsync = (url) => {
  return new Promise((resolve, reject) => {
    uni.downloadFile({
      url,
      success: resolve,
      fail: reject,
    })
  })
}

const remoteImageSizeCache = Object.create(null)

const resolveRemoteImageSize = async (rawUrl) => {
  const normalizedUrl = String(rawUrl || '').trim()
  if (!normalizedUrl) return 0

  if (Object.prototype.hasOwnProperty.call(remoteImageSizeCache, normalizedUrl)) {
    return Number(remoteImageSizeCache[normalizedUrl] || 0)
  }

  try {
    const downloadRes = await downloadFileAsync(normalizedUrl)
    const tempPath = String(downloadRes?.tempFilePath || '').trim()
    if (!tempPath) {
      remoteImageSizeCache[normalizedUrl] = 0
      return 0
    }
    const size = await getFileSizeAsync(tempPath)
    remoteImageSizeCache[normalizedUrl] = Number(size || 0)
    return Number(size || 0)
  } catch (e) {
    remoteImageSizeCache[normalizedUrl] = 0
    return 0
  }
}

const resolveImageLinkByFields = (item, primaryFields = [], externalFields = []) => {
  for (const field of primaryFields) {
    const raw = String(item?.[field] || '').trim()
    if (!raw) continue
    return /^https?:\/\//i.test(raw) || /^data:/i.test(raw) ? raw : getUrl(raw)
  }

  for (const field of externalFields) {
    const raw = String(item?.[field] || '').trim()
    if (!raw) continue
    return /^https?:\/\//i.test(raw) || /^data:/i.test(raw) ? raw : getExternalUrl(raw)
  }

  return ''
}

const getUpperImageLink = (item) => {
  return resolveImageLinkByFields(item, ['upperImage', 'upper_image', 'tryonUpperImage'], ['upperExternalImagePath', 'upper_external_image_path'])
}

const getLowerImageLink = (item) => {
  return resolveImageLinkByFields(item, ['lowerImage', 'lower_image', 'tryonLowerImage'], ['lowerExternalImagePath', 'lower_external_image_path'])
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
      name: 'all',
    }
    categoryList.value = [allCategory, ...res.data]
    if (!categoryList.value.find(item => Number(item?.ID || item?.id || 0) === Number(activeCategoryID.value))) {
      activeCategoryID.value = 0
    }

    await nextTick()
    const activeItem = categoryList.value.find(item => Number(item?.ID || item?.id || 0) === Number(activeCategoryID.value))
    if (activeItem) {
      ensureCategoryCentered(activeItem)
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
  const keyword = String(searchName.value || '').trim()
  if (!keyword) {
    uni.showToast({ title: $t('searchContentRequired'), icon: 'none' })
    return
  }
  searchName.value = keyword
  loadList(true)
}

const switchCategory = (category) => {
  const nextID = Number(category?.ID || category?.id || 0)
  if (nextID === Number(activeCategoryID.value)) return
  activeCategoryID.value = nextID
  ensureCategoryCentered(category)
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

const chooseTryon = async (item) => {
  const upperImage = getUpperImageLink(item)
  const lowerImage = getLowerImageLink(item)

  if (upperImage && lowerImage) {
    uni.showLoading({ title: $t('loading'), mask: true })
    try {
      const [upperSizeBytes, lowerSizeBytes] = await Promise.all([
        resolveRemoteImageSize(upperImage),
        resolveRemoteImageSize(lowerImage),
      ])
      setSelectedClothes({
        upperImage,
        lowerImage,
        upperSizeBytes,
        lowerSizeBytes,
      })
    } finally {
      uni.hideLoading()
    }
    uni.switchTab({ url: '/pages/tabBar/index' })
    return
  }

  const image = mainImage(item)
  if (!image) {
    uni.showToast({ title: $t('goodsImageMissing'), icon: 'none' })
    return
  }

  const tryonPart = inferTryonPart(item)
  uni.showLoading({ title: $t('loading'), mask: true })
  try {
    const sizeBytes = await resolveRemoteImageSize(image)
    setSelectedClothes(
      tryonPart === 'lower'
        ? { lowerImage: image, lowerSizeBytes: sizeBytes }
        : { upperImage: image, upperSizeBytes: sizeBytes }
    )
  } finally {
    uni.hideLoading()
  }

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

const previewImage = (item) => {
  const current = mainImage(item)
  if (!current) {
    uni.showToast({ title: $t('goodsImageMissing'), icon: 'none' })
    return
  }

  const urls = displayList.value
    .map((entry) => mainImage(entry))
    .filter((url) => !!url)

  if (!urls.length) {
    uni.showToast({ title: $t('goodsImageMissing'), icon: 'none' })
    return
  }

  uni.previewImage({
    current,
    urls,
    loop: true,
  })
}

const reloadLocaleSensitiveData = async () => {
  await loadCategories()
  await loadList(true)
}

const markLocaleLoaded = () => {
  lastLoadedLocale.value = locale.value
}

onShow(() => {
  const localeChanged = !!lastLoadedLocale.value && lastLoadedLocale.value !== locale.value
  if (localeChanged) {
    reloadLocaleSensitiveData().finally(() => {
      markLocaleLoaded()
    })
    return
  }

  if (categoryList.value.length === 0) {
    loadCategories()
      .then(() => loadList(true))
      .finally(() => {
        markLocaleLoaded()
      })
    return
  }

  if (goodsList.value.length === 0) {
    loadList(true).finally(() => {
      markLocaleLoaded()
    })
    return
  }

  markLocaleLoaded()
})
</script>

<style lang="scss">
page {
  background: #f4f7fb;
  overflow-x: hidden;
}

.clothes-page {
  --fixed-header-height: 184rpx;
  min-height: 100vh;
  background: radial-gradient(120% 80% at 100% -10%, #dbeafe 0%, transparent 60%), #f4f7fb;
  color: #0f172a;
  padding-top: 0;
  overflow-x: hidden;
}

.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 60;
  padding: calc(var(--status-bar-height, 0px) + 16rpx) 20rpx 16rpx;
  background: linear-gradient(180deg, rgba(244, 247, 251, 0.96) 0%, rgba(244, 247, 251, 0.9) 72%, rgba(244, 247, 251, 0));
  backdrop-filter: blur(12px);
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
  margin-top: calc(var(--status-bar-height, 0px) + var(--fixed-header-height));
  height: calc(100vh - var(--status-bar-height, 0px) - var(--fixed-header-height));
  padding: 0 20rpx;
  width: 100%;
  box-sizing: border-box;
  overflow-x: hidden;
}

.mobile-scroller-wrapper {
  display: flex;
  align-items: flex-start;
  gap: 18rpx;
  width: 100%;
  box-sizing: border-box;
  overflow-x: hidden;
}

.waterfall-column {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 28rpx;
}

.card {
  border-radius: 24rpx;
  overflow: hidden;
  box-shadow: 0 20rpx 40rpx rgba(15, 23, 42, 0.18);
  background: #ffffff;
  position: relative;
}

.card.card-flat {
  transform: translateY(0) rotate(0.18deg);
}

.card.card-shift-a {
  transform: rotate(-0.16deg);
}

.card.card-shift-b {
  transform: rotate(0.14deg);
}

.card.card-shift-c {
  transform: rotate(-0.18deg);
}

.card.card-shift-d {
  transform: rotate(0.12deg);
}

.card-media {
  position: relative;
  height: 560rpx;
  overflow: hidden;
  background: transparent;
}

.card-media.media-base {
  height: 560rpx;
}

.card-media.media-short {
  height: 520rpx;
}

.card-media.media-tall {
  height: 560rpx;
}

.card-image {
  width: 100%;
  height: 100%;
  filter: saturate(1.04) contrast(1.01);
}

.card-cinema-shadow,
.card-cinema-glow {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.card-cinema-shadow {
  background:
    radial-gradient(120% 88% at 50% 10%, rgba(255, 255, 255, 0.18) 0%, transparent 62%),
    linear-gradient(180deg, rgba(15, 23, 42, 0) 56%, rgba(15, 23, 42, 0.18) 100%);
}

.card-cinema-glow {
  background: radial-gradient(84% 36% at 50% 0%, rgba(255, 255, 255, 0.2) 0%, rgba(255, 255, 255, 0) 100%);
}

.clothes-page .gva-lazy-image {
  background: transparent;
}

.clothes-page .gva-lazy-image__status {
  background: transparent !important;
}

.clothes-page .gva-lazy-image__status--error {
  background: rgba(71, 85, 105, 0.12) !important;
}

.clothes-page .gva-lazy-image__status-text {
  opacity: 0;
}

.card-body {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 3;
  padding: 14rpx 16rpx;
  box-sizing: border-box;
  background: linear-gradient(180deg, rgba(100, 116, 139, 0.06) 0%, rgba(100, 116, 139, 0.34) 48%, rgba(71, 85, 105, 0.56) 100%);
  backdrop-filter: blur(11px);
  border-top: 1rpx solid rgba(226, 232, 240, 0.28);
}

.card-name {
  display: -webkit-box;
  color: #f8fafc;
  font-size: 24rpx;
  min-height: 56rpx;
  line-height: 33rpx;
  line-clamp: 2;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-shadow: 0 2rpx 10rpx rgba(15, 23, 42, 0.35);
}

.btn-row {
  display: flex;
  gap: 12rpx;
  margin-top: 8rpx;
}

.btn {
  flex: 1;
  height: 56rpx;
  border-radius: 999rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 21rpx;
  font-weight: 600;
}

.btn.tryon {
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  color: #fff;
  box-shadow: 0 8rpx 20rpx rgba(37, 99, 235, 0.3);
}

.btn.buy {
  background: rgba(241, 245, 249, 0.86);
  color: #0f172a;
}

.load-more-wrap {
  margin-top: 80rpx;
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
