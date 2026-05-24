'use strict';
const dev = {
	baseUrl: "https://www.xxxxxxxx.com/api/wechat", // dev server
	uploadImageUrl: "",
}
const prod = {
	baseUrl: "https://www.xxxxxxxx.com/api/wechat", // dev server
	uploadImageUrl: "",
}
// devBase是开发环境 prodBase是正式环境
export default process.env.NODE_ENV === 'development' ? dev : prod
