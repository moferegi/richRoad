<template>
  <view class="history-page">
    <view class="nav">
      <view class="nav-btn" @tap="goBack">
        <uni-icons type="left" size="20" color="#0f172a" />
      </view>
      <text class="nav-title">{{ $t('tryonHistory') }}</text>
      <view class="nav-btn" @tap="loadHistory(true)">
        <uni-icons type="reload" size="18" color="#0f172a" />
      </view>
    </view>

    <scroll-view class="filter-row" scroll-x :show-scrollbar="false">
      <view class="filter-chip" :class="{ active: currentFilter === 'all' }" @tap="currentFilter = 'all'">{{ $t('all') }}</view>
      <view class="filter-chip" :class="{ active: currentFilter === 'tryon' }" @tap="currentFilter = 'tryon'">{{ $t('tryonClothes') }}</view>
      <view class="filter-chip" :class="{ active: currentFilter === 'shoe' }" @tap="currentFilter = 'shoe'">{{ $t('tryonShoes') }}</view>
    </scroll-view>

    <scroll-view class="list-wrap" scroll-y @scrolltolower="onReachBottom">
      <view class="card" v-for="item in filteredList" :key="item._key" @tap="handleCardTap(item)">
        <view class="thumb-wrap">
          <LazyImage class="thumb" :src="thumbnailImage(item)" mode="aspectFit" />
        </view>
        <view class="card-main">
          <view class="line">
            <text class="type">{{ typeLabel(item) }}</text>
            <text class="status" :class="statusClass(item.status)">{{ statusLabel(item.status) }}</text>
          </view>
          <text class="no">{{ $t('tryonTaskNo') }}：{{ item.taskNo || '-' }}</text>
          <text class="time">{{ formatTime(item.CreatedAt || item.savedAt) }}</text>
          <text class="error" v-if="item.status === 'failed' && item.errorMessage">{{ item.errorMessage }}</text>
          <view class="card-actions">
            <view class="mini-btn danger" @tap.stop="removeHistory(item)">{{ $t('delete') }}</view>
          </view>
        </view>
      </view>

      <view class="empty" v-if="filteredList.length === 0">
        <text>{{ $t('tryonNoTask') }}</text>
      </view>

      <view class="load-more-row" v-if="showLoadMoreButton">
        <view class="load-more-btn" :class="{ disabled: loadingMoreHistory }" @tap="loadMoreHistory">
          {{ loadingMoreHistory ? $t('loading') : $t('loadMore') }}
        </view>
      </view>

      <view style="height: 30rpx"></view>
    </scroll-view>

    <view class="preview-mask" :class="{ 'preview-mask-passthrough': previewNativeOpening }" v-if="previewVisible" @tap="closePreview">
      <view class="preview-panel" @tap.stop>
        <view class="preview-head">
          <text class="preview-title">{{ $t('imagePreviewTitle') }}</text>
          <view class="preview-close" @tap="closePreview">
            <uni-icons type="closeempty" size="20" color="#0f172a" />
          </view>
        </view>

        <view class="preview-tabs">
          <view
            class="preview-tab"
            :class="{ active: previewTab === 'result', disabled: !previewResultImages.length }"
            @tap="switchPreviewTab('result')"
          >
            {{ $t('previewResultTab') }}
          </view>
          <view
            class="preview-tab"
            :class="{ active: previewTab === 'origin', disabled: !previewOriginImages.length }"
            @tap="switchPreviewTab('origin')"
          >
            {{ $t('previewOriginTab') }}
          </view>
        </view>

        <view class="preview-body" v-if="currentPreviewImages.length">
          <swiper class="preview-swiper" :current="previewCurrentIndex" @change="onPreviewSwiperChange">
            <swiper-item v-for="(img, idx) in currentPreviewImages" :key="`${img}_${idx}`">
              <view class="preview-image-wrap" @tap="previewCurrentImageByIndex(idx)">
                <LazyImage class="preview-image" :src="getUrl(img)" mode="aspectFit" />
              </view>
            </swiper-item>
          </swiper>
          <text class="preview-index">{{ previewCurrentIndex + 1 }}/{{ currentPreviewImages.length }}</text>
        </view>

        <view class="preview-empty" v-else>
          <text>{{ $t('noImagePreview') }}</text>
        </view>

        <view class="preview-actions">
          <text class="preview-save-tip">{{ $t('previewLongPressSaveHint') }}</text>
          <view class="preview-actions-row">
            <view class="preview-action-btn secondary" v-if="canComparePreview" @tap="openCompareFromPreview">{{ compareButtonText }}</view>
            <view class="preview-action-btn" v-if="canDownloadPreview" @tap="downloadCurrentImage">{{ $t('downloadAction') }}</view>
          </view>
        </view>
      </view>
    </view>

    <view class="h5-image-preview-mask" v-if="h5ImagePreviewVisible" @tap="closeH5ImagePreview">
      <view class="h5-image-preview-head" @tap.stop>
        <text class="h5-image-preview-index">{{ h5ImagePreviewCurrent + 1 }}/{{ h5ImagePreviewList.length }}</text>
        <view class="h5-image-preview-actions">
          <view class="h5-image-preview-save" @tap="downloadH5PreviewCurrentImage">{{ $t('downloadAction') }}</view>
          <view class="h5-image-preview-close" @tap="closeH5ImagePreview">
            <uni-icons type="closeempty" size="22" color="#ffffff" />
          </view>
        </view>
      </view>
      <swiper class="h5-image-preview-swiper" :current="h5ImagePreviewCurrent" @change="onH5ImagePreviewChange">
        <swiper-item v-for="(img, idx) in h5ImagePreviewList" :key="`h5_preview_${idx}`">
          <view class="h5-image-preview-wrap" @tap.stop @longpress.stop.prevent="downloadH5PreviewImage(img)">
            <LazyImage class="h5-image-preview-img" :src="img" mode="aspectFit" />
          </view>
        </swiper-item>
      </swiper>
      <view class="h5-image-preview-tip" @tap.stop>{{ $t('previewLongPressSaveHint') }}</view>
    </view>

    <view class="preview-mask" v-if="compareVisible" @tap="closeCompare">
      <view class="preview-panel compare-panel" @tap.stop>
        <view class="preview-head compare-head--compact">
          <view class="preview-close" @tap="closeCompare">
            <uni-icons type="closeempty" size="20" color="#0f172a" />
          </view>
        </view>

        <view
          class="compare-stage"
          @touchstart.stop.prevent="onCompareStageTouchStart"
          @touchmove.stop.prevent="onCompareStageTouchMove"
          @touchend.stop="onCompareStageTouchEnd"
          @touchcancel.stop="onCompareStageTouchEnd"
        >
          <LazyImage class="compare-image" :src="getUrl(compareOriginPreview)" mode="aspectFit" :style="compareImageStyle" />
          <view class="compare-result-layer" :style="{ width: `${comparePercent}%` }">
            <LazyImage class="compare-image compare-result-image" :src="getUrl(compareResultPreview)" mode="aspectFit" :style="compareResultInnerStyle" />
          </view>
          <view class="compare-divider" :style="{ left: `${comparePercent}%` }"></view>

          <view v-if="compareGuideVisible" class="compare-guide-stage-overlay">
            <view class="compare-guide-bubble compare-guide-bubble--stage">
              <text class="compare-guide-title">{{ $t('compareGuideTitle') }}</text>
              <text>{{ $t('compareGuideStageTip') }}</text>
            </view>
            <view class="compare-guide-pointer compare-guide-pointer--stage"></view>
            <view class="compare-guide-divider-hint" :style="{ left: `${comparePercent}%` }">
              <view class="compare-guide-hand">
                <view class="compare-guide-hand-dot"></view>
              </view>
              <text class="compare-guide-divider-text">{{ $t('compareGuideDividerTip') }}</text>
            </view>
          </view>
        </view>

        <view class="compare-slider-wrap">
          <view v-if="compareGuideVisible" class="compare-guide-slider-tip">
            <view class="compare-guide-bubble compare-guide-bubble--slider">
              <text>{{ $t('compareGuideSliderTip') }}</text>
            </view>
            <view class="compare-guide-pointer compare-guide-pointer--slider"></view>
          </view>

          <view class="compare-slider-label">
            <text>{{ compareLeftLabel }}</text>
            <text>{{ compareRightLabel }}</text>
          </view>
          <slider
            class="compare-slider"
            :value="comparePercent"
            :min="0"
            :max="100"
            :step="1"
            activeColor="#2563eb"
            backgroundColor="rgba(15,23,42,0.12)"
            block-color="#ffffff"
            :block-size="20"
            @changing="onCompareSliderChange"
            @change="onCompareSliderChange"
          />

          <view class="compare-slider-label compare-slider-label--zoom">
            <text>1x</text>
            <text>{{ (compareZoomPercent / 100).toFixed(2) }}x</text>
            <text class="compare-zoom-max" @tap="setCompareZoom(240)">MAX</text>
          </view>
          <slider
            class="compare-slider compare-zoom-slider"
            :value="compareZoomPercent"
            :min="100"
            :max="240"
            :step="5"
            activeColor="#0ea5e9"
            backgroundColor="rgba(15,23,42,0.12)"
            block-color="#ffffff"
            :block-size="18"
            @changing="onCompareZoomChange"
            @change="onCompareZoomChange"
          />

          <view class="compare-guide-action-row" v-if="compareGuideVisible">
            <view class="compare-guide-action-btn" @tap.stop="dismissCompareGuide">{{ $t('compareGuideGotIt') }}</view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, nextTick, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { getMyTryonTaskList, deleteMyTryonTask } from '@/api/tryonTask.js'
import { resolveApiMessage } from '@/utils/i18n.js'
import { getUrl } from '@/utils/url.js'
import { getTryonLocalHistory, removeTryonLocalHistoryByTaskKey } from '@/utils/tryon.js'
import LazyImage from '@/components/lazy-image/lazy-image.vue'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const PAGE_SIZE = 10
const currentFilter = ref('all')
const allList = ref([])
const localTaskKeys = ref([])
const localListCache = ref([])
const serverList = ref([])
const mergedList = ref([])
const serverPage = ref(1)
const serverHasMore = ref(false)
const displayPage = ref(1)
const hasMoreHistory = ref(false)
const reachedBottom = ref(false)
const loadingMoreHistory = ref(false)

const previewVisible = ref(false)
const previewTab = ref('result')
const previewResultImages = ref([])
const previewOriginImages = ref([])
const previewCurrentIndex = ref(0)
const compareVisible = ref(false)
const compareMode = ref('tryon')
const comparePercent = ref(50)
const compareZoomPercent = ref(100)
const compareStageWidthPx = ref(0)
const compareStageHeightPx = ref(0)
const previewNativeOpening = ref(false)
const h5ImagePreviewVisible = ref(false)
const h5ImagePreviewList = ref([])
const h5ImagePreviewCurrent = ref(0)
const compareDragging = ref(false)
const compareDragStartX = ref(0)
const compareDragStartY = ref(0)
const compareBaseOffsetX = ref(0)
const compareBaseOffsetY = ref(0)
const compareOffsetX = ref(0)
const compareOffsetY = ref(0)
const comparePinching = ref(false)
const comparePinchStartDistance = ref(0)
const comparePinchStartZoomPercent = ref(100)
const comparePinchStartCenterX = ref(0)
const comparePinchStartCenterY = ref(0)
const comparePinchBaseOffsetX = ref(0)
const comparePinchBaseOffsetY = ref(0)
const compareGuideVisible = ref(false)
const COMPARE_GUIDE_SEEN_KEY = 'tryon-compare-guide-seen-v1'

const getTaskKey = (item) => String(item?.taskNo || item?.requestID || '').trim()

const normalizeServerType = (item) => {
  if (item.sceneType === 'shoes') return 'shoe'
  return 'tryon'
}

const recordType = (item) => {
  if (item.roomType === 'shoe' || item.sceneType === 'shoes') return 'shoe'
  return normalizeServerType(item)
}

const filteredList = computed(() => {
  if (currentFilter.value === 'all') return allList.value
  return allList.value.filter(item => recordType(item) === currentFilter.value)
})

const currentPreviewImages = computed(() => {
  if (previewTab.value === 'origin') {
    return previewOriginImages.value
  }
  return previewResultImages.value
})

const canCompareTryonPreview = computed(() => previewOriginImages.value.length > 0 && previewResultImages.value.length > 0)
const canComparePreview = computed(() => canCompareTryonPreview.value)
const canDownloadPreview = computed(() => currentPreviewImages.value.length > 0)
const compareButtonText = computed(() => $t.value('compareImageButton'))
const compareOriginPreview = computed(() => previewOriginImages.value[0] || '')
const compareResultPreview = computed(() => previewResultImages.value[0] || '')
const compareLeftLabel = computed(() => $t.value('previewOriginTab'))
const compareRightLabel = computed(() => $t.value('previewResultTab'))
const compareZoomScale = computed(() => Math.max(1, Number(compareZoomPercent.value || 100) / 100))
const comparePanRangeX = computed(() => {
  const width = Number(compareStageWidthPx.value || 0)
  if (width <= 0) return 0
  return Math.max(0, ((compareZoomScale.value - 1) * width) / 2)
})
const comparePanRangeY = computed(() => {
  const height = Number(compareStageHeightPx.value || 0)
  if (height <= 0) return 0
  return Math.max(0, ((compareZoomScale.value - 1) * height) / 2)
})
const compareImageStyle = computed(() => ({
  transform: `translate3d(${compareOffsetX.value}px, ${compareOffsetY.value}px, 0) scale(${compareZoomScale.value})`,
  transformOrigin: 'center center',
}))
const compareResultInnerStyle = computed(() => {
  const width = compareStageWidthPx.value > 0 ? `${compareStageWidthPx.value}px` : '100%'
  return {
    width,
    height: '100%',
    transform: `translate3d(${compareOffsetX.value}px, ${compareOffsetY.value}px, 0) scale(${compareZoomScale.value})`,
    transformOrigin: 'center center',
  }
})

const showLoadMoreButton = computed(() => hasMoreHistory.value && reachedBottom.value)

const typeLabel = (item) => {
  const type = recordType(item)
  if (type === 'shoe') return $t.value('tryonShoes')
  return $t.value('tryonClothes')
}

const statusLabel = (status) => {
  if (status === 'success') return $t.value('tryonStatusSuccess')
  if (status === 'failed') return $t.value('tryonStatusFailed')
  return $t.value('tryonStatusProcessing')
}

const statusClass = (status) => {
  if (status === 'success') return 'success'
  if (status === 'failed') return 'failed'
  return 'processing'
}

const formatTime = (time) => {
  if (!time) return '-'
  const date = new Date(time)
  if (Number.isNaN(date.getTime())) return '-'
  const Y = date.getFullYear()
  const M = String(date.getMonth() + 1).padStart(2, '0')
  const D = String(date.getDate()).padStart(2, '0')
  const h = String(date.getHours()).padStart(2, '0')
  const m = String(date.getMinutes()).padStart(2, '0')
  return `${Y}-${M}-${D} ${h}:${m}`
}

const mergeHistory = (serverList, localList) => {
  const map = new Map()

  ;(localList || []).forEach((item, index) => {
    const key = item.taskNo || item.requestID || `local_${index}`
    map.set(key, {
      ...item,
      _key: key,
    })
  })

  ;(serverList || []).forEach((item) => {
    const key = item.taskNo || item.requestID || String(item.ID || Math.random())
    const old = map.get(key)
    map.set(key, {
      ...old,
      ...item,
      _key: key,
    })
  })

  return Array.from(map.values()).sort((a, b) => {
    const ta = new Date(a.CreatedAt || a.savedAt || 0).getTime()
    const tb = new Date(b.CreatedAt || b.savedAt || 0).getTime()
    return tb - ta
  })
}

const extractImageValues = (value) => {
  if (!value) return []

  if (Array.isArray(value)) {
    return value.flatMap(item => extractImageValues(item))
  }

  if (typeof value === 'string') {
    const text = value.trim()
    if (!text) return []

    if (text.startsWith('[')) {
      try {
        const parsed = JSON.parse(text)
        return extractImageValues(parsed)
      } catch {
        // ignore parse error and keep splitting below
      }
    }

    return text
      .split(/[\n\r,;|，]+/)
      .map(item => item.trim())
      .filter(Boolean)
  }

  return []
}

const uniqueImageList = (list) => {
  const dedup = new Set()
  ;(list || []).forEach((item) => {
    const val = String(item || '').trim()
    if (val) dedup.add(val)
  })
  return Array.from(dedup)
}

const thumbnailImage = (item) => {
  const resultList = uniqueImageList(extractImageValues(item?.resultImage))
  if (resultList.length) {
    return getUrl(resultList[0])
  }

  const sourceList = uniqueImageList([
    ...extractImageValues(item?.templateImage),
    ...extractImageValues(item?.sourceImage),
  ])
  if (sourceList.length) {
    return getUrl(sourceList[0])
  }

  return ''
}

const syncVisibleList = () => {
  mergedList.value = mergeHistory(serverList.value, localListCache.value)
  const visibleCount = displayPage.value * PAGE_SIZE
  allList.value = mergedList.value.slice(0, visibleCount)
  hasMoreHistory.value = mergedList.value.length > allList.value.length || serverHasMore.value
}

const fetchServerPage = async (pageNo) => {
  const res = await getMyTryonTaskList({
    page: pageNo,
    pageSize: PAGE_SIZE,
  })

  if (res.code !== 0 || !res.data) {
    if (pageNo === 1) {
      serverList.value = []
    }
    serverHasMore.value = false
    return
  }

  if (!Array.isArray(res.data.list)) {
    uni.showToast({ title: $t.value('apiResponseInvalid'), icon: 'none' })
    if (pageNo === 1) {
      serverList.value = []
    }
    serverHasMore.value = false
    return
  }

  const list = res.data.list
  if (pageNo === 1) {
    serverList.value = list
  } else {
    serverList.value = serverList.value.concat(list)
  }

  serverPage.value = pageNo
  const total = Number(res.data.total || 0)
  if (total > 0) {
    serverHasMore.value = serverList.value.length < total
  } else {
    serverHasMore.value = list.length >= PAGE_SIZE
  }
}

const resetHistoryState = () => {
  displayPage.value = 1
  serverPage.value = 1
  serverHasMore.value = false
  serverList.value = []
  mergedList.value = []
  allList.value = []
  hasMoreHistory.value = false
  reachedBottom.value = false
}

const loadHistory = async (reset = false) => {
  if (loadingMoreHistory.value) return
  loadingMoreHistory.value = true

  if (reset) {
    resetHistoryState()
  }

  const token = uni.getStorageSync('x-token') || ''
  localListCache.value = getTryonLocalHistory()
  localTaskKeys.value = localListCache.value.map(v => getTaskKey(v)).filter(Boolean)

  try {
    if (!token) {
      serverList.value = []
      serverHasMore.value = false
      syncVisibleList()
      return
    }

    await fetchServerPage(1)
    syncVisibleList()
  } finally {
    loadingMoreHistory.value = false
  }
}

const onReachBottom = () => {
  if (!hasMoreHistory.value) return
  reachedBottom.value = true
}

const loadMoreHistory = async () => {
  if (loadingMoreHistory.value || !hasMoreHistory.value) return
  loadingMoreHistory.value = true

  try {
    const token = uni.getStorageSync('x-token') || ''
    const nextDisplayPage = displayPage.value + 1
    const requiredCount = nextDisplayPage * PAGE_SIZE

    if (mergedList.value.length < requiredCount && token && serverHasMore.value) {
      await fetchServerPage(serverPage.value + 1)
    }

    displayPage.value = nextDisplayPage
    syncVisibleList()
    reachedBottom.value = false
  } finally {
    loadingMoreHistory.value = false
  }
}

const switchPreviewTab = (tab) => {
  if (tab === 'result' && !previewResultImages.value.length) return
  if (tab === 'origin' && !previewOriginImages.value.length) return
  previewTab.value = tab === 'origin' ? 'origin' : 'result'
  previewCurrentIndex.value = 0
}

const onPreviewSwiperChange = (e) => {
  previewCurrentIndex.value = Number(e?.detail?.current || 0)
}

const closePreview = () => {
  closeH5ImagePreview()
  previewNativeOpening.value = false
  previewVisible.value = false
  previewCurrentIndex.value = 0
  previewResultImages.value = []
  previewOriginImages.value = []
  closeCompare()
}

const closeCompare = () => {
  compareVisible.value = false
  dismissCompareGuide()
  comparePercent.value = 50
  compareZoomPercent.value = 100
  compareOffsetX.value = 0
  compareOffsetY.value = 0
  compareDragging.value = false
  comparePinching.value = false
  compareMode.value = 'tryon'
}

const dismissCompareGuide = () => {
  compareGuideVisible.value = false
}

const showCompareGuide = () => {
  const hasSeen = String(uni.getStorageSync(COMPARE_GUIDE_SEEN_KEY) || '') === '1'
  if (hasSeen) {
    compareGuideVisible.value = false
    return
  }
  compareGuideVisible.value = true
  uni.setStorageSync(COMPARE_GUIDE_SEEN_KEY, '1')
}

const onCompareSliderChange = (e) => {
  dismissCompareGuide()
  comparePercent.value = Math.max(0, Math.min(100, Number(e?.detail?.value ?? 50)))
}

const clampNumber = (value, min, max) => {
  const numeric = Number(value || 0)
  if (numeric < min) return min
  if (numeric > max) return max
  return numeric
}

const clampCompareOffset = (x, y) => {
  compareOffsetX.value = clampNumber(x, -comparePanRangeX.value, comparePanRangeX.value)
  compareOffsetY.value = clampNumber(y, -comparePanRangeY.value, comparePanRangeY.value)
}

const setCompareZoom = (value) => {
  compareZoomPercent.value = Math.max(100, Math.min(240, Number(value || 100)))
  if (compareZoomPercent.value <= 100) {
    compareOffsetX.value = 0
    compareOffsetY.value = 0
    return
  }
  clampCompareOffset(compareOffsetX.value, compareOffsetY.value)
}

const onCompareZoomChange = (e) => {
  dismissCompareGuide()
  setCompareZoom(e?.detail?.value)
}

const getTouchPoint = (event) => {
  const touch = event?.touches?.[0] || event?.changedTouches?.[0]
  if (!touch) return null
  return {
    x: Number(touch.clientX ?? touch.pageX ?? 0),
    y: Number(touch.clientY ?? touch.pageY ?? 0),
  }
}

const getTouchPoints = (event) => {
  const touches = event?.touches || []
  const points = []
  const length = Math.min(2, touches.length || 0)
  for (let i = 0; i < length; i += 1) {
    const touch = touches[i]
    points.push({
      x: Number(touch?.clientX ?? touch?.pageX ?? 0),
      y: Number(touch?.clientY ?? touch?.pageY ?? 0),
    })
  }
  return points
}

const getTouchDistance = (p1, p2) => {
  if (!p1 || !p2) return 0
  return Math.hypot(p1.x - p2.x, p1.y - p2.y)
}

const getTouchCenter = (p1, p2) => {
  if (!p1 || !p2) return null
  return {
    x: (p1.x + p2.x) / 2,
    y: (p1.y + p2.y) / 2,
  }
}

const startComparePinch = (event) => {
  const points = getTouchPoints(event)
  if (points.length < 2) return false

  const distance = getTouchDistance(points[0], points[1])
  if (!(distance > 0)) return false

  const center = getTouchCenter(points[0], points[1])
  if (!center) return false

  comparePinching.value = true
  compareDragging.value = false
  comparePinchStartDistance.value = distance
  comparePinchStartZoomPercent.value = compareZoomPercent.value
  comparePinchStartCenterX.value = center.x
  comparePinchStartCenterY.value = center.y
  comparePinchBaseOffsetX.value = compareOffsetX.value
  comparePinchBaseOffsetY.value = compareOffsetY.value
  return true
}

const onCompareStageTouchStart = (event) => {
  if (!compareVisible.value) return
  dismissCompareGuide()
  if (startComparePinch(event)) return
  if (comparePinching.value || compareZoomScale.value <= 1) return

  const point = getTouchPoint(event)
  if (!point) return

  compareDragging.value = true
  compareDragStartX.value = point.x
  compareDragStartY.value = point.y
  compareBaseOffsetX.value = compareOffsetX.value
  compareBaseOffsetY.value = compareOffsetY.value
}

const onCompareStageTouchMove = (event) => {
  if (!compareVisible.value) return

  const points = getTouchPoints(event)
  if (comparePinching.value && points.length >= 2) {
    const distance = getTouchDistance(points[0], points[1])
    if (!(distance > 0) || !(comparePinchStartDistance.value > 0)) return

    const center = getTouchCenter(points[0], points[1])
    const targetZoom = (comparePinchStartZoomPercent.value * distance) / comparePinchStartDistance.value
    compareZoomPercent.value = Math.max(100, Math.min(240, targetZoom))

    if (compareZoomPercent.value <= 100) {
      compareOffsetX.value = 0
      compareOffsetY.value = 0
      return
    }

    if (center) {
      const deltaX = center.x - comparePinchStartCenterX.value
      const deltaY = center.y - comparePinchStartCenterY.value
      clampCompareOffset(comparePinchBaseOffsetX.value + deltaX, comparePinchBaseOffsetY.value + deltaY)
    }
    return
  }

  if (points.length >= 2) {
    startComparePinch(event)
    return
  }

  if (!compareDragging.value || comparePinching.value || compareZoomScale.value <= 1) return

  const point = getTouchPoint(event)
  if (!point) return

  const deltaX = point.x - compareDragStartX.value
  const deltaY = point.y - compareDragStartY.value
  clampCompareOffset(compareBaseOffsetX.value + deltaX, compareBaseOffsetY.value + deltaY)
}

const onCompareStageTouchEnd = (event) => {
  const remainingTouches = Number(event?.touches?.length || 0)
  if (remainingTouches >= 2) return

  if (comparePinching.value && remainingTouches === 1) {
    comparePinching.value = false
    if (compareZoomScale.value > 1) {
      const point = getTouchPoint(event)
      if (point) {
        compareDragging.value = true
        compareDragStartX.value = point.x
        compareDragStartY.value = point.y
        compareBaseOffsetX.value = compareOffsetX.value
        compareBaseOffsetY.value = compareOffsetY.value
        return
      }
    }
  }

  if (remainingTouches === 0) {
    comparePinching.value = false
  }
  compareDragging.value = false
}

const measureCompareStage = () => {
  nextTick(() => {
    const query = uni.createSelectorQuery()
    query.select('.compare-stage').boundingClientRect((rect) => {
      compareStageWidthPx.value = Number(rect?.width || 0)
      compareStageHeightPx.value = Number(rect?.height || 0)
    }).exec()
  })
}

const openCompareFromPreview = () => {
  if (!canCompareTryonPreview.value) {
    uni.showToast({ title: $t.value('noImagePreview'), icon: 'none' })
    return
  }

  compareMode.value = 'tryon'
  comparePercent.value = 50
  compareZoomPercent.value = 100
  compareOffsetX.value = 0
  compareOffsetY.value = 0
  compareDragging.value = false
  comparePinching.value = false
  compareVisible.value = true
  showCompareGuide()
  measureCompareStage()
}

const isProcessingTask = (item) => String(item?.status || '').trim().toLowerCase() === 'processing'

const handleCardTap = (item) => {
  if (isProcessingTask(item) && Number(item?.ID || 0) > 0) {
    uni.navigateTo({ url: `/pages/tryon/generate?taskID=${Number(item.ID)}` })
    return
  }
  preview(item)
}

const currentPreviewImage = () => {
  const list = currentPreviewImages.value
  return list[previewCurrentIndex.value] || ''
}

const closeH5ImagePreview = () => {
  h5ImagePreviewVisible.value = false
  h5ImagePreviewCurrent.value = 0
  h5ImagePreviewList.value = []
}

const downloadH5PreviewImage = (url) => {
  const target = String(url || '').trim()
  if (!target) {
    uni.showToast({ title: $t.value('noImagePreview'), icon: 'none' })
    return
  }

  const anchor = document.createElement('a')
  anchor.href = target
  anchor.target = '_blank'
  anchor.rel = 'noopener'
  anchor.download = `tryon-${Date.now()}`
  document.body.appendChild(anchor)
  anchor.click()
  document.body.removeChild(anchor)
  uni.showToast({ title: $t.value('downloadStarted'), icon: 'none' })
}

const downloadH5PreviewCurrentImage = () => {
  const target = h5ImagePreviewList.value[h5ImagePreviewCurrent.value] || ''
  downloadH5PreviewImage(target)
}

const onH5ImagePreviewChange = (event) => {
  h5ImagePreviewCurrent.value = Number(event?.detail?.current || 0)
}

const previewCurrentImageByIndex = (index) => {
  const list = currentPreviewImages.value
    .map((item) => getUrl(item))
    .filter(Boolean)

  if (!list.length) {
    uni.showToast({ title: $t.value('noImagePreview'), icon: 'none' })
    return
  }

  const current = list[index] || list[previewCurrentIndex.value] || list[0]
  if (previewNativeOpening.value) return

  const isH5Runtime = typeof window !== 'undefined' && typeof document !== 'undefined'
  if (isH5Runtime) {
    h5ImagePreviewList.value = list
    h5ImagePreviewCurrent.value = Math.max(0, list.findIndex(item => item === current))
    h5ImagePreviewVisible.value = true
    return
  }

  previewNativeOpening.value = true

  uni.previewImage({
    urls: list,
    current,
    // #ifdef H5
    complete: () => {
      previewNativeOpening.value = false
    },
    // #endif
    fail: () => {
      previewNativeOpening.value = false
    },
  })
}

const ensureAlbumPermission = async () => {
  const isH5Runtime = typeof window !== 'undefined' && typeof document !== 'undefined'
  if (isH5Runtime) return true

  try {
    const settingRes = await uni.getSetting()
    const authState = settingRes?.authSetting?.['scope.writePhotosAlbum']
    if (authState === false) {
      const openRes = await uni.openSetting()
      return !!openRes?.authSetting?.['scope.writePhotosAlbum']
    }
    return true
  } catch {
    return true
  }
}

const saveImageToAlbumWithRetry = async (filePath) => {
  const canSave = await ensureAlbumPermission()
  if (!canSave) {
    throw new Error('NO_ALBUM_PERMISSION')
  }

  try {
    await uni.saveImageToPhotosAlbum({ filePath })
  } catch (e) {
    const errMsg = String(e?.errMsg || '').toLowerCase()
    if (!/auth|permission/.test(errMsg)) {
      throw e
    }

    const openRes = await uni.openSetting()
    const granted = !!openRes?.authSetting?.['scope.writePhotosAlbum']
    if (!granted) {
      throw new Error('NO_ALBUM_PERMISSION')
    }

    await uni.saveImageToPhotosAlbum({ filePath })
  }
}

const downloadCurrentImage = async () => {
  const target = currentPreviewImage()
  if (!target) {
    uni.showToast({ title: $t.value('noImagePreview'), icon: 'none' })
    return
  }

  const fullUrl = getUrl(target)
  if (!fullUrl) {
    uni.showToast({ title: $t.value('imageUrlInvalid'), icon: 'none' })
    return
  }

  const isH5Runtime = typeof window !== 'undefined' && typeof document !== 'undefined'
  if (isH5Runtime) {
    const anchor = document.createElement('a')
    anchor.href = fullUrl
    anchor.target = '_blank'
    anchor.rel = 'noopener'
    anchor.download = `tryon-${Date.now()}`
    document.body.appendChild(anchor)
    anchor.click()
    document.body.removeChild(anchor)
    uni.showToast({ title: $t.value('downloadStarted'), icon: 'none' })
    return
  }

  uni.showLoading({ title: $t.value('loading'), mask: true })
  try {
    const res = await uni.downloadFile({ url: fullUrl })
    if (res.statusCode !== 200 || !res.tempFilePath) {
      uni.showToast({ title: $t.value('downloadFailed'), icon: 'none' })
      return
    }

    await saveImageToAlbumWithRetry(res.tempFilePath)
    uni.showToast({ title: $t.value('savedToAlbum'), icon: 'none' })
  } catch (e) {
    if (String(e?.message || '') === 'NO_ALBUM_PERMISSION') {
      uni.showToast({ title: $t.value('saveToAlbumFailed'), icon: 'none' })
      return
    }
    uni.showToast({ title: $t.value('downloadFailed'), icon: 'none' })
  } finally {
    uni.hideLoading()
  }
}

const removeHistory = (item) => {
  const taskKey = getTaskKey(item)
  if (!taskKey) {
    uni.showToast({ title: $t.value('historyRecordNotFound'), icon: 'none' })
    return
  }

  uni.showModal({
    title: $t.value('deleteHistoryConfirmTitle'),
    content: $t.value('deleteHistoryConfirmContent'),
    cancelText: $t.value('cancel'),
    confirmText: $t.value('confirm'),
    success: async (res) => {
      if (!res.confirm) return

      const token = uni.getStorageSync('x-token') || ''
      if (token && item.taskNo) {
        const deleteRes = await deleteMyTryonTask({ taskNo: item.taskNo })
        if (deleteRes.code !== 0) {
          uni.showToast({ title: resolveApiMessage(deleteRes.msg, 'operationFailed'), icon: 'none' })
          return
        }
      }

      removeTryonLocalHistoryByTaskKey(taskKey)
      await loadHistory(true)
      uni.showToast({ title: $t.value('deleteSuccess'), icon: 'none' })
    },
  })
}

const preview = (item) => {
  const resultImages = uniqueImageList(extractImageValues(item.resultImage))
  const originImages = uniqueImageList([
    ...extractImageValues(item.sourceImage),
    ...extractImageValues(item.templateImage),
  ])

  if (!resultImages.length && !originImages.length) {
    uni.showToast({ title: $t.value('noImagePreview'), icon: 'none' })
    return
  }

  previewResultImages.value = resultImages
  previewOriginImages.value = originImages
  previewTab.value = resultImages.length ? 'result' : 'origin'
  previewCurrentIndex.value = 0
  compareMode.value = 'tryon'
  previewNativeOpening.value = false
  previewVisible.value = true
}

const goBack = () => {
  uni.navigateBack({ delta: 1 })
}

onShow(() => {
  previewNativeOpening.value = false
  loadHistory(true)
})
</script>

<style lang="scss">
page {
  background: #f4f7fb;
}

.history-page {
  min-height: 100vh;
  padding: calc(var(--status-bar-height, 0px) + 14rpx) 20rpx 20rpx;
  color: #0f172a;
  background: radial-gradient(120% 80% at 100% -10%, #dbeafe 0%, transparent 60%), #f4f7fb;
}

.nav {
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-btn {
  width: 90rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-title {
  font-size: 30rpx;
  font-weight: 700;
}

.filter-row {
  margin-top: 12rpx;
  display: flex;
  gap: 10rpx;
  overflow-x: auto;
  white-space: nowrap;
}

.filter-chip {
  min-width: 110rpx;
  height: 58rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(15,23,42,0.12);
  background: rgba(255,255,255,0.9);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 22rpx;
  color: rgba(15,23,42,0.62);
}

.filter-chip.active {
  color: #0f172a;
  border-color: rgba(37, 99, 235, 0.45);
  background: rgba(219, 234, 254, 0.82);
}

.list-wrap {
  margin-top: 14rpx;
  height: calc(100vh - var(--status-bar-height, 0px) - 170rpx);
}

.card {
  padding: 12rpx;
  border-radius: 14rpx;
  border: 1rpx solid rgba(15,23,42,0.08);
  background: rgba(255,255,255,0.95);
  box-shadow: 0 12rpx 26rpx rgba(15, 23, 42, 0.06);
  margin-bottom: 12rpx;
  display: flex;
  gap: 12rpx;
}

.thumb-wrap {
  width: 156rpx;
  height: 156rpx;
  min-width: 156rpx;
  border-radius: 10rpx;
  overflow: hidden;
  background: #f8fafc;
  display: flex;
  align-items: center;
  justify-content: center;
}

.thumb {
  width: 100%;
  height: 100%;
}

.card-main {
  flex: 1;
  min-width: 0;
}

.line {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8rpx;
}

.type {
  font-size: 24rpx;
  font-weight: 600;
}

.status {
  padding: 4rpx 12rpx;
  border-radius: 999rpx;
  font-size: 20rpx;
}

.status.processing {
  background: rgba(255,205,105,0.2);
}

.status.success {
  background: rgba(86,199,102,0.22);
}

.status.failed {
  background: rgba(255,107,107,0.22);
}

.no,
.time,
.error {
  margin-top: 7rpx;
  display: block;
  font-size: 21rpx;
  color: rgba(15,23,42,0.62);
}

.error {
  color: #ffaaaa;
}

.card-actions {
  margin-top: 8rpx;
  display: flex;
  justify-content: flex-end;
  gap: 10rpx;
}

.mini-btn {
  height: 42rpx;
  padding: 0 14rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(37, 99, 235, 0.22);
  background: rgba(37, 99, 235, 0.1);
  color: #1d4ed8;
  font-size: 20rpx;
  display: inline-flex;
  align-items: center;
}

.mini-btn.danger {
  border-color: rgba(220, 38, 38, 0.25);
  background: rgba(220, 38, 38, 0.1);
  color: #b91c1c;
}

.empty {
  margin-top: 120rpx;
  text-align: center;
  color: rgba(15,23,42,0.55);
}

.load-more-row {
  margin-top: 8rpx;
  margin-bottom: 8rpx;
  display: flex;
  justify-content: center;
}

.load-more-btn {
  width: 260rpx;
  height: 64rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.12);
  background: rgba(255, 255, 255, 0.96);
  color: #0f172a;
  font-size: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.load-more-btn.disabled {
  opacity: 0.65;
}

.preview-mask {
  position: fixed;
  inset: 0;
  z-index: 1200;
  background: rgba(15, 23, 42, 0.58);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24rpx;
}

.preview-mask-passthrough {
  pointer-events: none;
  background: transparent;
}

.preview-mask-passthrough .preview-panel {
  opacity: 0;
}

.h5-image-preview-mask {
  position: fixed;
  inset: 0;
  z-index: 2600;
  background: rgba(0, 0, 0, 0.92);
}

.h5-image-preview-head {
  position: absolute;
  left: 0;
  right: 0;
  top: calc(var(--status-bar-height, 0px) + 12rpx);
  z-index: 2;
  height: 76rpx;
  padding: 0 20rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.h5-image-preview-actions {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.h5-image-preview-save {
  height: 48rpx;
  padding: 0 18rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.42);
  background: rgba(255, 255, 255, 0.2);
  color: #ffffff;
  font-size: 20rpx;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.h5-image-preview-index {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.86);
}

.h5-image-preview-close {
  width: 56rpx;
  height: 56rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.18);
}

.h5-image-preview-swiper {
  width: 100%;
  height: 100%;
}

.h5-image-preview-wrap {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24rpx;
  box-sizing: border-box;
}

.h5-image-preview-img {
  width: 100%;
  height: 100%;
}

.h5-image-preview-tip {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 36rpx;
  text-align: center;
  color: rgba(255, 255, 255, 0.8);
  font-size: 22rpx;
}

.preview-panel {
  width: 100%;
  max-width: 680rpx;
  border-radius: 18rpx;
  background: #ffffff;
  overflow: hidden;
}

.preview-head {
  height: 88rpx;
  padding: 0 20rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1rpx solid rgba(15, 23, 42, 0.08);
}

.preview-title {
  font-size: 28rpx;
  font-weight: 700;
  color: #0f172a;
}

.preview-close {
  width: 54rpx;
  height: 54rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(15, 23, 42, 0.06);
}

.preview-tabs {
  padding: 14rpx 16rpx;
  display: flex;
  gap: 12rpx;
}

.preview-tab {
  flex: 1;
  height: 56rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.12);
  color: rgba(15, 23, 42, 0.64);
  background: rgba(15, 23, 42, 0.04);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22rpx;
}

.preview-tab.active {
  color: #ffffff;
  border-color: transparent;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
}

.preview-tab.disabled {
  opacity: 0.45;
}

.preview-body {
  padding: 0 16rpx 10rpx;
}

.preview-swiper {
  width: 100%;
  height: 680rpx;
}

.preview-image-wrap {
  width: 100%;
  height: 680rpx;
  border-radius: 12rpx;
  overflow: hidden;
  background: #f8fafc;
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview-image {
  width: 100%;
  height: 100%;
}

.preview-index {
  margin-top: 8rpx;
  text-align: center;
  display: block;
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.58);
}

.preview-empty {
  padding: 90rpx 16rpx;
  text-align: center;
  color: rgba(15, 23, 42, 0.52);
}

.preview-actions {
  padding: 10rpx 16rpx 18rpx;
  display: flex;
  flex-direction: column;
  gap: 10rpx;
}

.preview-actions-row {
  display: flex;
  gap: 12rpx;
}

.preview-save-tip {
  display: block;
  text-align: center;
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.56);
}

.preview-action-btn {
  flex: 1;
  height: 66rpx;
  border-radius: 999rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ffffff;
  font-size: 24rpx;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
}

.preview-action-btn.secondary {
  color: #0f172a;
  background: rgba(15, 23, 42, 0.06);
  border: 1rpx solid rgba(15, 23, 42, 0.12);
}

.compare-panel {
  width: calc(100vw - 32rpx);
  max-width: 980rpx;
  max-height: calc(100vh - 36rpx);
  display: flex;
  flex-direction: column;
}

.compare-head--compact {
  justify-content: flex-end;
}

.compare-stage {
  margin: 12rpx 16rpx 0;
  height: 56vh;
  min-height: 420rpx;
  max-height: none;
  border-radius: 12rpx;
  overflow: hidden;
  position: relative;
  background: #f8fafc;
}

.compare-image {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  will-change: transform;
}

.compare-result-layer {
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  overflow: hidden;
  z-index: 2;
  pointer-events: none;
}

.compare-result-image {
  right: auto;
  bottom: auto;
}

.compare-divider {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 4rpx;
  margin-left: -2rpx;
  background: #ffffff;
  box-shadow: 0 0 0 1rpx rgba(37, 99, 235, 0.35);
  z-index: 3;
  pointer-events: none;
}

.compare-slider-wrap {
  position: relative;
  padding: 16rpx 20rpx 20rpx;
}

.compare-guide-stage-overlay {
  position: absolute;
  inset: 0;
  z-index: 8;
  pointer-events: none;
}

.compare-guide-bubble {
  max-width: 86%;
  border-radius: 14rpx;
  padding: 12rpx 16rpx;
  color: #ffffff;
  font-size: 22rpx;
  line-height: 32rpx;
  box-shadow: 0 10rpx 28rpx rgba(37, 99, 235, 0.3);
}

.compare-guide-title {
  display: block;
  margin-bottom: 4rpx;
  font-size: 20rpx;
  opacity: 0.9;
}

.compare-guide-bubble--stage {
  position: absolute;
  top: 14rpx;
  left: 14rpx;
  background: linear-gradient(135deg, rgba(37, 99, 235, 0.96), rgba(14, 165, 233, 0.94));
}

.compare-guide-pointer {
  width: 0;
  height: 0;
  position: absolute;
}

.compare-guide-pointer--stage {
  top: 100rpx;
  left: 56rpx;
  border-left: 10rpx solid transparent;
  border-right: 10rpx solid transparent;
  border-top: 14rpx solid rgba(37, 99, 235, 0.94);
}

.compare-guide-divider-hint {
  position: absolute;
  top: 52%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
}

.compare-guide-hand {
  width: 50rpx;
  height: 30rpx;
  border-radius: 999rpx;
  border: 2rpx solid rgba(255, 255, 255, 0.92);
  background: rgba(15, 23, 42, 0.24);
  display: flex;
  align-items: center;
  justify-content: center;
  animation: compareGuideHandMove 0.9s ease-in-out infinite alternate;
}

.compare-guide-hand-dot {
  width: 10rpx;
  height: 10rpx;
  border-radius: 50%;
  background: #ffffff;
}

.compare-guide-divider-text {
  color: #ffffff;
  font-size: 20rpx;
  text-shadow: 0 2rpx 8rpx rgba(15, 23, 42, 0.5);
  white-space: nowrap;
}

.compare-guide-slider-tip {
  margin-bottom: 8rpx;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  pointer-events: none;
}

.compare-guide-bubble--slider {
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.92), rgba(30, 41, 59, 0.94));
}

.compare-guide-pointer--slider {
  position: static;
  margin-left: 26rpx;
  border-left: 10rpx solid transparent;
  border-right: 10rpx solid transparent;
  border-top: 14rpx solid rgba(30, 41, 59, 0.94);
}

.compare-guide-action-row {
  margin-top: 12rpx;
  display: flex;
  justify-content: flex-end;
}

.compare-guide-action-btn {
  min-width: 132rpx;
  height: 46rpx;
  border-radius: 999rpx;
  padding: 0 20rpx;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  color: #ffffff;
  font-size: 20rpx;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 18rpx rgba(37, 99, 235, 0.28);
}

.compare-slider-label {
  display: flex;
  justify-content: space-between;
  color: rgba(15, 23, 42, 0.6);
  font-size: 20rpx;
}

.compare-slider {
  margin-top: 6rpx;
}

.compare-slider-label--zoom {
  margin-top: 12rpx;
  align-items: center;
}

.compare-zoom-max {
  min-width: 72rpx;
  height: 42rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.16);
  background: rgba(255, 255, 255, 0.95);
  color: #0f172a;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 20rpx;
}

.compare-zoom-slider {
  margin-top: 4rpx;
}

@keyframes compareGuideHandMove {
  0% {
    transform: translateX(-10rpx);
  }
  100% {
    transform: translateX(10rpx);
  }
}
</style>
