<template>
	<view>
		<view @tap="toAddress" class="flex-aic flexr-jsb bgc_fff address_box m_t_24 m_b_24">
			<image class="location_icon m_r_16"  src="./../../static/images/dizhis-icons.png" mode="">

			</image>
			<view class="flex-fitem">
				<text class="color_333 font_28 m_b_4">{{data.name || "**"}} {{data.phone || "**"}}</text>
				<view class="text_nowrap color_999 font_24" style="max-width: 500rpx;">
				{{data.province || "**"}} - {{data.city || "**"}} - {{data.area || "**"}} - {{data.street || "**"}}
				</view>
			</view>
		</view>
		<view class="order_goods_box bgc_fff">
			<view class="flex p_b_16" v-for="(item, index) in data.detail" :key="index">
				<image class="item_goods_img m_r_16" :src="getUrl(item.sku.picture)" mode=""></image>
				<view class="flex-fitem">
					<view class="m_b_8 flex-aic flexr-jsb">
						<text class="color_333 font_28">￥{{item.sku.price/100}}</text>
					</view>
					<view class="color_999 font_24 m_b_4 text_nowrap" style="max-width: 560rpx;"
						v-for="(e, i) in item.sku.attrs">{{e.label}}：{{e.value}}</view>
				</view>
			</view>
		</view>
		<view class="order_info_box bgc_fff color_333">
			<view class="b_t_2 flex-aic flexr-jsb order_info_item_p24">
				<text class="font_28">运费</text>
				<text class="font_28 font_bold">快递 免邮</text>
			</view>
			<view class="b_t_2 order_info_item_p24 tar ">
				<text class="font_28">金额：¥</text>
				<text class="font_40">{{totalPrice/100}}</text>
			</view>
		</view>
		<view class="bgc_fff m_t_24 wx_pay_box m_b_24">
			<view class="color_333 font_28" style="margin-bottom: 32rpx;">支付方式</view>
			<view class="flex-aic flexr-jsb">
				<image class="wx_icon m_r_16" src="http://www.liwanying.top/applate-icon/weixinzhifu.png" mode="">
				</image>
				<text class="color_333 font_28 flex-fitem">微信支付</text>
				<image class="wx_select_icon" src="http://www.liwanying.top/applate-icon/xuanzhong.png" mode=""></image>
			</view>
		</view>
    <view class="tac">
      <button class="confirm" @tap="tapPay">确认支付</button>
    </view>

	</view>
</template>

<script setup>
	import {
		ref
	} from 'vue'
	import { onLoad } from '@dcloudio/uni-app'
	import { selfOrder } from '@/api/order.js'
	import { getPayParams, getOrderById } from '@/api/base.js'
import {getUrl} from "@/utils/url.js"
	const toAddress = () => {
    uni.navigateTo({
      url: `/pages/address/address?ID=${data.value.ID}`,
      })
	}

	const totalPrice = ref(0)
	const data = ref({})
	const orderID = ref("")
	onLoad((options) => {
    orderID.value = options.orderID
    initSingleOrder(orderID.value)
	})

  const initSingleOrder = async (orderID) => {
    const order = await selfOrder(orderID)
    if(order.code === 0) {
      data.value = order.data
      totalPrice.value = order.data.totalPrice
    }
  }


	let timer = null
	const clreatTimer = () =>{
		clearInterval(timer)
		timer = null
	}
	const checkOrder = () => {
		timer = setInterval(async ()=>{
		 const res = await getOrderById(orderID.value)
			if(res.data.TradeState === "SUCCESS"){
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
		},1000)
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
				"openID": uni.getStorageSync('openid')
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
	}

	.address_box {
		padding: 28rpx 32rpx;
	}

	.location_icon {
		width: 40rpx;
		height: 40rpx;
	}

	.wx_icon {
		width: 46rpx;
		height: 40rpx;
	}

	.wx_select_icon {
		width: 32rpx;
		height: 32rpx;
	}

	.address_icon {
		width: 12rpx;
		height: 24rpx;
	}

	.order_info_box {
		padding: 0 32rpx;
	}

	.order_info_item_p24 {
		padding: 24rpx 0;
	}

	.wx_pay_box {
		padding: 20rpx 32rpx;
	}

	.order_info_nav_view {
		height: 100rpx;
		width: 100%;
		padding-bottom: constant(safe-area-inset-bottom);
		padding-bottom: env(safe-area-inset-bottom);
	}

	.order_info_nav_box {
		bottom: 0;
		left: 0;
		right: 0;
		height: 100rpx;
		width: 100%;
		padding-bottom: constant(safe-area-inset-bottom);
		padding-bottom: env(safe-area-inset-bottom);
	}

	.item_goods_img {
		width: 108rpx;
		height: 108rpx;
		border-radius: 2rpx;
		overflow: hidden;
	}

	.order_goods_box {
		padding: 24rpx 32rpx;
	}
  .confirm{
    width: 60%;
    border-color: #fe5572;
    background-color: #fe5572;
    color: white;
  }
</style>
