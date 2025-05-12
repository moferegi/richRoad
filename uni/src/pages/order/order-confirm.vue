<template>
  <view class="order-confirm">
    <scroll-view class="content-scroll" scroll-y>
    <!-- 收货地址 -->
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

    <!-- 店铺商品信息 -->
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
              <text class="goods-price">¥{{ item.price }}</text>
              <text class="goods-count">x{{ item.count }}</text>
            </view>
          </view>
        </view>
      </view>
    </view>

    <!-- 优惠券选择 -->
    <view class="coupon-section" @click="chooseCoupon">
      <text class="section-title">优惠券</text>
      <view class="section-content">
        <text class="coupon-text">{{ selectedCoupon ? `省${selectedCoupon.amount}元` : '选择优惠券' }}</text>
        <uni-icons type="right" size="16" color="#999"></uni-icons>
      </view>
    </view>

    <!-- 商家促销 -->
    <view class="promotion-section">
      <text class="section-title">商家促销</text>
      <view class="section-content">
        <text class="promotion-text">暂无可用优惠</text>
      </view>
    </view>

    <!-- 订单备注 -->
    <view class="remark-section">
      <text class="section-title">备注</text>
      <input class="remark-input" type="text" v-model="remark" placeholder="请填写备注信息" />
    </view>

    <!-- 金额明细 -->
    <view class="amount-section">
      <view class="amount-item">
        <text>商品金额</text>
        <text>¥{{ totalAmount }}</text>
      </view>
      <view class="amount-item">
        <text>优惠券</text>
        <text class="discount">-¥{{ couponAmount }}</text>
      </view>
      <view class="amount-item">
        <text>运费</text>
        <text>{{ shipping === 0 ? '免运费' : `¥${shipping}` }}</text>
      </view>
    </view>
    </scroll-view>

    <!-- 底部提交栏 -->
    <view class="submit-bar">
      <view class="total-wrap">
        <text>实付款：</text>
        <text class="total-amount">¥{{ finalAmount }}</text>
      </view>
      <button class="submit-btn" @click="submitOrder">提交订单</button>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue';

// Mock数据
const address = ref({
  name: '许小星',
  phone: '13853989563',
  detail: '山东省济南市历城区 149号'
});

const store = ref({
  name: '西城小店铺',
  goods: [
    {
      title: '古缘妃 短袖t恤女夏装2019新款',
      spec: '春装款 L',
      price: '17.8',
      count: 1
    },
    {
      title: '韩版牛星洞洞裙鞋 夏季浴室防滑简约',
      spec: '春装款 L',
      price: '17.8',
      count: 1
    }, {
      title: '古缘妃 短袖t恤女夏装2019新款',
      spec: '春装款 L',
      price: '17.8',
      count: 1
    },
    {
      title: '韩版牛星洞洞裙鞋 夏季浴室防滑简约',
      spec: '春装款 L',
      price: '17.8',
      count: 1
    }, {
      title: '古缘妃 短袖t恤女夏装2019新款',
      spec: '春装款 L',
      price: '17.8',
      count: 1
    },
    {
      title: '韩版牛星洞洞裙鞋 夏季浴室防滑简约',
      spec: '春装款 L',
      price: '17.8',
      count: 1
    }, {
      title: '古缘妃 短袖t恤女夏装2019新款',
      spec: '春装款 L',
      price: '17.8',
      count: 1
    },
    {
      title: '韩版牛星洞洞裙鞋 夏季浴室防滑简约',
      spec: '春装款 L',
      price: '17.8',
      count: 1
    },
  ]
});

const selectedCoupon = ref(null);
const remark = ref('');
const shipping = ref(0);

// 计算总金额
const totalAmount = computed(() => {
  return store.value.goods.reduce((total, item) => {
    return total + Number(item.price) * item.count;
  }, 0).toFixed(2);
});

// 优惠券金额
const couponAmount = computed(() => {
  return selectedCoupon.value ? selectedCoupon.value.amount : 0;
});

// 最终支付金额
const finalAmount = computed(() => {
  return (Number(totalAmount.value) - couponAmount.value + shipping.value).toFixed(2);
});

// 选择地址
const chooseAddress = () => {
  console.log('选择地址');
};

// 选择优惠券
const chooseCoupon = () => {
  console.log('选择优惠券');
};

// 提交订单
const submitOrder = () => {
  uni.navigateTo({
    url: '/pages/pay/index',
    success: () => console.log('跳转成功'),
    fail: (err) => console.error('跳转失败', err)
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
  height: 0; // 这一行很重要，确保内容区域可以正确滚动
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
  position: relative; // 改为相对定位
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
