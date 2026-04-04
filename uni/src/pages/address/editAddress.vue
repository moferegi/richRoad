<template>
  <view class="nf-edit-addr">
    <view class="nf-edit-addr-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff" />
        </view>
        <text class="nf-navbar-title">{{ $t('editAddr') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <view class="nf-form-content">
      <!-- 收件人 -->
      <view class="nf-form-item">
        <text class="nf-form-label">{{ $t('recipientName') }}</text>
        <input class="nf-form-input" v-model="formData.name" :placeholder="$t('nameRequired')" />
      </view>

      <!-- 手机号 -->
      <view class="nf-form-item">
        <text class="nf-form-label">{{ $t('recipientPhone') }}</text>
        <input class="nf-form-input" v-model="formData.phone" :placeholder="$t('phoneRequired')" :maxlength="20" />
      </view>

      <!-- 地区 -->
      <view class="nf-form-item">
        <text class="nf-form-label">{{ $t('selectRegion') }}</text>
        <uni-data-select class="nf-select" v-model="formData.provinceSelect" :placeholder="$t('regionRequired')"
          :localdata="areaProvince" @change="changeProvince" :clear="false" />
      </view>

      <!-- 城市 -->
      <view class="nf-form-item" v-if="isCity">
        <text class="nf-form-label">{{ $t('selectCity') }}</text>
        <uni-data-select class="nf-select" v-model="formData.citySelect" :placeholder="$t('cityRequired')"
          :localdata="areaCity" @change="changeCity" :clear="false" />
      </view>

      <!-- 区县 -->
      <view class="nf-form-item" v-if="isCounty">
        <text class="nf-form-label">{{ $t('selectDistrict') }}</text>
        <uni-data-select class="nf-select" v-model="formData.countySelect" :placeholder="$t('districtRequired')"
          :localdata="areaCounty" @change="changeCounty" :clear="false" />
      </view>

      <!-- 详细地址 -->
      <view class="nf-form-item">
        <text class="nf-form-label">{{ $t('detailAddress') }}</text>
        <input class="nf-form-input" v-model="formData.street" :placeholder="$t('streetRequired')" />
      </view>

      <!-- 默认地址 -->
      <view class="nf-form-item nf-form-item-switch">
        <text class="nf-form-label">{{ $t('setDefault') }}</text>
        <switch @change="switchChange" :checked="formData.active" color="#e50914" style="transform:scale(0.8)" />
      </view>

      <!-- 保存按钮 -->
      <view class="nf-form-save" @tap="confirm">
        <text>{{ $t('saveAddr') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getGeos, updateAddress, findAddress } from '@/api/address.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { localText } from '@/utils/i18n'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const formData = ref({})
const areaProvince = ref([])
const areaCity = ref([])
const areaCounty = ref([])
const isCity = ref(false)
const isCounty = ref(false)
const orderID = ref('')
const addrID = ref('')

const changeKey = (data) => {
  return data.map(item => ({
    ...item,
    value: item.code,
    text: localText(item.nameI18n) || item.name
  }))
}

const findCodeByValue = (data, selectedValue) => {
  return data.find(item => item.value === selectedValue) || null
}

onLoad(async (e) => {
  orderID.value = e.orderID
  const province = await getGeos({ level: 2, code: 0 })
  areaProvince.value = changeKey(province.data)

  const res = await findAddress({ id: e.ID })
  if (res.code === 0) {
    isCity.value = true
    isCounty.value = true
    const info = res.data.readdress
    formData.value = info
    formData.value.provinceSelect = String(info.province)
    formData.value.citySelect = String(info.city)
    formData.value.countySelect = String(info.area)
    addrID.value = info.ID
  }
})

watch(() => formData.value.provinceSelect, (newVal) => {
  if (newVal) loadCities(newVal)
})

watch(() => formData.value.citySelect, (newVal) => {
  if (newVal) loadCounties(newVal)
})

const loadCities = async (provinceId) => {
  const city = await getGeos({ level: 0, code: provinceId })
  areaCity.value = changeKey(city.data)
}

const loadCounties = async (cityId) => {
  const counties = await getGeos({ level: 1, code: cityId })
  areaCounty.value = changeKey(counties.data)
}

const changeProvince = async (e) => {
  const selectedItems = findCodeByValue(areaProvince.value, e)
  if (formData.value.provinceSelect) {
    isCity.value = true
    formData.value.province = selectedItems.name
    formData.value.citySelect = null
    formData.value.countySelect = null
    const city = await getGeos({ level: 0, code: selectedItems.code })
    areaCity.value = changeKey(city.data)
  } else {
    isCity.value = false
  }
}

const changeCity = async (e) => {
  const selectedItems = findCodeByValue(areaCity.value, e)
  if (formData.value.citySelect) {
    isCounty.value = true
    formData.value.city = selectedItems.name
    const counties = await getGeos({ level: 1, code: selectedItems.code })
    areaCounty.value = changeKey(counties.data)
  } else {
    isCounty.value = false
  }
}

const changeCounty = (e) => {
  const selectedItems = findCodeByValue(areaCounty.value, e)
  if (formData.value.countySelect) {
    formData.value.county = selectedItems.name
  }
}

const switchChange = (e) => {
  formData.value.active = e.detail.value
}

const confirm = async () => {
  if (!formData.value.name || formData.value.name.trim() === '') {
    uni.showToast({ icon: 'none', title: $t.value('nameRequired') }); return
  }
  if (!formData.value.phone || formData.value.phone.trim() === '') {
    uni.showToast({ icon: 'none', title: $t.value('phoneRequired') }); return
  }
  if (!formData.value.provinceSelect || !formData.value.citySelect || !formData.value.countySelect) {
    uni.showToast({ icon: 'none', title: $t.value('regionRequired') }); return
  }
  if (!formData.value.street || formData.value.street.trim() === '') {
    uni.showToast({ icon: 'none', title: $t.value('streetRequired') }); return
  }

  const data = {
    ID: addrID.value,
    name: formData.value.name,
    phone: formData.value.phone,
    province: Number(formData.value.provinceSelect),
    provinceStr: formData.value.provinceStr,
    city: Number(formData.value.citySelect),
    cityStr: formData.value.cityStr,
    area: Number(formData.value.countySelect),
    areaStr: formData.value.county,
    street: formData.value.street,
    active: formData.value.active,
  }

  const res = await updateAddress(data)
  if (res.code === 0) {
    uni.showToast({ icon: 'none', title: $t.value('saveSuccess') })
    uni.navigateTo({ url: `/pages/address/address?ID=${orderID.value}` })
  }
}

const goBack = () => {
  uni.navigateBack({ fail: () => uni.navigateTo({ url: `/pages/address/address?ID=${orderID.value}` }) })
}
</script>

<style lang="scss" scoped>
.nf-edit-addr {
  min-height: 100vh;
  background: #141414;
  position: relative;
}
.nf-edit-addr-bg {
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
.nf-form-content {
  position: relative;
  z-index: 1;
  padding: calc(var(--status-bar-height, 44rpx) + 88rpx + 30rpx) 24rpx 60rpx;
}
.nf-form-item {
  background: rgba(255, 255, 255, 0.06);
  border-radius: 12rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
}
.nf-form-label {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.5);
  margin-bottom: 12rpx;
  display: block;
}
.nf-form-input {
  font-size: 28rpx;
  color: #fff;
  background: transparent;
  border: none;
  width: 100%;
  padding: 8rpx 0;
}
.nf-form-item-switch {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.nf-form-item-switch .nf-form-label {
  margin-bottom: 0;
}
.nf-select {
  width: 100%;
}
::v-deep .uni-select {
  background: rgba(255, 255, 255, 0.08) !important;
  border-color: rgba(255, 255, 255, 0.1) !important;
  border-radius: 8rpx !important;
}
::v-deep .uni-select__input-text {
  color: #fff !important;
}
::v-deep .uni-select__input-placeholder {
  color: rgba(255, 255, 255, 0.3) !important;
}
::v-deep .uni-select__selector {
  background: #2a2a2a !important;
  border-color: rgba(255, 255, 255, 0.1) !important;
}
::v-deep .uni-select__selector-item {
  color: #fff !important;
}
::v-deep .uni-select__selector-item:hover {
  background: rgba(229, 9, 20, 0.1) !important;
}
.nf-form-save {
  margin-top: 40rpx;
  height: 88rpx;
  border-radius: 12rpx;
  background: #e50914;
  display: flex;
  align-items: center;
  justify-content: center;
}
.nf-form-save text {
  font-size: 30rpx;
  font-weight: bold;
  color: #fff;
}
</style>
