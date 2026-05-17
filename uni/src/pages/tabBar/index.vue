<template>
  <view class="room-page">
    <view class="room-header">
      <text class="room-title">{{ roomTitle }}</text>
      <view class="room-header-action" @tap="showTutorialPopup = true">
        <text>{{ $t('tryonTutorialAction') }}</text>
      </view>
    </view>

    <view class="room-body">
      <view class="room-left" @tap="openUploadDrawer('person')">
        <LazyImage v-if="personPreview" class="room-preview" :src="personPreview" mode="aspectFit" />
        <view class="slot-actions" v-if="personPreview">
          <view class="slot-action-btn" @tap.stop="previewRoomImage('person')">
            <uni-icons type="search" size="16" color="#2563eb" />
          </view>
          <view class="slot-action-btn danger" @tap.stop="clearRoomImage('person')">
            <uni-icons type="trash" size="16" color="#dc2626" />
          </view>
        </view>
        <view v-if="personPreview && personSizeBytes > 0" class="preview-size-mask">{{ formatPreviewSize(personSizeBytes) }}</view>
        <view v-else class="room-upload-empty">
          <uni-icons type="camera" size="26" color="rgba(15,23,42,0.45)" />
          <text class="room-upload-text">{{ $t('uploadModelImage') }}</text>
        </view>
      </view>
      <view class="room-right">
        <view class="cloth-slot" @tap="openUploadDrawer('upper')">
          <LazyImage v-if="upperPreview" class="cloth-preview" :src="upperPreview" mode="aspectFit" />
          <view class="slot-actions" v-if="upperPreview">
            <view class="slot-action-btn" @tap.stop="previewRoomImage('upper')">
              <uni-icons type="search" size="16" color="#2563eb" />
            </view>
            <view class="slot-action-btn danger" @tap.stop="clearRoomImage('upper')">
              <uni-icons type="trash" size="16" color="#dc2626" />
            </view>
          </view>
          <view v-if="upperPreview && upperSizeBytes > 0" class="preview-size-mask">{{ formatPreviewSize(upperSizeBytes) }}</view>
          <text v-else class="cloth-text">{{ $t('uploadUpperImage') }}</text>
        </view>
        <view class="cloth-slot" @tap="openUploadDrawer('lower')">
          <LazyImage v-if="lowerPreview" class="cloth-preview" :src="lowerPreview" mode="aspectFit" />
          <view class="slot-actions" v-if="lowerPreview">
            <view class="slot-action-btn" @tap.stop="previewRoomImage('lower')">
              <uni-icons type="search" size="16" color="#2563eb" />
            </view>
            <view class="slot-action-btn danger" @tap.stop="clearRoomImage('lower')">
              <uni-icons type="trash" size="16" color="#dc2626" />
            </view>
          </view>
          <view v-if="lowerPreview && lowerSizeBytes > 0" class="preview-size-mask">{{ formatPreviewSize(lowerSizeBytes) }}</view>
          <text v-else class="cloth-text">{{ $t('uploadLowerImage') }}</text>
        </view>
        <view class="cloth-slot suit-slot" @tap="goClothesPage">
          <text class="cloth-text">{{ $t('suitSet') }}</text>
          <text class="cloth-tip">{{ $t('goClothesPageTip') }}</text>
        </view>
      </view>
    </view>

    <view class="image-size-rule-card">
      <text class="image-size-rule-text">{{ $t('tryonImageSizeRuleHint') }}</text>
    </view>

    <view class="refiner-card" :class="{ disabled: !currentModelSupportsRefiner }">
      <view class="refiner-left">
        <text class="refiner-label">{{ currentRefinerLabel }}</text>
        <text class="refiner-hint" v-if="currentModelSupportsRefiner">{{ currentRefinerHint }}</text>
        <text class="refiner-extra-cost" v-if="currentModelSupportsRefiner">
          {{ currentRefinerExtraCost > 0 ? $t('tryonRefinerExtraCostHint').replace('{cost}', String(currentRefinerExtraCost)) : $t('tryonFeatureFreeHint') }}
        </text>
        <text class="refiner-hint" v-if="!currentModelSupportsRefiner">{{ $t('tryonRefinerUnsupportedHint') }}</text>
      </view>
      <view class="refiner-right">
        <switch
          class="refiner-switch"
          :checked="refinerEnabled"
          :disabled="!currentModelSupportsRefiner"
          color="#2563eb"
          @change="onRefinerSwitchChange"
        />
        <view class="refiner-help" @tap="openFeatureHelp">?</view>
      </view>
    </view>

    <view class="model-card">
      <view class="model-left" @tap="showModelPopup = true">
        <text class="model-label">{{ $t('currentModelLabel') }}</text>
        <view class="model-info">
          <text class="model-value">{{ currentModel.name }}</text>
          <text class="model-cost-hint">{{ $t('pointsCostEach').replace('{cost}', String(currentBaseCost)) }}</text>
          <text
            class="model-cost-total"
            v-if="showModelCostBreakdown"
          >
            {{ currentCostFormulaText }}
          </text>
          <text class="model-desc-hint" v-if="currentModelDesc">{{ currentModelDesc }}</text>
        </view>
      </view>
      <view class="model-switch" @tap="showModelPopup = true">
        <text>{{ $t('switchAction') }}</text>
      </view>
    </view>

    <view class="upper-only-split-card" v-if="showUpperOnlySplitSelector">
      <view class="upper-only-split-options">
        <view
          class="upper-only-split-option"
          :class="{ active: upperOnlySplitMode === 'full_outfit' }"
          @tap="selectUpperOnlySplitMode('full_outfit')"
        >
          <text>{{ $t('tryonParsingOptionDress') }}</text>
        </view>
        <view
          class="upper-only-split-option"
          :class="{ active: upperOnlySplitMode === 'upper_only' }"
          @tap="selectUpperOnlySplitMode('upper_only')"
        >
          <text>{{ $t('tryonParsingOptionUpper') }}</text>
        </view>
      </view>
      <view class="upper-only-split-help" @tap="showUpperOnlySplitHelpPopup = true">?</view>
    </view>

    <view class="generate-btn" @tap="goGenerate">
      <text class="generate-text">{{ $t('generateWithCost').replace('{cost}', String(currentCost)) }}</text>
    </view>

    <announcement-marquee
      :enabled="announcementEnabled"
      :text="announcementText"
      :text-color="announcementTextColor"
      :speed="announcementSpeed"
    />

    <view class="footer-tips">
      <text>{{ $t('tryonModelHint') }}</text>
    </view>

    <view class="popup-mask" v-if="showModelPopup" @tap="showModelPopup = false">
      <view class="popup-panel" @tap.stop>
        <view class="popup-title">{{ $t('selectModel') }}</view>
        <scroll-view class="popup-list" scroll-y>
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

    <view class="popup-mask" v-if="showRefinerHelpPopup" @tap="showRefinerHelpPopup = false">
      <view class="popup-panel popup-panel-help" @tap.stop>
        <view class="popup-title">{{ currentHelpTitle }}</view>
        <view class="refiner-help-content">{{ currentHelpContent }}</view>
        <view class="refiner-help-actions">
          <view class="model-switch" @tap="showRefinerHelpPopup = false">
            <text>{{ $t('doneText') }}</text>
          </view>
        </view>
      </view>
    </view>

    <view class="popup-mask" v-if="showUpperOnlySplitHelpPopup" @tap="showUpperOnlySplitHelpPopup = false">
      <view class="popup-panel popup-panel-help" @tap.stop>
        <view class="refiner-help-content">{{ upperOnlySplitGuideText }}</view>
        <view class="refiner-help-actions">
          <view class="model-switch" @tap="showUpperOnlySplitHelpPopup = false">
            <text>{{ $t('doneText') }}</text>
          </view>
        </view>
      </view>
    </view>

    <view class="popup-mask" v-if="showTutorialPopup" @tap="showTutorialPopup = false" @touchmove.stop>
      <view class="popup-panel popup-panel-help popup-panel-tutorial" @tap.stop @touchmove.stop>
        <view class="tutorial-popup-head">
          <view class="popup-title tutorial-popup-title">{{ $t('tryonTutorialTitle') }}</view>
          <view class="tutorial-popup-close" @tap="showTutorialPopup = false">
            <text>×</text>
          </view>
        </view>
        <scroll-view class="tutorial-scroll" scroll-y @touchmove.stop>
          <view class="tutorial-list">
            <view class="tutorial-item" v-for="(item, index) in tutorialItems" :key="`tutorial-${index}`">
              <text class="tutorial-index">{{ index + 1 }}.</text>
              <text class="tutorial-text">{{ item }}</text>
            </view>
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
            <text class="upload-highlight-tip">{{ $t('uploadPersonDrawerStrongTip') }}</text>
            <text class="upload-highlight-tip">{{ $t('uploadPersonSidePoseHint') }}</text>

            <view class="example-grid three">
              <view
                v-for="item in personSidePoseExamples"
                :key="item.url"
                class="example-card"
                @tap="applyRemoteExample(item)"
              >
                <LazyImage class="example-image" :src="getExamplePreview(item.url)" mode="aspectFit" @error="handleExampleImageError(item.url)" />
                <text class="example-label">{{ item.label }}</text>
              </view>
            </view>

            <text class="upload-center-tip">{{ $t('uploadSingleFullBodyTip') }}</text>

            <view class="example-grid two">
              <view
                v-for="item in personGoodExamples"
                :key="item.url"
                class="example-card"
                @tap="applyRemoteExample(item)"
              >
                <LazyImage class="example-image" :src="getExamplePreview(item.url)" mode="aspectFit" @error="handleExampleImageError(item.url)" />
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
                <LazyImage class="example-image" :src="getExamplePreview(item.url)" mode="aspectFit" @error="handleExampleImageError(item.url)" />
                <text class="example-label">{{ item.label }}</text>
              </view>
            </view>
          </view>

          <view class="upload-section" v-if="isPersonDrawer && drawerTab === 'myModel'">
            <view v-if="myModelList.length === 0" class="drawer-empty">
              <text>{{ $t('uploadDrawerEmptyPersonToMyModel') }}</text>
            </view>
            <view v-else class="example-grid two">
              <view
                v-for="item in myModelList"
                :key="item.id"
                class="example-card"
                @tap="applyMyModel(item)"
              >
                <LazyImage class="example-image" :src="getExamplePreview(item.url)" mode="aspectFit" @error="handleExampleImageError(item.url)" />
                <text class="example-label">{{ item.name || $t('unnamedModel') }}</text>
              </view>
            </view>
          </view>

          <view class="upload-section" v-if="isPersonDrawer && drawerTab === 'official'">
            <text class="group-title">{{ $t('clothesGenderFemale') }}</text>
            <view class="example-grid two">
              <view
                v-for="item in officialFemaleExamples"
                :key="item.url"
                class="example-card"
                @tap="applyRemoteExample(item)"
              >
                <LazyImage class="example-image" :src="getExamplePreview(item.url)" mode="aspectFit" @error="handleExampleImageError(item.url)" />
                <text class="example-label">{{ item.label }}</text>
              </view>
            </view>

            <text class="group-title">{{ $t('clothesGenderMale') }}</text>
            <view class="example-grid two">
              <view
                v-for="item in officialMaleExamples"
                :key="item.url"
                class="example-card"
                @tap="applyRemoteExample(item)"
              >
                <LazyImage class="example-image" :src="getExamplePreview(item.url)" mode="aspectFit" @error="handleExampleImageError(item.url)" />
                <text class="example-label">{{ item.label }}</text>
              </view>
            </view>
          </view>

          <view class="upload-section" v-if="isClothDrawer && drawerTab === 'custom'">
            <view class="upload-mode-actions">
              <view class="upload-main-btn" @tap="uploadDrawerChooseImage('original')">{{ $t('uploadModeOriginal') }}</view>
              <view class="upload-main-btn" @tap="uploadDrawerChooseImage('crop')">{{ $t('uploadModeCrop') }}</view>
              <view class="upload-main-btn" @tap="uploadDrawerChooseImage('compress')">{{ $t('uploadModeCompress') }}</view>
            </view>
            <text class="upload-highlight-tip" v-if="uploadTarget === 'upper'">{{ $t('uploadUpperDrawerStrongTip') }}</text>
            <text class="upload-highlight-tip" v-if="uploadTarget === 'upper'">{{ $t('uploadUpperPersonWearHint') }}</text>

            <view class="example-grid three" v-if="uploadTarget === 'upper'">
              <view
                v-for="item in upperPersonWearExamples"
                :key="item.url"
                class="example-card"
                @tap="applyRemoteExample(item)"
              >
                <LazyImage class="example-image" :src="getExamplePreview(item.url)" mode="aspectFit" @error="handleExampleImageError(item.url)" />
                <text class="example-label">{{ item.label }}</text>
              </view>
            </view>

            <text class="upload-center-tip">{{ $t(uploadTarget === 'lower' ? 'uploadLowerOptionalTip' : 'uploadUpperRequiredTip') }}</text>

            <view class="example-grid three">
              <view
                v-for="item in clothGoodExamples"
                :key="item.url"
                class="example-card"
                @tap="applyRemoteExample(item)"
              >
                <LazyImage class="example-image" :src="getExamplePreview(item.url)" mode="aspectFit" @error="handleExampleImageError(item.url)" />
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
                v-for="item in clothBadExamples"
                :key="item.url"
                class="example-card"
                @tap="applyRemoteExample(item)"
              >
                <LazyImage class="example-image" :src="getExamplePreview(item.url)" mode="aspectFit" @error="handleExampleImageError(item.url)" />
                <text class="example-label">{{ item.label }}</text>
              </view>
            </view>
          </view>

          <view class="upload-section" v-if="isClothDrawer && drawerTab === 'myCloset'">
            <view v-if="filteredMyClothList.length === 0" class="drawer-empty">
              <text>{{ $t('uploadDrawerEmptyClothesToMyCloset') }}</text>
            </view>
            <view v-else class="example-grid three">
              <view
                v-for="item in filteredMyClothList"
                :key="item.id"
                class="example-card"
                @tap="applyMyCloth(item)"
              >
                <LazyImage class="example-image" :src="getExamplePreview(item.url)" mode="aspectFit" @error="handleExampleImageError(item.url)" />
                <text class="example-label">{{ item.name || myClothDisplayName(item.category) }}</text>
              </view>
            </view>
          </view>

          <view class="upload-section" v-if="isClothDrawer && drawerTab === 'recommended'">
            <view v-for="group in clothRecommendedGroups" :key="group.key" class="recommend-group">
              <text class="group-title">{{ group.title }}</text>
              <view class="example-grid three">
                <view
                  v-for="item in group.items"
                  :key="item.url"
                  class="example-card"
                  @tap="applyRemoteExample(item)"
                >
                  <LazyImage class="example-image" :src="getExamplePreview(item.url)" mode="aspectFit" @error="handleExampleImageError(item.url)" />
                  <text class="example-label">{{ item.label }}</text>
                </view>
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
            <LazyImage class="crop-editor-image" :src="cropSourcePath" mode="scaleToFill" />
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
          <view class="upload-drawer-btn" :class="{ disabled: cropConfirming }" @tap="confirmCropEditor">{{ cropConfirming ? $t('loading') : $t('doneText') }}</view>
        </view>
      </view>
    </view>

    <canvas
      canvas-id="tryonCropCanvas"
      id="tryonCropCanvas"
      class="crop-canvas-hidden"
      :style="{ width: `${cropCanvasSize.width}px`, height: `${cropCanvasSize.height}px` }"
    ></canvas>

    <lang-switch v-model="showLangPicker" />
  </view>
</template>

<script setup>
import { computed, nextTick, ref, watch } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import AnnouncementMarquee from '@/components/announcement-marquee/announcement-marquee.vue'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'
import { useUserStore } from '@/pinia/modules/user.js'
import { getTryonConfig, getDefaultDomain, getAnnouncementConfig } from '@/api/sysConfig.js'
import { getMyTryonModelList, getMyTryonClothList } from '@/api/tryonTask.js'
import { getUrl } from '@/utils/url.js'
import { localText, resolveApiMessage } from '@/utils/i18n.js'
import LazyImage from '@/components/lazy-image/lazy-image.vue'
import langSwitch from '@/components/lang-switch/lang-switch.vue'
import {
  createTryonRequestId,
  saveTryonDraft,
  getSelectedClothes,
  clearSelectedClothes,
  getSelectedTryonModel,
  clearSelectedTryonModel,
  parseTryonBeautifyModels,
  parseTryonModels,
  uploadTryonImage,
} from '@/utils/tryon.js'

const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
const userStore = useUserStore()
const $t = computed(() => langStore.$t)
const roomTitle = computed(() => appConfigStore.appName || $t.value('tryonRoom'))
const tutorialItems = computed(() => {
  return [
    $t.value('tryonTutorialWebBrowserTip'),
    $t.value('tryonTutorialItem1'),
    $t.value('tryonTutorialItem2'),
    $t.value('tryonTutorialItem3'),
    $t.value('tryonTutorialItem4'),
    $t.value('tryonTutorialItem5'),
    $t.value('tryonTutorialItem6'),
    $t.value('tryonTutorialItem7'),
    $t.value('tryonTutorialItem8'),
    $t.value('tryonTutorialItem9'),
    $t.value('tryonTutorialItem10'),
  ]
})

const showModelPopup = ref(false)
const showRefinerHelpPopup = ref(false)
const showUpperOnlySplitHelpPopup = ref(false)
const showTutorialPopup = ref(false)
const showLangPicker = ref(false)
const showUploadDrawer = ref(false)
const refinerEnabled = ref(false)
const upperOnlySplitMode = ref('')
const uploadTarget = ref('')
const drawerTab = ref('custom')
const selectedModelKey = ref('')
const modelList = ref([])
const myModelList = ref([])
const myClothList = ref([])
const defaultExternalDomain = ref('')
const exampleImageFailIndex = ref({})
const announcementConfig = ref({
  announcement_enabled: 'false',
  announcement_content: '',
  announcement_text_color: '#ff6600',
  announcement_speed: '50',
})
const tryonConfig = ref({
  tryon_cost_points: '1',
  tryon_models: '',
})

const personLocal = ref('')
const personRemote = ref('')
const upperLocal = ref('')
const upperRemote = ref('')
const lowerLocal = ref('')
const lowerRemote = ref('')
const personSizeBytes = ref(0)
const upperSizeBytes = ref(0)
const lowerSizeBytes = ref(0)
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

const personPreview = computed(() => personRemote.value ? getUrl(personRemote.value) : personLocal.value)
const upperPreview = computed(() => upperRemote.value ? getUrl(upperRemote.value) : upperLocal.value)
const lowerPreview = computed(() => lowerRemote.value ? getUrl(lowerRemote.value) : lowerLocal.value)

const activeSceneType = computed(() => 'clothes')

const parseBoolFlag = (value, fallback = false) => {
  if (typeof value === 'boolean') return value
  if (value === undefined || value === null || value === '') return fallback
  const text = String(value).trim().toLowerCase()
  if (['1', 'true', 'yes', 'on'].includes(text)) return true
  if (['0', 'false', 'no', 'off'].includes(text)) return false
  return fallback
}

const normalizeOptionalBeautifyDegree = (value) => {
  const numeric = Number(value)
  if (!Number.isFinite(numeric) || numeric <= 0) {
    return 0
  }
  return numeric
}

const supportsRefinerByModel = (model = {}) => {
  const refinerModelKey = String(model?.refinerModelKey || '').trim()
  if (refinerModelKey) {
    return true
  }
  const explicit = model?.supportsRefiner
  if (explicit !== undefined && explicit !== null && explicit !== '') {
    return parseBoolFlag(explicit, false)
  }
  return false
}

const supportsParsingByModel = (model = {}) => {
  const usage = String(model?.modelUsage || 'tryon').trim().toLowerCase()
  if (usage && usage !== 'tryon') {
    return false
  }
  const parsingModelKey = String(model?.parsingModelKey || '').trim()
  if (parsingModelKey) {
    return true
  }
  const provider = String(model?.provider || '').trim().toLowerCase()
  const modelName = String(model?.model || '').trim().toLowerCase()
  return (provider.includes('aliyun') || provider.includes('dashscope')) && (modelName === 'aitryon' || modelName === 'aitryon-plus')
}

const currentModel = computed(() => {
  const selected = modelList.value.find(v => v.key === selectedModelKey.value)
  return selected || modelList.value[0] || { key: 'aitryon', name: 'aitryon', cost: Number(tryonConfig.value.tryon_cost_points || 1), desc: {} }
})

const currentBaseCost = computed(() => Number(currentModel.value.cost || 1))
const currentModelDesc = computed(() => currentModel.value.descText || localText(currentModel.value.desc, langStore.locale) || '')
const currentRefinerLabel = computed(() => currentModel.value?.refinerName || $t.value('tryonRefinerLabel'))
const currentModelSupportsRefiner = computed(() => supportsRefinerByModel(currentModel.value))
const currentModelSupportsParsing = computed(() => supportsParsingByModel(currentModel.value))
const currentRefinerExtraCost = computed(() => {
  if (!currentModelSupportsRefiner.value) {
    return 0
  }
  return Math.max(0, Number(currentModel.value.refinerExtraCost || 0))
})
const currentParsingExtraCost = computed(() => {
  if (!currentModelSupportsParsing.value) {
    return 0
  }
  return Math.max(0, Number(currentModel.value.parsingExtraCost || 0))
})
const showUpperOnlySplitSelector = computed(() => hasUpperTemplate() && !hasLowerTemplate())
const upperOnlySplitGuideText = computed(() => $t.value('tryonParsingRecommendContent'))
const parsingEnabledBySelection = computed(() => {
  const hasUpper = hasUpperTemplate()
  const hasLower = hasLowerTemplate()
  if (hasLower && !hasUpper) {
    return true
  }
  if (hasUpper && !hasLower) {
    return upperOnlySplitMode.value === 'upper_only'
  }
  return false
})
const currentCost = computed(() => {
  let total = currentBaseCost.value
  if (currentModelSupportsRefiner.value && refinerEnabled.value) {
    total += currentRefinerExtraCost.value
  }
  if (parsingEnabledBySelection.value) {
    total += currentParsingExtraCost.value
  }
  return total
})
const showRefinerExtraCostInFormula = computed(() => {
  return currentModelSupportsRefiner.value && refinerEnabled.value && currentRefinerExtraCost.value > 0
})
const showParsingExtraCostInFormula = computed(() => {
  return parsingEnabledBySelection.value && currentParsingExtraCost.value > 0
})
const showModelCostBreakdown = computed(() => {
  return showRefinerExtraCostInFormula.value || showParsingExtraCostInFormula.value
})
const currentCostFormulaText = computed(() => {
  const parts = [$t.value('pointsCostEach').replace('{cost}', String(currentBaseCost.value))]
  if (showRefinerExtraCostInFormula.value) {
    parts.push($t.value('pointsCostEach').replace('{cost}', String(currentRefinerExtraCost.value)))
  }
  if (showParsingExtraCostInFormula.value) {
    parts.push($t.value('pointsCostEach').replace('{cost}', String(currentParsingExtraCost.value)))
  }
  return `${parts.join(' + ')} = ${$t.value('pointsCostEach').replace('{cost}', String(currentCost.value))}`
})
const currentRefinerDesc = computed(() => {
  const modelDesc = currentModel.value?.refinerDescText || localText(currentModel.value?.refinerDesc, langStore.locale) || ''
  if (modelDesc) {
    return modelDesc
  }
  return $t.value('tryonRefinerHelpFallback')
})
const currentRefinerHint = computed(() => currentRefinerDesc.value || $t.value('tryonRefinerHint'))
const currentHelpTitle = computed(() => currentRefinerLabel.value || $t.value('tryonRefinerHelpTitle'))
const currentHelpContent = computed(() => currentRefinerDesc.value || $t.value('tryonRefinerHelpFallback'))
const cropRatioLabel = computed(() => `${cropRatio.value.width}:${cropRatio.value.height}`)
const cropRatioOptions = [
  { key: '1:1', label: '1:1', width: 1, height: 1 },
  { key: '3:4', label: '3:4', width: 3, height: 4 },
  { key: '4:3', label: '4:3', width: 4, height: 3 },
  { key: '9:16', label: '9:16', width: 9, height: 16 },
  { key: '16:9', label: '16:9', width: 16, height: 9 },
]

const CROP_CANVAS_ID = 'tryonCropCanvas'
const CROP_MIN_EDGE_SIZE = 48

const setImageSizeForTarget = (target, bytes = 0) => {
  const safeSize = Number.isFinite(Number(bytes)) && Number(bytes) > 0 ? Math.round(Number(bytes)) : 0
  if (target === 'person') personSizeBytes.value = safeSize
  if (target === 'upper') upperSizeBytes.value = safeSize
  if (target === 'lower') lowerSizeBytes.value = safeSize
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

const hasUpperTemplate = () => !!(upperRemote.value || upperLocal.value)
const hasLowerTemplate = () => !!(lowerRemote.value || lowerLocal.value)

const resolveTemplatePartBySelection = () => {
  const hasUpper = hasUpperTemplate()
  const hasLower = hasLowerTemplate()
  if (hasLower && !hasUpper) return 'lower'
  if (hasUpper && !hasLower) {
    if (upperOnlySplitMode.value === 'full_outfit') return 'dress'
    if (upperOnlySplitMode.value === 'upper_only') return 'upper'
  }
  return ''
}

const resolveParsingPartsBySelection = () => {
  const templatePart = resolveTemplatePartBySelection()
  if (templatePart === 'upper' || templatePart === 'lower') {
    return [templatePart]
  }
  return []
}

const resolveTryonPointBalance = () => {
  const storeInfo = userStore?.userInfo || {}
  const cacheInfo = uni.getStorageSync('userInfo') || {}
  const storeValue = Number(storeInfo.tryonPoint ?? storeInfo.point)
  if (Number.isFinite(storeValue)) {
    return Math.max(0, storeValue)
  }
  const cacheValue = Number(cacheInfo.tryonPoint ?? cacheInfo.point)
  if (Number.isFinite(cacheValue)) {
    return Math.max(0, cacheValue)
  }
  return 0
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
  const normalizedUrl = String(getUrl(rawUrl) || '').trim()
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

const getRemoteValueByTarget = (target) => {
  if (target === 'person') return personRemote.value
  if (target === 'upper') return upperRemote.value
  if (target === 'lower') return lowerRemote.value
  return ''
}

const applyRemoteSizeForTarget = async (target, rawUrl) => {
  const normalizedUrl = String(getUrl(rawUrl) || '').trim()
  if (!target || !normalizedUrl) return

  const size = await resolveRemoteImageSize(normalizedUrl)
  if (size <= 0) return

  const currentTargetUrl = String(getUrl(getRemoteValueByTarget(target)) || '').trim()
  if (!currentTargetUrl || currentTargetUrl !== normalizedUrl) return
  setImageSizeForTarget(target, size)
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

const parseJsonSafe = (raw) => {
  try {
    return JSON.parse(raw)
  } catch {
    return null
  }
}

const resolveI18nText = (value) => {
  if (value === null || value === undefined) return ''

  if (typeof value === 'string') {
    const text = value.trim()
    if (!text) return ''

    if (text.startsWith('{') || text.startsWith('[')) {
      const parsed = parseJsonSafe(text)
      if (parsed) {
        return resolveI18nText(parsed)
      }
    }

    return text
  }

  if (Array.isArray(value)) {
    return resolveI18nText(value[0])
  }

  if (typeof value === 'object') {
    return localText(value, langStore.locale)
  }

  return String(value)
}

const announcementEnabled = computed(() => {
  const value = String(announcementConfig.value.announcement_enabled || '').trim().toLowerCase()
  return ['1', 'true', 'yes', 'on'].includes(value)
})
const announcementText = computed(() => resolveI18nText(announcementConfig.value.announcement_content))
const announcementTextColor = computed(() => String(announcementConfig.value.announcement_text_color || '#ff6600').trim() || '#ff6600')
const announcementSpeed = computed(() => {
  const speed = Number.parseInt(String(announcementConfig.value.announcement_speed || '50'), 10)
  if (!Number.isFinite(speed) || speed <= 0) return 50
  return speed
})

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
  if (currentIndex >= candidates.length - 1) {
    console.warn('[tryon-example-image] all candidates failed', { key, candidates })
    return
  }

  console.warn('[tryon-example-image] candidate failed, switching', {
    key,
    failedUrl: candidates[currentIndex],
    nextUrl: candidates[currentIndex + 1],
  })

  exampleImageFailIndex.value = {
    ...exampleImageFailIndex.value,
    [key]: currentIndex + 1,
  }
}

const personGoodExamples = computed(() => {
  return buildIndexedExampleItems(['model-baby.jpg', 'secondcan.jpg', 'thridcan.jpg', 'fourthcan.jpg'], 'tryonExampleFullBodyPrefix')
})

const personSidePoseExamples = computed(() => {
  return buildIndexedExampleItems(['model-ertong2.jpg', 'model-man.jpg', 'model-women2.jpg'], 'tryonExampleSidePosePrefix')
})

const personBadExamples = computed(() => {
  return buildIndexedExampleItems(['firstcant.jpg', 'fourthcant.jpg'], 'tryonExampleBadPrefix')
})

const officialMaleExamples = computed(() => {
  return buildIndexedExampleItems(['secondcan.jpg', 'fourthcan.jpg', 'model-man2.jpg', 'model-baby2.jpg'], 'tryonExampleMaleModelPrefix')
})

const officialFemaleExamples = computed(() => {
  return buildIndexedExampleItems(['firstcan.jpg', 'thridcan.jpg', 'model-sex5.jpg', 'model-sex.jpg'], 'tryonExampleFemaleModelPrefix')
})

const upperPersonWearExamples = computed(() => {
  return buildIndexedExampleItems(['model-ertong.jpg', 'model-man.jpg', 'model-sex2.jpg'], 'tryonExamplePersonWearPrefix')
})

const clothGoodExamples = computed(() => {
  if (uploadTarget.value === 'upper') {
    return buildIndexedExampleItems(['upcan.jpg', 'upcan2.jpg', 'upcan3.jpg', 'allcan.jpg', 'allcan2.jpg', 'allcan3.jpg'], 'tryonExampleGoodPrefix')
  }
  if (uploadTarget.value === 'lower') {
    return buildIndexedExampleItems(['downcan.jpg', 'downcan2.jpg', 'downcan3.jpg'], 'tryonExampleGoodPrefix')
  }
  return []
})

const clothBadExamples = computed(() => {
  if (uploadTarget.value === 'upper') {
    return buildIndexedExampleItems(['clothcant.jpg', 'clothcant2.jpg', 'clothcant3.jpg'], 'tryonExampleBadPrefix')
  }
  if (uploadTarget.value === 'lower') {
    return buildIndexedExampleItems(['clothcant.jpg', 'clothcant2.jpg', 'clothcant3.jpg'], 'tryonExampleBadPrefix')
  }
  return []
})

const clothRecommendedGroups = computed(() => {
  if (uploadTarget.value === 'upper') {
    return [
      {
        key: 'person-wear',
        title: $t.value('uploadPersonWearingImage'),
        items: buildIndexedExampleItems(['model-man2.jpg', 'model-baby.jpg', 'model-women.jpg'], 'tryonExamplePersonWearPrefix'),
      },
      {
        key: 'upper',
        title: $t.value('uploadUpperImage'),
        items: buildIndexedExampleItems(['upcan.jpg', 'upcan2.jpg', 'upcan3.jpg'], 'tryonExampleUpperPrefix'),
      },
      {
        key: 'all',
        title: $t.value('suitSet'),
        items: buildIndexedExampleItems(['allcan.jpg', 'allcan2.jpg', 'allcan3.jpg'], 'tryonExampleAllPrefix'),
      },
    ]
  }

  if (uploadTarget.value === 'lower') {
    return [
      {
        key: 'lower',
        title: $t.value('uploadLowerImage'),
        items: buildIndexedExampleItems(['downcan.jpg', 'downcan2.jpg', 'downcan3.jpg'], 'tryonExampleLowerPrefix'),
      },
    ]
  }

  return []
})

const myClothDisplayName = (category) => {
  if (category === 'upper') return $t.value('clothCategoryUpper')
  if (category === 'lower') return $t.value('clothCategoryLower')
  if (category === 'onepiece') return $t.value('clothCategoryOnepiece')
  if (category === 'shoes') return $t.value('clothCategoryShoes')
  return $t.value('myCloset')
}

const filteredMyClothList = computed(() => {
  if (uploadTarget.value === 'upper') {
    return myClothList.value.filter(item => item.category === 'upper' || item.category === 'onepiece')
  }
  if (uploadTarget.value === 'lower') {
    return myClothList.value.filter(item => item.category === 'lower')
  }
  return []
})

const uploadSceneKey = computed(() => {
  return uploadTarget.value || 'person'
})

const isPersonDrawer = computed(() => uploadSceneKey.value === 'person')
const isClothDrawer = computed(() => uploadSceneKey.value === 'upper' || uploadSceneKey.value === 'lower')

const drawerTabs = computed(() => {
  if (isPersonDrawer.value) {
    return [
      { key: 'custom', labelKey: 'drawerTabCustomUpload' },
      { key: 'official', labelKey: 'drawerTabOfficialModel' },
      { key: 'myModel', labelKey: 'drawerTabMyModel' },
    ]
  }

  return [
    { key: 'custom', labelKey: 'drawerTabCustomUpload' },
    { key: 'recommended', labelKey: 'drawerTabRecommended' },
    { key: 'myCloset', labelKey: 'drawerTabMyCloset' },
  ]
})

const uploadDrawerTitle = computed(() => {
  const titleKeyMap = {
    person: 'uploadModelImage',
    upper: 'uploadUpperImage',
    lower: 'uploadLowerImage',
  }
  return $t.value(titleKeyMap[uploadSceneKey.value] || 'uploadImages')
})

const uploadDrawerTip = computed(() => {
  if (isPersonDrawer.value) {
    return $t.value('uploadDrawerTipPerson')
  }
  return $t.value('uploadDrawerTipClothes')
})

const rebuildModelList = () => {
  modelList.value = parseTryonModels(
    tryonConfig.value.tryon_models,
    activeSceneType.value,
    Number(tryonConfig.value.tryon_cost_points || 1),
    langStore.locale
  )
  if (!modelList.value.find(item => item.key === selectedModelKey.value)) {
    selectedModelKey.value = modelList.value[0]?.key || ''
  }
  const selected = modelList.value.find(item => item.key === selectedModelKey.value) || modelList.value[0]
  if (!supportsRefinerByModel(selected)) {
    refinerEnabled.value = false
  }
}

const resolveDraftBeautifyMeta = () => {
  const beautifyModels = parseTryonBeautifyModels(
    tryonConfig.value.tryon_models,
    activeSceneType.value,
    Number(tryonConfig.value.tryon_cost_points || 0),
    langStore.locale
  )

  const selectedBeautifyModel = beautifyModels[0] || null

  if (selectedBeautifyModel) {
    return {
      supportsBeautify: true,
      beautifyModelKey: String(selectedBeautifyModel.key || ''),
      beautifyModel: String(selectedBeautifyModel.beautifyModel || 'RetouchSkin'),
      beautifyExtraCost: Math.max(0, Number(selectedBeautifyModel.beautifyExtraCost || 0)),
      beautifyRetouchDegree: normalizeOptionalBeautifyDegree(selectedBeautifyModel.beautifyRetouchDegree),
      beautifyWhiteningDegree: normalizeOptionalBeautifyDegree(selectedBeautifyModel.beautifyWhiteningDegree),
      beautifyDesc: selectedBeautifyModel.beautifyDesc || '',
      beautifyDescText: selectedBeautifyModel.beautifyDescText || '',
    }
  }

  return {
    supportsBeautify: false,
    beautifyModelKey: '',
    beautifyModel: 'custom_beautify',
    beautifyExtraCost: 0,
    beautifyRetouchDegree: 0,
    beautifyWhiteningDegree: 0,
    beautifyDesc: '',
    beautifyDescText: '',
  }
}

const assignImageToTarget = (target, value, isRemote = false, sizeBytes = 0) => {
  const localPath = isRemote ? '' : value
  const remotePath = isRemote ? value : ''

  if (target === 'person') {
    personLocal.value = localPath
    personRemote.value = remotePath
  }
  if (target === 'upper') {
    upperLocal.value = localPath
    upperRemote.value = remotePath
  }
  if (target === 'lower') {
    lowerLocal.value = localPath
    lowerRemote.value = remotePath
  }

  setImageSizeForTarget(target, sizeBytes)
}

const previewRoomImage = (target) => {
  const scene = String(target || '').trim()
  let preview = ''
  if (scene === 'person') {
    preview = personPreview.value
  } else if (scene === 'upper') {
    preview = upperPreview.value
  } else if (scene === 'lower') {
    preview = lowerPreview.value
  }
  if (!preview) {
    uni.showToast({ title: $t.value('noImagePreview'), icon: 'none' })
    return
  }
  uni.previewImage({ urls: [preview] })
}

const clearRoomImage = (target) => {
  const scene = String(target || '').trim()
  if (!scene) return
  assignImageToTarget(scene, '', false, 0)
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

const normalizeMyClothItem = (item) => {
  const id = item?.ID || item?.id || ''
  const rawUrl = item?.image || item?.url || ''
  const category = String(item?.category || item?.Category || '').trim().toLowerCase()
  return {
    id: String(id),
    name: item?.name || '',
    category,
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

const loadMyClothList = async () => {
  try {
    const res = await getMyTryonClothList()
    if (res.code !== 0) return
    const list = Array.isArray(res?.data?.list) ? res.data.list : []
    myClothList.value = list
      .map(normalizeMyClothItem)
      .filter(item => item.id && item.url && ['upper', 'lower', 'onepiece'].includes(item.category))
  } catch (e) {
    myClothList.value = []
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

const loadAnnouncement = async () => {
  try {
    const res = await getAnnouncementConfig({ includeI18n: true })
    if (res.code !== 0 || !res.data || typeof res.data !== 'object') return
    announcementConfig.value = {
      ...announcementConfig.value,
      ...res.data,
    }
  } catch (e) {
    // ignore announcement fetch failure
  }
}

const uploadFolderByTarget = (target) => {
  if (target === 'person') return 'cloth-on/up-try-cloth'
  if (target === 'upper') return 'cloth-on/up-try-cloth'
  if (target === 'lower') return 'cloth-on/up-try-cloth'
  return 'cloth-on/up-try-cloth'
}

const uploadTypeByTarget = (target) => {
  if (target === 'person') return 'person'
  if (target === 'upper' || target === 'lower') return 'cloth'
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
      uni.showToast({ title: resolveApiMessage(e?.message, 'uploadFail'), icon: 'none' })
    } finally {
      uni.hideLoading()
    }
  } catch (e) {
    const errMsg = String(e?.errMsg || '')
    if (!errMsg.includes('cancel')) {
      uni.showToast({ title: resolveApiMessage(e?.message, 'uploadFail'), icon: 'none' })
    }
  }
}

const openUploadDrawer = (target) => {
  uploadTarget.value = target
  drawerTab.value = 'custom'
  showUploadDrawer.value = true
  if (target === 'person') {
    loadMyModelList()
  } else if (target === 'upper' || target === 'lower') {
    loadMyClothList()
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
  applyRemoteSizeForTarget(uploadTarget.value, previewUrl)
  showUploadDrawer.value = false
  uni.showToast({ title: $t.value('autoFillApplied'), icon: 'none' })
}

const applyMyModel = (item) => {
  if (!item?.url) return
  assignImageToTarget('person', item.url, true)
  applyRemoteSizeForTarget('person', item.url)
  showUploadDrawer.value = false
  uni.showToast({ title: $t.value('autoFillApplied'), icon: 'none' })
}

const applyMyCloth = (item) => {
  if (!item?.url || !uploadTarget.value) return
  assignImageToTarget(uploadTarget.value, item.url, true)
  applyRemoteSizeForTarget(uploadTarget.value, item.url)
  showUploadDrawer.value = false
  uni.showToast({ title: $t.value('autoFillApplied'), icon: 'none' })
}

const onRefinerSwitchChange = (event) => {
  const enabled = !!event?.detail?.value
  if (!currentModelSupportsRefiner.value) {
    refinerEnabled.value = false
    uni.showToast({ title: $t.value('tryonRefinerUnsupportedHint'), icon: 'none' })
    return
  }
  refinerEnabled.value = enabled
}

const openFeatureHelp = () => {
  showRefinerHelpPopup.value = true
}

const selectUpperOnlySplitMode = (mode) => {
  if (mode !== 'full_outfit' && mode !== 'upper_only') {
    return
  }
  upperOnlySplitMode.value = mode
}

const selectModel = (key) => {
  selectedModelKey.value = key
  const selected = modelList.value.find(v => v.key === key)
  if (!supportsRefinerByModel(selected)) {
    refinerEnabled.value = false
  }
  showModelPopup.value = false
}

const loadConfig = async () => {
  const res = await getTryonConfig({ includeI18n: true })
  if (res.code === 0 && res.data) {
    tryonConfig.value = { ...tryonConfig.value, ...res.data }
  }
  rebuildModelList()
}

const applySelectedClothes = () => {
  const selected = getSelectedClothes()
  if (!selected || typeof selected !== 'object') return
  if (selected.upperImage) {
    const upperSizeBytes = Number(selected.upperSizeBytes || 0)
    assignImageToTarget('upper', selected.upperImage, true, upperSizeBytes)
    if (upperSizeBytes <= 0) {
      applyRemoteSizeForTarget('upper', selected.upperImage)
    }
  }
  if (selected.lowerImage) {
    const lowerSizeBytes = Number(selected.lowerSizeBytes || 0)
    assignImageToTarget('lower', selected.lowerImage, true, lowerSizeBytes)
    if (lowerSizeBytes <= 0) {
      applyRemoteSizeForTarget('lower', selected.lowerImage)
    }
  }
  clearSelectedClothes()
}

const applySelectedModel = () => {
  const selectedModel = getSelectedTryonModel()
  if (!selectedModel || typeof selectedModel !== 'object') return
  if (selectedModel.roomType && selectedModel.roomType !== 'tryon') return

  if (selectedModel.remoteUrl) {
    assignImageToTarget('person', selectedModel.remoteUrl, true, 0)
    applyRemoteSizeForTarget('person', selectedModel.remoteUrl)
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

const resolveTemplateDraftPayload = () => {
  const upperRemoteUrl = upperRemote.value
  const upperLocalPath = upperLocal.value
  const lowerRemoteUrl = lowerRemote.value
  const lowerLocalPath = lowerLocal.value
  const templatePart = resolveTemplatePartBySelection()
  const parsingParts = resolveParsingPartsBySelection()

  let templateRemoteUrl = templatePart === 'lower' ? lowerRemoteUrl : upperRemoteUrl
  let templateLocalPath = templatePart === 'lower' ? lowerLocalPath : upperLocalPath

  if (!templateRemoteUrl && !templateLocalPath) {
    templateRemoteUrl = upperRemoteUrl || lowerRemoteUrl
    templateLocalPath = upperLocalPath || lowerLocalPath
  }

  return {
    templatePart,
    parsingParts,
    templateRemoteUrl,
    templateLocalPath,
    templateUpperRemoteUrl: upperRemoteUrl,
    templateUpperLocalPath: upperLocalPath,
    templateLowerRemoteUrl: lowerRemoteUrl,
    templateLowerLocalPath: lowerLocalPath,
  }
}

const saveTryonDraftAndGoGenerate = () => {
  const templatePayload = resolveTemplateDraftPayload()
  if (!templatePayload.templateRemoteUrl && !templatePayload.templateLocalPath) {
    uni.showToast({ title: $t.value('uploadClothesFirst'), icon: 'none' })
    return
  }

  const sourceUploadFolder = 'cloth-on/up-try-cloth'
  const templateUploadFolder = 'cloth-on/up-try-cloth'
  const beautifyMeta = resolveDraftBeautifyMeta()
  const beautifyRetouchDegree = normalizeOptionalBeautifyDegree(beautifyMeta.beautifyRetouchDegree)
  const beautifyWhiteningDegree = normalizeOptionalBeautifyDegree(beautifyMeta.beautifyWhiteningDegree)
  const parsingEnabled = parsingEnabledBySelection.value
  const parsingParts = parsingEnabled ? templatePayload.parsingParts : []

  saveTryonDraft({
    roomType: 'tryon',
    sceneType: 'clothes',
    templatePart: templatePayload.templatePart,
    parsingParts,
    operationType: 'tryon',
    sourceLocalPath: personLocal.value,
    sourceRemoteUrl: personRemote.value,
    sourceUploadFolder,
    templateLocalPath: templatePayload.templateLocalPath,
    templateRemoteUrl: templatePayload.templateRemoteUrl,
    templateUpperLocalPath: templatePayload.templateUpperLocalPath,
    templateUpperRemoteUrl: templatePayload.templateUpperRemoteUrl,
    templateLowerLocalPath: templatePayload.templateLowerLocalPath,
    templateLowerRemoteUrl: templatePayload.templateLowerRemoteUrl,
    templateUploadFolder,
    modelKey: currentModel.value.key,
    modelName: currentModel.value.name,
    parsingModelKey: String(currentModel.value.parsingModelKey || ''),
    parsingExtraCost: parsingEnabled ? currentParsingExtraCost.value : 0,
    baseModelCost: currentBaseCost.value,
    refinerExtraCost: currentRefinerExtraCost.value,
    modelCost: currentCost.value,
    enableRefiner: currentModelSupportsRefiner.value && refinerEnabled.value,
    enableParsing: parsingEnabled,
    refinerModelKey: String(currentModel.value.refinerModelKey || ''),
    refinerModel: String(currentModel.value.refinerModel || 'aitryon-refiner'),
    supportsBeautify: !!beautifyMeta.supportsBeautify,
    beautifyModelKey: String(beautifyMeta.beautifyModelKey || ''),
    beautifyModel: String(beautifyMeta.beautifyModel || 'RetouchSkin'),
    beautifyExtraCost: Math.max(0, Number(beautifyMeta.beautifyExtraCost || 0)),
    beautifyRetouchDegree,
    beautifyWhiteningDegree,
    beautifyDesc: beautifyMeta.beautifyDesc || '',
    beautifyDescText: beautifyMeta.beautifyDescText || '',
    requestID: createTryonRequestId(),
  })

  uni.navigateTo({ url: '/pages/tryon/generate' })
}

const goGenerate = () => {
  if (!personLocal.value && !personRemote.value) {
    uni.showToast({ title: $t.value('uploadModelFirst'), icon: 'none' })
    return
  }

  if (!hasUpperTemplate() && !hasLowerTemplate()) {
    uni.showToast({ title: $t.value('uploadClothesFirst'), icon: 'none' })
    return
  }

  if (showUpperOnlySplitSelector.value && !upperOnlySplitMode.value) {
    uni.showModal({
      title: '',
      content: upperOnlySplitGuideText.value,
      showCancel: false,
      confirmText: $t.value('doneText'),
    })
    return
  }

  const token = uni.getStorageSync('x-token') || ''
  if (token) {
    const balance = resolveTryonPointBalance()
    if (balance < Math.max(0, Number(currentCost.value || 0))) {
      uni.showToast({ title: $t.value('tryonCoinsNotEnough'), icon: 'none' })
      return
    }
  }

  uni.showModal({
    title: $t.value('generateConfirmTitle'),
    content: $t.value('generateConfirmContent'),
    cancelText: $t.value('cancel'),
    confirmText: $t.value('confirm'),
    success: (confirmRes) => {
      if (!confirmRes.confirm) return
      saveTryonDraftAndGoGenerate()
    },
  })
}

const tryAutoShowLangPicker = () => {
  if (showLangPicker.value) return
  if (!langStore.shouldAutoShowLanguagePicker()) return
  showLangPicker.value = true
}

watch(
  [() => tryonConfig.value.tryon_models, () => tryonConfig.value.tryon_cost_points, () => langStore.locale],
  () => {
    rebuildModelList()
  }
)

watch([isPersonDrawer, drawerTab], ([isPerson, tab]) => {
  if (isPerson && tab === 'myModel') {
    loadMyModelList()
  }
})

watch([isClothDrawer, drawerTab], ([isCloth, tab]) => {
  if (isCloth && tab === 'myCloset') {
    loadMyClothList()
  }
})

watch(showUpperOnlySplitSelector, (visible, prevVisible) => {
  if (visible !== prevVisible) {
    upperOnlySplitMode.value = ''
  }
})

onShow(() => {
  tryAutoShowLangPicker()
  appConfigStore.loadConfig()
  loadExampleDomain()
  loadAnnouncement()
  loadConfig()
  loadMyModelList()
  loadMyClothList()
  applySelectedModel()
  applySelectedClothes()
})
</script>

<style lang="scss">
page {
  background: #f4f7fb;
}

.room-page {
  min-height: 100vh;
  padding: calc(var(--status-bar-height, 0px) + 24rpx) 24rpx 24rpx;
  color: #0f172a;
  background: radial-gradient(120% 80% at 100% -10%, #dbeafe 0%, transparent 60%), #f4f7fb;
}

.room-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20rpx;
}

.room-title {
  font-size: 38rpx;
  font-weight: 700;
}

.room-header-action {
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
  background: rgba(255, 255, 255, 0.92);
}

.room-tabs {
  margin-top: 16rpx;
  display: flex;
  gap: 12rpx;
}

.room-tab {
  flex: 1;
  text-align: center;
  padding: 16rpx 0;
  border-radius: 999rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.12);
  color: rgba(15, 23, 42, 0.55);
  background: rgba(255, 255, 255, 0.75);
}

.room-tab.active {
  color: #0f172a;
  border-color: rgba(37, 99, 235, 0.45);
  background: linear-gradient(90deg, rgba(191, 219, 254, 0.95), rgba(219, 234, 254, 0.95));
}

.room-body {
  display: flex;
  gap: 16rpx;
}

.room-left {
  width: 70%;
  height: 660rpx;
  position: relative;
  border-radius: 20rpx;
  border: 1rpx dashed rgba(15, 23, 42, 0.2);
  overflow: hidden;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 12rpx 36rpx rgba(15, 23, 42, 0.06);
}

.room-right {
  width: 30%;
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.cloth-slot {
  flex: 1;
  position: relative;
  border-radius: 16rpx;
  border: 1rpx dashed rgba(15, 23, 42, 0.2);
  background: rgba(255, 255, 255, 0.92);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 10rpx 26rpx rgba(15, 23, 42, 0.05);
}

.suit-slot {
  border-style: solid;
  border-color: rgba(37, 99, 235, 0.45);
}

.cloth-tip {
  margin-top: 8rpx;
  font-size: 20rpx;
  color: rgba(15, 23, 42, 0.55);
}

.cloth-text {
  color: rgba(15, 23, 42, 0.82);
  font-size: 24rpx;
}

.room-preview,
.takeoff-preview,
.cloth-preview {
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

.slot-actions {
  position: absolute;
  top: 12rpx;
  right: 12rpx;
  z-index: 4;
  display: flex;
  gap: 8rpx;
}

.slot-action-btn {
  width: 48rpx;
  height: 48rpx;
  border-radius: 999rpx;
  background: rgba(255, 255, 255, 0.92);
  border: 1rpx solid rgba(37, 99, 235, 0.28);
  display: flex;
  align-items: center;
  justify-content: center;
}

.slot-action-btn.danger {
  border-color: rgba(220, 38, 38, 0.28);
}

.takeoff-body {
  height: 660rpx;
  border-radius: 20rpx;
  border: 1rpx dashed rgba(15, 23, 42, 0.2);
  overflow: hidden;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 12rpx 36rpx rgba(15, 23, 42, 0.06);
}

.room-upload-empty {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 12rpx;
}

.room-upload-text {
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.6);
}

.refiner-card {
  margin-top: 18rpx;
  border: 1rpx solid rgba(37, 99, 235, 0.2);
  border-radius: 16rpx;
  padding: 18rpx 20rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(239, 246, 255, 0.92);
}

.image-size-rule-card {
  margin-top: 18rpx;
  padding: 14rpx 18rpx;
  border-radius: 14rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.12);
  background: rgba(255, 255, 255, 0.94);
}

.image-size-rule-text {
  display: block;
  font-size: 22rpx;
  line-height: 1.45;
  color: rgba(15, 23, 42, 0.72);
}

.refiner-card.disabled {
  border-color: rgba(148, 163, 184, 0.3);
  background: rgba(248, 250, 252, 0.96);
}

.refiner-left {
  display: flex;
  flex-direction: column;
  gap: 6rpx;
  min-width: 0;
}

.refiner-label {
  color: #0f172a;
  font-size: 24rpx;
  font-weight: 600;
}

.refiner-hint {
  font-size: 20rpx;
  color: rgba(15, 23, 42, 0.58);
}

.refiner-extra-cost {
  font-size: 20rpx;
  color: rgba(37, 99, 235, 0.92);
}

.refiner-right {
  display: flex;
  align-items: center;
  gap: 14rpx;
  flex-shrink: 0;
}

.refiner-help {
  width: 42rpx;
  height: 42rpx;
  border-radius: 50%;
  border: 1rpx solid rgba(37, 99, 235, 0.45);
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24rpx;
  font-weight: 700;
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

.model-cost-total {
  display: block;
  margin-top: 6rpx;
  font-size: 20rpx;
  color: #1d4ed8;
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

.upper-only-split-card {
  margin-top: 18rpx;
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.upper-only-split-options {
  flex: 1;
  display: flex;
  gap: 10rpx;
}

.upper-only-split-option {
  flex: 1;
  height: 66rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.12);
  background: rgba(255, 255, 255, 0.95);
  color: rgba(15, 23, 42, 0.7);
  font-size: 23rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.upper-only-split-option.active {
  color: #ffffff;
  border-color: transparent;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
}

.upper-only-split-help {
  width: 46rpx;
  height: 46rpx;
  border-radius: 50%;
  border: 1rpx solid rgba(37, 99, 235, 0.45);
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24rpx;
  font-weight: 700;
  flex-shrink: 0;
}

.generate-btn {
  margin-top: 20rpx;
  height: 88rpx;
  border-radius: 999rpx;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 14rpx 30rpx rgba(37, 99, 235, 0.32);
}

.generate-text {
  font-size: 28rpx;
  font-weight: 700;
  color: #fff;
}

.footer-tips {
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

.popup-panel-help {
  max-height: 72vh;
  min-height: 240rpx;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.popup-panel-tutorial {
  height: 72vh;
}

.refiner-help-content {
  margin-top: 8rpx;
  color: rgba(15, 23, 42, 0.72);
  font-size: 24rpx;
  line-height: 1.6;
  white-space: pre-wrap;
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

.refiner-help-actions {
  margin-top: 26rpx;
  display: flex;
  justify-content: flex-end;
}

.tutorial-list {
  margin-top: 8rpx;
  display: flex;
  flex-direction: column;
  gap: 10rpx;
}

.tutorial-popup-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
  margin-bottom: 8rpx;
}

.tutorial-popup-title {
  margin-bottom: 0;
}

.tutorial-popup-close {
  width: 56rpx;
  height: 56rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(15, 23, 42, 0.06);
  color: rgba(15, 23, 42, 0.68);
  font-size: 44rpx;
  line-height: 1;
  flex-shrink: 0;
}

.tutorial-scroll {
  flex: 1;
  height: 0;
  min-height: 0;
  overscroll-behavior: contain;
}

.tutorial-item {
  display: flex;
  align-items: flex-start;
  gap: 8rpx;
}

.tutorial-index {
  width: 30rpx;
  flex-shrink: 0;
  font-size: 23rpx;
  line-height: 1.55;
  color: #1d4ed8;
  font-weight: 600;
}

.tutorial-text {
  flex: 1;
  font-size: 23rpx;
  line-height: 1.55;
  color: rgba(15, 23, 42, 0.76);
}

.popup-title {
  font-size: 30rpx;
  font-weight: 700;
  margin-bottom: 16rpx;
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

.popup-item-extra-cost {
  display: block;
  margin-top: 6rpx;
  font-size: 20rpx;
  color: rgba(245, 158, 11, 0.95);
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

.upload-highlight-tip {
  display: block;
  margin-bottom: 12rpx;
  padding: 12rpx 14rpx;
  border-radius: 10rpx;
  border: 1rpx solid rgba(245, 158, 11, 0.38);
  background: rgba(254, 243, 199, 0.62);
  color: #92400e;
  font-size: 21rpx;
  line-height: 1.55;
}

.example-grid {
  display: grid;
  gap: 12rpx;
}

.example-grid.two {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.example-grid.three {
  grid-template-columns: repeat(3, minmax(0, 1fr));
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

.group-title {
  display: block;
  font-size: 24rpx;
  font-weight: 700;
  color: #0f172a;
  margin: 8rpx 0 12rpx;
}

.recommend-group {
  margin-bottom: 8rpx;
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

.upload-drawer-btn.disabled {
  opacity: 0.65;
  pointer-events: none;
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
