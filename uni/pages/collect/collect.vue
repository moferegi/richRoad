<template>
	<view v-if="collectList.length">
    <scroll-view scroll-y="true" class="scroll-Y"
                 @scrolltolower="debouncedLower">
		<uni-swipe-action autoClose>
			<uni-swipe-action-item @click="cancelCollect(item.ID)" v-for="(item, index) in collectList" :right-options="options" :key="index">
          <view class="collect_goods flex" @tap="goTo(item)">
            <image class="collect_goods_img m_r_24" :src="getUrl(item.imageUrl)" mode="aspectFill"></image>
            <view class="flex-fitem flexc-jsb">
              <view class="color_333 font_30 m_b_4 text_nowrap_2">{{ item.title }}</view>
              <view class="flex flex-aife color_ff0003">
                <text class="font_28">¥</text>
                <text class="font_40">{{item.price/100}}</text>
              </view>
              <view class="desc">
                <p class="color_b5b5b5 font_24">已售：{{item.saleNum}}</p>
              </view>
            </view>
          </view>
        <up-divider></up-divider>
			</uni-swipe-action-item>
		</uni-swipe-action>
    </scroll-view>
	</view>
  <view v-if="!collectList.length" class="non-collect flexr-jsc">
    <up-empty
        mode="favor"
        icon="http://cdn.uviewui.com/uview/empty/car.png"
    >
    </up-empty>
  </view>
</template>

<script setup>
	import { ref } from 'vue'
  import { onShow } from '@dcloudio/uni-app'
  import { getCollectList } from '@/api/collect'
  import {useUserStore} from "@/pinia/modules/user";
  import {findCollect, createCollect} from '@/api/collect.js'
  	import {getUrl} from "@/utils/url.js"
	const options = ref([{
		text: '删除',
		style: {
			backgroundColor: '#F56C6C'
		}
	}])

  const userStore = useUserStore()
  const token = userStore.token || ''
  const collectList = ref([])
  const collectionFlag = ref('')
  const init = async () => {
    const pageInfo = {
      page: 1,
      pageSize: 10
    }
    const res = await getCollectList(pageInfo)
    if(res.code === 0) {
      collectList.value = res.data.list
    }
  }

  onShow(() => {
    init()
  })
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

  let params = {
    page: 1,
    pageSize: 10
  }
  const isBottom = ref(false)
  const lower = async (e) => {
    if(isBottom.value) {
      uni.showToast({
        title: '没有更多数据了',
        icon: 'none'
      })
      return
    } else {
      // 滑动到底了，然后每次给page+1 调接口继续加载下一页 如果接口已经没有数据了，给出提示并且不允许再次加载
      params.page += 1
      const res = await getCollectList(params)
      if (res.code === 0 && res.data.list.length) {
        collectList.value.push(...res.data.list)
        isBottom.value = false
      } else {
        isBottom.value = true
      }
    }
  }
  // 防抖包装的 lower 方法
  const debouncedLower = debounce(lower, 300);


  const cancelCollect = async (ID) => {
    if (token) {
      // 先查看当前商品收藏状态
      const status = await findCollect({
        goodID: ID
      })
      status.code === 0 ? collectionFlag.value = status.data : ''
      // 如果已登录并且未收藏 则允许进行收藏操作
      const res = await createCollect({
        goodID: Number(ID)
      })
      if (res.code === 0) {
        collectionFlag.value = !collectionFlag.value
        init()
        uni.showToast({
          title: collectionFlag.value ? '已收藏' : '已取消收藏',
          mask: true,
          icon: 'none'
        });
      }
    } else {
      uni.showToast({
        title: '请登录后进行操作',
        mask: true,
        icon: 'none'
      });
      uni.redirectTo({
        url: '/pages/user/login'
      })
    }
  }

  const goTo = (item) => {
    uni.navigateTo({
      url: `/pages/goodsDetails/goodsDetails?id=${item.ID}`,
    })
  }


</script>

<style lang="scss" scoped>
	.collect_nav_view {
		width: 100%;
		height: 124rpx;
	}
  .scroll-Y {
    height: 100vh;
  }
	.collect_nav_box {
		width: 100%;
		height: 124rpx;
		z-index: 1;
		/* #ifdef H5 */
		top: var(--window-top);
		/* #endif */
		/* #ifndef H5 */
		top: 0;
		/* #endif */
		left: 0;
		right: 0;
	}

	.collect_search_box {
		width: 100%;
		padding: 24rpx 32rpx;
	}

	.search_icon {
		top: 42rpx;
		left: 64rpx;
		width: 40rpx;
		height: 40rpx;
	}

	.search_input {
		width: 100%;
		display: block;
		height: 76rpx;
		border-radius: 76rpx;
		padding-left: 88rpx;
	}
  .collect_goods {
    padding: 24rpx 32rpx 0;
  }
  .collect_goods_img {
    width: 200rpx;
    height: 172rpx;
    border-radius: 12rpx;
    overflow: hidden;
  }
	.non-collect{
    height: 100vh;
  }

  .desc{
    width: 100%;
    margin-top: 20rpx;
  }
</style>
