<template>
  <view class="nf-integral">
    <view class="nf-integral-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#0f172a"></uni-icons>
        </view>
        <text class="nf-navbar-title">{{ pageTitle }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <view class="nf-body">
      <!-- 积分概览卡片 -->
      <view class="nf-points-hero">
        <view class="nf-points-hero-glow"></view>
        <text class="nf-points-hero-label">{{ heroLabel }}</text>
        <text class="nf-points-hero-value">{{ userPoints }}</text>
        <text class="nf-points-hero-tip">{{ heroTip }}</text>
      </view>

      <!-- 积分记录列表 -->
      <view class="nf-section-title">
        <view class="nf-section-line"></view>
        <text>{{ detailTitle }}</text>
        <view class="nf-section-line"></view>
      </view>

      <view class="nf-empty" v-if="recordList.length === 0 && !loading">
        <view class="nf-empty-icon">💎</view>
        <text class="nf-empty-text">{{ emptyText }}</text>
      </view>

      <view class="nf-record-list" v-else>
        <view class="nf-record-card" v-for="(item, index) in recordList" :key="index">
          <view class="nf-record-icon" :class="item.changeType === 'increase' ? 'increase' : 'decrease'">
            {{ item.changeType === 'increase' ? '+' : '-' }}
          </view>
          <view class="nf-record-info">
            <text class="nf-record-reason">{{ translateReason(item) }}</text>
            <text class="nf-record-time">{{ formatTime(item.CreatedAt) }}</text>
          </view>
          <text class="nf-record-amount" :class="item.changeType === 'increase' ? 'nf-add' : 'nf-sub'">
            {{ item.pointChange > 0 ? '+' : '' }}{{ item.pointChange }}
          </text>
        </view>
      </view>

      <view class="nf-load-more" v-if="recordList.length > 0">
        <text class="nf-load-more-text" v-if="loading">...</text>
        <text class="nf-load-more-text" v-else-if="noMore">{{ $t('reachedBottom') }}</text>
        <text class="nf-load-more-text" v-else-if="reachedBottom" @tap="loadMore">{{ $t('loadMore') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onLoad, onReachBottom, onShow } from '@dcloudio/uni-app'
import { getUserInfo } from '@/api/base.js'
import { request } from '@/utils/request.js'
import { localText, t as i18nT } from '@/utils/i18n.js'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)
const mode = ref('point')
const assetType = ref('point')

const REASON_KEYS = [
  'reason_newUserReward',
  'reason_orderReward',
  'reason_subUserReward',
  'reason_pointDeduct',
  'reason_orderCancelRefund',
  'reason_signInReward',
  'reason_orderCancelDeduct',
  'reason_tryonTaskDeduct',
  'reason_tryonTaskRefund',
  'reason_tryonGuestInit',
  'reason_tryonRegisterReward',
  'reason_tryonInviteReward',
]

const OPERATION_REASON_KEY_MAP = {
  tryon_consume: 'reason_tryonTaskDeduct',
  tryon_refund: 'reason_tryonTaskRefund',
  tryon_guest_init: 'reason_tryonGuestInit',
  tryon_register_reward: 'reason_tryonRegisterReward',
  tryon_invite_register_reward: 'reason_tryonInviteReward',
}

const REASON_ALIAS_LANGS = ['zh', 'zh-TW', 'en', 'mn', 'th', 'hi', 'id']

const buildReasonAliasMap = () => {
  const aliasMap = {}

  REASON_KEYS.forEach((key) => {
    aliasMap[key] = key

    REASON_ALIAS_LANGS.forEach((lang) => {
      const text = String(i18nT(key, lang) || '').trim()
      if (!text) return
      aliasMap[text] = key
      aliasMap[text.replace(/\s+/g, '')] = key
    })
  })

  return aliasMap
}

const reasonAliasMap = buildReasonAliasMap()

const translateReason = (record) => {
  const reason = record?.reason ?? record?.Reason
  const operationType = String(record?.operationType ?? record?.OperationType ?? '').trim().toLowerCase()

  if (reason && typeof reason === 'object') {
    const text = localText(reason, langStore.locale)
    if (text) return text
  }

  let rawReason = String(reason || '').trim()
  if (!rawReason) return $t.value('pointsChange')

  if ((rawReason.startsWith('{') || rawReason.startsWith('['))) {
    try {
      const parsed = JSON.parse(rawReason)
      const text = localText(parsed, langStore.locale)
      if (text) return text
      rawReason = String(text || '').trim() || rawReason
    } catch {
      // ignore parse error
    }
  }

  const normalizedReason = rawReason.replace(/\s+/g, '')
  const reasonKey = reasonAliasMap[rawReason] || reasonAliasMap[normalizedReason]
  if (reasonKey) return $t.value(reasonKey)

  if (operationType && OPERATION_REASON_KEY_MAP[operationType]) {
    return $t.value(OPERATION_REASON_KEY_MAP[operationType])
  }

  return rawReason
}

const userPoints = ref(0)
const recordList = ref([])
const loading = ref(false)
const noMore = ref(false)
const reachedBottom = ref(false)
const page = ref(1)
const pageSize = 10

const isTryonMode = computed(() => mode.value === 'tryon' || assetType.value === 'tryon_point')
const pageTitle = computed(() => isTryonMode.value ? $t.value('tryonPointRecord') : $t.value('myPoints'))
const heroLabel = computed(() => isTryonMode.value ? $t.value('availableTryonCoins') : $t.value('availablePoints'))
const heroTip = computed(() => isTryonMode.value ? $t.value('tryonCoinDeductTip') : $t.value('pointsDeductTip'))
const detailTitle = computed(() => isTryonMode.value ? $t.value('tryonPointDetails') : $t.value('pointsDetail'))
const emptyText = computed(() => isTryonMode.value ? $t.value('noTryonPointRecord') : $t.value('noPointsRecord'))

const goBack = () => { uni.navigateBack() }

const loadUserPoints = async () => {
  try {
    const res = await getUserInfo()
    if (res.code === 0) {
      userPoints.value = isTryonMode.value
        ? Number((res.data?.tryonPoint ?? 0) || 0)
        : Number((res.data?.point ?? 0) || 0)
    }
  } catch (e) { console.error('Failed to load points', e) }
}

const loadRecords = async (isLoadMore = false) => {
  if (loading.value) return
  loading.value = true
  try {
    const res = await request({
      url: '/cpr/getPointRecordList',
      method: 'get',
      params: {
        page: page.value,
        pageSize,
        sort: 'created_at',
        order: 'descending',
        assetType: assetType.value,
      }
    })
    if (res.code === 0 && res.data && res.data.list) {
      if (isLoadMore) {
        recordList.value = [...recordList.value, ...res.data.list]
      } else {
        recordList.value = res.data.list
      }
      if (res.data.list.length < pageSize || recordList.value.length >= (res.data.total || Infinity)) {
        noMore.value = true
      }
    }
  } catch (e) { console.error('Failed to load point records', e) }
  loading.value = false
}

const loadMore = () => {
  if (noMore.value || loading.value) return
  reachedBottom.value = false
  page.value++
  loadRecords(true)
}

const resetAndLoad = () => {
  page.value = 1
  noMore.value = false
  reachedBottom.value = false
  recordList.value = []
  loadUserPoints()
  loadRecords(false)
}

const formatTime = (t) => {
  if (!t) return ''
  const d = new Date(t)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

onLoad((options = {}) => {
  const nextMode = String(options.mode || '').trim().toLowerCase()
  const nextAssetType = String(options.assetType || '').trim().toLowerCase()
  mode.value = nextMode === 'tryon' ? 'tryon' : 'point'
  assetType.value = nextAssetType === 'tryon_point' ? 'tryon_point' : (mode.value === 'tryon' ? 'tryon_point' : 'point')
  resetAndLoad()
})

onShow(() => {
  if (recordList.value.length === 0) return
  resetAndLoad()
})

onReachBottom(() => {
  if (noMore.value || loading.value) return
  reachedBottom.value = true
})
</script>

<style lang="scss">
page {
  background: #f4f7fb;
}

.nf-integral {
  min-height: 100vh;
  background: radial-gradient(120% 80% at 100% -10%, #dbeafe 0%, transparent 62%), #f4f7fb;
  position: relative;
}

.nf-integral-bg {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 620rpx;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 8% 10%, rgba(59, 130, 246, 0.12), transparent 50%),
    radial-gradient(circle at 92% 20%, rgba(14, 165, 233, 0.1), transparent 46%);
}

.nf-navbar {
  position: sticky;
  top: 0;
  z-index: 99;
  padding: 0 24rpx 12rpx;
  backdrop-filter: blur(16rpx);
  background: rgba(244, 247, 251, 0.72);
  border-bottom: 1rpx solid rgba(15, 23, 42, 0.06);
}

.nf-navbar-status {
  height: var(--status-bar-height, 0px);
}

.nf-navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 88rpx;
}

.nf-navbar-back {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  border: 1rpx solid rgba(15, 23, 42, 0.1);
  background: rgba(255, 255, 255, 0.95);
  display: flex;
  align-items: center;
  justify-content: center;
}

.nf-navbar-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #0f172a;
}

.nf-body {
  position: relative;
  z-index: 1;
  padding: 24rpx;
}

.nf-points-hero {
  position: relative;
  text-align: center;
  padding: 56rpx 36rpx;
  margin-bottom: 30rpx;
  border-radius: 22rpx;
  border: 1rpx solid rgba(37, 99, 235, 0.14);
  background: linear-gradient(135deg, rgba(219, 234, 254, 0.92), rgba(239, 246, 255, 0.98));
  box-shadow: 0 16rpx 32rpx rgba(59, 130, 246, 0.12);
  overflow: hidden;
}

.nf-points-hero-glow {
  position: absolute;
  top: -30%;
  left: 50%;
  transform: translateX(-50%);
  width: 360rpx;
  height: 360rpx;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(37, 99, 235, 0.18) 0%, transparent 70%);
  pointer-events: none;
}

.nf-points-hero-label {
  display: block;
  margin-bottom: 10rpx;
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.64);
}

.nf-points-hero-value {
  display: block;
  margin-bottom: 12rpx;
  font-size: 80rpx;
  line-height: 1;
  letter-spacing: 2rpx;
  font-weight: 800;
  color: #0f172a;
}

.nf-points-hero-tip {
  display: block;
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.55);
}

.nf-section-title {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16rpx;
  margin-bottom: 22rpx;

  text {
    font-size: 24rpx;
    color: rgba(15, 23, 42, 0.55);
  }
}

.nf-section-line {
  width: 72rpx;
  height: 1rpx;
  background: linear-gradient(90deg, transparent, rgba(15, 23, 42, 0.2), transparent);
}

.nf-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 90rpx 0;
}

.nf-empty-icon {
  font-size: 90rpx;
  margin-bottom: 14rpx;
}

.nf-empty-text {
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.5);
}

.nf-record-card {
  display: flex;
  align-items: center;
  gap: 18rpx;
  margin-bottom: 12rpx;
  padding: 22rpx 24rpx;
  border-radius: 16rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.08);
  background: rgba(255, 255, 255, 0.95);
  box-shadow: 0 10rpx 24rpx rgba(15, 23, 42, 0.06);
}

.nf-record-icon {
  width: 62rpx;
  height: 62rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 30rpx;
  font-weight: 700;
}

.nf-record-icon.increase {
  color: #16a34a;
  background: rgba(22, 163, 74, 0.14);
}

.nf-record-icon.decrease {
  color: #dc2626;
  background: rgba(220, 38, 38, 0.14);
}

.nf-record-info {
  flex: 1;
  min-width: 0;
}

.nf-record-reason {
  display: block;
  margin-bottom: 4rpx;
  font-size: 26rpx;
  color: #0f172a;
  font-weight: 600;
}

.nf-record-time {
  display: block;
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.48);
}

.nf-record-amount {
  font-size: 32rpx;
  font-weight: 700;
}

.nf-add {
  color: #16a34a;
}

.nf-sub {
  color: #dc2626;
}

.nf-load-more {
  text-align: center;
  padding: 22rpx 0 60rpx;
}

.nf-load-more-text {
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.5);
}
</style>
