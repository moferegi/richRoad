<template>
  <view class="cate-content">
    <view class="cate-box">
      <view class="category-section">
        <view
            class="category-item"
            :class="{ 'active': selectedIndex === index }"
            v-for="(item, index) in gridList"
            :key="index"
            @tap="goto(item, index)"
        >
          <view class="">
            <image class="category-icon" :src="getUrl(item.icons)" mode="aspectFill"></image>
          </view>
          <text class="category-title">{{ item.title }}</text>
        </view>
      </view>
    </view>

    <!--  商品展示-->
    <view>
      <noPaginGridGoodList
          :goodsList="flowData"
          @loadMore="loadMoreGoods"
          :loading="loading"
          :noMore="noMore"
      ></noPaginGridGoodList>
    </view>
  </view>
</template>

<script setup>
import {ref, nextTick} from 'vue'
import noPaginGridGoodList from '@/components/good-list/no-pagin-grid-good-list.vue'
import { getCategoryMobile, getGoodList } from '@/api/homePage.js'
import {onLoad} from "@dcloudio/uni-app"
import {getUrl} from "@/utils/url"

const selectedIndex = ref(0)
const optionID = ref('')
const currentCategoryID = ref('') // 当前选中的分类ID

// 分页相关状态
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const noMore = ref(false)

onLoad((options) => {
  if (options.id) {
    optionID.value = options.id
    currentCategoryID.value = options.id
  }
})

const gridList = ref([])
const flowData = ref([])

const initCategory = async () => {
  const res = await getCategoryMobile()
  if (res.code === 0 && res.data.length) {
    gridList.value = res.data
  }

  // 匹配选中项的颜色
  const targetIndex = gridList.value.findIndex(item => item.ID == optionID.value)
  if (targetIndex !== -1) {
    selectedIndex.value = targetIndex
  }

  // 如果有传入的分类ID，则加载对应的商品数据
  if (currentCategoryID.value) {
    await loadGoodsList(currentCategoryID.value, true)
  }
}


// 加载商品列表
const loadGoodsList = async (categoryID, isReset = false) => {
  if (loading.value) return

  loading.value = true

  try {
    // 如果是重置（切换分类），则重置分页状态
    if (isReset) {
      currentPage.value = 1
      flowData.value = []
      noMore.value = false
    }

    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      categoryID: categoryID
    }

    const res = await getGoodList(params)
    if (res.code === 0) {
      // 获取商品列表数据
      const newGoodsList = res.data?.list || []

      if (isReset || currentPage.value === 1) {
        // 初始化或重置时直接赋值
        flowData.value = newGoodsList
      } else {
        // 翻页时追加数据
        flowData.value = [...flowData.value, ...newGoodsList]
      }

      // 判断是否还有更多数据
      if (newGoodsList.length < pageSize.value) {
        noMore.value = true
      } else {
        currentPage.value += 1
      }
    } else {
      // API返回非0状态码时的处理
      uni.showToast({
        title: res.msg || '加载失败，请重试',
        icon: 'none'
      })
    }
  } catch (error) {
    console.error('加载商品列表失败:', error)
    uni.showToast({
      title: '加载失败，请重试',
      icon: 'none'
    })
  } finally {
    loading.value = false
  }
}


// 加载更多商品（由子组件触发）
const loadMoreGoods = () => {
  if (!loading.value && !noMore.value && currentCategoryID.value) {
    loadGoodsList(currentCategoryID.value, false)
  }
}

// 点击分类切换
const goto = async (item, index) => {
  selectedIndex.value = index
  currentCategoryID.value = item.ID

  // 切换分类时重新加载商品数据
  await loadGoodsList(item.ID, true)
}

// 初始化
initCategory()
</script>

<style scoped lang="scss">
.cate-box{
  height: 200rpx;
  width: 100%;
  background:  rgba(254, 169, 74, 0.9)
}

/* 分类导航 - 横向滚动版本 */
.category-section {
  display: flex;
  padding: 36rpx 0;
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
/* 选中状态样式 */
.category-item.active {
  transform: translateY(-4rpx); /* 轻微上移效果 */
}

.category-item.active .category-icon {
  border: 4rpx solid #ffffff; /* 选中时图标边框 */
  box-shadow: 0 4rpx 12rpx rgba(0, 122, 255, 0.3); /* 添加阴影 */
}

.category-item.active .category-title {
  color: #ffffff; /* 选中时文字颜色 */
  font-weight: bold; /* 选中时文字加粗 */
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
</style>
