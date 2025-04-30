// import lyzDebounce from '../js_sdk/lyz-debounce/lyz-debounce.js'
export const login = (callback) => {
	let token = uni.getStorageSync('token')
	if (token) {
		callback(token, JSON.parse(uni.getStorageSync('userInfo')))
		return
	}
	uni.login({
		provider: 'weixin',
		success: (res) => {
			callback(res)
		}
	})
}
//
// export const lyzLocation = (callback) => {
// 	uni.showLoading({
// 		title: '正在定位'
// 	})
// 	lyzDebounce.canDoFunction({
// 		key: "lyzLocation",
// 		time: 40000,
// 		success: () => {
// 			uni.getLocation({
// 				type: 'gcj02',
// 				isHighAccuracy: true,
// 				highAccuracyExpireTime: 4000,
// 				success: (res) => {
// 					uni.setStorageSync('location', JSON.stringify(res))
// 					uni.hideLoading()
// 					callback(res)
// 				},
// 				fail: (err) => {
// 					uni.hideLoading()
// 					uni.removeStorageSync('location')
// 					uni.showToast({
// 						icon: "none",
// 						title: '定位失败，请打开定位权限'
// 					})
// 				}
// 			});
// 		},
// 		fail: () => {
// 			uni.hideLoading()
// 			callback(JSON.parse(uni.getStorageSync('location') ||
// 				'{"latitude":null,"longitude":null,"speed":-1,"accuracy":65,"verticalAccuracy":65,"horizontalAccuracy":65,"errMsg":"getLocation:ok"}'
// 				) || {
// 				"latitude": null,
// 				"longitude": null,
// 				"speed": -1,
// 				"accuracy": 65,
// 				"verticalAccuracy": 65,
// 				"horizontalAccuracy": 65,
// 				"errMsg": "getLocation:ok"
// 			})
// 		}
// 	})
// }
export const isAdmin = (action) => {
	const perList = JSON.parse(uni.getStorageSync('userInfo')).roles
	if (!perList || perList.length == 0) return false
	if (perList[0] == action) return true
	if (perList[0] !== action) return false
}
// 判断平台 苹果返回true 安卓返回false
export const isIos = () => {
	// #ifdef APP-PLUS
	if (plus.os.name.toLowerCase() === 'android') {
		return false
	} else {
		return true
	}
	// #endif
}
export const getOssImgName = () => {
	return Math.random().toString(36).slice(2)
}
// 时间戳转换 传入毫秒级时间戳 返回YYYY-MM-DD hh-mm-ss
export const timeFormat = (time, type = 'YYYY-MM-DD hh-mm-ss', line = "-") => {
	if (!time) return ''
	let date = new Date(time.replace(/-/g, "/"));
	let Y = date.getFullYear() + line;
	let M = (date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1) + line;
	let D = (date.getDate() < 10 ? '0' + date.getDate() : date.getDate());
	let h = (date.getHours() < 10 ? '0' + date.getHours() : date.getHours()) + ':';
	let m = (date.getMinutes() < 10 ? '0' + date.getMinutes() : date.getMinutes()) + ':';
	let s = (date.getSeconds() < 10 ? '0' + date.getSeconds() : date.getSeconds());
	if (type == 'YYYY-MM-DD hh-mm-ss') {
		return Y + M + D + ' ' + h + m + s
	}
	if (type == 'YYYY-MM-DD') {
		return Y + M + D
	} else {
		return h + m + s
	}
}
/**
 * @param {Number} timeStamp 传入秒级的时间戳
 */
export const getDate = (timeStamp) => {
	const d = new Date(timeStamp * 1000)
	let tiem = Math.ceil(new Date().getTime() / 1000)
	const year = d.getFullYear()
	const month = getHandledValue(d.getMonth() + 1)
	const date = getHandledValue(d.getDate())
	const hours = getHandledValue(d.getHours())
	const minutes = getHandledValue(d.getMinutes())
	const second = getHandledValue(d.getSeconds())
	let resStr = ''
	if (!isToday(timeStamp)) {
		// 大于一天
		resStr = month + '-' + date + ' ' + hours + ':' + minutes
	} else {
		//小于一天
		resStr = hours + ':' + minutes
	}
	return resStr
}
/**
 * @param {Number} num 数值
 * @returns {String} 处理后的字符串
 * @description 如果传入的数值小于10，即位数只有1位，则在前面补充0
 */
const getHandledValue = num => {
	return num < 10 ? '0' + num : num
}
/**
 * 判断是否为今天
 * @param {Number} timeStamp 秒级时间戳
 * 
 */
const isToday = (timeStamp) => {
	var d = new Date(timeStamp * 1000);
	var todaysDate = new Date();
	if (d.setHours(0, 0, 0, 0) == todaysDate.setHours(0, 0, 0, 0)) {
		return true;
	} else {
		return false;
	}
}
// 多久之前 传入YYYY-MM-DD hh-mm-ss 返回多久之前
export const dateTimeStamp = (dateTime) => {
	let date = new Date(dateTime.replace(/-/g, "/")).getTime();
	let minute = 1000 * 60; //把分，时，天，周，半个月，一个月用毫秒表示
	let hour = minute * 60;
	let day = hour * 24;
	let week = day * 7;
	let halfamonth = day * 15;
	let month = day * 30;
	let result = ''
	let now = new Date().getTime(); //获取当前时间毫秒
	let diffValue = now - date; //时间差

	if (diffValue < 0) {
		return dateTime;
	}
	let minC = diffValue / minute; //计算时间差的分，时，天，周，月
	let hourC = diffValue / hour;
	let dayC = diffValue / day;
	let weekC = diffValue / week;
	let monthC = diffValue / month;
	if (dayC > 3) {
		let datestr = dateTime
		result = datestr.split(' ')[0];
	} else if (monthC >= 1) {
		result = "" + parseInt(monthC) + "月前";
	} else if (weekC >= 1) {
		result = "" + parseInt(weekC) + "周前";
	} else if (dayC >= 1) {
		result = "" + parseInt(dayC) + "天前";
	} else if (hourC >= 1) {
		result = "" + parseInt(hourC) + "小时前";
	} else if (minC >= 1) {
		result = "" + parseInt(minC) + "分钟前";
	} else {
		result = "刚刚";
	}
	return result;

}
// 获取三级联动日历数据
export const getYearListMy = (min, max) => {
	let days = [];
	let years = [];
	let months = []
	// let months = [1,2,3,4,5,6,7,8,9,10,11,12];
	// let date = [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31]
	for (let i = 1; i <= 31; i++) {
		days.push(i + '日');
	}
	for (let monthsi = 1; monthsi <= 12; monthsi++) {
		months.push(monthsi + '月');
	}
	for (let Yeari = min; Yeari <= max; Yeari++) {
		years.push(Yeari + '年');
	}
	return [years, months, days]
}
// 获取三级联动日历数据
export const getYearList = (min, max, callback) => {
	let days = [];
	let months = [];
	let years = [];
	let date = []
	for (let i = 1; i <= 31; i++) {
		let obj = {}
		obj.label = i + '号'
		obj.value = i
		days.push(obj);
	}
	for (let monthsi = 1; monthsi <= 12; monthsi++) {
		let obj = {}
		obj.label = monthsi + '月'
		obj.value = monthsi
		obj.children = days
		months.push(obj);
	}
	for (let Yeari = min; Yeari <= max; Yeari++) {
		let obj = {}
		obj.label = Yeari + '年'
		obj.value = Yeari
		obj.children = months
		years.push(obj);
	}
	callback(years)
}
export const urlHttp = (url) => {
	if (url && url !== '') {
		if (url.substr(0, 7).toLowerCase() == "http://" || url.substr(0, 8).toLowerCase() == "https://") {
			return url;
		} else if (url.substr(0, 2).toLowerCase() == "//") {
			return "https:" + url;
		} else {
			return "https://" + url;
		}
	} else {
		return '/static/temp/avatarNo.jpg'
	}
}
// 将rpx转px
export const sizeDeal = (size) => {
	const info = uni.getSystemInfoSync()
	let scale = 750 / info.windowWidth;
	// 分离字体大小和单位,rpx 转 px
	let s = Number.isNaN(parseFloat(size)) ? 0 : parseFloat(size)
	let u = size.toString().replace(/[0-9]/g, '').replace('-', '')
	if (u == 'rpx') {
		s /= scale
		u = 'px'
	} else if (u == '') {
		u = 'px'
	} else if (u == 'vw') {
		u = 'px'
		s = s / 100 * 750 / scale
	}
	return [s, u, s + u]
}
// 区分是否是微信内置浏览器
export const isMicroMessenger = () => {
	let result = false;
	// #ifdef APP-PLUS
	result = true
	// #endif
	// #ifdef H5
	let userAgent = window.navigator.userAgent;
	if (userAgent.indexOf('MicroMessenger') > -1) {
		result = true;
	}
	// #endif
	return result;
}
// 深拷贝函数  接收目标target参数
export const deepClone = (target) => {
	//  var obj1 = {
	//     a: 1,
	//     b: 2,
	//     c: ['a','b','c']
	// }
	// var obj2 = Object.assign({}, obj1);
	// obj2.c[1] = 5;
	// console.log(obj1.c); // ["a", 5, "c"]
	// console.log(obj2.c); // ["a", 5, "c"]
	// Object.assign()对于一层对象来说是没有任何问题的，但是如果对象的属性对应的是其它的引用类型的话，还是只拷贝了引用，修改的话还是会有问题
	// 所以要使用递归
	// 定义一个变量
	let result;
	// 如果当前需要深拷贝的是一个对象的话
	if (typeof target === 'object') {
		// 如果是一个数组的话
		if (Array.isArray(target)) {
			result = []; // 将result赋值为一个数组，并且执行遍历
			for (let i in target) {
				// 递归克隆数组中的每一项
				result.push(deepClone(target[i]))
			}
			// 判断如果当前的值是null的话；直接赋值为null
		} else if (target === null) {
			result = null;
			// 判断如果当前的值是一个RegExp对象的话，直接赋值    
		} else if (target.constructor === RegExp) {
			result = target;
		} else {
			// 否则是普通对象，直接for in循环，递归赋值对象的所有值
			result = {};
			for (let i in target) {
				result[i] = deepClone(target[i]);
			}
		}
		// 如果不是对象的话，就是基本数据类型，那么直接赋值
	} else {
		result = target;
	}
	// 返回最终结果
	return result;
}
// 获取 00:00-00:59 24小时制
export const getHour = () => {
	let getHourList = []
	for (let i = 0; i <= 23; i++) {
		if (i == 0) {
			getHourList.push('未知')
		}
		getHourList.push(`${i >= 10 ? '' : '0'}${i}:00-${i >= 10 ? '' : '0'}${i}:59`)
	}
	return [getHourList]
}
// 限制输入框只能输入数字 保留两位小数
export const getNumTo = (value) => {
	let val = value.replace(/(^\s*)|(\s*$)/g, "")
	console.log(val)
	if (!val) {
		return null
	}
	let reg = /[^\d.]/g

	// 只能是数字和小数点，不能是其他输入
	val = val.replace(reg, "")
	// // 保证第一位只能是数字，不能是点
	val = val.replace(/^\./g, "");
	// // 小数只能出现1位
	val = val.replace(".", "$#$").replace(/\./g, "").replace("$#$", ".");
	// // 小数点后面保留2位
	val = val.replace(/^(\-)*(\d+)\.(\d\d).*$/, '$1$2.$3');
	return val
}
// 限制输入框只能输入数字 
export const getNum = (value) => {
	let val = value.replace(/(^\s*)|(\s*$)/g, "")
	console.log(val)
	if (!val) {
		return null
	}
	let reg = /[^\d]/g
	val = val.replace(reg, "")
	return val
}

export const checkAuditTime = (beginTime, endTime) => {
	var nowDate = new Date();
	var beginDate = new Date(nowDate);
	var endDate = new Date(nowDate);

	var beginIndex = beginTime.lastIndexOf("\:");
	var beginHour = beginTime.substring(0, beginIndex);
	var beginMinue = beginTime.substring(beginIndex + 1, beginTime.length);
	beginDate.setHours(beginHour, beginMinue, 0, 0);

	var endIndex = endTime.lastIndexOf("\:");
	var endHour = endTime.substring(0, endIndex);
	var endMinue = endTime.substring(endIndex + 1, endTime.length);
	endDate.setHours(endHour, endMinue, 0, 0);
	return nowDate.getTime() - beginDate.getTime() >= 0 && nowDate.getTime() <= endDate.getTime();
}

export const initTime = (time) => {
	if (time <= 59) {
		let second = time
		return [0, 0, second]
	} else if (time > 59 && time <= 3599) {
		let minute = parseInt(time / 60)
		let second = parseInt(time - (minute * 60))
		return [0, minute, second]
	} else if (time > 3599) {
		let hour = parseInt(time / 3600)
		let minute = parseInt((time - (hour * 60 * 60)) / 60)
		let second = parseInt((time - (hour * 60 * 60)) - (minute * 60))
		return [hour, minute, second]
	}
}

export const formatSeconds = (value) => {
	let theTime = parseInt(value); // 秒
	let middle = 0; // 分
	let hour = 0; // 小时
	if (theTime > 60) {
		middle = parseInt(theTime / 60);
		theTime = parseInt(theTime % 60);
		if (middle > 60) {
			hour = parseInt(middle / 60);
			middle = parseInt(middle % 60);
		}
	}
	let result = "" + parseInt(theTime) + "秒";
	if (middle > 0) {
		result = "" + parseInt(middle) + "分" + result;
	}
	if (hour > 0) {
		result = "" + parseInt(hour) + "小时" + result;
	}
	return result;
}

export const insertStr = (str) => {
	let re = '';
	let length = str.length;
	for (let i = 0, j = 1; i < length; i++, j++) {
		if (j && j % 11 == 0) {
			re += '\n';
		} else {
			re += str[i];
		}
	}
	return re;
}

export const richFilters = (content) => {
	if (content) return content.replace(/<img/g, "<img style='width:100%'")
	return ''
}
