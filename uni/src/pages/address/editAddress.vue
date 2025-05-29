<template>
	<view style="padding: 0 32rpx;">
		<view class="form_item flexr-jsb b_b_2 flex-aic">
			<view class="font_28 color_333" style="width: 132rpx;">收件人</view>
      <wu-input
          placeholder="请输入收件人姓名"
          border="surround"
          v-model="formData.name"
      ></wu-input>
		</view>
		<view class="form_item flexr-jsb b_b_2 flex-aic">
			<view class="font_28 color_333" style="width: 132rpx;">手机号</view>
      <wu-input
          placeholder="请输入手机号码"
          border="surround"
          :maxlength="11"
          v-model="formData.phone"
      ></wu-input>
		</view>
    <view class="form_item b_b_2 flex flex-aic">
      <view class="font_28 color_333" style="width: 132rpx;">地区</view>
      <view style="width: 80%;">
        <view>
          <uni-data-select class="selects" v-model="formData.provinceSelect" placeholder="请选择收件地址"
                           :localdata="areaProvince" @change="changeProvince" :clear="false"></uni-data-select>
        </view>
      </view>
    </view>
    <view v-if="isCity" class="form_item flexr-jsb b_b_2 flex-aic">
      <view class="font_28 color_333" style="width: 132rpx;">城市</view>
      <view style="width: 80%;">
        <view>
          <uni-data-select class="selects" v-model="formData.citySelect" placeholder="请选择收件城市" :localdata="areaCity"
                           @change="changeCity" :clear="false"></uni-data-select>
        </view>
      </view>
    </view>
    <view v-if="isCounty" class="form_item flexr-jsb b_b_2 flex-aic">
      <view class="font_28 color_333" style="width: 132rpx;">区县</view>
      <view style="width: 80%;">
        <view>
          <uni-data-select class="selects" v-model="formData.countySelect" placeholder="请选择收件区县"
                           :localdata="areaCounty" @change="changeCounty"  :clear="false"></uni-data-select>
        </view>
      </view>
    </view>
    <view class="form_item flexr-jsb b_b_2 flex-aic">
      <view class="font_28 color_333" style="width: 132rpx;">详细地址</view>
      <wu-input
          placeholder="请输入详细地址"
          border="surround"
          v-model="formData.street"
      ></wu-input>
    </view>
		<view class="form_item flexr-jsb b_b_2 flex-aic">
			<view class="font_28 color_333" style="width: 172rpx;">设为默认地址</view>
			<switch @change="switchChange" :checked="formData.checked" color="#fe5572" style="transform:scale(0.7)" />
		</view>


		<view class="address_nav_view">
			<view class="address_nav_box pos_f bgc_fff flex-aic flexr-jsc">
				<button class="btns" @tap="confirm">保存</button>
			</view>
		</view>
	</view>
</template>

<script setup>
	import {
		ref,
		watch
	} from 'vue'
	import { onLoad } from '@dcloudio/uni-app'
	import {
		getGeos,
		updateAddress,
		findAddress
	} from '@/api/address.js'

	const formData = ref({})

	const areaProvince = ref([])
	// 洗数据改变id为value改变name为text
	const changeKey = (data) => {
		const newData = data.map(item => {
			// 创建一个新的对象，将id属性复制为value属性，name属性复制为text属性，其他属性保持不变
			const newItem = {
				...item,
				value: item.code,
				text: item.name
			};
			return newItem;
		});
		return newData
	}

	function findCodeByValue(data, selectedValue) {
		// 使用find方法查找匹配的项
		const selectedItem = data.find(item => item.value === selectedValue);
		// 如果找到了匹配项，返回其code，否则返回null
		return selectedItem ? selectedItem : null;
	}

	const isCity = ref(false)
	const isCounty = ref(false)
	const areaCity = ref([])
  const orderID = ref('')
  const addrID = ref('')
	onLoad(async (e) => {
      orderID.value = e.orderID
		const province = await getGeos({
			level: 2,
			code: 0
		})
		areaProvince.value = changeKey(province.data)
		// 根据ID拉取用户输入的信息
		const res = await findAddress({
			id: e.ID
		})
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

	const loadCities = async (provinceId) => {
		const city = await getGeos({
			level: 0,
			code: provinceId
		});
		areaCity.value = changeKey(city.data);
		console.log(areaCity.value);
	};

	// 根据城市ID获取县区数据
	const loadCounties = async (cityId) => {
		const counties = await getGeos({
			level: 1,
			code: cityId
		})
		areaCounty.value = changeKey(counties.data)
	};

	// 监听 provinceSelect 的变化
	watch(() => formData.value.provinceSelect, (newVal) => {
		if (newVal) {
			loadCities(newVal)
		}
	});

	// 监听 citySelect 的变化
	watch(() => formData.value.citySelect, (newVal) => {
		if (newVal) {
			loadCounties(newVal);
		}
	});

	const changeProvince = async (e) => {
		const selectedItems = findCodeByValue(areaProvince.value, e);
		if (formData.value.provinceSelect) {
			isCity.value = true
			formData.value.province = selectedItems.name
			console.log(formData.value)
			formData.value.citySelect = null
			formData.value.countySelect = null
			// loadCities(selectedItems.code)
			const city = await getGeos({
				level: 0,
				code: selectedItems.code
			});
			areaCity.value = changeKey(city.data);
		} else {
			isCity.value = false
		}
	}

	const areaCounty = ref([])
	const changeCity = async (e) => {
		const selectedItems = findCodeByValue(areaCity.value, e);
		if (formData.value.citySelect) {
			isCounty.value = true
			formData.value.city = selectedItems.name
			// loadCounties(selectedItems.code)
			const counties = await getGeos({
				level: 1,
				code: selectedItems.code
			})
			areaCounty.value = changeKey(counties.data)
		} else {
			isCounty.value = false
		}
	}
	const changeCounty = async (e) => {
		const selectedItems = findCodeByValue(areaCounty.value, e);
		if (formData.value.countySelect) {
			formData.value.county = selectedItems.name
		}
	}

	const switchChange = (e) => {
		formData.value.active = e.detail.value
	}

	const confirm = async () => {
		if (formData.value.name.trim() == '') {
			uni.showToast({
				icon: "none",
				title: '请输入收件人'
			})
			return
		} else if (formData.value.phone.trim() == '') {
			uni.showToast({
				icon: "none",
				title: '请输入手机号码'
			})
			return
		} else if (!(/^1[3456789]\d{9}$/.test(formData.value.phone))) {
			uni.showToast({
				icon: "none",
				title: '请输入正确的手机号码'
			})
			return
		} else if (!formData.value.provinceSelect || !formData.value.citySelect || !formData.value.countySelect) {
			uni.showToast({
				icon: "none",
				title: '请选择地区'
			})
			return
		} else if (formData.value.street.trim() == '') {
			uni.showToast({
				icon: "none",
				title: '请输入详细地址'
			})
			return
		} else {
			const data = {
        ID: addrID.value,
				name: formData.value.name,
				phone: formData.value.phone,
				province: Number(formData.value.provinceSelect),
				provinceStr: formData.value.province,
				city: Number(formData.value.citySelect),
				cityStr: formData.value.city,
				area: Number(formData.value.countySelect),
				street: formData.value.street,
				areaStr: formData.value.county,
				active: formData.value.active,
			}
      console.log(data);

      const res = await updateAddress(data)
      if (res.code === 0) {
        uni.showToast({
          icon: "none",
          title: '编辑成功！'
        })
        uni.navigateTo({
          url: `/pages/address/address?ID=${orderID.value}`
        })
      }
		}

	}
</script>

<style lang="scss" scoped>
 	page{
		background:#fff;
	}
	.form_item {
		padding: 28rpx 0;
	}

	.address_nav_view {
		height: 100rpx;
		width: 100%;
		padding-bottom: constant(safe-area-inset-bottom);
		padding-bottom: env(safe-area-inset-bottom);
	}

	.address_nav_box {
    bottom: 30rpx;
		left: 0;
		right: 0;
		height: 100rpx;
		width: 100%;
		padding-bottom: constant(safe-area-inset-bottom);
		padding-bottom: env(safe-area-inset-bottom);
	}

	::v-deep .uni-select {
    width: 100%;
	}

	.text-ellipsis {
		overflow: hidden;
		/* 内容超出容器时不显示 */
		text-overflow: ellipsis;
		/* 超出部分以省略号显示 */
		white-space: nowrap;
		/* 强制文本在同一行内显示，不换行 */
		width: 280rpx;
		/* 你希望设置的宽度 */
		;
		/* 必须设置一个宽度或者max-width */
	}

	.btns {
		width: 100%;
		padding: 6rpx;
		margin: 0 20rpx;
		color: white;
		border-color: #fe5572;
		background-color: #fe5572;
	}
</style>
