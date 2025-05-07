<template>
    <view class="custom-divider">
      <view class="divider-line"></view>
      <text class="divider-text" v-if="text">{{ text }}</text>
      <view class="divider-line"></view>
    </view>
  </template>
  
  <script setup>
  /**
   * 自定义分割线组件
   * @prop {String} text - 分割线中间显示的文字
   * @prop {String} textColor - 文字颜色，默认 #999
   * @prop {String} lineColor - 线条颜色，默认 #EBEDF0
   * @prop {Number} lineHeight - 线条高度，默认 1px
   */
  
  // 定义组件属性
  const props = defineProps({
    text: {
      type: String,
      default: ''
    },
    textColor: {
      type: String,
      default: '#999'
    },
    lineColor: {
      type: String,
      default: '#EBEDF0'
    },
    lineHeight: {
      type: Number,
      default: 1
    }
  })
  </script>
  
  <style scoped>
  .custom-divider {
    display: flex;
    align-items: center;
    padding: 0 6rpx; /* 两侧6rpx留白 */
    width: 100%;
  }
  
  .divider-line {
    flex: 1;
    height: v-bind('props.lineHeight + "px"');
    background-color: v-bind('props.lineColor');
  }
  
  .divider-text {
    padding: 0 20rpx; /* 文字两侧留白 */
    font-size: 28rpx;
    color: v-bind('props.textColor');
    white-space: nowrap;
  }
  </style>