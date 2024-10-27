<template>
	<view class="flow_box">
    <scroll-view scroll-y="true" class="scroll-Y"
                 @scrolltolower="debouncedLower">
      <view class="m_t_24 bgc_fff address_item_box" v-for="(items, index) in addressList" :key="index">
        <view class="b_b_2 flex flex-aic p_b_24" @tap="selectAddr(items)">
          <view class="flex-fitem">
            <view class="m_b_4 color_333 font_28">{{items.name}} {{items.phone}}</view>
            <view class="color_999 font_24 text_nowrap" style="max-width: 580rpx;">{{items.provinceTrans}}{{items.cityTrans}} {{items.areaTrans}} {{items.street}}</view>
          </view>
          <text v-if="isShow" class="font_24 color_fe5572">使用</text>
        </view>
        <view class="p_t_24 flex-aic flexr-jsb">
          <view class="font_24 color_fe5572">
            <text @tap="editAddress(items)" class="m_r_32">编辑</text>
            <text @tap="delAddress(items)">删除</text>
          </view>
        </view>
      </view>
      <up-divider text="分割线" :dot="true" v-if="isBottom"></up-divider>
    </scroll-view>
		<view class="address_nav_view">
			<view class="address_nav_box pos_f bgc_fff flex-aic flexr-jsc">
        <button class="addBtn" @tap="toAddress">新增地址</button>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref } from "vue"
import {  updateOrder } from '@/api/order.js'
import { getAddressList, getAddressDataSource, deleteAddress } from '@/api/address.js'
	import { onLoad } from '@dcloudio/uni-app'

  const addressList = ref([])
  const isShow = ref(false)
  const orderID = ref('')
onLoad(async (options) => {
    if(options.ID) {
      orderID.value = options.ID
      isShow.value = true
    } else {
      isShow.value = false
    }
    getAddress({
      page: 1,
      pageSize: 10,
    })
    getAddressDataSources()
  });

  const addressSource = ref([])
  const getAddressDataSources = async () => {
    try {
      const res = await getAddressDataSource();
      if (res.code === 0) {
        addressSource.value = res.data;
      } else {
        uni.showToast({ title: res.msg, icon: "none" });
      }
    } catch (error) {
      console.error("获取地址数据源失败", error);
    }
  };

  const formatt = async (value, type) => {
    const source = addressSource.value
    let province = source['province']
    let city = source['city']
    let area = source['area']

    if (type === 'province') {
      // 根据 id 查找
      let result = province.find(item => Number(item.value) === Number(value));
      return result ? result : null;
    }

    if (type === 'city') {
      // 根据 id 查找
      let result = city.find(item => Number(item.value) === Number(value));
      return result ? result : null;
    }

    if (type === 'area') {
      // 根据 id 查找
      let result = area.find(item => Number(item.value) === Number(value));
      return result ? result : null;
    }
  };

  const isBottom = ref(false)
  // 初始化加载和滚动加载都会用到这个方法
  const getAddress = async (params) => {
    await getAddressDataSources(); // 明确等待数据源获取完成
    const res = await getAddressList(params);
    if(res.code === 0) {
      // 如果是滚动加载且滑到尽头
      if(res.data.list.length === 0) {
        isBottom.value = true
        uni.showToast({
          title: '没有更多数据了',
          icon: 'none'
        })
        return
      } else {
        addressList.value.push(...res.data.list)
        isBottom.value = false
      }
      for (const item of addressList.value) {
        const province = await formatt(item.province, 'province');
        item.provinceTrans = province?.label || '';

        const city = await formatt(item.city, 'city');
        item.cityTrans = city?.label || '';

        const area = await formatt(item.area, 'area');
        item.areaTrans = area?.label || '';
      }
    } else {
      uni.showToast({ title: res.msg, icon: "none" });
    }

  };


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
const lower = async () => {
  if(isBottom.value) {
    uni.showToast({
      title: '没有更多数据了',
      icon: 'none'
    })
    return
  } else {
    // 滑动到底了，然后每次给page+1 调接口继续加载下一页 如果接口已经没有数据了，给出提示并且不允许再次加载
    params.page += 1
    await getAddress(params)
  }
}

// 防抖包装的 lower 方法
const debouncedLower = debounce(lower, 300);

  const delAddress = (item) => {
    uni.showModal({
      title: '收货地址',
      content: '确定删除收货地址吗？',
      success: async function  (res) {
        if (res.confirm) {
          const del = await deleteAddress(item.ID)
          if(del.code === 0){
            uni.showToast({
              title: "删除成功",
              icon: "none"
            })
            getAddress()
          }
        } else if (res.cancel) {
          console.log('用户点击取消');
        }
      }
    });
  }

  const editAddress = (item) => {
    // 携带当前地址ID或者其他参数跳转到编辑页面反填
    uni.redirectTo({
      url:`/pages/address/editAddress?ID=${item.ID}&orderID=${orderID.value}`
    })
  }

const toAddress = () => {
    // 点击跳转到新增地址页面并携带订单编号
  uni.redirectTo({
    url: `/pages/address/addAddress?ID=${orderID.value}`
  })
}

const selectAddr = async (item) => {
  // 如果从订单页面进入，选择地址后，更新订单信息并带回
  const req = {
    ID: Number(orderID.value),
    userID: item.userID,
    name: item.name,
    phone: item.phone,
    province: item.provinceStr,
    city: item.cityStr,
    area: item.areaStr,
    Street: item.street,
    active: item.active
  }
  if(isShow) {
    const res = await updateOrder(req)
    if(res.code === 0){
      // 返回订单详情页并刷新，并携带orderID
      uni.redirectTo({
        url: `/pages/orderInfo/orderInfo?orderID=${orderID.value}`
      })

    }
  }
}

</script>

<style lang="scss">
	page {
		background-color: #f8f8f8;
	}
  .scroll-Y {
    height: 100vh;
  }
	.address_nav_view {
		height: 100rpx;
		width: 100%;
		padding-bottom: constant(safe-area-inset-bottom);
		padding-bottom: env(safe-area-inset-bottom);
	}

	.address_nav_box {
		bottom: 0;
		left: 0;
		right: 0;
		height: 100rpx;
		width: 100%;
		padding-bottom: constant(safe-area-inset-bottom);
		padding-bottom: env(safe-area-inset-bottom);
	}
  .address_item_box {
    padding: 24rpx 32rpx;
  }

  .address_item_location {
    width: 40rpx;
    height: 40rpx;
  }

  .address_item_check {
    width: 32rpx;
    height: 32rpx;
    border: 2rpx solid #DDDDDD;
    border-radius: 50%;
  }

  .address_item_check_active {
    border: none;
    width: 36rpx;
    height: 36rpx;
    background-image: url(http://www.liwanying.top/applate-icon/xuanzhong.png);
    background-repeat: no-repeat;
    background-size: 100%;
  }

  .addBtn {
    width: 90%;
    color: #fff;
    border-color: #fe5572;
    background-color: #fe5572;
  }
</style>
