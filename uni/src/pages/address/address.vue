<template>
  <view class="nf-address">
    <view class="nf-address-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff" />
        </view>
        <text class="nf-navbar-title">{{ $t('myAddresses') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <scroll-view scroll-y class="nf-address-scroll" @scrolltolower="debouncedLower">
      <view class="nf-address-list" v-if="addressList.length > 0">
        <view class="nf-addr-card" v-for="(item, index) in addressList" :key="index" @tap="selectAddr(item)">
          <view class="nf-addr-info">
            <view class="nf-addr-header">
              <view class="nf-addr-name-phone">
                <text class="nf-addr-name">{{ item.name }}</text>
                <text class="nf-addr-phone">{{ item.areaCode && item.areaCode !== '+86' ? item.areaCode + ' ' : '' }}{{ item.phone }}</text>
              </view>
              <view class="nf-addr-default-tag" v-if="item.active">
                <text>{{ $t('defaultAddr') }}</text>
              </view>
            </view>
            <text class="nf-addr-detail">{{ item.provinceTrans }}{{ item.cityTrans }}{{ item.areaTrans }}{{ item.street }}</text>
          </view>
          <view class="nf-addr-actions">
            <view class="nf-addr-action" @tap.stop="editAddress(item)">
              <uni-icons type="compose" size="16" color="rgba(255,255,255,0.6)" />
              <text>{{ $t('editAddr') }}</text>
            </view>
            <view class="nf-addr-action nf-addr-action-del" @tap.stop="delAddress(item)">
              <uni-icons type="trash" size="16" color="#e50914" />
              <text>{{ $t('deleteAddr') }}</text>
            </view>
          </view>
        </view>

        <view class="nf-addr-bottom" v-if="isBottom">
          <text>{{ $t('reachedBottom') }}</text>
        </view>
      </view>

      <!-- 空状态 -->
      <view class="nf-addr-empty" v-if="addressList.length === 0">
        <uni-icons type="location" size="48" color="rgba(229,9,20,0.4)" />
        <text class="nf-addr-empty-text">{{ $t('noAddress') || '暂无收货地址' }}</text>
      </view>
    </scroll-view>

    <!-- 底部添加按钮 -->
    <view class="nf-addr-add-wrap">
      <view class="nf-addr-add-btn" @tap="toAddress">
        <uni-icons type="plusempty" size="18" color="#fff" />
        <text>{{ $t('addNewAddr') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { updateOrder } from '@/api/order.js'
import { getAddressList, getAddressDataSource, deleteAddress } from '@/api/address.js'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { localText } from '@/utils/i18n'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const addressList = ref([])
const isShow = ref(false)
const orderID = ref('')
const isBottom = ref(false)
const addressSource = ref([])
const fromPage = ref('')

onLoad(async (options) => {
  if (options.ID) {
    orderID.value = options.ID
    isShow.value = true
  } else {
    isShow.value = false
  }
  if (options.from) {
    fromPage.value = options.from
  }
  addressList.value = []
  await getAddressDataSources()
  getAddress({ page: 1, pageSize: 10 })
})

let loaded = false
onShow(() => {
  if (loaded) {
    addressList.value = []
    getAddress({ page: 1, pageSize: 10 })
  }
  loaded = true
})

const getAddressDataSources = async () => {
  try {
    const res = await getAddressDataSource()
    if (res.code === 0) {
      addressSource.value = res.data
    }
  } catch (error) {
    console.error('获取地址数据源失败', error)
  }
}

const formatt = (value, type) => {
  const source = addressSource.value
  const list = source[type]
  if (!list) return null
  return list.find(item => Number(item.value) === Number(value)) || null
}

const getGeoLabel = (item) => {
  if (!item) return ''
  // 优先使用多语言字段
  if (item.labelI18n) {
    const translated = localText(item.labelI18n)
    if (translated) return translated
  }
  return item.label || ''
}

const getAddress = async (params) => {
  try {
    const res = await getAddressList(params)
    if (res.code === 0) {
      if (res.data.list.length === 0) {
        isBottom.value = true
        if (params.page === 1) addressList.value = []
        return
      }
      addressList.value.push(...res.data.list)
      isBottom.value = false

      for (const item of addressList.value) {
        const province = formatt(item.province, 'province')
        item.provinceTrans = getGeoLabel(province)

        const city = formatt(item.city, 'city')
        item.cityTrans = getGeoLabel(city)

        const area = formatt(item.area, 'area')
        item.areaTrans = getGeoLabel(area)
      }
    }
  } catch (error) {
    console.error('获取地址列表失败', error)
  }
}

const debounce = (func, delay) => {
  let timer
  return function (...args) {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => func.apply(this, args), delay)
  }
}

let params = { page: 1, pageSize: 10 }

const lower = async () => {
  if (isBottom.value) return
  params.page += 1
  await getAddress(params)
}
const debouncedLower = debounce(lower, 300)

const delAddress = (item) => {
  uni.showModal({
    title: $t.value('confirmDeleteAddr') || '确定删除收货地址吗？',
    success: async (res) => {
      if (res.confirm) {
        const del = await deleteAddress(item.ID)
        if (del.code === 0) {
          uni.showToast({ title: $t.value('deleteSuccess') || '删除成功', icon: 'none' })
          addressList.value = []
          params.page = 1
          getAddress(params)
        }
      }
    }
  })
}

const editAddress = (item) => {
  uni.navigateTo({
    url: `/pages/address/editAddress?ID=${item.ID}&orderID=${orderID.value}`
  })
}

const toAddress = () => {
  uni.navigateTo({
    url: `/pages/address/addAddress?ID=${orderID.value}`
  })
}

const selectAddr = async (item) => {
  // 从 orderInfo 页面进入（订单尚未创建），将选中地址存入 storage 后返回
  if (fromPage.value === 'orderInfo') {
    uni.setStorageSync('selectedAddress', {
      ID: item.ID,
      name: item.name,
      phone: item.phone,
      province: item.province,
      city: item.city,
      area: item.area,
      provinceStr: item.provinceTrans || item.provinceStr || '',
      cityStr: item.cityTrans || item.cityStr || '',
      areaStr: item.areaTrans || item.areaStr || '',
      street: item.street,
      active: item.active,
    })
    uni.navigateBack()
    return
  }
  // 从 orderDetail 页面进入，将选中地址存入 storage 后返回（不立即提交后端）
  if (fromPage.value === 'orderDetail') {
    uni.setStorageSync('selectedAddress', {
      ID: item.ID,
      name: item.name,
      phone: item.phone,
      province: item.province,
      city: item.city,
      area: item.area,
      provinceStr: item.provinceTrans || item.provinceStr || '',
      cityStr: item.cityTrans || item.cityStr || '',
      areaStr: item.areaTrans || item.areaStr || '',
      street: item.street,
      active: item.active,
    })
    uni.navigateBack()
    return
  }
  if (isShow.value) {
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
    if (res.code === 0) {
      // 根据来源页面返回
      if (fromPage.value === 'orderDetail') {
        uni.navigateBack()
      } else {
        uni.redirectTo({
          url: `/pages/orderInfo/orderInfo?orderID=${orderID.value}`
        })
      }
    }
  }
}

const goBack = () => {
  uni.navigateBack({ fail: () => uni.switchTab({ url: '/pages/tabBar/index' }) })
}
</script>

<style lang="scss" scoped>
.nf-address {
  min-height: 100vh;
  background: #141414;
  position: relative;
}
.nf-address-bg {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: linear-gradient(180deg, #1a1a2e 0%, #141414 100%);
  z-index: 0;
}
.nf-navbar {
  position: fixed;
  top: 0; left: 0; right: 0;
  z-index: 100;
  background: rgba(20, 20, 20, 0.95);
  backdrop-filter: blur(20rpx);
}
.nf-navbar-status {
  height: var(--status-bar-height, 44rpx);
}
.nf-navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 88rpx;
  padding: 0 24rpx;
}
.nf-navbar-back {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.06);
  border: 1rpx solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}
.nf-navbar-title {
  font-size: 34rpx;
  font-weight: bold;
  color: #fff;
}
.nf-address-scroll {
  position: relative;
  z-index: 1;
  height: calc(100vh - var(--status-bar-height, 44rpx) - 88rpx - 120rpx);
  margin-top: calc(var(--status-bar-height, 44rpx) + 88rpx);
  padding: 0 24rpx;
}
.nf-address-list {
  padding-top: 20rpx;
  padding-bottom: 30rpx;
}
.nf-addr-card {
  background: rgba(255, 255, 255, 0.06);
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.08);
}
.nf-addr-info {
  margin-bottom: 20rpx;
}
.nf-addr-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12rpx;
}
.nf-addr-name-phone {
  display: flex;
  align-items: center;
  gap: 16rpx;
}
.nf-addr-name {
  font-size: 30rpx;
  font-weight: bold;
  color: #fff;
}
.nf-addr-phone {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.5);
}
.nf-addr-default-tag {
  background: rgba(229, 9, 20, 0.2);
  border: 1rpx solid rgba(229, 9, 20, 0.4);
  border-radius: 6rpx;
  padding: 4rpx 12rpx;
}
.nf-addr-default-tag text {
  font-size: 20rpx;
  color: #e50914;
}
.nf-addr-detail {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.6);
  line-height: 1.5;
}
.nf-addr-actions {
  display: flex;
  gap: 30rpx;
  padding-top: 20rpx;
  border-top: 1rpx solid rgba(255, 255, 255, 0.08);
}
.nf-addr-action {
  display: flex;
  align-items: center;
  gap: 8rpx;
}
.nf-addr-action text {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.6);
}
.nf-addr-action-del text {
  color: #e50914;
}
.nf-addr-bottom {
  text-align: center;
  padding: 20rpx 0;
}
.nf-addr-bottom text {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.3);
}
.nf-addr-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding-top: 300rpx;
}
.nf-addr-empty-text {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.4);
  margin-top: 20rpx;
}
.nf-addr-add-wrap {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20rpx 24rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background: rgba(20, 20, 20, 0.95);
  backdrop-filter: blur(20rpx);
  z-index: 100;
}
.nf-addr-add-btn {
  height: 88rpx;
  border-radius: 12rpx;
  background: #e50914;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8rpx;
}
.nf-addr-add-btn text {
  font-size: 30rpx;
  font-weight: bold;
  color: #fff;
}
</style>
