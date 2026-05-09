<template>
  <view class="generate-page">
    <view class="nav">
      <view class="nav-btn" @tap="goBack">
        <uni-icons type="left" size="20" color="#0f172a" />
      </view>
      <text class="nav-title">{{ $t('tryonResultImage') }}</text>
      <view class="nav-btn" @tap="goHistory">
        <text class="nav-link">{{ $t('tryonMyTasks') }}</text>
      </view>
    </view>

    <view class="preview-panel">
      <view class="preview-col">
        <text class="preview-label">{{ $t('tryonSourceImage') }}</text>
        <view class="preview-image-wrap" @tap="previewSource">
          <LazyImage class="preview-image" :src="sourcePreview" mode="aspectFit" />
        </view>
      </view>
      <view class="preview-col">
        <text class="preview-label">{{ $t('tryonTemplateImage') }}</text>
        <view class="preview-image-wrap" @tap="previewTemplate">
          <LazyImage class="preview-image" :src="templatePreview" mode="aspectFit" />
        </view>
      </view>
    </view>

    <view class="generate-notice">
      <text>{{ $t('tryonGenerateNotice') }}</text>
    </view>

    <view class="status-card">
      <view class="status-top">
        <text class="status-title">{{ $t('taskStatusTitle') }}</text>
        <text class="status-pill" :class="taskStatusClass">{{ taskStatusText }}</text>
      </view>
      <view class="task-no" v-if="taskNo">{{ $t('tryonTaskNo') }}：{{ taskNo }}</view>
      <view class="task-cost" v-if="draftCost > 0">
        <text class="task-cost-line">{{ $t('generateWithCost').replace('{cost}', String(draftCost)) }}</text>
        <text class="task-cost-line" v-if="draftEnableRefiner && draftRefinerExtraCost > 0">
          {{ $t('tryonRefinerExtraCostHint').replace('{cost}', String(draftRefinerExtraCost)) }}
        </text>
      </view>
      <view class="task-cost" v-if="showBeautifyMeta">
        <text class="task-cost-line" v-if="beautifyTaskNo">{{ $t('tryonBeautifyTaskNo') }}：{{ beautifyTaskNo }}</text>
        <text class="task-cost-line">{{ beautifyHintText }}</text>
        <text class="task-cost-line" v-if="beautifyStatusText">{{ beautifyStatusText }}</text>
      </view>
      <view class="progress-wrap" v-if="isGenerating">
        <view class="progress-track">
          <view class="progress-fill" :style="{ width: progress + '%' }"></view>
        </view>
        <text class="progress-text">{{ $t('tryonStatusProcessing') }} {{ progress }}%</text>
      </view>
      <text class="error-msg" v-if="errorMessage">{{ errorMessage }}</text>
      <view class="status-actions" v-if="showStatusActions">
        <view class="status-action-btn highlight" v-if="canCompare" @tap="openCompare">{{ compareActionText }}</view>
        <view class="status-action-btn" :class="{ disabled: !canUseBeautify || isBeautifyBusy }" @tap="handleBeautify">
          {{ beautifyButtonText }}
        </view>
        <text class="status-action-tip" v-if="showBeautifyMeta">{{ beautifyHintText }}</text>
      </view>
    </view>

    <view class="result-card">
      <text class="result-save-tip">{{ $t('previewLongPressSaveHint') }}</text>
      <view class="result-tabs">
        <view class="result-tab" :class="{ active: resultTab === 'tryon', disabled: !resultPreview }" @tap="switchResultTab('tryon')">
          {{ $t('tryonResultTab') }}
        </view>
        <view class="result-tab" :class="{ active: resultTab === 'beautify', disabled: !beautifyResultPreview }" @tap="switchResultTab('beautify')">
          {{ $t('beautifyResultTab') }}
        </view>
      </view>
      <view class="result-box" v-if="activeResultPreview">
        <image class="result-image" :src="activeResultImageSrc" mode="aspectFit" @tap="previewResult" @load="onResultImageLoad" @error="onResultImageError" />
        <view class="result-loading" v-if="resultTab === 'tryon' && taskStatus === 'success' && resultImageLoading">
          <text>{{ $t('loading') }}</text>
        </view>
      </view>
      <view class="result-empty" v-else>
        <text>{{ resultEmptyHint }}</text>
      </view>
    </view>

    <view class="action-row">
      <view class="action-btn secondary" v-if="canDownloadResult" @tap="downloadResult">{{ $t('downloadAction') }}</view>
      <view class="action-btn" @tap="goContinue">{{ $t('continueTryonAction') }}</view>
    </view>

    <view class="compare-mask" v-if="compareVisible" @tap="closeCompare">
      <view class="compare-panel" @tap.stop>
        <view class="compare-head">
          <text class="compare-title">{{ compareDialogTitle }}</text>
          <view class="compare-close" @tap="closeCompare">
            <uni-icons type="closeempty" size="20" color="#0f172a" />
          </view>
        </view>

        <view class="compare-mode-row">
          <view class="compare-mode-tab" :class="{ active: compareMode === 'tryon', disabled: !canCompareTryon }" @tap="switchCompareMode('tryon')">
            {{ $t('compareTryonTab') }}
          </view>
          <view class="compare-mode-tab" :class="{ active: compareMode === 'beautify', disabled: !canCompareBeautify }" @tap="switchCompareMode('beautify')">
            {{ $t('compareBeautifyTab') }}
          </view>
        </view>

        <view
          class="compare-stage"
          @touchstart.stop.prevent="onCompareStageTouchStart"
          @touchmove.stop.prevent="onCompareStageTouchMove"
          @touchend.stop="onCompareStageTouchEnd"
          @touchcancel.stop="onCompareStageTouchEnd"
        >
          <LazyImage class="compare-image" :src="compareOriginPreview" mode="aspectFit" :style="compareImageStyle" />
          <view class="compare-result-layer" :style="{ width: `${comparePercent}%` }">
            <LazyImage class="compare-image compare-result-image" :src="compareResultPreview" mode="aspectFit" :style="compareResultInnerStyle" />
          </view>
          <view class="compare-divider" :style="{ left: `${comparePercent}%` }"></view>
        </view>

        <view class="compare-slider-wrap">
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
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, nextTick, ref } from 'vue'
import { onLoad, onUnload } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { createTryonTask, findTryonTask, applyTryonBeautify } from '@/api/tryonTask.js'
import { getTryonConfig } from '@/api/sysConfig.js'
import { resolveApiMessage } from '@/utils/i18n.js'
import { getUrl } from '@/utils/url.js'
import LazyImage from '@/components/lazy-image/lazy-image.vue'
import {
  appendTryonLocalHistory,
  clearTryonDraft,
  getTryonDraft,
  getTryonUploadFolder,
  parseTryonBeautifyModels,
  uploadTryonImage,
} from '@/utils/tryon.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const draft = ref({})
const task = ref(null)
const sourcePreview = ref('')
const templatePreview = ref('')
const resultPreview = ref('')
const beautifyResultPreview = ref('')
const resultTab = ref('tryon')
const progress = ref(5)
const isGenerating = ref(false)
const errorMessage = ref('')
const isHistoryMode = ref(false)
const compareVisible = ref(false)
const compareMode = ref('tryon')
const comparePercent = ref(50)
const compareZoomPercent = ref(100)
const compareStageWidthPx = ref(0)
const compareStageHeightPx = ref(0)
const resultImageLoading = ref(false)
const resultImageRetryKey = ref(0)
const beautifyLoading = ref(false)
const modelMeta = ref({})
const beautifyPolling = ref(false)
const autoBeautifyRequested = ref(false)
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

let pollTimer = null
let pollCount = 0
let resultImageRetryTimer = null
let beautifyPollTimer = null
let beautifyPollCount = 0

const isTempLocalPath = (value) => {
  const text = String(value || '').trim().toLowerCase()
  if (!text) return false
  return text.startsWith('blob:') || text.startsWith('file:') || text.startsWith('wxfile:') || text.startsWith('content:')
}
const POLL_INTERVAL = 2000
// Plus model tasks can take over 1 minute on provider side.
const POLL_MAX_COUNT = 90

const taskNo = computed(() => task.value?.taskNo || '')
const taskStatus = computed(() => task.value?.status || '')
const beautifyStatus = computed(() => {
  const status = String(task.value?.beautifyStatus || '').trim().toLowerCase()
  return status || 'disabled'
})
const beautifyTaskNo = computed(() => task.value?.beautifyTaskNo || '')
const draftCost = computed(() => Math.max(0, Number(draft.value?.modelCost || 0)))
const draftEnableRefiner = computed(() => !!draft.value?.enableRefiner)
const draftRefinerExtraCost = computed(() => Math.max(0, Number(draft.value?.refinerExtraCost || 0)))
const beautifySupported = computed(() => {
  const fromMeta = modelMeta.value?.supportsBeautify
  if (fromMeta !== undefined && fromMeta !== null) {
    return !!fromMeta
  }
  if (draft.value?.supportsBeautify !== undefined && draft.value?.supportsBeautify !== null) {
    return !!draft.value?.supportsBeautify
  }
  return true
})
const beautifyConfiguredCost = computed(() => {
  const fromMeta = Number(modelMeta.value?.beautifyExtraCost || 0)
  if (fromMeta > 0) return Math.max(0, fromMeta)
  return Math.max(0, Number(draft.value?.beautifyExtraCost || 0))
})
const beautifyCostActual = computed(() => Math.max(0, Number(task.value?.beautifyCost || 0)))
const beautifyCostDisplay = computed(() => {
  if (beautifyCostActual.value > 0) {
    return beautifyCostActual.value
  }
  return beautifyConfiguredCost.value
})
const hasBeautifyResult = computed(() => !!beautifyResultPreview.value)
const showBeautifyMeta = computed(() => beautifySupported.value || hasBeautifyResult.value || beautifyStatus.value !== 'disabled')
const beautifyStatusText = computed(() => {
  if (beautifyStatus.value === 'processing') return $t.value('tryonBeautifyProcessing')
  if (beautifyStatus.value === 'success') return $t.value('tryonBeautifySuccess')
  if (beautifyStatus.value === 'failed') return $t.value('tryonBeautifyFailed')
  return ''
})
const beautifyHintText = computed(() => {
  if (!beautifySupported.value && !hasBeautifyResult.value) {
    return $t.value('tryonBeautifyUnsupportedHint')
  }
  if (beautifyStatus.value === 'success' || beautifyStatus.value === 'failed') {
    return $t.value('tryonBeautifyUsedHint')
  }
  if (beautifyCostDisplay.value > 0) {
    return $t.value('tryonBeautifyCostHint').replace('{cost}', String(beautifyCostDisplay.value))
  }
  return $t.value('tryonBeautifyFreeHint')
})
const isBeautifyBusy = computed(() => beautifyLoading.value || beautifyPolling.value || beautifyStatus.value === 'processing')
const canUseBeautify = computed(() => {
  if (!task.value?.ID) return false
  if (taskStatus.value !== 'success') return false
  if (!resultPreview.value) return false
  if (!beautifySupported.value) return false
  return beautifyStatus.value === 'disabled'
})
const beautifyButtonText = computed(() => {
  if (isBeautifyBusy.value) return $t.value('tryonBeautifyProcessing')
  if (beautifyStatus.value === 'success') return $t.value('tryonBeautifyDone')
  if (beautifyStatus.value === 'failed') return $t.value('tryonBeautifyUsedHint')
  return $t.value('tryonBeautifyAction')
})
const taskStatusText = computed(() => {
  if (taskStatus.value === 'success') return $t.value('tryonStatusSuccess')
  if (taskStatus.value === 'failed') return $t.value('tryonStatusFailed')
  if (isGenerating.value) return $t.value('tryonStatusProcessing')
  return $t.value('loading')
})
const taskStatusClass = computed(() => {
  if (taskStatus.value === 'success') return 'success'
  if (taskStatus.value === 'failed') return 'failed'
  if (isGenerating.value) return 'processing'
  return ''
})
const activeResultPreview = computed(() => {
  if (resultTab.value === 'beautify') {
    return beautifyResultPreview.value
  }
  return resultPreview.value
})
const resultEmptyHint = computed(() => {
  if (resultTab.value === 'beautify') {
    return $t.value('tryonBeautifyEmptyHint')
  }
  return $t.value('resultReadyHint')
})
const canDownloadResult = computed(() => !!activeResultPreview.value)
const canCompareTryon = computed(() => !!sourcePreview.value && !!resultPreview.value)
const canCompareBeautify = computed(() => !!resultPreview.value && !!beautifyResultPreview.value)
const canCompare = computed(() => {
  if (resultTab.value === 'beautify') {
    return canCompareBeautify.value
  }
  return canCompareTryon.value
})
const showStatusActions = computed(() => canCompare.value || showBeautifyMeta.value)
const compareActionText = computed(() => {
  if (resultTab.value === 'beautify') {
    return $t.value('compareBeautifyButton')
  }
  return $t.value('compareImageButton')
})
const compareDialogTitle = computed(() => {
  if (compareMode.value === 'beautify') {
    return $t.value('compareBeautifyButton')
  }
  return $t.value('compareImageButton')
})
const compareOriginPreview = computed(() => {
  if (compareMode.value === 'beautify') {
    return resultPreview.value
  }
  return sourcePreview.value
})
const compareResultPreview = computed(() => {
  if (compareMode.value === 'beautify') {
    return beautifyResultPreview.value
  }
  return resultPreview.value
})
const compareLeftLabel = computed(() => {
  if (compareMode.value === 'beautify') {
    return $t.value('tryonResultTab')
  }
  return $t.value('previewOriginTab')
})
const compareRightLabel = computed(() => {
  if (compareMode.value === 'beautify') {
    return $t.value('beautifyResultTab')
  }
  return $t.value('previewResultTab')
})
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
const resultImageSrc = computed(() => {
  const base = String(resultPreview.value || '').trim()
  if (!base) return ''
  const separator = base.includes('?') ? '&' : '?'
  return `${base}${separator}_ri=${resultImageRetryKey.value}`
})
const activeResultImageSrc = computed(() => {
  if (resultTab.value === 'beautify') {
    return beautifyResultPreview.value
  }
  return resultImageSrc.value
})

const normalizeTask = (res) => {
  if (!res || !res.data) return null
  return res.data.task || res.data.reTryonTask || res.data.tryonTask || null
}

const normalizeBool = (value, fallback = false) => {
  if (typeof value === 'boolean') return value
  if (value === undefined || value === null || value === '') return fallback
  const text = String(value).trim().toLowerCase()
  if (['1', 'true', 'yes', 'on'].includes(text)) return true
  if (['0', 'false', 'no', 'off'].includes(text)) return false
  return fallback
}

const normalizeModelMeta = (source = {}) => {
  return {
    supportsBeautify: normalizeBool(source?.supportsBeautify, false),
    beautifyModelKey: String(source?.beautifyModelKey || source?.key || ''),
    beautifyExtraCost: Math.max(0, Number(source?.beautifyExtraCost || source?.beautifyExtraPoints || 0)),
    beautifyRetouchDegree: Number(source?.beautifyRetouchDegree || 70),
    beautifyWhiteningDegree: Number(source?.beautifyWhiteningDegree || 30),
  }
}

const applyModelMeta = (source = {}) => {
  modelMeta.value = normalizeModelMeta(source)
}

const loadModelMetaByTask = async (taskData) => {
  try {
    const configRes = await getTryonConfig()
    if (configRes.code !== 0 || !configRes.data) return
    const sceneType = taskData?.sceneType || draft.value?.sceneType || 'clothes'
    const beautifyModels = parseTryonBeautifyModels(
      configRes.data.tryon_models,
      sceneType,
      0,
      langStore.locale
    )

    const selectedBeautifyModel = beautifyModels[0] || null

    if (selectedBeautifyModel) {
      applyModelMeta({
        ...selectedBeautifyModel,
        supportsBeautify: true,
      })
      return
    }

    applyModelMeta({
      supportsBeautify: false,
      beautifyModelKey: '',
      beautifyExtraCost: 0,
      beautifyRetouchDegree: 70,
      beautifyWhiteningDegree: 30,
    })
  } catch {
    // ignore model meta resolve error
  }
}

const stopPolling = () => {
  if (pollTimer) {
    clearTimeout(pollTimer)
    pollTimer = null
  }
  isGenerating.value = false
}

const stopBeautifyPolling = () => {
  if (beautifyPollTimer) {
    clearTimeout(beautifyPollTimer)
    beautifyPollTimer = null
  }
  beautifyPolling.value = false
}

const clearResultImageRetry = () => {
  if (resultImageRetryTimer) {
    clearTimeout(resultImageRetryTimer)
    resultImageRetryTimer = null
  }
}

const prepareResultPreview = (raw) => {
  const normalized = getUrl(raw)
  resultPreview.value = normalized
  if (normalized) {
    resultImageLoading.value = true
    resultImageRetryKey.value += 1
    clearResultImageRetry()
  } else {
    resultImageLoading.value = false
  }
}

const prepareBeautifyPreview = (raw) => {
  beautifyResultPreview.value = getUrl(raw)
}

const onResultImageLoad = () => {
  if (resultTab.value !== 'tryon') return
  resultImageLoading.value = false
  clearResultImageRetry()
}

const scheduleResultImageRetry = () => {
  if (!resultPreview.value || taskStatus.value !== 'success') return
  clearResultImageRetry()
  resultImageRetryTimer = setTimeout(() => {
    resultImageRetryKey.value += 1
  }, 2200)
}

const onResultImageError = () => {
  if (resultTab.value !== 'tryon') return
  if (taskStatus.value === 'success') {
    resultImageLoading.value = true
    scheduleResultImageRetry()
  }
}

const finishTask = (latestTask, options = {}) => {
  const persistHistory = options.persistHistory !== false
  const clearDraftAfterFinish = options.clearDraft !== false
  task.value = latestTask
  isGenerating.value = false
  progress.value = 100
  if (latestTask?.resultImage) {
    prepareResultPreview(latestTask.resultImage)
  }
  prepareBeautifyPreview(latestTask?.beautifyResult || '')
  if (resultTab.value === 'beautify' && !latestTask?.beautifyResult) {
    resultTab.value = 'tryon'
  }
  if (latestTask?.status === 'failed') {
    errorMessage.value = latestTask.errorMessage || $t.value('operationFailed')
  }
  if (persistHistory) {
    appendTryonLocalHistory({
      ID: latestTask?.ID,
      taskNo: latestTask?.taskNo,
      requestID: draft.value.requestID,
      roomType: draft.value.roomType,
      sceneType: latestTask?.sceneType || draft.value.sceneType,
      status: latestTask?.status,
      sourceImage: latestTask?.sourceImage || draft.value.sourceRemoteUrl,
      templateImage: latestTask?.templateImage || draft.value.templateRemoteUrl,
      resultImage: latestTask?.resultImage || '',
      beautifyStatus: latestTask?.beautifyStatus || '',
      beautifyTaskNo: latestTask?.beautifyTaskNo || '',
      beautifyResult: latestTask?.beautifyResult || '',
      beautifyCost: latestTask?.beautifyCost || 0,
      beautifyRefund: latestTask?.beautifyRefund || 0,
      beautifyError: latestTask?.beautifyError || '',
      costPoints: latestTask?.costPoints || draft.value.modelCost || 0,
      errorMessage: latestTask?.errorMessage || '',
      CreatedAt: latestTask?.CreatedAt,
      CompletedAt: latestTask?.completedAt,
    })
  }
  if (clearDraftAfterFinish) {
    clearTryonDraft()
  }
}

const startPolling = (taskID, options = {}) => {
  const persistHistory = options.persistHistory !== false
  const clearDraftAfterFinish = options.clearDraft !== false
  if (!taskID) {
    isGenerating.value = false
    return
  }
  pollCount = 0
  isGenerating.value = true

  const loop = async () => {
    pollCount += 1
    progress.value = Math.min(95, 10 + pollCount * 3)
    try {
      const res = await findTryonTask({ ID: taskID })
      if (res.code === 0) {
        const latestTask = normalizeTask(res)
        if (latestTask) {
          task.value = latestTask
          if (latestTask.resultImage) {
            prepareResultPreview(latestTask.resultImage)
          }
          prepareBeautifyPreview(latestTask.beautifyResult || '')
          if (latestTask.status && latestTask.status !== 'processing') {
            stopPolling()
            finishTask(latestTask, { persistHistory, clearDraft: clearDraftAfterFinish })
            return
          }
        }
      }
    } catch (e) {
      // 忽略单次轮询错误，继续重试
    }

    if (pollCount >= POLL_MAX_COUNT) {
      stopPolling()
      errorMessage.value = $t.value('tryonPollingTimeout')
      return
    }

    pollTimer = setTimeout(loop, POLL_INTERVAL)
  }

  pollTimer = setTimeout(loop, POLL_INTERVAL)
}

const ensureRemoteImage = async (remoteUrl, localPath, folder = '', uploadType = '') => {
  const remote = String(remoteUrl || '').trim()
  if (remote && !isTempLocalPath(remote)) return remote
  if (!localPath) return ''
  return uploadTryonImage(localPath, folder, uploadType)
}

const getDraftUploadFolder = (role) => {
  const sourceFolder = String(draft.value.sourceUploadFolder || '').trim()
  const templateFolder = String(draft.value.templateUploadFolder || '').trim()
  if (role === 'source' && sourceFolder) {
    return sourceFolder
  }
  if (role === 'template' && templateFolder) {
    return templateFolder
  }

  return getTryonUploadFolder({
    sceneType: draft.value.sceneType || 'clothes',
    operationType: draft.value.operationType || 'tryon',
    role,
  })
}

const getDraftUploadType = (role) => {
  const roomType = String(draft.value.roomType || '').trim().toLowerCase()
  if (role === 'source') {
    return 'person'
  }
  if (roomType === 'shoe') {
    return 'shoe'
  }
  return 'cloth'
}

const createTask = async () => {
  if (isHistoryMode.value) return

  const token = uni.getStorageSync('x-token') || ''
  if (!token) {
    uni.showToast({ title: $t.value('tryonNeedLogin'), icon: 'none' })
    uni.navigateTo({ url: '/pages/user/login' })
    return
  }

  if (!draft.value.sourceLocalPath && !draft.value.sourceRemoteUrl) {
    uni.showToast({ title: $t.value('tryonPickSourceFirst'), icon: 'none' })
    return
  }
  if (!draft.value.templateLocalPath && !draft.value.templateRemoteUrl) {
    uni.showToast({ title: $t.value('tryonPickTemplateFirst'), icon: 'none' })
    return
  }

  uni.showLoading({ title: $t.value('submitting'), mask: true })
  try {
    const sourceImage = await ensureRemoteImage(
      draft.value.sourceRemoteUrl,
      draft.value.sourceLocalPath,
      getDraftUploadFolder('source'),
      getDraftUploadType('source')
    )
    const templateImage = await ensureRemoteImage(
      draft.value.templateRemoteUrl,
      draft.value.templateLocalPath,
      getDraftUploadFolder('template'),
      getDraftUploadType('template')
    )

    if (!sourceImage) {
      throw new Error($t.value('tryonPickSourceFirst'))
    }
    if (!templateImage) {
      throw new Error($t.value('tryonPickTemplateFirst'))
    }

    const res = await createTryonTask({
      requestID: draft.value.requestID,
      sceneType: draft.value.sceneType || 'clothes',
      templatePart: String(draft.value.templatePart || '').trim(),
      sourceImage,
      templateImage,
      modelKey: draft.value.modelKey || '',
      enableRefiner: !!draft.value.enableRefiner,
    })

    const taskData = normalizeTask(res)
    if (!taskData) {
      errorMessage.value = resolveApiMessage(res.msg, 'operationFailed')
      return
    }

    task.value = taskData
    sourcePreview.value = getUrl(taskData.sourceImage || sourceImage)
    templatePreview.value = getUrl(taskData.templateImage || templateImage)
    if (!beautifySupported.value) {
      loadModelMetaByTask(taskData)
    }

    if (taskData.status === 'processing') {
      progress.value = 12
      startPolling(taskData.ID, { persistHistory: true, clearDraft: true })
      return
    }

    finishTask(taskData, { persistHistory: true, clearDraft: true })
  } catch (e) {
    errorMessage.value = resolveApiMessage(e?.message, 'operationFailed')
  } finally {
    uni.hideLoading()
  }
}

const goBack = () => {
  uni.navigateBack({ delta: 1 })
}

const goHistory = () => {
  uni.navigateTo({ url: '/pages/tryon/history' })
}

const previewResult = () => {
  if (!activeResultPreview.value) return
  uni.previewImage({ urls: [activeResultPreview.value] })
}

const previewSource = () => {
  if (!sourcePreview.value) {
    uni.showToast({ title: $t.value('noImagePreview'), icon: 'none' })
    return
  }
  uni.previewImage({ urls: [sourcePreview.value] })
}

const previewTemplate = () => {
  if (!templatePreview.value) {
    uni.showToast({ title: $t.value('noImagePreview'), icon: 'none' })
    return
  }
  uni.previewImage({ urls: [templatePreview.value] })
}

const switchResultTab = (tab) => {
  if (tab === 'beautify' && !beautifyResultPreview.value) {
    uni.showToast({ title: $t.value('tryonBeautifyEmptyHint'), icon: 'none' })
    return
  }
  if (tab === 'tryon' && !resultPreview.value) {
    uni.showToast({ title: $t.value('resultReadyHint'), icon: 'none' })
    return
  }
  resultTab.value = tab === 'beautify' ? 'beautify' : 'tryon'
}

const switchCompareMode = (mode) => {
  if (mode === 'beautify') {
    if (!canCompareBeautify.value) {
      uni.showToast({ title: $t.value('tryonBeautifyCompareEmptyHint'), icon: 'none' })
      return
    }
    compareMode.value = 'beautify'
    return
  }

  if (!canCompareTryon.value) {
    uni.showToast({ title: $t.value('noImagePreview'), icon: 'none' })
    return
  }
  compareMode.value = 'tryon'
}

const startBeautifyPolling = (taskID) => {
  if (!taskID) return

  stopBeautifyPolling()
  beautifyPolling.value = true
  beautifyPollCount = 0

  const loop = async () => {
    beautifyPollCount += 1
    try {
      const res = await findTryonTask({ ID: taskID })
      if (res.code === 0) {
        const latestTask = normalizeTask(res)
        if (latestTask) {
          task.value = latestTask
          if (latestTask.resultImage) {
            prepareResultPreview(latestTask.resultImage)
          }
          prepareBeautifyPreview(latestTask.beautifyResult || '')
          if (latestTask.beautifyStatus && latestTask.beautifyStatus !== 'processing') {
            stopBeautifyPolling()
            finishTask(latestTask, { persistHistory: true, clearDraft: false })
            if (latestTask.beautifyStatus === 'success' && latestTask.beautifyResult) {
              resultTab.value = 'beautify'
            }
            return
          }
        }
      }
    } catch {
      // ignore one-shot error while polling
    }

    if (beautifyPollCount >= 60) {
      stopBeautifyPolling()
      uni.showToast({ title: $t.value('tryonBeautifyPollingTimeout'), icon: 'none' })
      return
    }

    beautifyPollTimer = setTimeout(loop, POLL_INTERVAL)
  }

  beautifyPollTimer = setTimeout(loop, POLL_INTERVAL)
}

const handleBeautify = async () => {
  if (!canUseBeautify.value || isBeautifyBusy.value) {
    if (beautifyStatus.value !== 'disabled') {
      uni.showToast({ title: $t.value('tryonBeautifyUsedHint'), icon: 'none' })
    } else if (!beautifySupported.value) {
      uni.showToast({ title: $t.value('tryonBeautifyUnsupportedHint'), icon: 'none' })
    }
    return
  }

  const taskID = Number(task.value?.ID || 0)
  if (!taskID) {
    uni.showToast({ title: $t.value('tryonTaskNotFound'), icon: 'none' })
    return
  }

  beautifyLoading.value = true
  errorMessage.value = ''
  try {
    const payload = {
      taskID,
      beautifyModelKey: String(modelMeta.value?.beautifyModelKey || draft.value?.beautifyModelKey || ''),
      retouchDegree: Number(modelMeta.value?.beautifyRetouchDegree || draft.value?.beautifyRetouchDegree || 70),
      whiteningDegree: Number(modelMeta.value?.beautifyWhiteningDegree || draft.value?.beautifyWhiteningDegree || 30),
    }
    const res = await applyTryonBeautify(payload)
    const latestTask = normalizeTask(res)
    if (!latestTask) {
      errorMessage.value = resolveApiMessage(res.msg, 'operationFailed')
      return
    }

    task.value = latestTask
    if (latestTask.resultImage) {
      prepareResultPreview(latestTask.resultImage)
    }
    prepareBeautifyPreview(latestTask.beautifyResult || '')
    finishTask(latestTask, { persistHistory: true, clearDraft: false })

    if (latestTask.beautifyStatus === 'processing') {
      startBeautifyPolling(latestTask.ID)
      return
    }

    if (latestTask.beautifyStatus === 'success' && latestTask.beautifyResult) {
      resultTab.value = 'beautify'
      uni.showToast({ title: $t.value('tryonBeautifyDone'), icon: 'none' })
      return
    }

    if (latestTask.beautifyStatus === 'failed') {
      errorMessage.value = latestTask.beautifyError || $t.value('tryonBeautifyFailed')
      uni.showToast({ title: errorMessage.value, icon: 'none' })
    }
  } catch (e) {
    errorMessage.value = resolveApiMessage(e?.message, 'operationFailed')
  } finally {
    beautifyLoading.value = false
  }
}

const openCompare = () => {
  let targetMode = resultTab.value === 'beautify' ? 'beautify' : 'tryon'
  if (targetMode === 'beautify' && !canCompareBeautify.value) {
    targetMode = 'tryon'
  }
  if (targetMode === 'tryon' && !canCompareTryon.value) {
    targetMode = 'beautify'
  }

  if (targetMode === 'beautify' && !canCompareBeautify.value) {
    uni.showToast({ title: $t.value('tryonBeautifyCompareEmptyHint'), icon: 'none' })
    return
  }
  if (targetMode === 'tryon' && !canCompareTryon.value) {
    uni.showToast({ title: $t.value('noImagePreview'), icon: 'none' })
    return
  }

  compareMode.value = targetMode
  comparePercent.value = 50
  compareZoomPercent.value = 100
  compareOffsetX.value = 0
  compareOffsetY.value = 0
  compareDragging.value = false
  comparePinching.value = false
  compareVisible.value = true
  measureCompareStage()
}

const closeCompare = () => {
  compareVisible.value = false
  comparePercent.value = 50
  compareZoomPercent.value = 100
  compareOffsetX.value = 0
  compareOffsetY.value = 0
  compareDragging.value = false
  comparePinching.value = false
}

const onCompareSliderChange = (e) => {
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

const initTaskFromHistory = async (taskID) => {
  const id = Number(taskID || 0)
  if (!id) return

  isHistoryMode.value = true
  uni.showLoading({ title: $t.value('loading'), mask: true })
  try {
    const res = await findTryonTask({ ID: id })
    const taskData = normalizeTask(res)
    if (!taskData) {
      errorMessage.value = resolveApiMessage(res.msg, 'operationFailed')
      return
    }

    task.value = taskData
    sourcePreview.value = getUrl(taskData.sourceImage || draft.value.sourceRemoteUrl || draft.value.sourceLocalPath || '')
    templatePreview.value = getUrl(taskData.templateImage || draft.value.templateRemoteUrl || draft.value.templateLocalPath || '')
    if (taskData.resultImage) {
      prepareResultPreview(taskData.resultImage)
    }
    prepareBeautifyPreview(taskData.beautifyResult || '')
    if (taskData.beautifyResult) {
      resultTab.value = 'beautify'
    }
    await loadModelMetaByTask(taskData)

    if (taskData.status === 'processing') {
      progress.value = 20
      startPolling(taskData.ID, { persistHistory: false, clearDraft: false })
      return
    }

    finishTask(taskData, { persistHistory: false, clearDraft: false })
  } catch (e) {
    errorMessage.value = resolveApiMessage(e?.message, 'operationFailed')
  } finally {
    uni.hideLoading()
  }
}

const ensureAlbumPermission = async () => {
  // #ifdef H5
  return true
  // #endif

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

const downloadResult = async () => {
  if (!activeResultPreview.value) {
    uni.showToast({ title: $t.value('noImagePreview'), icon: 'none' })
    return
  }

  // #ifdef H5
  const anchor = document.createElement('a')
  anchor.href = activeResultPreview.value
  anchor.target = '_blank'
  anchor.rel = 'noopener'
  const filePrefix = resultTab.value === 'beautify' ? 'beautify' : 'tryon'
  anchor.download = `${filePrefix}-${Date.now()}`
  document.body.appendChild(anchor)
  anchor.click()
  document.body.removeChild(anchor)
  uni.showToast({ title: $t.value('downloadAction'), icon: 'none' })
  return
  // #endif

  try {
    uni.showLoading({ title: $t.value('loading'), mask: true })
    const downloadRes = await uni.downloadFile({ url: activeResultPreview.value })
    const filePath = downloadRes?.tempFilePath
    if (!filePath) throw new Error($t.value('downloadFailed'))
    await saveImageToAlbumWithRetry(filePath)
    uni.showToast({ title: $t.value('savedToAlbum'), icon: 'none' })
  } catch (e) {
    if (String(e?.message || '') === 'NO_ALBUM_PERMISSION') {
      uni.showToast({ title: $t.value('saveToAlbumFailed'), icon: 'none' })
      return
    }
    uni.showToast({ title: $t.value('saveToAlbumFailed'), icon: 'none' })
  } finally {
    uni.hideLoading()
  }
}

const goContinue = () => {
  if (draft.value.roomType === 'shoe') {
    uni.switchTab({ url: '/pages/tabBar/shop/shop' })
    return
  }
  uni.switchTab({ url: '/pages/tabBar/index' })
}

onLoad(async (options) => {
  draft.value = getTryonDraft()
  applyModelMeta(draft.value || {})
  sourcePreview.value = draft.value.sourceRemoteUrl ? getUrl(draft.value.sourceRemoteUrl) : draft.value.sourceLocalPath
  templatePreview.value = draft.value.templateRemoteUrl ? getUrl(draft.value.templateRemoteUrl) : draft.value.templateLocalPath
  autoBeautifyRequested.value = String(options?.autoBeautify || options?.beautify || '').trim() === '1'

  const historyTaskID = Number(options?.taskID || 0)
  if (historyTaskID > 0) {
    await initTaskFromHistory(historyTaskID)
    if (autoBeautifyRequested.value) {
      setTimeout(() => {
        handleBeautify()
      }, 120)
    }
    return
  }

  createTask()
})

onUnload(() => {
  stopPolling()
  stopBeautifyPolling()
  clearResultImageRetry()
})
</script>

<style lang="scss">
page {
  background: #f4f7fb;
}

.generate-page {
  min-height: 100vh;
  padding: calc(var(--status-bar-height, 0px) + 14rpx) 20rpx 24rpx;
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
  min-width: 90rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-title {
  font-size: 30rpx;
  font-weight: 700;
}

.nav-link {
  font-size: 24rpx;
  color: rgba(15,23,42,0.62);
}

.preview-panel {
  margin-top: 16rpx;
  display: flex;
  gap: 12rpx;
}

.generate-notice {
  margin-top: 12rpx;
  padding: 14rpx 16rpx;
  border-radius: 12rpx;
  border: 1rpx solid rgba(14, 165, 233, 0.18);
  background: rgba(239, 246, 255, 0.9);
  color: rgba(15, 23, 42, 0.65);
  font-size: 22rpx;
  line-height: 34rpx;
}

.preview-col {
  flex: 1;
  background: rgba(255,255,255,0.95);
  border: 1rpx solid rgba(15,23,42,0.08);
  border-radius: 14rpx;
  padding: 12rpx;
  box-shadow: 0 12rpx 26rpx rgba(15, 23, 42, 0.06);
}

.preview-image-wrap {
  margin-top: 8rpx;
  width: 100%;
  height: 280rpx;
  border-radius: 10rpx;
  overflow: hidden;
  background: #f8fafc;
}

.preview-label {
  font-size: 22rpx;
  color: rgba(15,23,42,0.58);
}

.preview-image {
  width: 100%;
  height: 100%;
}

.status-card,
.result-card {
  margin-top: 16rpx;
  background: rgba(255,255,255,0.95);
  border: 1rpx solid rgba(15,23,42,0.08);
  border-radius: 14rpx;
  padding: 16rpx;
  box-shadow: 0 12rpx 26rpx rgba(15, 23, 42, 0.06);
}

.status-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-title,
.result-title {
  font-size: 26rpx;
  font-weight: 600;
}

.status-pill {
  padding: 6rpx 14rpx;
  border-radius: 999rpx;
  font-size: 20rpx;
  background: rgba(15,23,42,0.08);
}

.status-pill.processing {
  background: rgba(255,205,105,0.24);
}

.status-pill.success {
  background: rgba(86,199,102,0.26);
}

.status-pill.failed {
  background: rgba(255,107,107,0.24);
}

.task-no {
  margin-top: 8rpx;
  font-size: 22rpx;
  color: rgba(15,23,42,0.62);
}

.task-cost {
  margin-top: 8rpx;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.task-cost-line {
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.62);
}

.progress-wrap {
  margin-top: 14rpx;
}

.progress-track {
  width: 100%;
  height: 12rpx;
  border-radius: 999rpx;
  background: rgba(15,23,42,0.08);
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
}

.progress-text {
  margin-top: 8rpx;
  color: rgba(15,23,42,0.62);
  font-size: 22rpx;
}

.error-msg {
  margin-top: 10rpx;
  color: #ff9c9c;
  font-size: 22rpx;
}

.status-actions {
  margin-top: 12rpx;
  display: flex;
  flex-direction: column;
  gap: 10rpx;
}

.status-action-tip {
  display: block;
  margin-bottom: 8rpx;
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.56);
}

.status-action-btn {
  height: 60rpx;
  border-radius: 999rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(15, 23, 42, 0.06);
  color: #0f172a;
  font-size: 22rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.12);
}

.status-action-btn.highlight {
  color: #ffffff;
  border-color: transparent;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  box-shadow: 0 10rpx 22rpx rgba(37, 99, 235, 0.26);
}

.status-action-btn.disabled {
  opacity: 0.55;
}

.result-save-tip {
  display: block;
  margin-bottom: 8rpx;
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.56);
}

.result-tabs {
  display: flex;
  gap: 10rpx;
}

.result-tab {
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

.result-tab.active {
  color: #ffffff;
  border-color: transparent;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
}

.result-tab.disabled {
  opacity: 0.5;
}

.result-box {
  margin-top: 12rpx;
  border-radius: 12rpx;
  overflow: hidden;
  position: relative;
  height: 680rpx;
  background: #f8fafc;
}

.result-image {
  width: 100%;
  height: 100%;
}

.result-loading {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(15, 23, 42, 0.58);
  font-size: 22rpx;
  background: rgba(248, 250, 252, 0.72);
}

.result-empty {
  margin-top: 12rpx;
  height: 220rpx;
  border-radius: 12rpx;
  border: 1rpx dashed rgba(15,23,42,0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(15,23,42,0.55);
  font-size: 22rpx;
}

.action-row {
  margin-top: 18rpx;
  display: flex;
  gap: 12rpx;
}

.action-btn {
  flex: 1;
  height: 84rpx;
  border-radius: 999rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26rpx;
  font-weight: 600;
  color: #fff;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  box-shadow: 0 12rpx 26rpx rgba(37, 99, 235, 0.28);
}

.action-btn.secondary {
  background: rgba(255,255,255,0.92);
  color: #0f172a;
  border: 1rpx solid rgba(15,23,42,0.1);
  box-shadow: none;
}

.compare-mask {
  position: fixed;
  inset: 0;
  z-index: 1200;
  background: rgba(15, 23, 42, 0.58);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24rpx;
}

.compare-panel {
  width: calc(100vw - 32rpx);
  max-width: 980rpx;
  border-radius: 18rpx;
  background: #ffffff;
  overflow: hidden;
}

.compare-head {
  height: 88rpx;
  padding: 0 20rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1rpx solid rgba(15, 23, 42, 0.08);
}

.compare-title {
  font-size: 28rpx;
  font-weight: 700;
  color: #0f172a;
}

.compare-close {
  width: 54rpx;
  height: 54rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(15, 23, 42, 0.06);
}

.compare-mode-row {
  margin: 12rpx 16rpx 0;
  display: flex;
  gap: 10rpx;
}

.compare-mode-tab {
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

.compare-mode-tab.active {
  color: #ffffff;
  border-color: transparent;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
}

.compare-mode-tab.disabled {
  opacity: 0.5;
}

.compare-stage {
  margin: 12rpx 16rpx 0;
  height: 78vh;
  min-height: 700rpx;
  max-height: 1120rpx;
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
  padding: 16rpx 20rpx 20rpx;
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
</style>
