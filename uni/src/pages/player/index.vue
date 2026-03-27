<template>
  <view class="nf-player">

    <!-- ===== 全屏播放器弹层 ===== -->
    <view class="nf-theater" :class="{ 'nf-theater-show': showTheater }" @tap.self="closeTheater">
      <view class="nf-theater-inner" :class="{ 'nf-theater-fade-in': theaterReady }">
        <!-- 关闭按钮 -->
        <view class="nf-theater-close" @tap="closeTheater">
          <uni-icons type="closeempty" size="22" color="#fff" />
        </view>
        <!-- 正在播放标签 -->
        <view class="nf-theater-badge">
          <view class="nf-theater-dot"></view>
          <text class="nf-theater-badge-text">{{ $t('playerNowPlaying') }}</text>
        </view>
        <!-- 视频播放器 -->
        <view class="nf-theater-video-wrap">
          <video
            v-if="videoUrl && showTheater"
            id="nfVideo"
            :src="videoUrl"
            class="nf-theater-video"
            :poster="getUrl(data.imageUrl)"
            :title="currentEpisode.name || $lt(data.title)"
            :autoplay="true"
            :show-fullscreen-btn="true"
            :show-play-btn="true"
            :enable-progress-gesture="true"
            object-fit="contain"
            @ended="onVideoEnded"
          />
        </view>
        <!-- 当前集信息 -->
        <view class="nf-theater-info">
          <text class="nf-theater-title">{{ currentEpisode.name || $lt(data.title) }}</text>
          <text class="nf-theater-sub" v-if="currentEpisode.description">{{ currentEpisode.description }}</text>
        </view>
        <!-- 集数快选 -->
        <view class="nf-theater-eps" v-if="episodes.length > 1">
          <scroll-view scroll-x :show-scrollbar="false" class="nf-theater-eps-scroll">
            <view class="nf-theater-eps-list">
              <view
                v-for="(ep, idx) in episodes"
                :key="ep.ID"
                class="nf-theater-ep-btn"
                :class="{ 'nf-theater-ep-active': currentIndex === idx }"
                @tap="switchEpisode(idx)"
              >
                <text class="nf-theater-ep-text">{{ idx + 1 }}</text>
              </view>
            </view>
          </scroll-view>
        </view>
      </view>
    </view>

    <!-- ===== 主页面内容 ===== -->
    <!-- 封面 Hero -->
    <view class="nf-hero">
      <view class="nf-hero-status"></view>
      <image class="nf-hero-img" :src="getUrl(data.imageUrl)" mode="aspectFill" @tap="playFirst" />
      <view class="nf-hero-gradient"></view>
      <!-- 返回按钮 -->
      <view class="nf-hero-back" @tap="goBack">
        <uni-icons type="left" size="20" color="#fff" />
      </view>

    </view>

    <!-- 内容区域 -->
    <scroll-view scroll-y :show-scrollbar="false" class="nf-content-scroll">
      <!-- 标题信息 -->
      <view class="nf-info-section">
        <text class="nf-title">{{ $lt(data.title) }}</text>
        <view class="nf-meta-row">
          <view class="nf-meta-item" v-if="data.rating">
            <text class="nf-meta-star">★</text>
            <text class="nf-meta-value">{{ data.rating }}</text>
            <text class="nf-meta-label">{{ $t('playerRating') }}</text>
          </view>
          <view class="nf-meta-divider" v-if="data.rating && data.view_num"></view>
          <view class="nf-meta-item" v-if="data.view_num">
            <text class="nf-meta-value">{{ formatNum(data.view_num) }}</text>
            <text class="nf-meta-label">{{ $t('playerViews') }}</text>
          </view>
          <view class="nf-meta-divider" v-if="episodes.length"></view>
          <view class="nf-meta-item" v-if="episodes.length">
            <text class="nf-meta-value">{{ episodes.length }}</text>
            <text class="nf-meta-label">{{ $t('playerEpisodes') }}</text>
          </view>
        </view>

        <!-- 操作栏 -->
        <view class="nf-actions">
          <view class="nf-action-btn" @tap="toggleCollect">
            <uni-icons
              :type="collectionFlag ? 'heart-filled' : 'heart'"
              size="22"
              :color="collectionFlag ? '#e50914' : 'rgba(255,255,255,0.6)'"
            />
            <text class="nf-action-text">{{ collectionFlag ? $t('playerCollected') : $t('playerCollect') }}</text>
          </view>
          <view class="nf-action-btn" @tap="shareVideo">
            <uni-icons type="redo" size="22" color="rgba(255,255,255,0.6)" />
            <text class="nf-action-text">{{ $t('playerShare') }}</text>
          </view>
        </view>
      </view>

      <!-- 简介 -->
      <view class="nf-desc-section" v-if="$lt(data.description)">
        <text class="nf-section-title">{{ $t('playerDesc') }}</text>
        <text class="nf-desc-text" :class="{ 'nf-desc-expand': descExpand }">{{ $lt(data.description) }}</text>
        <text class="nf-desc-toggle" @tap="descExpand = !descExpand">
          {{ descExpand ? '收起' : '展开' }}
        </text>
      </view>

      <!-- 选集 -->
      <view class="nf-episodes-section" v-if="episodes.length">
        <view class="nf-section-header">
          <text class="nf-section-title">{{ $t('playerEpisodes') }}</text>
          <text class="nf-episode-count">{{ episodes.length }}集</text>
        </view>
        <scroll-view scroll-x :show-scrollbar="false" class="nf-episodes-scroll">
          <view class="nf-episodes-list">
            <view
              v-for="(ep, idx) in episodes"
              :key="ep.ID"
              class="nf-episode-card"
              :class="{ 'nf-episode-active': currentIndex === idx }"
              @tap="playEpisode(idx)"
            >
              <view class="nf-ep-thumb-wrap">
                <image
                  class="nf-ep-thumb"
                  :src="ep.picture ? getUrl(ep.picture) : getUrl(data.imageUrl)"
                  mode="aspectFill"
                />
                <view class="nf-ep-play-icon" v-if="currentIndex !== idx">
                  <text class="nf-ep-play-text">▶</text>
                </view>
                <view class="nf-ep-playing" v-else>
                  <view class="nf-eq-bar"></view>
                  <view class="nf-eq-bar"></view>
                  <view class="nf-eq-bar"></view>
                </view>
              </view>
              <view class="nf-ep-info">
                <text class="nf-ep-name">{{ ep.name || epLabel(idx + 1) }}</text>
                <text class="nf-ep-desc" v-if="ep.description">{{ ep.description }}</text>
                <text class="nf-ep-price" v-if="ep.price > 0">¥{{ (ep.price / 100).toFixed(2) }}</text>
                <text class="nf-ep-free" v-else>{{ $t('playerFree') }}</text>
              </view>
            </view>
          </view>
        </scroll-view>
      </view>

      <!-- 无剧集时的占位 -->
      <view class="nf-no-episodes" v-if="!episodes.length && !isLoading">
        <text class="nf-no-episodes-text">{{ $t('playerNoEpisodes') }}</text>
      </view>

      <!-- 底部安全区 -->
      <view style="height: 60rpx;"></view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref, computed, nextTick } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { findGood } from '@/api/product.js'
import { findCollect, createCollect } from '@/api/collect.js'
import { getUrl } from '@/utils/url.js'
import { useUserStore } from '@/pinia/modules/user'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)

const userStore = useUserStore()
const token = userStore.token || ''

const data = ref({})
const episodes = ref([])
const currentIndex = ref(-1)
const currentEpisode = ref({})
const videoUrl = ref('')
const collectionFlag = ref(false)
const isLoading = ref(true)
const descExpand = ref(false)
const goodID = ref(0)
const showTheater = ref(false)
const theaterReady = ref(false)

const epLabel = (n) => {
  return $t.value('playerEp').replace('{{n}}', n)
}

const formatNum = (num) => {
  if (!num) return '0'
  if (num >= 10000) return (num / 10000).toFixed(1) + 'w'
  return String(num)
}

onLoad((options) => {
  if (options.id) {
    goodID.value = options.id
    init()
  }
})

const init = async () => {
  isLoading.value = true
  try {
    const res = await findGood(goodID.value)
    if (res.code === 0 && res.data.regood) {
      data.value = res.data.regood
      episodes.value = res.data.regood.skus || []
    }
    if (token) {
      const status = await findCollect({ goodID: goodID.value })
      if (status.code === 0) collectionFlag.value = status.data
    }
  } catch (e) {
    console.error('加载失败', e)
  } finally {
    isLoading.value = false
  }
}

const openTheater = () => {
  showTheater.value = true
  nextTick(() => {
    setTimeout(() => { theaterReady.value = true }, 50)
  })
}

const closeTheater = () => {
  theaterReady.value = false
  setTimeout(() => {
    showTheater.value = false
    videoUrl.value = ''
  }, 300)
}

const resolveVideoUrl = (ep) => {
  let url = ''
  if (ep.specs) {
    try {
      const specs = typeof ep.specs === 'string' ? JSON.parse(ep.specs) : ep.specs
      url = specs.videoUrl || specs.video_url || specs.url || ''
    } catch (e) {}
  }
  return url
}

const playFirst = () => {
  if (episodes.value.length > 0) {
    playEpisode(0)
  }
}

const playEpisode = (idx) => {
  const ep = episodes.value[idx]
  if (!ep) return
  currentIndex.value = idx
  currentEpisode.value = ep
  videoUrl.value = resolveVideoUrl(ep)
  openTheater()
}

const switchEpisode = (idx) => {
  const ep = episodes.value[idx]
  if (!ep) return
  currentIndex.value = idx
  currentEpisode.value = ep
  videoUrl.value = resolveVideoUrl(ep)
}

const onVideoEnded = () => {
  // 自动播放下一集
  if (currentIndex.value < episodes.value.length - 1) {
    switchEpisode(currentIndex.value + 1)
  }
}

const toggleCollect = async () => {
  if (!token) {
    uni.showToast({ title: $t.value('loginFirst') || '请登录', icon: 'none' })
    uni.redirectTo({ url: '/pages/user/login' })
    return
  }
  const res = await createCollect({ goodID: Number(goodID.value) })
  if (res.code === 0) {
    collectionFlag.value = !collectionFlag.value
    uni.showToast({
      title: collectionFlag.value ? $t.value('playerCollected') : $t.value('playerCollect'),
      icon: 'none'
    })
  }
}

const shareVideo = () => {
  uni.showToast({ title: $t.value('playerShare'), icon: 'none' })
}

const goBack = () => {
  uni.navigateBack()
}
</script>

<style lang="scss">
page {
  background-color: #000;
}

.nf-player {
  width: 100%;
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #000;
  position: relative;
}

/* ========== 全屏影院弹层 ========== */
.nf-theater {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  z-index: 9999;
  background: rgba(0, 0, 0, 0.96);
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.nf-theater-show {
  pointer-events: auto;
  opacity: 1;
}

.nf-theater-inner {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transform: scale(0.92);
  transition: all 0.35s cubic-bezier(0.16, 1, 0.3, 1);
}

.nf-theater-fade-in {
  opacity: 1;
  transform: scale(1);
}

.nf-theater-close {
  position: absolute;
  top: 0;
  right: 0;
  padding-top: var(--status-bar-height, 0px);
  width: 88rpx;
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
  margin-top: 16rpx;
  margin-right: 16rpx;

  &:active {
    opacity: 0.6;
    transform: scale(0.9);
  }
}

/* #ifdef MP-WEIXIN */
.nf-theater-close {
  padding-top: var(--status-bar-height, 44px);
}
/* #endif */

.nf-theater-badge {
  display: flex;
  align-items: center;
  gap: 10rpx;
  margin-bottom: 24rpx;
}

.nf-theater-dot {
  width: 14rpx;
  height: 14rpx;
  border-radius: 50%;
  background: #e50914;
  box-shadow: 0 0 10rpx rgba(229, 9, 20, 0.8);
  animation: pulseDot 1.5s ease-in-out infinite;
}

@keyframes pulseDot {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.5; transform: scale(0.8); }
}

.nf-theater-badge-text {
  font-size: 22rpx;
  color: rgba(255, 255, 255, 0.5);
  letter-spacing: 3rpx;
  text-transform: uppercase;
}

.nf-theater-video-wrap {
  width: 100%;
  height: 422rpx;
  background: #000;
  border-radius: 0;
  overflow: hidden;
  box-shadow: 0 20rpx 80rpx rgba(0, 0, 0, 0.6);
}

.nf-theater-video {
  width: 100%;
  height: 100%;
}

.nf-theater-info {
  padding: 28rpx 40rpx 16rpx;
  width: 100%;
  box-sizing: border-box;
}

.nf-theater-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #fff;
  letter-spacing: 1rpx;
  margin-bottom: 8rpx;
}

.nf-theater-sub {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.4);
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

/* 弹层内集数快选 */
.nf-theater-eps {
  width: 100%;
  padding: 16rpx 0 0;
}

.nf-theater-eps-scroll {
  width: 100%;
  white-space: nowrap;
}

.nf-theater-eps-list {
  display: flex;
  padding: 0 40rpx 24rpx;
  gap: 16rpx;
}

.nf-theater-ep-btn {
  width: 72rpx;
  height: 72rpx;
  border-radius: 16rpx;
  background: rgba(255, 255, 255, 0.06);
  border: 1rpx solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.25s;

  &:active {
    transform: scale(0.92);
  }
}

.nf-theater-ep-active {
  background: linear-gradient(135deg, #e50914, #b20710);
  border-color: #e50914;
  box-shadow: 0 0 20rpx rgba(229, 9, 20, 0.4);
}

.nf-theater-ep-text {
  font-size: 26rpx;
  font-weight: 700;
  color: #fff;
}

/* ========== 封面 Hero ========== */
.nf-hero {
  width: 100%;
  height: 560rpx;
  position: relative;
  flex-shrink: 0;
  overflow: hidden;
}

.nf-hero-status {
  position: absolute;
  top: 0; left: 0; right: 0;
  height: var(--status-bar-height, 0px);
  z-index: 5;
}

/* #ifdef MP-WEIXIN */
.nf-hero-status {
  height: var(--status-bar-height, 44px);
}
/* #endif */

.nf-hero-img {
  width: 100%;
  height: 100%;
  position: absolute;
  top: 0; left: 0;
}

.nf-hero-gradient {
  position: absolute;
  bottom: 0; left: 0; right: 0;
  height: 320rpx;
  background: linear-gradient(to top, #000 0%, rgba(0,0,0,0.7) 50%, transparent 100%);
  pointer-events: none;
}

.nf-hero-back {
  position: absolute;
  top: var(--status-bar-height, 0px);
  left: 20rpx;
  margin-top: 16rpx;
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;

  &:active {
    background: rgba(0, 0, 0, 0.7);
    transform: scale(0.92);
  }
}

/* #ifdef MP-WEIXIN */
.nf-hero-back {
  top: var(--status-bar-height, 44px);
}
/* #endif */



/* ===== 内容滚动区 ===== */
.nf-content-scroll {
  flex: 1;
  width: 100%;
  scrollbar-width: none;
  -ms-overflow-style: none;

  &::-webkit-scrollbar {
    display: none;
    width: 0;
    height: 0;
  }
}

/* ===== 标题信息区 ===== */
.nf-info-section {
  padding: 28rpx 32rpx 0;
}

.nf-title {
  font-size: 40rpx;
  font-weight: 800;
  color: #fff;
  letter-spacing: 1rpx;
  line-height: 1.3;
  margin-bottom: 16rpx;
}

.nf-meta-row {
  display: flex;
  align-items: center;
  margin-bottom: 24rpx;
}

.nf-meta-item {
  display: flex;
  align-items: baseline;
  gap: 6rpx;
}

.nf-meta-star {
  color: #ffd700;
  font-size: 26rpx;
}

.nf-meta-value {
  color: #fff;
  font-size: 28rpx;
  font-weight: 700;
}

.nf-meta-label {
  color: rgba(255, 255, 255, 0.4);
  font-size: 22rpx;
  margin-left: 4rpx;
}

.nf-meta-divider {
  width: 1rpx;
  height: 24rpx;
  background: rgba(255, 255, 255, 0.15);
  margin: 0 24rpx;
}

/* ===== 操作栏 ===== */
.nf-actions {
  display: flex;
  gap: 48rpx;
  padding: 16rpx 0 24rpx;
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
}

.nf-action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;

  &:active {
    opacity: 0.7;
    transform: scale(0.95);
  }
}

.nf-action-text {
  font-size: 22rpx;
  color: rgba(255, 255, 255, 0.5);
}

/* ===== 简介 ===== */
.nf-desc-section {
  padding: 24rpx 32rpx;
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
}

.nf-section-title {
  font-size: 30rpx;
  font-weight: 700;
  color: #fff;
  letter-spacing: 1rpx;
  margin-bottom: 12rpx;
}

.nf-desc-text {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.55);
  line-height: 1.65;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3;
}

.nf-desc-expand {
  -webkit-line-clamp: unset;
}

.nf-desc-toggle {
  font-size: 24rpx;
  color: #e50914;
  margin-top: 8rpx;

  &:active {
    opacity: 0.7;
  }
}

/* ===== 选集 ===== */
.nf-episodes-section {
  padding: 24rpx 0 0;
}

.nf-section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32rpx 16rpx;
}

.nf-episode-count {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.35);
}

.nf-episodes-scroll {
  width: 100%;
  white-space: nowrap;
}

.nf-episodes-list {
  display: flex;
  padding: 0 32rpx 24rpx;
  gap: 20rpx;
}

.nf-episode-card {
  width: 280rpx;
  flex-shrink: 0;
  border-radius: 16rpx;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.06);
  transition: all 0.3s;

  &:active {
    transform: scale(0.97);
  }
}

.nf-episode-active {
  border-color: rgba(229, 9, 20, 0.5);
  background: rgba(229, 9, 20, 0.08);
  box-shadow: 0 0 20rpx rgba(229, 9, 20, 0.15);
}

.nf-ep-thumb-wrap {
  width: 100%;
  height: 160rpx;
  position: relative;
  overflow: hidden;
}

.nf-ep-thumb {
  width: 100%;
  height: 100%;
}

.nf-ep-play-icon {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 50rpx;
  height: 50rpx;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
}

.nf-ep-play-text {
  font-size: 22rpx;
  color: #fff;
  margin-left: 3rpx;
}

.nf-ep-playing {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  align-items: flex-end;
  gap: 4rpx;
  height: 36rpx;
}

.nf-eq-bar {
  width: 6rpx;
  background: #e50914;
  border-radius: 3rpx;
  animation: eqBounce 0.8s ease-in-out infinite alternate;

  &:nth-child(1) { height: 60%; animation-delay: 0s; }
  &:nth-child(2) { height: 100%; animation-delay: 0.2s; }
  &:nth-child(3) { height: 40%; animation-delay: 0.4s; }
}

@keyframes eqBounce {
  0% { height: 30%; }
  100% { height: 100%; }
}

.nf-ep-info {
  padding: 14rpx 16rpx;
  white-space: normal;
}

.nf-ep-name {
  font-size: 24rpx;
  font-weight: 600;
  color: #fff;
  margin-bottom: 6rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.nf-ep-desc {
  font-size: 20rpx;
  color: rgba(255, 255, 255, 0.4);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 6rpx;
}

.nf-ep-price {
  font-size: 24rpx;
  color: #e50914;
  font-weight: 700;
}

.nf-ep-free {
  font-size: 22rpx;
  color: #22c55e;
  font-weight: 600;
}

/* ===== 无剧集 ===== */
.nf-no-episodes {
  padding: 80rpx 0;
  text-align: center;
}

.nf-no-episodes-text {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.3);
}
</style>
