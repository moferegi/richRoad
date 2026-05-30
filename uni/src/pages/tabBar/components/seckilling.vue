<template>
  <view class="flash-sale">
    <view class="flash-header">
      <image class="top-sell" src="@/static/top.png" mode="aspectFill"></image>
      <text class="flash-title">{{ $t('hotSelling') }}</text>
    </view>

    <view class="product-list">
      <view class="product-item" v-for="item in productViews" :key="item._productKey" @click="goto(item)">
          <image class="product-image" :src="item._imageUrl" mode="aspectFill"></image>
        <view class="desc">
          <text class="product-title">{{ item._titleText }}</text>
          <text class="product-price">{{ cs }} {{ item._priceText }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed } from 'vue'
import {getUrl} from "@/utils/url";
import { formatLocalizedPrice } from '@/utils/price-i18n.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'

const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
const cs = computed(() => appConfigStore.currencySymbol)
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)
const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')

const props = defineProps({
  productData: {
    type: Array,
    default: () => [] // 修改为函数返回空数组
  }
})

const formatPrice = (item) => formatLocalizedPrice(item?.price, item?.priceI18n, locale.value)

// 热卖横滑列表只读展示字段；提前生成图片、标题、价格和稳定 key，降低横向滚动时的重复计算。
const productViews = computed(() => (props.productData || []).map((item, index) => ({
  ...item,
  _productKey: `${item?.ID || item?.id || item?.imageUrl || 'product'}-${index}`,
  _imageUrl: getUrl(item?.imageUrl || ''),
  _titleText: $lt.value(item?.title) || '',
  _priceText: formatPrice(item),
})))

const goto = (item) => {
  uni.navigateTo({
    url: '/pages/player/index?id=' + item.ID
  })
}
const init = () => {
  // 初始化逻辑，如果需要的话
  console.log(props.productData);
}
setTimeout(() => {
  init();
}, 1000);
</script>

<style scoped lang="scss">
.flash-sale {
  margin: 16rpx 28rpx;
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.06);
  border-radius: 28rpx;
  padding: 28rpx;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
}

.flash-header {
  display: flex;
  align-items: center;
  margin-bottom: 24rpx;
  gap: 16rpx;

  .top-sell {
    width: 40rpx;
    height: 40rpx;
  }
}

.flash-title {
  font-size: 32rpx;
  font-weight: 800;
  color: #e50914;
  letter-spacing: 2rpx;
  margin-right: auto;
}

.product-list {
  display: flex;
  overflow-x: auto;
  padding-bottom: 8rpx;
  gap: 20rpx;
  scrollbar-width: none;
  &::-webkit-scrollbar { display: none; }

  .desc {
    padding: 12rpx 4rpx;
    display: flex;
    flex-direction: column;
    gap: 6rpx;
  }
}

.product-item {
  width: 240rpx;
  flex-shrink: 0;
  transition: transform 0.3s;

  &:active { transform: scale(0.97); }
}

.product-image {
  width: 240rpx;
  height: 300rpx;
  background: #111;
  border-radius: 16rpx;
  box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.4);
}

.product-title {
  font-size: 26rpx;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-weight: 500;
}

.product-price {
  font-size: 26rpx;
  color: #e50914;
  font-weight: 800;
}
</style>

