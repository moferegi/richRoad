<template>
  <view class="order-confirm">
    <scroll-view class="content-scroll" scroll-y>
    <!-- Address -->
    <view class="address-section">
      <view class="address-content" @click="chooseAddress">
        <view class="left">
          <view class="icon" :style="{ backgroundColor: '#ff4757' }"></view>
        </view>
        <view class="center">
          <view class="user-info">
            <text class="name">{{ address.name }}</text>
            <text class="phone">{{ address.phone }}</text>
          </view>
          <view class="address-detail">{{ address.detail }}</view>
        </view>
        <view class="right">
          <uni-icons type="right" size="16" color="#999"></uni-icons>
        </view>
      </view>
    </view>

    <!-- Store goods -->
    <view class="store-section">
      <view class="store-header">
        <view class="store-avatar" :style="{ backgroundColor: '#ddd' }"></view>
        <text class="store-name">{{ store.name }}</text>
      </view>

      <view class="goods-list">
        <view class="goods-item" v-for="(item, index) in store.goods" :key="index">
          <view class="goods-image" :style="{ backgroundColor: '#f5f5f5' }"></view>
          <view class="goods-info">
            <text class="goods-title">{{ item.title }}</text>
            <text class="goods-spec">{{ item.spec }}</text>
            <view class="goods-price-wrap">
              <text class="goods-price">{{ cs }}{{ item.price }}</text>
              <text class="goods-count">x{{ item.count }}</text>
            </view>
          </view>
        </view>
      </view>
    </view>

    <!-- Coupon -->
    <view class="coupon-section" @click="chooseCoupon">
      <text class="section-title">{{ $t('couponSection') }}</text>
      <view class="section-content">
        <text class="coupon-text">{{ selectedCoupon ? $t('saveCouponAmount').replace('{}', selectedCoupon.amount) : $t('selectCouponText') }}</text>
        <uni-icons type="right" size="16" color="#999"></uni-icons>
      </view>
    </view>

    <!-- Promotion -->
    <view class="promotion-section">
      <text class="section-title">{{ $t('merchantPromotion') }}</text>
      <view class="section-content">
        <text class="promotion-text">{{ $t('noPromotion') }}</text>
      </view>
    </view>

    <!-- Remark -->
    <view class="remark-section">
      <text class="section-title">{{ $t('remark') }}</text>
      <input class="remark-input" type="text" v-model="remark" :placeholder="$t('remarkPlaceholder')" />
    </view>

    <!-- Amount summary -->
    <view class="amount-section">
      <view class="amount-item">
        <text>{{ $t('goodsAmount') }}</text>
        <text>{{ cs }}{{ totalAmount }}</text>
      </view>
      <view class="amount-item">
        <text>{{ $t('couponDiscount') }}</text>
        <text class="discount">-{{ cs }}{{ couponAmount }}</text>
      </view>
      <view class="amount-item">
        <text>{{ $t('shippingFee') }}</text>
        <text>{{ shipping === 0 ? $t('freeShippingText') : `${cs}${shipping}` }}</text>
      </view>
    </view>
    </scroll-view>

    <!-- Submit bar -->
    <view class="submit-bar">
      <view class="total-wrap">
        <text>{{ $t('actualPay') }}</text>
        <text class="total-amount">{{ cs }}{{ finalAmount }}</text>
      </view>
      <button class="submit-btn" @click="submitOrder">{{ $t('submitOrder') }}</button>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'
import { useLangStore } from '@/pinia/modules/lang.js'
const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const appConfigStore = useAppConfigStore()
const cs = computed(() => appConfigStore.currencySymbol)

const address = computed(() => ({
  name: $t.value('orderConfirmMockName'),
  phone: '13853989563',
  detail: $t.value('orderConfirmMockAddress')
}));

const buildMockGood = (titleKey) => ({
  title: $t.value(titleKey),
  spec: $t.value('orderConfirmMockGoodsSpec'),
  price: '17.8',
  count: 1
});

const mockGoods = computed(() => ([
  buildMockGood('orderConfirmMockGoodsTitleA'),
  buildMockGood('orderConfirmMockGoodsTitleB'),
  buildMockGood('orderConfirmMockGoodsTitleA'),
  buildMockGood('orderConfirmMockGoodsTitleB')
]));

const store = computed(() => ({
  name: $t.value('orderConfirmMockStoreName'),
  goods: mockGoods.value
}));

const selectedCoupon = ref(null);
const remark = ref('');
const shipping = ref(0);

// Total amount
const totalAmount = computed(() => {
  return store.value.goods.reduce((total, item) => {
    return total + Number(item.price) * item.count;
  }, 0).toFixed(2);
});

// Coupon amount
const couponAmount = computed(() => {
  return selectedCoupon.value ? selectedCoupon.value.amount : 0;
});

// Final payment amount
const finalAmount = computed(() => {
  return (Number(totalAmount.value) - couponAmount.value + shipping.value).toFixed(2);
});

// Choose address
const chooseAddress = () => {
  uni.showToast({ title: $t.value('selectAddress'), icon: 'none' });
};

// Choose coupon
const chooseCoupon = () => {
  uni.showToast({ title: $t.value('selectCoupon'), icon: 'none' });
};

// Submit order
const submitOrder = () => {
  uni.navigateTo({
    url: '/pages/pay/index'
  });
};
</script>

<style lang="scss">
.order-confirm {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f8f8f8;
}

.content-scroll {
  flex: 1;
  height: 0; // Required to keep the scroll area working
}

.address-section {
  background-color: #fff;
  margin-bottom: 20rpx;
  padding: 30rpx;
}

.address-content {
  display: flex;
  align-items: center;
}

.left {
  margin-right: 20rpx;
}

.icon {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
}

.center {
  flex: 1;
}

.user-info {
  margin-bottom: 10rpx;
}

.name {
  margin-right: 20rpx;
  font-weight: bold;
}

.address-detail {
  color: #666;
  font-size: 28rpx;
}

.submit-bar {
  position: relative; // Keep layout stable before fixed bar styles below
  height: 100rpx;
  background-color: #fff;
  display: flex;
  align-items: center;
  padding: 0 30rpx;
  box-shadow: 0 -2rpx 10rpx rgba(0,0,0,0.05);
}

.store-section {
  background-color: #fff;
  margin-bottom: 20rpx;
  padding: 30rpx;
}

.store-header {
  display: flex;
  align-items: center;
  margin-bottom: 30rpx;
}

.store-avatar {
  width: 60rpx;
  height: 60rpx;
  border-radius: 50%;
  margin-right: 20rpx;
}

.store-name {
  font-size: 28rpx;
  font-weight: bold;
}

.goods-list {
  .goods-item {
    display: flex;
    margin-bottom: 30rpx;

    &:last-child {
      margin-bottom: 0;
    }
  }
}

.goods-image {
  width: 160rpx;
  height: 160rpx;
  border-radius: 12rpx;
  margin-right: 20rpx;
}

.goods-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.goods-title {
  font-size: 28rpx;
  color: #333;
}

.goods-spec {
  font-size: 24rpx;
  color: #999;
}

.goods-price-wrap {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.goods-price {
  color: #ff4757;
  font-weight: bold;
}

.goods-count {
  color: #999;
  font-size: 24rpx;
}

.coupon-section, .promotion-section, .remark-section {
  background-color: #fff;
  margin-bottom: 20rpx;
  padding: 30rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  color: #333;
  font-size: 28rpx;
}

.section-content {
  flex: 1;
  text-align: right;
  color: #666;
  font-size: 28rpx;
}

.remark-input {
  flex: 1;
  text-align: right;
  font-size: 28rpx;
}

.amount-section {
  background-color: #fff;
  padding: 30rpx 30rpx 200rpx 30rpx;
}

.amount-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20rpx;
  font-size: 28rpx;

  &:last-child {
    margin-bottom: 0;
  }
}

.discount {
  color: #ff4757;
}

.submit-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 100rpx;
  background-color: #fff;
  display: flex;
  align-items: center;
  padding: 0 30rpx;
  box-shadow: 0 -2rpx 10rpx rgba(0,0,0,0.05);
}

.total-wrap {
  flex: 1;
  font-size: 28rpx;
}

.total-amount {
  color: #ff4757;
  font-size: 36rpx;
  font-weight: bold;
}

.submit-btn {
  width: 240rpx;
  height: 80rpx;
  line-height: 80rpx;
  text-align: center;
  background-color: #ff4757;
  color: #fff;
  border-radius: 40rpx;
  margin-left: 30rpx;
  font-size: 28rpx;
}
</style>
