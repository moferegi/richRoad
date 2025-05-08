<template>
	<view style="padding: 0 32rpx">
		<view class="form_item flexr-jsb b_b_2 flex-aic">
			<view class="font_28 color_333" style="width: 132rpx;">收件人</view>
			<view class="flex-fitem">
				<input class="color_333 font_28 tar" v-model="formData.name" :maxlength="30" placeholder="请输入收件人" />
			</view>
		</view>
		<view class="form_item flexr-jsb b_b_2 flex-aic">
			<view class="font_28 color_333" style="width: 132rpx;">手机号码</view>
			<view class="flex-fitem">
				<input class="color_333 font_28 tar" v-model="formData.phone" type="number" :maxlength="11"
					placeholder="请输入手机号码" />
			</view>
		</view>
		<view class="form_item flexr-jsb b_b_2 flex-aic">
			<view class="font_28 color_333" style="width: 132rpx;">选择地区</view>
			<view style="width: 280rpx;">
				<view>
					<uni-data-select class="selects" v-model="formData.provinceSelect" placeholder="请选择收件地址"
						:localdata="areaProvince" @change="changeProvince" :clear="false"></uni-data-select>

				</view>
			</view>
		</view>
		<view v-if="isCity" class="form_item flexr-jsb b_b_2 flex-aic">
			<view class="font_28 color_333" style="width: 132rpx;">选择城市</view>
			<view style="width: 360rpx;">
				<view>
					<uni-data-select class="selects" v-model="formData.citySelect" placeholder="请选择收件城市" :localdata="areaCity"
						@change="changeCity" :clear="false"></uni-data-select>
				</view>
			</view>
		</view>
		<view v-if="isCounty" class="form_item flexr-jsb b_b_2 flex-aic">
			<view class="font_28 color_333" style="width: 132rpx;">选择区县</view>
			<view style="width: 360rpx;">
				<view>
					<uni-data-select class="selects" v-model="formData.countySelect" placeholder="请选择收件城市"
						:localdata="areaCounty" @change="changeCounty"  :clear="false"></uni-data-select>
				</view>
			</view>
		</view>
		<view class="form_item flexr-jsb b_b_2 flex-aic">
			<view class="font_28 color_333" style="width: 132rpx;">详细地址</view>
			<view class="flex-fitem">
				<input class="color_333 font_28 tar" v-model="formData.address" placeholder="请输入详细地址" />
			</view>
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
		reactive
	} from 'vue'
	import {
		onLoad
	} from '@dcloudio/uni-app'

	import {
		getGeos,
		createAddress
	} from '@/api/address.js'

	const formData = reactive({
		phone: '',
		name: '',
		address: '',
		province: '',
		provinceSelect: '',
		citySelect: '',
		countySelect: '',
		city: '',
		county: '',
		regionLabel: '',
		checked: false
	})

  const orderID = ref('')
  const areaProvince = ref([])
  const init = async () => {
    const province = await getGeos({
      level: 2,
      code: 0
    })
    areaProvince.value = changeKey(province.data)
  }

  onLoad((options) => {
    // 拿到订单编号
    if(options.ID) {
      orderID.value = options.ID
    }
    init()
  })

	// uniUI 规定下拉控件只能以固定格式固定key  这里洗数据改变id为value改变name为text
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

	const findCodeByValue = (data, selectedValue) => {
		// 使用find方法查找匹配的项
		const selectedItem = data.find(item => item.value === selectedValue);
		// 如果找到了匹配项，返回其code，否则返回null
		return selectedItem ? selectedItem : null;
	}

	const isCity = ref(false)
	const areaCity = ref([])
	const changeProvince = async (e) => {
		const selectedItems = findCodeByValue(areaProvince.value, e);
		if (formData.provinceSelect) {
			isCity.value = true
			formData.province = selectedItems.name
			formData.citySelect = null
			formData.countySelect = null
			const cities = await getGeos({
				level: 0,
				code: selectedItems.code
			})
			areaCity.value = changeKey(cities.data)
		} else {
			isCity.value = false
		}
	}

	const isCounty = ref(false)
	const areaCounty = ref([])
	const changeCity = async (e) => {
		const selectedItems = findCodeByValue(areaCity.value, e);
		if (formData.citySelect) {
			isCounty.value = true
			formData.city = selectedItems.name
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
		if (formData.countySelect) {
			formData.county = selectedItems.name
		}
	}
	const switchChange = (e) => {
		formData.checked = e.detail.value
	}

	const confirm = async() => {
		if (formData.name.trim() == '') {
			uni.showToast({
				icon: "none",
				title: '请输入收件人'
			})
			return
		} else if (formData.phone.trim() == '') {
			uni.showToast({
				icon: "none",
				title: '请输入手机号码'
			})
			return
		} else if (!(/^1[3456789]\d{9}$/.test(formData.phone))) {
			uni.showToast({
				icon: "none",
				title: '请输入正确的手机号码'
			})
			return
		} else if (!formData.provinceSelect || !formData.citySelect || !formData.countySelect) {
			uni.showToast({
				icon: "none",
				title: '请选择地区'
			})
			return
		} else if (formData.address.trim() == '') {
			uni.showToast({
				icon: "none",
				title: '请输入详细地址'
			})
			return
		} else {
			console.log(formData)
		}
		const data = {
			Name: formData.name,
			Phone: formData.phone,
			Province: Number(formData.provinceSelect),
			ProvinceStr: formData.province,
			City: Number(formData.citySelect),
			CityStr: formData.city,
			Area: Number(formData.countySelect),
			AreaStr: formData.county,
			Street: formData.address,
			Active: formData.checked,
		}
		const res = await createAddress (data)
		if(res.code === 0) {
			uni.showToast({
				icon: "none",
				title: '创建成功！'
			})
      // 如果新增地址成功，就携带参数跳转回地址页面并携带参数，因为地址页面根据参数判断是否展示radio
			uni.navigateTo({
				url: `/pages/address/address?ID=${orderID.value}`
			})
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
		bottom: 0;
		left: 0;
		right: 0;
		height: 100rpx;
		width: 100%;
		padding-bottom: constant(safe-area-inset-bottom);
		padding-bottom: env(safe-area-inset-bottom);
	}

	::v-deep .uni-select {
		border: none !important;
		border-bottom: 2rpx solid #c0c0c0 !important;
		// width: 60% !important;
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
