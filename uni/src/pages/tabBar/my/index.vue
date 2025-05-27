<template>
  <view class="my-page">
    <!-- 头部个人信息区域 -->
    <view class="header-section">
      <!-- 已登录状态 -->
      <view class="user-info" v-if="isShow">
        <button class="avatar-wrapper" open-type="chooseAvatar" @chooseavatar="onChooseAvatar">
          <image class="avatar" :src="info.avatar || avatarUrl"></image>
        </button>
        <view class="user-details">
          <text class="username">{{ info.nickname || '默认用户' }}</text>
          <input name="nickName" type="nickname" placeholder="请填写昵称" class="nickname-input" @change="onInputNickName" v-model="info.nickname" />
        </view>
      </view>

      <!-- 未登录状态 -->
      <view class="unlogin-box" v-if="!isShow" @tap="logins">
        <text class="login-prompt">请登录后查看</text>
      </view>

      <!-- 账户信息 -->
      <view class="account-info">
        <view class="account-item">
          <text class="account-value">0</text>
          <text class="account-label">优惠券</text>
        </view>
        <view class="account-item">
          <text class="account-value">20</text>
          <text class="account-label">积分</text>
        </view>
      </view>
    </view>

    <!-- 我的订单区域 -->
    <view class="order-section">
      <view class="order-title">
        <text>我的订单</text>
        <view class="view-all" @tap="toOrder('all')">
          <text>查看全部</text>
          <uni-icons color="#999999" type="forward" size="14"></uni-icons>
        </view>
      </view>

      <view class="order-items">
        <view class="order-item" @tap="toNav(`/pages/order/order?status=${0}`)">
          <image src="../../../static/pay.png" mode="aspectFit"></image>
          <text>待付款</text>
        </view>
        <view class="order-item" @tap="toNav(`/pages/order/order?status=${1}`)">
          <image src="../../../static/delivery.png" mode="aspectFit"></image>
          <text>待发货</text>
        </view>
        <view class="order-item" @tap="toNav(`/pages/order/order?status=${2}`)">
          <image src="../../../static/wait.png" mode="aspectFit"></image>
          <text>待收货</text>
        </view>
        <view class="order-item" @tap="toNav(`/pages/order/order?status=${3}`)">
          <image src="../../../static/refound.png" mode="aspectFit"></image>
          <text>退款/售后</text>
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

    <!-- 其他功能区域 -->
    <view class="tools-section" v-if="columns.length >=1">
      <template v-for="(item, index) in columns.filter(c=>!c.hidden)">
        <view class="tool-item" @tap="toPages(item.pages)" :class="columns.length-1 == index ? '' : 'border-bottom'">
          <view class="tool-left">
            <image :src="item.icon" mode="aspectFit"></image>
            <text>{{item.title}}</text>
          </view>
          <uni-icons color="#999999" type="forward" size="14"></uni-icons>
        </view>
      </template>
    </view>
  </view>
</template>

<script setup>
import { myRouter } from "@/utils/permission";
import { setClientUserInfo } from "@/api/base";
import {useUserStore} from "@/pinia/modules/user.js"
import { onShow } from '@dcloudio/uni-app'
const isShow = ref(false)
import { ref } from 'vue'
const columns = ref([
  {
    title: '地址管理',
    pages: '/pages/address/address',
    icon: '/static/address.png'
  },{
    title: '我的收藏',
    pages: '/pages/collect/collect',
    icon: '/static/MYcollect.png'
  },{
    title: '退出登录',
    pages: 'exit',
    icon: '/static/images/exit.png',
    hidden: true
  }
])
// 浏览历史
const historyList = ref(Array(5).fill(null));
const defaultAvatarUrl = 'https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0'
const avatarUrl = ref('')
const nickName = ref('默认用户')
avatarUrl.value = defaultAvatarUrl

const info = ref({})
const userStore = useUserStore()
onShow(() => {
  const token = userStore.token || ''
  if (token) {
    isShow.value = true
    info.value = uni.getStorageSync("userInfo")
  } else {
    isShow.value = false
  }
  const c = columns.value.find(item=>item.pages==='exit')
  if(c){
    c.hidden = !isShow.value
  }
})

const logins = () => {
  uni.redirectTo({
    url: '/pages/user/login'
  })
}

// 修改头像的方法
const onChooseAvatar = async (e) => {
  avatarUrl.value = e.detail.avatarUrl
  let tmpFilePath = avatarUrl.value
  // 对微信返回的临时图片链接进行base64编码
  var avatarUrl_base64 = 'data:image/jpeg;base64,' + wx.getFileSystemManager().readFileSync(tmpFilePath, 'base64');
  const res = await setClientUserInfo({
    key: "avatar",
    value: avatarUrl_base64
  })
  if(res.code === 0) {
    uni.showToast({
      title: res.msg,
      icon: 'success'
    })
  } else {
    uni.showToast({
      title: res.msg,
      icon: 'none'
    })
  }
}

// 修改昵称的方法
const onInputNickName = async (e) => {
  nickName.value = e.detail.value
  const res = await setClientUserInfo({
    key: "nickname",
    value: nickName.value
  })
  if (res.code === 0) {
    uni.showToast({
      title: res.msg,
      icon: 'success'
    })
  } else {
    uni.showToast({
      title: res.msg,
      icon: 'none'
    })
  }
}


const toPages = (pages) => {
  if(pages === 'exit'){
    userStore.loginOut()
    uni.showToast({
      icon: 'none',
      title: '退出成功'
    })

    return
  }
  if (!pages) {
    uni.showToast({
      icon: 'none',
      title: '正在开发中...'
    })
  } else {
    uni.navigateTo({
      url: pages
    })
  }
}

const toOrder = (val) => {
  myRouter("/pages/order/order",true)
}
const toNav = (pages) => {
  uni.navigateTo({
    url: pages
  })
}
const toRetreatOrder = () => {
  uni.navigateTo({
    url: '/pages/retreat/retreatOrder'
  })
}
</script>

<style lang="scss">
page {
  background-color: #f8f8f8;
}

.my-page {
  width: 100%;
  min-height: 100vh;
}

/* 头部样式 */
.header-section {
  position: relative;
  padding: 20rpx 30rpx;
  background: linear-gradient(135deg, #fea94a, #ff8c55);
  border-radius: 0 0 20rpx 20rpx;
  overflow: hidden;
}

/* 用户信息样式 */
.user-info {
  display: flex;
  align-items: center;
  padding: 20rpx 0;
}

.avatar-wrapper {
  padding: 0;
  margin: 0;
  background: none;
  border: none;
  width: auto;
  line-height: normal;
  &::after {
    border: none;
  }
}

.avatar {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  background-color: #fff;
  margin-right: 20rpx;
}

.user-details {
  display: flex;
  flex-direction: column;
}

.username {
  font-size: 36rpx;
  color: #fff;
  font-weight: bold;
  margin-bottom: 10rpx;
}

.nickname-input {
  background: rgba(255, 255, 255, 0.3);
  border-radius: 8rpx;
  font-size: 24rpx;
  color: #fff;
  padding: 6rpx 10rpx;
  width: 200rpx;
}

/* 未登录样式 */
.unlogin-box {
  height: 120rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-prompt {
  color: #fff;
  font-size: 32rpx;
  font-weight: bold;
}

/* VIP卡片样式 */
.vip-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: rgba(0, 0, 0, 0.15);
  border-radius: 12rpx;
  padding: 16rpx 24rpx;
  margin-top: 20rpx;
}

.vip-title {
  display: flex;
  align-items: center;
}

.vip-icon {
  width: 40rpx;
  height: 40rpx;
  margin-right: 10rpx;
}

.vip-title text {
  color: #fff;
  font-size: 28rpx;
}

.vip-btn {
  background-color: #3c3c3c;
  color: #e9cd77;
  font-size: 24rpx;
  padding: 6rpx 16rpx;
  border-radius: 30rpx;
  line-height: normal;
  min-height: 0;
}

/* 账户信息样式 */
.account-info {
  display: flex;
  justify-content: space-between;
  padding: 30rpx 20rpx;
  margin-top: 20rpx;
  background-color: #ffffff;
  border-radius: 12rpx;
}

.account-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
}

.account-value {
  font-size: 32rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 8rpx;
}

.account-label {
  font-size: 24rpx;
  color: #999;
}

/* 订单区域样式 */
.order-section {
  margin: 20rpx;
  padding: 30rpx;
  background-color: #fff;
  border-radius: 12rpx;
  box-shadow: 0 4rpx 8rpx rgba(0, 0, 0, 0.05);
}

.order-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30rpx;
}

.order-title text {
  font-size: 32rpx;
  font-weight: bold;
  color: #333;
}

.view-all {
  display: flex;
  align-items: center;
}

.view-all text {
  font-size: 24rpx;
  color: #999;
  margin-right: 8rpx;
}

.order-items {
  display: flex;
  justify-content: space-between;
}

.order-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
}

.order-item image {
  width: 60rpx;
  height: 60rpx;
  margin-bottom: 16rpx;
}

.order-item text {
  font-size: 24rpx;
  color: #333;
}

/* 浏览历史样式 */
.history-section {
  margin: 20rpx;
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


.section-title {
  display: flex;
  align-items: center;
  margin-bottom: 20rpx;
}

.title-icon {
  width: 40rpx;
  height: 40rpx;
  margin-right: 10rpx;
}

.section-title text {
  font-size: 30rpx;
  color: #333;
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
.history-items {
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
}

.history-items image {
  width: 150rpx;
  height: 150rpx;
  margin-bottom: 10rpx;
  border-radius: 8rpx;
}

/* 钱包区域样式 */
.wallet-section {
  margin: 20rpx;
  padding: 30rpx;
  background-color: #fff;
  border-radius: 12rpx;
  box-shadow: 0 4rpx 8rpx rgba(0, 0, 0, 0.05);
}

.wallet-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx 0 0;
}

.wallet-info text {
  font-size: 26rpx;
  color: #666;
}

/* 工具区域样式 */
.tools-section {
  margin: 20rpx;
  padding: 0 30rpx;
  background-color: #fff;
  border-radius: 12rpx;
  box-shadow: 0 4rpx 8rpx rgba(0, 0, 0, 0.05);
}

.tool-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 30rpx 0;
}

.border-bottom {
  border-bottom: 1rpx solid #f0f0f0;
}

.tool-left {
  display: flex;
  align-items: center;
}

.tool-left image {
  width: 40rpx;
  height: 40rpx;
  margin-right: 20rpx;
}

.tool-left text {
  font-size: 28rpx;
  color: #333;
}

.tool-right {
  display: flex;
  align-items: center;
}

.share-info {
  font-size: 24rpx;
  color: #999;
  margin-right: 10rpx;
}
</style>
