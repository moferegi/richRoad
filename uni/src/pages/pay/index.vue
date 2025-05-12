<template>
  <view class="payment-page">
    <!-- 支付金额 -->
    <view class="payment-amount">
      <text class="amount-label">支付金额</text>
      <text class="amount-value">¥{{ amount }}</text>
    </view>

    <!-- 支付方式列表 -->
    <view class="payment-methods">
      <view
          class="payment-method-item"
          v-for="(method, index) in paymentMethods"
          :key="index"
          @click="selectPayMethod(index)"
      >
        <view class="method-left">
          <image class="method-icon" :src="method.iconUrl" mode="aspectFill" />
          <view class="method-info">
            <text class="method-name">{{ method.name }}</text>
            <text class="method-desc" v-if="method.desc">{{ method.desc }}</text>
          </view>
        </view>
        <view class="method-right">
          <radio :checked="selectedMethod === index" color="#ff4757" />
        </view>
      </view>
    </view>

    <!-- 确认支付按钮 -->
    <view class="payment-submit">
      <button class="submit-btn" @click="confirmPayment">确认支付</button>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue';

// 支付金额
const amount = ref('38.88');

// 选中的支付方式
const selectedMethod = ref(0);

// 支付方式列表
const paymentMethods = ref([
  {
    name: '微信支付',
    desc: '推荐使用微信支付',
    iconUrl: '/static/wxPay.png'
  },
  {
    name: '支付宝支付',
    iconUrl: '/static/alipay.png'
  },
  {
    name: '预存款支付',
    desc: '可用余额 ¥198.5',
    iconUrl: '/static/wallet.png'
  }
]);

// 选择支付方式
const selectPayMethod = (index) => {
  selectedMethod.value = index;
};

// 确认支付
const confirmPayment = () => {
  console.log('确认支付', {
    amount: amount.value,
    method: paymentMethods.value[selectedMethod.value].name
  });
};
</script>

<style lang="scss">
.payment-page {
  min-height: 100vh;
  background-color: #f8f8f8;
  padding: 30rpx;
}

.payment-amount {
  background-color: #fff;
  padding: 40rpx;
  border-radius: 12rpx;
  text-align: center;
  margin-bottom: 30rpx;
}

.amount-label {
  font-size: 28rpx;
  color: #666;
  margin-bottom: 20rpx;
  display: block;
}

.amount-value {
  font-size: 48rpx;
  color: #333;
  font-weight: bold;
}

.payment-methods {
  background-color: #fff;
  border-radius: 12rpx;
  padding: 0 30rpx;
}

.payment-method-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 30rpx 0;
  border-bottom: 1rpx solid #f5f5f5;

  &:last-child {
    border-bottom: none;
  }
}

.method-left {
  display: flex;
  align-items: center;
}

.method-icon {
  width: 80rpx;
  height: 80rpx;
  border-radius: 12rpx;
  margin-right: 20rpx;
  overflow: hidden;
  background-color: #f0f0f0; // 默认底色防止图片加载失败时留白
}

.method-info {
  display: flex;
  flex-direction: column;
}

.method-name {
  font-size: 28rpx;
  color: #333;
  margin-bottom: 6rpx;
}

.method-desc {
  font-size: 24rpx;
  color: #999;
}

.payment-submit {
  margin-top: 60rpx;
  padding: 0 30rpx;
}

.submit-btn {
  width: 100%;
  height: 88rpx;
  line-height: 88rpx;
  background-color: #ff4757;
  color: #fff;
  border-radius: 44rpx;
  font-size: 32rpx;
}
</style>
