<template>
  <view class="nf-my">
    <!-- 背景光效 -->
    <view class="nf-my-bg"></view>

    <!-- 头部个人信息区域 -->
    <view class="nf-profile">
      <view class="nf-profile-glow"></view>

      <!-- 已登录 -->
      <view class="nf-user" v-if="isShow">
        <view class="nf-avatar-wrap">
          <button class="nf-avatar-btn" open-type="chooseAvatar" @chooseavatar="onChooseAvatar">
            <image class="nf-avatar" :src="info.avatar || avatarUrl"></image>
          </button>
          <view class="nf-avatar-ring"></view>
        </view>
        <view class="nf-user-info">
          <text class="nf-username">{{ info.nickname || $t('defaultUser') }}</text>
          <input name="nickName" type="nickname" :placeholder="$t('nicknamePlaceholder')" class="nf-nickname-input" @change="onInputNickName" v-model="info.nickname" />
        </view>
      </view>

      <!-- 未登录 -->
      <view class="nf-unlogin" v-if="!isShow" @tap="logins">
        <view class="nf-unlogin-icon">
          <text class="nf-unlogin-icon-text">→</text>
        </view>
        <text class="nf-unlogin-text">{{ $t('loginPrompt') }}</text>
      </view>

      <!-- 数据统计 -->
      <view class="nf-stats">
        <view class="nf-stat-item">
          <text class="nf-stat-value">0</text>
          <text class="nf-stat-label">{{ $t('coupon') }}</text>
        </view>
        <view class="nf-stat-divider"></view>
        <view class="nf-stat-item">
          <text class="nf-stat-value">20</text>
          <text class="nf-stat-label">{{ $t('points') }}</text>
        </view>
      </view>
    </view>

    <!-- 浏览历史 -->
    <view class="nf-section">
      <view class="nf-section-header">
        <view class="nf-section-icon">
          <image class="nf-section-icon-img" src="/static/history.png" mode="aspectFill" />
        </view>
        <text class="nf-section-title">{{ $t('browseHistory') }}</text>
        <view class="nf-section-line"></view>
      </view>
      <scroll-view class="nf-history-scroll" scroll-x>
        <view class="nf-history-list">
          <view
              class="nf-history-item"
              v-for="(item, index) in historyList"
              :key="index"
              @tap="goto(item)"
          >
            <image class="nf-history-img" :src="item && item.imageUrl ? getUrl(item.imageUrl) : ''" mode="aspectFill"></image>
            <view class="nf-history-overlay"></view>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- 功能菜单 -->
    <view class="nf-menu" v-if="columns.length >= 1">
      <template v-for="(item, index) in columns.filter(c => !c.hidden)" :key="index">
        <button v-if="item.pages === 'contact'" class="nf-menu-item nf-contact-btn" open-type="contact">
          <view class="nf-menu-left">
            <view class="nf-menu-icon-wrap">
              <image class="nf-menu-icon" :src="item.icon" mode="aspectFit"></image>
            </view>
            <text class="nf-menu-text">{{ item.title }}</text>
          </view>
          <view class="nf-menu-arrow">
            <uni-icons color="rgba(255,255,255,0.3)" type="forward" size="14"></uni-icons>
          </view>
        </button>
        <view v-else class="nf-menu-item" @tap="toPages(item.pages)">
          <view class="nf-menu-left">
            <view class="nf-menu-icon-wrap">
              <image class="nf-menu-icon" :src="item.icon" mode="aspectFit"></image>
            </view>
            <text class="nf-menu-text">{{ item.title }}</text>
          </view>
          <view class="nf-menu-arrow">
            <uni-icons color="rgba(255,255,255,0.3)" type="forward" size="14"></uni-icons>
          </view>
        </view>
      </template>
    </view>

    <!-- 底部留白 -->
    <view style="height: 60rpx;"></view>

    <!-- 语言切换弹窗 -->
    <lang-switch v-model="showLangPicker" @change="onLangChange" />
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import {getUrl} from "@/utils/url";
import { myRouter } from "@/utils/permission";
import { setClientUserInfo } from "@/api/base";
import {useUserStore} from "@/pinia/modules/user.js"
import { useLangStore } from '@/pinia/modules/lang.js'
import { onShow } from '@dcloudio/uni-app'
import { getGoodHistory } from '@/api/order.js'
import langSwitch from '@/components/lang-switch/lang-switch.vue'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const showLangPicker = ref(false)
const onLangChange = () => {}
const openLangPicker = () => { showLangPicker.value = true }

const isShow = ref(false)

const columns = computed(() => [
  {
    title: $t.value('myCollection'),
    pages: '/pages/collect/collect',
    icon: '/static/MYcollect.png'
  },{
    title: $t.value('onlineService'),
    pages: 'contact',
    icon: '/static/wx.png'
  },{
    title: $t.value('logout'),
    pages: 'exit',
    icon: '/static/images/exit.png',
    hidden: !isShow.value
  },{
    title: $t.value('switchLang'),
    pages: 'lang',
    icon: '/static/images/exit.png'
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
  getHistory()
})

// 获取浏览历史的方法
const getHistory = async () => {
  const res = await getGoodHistory()
  if (res.code === 0) {
    historyList.value = res.data
  } else {
    uni.showToast({
      title: res.msg,
      icon: 'none'
    })
  }
}

// 跳转到商品详情页
const goto = (item) => {
  console.log(item);
  if (!item || !item.ID) {
    uni.showToast({
      title: $t.value('goodsInfoIncomplete'),
      icon: 'none'
    })
    return
  }
  myRouter(`/pages/goodsDetails/goodsDetails?id=${item.ID}`, true)
}

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
  if(pages === 'lang'){
    openLangPicker()
    return
  }
  if(pages === 'exit'){
    userStore.loginOut()
    uni.showToast({
      icon: 'none',
      title: $t.value('logoutSuccess')
    })

    return
  }
  if (!pages) {
    uni.showToast({
      icon: 'none',
      title: $t.value('developing')
    })
  } else {
    // 如果是收藏页面，使用switchTab切换到tab
    if (pages === '/pages/collect/collect') {
      uni.switchTab({
        url: pages
      })
    } else {
      uni.navigateTo({
        url: pages
      })
    }
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
  background-color: #000;
}

.nf-my {
  width: 100%;
  min-height: 100vh;
  background: #000;
  position: relative;
}

/* ===== 背景光效 ===== */
.nf-my-bg {
  position: fixed;
  top: 0; left: 0; right: 0;
  height: 600rpx;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(ellipse at 50% -10%, rgba(229, 9, 20, 0.18) 0%, transparent 55%),
    radial-gradient(ellipse at 80% 30%, rgba(229, 9, 20, 0.08) 0%, transparent 45%);
}

/* ===== 个人信息头部 ===== */
.nf-profile {
  position: relative;
  z-index: 1;
  margin: 0 28rpx;
  padding: 60rpx 36rpx 40rpx;
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.06);
  border-radius: 32rpx;
  margin-top: calc(var(--status-bar-height, 44px) + 20rpx);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  overflow: hidden;
}

.nf-profile-glow {
  position: absolute;
  top: -60rpx; left: 50%;
  transform: translateX(-50%);
  width: 300rpx;
  height: 300rpx;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(229, 9, 20, 0.2) 0%, transparent 70%);
  pointer-events: none;
}

/* 已登录用户 */
.nf-user {
  display: flex;
  align-items: center;
  gap: 28rpx;
  margin-bottom: 36rpx;
}

.nf-avatar-wrap {
  position: relative;
  width: 130rpx;
  height: 130rpx;
  flex-shrink: 0;
}

.nf-avatar-btn {
  padding: 0;
  margin: 0;
  background: none;
  border: none;
  width: 130rpx;
  height: 130rpx;
  line-height: normal;
  &::after { border: none; }
}

.nf-avatar {
  width: 130rpx;
  height: 130rpx;
  border-radius: 50%;
  border: 3rpx solid rgba(229, 9, 20, 0.5);
}

.nf-avatar-ring {
  position: absolute;
  top: -6rpx; left: -6rpx;
  width: 142rpx;
  height: 142rpx;
  border-radius: 50%;
  border: 2rpx solid rgba(229, 9, 20, 0.3);
  pointer-events: none;
}

.nf-user-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 12rpx;
  overflow: hidden;
}

.nf-username {
  font-size: 38rpx;
  font-weight: 800;
  color: #fff;
  letter-spacing: 2rpx;
}

.nf-nickname-input {
  background: rgba(255, 255, 255, 0.06);
  border: 1rpx solid rgba(255, 255, 255, 0.1);
  border-radius: 12rpx;
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.7);
  padding: 12rpx 16rpx;
  width: 100%;
  height: auto;
  min-height: 56rpx;
  line-height: 1.4;
  box-sizing: border-box;
  overflow: visible;
  text-overflow: ellipsis;
  transition: border-color 0.3s;

  &:focus {
    border-color: rgba(229, 9, 20, 0.4);
  }
}

/* 未登录 */
.nf-unlogin {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 20rpx 0;
  margin-bottom: 36rpx;
}

.nf-unlogin-icon {
  width: 88rpx;
  height: 88rpx;
  border-radius: 50%;
  background: linear-gradient(135deg, #e50914, #b20710);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 6rpx 24rpx rgba(229, 9, 20, 0.4);
}

.nf-unlogin-icon-text {
  font-size: 36rpx;
  color: #fff;
  font-weight: bold;
}

.nf-unlogin-text {
  color: #fff;
  font-size: 34rpx;
  font-weight: 700;
  letter-spacing: 2rpx;
}

/* 数据统计 */
.nf-stats {
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.06);
  border-radius: 20rpx;
  padding: 28rpx 0;
}

.nf-stat-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
}

.nf-stat-value {
  font-size: 40rpx;
  font-weight: 800;
  color: #e50914;
  line-height: 1;
}

.nf-stat-label {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.5);
  letter-spacing: 2rpx;
}

.nf-stat-divider {
  width: 1rpx;
  height: 56rpx;
  background: rgba(255, 255, 255, 0.1);
}

/* ===== 通用板块 ===== */
.nf-section {
  position: relative;
  z-index: 1;
  margin: 24rpx 28rpx;
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.06);
  border-radius: 28rpx;
  padding: 32rpx;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
}

.nf-section-header {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-bottom: 28rpx;
}

.nf-section-icon {
  width: 48rpx;
  height: 48rpx;
  border-radius: 12rpx;
  background: rgba(229, 9, 20, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.nf-section-icon-img {
  width: 32rpx;
  height: 32rpx;
}

.nf-section-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #fff;
  letter-spacing: 2rpx;
}

.nf-section-line {
  flex: 1;
  height: 1rpx;
  background: linear-gradient(90deg, rgba(229, 9, 20, 0.3), transparent);
}

/* 浏览历史 */
.nf-history-scroll {
  width: 100%;
}

.nf-history-list {
  display: flex;
  gap: 16rpx;
  padding: 4rpx 0;
}

.nf-history-item {
  position: relative;
  flex-shrink: 0;
  width: 200rpx;
  height: 260rpx;
  border-radius: 16rpx;
  overflow: hidden;
  transition: transform 0.3s;

  &:active {
    transform: scale(0.96);
  }
}

.nf-history-img {
  width: 100%;
  height: 100%;
  display: block;
}

.nf-history-overlay {
  position: absolute;
  bottom: 0; left: 0; right: 0;
  height: 50%;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.6));
  pointer-events: none;
}

/* ===== 功能菜单 ===== */
.nf-menu {
  position: relative;
  z-index: 1;
  margin: 24rpx 28rpx;
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.06);
  border-radius: 28rpx;
  padding: 8rpx 32rpx;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
}

.nf-menu-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 32rpx 0;
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
  transition: background 0.3s;

  &:last-child {
    border-bottom: none;
  }

  &:active {
    background: rgba(255, 255, 255, 0.03);
  }
}

.nf-menu-left {
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.nf-menu-icon-wrap {
  width: 48rpx;
  height: 48rpx;
  border-radius: 12rpx;
  background: rgba(255, 255, 255, 0.06);
  display: flex;
  align-items: center;
  justify-content: center;
}

.nf-menu-icon {
  width: 28rpx;
  height: 28rpx;
  opacity: 0.8;
}

.nf-menu-text {
  font-size: 30rpx;
  color: #fff;
  font-weight: 500;
  letter-spacing: 1rpx;
}

.nf-menu-arrow {
  opacity: 0.5;
}

.nf-contact-btn {
  background: transparent;
  margin: 0;
  padding: 32rpx 0;
  line-height: normal;
  border-radius: 0;
  text-align: left;
  font-size: inherit;
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);

  &::after { border: none; }
  &:last-child { border-bottom: none; }
}
</style>
