<template>
  <view class="nf-reviews">
    <view class="nf-reviews-bg"></view>

    <!-- navbar -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff"></uni-icons>
        </view>
        <text class="nf-navbar-title">{{ $t('productReview') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <!-- tabs -->
    <view class="nf-tabs">
      <view :class="['nf-tab', { active: tab === 'all' }]" @tap="choose('all')">
        {{ $t('allReviews') }}（{{ totalCount }}）
      </view>
      <view :class="['nf-tab', { active: tab === 'pics' }]" @tap="choose('pics')">
        {{ $t('imageReviews') }}（{{ picCount }}）
      </view>
    </view>

    <view class="nf-body">
      <view class="nf-empty" v-if="commentInfo.length === 0">
        <text class="nf-empty-icon">💬</text>
        <text class="nf-empty-text">{{ $t('noData') || '暂无评价' }}</text>
      </view>

      <view class="nf-comment-card" v-for="(item, key) in commentInfo" :key="key">
        <!-- user info -->
        <view class="nf-comment-header">
          <image class="nf-avatar" :src="getUrl(item.user.avatar)" mode="aspectFill" />
          <view class="nf-comment-meta">
            <text class="nf-nickname">{{ item.user.nickname }}</text>
            <text class="nf-time">{{ formatTimeToStr(item.CreatedAt, 'yyyy-MM-dd') }}</text>
          </view>
        </view>

        <!-- rating -->
        <view class="nf-comment-stars">
          <text v-for="s in 5" :key="s" :class="['nf-mini-star', { active: s <= item.rating }]"></text>
        </view>

        <!-- content -->
        <text class="nf-comment-text" v-if="item.content">{{ item.content }}</text>

        <!-- pics -->
        <view class="nf-comment-pics" v-if="item.pics && item.pics.length > 0">
          <view
            v-for="(pic, index) in item.pics"
            :key="index"
            class="nf-comment-pic"
            @tap="previewImage(pic, index, item.pics)"
          >
            <image :src="getUrl(pic)" class="nf-comment-pic-img" mode="aspectFill" />
          </view>
        </view>

        <!-- shop reply -->
        <view class="nf-shop-reply" v-if="item.shopReply">
          <text class="nf-shop-reply-label">{{ $t('shopReplyLabel') }}：</text>
          <text class="nf-shop-reply-text">{{ item.shopReply }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { findComment } from '@/api/comment.js'
import { formatTimeToStr } from '@/utils/date.js'
import { onLoad } from '@dcloudio/uni-app'
import { getUrl } from '@/utils/url.js'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const commentInfo = ref([])
const allComments = ref([])
const picCount = ref(0)
const totalCount = ref(0)
const tab = ref('all')
let ID = ''

const goBack = () => {
  uni.navigateBack({ delta: 1 })
}

const findFunc = async (params) => {
  const res = await findComment({ ID: params })
  if (res.code === 0) {
    allComments.value = res.data || []
    commentInfo.value = allComments.value
    totalCount.value = allComments.value.length
    picCount.value = allComments.value.filter(i => i.pics && i.pics.length > 0).length
  }
}

onLoad((options) => {
  ID = options.goodsID
  setTimeout(() => findFunc(ID), 300)
})

const choose = (type) => {
  tab.value = type
  if (type === 'all') {
    commentInfo.value = allComments.value
  } else {
    commentInfo.value = allComments.value.filter(item => item.pics && item.pics.length > 0)
  }
}

const previewImage = (currentPic, index, allPics) => {
  uni.previewImage({
    current: getUrl(currentPic),
    urls: allPics.map(pic => getUrl(pic))
  })
}
</script>

<style lang="scss" scoped>
page { background: #000; }
.nf-reviews {
  min-height: 100vh;
  background: #000;
  color: #fff;
}
.nf-reviews-bg {
  position: fixed; top: 0; left: 0; right: 0;
  height: 400rpx;
  background: radial-gradient(ellipse at 50% 0%, rgba(229, 9, 20, 0.12) 0%, transparent 70%);
  pointer-events: none; z-index: 0;
}
.nf-navbar {
  position: sticky; top: 0; z-index: 100;
  background: rgba(0, 0, 0, 0.85);
  backdrop-filter: blur(20rpx);
}
.nf-navbar-status { height: var(--status-bar-height); }
.nf-navbar-content {
  display: flex; align-items: center; justify-content: space-between;
  height: 88rpx; padding: 0 24rpx;
}
.nf-navbar-back {
  width: 64rpx; height: 64rpx;
  display: flex; align-items: center; justify-content: center;
}
.nf-navbar-title { font-size: 34rpx; font-weight: 600; color: #fff; }

.nf-tabs {
  display: flex; padding: 0 24rpx;
  position: relative; z-index: 1;
  background: rgba(0, 0, 0, 0.6);
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
}
.nf-tab {
  padding: 20rpx 32rpx;
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.5);
  border-bottom: 3rpx solid transparent;
  transition: all 0.2s;
  &.active {
    color: #e50914;
    border-bottom-color: #e50914;
    font-weight: 600;
  }
}

.nf-body { position: relative; z-index: 1; padding: 24rpx; }

.nf-empty {
  display: flex; flex-direction: column; align-items: center;
  padding: 120rpx 0;
}
.nf-empty-icon { font-size: 80rpx; margin-bottom: 20rpx; }
.nf-empty-text { font-size: 28rpx; color: rgba(255, 255, 255, 0.4); }

.nf-comment-card {
  background: rgba(255, 255, 255, 0.06);
  border-radius: 16rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.08);
  padding: 28rpx;
  margin-bottom: 20rpx;
}
.nf-comment-header {
  display: flex; align-items: center; margin-bottom: 16rpx;
}
.nf-avatar {
  width: 64rpx; height: 64rpx;
  border-radius: 50%; margin-right: 16rpx;
  border: 2rpx solid rgba(255, 255, 255, 0.1);
}
.nf-comment-meta { flex: 1; display: flex; justify-content: space-between; align-items: center; }
.nf-nickname { font-size: 28rpx; color: #fff; font-weight: 500; }
.nf-time { font-size: 22rpx; color: rgba(255, 255, 255, 0.35); }

.nf-comment-stars { display: flex; margin-bottom: 12rpx; }
.nf-mini-star {
  font-size: 28rpx; color: rgba(255, 255, 255, 0.15); margin-right: 4rpx;
  &.active { color: #e50914; }
}

.nf-comment-text {
  font-size: 28rpx; color: rgba(255, 255, 255, 0.8);
  line-height: 1.6; margin-bottom: 16rpx;
}

.nf-comment-pics {
  display: flex; flex-wrap: wrap; gap: 12rpx; margin-bottom: 16rpx;
}
.nf-comment-pic {
  width: 140rpx; height: 140rpx;
  border-radius: 10rpx; overflow: hidden;
}
.nf-comment-pic-img { width: 100%; height: 100%; }

.nf-shop-reply {
  background: rgba(229, 9, 20, 0.08);
  border-radius: 10rpx;
  padding: 16rpx 20rpx;
  border-left: 4rpx solid #e50914;
  margin-top: 12rpx;
}
.nf-shop-reply-label {
  font-size: 24rpx; color: #e50914; font-weight: 500;
}
.nf-shop-reply-text {
  font-size: 26rpx; color: rgba(255, 255, 255, 0.7); line-height: 1.5;
}
</style>
