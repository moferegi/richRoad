<template>
  <view class="nf-integral">
    <view class="nf-integral-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff"></uni-icons>
        </view>
        <text class="nf-navbar-title">{{ $t('myPoints') || '我的积分' }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <view class="nf-body">
      <!-- 积分概览卡片 -->
      <view class="nf-points-hero">
        <view class="nf-points-hero-glow"></view>
        <text class="nf-points-hero-label">可用积分</text>
        <text class="nf-points-hero-value">{{ userPoints }}</text>
        <text class="nf-points-hero-tip">积分可用于订单抵扣</text>
      </view>

      <!-- 积分记录列表 -->
      <view class="nf-section-title">
        <view class="nf-section-line"></view>
        <text>积分明细</text>
        <view class="nf-section-line"></view>
      </view>

      <view class="nf-empty" v-if="recordList.length === 0">
        <view class="nf-empty-icon">💎</view>
        <text class="nf-empty-text">暂无积分记录</text>
      </view>

      <view class="nf-record-list" v-else>
        <view class="nf-record-card" v-for="(item, index) in recordList" :key="index">
          <view class="nf-record-icon">{{ item.operationType === 'add' ? '⬆' : '⬇' }}</view>
          <view class="nf-record-info">
            <text class="nf-record-reason">{{ item.reason || '积分变动' }}</text>
            <text class="nf-record-time">{{ formatTime(item.CreatedAt) }}</text>
          </view>
          <text class="nf-record-amount" :class="item.operationType === 'add' ? 'nf-add' : 'nf-sub'">
            {{ item.operationType === 'add' ? '+' : '-' }}{{ item.points }}
          </text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getUserInfo } from '@/api/base.js'
import { request } from '@/utils/request.js'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const userPoints = ref(0)
const recordList = ref([])

const goBack = () => { uni.navigateBack() }

const loadUserPoints = async () => {
  try {
    const res = await getUserInfo()
    if (res.code === 0) userPoints.value = res.data.point || 0
  } catch (e) { console.error('获取积分失败', e) }
}

const loadRecords = async () => {
  try {
    const res = await request({
      url: '/cpr/getPointRecordList',
      method: 'get',
      params: { page: 1, pageSize: 50 }
    })
    if (res.code === 0 && res.data && res.data.list) {
      recordList.value = res.data.list
    }
  } catch (e) { console.error('获取积分记录失败', e) }
}

const formatTime = (t) => {
  if (!t) return ''
  const d = new Date(t)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

onMounted(() => {
  loadUserPoints()
  loadRecords()
})
</script>

<style lang="scss">
page { background-color: #000; }

.nf-integral { min-height: 100vh; background: #000; position: relative; }

.nf-integral-bg {
  position: fixed; top: 0; left: 0; right: 0; height: 600rpx; z-index: 0; pointer-events: none;
  background: radial-gradient(ellipse at 50% 0%, rgba(229, 9, 20, 0.12) 0%, transparent 60%),
    radial-gradient(ellipse at 80% 10%, rgba(229, 9, 20, 0.06) 0%, transparent 40%);
}

.nf-navbar {
  background: rgba(0, 0, 0, 0.85); backdrop-filter: blur(24px);
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
  padding: 0 28rpx 16rpx; position: sticky; top: 0; z-index: 99;
}
.nf-navbar-status { height: var(--status-bar-height, 0px); }
.nf-navbar-content { display: flex; align-items: center; justify-content: space-between; height: 88rpx; }
.nf-navbar-back {
  width: 64rpx; height: 64rpx; border-radius: 50%;
  background: rgba(255, 255, 255, 0.06); border: 1rpx solid rgba(255, 255, 255, 0.1);
  display: flex; align-items: center; justify-content: center;
}
.nf-navbar-title { font-size: 34rpx; font-weight: 700; color: #fff; letter-spacing: 2rpx; }

.nf-body { padding: 24rpx; position: relative; z-index: 1; }

/* 积分概览 */
.nf-points-hero {
  position: relative; text-align: center; padding: 60rpx 40rpx;
  background: rgba(255, 255, 255, 0.04); border: 1rpx solid rgba(255, 255, 255, 0.06);
  border-radius: 24rpx; margin-bottom: 40rpx; overflow: hidden;
}
.nf-points-hero-glow {
  position: absolute; top: -50%; left: 50%; transform: translateX(-50%);
  width: 300rpx; height: 300rpx; border-radius: 50%;
  background: radial-gradient(circle, rgba(229, 9, 20, 0.20) 0%, transparent 70%);
  pointer-events: none;
}
.nf-points-hero-label { display: block; font-size: 26rpx; color: rgba(255, 255, 255, 0.5); margin-bottom: 12rpx; }
.nf-points-hero-value {
  display: block; font-size: 80rpx; font-weight: 800; color: #fff;
  letter-spacing: 4rpx; line-height: 1; margin-bottom: 16rpx;
}
.nf-points-hero-tip { display: block; font-size: 22rpx; color: rgba(255, 255, 255, 0.3); }

/* 区块标题 */
.nf-section-title {
  display: flex; align-items: center; justify-content: center;
  gap: 20rpx; margin-bottom: 24rpx;
  text { font-size: 26rpx; color: rgba(255, 255, 255, 0.5); }
}
.nf-section-line {
  height: 1rpx; width: 80rpx;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.15), transparent);
}

/* 空状态 */
.nf-empty {
  display: flex; flex-direction: column; align-items: center; padding: 100rpx 0;
}
.nf-empty-icon { font-size: 100rpx; margin-bottom: 16rpx; opacity: 0.5; }
.nf-empty-text { font-size: 26rpx; color: rgba(255, 255, 255, 0.3); }

/* 记录列表 */
.nf-record-list { }

.nf-record-card {
  display: flex; align-items: center; gap: 20rpx;
  padding: 24rpx 28rpx;
  background: rgba(255, 255, 255, 0.04); border: 1rpx solid rgba(255, 255, 255, 0.06);
  border-radius: 16rpx; margin-bottom: 12rpx;
}
.nf-record-icon {
  width: 64rpx; height: 64rpx; border-radius: 50%;
  background: rgba(255, 255, 255, 0.06); display: flex; align-items: center; justify-content: center;
  font-size: 28rpx; flex-shrink: 0;
}
.nf-record-info { flex: 1; }
.nf-record-reason { display: block; font-size: 28rpx; color: #fff; font-weight: 500; margin-bottom: 4rpx; }
.nf-record-time { display: block; font-size: 22rpx; color: rgba(255, 255, 255, 0.3); }
.nf-record-amount { font-size: 32rpx; font-weight: 700; }
.nf-add { color: #22c55e; }
.nf-sub { color: #e50914; }
</style>
