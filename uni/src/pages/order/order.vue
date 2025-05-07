<template>
  <view style="padding-bottom: 100rpx;">
    <view class="order_nav_view">
      <view class="bgc_fff order_nav_box pos_f">
        <view class="flexr-jsa flex-aic font_28 color_333 p_t_16">
          <view v-for="(item,index) in tabColumns" @tap="tapBtn(item)">
            <text class="font_bold" :class="[item.id === activeSataus?'color: color_fe5572': 'color: #000']">{{ item.title }}</text>
            <view class="line_box line_active"></view>
          </view>
        </view>
      </view>
    </view>
    <view v-for="(item, index) in orderList" :key="index" class="order_goods_card">
      <view class="bgc_fff order_goods_card_item">
        <view class="color_fe5572 font_28 m_b_24">
			<span>{{ tabColumns.find(tab=>tab.id === item.status)?.title}}</span>
        </view>
        <view v-for="detail in item.detail">
          <view class="flex m_b_24 m_t_16">
            <image :src="getUrl(detail.sku.picture)"
                   class="order_goods_card_img m_r_24"
                   mode="aspectFill"></image>
            <view class="flex-fitem">
              <view class="color_333 font_32 text_nowrap" style="max-width: 420rpx;">{{ detail.sku.name }}</view>
              <view class="font_24 color_999 m_b_24">
                <view class="color_b7bed0 m_t_24 flex">
                  <text v-for="(sku, index) in detail.sku.attrs" :key="index"
                        class="text_nowrap color_b7bed0 m_b_16 texts">{{ sku.label }}:{{ sku.value }}
                  </text>
                </view>
                <view class="flex flex-aic" style="justify-content: space-between;">
                  <text class="font_28">x{{ detail.quantity }}</text>
                  <view class="font_40 color_ff0003">
                    <span class="font_28">¥</span>
                    {{ detail.quantity * detail.price / 100 }}</view>
                </view>
              </view>
              <view class="btn-box">
                <button v-if="(item.status==='3'||item.status==='7')&& !detail.isComment" class="comment-box redBtn" @tap="goComment(item, detail)">评价</button>
                <button v-if="(item.status==='3'||item.status==='7')&& detail.isComment" class="comment-box grayBtn" @tap="goComment(item, detail)">查看评价</button>
			  </view>
            </view>
          </view>
        </view>
        <view class="btn-box">
          <button v-if="item.status==='0'" class="grayBtn" @tap="cancelOrder(item)">取消订单</button>
          <button v-if="item.status==='2'" class="grayBtn" @tap="trackLogistics(item)">查看物流</button>
          <button v-if="item.status==='2'" class="redBtn" @tap="confirm(item)">确认收货</button>
          <button v-if="item.status==='0'" class="redBtn" @tap="payOff(item.ID)">支付订单</button>
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
    orderList.value = res.data.list
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
    orderList.value = res.data.list
  }
}

</script>

<style lang="scss">
page {
  background-color: #f8f8f8;
}

.order_nav_view {
  width: 100%;
  height: 80rpx;
}

.order_nav_box {
  width: 100%;
  height: 80rpx;
  /* #ifdef H5 */
  top: var(--window-top);
  /* #endif */
  /* #ifndef H5 */
  top: 0;
  /* #endif */
  z-index: 1;
  left: 0;
  right: 0;
}

.order_goods_card {
  padding: 24rpx 32rpx 0;
}

.order_goods_card_item {
  padding: 24rpx;
  border-radius: 12rpx;
}

.order_goods_card_img {
  width: 151rpx;
  height: 152rpx;
  border-radius: 12rpx;
  overflow: hidden;
}

.comment-box{
  width: 200rpx;
  height: 56rpx;
  line-height: 56rpx;
  font-size: 28rpx;
  margin-left: 10rpx;
}
.btn-box {
  width: 100%;
  display: flex;
  justify-content: flex-end;
  margin-top: 24rpx;
  .grayBtn, .redBtn {
    width: 200rpx;
    border-color: #ddd;
    background-color: #fff;
    height: 56rpx;
    line-height: 56rpx;
    color: #333;
    font-size: 28rpx;
	margin-right: 0;
	margin-left: 10rpx;
  }

  .redBtn {
    background-color: #fe5572 !important;
    border-color: #fe5572;
    color: #fff;
  }
}

.texts {
  max-width: 340rpx;
  font-size: 24rpx;
  padding-right: 18rpx;
}

</style>
