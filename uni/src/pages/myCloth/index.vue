<template>
  <view class="cloth-page">
    <view class="nav">
      <view class="nav-btn" @tap="goBack">
        <uni-icons type="left" size="20" color="#0f172a" />
      </view>
      <text class="nav-title">{{ $t('myCloset') }}</text>
      <view class="nav-btn" @tap="addCloth">
        <uni-icons type="plusempty" size="22" color="#0f172a" />
      </view>
    </view>

    <view class="tips">{{ $t('myClosetTips') }}</view>

    <view class="category-tabs">
      <view
        v-for="item in categoryOptions"
        :key="item.key"
        class="category-tab"
        :class="{ active: activeCategory === item.key }"
        @tap="activeCategory = item.key"
      >
        {{ $t(item.labelKey) }}
      </view>
    </view>

    <scroll-view class="list-wrap" scroll-y>
      <view class="grid">
        <view class="card" v-for="item in filteredClothList" :key="item.id">
          <view class="card-image-wrap">
            <image class="card-image" :src="getUrl(item.url)" mode="aspectFit" @tap="preview(item)" />
            <view v-if="item.sizeText" class="card-size-badge">{{ item.sizeText }}</view>
          </view>
          <view class="card-foot">
            <view class="name-row">
              <text class="card-name">{{ item.name || categoryLabel(item.category) }}</text>
              <text class="category-badge">{{ categoryLabel(item.category) }}</text>
            </view>
            <view class="card-actions card-actions-secondary">
              <view class="mini-btn" @tap.stop="renameCloth(item)">{{ $t('rename') }}</view>
              <view class="mini-btn danger" @tap.stop="removeCloth(item)">{{ $t('delete') }}</view>
            </view>
          </view>
        </view>
      </view>

      <view class="empty" v-if="filteredClothList.length === 0">
        <text>{{ $t('myClosetEmpty') }}</text>
      </view>
      <view style="height: 40rpx"></view>
    </scroll-view>

    <view class="crop-editor-mask" v-if="cropEditorVisible" @tap.stop>
      <view class="crop-editor-panel" @tap.stop>
        <view class="crop-editor-head">
          <text class="crop-editor-title">{{ $t('uploadModeCrop') }}</text>
          <text class="crop-editor-ratio">{{ cropRatioLabel }}</text>
        </view>
        <view class="crop-ratio-options">
          <view
            v-for="ratio in cropRatioOptions"
            :key="ratio.key"
            class="crop-ratio-chip"
            :class="{ active: isCropRatioActive(ratio) }"
            @tap="changeCropRatio(ratio)"
          >
            {{ ratio.label }}
          </view>
        </view>
        <view class="crop-editor-stage-wrap">
          <view class="crop-editor-stage" :style="{ width: `${cropStageSize.width}px`, height: `${cropStageSize.height}px` }">
            <image class="crop-editor-image" :src="cropSourcePath" mode="scaleToFill" />
            <movable-area class="crop-editor-area" :style="{ width: `${cropStageSize.width}px`, height: `${cropStageSize.height}px` }">
              <movable-view
                class="crop-editor-box"
                direction="all"
                :x="cropBoxPosition.x"
                :y="cropBoxPosition.y"
                :style="{ width: `${cropBoxSize.width}px`, height: `${cropBoxSize.height}px` }"
                @change="onCropBoxChange"
              >
                <view class="crop-editor-box-inner"></view>
              </movable-view>
            </movable-area>
          </view>
        </view>
        <view class="crop-editor-edge-tools">
          <text class="crop-editor-edge-tip">{{ $t('cropEdgeAdjustTip') }}</text>
          <view class="crop-editor-edge-list">
            <view class="crop-editor-edge-chip" :class="{ active: activeCropEdge === 'all' }" @tap="activeCropEdge = 'all'">{{ $t('cropEdgeAll') }}</view>
            <view class="crop-editor-edge-chip" :class="{ active: activeCropEdge === 'left' }" @tap="activeCropEdge = 'left'">{{ $t('cropEdgeLeft') }}</view>
            <view class="crop-editor-edge-chip" :class="{ active: activeCropEdge === 'right' }" @tap="activeCropEdge = 'right'">{{ $t('cropEdgeRight') }}</view>
            <view class="crop-editor-edge-chip" :class="{ active: activeCropEdge === 'top' }" @tap="activeCropEdge = 'top'">{{ $t('cropEdgeTop') }}</view>
            <view class="crop-editor-edge-chip" :class="{ active: activeCropEdge === 'bottom' }" @tap="activeCropEdge = 'bottom'">{{ $t('cropEdgeBottom') }}</view>
          </view>
        </view>
        <view class="crop-editor-zoom">
          <view class="crop-editor-zoom-btn" @tap="zoomOutCropBox">-</view>
          <view class="crop-editor-zoom-btn" @tap="zoomInCropBox">+</view>
        </view>
        <view class="crop-editor-actions">
          <view class="crop-editor-btn ghost" @tap="cancelCropEditor">{{ $t('cancel') }}</view>
          <view class="crop-editor-btn" :class="{ disabled: cropConfirming }" @tap="confirmCropEditor">{{ cropConfirming ? $t('loading') : $t('doneText') }}</view>
        </view>
      </view>
    </view>

    <view class="mode-popup-mask" v-if="modePopupVisible" @tap="closeModePopup('')">
      <view class="mode-popup-panel" @tap.stop>
        <text v-if="modePopupTitle" class="mode-popup-title">{{ modePopupTitle }}</text>
        <view
          v-for="item in modePopupOptions"
          :key="item.key"
          class="mode-popup-item"
          @tap="closeModePopup(item.key)"
        >
          {{ item.label }}
        </view>
        <view class="mode-popup-cancel" @tap="closeModePopup('')">{{ $t('cancel') }}</view>
      </view>
    </view>

    <canvas
      canvas-id="myClothCropCanvas"
      id="myClothCropCanvas"
      class="crop-canvas-hidden"
      :style="{ width: `${cropCanvasSize.width}px`, height: `${cropCanvasSize.height}px` }"
    ></canvas>
  </view>
</template>

<script setup>
import { ref, computed, nextTick } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { resolveApiMessage } from '@/utils/i18n.js'
import { getUrl } from '@/utils/url.js'
import { uploadTryonImage } from '@/utils/tryon.js'
import {
  createTryonCloth,
  updateTryonCloth,
  deleteTryonCloth,
  getMyTryonClothList,
} from '@/api/tryonTask.js'
import { useLangStore } from '@/pinia/modules/lang.js'

const CATEGORY_KEYS = ['upper', 'lower', 'onepiece']

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const clothList = ref([])
const activeCategory = ref('all')
const modePopupVisible = ref(false)
const modePopupTitle = ref('')
const modePopupOptions = ref([])
const cropCanvasSize = ref({ width: 1, height: 1 })
const cropEditorVisible = ref(false)
const cropSourcePath = ref('')
const cropSourceSize = ref({ width: 0, height: 0 })
const cropStageSize = ref({ width: 1, height: 1 })
const cropBoxSize = ref({ width: 1, height: 1 })
const cropBoxPosition = ref({ x: 0, y: 0 })
const cropRatio = ref({ width: 3, height: 4 })
const activeCropEdge = ref('all')
const cropConfirming = ref(false)
let cropResolve = null

const CROP_CANVAS_ID = 'myClothCropCanvas'
const CROP_MIN_EDGE_SIZE = 48
const cropRatioLabel = computed(() => `${cropRatio.value.width}:${cropRatio.value.height}`)
const cropRatioOptions = [
  { key: '1:1', label: '1:1', width: 1, height: 1 },
  { key: '3:4', label: '3:4', width: 3, height: 4 },
  { key: '4:3', label: '4:3', width: 4, height: 3 },
  { key: '9:16', label: '9:16', width: 9, height: 16 },
  { key: '16:9', label: '16:9', width: 16, height: 9 },
]

const categoryOptions = [
  { key: 'all', labelKey: 'clothCategoryAll' },
  { key: 'upper', labelKey: 'clothCategoryUpper' },
  { key: 'lower', labelKey: 'clothCategoryLower' },
  { key: 'onepiece', labelKey: 'clothCategoryOnepiece' },
]

const CATEGORY_LABEL_KEY_MAP = {
  upper: 'clothCategoryUpper',
  lower: 'clothCategoryLower',
  onepiece: 'clothCategoryOnepiece',
}

const waitFrame = (delay = 30) => new Promise((resolve) => setTimeout(resolve, delay))


let modePopupResolve = null

const openModePopup = ({ title = '', options = [] } = {}) => {
  return new Promise((resolve) => {
    modePopupTitle.value = String(title || '').trim()
    modePopupOptions.value = Array.isArray(options) ? options.filter((item) => item && item.key && item.label) : []
    modePopupVisible.value = true
    modePopupResolve = resolve
  })
}

const closeModePopup = (selectedKey = '') => {
  modePopupVisible.value = false
  const resolver = modePopupResolve
  modePopupResolve = null
  if (typeof resolver === 'function') {
    resolver(String(selectedKey || ''))
  }
}

const chooseImageAsync = (sizeType = ['original']) => {
  return new Promise((resolve, reject) => {
    uni.chooseImage({
      count: 1,
      sizeType,
      sourceType: ['album', 'camera'],
      success: resolve,
      fail: reject,
    })
  })
}

const getImageInfoAsync = (src) => {
  return new Promise((resolve, reject) => {
    uni.getImageInfo({
      src,
      success: resolve,
      fail: reject,
    })
  })
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

const sizeTextCache = Object.create(null)

const formatBadgeSize = (bytes) => {
  const size = Number(bytes || 0)
  if (!Number.isFinite(size) || size <= 0) return ''

  const kb = size / 1024
  if (kb < 1024) {
    const value = kb >= 10 ? kb.toFixed(0) : kb.toFixed(1)
    return `${value}KB`
  }

  const mb = kb / 1024
  const value = mb >= 10 ? mb.toFixed(0) : mb.toFixed(1)
  return `${value}MB`
}

const resolveRemoteFileSizeText = async (rawUrl) => {
  const normalizedUrl = String(getUrl(rawUrl) || '').trim()
  if (!normalizedUrl) return ''

  if (Object.prototype.hasOwnProperty.call(sizeTextCache, normalizedUrl)) {
    return sizeTextCache[normalizedUrl]
  }

  try {
    const downloadRes = await downloadFileAsync(normalizedUrl)
    const tempPath = String(downloadRes?.tempFilePath || '').trim()
    if (!tempPath) {
      sizeTextCache[normalizedUrl] = ''
      return ''
    }
    const bytes = await getFileSizeAsync(tempPath)
    const text = formatBadgeSize(bytes)
    sizeTextCache[normalizedUrl] = text
    return text
  } catch (e) {
    sizeTextCache[normalizedUrl] = ''
    return ''
  }
}

const fillClothItemSizeText = (item) => {
  if (!item?.url || item.sizeText) return
  const currentUrl = String(item.url || '').trim()
  resolveRemoteFileSizeText(currentUrl).then((text) => {
    if (!text) return
    if (String(item.url || '').trim() !== currentUrl) return
    item.sizeText = text
  })
}

const drawCanvasAsync = (ctx) => {
  return new Promise((resolve) => {
    ctx.draw(false, resolve)
  })
}

const canvasToTempFilePathAsync = (options) => {
  return new Promise((resolve, reject) => {
    uni.canvasToTempFilePath({
      ...options,
      success: resolve,
      fail: reject,
    })
  })
}

const clamp = (value, min, max) => {
  if (!Number.isFinite(value)) return min
  return Math.min(Math.max(value, min), max)
}

const isCropRatioActive = (ratio) => {
  const left = Number(cropRatio.value.width || 0) * Number(ratio.height || 0)
  const right = Number(cropRatio.value.height || 0) * Number(ratio.width || 0)
  return Math.abs(left - right) < 0.001
}

const applyCropBoxRect = (nextWidth, nextHeight, nextX = cropBoxPosition.value.x, nextY = cropBoxPosition.value.y) => {
  const stageWidth = Number(cropStageSize.value.width || 0)
  const stageHeight = Number(cropStageSize.value.height || 0)
  if (stageWidth <= 0 || stageHeight <= 0) return

  const minWidth = Math.min(stageWidth, CROP_MIN_EDGE_SIZE)
  const minHeight = Math.min(stageHeight, CROP_MIN_EDGE_SIZE)
  const width = clamp(Math.round(nextWidth), minWidth, stageWidth)
  const height = clamp(Math.round(nextHeight), minHeight, stageHeight)
  const maxX = Math.max(0, stageWidth - width)
  const maxY = Math.max(0, stageHeight - height)

  cropBoxSize.value = { width, height }
  cropBoxPosition.value = {
    x: clamp(Math.round(nextX), 0, maxX),
    y: clamp(Math.round(nextY), 0, maxY),
  }
}

const changeCropRatio = (ratio) => {
  if (!ratio) return
  const stageWidth = Number(cropStageSize.value.width || 0)
  const stageHeight = Number(cropStageSize.value.height || 0)
  if (stageWidth <= 0 || stageHeight <= 0) return

  const targetRatio = Number(ratio.width || 1) / Number(ratio.height || 1)
  const currentWidth = Number(cropBoxSize.value.width || 0)
  const currentHeight = Number(cropBoxSize.value.height || 0)
  const centerX = Number(cropBoxPosition.value.x || 0) + currentWidth / 2
  const centerY = Number(cropBoxPosition.value.y || 0) + currentHeight / 2

  const area = Math.max(1, currentWidth * currentHeight)
  let nextWidth = Math.sqrt(area * targetRatio)
  let nextHeight = nextWidth / targetRatio

  const maxWidth = Math.min(stageWidth, stageHeight * targetRatio)
  if (nextWidth > maxWidth) {
    nextWidth = maxWidth
    nextHeight = nextWidth / targetRatio
  }

  const minWidth = Math.min(stageWidth, Math.max(72, Math.round(stageWidth * 0.2), Math.round(stageHeight * 0.2 * targetRatio)))
  if (nextWidth < minWidth) {
    nextWidth = minWidth
    nextHeight = nextWidth / targetRatio
  }

  cropRatio.value = { width: Number(ratio.width), height: Number(ratio.height) }
  applyCropBoxRect(nextWidth, nextHeight, Math.round(centerX - nextWidth / 2), Math.round(centerY - nextHeight / 2))
}

const buildCropStageLayout = (imageWidth, imageHeight, ratio) => {
  const systemInfo = uni.getSystemInfoSync()
  const maxStageWidth = Math.max(220, Math.min(systemInfo.windowWidth - 40, 420))
  const maxStageHeight = Math.max(260, Math.min(Math.round(systemInfo.windowHeight * 0.58), 620))
  const scale = Math.min(maxStageWidth / imageWidth, maxStageHeight / imageHeight, 1)

  const stageWidth = Math.max(120, Math.round(imageWidth * scale))
  const stageHeight = Math.max(120, Math.round(imageHeight * scale))
  const targetRatio = ratio.width / ratio.height

  let boxWidth = Math.round(stageWidth * 0.72)
  let boxHeight = Math.round(boxWidth / targetRatio)

  if (boxHeight > Math.round(stageHeight * 0.82)) {
    boxHeight = Math.round(stageHeight * 0.82)
    boxWidth = Math.round(boxHeight * targetRatio)
  }

  if (boxWidth > stageWidth) {
    boxWidth = stageWidth
    boxHeight = Math.round(boxWidth / targetRatio)
  }

  if (boxHeight > stageHeight) {
    boxHeight = stageHeight
    boxWidth = Math.round(boxHeight * targetRatio)
  }

  const boxX = Math.round((stageWidth - boxWidth) / 2)
  const boxY = Math.round((stageHeight - boxHeight) / 2)

  return {
    stageWidth,
    stageHeight,
    boxWidth: Math.max(1, boxWidth),
    boxHeight: Math.max(1, boxHeight),
    boxX: Math.max(0, boxX),
    boxY: Math.max(0, boxY),
  }
}

const finishCropEditor = (croppedPath = '') => {
  cropEditorVisible.value = false
  const resolver = cropResolve
  cropResolve = null
  if (typeof resolver === 'function') {
    resolver(croppedPath)
  }
}

const cancelCropEditor = () => {
  finishCropEditor('')
}

const onCropBoxChange = (event) => {
  const maxX = Math.max(0, cropStageSize.value.width - cropBoxSize.value.width)
  const maxY = Math.max(0, cropStageSize.value.height - cropBoxSize.value.height)
  const x = clamp(Number(event?.detail?.x || 0), 0, maxX)
  const y = clamp(Number(event?.detail?.y || 0), 0, maxY)
  cropBoxPosition.value = {
    x: Math.round(x),
    y: Math.round(y),
  }
}

const adjustCropEdge = (edge, direction = 1) => {
  const stageWidth = Number(cropStageSize.value.width || 0)
  const stageHeight = Number(cropStageSize.value.height || 0)
  if (stageWidth <= 0 || stageHeight <= 0) return

  const stepX = Math.max(8, Math.round(Number(cropBoxSize.value.width || 0) * 0.08))
  const stepY = Math.max(8, Math.round(Number(cropBoxSize.value.height || 0) * 0.08))
  const deltaX = stepX * (direction >= 0 ? 1 : -1)
  const deltaY = stepY * (direction >= 0 ? 1 : -1)

  let nextX = Number(cropBoxPosition.value.x || 0)
  let nextY = Number(cropBoxPosition.value.y || 0)
  let nextWidth = Number(cropBoxSize.value.width || 0)
  let nextHeight = Number(cropBoxSize.value.height || 0)

  if (edge === 'left') {
    nextX -= deltaX
    nextWidth += deltaX
  } else if (edge === 'right') {
    nextWidth += deltaX
  } else if (edge === 'top') {
    nextY -= deltaY
    nextHeight += deltaY
  } else if (edge === 'bottom') {
    nextHeight += deltaY
  }

  applyCropBoxRect(nextWidth, nextHeight, nextX, nextY)
}

const scaleCropBox = (deltaScale = 0) => {
  const stageWidth = Number(cropStageSize.value.width || 0)
  const stageHeight = Number(cropStageSize.value.height || 0)
  if (stageWidth <= 0 || stageHeight <= 0) return

  if (activeCropEdge.value !== 'all') {
    adjustCropEdge(activeCropEdge.value, deltaScale >= 0 ? 1 : -1)
    return
  }

  const ratioValue = Number(cropRatio.value.width || 1) / Number(cropRatio.value.height || 1)
  const minWidth = Math.min(stageWidth, Math.max(72, Math.round(stageWidth * 0.28), Math.round(stageHeight * 0.28 * ratioValue)))
  const maxWidth = Math.max(minWidth, Math.min(Math.round(stageWidth * 0.96), Math.round(stageHeight * 0.96 * ratioValue)))

  const currentWidth = Number(cropBoxSize.value.width || minWidth)
  const centerX = Number(cropBoxPosition.value.x || 0) + currentWidth / 2
  const centerY = Number(cropBoxPosition.value.y || 0) + Number(cropBoxSize.value.height || 0) / 2

  const targetWidth = clamp(Math.round(currentWidth * (1 + deltaScale)), minWidth, maxWidth)
  const targetHeight = Math.max(1, Math.round(targetWidth / ratioValue))
  const targetX = Math.round(centerX - targetWidth / 2)
  const targetY = Math.round(centerY - targetHeight / 2)

  applyCropBoxRect(targetWidth, targetHeight, targetX, targetY)
}

const zoomInCropBox = () => {
  scaleCropBox(0.12)
}

const zoomOutCropBox = () => {
  scaleCropBox(-0.12)
}

const openCropEditor = async (filePath, ratio = { width: 3, height: 4 }) => {
  const info = await getImageInfoAsync(filePath)
  const imageWidth = Number(info.width || 0)
  const imageHeight = Number(info.height || 0)
  if (imageWidth <= 0 || imageHeight <= 0) {
    return false
  }

  const safeRatio = Number(ratio?.width || 0) > 0 && Number(ratio?.height || 0) > 0
    ? ratio
    : { width: 3, height: 4 }

  const layout = buildCropStageLayout(imageWidth, imageHeight, safeRatio)
  cropSourcePath.value = filePath
  cropSourceSize.value = {
    width: imageWidth,
    height: imageHeight,
  }
  cropRatio.value = {
    width: safeRatio.width,
    height: safeRatio.height,
  }
  cropStageSize.value = {
    width: layout.stageWidth,
    height: layout.stageHeight,
  }
  cropBoxSize.value = {
    width: layout.boxWidth,
    height: layout.boxHeight,
  }
  cropBoxPosition.value = {
    x: layout.boxX,
    y: layout.boxY,
  }
  activeCropEdge.value = 'all'
  cropEditorVisible.value = true
  return true
}

const requestCropImage = async (filePath, ratio = { width: 3, height: 4 }) => {
  if (typeof cropResolve === 'function') {
    cropResolve('')
    cropResolve = null
  }

  const opened = await openCropEditor(filePath, ratio)
  if (!opened) return ''

  return new Promise((resolve) => {
    cropResolve = resolve
  })
}

const renderCanvasFromSource = async ({ filePath, srcX = 0, srcY = 0, srcWidth, srcHeight, destWidth, destHeight, quality = 0.9 }) => {
  cropCanvasSize.value = {
    width: destWidth,
    height: destHeight,
  }
  await nextTick()
  await waitFrame()

  const ctx = uni.createCanvasContext(CROP_CANVAS_ID)
  ctx.clearRect(0, 0, destWidth, destHeight)
  ctx.drawImage(filePath, srcX, srcY, srcWidth, srcHeight, 0, 0, destWidth, destHeight)
  await drawCanvasAsync(ctx)

  const temp = await canvasToTempFilePathAsync({
    canvasId: CROP_CANVAS_ID,
    x: 0,
    y: 0,
    width: destWidth,
    height: destHeight,
    destWidth,
    destHeight,
    fileType: 'jpg',
    quality,
  })
  return temp.tempFilePath || ''
}

const confirmCropEditor = async () => {
  if (cropConfirming.value) return
  cropConfirming.value = true

  try {
    const srcWidth = Number(cropSourceSize.value.width || 0)
    const srcHeight = Number(cropSourceSize.value.height || 0)
    const stageWidth = Number(cropStageSize.value.width || 0)
    const stageHeight = Number(cropStageSize.value.height || 0)
    if (!cropSourcePath.value || srcWidth <= 0 || srcHeight <= 0 || stageWidth <= 0 || stageHeight <= 0) {
      uni.showToast({ title: $t.value('uploadFail'), icon: 'none' })
      return
    }

    const scaleX = srcWidth / stageWidth
    const scaleY = srcHeight / stageHeight
    const maxX = Math.max(0, stageWidth - cropBoxSize.value.width)
    const maxY = Math.max(0, stageHeight - cropBoxSize.value.height)

    const safeX = clamp(cropBoxPosition.value.x, 0, maxX)
    const safeY = clamp(cropBoxPosition.value.y, 0, maxY)

    let cropX = Math.round(safeX * scaleX)
    let cropY = Math.round(safeY * scaleY)
    let cropWidth = Math.round(cropBoxSize.value.width * scaleX)
    let cropHeight = Math.round(cropBoxSize.value.height * scaleY)

    cropX = clamp(cropX, 0, Math.max(0, srcWidth - 1))
    cropY = clamp(cropY, 0, Math.max(0, srcHeight - 1))
    cropWidth = clamp(cropWidth, 1, srcWidth - cropX)
    cropHeight = clamp(cropHeight, 1, srcHeight - cropY)

    const maxEdge = 1600
    const scale = Math.min(1, maxEdge / Math.max(cropWidth, cropHeight))
    const destWidth = Math.max(1, Math.round(cropWidth * scale))
    const destHeight = Math.max(1, Math.round(cropHeight * scale))

    const croppedPath = await renderCanvasFromSource({
      filePath: cropSourcePath.value,
      srcX: cropX,
      srcY: cropY,
      srcWidth: cropWidth,
      srcHeight: cropHeight,
      destWidth,
      destHeight,
      quality: 0.92,
    })

    if (!croppedPath) {
      uni.showToast({ title: $t.value('uploadFail'), icon: 'none' })
      return
    }

    finishCropEditor(croppedPath)
  } catch (error) {
    uni.showToast({ title: $t.value('uploadFail'), icon: 'none' })
  } finally {
    cropConfirming.value = false
  }
}

const compressImage = async (filePath) => {
  try {
    const res = await new Promise((resolve, reject) => {
      uni.compressImage({
        src: filePath,
        quality: 55,
        success: resolve,
        fail: reject,
      })
    })
    return res.tempFilePath || filePath
  } catch (e) {
    return filePath
  }
}

const pickAndProcessImage = async (mode = 'original') => {
  const chooseSizeType = mode === 'compress' ? ['compressed'] : ['original']
  const chooseRes = await chooseImageAsync(chooseSizeType)
  const selectedPath = chooseRes.tempFilePaths && chooseRes.tempFilePaths[0]
  if (!selectedPath) return null

  let uploadPath = selectedPath
  if (mode === 'crop') {
    uploadPath = await requestCropImage(selectedPath, { width: 3, height: 4 })
    if (!uploadPath) return null
  } else if (mode === 'compress') {
    uploadPath = await compressImage(selectedPath)
  }

  const sizeBytes = await getFileSizeAsync(uploadPath)
  return { uploadPath, sizeBytes }
}

const chooseUploadMode = () => {
  return openModePopup({
    title: '',
    options: [
      { key: 'original', label: $t.value('uploadModeOriginal') },
      { key: 'crop', label: $t.value('uploadModeCrop') },
      { key: 'compress', label: $t.value('uploadModeCompress') },
    ],
  })
}

const normalizeCategory = (value) => {
  const category = String(value || '').trim().toLowerCase()
  const aliasMap = {
    '上装': 'upper',
    '上衣': 'upper',
    '下装': 'lower',
    '下衣': 'lower',
    '裤子': 'lower',
    '裤': 'lower',
    '裙子': 'lower',
    '裙': 'lower',
    '连体': 'onepiece',
    '连体装': 'onepiece',
    '连衣裙': 'onepiece',
    dress: 'onepiece',
    shoe: 'shoes',
    '鞋': 'shoes',
  }
  if (aliasMap[category]) {
    return aliasMap[category]
  }
  return CATEGORY_KEYS.includes(category) ? category : ''
}

const categoryLabel = (category) => {
  const key = CATEGORY_LABEL_KEY_MAP[normalizeCategory(category)]
  return key ? $t.value(key) : '-'
}

const buildDefaultClothName = (category) => {
  const prefix = categoryLabel(category)
  const safePrefix = prefix && prefix !== '-' ? prefix : $t.value('myCloset')
  return `${safePrefix}${Date.now()}`
}

const normalizeClothItem = (item) => {
  const id = item?.ID || item?.id || ''
  const category = normalizeCategory(item?.category || item?.Category)
  return {
    id: String(id),
    name: item?.name || '',
    category,
    url: item?.image || item?.url || '',
    sizeText: '',
  }
}

const filteredClothList = computed(() => {
  if (activeCategory.value === 'all') {
    return clothList.value
  }
  return clothList.value.filter(item => item.category === activeCategory.value)
})

const loadClothes = async () => {
  const res = await getMyTryonClothList()
  if (res.code !== 0) {
    return
  }
  const list = Array.isArray(res?.data?.list) ? res.data.list : []
  clothList.value = list.map(normalizeClothItem).filter(item => item.id && item.category && item.url)
  clothList.value.forEach((item) => fillClothItemSizeText(item))
}

const buildCreateCategoryOptions = () => {
  return categoryOptions
    .filter(item => item.key !== 'all')
    .map(item => ({ key: item.key, label: $t.value(item.labelKey) }))
}

const chooseCreateCategory = async () => {
  return openModePopup({
    title: '',
    options: buildCreateCategoryOptions(),
  })
}

const addCloth = async () => {
  const category = await chooseCreateCategory()
  if (!category) return

  const mode = await chooseUploadMode()
  if (!mode) return

  try {
    const selected = await pickAndProcessImage(mode)
    if (!selected?.uploadPath) return

    uni.showLoading({ title: $t.value('uploading'), mask: true })
    let uploadedUrl = ''
    try {
      const uploadType = category === 'shoes' ? 'shoe' : 'cloth'
      uploadedUrl = await uploadTryonImage(selected.uploadPath, 'cloth-on/my-cloth', uploadType)
    } catch (e) {
      uni.showToast({ title: resolveApiMessage(e?.message, 'uploadFail'), icon: 'none' })
      return
    } finally {
      uni.hideLoading()
    }

    const createRes = await createTryonCloth({
      name: buildDefaultClothName(category),
      category,
      image: uploadedUrl,
    })
    if (createRes.code !== 0) {
      return
    }

    await loadClothes()
    uni.showToast({ title: $t.value('createSuccess'), icon: 'none' })
  } catch (e) {
    const errMsg = String(e?.errMsg || '')
    if (!errMsg.includes('cancel')) {
      uni.showToast({ title: resolveApiMessage(e?.message, 'uploadFail'), icon: 'none' })
    }
  }
}

const renameCloth = (item) => {
  uni.showModal({
    title: $t.value('rename'),
    editable: true,
    placeholderText: $t.value('myClosetNamePlaceholder'),
    content: item.name || '',
    cancelText: $t.value('cancel'),
    confirmText: $t.value('confirm'),
    success: async (res) => {
      if (!res.confirm) return
      const value = (res.content || '').trim()
      if (!value) return

      const renameRes = await updateTryonCloth({
        ID: Number(item.id),
        name: value,
      })
      if (renameRes.code !== 0) {
        return
      }

      item.name = value
      uni.showToast({ title: $t.value('updateSuccess'), icon: 'none' })
    },
  })
}

const removeCloth = (item) => {
  uni.showModal({
    title: $t.value('pendingOrderTitle'),
    content: $t.value('confirmDeleteCloth'),
    cancelText: $t.value('cancel'),
    confirmText: $t.value('confirm'),
    success: async (res) => {
      if (!res.confirm) return

      const deleteRes = await deleteTryonCloth({ ID: Number(item.id) })
      if (deleteRes.code !== 0) {
        return
      }

      clothList.value = clothList.value.filter(v => v.id !== item.id)
      uni.showToast({ title: $t.value('deleteSuccess'), icon: 'none' })
    },
  })
}

const preview = (item) => {
  if (!item?.url) return
  uni.previewImage({ urls: [getUrl(item.url)] })
}

const goBack = () => {
  uni.navigateBack({ delta: 1 })
}

onShow(() => {
  loadClothes()
})
</script>

<style lang="scss">
page {
  background: #f4f7fb;
}

.cloth-page {
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

.tips {
  margin-top: 10rpx;
  margin-bottom: 12rpx;
  font-size: 22rpx;
  color: rgba(15,23,42,0.62);
}

.category-tabs {
  margin-bottom: 12rpx;
  display: flex;
  gap: 8rpx;
  overflow-x: auto;
}

.category-tab {
  flex-shrink: 0;
  min-width: 120rpx;
  height: 56rpx;
  padding: 0 16rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(15,23,42,0.12);
  background: rgba(255,255,255,0.8);
  color: rgba(15,23,42,0.62);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22rpx;
}

.category-tab.active {
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  border: none;
  color: #fff;
}

.list-wrap {
  height: calc(100vh - var(--status-bar-height, 0px) - 190rpx);
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12rpx;
}

.card {
  background: rgba(255,255,255,0.95);
  border: 1rpx solid rgba(15,23,42,0.08);
  border-radius: 12rpx;
  overflow: hidden;
  box-shadow: 0 12rpx 26rpx rgba(15, 23, 42, 0.06);
}

.card-image {
  width: 100%;
  height: 260rpx;
}

.card-image-wrap {
  position: relative;
}

.card-size-badge {
  position: absolute;
  top: 8rpx;
  right: 8rpx;
  z-index: 2;
  padding: 4rpx 10rpx;
  border-radius: 999rpx;
  background: rgba(15, 23, 42, 0.68);
  color: #ffffff;
  font-size: 18rpx;
  line-height: 1.2;
}

.card-foot {
  padding: 10rpx;
}

.name-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8rpx;
}

.card-name {
  flex: 1;
  min-width: 0;
  font-size: 22rpx;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.category-badge {
  flex-shrink: 0;
  height: 40rpx;
  line-height: 40rpx;
  padding: 0 12rpx;
  border-radius: 999rpx;
  background: rgba(37, 99, 235, 0.12);
  color: #1d4ed8;
  font-size: 20rpx;
}

.card-actions {
  margin-top: 10rpx;
  display: flex;
  gap: 8rpx;
}

.card-actions-secondary {
  margin-top: 8rpx;
}

.mini-btn {
  flex: 1;
  height: 50rpx;
  border-radius: 999rpx;
  background: rgba(248,250,252,0.95);
  border: 1rpx solid rgba(15,23,42,0.08);
  color: #0f172a;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20rpx;
}

.mini-btn.danger {
  background: rgba(220, 38, 38, 0.12);
  border-color: rgba(220, 38, 38, 0.24);
  color: #b91c1c;
}

.empty {
  margin-top: 140rpx;
  text-align: center;
  color: rgba(15,23,42,0.55);
}

.crop-editor-mask {
  position: fixed;
  inset: 0;
  z-index: 1200;
  background: rgba(15, 23, 42, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24rpx;
}

.crop-editor-panel {
  width: 100%;
  max-width: 700rpx;
  background: #fff;
  border-radius: 24rpx;
  padding: 24rpx;
  box-shadow: 0 20rpx 48rpx rgba(15, 23, 42, 0.24);
}

.crop-editor-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14rpx;
}

.crop-editor-title {
  font-size: 30rpx;
  font-weight: 700;
  color: #0f172a;
}

.crop-editor-ratio {
  font-size: 22rpx;
  color: rgba(37, 99, 235, 0.92);
  font-weight: 600;
}

.crop-ratio-options {
  display: flex;
  flex-wrap: wrap;
  gap: 10rpx;
  margin-bottom: 12rpx;
}

.crop-ratio-chip {
  padding: 8rpx 16rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.14);
  background: rgba(248, 250, 252, 0.95);
  color: rgba(15, 23, 42, 0.64);
  font-size: 22rpx;
}

.crop-ratio-chip.active {
  border-color: rgba(37, 99, 235, 0.55);
  background: rgba(219, 234, 254, 0.8);
  color: #1d4ed8;
}

.crop-editor-stage-wrap {
  display: flex;
  justify-content: center;
}

.crop-editor-edge-tools {
  margin-top: 12rpx;
}

.crop-editor-edge-tip {
  display: block;
  margin-bottom: 8rpx;
  text-align: center;
  font-size: 20rpx;
  color: rgba(15, 23, 42, 0.56);
}

.crop-editor-edge-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8rpx;
  justify-content: center;
}

.crop-editor-edge-chip {
  min-width: 86rpx;
  height: 46rpx;
  padding: 0 14rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.14);
  background: rgba(248, 250, 252, 0.95);
  color: rgba(15, 23, 42, 0.64);
  font-size: 20rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.crop-editor-edge-chip.active {
  border-color: rgba(37, 99, 235, 0.55);
  background: rgba(219, 234, 254, 0.8);
  color: #1d4ed8;
}

.crop-editor-stage {
  position: relative;
  overflow: hidden;
  border-radius: 12rpx;
  background: #000;
}

.crop-editor-image {
  width: 100%;
  height: 100%;
  display: block;
  opacity: 0.92;
}

.crop-editor-area {
  position: absolute;
  inset: 0;
}

.crop-editor-box {
  box-sizing: border-box;
  border: 2px solid #ffffff;
  background: rgba(37, 99, 235, 0.18);
}

.crop-editor-box-inner {
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  border: 1px dashed rgba(255, 255, 255, 0.78);
}

.crop-editor-zoom {
  margin-top: 12rpx;
  display: flex;
  justify-content: center;
  gap: 12rpx;
}

.crop-editor-zoom-btn {
  width: 70rpx;
  height: 52rpx;
  border-radius: 999rpx;
  background: rgba(15, 23, 42, 0.08);
  color: #0f172a;
  font-size: 34rpx;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
}

.crop-editor-actions {
  margin-top: 16rpx;
  display: flex;
  gap: 12rpx;
}

.crop-editor-btn {
  flex: 1;
  height: 72rpx;
  border-radius: 999rpx;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  color: #fff;
  font-size: 24rpx;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
}

.crop-editor-btn.ghost {
  color: rgba(15, 23, 42, 0.76);
  background: rgba(15, 23, 42, 0.08);
}

.crop-editor-btn.disabled {
  opacity: 0.65;
  pointer-events: none;
}

.crop-canvas-hidden {
  position: fixed;
  left: -9999px;
  top: -9999px;
  opacity: 0;
  pointer-events: none;
}

.mode-popup-mask {
  position: fixed;
  inset: 0;
  z-index: 1300;
  background: rgba(15, 23, 42, 0.48);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding: 24rpx;
}

.mode-popup-panel {
  width: 100%;
  max-width: 700rpx;
  background: #ffffff;
  border-radius: 24rpx;
  padding: 20rpx;
  box-shadow: 0 20rpx 48rpx rgba(15, 23, 42, 0.2);
}

.mode-popup-title {
  display: block;
  margin-bottom: 12rpx;
  text-align: center;
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.62);
}

.mode-popup-item,
.mode-popup-cancel {
  height: 76rpx;
  border-radius: 14rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24rpx;
}

.mode-popup-item {
  margin-bottom: 10rpx;
  background: rgba(37, 99, 235, 0.08);
  color: #1d4ed8;
}

.mode-popup-cancel {
  margin-top: 6rpx;
  background: rgba(15, 23, 42, 0.08);
  color: rgba(15, 23, 42, 0.72);
}
</style>
