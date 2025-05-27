<template>
  <view class="container">
    <!-- 收货地址信息 -->
    <view @tap="toAddress" class="flex-aic flexr-jsb bgc_fff address_box">
      <image class="location_icon m_r_16" src="./../../static/images/dizhis-icons.png" mode=""></image>
      <view class="flex-fitem" v-if="hasAddress">
        <text class="color_333 font_28 m_b_4">{{data.name}} {{data.phone}}</text>
        <view class="text_nowrap color_999 font_24" style="max-width: 500rpx;">
          {{data.province}}{{data.city}}{{data.area}} {{data.street}}
        </view>
      </view>
      <view class="flex-fitem address-placeholder" v-else>
        <text class="color_333 font_28 m_b_4">请选择收货地址</text>
        <view class="text_nowrap color_999 font_24" style="max-width: 500rpx;">
          添加收货人信息及地址
        </view>
      </view>
      <text class="forward-icon">›</text>
    </view>

    <!-- 分割线 -->
    <view class="dashed-line"></view>

    <!-- 商品列表 -->
    <view class="order_goods_box bgc_fff">

      <!-- 商品1 -->
      <view class="flex p_b_16" v-for="d in data.detail">
        <image class="item_goods_img m_r_16" :src="d.sku.picture" mode="cover"	></image>
        <view class="flex-fitem">
          <view class="goods-name">{{d?.sku?.name}}</view>
          <view class="goods-desc">{{d?.good?.description}}</view>
          <view class="goods-attrs">{{d?.sku?.specs?.map(i=>i.value).join(" ")}}</view>
          <view class="goods-price-row">
            <text class="goods-price">¥{{d.sku.price / 100}}</text>
            <text class="goods-count">×{{d.quantity}}</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 优惠券和备注 -->
    <view class="order-options">
      <view class="option-item" @click="opencoupon">
        <view class="option-label red-icon">优惠券</view>
        <view class="option-value">
          <text class="discount-text" v-if="data.discount > 0">¥{{data.discount/100}}</text>
          <text class="discount-text" v-else>选择优惠券</text>
          <text class=""><wu-icon name="arrow-right"></wu-icon></text>
        </view>
      </view>
    </view>

    <!-- 订单金额信息 -->
    <view class="price-summary">
      <view class="price-row">
        <text>商品金额</text>
        <text>¥{{ data.originPrice / 100 }}</text>
      </view>
      <view class="price-row discount">
        <text>优惠金额</text>
        <text>-¥{{ data.discount / 100 }}</text>
      </view>
      <!-- <view class="price-row">
        <text>运费</text>
        <text>免运费</text>
      </view> -->
    </view>

    <!-- 底部支付栏 -->
    <view class="footer">
      <view class="total-container">
        <text class="total-label">实付款</text>
        <text class="total-price">¥{{totalPrice / 100}}</text>
      </view>
      <view class="pay-btn" @tap="tapPay">
        <text>提交订单</text>
      </view>
    </view>


    <!-- 选择优惠券弹出层 -->
		<view class="mask" catchtouchmove="preventTouchMove" v-if="couponshow == true" @tap="hidecoupon"></view>
		<view class="coupon" :style="'bottom:' + (couponshow == true ? '0px':'')">
			<scroll-view class="scrolls" scroll-y>
				<!-- colors:按钮颜色 couponList:优惠卷列表数据  @onReceive：领取或立即使用按钮事件 -->
				<cc-defineCoupon v-if="couponshow" colors="#fa436a" @onReceive="onReceive"></cc-defineCoupon>
			</scroll-view>
		</view>

  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { selfOrder,changeOrderCoupon } from '@/api/order.js'
import { claimCouponByUser } from '@/api/coupon.js'
import { getPayParams, getOrderById,checkNeedPay } from '@/api/base.js'
import { getUrl } from "@/utils/url.js"

const toAddress = () => {
  uni.navigateTo({
    url: `/pages/address/address?ID=${data.value.ID}`,
  })
}



const totalPrice = ref(475)
const hasAddress = ref(false)
const data = ref({
  name: '许小贤',
  phone: '13685395563',
  province: '山东省',
  city: '济南市',
  area: '历城区',
  street: '149号',
  detail: [
    {
      sku: {
        name: '西域小洛镇',
        description: '古瓷妃 短袖t恤女夏装2019新款',
        price: 1780,
        picture: 'http://www.liwanying.top/applate-icon/product-img1.png',
        attrs: [{ label: '颜色', value: '粉紫色' }, { label: '尺码', value: 'L' }]
      }
    },
    {
      sku: {
        name: '西域小洛镇',
        description: '韩版子是网络鞋 夏季清凉防滑简约百搭',
        price: 1780,
        picture: 'http://www.liwanying.top/applate-icon/product-img2.png',
        attrs: [{ label: '颜色', value: '粉紫色' }, { label: '尺码', value: 'L' }]
      }
    }
  ]
})
const orderID = ref("")

const couponshow = ref(false)

const opencoupon = () => {
				couponshow.value = true
			}
			// 关闭优惠券弹框
    const hidecoupon = () => {
				couponshow.value = false
			}
			//领取优惠券 立即使用事件
			//领取优惠券 立即使用事件
			const onReceive = async (item, index) => {
        // 添加loading遮罩防止多次点击
        uni.showLoading({
          title: item.couponNum == 0 ? '领取中...' : '选择中...',
          mask: true
        })

        try {
          if (item.couponNum == 0) {
            const res = await claimCouponByUser({
              couponID: item.couponID,
            })
            item.couponNum = res.data
            // 领取成功提示
            uni.showToast({
              title: '领取成功',
              icon: 'success',
              duration: 1500
            })
          }
          await changeOrderCoupon({
            orderID: orderID.value,
            couponNum: item.couponNum
          })

          // 关闭优惠券弹窗
          setTimeout(() => {
            initSingleOrder()
            hidecoupon()
          }, 500)
        } catch (error) {
          uni.showToast({
            title: '操作失败，请重试',
            icon: 'none'
          })
        } finally {
          uni.hideLoading()
        }
			}

onLoad((options) => {
  orderID.value = options.orderID
  initSingleOrder()
})

const initSingleOrder = async () => {
  const order = await selfOrder(orderID.value)
  if (order.code === 0) {
    data.value = order.data
    totalPrice.value = order.data.totalPrice
    if (order.data.city) {
      hasAddress.value = true
    }
  }
}

let timer = null
const clreatTimer = () => {
  clearInterval(timer)
  timer = null
}

const checkOrder = () => {
  timer = setInterval(async () => {
    const res = await getOrderById(orderID.value)
    if (res.data.TradeState === "SUCCESS") {
      uni.showToast({
        title: "支付成功",
        icon: "none"
      });
      clreatTimer()
      // 这里要重新获取当前订单信息改变状态
      await selfOrder(orderID.value)
      // 更新成功后跳转到订单页
      uni.navigateTo({
        url: `/pages/order/order?orderID=${orderID.value}`,
      })
    }
  }, 1000)
}

const tapPay = async () => {
  // 遍历data，如果地址和用户信息任一一项缺失则提示用户手动填写并跳转到地址页面
  if (!data.value.name || !data.value.phone || !data.value.province || !data.value.city || !data.value.area || !data.value.detail) {
    uni.navigateTo({
      url: `/pages/address/address?ID=${data.value.goodID || data.value.ID}`,
    })
    uni.showToast({
      title: '地址信息不全！请填写收货地址',
      icon: 'none'
    })
    return
  }
  // 获取参数
  const params = {
    "orderID": Number(orderID.value),
    "openid": uni.getStorageSync('openid')
  }

  const needPayRes = await checkNeedPay(params)
  if (!needPayRes.data) {
    uni.showToast({
      title: '订单已支付！',
      icon: 'none'
    })
    uni.navigateTo({
        url: `/pages/order/order?orderID=${orderID.value}`,
      })
      return
  }

  const res = await getPayParams(params)
  if (res.code === 0){
    uni.requestPayment({
      provider: 'wxpay',
      timeStamp: res.data.timeStamp,
      nonceStr: res.data.nonceStr,
      package: res.data.package,
      signType: res.data.signType,
      paySign: res.data.paySign,
      success: function (res) {
        checkOrder()
      },
      fail: function (err) {
        console.log('fail:' + JSON.stringify(err));
      }
    });
  }
}
</script>

<style lang="scss">
page {
  background-color: #f8f8f8;
  color: #333;
  font-size: 28rpx;
}

.container {
  padding-bottom: 120rpx;
}

.header {
  display: flex;
  align-items: center;
  padding: 20rpx 32rpx;
  background-color: #fff;
}

.back-icon {
  font-size: 40rpx;
  margin-right: 20rpx;
}

.header-title {
  font-size: 32rpx;
  font-weight: bold;
}

.address_box {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 32rpx;
  background-color: #fff;
  margin-top: 20rpx;
}

.location_icon {
  width: 40rpx;
  height: 40rpx;
  margin-right: 16rpx;
}

.forward-icon {
  font-size: 36rpx;
  color: #ccc;
}

.dashed-line {
  height: 2rpx;
  background-image: linear-gradient(to right, #ddd 50%, transparent 50%);
  background-size: 10rpx 1rpx;
  background-repeat: repeat-x;
  margin: 0 32rpx;
}

.flex-aic {
  display: flex;
  align-items: center;
}

.flexr-jsb {
  justify-content: space-between;
}

.flex-fitem {
  flex: 1;
}

.color_333 {
  color: #333;
}

.color_999 {
  color: #999;
}

.font_28 {
  font-size: 28rpx;
}

.font_24 {
  font-size: 24rpx;
}

.m_r_16 {
  margin-right: 16rpx;
}

.m_b_4 {
  margin-bottom: 4rpx;
}

.text_nowrap {
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}

.order_goods_box {
  padding: 24rpx 32rpx;
  background-color: #fff;
}

.flex {
  display: flex;
}

.p_b_16 {
  padding-bottom: 16rpx;
}

.item_goods_img {
  min-width: 168rpx;
  max-width: 168rpx;
  height: 168rpx;
  border-radius: 8rpx;
  margin-right: 16rpx;
}

.goods-name {
  font-size: 24rpx;
  color: #999;
  margin-bottom: 4rpx;
}

.goods-desc {
  font-size: 28rpx;
  color: #333;
  margin-bottom: 8rpx;
}

.goods-attrs {
  font-size: 24rpx;
  color: #999;
  margin-bottom: 8rpx;
}

.goods-price-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.goods-price {
  font-size: 28rpx;
  color: #333;
}

.goods-count {
  font-size: 24rpx;
  color: #999;
}

.order-options {
  margin-top: 20rpx;
  background-color: #fff;
}

.option-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx 32rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.option-label {
  display: flex;
  align-items: center;
}

.option-label::before {
  content: "";
  display: inline-block;
  width: 30rpx;
  height: 30rpx;
  margin-right: 10rpx;
  border-radius: 6rpx;
}

.red-icon::before {
  background-color: #ff6b6b;
}

.yellow-icon::before {
  background-color: #ffb73b;
}

.option-value {
  display: flex;
  align-items: center;
}

.discount-text {
  color: #ff5572;
  margin-right: 10rpx;
}

.message-placeholder {
  color: #999;
  margin-right: 10rpx;
}

.price-summary {
  margin-top: 20rpx;
  background-color: #fff;
  padding: 0 32rpx;
}

.price-row {
  display: flex;
  justify-content: space-between;
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.discount {
  color: #ff5572;
}

.footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 100rpx;
  display: flex;
  background-color: #fff;
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
}

.total-container {
  flex: 1;
  display: flex;
  align-items: center;
  padding-left: 32rpx;
}

.total-label {
  font-size: 28rpx;
  margin-right: 8rpx;
}

.total-price {
  font-size: 36rpx;
  font-weight: bold;
  color: #ff5572;
}

.pay-btn {
  width: 240rpx;
  height: 100%;
  background-color: #ff5572;
  color: #fff;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 32rpx;
}


.content {
		display: flex;
		flex-direction: column;

	}

	.mask {
		width: 100%;
		height: 100vh;
		position: fixed;
		top: 0;
		left: 0;
		background: #000;
		z-index: 900;
		opacity: 0.7;
	}

	/* 优惠券 */
	.coupon {
		background-color: #fff;
		border-radius: 10upx 10upx 0 0;
		position: fixed;
		left: 0;
		bottom: -1000upx;
		z-index: 999;
		transition: all 0.3s;
	}

	.scrolls {
		width: 100vw;
		height: 60vh;
		padding-top: 10upx;
		z-index: 500;
	}
</style>
