<template>
  <view class="comment-section section-card">
    <view class="comment-header">
      <text class="section-title">评价({{ commentData.commentCount }})</text>
      <view class="comment-rate">
        <text>好评率 {{ commentData.goodRate }}%</text>
        <uni-icons type="right" size="16" color="#999"></uni-icons>
      </view>
    </view>

    <view class="comment-item" v-if="commentData.comments.length > 0">
      <view class="comment-user">
        <view class="user-avatar" :style="{ backgroundColor: '#ddd' }"></view>
        <text class="user-name">{{ commentData.comments[0].userName }}</text>
      </view>
      <view class="comment-content">
        <text class="comment-text">{{ commentData.comments[0].content }}</text>
        <view class="comment-images">
          <view class="comment-image" :style="{ backgroundColor: '#f0f0f0' }"></view>
        </view>
        <view class="comment-info">
          <text class="comment-spec">{{ commentData.comments[0].spec }}</text>
          <text class="comment-time">{{ commentData.comments[0].time }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
const props = defineProps({
  commentData: {
    type: Object, // 修改为Object类型，因为传入的是一个对象
    default: () => ({}) // 默认值应该是空对象，而不是空数组
  }
})
</script>


<style scoped lang="scss">
/* 评论区域 */
.comment-section {
  flex-direction: column;
  align-items: flex-start;
}

.comment-header {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.comment-rate {
  display: flex;
  align-items: center;
  font-size: 12px;
  color: #666;
}

.comment-item {
  width: 100%;
}

.comment-user {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.user-avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  margin-right: 8px;
}

.user-name {
  font-size: 14px;
  color: #333;
}

.comment-content {
  padding-left: 38px;
}

.comment-text {
  font-size: 14px;
  color: #333;
  line-height: 1.4;
  margin-bottom: 10px;
}

.comment-images {
  display: flex;
  margin-bottom: 10px;
}

.comment-image {
  width: 80px;
  height: 80px;
  margin-right: 8px;
  border-radius: 4px;
}

.comment-info {
  display: flex;
  font-size: 12px;
  color: #999;
}

.comment-spec {
  margin-right: 15px;
}
</style>
