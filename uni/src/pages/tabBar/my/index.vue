<template>
  <view class="user-center">
    <!-- 顶部用户信息 -->
    <view class="user-header">
      <view class="user-info" @click="handleLogin">
        <image
            class="avatar"
            :src="isLogin ? userInfo.avatar : defaultAvatar"
            mode="aspectFill"
        ></image>
        <text class="nickname">{{ isLogin ? userInfo.nickname : '点击登录' }}</text>
      </view>
    </view>

    <!-- 用户数据 -->
    <view class="user-data">
      <view class="data-item">
        <text class="number">{{ userInfo.balance || '0.00' }}</text>
        <text class="label">余额</text>
      </view>
      <view class="data-item">
        <text class="number">{{ userInfo.coupons || 0 }}</text>
        <text class="label">优惠券</text>
      </view>
      <view class="data-item">
        <text class="number">{{ userInfo.points || 0 }}</text>
        <text class="label">积分</text>
      </view>
    </view>

    <!-- 订单导航 -->
    <view class="order-section">
      <view class="order-types">
        <view class="type-item" v-for="(item, index) in orderTypes" :key="index">
          <image class="order-icon" :src="item.iconUrl" mode="aspectFill" />
          <text class="name">{{ item.name }}</text>
        </view>
      </view>
    </view>

    <!-- 浏览历史 -->
    <view class="history-section">
      <view class="section-header">
        <image class="history-icon" src="/static/history.png" mode="aspectFill" />
        <text class="title">浏览历史</text>
      </view>
      <scroll-view class="history-scroll" scroll-x>
        <view class="history-list">
          <view
              class="history-item"
              v-for="(item, index) in historyList"
              :key="index"
          >
            <view class="history-image" :style="{ backgroundColor: '#f5f5f5' }"></view>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- 功能列表 -->
    <view class="feature-list">
      <view class="feature-item" @tap="navigateTo('/pages/address/list')">
        <view class="feature-left">
          <view class="feature-icon location" :style="{ backgroundColor: '#4CAF50' }"></view>
          <text class="feature-name">地址管理</text>
        </view>
        <wu-icon name="arrow-right"></wu-icon>
      </view>
      <view class="feature-item" @tap="navigateTo('/pages/collect/list')">
        <view class="feature-left">
          <view class="feature-icon star" :style="{ backgroundColor: '#FF9800' }"></view>
          <text class="feature-name">我的收藏</text>
        </view>
        <wu-icon name="arrow-right"></wu-icon>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue';

// 登录状态
const isLogin = ref(false);

// 默认头像
const defaultAvatar = 'https://example.com/default-avatar.png';

// 用户信息
const userInfo = ref({
  avatar: 'https://example.com/user-avatar.png',
  nickname: 'Leo yo',
  balance: '128.8',
  coupons: 0,
  points: 20
});

// 订单类型
const orderTypes = ref([
  { name: '全部订单', iconUrl: '/static/all.png' },
  { name: '待付款', iconUrl: '/static/pay.png' },
  { name: '待收货', iconUrl: '/static/delivery.png' },
  { name: '退款/售后', iconUrl: '/static/refound.png' }
]);

// 浏览历史
const historyList = ref(Array(5).fill(null));

// 处理登录
const handleLogin = () => {
  if (!isLogin.value) {
    isLogin.value = true;
    // 这里后期替换为实际的登录逻辑
    console.log('执行登录');
  }
};
</script>

<style lang="scss">
.user-center {
  min-height: 100vh;
  background-color: #f8f8f8;
}

.user-header {
  height: 200rpx;
  background: linear-gradient(to bottom right, #ff6b6b, #ff4757);
  padding: 60rpx 30rpx;
  border-radius: 0 0 30rpx 30rpx;
}

.user-info {
  display: flex;
  align-items: center;
}

.avatar {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  border: 4rpx solid #fff;
}

.nickname {
  margin-left: 20rpx;
  color: #fff;
  font-size: 32rpx;
  font-weight: bold;
}

.user-data {
  margin: -60rpx 30rpx 0;
  border-radius: 20rpx;
  padding: 30rpx;
  display: flex;
  justify-content: space-around;
  background: #ffffff;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.1);
}

.data-item {
  text-align: center;
}

.number {
  font-size: 36rpx;
  color: #333;
  font-weight: bold;
  margin-bottom: 10rpx;
  display: block;
}

.label {
  font-size: 24rpx;
  color: #666;
}

.order-section, .history-section {
  margin: 30rpx;
  background-color: #fff;
  border-radius: 20rpx;
}
.history-section{
  padding: 30rpx;
}
.history-icon{
  width: 50rpx;
  height: 50rpx;
  border-radius: 50%;
  margin-right: 20rpx;
}
.section-header {
  display: flex;
  margin-bottom: 30rpx;
}

.title {
  font-size: 32rpx;
  color: #333;
  font-weight: bold;
}

.more {
  font-size: 24rpx;
  color: #999;
}


.order-types {
  display: flex;
  justify-content: space-around; /* 改为space-around，确保每个项目有均等的空间 */
  padding: 20rpx 0; /* 添加上下内边距 */
}

.type-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1; /* 添加flex: 1，确保每个项目占据相等的空间 */
  text-align: center; /* 确保文字居中 */
}

.order-icon {
  width: 60rpx;
  height: 60rpx;
  border-radius: 12rpx;
  margin-bottom: 10rpx; /* 将margin-right改为margin-bottom */
  overflow: hidden;
}

.name {
  font-size: 24rpx;
  color: #666;
  width: 100%; /* 确保文字容器占满宽度 */
  text-align: center; /* 确保文字居中 */
}

.icon {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  margin-bottom: 10rpx;
}


.history-scroll {
  width: 100%;
}

.history-list {
  display: flex;
  padding: 10rpx 0;
}

.history-item {
  margin-right: 20rpx;
  flex-shrink: 0;
}

.history-image {
  width: 200rpx;
  height: 200rpx;
  border-radius: 12rpx;
}

.feature-list {
  margin: 30rpx;
  background-color: #fff;
  border-radius: 20rpx;
  overflow: hidden; /* 确保圆角不被子元素溢出 */
}

.feature-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 30rpx;
  border-bottom: 1rpx solid #f5f5f5;

  &:last-child {
    border-bottom: none;
  }
}

.feature-left {
  display: flex;
  align-items: center;
  flex: 1;
}

.feature-icon {
  width: 60rpx;
  height: 60rpx;
  border-radius: 50%;
  margin-right: 20rpx;
  flex-shrink: 0; /* 防止图标被压缩 */
}

.feature-name {
  font-size: 28rpx;
  color: #333;
}

</style>
