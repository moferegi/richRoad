<template>
  <view class="nf-page">
    <view class="nf-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff" />
        </view>
        <text class="nf-navbar-title">{{ $t('inviteFriends') }}</text>
      </view>
    </view>

    <view class="nf-container">
      <!-- 邀请码卡片 -->
      <view class="nf-card invite-card">
        <text class="invite-title">{{ $t('myInviteCode') }}</text>
        <text class="invite-code">{{ inviteInfo.inviteCode || '--' }}</text>
        <view class="invite-btn-row">
          <button class="nf-btn nf-btn-primary invite-copy-btn" @tap="copyCode">{{ $t('copyInviteCode') }}</button>
          <button class="nf-btn nf-btn-ghost invite-share-btn" @tap="shareLink">{{ $t('shareLink') }}</button>
        </view>
      </view>

      <!-- 统计 -->
      <view class="nf-card stats-card">
        <view class="stats-row">
          <view class="stats-item">
            <text class="stats-value">{{ inviteInfo.subordinateCount || 0 }}</text>
            <text class="stats-label">{{ $t('invited') }}</text>
          </view>
          <view class="stats-divider"></view>
          <view class="stats-item">
            <text class="stats-value">{{ inviteInfo.point || 0 }}</text>
            <text class="stats-label">{{ $t('myPoints') }}</text>
          </view>
        </view>
      </view>

      <!-- 下级列表 -->
      <view class="nf-card">
        <text class="section-title">{{ $t('myFriends') }}</text>
        <view v-if="subordinateList.length === 0" class="empty-tip">
          <text>{{ $t('noInviteRecord') }}</text>
        </view>
        <view v-for="(item, index) in subordinateList" :key="index" class="sub-item">
          <view class="sub-left">
            <image class="sub-avatar" :src="item.avatar || defaultAvatar" mode="aspectFill"></image>
            <view class="sub-info">
              <text class="sub-name">{{ item.nickname || item.username }}</text>
              <text class="sub-time">{{ formatTime(item.CreatedAt) }}</text>
            </view>
          </view>
        </view>
        <view v-if="hasMore" class="load-more" @tap="loadMore">
          <text>{{ $t('loadMore') }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getMyInviteInfo, getMySubordinates } from '@/api/base.js'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const defaultAvatar = 'https://qmplusimg.henrongyi.top/gva_header.jpg'

const inviteInfo = ref({
  inviteCode: '',
  subordinateCount: 0,
  point: 0
})

const subordinateList = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const hasMore = computed(() => subordinateList.value.length < total.value)

const goBack = () => {
  uni.navigateBack({ fail: () => uni.switchTab({ url: '/pages/tabBar/my/index' }) })
}

const loadInviteInfo = async () => {
  const res = await getMyInviteInfo()
  if (res.code === 0) {
    inviteInfo.value = res.data
  }
}

const loadSubordinates = async () => {
  const res = await getMySubordinates({ page: page.value, pageSize: pageSize.value })
  if (res.code === 0) {
    if (page.value === 1) {
      subordinateList.value = res.data.list || []
    } else {
      subordinateList.value.push(...(res.data.list || []))
    }
    total.value = res.data.total
  }
}

const loadMore = () => {
  page.value++
  loadSubordinates()
}

const copyCode = () => {
  if (!inviteInfo.value.inviteCode) return
  uni.setClipboardData({
    data: inviteInfo.value.inviteCode,
    success: () => {
      uni.showToast({ title: $t.value('inviteCodeCopied'), icon: 'success' })
    }
  })
}

const shareLink = () => {
  const code = inviteInfo.value.inviteCode
  if (!code) return
  // #ifdef H5
  const origin = window.location.origin
  // #endif
  // #ifndef H5
  const origin = 'https://your-domain.com'
  // #endif
  const link = `${origin}/#/pages/user/register?inviteCode=${code}`
  uni.setClipboardData({
    data: link,
    success: () => {
      uni.showToast({ title: $t.value('shareLinkCopied'), icon: 'success' })
    }
  })
}

const formatTime = (t) => {
  if (!t) return ''
  return t.substring(0, 10)
}

onMounted(() => {
  loadInviteInfo()
  loadSubordinates()
})
</script>

<style lang="scss" scoped>
page { background-color: #000; }

.nf-page {
  min-height: 100vh;
  background: #000;
  position: relative;
}

.nf-bg {
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 500rpx;
  background: linear-gradient(180deg, rgba(229,9,20,0.3) 0%, transparent 100%);
}

.nf-navbar {
  position: relative;
  z-index: 10;
}

.nf-navbar-status {
  height: var(--status-bar-height, 44px);
}

.nf-navbar-content {
  display: flex;
  align-items: center;
  height: 88rpx;
  padding: 0 24rpx;
}

.nf-navbar-back {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.06);
  border: 1rpx solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.nf-navbar-title {
  flex: 1;
  text-align: center;
  font-size: 34rpx;
  font-weight: 600;
  color: #fff;
}

.nf-container {
  position: relative;
  z-index: 5;
  padding: 0 30rpx;
}

.nf-card {
  background: rgba(255,255,255,0.08);
  border-radius: 20rpx;
  padding: 36rpx;
  margin-bottom: 24rpx;
  backdrop-filter: blur(10px);
}

.invite-card {
  text-align: center;
}

.invite-title {
  font-size: 28rpx;
  color: rgba(255,255,255,0.6);
}

.invite-code {
  display: block;
  font-size: 56rpx;
  font-weight: 700;
  color: #e50914;
  letter-spacing: 8rpx;
  margin: 20rpx 0 30rpx;
}

.invite-btn-row {
  display: flex;
  gap: 20rpx;
}

.nf-btn {
  flex: 1;
  height: 80rpx;
  line-height: 80rpx;
  border-radius: 12rpx;
  font-size: 28rpx;
  text-align: center;
}

.nf-btn-primary {
  background: #e50914;
  color: #fff;
}

.nf-btn-ghost {
  background: transparent;
  border: 1px solid rgba(255,255,255,0.3);
  color: #fff;
}

.stats-card {
  .stats-row {
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .stats-item {
    flex: 1;
    text-align: center;
  }
  .stats-value {
    display: block;
    font-size: 48rpx;
    font-weight: 700;
    color: #fff;
  }
  .stats-label {
    font-size: 24rpx;
    color: rgba(255,255,255,0.5);
  }
  .stats-divider {
    width: 1px;
    height: 60rpx;
    background: rgba(255,255,255,0.15);
  }
}

.section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #fff;
  margin-bottom: 24rpx;
}

.empty-tip {
  text-align: center;
  padding: 40rpx 0;
  color: rgba(255,255,255,0.4);
  font-size: 26rpx;
}

.sub-item {
  display: flex;
  align-items: center;
  padding: 20rpx 0;
  border-bottom: 1px solid rgba(255,255,255,0.06);

  &:last-child {
    border-bottom: none;
  }
}

.sub-left {
  display: flex;
  align-items: center;
  flex: 1;
}

.sub-avatar {
  width: 72rpx;
  height: 72rpx;
  border-radius: 50%;
  margin-right: 20rpx;
}

.sub-info {
  display: flex;
  flex-direction: column;
}

.sub-name {
  font-size: 28rpx;
  color: #fff;
}

.sub-time {
  font-size: 22rpx;
  color: rgba(255,255,255,0.4);
  margin-top: 4rpx;
}

.load-more {
  text-align: center;
  padding: 20rpx 0;
  color: rgba(255,255,255,0.5);
  font-size: 26rpx;
}
</style>
