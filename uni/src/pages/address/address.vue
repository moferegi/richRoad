<template>
  <view class="address-container">
    <scroll-view scroll-y="true" class="address-scroll" @scrolltolower="debouncedLower">
      <!-- 有地址时显示列表 -->
      <block v-if="addressList.length > 0">
        <view class="address-item" v-for="(item, index) in addressList" :key="index" @tap="selectAddr(item)">
          <view class="address-content">
            <view class="address-header">
              <view class="name-phone">
                <text class="name">{{item.name}}</text>
                <text class="phone">{{item.phone}}</text>
              </view>
              <view class="tag-box" v-if="item.active">
                <text class="default-tag">默认</text>
              </view>
            </view>
            <view class="address-detail">
              {{item.provinceTrans}}{{item.cityTrans}}{{item.areaTrans}}{{item.street}}
            </view>
          </view>

          <!-- 每个地址下面的编辑删除操作栏 -->
          <view class="address-actions" >
            <view class="action-divider"></view>
            <view class="action-buttons">
              <view class="action-btn edit-btn" @tap="editAddress(item)">
                <text>编辑</text>
              </view>
              <view class="action-divider-vertical"></view>
              <view class="action-btn delete-btn" @tap="delAddress(item)">
                <text>删除</text>
              </view>
            </view>
          </view>
        </view>

        <view class="bottom-text" v-if="isBottom">已经到底了</view>
      </block>

      <!-- 空地址状态 -->
      <view class="empty-address" v-if="addressList.length === 0">
        <image class="empty-icon" src="/static/images/empty-address.png"></image>
        <text class="empty-text">暂无收货地址，请添加</text>
      </view>
    </scroll-view>

    <!-- 底部添加按钮 -->
    <view class="add-address-btn-wrapper">
      <button class="add-address-btn" @tap="toAddress">新增地址</button>
    </view>
  </view>
</template>

<script setup>
import { ref } from "vue"
import { updateOrder } from '@/api/order.js'
import { getAddressList, getAddressDataSource, deleteAddress } from '@/api/address.js'
import { onLoad } from '@dcloudio/uni-app'

const addressList = ref([])
const isShow = ref(false)
const orderID = ref('')
const isBottom = ref(false)
const addressSource = ref([])

onLoad(async (options) => {
  if(options.ID) {
    orderID.value = options.ID
    isShow.value = true
  } else {
    isShow.value = false
  }
  // 重置列表避免重复加载
  addressList.value = []
  await getAddressDataSources()
  getAddress({
    page: 1,
    pageSize: 10,
  })
})

const getAddressDataSources = async () => {
  try {
    const res = await getAddressDataSource()
    if (res.code === 0) {
      addressSource.value = res.data
    } else {
      uni.showToast({ title: res.msg, icon: "none" })
    }
  } catch (error) {
    console.error("获取地址数据源失败", error)
  }
}

const formatt = async (value, type) => {
  const source = addressSource.value
  let province = source['province']
  let city = source['city']
  let area = source['area']

  if (type === 'province') {
    let result = province.find(item => Number(item.value) === Number(value))
    return result ? result : null
  }

  if (type === 'city') {
    let result = city.find(item => Number(item.value) === Number(value))
    return result ? result : null
  }

  if (type === 'area') {
    let result = area.find(item => Number(item.value) === Number(value))
    return result ? result : null
  }
}

const getAddress = async (params) => {
  try {
    const res = await getAddressList(params)
    if(res.code === 0) {
      // 如果是滚动加载且滑到尽头
      if(res.data.list.length === 0) {
        isBottom.value = true
        if (params.page > 1) {
          uni.showToast({
            title: '没有更多地址了',
            icon: 'none'
          })
        }
        return
      } else {
        addressList.value.push(...res.data.list)
        isBottom.value = false
      }

      // 处理地址翻译
      for (const item of addressList.value) {
        const province = await formatt(item.province, 'province')
        item.provinceTrans = province?.label || ''

        const city = await formatt(item.city, 'city')
        item.cityTrans = city?.label || ''

        const area = await formatt(item.area, 'area')
        item.areaTrans = area?.label || ''
      }
    } else {
      uni.showToast({ title: res.msg, icon: "none" })
    }
  } catch (error) {
    console.error("获取地址列表失败", error)
    uni.showToast({ title: "获取地址列表失败", icon: "none" })
  }
}

const debounce = (func, delay) => {
  let debounceTimer
  return function(...args) {
    if (debounceTimer) clearTimeout(debounceTimer)
    debounceTimer = setTimeout(() => {
      func.apply(this, args)
    }, delay)
  }
}

let params = {
  page: 1,
  pageSize: 10
}

const lower = async () => {
  if(isBottom.value) {
    return
  } else {
    params.page += 1
    await getAddress(params)
  }
}

// 防抖包装的 lower 方法
const debouncedLower = debounce(lower, 300)

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
  if(isShow.value) {
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
  background-color: #f5f5f5;
}

.address-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.address-scroll {
  flex: 1;
  padding-bottom: 100rpx;
}

.address-item {
  background-color: #fff;
  margin: 20rpx;
  border-radius: 12rpx;
  padding: 30rpx;
  display: flex;
  flex-direction: column;
}

.address-content {
  flex: 1;
}

.address-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10rpx;
}

.name-phone {
  display: flex;
  align-items: center;
}

.name {
  font-size: 32rpx;
  font-weight: 500;
  color: #333;
  margin-right: 20rpx;
}

.phone {
  font-size: 28rpx;
  color: #666;
}

.tag-box {
  display: flex;
}

.default-tag {
  font-size: 22rpx;
  background-color: #ff4c7d;
  color: #fff;
  padding: 2rpx 10rpx;
  border-radius: 4rpx;
}

.address-detail {
  font-size: 28rpx;
  color: #666;
  line-height: 1.4;
  margin-top: 10rpx;
}

.address-actions {
  background-color: #fff;
  border-radius: 0 0 12rpx 12rpx;
  margin-top: 40rpx;
}

.action-divider {
  height: 1rpx;
  background-color: #eee;
  margin: 0 0 20rpx 0;
}

.action-divider-vertical {
  width: 1rpx;
  height: 40rpx;
  background-color: #eee;
}

.action-buttons {
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding-bottom: 20rpx;
}

.action-btn {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 60rpx;
}

.action-btn text {
  font-size: 28rpx;
  color: #666;
}

.edit-btn text {
  color: #333;
}

.delete-btn text {
  color: #ff4c7d;
}

.empty-address {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding-top: 200rpx;
}

.empty-icon {
  width: 200rpx;
  height: 200rpx;
  margin-bottom: 30rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #999;
}

.add-address-btn-wrapper {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20rpx;
  background-color: #fff;
  box-shadow: 0 -2rpx 10rpx rgba(0,0,0,0.05);
}

.add-address-btn {
  background-color: #ff4c7d;
  color: #fff;
  border-radius: 8rpx;
  font-size: 32rpx;
  height: 90rpx;
  line-height: 90rpx;
}

.bottom-text {
  text-align: center;
  color: #999;
  font-size: 24rpx;
  padding: 30rpx 0;
}
</style>
