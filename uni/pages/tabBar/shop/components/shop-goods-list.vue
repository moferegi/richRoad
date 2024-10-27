<template>
	<view class="shop_tool_view">
		<view class="shop_tool_box flex-aic flexr-jfe boxs_bb pos_f bgc_fff">
			<text v-if="isEdit" @tap="manage">管理</text>
			<text v-if="!isEdit" @tap="finishEdit">退出管理</text>
		</view>
	</view>
	<view class="shop_list_view">
		<view class="bgc_fff shop_list_box">
      <up-empty
          v-if="cartList.length<1"
          mode="car"
          icon="http://cdn.uviewui.com/uview/empty/car.png"
      >
        <template #default>
          <text class="p_t_16 p_b_40 font_24 color_fe5572" @tap="goTo" >挑选商品</text>
        </template>
      </up-empty>
			<template>
					<view class="flex flex-aic p_t_24" v-for="(item, index) in cartList" :key="index">
						<image class="shop_item_img m_r_16" :src="getUrl(item.sku.picture)" mode="aspectFill"></image>
						<view class="flex-fitem">
							<view class="color_333 font_26 m_b_8 title">{{item.sku.name}}</view>
							<view class="flex flex-aic m_b_4">
								<view class="color_999 m_r_24 title size-sm">{{item.sku.description}}</view>
							</view>
							<view class="flex-aic flexr-jsb">
								<view class="color_ff0003">
									<text class="font_28">¥</text>
									<text class="font_40">{{item.sku.price/100}}</text>
								</view>
								<text v-if="!isEdit" class="color_fe5572 font_24" @tap="deleteItem(item)">删除</text>
							</view>
						</view>
					</view>
			</template>
		</view>
	</view>
	<view class="shop_nav_box pos_f bgc_fff flex-aic flexr-jsb boxs_bb">
		<view class="flex-fitem flex-aic flexr-jfe m_r_16"></view>
    <view class="btns_box color_fff font_32 flex flex-aic">
      <view class="bgc_ff7000 tac" v-if="!isDeleteAll" @tap="toSettlement">结算</view>
      <view class="bgc_ff7000 tac" v-if="isDeleteAll" @tap="clearAllCart">删除</view>
    </view>
	</view>
</template>

<script setup>
	import { ref } from 'vue'
	import { getSelfCart, cutCart, clearCart } from "@/api/cart.js"
	import { onShow } from '@dcloudio/uni-app'
	import { useUserStore } from "@/pinia/modules/user";
  import { placeOrderByCart } from '@/api/order.js'
import {getUrl} from "@/utils/url.js"

	const cartList = ref([])
  const isEdit = ref(true)

	const initPage = async () => {
    isEdit.value = true
		const userStore = useUserStore()
		const token = userStore.token || ''

		if (token) {
			const res = await getSelfCart()
			if (res.code === 0) {
				cartList.value = res.data
			}
		} else {
			uni.showToast({
				title: '请登录后访问',
				mask: true,
				icon: 'none'
			});
			uni.redirectTo({
				url: '/pages/user/login'
			})
		}
	}
	onShow(() => {
		initPage()
	})

	const isDeleteAll = ref(false)

  // 跳转到首页挑选商品
  const goTo = () => {
    uni.switchTab({
      url: '/pages/tabBar/index'
    })
  }

	// 删除购物车某项
	const deleteItem = async (params) => {
		const temp = {
			goodID: params.goodID,
			skuID: params.skuID
		}
		const res = await cutCart(temp)
		if (res.code === 0) {
			uni.showToast({
				title: '删除成功',
				mask: true,
				icon: 'none'
			});
			initPage()
		} else {
			uni.showToast({
				title: '删除失败，请稍后重试',
				mask: true,
				icon: 'none'
			});
		}
	}

	// 删除购物车全部内容
	const clearAllCart = async () => {
		if (!selectedItemIds.value.length) {
			uni.showToast({
				title: '您还没有选择宝贝哦',
				mask: true,
				icon: 'none'
			});
			return
		}
		const res = await clearCart()
		if (res.code === 0) {
			uni.showToast({
				title: '删除成功',
				mask: true,
				icon: 'none'
			});
		} else {
			uni.showToast({
				title: '删除失败，请稍后重试',
				mask: true,
				icon: 'none'
			});
		}
	}

	//购物车编辑
	const manage = () => {
		isEdit.value = false
		// 购物车为正在编辑状态时,如已全选则展示删除按钮
		isDeleteAll.value = true
	}
	// 购物车编辑提交
	const finishEdit = () => {
		isEdit.value = true
		isDeleteAll.value = false
	}

	// 结算
	const toSettlement = async () => {
    // 先调用placeOrderByCart生成订单成功后再跳转
    if(cartList.value.length) {
      const data = cartList.value.map(item => {
        return {
          goodID: item.goodID,
          skuID: item.skuID,
          quantity: item.quantity
        }
      })
      let detail = {
        "detail": data
      };
      const res = await placeOrderByCart(detail);
      if (res.code === 0) {
        uni.navigateTo({
          url: `/pages/orderInfo/orderInfo?orderID=${res.data.orderID}&type=cart`
        })
      } else {
        uni.showToast({
          title: '生成订单失败，请稍后重试',
          mask: true,
          icon: 'none'
        });
      }
    } else {
      uni.showToast({
        title: '请先添加商品到购物车',
        mask: true,
        icon: 'none'
      });
    }
	}
</script>

<style lang="scss" scoped>
	.shop_list_view {
		padding: 24rpx 32rpx;
	}

	.shop_list_box {
		padding: 20rpx;
		border-radius: 12rpx;
	}

	.shop_item_img {
		width: 152rpx;
		height: 152rpx;
		border-radius: 6rpx;
	}

	.shop_tool_view {
		width: 100%;
		height: 88rpx;
	}

	.shop_tool_box {
		height: 88rpx;
		left: 0;
		right: 0;
		z-index: 1;
		font-size: 24rpx;
		padding-right: 32rpx;
	}

	.shop_nav_box {
		width: 100%;
		height: 100rpx;
		box-shadow: 0rpx -2rpx 6rpx 0rpx rgba(215, 215, 215, 0.5);
		left: 0;
		bottom: 0;
		z-index: 1;
		padding: 0 32rpx 0 16rpx;
	}

	.btn_shop_a {
		padding: 0 8rpx;
		height: 24rpx;
		top: 10rpx;
		right: -16rpx;
		line-height: 26rpx;
		border-radius: 24rpx;
	}

  .btns_box {
    border-radius: 6rpx;
    overflow: hidden;

    view {
      width: 200rpx;
      height: 80rpx;
      line-height: 80rpx;
    }
  }
  .title{
    width: 100%;
    font-size: 30rpx;
    max-hight: 30rpx;
    overflow: hidden;				//溢出内容隐藏
    text-overflow: ellipsis;		//文本溢出部分用省略号表示
    display: -webkit-box;			//特别显示模式
    -webkit-line-clamp: 1;			//行数
    line-clamp: 1;
    -webkit-box-orient: vertical;	//盒子中内容竖直排列
  }
  .size-sm{
    font-size: 24rpx !important;
  }

</style>
