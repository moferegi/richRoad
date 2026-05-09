package initialize

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/songzhibin97/gkit/cache/local_cache"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
)

func OtherInit() {
	// 初始化 IP 归属地查询
	if err := utils.InitIPSearcher("resource/ip2region/ip2region.xdb"); err != nil {
		global.GVA_LOG.Warn("IP归属地查询初始化失败: " + err.Error())
	}
	dr, err := utils.ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
	if err != nil {
		panic(err)
	}
	_, err = utils.ParseDuration(global.GVA_CONFIG.JWT.BufferTime)
	if err != nil {
		panic(err)
	}

	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(dr),
	)
	file, err := os.Open("go.mod")
	if err == nil && global.GVA_CONFIG.AutoCode.Module == "" {
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		global.GVA_CONFIG.AutoCode.Module = strings.TrimPrefix(scanner.Text(), "module ")
	}

	// 初始化默认系统参数（需要DB就绪后调用）
}

// InitInviteSystem 初始化邀请系统默认配置，需在DB初始化后调用
func InitInviteSystem() {
	if global.GVA_DB == nil {
		return
	}
	initDefaultSysConfigs()
}

func initDefaultSysConfigs() {
	defaults := []client.SysConfig{
		// invite 分组
		{ConfigKey: "invite_reward_points", ConfigValue: "10", ConfigName: "邀请奖励积分", ConfigGroup: "invite", Remark: "用户邀请好友注册后获得的积分奖励"},
		// payment 分组
		{ConfigKey: "payment_tip_text", ConfigValue: "{\"zh\":\"请在规定时间内完成付款\",\"en\":\"Please complete payment within the specified time\",\"mn\":\"Заасан хугацаанд төлбөрөө хийнэ үү\"}", ConfigName: "付款提示文本", ConfigGroup: "payment", Remark: "二维码付款弹窗提示文本(JSON多语言)"},
		{ConfigKey: "payment_tip_text_size", ConfigValue: "14", ConfigName: "付款提示文字大小", ConfigGroup: "payment", Remark: "付款提示文本字体大小(px)"},
		{ConfigKey: "payment_tip_text_color", ConfigValue: "#ff0000", ConfigName: "付款提示文字颜色", ConfigGroup: "payment", Remark: "付款提示文本颜色(hex)"},
		{ConfigKey: "payment_auto_enabled", ConfigValue: "false", ConfigName: "自动支付总开关", ConfigGroup: "payment", Remark: "是否启用自动支付渠道(true/false)，默认关闭保留人工流程"},
		{ConfigKey: "payment_manual_qrcode_enabled", ConfigValue: "true", ConfigName: "人工二维码支付开关", ConfigGroup: "payment", Remark: "是否展示人工二维码付款(true/false)"},
		{ConfigKey: "payment_manual_contact_enabled", ConfigValue: "true", ConfigName: "人工联系客服支付开关", ConfigGroup: "payment", Remark: "是否展示联系客服付款(true/false)"},
		{ConfigKey: "payment_manual_methods", ConfigValue: "[{\"key\":\"qrcode\",\"enabled\":true,\"manual\":true,\"sort\":10,\"name\":{\"zh\":\"二维码支付\",\"en\":\"QR Payment\",\"mn\":\"QR төлбөр\"},\"copyText\":{\"zh\":\"请扫码后备注订单号\",\"en\":\"Please include your order number after payment\",\"mn\":\"Төлбөр хийхдээ захиалгын дугаараа тэмдэглэнэ үү\"},\"image\":\"\",\"externalPath\":\"\"},{\"key\":\"contact\",\"enabled\":true,\"manual\":true,\"sort\":20,\"name\":{\"zh\":\"联系客服\",\"en\":\"Contact Support\",\"mn\":\"Хэрэглэгчийн дэмжлэг\"},\"copyText\":{\"zh\":\"付款后请联系客服并提供订单号\",\"en\":\"Please contact support with your order number after payment\",\"mn\":\"Төлбөр хийсний дараа захиалгын дугаартайгаа客服-т хандана уу\"},\"image\":\"\",\"externalPath\":\"\"},{\"key\":\"wechat\",\"enabled\":false,\"manual\":false,\"sort\":30,\"name\":{\"zh\":\"微信支付\",\"en\":\"WeChat Pay\",\"mn\":\"WeChat Pay\"},\"copyText\":{},\"image\":\"\",\"externalPath\":\"\"},{\"key\":\"alipay\",\"enabled\":false,\"manual\":false,\"sort\":40,\"name\":{\"zh\":\"支付宝\",\"en\":\"Alipay\",\"mn\":\"Alipay\"},\"copyText\":{},\"image\":\"\",\"externalPath\":\"\"},{\"key\":\"bank_card_cn\",\"enabled\":false,\"manual\":false,\"sort\":50,\"name\":{\"zh\":\"银行卡(国内)\",\"en\":\"Bank Card (CN)\",\"mn\":\"Банкны карт (CN)\"},\"copyText\":{},\"image\":\"\",\"externalPath\":\"\"},{\"key\":\"bank_card_us\",\"enabled\":false,\"manual\":false,\"sort\":60,\"name\":{\"zh\":\"银行卡(美国)\",\"en\":\"Bank Card (US)\",\"mn\":\"Банкны карт (US)\"},\"copyText\":{},\"image\":\"\",\"externalPath\":\"\"},{\"key\":\"bank_card_mn\",\"enabled\":false,\"manual\":false,\"sort\":70,\"name\":{\"zh\":\"银行卡(蒙古)\",\"en\":\"Bank Card (MN)\",\"mn\":\"Банкны карт (MN)\"},\"copyText\":{},\"image\":\"\",\"externalPath\":\"\"},{\"key\":\"paypal\",\"enabled\":false,\"manual\":false,\"sort\":80,\"name\":{\"zh\":\"PayPal\",\"en\":\"PayPal\",\"mn\":\"PayPal\"},\"copyText\":{},\"image\":\"\",\"externalPath\":\"\"}]", ConfigName: "支付方式配置", ConfigGroup: "payment", Remark: "支付中间页方式配置(JSON)：支持多语言名称、图片、复制信息、开关和排序"},
		{ConfigKey: "payment_uni_preferred_methods", ConfigValue: "[{\"key\":\"wechat\",\"enabled\":true,\"sort\":10,\"name\":{\"zh\":\"微信支付\",\"en\":\"WeChat Pay\",\"mn\":\"WeChat Pay\"},\"copyText\":{\"zh\":\"您好，我的订单号是 {orderID}，期望使用 {payMethod} 支付，请协助提供收款方式并处理订单。\",\"en\":\"Hi, my order number is {orderID}. I expect to pay via {payMethod}. Please provide the receiving method and help process this order.\",\"mn\":\"Сайн байна уу, миний захиалгын дугаар {orderID}. Би {payMethod} аргаар төлөхийг хүсэж байна. Хүлээн авах мэдээлэл өгч, захиалгыг боловсруулж өгнө үү.\"},\"image\":\"cloth-on/up/wechat.png\",\"externalPath\":\"\"},{\"key\":\"alipay\",\"enabled\":true,\"sort\":20,\"name\":{\"zh\":\"支付宝\",\"en\":\"Alipay\",\"mn\":\"Alipay\"},\"copyText\":{\"zh\":\"您好，我的订单号是 {orderID}，期望使用 {payMethod} 支付，请协助提供收款方式并处理订单。\",\"en\":\"Hi, my order number is {orderID}. I expect to pay via {payMethod}. Please provide the receiving method and help process this order.\",\"mn\":\"Сайн байна уу, миний захиалгын дугаар {orderID}. Би {payMethod} аргаар төлөхийг хүсэж байна. Хүлээн авах мэдээлэл өгч, захиалгыг боловсруулж өгнө үү.\"},\"image\":\"cloth-on/up/alipay.png\",\"externalPath\":\"\"},{\"key\":\"bank_card_cn\",\"enabled\":true,\"sort\":30,\"name\":{\"zh\":\"银行卡(国内)\",\"en\":\"Bank Card (CN)\",\"mn\":\"Банкны карт (CN)\"},\"copyText\":{\"zh\":\"您好，我的订单号是 {orderID}，期望使用 {payMethod} 支付，请协助提供收款方式并处理订单。\",\"en\":\"Hi, my order number is {orderID}. I expect to pay via {payMethod}. Please provide the receiving method and help process this order.\",\"mn\":\"Сайн байна уу, миний захиалгын дугаар {orderID}. Би {payMethod} аргаар төлөхийг хүсэж байна. Хүлээн авах мэдээлэл өгч, захиалгыг боловсруулж өгнө үү.\"},\"image\":\"cloth-on/up/bank-card-cn.png\",\"externalPath\":\"\"}]", ConfigName: "Uni期望支付方式配置", ConfigGroup: "payment", Remark: "仅用于uni支付页“联系客服付款”场景(JSON)：支持多语言名称、图片、复制话术、开关和排序，独立于payment_manual_methods"},
		{ConfigKey: "payment_wechat_enabled", ConfigValue: "false", ConfigName: "微信支付开关", ConfigGroup: "payment", Remark: "是否启用微信自动支付渠道(true/false)"},
		{ConfigKey: "payment_alipay_enabled", ConfigValue: "false", ConfigName: "支付宝支付开关", ConfigGroup: "payment", Remark: "是否启用支付宝自动支付渠道(true/false)"},
		{ConfigKey: "payment_bank_cn_enabled", ConfigValue: "false", ConfigName: "国内银行卡支付开关", ConfigGroup: "payment", Remark: "是否启用国内银行卡自动支付渠道(true/false)"},
		{ConfigKey: "payment_bank_us_enabled", ConfigValue: "false", ConfigName: "美国银行卡支付开关", ConfigGroup: "payment", Remark: "是否启用美国银行卡自动支付渠道(true/false)"},
		{ConfigKey: "payment_bank_mn_enabled", ConfigValue: "false", ConfigName: "蒙古国银行卡支付开关", ConfigGroup: "payment", Remark: "是否启用蒙古国银行卡自动支付渠道(true/false)"},
		{ConfigKey: "payment_paypal_enabled", ConfigValue: "false", ConfigName: "PayPal支付开关", ConfigGroup: "payment", Remark: "是否启用PayPal自动支付渠道(true/false)"},
		// system 分组
		{ConfigKey: "maintenance_enabled", ConfigValue: "false", ConfigName: "维护模式开关", ConfigGroup: "maintenance", Remark: "是否开启全站维护模式(true/false)"},
		{ConfigKey: "maintenance_message", ConfigValue: "系统维护中，请稍后再试", ConfigName: "维护提示语", ConfigGroup: "maintenance", Remark: "维护模式下的提示信息"},
		{ConfigKey: "maintenance_bg_image", ConfigValue: "", ConfigName: "维护背景图", ConfigGroup: "maintenance", Remark: "维护页面的背景图片地址"},
		{ConfigKey: "maintenance_popup_enabled", ConfigValue: "false", ConfigName: "维护弹窗开关", ConfigGroup: "maintenance", Remark: "是否显示维护弹窗(true/false)"},
		{ConfigKey: "maintenance_popup_title", ConfigValue: "{\"zh\":\"系统维护中\",\"en\":\"System Maintenance\",\"mn\":\"Систем засвар үйлчилгээнд\"}", ConfigName: "维护弹窗标题", ConfigGroup: "maintenance", Remark: "维护弹窗标题(JSON多语言)"},
		{ConfigKey: "maintenance_popup_content", ConfigValue: "{\"zh\":\"系统正在维护，请稍后再试\",\"en\":\"System is under maintenance, please try again later\",\"mn\":\"Систем засвар хийгдэж байна, дараа дахин оролдоно уу\"}", ConfigName: "维护弹窗内容", ConfigGroup: "maintenance", Remark: "维护弹窗内容(JSON多语言)"},
		{ConfigKey: "maintenance_home_btn_enabled", ConfigValue: "false", ConfigName: "进入首页按钮", ConfigGroup: "maintenance", Remark: "维护模式下是否允许进入首页(true/false)"},
		{ConfigKey: "currency_symbol", ConfigValue: "¥", ConfigName: "货币符号", ConfigGroup: "system", Remark: "前端展示的货币符号"},
		{ConfigKey: "currency_unit", ConfigValue: "CNY", ConfigName: "货币单位", ConfigGroup: "system", Remark: "货币单位编码"},
		{ConfigKey: "currency_suffix", ConfigValue: "rmb", ConfigName: "货币后缀", ConfigGroup: "system", Remark: "uni端显示的货币后缀(如rmb, usd等)"},
		{ConfigKey: "app_name", ConfigValue: "RichRoad", ConfigName: "应用名称", ConfigGroup: "system", Remark: "uni端显示的应用名称"},
		{ConfigKey: "app_logo", ConfigValue: "", ConfigName: "应用Logo", ConfigGroup: "system", Remark: "uni端显示的应用Logo图片地址"},
		// security 分组
		{ConfigKey: "captcha_expiry_seconds", ConfigValue: "300", ConfigName: "验证码有效期", ConfigGroup: "security", Remark: "图形验证码有效期(秒)，默认300秒"},
		{ConfigKey: "captcha_rate_limit", ConfigValue: "10", ConfigName: "验证码请求频率限制", ConfigGroup: "security", Remark: "每分钟每IP最多请求验证码次数"},
		{ConfigKey: "register_ip_limit", ConfigValue: "3", ConfigName: "注册IP限制", ConfigGroup: "security", Remark: "同一IP最多可注册账号数量"},
		{ConfigKey: "login_fail_max", ConfigValue: "5", ConfigName: "登录失败上限", ConfigGroup: "security", Remark: "连续登录失败N次后锁定"},
		{ConfigKey: "login_fail_wait_seconds", ConfigValue: "900", ConfigName: "登录锁定等待时间", ConfigGroup: "security", Remark: "登录锁定后等待时间(秒)，默认900秒(15分钟)"},
		// auth 分组
		{ConfigKey: "default_login_method", ConfigValue: "username", ConfigName: "默认登录方式", ConfigGroup: "auth", Remark: "默认登录方式: username(用户名+密码) 或 phone(手机号+密码)"},
		{ConfigKey: "phone_login_enabled", ConfigValue: "false", ConfigName: "手机号登录开关", ConfigGroup: "auth", Remark: "是否开启手机号+密码登录(true/false)"},
		{ConfigKey: "password_change_enabled", ConfigValue: "true", ConfigName: "密码修改开关", ConfigGroup: "auth", Remark: "uni端是否允许用户修改密码(true/false)"},
		{ConfigKey: "username_login_enabled", ConfigValue: "true", ConfigName: "用户名登录开关", ConfigGroup: "auth", Remark: "是否开启用户名+密码登录注册(true/false)"},
		{ConfigKey: "username_regex", ConfigValue: "", ConfigName: "用户名格式正则式", ConfigGroup: "auth", Remark: "用户名格式正则验证，空则不验证"},
		{ConfigKey: "password_regex", ConfigValue: "", ConfigName: "密码格式正则式", ConfigGroup: "auth", Remark: "密码格式正则验证，空则不验证"},
		{ConfigKey: "username_regex_tip", ConfigValue: "", ConfigName: "用户名提示文本", ConfigGroup: "auth", Remark: "用户名输入框内提示文本(JSON多语言)"},
		{ConfigKey: "password_regex_tip", ConfigValue: "", ConfigName: "密码提示文本", ConfigGroup: "auth", Remark: "密码输入框内提示文本(JSON多语言)"},
		// points 分组
		{ConfigKey: "points_exchange_rate", ConfigValue: "100", ConfigName: "积分兑换比率", ConfigGroup: "points", Remark: "多少积分兑换1货币单位，如100积分=1元"},
		// tryon 分组
		{ConfigKey: "tryon_guest_init_points", ConfigValue: "0", ConfigName: "游客初始试衣币", ConfigGroup: "tryon", Remark: "游客模式不再赠送试衣币，保留为兼容参数"},
		{ConfigKey: "tryon_register_reward_points", ConfigValue: "8", ConfigName: "注册奖励试衣币", ConfigGroup: "tryon", Remark: "用户注册成功后奖励的试衣币数量"},
		{ConfigKey: "tryon_invite_register_reward_points", ConfigValue: "0", ConfigName: "邀请注册奖励试衣币", ConfigGroup: "tryon", Remark: "邀请下级注册成功后奖励给邀请人的试衣币数量，0表示不赠送"},
		{ConfigKey: "tryon_cost_points", ConfigValue: "1", ConfigName: "单次试衣消耗", ConfigGroup: "tryon", Remark: "每次发起试衣或试鞋消耗的试衣币数量"},
		{ConfigKey: "tryon_fail_refund_percent", ConfigValue: "100", ConfigName: "试衣失败退币比例", ConfigGroup: "tryon", Remark: "试衣失败时退回试衣币百分比，默认100"},
		{ConfigKey: "tryon_models", ConfigValue: defaultTryonModelsSysConfigValue(), ConfigName: "试衣模型列表", ConfigGroup: "tryon", Remark: "JSON数组：可配置多模型及开关、多语言名称与说明、接口地址等"},
		{ConfigKey: "shoe_models", ConfigValue: "[{\"key\":\"aliyun_shoes_and_boots\",\"enabled\":true,\"scenes\":[\"shoes\"],\"model\":\"shoes-and-boots\",\"name\":{\"zh\":\"阿里AI试鞋\",\"en\":\"Aliyun Shoes Try-On\",\"mn\":\"Aliyun гутлын туршилт\"},\"desc\":{\"zh\":\"适用于鞋靴类虚拟试穿。\",\"en\":\"Suitable for virtual try-on of shoes and boots.\",\"mn\":\"Гутал, түрийвчний виртуал туршилтад тохиромжтой.\"},\"cost\":1,\"provider\":\"aliyun\",\"mode\":\"prod\",\"url\":\"https://dashscope.aliyuncs.com/api/v1/services/aigc/image2image/image-synthesis\",\"taskQueryUrl\":\"https://dashscope.aliyuncs.com/api/v1/tasks/{task_id}\",\"token\":\"\",\"resolution\":-1,\"restoreFace\":true}]", ConfigName: "试鞋模型列表", ConfigGroup: "tryon", Remark: "JSON数组：专用于试鞋场景的模型配置，支持多语言名称与说明、接口地址等"},
		{ConfigKey: "tryon_recharge_plans", ConfigValue: "[{\"points\":{\"zh\":\"50\",\"en\":\"50\",\"mn\":\"50\"},\"coinLabel\":{\"zh\":\"试衣币\",\"en\":\"Try-on Coins\",\"mn\":\"Туршилтын зоос\"},\"currencySymbol\":{\"zh\":\"￥\",\"en\":\"CNY \"},\"price\":{\"zh\":\"9.9\",\"en\":\"9.9\",\"mn\":\"9.9\"},\"currencySuffix\":{\"zh\":\"元\",\"en\":\"\",\"mn\":\"\"}},{\"points\":{\"zh\":\"180\",\"en\":\"180\",\"mn\":\"180\"},\"coinLabel\":{\"zh\":\"试衣币\",\"en\":\"Try-on Coins\",\"mn\":\"Туршилтын зоос\"},\"currencySymbol\":{\"zh\":\"￥\",\"en\":\"CNY \"},\"price\":{\"zh\":\"29.9\",\"en\":\"29.9\",\"mn\":\"29.9\"},\"currencySuffix\":{\"zh\":\"元\",\"en\":\"\",\"mn\":\"\"}},{\"points\":{\"zh\":\"680\",\"en\":\"680\",\"mn\":\"680\"},\"coinLabel\":{\"zh\":\"试衣币\",\"en\":\"Try-on Coins\",\"mn\":\"Туршилтын зоос\"},\"currencySymbol\":{\"zh\":\"￥\",\"en\":\"CNY \"},\"price\":{\"zh\":\"99.9\",\"en\":\"99.9\",\"mn\":\"99.9\"},\"currencySuffix\":{\"zh\":\"元\",\"en\":\"\",\"mn\":\"\"}}]", ConfigName: "试衣币充值套餐", ConfigGroup: "tryon", Remark: "后台可视化维护：点数、币名、货币符号、价格和单位均支持多语言"},
		{ConfigKey: "tryon_provider_mode", ConfigValue: "mock_success", ConfigName: "试衣模型模式", ConfigGroup: "tryon", Remark: "mock_success表示本地联调直接返回原图，prod表示调用真实模型"},
		{ConfigKey: "tryon_provider_url", ConfigValue: "", ConfigName: "试衣模型地址", ConfigGroup: "tryon", Remark: "真实模型推理服务地址(URL)"},
		{ConfigKey: "tryon_provider_token", ConfigValue: "", ConfigName: "试衣模型令牌", ConfigGroup: "tryon", Remark: "调用真实模型服务的Bearer Token"},
		{ConfigKey: "tryon_media_public_base_url", ConfigValue: "", ConfigName: "试衣图片公网前缀", ConfigGroup: "tryon", Remark: "当上传返回相对路径或localhost/内网URL时，自动拼接为公网地址，如https://back.example.com"},
		// order 分组
		{ConfigKey: "order_close_minutes", ConfigValue: "20", ConfigName: "订单自动关闭时间", ConfigGroup: "order", Remark: "未支付订单自动关闭的分钟数"},
		// display 分组
		{ConfigKey: "presale_home_count", ConfigValue: "4", ConfigName: "首页预售展示数", ConfigGroup: "display", Remark: "uni首页展示的预售商品数量"},
		{ConfigKey: "sign_in_enabled", ConfigValue: "true", ConfigName: "签到功能开关", ConfigGroup: "display", Remark: "uni端是否显示签到悬浮按钮(true/false)"},
		{ConfigKey: "shop_kefu_enabled", ConfigValue: "true", ConfigName: "外部客服展示开关", ConfigGroup: "display", Remark: "uni客服页是否展示外部客服列表(true/false)"},
		// announcement 分组
		{ConfigKey: "announcement_enabled", ConfigValue: "false", ConfigName: "公告开关", ConfigGroup: "announcement", Remark: "是否开启首页公告走马灯(true/false)"},
		{ConfigKey: "announcement_content", ConfigValue: "", ConfigName: "公告内容", ConfigGroup: "announcement", Remark: "公告走马灯文字内容，支持表情符号"},
		{ConfigKey: "announcement_text_color", ConfigValue: "#ff6600", ConfigName: "公告文字颜色", ConfigGroup: "announcement", Remark: "公告走马灯文字颜色(hex)"},
		{ConfigKey: "announcement_speed", ConfigValue: "50", ConfigName: "公告走马灯速度", ConfigGroup: "announcement", Remark: "走马灯滚动速度(px/s)"},
	}
	for _, cfg := range defaults {
		var count int64
		global.GVA_DB.Model(&client.SysConfig{}).Where("config_key = ?", cfg.ConfigKey).Count(&count)
		if count == 0 {
			if err := global.GVA_DB.Create(&cfg).Error; err != nil {
				global.GVA_LOG.Error("初始化默认系统参数失败", zap.Error(err), zap.String("configKey", cfg.ConfigKey))
			}
		}
	}
	appendIDMVTONTryonModelIfMissing()

	// 为已存在的用户生成邀请码（如果缺失）
	var users []client.ClientUser
	global.GVA_DB.Where("invite_code = '' OR invite_code IS NULL").Find(&users)
	for _, u := range users {
		b := make([]byte, 4)
		_, _ = rand.Read(b)
		code := hex.EncodeToString(b)
		global.GVA_DB.Model(&client.ClientUser{}).Where("id = ?", u.ID).Update("invite_code", code)
	}
}

func defaultTryonModelsSysConfigValue() string {
	return `[
	  {
	    "key": "aliyun_aitryon",
	    "modelUsage": "tryon",
	    "enabled": true,
	    "scenes": ["clothes", "shoes"],
	    "model": "aitryon",
	    "name": {"zh": "阿里AI试衣（基础）", "en": "Aliyun AI Try-On (Basic)", "mn": "Aliyun AI өмсгөл (Суурь)"},
	    "desc": {"zh": "基础版试衣模型，速度更快，适合日常试衣。", "en": "Basic try-on model with faster generation for everyday use.", "mn": "Өдөр тутмын туршилтад тохирох, хурдан суурь загвар."},
	    "cost": 1,
	    "provider": "aliyun",
	    "mode": "prod",
	    "url": "https://dashscope.aliyuncs.com/api/v1/services/aigc/image2image/image-synthesis",
	    "taskQueryUrl": "https://dashscope.aliyuncs.com/api/v1/tasks/{task_id}",
	    "token": "",
	    "resolution": -1,
	    "restoreFace": true,
	    "supportsRefiner": true,
	    "refinerModel": "aitryon-refiner",
	    "refinerGender": "woman",
	    "refinerExtraCost": 1,
	    "autoEnableAliyunParsingUpperOnly": true,
	    "autoEnableAliyunParsingLowerOnly": true,
	    "supportsBeautify": false
	  },
	  {
	    "key": "aliyun_aitryon_plus",
	    "modelUsage": "tryon",
	    "enabled": true,
	    "scenes": ["clothes", "shoes"],
	    "model": "aitryon-plus",
	    "name": {"zh": "阿里AI试衣（Plus）", "en": "Aliyun AI Try-On (Plus)", "mn": "Aliyun AI өмсгөл (Plus)"},
	    "desc": {"zh": "Plus版细节更好，适合高质量试衣图。", "en": "Higher quality rendering with better texture and logo details.", "mn": "Нэхмэл, логог илүү сайн сэргээдэг өндөр чанарын загвар."},
	    "cost": 1,
	    "provider": "aliyun",
	    "mode": "prod",
	    "url": "https://dashscope.aliyuncs.com/api/v1/services/aigc/image2image/image-synthesis",
	    "taskQueryUrl": "https://dashscope.aliyuncs.com/api/v1/tasks/{task_id}",
	    "token": "",
	    "resolution": -1,
	    "restoreFace": true,
	    "supportsRefiner": true,
	    "refinerModel": "aitryon-refiner",
	    "refinerGender": "woman",
	    "refinerExtraCost": 1,
	    "autoEnableAliyunParsingUpperOnly": true,
	    "autoEnableAliyunParsingLowerOnly": true,
	    "supportsBeautify": false
	  },
	  {
	    "key": "yisol_idm_vton",
	    "modelUsage": "tryon",
	    "enabled": true,
	    "scenes": ["clothes"],
	    "model": "IDM-VTON",
	    "name": {"zh": "IDM-VTON 开源试衣", "en": "IDM-VTON Open Try-On", "mn": "IDM-VTON нээлттэй өмсгөл"},
	    "desc": {"zh": "HuggingFace Space yisol/IDM-VTON，使用 Gradio /tryon 接口，支持自动蒙版与裁剪参数。", "en": "HuggingFace Space yisol/IDM-VTON via Gradio /tryon API with auto-mask and crop options.", "mn": "HuggingFace Space yisol/IDM-VTON Gradio /tryon API ашиглана."},
	    "cost": 1,
	    "provider": "gradio",
	    "mode": "prod",
	    "url": "https://yisol-idm-vton.hf.space",
	    "apiName": "/tryon",
	    "garmentDes": "clothing item",
	    "isChecked": true,
	    "isCheckedCrop": false,
	    "denoiseSteps": 30,
	    "seed": 42,
	    "token": "",
	    "resolution": -1,
	    "restoreFace": true,
	    "supportsBeautify": false
	  },
	  {
	    "key": "aliyun_aitryon_parsing",
	    "modelUsage": "tryon",
	    "enabled": false,
	    "scenes": ["takeoff"],
	    "model": "aitryon-parsing-v1",
	    "name": {"zh": "阿里取衣分割", "en": "Aliyun Takeoff Parsing", "mn": "Aliyun хувцас салгах"},
	    "desc": {"zh": "用于取衣区分割模特服饰并输出可用服饰图。", "en": "Segments garment regions for takeoff area and outputs reusable garment images.", "mn": "Загварын хувцсыг ялган авч, дахин ашиглах зургийг гаргана."},
	    "cost": 1,
	    "provider": "aliyun",
	    "mode": "prod",
	    "url": "https://dashscope.aliyuncs.com/api/v1/services/vision/image-process/process",
	    "token": "",
	    "clothesType": ["upper"],
	    "supportsBeautify": false
	  },
	  {
	    "key": "smart_beautify_default",
	    "modelUsage": "beautify",
	    "enabled": true,
	    "scenes": ["clothes", "shoes", "takeoff"],
	    "model": "custom_beautify",
	    "beautifyModel": "custom_beautify",
	    "name": {"zh": "智能美肤模型", "en": "Smart Beautify Model", "mn": "Ухаалаг арьс гоёжуулах загвар"},
	    "beautifyDesc": {"zh": "仅用于生成页对试衣结果图进行二次美肤。", "en": "Used only on generate page to beautify try-on results.", "mn": "Зөвхөн туршилтын үр дүнгийн зургийг арьс сайжруулахад ашиглана."},
	    "cost": 0,
	    "beautifyExtraCost": 0,
	    "provider": "custom",
	    "mode": "prod",
	    "url": "",
	    "token": "",
	    "supportsBeautify": true,
	    "beautifyRetouchDegree": 70,
	    "beautifyWhiteningDegree": 30
	  }
	]`
}

func appendIDMVTONTryonModelIfMissing() {
	var cfg client.SysConfig
	if err := global.GVA_DB.Where("config_key = ?", "tryon_models").First(&cfg).Error; err != nil {
		return
	}

	models := make([]map[string]interface{}, 0)
	if err := json.Unmarshal([]byte(cfg.ConfigValue), &models); err != nil {
		global.GVA_LOG.Warn("试衣模型配置解析失败，跳过IDM-VTON自动追加", zap.Error(err))
		return
	}

	changed := false
	hasIDM := false
	hasBeautify := false
	for i := range models {
		model := models[i]
		key := strings.ToLower(strings.TrimSpace(toString(model["key"])))
		usage := strings.ToLower(strings.TrimSpace(toString(model["modelUsage"])))
		if usage == "" {
			usage = "tryon"
		}

		if key == "yisol_idm_vton" {
			hasIDM = true
		}

		if usage == "beautify" {
			hasBeautify = true
			if !toBoolLoose(model["supportsBeautify"], true) {
				model["supportsBeautify"] = true
				changed = true
			}
			if strings.TrimSpace(toString(model["beautifyModel"])) == "" {
				model["beautifyModel"] = "custom_beautify"
				changed = true
			}
			if strings.TrimSpace(toString(model["model"])) == "" {
				model["model"] = "custom_beautify"
				changed = true
			}
			continue
		}

		if strings.TrimSpace(toString(model["beautifyModelKey"])) != "" {
			model["beautifyModelKey"] = ""
			changed = true
		}
		if toBoolLoose(model["supportsBeautify"], false) {
			model["supportsBeautify"] = false
			changed = true
		}

		if isAliyunAitryonSeriesLoose(model) {
			if _, ok := model["autoEnableAliyunParsingUpperOnly"]; !ok {
				model["autoEnableAliyunParsingUpperOnly"] = true
				changed = true
			}
			if _, ok := model["autoEnableAliyunParsingLowerOnly"]; !ok {
				model["autoEnableAliyunParsingLowerOnly"] = true
				changed = true
			}
		}
	}

	if !hasIDM {
		models = append(models, map[string]interface{}{
			"key":              "yisol_idm_vton",
			"modelUsage":       "tryon",
			"enabled":          true,
			"scenes":           []string{"clothes"},
			"model":            "IDM-VTON",
			"name":             map[string]string{"zh": "IDM-VTON 开源试衣", "en": "IDM-VTON Open Try-On", "mn": "IDM-VTON нээлттэй өмсгөл"},
			"desc":             map[string]string{"zh": "HuggingFace Space yisol/IDM-VTON，使用 Gradio /tryon 接口，支持自动蒙版与裁剪参数。", "en": "HuggingFace Space yisol/IDM-VTON via Gradio /tryon API with auto-mask and crop options.", "mn": "HuggingFace Space yisol/IDM-VTON Gradio /tryon API ашиглана."},
			"cost":             1,
			"provider":         "gradio",
			"mode":             "prod",
			"url":              "https://yisol-idm-vton.hf.space",
			"apiName":          "/tryon",
			"garmentDes":       "clothing item",
			"isChecked":        true,
			"isCheckedCrop":    false,
			"denoiseSteps":     30,
			"seed":             42,
			"token":            "",
			"resolution":       -1,
			"restoreFace":      true,
			"supportsBeautify": false,
		})
		changed = true
	}

	if !hasBeautify {
		models = append(models, map[string]interface{}{
			"key":                     "smart_beautify_default",
			"modelUsage":              "beautify",
			"enabled":                 true,
			"scenes":                  []string{"clothes", "shoes", "takeoff"},
			"model":                   "custom_beautify",
			"beautifyModel":           "custom_beautify",
			"name":                    map[string]string{"zh": "智能美肤模型", "en": "Smart Beautify Model", "mn": "Ухаалаг арьс гоёжуулах загвар"},
			"beautifyDesc":            map[string]string{"zh": "仅用于生成页对试衣结果图进行二次美肤。", "en": "Used only on generate page to beautify try-on results.", "mn": "Зөвхөн туршилтын үр дүнгийн зургийг арьс сайжруулахад ашиглана."},
			"cost":                    0,
			"beautifyExtraCost":       0,
			"provider":                "custom",
			"mode":                    "prod",
			"url":                     "",
			"token":                   "",
			"supportsBeautify":        true,
			"beautifyRetouchDegree":   70,
			"beautifyWhiteningDegree": 30,
		})
		changed = true
	}

	if !changed {
		return
	}

	encoded, err := json.Marshal(models)
	if err != nil {
		global.GVA_LOG.Warn("试衣模型配置编码失败，跳过自动修复", zap.Error(err))
		return
	}

	if err := global.GVA_DB.Model(&client.SysConfig{}).Where("id = ?", cfg.ID).Update("config_value", string(encoded)).Error; err != nil {
		global.GVA_LOG.Warn("试衣模型配置自动修复失败", zap.Error(err))
	}
}

func toBoolLoose(value interface{}, fallback bool) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		text := strings.ToLower(strings.TrimSpace(v))
		switch text {
		case "1", "true", "yes", "on":
			return true
		case "0", "false", "no", "off":
			return false
		}
	case float64:
		return v != 0
	case float32:
		return v != 0
	case int:
		return v != 0
	case int8:
		return v != 0
	case int16:
		return v != 0
	case int32:
		return v != 0
	case int64:
		return v != 0
	case uint:
		return v != 0
	case uint8:
		return v != 0
	case uint16:
		return v != 0
	case uint32:
		return v != 0
	case uint64:
		return v != 0
	}
	return fallback
}

func isAliyunAitryonSeriesLoose(model map[string]interface{}) bool {
	if model == nil {
		return false
	}
	modelName := strings.ToLower(strings.TrimSpace(toString(model["model"])))
	if modelName == "aitryon" || modelName == "aitryon-plus" {
		return true
	}

	key := strings.ToLower(strings.TrimSpace(toString(model["key"])))
	provider := strings.ToLower(strings.TrimSpace(toString(model["provider"])))
	if strings.Contains(key, "aliyun_aitryon") {
		return true
	}
	if (strings.Contains(provider, "aliyun") || strings.Contains(provider, "dashscope")) && strings.HasPrefix(modelName, "aitryon") {
		return true
	}

	return false
}

func toString(value interface{}) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(value))
}
