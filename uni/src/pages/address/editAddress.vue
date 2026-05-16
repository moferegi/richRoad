<template>
  <view class="nf-edit-addr">
    <view class="nf-edit-addr-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="rgba(15,23,42,0.78)" />
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
        <view class="nf-phone-row">
          <view class="nf-area-code-btn" @tap="showAreaCodePicker = true">
            <text class="nf-area-code-text">{{ selectedAreaCode }}</text>
            <text class="nf-area-code-arrow">▼</text>
          </view>
          <input class="nf-form-input nf-phone-input" v-model="formData.phone" :placeholder="$t('phoneRequired')" :maxlength="20" />
        </view>
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
        <switch @change="switchChange" :checked="formData.active" color="#2563eb" style="transform:scale(0.8)" />
      </view>

      <!-- 保存按钮 -->
      <view class="nf-form-save" @tap="confirm">
        <text>{{ $t('saveAddr') }}</text>
      </view>
    </view>

    <!-- 区号选择弹窗 -->
    <view class="nf-popup-mask" v-if="showAreaCodePicker" @tap="showAreaCodePicker = false">
      <view class="nf-popup-content" @tap.stop>
        <view class="nf-popup-title">{{ $t('selectAreaCode') }}</view>
        <scroll-view scroll-y class="nf-popup-scroll">
          <view class="nf-area-item" v-for="item in areaCodes" :key="item.ID" @tap="selectArea(item)">
            <text class="nf-area-name">{{ $lt(item.countryName) }}</text>
            <text class="nf-area-code-val">{{ item.areaCode }}</text>
          </view>
        </scroll-view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getGeos, updateAddress, findAddress } from '@/api/address.js'
import { getEnabledPhoneAreaCodes } from '@/api/phoneAreaCode.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { localText } from '@/utils/i18n'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)

const formData = ref({})
const areaProvince = ref([])
const areaCity = ref([])
const areaCounty = ref([])
const isCity = ref(false)
const isCounty = ref(false)
const orderID = ref('')
const addrID = ref('')
const showAreaCodePicker = ref(false)
const areaCodes = ref([])
const selectedAreaCode = ref('+86')

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

const loadAreaCodes = async () => {
  try {
    const res = await getEnabledPhoneAreaCodes()
    if (res.code === 0 && res.data) {
      const list = Array.isArray(res.data) ? res.data : (res.data.list || [])
      areaCodes.value = list
    }
  } catch(e) {}
}

const selectArea = (item) => {
  selectedAreaCode.value = item.areaCode
  formData.value.areaCode = item.areaCode
  showAreaCodePicker.value = false
}

onLoad(async (e) => {
  orderID.value = e.orderID
  const province = await getGeos({ level: 2, code: 0 })
  areaProvince.value = changeKey(province.data)
  loadAreaCodes()

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
    // 回显区号
    selectedAreaCode.value = info.areaCode || '+86'
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
    formData.value.provinceStr = selectedItems.name
    formData.value.citySelect = null
    formData.value.countySelect = null
    formData.value.cityStr = ''
    formData.value.county = ''
    formData.value.areaStr = ''
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
    formData.value.cityStr = selectedItems.name
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
    formData.value.areaStr = selectedItems.name
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
    areaCode: formData.value.areaCode || selectedAreaCode.value,
    province: Number(formData.value.provinceSelect),
    provinceStr: formData.value.provinceStr || formData.value.province,
    city: Number(formData.value.citySelect),
    cityStr: formData.value.cityStr || formData.value.city,
    area: Number(formData.value.countySelect),
    areaStr: formData.value.areaStr || formData.value.county,
    street: formData.value.street,
    active: formData.value.active,
  }

  const res = await updateAddress(data)
  if (res.code === 0) {
    uni.showToast({ icon: 'none', title: $t.value('saveSuccess') })
    uni.navigateBack()
  }
}

const goBack = () => {
  uni.navigateBack({ fail: () => uni.navigateTo({ url: `/pages/address/address?ID=${orderID.value}` }) })
}
</script>

<style lang="scss" scoped>
.nf-edit-addr {
  min-height: 100vh;
  background: #f4f7fb;
  position: relative;
}
.nf-edit-addr-bg {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: radial-gradient(120% 80% at 100% -10%, #dbeafe 0%, transparent 60%), #f4f7fb;
  z-index: 0;
}
.nf-navbar {
  position: fixed;
  top: 0; left: 0; right: 0;
  z-index: 100;
  background: rgba(244, 247, 251, 0.94);
  border-bottom: 1rpx solid rgba(15, 23, 42, 0.08);
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
  background: rgba(255, 255, 255, 0.9);
  border: 1rpx solid rgba(15, 23, 42, 0.08);
  display: flex;
  align-items: center;
  justify-content: center;
}
.nf-navbar-title {
  font-size: 34rpx;
  font-weight: bold;
  color: #0f172a;
}
.nf-form-content {
  position: relative;
  z-index: 1;
  padding: calc(var(--status-bar-height, 44rpx) + 88rpx + 30rpx) 24rpx 60rpx;
}
.nf-form-item {
  background: rgba(255, 255, 255, 0.94);
  border: 1rpx solid rgba(15, 23, 42, 0.08);
  border-radius: 12rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
}
.nf-form-label {
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.56);
  margin-bottom: 12rpx;
  display: block;
}
.nf-form-input {
  font-size: 28rpx;
  color: #0f172a;
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
  background: rgba(255, 255, 255, 0.95) !important;
  border-color: rgba(15, 23, 42, 0.12) !important;
  border-radius: 8rpx !important;
}
::v-deep .uni-select__input-text {
  color: #0f172a !important;
}
::v-deep .uni-select__input-placeholder {
  color: rgba(15, 23, 42, 0.38) !important;
}
::v-deep .uni-select__selector {
  background: #ffffff !important;
  border-color: rgba(15, 23, 42, 0.12) !important;
}
::v-deep .uni-select__selector-item {
  color: #0f172a !important;
}
::v-deep .uni-select__selector-item:hover {
  background: rgba(219, 234, 254, 0.5) !important;
}
.nf-form-save {
  margin-top: 40rpx;
  height: 88rpx;
  border-radius: 12rpx;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  display: flex;
  align-items: center;
  justify-content: center;
}
.nf-form-save text {
  font-size: 30rpx;
  font-weight: bold;
  color: #fff;
}

/* 手机号输入行 */
.nf-phone-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
}
.nf-area-code-btn {
  height: 72rpx;
  padding: 0 20rpx;
  background: rgba(255, 255, 255, 0.92);
  border: 1rpx solid rgba(15, 23, 42, 0.12);
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  gap: 8rpx;
  flex-shrink: 0;
}
.nf-area-code-text {
  color: #0f172a;
  font-size: 28rpx;
  font-weight: 600;
}
.nf-area-code-arrow {
  color: rgba(15, 23, 42, 0.4);
  font-size: 20rpx;
}
.nf-phone-input {
  flex: 1;
}

/* 区号弹窗 */
.nf-popup-mask {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  z-index: 200;
  background: rgba(15, 23, 42, 0.36);
  display: flex;
  align-items: flex-end;
  justify-content: center;
}
.nf-popup-content {
  width: 100%;
  max-height: 60vh;
  background: #ffffff;
  border-radius: 28rpx 28rpx 0 0;
  padding: 32rpx 0;
  border: 1rpx solid rgba(15, 23, 42, 0.08);
}
.nf-popup-title {
  text-align: center;
  font-size: 32rpx;
  font-weight: 700;
  color: #0f172a;
  padding-bottom: 24rpx;
  border-bottom: 1rpx solid rgba(15, 23, 42, 0.08);
}
.nf-popup-scroll {
  max-height: 50vh;
}
.nf-area-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 40rpx;
  border-bottom: 1rpx solid rgba(15, 23, 42, 0.06);
  &:active { background: rgba(219, 234, 254, 0.4); }
}
.nf-area-name {
  color: #0f172a;
  font-size: 28rpx;
}
.nf-area-code-val {
  color: rgba(37, 99, 235, 0.9);
  font-size: 28rpx;
  font-weight: 600;
}
</style>
