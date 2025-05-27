<template>
  <view class="category-section">
    <view class="category-item" v-for="(item, index) in categoriesData" :key="index">
      <view class="" @tap="goto(item)">
        <image class="category-icon" :src="getUrl(item.icons)" mode="aspectFill"></image>
      </view>
      <text class="category-title">{{ item.title }}</text>
    </view>
  </view>

  <!-- 促销横幅 -->
  <view class="promo-banner">
    <view
        class="promo-left"
        :style="promotionData && promotionData.promotionImage ?
        `background-image: url(${getUrl(promotionData.promotionImage)})` :
        'background: linear-gradient(90deg, #FF8C69, #FFD700)'"
    >
      <view class="promo-content">
        <text class="promo-text">{{ promotionData?.title || '暂无促销活动' }}</text>
        <text class="promo-subtitle">{{ promotionData?.description || '敬请期待' }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import {getUrl} from "@/utils/url";
import { getPromotionPublic } from "@/api/homePage";

// 响应式数据
const promotionData = ref(null)

const init = async () => {
  try {
    const res = await getPromotionPublic();
    if (res.code === 0 && res.data) {
      promotionData.value = res.data;
      console.log('促销数据:', res.data);
    } else {
      console.log('暂无促销数据');
    }
  } catch (error) {
    console.error('获取促销信息失败:', error);
  }
};
init();

const props = defineProps({
  categoriesData: {
    type: Array,
    default: () => [] // 修改为函数返回空数组
  }
})

const goto = (item) => {
  if (item.ID) {
    uni.navigateTo({
      url: '/pages/tabBar/components/category-page?id=' + item.ID
    });
  }
}
</script>

<style scoped lang="scss">
/* 分类导航 - 横向滚动版本 */
.category-section {
  display: flex;
  padding: 36rpx 0;
  background-color: #fff;
  overflow-x: auto;
  overflow-y: hidden;
  white-space: nowrap;
  /* 隐藏滚动条但保持滚动功能 */
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE 10+ */
}

/* 隐藏 Webkit 内核浏览器的滚动条 */
.category-section::-webkit-scrollbar {
  display: none;
}

.category-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex-shrink: 0; /* 防止项目被压缩 */
  min-width: 120rpx; /* 设置最小宽度确保内容完整显示 */
  margin: 0 24rpx; /* 左右间距 */
}

/* 第一个和最后一个项目的特殊间距处理 */
.category-item:first-child {
  margin-left: 36rpx;
}

.category-item:last-child {
  margin-right: 36rpx;
}

.category-icon {
  width: 88rpx;
  height: 88rpx;
  border-radius: 50%;
  margin-bottom: 16rpx;
}

.category-title {
  font-size: 24rpx;
  color: #333;
  text-align: center;
  white-space: nowrap; /* 防止文字换行 */
  overflow: hidden;
  text-overflow: ellipsis; /* 文字过长时显示省略号 */
  max-width: 120rpx; /* 限制文字最大宽度 */
}

// 胶囊形状的促销横幅
.promo-banner {
  padding: 24rpx;
  background-color: #fff;
}

.promo-left {
  border-radius: 80rpx;
  height: 140rpx;
  padding: 0 40rpx;
  position: relative;
  display: flex;
  align-items: center;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  overflow: hidden;

  // 添加一个半透明遮罩，确保文字清晰可见
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.3);
    border-radius: 80rpx;
    z-index: 1;
  }
}

.promo-content {
  display: flex;
  flex-direction: column;
  position: relative;
  z-index: 2;
}

.promo-text {
  font-size: 36rpx;
  font-weight: bold;
  color: white;
  text-shadow: 0 2rpx 4rpx rgba(0, 0, 0, 0.5);
}

.promo-subtitle {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.9);
  margin-top: 8rpx;
  text-shadow: 0 1rpx 2rpx rgba(0, 0, 0, 0.3);
}

// 删除不需要的图片容器样式
</style>
