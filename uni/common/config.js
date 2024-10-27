'use strict';
const dev = {
	baseUrl: "https://www.xxxxxxxx.com/api/wechat", // dev server
	uploadImageUrl: "https://xxxxxxx.oss-cn-hangzhou.aliyuncs.com/",
	AccessKeySecret: "318ESFytpYSyrdDYfxhjHkcmmCSrxd",
	OSSAccessKeyId: "LTAI7yDacD6sud82oYiB2ChT",
}
const prod = {
	baseUrl: "https://www.xxxxxxxx.com/api/wechat", // dev server
	uploadImageUrl: "https://xxxxxxxx.oss-accelerate.aliyuncs.com/",
	AccessKeySecret: "318ESFytpYSyrdDYfxhjHkcmmCSrxd",
	OSSAccessKeyId: "LTAI7yDacD6sud82oYiB2ChT",
}
// devBase是开发环境 prodBase是正式环境
export default process.env.NODE_ENV === 'development' ? dev : prod
