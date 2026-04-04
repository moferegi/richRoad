<template>
  <view class="nf-signin">
    <view class="nf-signin-bg"></view>
    <view class="nf-signin-card">
      <view class="nf-signin-header">
        <text class="nf-signin-title">{{ $t('signInTitle') }}</text>
        <text class="nf-signin-sub">{{ $t('signInConsecutive').replace('{n}', status.consecutiveDays || 0) }}</text>
      </view>

      <!-- 签到按钮 -->
      <view class="nf-signin-btn-wrap">
        <view
          class="nf-signin-btn"
          :class="{ 'nf-signin-done': status.signedToday }"
          @tap="onSignIn"
        >
          <text class="nf-signin-btn-text">
            {{ status.signedToday ? $t('signedToday') : $t('signInNow') }}
          </text>
        </view>
      </view>

      <!-- 奖励提示 -->
      <view class="nf-signin-reward" v-if="status.rewardPoints > 0">
        <text class="nf-signin-reward-text">
          {{ $t('signInReward').replace('{n}', status.rewardPoints) }}
        </text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { doSignIn, getSignInStatus } from '@/api/signIn.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { onShow } from '@dcloudio/uni-app'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const status = ref({
  signedToday: false,
  consecutiveDays: 0,
  rewardPoints: 0,
})

const loadStatus = async () => {
  const res = await getSignInStatus()
  if (res.code === 0) {
    status.value = res.data
  }
}

const onSignIn = async () => {
  if (status.value.signedToday) {
    uni.showToast({ title: $t.value('alreadySigned'), icon: 'none' })
    return
  }
  const res = await doSignIn()
  if (res.code === 0) {
    uni.showToast({ title: $t.value('signInSuccess'), icon: 'success' })
    loadStatus()
  }
}

onShow(() => {
  loadStatus()
})
</script>

<style scoped>
.nf-signin {
  min-height: 100vh;
  background: #141414;
  padding: 40rpx;
  position: relative;
}
.nf-signin-bg {
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 400rpx;
  background: linear-gradient(180deg, rgba(229,9,20,0.3) 0%, transparent 100%);
}
.nf-signin-card {
  position: relative;
  margin-top: 120rpx;
  background: rgba(255,255,255,0.05);
  border-radius: 24rpx;
  padding: 60rpx 40rpx;
  text-align: center;
}
.nf-signin-title {
  font-size: 44rpx;
  font-weight: bold;
  color: #fff;
  display: block;
}
.nf-signin-sub {
  font-size: 26rpx;
  color: rgba(255,255,255,0.5);
  margin-top: 10rpx;
  display: block;
}
.nf-signin-btn-wrap {
  margin-top: 60rpx;
}
.nf-signin-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 300rpx;
  height: 100rpx;
  border-radius: 50rpx;
  background: linear-gradient(135deg, #e50914, #b20710);
}
.nf-signin-btn.nf-signin-done {
  background: rgba(255,255,255,0.1);
}
.nf-signin-btn-text {
  font-size: 32rpx;
  font-weight: bold;
  color: #fff;
}
.nf-signin-reward {
  margin-top: 30rpx;
}
.nf-signin-reward-text {
  font-size: 24rpx;
  color: #e5b875;
}
</style>
