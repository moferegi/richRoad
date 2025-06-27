<template>
  <view class="goods-list" ref="goodsListRef">
    <view v-for="(item, index) in props.goodsList" :key="index" class="goods-item" @tap="handleGoodsClick(item)">
      <image :src="getUrl(item.imageUrl)" class="goods-image" mode="aspectFill"></image>
      <view class="goods-info">
        <text class="goods-name">{{ item.title }}</text>
        <view v-if="item.tags && item.tags.length > 0" class="merchant-tags">
          <text
              v-for="tag in item.tags"
              :key="tag.ID"
              :style="{
            color: tag.color,
            background: `${tag.color}1A`,
            border: `1px solid ${tag.color}33`
        }"
              class="merchant-tag"
          >
            {{ tag.name }}
          </text>
        </view>
        <view class="price-container">
          <text class="discount-price">¥{{ item.price / 100 }}</text>
          <text class="original-price">¥{{ item.originalPrice }}</text>
          <text class="discount-tag">{{ getDiscountText(item.discount) }}</text>
        </view>
        <view class="goods-extra">
          <view class="rating">
            <text class="rating-score">{{ item.rating }}</text>
            <text class="rating-stars">{{ getRatingStars(item.rating) }}</text>
            <text class="rating-count">({{ item.ratingCount || 0 }})</text>
          </view>
          <view class="sales">
            <text>售出 {{ item.saleNum }} 件</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 滚动检测触发区域 -->
    <view
      class="scroll-trigger"
      ref="scrollTriggerRef"
      v-if="!isLastPage && goodsList.length > 0"
    >
      <!-- 这个区域用于触发滚动检测，当它进入可视区域时自动加载下一页 -->
    </view>
  </view>
</template>

<script setup>
import {ref, computed, onMounted, onUnmounted, nextTick, watch} from 'vue'
import {getUrl} from "@/utils/url.js"
import { onReachBottom } from '@dcloudio/uni-app'

const goodsList = ref([])
const props = defineProps({
  goodsList: {
    type: Array,
    default: () => []
  },
  currentPage: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 10
  },
  total: {
    type: Number,
    default: 0
  },
  // 距离底部多少rpx时触发加载
  loadOffset: {
    type: Number,
    default: 400
  }
})
const getDiscountText = (discount) => {
  if (discount >= 9.5) return '小降'
  if (discount >= 9.0) return '优惠'
  if (discount >= 8.0) return '特惠'
  if (discount >= 7.0) return '好价'
  if (discount >= 6.0) return '低价'
  if (discount >= 5.0) return '特价'
  return '折扣'
}
const queryList = async (pageNo, pageSize) => {

// 加载数据方法

}
// 向父组件发送事件
const emit = defineEmits(['load-more'])

// 组件引用
const goodsListRef = ref(null)

// 计算是否为最后一页
const isLastPage = computed(() => {
  if (props.total % props.pageSize === 0) {
    return props.currentPage >= props.total / props.pageSize
  }
  return props.currentPage >= Math.floor(props.total / props.pageSize) + 1
})

// 页面滚动到底部触发加载
onReachBottom(() => {
  if (!isLastPage.value) {
    emit('load-more')
  }
})

/*const emit = defineEmits(['load-more'])
// 组件引用
const goodsListRef = ref(null)
const scrollTriggerRef = ref(null)

// 滚动检测相关
const isLoadingMore = ref(false)
const intersectionObserver = ref(null)

// 计算是否为最后一页
const isLastPage = computed(() => {
  if (props.total === 0) return false
  return props.currentPage * props.pageSize >= props.total
})*/
// 根据评分生成星星
const getRatingStars = (rating) => {
  // 如果评分为0，返回1颗星
  if (rating === 0) return '★☆☆☆☆';

  // 计算实心星星数量（最大5颗）
  const fullStars = Math.min(Math.floor(rating), 5);
  // 计算空心星星数量
  const emptyStars = 5 - fullStars;

  // 返回对应数量的星星
  return '★'.repeat(fullStars) + '☆'.repeat(emptyStars);
}

// 跳转到商品详情页
const handleGoodsClick = (item) => {
  uni.navigateTo({
    url: '/pages/goodsDetails/goodsDetails?id=' + item.ID
  })
}

/*// 初始化滚动检测
const initScrollDetection = () => {
  // 使用 IntersectionObserver 检测触发区域是否进入可视区域
  intersectionObserver.value = uni.createIntersectionObserver()

  intersectionObserver.value
    .relativeToViewport({ bottom: props.loadOffset })
    .observe('.scroll-trigger', (res) => {
      if (res.intersectionRatio > 0 && !isLoadingMore.value && !props.loading && !isLastPage.value) {
        handleAutoLoadMore()
      }
    })
}

// 自动加载更多
const handleAutoLoadMore = async () => {
  if (isLoadingMore.value || props.loading || isLastPage.value) {
    return
  }

  isLoadingMore.value = true

  try {
    // 发送加载更多事件给父组件
    emit('load-more')
  } finally {
    // 延迟重置状态，避免重复触发
    setTimeout(() => {
      isLoadingMore.value = false
    }, 500)
  }
}

// 监听商品列表变化，重新初始化滚动检测
watch(() => props.goodsList.length, async (newLength, oldLength) => {
  if (newLength > oldLength) {
    // 新数据加载完成后，重新初始化滚动检测
    await nextTick()
    if (intersectionObserver.value) {
      intersectionObserver.value.disconnect()
      initScrollDetection()
    }
  }
})

// 组件挂载时初始化
onMounted(() => {
  nextTick(() => {
    initScrollDetection()
  })
})

// 组件卸载时清理
onUnmounted(() => {
  if (intersectionObserver.value) {
    intersectionObserver.value.disconnect()
  }
})*/

</script>

<style lang="scss" scoped>
page {
  background: #f5f7fa;
}

.goods-list {
  padding: 12rpx;
  background: #f5f7fa;
}

.goods-item {
  display: flex;
  background: #ffffff;
  margin-bottom: 12rpx;
  border-radius: 12rpx;
  padding: 12rpx;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);

  &:active {
    transform: scale(0.98);
  }

  .goods-image {
    width: 240rpx;
    height: 240rpx;
    border-radius: 8rpx;
    margin-right: 16rpx;
  }
}

.goods-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  justify-content: space-around;

  .goods-name {
    font-size: 28rpx;
    color: #333;
    line-height: 1.4;
    margin-bottom: 8rpx;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 2;
  }
}

.merchant-tags {
  display: flex;
  flex-wrap: wrap;
  margin: 4rpx 0;

  .merchant-tag {
    font-size: 18rpx;
    padding: 0 6rpx;
    border-radius: 4rpx;
    height: 26rpx;
    line-height: 26rpx;
    margin-right: 4rpx;
    margin-bottom: 4rpx;

    &.self-operated {
      color: #ff6b6b;
      background: rgba(255, 107, 107, 0.1);
      border: 1px solid rgba(255, 107, 107, 0.2);
    }

    &.quality-assured {
      color: #2196f3;
      background: rgba(33, 150, 243, 0.1);
      border: 1px solid rgba(33, 150, 243, 0.2);
    }

    &.plus-delivery {
      color: #4caf50;
      background: rgba(76, 175, 80, 0.1);
      border: 1px solid rgba(76, 175, 80, 0.2);
    }
  }
}

.price-container {
  display: flex;
  align-items: center;
  margin: 8rpx 0;

  .discount-price {
    font-size: 32rpx;
    color: #ff4444;
    font-weight: bold;
    margin-right: 8rpx;
  }

  .original-price {
    font-size: 22rpx;
    color: #999;
    text-decoration: line-through;
    margin-right: 8rpx;
  }

  .discount-tag {
    font-size: 20rpx;
    color: #fff;
    background: #ff4444;
    padding: 2rpx 8rpx;
    border-radius: 4rpx;
  }
}

.goods-extra {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 20rpx;
  color: #666;
  margin-top: 8px;

  .rating {
    display: flex;
    align-items: center;

    .rating-score {
      color: #ff4444;
      font-weight: bold;
      margin-right: 2rpx;
    }

    .rating-stars {
      color: #ffd700;
      font-size: 18rpx;
      margin-right: 2rpx;
    }

    .rating-count {
      color: #999;
    }
  }

  .sales {
    color: #999;
  }
}

.goods-tags {
  display: flex;
  flex-wrap: wrap;
  margin: 4rpx 0;

  .tag {
    font-size: 18rpx;
    color: #666;
    background: #f7f7f7;
    padding: 0 6rpx;
    border-radius: 2rpx;
    margin-right: 4rpx;
    margin-bottom: 4rpx;
  }
}

.scroll-trigger {
  height: 1rpx;
  width: 100%;
  // 这个区域用于触发滚动检测，不需要可见
}

.loading-indicator {
  padding: 30rpx;
  text-align: center;
  color: #666;
  font-size: 28rpx;
}

.no-more-data {
  padding: 30rpx;
  text-align: center;
  color: #999;
  font-size: 24rpx;
}
</style>
