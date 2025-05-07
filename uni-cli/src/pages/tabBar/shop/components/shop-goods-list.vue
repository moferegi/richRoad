<template>
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
		padding: 12rpx 12rpx;
    margin-bottom: 12rpx;
	}

	.shop_list_box {
		padding: 12rpx;
		border-radius: 6rpx;
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
    border-bottom: 1px solid #f5f5f5;
    align-items: center;

    &:last-child {
      border-bottom: none;
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
		width: 220rpx;
		height: 220rpx;
		border-radius: 6rpx;
    margin-right: 24rpx;
    box-shadow: 0 2rpx 6rpx rgba(0, 0, 0, 0.05);
	}

  .cart-item-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    height: 220rpx;
  }

  .title {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 28rpx;
    color: #333;
    font-weight: 500;
    margin-bottom: 12rpx;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .desc {
    margin-bottom: 24rpx;

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
  
  .delete-btn {
    padding: 8rpx 20rpx;
    background-color: #FE5572;
    color: #fff;
    border-radius: 30rpx;
    font-size: 24rpx;
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
      background-color: #f5f5f5;
      border-radius: 4rpx;
      
      .wu-numberbox__minus,
      .wu-numberbox__plus {
        width: 60rpx;
        background-color: #f5f5f5;
        color: #333;
        border: none;
      }
      
      .wu-numberbox__value {
        width: 80rpx;
        background-color: #fff;
        color: #333;
        margin: 0 2rpx;
        font-size: 28rpx;
      }
    }
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

  .shop-title {
    font-size: 36rpx;
    font-weight: bold;
    text-align: center;
    padding: 20rpx 0;
    margin-bottom: 10rpx;
  }

  .total_money{
      font-size: 32rpx;
      color: #FF4141;
      font-weight: bold;
      display: flex;
      align-items: center;
      .total_money_label{
        font-size: 24rpx;
        color: #999;
        margin-right: 8rpx;
      }
      .discount-text {
        font-size: 22rpx;
        color: #999;
        font-weight: normal;
        margin-top: 4rpx;
      }
  }

	.shop_nav_box {
		position: fixed;
		display: flex;
		align-items: center;
		justify-content: space-between;
		width: 100%;
		height: 100rpx;
		box-shadow: 0rpx -2rpx 8rpx 0rpx rgba(0, 0, 0, 0.08);
		left: 0;
		z-index: 10;
		padding: 0 12rpx;
      /* #ifdef H5 */
      bottom: 90rpx;
      /* #endif */
      /* #ifdef MP */
      bottom: 0rpx; /* 或者你需要的小程序特定值 */
      /* #endif */
	}

  .btns_box {
    display: flex;
    align-items: center;
    gap: 20rpx;
  }
  
  .edit-mode {
    font-size: 28rpx;
    color: #333;
    font-weight: 500;
  }

  .checkout-btn {
    width: 200rpx;
    height: 80rpx;
    line-height: 80rpx;
    background: linear-gradient(to right, #FF5B8D, #FE5572);
    color: #fff;
    font-size: 30rpx;
    border-radius: 40rpx;
    text-align: center;
    border: none;
  }
  
  .edit-btn {
    width: 120rpx;
    height: 80rpx;
    line-height: 80rpx;
    background-color: rgba(254, 85, 114, 0.1);
    color: #FE5572;
    font-size: 28rpx;
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
  
  // 空购物车样式优化
  .empty-cart {
    &-icon {
      width: 240rpx;
      height: 240rpx;
    }
    
    &-btn {
      background: linear-gradient(to right, #FF5B8D, #FE5572);
    }
  }
</style>

