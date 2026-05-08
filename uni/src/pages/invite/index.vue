<template>
  <view class="invite-page">
    <view class="invite-bg"></view>

    <view class="invite-nav">
      <view class="invite-nav-status"></view>
      <view class="invite-nav-content">
        <view class="invite-nav-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#0f172a" />
        </view>
        <text class="invite-nav-title">{{ $t('inviteFriends') }}</text>
        <view class="invite-nav-right-spacer"></view>
      </view>
    </view>

    <scroll-view class="invite-scroll" scroll-y>
      <view class="invite-body">
        <view class="invite-card code-card">
          <text class="code-title">{{ $t('myInviteCode') }}</text>
          <text class="code-value">{{ inviteInfo.inviteCode || '--' }}</text>
          <view class="code-actions">
            <view class="code-btn primary" @tap="copyCode">{{ $t('copyInviteCode') }}</view>
            <view class="code-btn" @tap="shareLink">{{ $t('shareLink') }}</view>
          </view>
        </view>

        <view class="invite-card stats-card">
          <view class="stats-grid">
            <view class="stats-item">
              <text class="stats-value">{{ inviteInfo.subordinateCount || 0 }}</text>
              <text class="stats-label">{{ $t('invited') }}</text>
            </view>
            <view class="stats-item">
              <text class="stats-value">{{ inviteRewardTotal }}</text>
              <text class="stats-label">{{ $t('inviteRewardTryonCoins') }}</text>
            </view>
            <view class="stats-item">
              <text class="stats-value">{{ inviteInfo.tryonPoint || 0 }}</text>
              <text class="stats-label">{{ $t('tryonCoins') }}</text>
            </view>
          </view>
          <text class="stats-tip">{{ $t('inviteRewardTimes') }}: {{ inviteRewardCount }}</text>
        </view>

        <view class="invite-card friends-card">
          <view class="friends-head">
            <text class="friends-title">{{ $t('myFriends') }}</text>
            <text class="friends-sub">{{ subordinateList.length }}/{{ total }}</text>
          </view>

          <view v-if="subordinateList.length === 0" class="empty-tip">
            <text>{{ $t('noInviteRecord') }}</text>
          </view>

          <view v-else>
            <view v-for="(item, index) in subordinateList" :key="item.ID || item.id || index" class="friend-item">
              <image class="friend-avatar" :src="item.avatar || defaultAvatar" mode="aspectFill" />
              <view class="friend-main">
                <text class="friend-name">{{ item.nickname || item.username }}</text>
                <text class="friend-time">{{ formatTime(item.CreatedAt || item.createdAt) }}</text>
              </view>
            </view>

            <view v-if="hasMore" class="load-more-btn" @tap="loadMore">{{ $t('loadMore') }}</view>
            <text v-else class="load-more-end">{{ $t('reachedBottom') }}</text>
          </view>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getMyInviteInfo, getMySubordinates } from '@/api/base.js'
import { request } from '@/utils/request.js'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const defaultAvatar = 'https://qmplusimg.henrongyi.top/gva_header.jpg'

const inviteInfo = ref({
  inviteCode: '',
  subordinateCount: 0,
  point: 0,
  tryonPoint: 0,
})

const inviteRewardStats = ref({
  totalTryonReward: 0,
  rewardCount: 0,
})

const subordinateList = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const hasMore = computed(() => subordinateList.value.length < total.value)
const inviteRewardTotal = computed(() => Number(inviteRewardStats.value.totalTryonReward || 0))
const inviteRewardCount = computed(() => Number(inviteRewardStats.value.rewardCount || 0))

const goBack = () => {
  uni.navigateBack({ fail: () => uni.switchTab({ url: '/pages/tabBar/my/index' }) })
}

const loadInviteInfo = async () => {
  const res = await getMyInviteInfo()
  if (res.code !== 0 || !res.data) return
  inviteInfo.value = {
    ...inviteInfo.value,
    ...res.data,
  }
}

const loadSubordinates = async () => {
  const res = await getMySubordinates({ page: page.value, pageSize: pageSize.value })
  if (res.code !== 0 || !res.data) return

  if (page.value === 1) {
    subordinateList.value = res.data.list || []
  } else {
    subordinateList.value.push(...(res.data.list || []))
  }
  total.value = Number(res.data.total || 0)
}

const loadInviteRewardStats = async () => {
  const stat = {
    totalTryonReward: 0,
    rewardCount: 0,
  }

  let nextPage = 1
  const maxPage = 30
  const fetchSize = 100

  while (nextPage <= maxPage) {
    const res = await request({
      url: '/cpr/getPointRecordList',
      method: 'get',
      params: {
        page: nextPage,
        pageSize: fetchSize,
        assetType: 'tryon_point',
        operationType: 'tryon_invite_register_reward',
        sort: 'created_at',
        order: 'descending',
      },
    })

    const list = Array.isArray(res?.data?.list) ? res.data.list : []
    list.forEach((item) => {
      const rawChange = Number(item?.pointChange ?? item?.PointChange ?? 0)
      if (!Number.isFinite(rawChange) || rawChange <= 0) return
      stat.totalTryonReward += rawChange
      stat.rewardCount += 1
    })

    const totalCount = Number(res?.data?.total || 0)
    if (list.length < fetchSize || nextPage * fetchSize >= totalCount) {
      break
    }
    nextPage += 1
  }

  inviteRewardStats.value = stat
}

const loadMore = () => {
  if (!hasMore.value) return
  page.value += 1
  loadSubordinates()
}

const copyCode = () => {
  if (!inviteInfo.value.inviteCode) return
  uni.setClipboardData({
    data: inviteInfo.value.inviteCode,
    success: () => {
      uni.showToast({ title: $t.value('inviteCodeCopied'), icon: 'none' })
    },
  })
}

const shareLink = () => {
  const code = String(inviteInfo.value.inviteCode || '').trim()
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
      uni.showToast({ title: $t.value('shareLinkCopied'), icon: 'none' })
    },
  })
}

const formatTime = (value) => {
  if (!value) return '--'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '--'
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

const refreshInvitePage = async () => {
  page.value = 1
  total.value = 0
  subordinateList.value = []
  inviteRewardStats.value = { totalTryonReward: 0, rewardCount: 0 }

  await Promise.all([
    loadInviteInfo(),
    loadSubordinates(),
    loadInviteRewardStats(),
  ])
}

onShow(() => {
  refreshInvitePage()
})
</script>

<style lang="scss" scoped>
page {
  background: #f4f7fb;
}

.invite-page {
  min-height: 100vh;
  background: radial-gradient(120% 80% at 100% -10%, #dbeafe 0%, transparent 62%), #f4f7fb;
  color: #0f172a;
}

.invite-bg {
  position: fixed;
  inset: 0 auto auto 0;
  width: 100%;
  height: 520rpx;
  pointer-events: none;
  background:
    radial-gradient(circle at 12% 10%, rgba(37, 99, 235, 0.16), transparent 52%),
    radial-gradient(circle at 90% 18%, rgba(14, 165, 233, 0.12), transparent 48%);
}

.invite-nav {
  position: sticky;
  top: 0;
  z-index: 99;
  padding: 0 24rpx 12rpx;
  backdrop-filter: blur(16rpx);
  background: rgba(244, 247, 251, 0.72);
  border-bottom: 1rpx solid rgba(15, 23, 42, 0.06);
}

.invite-nav-status {
  height: var(--status-bar-height, 44rpx);
}

.invite-nav-content {
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.invite-nav-back {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  border: 1rpx solid rgba(15, 23, 42, 0.12);
  background: rgba(255, 255, 255, 0.95);
  display: flex;
  align-items: center;
  justify-content: center;
}

.invite-nav-right-spacer {
  width: 64rpx;
}

.invite-nav-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #0f172a;
}

.invite-scroll {
  height: calc(100vh - var(--status-bar-height, 44rpx) - 100rpx);
}

.invite-body {
  position: relative;
  z-index: 1;
  padding: 20rpx;
}

.invite-card {
  border-radius: 18rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.08);
  background: rgba(255, 255, 255, 0.95);
  box-shadow: 0 12rpx 28rpx rgba(15, 23, 42, 0.06);
  padding: 20rpx;
  margin-bottom: 14rpx;
}

.code-card {
  text-align: center;
}

.code-title {
  display: block;
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.56);
}

.code-value {
  display: block;
  margin: 12rpx 0 16rpx;
  font-size: 56rpx;
  line-height: 1;
  letter-spacing: 6rpx;
  font-weight: 800;
  color: #1d4ed8;
}

.code-actions {
  display: flex;
  gap: 10rpx;
}

.code-btn {
  flex: 1;
  height: 68rpx;
  border-radius: 999rpx;
  font-size: 24rpx;
  color: #0f172a;
  background: rgba(241, 245, 249, 0.9);
  border: 1rpx solid rgba(15, 23, 42, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.code-btn.primary {
  color: #ffffff;
  border-color: transparent;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8rpx;
}

.stats-item {
  text-align: center;
  padding: 10rpx 6rpx;
  border-radius: 12rpx;
  background: rgba(241, 245, 249, 0.86);
}

.stats-value {
  display: block;
  font-size: 34rpx;
  font-weight: 700;
  color: #0f172a;
}

.stats-label {
  display: block;
  margin-top: 4rpx;
  font-size: 20rpx;
  color: rgba(15, 23, 42, 0.56);
}

.stats-tip {
  display: block;
  margin-top: 12rpx;
  text-align: center;
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.56);
}

.friends-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10rpx;
}

.friends-title {
  font-size: 28rpx;
  font-weight: 700;
  color: #0f172a;
}

.friends-sub {
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.52);
}

.empty-tip {
  padding: 54rpx 0;
  text-align: center;
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.5);
}

.friend-item {
  display: flex;
  align-items: center;
  gap: 12rpx;
  padding: 14rpx 0;
  border-bottom: 1rpx solid rgba(15, 23, 42, 0.06);
}

.friend-item:last-child {
  border-bottom: none;
}

.friend-avatar {
  width: 72rpx;
  height: 72rpx;
  border-radius: 50%;
  background: #e2e8f0;
}

.friend-main {
  flex: 1;
  min-width: 0;
}

.friend-name {
  display: block;
  font-size: 26rpx;
  color: #0f172a;
  font-weight: 600;
}

.friend-time {
  display: block;
  margin-top: 4rpx;
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.5);
}

.load-more-btn {
  margin-top: 12rpx;
  width: 260rpx;
  height: 62rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.14);
  background: rgba(241, 245, 249, 0.96);
  color: #0f172a;
  font-size: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-left: auto;
  margin-right: auto;
}

.load-more-end {
  display: block;
  margin-top: 12rpx;
  text-align: center;
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.46);
}
</style>
