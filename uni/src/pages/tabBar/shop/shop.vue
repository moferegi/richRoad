<template>
  <view class="shoe-page">
    <view class="shoe-header">
      <text class="shoe-title">{{ $t('shoeRoom') }}</text>
    </view>

    <view class="shoe-body">
      <view class="shoe-left" @tap="openUploadDrawer('person')">
        <image v-if="personPreview" class="shoe-preview" :src="personPreview" mode="aspectFit" />
        <view v-if="personPreview && personSizeBytes > 0" class="preview-size-mask">{{ formatPreviewSize(personSizeBytes) }}</view>
        <view v-else class="upload-empty">
          <uni-icons type="camera" size="26" color="rgba(15,23,42,0.45)" />
          <text class="upload-text">{{ $t('uploadPersonImage') }}</text>
        </view>
      </view>
      <view class="shoe-right">
        <view class="slot" @tap="openUploadDrawer('shoe')">
          <image v-if="shoePreview" class="slot-preview" :src="shoePreview" mode="aspectFit" />
          <view v-if="shoePreview && shoeSizeBytes > 0" class="preview-size-mask">{{ formatPreviewSize(shoeSizeBytes) }}</view>
          <text v-else class="slot-text">{{ $t('uploadShoeImage') }}</text>
        </view>
        <view class="slot suit-slot" @tap="goClothesPage">
          <text class="slot-text">{{ $t('suitSet') }}</text>
          <text class="slot-tip">{{ $t('goClothesPageTip') }}</text>
        </view>
      </view>
    </view>

    <view class="model-card">
      <view class="model-left" @tap="showModelPopup = true">
        <text class="model-label">{{ $t('currentModelLabel') }}</text>
        <view class="model-info">
          <text class="model-value">{{ currentModel.name }}</text>
          <text class="model-cost-hint">{{ $t('pointsCostEach').replace('{cost}', String(currentCost)) }}</text>
          <text class="model-desc-hint" v-if="currentModelDesc">{{ currentModelDesc }}</text>
        </view>
      </view>
      <view class="model-switch" @tap="showModelPopup = true">
        <text>{{ $t('switchAction') }}</text>
      </view>
    </view>

    <view class="action-btn" @tap="goGenerate">
      <text class="action-text">{{ $t('tryOnShoesWithCost').replace('{cost}', String(currentCost)) }}</text>
    </view>

    <!-- <view class="tips">{{ $t('shoeTryonTip') }}</view> -->

    <view class="popup-mask" v-if="showModelPopup" @tap="showModelPopup = false">
      <view class="popup-panel" @tap.stop>
        <view class="popup-title">{{ $t('selectModel') }}</view>
        <scroll-view scroll-y class="popup-list">
          <view
            v-for="item in modelList"
            :key="item.key"
            class="popup-item"
            :class="{ active: item.key === selectedModelKey }"
            @tap="selectModel(item.key)"
          >
            <view>
              <text class="popup-item-name">{{ item.name }}</text>
              <text class="popup-item-cost">{{ $t('pointsCostEach').replace('{cost}', String(item.cost)) }}</text>
              <text class="popup-item-desc" v-if="item.descText">{{ item.descText }}</text>
            </view>
            <uni-icons type="checkmarkempty" size="20" color="#2563eb" v-if="item.key === selectedModelKey" />
          </view>
        </scroll-view>
      </view>
    </view>

    <view class="upload-mask" v-if="showUploadDrawer" @tap="closeUploadDrawer">
      <view class="upload-drawer" @tap.stop>
        <view class="upload-drawer-head">
          <text class="upload-drawer-title">{{ uploadDrawerTitle }}</text>
          <text class="upload-drawer-tip">{{ uploadDrawerTip }}</text>
        </view>

        <view class="upload-tabs">
          <view
            v-for="tab in drawerTabs"
            :key="tab.key"
            class="upload-tab"
            :class="{ active: drawerTab === tab.key }"
            @tap="drawerTab = tab.key"
          >
            {{ $t(tab.labelKey) }}
          </view>
        </view>

        <scroll-view class="upload-scroll" scroll-y>
          <view class="upload-section" v-if="isPersonDrawer && drawerTab === 'custom'">
            <view class="upload-mode-actions">
              <view class="upload-main-btn" @tap="uploadDrawerChooseImage('original')">{{ $t('uploadModeOriginal') }}</view>
              <view class="upload-main-btn" @tap="uploadDrawerChooseImage('crop')">{{ $t('uploadModeCrop') }}</view>
              <view class="upload-main-btn" @tap="uploadDrawerChooseImage('compress')">{{ $t('uploadModeCompress') }}</view>
            </view>
            <text class="upload-center-tip">{{ $t('uploadSinglePersonTip') }}</text>

            <view class="example-grid two">
              <view
                v-for="item in personGoodExamples"
                :key="item.url"
                class="example-card"
                @tap="applyRemoteExample(item)"
              >
                <image class="example-image" :src="getExamplePreview(item.url)" mode="aspectFit" @error="handleExampleImageError(item.url)" />
                <text class="example-label">{{ item.label }}</text>
              </view>
            </view>

            <view class="divider-line">
              <view class="divider-side"></view>
              <text class="divider-text">{{ $t('errorExampleTitle') }}</text>
              <view class="divider-side"></view>
            </view>

            <view class="example-grid two">
              <view
                v-for="item in personBadExamples"
                :key="item.url"
                class="example-card"
                @tap="applyRemoteExample(item)"
              >
                <image class="example-image" :src="getExamplePreview(item.url)" mode="aspectFit" @error="handleExampleImageError(item.url)" />
                <text class="example-label">{{ item.label }}</text>
              </view>
            </view>
          </view>

          <view class="upload-section" v-if="isPersonDrawer && drawerTab === 'myModel'">
            <view v-if="myModelList.length === 0" class="drawer-empty">
              <text>{{ $t('myModelEmpty') }}</text>
            </view>
            <view v-else class="example-grid two">
              <view
                v-for="item in myModelList"
                :key="item.id"
                class="example-card"
                @tap="applyMyModel(item)"
              >
                <image class="example-image" :src="getExamplePreview(item.url)" mode="aspectFit" @error="handleExampleImageError(item.url)" />
                <text class="example-label">{{ item.name || $t('unnamedModel') }}</text>
              </view>
            </view>
          </view>

          <view class="upload-section" v-if="isPersonDrawer && drawerTab === 'official'">
            <view class="example-grid two">
              <view
                v-for="item in personOfficialExamples"
                :key="item.url"
                class="example-card"
                @tap="applyRemoteExample(item)"
              >
                <image class="example-image" :src="getExamplePreview(item.url)" mode="aspectFit" @error="handleExampleImageError(item.url)" />
                <text class="example-label">{{ item.label }}</text>
              </view>
            </view>
          </view>

          <view class="upload-section" v-if="isShoeDrawer && drawerTab === 'custom'">
            <view class="upload-mode-actions">
              <view class="upload-main-btn" @tap="uploadDrawerChooseImage('original')">{{ $t('uploadModeOriginal') }}</view>
              <view class="upload-main-btn" @tap="uploadDrawerChooseImage('crop')">{{ $t('uploadModeCrop') }}</view>
              <view class="upload-main-btn" @tap="uploadDrawerChooseImage('compress')">{{ $t('uploadModeCompress') }}</view>
            </view>
            <text class="upload-center-tip">{{ $t('uploadSingleShoeTip') }}</text>

            <view class="example-grid two">
              <view
                v-for="item in shoeGoodExamples"
                :key="item.url"
                class="example-card"
                @tap="applyRemoteExample(item)"
              >
                <image class="example-image" :src="getExamplePreview(item.url)" mode="aspectFit" @error="handleExampleImageError(item.url)" />
                <text class="example-label">{{ item.label }}</text>
              </view>
            </view>
          </view>

          <view class="upload-section" v-if="isShoeDrawer && drawerTab === 'recommended'">
            <view class="example-grid two">
              <view
                v-for="item in shoeRecommendedExamples"
                :key="item.url"
                class="example-card"
                @tap="applyRemoteExample(item)"
              >
                <image class="example-image" :src="getExamplePreview(item.url)" mode="aspectFit" @error="handleExampleImageError(item.url)" />
                <text class="example-label">{{ item.label }}</text>
              </view>
            </view>
          </view>
        </scroll-view>

        <view class="drawer-foot-actions">
          <view class="upload-drawer-btn ghost" @tap="clearUploadTarget">{{ $t('clearCurrentAction') }}</view>
          <view class="upload-drawer-btn" @tap="closeUploadDrawer">{{ $t('doneText') }}</view>
        </view>
      </view>
    </view>

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
          <view class="upload-drawer-btn ghost" @tap="cancelCropEditor">{{ $t('cancel') }}</view>
          <view class="upload-drawer-btn" @tap="confirmCropEditor">{{ $t('doneText') }}</view>
        </view>
      </view>
    </view>

    <canvas
      canvas-id="shoeCropCanvas"
      id="shoeCropCanvas"
      class="crop-canvas-hidden"
      :style="{ width: `${cropCanvasSize.width}px`, height: `${cropCanvasSize.height}px` }"
    ></canvas>
  </view>
</template>

<script setup>
import { computed, nextTick, ref, watch } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { getTryonConfig, getDefaultDomain } from '@/api/sysConfig.js'
import { getMyTryonModelList } from '@/api/tryonTask.js'
import { getUrl } from '@/utils/url.js'
import { localText } from '@/utils/i18n.js'
import {
  createTryonRequestId,
  saveTryonDraft,
  getSelectedTryonModel,
  clearSelectedTryonModel,
  parseTryonModels,
  uploadTryonImage,
} from '@/utils/tryon.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const showModelPopup = ref(false)
const showUploadDrawer = ref(false)
const uploadTarget = ref('')
const drawerTab = ref('custom')
const selectedModelKey = ref('')
const modelList = ref([])
const myModelList = ref([])
const defaultExternalDomain = ref('')
const exampleImageFailIndex = ref({})
const tryonConfig = ref({
  tryon_cost_points: '1',
  tryon_models: '',
  shoe_models: '',
})

const personLocal = ref('')
const personRemote = ref('')
const shoeLocal = ref('')
const shoeRemote = ref('')
const personSizeBytes = ref(0)
const shoeSizeBytes = ref(0)
const cropCanvasSize = ref({ width: 1, height: 1 })
const cropEditorVisible = ref(false)
const cropSourcePath = ref('')
const cropSourceSize = ref({ width: 0, height: 0 })
const cropStageSize = ref({ width: 1, height: 1 })
const cropBoxSize = ref({ width: 1, height: 1 })
const cropBoxPosition = ref({ x: 0, y: 0 })
const cropRatio = ref({ width: 3, height: 4 })
const activeCropEdge = ref('all')
let cropResolve = null

const personPreview = computed(() => personRemote.value ? getUrl(personRemote.value) : personLocal.value)
const shoePreview = computed(() => shoeRemote.value ? getUrl(shoeRemote.value) : shoeLocal.value)

const currentModel = computed(() => {
  const selected = modelList.value.find(v => v.key === selectedModelKey.value)
  return selected || modelList.value[0] || { key: 'shoes-and-boots', name: 'shoes-and-boots', cost: Number(tryonConfig.value.tryon_cost_points || 1), desc: {} }
})
const currentCost = computed(() => Number(currentModel.value.cost || 1))
const currentModelDesc = computed(() => currentModel.value.descText || localText(currentModel.value.desc, langStore.locale) || '')
const cropRatioLabel = computed(() => `${cropRatio.value.width}:${cropRatio.value.height}`)
const cropRatioOptions = [
  { key: '1:1', label: '1:1', width: 1, height: 1 },
  { key: '3:4', label: '3:4', width: 3, height: 4 },
  { key: '4:3', label: '4:3', width: 4, height: 3 },
  { key: '9:16', label: '9:16', width: 9, height: 16 },
  { key: '16:9', label: '16:9', width: 16, height: 9 },
]

const CROP_CANVAS_ID = 'shoeCropCanvas'
const CROP_MIN_EDGE_SIZE = 48

const setImageSizeForTarget = (target, bytes = 0) => {
  const safeSize = Number.isFinite(Number(bytes)) && Number(bytes) > 0 ? Math.round(Number(bytes)) : 0
  if (target === 'person') personSizeBytes.value = safeSize
  if (target === 'shoe') shoeSizeBytes.value = safeSize
}

const formatPreviewSize = (bytes) => {
  const size = Number(bytes || 0)
  if (!Number.isFinite(size) || size <= 0) return ''

  const kb = size / 1024
  if (kb < 1024) {
    const value = kb >= 10 ? kb.toFixed(0) : kb.toFixed(1)
    return `${$t.value('uploadImageSizePrefix')}${value}kb`
  }

  const mb = kb / 1024
  const value = mb >= 10 ? mb.toFixed(0) : mb.toFixed(1)
  return `${$t.value('uploadImageSizePrefix')}${value}mb`
}

const waitFrame = (delay = 30) => new Promise((resolve) => setTimeout(resolve, delay))

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
  try {
    const srcWidth = Number(cropSourceSize.value.width || 0)
    const srcHeight = Number(cropSourceSize.value.height || 0)
    const stageWidth = Number(cropStageSize.value.width || 0)
    const stageHeight = Number(cropStageSize.value.height || 0)
    if (!cropSourcePath.value || srcWidth <= 0 || srcHeight <= 0 || stageWidth <= 0 || stageHeight <= 0) {
      finishCropEditor('')
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

    finishCropEditor(croppedPath)
  } catch (error) {
    uni.showToast({ title: $t.value('uploadFail'), icon: 'none' })
    finishCropEditor('')
  }
}

const exportCompressedCanvas = async (filePath, quality, maxEdge) => {
  const info = await getImageInfoAsync(filePath)
  const srcWidth = Number(info.width || 0)
  const srcHeight = Number(info.height || 0)
  if (srcWidth <= 0 || srcHeight <= 0) {
    return ''
  }

  const scale = Math.min(1, maxEdge / Math.max(srcWidth, srcHeight))
  const destWidth = Math.max(1, Math.round(srcWidth * scale))
  const destHeight = Math.max(1, Math.round(srcHeight * scale))

  return renderCanvasFromSource({
    filePath,
    srcWidth,
    srcHeight,
    destWidth,
    destHeight,
    quality,
  })
}

const compressImage = async (filePath) => {
  const originalSize = await getFileSizeAsync(filePath)
  let bestPath = filePath
  let bestSize = originalSize > 0 ? originalSize : Number.MAX_SAFE_INTEGER

  const attempts = [
    { quality: 0.72, maxEdge: 1680 },
    { quality: 0.6, maxEdge: 1360 },
    { quality: 0.48, maxEdge: 1080 },
  ]

  for (const attempt of attempts) {
    try {
      const tempPath = await exportCompressedCanvas(filePath, attempt.quality, attempt.maxEdge)
      if (!tempPath) continue
      const tempSize = await getFileSizeAsync(tempPath)
      if (tempSize > 0 && tempSize < bestSize) {
        bestPath = tempPath
        bestSize = tempSize
      }
    } catch (error) {
      // ignore this attempt and continue
    }
  }

  if (bestPath !== filePath) {
    return bestPath
  }

  try {
    const fallback = await new Promise((resolve, reject) => {
      uni.compressImage({
        src: filePath,
        quality: 55,
        success: resolve,
        fail: reject,
      })
    })
    const fallbackPath = fallback.tempFilePath || filePath
    if (fallbackPath !== filePath) {
      const fallbackSize = await getFileSizeAsync(fallbackPath)
      if (originalSize <= 0 || (fallbackSize > 0 && fallbackSize < originalSize)) {
        return fallbackPath
      }
    }
  } catch (error) {
    // ignore fallback errors
  }

  return filePath
}

const pickAndProcessImage = async (mode = 'original') => {
  const chooseRes = await chooseImageAsync(['original'])
  const selectedPath = chooseRes.tempFilePaths && chooseRes.tempFilePaths[0]
  if (!selectedPath) return null

  let uploadPath = selectedPath
  if (mode === 'crop') {
    uploadPath = await requestCropImage(selectedPath)
    if (!uploadPath) return null
  } else if (mode === 'compress') {
    uploadPath = await compressImage(selectedPath)
  }

  const sizeBytes = await getFileSizeAsync(uploadPath)
  return { uploadPath, sizeBytes }
}

const EXAMPLE_PRIMARY_DOMAIN = 'https://video.mnmovie.icu'
const EXAMPLE_PRIMARY_PATH = '/cloth-on/uni-set'
const EXAMPLE_LEGACY_PATH = '/file/Moffuu/cloth-on/uni-set'
const EXAMPLE_FALLBACK_DOMAIN = 'https://f005.backblazeb2.com'
const b2Image = (name) => `${EXAMPLE_PRIMARY_PATH}/${name}`

const buildIndexedExampleItems = (names, prefixKey) => {
  return names.map((name, index) => ({
    label: `${$t.value(prefixKey)}${index + 1}`,
    url: b2Image(name),
  }))
}

const normalizeDomain = (rawDomain) => {
  const value = String(rawDomain || '').trim().replace(/\/+$/, '')
  if (!value) return ''
  if (value.startsWith('//')) return `https:${value}`
  if (value.startsWith('http://') || value.startsWith('https://')) return value
  return `https://${value}`
}

const joinDomainPath = (domain, path) => {
  const safeDomain = normalizeDomain(domain)
  const safePath = String(path || '').trim()
  if (!safeDomain || !safePath) return ''
  const sep = safePath.startsWith('/') ? '' : '/'
  return `${safeDomain}${sep}${safePath}`
}

const toPrimaryExamplePath = (path) => {
  const value = String(path || '').trim()
  if (value.startsWith(`${EXAMPLE_LEGACY_PATH}/`)) {
    return value.replace(EXAMPLE_LEGACY_PATH, EXAMPLE_PRIMARY_PATH)
  }
  return value
}

const toLegacyExamplePath = (path) => {
  const value = String(path || '').trim()
  if (value.startsWith(`${EXAMPLE_PRIMARY_PATH}/`)) {
    return value.replace(EXAMPLE_PRIMARY_PATH, EXAMPLE_LEGACY_PATH)
  }
  return value
}

const getExamplePathVariants = (path) => {
  const primaryPath = toPrimaryExamplePath(path)
  const legacyPath = toLegacyExamplePath(primaryPath)
  return Array.from(new Set([primaryPath, legacyPath].filter(Boolean)))
}

const getExampleCandidates = (url) => {
  const value = String(url || '').trim()
  if (!value) return []
  if (value.startsWith('http://') || value.startsWith('https://') || value.startsWith('data:')) {
    return [value]
  }

  const pathVariants = getExamplePathVariants(value)

  const fromPreferredDomain = pathVariants
    .map(path => joinDomainPath(EXAMPLE_PRIMARY_DOMAIN, path))
    .filter(Boolean)

  const fromDefaultDomain = pathVariants
    .map(path => joinDomainPath(defaultExternalDomain.value, path))
    .filter(Boolean)

  const fromFallbackDomain = pathVariants
    .map(path => joinDomainPath(EXAMPLE_FALLBACK_DOMAIN, path))
    .filter(Boolean)

  const fromGetUrl = pathVariants
    .map(path => getUrl(path))
    .filter(Boolean)

  const candidates = [
    ...fromPreferredDomain,
    ...fromDefaultDomain,
    ...fromFallbackDomain,
    ...fromGetUrl,
  ]

  return Array.from(new Set(candidates.filter(Boolean)))
}

const getExamplePreview = (url) => {
  const key = String(url || '')
  const candidates = getExampleCandidates(key)
  if (!candidates.length) return ''

  const currentIndex = Number(exampleImageFailIndex.value[key] || 0)
  return candidates[Math.min(currentIndex, candidates.length - 1)]
}

const handleExampleImageError = (url) => {
  const key = String(url || '')
  if (!key) return

  const candidates = getExampleCandidates(key)
  if (candidates.length <= 1) return

  const currentIndex = Number(exampleImageFailIndex.value[key] || 0)
  if (currentIndex >= candidates.length - 1) return

  exampleImageFailIndex.value = {
    ...exampleImageFailIndex.value,
    [key]: currentIndex + 1,
  }
}

const personGoodExamples = computed(() => buildIndexedExampleItems(
  ['shoeperson_0001.jpg', 'shoeperson2_0001.jpg'],
  'tryonExampleSinglePersonPrefix'
))

const personBadExamples = computed(() => buildIndexedExampleItems(
  ['shoecant_0001.jpg', 'shoecant2_0001.jpg'],
  'tryonExampleBadPrefix'
))

const personOfficialExamples = computed(() => buildIndexedExampleItems(
  ['shoeperson_0001.jpg', 'shoeperson2_0001.jpg'],
  'tryonExampleOfficialModelPrefix'
))

const shoeGoodExamples = computed(() => buildIndexedExampleItems(
  ['shoepic_0001.jpg', 'shoepic2_0001.jpg', 'shoepic3_0001.jpg', 'shoepic5_0001.jpg'],
  'tryonExampleShoePrefix'
))

const shoeRecommendedExamples = computed(() => buildIndexedExampleItems(
  ['shoepic_0001.jpg', 'shoepic2_0001.jpg', 'shoepic3_0001.jpg', 'shoepic5_0001.jpg'],
  'tryonExampleShoePrefix'
))

const uploadSceneKey = computed(() => uploadTarget.value || 'person')
const isPersonDrawer = computed(() => uploadSceneKey.value === 'person')
const isShoeDrawer = computed(() => uploadSceneKey.value === 'shoe')

const drawerTabs = computed(() => {
  if (isPersonDrawer.value) {
    return [
      { key: 'custom', labelKey: 'drawerTabCustomUpload' },
      { key: 'myModel', labelKey: 'drawerTabMyModel' },
      { key: 'official', labelKey: 'drawerTabOfficialModel' },
    ]
  }
  return [
    { key: 'custom', labelKey: 'drawerTabCustomUpload' },
    { key: 'recommended', labelKey: 'drawerTabRecommended' },
  ]
})

const uploadDrawerTitle = computed(() => {
  return uploadTarget.value === 'shoe' ? $t.value('uploadShoeImage') : $t.value('uploadPersonImage')
})

const uploadDrawerTip = computed(() => {
  if (isPersonDrawer.value) {
    return $t.value('uploadDrawerTipPerson')
  }
  return $t.value('uploadDrawerTipShoes')
})

const rebuildModelList = () => {
  const shoesModelsRaw = tryonConfig.value.shoe_models || tryonConfig.value.tryon_models
  modelList.value = parseTryonModels(
    shoesModelsRaw,
    'shoes',
    Number(tryonConfig.value.tryon_cost_points || 1),
    langStore.locale
  )
  if (!modelList.value.find(item => item.key === selectedModelKey.value)) {
    selectedModelKey.value = modelList.value[0]?.key || ''
  }
}

const assignImageToTarget = (target, value, isRemote = false, sizeBytes = 0) => {
  const localPath = isRemote ? '' : value
  const remotePath = isRemote ? value : ''

  if (target === 'person') {
    personLocal.value = localPath
    personRemote.value = remotePath
  }
  if (target === 'shoe') {
    shoeLocal.value = localPath
    shoeRemote.value = remotePath
  }

  setImageSizeForTarget(target, sizeBytes)
}

const normalizeMyModelItem = (item) => {
  const id = item?.ID || item?.id || ''
  const rawUrl = item?.image || item?.url || ''
  return {
    id: String(id),
    name: item?.name || '',
    url: getUrl(rawUrl),
  }
}

const loadMyModelList = async () => {
  try {
    const res = await getMyTryonModelList()
    if (res.code !== 0) return
    const list = Array.isArray(res?.data?.list) ? res.data.list : []
    myModelList.value = list.map(normalizeMyModelItem).filter(item => item.id && item.url)
  } catch (e) {
    myModelList.value = []
  }
}

const loadExampleDomain = async () => {
  try {
    const res = await getDefaultDomain()
    if (res.code !== 0 || !res.data) {
      defaultExternalDomain.value = ''
      return
    }
    const rawDomain = typeof res.data === 'object'
      ? (res.data.domain || res.data.configValue || '')
      : res.data
    defaultExternalDomain.value = normalizeDomain(rawDomain)
    exampleImageFailIndex.value = {}
  } catch (e) {
    defaultExternalDomain.value = ''
    exampleImageFailIndex.value = {}
  }
}

const uploadFolderByTarget = () => 'cloth-on/uni-up'

const uploadTypeByTarget = (target) => {
  if (target === 'person') return 'person'
  if (target === 'shoe') return 'shoe'
  return 'cloth'
}

const chooseAndUploadImage = async (target, mode = 'original') => {
  if (!target) return
  try {
    const selected = await pickAndProcessImage(mode)
    if (!selected?.uploadPath) return

    uni.showLoading({ title: $t.value('uploading'), mask: true })
    try {
      const remoteUrl = await uploadTryonImage(selected.uploadPath, uploadFolderByTarget(target), uploadTypeByTarget(target))
      assignImageToTarget(target, remoteUrl, true, selected.sizeBytes)
      uni.showToast({ title: $t.value('uploadSuccess'), icon: 'none' })
    } catch (e) {
      uni.showToast({ title: e.message || $t.value('uploadFail'), icon: 'none' })
    } finally {
      uni.hideLoading()
    }
  } catch (e) {
    const errMsg = String(e?.errMsg || '')
    if (!errMsg.includes('cancel')) {
      uni.showToast({ title: e.message || $t.value('uploadFail'), icon: 'none' })
    }
  }
}

const openUploadDrawer = (target) => {
  uploadTarget.value = target
  drawerTab.value = 'custom'
  showUploadDrawer.value = true
  if (target === 'person') {
    loadMyModelList()
  }
}

const closeUploadDrawer = () => {
  showUploadDrawer.value = false
}

const uploadDrawerChooseImage = (mode = 'original') => {
  const target = uploadTarget.value
  if (!target) return
  showUploadDrawer.value = false
  chooseAndUploadImage(target, mode)
}

const clearUploadTarget = () => {
  if (!uploadTarget.value) return
  assignImageToTarget(uploadTarget.value, '', false)
  showUploadDrawer.value = false
}

const applyRemoteExample = (item) => {
  if (!uploadTarget.value || !item?.url) return
  const previewUrl = getExamplePreview(item.url)
  if (!previewUrl) return
  assignImageToTarget(uploadTarget.value, previewUrl, true)
  showUploadDrawer.value = false
  uni.showToast({ title: $t.value('autoFillApplied'), icon: 'none' })
}

const applyMyModel = (item) => {
  if (!item?.url) return
  assignImageToTarget('person', item.url, true)
  showUploadDrawer.value = false
  uni.showToast({ title: $t.value('autoFillApplied'), icon: 'none' })
}

const selectModel = (key) => {
  selectedModelKey.value = key
  showModelPopup.value = false
}

const loadConfig = async () => {
  const res = await getTryonConfig()
  if (res.code === 0 && res.data) {
    tryonConfig.value = { ...tryonConfig.value, ...res.data }
  }
  rebuildModelList()
}

const applySelectedModel = () => {
  const selectedModel = getSelectedTryonModel()
  if (!selectedModel || typeof selectedModel !== 'object') return
  if (selectedModel.roomType && selectedModel.roomType !== 'shoe') return

  if (selectedModel.remoteUrl) {
    personRemote.value = selectedModel.remoteUrl
    personLocal.value = ''
    personSizeBytes.value = 0
    clearSelectedTryonModel()
    return
  }

  if (selectedModel.localPath) {
    personLocal.value = selectedModel.localPath
    personRemote.value = ''
    personSizeBytes.value = 0
    clearSelectedTryonModel()
  }
}

const goClothesPage = () => {
  uni.switchTab({ url: '/pages/tabBar/clothes/index' })
}

const goGenerate = () => {
  if (!personLocal.value && !personRemote.value) {
    uni.showToast({ title: $t.value('uploadPersonFirst'), icon: 'none' })
    return
  }
  if (!shoeLocal.value && !shoeRemote.value) {
    uni.showToast({ title: $t.value('uploadShoeFirst'), icon: 'none' })
    return
  }

  saveTryonDraft({
    roomType: 'shoe',
    sceneType: 'shoes',
    operationType: 'tryon',
    sourceLocalPath: personLocal.value,
    sourceRemoteUrl: personRemote.value,
    sourceUploadFolder: 'cloth-on/uni-up',
    templateLocalPath: shoeLocal.value,
    templateRemoteUrl: shoeRemote.value,
    templateUploadFolder: 'cloth-on/uni-up',
    modelKey: currentModel.value.key,
    modelName: currentModel.value.name,
    modelCost: currentCost.value,
    requestID: createTryonRequestId(),
  })

  uni.navigateTo({ url: '/pages/tryon/generate' })
}

watch(
  [() => tryonConfig.value.shoe_models, () => tryonConfig.value.tryon_models, () => tryonConfig.value.tryon_cost_points, () => langStore.locale],
  () => {
    rebuildModelList()
  }
)

watch([isPersonDrawer, drawerTab], ([isPerson, tab]) => {
  if (isPerson && tab === 'myModel') {
    loadMyModelList()
  }
})

onShow(() => {
  loadExampleDomain()
  loadConfig()
  loadMyModelList()
  applySelectedModel()
})
</script>

<style lang="scss">
page {
  background: #f4f7fb;
}

.shoe-page {
  min-height: 100vh;
  padding: calc(var(--status-bar-height, 0px) + 24rpx) 24rpx 24rpx;
  color: #0f172a;
  background: radial-gradient(120% 80% at 100% -10%, #dbeafe 0%, transparent 60%), #f4f7fb;
}

.shoe-title {
  font-size: 38rpx;
  font-weight: 700;
  margin-bottom: 20rpx;
}

.shoe-body {
  display: flex;
  gap: 16rpx;
}

.shoe-left {
  width: 70%;
  height: 660rpx;
  position: relative;
  border-radius: 20rpx;
  border: 1rpx dashed rgba(15, 23, 42, 0.2);
  overflow: hidden;
  background: rgba(255, 255, 255, 0.95);
  box-shadow: 0 12rpx 34rpx rgba(15, 23, 42, 0.06);
}

.shoe-right {
  width: 30%;
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.slot {
  flex: 1;
  position: relative;
  border-radius: 16rpx;
  border: 1rpx dashed rgba(15, 23, 42, 0.2);
  background: rgba(255, 255, 255, 0.95);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  overflow: hidden;
}

.suit-slot {
  border-style: solid;
  border-color: rgba(37, 99, 235, 0.45);
}

.slot-tip {
  margin-top: 8rpx;
  font-size: 20rpx;
  color: rgba(15, 23, 42, 0.55);
}

.slot-text {
  color: rgba(15, 23, 42, 0.82);
  font-size: 24rpx;
}

.shoe-preview,
.slot-preview {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.preview-size-mask {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20rpx;
  color: #ffffff;
  background: rgba(15, 23, 42, 0.38);
  pointer-events: none;
}

.upload-empty {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 12rpx;
}

.upload-text {
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.6);
}

.model-card {
  margin-top: 18rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.08);
  border-radius: 16rpx;
  padding: 20rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 10rpx 28rpx rgba(15, 23, 42, 0.05);
}

.model-left {
  display: flex;
  align-items: flex-start;
  gap: 8rpx;
  flex: 1;
}

.model-info {
  min-width: 0;
}

.model-label {
  color: rgba(15, 23, 42, 0.55);
  font-size: 24rpx;
  margin-top: 2rpx;
}

.model-value {
  display: block;
  color: #0f172a;
  font-size: 26rpx;
  font-weight: 600;
}

.model-cost-hint {
  display: block;
  margin-top: 4rpx;
  font-size: 21rpx;
  color: rgba(37, 99, 235, 0.9);
}

.model-desc-hint {
  display: -webkit-box;
  margin-top: 4rpx;
  font-size: 21rpx;
  color: rgba(15, 23, 42, 0.6);
  overflow: hidden;
  text-overflow: ellipsis;
  line-clamp: 1;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
}

.model-switch {
  min-width: 92rpx;
  height: 52rpx;
  border-radius: 999rpx;
  padding: 0 18rpx;
  border: 1rpx solid rgba(37, 99, 235, 0.35);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #2563eb;
  font-size: 22rpx;
  font-weight: 600;
}

.action-btn {
  margin-top: 20rpx;
  height: 88rpx;
  border-radius: 999rpx;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 14rpx 30rpx rgba(37, 99, 235, 0.32);
}

.action-text {
  font-size: 28rpx;
  font-weight: 700;
  color: #fff;
}

.tips {
  margin-top: 14rpx;
  color: rgba(15, 23, 42, 0.55);
  font-size: 22rpx;
}

.popup-mask {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.34);
  display: flex;
  justify-content: center;
  align-items: flex-end;
  z-index: 999;
}

.popup-panel {
  width: 100%;
  border-top-left-radius: 24rpx;
  border-top-right-radius: 24rpx;
  background: #ffffff;
  padding: 24rpx;
}

.popup-title {
  font-size: 30rpx;
  font-weight: 700;
  margin-bottom: 16rpx;
  color: #0f172a;
}

.popup-list {
  max-height: 520rpx;
}

.popup-item {
  border: 1rpx solid rgba(15, 23, 42, 0.08);
  border-radius: 14rpx;
  padding: 18rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12rpx;
  background: #fff;
}

.popup-item.active {
  border-color: rgba(37, 99, 235, 0.45);
  background: rgba(219, 234, 254, 0.75);
}

.popup-item-name {
  display: block;
  font-size: 26rpx;
  margin-bottom: 6rpx;
  color: #0f172a;
}

.popup-item-cost {
  display: block;
  color: rgba(15, 23, 42, 0.55);
  font-size: 22rpx;
}

.popup-item-desc {
  display: -webkit-box;
  margin-top: 6rpx;
  max-width: 460rpx;
  color: rgba(15, 23, 42, 0.56);
  font-size: 20rpx;
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  line-clamp: 2;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.upload-mask {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.34);
  display: flex;
  justify-content: center;
  align-items: flex-end;
  z-index: 1000;
}

.upload-drawer {
  width: 100%;
  background: #ffffff;
  border-top-left-radius: 28rpx;
  border-top-right-radius: 28rpx;
  padding: 24rpx;
  max-height: 88vh;
  display: flex;
  flex-direction: column;
}

.upload-drawer-head {
  margin-bottom: 16rpx;
}

.upload-drawer-title {
  display: block;
  font-size: 30rpx;
  font-weight: 700;
  color: #0f172a;
}

.upload-drawer-tip {
  display: block;
  margin-top: 8rpx;
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.56);
}

.upload-tabs {
  display: flex;
  gap: 12rpx;
  margin-bottom: 16rpx;
}

.upload-tab {
  flex: 1;
  text-align: center;
  padding: 14rpx 0;
  border-radius: 999rpx;
  background: rgba(15, 23, 42, 0.07);
  color: rgba(15, 23, 42, 0.65);
  font-size: 22rpx;
}

.upload-tab.active {
  color: #ffffff;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
}

.upload-scroll {
  max-height: 60vh;
}

.upload-section {
  padding-bottom: 10rpx;
}

.upload-main-btn {
  height: 72rpx;
  border-radius: 999rpx;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  color: #fff;
  font-size: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.upload-mode-actions {
  display: flex;
  flex-direction: column;
  gap: 10rpx;
  margin-bottom: 14rpx;
}

.upload-center-tip {
  display: block;
  text-align: center;
  color: rgba(15, 23, 42, 0.6);
  font-size: 22rpx;
  margin-bottom: 14rpx;
}

.example-grid {
  display: grid;
  gap: 12rpx;
}

.example-grid.two {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.example-card {
  border: 1rpx solid rgba(15, 23, 42, 0.1);
  border-radius: 12rpx;
  overflow: hidden;
  background: #fff;
}

.example-image {
  width: 100%;
  height: 180rpx;
  display: block;
  background: #f8fafc;
}

.example-label {
  display: block;
  padding: 8rpx;
  text-align: center;
  font-size: 20rpx;
  color: rgba(15, 23, 42, 0.75);
  line-height: 1.45;
}

.divider-line {
  display: flex;
  align-items: center;
  gap: 10rpx;
  margin: 14rpx 0;
}

.divider-side {
  flex: 1;
  height: 1rpx;
  background: rgba(15, 23, 42, 0.15);
}

.divider-text {
  font-size: 20rpx;
  color: rgba(15, 23, 42, 0.55);
}

.drawer-empty {
  padding: 36rpx 20rpx;
  text-align: center;
  color: rgba(15, 23, 42, 0.5);
  font-size: 22rpx;
}

.drawer-foot-actions {
  margin-top: 14rpx;
  display: flex;
  gap: 12rpx;
}

.upload-drawer-btn {
  flex: 1;
  text-align: center;
  border-radius: 999rpx;
  padding: 14rpx 10rpx;
  font-size: 24rpx;
  color: #fff;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
}

.upload-drawer-btn.ghost {
  color: rgba(15, 23, 42, 0.76);
  background: rgba(15, 23, 42, 0.08);
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

.crop-canvas-hidden {
  position: fixed;
  left: -9999px;
  top: -9999px;
  opacity: 0;
  pointer-events: none;
}
</style>
