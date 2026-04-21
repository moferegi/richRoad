<template>
  <view class="nf-coupon-page">
    <view class="nf-coupon-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff"></uni-icons>
        </view>
        <text class="nf-navbar-title">{{ t('couponCenter') || '优惠券中心' }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <!-- 标签切换 -->
    <view class="nf-tabs">
      <view class="nf-tab" :class="{ active: activeTab === 'available' }" @tap="activeTab = 'available'">
        {{ t('couponAvailable') }}
      </view>
      <view class="nf-tab" :class="{ active: activeTab === 'used' }" @tap="activeTab = 'used'">
        {{ t('couponUsed') }}
      </view>
    </view>

    <!-- 空状态 -->
    <view class="nf-empty" v-if="filteredList.length === 0">
      <view class="nf-empty-icon">🎫</view>
      <text class="nf-empty-text">{{ t('couponEmpty') }}</text>
    </view>

    <!-- 优惠券列表 -->
    <view class="nf-coupon-list" v-else>
      <view
        v-for="(item, index) in filteredList"
        :key="index"
        class="nf-coupon-card"
        :style="couponBgStyle(item)"
      >
        <!-- 优惠券内容 -->
        <view class="nf-coupon-body">
          <view class="nf-coupon-left">
            <view class="nf-coupon-tag">{{ t('couponTag') }}</view>
            <text class="nf-coupon-name">{{ couponName(item) }}</text>
            <text class="nf-coupon-desc" v-if="couponDesc(item)">{{ couponDesc(item) }}</text>
            <text class="nf-coupon-date">{{ t('couponValidity') }}{{ item.startTime }} - {{ item.endTime }}</text>
          </view>
          <view class="nf-coupon-right" :class="{ disabled: item.status === 1 }">
            <text class="nf-coupon-symbol">{{ cs }}</text>
            <text class="nf-coupon-discount">{{ (item.discount / 100).toFixed(0) }}</text>
            <text class="nf-coupon-condition">{{ minSpendText(item) }}</text>
          </view>
        </view>

        <!-- 操作按钮 -->
        <view class="nf-coupon-footer" v-if="activeTab === 'available'">
          <view v-if="item.couponNum == 0" class="nf-coupon-btn nf-coupon-btn-claim" @tap="onClaim(item)">
            {{ t('couponClaim') }}
          </view>
          <view v-else-if="item.status === 0" class="nf-coupon-btn nf-coupon-btn-use" @tap="onUse(item)">
            {{ t('couponUse') }}
          </view>
        </view>
        <view class="nf-coupon-footer" v-else>
          <view class="nf-coupon-used-tag">{{ t('couponUsedStatus') }}</view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getAllClaimCoupon, claimCouponByUser } from '@/api/coupon'
import { t, localText } from '@/utils/i18n'
import { getUrl } from '@/utils/url'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'

const appConfigStore = useAppConfigStore()
const cs = computed(() => appConfigStore.currencySymbol)

const activeTab = ref('available')
const couponList = ref([])

const filteredList = computed(() => {
  if (activeTab.value === 'available') return couponList.value.filter(c => c.status !== 1)
  return couponList.value.filter(c => c.status === 1)
})

const couponName = (item) => {
  if (item.nameI18n) return localText(item.nameI18n) || item.name
  return item.name
}

const couponDesc = (item) => {
  if (item.descriptionI18n) return localText(item.descriptionI18n) || item.description || ''
  return item.description || ''
}

const minSpendText = (item) => {
  if (item.minSpend > 0) return t('couponMinSpend').replace('{n}', item.minSpend / 100)
  return t('couponNoLimit')
}

const couponBgStyle = (item) => {
  const url = getUrl(item.externalBgPath || item.backgroundImage || '')
  if (url) {
    return { backgroundImage: `url(${url})`, backgroundSize: 'cover', backgroundPosition: 'center' }
  }
  // Premium default gradient background for coupons without custom bg
  return {
    background: 'linear-gradient(135deg, rgba(229, 9, 20, 0.15) 0%, rgba(40, 10, 10, 0.6) 50%, rgba(20, 5, 5, 0.8) 100%)'
  }
}

const loadCoupons = async () => {
  try {
    const res = await getAllClaimCoupon({ goodIds: [] })
    if (res.code === 0) couponList.value = res.data || []
  } catch (e) { console.error('加载优惠券失败', e) }
}

const onClaim = async (item) => {
  try {
    const res = await claimCouponByUser({ couponID: item.couponID })
    if (res.code === 0) {
      uni.showToast({ title: t('couponClaimSuccess'), icon: 'success' })
      loadCoupons()
    } else {
      uni.showToast({ title: res.msg || t('couponClaimFail'), icon: 'none' })
    }
  } catch (e) { uni.showToast({ title: t('couponClaimFail'), icon: 'none' }) }
}

const onUse = () => { uni.switchTab({ url: '/pages/tabBar/index' }) }
const goBack = () => { uni.navigateBack() }

onMounted(() => { loadCoupons() })
</script>

<style lang="scss">
page { background-color: #000; }

.nf-coupon-page { min-height: 100vh; background: #000; position: relative; }

.nf-coupon-bg {
  position: fixed; top: 0; left: 0; right: 0; height: 500rpx; z-index: 0; pointer-events: none;
  background: radial-gradient(ellipse at 50% 0%, rgba(229, 9, 20, 0.10) 0%, transparent 60%);
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

/* 标签切换 */
.nf-tabs {
  display: flex; position: sticky; top: 0; z-index: 98;
  background: rgba(0, 0, 0, 0.9); backdrop-filter: blur(16px);
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
}
.nf-tab {
  flex: 1; text-align: center; padding: 24rpx 0;
  font-size: 28rpx; color: rgba(255, 255, 255, 0.5); font-weight: 500;
  position: relative;
  &.active {
    color: #fff; font-weight: 700;
    &::after {
      content: ''; position: absolute; bottom: 0;
      left: 50%; transform: translateX(-50%);
      width: 48rpx; height: 4rpx;
      background: #e50914; border-radius: 2rpx;
    }
  }
}

/* 空状态 */
.nf-empty {
  display: flex; flex-direction: column; align-items: center; justify-content: center; min-height: 60vh;
}
.nf-empty-icon { font-size: 120rpx; margin-bottom: 20rpx; opacity: 0.5; }
.nf-empty-text { font-size: 28rpx; color: rgba(255, 255, 255, 0.4); }

/* 优惠券列表 */
.nf-coupon-list { padding: 24rpx; position: relative; z-index: 1; }

.nf-coupon-card {
  border-radius: 20rpx; margin-bottom: 24rpx; overflow: hidden;
  border: 1rpx solid rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(8px);
}

.nf-coupon-body {
  display: flex; padding: 28rpx; gap: 20rpx;
}

.nf-coupon-left { flex: 1; display: flex; flex-direction: column; gap: 8rpx; }

.nf-coupon-tag {
  display: inline-block; width: fit-content;
  padding: 4rpx 14rpx; border-radius: 6rpx;
  background: #e50914; color: #fff;
  font-size: 20rpx; font-weight: 700;
}

.nf-coupon-name { font-size: 28rpx; font-weight: 700; color: #fff; line-height: 1.3; }
.nf-coupon-desc { font-size: 22rpx; color: rgba(255, 255, 255, 0.5); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.nf-coupon-date { font-size: 20rpx; color: rgba(255, 255, 255, 0.3); }

.nf-coupon-right {
  width: 160rpx; display: flex; flex-direction: column; align-items: center; justify-content: center;
  background: rgba(229, 9, 20, 0.15); border-radius: 16rpx;
  padding: 16rpx 8rpx; flex-shrink: 0;
  &.disabled { background: rgba(255, 255, 255, 0.06); }
  &.disabled .nf-coupon-symbol, &.disabled .nf-coupon-discount { color: rgba(255, 255, 255, 0.3); }
}
.nf-coupon-symbol { font-size: 24rpx; font-weight: 700; color: #e50914; }
.nf-coupon-discount { font-size: 52rpx; font-weight: 800; color: #e50914; line-height: 1; }
.nf-coupon-condition { font-size: 20rpx; color: rgba(255, 255, 255, 0.4); margin-top: 6rpx; text-align: center; }

.nf-coupon-footer {
  display: flex; justify-content: flex-end; padding: 0 28rpx 20rpx;
}

.nf-coupon-btn {
  padding: 12rpx 36rpx; border-radius: 30rpx;
  font-size: 24rpx; font-weight: 700; text-align: center;
  &:active { transform: scale(0.96); }
}
.nf-coupon-btn-claim { background: #e50914; color: #fff; }
.nf-coupon-btn-use {
  background: rgba(255, 255, 255, 0.06); border: 1rpx solid rgba(255, 255, 255, 0.15); color: #fff;
}
.nf-coupon-used-tag {
  font-size: 22rpx; color: rgba(255, 255, 255, 0.3);
  padding: 8rpx 20rpx; border-radius: 20rpx;
  background: rgba(255, 255, 255, 0.04);
}
</style>
