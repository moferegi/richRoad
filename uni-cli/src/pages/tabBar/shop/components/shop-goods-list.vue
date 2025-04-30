<template>
	<view class="shop_tool_view">
		<view class="shop_tool_box flex-aic flexr-jfe boxs_bb pos_f bgc_fff">
			<text v-if="isEdit" @tap="manage" class="manage-btn">管理</text>
			<text v-if="!isEdit" @tap="finishEdit" class="manage-btn">退出管理</text>
		</view>
	</view>
	<view class="shop_list_view">
		<view class="bgc_fff shop_list_box">
      <!-- 自定义空购物车提示，替换原来的up-empty组件 -->
      <view v-if="cartList.length < 1" class="empty-cart">
        <image class="empty-cart-icon" src="http://cdn.uviewui.com/uview/empty/car.png" mode="aspectFit"></image>
        <text class="empty-cart-text">购物车空空如也</text>
        <button class="empty-cart-btn" @tap="goTo">去挑选商品</button>
      </view>

			<view class="cart-item" v-for="(item, index) in cartList" :key="index">
				<image class="shop_item_img" :src="getUrl(item.sku.picture)" mode="aspectFill"></image>
				<view class="cart-item-info">
					<view class="title">{{item.sku.name}}</view>
					<view class="desc">
						<view class="desc-text">{{item.sku.description}}</view>
					</view>
					<view class="price-action">
						<view class="price">
							<text class="price-symbol">¥</text>
							<text class="price-value">{{item.sku.price/100}}</text>
						</view>
						<text v-if="!isEdit" class="delete-btn" @tap="deleteItem(item)">删除</text>
					</view>
				</view>
			</view>
		</view>
	</view>
	<view class="shop_nav_box pos_f bgc_fff flex-aic flexr-jsb boxs_bb" v-if="cartList.length > 0">
		<view class="flex-fitem flex-aic flexr-jfe m_r_16"></view>
    <view class="btns_box">
      <button class="checkout-btn" v-if="!isDeleteAll" @tap="toSettlement">去结算</button>
      <button class="delete-all-btn" v-if="isDeleteAll" @tap="clearAllCart">删除所选</button>
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
    margin-bottom: 120rpx;
	}

	.shop_list_box {
		padding: 20rpx;
		border-radius: 16rpx;
    box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);
	}

  // 自定义的空购物车组件样式
  .empty-cart {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 60rpx 0;

    &-icon {
      width: 200rpx;
      height: 200rpx;
      margin-bottom: 30rpx;
    }

    &-text {
      font-size: 28rpx;
      color: #999;
      margin-bottom: 40rpx;
    }

    &-btn {
      background-color: #FE5572;
      color: #fff;
      font-size: 28rpx;
      padding: 16rpx 40rpx;
      border-radius: 40rpx;
      border: none;
    }
  }

  // 购物车项目样式优化
  .cart-item {
    display: flex;
    padding: 24rpx 0;
    border-bottom: 1px solid #f5f5f5;

    &:last-child {
      border-bottom: none;
    }
  }

	.shop_item_img {
		width: 160rpx;
		height: 160rpx;
		border-radius: 12rpx;
    margin-right: 24rpx;
    box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
	}

  .cart-item-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
  }

  .title {
    font-size: 28rpx;
    color: #333;
    font-weight: 500;
    margin-bottom: 12rpx;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 1;
    line-clamp: 1;
    -webkit-box-orient: vertical;
  }

  .desc {
    margin-bottom: 12rpx;

    &-text {
      font-size: 24rpx;
      color: #999;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 1;
      line-clamp: 1;
      -webkit-box-orient: vertical;
    }
  }

  .price-action {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .price {
    color: #FF4141;

    &-symbol {
      font-size: 24rpx;
    }

    &-value {
      font-size: 32rpx;
      font-weight: bold;
    }
  }

  .delete-btn {
    font-size: 24rpx;
    color: #FE5572;
    background-color: #fff;
    padding: 6rpx 20rpx;
    border-radius: 24rpx;
    border: 1px solid #FE5572;
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
		padding-right: 32rpx;
    border-bottom: 1px solid #f5f5f5;
    border-top: 1px solid #f5f5f5;
	}

  .manage-btn {
    font-size: 26rpx;
    color: #666;
    padding: 8rpx 24rpx;
  }

	.shop_nav_box {
		width: 100%;
		height: 100rpx;
		box-shadow: 0rpx -2rpx 8rpx 0rpx rgba(0, 0, 0, 0.08);
		left: 0;
		bottom: 0rpx;
		z-index: 10;
		padding: 0 32rpx;
	}

  .btns_box {
    display: flex;
    align-items: center;
  }

  .checkout-btn {
    width: 200rpx;
    height: 80rpx;
    line-height: 80rpx;
    background: linear-gradient(to right, #FF7000, #FF4141);
    color: #fff;
    font-size: 30rpx;
    border-radius: 40rpx;
    text-align: center;
    border: none;
  }

  .delete-all-btn {
    width: 200rpx;
    height: 80rpx;
    line-height: 80rpx;
    background-color: #FE5572;
    color: #fff;
    font-size: 30rpx;
    border-radius: 40rpx;
    text-align: center;
    border: none;
  }
</style>

