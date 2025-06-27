<template>
  <view class="shop_list_view">
    <view class="bgc_fff shop_list_box">
      <view v-if="cartList.length < 1" class="empty-cart">
        <text class="empty-cart-title">购物车竟然是空的</text>
        <text class="empty-cart-text">"再忙，也要记得买点什么犒赏自己~"</text>
        <view class="empty-cart-buttons">
          <button class="empty-cart-btn go-shopping" @tap="goTo">去逛逛</button>
        </view>
      </view>

      <view class="cart-item" v-for="(item, index) in cartList" :key="index">
        <image class="shop_item_img" :src="getUrl(item.sku.picture)" mode="aspectFill"></image>
        <view class="cart-item-info">
          <view>
            <view class="title">
              <text>{{item.sku.name}}</text>
            </view>
            <view class="desc">
              <view class="desc-text">{{item.sku.description}}</view>
            </view>
          </view>

          <view class="price-action">
            <view class="price">
              <text class="price-symbol">¥</text>
              <text class="price-value">{{item.sku.price /100}}</text>
            </view>
            <view class="quantity-control">
              <view v-if="!isDeleteAll" class="number-box-container">
                <wu-number-box :asyncChange="true" :min="0" @change="(e)=>onChange(item,e)" integer v-model="item.quantity"></wu-number-box>
              </view>
              <view v-else class="delete-btn" @tap="deleteItem(item)">删除</view>
            </view>
          </view>
        </view>
      </view>
    </view>
  </view>
  <view class="shop_nav_box pos_f bgc_fff flex-aic flexr-jsb boxs_bb" v-if="cartList.length > 0">
    <view class="total_money" v-if="!isDeleteAll">
      <text class="total_money_label">合计</text>￥{{cartList.reduce((total, item) => total + item.sku.price * item.quantity, 0) / 100}}
    </view>
    <view class="edit-mode" v-else @tap="toggleDeleteMode">完成</view>
    <view class="btns_box">
      <button class="edit-btn" v-if="!isDeleteAll" @tap="toggleDeleteMode">编辑</button>
      <button class="checkout-btn" v-if="!isDeleteAll" @tap="toSettlement">去结算</button>
      <button class="delete-all-btn" v-if="isDeleteAll" @tap="clearAllCart">清空购物车</button>
    </view>
  </view>
</template>

<script setup>
	import { ref } from 'vue'
	import { getSelfCart, cutCart, addCart, clearCart } from "@/api/cart.js"
	import { onShow } from '@dcloudio/uni-app'
	import { useUserStore } from "@/pinia/modules/user";
  import { placeOrderByCart } from '@/api/order.js'
  import {getUrl} from "@/utils/url.js"

	const cartList = ref([])

	const initPage = async () => {
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

	const onChange = async (item,e) => {
    console.log(item)
    uni.showLoading({
					title: '加载中',
					mask: true
				})
    let res
    if(item.quantity > e.value) {
       res = await cutCart({
        goodID: item.goodID,
        skuID: item.skuID,
        quantity:1,
      })
    }else{
      res = await addCart({
        goodID: item.goodID,
        skuID: item.skuID,
        quantity:1,
      })
    }
    uni.hideLoading();
    if (res.code !== 0) {
      uni.showToast({
        title: '调整失败',
        mask: true,
        icon: 'none'
      })
      return
    }
    item.quantity = e.value
    if(e.value == 0) {
      cartList.value = cartList.value.filter(i => i.ID !== item.ID)
    }
	}

	// 切换删除模式
	const toggleDeleteMode = () => {
		isDeleteAll.value = !isDeleteAll.value
	}

	// 删除单个商品
	const deleteItem = async (item) => {
		uni.showModal({
			title: '提示',
			content: '确定要删除这个商品吗？',
			success: async function (res) {
				if (res.confirm) {
					uni.showLoading({
						title: '删除中',
						mask: true
					})
					const result = await cutCart({
						goodID: item.goodID,
						skuID: item.skuID,
						quantity: item.quantity,
					})
					uni.hideLoading()
					if (result.code === 0) {
						cartList.value = cartList.value.filter(i => i.ID !== item.ID)
						uni.showToast({
							title: '删除成功',
							mask: true,
							icon: 'success'
						})
						// 如果购物车为空，自动退出删除模式
						if (cartList.value.length === 0) {
							isDeleteAll.value = false
						}
					} else {
						uni.showToast({
							title: '删除失败，请稍后重试',
							mask: true,
							icon: 'none'
						})
					}
				}
			}
		})
	}

	// 删除购物车全部内容
	const clearAllCart = async () => {
		uni.showModal({
			title: '提示',
			content: '确定要清空购物车吗？',
			success: async function (res) {
				if (res.confirm) {
					uni.showLoading({
						title: '清空中',
						mask: true
					})
					const result = await clearCart()
					uni.hideLoading()
					if (result.code === 0) {
						cartList.value = []
						isDeleteAll.value = false
						uni.showToast({
							title: '清空成功',
							mask: true,
							icon: 'success'
						})
					} else {
						uni.showToast({
							title: '清空失败，请稍后重试',
							mask: true,
							icon: 'none'
						})
					}
				}
			}
		})
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
  padding: 20rpx;
  margin-bottom: 20rpx;
}

.shop_list_box {
  border-radius: 24rpx;
  background: linear-gradient(135deg, #ffffff 0%, #fafafa 100%);
  box-shadow: 0 8rpx 32rpx rgba(255, 76, 125, 0.1);
  backdrop-filter: blur(10rpx);
  border: 1rpx solid rgba(255, 255, 255, 0.3);
  overflow: hidden;
}

// 自定义的空购物车组件样式
.empty-cart {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 100rpx 40rpx;
  height: 66vh;
  justify-content: center;
  background: linear-gradient(135deg, #fef7f0 0%, #fff5f5 100%);
  border-radius: 24rpx;
  margin: 20rpx;

  &-icon-container {
    width: 240rpx;
    height: 240rpx;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 48rpx;
    background: linear-gradient(135deg, #ff4c7d, #ff6b9d);
    box-shadow: 0 16rpx 48rpx rgba(255, 76, 125, 0.3);
    position: relative;

    &::before {
      content: '';
      position: absolute;
      width: 100%;
      height: 100%;
      border-radius: 50%;
      background: linear-gradient(135deg, rgba(255, 255, 255, 0.2), transparent);
      top: 0;
      left: 0;
    }
  }

  &-icon {
    width: 160rpx;
    height: 160rpx;
    filter: brightness(0) invert(1);
  }

  &-title {
    font-size: 36rpx;
    color: #333;
    font-weight: 600;
    margin-bottom: 20rpx;
    background: linear-gradient(135deg, #ff4c7d, #ff6b9d);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  &-text {
    font-size: 28rpx;
    color: #666;
    margin-bottom: 80rpx;
    text-align: center;
    line-height: 1.6;
  }

  &-buttons {
    display: flex;
    gap: 30rpx;
  }

  &-btn {
    min-width: 200rpx;
    height: 80rpx;
    line-height: 80rpx;
    font-size: 30rpx;
    border-radius: 40rpx;
    text-align: center;
    border: none;
    font-weight: 500;
    transition: all 0.3s ease;
    box-shadow: 0 8rpx 24rpx rgba(255, 76, 125, 0.2);

    &.go-shopping {
      background: linear-gradient(135deg, #ff4c7d, #ff6b9d);
      color: #fff;

      &:active {
        transform: translateY(2rpx);
        box-shadow: 0 4rpx 12rpx rgba(255, 76, 125, 0.3);
      }
    }
  }
}

// 购物车项目样式优化
.cart-item {
  display: flex;
  align-items: center;
  padding: 24rpx;
  margin: 16rpx 20rpx;
  background: linear-gradient(135deg, #ffffff 0%, #fafafa 100%);
  border-radius: 20rpx;
  box-shadow: 0 4rpx 16rpx rgba(255, 76, 125, 0.08);
  border: 1rpx solid rgba(255, 255, 255, 0.5);
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;

  &:hover {
    transform: translateY(-2rpx);
    box-shadow: 0 8rpx 24rpx rgba(255, 76, 125, 0.15);
  }

  &:last-child {
    margin-bottom: 24rpx;
  }
}


.select-circle {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  border: 2rpx solid #FE5572;
  background-color: #FE5572;
  position: relative;

  &::after {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 20rpx;
    height: 20rpx;
    border-radius: 50%;
    background-color: #fff;
  }
}

.shop_item_img {
  width: 200rpx;
  height: 200rpx;
  border-radius: 16rpx;
  margin-right: 24rpx;
  box-shadow: 0 8rpx 24rpx rgba(255, 76, 125, 0.15);
  border: 2rpx solid rgba(255, 255, 255, 0.8);
  transition: all 0.3s ease;

  &:hover {
    transform: scale(1.02);
    box-shadow: 0 12rpx 32rpx rgba(255, 76, 125, 0.2);
  }
}

.cart-item-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 200rpx;
  padding: 8rpx 0;
}

.title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 30rpx;
  color: #333;
  font-weight: 600;
  margin-bottom: 16rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.4;
}

.desc {
  margin-bottom: 24rpx;

  &-text {
    font-size: 26rpx;
    color: #666;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    line-height: 1.5;
    background: linear-gradient(135deg, #f8f9fa, #e9ecef);
    padding: 8rpx 12rpx;
    border-radius: 8rpx;
  }
}

.price-action {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 8rpx;
}

.delete-btn {
  padding: 12rpx 24rpx;
  background: linear-gradient(135deg, #ff4c7d, #ff6b9d);
  color: #fff;
  border-radius: 20rpx;
  font-size: 26rpx;
  font-weight: 500;
  box-shadow: 0 4rpx 12rpx rgba(255, 76, 125, 0.3);
  transition: all 0.3s ease;

  &:active {
    transform: translateY(1rpx);
    box-shadow: 0 2rpx 8rpx rgba(255, 76, 125, 0.4);
  }
}

.quantity-control {
  display: flex;
  align-items: center;
}

.number-box-container {
  display: flex;
  align-items: center;

  :deep(.wu-numberbox) {
    border: none;
    background: linear-gradient(135deg, #f8f9fa, #e9ecef);
    border-radius: 12rpx;
    box-shadow: 0 2rpx 8rpx rgba(255, 76, 125, 0.1);
    overflow: hidden;

    .wu-numberbox__minus,
    .wu-numberbox__plus {
      width: 64rpx;
      height: 64rpx;
      background: linear-gradient(135deg, #ff4c7d, #ff6b9d);
      color: #fff;
      border: none;
      font-weight: 600;
      transition: all 0.3s ease;

      &:active {
        background: linear-gradient(135deg, #e63946, #ff4757);
      }
    }

    .wu-numberbox__value {
      width: 80rpx;
      height: 64rpx;
      background: #fff;
      color: #333;
      margin: 0;
      font-size: 28rpx;
      font-weight: 600;
      border: 2rpx solid rgba(255, 76, 125, 0.1);
    }
  }
}

.price {
  background: linear-gradient(135deg, #ff4c7d, #ff6b9d);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  display: flex;
  align-items: baseline;
  gap: 4rpx;

  &-symbol {
    font-size: 26rpx;
    font-weight: 500;
  }

  &-value {
    font-size: 36rpx;
    font-weight: 700;
  }
}

.shop-title {
  font-size: 36rpx;
  font-weight: bold;
  text-align: center;
  padding: 20rpx 0;
  margin-bottom: 10rpx;
}

.total_money{
  font-size: 36rpx;
  background: linear-gradient(135deg, #ff4c7d, #ff6b9d);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-weight: 700;
  display: flex;
  align-items: center;
  .total_money_label{
    font-size: 26rpx;
    color: #666;
    margin-right: 8rpx;
    background: none;
    -webkit-text-fill-color: #666;
  }
  .discount-text {
    font-size: 22rpx;
    color: #999;
    font-weight: normal;
    margin-top: 4rpx;
    background: none;
    -webkit-text-fill-color: #999;
  }
}

.shop_nav_box {
  position: fixed;
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  height: 120rpx;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.95), rgba(248, 249, 250, 0.95));
  backdrop-filter: blur(20rpx);
  box-shadow: 0 -8rpx 32rpx rgba(255, 76, 125, 0.15);
  border-top: 1rpx solid rgba(255, 255, 255, 0.3);
  left: 0;
  z-index: 100;
  padding: 0 24rpx;
  /* #ifdef H5 */
  bottom: 88rpx;
  /* #endif */
  /* #ifdef MP */
  bottom: 0rpx;
  /* #endif */
}

.btns_box {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.edit-mode {
  font-size: 30rpx;
  color: #333;
  font-weight: 600;
  background: linear-gradient(135deg, #ff4c7d, #ff6b9d);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.checkout-btn {
  width: 220rpx;
  height: 88rpx;
  line-height: 88rpx;
  background: linear-gradient(135deg, #ff4c7d, #ff6b9d);
  color: #fff;
  font-size: 32rpx;
  font-weight: 600;
  border-radius: 44rpx;
  text-align: center;
  border: none;
  box-shadow: 0 8rpx 24rpx rgba(255, 76, 125, 0.3);
  transition: all 0.3s ease;

  &:active {
    transform: translateY(2rpx);
    box-shadow: 0 4rpx 16rpx rgba(255, 76, 125, 0.4);
  }
}

.edit-btn {
  width: 140rpx;
  height: 88rpx;
  line-height: 88rpx;
  background: linear-gradient(135deg, rgba(255, 76, 125, 0.1), rgba(255, 107, 157, 0.1));
  color: #ff4c7d;
  font-size: 28rpx;
  font-weight: 500;
  border-radius: 44rpx;
  text-align: center;
  border: 2rpx solid rgba(255, 76, 125, 0.2);
  transition: all 0.3s ease;

  &:active {
    background: linear-gradient(135deg, rgba(255, 76, 125, 0.2), rgba(255, 107, 157, 0.2));
    transform: translateY(1rpx);
  }
}

.delete-all-btn {
  width: 220rpx;
  height: 88rpx;
  line-height: 88rpx;
  background: linear-gradient(135deg, #ff4c7d, #ff6b9d);
  color: #fff;
  font-size: 32rpx;
  font-weight: 600;
  border-radius: 44rpx;
  text-align: center;
  border: none;
  box-shadow: 0 8rpx 24rpx rgba(255, 76, 125, 0.3);
  transition: all 0.3s ease;

  &:active {
    transform: translateY(2rpx);
    box-shadow: 0 4rpx 16rpx rgba(255, 76, 125, 0.4);
  }
}
</style>

