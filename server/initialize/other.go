package initialize

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

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
		{ConfigKey: "invite_share_link_tip_text", ConfigValue: "{\"zh\":\"将分享链接发送给好友，好友注册后你将获得试衣币奖励。\",\"en\":\"Share the link with friends. You will receive try-on coin rewards after they register.\",\"mn\":\"Урилгын холбоосоо найзууддаа илгээснээр тэд бүртгүүлэхэд та туршилтын зоосны урамшуулал авна.\"}", ConfigName: "邀请分享提示文案", ConfigGroup: "invite", Remark: "邀请页分享链接下方提示文案(JSON多语言)"},
		{ConfigKey: "invite_share_link_tip_text_color", ConfigValue: "#475569", ConfigName: "邀请分享提示颜色", ConfigGroup: "invite", Remark: "邀请页分享链接下方提示文字颜色(hex)"},
		// payment 分组
		{ConfigKey: "payment_tip_text", ConfigValue: "{\"zh\":\"请在规定时间内完成付款\",\"en\":\"Please complete payment within the specified time\",\"mn\":\"Заасан хугацаанд төлбөрөө хийнэ үү\"}", ConfigName: "付款提示文本", ConfigGroup: "payment", Remark: "二维码付款弹窗提示文本(JSON多语言)"},
		{ConfigKey: "payment_tip_text_size", ConfigValue: "14", ConfigName: "付款提示文字大小", ConfigGroup: "payment", Remark: "付款提示文本字体大小(px)"},
		{ConfigKey: "payment_tip_text_color", ConfigValue: "#ff0000", ConfigName: "付款提示文字颜色", ConfigGroup: "payment", Remark: "付款提示文本颜色(hex)"},
		{ConfigKey: "payment_auto_enabled", ConfigValue: "false", ConfigName: "自动支付总开关", ConfigGroup: "payment", Remark: "是否启用自动支付渠道(true/false)，默认关闭保留人工流程"},
		{ConfigKey: "payment_manual_qrcode_enabled", ConfigValue: "true", ConfigName: "人工二维码支付开关", ConfigGroup: "payment", Remark: "是否展示人工二维码付款(true/false)"},
		{ConfigKey: "payment_manual_contact_enabled", ConfigValue: "true", ConfigName: "人工联系客服支付开关", ConfigGroup: "payment", Remark: "是否展示联系客服付款(true/false)"},
		{ConfigKey: "payment_manual_methods", ConfigValue: "[{\"key\":\"qrcode\",\"enabled\":true,\"manual\":true,\"sort\":10,\"name\":{\"zh\":\"二维码支付\",\"en\":\"QR Payment\",\"mn\":\"QR төлбөр\"},\"copyText\":{\"zh\":\"请扫码后备注订单号\",\"en\":\"Please include your order number after payment\",\"mn\":\"Төлбөр хийхдээ захиалгын дугаараа тэмдэглэнэ үү\"},\"image\":\"\"},{\"key\":\"contact\",\"enabled\":true,\"manual\":true,\"sort\":20,\"name\":{\"zh\":\"联系客服\",\"en\":\"Contact Support\",\"mn\":\"Хэрэглэгчийн дэмжлэг\"},\"copyText\":{\"zh\":\"付款后请联系客服并提供订单号\",\"en\":\"Please contact support with your order number after payment\",\"mn\":\"Төлбөр хийсний дараа захиалгын дугаартайгаа客服-т хандана уу\"},\"image\":\"\"},{\"key\":\"wechat\",\"enabled\":false,\"manual\":false,\"sort\":30,\"name\":{\"zh\":\"微信支付\",\"en\":\"WeChat Pay\",\"mn\":\"WeChat Pay\"},\"copyText\":{},\"image\":\"\"},{\"key\":\"alipay\",\"enabled\":false,\"manual\":false,\"sort\":40,\"name\":{\"zh\":\"支付宝\",\"en\":\"Alipay\",\"mn\":\"Alipay\"},\"copyText\":{},\"image\":\"\"},{\"key\":\"bank_card_cn\",\"enabled\":false,\"manual\":false,\"sort\":50,\"name\":{\"zh\":\"银行卡(国内)\",\"en\":\"Bank Card (CN)\",\"mn\":\"Банкны карт (CN)\"},\"copyText\":{},\"image\":\"\"},{\"key\":\"bank_card_us\",\"enabled\":false,\"manual\":false,\"sort\":60,\"name\":{\"zh\":\"银行卡(美国)\",\"en\":\"Bank Card (US)\",\"mn\":\"Банкны карт (US)\"},\"copyText\":{},\"image\":\"\"},{\"key\":\"bank_card_mn\",\"enabled\":false,\"manual\":false,\"sort\":70,\"name\":{\"zh\":\"银行卡(蒙古)\",\"en\":\"Bank Card (MN)\",\"mn\":\"Банкны карт (MN)\"},\"copyText\":{},\"image\":\"\"},{\"key\":\"paypal\",\"enabled\":false,\"manual\":false,\"sort\":80,\"name\":{\"zh\":\"PayPal\",\"en\":\"PayPal\",\"mn\":\"PayPal\"},\"copyText\":{},\"image\":\"\"}]", ConfigName: "支付方式配置", ConfigGroup: "payment", Remark: "支付中间页方式配置(JSON)：支持多语言名称、图片、复制信息、开关和排序"},
		{ConfigKey: "payment_uni_preferred_methods", ConfigValue: "[{\"key\":\"wechat\",\"enabled\":true,\"sort\":10,\"name\":{\"zh\":\"微信支付\",\"en\":\"WeChat Pay\",\"mn\":\"WeChat Pay\"},\"copyText\":{\"zh\":\"您好，我的订单号是 {orderID}，期望使用 {payMethod} 支付，请协助提供收款方式并处理订单。\",\"en\":\"Hi, my order number is {orderID}. I expect to pay via {payMethod}. Please provide the receiving method and help process this order.\",\"mn\":\"Сайн байна уу, миний захиалгын дугаар {orderID}. Би {payMethod} аргаар төлөхийг хүсэж байна. Хүлээн авах мэдээлэл өгч, захиалгыг боловсруулж өгнө үү.\"},\"image\":\"cloth-on/web-else/hope-pay/wechat.png\"},{\"key\":\"alipay\",\"enabled\":true,\"sort\":20,\"name\":{\"zh\":\"支付宝\",\"en\":\"Alipay\",\"mn\":\"Alipay\"},\"copyText\":{\"zh\":\"您好，我的订单号是 {orderID}，期望使用 {payMethod} 支付，请协助提供收款方式并处理订单。\",\"en\":\"Hi, my order number is {orderID}. I expect to pay via {payMethod}. Please provide the receiving method and help process this order.\",\"mn\":\"Сайн байна уу, миний захиалгын дугаар {orderID}. Би {payMethod} аргаар төлөхийг хүсэж байна. Хүлээн авах мэдээлэл өгч, захиалгыг боловсруулж өгнө үү.\"},\"image\":\"cloth-on/web-else/hope-pay/alipay.png\"},{\"key\":\"bank_card_cn\",\"enabled\":true,\"sort\":30,\"name\":{\"zh\":\"银行卡(国内)\",\"en\":\"Bank Card (CN)\",\"mn\":\"Банкны карт (CN)\"},\"copyText\":{\"zh\":\"您好，我的订单号是 {orderID}，期望使用 {payMethod} 支付，请协助提供收款方式并处理订单。\",\"en\":\"Hi, my order number is {orderID}. I expect to pay via {payMethod}. Please provide the receiving method and help process this order.\",\"mn\":\"Сайн байна уу, миний захиалгын дугаар {orderID}. Би {payMethod} аргаар төлөхийг хүсэж байна. Хүлээн авах мэдээлэл өгч, захиалгыг боловсруулж өгнө үү.\"},\"image\":\"cloth-on/web-else/hope-pay/bank-card-cn.png\"}]", ConfigName: "Uni期望支付方式配置", ConfigGroup: "payment", Remark: "仅用于uni支付页“联系客服付款”场景(JSON)：支持多语言名称、图片、复制话术、开关和排序，独立于payment_manual_methods"},
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
		{ConfigKey: "app_name", ConfigValue: "{\"zh\":\"RichRoad\",\"en\":\"RichRoad\",\"mn\":\"RichRoad\"}", ConfigName: "应用名称", ConfigGroup: "system", Remark: "uni端显示的应用名称(JSON多语言)"},
		{ConfigKey: "app_logo", ConfigValue: "", ConfigName: "应用Logo", ConfigGroup: "system", Remark: "uni端显示的应用Logo图片地址"},
		{ConfigKey: "language_translate_timeout_ms", ConfigValue: "45000", ConfigName: "翻译超时(毫秒)", ConfigGroup: "system", Remark: "语言自动翻译接口请求超时，范围5000-180000"},
		{ConfigKey: "language_translate_retry_count", ConfigValue: "2", ConfigName: "翻译重试次数", ConfigGroup: "system", Remark: "Google翻译每个网关的重试次数，范围1-5"},
		{ConfigKey: "language_translate_endpoints", ConfigValue: "https://translate.googleapis.com/translate_a/single,https://translate.google.com/translate_a/single,https://translate.google.com.hk/translate_a/single", ConfigName: "翻译网关列表", ConfigGroup: "system", Remark: "逗号分隔的翻译网关URL列表，按顺序回退"},
		{ConfigKey: "language_translate_provider_url", ConfigValue: "", ConfigName: "外部翻译服务地址", ConfigGroup: "system", Remark: "非Google翻译回退服务地址（可接入自建服务或兼容网关）"},
		{ConfigKey: "language_translate_provider_api_key", ConfigValue: "", ConfigName: "外部翻译服务密钥", ConfigGroup: "system", Remark: "外部翻译服务鉴权密钥（可空）"},
		{ConfigKey: "language_translate_provider_api_key_header", ConfigValue: "Authorization", ConfigName: "外部翻译密钥请求头", ConfigGroup: "system", Remark: "外部翻译服务密钥请求头名称，默认Authorization"},
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
		// english_learning 分组
		{ConfigKey: "learning_daily_target", ConfigValue: "10", ConfigName: "英语学习每日目标", ConfigGroup: "english_learning", Remark: "学习首页每日目标展示值"},
		{ConfigKey: "learning_new_user_free_hours", ConfigValue: "24", ConfigName: "英语学习新用户全站免费时长(小时)", ConfigGroup: "english_learning", Remark: "首次创建学习资产时自动发放的全站免费权益时长，0表示关闭"},
		{ConfigKey: "learning_checkin_base_point", ConfigValue: "100", ConfigName: "英语签到基础积分", ConfigGroup: "english_learning", Remark: "英语学习签到第1天基础积分"},
		{ConfigKey: "learning_checkin_increment", ConfigValue: "5", ConfigName: "英语签到每日递增", ConfigGroup: "english_learning", Remark: "英语学习签到每天递增积分"},
		{ConfigKey: "learning_checkin_cycle_days", ConfigValue: "10", ConfigName: "英语签到轮回天数", ConfigGroup: "english_learning", Remark: "英语学习签到积分轮回周期天数"},
		{ConfigKey: "learning_tts_enabled", ConfigValue: "false", ConfigName: "英语TTS自动发音开关", ConfigGroup: "english_learning", Remark: "是否启用英语单词自动发音生成(true/false)"},
		{ConfigKey: "learning_tts_provider_url", ConfigValue: "", ConfigName: "英语TTS服务地址", ConfigGroup: "english_learning", Remark: "外部TTS服务HTTP地址，需返回可访问音频URL"},
		{ConfigKey: "learning_tts_api_key", ConfigValue: "", ConfigName: "英语TTS服务密钥", ConfigGroup: "english_learning", Remark: "外部TTS服务鉴权密钥（可空）"},
		{ConfigKey: "learning_tts_timeout_ms", ConfigValue: "8000", ConfigName: "英语TTS请求超时(毫秒)", ConfigGroup: "english_learning", Remark: "调用外部TTS服务的超时时间，范围1000-60000"},
		{ConfigKey: "learning_tts_voice_us", ConfigValue: "en-US-JennyNeural", ConfigName: "英语TTS美式音色", ConfigGroup: "english_learning", Remark: "自动生成美式发音使用的voice参数"},
		{ConfigKey: "learning_tts_voice_uk", ConfigValue: "en-GB-SoniaNeural", ConfigName: "英语TTS英式音色", ConfigGroup: "english_learning", Remark: "自动生成英式发音使用的voice参数"},
		{ConfigKey: "learning_api_encrypt_enabled", ConfigValue: "true", ConfigName: "英语学习UNI接口响应加密开关", ConfigGroup: "english_learning", Remark: "是否启用英语学习UNI接口JSON响应加密(true/false)，便于调试可关闭"},
		{ConfigKey: "learning_api_sign_enabled", ConfigValue: "true", ConfigName: "英语学习UNI接口签名校验开关", ConfigGroup: "english_learning", Remark: "是否启用英语学习UNI接口请求签名校验(true/false)，便于调试可关闭"},
		// points 分组
		{ConfigKey: "points_exchange_rate", ConfigValue: "100", ConfigName: "积分兑换比率", ConfigGroup: "points", Remark: "多少积分兑换1货币单位，如100积分=1元"},
		// tryon 分组
		{ConfigKey: "tryon_guest_init_points", ConfigValue: "0", ConfigName: "游客初始试衣币", ConfigGroup: "tryon", Remark: "游客模式不再赠送试衣币，保留为兼容参数"},
		{ConfigKey: "tryon_register_reward_points", ConfigValue: "8", ConfigName: "注册奖励试衣币", ConfigGroup: "tryon", Remark: "用户注册成功后奖励的试衣币数量"},
		{ConfigKey: "tryon_invite_register_reward_points", ConfigValue: "0", ConfigName: "邀请注册奖励试衣币", ConfigGroup: "tryon", Remark: "邀请下级注册成功后奖励给邀请人的试衣币数量，0表示不赠送"},
		{ConfigKey: "tryon_cost_points", ConfigValue: "1", ConfigName: "单次试衣消耗", ConfigGroup: "tryon", Remark: "每次发起试衣或试鞋消耗的试衣币数量"},
		{ConfigKey: "tryon_fail_refund_percent", ConfigValue: "100", ConfigName: "试衣失败退币比例", ConfigGroup: "tryon", Remark: "试衣失败时退回试衣币百分比，默认100"},
		{ConfigKey: "tryon_append_parsing_failed_tip", ConfigValue: "false", ConfigName: "分割失败提示拼接开关", ConfigGroup: "tryon", Remark: "试衣成功但分割增强失败时，是否在客户端拼接分割失败提示(true/false)"},
		{ConfigKey: "tryon_parsing_failed_tip_text", ConfigValue: "{\"zh\":\"试衣成功，但分割增强失败，相关金币已退回\",\"en\":\"Try-on succeeded, but parsing enhancement failed. Related coins have been refunded.\",\"mn\":\"Туршилт амжилттай боловч segmentation enhancement амжилтгүй боллоо. Холбогдох зоос буцаан олгогдлоо.\"}", ConfigName: "分割失败提示文本", ConfigGroup: "tryon", Remark: "试衣成功但分割增强失败时拼接的提示(JSON多语言)"},
		{ConfigKey: "tryon_append_refiner_failed_tip", ConfigValue: "true", ConfigName: "精修失败提示拼接开关", ConfigGroup: "tryon", Remark: "试衣成功且精修失败时，是否在客户端拼接“金币已退回”提示(true/false)"},
		{ConfigKey: "tryon_refiner_failed_tip_text", ConfigValue: "{\"zh\":\"试衣成功，但精修失败，金币已退回\",\"en\":\"Try-on succeeded, but refiner failed. Coins have been refunded.\",\"mn\":\"Туршилт амжилттай боловч нарийвчлал амжилтгүй боллоо. Зоос буцаан олгогдсон.\"}", ConfigName: "精修失败提示文本", ConfigGroup: "tryon", Remark: "试衣成功但精修失败时拼接的提示(JSON多语言)"},
		{ConfigKey: "tryon_models", ConfigValue: defaultTryonModelsSysConfigValue(), ConfigName: "试衣模型列表", ConfigGroup: "tryon", Remark: "JSON数组：统一承载试衣/试鞋/取衣/美肤模型配置，按模型独立设置参数"},
		{ConfigKey: "tryon_recharge_plans", ConfigValue: "[{\"points\":{\"zh\":\"50\",\"en\":\"50\",\"mn\":\"50\"},\"coinLabel\":{\"zh\":\"试衣币\",\"en\":\"Try-on Coins\",\"mn\":\"Туршилтын зоос\"},\"price\":990,\"priceI18n\":{\"zh\":990,\"en\":990,\"mn\":990}},{\"points\":{\"zh\":\"180\",\"en\":\"180\",\"mn\":\"180\"},\"coinLabel\":{\"zh\":\"试衣币\",\"en\":\"Try-on Coins\",\"mn\":\"Туршилтын зоос\"},\"price\":2990,\"priceI18n\":{\"zh\":2990,\"en\":2990,\"mn\":2990}},{\"points\":{\"zh\":\"680\",\"en\":\"680\",\"mn\":\"680\"},\"coinLabel\":{\"zh\":\"试衣币\",\"en\":\"Try-on Coins\",\"mn\":\"Туршилтын зоос\"},\"price\":9990,\"priceI18n\":{\"zh\":9990,\"en\":9990,\"mn\":9990}}]", ConfigName: "试衣币充值套餐", ConfigGroup: "tryon", Remark: "后台可视化维护：点数、币名与多语言价格（基础分价 + 汇率换算）"},
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
	purgeDeprecatedTryonSysConfigs()
	appendIDMVTONTryonModelIfMissing()
	normalizeTryonModelsSysConfig()
	normalizeLocalizedSysConfigDefaults()

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

func normalizeLocalizedSysConfigDefaults() {
	normalizeConfigValueToMultilingualObject("app_name")
}

func normalizeConfigValueToMultilingualObject(configKey string) {
	if strings.TrimSpace(configKey) == "" {
		return
	}

	var cfg client.SysConfig
	if err := global.GVA_DB.Where("config_key = ?", configKey).First(&cfg).Error; err != nil {
		return
	}

	rawValue := strings.TrimSpace(cfg.ConfigValue)
	if rawValue == "" {
		return
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(rawValue), &obj); err == nil {
		if _, hasZh := obj["zh"]; hasZh {
			return
		}
		if _, hasEn := obj["en"]; hasEn {
			return
		}
		if _, hasMn := obj["mn"]; hasMn {
			return
		}
	}

	normalized := map[string]string{
		"zh": rawValue,
		"en": rawValue,
		"mn": rawValue,
	}
	payload, err := json.Marshal(normalized)
	if err != nil {
		return
	}

	if err := global.GVA_DB.Model(&client.SysConfig{}).
		Where("id = ?", cfg.ID).
		Update("config_value", string(payload)).Error; err != nil {
		global.GVA_LOG.Warn("规范化系统参数多语言值失败", zap.String("configKey", configKey), zap.Error(err))
	}
}

func purgeDeprecatedTryonSysConfigs() {
	deprecatedKeys := []string{
		"shoe_models",
		"tryon_provider_mode",
		"tryon_provider_url",
		"tryon_provider_token",
	}
	if err := global.GVA_DB.Where("config_key IN ?", deprecatedKeys).Delete(&client.SysConfig{}).Error; err != nil {
		global.GVA_LOG.Warn("清理试衣历史参数失败", zap.Error(err))
	}
}

func defaultTryonModelsSysConfigValue() string {
	return `[
	  {
	    "key": "aliyun_aitryon",
	    "modelUsage": "tryon",
	    "enabled": true,
	    "scenes": ["clothes"],
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
	    "refinerModelKey": "aliyun_aitryon_refiner",
	    "parsingModelKey": "aliyun_aitryon_parsing"
	  },
	  {
	    "key": "aliyun_aitryon_plus",
	    "modelUsage": "tryon",
	    "enabled": true,
	    "scenes": ["clothes"],
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
	    "refinerModelKey": "aliyun_aitryon_refiner",
	    "parsingModelKey": "aliyun_aitryon_parsing"
	  },
	  {
	    "key": "aliyun_shoemodel_v1",
	    "modelUsage": "tryon",
	    "enabled": true,
	    "scenes": ["shoes"],
	    "model": "shoemodel-v1",
	    "name": {"zh": "阿里鞋靴模特", "en": "Aliyun Shoes Virtual Model", "mn": "Aliyun гутлын загвар"},
	    "desc": {"zh": "阿里鞋靴模特专用模型，仅用于试鞋间生成。", "en": "Dedicated Aliyun shoes model for shoe try-on room.", "mn": "Гутлын туршилтын өрөөнд ашиглах тусгай Aliyun загвар."},
	    "cost": 1,
	    "provider": "aliyun",
	    "mode": "prod",
	    "url": "https://dashscope.aliyuncs.com/api/v1/services/aigc/virtualmodel/generation",
	    "taskQueryUrl": "https://dashscope.aliyuncs.com/api/v1/tasks/{task_id}",
	    "token": ""
	  },
	  {
	    "key": "aliyun_aitryon_refiner",
	    "modelUsage": "refiner",
	    "enabled": true,
	    "scenes": ["clothes", "shoes"],
	    "model": "aitryon-refiner",
	    "name": {"zh": "阿里图片精修", "en": "Aliyun Image Refiner", "mn": "Aliyun зураг нарийвчлал"},
	    "desc": {"zh": "独立的图片精修模型，可被试衣模型按 key 引用。", "en": "Independent image refiner model that can be referenced by try-on models via key.", "mn": "Туршилтын загварууд key-ээр зааж ашиглах бие даасан нарийвчлалын загвар."},
	    "cost": 1,
	    "refinerExtraCost": 1,
	    "provider": "aliyun",
	    "mode": "prod",
	    "url": "https://dashscope.aliyuncs.com/api/v1/services/aigc/image2image/image-synthesis",
	    "taskQueryUrl": "https://dashscope.aliyuncs.com/api/v1/tasks/{task_id}",
	    "token": "",
	    "refinerGender": "woman"
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
	    "restoreFace": true
	  },
	  {
	    "key": "aliyun_aitryon_parsing",
	    "modelUsage": "parsing",
	    "enabled": false,
	    "scenes": ["takeoff"],
	    "model": "aitryon-parsing-v1",
	    "name": {"zh": "阿里取衣分割", "en": "Aliyun Takeoff Parsing", "mn": "Aliyun хувцас салгах"},
	    "desc": {"zh": "用于取衣区分割模特服饰并输出可用服饰图。", "en": "Segments garment regions for takeoff area and outputs reusable garment images.", "mn": "Загварын хувцсыг ялган авч, дахин ашиглах зургийг гаргана."},
	    "cost": 0,
	    "parsingExtraCost": 0,
	    "provider": "aliyun",
	    "mode": "prod",
	    "url": "https://dashscope.aliyuncs.com/api/v1/services/vision/image-process/process",
	    "token": "",
	    "clothesType": ["upper"]
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
	    "token": ""
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
	hasRefiner := false
	hasParsing := false
	hasShoeTryon := false
	for i := range models {
		model := models[i]
		key := strings.ToLower(strings.TrimSpace(toString(model["key"])))
		usage := strings.ToLower(strings.TrimSpace(toString(model["modelUsage"])))
		modelName := strings.ToLower(strings.TrimSpace(toString(model["model"])))
		if usage == "" {
			usage = "tryon"
		}

		if key == "yisol_idm_vton" {
			hasIDM = true
		}
		if key == "aliyun_aitryon_refiner" || usage == "refiner" {
			hasRefiner = true
		}
		if key == "aliyun_aitryon_parsing" || usage == "parsing" {
			hasParsing = true
		}
		if key == "aliyun_shoemodel_v1" || modelName == "shoemodel-v1" {
			hasShoeTryon = true
		}

		if key == "aliyun_aitryon" || key == "aliyun_aitryon_plus" {
			scenes := toStringArrayLoose(model["scenes"])
			if len(scenes) != 1 || !strings.EqualFold(strings.TrimSpace(scenes[0]), "clothes") {
				model["scenes"] = []string{"clothes"}
				changed = true
			}
		}

		if usage == "beautify" {
			hasBeautify = true
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

	}

	if !hasShoeTryon {
		models = append(models, map[string]interface{}{
			"key":          "aliyun_shoemodel_v1",
			"modelUsage":   "tryon",
			"enabled":      true,
			"scenes":       []string{"shoes"},
			"model":        "shoemodel-v1",
			"name":         map[string]string{"zh": "阿里鞋靴模特", "en": "Aliyun Shoes Virtual Model", "mn": "Aliyun гутлын загвар"},
			"desc":         map[string]string{"zh": "阿里鞋靴模特专用模型，仅用于试鞋间生成。", "en": "Dedicated Aliyun shoes model for shoe try-on room.", "mn": "Гутлын туршилтын өрөөнд ашиглах тусгай Aliyun загвар."},
			"cost":         1,
			"provider":     "aliyun",
			"mode":         "prod",
			"url":          "https://dashscope.aliyuncs.com/api/v1/services/aigc/virtualmodel/generation",
			"taskQueryUrl": "https://dashscope.aliyuncs.com/api/v1/tasks/{task_id}",
			"token":        "",
		})
		changed = true
	}

	if !hasIDM {
		models = append(models, map[string]interface{}{
			"key":           "yisol_idm_vton",
			"modelUsage":    "tryon",
			"enabled":       true,
			"scenes":        []string{"clothes"},
			"model":         "IDM-VTON",
			"name":          map[string]string{"zh": "IDM-VTON 开源试衣", "en": "IDM-VTON Open Try-On", "mn": "IDM-VTON нээлттэй өмсгөл"},
			"desc":          map[string]string{"zh": "HuggingFace Space yisol/IDM-VTON，使用 Gradio /tryon 接口，支持自动蒙版与裁剪参数。", "en": "HuggingFace Space yisol/IDM-VTON via Gradio /tryon API with auto-mask and crop options.", "mn": "HuggingFace Space yisol/IDM-VTON Gradio /tryon API ашиглана."},
			"cost":          1,
			"provider":      "gradio",
			"mode":          "prod",
			"url":           "https://yisol-idm-vton.hf.space",
			"apiName":       "/tryon",
			"garmentDes":    "clothing item",
			"isChecked":     true,
			"isCheckedCrop": false,
			"denoiseSteps":  30,
			"seed":          42,
			"token":         "",
			"resolution":    -1,
			"restoreFace":   true,
		})
		changed = true
	}

	if !hasBeautify {
		models = append(models, map[string]interface{}{
			"key":               "smart_beautify_default",
			"modelUsage":        "beautify",
			"enabled":           true,
			"scenes":            []string{"clothes", "shoes", "takeoff"},
			"model":             "custom_beautify",
			"beautifyModel":     "custom_beautify",
			"name":              map[string]string{"zh": "智能美肤模型", "en": "Smart Beautify Model", "mn": "Ухаалаг арьс гоёжуулах загвар"},
			"beautifyDesc":      map[string]string{"zh": "仅用于生成页对试衣结果图进行二次美肤。", "en": "Used only on generate page to beautify try-on results.", "mn": "Зөвхөн туршилтын үр дүнгийн зургийг арьс сайжруулахад ашиглана."},
			"cost":              0,
			"beautifyExtraCost": 0,
			"provider":          "custom",
			"mode":              "prod",
			"url":               "",
			"token":             "",
		})
		changed = true
	}

	if !hasRefiner {
		models = append(models, map[string]interface{}{
			"key":              "aliyun_aitryon_refiner",
			"modelUsage":       "refiner",
			"enabled":          true,
			"scenes":           []string{"clothes", "shoes"},
			"model":            "aitryon-refiner",
			"name":             map[string]string{"zh": "阿里图片精修", "en": "Aliyun Image Refiner", "mn": "Aliyun зураг нарийвчлал"},
			"desc":             map[string]string{"zh": "独立的图片精修模型，可被试衣模型按 key 引用。", "en": "Independent image refiner model referenced by try-on models.", "mn": "Туршилтын загваруудаас key-ээр дуудагдах бие даасан нарийвчлалын загвар."},
			"cost":             1,
			"refinerExtraCost": 1,
			"provider":         "aliyun",
			"mode":             "prod",
			"url":              "https://dashscope.aliyuncs.com/api/v1/services/aigc/image2image/image-synthesis",
			"taskQueryUrl":     "https://dashscope.aliyuncs.com/api/v1/tasks/{task_id}",
			"token":            "",
			"refinerGender":    "woman",
		})
		changed = true
	}

	if !hasParsing {
		models = append(models, map[string]interface{}{
			"key":              "aliyun_aitryon_parsing",
			"modelUsage":       "parsing",
			"enabled":          false,
			"scenes":           []string{"takeoff"},
			"model":            "aitryon-parsing-v1",
			"name":             map[string]string{"zh": "阿里取衣分割", "en": "Aliyun Takeoff Parsing", "mn": "Aliyun хувцас салгах"},
			"desc":             map[string]string{"zh": "用于取衣区分割模特服饰并输出可用服饰图。", "en": "Segments garment regions for takeoff area and outputs reusable garment images.", "mn": "Загварын хувцсыг ялган авч, дахин ашиглах зургийг гаргана."},
			"cost":             0,
			"parsingExtraCost": 0,
			"provider":         "aliyun",
			"mode":             "prod",
			"url":              "https://dashscope.aliyuncs.com/api/v1/services/vision/image-process/process",
			"token":            "",
			"clothesType":      []string{"upper"},
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

func normalizeTryonModelsSysConfig() {
	startedAt := time.Now()

	var cfg client.SysConfig
	if err := global.GVA_DB.Where("config_key = ?", "tryon_models").First(&cfg).Error; err != nil {
		return
	}

	models := make([]map[string]interface{}, 0)
	if err := json.Unmarshal([]byte(cfg.ConfigValue), &models); err != nil {
		global.GVA_LOG.Warn("试衣模型配置解析失败，跳过存量规范化", zap.Error(err))
		return
	}

	normalizedModels := make([]map[string]interface{}, 0, len(models))
	for i := range models {
		normalizedModels = append(normalizedModels, normalizeTryonModelConfigItem(models[i], i))
	}

	currentEncoded, err := json.Marshal(models)
	if err != nil {
		global.GVA_LOG.Warn("试衣模型配置编码失败，跳过存量规范化", zap.Error(err))
		return
	}
	normalizedEncoded, err := json.Marshal(normalizedModels)
	if err != nil {
		global.GVA_LOG.Warn("试衣模型配置规范化编码失败", zap.Error(err))
		return
	}

	if string(currentEncoded) == string(normalizedEncoded) {
		return
	}

	changedModelKeys := collectChangedTryonModelKeys(models, normalizedModels)
	keysForLog := changedModelKeys
	if len(keysForLog) > 30 {
		keysForLog = append(append([]string{}, keysForLog[:30]...), fmt.Sprintf("...+%d", len(keysForLog)-30))
	}

	if err := global.GVA_DB.Model(&client.SysConfig{}).Where("id = ?", cfg.ID).Update("config_value", string(normalizedEncoded)).Error; err != nil {
		global.GVA_LOG.Warn("试衣模型配置规范化失败", zap.Error(err))
		return
	}

	nodeName := ""
	if host, err := os.Hostname(); err == nil {
		nodeName = strings.TrimSpace(host)
	}

	global.GVA_LOG.Info("试衣模型配置规范化完成",
		zap.Uint("configID", cfg.ID),
		zap.Int("modelCount", len(models)),
		zap.Int("changedModelCount", len(changedModelKeys)),
		zap.Strings("changedModelKeys", keysForLog),
		zap.String("node", nodeName),
		zap.Duration("duration", time.Since(startedAt)),
	)
}

func collectChangedTryonModelKeys(original []map[string]interface{}, normalized []map[string]interface{}) []string {
	maxLen := len(original)
	if len(normalized) > maxLen {
		maxLen = len(normalized)
	}

	changedKeys := make([]string, 0)
	seen := make(map[string]struct{})

	for i := 0; i < maxLen; i++ {
		originalJSON := []byte("null")
		normalizedJSON := []byte("null")

		if i < len(original) {
			if encoded, err := json.Marshal(original[i]); err == nil {
				originalJSON = encoded
			}
		}
		if i < len(normalized) {
			if encoded, err := json.Marshal(normalized[i]); err == nil {
				normalizedJSON = encoded
			}
		}

		if string(originalJSON) == string(normalizedJSON) {
			continue
		}

		modelKey := ""
		if i < len(normalized) {
			modelKey = toString(normalized[i]["key"])
		}
		if modelKey == "" && i < len(original) {
			modelKey = toString(original[i]["key"])
			if modelKey == "" {
				modelKey = toString(original[i]["modelKey"])
			}
		}
		if modelKey == "" {
			modelKey = fmt.Sprintf("index_%d", i)
		}

		if _, ok := seen[modelKey]; ok {
			continue
		}
		seen[modelKey] = struct{}{}
		changedKeys = append(changedKeys, modelKey)
	}

	return changedKeys
}

func normalizeTryonModelConfigItem(model map[string]interface{}, index int) map[string]interface{} {
	if model == nil {
		model = map[string]interface{}{}
	}

	modelUsage := normalizeTryonModelUsage(
		toString(model["modelUsage"]),
		toString(model["key"]),
		toString(model["model"]),
		toBoolLoose(model["supportsBeautify"], false),
	)
	modelKey := toString(model["key"])
	if modelKey == "" {
		modelKey = toString(model["modelKey"])
	}
	if modelKey == "" {
		if modelUsage == "beautify" {
			modelKey = fmt.Sprintf("beautify_%d", index+1)
		} else {
			modelKey = fmt.Sprintf("model_%d", index+1)
		}
	}

	nameValue := pickFirstValue(model, "name", "nameI18n", "title", "titleI18n")
	if nameValue == nil {
		nameValue = map[string]string{"zh": modelKey, "en": modelKey, "mn": modelKey}
	}
	descValue := pickFirstValue(model, "desc", "descI18n", "description", "descriptionI18n")
	if descValue == nil {
		descValue = map[string]string{}
	}

	scenes := toStringArrayLoose(model["scenes"])
	if len(scenes) == 0 {
		sceneType := toString(pickFirstValue(model, "sceneType", "scene", "roomType"))
		if sceneType != "" {
			scenes = []string{sceneType}
		} else if modelUsage == "beautify" {
			scenes = []string{"clothes", "shoes", "takeoff"}
		} else {
			scenes = []string{"clothes"}
		}
	}

	baseModel := toString(model["model"])
	if modelUsage == "beautify" && baseModel == "" {
		baseModel = toString(model["beautifyModel"])
	}
	if modelUsage == "beautify" && baseModel == "" {
		baseModel = "custom_beautify"
	}

	tokenBackups := toStringArrayLoose(pickFirstValue(model, "tokenBackups", "backupTokens"))

	basePayload := map[string]interface{}{
		"key":          modelKey,
		"modelUsage":   modelUsage,
		"enabled":      toBoolLoose(model["enabled"], true),
		"scenes":       scenes,
		"model":        baseModel,
		"name":         nameValue,
		"desc":         descValue,
		"cost":         clampNonNegativeInt(toIntLoose(model["cost"], 1)),
		"provider":     toString(model["provider"]),
		"mode":         firstNonEmptyStringValue(toString(model["mode"]), "prod"),
		"url":          firstNonEmptyStringValue(toString(model["url"]), toString(model["providerUrl"])),
		"taskQueryUrl": toString(model["taskQueryUrl"]),
		"token":        firstNonEmptyStringValue(toString(model["token"]), toString(model["providerToken"])),
		"tokenBackups": tokenBackups,
	}

	providerLower := strings.ToLower(strings.TrimSpace(toString(model["provider"])))
	autoParsingDefault := isAliyunAitryonSeriesLoose(model)

	switch modelUsage {
	case "tryon":
		basePayload["cost"] = clampNonNegativeInt(toIntLoose(model["cost"], 1))
		basePayload["resolution"] = toIntLoose(model["resolution"], -1)
		basePayload["restoreFace"] = toBoolLoose(model["restoreFace"], true)
		basePayload["freeQuotaTotal"] = clampNonNegativeInt(toIntLoose(model["freeQuotaTotal"], 0))

		parsingModelKey := toString(model["parsingModelKey"])
		if parsingModelKey == "" && autoParsingDefault {
			parsingModelKey = "aliyun_aitryon_parsing"
		}
		basePayload["parsingModelKey"] = parsingModelKey

		refinerModelKey := toString(model["refinerModelKey"])
		if refinerModelKey == "" && toBoolLoose(model["supportsRefiner"], autoParsingDefault) {
			refinerModelKey = "aliyun_aitryon_refiner"
		}
		basePayload["refinerModelKey"] = refinerModelKey

		if strings.Contains(providerLower, "gradio") || strings.Contains(providerLower, "huggingface") || strings.Contains(providerLower, "hf") {
			basePayload["apiName"] = firstNonEmptyStringValue(toString(model["apiName"]), "/tryon")
			basePayload["garmentDes"] = firstNonEmptyStringValue(toString(model["garmentDes"]), "clothing item")
			basePayload["isChecked"] = toBoolLoose(model["isChecked"], true)
			basePayload["isCheckedCrop"] = toBoolLoose(model["isCheckedCrop"], false)
			basePayload["denoiseSteps"] = clampPositiveInt(toIntLoose(model["denoiseSteps"], 30), 30)
			basePayload["seed"] = toIntLoose(model["seed"], 42)
		}
		return basePayload

	case "refiner":
		refinerExtraCost := toIntLooseWithFallback(model, []string{"refinerExtraCost", "refinerExtraPoints", "cost"}, 1)
		basePayload["cost"] = clampNonNegativeInt(toIntLoose(model["cost"], refinerExtraCost))
		basePayload["refinerExtraCost"] = clampNonNegativeInt(refinerExtraCost)
		basePayload["freeQuotaTotal"] = clampNonNegativeInt(toIntLooseWithFallback(model, []string{"freeQuotaTotal", "refinerFreeQuotaTotal"}, 0))
		basePayload["refinerGender"] = normalizeRefinerGender(toString(model["refinerGender"]))

		if strings.TrimSpace(toString(basePayload["url"])) == "" {
			basePayload["url"] = toString(model["refinerUrl"])
		}
		if strings.TrimSpace(toString(basePayload["taskQueryUrl"])) == "" {
			basePayload["taskQueryUrl"] = toString(model["refinerTaskQueryUrl"])
		}
		if strings.TrimSpace(toString(basePayload["token"])) == "" {
			basePayload["token"] = toString(model["refinerToken"])
		}
		return basePayload

	case "parsing":
		parsingExtraCost := clampNonNegativeInt(toIntLooseWithFallback(model, []string{"parsingExtraCost", "parsingExtraPoints", "cost"}, 0))
		basePayload["cost"] = clampNonNegativeInt(toIntLoose(model["cost"], parsingExtraCost))
		basePayload["parsingExtraCost"] = parsingExtraCost
		basePayload["freeQuotaTotal"] = clampNonNegativeInt(toIntLoose(model["freeQuotaTotal"], 0))
		basePayload["clothesType"] = toStringArrayLooseWithDefault(model["clothesType"], []string{"upper"})
		return basePayload

	case "beautify":
		beautifyExtraCost := toIntLooseWithFallback(model, []string{"beautifyExtraCost", "beautifyExtraPoints", "cost"}, 0)
		basePayload["cost"] = clampNonNegativeInt(toIntLoose(model["cost"], 0))
		basePayload["beautifyExtraCost"] = clampNonNegativeInt(beautifyExtraCost)
		basePayload["beautifyModel"] = firstNonEmptyStringValue(toString(model["beautifyModel"]), baseModel)
		if degree, ok := parseOptionalBeautifyDegree(model["beautifyRetouchDegree"]); ok {
			basePayload["beautifyRetouchDegree"] = degree
		}
		if degree, ok := parseOptionalBeautifyDegree(model["beautifyWhiteningDegree"]); ok {
			basePayload["beautifyWhiteningDegree"] = degree
		}
		basePayload["beautifyUrl"] = toString(model["beautifyUrl"])
		basePayload["beautifyAccessKeyId"] = toString(model["beautifyAccessKeyId"])
		basePayload["beautifyAccessKeySecret"] = toString(model["beautifyAccessKeySecret"])
		basePayload["beautifySecurityToken"] = toString(model["beautifySecurityToken"])
		basePayload["beautifyToken"] = toString(model["beautifyToken"])
		basePayload["beautifyDesc"] = firstNonNilValue(model["beautifyDesc"], map[string]string{})
		return basePayload
	}

	return basePayload
}

func normalizeTryonModelUsage(value string, key string, model string, supportsBeautify bool) string {
	usage := strings.ToLower(strings.TrimSpace(value))
	switch usage {
	case "tryon", "refiner", "parsing", "beautify":
		return usage
	}

	if supportsBeautify {
		return "beautify"
	}

	hint := strings.ToLower(strings.TrimSpace(key)) + " " + strings.ToLower(strings.TrimSpace(model))
	if strings.Contains(hint, "refiner") {
		return "refiner"
	}
	if strings.Contains(hint, "parsing") || strings.Contains(hint, "takeoff") {
		return "parsing"
	}
	return "tryon"
}

func toIntLoose(value interface{}, fallback int) int {
	switch v := value.(type) {
	case float64:
		return int(v)
	case float32:
		return int(v)
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	case string:
		text := strings.TrimSpace(v)
		if text == "" {
			return fallback
		}
		var parsed int
		if _, err := fmt.Sscanf(text, "%d", &parsed); err == nil {
			return parsed
		}
	}
	return fallback
}

func parseOptionalBeautifyDegree(value interface{}) (float64, bool) {
	text := strings.TrimSpace(toString(value))
	if text == "" {
		return 0, false
	}

	var degree float64
	if _, err := fmt.Sscanf(text, "%f", &degree); err != nil {
		return 0, false
	}
	if degree <= 0 {
		return 0, false
	}
	if degree > 1.5 {
		return 0, false
	}
	return degree, true
}

func toIntLooseWithFallback(model map[string]interface{}, keys []string, fallback int) int {
	for _, key := range keys {
		if model == nil {
			break
		}
		value, ok := model[key]
		if !ok || value == nil {
			continue
		}
		if text, isString := value.(string); isString && strings.TrimSpace(text) == "" {
			continue
		}
		return toIntLoose(value, fallback)
	}
	return fallback
}

func toStringArrayLoose(value interface{}) []string {
	result := make([]string, 0)
	appendValue := func(raw string) {
		item := strings.TrimSpace(raw)
		if item == "" {
			return
		}
		for _, existing := range result {
			if strings.EqualFold(existing, item) {
				return
			}
		}
		result = append(result, item)
	}

	switch v := value.(type) {
	case []interface{}:
		for _, item := range v {
			appendValue(toString(item))
		}
	case []string:
		for _, item := range v {
			appendValue(item)
		}
	case string:
		text := strings.TrimSpace(v)
		if text == "" {
			return result
		}
		if strings.HasPrefix(text, "[") {
			parsed := make([]string, 0)
			if err := json.Unmarshal([]byte(text), &parsed); err == nil {
				for _, item := range parsed {
					appendValue(item)
				}
				return result
			}
		}
		for _, item := range strings.FieldsFunc(text, func(r rune) bool {
			return r == ',' || r == ';' || r == '\n' || r == '\r' || r == '\t'
		}) {
			appendValue(item)
		}
	}

	return result
}

func toStringArrayLooseWithDefault(value interface{}, fallback []string) []string {
	items := toStringArrayLoose(value)
	if len(items) > 0 {
		return items
	}
	if len(fallback) == 0 {
		return []string{}
	}
	cloned := make([]string, 0, len(fallback))
	for _, item := range fallback {
		item = strings.TrimSpace(item)
		if item != "" {
			cloned = append(cloned, item)
		}
	}
	return cloned
}

func pickFirstValue(model map[string]interface{}, keys ...string) interface{} {
	if model == nil {
		return nil
	}
	for _, key := range keys {
		value, ok := model[key]
		if !ok || value == nil {
			continue
		}
		if text, isString := value.(string); isString && strings.TrimSpace(text) == "" {
			continue
		}
		return value
	}
	return nil
}

func firstNonNilValue(value interface{}, fallback interface{}) interface{} {
	if value == nil {
		return fallback
	}
	return value
}

func firstNonEmptyStringValue(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func clampNonNegativeInt(value int) int {
	if value < 0 {
		return 0
	}
	return value
}

func clampPositiveInt(value int, fallback int) int {
	if value <= 0 {
		return fallback
	}
	return value
}

func clampIntRange(value int, minValue int, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func normalizeRefinerGender(value string) string {
	gender := strings.ToLower(strings.TrimSpace(value))
	if gender == "man" || gender == "woman" {
		return gender
	}
	return "woman"
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
