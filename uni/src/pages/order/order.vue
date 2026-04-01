<template>
  <view class="nf-order">
    <view class="nf-order-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff"></uni-icons>
        </view>
        <text class="nf-navbar-title">{{ $t('myOrders') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <!-- 订单状态标签 -->
    <view class="nf-tabs">
      <scroll-view scroll-x class="nf-tabs-scroll" :show-scrollbar="false">
        <view class="nf-tabs-inner">
          <view
            v-for="(item, index) in tabColumns"
            :key="index"
            class="nf-tab"
            :class="{ active: item.id === activeSataus }"
            @tap="tapBtn(item)"
          >
            <text>{{ item.title }}</text>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- 空状态 -->
    <view class="nf-empty" v-if="orderList.length === 0">
      <view class="nf-empty-icon">📦</view>
      <text class="nf-empty-text">{{ $t('noOrderData') }}</text>
    </view>

    <!-- 订单列表 -->
    <view class="nf-order-list" v-else>
      <view class="nf-order-card" v-for="(item, index) in orderList" :key="index">
        <!-- 订单头部 -->
        <view class="nf-order-header">
          <view class="nf-order-date">
            {{ item.CreatedAt.split('T')[0] }}
            <text v-if="item.isPresale" class="nf-presale-tag">{{ $t('presale') }}</text>
          </view>
          <view class="nf-order-status" :class="'nf-st-' + item.status">
            {{ tabColumns.find(tab => tab.id === item.status)?.title }}
            <text v-if="item.status === '0' && countdownMap[item.ID]" class="nf-countdown"> {{ countdownMap[item.ID] }}</text>
          </view>
        </view>

        <!-- 商品列表（多件横滑） -->
        <scroll-view scroll-x class="nf-goods-scroll" :show-scrollbar="false"
          v-if="item.detail && item.detail.length > 1">
          <view class="nf-goods-row">
            <image
              v-for="(d, di) in item.detail" :key="di"
              :src="getUrl(d.sku.picture)"
              class="nf-goods-thumb"
              mode="aspectFill"
            />
          </view>
        </scroll-view>

        <!-- 单件商品 -->
        <view class="nf-goods-single" v-if="item.detail && item.detail.length === 1">
          <image :src="getUrl(item.detail[0].sku.picture)" class="nf-goods-thumb-lg" mode="aspectFill" />
          <view class="nf-goods-single-info">
            <text class="nf-goods-single-name">{{ item.detail[0].sku.name }}</text>
            <text class="nf-goods-single-desc">{{ item.detail[0].sku.description }}</text>
          </view>
        </view>

        <!-- 金额统计 -->
        <view class="nf-order-summary">
          <text class="nf-order-count">共 {{ item.detail ? item.detail.length : 0 }} 件商品 · 实付</text>
          <text class="nf-order-price">¥{{ (item.totalPrice / 100).toFixed(2) }}</text>
        </view>

        <!-- 操作按钮 -->
        <view class="nf-order-actions">
          <view v-if="item.status === '4'" class="nf-action-tag cancelled">已取消</view>
          <view v-if="item.status === '0'" class="nf-action nf-action-ghost" @tap="cancelOrder(item)">取消订单</view>
          <view v-if="item.status === '0'" class="nf-action nf-action-primary" @tap="payOff(item.ID)">立即支付</view>
          <view v-if="item.status === '2'" class="nf-action nf-action-ghost" @tap="trackLogistics(item)">查看物流</view>
          <view v-if="canApplyRefund(item)" class="nf-action nf-action-warn" @tap="openRefund(item)">申请退款</view>
          <view v-else-if="item.status === '6'" class="nf-action-tag refunding">退款中</view>
          <view v-else-if="item.status === '5'" class="nf-action-tag refunded">已退款</view>
          <!-- 评价 -->
          <template v-if="item.status === '3' || item.status === '7'">
            <template v-if="item.detail && item.detail.length > 1">
              <view v-if="hasUncommentedItems(item)" class="nf-action nf-action-primary" @tap="goCommentAll(item)">评价订单</view>
              <view v-else-if="hasCommentedItems(item)" class="nf-action nf-action-ghost" @tap="goCommentAll(item)">查看评价</view>
            </template>
            <template v-else-if="item.detail && item.detail.length === 1">
              <view v-if="!item.detail[0].isComment" class="nf-action nf-action-primary" @tap="goComment(item, item.detail[0])">评价订单</view>
              <view v-else class="nf-action nf-action-ghost" @tap="goComment(item, item.detail[0])">查看评价</view>
            </template>
          </template>
          <view v-if="item.status === '2'" class="nf-action nf-action-primary" @tap="confirm(item)">确认收货</view>
        </view>
      </view>
    </view>

    <view style="height: 40rpx;"></view>

    <refund-apply-popup
      v-model:visible="refundVisible"
      :order-id="refundOrderId"
      @success="onRefundSuccess"
    />
  </view>
</template>

<script setup>
import { ref, computed, onUnmounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { updateOrderStatus, SelfOrderList } from "../../api/order"
import { getUrl } from "@/utils/url.js"
import RefundApplyPopup from '@/components/refund-apply-popup/refund-apply-popup.vue'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const activeSataus = ref("")
const tabColumns = ref([
  { title: '全部', id: '' },
  { title: '待付款', id: '0' },
  { title: '待发货', id: '1' },
  { title: '待收货', id: '2' },
  { title: '待评价', id: '3' },
  { title: '退款中', id: '6' },
  { title: '已退款', id: '5' },
  { title: '取消', id: '4' },
  { title: '已评价', id: '7' },
])

const orderList = ref([])
const refundVisible = ref(false)
const refundOrderId = ref('')

const init = async (params) => {
  activeSataus.value = params || ''
  const res = await SelfOrderList(activeSataus.value)
  if (res.code === 0) { orderList.value = res.data.list || [] }
  startCountdown()
}

const countdownMap = ref({})
let countdownTimer = null

const updateCountdowns = () => {
  const map = {}
  orderList.value.forEach(item => {
    if (item.status === '0' && item.closeTime) {
      const remain = Math.max(0, Math.floor((new Date(item.closeTime).getTime() - Date.now()) / 1000))
      if (remain > 0) {
        const m = Math.floor(remain / 60), s = remain % 60
        map[item.ID] = `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
      }
    }
  })
  countdownMap.value = map
}

const startCountdown = () => {
  if (countdownTimer) clearInterval(countdownTimer)
  updateCountdowns()
  countdownTimer = setInterval(updateCountdowns, 1000)
}

onUnmounted(() => { if (countdownTimer) clearInterval(countdownTimer) })
onLoad((options) => { init(options.status) })

const cancelOrder = (item) => {
  uni.showModal({
    title: "取消提示", content: "是否取消该订单？", confirmColor: "#e50914",
    success: async (res) => {
      if (res.confirm) {
        const r = await updateOrderStatus({ ID: item.ID, status: "4" })
        if (r.code === 0) { uni.showToast({ title: "取消成功", icon: "none" }); init() }
      }
    }
  })
}

const payOff = (ID) => { uni.navigateTo({ url: `/pages/orderInfo/orderInfo?orderID=${ID}` }) }
const trackLogistics = (item) => { uni.navigateTo({ url: `/pages/logistics/logistics?express=${item.express}` }) }
const canApplyRefund = (item) => ['1', '2', '3', '7'].includes(item.status)
const openRefund = (item) => { refundOrderId.value = item.ID; refundVisible.value = true }
const onRefundSuccess = () => { init(activeSataus.value) }

const confirm = async (item) => {
  const res = await updateOrderStatus({ ID: item.ID, status: '3' })
  if (res.code === 0) { uni.showToast({ title: "确认收货成功", icon: "none" }); init() }
}

const goComment = (order, detail) => {
  uni.navigateTo({ url: `/pages/evaluate/addEvaluate?orderID=${order.ID}&goodID=${detail.goodID}&SKUID=${detail.skuID}` })
}
const hasUncommentedItems = (order) => order.detail && order.detail.some(i => !i.isComment)
const hasCommentedItems = (order) => order.detail && order.detail.some(i => i.isComment)

const goCommentAll = (order) => {
  if (order.detail && order.detail.length > 1) {
    uni.navigateTo({ url: hasUncommentedItems(order)
      ? `/pages/evaluate/orderEvaluate?orderID=${order.ID}`
      : `/pages/evaluate/orderEvaluate?orderID=${order.ID}&mode=view`
    })
  } else if (order.detail && order.detail.length === 1) {
    const d = order.detail[0]
    uni.navigateTo({ url: `/pages/evaluate/addEvaluate?orderID=${order.ID}&goodID=${d.goodID}&SKUID=${d.skuID}` })
  }
}

const tapBtn = async (item) => {
  activeSataus.value = item.id
  const res = await SelfOrderList(item.id)
  if (res.code === 0) { orderList.value = res.data.list || [] }
}

const goBack = () => { uni.navigateBack() }
</script>

<style lang="scss">
page { background-color: #000; }

.nf-order { min-height: 100vh; background: #000; position: relative; }

.nf-order-bg {
  position: fixed; top: 0; left: 0; right: 0; height: 500rpx; z-index: 0; pointer-events: none;
  background: radial-gradient(ellipse at 30% 0%, rgba(229, 9, 20, 0.10) 0%, transparent 60%),
    radial-gradient(ellipse at 70% 10%, rgba(229, 9, 20, 0.06) 0%, transparent 50%);
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

.nf-tabs {
  position: sticky; top: 0; z-index: 98;
  background: rgba(0, 0, 0, 0.9); backdrop-filter: blur(16px);
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
}
.nf-tabs-scroll { white-space: nowrap; }
.nf-tabs-inner { display: inline-flex; padding: 0 16rpx; }
.nf-tab {
  display: inline-flex; padding: 20rpx 28rpx;
  font-size: 26rpx; color: rgba(255, 255, 255, 0.5); font-weight: 500;
  position: relative;
  &.active {
    color: #fff; font-weight: 700;
    &::after {
      content: ''; position: absolute; bottom: 0;
      left: 28rpx; right: 28rpx; height: 4rpx;
      background: #e50914; border-radius: 2rpx;
    }
  }
}

.nf-empty {
  display: flex; flex-direction: column; align-items: center; justify-content: center; min-height: 60vh;
}
.nf-empty-icon { font-size: 120rpx; margin-bottom: 20rpx; opacity: 0.5; }
.nf-empty-text { font-size: 28rpx; color: rgba(255, 255, 255, 0.4); }

.nf-order-list { padding: 20rpx 24rpx; position: relative; z-index: 1; }

.nf-order-card {
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.06);
  border-radius: 20rpx; margin-bottom: 20rpx;
  overflow: hidden; backdrop-filter: blur(8px);
}

.nf-order-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 24rpx 28rpx; border-bottom: 1rpx solid rgba(255, 255, 255, 0.04);
}
.nf-order-date {
  font-size: 24rpx; color: rgba(255, 255, 255, 0.5);
  display: flex; align-items: center; gap: 10rpx;
}
.nf-presale-tag {
  font-size: 20rpx; color: #fff;
  background: linear-gradient(135deg, #ff6b35, #e50914);
  padding: 2rpx 12rpx; border-radius: 6rpx;
}
.nf-order-status { font-size: 24rpx; font-weight: 600; }
.nf-st-0 { color: #e50914; }
.nf-st-1 { color: #22c55e; }
.nf-st-2 { color: #3b82f6; }
.nf-st-3 { color: #f59e0b; }
.nf-st-4 { color: rgba(255,255,255,0.3); }
.nf-st-5 { color: rgba(255,255,255,0.4); }
.nf-st-6 { color: #f59e0b; }
.nf-st-7 { color: #22c55e; }
.nf-countdown { font-size: 22rpx; color: #e50914; margin-left: 8rpx; }

.nf-goods-scroll { padding: 20rpx 28rpx; }
.nf-goods-row { display: inline-flex; gap: 12rpx; }
.nf-goods-thumb {
  width: 140rpx; height: 140rpx; border-radius: 12rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.06); flex-shrink: 0;
}
.nf-goods-single { display: flex; padding: 20rpx 28rpx; gap: 20rpx; }
.nf-goods-thumb-lg {
  width: 160rpx; height: 160rpx; border-radius: 12rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.06); flex-shrink: 0;
}
.nf-goods-single-info { flex: 1; display: flex; flex-direction: column; justify-content: center; }
.nf-goods-single-name {
  font-size: 28rpx; font-weight: 600; color: #fff; margin-bottom: 8rpx;
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.nf-goods-single-desc {
  font-size: 24rpx; color: rgba(255, 255, 255, 0.4);
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}

.nf-order-summary {
  display: flex; justify-content: flex-end; align-items: center;
  padding: 16rpx 28rpx; border-top: 1rpx solid rgba(255, 255, 255, 0.04); gap: 8rpx;
}
.nf-order-count { font-size: 24rpx; color: rgba(255, 255, 255, 0.4); }
.nf-order-price { font-size: 30rpx; font-weight: 700; color: #e50914; }

.nf-order-actions {
  display: flex; justify-content: flex-end; align-items: center;
  padding: 16rpx 28rpx 20rpx; gap: 12rpx; flex-wrap: wrap;
}
.nf-action {
  height: 60rpx; line-height: 58rpx; padding: 0 28rpx;
  border-radius: 30rpx; font-size: 24rpx; font-weight: 600;
  text-align: center; transition: all 0.2s;
  &:active { transform: scale(0.96); }
}
.nf-action-primary { background: #e50914; color: #fff; }
.nf-action-ghost {
  background: rgba(255, 255, 255, 0.06);
  border: 1rpx solid rgba(255, 255, 255, 0.12); color: rgba(255, 255, 255, 0.7);
}
.nf-action-warn {
  background: rgba(245, 158, 11, 0.15);
  border: 1rpx solid rgba(245, 158, 11, 0.3); color: #f59e0b;
}
.nf-action-tag {
  font-size: 22rpx; padding: 6rpx 20rpx; border-radius: 20rpx;
  &.cancelled { color: rgba(255,255,255,0.3); background: rgba(255,255,255,0.04); }
  &.refunding { color: #f59e0b; background: rgba(245, 158, 11, 0.1); }
  &.refunded { color: rgba(255,255,255,0.4); background: rgba(255,255,255,0.04); }
}
</style>
