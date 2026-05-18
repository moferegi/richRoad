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
        <image class="nf-cart-img" :src="item.sku.externalPicturePath ? getExternalUrl(item.sku.externalPicturePath) : getUrl(item.sku.picture)" mode="aspectFill"></image>
        <!-- 商品信息 -->
        <view class="nf-cart-info">
          <text class="nf-cart-name">{{ resolveDisplayText(item?.sku?.nameI18n || item?.sku?.name, item?.sku?.name) }}</text>
          <view class="nf-cart-specs" v-if="getItemSpecsText(item)">
            <text class="nf-cart-spec-tag">{{ getItemSpecsText(item) }}</text>
          </view>
          <view class="nf-cart-specs" v-else-if="item?.sku?.description || item?.sku?.descriptionI18n">
            <text class="nf-cart-spec-tag">{{ resolveDisplayText(item?.sku?.descriptionI18n || item?.sku?.description, item?.sku?.description) }}</text>
          </view>
          <view class="nf-cart-bottom">
            <text class="nf-cart-price">{{ cs }}{{ formatCartItemPrice(item) }}</text>
            <view v-if="!isDeleteAll" class="nf-cart-qty">
              <wu-number-box :asyncChange="true" :min="0" @change="(e)=>onChange(item,e)" integer v-model="item.quantity"></wu-number-box>
            </view>
            <view v-else class="nf-btn-del" @tap="deleteItem(item)">{{ $t('deleteText') }}</view>
          </view>
        </view>
      </view>
    </view>
  </view>

  <!-- 底部操作栏 -->
  <view class="nf-cart-bar" :class="{ 'nf-cart-bar--stack': !isDeleteAll && isBarStacked }" v-if="isLoggedIn && cartList.length > 0">
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

    <view class="nf-cart-bar-right" :class="{ 'nf-cart-bar-right--stack': isBarStacked }" v-if="!isDeleteAll">
      <view class="nf-cart-bar-total" :class="{ 'nf-cart-bar-total--stack': isBarStacked }">
        <text class="nf-cart-bar-label">{{ $t('totalText') }}</text>
        <text class="nf-cart-bar-price">{{ cs }}{{ selectedTotal }}</text>
        <text class="nf-cart-bar-count">({{ selectedCount }}{{ $t('itemCount') }})</text>
      </view>
      <view class="nf-cart-bar-btns">
        <view class="nf-btn nf-btn-ghost-sm nf-btn-edit" @tap="toggleDeleteMode">{{ $t('editCart') }}</view>
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
  import { resolveApiMessage } from '@/utils/i18n.js'
  import { useI18nDisplay } from '@/composables/useI18nDisplay.js'
  import { formatLocalizedPrice, resolveLocalizedPriceFen } from '@/utils/price-i18n.js'
	import { getUrl, getExternalUrl } from "@/utils/url.js"

  const langStore = useLangStore()
  const { $t } = langStore
  const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')

  const { resolveDisplayText } = useI18nDisplay(locale)

  const parseSpecItems = (payload) => {
    if (Array.isArray(payload)) {
      return payload
    }
    if (typeof payload === 'string' && payload.trim()) {
      try {
        const parsed = JSON.parse(payload)
        return Array.isArray(parsed) ? parsed : []
      } catch {
        return []
      }
    }
    return []
  }

  const getItemSpecsText = (item) => {
    const specs = parseSpecItems(item?.sku?.specs)
    if (!specs.length) return ''
    return specs
      .map(s => resolveDisplayText(s?.valueI18n || s?.value, s?.value || ''))
      .filter(Boolean)
      .join(' / ')
  }
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
  const getCartItemPriceFen = (item) => resolveLocalizedPriceFen(item?.sku?.price, item?.sku?.priceI18n, locale.value)
  const selectedTotalFen = computed(() => selectedItems.value.reduce((total, item) => {
    return total + getCartItemPriceFen(item) * Number(item?.quantity || 0)
  }, 0))
  const selectedTotal = computed(() => (selectedTotalFen.value / 100).toFixed(2))
  const formatCartItemPrice = (item) => formatLocalizedPrice(item?.sku?.price, item?.sku?.priceI18n, locale.value)
  const isBarStacked = computed(() => {
    const digits = String(selectedTotal.value || '').replace(/[^0-9]/g, '').length
    return digits >= 6
  })

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
    if (res.code !== 0) { uni.showToast({ title: resolveApiMessage(res.msg, 'adjustFail'), icon: 'none' }); return }
		item.quantity = e.value
		if (e.value == 0) { cartList.value = cartList.value.filter(i => i.ID !== item.ID) }
	}

	const toggleDeleteMode = () => { isDeleteAll.value = !isDeleteAll.value }

	const deleteItem = async (item) => {
		uni.showModal({
			title: '', content: $t('deleteConfirmCart'),
      cancelText: $t('cancel'),
      confirmText: $t('confirm'),
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
      cancelText: $t('cancel'),
      confirmText: $t('confirm'),
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
		if (!items.length) { uni.showToast({ title: $t('selectGoods'), icon: 'none' }); return }
		uni.showLoading({ title: '', mask: true })
    const res = await placeOrderByCart({
      cartIDs: items.map(item => item.ID),
      settlementCurrency: locale.value,
      settlementCurrencySymbol: cs.value,
    })
		uni.hideLoading()
		if (res.code === 0) {
			const orderedIDs = new Set(items.map(i => i.ID))
			cartList.value = cartList.value.filter(i => !orderedIDs.has(i.ID))
			uni.navigateTo({ url: `/pages/orderInfo/orderInfo?orderID=${res.data.orderID}&type=cart` })
		} else {
      uni.showToast({ title: resolveApiMessage(res.msg, 'createOrderFail'), icon: 'none' })
		}
	}
</script>

<style lang="scss" scoped>
/* ===== Light Premium Cart ===== */
.nf-cart {
  padding: 6rpx 8rpx 24rpx;
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
  opacity: 0.75;
}

.nf-cart-empty-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #0f172a;
  margin-bottom: 16rpx;
}

.nf-cart-empty-sub {
  font-size: 26rpx;
  color: rgba(15, 23, 42, 0.56);
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
  background: rgba(255, 255, 255, 0.92);
  border: 1rpx solid rgba(15, 23, 42, 0.08);
  border-radius: 20rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  backdrop-filter: blur(8px);
  box-shadow: 0 14rpx 28rpx rgba(15, 23, 42, 0.07);
}

.nf-cart-check {
  margin-right: 16rpx;
  flex-shrink: 0;
}

.nf-checkbox {
  width: 44rpx;
  height: 44rpx;
  border-radius: 50%;
  border: 2rpx solid rgba(15, 23, 42, 0.22);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;

  &.checked {
    background: linear-gradient(135deg, #2563eb, #0ea5e9);
    border-color: transparent;
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
  border-radius: 14rpx;
  margin-right: 20rpx;
  flex-shrink: 0;
  border: 1rpx solid rgba(15, 23, 42, 0.08);
  box-shadow: 0 10rpx 20rpx rgba(15, 23, 42, 0.08);
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
  color: #0f172a;
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
  color: rgba(15, 23, 42, 0.54);
  background: rgba(15, 23, 42, 0.06);
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
  color: #1d4ed8;
}

.nf-cart-qty {
  :deep(.wu-numberbox) {
    border: 1rpx solid rgba(15, 23, 42, 0.1);
    background: rgba(248, 250, 252, 0.96);
    border-radius: 10rpx;
    overflow: hidden;

    .wu-numberbox__minus,
    .wu-numberbox__plus {
      width: 56rpx;
      height: 56rpx;
      background: rgba(219, 234, 254, 0.9);
      color: #2563eb;
      border: none;
      font-weight: 600;
    }

    .wu-numberbox__value {
      width: 64rpx;
      height: 56rpx;
      background: rgba(255, 255, 255, 0.96);
      color: #0f172a;
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
  background: rgba(248, 250, 252, 0.94);
  backdrop-filter: blur(24px);
  border-top: 1rpx solid rgba(15, 23, 42, 0.08);
  z-index: 100;
  /* #ifdef H5 */
  bottom: 88rpx;
  /* #endif */
}

.nf-cart-bar--stack {
  height: auto;
  min-height: 214rpx;
  align-items: flex-start;
  padding-top: 14rpx;
  padding-bottom: 14rpx;
}

.nf-cart-bar--stack .nf-cart-bar-left {
  padding-top: 6rpx;
}

.nf-cart-bar-left {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.nf-cart-bar-all {
  font-size: 26rpx;
  color: rgba(15, 23, 42, 0.64);
  margin-left: 8rpx;
}

.nf-cart-bar-done {
  font-size: 28rpx;
  font-weight: 600;
  color: #2563eb;
}

.nf-cart-bar-right {
  display: flex;
  align-items: center;
  flex: 1;
  justify-content: flex-end;
  gap: 16rpx;
}

.nf-cart-bar-right--stack {
  flex-direction: column;
  align-items: stretch;
  justify-content: center;
  gap: 12rpx;
}

.nf-cart-bar-total {
  display: flex;
  align-items: baseline;
  margin-right: 12rpx;
}

.nf-cart-bar-total--stack {
  margin-right: 0;
  width: 100%;
  justify-content: flex-end;
}

.nf-cart-bar-label {
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.48);
  margin-right: 6rpx;
}

.nf-cart-bar-price {
  font-size: 34rpx;
  font-weight: 700;
  color: #1d4ed8;
}

.nf-cart-bar-count {
  font-size: 20rpx;
  color: rgba(15, 23, 42, 0.44);
  margin-left: 6rpx;
}

.nf-cart-bar-btns {
  display: flex;
  gap: 12rpx;
}

.nf-cart-bar-right--stack .nf-cart-bar-btns {
  width: 280rpx;
  margin-left: auto;
  display: flex;
  flex-direction: column;
  gap: 10rpx;
}

.nf-cart-bar-right--stack .nf-cart-bar-btns .nf-btn,
.nf-cart-bar-right--stack .nf-cart-bar-btns .nf-btn-ghost-sm,
.nf-cart-bar-right--stack .nf-cart-bar-btns .nf-btn-primary {
  width: 100%;
  box-sizing: border-box;
  text-align: center;
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
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  color: #fff;
  box-shadow: 0 10rpx 22rpx rgba(37, 99, 235, 0.26);
}

.nf-btn-ghost {
  background: rgba(255, 255, 255, 0.82);
  border: 1rpx solid rgba(15, 23, 42, 0.12);
  color: #0f172a;
}

.nf-btn-ghost-sm {
  height: 64rpx;
  line-height: 64rpx;
  padding: 0 24rpx;
  border-radius: 32rpx;
  font-size: 24rpx;
  background: rgba(255, 255, 255, 0.92);
  border: 1rpx solid rgba(15, 23, 42, 0.14);
  color: rgba(15, 23, 42, 0.72);
}

.nf-btn-edit {
  background: linear-gradient(90deg, #f59e0b, #f97316);
  border: none;
  color: #fff;
  box-shadow: 0 8rpx 18rpx rgba(249, 115, 22, 0.28);
}

.nf-btn-danger {
  background: #dc2626;
  color: #fff;
}

.nf-btn-del {
  padding: 8rpx 24rpx;
  background: rgba(220, 38, 38, 0.12);
  color: #dc2626;
  border-radius: 16rpx;
  font-size: 24rpx;
  font-weight: 500;
  border: 1rpx solid rgba(220, 38, 38, 0.25);
}
</style>
