<template>
  <view class="nf-add-addr">
    <view class="nf-add-addr-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff" />
        </view>
        <text class="nf-navbar-title">{{ $t('addNewAddr') }}</text>
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
        <switch @change="switchChange" :checked="formData.checked" color="#e50914" style="transform:scale(0.8)" />
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
import { ref, reactive, computed, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getGeos, createAddress } from '@/api/address.js'
import { getEnabledPhoneAreaCodes } from '@/api/phoneAreaCode.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { localText } from '@/utils/i18n'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)

const formData = reactive({
  phone: '',
  name: '',
  street: '',
  province: '',
  provinceSelect: '',
  citySelect: '',
  countySelect: '',
  city: '',
  county: '',
  checked: false,
  areaCode: '+86'
})

const orderID = ref('')
const areaProvince = ref([])
const areaCity = ref([])
const areaCounty = ref([])
const isCity = ref(false)
const isCounty = ref(false)
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

const init = async () => {
  const province = await getGeos({ level: 2, code: 0 })
  areaProvince.value = changeKey(province.data)
}

const loadAreaCodes = async () => {
  try {
    const res = await getEnabledPhoneAreaCodes()
    if (res.code === 0 && res.data) {
      const list = Array.isArray(res.data) ? res.data : (res.data.list || [])
      areaCodes.value = list
      if (areaCodes.value.length > 0) {
        selectedAreaCode.value = areaCodes.value[0].areaCode
        formData.areaCode = areaCodes.value[0].areaCode
      }
    }
  } catch(e) {}
}

const selectArea = (item) => {
  selectedAreaCode.value = item.areaCode
  formData.areaCode = item.areaCode
  showAreaCodePicker.value = false
}

onLoad((options) => {
  if (options.ID) orderID.value = options.ID
  init()
  loadAreaCodes()
})

const changeProvince = async (e) => {
  const selectedItems = findCodeByValue(areaProvince.value, e)
  if (formData.provinceSelect) {
    isCity.value = true
    formData.province = selectedItems.name
    formData.citySelect = null
    formData.countySelect = null
    const cities = await getGeos({ level: 0, code: selectedItems.code })
    areaCity.value = changeKey(cities.data)
  } else {
    isCity.value = false
  }
}

const changeCity = async (e) => {
  const selectedItems = findCodeByValue(areaCity.value, e)
  if (formData.citySelect) {
    isCounty.value = true
    formData.city = selectedItems.name
    const counties = await getGeos({ level: 1, code: selectedItems.code })
    areaCounty.value = changeKey(counties.data)
  } else {
    isCounty.value = false
  }
}

const changeCounty = (e) => {
  const selectedItems = findCodeByValue(areaCounty.value, e)
  if (formData.countySelect) {
    formData.county = selectedItems.name
  }
}

const switchChange = (e) => {
  formData.checked = e.detail.value
}

const confirm = async () => {
  if (!formData.name || formData.name.trim() === '') {
    uni.showToast({ icon: 'none', title: $t.value('nameRequired') }); return
  }
  if (!formData.phone || formData.phone.trim() === '') {
    uni.showToast({ icon: 'none', title: $t.value('phoneRequired') }); return
  }
  if (!formData.provinceSelect || !formData.citySelect || !formData.countySelect) {
    uni.showToast({ icon: 'none', title: $t.value('regionRequired') }); return
  }
  if (!formData.street || formData.street.trim() === '') {
    uni.showToast({ icon: 'none', title: $t.value('streetRequired') }); return
  }

  const data = {
    Name: formData.name,
    Phone: formData.phone,
    AreaCode: formData.areaCode,
    Province: Number(formData.provinceSelect),
    ProvinceStr: formData.province,
    City: Number(formData.citySelect),
    CityStr: formData.city,
    Area: Number(formData.countySelect),
    AreaStr: formData.county,
    street: formData.street,
    Active: formData.checked,
  }

  const res = await createAddress(data)
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
.nf-add-addr {
  min-height: 100vh;
  background: #141414;
  position: relative;
}
.nf-add-addr-bg {
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

/* 手机号输入行 */
.nf-phone-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
}
.nf-area-code-btn {
  height: 72rpx;
  padding: 0 20rpx;
  background: rgba(255, 255, 255, 0.06);
  border: 1rpx solid rgba(255, 255, 255, 0.1);
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  gap: 8rpx;
  flex-shrink: 0;
}
.nf-area-code-text {
  color: #fff;
  font-size: 28rpx;
  font-weight: 600;
}
.nf-area-code-arrow {
  color: rgba(255, 255, 255, 0.4);
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
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: flex-end;
  justify-content: center;
}
.nf-popup-content {
  width: 100%;
  max-height: 60vh;
  background: #1a1a1a;
  border-radius: 28rpx 28rpx 0 0;
  padding: 32rpx 0;
}
.nf-popup-title {
  text-align: center;
  font-size: 32rpx;
  font-weight: 700;
  color: #fff;
  padding-bottom: 24rpx;
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.08);
}
.nf-popup-scroll {
  max-height: 50vh;
}
.nf-area-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 40rpx;
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.05);
  &:active { background: rgba(255, 255, 255, 0.05); }
}
.nf-area-name {
  color: #fff;
  font-size: 28rpx;
}
.nf-area-code-val {
  color: rgba(229, 9, 20, 0.8);
  font-size: 28rpx;
  font-weight: 600;
}
</style>
