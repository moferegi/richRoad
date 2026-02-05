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
        <image class="item_goods_img m_r_16" :src="getUrl(d.sku.picture)" mode="cover"	></image>
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
      <!-- 积分抵扣选项 -->
      <view class="price-row points-row">
        <view class="points-left">
          <checkbox-group @change="onPointsChange">
            <view class="flex-aic">
              <checkbox value="points" :checked="usePoints" color="#fa436a" style="transform: scale(0.8); margin-right: 8rpx;"/> <text>积分抵扣</text>
            </view>
          </checkbox-group>
         
        </view>
        <text class="points-amount">可抵扣 ¥{{ Math.min(availablePointsAmount, (data.originPrice - data.discount)) / 100 }}</text>
      </view>
      <!-- <view class="price-row">
        <text>运费</text>
        <text>免运费</text>
      </view> -->
    </view>

    <!-- 退款操作 -->
    <view class="refund-section" v-if="canApplyRefund || isRefunding || isRefunded">
      <view v-if="canApplyRefund" class="refund-btn" @tap="openRefund">申请退款</view>
      <view v-else class="refund-status">{{ isRefunding ? '退款处理中' : '已退款' }}</view>
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

    <refund-apply-popup
      v-model:visible="refundVisible"
      :order-id="orderID"
      @success="initSingleOrder"
    />
  </view>
</template>

<script setup>
import { ref, nextTick, watch, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { selfOrder,changeOrderCoupon } from '@/api/order.js'
import { claimCouponByUser } from '@/api/coupon.js'
import { getPayParams, getOrderById,checkNeedPay } from '@/api/base.js'
import { getUserInfo } from '@/api/base.js'
import { getUrl } from "@/utils/url.js"
import { changeOrderPoints } from '@/api/order.js'
import RefundApplyPopup from '@/components/refund-apply-popup/refund-apply-popup.vue'

const toAddress = () => {
  uni.navigateTo({
    url: `/pages/address/address?ID=${data.value.ID}`,
  })
}



const totalPrice = ref(475)
const hasAddress = ref(false)
const usePoints = ref(false)
const availablePointsAmount = ref(0)
const userPoints = ref(0)
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
const refundVisible = ref(false)

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
            /*uni.showToast({
              title: '领取成功',
              icon: 'success',
              duration: 1500
            })*/
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
          /*uni.showToast({
            title: '操作失败，请重试',
            icon: 'none'
          })*/
        } finally {
          uni.hideLoading()
        }
			}

// 注意：价格计算现在由后端处理，前端只需要显示从后端获取的价格

// 获取用户积分信息
const getUserPoints = async () => {
  try {
    const res = await getUserInfo()
    if (res.code === 0) {
      userPoints.value = res.data.point || 0
      // 计算可抵扣金额（积分按1:1抵扣分，即100积分=1元）
      // 积分不能超过商品价格，需要在计算总价时动态限制
      availablePointsAmount.value = userPoints.value

    }
  } catch (error) {
    console.error('获取用户积分失败:', error)
  }
}

onLoad((options) => {
  orderID.value = options.orderID
  initSingleOrder()
  getUserPoints()
})

const canApplyRefund = computed(() => {
  const status = String(data.value.status || '')
  return ['1', '2', '3', '7'].includes(status)
})

const isRefunding = computed(() => String(data.value.status || '') === '6')
const isRefunded = computed(() => String(data.value.status || '') === '5')

const openRefund = () => {
  refundVisible.value = true
}

// 处理积分抵扣checkbox变化
const onPointsChange = async (e) => {
  console.log('checkbox change event:', e)
  const isChecked = e.detail.value.includes('points')
  usePoints.value = isChecked
  console.log('usePoints updated to:', usePoints.value)
  
  // 调用后端接口变更积分抵扣
  try {
    const res = await changeOrderPoints({
      orderID: orderID.value,
      usePoints: usePoints.value
    })
    
    if (res.code === 0) {
      // 重新获取订单信息以更新价格
      await initSingleOrder()
      console.log('积分抵扣状态更新成功')
    } else {
      uni.showToast({
        title: '积分抵扣更新失败',
        icon: 'none'
      })
      // 恢复checkbox状态
      usePoints.value = !isChecked
    }
  } catch (error) {
    console.error('积分抵扣更新失败:', error)
    uni.showToast({
      title: '积分抵扣更新失败',
      icon: 'none'
    })
    // 恢复checkbox状态
    usePoints.value = !isChecked
  }
}



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

let isChecking = false
const stopChecking = () => {
  isChecking = false
}

const checkOrder = () => {
  isChecking = true
  checkOrderRecursive()
}

const checkOrderRecursive = async () => {
  if (!isChecking) return
  
  try {
    const res = await getOrderById(orderID.value)
    if (res.data.TradeState === "SUCCESS") {
      // uni.showToast({
      //   title: "支付成功",
      //   icon: "none"
      // });
      stopChecking()
      // 这里要重新获取当前订单信息改变状态
      await selfOrder(orderID.value)
      // 更新成功后跳转到订单页
      uni.navigateTo({
        url: `/pages/order/order?orderID=${orderID.value}`,
      })
      return
    }
  } catch (error) {
    console.error('检查订单状态失败:', error)
  }
  
  // 如果还在检查中，1秒后继续检查
  if (isChecking) {
    setTimeout(checkOrderRecursive, 1000)
  }
}

const tapPay = async () => {
  // 遍历data，如果地址和用户信息任一一项缺失则提示用户手动填写并跳转到地址页面
  if (!data.value.name || !data.value.phone || !data.value.province || !data.value.city || !data.value.area || !data.value.detail) {
    uni.navigateTo({
      url: `/pages/address/address?ID=${data.value.goodID || data.value.ID}`,
    })
    /*uni.showToast({
      title: '地址信息不全！请填写收货地址',
      icon: 'none'
    })*/
    return
  }

  // 获取参数
  const params = {
    "orderID": Number(orderID.value),
    "openid": uni.getStorageSync('openid')
  }

  const needPayRes = await checkNeedPay(params)
  if (!needPayRes.data) {
    // 0元购成功，显示成功提示
    uni.showToast({
      title: '订单提交成功！',
      icon: 'success',
      duration: 2000
    })
    // 延迟跳转到订单页面
    setTimeout(() => {
      uni.navigateTo({
        url: `/pages/order/order?orderID=${orderID.value}`,
      })
    }, 2000)
    return
  }

  console.log(params)

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

.refund-section {
  margin-top: 20rpx;
  background-color: #fff;
  padding: 20rpx 32rpx;
  display: flex;
  justify-content: flex-end;
  border-radius: 12rpx;
}

.refund-btn {
  padding: 12rpx 28rpx;
  border-radius: 30rpx;
  border: 1rpx solid #fa8c16;
  color: #fa8c16;
  background-color: #fff7e6;
  font-size: 26rpx;
}

.refund-status {
  padding: 12rpx 28rpx;
  border-radius: 30rpx;
  border: 1rpx solid #ccc;
  color: #999;
  background-color: #fff;
  font-size: 26rpx;
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
