<template>
  <view class="swiper-section">
    <swiper class="swiper" circular autoplay :indicator-dots="true" indicator-color="rgba(255,255,255,0.6)" indicator-active-color="#fff">
      <swiper-item v-for="item in normalizedSwiperList" :key="item._swiperKey" @click="handleSwiperClick(item)">
        <view class="swiper-item">
          <image class="swiper-image" :src="item._imageUrl" mode="aspectFill"></image>
          <!-- 遮罩文字层 -->
          <view v-if="item._maskText" class="swiper-mask" :style="item._maskStyle">
            <text class="swiper-mask-text"
              :style="item._maskTextStyle">{{ item._maskText }}</text>
          </view>
        </view>
      </swiper-item>
    </swiper>
  </view>
</template>

<script setup>
import {getUrl, getExternalUrl} from "@/utils/url.js"
import { useLangStore } from '@/pinia/modules/lang.js'
import { computed } from 'vue'

const langStore = useLangStore()
const $lt = computed(() => langStore.$lt)

const props = defineProps({
  lists: {
    type: Array,
    default: () => []
  }
})

// 首页轮播只读展示字段；集中归一化可避免模板每次渲染重复解析图片、多语言遮罩和样式对象。
const normalizedSwiperList = computed(() => (props.lists || []).map((item, index) => {
  const maskText = item?.maskEnabled ? ($lt.value(item?.maskText) || '') : ''
  const maskTextAlign = item?.maskTextAlign || 'center'
  return {
    ...item,
    _swiperKey: `${item?.ID || item?.id || item?.src || item?.externalPath || 'swiper'}-${index}`,
    _imageUrl: item?.externalPath ? getExternalUrl(item.externalPath) : getUrl(item?.src || ''),
    _maskText: maskText,
    _maskStyle: {
      height: `${item?.maskHeight || 40}px`,
      background: item?.maskBgColor || 'rgba(0,0,0,0.5)',
      justifyContent: maskTextAlign === 'left' ? 'flex-start' : maskTextAlign === 'right' ? 'flex-end' : 'center',
      paddingLeft: maskTextAlign === 'left' ? '24rpx' : '0',
      paddingRight: maskTextAlign === 'right' ? '24rpx' : '0'
    },
    _maskTextStyle: {
      color: item?.maskTextColor || '#FFFFFF',
      fontSize: `${item?.maskTextSize || 14}px`
    }
  }
}))

// 处理轮播图点击事件
const handleSwiperClick = (item) => {
  // 优先跳转链接，其次关联商品
  if (item.href && item.href.trim() !== '') {
    uni.navigateTo({
      url: item.href,
      fail: () => {}
    });
  } else if (item.goodID) {
    uni.navigateTo({
      url: '/pages/goodsDetails/goodsDetails?id=' + item.goodID
    });
  }
};

</script>

<style scoped lang="scss">
.swiper-section {
  padding: 24rpx 28rpx;
  background: #000;
}

.swiper {
  height: 340rpx;
  border-radius: 20rpx;
  overflow: hidden;
  box-shadow: 0 12rpx 48rpx rgba(0, 0, 0, 0.5), 0 0 60rpx rgba(229, 9, 20, 0.06);
}

.swiper-item {
  width: 100%;
  height: 100%;
  position: relative;
  background: #111;
}

.swiper-image {
  width: 100%;
  height: 100%;
  display: block;
  border-radius: 20rpx;
}

.swiper-mask {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  z-index: 2;
  border-radius: 0 0 20rpx 20rpx;
}

.swiper-mask-text {
  font-weight: 600;
  letter-spacing: 2rpx;
  text-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.5);
}
</style>
