<template>
  <!-- 遮罩 -->
  <view class="nf-sku-mask" v-if="skuShow" @tap="closeSku"></view>
  <!-- SKU弹窗 -->
  <view class="nf-sku-popup" :class="{ show: skuShow }">
    <!-- 顶部商品信息 -->
    <view class="nf-sku-header">
      <image class="nf-sku-cover" :src="currentCover" mode="aspectFill"></image>
      <view class="nf-sku-info">
        <text class="nf-sku-price">{{ cs }}{{ currentPrice }}</text>
        <text class="nf-sku-stock">{{ $t('stock') }}: {{ currentStock }}</text>
        <text class="nf-sku-selected" v-if="selectedSpecText">{{ selectedSpecText }}</text>
      </view>
      <view class="nf-sku-close" @tap="closeSku">
        <uni-icons type="close" size="18" color="rgba(255,255,255,0.6)"></uni-icons>
      </view>
    </view>

    <!-- 规格选择区 -->
    <scroll-view class="nf-sku-body" scroll-y>
      <view v-for="(group, gIdx) in specGroups" :key="gIdx" class="nf-spec-group">
        <text class="nf-spec-title">{{ $lt(group.label) }}</text>
        <view class="nf-spec-options">
          <view v-for="(opt, oIdx) in group.values" :key="oIdx"
            class="nf-spec-tag"
            :class="{
              'nf-spec-active': group.selected === opt.value,
              'nf-spec-disabled': opt.disabled
            }"
            @tap="selectSpec(gIdx, opt)">
            <text>{{ $lt(opt.value) }}</text>
          </view>
        </view>
      </view>

      <!-- 数量选择 -->
      <view class="nf-qty-section">
        <text class="nf-spec-title">{{ $t('itemCount') }}</text>
        <view class="nf-qty-stepper">
          <view class="nf-qty-btn" :class="{ 'nf-qty-disabled': quantity <= 1 }" @tap="changeQty(-1)">-</view>
          <text class="nf-qty-num">{{ quantity }}</text>
          <view class="nf-qty-btn" :class="{ 'nf-qty-disabled': quantity >= maxStock }" @tap="changeQty(1)">+</view>
        </view>
      </view>
    </scroll-view>

    <!-- 底部按钮 -->
    <view class="nf-sku-footer">
      <view class="nf-sku-btn-cart" @tap="handleAddCart">{{ $t('addToCart') }}</view>
      <view class="nf-sku-btn-buy" @tap="handleCheckout">{{ $t('goCheckout') }}</view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useUserStore } from '@/pinia/modules/user'
import { addCart } from '@/api/cart.js'
import { getUrl } from '@/utils/url.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'

const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
const cs = computed(() => appConfigStore.currencySymbol)
const $lt = computed(() => langStore.$lt)
const $t = computed(() => langStore.$t)

const props = defineProps({
  isCart: { type: Boolean, default: false },
  good: { type: Object, default: () => ({}) },
  selectedCoupon: { type: Object, default: () => ({}) }
})

const userStore = useUserStore()
const token = userStore.token || ''

const skuShow = ref(false)
const quantity = ref(1)
const specGroups = ref([])
const matchedSku = ref(null)

// 构建规格分组
const buildSpecGroups = () => {
  if (!props.good.skus || props.good.skus.length === 0) return
  const groupMap = {}
  props.good.skus.forEach(sku => {
    if (!sku.specs) return
    const specs = typeof sku.specs === 'string' ? JSON.parse(sku.specs) : sku.specs
    if (!Array.isArray(specs)) return
    specs.forEach(spec => {
      const key = JSON.stringify(spec.label)
      if (!groupMap[key]) {
        groupMap[key] = { label: spec.label, values: [], selected: null }
      }
      const valKey = JSON.stringify(spec.value)
      if (!groupMap[key].values.find(v => JSON.stringify(v.value) === valKey)) {
        groupMap[key].values.push({ value: spec.value, disabled: false })
      }
    })
  })
  specGroups.value = Object.values(groupMap)
  matchSku()
}

// 选择规格
const selectSpec = (gIdx, opt) => {
  if (opt.disabled) {
    const specText = specGroups.value.map((g, i) =>
      i === gIdx ? $lt.value(opt.value) : (g.selected ? $lt.value(g.selected) : '?')
    ).join(' + ')
    uni.showToast({ title: specText + ' ' + $t.value('skuSoldOut'), icon: 'none' })
    return
  }
  const group = specGroups.value[gIdx]
  group.selected = group.selected === opt.value ? null : opt.value
  matchSku()
}

// 匹配SKU
const matchSku = () => {
  const selections = specGroups.value.map(g => g.selected)
  const allSelected = selections.every(s => s !== null)

  if (allSelected && props.good.skus) {
    const found = props.good.skus.find(sku => {
      const specs = typeof sku.specs === 'string' ? JSON.parse(sku.specs) : sku.specs
      if (!Array.isArray(specs)) return false
      return specs.every((spec, i) => JSON.stringify(spec.value) === JSON.stringify(selections[i]))
    })
    matchedSku.value = found || null
  } else {
    matchedSku.value = null
  }

  updateDisabledState()
  if (matchedSku.value && quantity.value > matchedSku.value.inventory) {
    quantity.value = Math.max(1, matchedSku.value.inventory)
  }
}

// 更新不可选状态
const updateDisabledState = () => {
  specGroups.value.forEach((group, gIdx) => {
    group.values.forEach(opt => {
      const testSelections = specGroups.value.map((g, i) => {
        if (i === gIdx) return opt.value
        return g.selected
      })
      const hasStock = props.good.skus?.some(sku => {
        const specs = typeof sku.specs === 'string' ? JSON.parse(sku.specs) : sku.specs
        if (!Array.isArray(specs)) return false
        const match = specs.every((spec, i) => {
          if (testSelections[i] === null) return true
          return JSON.stringify(spec.value) === JSON.stringify(testSelections[i])
        })
        return match && sku.inventory > 0
      })
      opt.disabled = !hasStock
    })
  })
}

const currentCover = computed(() => {
  if (matchedSku.value?.picture) return getUrl(matchedSku.value.picture)
  return getUrl(props.good.imageUrl)
})

const currentPrice = computed(() => {
  if (matchedSku.value) return (matchedSku.value.price / 100).toFixed(2)
  if (props.good.price) return (props.good.price / 100).toFixed(2)
  return '0.00'
})

const currentStock = computed(() => {
  if (matchedSku.value) return matchedSku.value.inventory
  if (!props.good.skus) return 0
  return props.good.skus.reduce((t, s) => t + s.inventory, 0)
})

const maxStock = computed(() => {
  if (matchedSku.value) return matchedSku.value.inventory
  return 999
})

const selectedSpecText = computed(() => {
  const selected = specGroups.value.filter(g => g.selected).map(g => $lt.value(g.selected))
  return selected.length > 0 ? selected.join(' / ') : ''
})

const changeQty = (delta) => {
  const newVal = quantity.value + delta
  if (newVal < 1) return
  if (newVal > maxStock.value) {
    uni.showToast({ title: $t.value('skuStockOnly').replace('{n}', maxStock.value), icon: 'none' })
    return
  }
  quantity.value = newVal
}

const validateSelection = () => {
  const allSelected = specGroups.value.every(g => g.selected !== null)
  if (!allSelected) {
    uni.showToast({ title: $t.value('skuSelectFull'), icon: 'none' })
    return false
  }
  if (!matchedSku.value || matchedSku.value.inventory <= 0) {
    uni.showToast({ title: $t.value('skuSoldOut'), icon: 'none' })
    return false
  }
  return true
}

const handleAddCart = async () => {
  if (!token) {
    uni.showToast({ title: $t.value('pleaseLogin'), mask: true, icon: 'none' })
    uni.redirectTo({ url: '/pages/user/login' })
    return
  }
  if (!validateSelection()) return
  const res = await addCart({
    goodID: props.good.ID,
    skuID: matchedSku.value.ID,
    quantity: quantity.value
  })
  if (res.code === 0) {
    uni.showToast({ title: $t.value('addCartSuccess'), mask: true, icon: 'none' })
    closeSku()
  } else {
    uni.showToast({ title: $t.value('addCartFail'), mask: true, icon: 'none' })
  }
}

const handleCheckout = () => {
  if (!token) {
    uni.showToast({ title: $t.value('pleaseLogin'), mask: true, icon: 'none' })
    uni.redirectTo({ url: '/pages/user/login' })
    return
  }
  if (!validateSelection()) return
  closeSku()
  const query = `goodID=${props.good.ID}&skuID=${matchedSku.value.ID}&quantity=${quantity.value}&couponNum=${props.selectedCoupon.couponNum || ''}`
  uni.navigateTo({ url: `/pages/orderInfo/orderInfo?${query}` })
}

const showSku = () => {
  quantity.value = 1
  buildSpecGroups()
  skuShow.value = true
}
const closeSku = () => { skuShow.value = false }

defineExpose({ showSku, closeSku })
</script>

<style lang="scss" scoped>
.nf-sku-mask {
  position: fixed; top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.7); z-index: 1000;
}
.nf-sku-popup {
  position: fixed; left: 0; right: 0; bottom: -100vh; z-index: 1001;
  background: #141414; border-radius: 24rpx 24rpx 0 0;
  transition: all 0.3s ease;
  max-height: 80vh; display: flex; flex-direction: column;
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
  &.show { bottom: 0; }
}
.nf-sku-header {
  display: flex; padding: 28rpx 32rpx; gap: 24rpx;
  border-bottom: 1rpx solid rgba(255,255,255,0.06);
  position: relative;
}
.nf-sku-cover {
  width: 180rpx; height: 180rpx; border-radius: 16rpx;
  border: 1rpx solid rgba(255,255,255,0.08);
}
.nf-sku-info { flex: 1; display: flex; flex-direction: column; justify-content: center; gap: 8rpx; }
.nf-sku-price { font-size: 44rpx; font-weight: 800; color: #e50914; }
.nf-sku-stock { font-size: 24rpx; color: rgba(255,255,255,0.4); }
.nf-sku-selected { font-size: 24rpx; color: rgba(255,255,255,0.6); }
.nf-sku-close {
  position: absolute; right: 24rpx; top: 24rpx;
  width: 56rpx; height: 56rpx; display: flex; align-items: center; justify-content: center;
  background: rgba(255,255,255,0.08); border-radius: 50%;
}
.nf-sku-body { flex: 1; padding: 24rpx 32rpx; overflow-y: auto; }
.nf-spec-group { margin-bottom: 32rpx; }
.nf-spec-title { display: block; font-size: 26rpx; color: rgba(255,255,255,0.6); margin-bottom: 16rpx; font-weight: 600; }
.nf-spec-options { display: flex; flex-wrap: wrap; gap: 16rpx; }
.nf-spec-tag {
  padding: 12rpx 28rpx; border-radius: 10rpx;
  background: rgba(255,255,255,0.06); border: 1rpx solid rgba(255,255,255,0.1);
  font-size: 26rpx; color: rgba(255,255,255,0.7);
  transition: all 0.2s;
}
.nf-spec-active {
  background: rgba(229,9,20,0.15); border-color: #e50914; color: #e50914; font-weight: 600;
}
.nf-spec-disabled { opacity: 0.3; text-decoration: line-through; }
.nf-qty-section {
  display: flex; align-items: center; justify-content: space-between;
  padding: 24rpx 0; border-top: 1rpx solid rgba(255,255,255,0.06);
}
.nf-qty-stepper { display: flex; align-items: center; gap: 4rpx; }
.nf-qty-btn {
  width: 56rpx; height: 56rpx; display: flex; align-items: center; justify-content: center;
  background: rgba(255,255,255,0.08); border-radius: 10rpx; font-size: 32rpx;
  color: #fff; font-weight: 600;
}
.nf-qty-disabled { opacity: 0.3; }
.nf-qty-num { min-width: 72rpx; text-align: center; font-size: 30rpx; color: #fff; font-weight: 600; }
.nf-sku-footer {
  display: flex; padding: 20rpx 32rpx; gap: 20rpx;
  border-top: 1rpx solid rgba(255,255,255,0.06);
}
.nf-sku-btn-cart, .nf-sku-btn-buy {
  flex: 1; height: 88rpx; display: flex; align-items: center; justify-content: center;
  border-radius: 12rpx; font-size: 30rpx; font-weight: 700; color: #fff;
}
.nf-sku-btn-cart { background: rgba(255,149,0,0.9); }
.nf-sku-btn-buy { background: #e50914; }
</style>

