<template>
  <view class="nf-cart">
    <!-- 未登录提示 -->
    <view v-if="!isLoggedIn" class="nf-cart-empty">
      <view class="nf-cart-empty-icon">🛒</view>
      <text class="nf-cart-empty-title">{{ $t('cartLoginHint') }}</text>
      <view class="nf-cart-empty-btns">
        <view class="nf-btn nf-btn-primary" @tap="goLogin">{{ $t('goLogin') }}</view>
        <view class="nf-btn nf-btn-ghost" @tap="goRegister">{{ $t('goRegister') }}</view>
      </view>
    </view>

    <!-- 空购物车 -->
    <view v-else-if="cartList.length < 1" class="nf-cart-empty">
      <view class="nf-cart-empty-icon">🛒</view>
      <text class="nf-cart-empty-title">{{ $t('emptyCart') }}</text>
      <text class="nf-cart-empty-sub">{{ $t('emptyCartHint') }}</text>
      <view class="nf-cart-empty-btns">
        <view class="nf-btn nf-btn-primary" @tap="goTo">{{ $t('goShopping') }}</view>
      </view>
    </view>

    <!-- 购物车列表 -->
    <view class="nf-cart-list" v-if="isLoggedIn && cartList.length > 0">
      <view class="nf-cart-item" v-for="(item, index) in cartList" :key="item.ID">
        <!-- 选择 -->
        <view class="nf-cart-check" @tap="toggleSelect(index)" v-if="!isDeleteAll">
          <view :class="['nf-checkbox', { checked: item._selected }]">
            <text v-if="item._selected" class="nf-check-icon">✓</text>
          </view>
        </view>
        <!-- 商品图片 -->
        <image class="nf-cart-img" :src="getUrl(item.sku.picture)" mode="aspectFill"></image>
        <!-- 商品信息 -->
        <view class="nf-cart-info">
          <text class="nf-cart-name">{{ item.sku.name }}</text>
          <view class="nf-cart-specs" v-if="item.sku.specs && item.sku.specs.length">
            <text class="nf-cart-spec-tag">{{ item.sku.specs.map(s => s.value).join(' / ') }}</text>
          </view>
          <view class="nf-cart-specs" v-else-if="item.sku.description">
            <text class="nf-cart-spec-tag">{{ item.sku.description }}</text>
          </view>
          <view class="nf-cart-bottom">
            <text class="nf-cart-price">{{ cs }}{{ (item.sku.price / 100).toFixed(2) }}</text>
            <view v-if="!isDeleteAll" class="nf-cart-qty">
              <wu-number-box :asyncChange="true" :min="0" @change="(e)=>onChange(item,e)" integer v-model="item.quantity"></wu-number-box>
            </view>
            <view v-else class="nf-btn-del" @tap="deleteItem(item)">{{ $t('delete') || '删除' }}</view>
          </view>
        </view>
      </view>
    </view>
  </view>

  <!-- 底部操作栏 -->
  <view class="nf-cart-bar" v-if="isLoggedIn && cartList.length > 0">
    <view class="nf-cart-bar-left" v-if="!isDeleteAll">
      <view class="nf-cart-check" @tap="toggleSelectAll">
        <view :class="['nf-checkbox', { checked: isAllSelected }]">
          <text v-if="isAllSelected" class="nf-check-icon">✓</text>
        </view>
      </view>
      <text class="nf-cart-bar-all">{{ $t('selectAll') }}</text>
    </view>
    <view class="nf-cart-bar-left" v-else>
      <view class="nf-cart-bar-done" @tap="toggleDeleteMode">{{ $t('doneText') }}</view>
    </view>

    <view class="nf-cart-bar-right" v-if="!isDeleteAll">
      <view class="nf-cart-bar-total">
        <text class="nf-cart-bar-label">{{ $t('totalText') }}</text>
        <text class="nf-cart-bar-price">{{ cs }}{{ selectedTotal }}</text>
        <text class="nf-cart-bar-count">({{ selectedCount }}{{ $t('itemCount') }})</text>
      </view>
      <view class="nf-cart-bar-btns">
        <view class="nf-btn nf-btn-ghost-sm" @tap="toggleDeleteMode">{{ $t('editCart') }}</view>
        <view class="nf-btn nf-btn-primary" @tap="toSettlement">{{ $t('checkout') }}({{ selectedCount }})</view>
      </view>
    </view>
    <view class="nf-cart-bar-right" v-else>
      <view class="nf-btn nf-btn-danger" @tap="clearAllCart">{{ $t('clearCart') }}</view>
    </view>
  </view>
</template>

<script setup>
	import { ref, computed } from 'vue'
	import { getSelfCart, cutCart, addCart, clearCart } from "@/api/cart.js"
	import { onShow } from '@dcloudio/uni-app'
	import { useUserStore } from "@/pinia/modules/user"
	import { useLangStore } from '@/pinia/modules/lang.js'
	import { useAppConfigStore } from '@/pinia/modules/appConfig.js'
	import { placeOrderByCart } from '@/api/order.js'
	import { getUrl } from "@/utils/url.js"

	const { $t } = useLangStore()
	const appConfigStore = useAppConfigStore()
	const cs = computed(() => appConfigStore.currencySymbol)
	const cartList = ref([])
	const isLoggedIn = ref(false)

	const initPage = async () => {
		const userStore = useUserStore()
		const token = userStore.token || ''
		isLoggedIn.value = !!token
		if (token) {
			const res = await getSelfCart()
			if (res.code === 0) {
				cartList.value = (res.data || []).map(item => ({ ...item, _selected: true }))
			}
		}
	}
	onShow(() => { initPage() })

	const toggleSelect = (index) => { cartList.value[index]._selected = !cartList.value[index]._selected }
	const isAllSelected = computed(() => cartList.value.length > 0 && cartList.value.every(item => item._selected))
	const toggleSelectAll = () => { const v = !isAllSelected.value; cartList.value.forEach(item => { item._selected = v }) }
	const selectedItems = computed(() => cartList.value.filter(item => item._selected))
	const selectedCount = computed(() => selectedItems.value.length)
	const selectedTotal = computed(() => (selectedItems.value.reduce((t, i) => t + i.sku.price * i.quantity, 0) / 100).toFixed(2))

	const isDeleteAll = ref(false)
	const goTo = () => { uni.switchTab({ url: '/pages/tabBar/index' }) }
	const goLogin = () => { uni.navigateTo({ url: '/pages/user/login' }) }
	const goRegister = () => { uni.navigateTo({ url: '/pages/user/register' }) }

	const onChange = async (item, e) => {
		uni.showLoading({ title: '', mask: true })
		let res
		if (item.quantity > e.value) {
			res = await cutCart({ goodID: item.goodID, skuID: item.skuID, quantity: 1 })
		} else {
			res = await addCart({ goodID: item.goodID, skuID: item.skuID, quantity: 1 })
		}
		uni.hideLoading()
		if (res.code !== 0) { uni.showToast({ title: res.msg || '调整失败', icon: 'none' }); return }
		item.quantity = e.value
		if (e.value == 0) { cartList.value = cartList.value.filter(i => i.ID !== item.ID) }
	}

	const toggleDeleteMode = () => { isDeleteAll.value = !isDeleteAll.value }

	const deleteItem = async (item) => {
		uni.showModal({
			title: '', content: $t('deleteConfirmCart'),
			success: async (res) => {
				if (res.confirm) {
					uni.showLoading({ title: '', mask: true })
					const result = await cutCart({ goodID: item.goodID, skuID: item.skuID, quantity: item.quantity })
					uni.hideLoading()
					if (result.code === 0) {
						cartList.value = cartList.value.filter(i => i.ID !== item.ID)
						if (cartList.value.length === 0) isDeleteAll.value = false
					}
				}
			}
		})
	}

	const clearAllCart = async () => {
		uni.showModal({
			title: '', content: $t('clearConfirmCart'),
			success: async (res) => {
				if (res.confirm) {
					uni.showLoading({ title: '', mask: true })
					const result = await clearCart()
					uni.hideLoading()
					if (result.code === 0) { cartList.value = []; isDeleteAll.value = false }
				}
			}
		})
	}

	const toSettlement = async () => {
		const items = selectedItems.value
		if (!items.length) { uni.showToast({ title: $t('selectAll') || '请选择商品', icon: 'none' }); return }
		uni.showLoading({ title: '', mask: true })
		const res = await placeOrderByCart({ cartIDs: items.map(item => item.ID) })
		uni.hideLoading()
		if (res.code === 0) {
			const orderedIDs = new Set(items.map(i => i.ID))
			cartList.value = cartList.value.filter(i => !orderedIDs.has(i.ID))
			uni.navigateTo({ url: `/pages/orderInfo/orderInfo?orderID=${res.data.orderID}&type=cart` })
		} else {
			uni.showToast({ title: res.msg || '生成订单失败', icon: 'none' })
		}
	}
</script>

<style lang="scss" scoped>
/* ===== Netflix Dark Cart ===== */
.nf-cart {
  padding: 0 24rpx;
}

/* 空状态 */
.nf-cart-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  padding: 60rpx 40rpx;
}

.nf-cart-empty-icon {
  font-size: 120rpx;
  margin-bottom: 30rpx;
  opacity: 0.6;
}

.nf-cart-empty-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #fff;
  margin-bottom: 16rpx;
}

.nf-cart-empty-sub {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.4);
  margin-bottom: 50rpx;
  text-align: center;
}

.nf-cart-empty-btns {
  display: flex;
  gap: 24rpx;
}

/* 购物车列表 */
.nf-cart-list {
  padding-bottom: 20rpx;
}

.nf-cart-item {
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.06);
  border-radius: 20rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  backdrop-filter: blur(8px);
}

.nf-cart-check {
  margin-right: 16rpx;
  flex-shrink: 0;
}

.nf-checkbox {
  width: 44rpx;
  height: 44rpx;
  border-radius: 50%;
  border: 2rpx solid rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;

  &.checked {
    background: #e50914;
    border-color: #e50914;
  }
}

.nf-check-icon {
  color: #fff;
  font-size: 24rpx;
  font-weight: bold;
}

.nf-cart-img {
  width: 180rpx;
  height: 180rpx;
  border-radius: 12rpx;
  margin-right: 20rpx;
  flex-shrink: 0;
  border: 1rpx solid rgba(255, 255, 255, 0.06);
}

.nf-cart-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 180rpx;
  overflow: hidden;
}

.nf-cart-name {
  font-size: 28rpx;
  font-weight: 600;
  color: #fff;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 8rpx;
}

.nf-cart-specs {
  margin-bottom: 8rpx;
}

.nf-cart-spec-tag {
  font-size: 22rpx;
  color: rgba(255, 255, 255, 0.4);
  background: rgba(255, 255, 255, 0.06);
  padding: 4rpx 12rpx;
  border-radius: 6rpx;
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.nf-cart-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
}

.nf-cart-price {
  font-size: 32rpx;
  font-weight: 700;
  color: #e50914;
}

.nf-cart-qty {
  :deep(.wu-numberbox) {
    border: none;
    background: rgba(255, 255, 255, 0.06);
    border-radius: 10rpx;
    overflow: hidden;

    .wu-numberbox__minus,
    .wu-numberbox__plus {
      width: 56rpx;
      height: 56rpx;
      background: rgba(229, 9, 20, 0.2);
      color: #e50914;
      border: none;
      font-weight: 600;
    }

    .wu-numberbox__value {
      width: 64rpx;
      height: 56rpx;
      background: rgba(255, 255, 255, 0.04);
      color: #fff;
      margin: 0;
      font-size: 26rpx;
      font-weight: 600;
      border: none;
    }
  }
}

/* 底部操作栏 */
.nf-cart-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 120rpx;
  padding: 0 24rpx;
  background: rgba(20, 20, 20, 0.95);
  backdrop-filter: blur(24px);
  border-top: 1rpx solid rgba(255, 255, 255, 0.06);
  z-index: 100;
  /* #ifdef H5 */
  bottom: 88rpx;
  /* #endif */
}

.nf-cart-bar-left {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.nf-cart-bar-all {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.6);
  margin-left: 8rpx;
}

.nf-cart-bar-done {
  font-size: 28rpx;
  font-weight: 600;
  color: #e50914;
}

.nf-cart-bar-right {
  display: flex;
  align-items: center;
  flex: 1;
  justify-content: flex-end;
  gap: 16rpx;
}

.nf-cart-bar-total {
  display: flex;
  align-items: baseline;
  margin-right: 12rpx;
}

.nf-cart-bar-label {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.5);
  margin-right: 6rpx;
}

.nf-cart-bar-price {
  font-size: 34rpx;
  font-weight: 700;
  color: #e50914;
}

.nf-cart-bar-count {
  font-size: 20rpx;
  color: rgba(255, 255, 255, 0.3);
  margin-left: 6rpx;
}

.nf-cart-bar-btns {
  display: flex;
  gap: 12rpx;
}

/* 按钮系统 */
.nf-btn {
  height: 72rpx;
  line-height: 72rpx;
  padding: 0 36rpx;
  border-radius: 36rpx;
  font-size: 28rpx;
  font-weight: 600;
  text-align: center;
  transition: all 0.2s;

  &:active { transform: scale(0.96); }
}

.nf-btn-primary {
  background: #e50914;
  color: #fff;
}

.nf-btn-ghost {
  background: rgba(229, 9, 20, 0.12);
  border: 1rpx solid rgba(229, 9, 20, 0.3);
  color: #e50914;
}

.nf-btn-ghost-sm {
  height: 64rpx;
  line-height: 64rpx;
  padding: 0 24rpx;
  border-radius: 32rpx;
  font-size: 24rpx;
  background: rgba(255, 255, 255, 0.06);
  border: 1rpx solid rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.6);
}

.nf-btn-danger {
  background: #e50914;
  color: #fff;
}

.nf-btn-del {
  padding: 8rpx 24rpx;
  background: rgba(229, 9, 20, 0.15);
  color: #e50914;
  border-radius: 16rpx;
  font-size: 24rpx;
  font-weight: 500;
  border: 1rpx solid rgba(229, 9, 20, 0.3);
}
</style>
