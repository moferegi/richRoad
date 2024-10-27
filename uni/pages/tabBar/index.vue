<template>
	<view class="content">
		<scroll-view scroll-y="true" class="scroll-Y" @scrolltolower="debouncedLower">
		<home-swiper :list="list"></home-swiper>
		<up-tabs
		    :list="arr"
		    lineWidth="30"
		    lineColor="#f56c6c"
		    @click="changeTabs"
		    :activeStyle="{
		        color: '#303133',
		        fontWeight: 'bold',
		        transform: 'scale(1.05)'
		    }"
		    :inactiveStyle="{
		        color: '#606266',
		        transform: 'scale(1)'
		    }"
		    itemStyle="padding-left: 15px; padding-right: 15px; height: 34px;"
		>
		</up-tabs>
		<home-flow :flowData="flowData"></home-flow>
		<view class="tabs_box">
			<up-divider :text="isBottom?'没有更多了':'正在加载中...'"></up-divider>
		</view>
		</scroll-view>
	</view>
</template>

<script setup>
	import { ref } from 'vue';
	import homeSwiper from './components/home-swiper.vue'
	import homeFlow from './components/home-flow.vue'
	import { getCategoryMobile, getGoodList } from '/api/homePage.js'
	import { getBannerList } from '/api/homePage.js'

// 轮播图相关业务逻辑
	const list = ref([])
	const initBanner = async () => {
		const res = await getBannerList()
		list.value = res.data.list
	}
	initBanner()
	
	// 商品相关属性
	const params = ref({
	  page: 1,
	  pageSize: 10,
	  categoryID: 0
	})
	
	const flowData = ref([])
	const isBottom = ref(false)
	
	
	// 获取商品相关业务逻辑
	const lower = async () => {
	  if(isBottom.value) {
	    return
	  } else {
	    // 滑动到底了，然后每次给page+1 调接口继续加载下一页 如果接口已经没有数据了，给出提示并且不允许再次加载
	    params.value.page += 1
	    const res = await getGoodList(params.value)
	    if (res.code === 0 && res.data.list.length) {
	      flowData.value.push(...res.data.list)
	      isBottom.value = false
	    } else {
	      isBottom.value = true
	    }
	  }
	}
	
	
	// 切换tabs
	const changeTabs = async (index) => {
	 params.value.categoryID = index.id
	 params.value.page = 0
	 isBottom.value = false
	 flowData.value = []
	 // 拿到index.id作为categoryID去调用/good/getGoodList接口
	 lower()
	}
	
	// 分类tabs相关业务逻辑
	const gridList = ref([])
	const arr = ref([])
	const initCategory = async () => {
	  const res = await getCategoryMobile()
	  if (res.code === 0 && res.data.length) {
	    gridList.value = res.data
	    // 赋值给新数组，新数组清洗数据改变desc为name{ name: '电影' },{ name: '科技' }的格式
	     arr.value = res.data.map(item => {
	      return {
	        name: item.title,
	        id: item.ID
	      }
	    })
	
	    arr.value.unshift({name: '全部', id: 0})
	  }
	}
	initCategory()
	changeTabs({id:0})
	
	
	// 防抖函数
	const debounce = (func, delay) => {
	  let debounceTimer;
	  return function(...args) {
	    if (debounceTimer) clearTimeout(debounceTimer);
	    debounceTimer = setTimeout(() => {
	      func.apply(this, args);
	    }, delay);
	  };
	};
	
	
	// 防抖包装的 lower 方法
	const debouncedLower = debounce(lower, 300);
	

</script>

<style lang="less">
	.content {
		width: 100%;
		.tabs_box {
			margin-top: 20rpx;
		}
	}
	.scroll-Y {
	  height: 100vh;
	}
</style>
