<template>
  <view class="order-container">
    <!-- 顶部导航栏 -->
    <view class="order-tabs">
      <scroll-view scroll-x class="nav-scroll" show-scrollbar="false">
        <view class="tab-container">
          <view
              v-for="(item, index) in tabColumns"
              :key="index"
              class="tab-item"
              :class="{ active: item.id === activeSataus }"
              @tap="tapBtn(item)"
          >
            {{ item.title }}
            <view class="tab-line" v-if="item.id === activeSataus"></view>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- 空状态展示 -->
    <view class="empty-state" v-if="orderList.length === 0">
      <view class="empty-image-container">
          <image class="empty-image-placeholder" src="./../../static/emptyStatus.png"></image>
      </view>
      <view class="empty-text">暂无订单数据</view>
    </view>

    <!-- 订单列表 -->
    <view class="order-list" v-else>
      <view
          v-for="(item, index) in orderList"
          :key="index"
          class="order-card"
      >
        <!-- 订单时间和状态 -->
        <view class="order-header">
          <view class="order-time">2019-04-06 11:37</view>
          <view class="order-status" :class="{'status-pending': item.status === '0', 'status-closed': item.status === '4'}">
            {{ tabColumns.find(tab=>tab.id === item.status)?.title }}
          </view>
        </view>

        <!-- 商品图片滑动区域 -->
        <scroll-view
            scroll-x
            class="goods-images-scroll"
            show-scrollbar="false"
            v-if="item.detail && item.detail.length > 0"
        >
          <view class="goods-images-container">
            <image
                v-for="(detail, detailIndex) in item.detail"
                :key="detailIndex"
                :src="getUrl(detail.sku.picture)"
                class="goods-thumbnail"
                mode="aspectFill"
            ></image>
          </view>
        </scroll-view>

        <!-- 订单商品统计信息 -->
        <view class="order-summary">
          <view class="total-count">共 {{ item.detail ? item.detail.length : 0 }} 件商品 实付款</view>
          <view class="total-price">¥ {{
              item.detail ? item.detail.reduce((total, curr) => total + (curr.quantity * curr.price / 100), 0).toFixed(1) : 0
            }}</view>
        </view>

        <!-- 订单操作按钮 -->
        <view class="order-actions">
          <view class="left-actions">
            <button
                v-if="item.status==='4'"
                class="action-btn delete-btn"
            >
              <text class="btn-icon">×</text>
              删除订单
            </button>
          </view>
          <view class="right-actions">
            <button
                v-if="item.status==='0'"
                class="action-btn cancel-btn"
                @tap="cancelOrder(item)"
            >取消订单</button>
            <button
                v-if="item.status==='0'"
                class="action-btn pay-btn"
                @tap="payOff(item.ID)"
            >立即支付</button>
            <button
                v-if="item.status==='2'"
                class="action-btn track-btn"
                @tap="trackLogistics(item)"
            >查看物流</button>
            <button
                v-if="item.status==='2'"
                class="action-btn confirm-btn"
                @tap="confirm(item)"
            >确认收货</button>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import {onLoad} from '@dcloudio/uni-app'
import {updateOrderStatus, SelfOrderList} from "../../api/order";
import {getUrl} from "@/utils/url.js"
const activeSataus = ref("")

const tabColumns = ref([{
  title: '全部',
  id: ''
}, {
  title: '待付款',
  id: '0'
}, {
  title: '待发货',
  id: '1'
}, {
  title: '待收货',
  id: '2'
}, {
  title: '待评价',
  id: '3'
},  {
  title: '取消',
  id: '4'
}, {
  title: '已评价',
  id: '7'
}])
// 获取订单列表
const orderList = ref([])
const init = async (params) => {
  activeSataus.value = params || ''
  const res = await SelfOrderList(activeSataus.value)
  if (res.code === 0) {
    orderList.value = res.data.list || []
  }
}

onLoad((options) => {
  init(options.status)
})

const cancelOrder = (item) => {
  uni.showModal({
    title: "取消提示",
    content: "是否取消该订单？",
    confirmColor: "#fe5572",
    success: async function (res) {
      if (res.confirm) {
        const req = {
          ID: item.ID,
          status: "4"
        }
        const res = await updateOrderStatus(req)
        if (res.code === 0) {
          uni.showToast({
            title: "取消订单成功",
            icon: "none"
          });
          tabColumns.value.forEach((item) => {
            item.active = false
          })
          init()
        }
      }
    }
  });
}
const payOff = (ID) => {
  uni.navigateTo({
    url: `/pages/orderInfo/orderInfo?orderID=${ID}`,
  })
}

const trackLogistics = async (item) => {
  uni.navigateTo({
    url: `/pages/logistics/logistics?express=${item.express}`
  })

}

// 确认收货
const confirm = async (item) => {
  const req = {
    ID: item.ID,
    status: '3'
  }
  const res = await updateOrderStatus(req)
  if (res.code === 0) {
    uni.showToast({
      title: "确认收货成功",
      icon: "none"
    });
    init()
  }
}

const goComment = (order, detail) => {
  uni.navigateTo({
    url: `/pages/evaluate/addEvaluate?orderID=${order.ID}&goodID=${detail.goodID}&SKUID=${detail.skuID}`
  })
}

const tapBtn = async (item) => {
  activeSataus.value = item.id
  const res = await SelfOrderList(item.id)
  if (res.code === 0) {
    orderList.value = res.data.list || []
  }
}

</script>

<style lang="scss">
page {
  background-color: #f7f7f7;
  font-family: -apple-system, BlinkMacSystemFont, 'Helvetica Neue', Helvetica, sans-serif;
}

/* 导航栏样式 */
.order-tabs {
  position: sticky;
  top: 0;
  left: 0;
  right: 0;
  background-color: #fff;
  z-index: 100;
  /* #ifdef H5 */
  top: var(--window-top);
  /* #endif */
}

.nav-scroll {
  white-space: nowrap;
  width: 100%;
}

.tab-container {
  display: flex;
  background-color: #fff;
}

.tab-item {
  position: relative;
  font-size: 28rpx;
  color: #666;
  padding: 20rpx 30rpx;
  display: inline-block;
}

.tab-item.active {
  color: #fe5572;
  font-weight: normal;
}

.tab-line {
  position: absolute;
  bottom: 0;
  left: 20rpx;
  right: 20rpx;
  height: 4rpx;
  background-color: #fe5572;
}

/* 空状态样式 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: calc(100vh - 88rpx);
  /* #ifdef H5 */
  height: calc(100vh - 88rpx - var(--window-top));
  /* #endif */
}

.empty-image-container {
  margin-bottom: 30rpx;
}

.empty-image-placeholder {
  width: 280rpx;
  height: 280rpx;
  border-radius: 8rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #999;
}

/* 订单列表样式 */
.order-list {
  padding: 16rpx;
}

.order-card {
  background-color: #fff;
  margin-bottom: 20rpx;
}

/* 订单头部 */
.order-header {
  display: flex;
  justify-content: space-between;
  padding: 20rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.order-time {
  font-size: 24rpx;
  color: #666;
}

.order-status {
  font-size: 24rpx;
  color: #fe5572;
}

.status-closed {
  color: #999;
}

/* 商品图片滑动区域 */
.goods-images-scroll {
  padding: 20rpx;
  white-space: nowrap;
}

.goods-images-container {
  display: inline-flex;
}

.goods-thumbnail {
  width: 140rpx;
  height: 140rpx;
  margin-right: 16rpx;
}

/* 订单商品统计 */
.order-summary {
  display: flex;
  justify-content: flex-end;
  padding: 20rpx;
  font-size: 24rpx;
  align-items: center;
  border-bottom: 1rpx solid #f0f0f0;
}

.total-count {
  color: #666;
  margin-right: 10rpx;
}

.total-price {
  color: #fe5572;
  font-weight: bold;
}

/* 订单操作按钮 */
.order-actions {
  display: flex;
  justify-content: space-between;
  padding: 20rpx;
}

.right-actions {
  display: flex;
  justify-content: flex-end;
}

.action-btn {
  height: 60rpx;
  line-height: 58rpx;
  font-size: 26rpx;
  padding: 0 24rpx;
  border-radius: 30rpx;
  margin-left: 16rpx;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
}

.btn-icon {
  margin-right: 4rpx;
  font-size: 28rpx;
}

.delete-btn {
  border: 1rpx solid #ccc;
  color: #666;
  background-color: #fff;
}

.cancel-btn {
  border: 1rpx solid #ccc;
  color: #666;
  background-color: #fff;
}

.pay-btn {
  border: 1rpx solid #fe5572;
  background-color: #fe5572;
  color: #fff;
}

.track-btn {
  border: 1rpx solid #ccc;
  color: #666;
  background-color: #fff;
}

.confirm-btn {
  border: 1rpx solid #fe5572;
  background-color: #fe5572;
  color: #fff;
}
</style>
