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
		// security 分组
		{ConfigKey: "captcha_expiry_seconds", ConfigValue: "300", ConfigName: "验证码有效期", ConfigGroup: "security", Remark: "图形验证码有效期(秒)，默认300秒"},
		{ConfigKey: "captcha_rate_limit", ConfigValue: "10", ConfigName: "验证码请求频率限制", ConfigGroup: "security", Remark: "每分钟每IP最多请求验证码次数"},
		{ConfigKey: "captcha_rate_limit_window_seconds", ConfigValue: "60", ConfigName: "验证码频率窗口(秒)", ConfigGroup: "security", Remark: "验证码频率限制统计窗口时长，单位秒"},
		{ConfigKey: "register_ip_limit", ConfigValue: "3", ConfigName: "注册IP限制", ConfigGroup: "security", Remark: "同一IP最多可注册账号数量"},
		{ConfigKey: "security_login_ip_rate_limit_per_minute", ConfigValue: "30", ConfigName: "登录IP每分钟上限", ConfigGroup: "security", Remark: "同一IP每分钟允许的登录请求次数（用户名登录+手机号登录）"},
		{ConfigKey: "security_login_ip_rate_limit_window_seconds", ConfigValue: "60", ConfigName: "登录IP限流窗口(秒)", ConfigGroup: "security", Remark: "登录IP频率限制统计窗口时长，单位秒"},
		{ConfigKey: "login_fail_max", ConfigValue: "5", ConfigName: "登录失败上限", ConfigGroup: "security", Remark: "连续登录失败N次后锁定"},
		{ConfigKey: "login_fail_wait_seconds", ConfigValue: "900", ConfigName: "登录锁定等待时间", ConfigGroup: "security", Remark: "登录锁定后等待时间(秒)，默认900秒(15分钟)"},
		{ConfigKey: "security_init_api_enabled", ConfigValue: "false", ConfigName: "初始化接口开关", ConfigGroup: "security", Remark: "是否启用 /init/checkdb 与 /init/initdb 接口(true/false)"},
		{ConfigKey: "security_init_api_private_network_only", ConfigValue: "true", ConfigName: "初始化接口内网限制", ConfigGroup: "security", Remark: "初始化接口是否仅允许内网/回环地址访问(true/false)"},
		{ConfigKey: "security_trusted_proxies", ConfigValue: "127.0.0.1,::1,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,fc00::/7", ConfigName: "可信代理白名单", ConfigGroup: "security", Remark: "用于解析真实客户端IP的可信代理列表，逗号分隔，支持IP/CIDR"},
		{ConfigKey: "security_upload_max_size_mb", ConfigValue: "20", ConfigName: "上传大小上限(MB)", ConfigGroup: "security", Remark: "文件上传大小上限，单位MB，默认20"},
		{ConfigKey: "security_upload_strict_validation_enabled", ConfigValue: "true", ConfigName: "上传严格校验开关", ConfigGroup: "security", Remark: "是否启用上传危险扩展名与MIME拦截(true/false)"},
		{ConfigKey: "security_export_allow_query_token", ConfigValue: "false", ConfigName: "导出URL Token兼容开关", ConfigGroup: "security", Remark: "是否允许导出接口通过URL query传token(true/false)"},
		{ConfigKey: "security_ws_allow_query_token", ConfigValue: "false", ConfigName: "WS URL Token兼容开关", ConfigGroup: "security", Remark: "是否允许WebSocket通过URL query传token(true/false)"},
		{ConfigKey: "security_ws_max_conns_per_ip", ConfigValue: "8", ConfigName: "WS单IP连接上限", ConfigGroup: "security", Remark: "单个IP允许的WebSocket最大连接数，最小值1"},
		{ConfigKey: "security_public_rate_limit_enabled", ConfigValue: "true", ConfigName: "公开接口限流开关", ConfigGroup: "security", Remark: "是否启用公开接口统一限流(true/false)"},
		{ConfigKey: "security_public_rate_limit_window_seconds", ConfigValue: "60", ConfigName: "公开接口限流窗口(秒)", ConfigGroup: "security", Remark: "公开接口统一限流统计窗口时长，单位秒"},
		{ConfigKey: "security_public_rate_limit_max_requests", ConfigValue: "300", ConfigName: "公开接口限流次数", ConfigGroup: "security", Remark: "同一IP在限流窗口内允许的公开接口请求次数"},
		{ConfigKey: "security_attack_auto_ban_enabled", ConfigValue: "true", ConfigName: "攻击自动封禁开关", ConfigGroup: "security", Remark: "是否启用攻击行为自动封禁(true/false)"},
		{ConfigKey: "security_attack_auto_ban_window_seconds", ConfigValue: "3600", ConfigName: "攻击自动封禁窗口(秒)", ConfigGroup: "security", Remark: "攻击累计计数窗口时长，单位秒"},
		{ConfigKey: "security_attack_auto_ban_duration_minutes", ConfigValue: "60", ConfigName: "攻击自动封禁时长(分钟)", ConfigGroup: "security", Remark: "达到阈值后自动封禁时长，单位分钟"},
		{ConfigKey: "security_attack_auto_ban_threshold_default", ConfigValue: "10", ConfigName: "攻击自动封禁默认阈值", ConfigGroup: "security", Remark: "默认攻击类型自动封禁阈值，设为0表示关闭"},
		{ConfigKey: "security_attack_auto_ban_threshold_sys_error_rate_limit", ConfigValue: "30", ConfigName: "错误上报自动封禁阈值", ConfigGroup: "security", Remark: "错误上报限流命中自动封禁阈值，设为0表示关闭"},
		{ConfigKey: "security_sys_error_create_rate_limit_per_minute", ConfigValue: "30", ConfigName: "错误上报每分钟上限", ConfigGroup: "security", Remark: "同一IP每分钟允许创建错误日志次数，设为0表示关闭"},
		{ConfigKey: "security_sys_error_create_rate_limit_window_seconds", ConfigValue: "60", ConfigName: "错误上报限流窗口(秒)", ConfigGroup: "security", Remark: "错误上报限流统计窗口时长，单位秒"},
		{ConfigKey: "security_visitor_heartbeat_rate_limit_per_minute", ConfigValue: "120", ConfigName: "访客心跳每分钟上限", ConfigGroup: "security", Remark: "同一IP每分钟允许的访客心跳请求次数"},
		{ConfigKey: "security_visitor_heartbeat_dedupe_seconds", ConfigValue: "3", ConfigName: "访客心跳去重窗口(秒)", ConfigGroup: "security", Remark: "同一visitor在窗口内重复心跳仅保留一次入库"},
		{ConfigKey: "security_tryon_create_rate_limit_per_minute", ConfigValue: "20", ConfigName: "试衣创建每分钟上限", ConfigGroup: "security", Remark: "单用户每分钟允许创建试衣任务次数"},
		{ConfigKey: "security_tryon_create_concurrency_limit", ConfigValue: "2", ConfigName: "试衣创建并发上限", ConfigGroup: "security", Remark: "单用户同时处于processing状态的试衣任务上限"},
		{ConfigKey: "security_order_create_rate_limit_per_minute", ConfigValue: "30", ConfigName: "商城下单每分钟上限", ConfigGroup: "security", Remark: "单用户每分钟允许创建商品订单次数，设为0表示关闭"},
		{ConfigKey: "security_order_pending_limit_per_user", ConfigValue: "10", ConfigName: "商城待支付订单上限", ConfigGroup: "security", Remark: "单用户处于待支付/待确认状态的商品订单最大数量，设为0表示关闭"},
		{ConfigKey: "security_tryon_recharge_create_rate_limit_per_minute", ConfigValue: "20", ConfigName: "试衣币充值每分钟建单上限", ConfigGroup: "security", Remark: "单用户每分钟允许创建试衣币充值订单次数，设为0表示关闭"},
		{ConfigKey: "security_tryon_recharge_pending_limit_per_user", ConfigValue: "8", ConfigName: "试衣币充值待支付订单上限", ConfigGroup: "security", Remark: "单用户处于待支付/待确认状态的试衣币充值订单最大数量，设为0表示关闭"},
		{ConfigKey: "security_cleanup_legacy_casbin_overgrant", ConfigValue: "false", ConfigName: "历史权限清理开关", ConfigGroup: "security", Remark: "是否启用历史Casbin过授权规则清理(true/false)，建议仅在维护窗口临时开启"},
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
		{ConfigKey: "tryon_append_parsing_failed_tip", ConfigValue: "false", ConfigName: "分割失败提示拼接开关", ConfigGroup: "tryon", Remark: "试衣成功但分割增强失败时，是否在客户端拼接分割失败提示(true/false)"},
		{ConfigKey: "tryon_parsing_failed_tip_text", ConfigValue: "{\"zh\":\"试衣成功，但分割增强失败，相关金币已退回\",\"en\":\"Try-on succeeded, but parsing enhancement failed. Related coins have been refunded.\",\"mn\":\"Туршилт амжилттай боловч segmentation enhancement амжилтгүй боллоо. Холбогдох зоос буцаан олгогдлоо.\"}", ConfigName: "分割失败提示文本", ConfigGroup: "tryon", Remark: "试衣成功但分割增强失败时拼接的提示(JSON多语言)"},
		{ConfigKey: "tryon_append_refiner_failed_tip", ConfigValue: "true", ConfigName: "精修失败提示拼接开关", ConfigGroup: "tryon", Remark: "试衣成功且精修失败时，是否在客户端拼接“金币已退回”提示(true/false)"},
		{ConfigKey: "tryon_refiner_failed_tip_text", ConfigValue: "{\"zh\":\"试衣成功，但精修失败，金币已退回\",\"en\":\"Try-on succeeded, but refiner failed. Coins have been refunded.\",\"mn\":\"Туршилт амжилттай боловч нарийвчлал амжилтгүй боллоо. Зоос буцаан олгогдсон.\"}", ConfigName: "精修失败提示文本", ConfigGroup: "tryon", Remark: "试衣成功但精修失败时拼接的提示(JSON多语言)"},
		{ConfigKey: "tryon_tutorial_title", ConfigValue: defaultTryonTutorialTitleSysConfigValue(), ConfigName: "试衣教程标题", ConfigGroup: "tryon", Remark: "试衣页面教程弹窗标题(JSON多语言)"},
		{ConfigKey: "tryon_tutorial_content", ConfigValue: defaultTryonTutorialContentSysConfigValue(), ConfigName: "试衣教程内容", ConfigGroup: "tryon", Remark: "试衣页面教程弹窗内容(JSON多语言，换行分段)"},
		{ConfigKey: "tryon_models", ConfigValue: defaultTryonModelsSysConfigValue(), ConfigName: "试衣模型列表", ConfigGroup: "tryon", Remark: "JSON数组：统一承载试衣/试鞋/取衣/美肤模型配置，按模型独立设置参数"},
		{ConfigKey: "tryon_recharge_plans", ConfigValue: "[{\"points\":{\"zh\":\"50\",\"en\":\"50\",\"mn\":\"50\"},\"coinLabel\":{\"zh\":\"试衣币\",\"en\":\"Try-on Coins\",\"mn\":\"Туршилтын зоос\"},\"price\":990,\"priceI18n\":{\"zh\":990,\"en\":990,\"mn\":990}},{\"points\":{\"zh\":\"180\",\"en\":\"180\",\"mn\":\"180\"},\"coinLabel\":{\"zh\":\"试衣币\",\"en\":\"Try-on Coins\",\"mn\":\"Туршилтын зоос\"},\"price\":2990,\"priceI18n\":{\"zh\":2990,\"en\":2990,\"mn\":2990}},{\"points\":{\"zh\":\"680\",\"en\":\"680\",\"mn\":\"680\"},\"coinLabel\":{\"zh\":\"试衣币\",\"en\":\"Try-on Coins\",\"mn\":\"Туршилтын зоос\"},\"price\":9990,\"priceI18n\":{\"zh\":9990,\"en\":9990,\"mn\":9990}}]", ConfigName: "试衣币充值套餐", ConfigGroup: "tryon", Remark: "后台可视化维护：点数、币名与多语言价格（基础分价 + 汇率换算）"},
		{ConfigKey: "tryon_recharge_pending_confirm_close_minutes", ConfigValue: "180", ConfigName: "试衣币待确认超时(分钟)", ConfigGroup: "tryon", Remark: "试衣币充值订单进入待确认后自动关闭的分钟数"},
		{ConfigKey: "tryon_media_public_base_url", ConfigValue: "", ConfigName: "试衣图片公网前缀", ConfigGroup: "tryon", Remark: "当上传返回相对路径或localhost/内网URL时，自动拼接为公网地址，如https://back.example.com"},
		// order 分组
		{ConfigKey: "order_close_minutes", ConfigValue: "20", ConfigName: "订单自动关闭时间", ConfigGroup: "order", Remark: "未支付订单自动关闭的分钟数"},
		{ConfigKey: "order_pending_confirm_close_minutes", ConfigValue: "180", ConfigName: "订单待确认超时(分钟)", ConfigGroup: "order", Remark: "商品订单进入待后台确认后自动关闭的分钟数"},
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
	normalizeConfigValueToMultilingualObject("tryon_tutorial_title")
	normalizeConfigValueToMultilingualObject("tryon_tutorial_content")
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

func defaultTryonTutorialTitleSysConfigValue() string {
	return "{\"zh\":\"\\u8bd5\\u8863\\u6559\\u7a0b\",\"en\":\"Try-On Guide\",\"mn\":\"\\u0422\\u0443\\u0440\\u0448\\u0438\\u043b\\u0442\\u044b\\u043d \\u0437\\u0430\\u0430\\u0432\\u0430\\u0440\",\"zh-TW\":\"\\u8a66\\u8863\\u6559\\u7a0b\",\"th\":\"\\u0e04\\u0e39\\u0e48\\u0e21\\u0e37\\u0e2d\\u0e17\\u0e14\\u0e25\\u0e2d\\u0e07\\u0e2a\\u0e27\\u0e21\",\"hi\":\"\\u091f\\u094d\\u0930\\u093e\\u092f-\\u0911\\u0928 \\u0917\\u093e\\u0907\\u0921\",\"id\":\"Panduan Try-On\",\"vi\":\"H\\u01b0\\u1edbng d\\u1eabn d\\u00f9ng th\\u1eed\",\"ar\":\"\\u062f\\u0644\\u064a\\u0644 \\u0627\\u0644\\u062a\\u062c\\u0631\\u0628\\u0629\",\"ja\":\"\\u8a66\\u7740\\u30ac\\u30a4\\u30c9\",\"ko\":\"\\uc2dc\\ucc29 \\uc548\\ub0b4\",\"ms\":\"Panduan Mencuba\"}"
}

func defaultTryonTutorialContentSysConfigValue() string {
	return "{\"zh\":\"\\u7f51\\u9875\\u7248\\u63a8\\u8350\\u4f7f\\u7528 Chrome \\u548c Safari \\u6d4f\\u89c8\\u5668\\u3002\\n\\u6a21\\u7279\\u56fe\\u53ef\\u4ee5\\u4e0a\\u4f20\\u4fa7\\u8eab\\u59ff\\u52bf\\u3001\\u5750\\u7740\\u59ff\\u52bf\\u7b49\\u56fe\\u7247\\uff0c\\u4f46\\u662f\\u6b63\\u9762\\u7ad9\\u7acb\\u59ff\\u52bf\\u751f\\u6210\\u6548\\u679c\\u66f4\\u4f73\\u3002\\n\\u5982\\u679c\\u6ca1\\u6709\\u5355\\u72ec\\u7684\\u670d\\u88c5\\u56fe\\u7247\\uff0c\\u53ef\\u76f4\\u63a5\\u5728\\u201c\\u4e0a\\u88c5\\u56fe\\u201d\\u4f4d\\u7f6e\\u4e0a\\u4f20\\u4eba\\u7269\\u7a7f\\u7740\\u7684\\u56fe\\u7247\\u3002\\n\\u5982\\u679c\\u670d\\u88c5\\u56fe\\u7247\\u5f88\\u96be\\u83b7\\u53d6\\uff0c\\u53ef\\u4ee5\\u4f7f\\u7528\\u201c\\u88c1\\u526a\\u4e0a\\u4f20\\u201d\\u88c1\\u526a\\uff0c\\u540e\\u7eed\\u4f1a\\u6dfb\\u52a0\\u201c\\u53d6\\u8863\\u201d\\u529f\\u80fd\\uff0c\\u656c\\u8bf7\\u671f\\u5f85\\u3002\\n\\u751f\\u6210\\u8d28\\u91cf\\u548c\\u6a21\\u578b\\u3001\\u56fe\\u7247\\u8d28\\u91cf\\u3001\\u80cc\\u666f\\u6709\\u5173\\uff0c\\u80cc\\u666f\\u8d8a\\u7b80\\u5355\\uff0c\\u4eba\\u7269\\u5728\\u56fe\\u7247\\u5360\\u6bd4\\u8d8a\\u5927\\uff0c\\u751f\\u6210\\u8d28\\u91cf\\u66f4\\u597d\\u3002\\n\\u56fe\\u7247\\u7cbe\\u4fee\\u6a21\\u578b\\u4f1a\\u63d0\\u9ad8\\u8863\\u7269\\u3001\\u76ae\\u80a4\\u7684\\u7ec6\\u8282\\u8d28\\u611f\\uff0c\\u4f46\\u4e0d\\u9002\\u7528\\u6240\\u6709\\uff0c\\u6709\\u65f6\\u5019\\u53ef\\u80fd\\u9002\\u5f97\\u5176\\u53cd\\uff0c\\u8bf7\\u6ce8\\u610f\\u4f7f\\u7528\\u3002\\n\\u53ef\\u4ee5\\u5728\\u201c\\u6211\\u7684\\u6a21\\u7279\\u201d\\u548c\\u201c\\u6211\\u7684\\u8863\\u6a71\\u201d\\u9875\\u9762\\u4e0a\\u4f20\\u4fdd\\u5b58\\u81ea\\u5df1\\u5e38\\u7528\\u7684\\u6a21\\u7279\\u56fe\\u548c\\u670d\\u88c5\\uff0c\\u8fd9\\u6837\\u5c31\\u4e0d\\u7528\\u6bcf\\u6b21\\u90fd\\u4e0a\\u4f20\\u540c\\u6837\\u7684\\u56fe\\u7247\\u3002\\n\\u5982\\u679c\\u8bd5\\u7a7f\\u7684\\u670d\\u88c5\\u592a\\u5c0f\\u5bfc\\u81f4\\u6ca1\\u6709\\u8986\\u76d6\\u6a21\\u7279\\u56fe\\u539f\\u6765\\u6ca1\\u9732\\u51fa\\u76ae\\u80a4\\u7684\\u4f4d\\u7f6e\\uff0c\\u5bf9\\u5e94\\u4f4d\\u7f6e\\u7684\\u76ae\\u80a4\\u53ef\\u80fd\\u751f\\u6210\\u4e0d\\u597d\\uff0c\\u53ef\\u4ee5\\u5207\\u6362\\u66f4\\u597d\\u7684\\u6a21\\u578b\\u6216\\u8005\\u66f4\\u6362\\u9732\\u51fa\\u5bf9\\u5e94\\u4f4d\\u7f6e\\u76ae\\u80a4\\u7684\\u6a21\\u7279\\u56fe\\u3002\\u8fd9\\u4e2a\\u573a\\u666f\\u4e00\\u822c\\u662f\\u5728\\u8bd5\\u7a7f\\u6bd4\\u57fa\\u5c3c\\u6216\\u8005\\u5185\\u8863\\u3002\\n\\u6a21\\u7279\\u56fe\\u4e0d\\u80fd\\u9732\\u51fa\\u654f\\u611f\\u4f4d\\u7f6e\\uff0c\\u592a\\u88f8\\u9732\\u4f1a\\u5ba1\\u6838\\u751f\\u6210\\u5931\\u8d25\\u3002\\n\\u4ee5\\u540e\\u4f1a\\u66f4\\u65b0\\u5bf9\\u4e0d\\u540c\\u59ff\\u52bf\\u6a21\\u7279\\u56fe\\u3001\\u5185\\u8863\\u3001\\u6bd4\\u57fa\\u5c3c\\u7b49\\u751f\\u6210\\u66f4\\u597d\\u7684\\u6a21\\u578b\\uff0c\\u8bf7\\u968f\\u65f6\\u5173\\u6ce8\\u5e73\\u53f0\\u516c\\u544a\\u3002\\n\\u6a21\\u7279\\u56fe\\u53ef\\u4ee5\\u4e0a\\u4f20\\u4fa7\\u8eab\\u56fe\\uff08\\u6b63\\u9762\\u7167\\u4e00\\u822c\\u6548\\u679c\\u66f4\\u597d\\uff09\\u3002\\u5982\\u679c\\u6ca1\\u6709\\u5b8c\\u6574\\u7684\\u4e0a\\u88c5\\u548c\\u4e0b\\u88c5\\u56fe\\uff0c\\u53ef\\u4ee5\\u76f4\\u63a5\\u4e0a\\u4f20\\u4eba\\u7269\\u7a7f\\u7740\\u56fe\\uff08\\u4e0a\\u4f20\\u5230\\u4e0a\\u88c5\\u56fe\\u5373\\u53ef\\uff0c\\u4e0d\\u7528\\u4f20\\u4e0b\\u88c5\\u56fe\\uff09\\uff0c\\u6a21\\u578b\\u4f1a\\u81ea\\u52a8\\u628a\\u8be5\\u4eba\\u7269\\u7684\\u4e0a\\u4e0b\\u8eab\\u670d\\u88c5\\u8bd5\\u7a7f\\u5230\\u6a21\\u7279\\u56fe\\u4e0a\\u3002\",\"en\":\"For the web version, Chrome and Safari are recommended.\\nYou can upload side-view or sitting model photos, but front-facing standing poses usually generate better results.\\nIf you do not have separate garment images, you can directly upload a dressed-person image in the \\\"Upper Image\\\" slot.\\nIf garment images are hard to obtain, use \\\"Crop Upload\\\" first. The \\\"Garment Extraction\\\" feature will be added later.\\nGeneration quality is related to the model, image quality, and background. Simpler backgrounds and a larger person ratio usually produce better quality.\\nThe Image Refiner model can improve garment and skin details, but it is not suitable for every case and may sometimes make results worse. Use it with care.\\nYou can upload and save frequently used model and garment images in \\\"My Models\\\" and \\\"My Closet\\\" so you do not need to upload the same images every time.\\nIf a try-on garment is too small and does not cover areas that were not exposed on the original model photo, skin in those areas may render poorly. In this case (common with bikini/underwear try-ons), switch to a better model or use a model photo that exposes the corresponding skin area.\\nModel photos must not expose sensitive body parts. Overly revealing images may fail moderation and generation.\\nWe will continue updating better models for different poses and for underwear/bikini scenarios. Please follow platform announcements.\\nYou can upload side-view model photos (front-facing photos usually work better). If you do not have complete upper and lower garment images, you can directly upload a dressed-person image (upload it to \\\"Upper Image\\\" only, no need to upload a lower image). The model will automatically try on that person's upper and lower garments onto the model photo.\",\"mn\":\"\\u0412\\u044d\\u0431 \\u0445\\u0443\\u0432\\u0438\\u043b\\u0431\\u0430\\u0440\\u0442 Chrome \\u0431\\u043e\\u043b\\u043e\\u043d Safari \\u0445\\u04e9\\u0442\\u0447\\u0438\\u0439\\u0433 \\u0430\\u0448\\u0438\\u0433\\u043b\\u0430\\u0445\\u044b\\u0433 \\u0437\\u04e9\\u0432\\u043b\\u04e9\\u0436 \\u0431\\u0430\\u0439\\u043d\\u0430.\\n\\u041c\\u043e\\u0434\\u0435\\u043b\\u0438\\u0439\\u043d \\u0437\\u0443\\u0440\\u0430\\u0433\\u0442 \\u0445\\u0430\\u0436\\u0443\\u0443, \\u0441\\u0443\\u0443\\u0436 \\u0431\\u0443\\u0439 \\u0437\\u044d\\u0440\\u044d\\u0433 \\u043f\\u043e\\u0437\\u0442\\u043e\\u0439 \\u0437\\u0443\\u0440\\u0430\\u0433 \\u043e\\u0440\\u0443\\u0443\\u043b\\u0436 \\u0431\\u043e\\u043b\\u043d\\u043e. \\u0413\\u044d\\u0445\\u0434\\u044d\\u044d \\u0443\\u0440\\u0434\\u0430\\u0430\\u0441 \\u044d\\u0433\\u0446 \\u0437\\u043e\\u0433\\u0441\\u0441\\u043e\\u043d \\u043f\\u043e\\u0437 \\u0438\\u0445\\u044d\\u0432\\u0447\\u043b\\u044d\\u043d \\u0438\\u043b\\u04af\\u04af \\u0441\\u0430\\u0439\\u043d \\u04af\\u0440 \\u0434\\u04af\\u043d\\u0442\\u044d\\u0439.\\n\\u0422\\u0443\\u0441\\u0434\\u0430\\u0430 \\u0445\\u0443\\u0432\\u0446\\u0430\\u0441\\u043d\\u044b \\u0437\\u0443\\u0440\\u0430\\u0433 \\u0431\\u0430\\u0439\\u0445\\u0433\\u04af\\u0439 \\u0431\\u043e\\u043b \\u201c\\u0414\\u044d\\u044d\\u0434 \\u0437\\u0443\\u0440\\u0430\\u0433\\u201d \\u0445\\u044d\\u0441\\u044d\\u0433\\u0442 \\u0445\\u0443\\u0432\\u0446\\u0430\\u0441\\u0442\\u0430\\u0439 \\u0445\\u04af\\u043d\\u0438\\u0439 \\u0437\\u0443\\u0440\\u0433\\u0438\\u0439\\u0433 \\u0448\\u0443\\u0443\\u0434 \\u043e\\u0440\\u0443\\u0443\\u043b\\u0436 \\u0431\\u043e\\u043b\\u043d\\u043e.\\n\\u0425\\u0443\\u0432\\u0446\\u0430\\u0441\\u043d\\u044b \\u0437\\u0443\\u0440\\u0430\\u0433 \\u043e\\u043b\\u043e\\u0445\\u043e\\u0434 \\u0445\\u044d\\u0446\\u04af\\u04af \\u0431\\u043e\\u043b \\u201c\\u0422\\u0430\\u0439\\u0440\\u0447 \\u043e\\u0440\\u0443\\u0443\\u043b\\u0430\\u0445\\u201d-\\u044b\\u0433 \\u0430\\u0448\\u0438\\u0433\\u043b\\u0430\\u0436 \\u0431\\u043e\\u043b\\u043d\\u043e. \\u201c\\u0425\\u0443\\u0432\\u0446\\u0430\\u0441 \\u044f\\u043b\\u0433\\u0430\\u0445\\u201d \\u0444\\u0443\\u043d\\u043a\\u0446\\u0438\\u0439\\u0433 \\u0434\\u0430\\u0440\\u0430\\u0430 \\u043d\\u044c \\u043d\\u044d\\u043c\\u044d\\u0445 \\u0442\\u0443\\u043b \\u0445\\u04af\\u043b\\u044d\\u044d\\u043d\\u044d \\u04af\\u04af.\\n\\u04ae\\u0440 \\u0434\\u04af\\u043d\\u0433\\u0438\\u0439\\u043d \\u0447\\u0430\\u043d\\u0430\\u0440 \\u043d\\u044c \\u043c\\u043e\\u0434\\u0435\\u043b\\u044c, \\u0437\\u0443\\u0440\\u0433\\u0438\\u0439\\u043d \\u0447\\u0430\\u043d\\u0430\\u0440, \\u0434\\u044d\\u0432\\u0441\\u0433\\u044d\\u0440\\u044d\\u044d\\u0441 \\u0448\\u0430\\u043b\\u0442\\u0433\\u0430\\u0430\\u043b\\u043d\\u0430. \\u0414\\u044d\\u0432\\u0441\\u0433\\u044d\\u0440 \\u044d\\u043d\\u0433\\u0438\\u0439\\u043d, \\u0445\\u04af\\u043d\\u0438\\u0439 \\u044d\\u0437\\u043b\\u044d\\u0445 \\u0445\\u0443\\u0432\\u044c \\u0438\\u0445 \\u0431\\u0430\\u0439\\u0445 \\u0442\\u0443\\u0441\\u0430\\u043c \\u0447\\u0430\\u043d\\u0430\\u0440 \\u0438\\u043b\\u04af\\u04af \\u0441\\u0430\\u0439\\u043d.\\n\\u0417\\u0443\\u0440\\u0433\\u0438\\u0439\\u043d \\u043d\\u0430\\u0440\\u0438\\u0439\\u0432\\u0447\\u043b\\u0430\\u043b\\u044b\\u043d \\u043c\\u043e\\u0434\\u0435\\u043b\\u044c \\u043d\\u044c \\u0445\\u0443\\u0432\\u0446\\u0430\\u0441, \\u0430\\u0440\\u044c\\u0441\\u043d\\u044b \\u0434\\u0435\\u0442\\u0430\\u043b\\u0438\\u0439\\u0433 \\u0441\\u0430\\u0439\\u0436\\u0440\\u0443\\u0443\\u043b\\u0436 \\u0447\\u0430\\u0434\\u043d\\u0430. \\u0413\\u044d\\u0445\\u0434\\u044d\\u044d \\u0431\\u04af\\u0445 \\u0442\\u043e\\u0445\\u0438\\u043e\\u043b\\u0434\\u043e\\u043b\\u0434 \\u0442\\u043e\\u0445\\u0438\\u0440\\u043e\\u0445\\u0433\\u04af\\u0439, \\u0437\\u0430\\u0440\\u0438\\u043c\\u0434\\u0430\\u0430 \\u044d\\u0441\\u0440\\u044d\\u0433 \\u043d\\u04e9\\u043b\\u04e9\\u04e9\\u0442\\u044d\\u0439 \\u0431\\u0430\\u0439\\u0436 \\u0431\\u043e\\u043b\\u043d\\u043e. \\u0411\\u043e\\u043b\\u0433\\u043e\\u043e\\u043c\\u0436\\u0442\\u043e\\u0439 \\u0430\\u0448\\u0438\\u0433\\u043b\\u0430\\u0430\\u0440\\u0430\\u0439.\\n\\u201c\\u041c\\u0438\\u043d\\u0438\\u0439 \\u043c\\u043e\\u0434\\u0435\\u043b\\u044c\\u201d, \\u201c\\u041c\\u0438\\u043d\\u0438\\u0439 \\u0448\\u04af\\u04af\\u0433\\u044d\\u044d\\u201d \\u0445\\u0443\\u0443\\u0434\\u0441\\u0430\\u043d\\u0434 \\u0431\\u0430\\u0439\\u043d\\u0433\\u0430 \\u0430\\u0448\\u0438\\u0433\\u043b\\u0430\\u0434\\u0430\\u0433 \\u043c\\u043e\\u0434\\u0435\\u043b\\u0438\\u0439\\u043d \\u0437\\u0443\\u0440\\u0430\\u0433, \\u0445\\u0443\\u0432\\u0446\\u0441\\u0430\\u0430 \\u0445\\u0430\\u0434\\u0433\\u0430\\u043b\\u0436 \\u0431\\u043e\\u043b\\u043d\\u043e. \\u0418\\u043d\\u0433\\u044d\\u0441\\u043d\\u044d\\u044d\\u0440 \\u0438\\u0436\\u0438\\u043b \\u0437\\u0443\\u0440\\u0433\\u0438\\u0439\\u0433 \\u0434\\u0430\\u0445\\u0438\\u043d \\u0434\\u0430\\u0445\\u0438\\u043d \\u043e\\u0440\\u0443\\u0443\\u043b\\u0430\\u0445 \\u0448\\u0430\\u0430\\u0440\\u0434\\u043b\\u0430\\u0433\\u0430\\u0433\\u04af\\u0439.\\n\\u0422\\u0443\\u0440\\u0448\\u0438\\u0436 \\u0431\\u0443\\u0439 \\u0445\\u0443\\u0432\\u0446\\u0430\\u0441 \\u0445\\u044d\\u0442 \\u0436\\u0438\\u0436\\u0438\\u0433 \\u0431\\u0430\\u0439\\u0436, \\u043c\\u043e\\u0434\\u0435\\u043b\\u0438\\u0439\\u043d \\u0437\\u0443\\u0440\\u0430\\u0433\\u0442 \\u04e9\\u043c\\u043d\\u04e9 \\u043d\\u044c \\u0438\\u043b \\u0433\\u0430\\u0440\\u0430\\u0430\\u0433\\u04af\\u0439 \\u0430\\u0440\\u044c\\u0441\\u043d\\u044b \\u0445\\u044d\\u0441\\u0433\\u0438\\u0439\\u0433 \\u0431\\u04af\\u0440\\u0445\\u044d\\u0445\\u0433\\u04af\\u0439 \\u0431\\u043e\\u043b \\u0442\\u0443\\u0445\\u0430\\u0439\\u043d \\u0445\\u044d\\u0441\\u0433\\u0438\\u0439\\u043d \\u0430\\u0440\\u044c\\u0441 \\u043c\\u0443\\u0443 \\u04af\\u04af\\u0441\\u0447 \\u0431\\u043e\\u043b\\u043d\\u043e. \\u0418\\u0439\\u043c \\u04af\\u0435\\u0434 \\u0438\\u043b\\u04af\\u04af \\u0441\\u0430\\u0439\\u043d \\u043c\\u043e\\u0434\\u0435\\u043b\\u044c \\u0441\\u043e\\u043d\\u0433\\u043e\\u0445 \\u044d\\u0441\\u0432\\u044d\\u043b \\u0442\\u0443\\u0445\\u0430\\u0439\\u043d \\u0445\\u044d\\u0441\\u0433\\u0438\\u0439\\u043d \\u0430\\u0440\\u044c\\u0441 \\u0438\\u043b \\u0433\\u0430\\u0440\\u0441\\u0430\\u043d \\u043c\\u043e\\u0434\\u0435\\u043b\\u0438\\u0439\\u043d \\u0437\\u0443\\u0440\\u0430\\u0433 \\u0430\\u0448\\u0438\\u0433\\u043b\\u0430\\u043d\\u0430. \\u042d\\u043d\\u044d \\u043d\\u044c \\u0438\\u0445\\u044d\\u0432\\u0447\\u043b\\u044d\\u043d \\u0431\\u0438\\u043a\\u0438\\u043d\\u0438 \\u044d\\u0441\\u0432\\u044d\\u043b \\u0434\\u043e\\u0442\\u0443\\u0443\\u0440 \\u0445\\u0443\\u0432\\u0446\\u0430\\u0441 \\u0442\\u0443\\u0440\\u0448\\u0438\\u0445 \\u04af\\u0435\\u0434 \\u0442\\u043e\\u0445\\u0438\\u043e\\u043b\\u0434\\u0434\\u043e\\u0433.\\n\\u041c\\u043e\\u0434\\u0435\\u043b\\u0438\\u0439\\u043d \\u0437\\u0443\\u0440\\u0430\\u0433 \\u044d\\u043c\\u0437\\u044d\\u0433 \\u0445\\u044d\\u0441\\u044d\\u0433 \\u0438\\u043b \\u0433\\u0430\\u0440\\u0433\\u0430\\u0441\\u0430\\u043d \\u0431\\u0430\\u0439\\u0436 \\u0431\\u043e\\u043b\\u043e\\u0445\\u0433\\u04af\\u0439. \\u0425\\u044d\\u0442 \\u0438\\u043b \\u0437\\u0430\\u0434\\u0433\\u0430\\u0439 \\u0431\\u043e\\u043b \\u0448\\u0430\\u043b\\u0433\\u0430\\u043b\\u0442\\u0430\\u0430\\u0440 \\u0443\\u043d\\u0430\\u0436, \\u04af\\u04af\\u0441\\u0433\\u044d\\u0445 \\u0431\\u043e\\u043b\\u043e\\u043c\\u0436\\u0433\\u04af\\u0439 \\u0431\\u043e\\u043b\\u043d\\u043e.\\n\\u0426\\u0430\\u0430\\u0448\\u0438\\u0434 \\u04e9\\u04e9\\u0440 \\u043f\\u043e\\u0437\\u044b\\u043d \\u043c\\u043e\\u0434\\u0435\\u043b\\u0438\\u0439\\u043d \\u0437\\u0443\\u0440\\u0430\\u0433, \\u0434\\u043e\\u0442\\u0443\\u0443\\u0440 \\u0445\\u0443\\u0432\\u0446\\u0430\\u0441, \\u0431\\u0438\\u043a\\u0438\\u043d\\u0438 \\u0437\\u044d\\u0440\\u044d\\u0433\\u0442 \\u0438\\u043b\\u04af\\u04af \\u0442\\u043e\\u0445\\u0438\\u0440\\u0441\\u043e\\u043d \\u043c\\u043e\\u0434\\u0435\\u043b\\u044c\\u04af\\u04af\\u0434\\u0438\\u0439\\u0433 \\u0448\\u0438\\u043d\\u044d\\u0447\\u0438\\u043b\\u043d\\u044d. \\u041f\\u043b\\u0430\\u0442\\u0444\\u043e\\u0440\\u043c\\u044b\\u043d \\u0437\\u0430\\u0440\\u043b\\u0430\\u043b\\u044b\\u0433 \\u0442\\u043e\\u0433\\u0442\\u043c\\u043e\\u043b \\u0434\\u0430\\u0433\\u0430\\u043d\\u0430 \\u0443\\u0443.\\n\\u041c\\u043e\\u0434\\u0435\\u043b\\u0438\\u0439\\u043d \\u0437\\u0443\\u0440\\u0430\\u0433\\u0442 \\u0445\\u0430\\u0436\\u0443\\u0443 \\u043f\\u043e\\u0437\\u0442\\u043e\\u0439 \\u0437\\u0443\\u0440\\u0430\\u0433 \\u043e\\u0440\\u0443\\u0443\\u043b\\u0436 \\u0431\\u043e\\u043b\\u043d\\u043e (\\u0443\\u0440\\u0434 \\u0437\\u0443\\u0440\\u0430\\u0433 \\u0438\\u0445\\u044d\\u0432\\u0447\\u043b\\u044d\\u043d \\u0438\\u043b\\u04af\\u04af \\u0441\\u0430\\u0439\\u043d). \\u0425\\u044d\\u0440\\u044d\\u0432 \\u0434\\u044d\\u044d\\u0434 \\u0431\\u0430 \\u0434\\u043e\\u043e\\u0434 \\u0445\\u0443\\u0432\\u0446\\u0430\\u0441\\u043d\\u044b \\u0431\\u04af\\u0440\\u044d\\u043d \\u0437\\u0443\\u0440\\u0430\\u0433 \\u0431\\u0430\\u0439\\u0445\\u0433\\u04af\\u0439 \\u0431\\u043e\\u043b \\u0445\\u0443\\u0432\\u0446\\u0430\\u0441\\u0442\\u0430\\u0439 \\u0445\\u04af\\u043d\\u0438\\u0439 \\u0437\\u0443\\u0440\\u0433\\u0438\\u0439\\u0433 \\u0448\\u0443\\u0443\\u0434 \\u043e\\u0440\\u0443\\u0443\\u043b\\u0436 \\u0431\\u043e\\u043b\\u043d\\u043e (\\u0437\\u04e9\\u0432\\u0445\\u04e9\\u043d \\u201c\\u0414\\u044d\\u044d\\u0434 \\u0437\\u0443\\u0440\\u0430\\u0433\\u201d-\\u0442 \\u043e\\u0440\\u0443\\u0443\\u043b\\u043d\\u0430, \\u201c\\u0414\\u043e\\u043e\\u0434 \\u0437\\u0443\\u0440\\u0430\\u0433\\u201d \\u0448\\u0430\\u0430\\u0440\\u0434\\u043b\\u0430\\u0433\\u0430\\u0433\\u04af\\u0439). \\u041c\\u043e\\u0434\\u0435\\u043b\\u044c \\u043d\\u044c \\u0442\\u0443\\u0445\\u0430\\u0439\\u043d \\u0445\\u04af\\u043d\\u0438\\u0439 \\u0434\\u044d\\u044d\\u0434 \\u0431\\u0430 \\u0434\\u043e\\u043e\\u0434 \\u0445\\u0443\\u0432\\u0446\\u0441\\u044b\\u0433 \\u0430\\u0432\\u0442\\u043e\\u043c\\u0430\\u0442\\u0430\\u0430\\u0440 \\u043c\\u043e\\u0434\\u0435\\u043b\\u0438\\u0439\\u043d \\u0437\\u0443\\u0440\\u0430\\u0433\\u0442 \\u04e9\\u043c\\u0441\\u0433\\u04e9\\u043d\\u04e9.\",\"zh-TW\":\"\\u7db2\\u9801\\u7248\\u5efa\\u8b70\\u4f7f\\u7528 Chrome \\u8207 Safari \\u700f\\u89bd\\u5668\\u3002\\n\\u6a21\\u7279\\u5716\\u53ef\\u4ee5\\u4e0a\\u50b3\\u5074\\u8eab\\u59ff\\u52e2\\u3001\\u5750\\u8457\\u59ff\\u52e2\\u7b49\\u5716\\u7247\\uff0c\\u4f46\\u662f\\u6b63\\u9762\\u7ad9\\u7acb\\u59ff\\u52e2\\u751f\\u6210\\u6548\\u679c\\u66f4\\u4f73\\u3002\\n\\u5982\\u679c\\u6c92\\u6709\\u55ae\\u7368\\u7684\\u670d\\u88dd\\u5716\\u7247\\uff0c\\u53ef\\u76f4\\u63a5\\u5728\\u300c\\u4e0a\\u88dd\\u5716\\u300d\\u4f4d\\u7f6e\\u4e0a\\u50b3\\u4eba\\u7269\\u7a7f\\u8457\\u7684\\u5716\\u7247\\u3002\\n\\u5982\\u679c\\u670d\\u88dd\\u5716\\u7247\\u5f88\\u96e3\\u7372\\u53d6\\uff0c\\u53ef\\u4ee5\\u4f7f\\u7528\\u300c\\u88c1\\u526a\\u4e0a\\u50b3\\u300d\\u88c1\\u526a\\uff0c\\u5f8c\\u7e8c\\u6703\\u6dfb\\u52a0\\u300c\\u53d6\\u8863\\u300d\\u529f\\u80fd\\uff0c\\u656c\\u8acb\\u671f\\u5f85\\u3002\\n\\u751f\\u6210\\u54c1\\u8cea\\u548c\\u6a21\\u578b\\u3001\\u5716\\u7247\\u54c1\\u8cea\\u3001\\u80cc\\u666f\\u6709\\u95dc\\uff0c\\u80cc\\u666f\\u8d8a\\u7c21\\u55ae\\uff0c\\u4eba\\u7269\\u5728\\u5716\\u7247\\u5360\\u6bd4\\u8d8a\\u5927\\uff0c\\u751f\\u6210\\u54c1\\u8cea\\u66f4\\u597d\\u3002\\n\\u5716\\u7247\\u7cbe\\u4fee\\u6a21\\u578b\\u6703\\u63d0\\u9ad8\\u8863\\u7269\\u3001\\u76ae\\u819a\\u7684\\u7d30\\u7bc0\\u8cea\\u611f\\uff0c\\u4f46\\u4e0d\\u9069\\u7528\\u6240\\u6709\\uff0c\\u6709\\u6642\\u5019\\u53ef\\u80fd\\u9069\\u5f97\\u5176\\u53cd\\uff0c\\u8acb\\u6ce8\\u610f\\u4f7f\\u7528\\u3002\\n\\u53ef\\u4ee5\\u5728\\u300c\\u6211\\u7684\\u6a21\\u7279\\u300d\\u548c\\u300c\\u6211\\u7684\\u8863\\u6ae5\\u300d\\u9801\\u9762\\u4e0a\\u50b3\\u4fdd\\u5b58\\u81ea\\u5df1\\u5e38\\u7528\\u7684\\u6a21\\u7279\\u5716\\u548c\\u670d\\u88dd\\uff0c\\u9019\\u6a23\\u5c31\\u4e0d\\u7528\\u6bcf\\u6b21\\u90fd\\u4e0a\\u50b3\\u540c\\u6a23\\u7684\\u5716\\u7247\\u3002\\n\\u5982\\u679c\\u8a66\\u7a7f\\u7684\\u670d\\u88dd\\u592a\\u5c0f\\u5c0e\\u81f4\\u6c92\\u6709\\u8986\\u84cb\\u6a21\\u7279\\u5716\\u539f\\u4f86\\u6c92\\u9732\\u51fa\\u76ae\\u819a\\u7684\\u4f4d\\u7f6e\\uff0c\\u5c0d\\u61c9\\u4f4d\\u7f6e\\u7684\\u76ae\\u819a\\u53ef\\u80fd\\u751f\\u6210\\u4e0d\\u597d\\uff0c\\u53ef\\u4ee5\\u5207\\u63db\\u66f4\\u597d\\u7684\\u6a21\\u578b\\u6216\\u8005\\u66f4\\u63db\\u9732\\u51fa\\u5c0d\\u61c9\\u4f4d\\u7f6e\\u76ae\\u819a\\u7684\\u6a21\\u7279\\u5716\\u3002\\u9019\\u500b\\u5834\\u666f\\u4e00\\u822c\\u662f\\u5728\\u8a66\\u7a7f\\u6bd4\\u57fa\\u5c3c\\u6216\\u8005\\u5167\\u8863\\u3002\\n\\u6a21\\u7279\\u5716\\u4e0d\\u80fd\\u9732\\u51fa\\u654f\\u611f\\u4f4d\\u7f6e\\uff0c\\u592a\\u88f8\\u9732\\u6703\\u5be9\\u6838\\u751f\\u6210\\u5931\\u6557\\u3002\\n\\u4ee5\\u5f8c\\u6703\\u66f4\\u65b0\\u5c0d\\u4e0d\\u540c\\u59ff\\u52e2\\u6a21\\u7279\\u5716\\u3001\\u5167\\u8863\\u3001\\u6bd4\\u57fa\\u5c3c\\u7b49\\u751f\\u6210\\u66f4\\u597d\\u7684\\u6a21\\u578b\\uff0c\\u8acb\\u96a8\\u6642\\u95dc\\u6ce8\\u5e73\\u53f0\\u516c\\u544a\\u3002\\n\\u6a21\\u7279\\u5716\\u53ef\\u4ee5\\u4e0a\\u50b3\\u5074\\u8eab\\u5716\\uff08\\u6b63\\u9762\\u7167\\u4e00\\u822c\\u6548\\u679c\\u66f4\\u597d\\uff09\\u3002\\u5982\\u679c\\u6c92\\u6709\\u5b8c\\u6574\\u7684\\u4e0a\\u88dd\\u548c\\u4e0b\\u88dd\\u5716\\uff0c\\u53ef\\u4ee5\\u76f4\\u63a5\\u4e0a\\u50b3\\u4eba\\u7269\\u7a7f\\u8457\\u5716\\uff08\\u4e0a\\u50b3\\u5230\\u4e0a\\u88dd\\u5716\\u5373\\u53ef\\uff0c\\u4e0d\\u7528\\u50b3\\u4e0b\\u88dd\\u5716\\uff09\\uff0c\\u6a21\\u578b\\u6703\\u81ea\\u52d5\\u628a\\u8a72\\u4eba\\u7269\\u7684\\u4e0a\\u4e0b\\u8eab\\u670d\\u88dd\\u8a66\\u7a7f\\u5230\\u6a21\\u7279\\u5716\\u4e0a\\u3002\",\"th\":\"\\u0e2a\\u0e33\\u0e2b\\u0e23\\u0e31\\u0e1a\\u0e40\\u0e27\\u0e2d\\u0e23\\u0e4c\\u0e0a\\u0e31\\u0e19\\u0e40\\u0e27\\u0e47\\u0e1a \\u0e41\\u0e19\\u0e30\\u0e19\\u0e33\\u0e43\\u0e2b\\u0e49\\u0e43\\u0e0a\\u0e49\\u0e40\\u0e1a\\u0e23\\u0e32\\u0e27\\u0e4c\\u0e40\\u0e0b\\u0e2d\\u0e23\\u0e4c Chrome \\u0e41\\u0e25\\u0e30 Safari\\n\\u0e23\\u0e39\\u0e1b\\u0e42\\u0e21\\u0e40\\u0e14\\u0e25\\u0e2a\\u0e32\\u0e21\\u0e32\\u0e23\\u0e16\\u0e2d\\u0e31\\u0e1b\\u0e42\\u0e2b\\u0e25\\u0e14\\u0e40\\u0e1b\\u0e47\\u0e19\\u0e17\\u0e48\\u0e32\\u0e2b\\u0e31\\u0e19\\u0e02\\u0e49\\u0e32\\u0e07 \\u0e17\\u0e48\\u0e32\\u0e19\\u0e31\\u0e48\\u0e07 \\u0e2f\\u0e25\\u0e2f \\u0e44\\u0e14\\u0e49 \\u0e41\\u0e15\\u0e48\\u0e17\\u0e48\\u0e32\\u0e22\\u0e37\\u0e19\\u0e15\\u0e23\\u0e07\\u0e2b\\u0e31\\u0e19\\u0e2b\\u0e19\\u0e49\\u0e32\\u0e21\\u0e31\\u0e01\\u0e43\\u0e2b\\u0e49\\u0e1c\\u0e25\\u0e25\\u0e31\\u0e1e\\u0e18\\u0e4c\\u0e14\\u0e35\\u0e01\\u0e27\\u0e48\\u0e32\\n\\u0e2b\\u0e32\\u0e01\\u0e44\\u0e21\\u0e48\\u0e21\\u0e35\\u0e23\\u0e39\\u0e1b\\u0e40\\u0e2a\\u0e37\\u0e49\\u0e2d\\u0e1c\\u0e49\\u0e32\\u0e41\\u0e22\\u0e01\\u0e0a\\u0e34\\u0e49\\u0e19 \\u0e2a\\u0e32\\u0e21\\u0e32\\u0e23\\u0e16\\u0e2d\\u0e31\\u0e1b\\u0e42\\u0e2b\\u0e25\\u0e14\\u0e23\\u0e39\\u0e1b\\u0e04\\u0e19\\u0e17\\u0e35\\u0e48\\u0e2a\\u0e27\\u0e21\\u0e0a\\u0e38\\u0e14\\u0e44\\u0e27\\u0e49\\u0e43\\u0e19\\u0e0a\\u0e48\\u0e2d\\u0e07 \\\"\\u0e23\\u0e39\\u0e1b\\u0e1a\\u0e19\\\" \\u0e44\\u0e14\\u0e49\\u0e42\\u0e14\\u0e22\\u0e15\\u0e23\\u0e07\\n\\u0e2b\\u0e32\\u0e01\\u0e2b\\u0e32\\u0e23\\u0e39\\u0e1b\\u0e40\\u0e2a\\u0e37\\u0e49\\u0e2d\\u0e1c\\u0e49\\u0e32\\u0e44\\u0e14\\u0e49\\u0e22\\u0e32\\u0e01 \\u0e2a\\u0e32\\u0e21\\u0e32\\u0e23\\u0e16\\u0e43\\u0e0a\\u0e49 \\\"\\u0e2d\\u0e31\\u0e1b\\u0e42\\u0e2b\\u0e25\\u0e14\\u0e41\\u0e1a\\u0e1a\\u0e04\\u0e23\\u0e2d\\u0e1b\\\" \\u0e01\\u0e48\\u0e2d\\u0e19 \\u0e42\\u0e14\\u0e22\\u0e1f\\u0e35\\u0e40\\u0e08\\u0e2d\\u0e23\\u0e4c \\\"\\u0e14\\u0e36\\u0e07\\u0e40\\u0e2a\\u0e37\\u0e49\\u0e2d\\u0e1c\\u0e49\\u0e32\\\" \\u0e08\\u0e30\\u0e40\\u0e1e\\u0e34\\u0e48\\u0e21\\u0e20\\u0e32\\u0e22\\u0e2b\\u0e25\\u0e31\\u0e07\\n\\u0e04\\u0e38\\u0e13\\u0e20\\u0e32\\u0e1e\\u0e1c\\u0e25\\u0e25\\u0e31\\u0e1e\\u0e18\\u0e4c\\u0e02\\u0e36\\u0e49\\u0e19\\u0e01\\u0e31\\u0e1a\\u0e42\\u0e21\\u0e40\\u0e14\\u0e25 \\u0e04\\u0e38\\u0e13\\u0e20\\u0e32\\u0e1e\\u0e23\\u0e39\\u0e1b \\u0e41\\u0e25\\u0e30\\u0e1e\\u0e37\\u0e49\\u0e19\\u0e2b\\u0e25\\u0e31\\u0e07 \\u0e22\\u0e34\\u0e48\\u0e07\\u0e1e\\u0e37\\u0e49\\u0e19\\u0e2b\\u0e25\\u0e31\\u0e07\\u0e40\\u0e23\\u0e35\\u0e22\\u0e1a\\u0e41\\u0e25\\u0e30\\u0e04\\u0e19\\u0e21\\u0e35\\u0e2a\\u0e31\\u0e14\\u0e2a\\u0e48\\u0e27\\u0e19\\u0e43\\u0e2b\\u0e0d\\u0e48\\u0e43\\u0e19\\u0e20\\u0e32\\u0e1e \\u0e1c\\u0e25\\u0e25\\u0e31\\u0e1e\\u0e18\\u0e4c\\u0e21\\u0e31\\u0e01\\u0e22\\u0e34\\u0e48\\u0e07\\u0e14\\u0e35\\n\\u0e42\\u0e21\\u0e40\\u0e14\\u0e25\\u0e1b\\u0e23\\u0e31\\u0e1a\\u0e20\\u0e32\\u0e1e\\u0e25\\u0e30\\u0e40\\u0e2d\\u0e35\\u0e22\\u0e14\\u0e0a\\u0e48\\u0e27\\u0e22\\u0e40\\u0e1e\\u0e34\\u0e48\\u0e21\\u0e23\\u0e32\\u0e22\\u0e25\\u0e30\\u0e40\\u0e2d\\u0e35\\u0e22\\u0e14\\u0e02\\u0e2d\\u0e07\\u0e40\\u0e2a\\u0e37\\u0e49\\u0e2d\\u0e1c\\u0e49\\u0e32\\u0e41\\u0e25\\u0e30\\u0e1c\\u0e34\\u0e27 \\u0e41\\u0e15\\u0e48\\u0e44\\u0e21\\u0e48\\u0e40\\u0e2b\\u0e21\\u0e32\\u0e30\\u0e01\\u0e31\\u0e1a\\u0e17\\u0e38\\u0e01\\u0e01\\u0e23\\u0e13\\u0e35 \\u0e41\\u0e25\\u0e30\\u0e1a\\u0e32\\u0e07\\u0e04\\u0e23\\u0e31\\u0e49\\u0e07\\u0e2d\\u0e32\\u0e08\\u0e43\\u0e2b\\u0e49\\u0e1c\\u0e25\\u0e15\\u0e23\\u0e07\\u0e02\\u0e49\\u0e32\\u0e21 \\u0e04\\u0e27\\u0e23\\u0e43\\u0e0a\\u0e49\\u0e2d\\u0e22\\u0e48\\u0e32\\u0e07\\u0e23\\u0e30\\u0e21\\u0e31\\u0e14\\u0e23\\u0e30\\u0e27\\u0e31\\u0e07\\n\\u0e04\\u0e38\\u0e13\\u0e2a\\u0e32\\u0e21\\u0e32\\u0e23\\u0e16\\u0e2d\\u0e31\\u0e1b\\u0e42\\u0e2b\\u0e25\\u0e14\\u0e41\\u0e25\\u0e30\\u0e1a\\u0e31\\u0e19\\u0e17\\u0e36\\u0e01\\u0e23\\u0e39\\u0e1b\\u0e42\\u0e21\\u0e40\\u0e14\\u0e25\\u0e41\\u0e25\\u0e30\\u0e40\\u0e2a\\u0e37\\u0e49\\u0e2d\\u0e1c\\u0e49\\u0e32\\u0e17\\u0e35\\u0e48\\u0e43\\u0e0a\\u0e49\\u0e1a\\u0e48\\u0e2d\\u0e22\\u0e44\\u0e27\\u0e49\\u0e43\\u0e19 \\\"\\u0e42\\u0e21\\u0e40\\u0e14\\u0e25\\u0e02\\u0e2d\\u0e07\\u0e09\\u0e31\\u0e19\\\" \\u0e41\\u0e25\\u0e30 \\\"\\u0e15\\u0e39\\u0e49\\u0e40\\u0e2a\\u0e37\\u0e49\\u0e2d\\u0e1c\\u0e49\\u0e32\\u0e02\\u0e2d\\u0e07\\u0e09\\u0e31\\u0e19\\\" \\u0e40\\u0e1e\\u0e37\\u0e48\\u0e2d\\u0e44\\u0e21\\u0e48\\u0e15\\u0e49\\u0e2d\\u0e07\\u0e2d\\u0e31\\u0e1b\\u0e42\\u0e2b\\u0e25\\u0e14\\u0e0b\\u0e49\\u0e33\\u0e17\\u0e38\\u0e01\\u0e04\\u0e23\\u0e31\\u0e49\\u0e07\\n\\u0e2b\\u0e32\\u0e01\\u0e40\\u0e2a\\u0e37\\u0e49\\u0e2d\\u0e1c\\u0e49\\u0e32\\u0e17\\u0e35\\u0e48\\u0e25\\u0e2d\\u0e07\\u0e21\\u0e35\\u0e02\\u0e19\\u0e32\\u0e14\\u0e40\\u0e25\\u0e47\\u0e01\\u0e08\\u0e19\\u0e44\\u0e21\\u0e48\\u0e04\\u0e23\\u0e2d\\u0e1a\\u0e04\\u0e25\\u0e38\\u0e21\\u0e1a\\u0e23\\u0e34\\u0e40\\u0e27\\u0e13\\u0e1c\\u0e34\\u0e27\\u0e17\\u0e35\\u0e48\\u0e40\\u0e14\\u0e34\\u0e21\\u0e44\\u0e21\\u0e48\\u0e44\\u0e14\\u0e49\\u0e40\\u0e1b\\u0e34\\u0e14\\u0e40\\u0e1c\\u0e22\\u0e43\\u0e19\\u0e23\\u0e39\\u0e1b\\u0e42\\u0e21\\u0e40\\u0e14\\u0e25 \\u0e1c\\u0e34\\u0e27\\u0e1a\\u0e23\\u0e34\\u0e40\\u0e27\\u0e13\\u0e19\\u0e31\\u0e49\\u0e19\\u0e2d\\u0e32\\u0e08\\u0e2a\\u0e23\\u0e49\\u0e32\\u0e07\\u0e44\\u0e14\\u0e49\\u0e44\\u0e21\\u0e48\\u0e14\\u0e35 \\u0e01\\u0e23\\u0e13\\u0e35\\u0e19\\u0e35\\u0e49 (\\u0e21\\u0e31\\u0e01\\u0e40\\u0e01\\u0e34\\u0e14\\u0e01\\u0e31\\u0e1a\\u0e1a\\u0e34\\u0e01\\u0e34\\u0e19\\u0e35/\\u0e0a\\u0e38\\u0e14\\u0e0a\\u0e31\\u0e49\\u0e19\\u0e43\\u0e19) \\u0e43\\u0e2b\\u0e49\\u0e40\\u0e1b\\u0e25\\u0e35\\u0e48\\u0e22\\u0e19\\u0e42\\u0e21\\u0e40\\u0e14\\u0e25\\u0e17\\u0e35\\u0e48\\u0e14\\u0e35\\u0e01\\u0e27\\u0e48\\u0e32 \\u0e2b\\u0e23\\u0e37\\u0e2d\\u0e43\\u0e0a\\u0e49\\u0e23\\u0e39\\u0e1b\\u0e42\\u0e21\\u0e40\\u0e14\\u0e25\\u0e17\\u0e35\\u0e48\\u0e40\\u0e1b\\u0e34\\u0e14\\u0e1c\\u0e34\\u0e27\\u0e43\\u0e19\\u0e15\\u0e33\\u0e41\\u0e2b\\u0e19\\u0e48\\u0e07\\u0e19\\u0e31\\u0e49\\u0e19\\n\\u0e23\\u0e39\\u0e1b\\u0e42\\u0e21\\u0e40\\u0e14\\u0e25\\u0e15\\u0e49\\u0e2d\\u0e07\\u0e44\\u0e21\\u0e48\\u0e40\\u0e1c\\u0e22\\u0e15\\u0e33\\u0e41\\u0e2b\\u0e19\\u0e48\\u0e07\\u0e17\\u0e35\\u0e48\\u0e2d\\u0e48\\u0e2d\\u0e19\\u0e44\\u0e2b\\u0e27 \\u0e2b\\u0e32\\u0e01\\u0e42\\u0e1b\\u0e4a\\u0e21\\u0e32\\u0e01\\u0e40\\u0e01\\u0e34\\u0e19\\u0e44\\u0e1b\\u0e2d\\u0e32\\u0e08\\u0e15\\u0e23\\u0e27\\u0e08\\u0e44\\u0e21\\u0e48\\u0e1c\\u0e48\\u0e32\\u0e19\\u0e41\\u0e25\\u0e30\\u0e2a\\u0e23\\u0e49\\u0e32\\u0e07\\u0e44\\u0e21\\u0e48\\u0e2a\\u0e33\\u0e40\\u0e23\\u0e47\\u0e08\\n\\u0e43\\u0e19\\u0e2d\\u0e19\\u0e32\\u0e04\\u0e15\\u0e08\\u0e30\\u0e21\\u0e35\\u0e01\\u0e32\\u0e23\\u0e2d\\u0e31\\u0e1b\\u0e40\\u0e14\\u0e15\\u0e42\\u0e21\\u0e40\\u0e14\\u0e25\\u0e17\\u0e35\\u0e48\\u0e14\\u0e35\\u0e01\\u0e27\\u0e48\\u0e32\\u0e2a\\u0e33\\u0e2b\\u0e23\\u0e31\\u0e1a\\u0e17\\u0e48\\u0e32\\u0e17\\u0e32\\u0e07\\u0e15\\u0e48\\u0e32\\u0e07 \\u0e46 \\u0e41\\u0e25\\u0e30\\u0e2a\\u0e33\\u0e2b\\u0e23\\u0e31\\u0e1a\\u0e0a\\u0e38\\u0e14\\u0e0a\\u0e31\\u0e49\\u0e19\\u0e43\\u0e19/\\u0e1a\\u0e34\\u0e01\\u0e34\\u0e19\\u0e35 \\u0e42\\u0e1b\\u0e23\\u0e14\\u0e15\\u0e34\\u0e14\\u0e15\\u0e32\\u0e21\\u0e1b\\u0e23\\u0e30\\u0e01\\u0e32\\u0e28\\u0e41\\u0e1e\\u0e25\\u0e15\\u0e1f\\u0e2d\\u0e23\\u0e4c\\u0e21\\n\\u0e2a\\u0e32\\u0e21\\u0e32\\u0e23\\u0e16\\u0e2d\\u0e31\\u0e1b\\u0e42\\u0e2b\\u0e25\\u0e14\\u0e23\\u0e39\\u0e1b\\u0e42\\u0e21\\u0e40\\u0e14\\u0e25\\u0e21\\u0e38\\u0e21\\u0e02\\u0e49\\u0e32\\u0e07\\u0e44\\u0e14\\u0e49 (\\u0e42\\u0e14\\u0e22\\u0e17\\u0e31\\u0e48\\u0e27\\u0e44\\u0e1b\\u0e23\\u0e39\\u0e1b\\u0e14\\u0e49\\u0e32\\u0e19\\u0e2b\\u0e19\\u0e49\\u0e32\\u0e08\\u0e30\\u0e44\\u0e14\\u0e49\\u0e1c\\u0e25\\u0e14\\u0e35\\u0e01\\u0e27\\u0e48\\u0e32) \\u0e2b\\u0e32\\u0e01\\u0e44\\u0e21\\u0e48\\u0e21\\u0e35\\u0e23\\u0e39\\u0e1b\\u0e17\\u0e48\\u0e2d\\u0e19\\u0e1a\\u0e19\\u0e41\\u0e25\\u0e30\\u0e17\\u0e48\\u0e2d\\u0e19\\u0e25\\u0e48\\u0e32\\u0e07\\u0e04\\u0e23\\u0e1a \\u0e2a\\u0e32\\u0e21\\u0e32\\u0e23\\u0e16\\u0e2d\\u0e31\\u0e1b\\u0e42\\u0e2b\\u0e25\\u0e14\\u0e23\\u0e39\\u0e1b\\u0e04\\u0e19\\u0e17\\u0e35\\u0e48\\u0e2a\\u0e27\\u0e21\\u0e0a\\u0e38\\u0e14\\u0e44\\u0e14\\u0e49\\u0e42\\u0e14\\u0e22\\u0e15\\u0e23\\u0e07 (\\u0e2d\\u0e31\\u0e1b\\u0e42\\u0e2b\\u0e25\\u0e14\\u0e40\\u0e09\\u0e1e\\u0e32\\u0e30\\u0e0a\\u0e48\\u0e2d\\u0e07 \\\"\\u0e23\\u0e39\\u0e1b\\u0e1a\\u0e19\\\" \\u0e44\\u0e21\\u0e48\\u0e15\\u0e49\\u0e2d\\u0e07\\u0e2d\\u0e31\\u0e1b\\u0e42\\u0e2b\\u0e25\\u0e14\\u0e17\\u0e48\\u0e2d\\u0e19\\u0e25\\u0e48\\u0e32\\u0e07) \\u0e23\\u0e30\\u0e1a\\u0e1a\\u0e08\\u0e30\\u0e25\\u0e2d\\u0e07\\u0e2a\\u0e27\\u0e21\\u0e40\\u0e2a\\u0e37\\u0e49\\u0e2d\\u0e1c\\u0e49\\u0e32\\u0e17\\u0e48\\u0e2d\\u0e19\\u0e1a\\u0e19\\u0e41\\u0e25\\u0e30\\u0e17\\u0e48\\u0e2d\\u0e19\\u0e25\\u0e48\\u0e32\\u0e07\\u0e02\\u0e2d\\u0e07\\u0e04\\u0e19\\u0e19\\u0e31\\u0e49\\u0e19\\u0e43\\u0e2b\\u0e49\\u0e42\\u0e21\\u0e40\\u0e14\\u0e25\\u0e42\\u0e14\\u0e22\\u0e2d\\u0e31\\u0e15\\u0e42\\u0e19\\u0e21\\u0e31\\u0e15\\u0e34\",\"hi\":\"\\u0935\\u0947\\u092c \\u0938\\u0902\\u0938\\u094d\\u0915\\u0930\\u0923 \\u0915\\u0947 \\u0932\\u093f\\u090f Chrome \\u0914\\u0930 Safari \\u092c\\u094d\\u0930\\u093e\\u0909\\u091c\\u093c\\u0930 \\u0915\\u093e \\u0909\\u092a\\u092f\\u094b\\u0917 \\u0915\\u0930\\u0928\\u0947 \\u0915\\u0940 \\u0938\\u093f\\u092b\\u093e\\u0930\\u093f\\u0936 \\u0915\\u0940 \\u091c\\u093e\\u0924\\u0940 \\u0939\\u0948\\u0964\\n\\u092e\\u0949\\u0921\\u0932 \\u092b\\u094b\\u091f\\u094b \\u092e\\u0947\\u0902 \\u0938\\u093e\\u0907\\u0921 \\u092a\\u094b\\u091c\\u093c, \\u092c\\u0948\\u0920\\u0940 \\u0939\\u0941\\u0908 \\u092a\\u094b\\u091c\\u093c \\u0906\\u0926\\u093f \\u0905\\u092a\\u0932\\u094b\\u0921 \\u0915\\u0940 \\u091c\\u093e \\u0938\\u0915\\u0924\\u0940 \\u0939\\u0948\\u0902, \\u0932\\u0947\\u0915\\u093f\\u0928 \\u0938\\u093e\\u092e\\u0928\\u0947 \\u0938\\u0940\\u0927\\u093e \\u0916\\u0921\\u093c\\u093e \\u092a\\u094b\\u091c\\u093c \\u0938\\u093e\\u092e\\u093e\\u0928\\u094d\\u092f\\u0924\\u0903 \\u092c\\u0947\\u0939\\u0924\\u0930 \\u092a\\u0930\\u093f\\u0923\\u093e\\u092e \\u0926\\u0947\\u0924\\u093e \\u0939\\u0948\\u0964\\n\\u0905\\u0917\\u0930 \\u0905\\u0932\\u0917 \\u0915\\u092a\\u0921\\u093c\\u094b\\u0902 \\u0915\\u0940 \\u092b\\u094b\\u091f\\u094b \\u0928\\u0939\\u0940\\u0902 \\u0939\\u0948, \\u0924\\u094b \\\"\\u090a\\u092a\\u0930\\u0940 \\u091a\\u093f\\u0924\\u094d\\u0930\\\" \\u0938\\u094d\\u0932\\u0949\\u091f \\u092e\\u0947\\u0902 \\u092a\\u0939\\u0928\\u0947 \\u0939\\u0941\\u090f \\u0935\\u094d\\u092f\\u0915\\u094d\\u0924\\u093f \\u0915\\u0940 \\u092b\\u094b\\u091f\\u094b \\u0938\\u0940\\u0927\\u0947 \\u0905\\u092a\\u0932\\u094b\\u0921 \\u0915\\u0930 \\u0938\\u0915\\u0924\\u0947 \\u0939\\u0948\\u0902\\u0964\\n\\u0905\\u0917\\u0930 \\u0915\\u092a\\u0921\\u093c\\u094b\\u0902 \\u0915\\u0940 \\u092b\\u094b\\u091f\\u094b \\u092e\\u093f\\u0932\\u0928\\u093e \\u0915\\u0920\\u093f\\u0928 \\u0939\\u094b, \\u0924\\u094b \\\"\\u0915\\u094d\\u0930\\u0949\\u092a \\u0905\\u092a\\u0932\\u094b\\u0921\\\" \\u0915\\u093e \\u0909\\u092a\\u092f\\u094b\\u0917 \\u0915\\u0930\\u0947\\u0902\\u0964 \\u0906\\u0917\\u0947 \\\"\\u0915\\u092a\\u0921\\u093c\\u093e \\u0928\\u093f\\u0915\\u093e\\u0932\\u0947\\u0902\\\" \\u092b\\u0940\\u091a\\u0930 \\u091c\\u094b\\u0921\\u093c\\u093e \\u091c\\u093e\\u090f\\u0917\\u093e\\u0964\\n\\u091c\\u0928\\u0930\\u0947\\u0936\\u0928 \\u0915\\u094d\\u0935\\u093e\\u0932\\u093f\\u091f\\u0940 \\u092e\\u0949\\u0921\\u0932, \\u0907\\u092e\\u0947\\u091c \\u0915\\u094d\\u0935\\u093e\\u0932\\u093f\\u091f\\u0940 \\u0914\\u0930 \\u092c\\u0948\\u0915\\u0917\\u094d\\u0930\\u093e\\u0909\\u0902\\u0921 \\u092a\\u0930 \\u0928\\u093f\\u0930\\u094d\\u092d\\u0930 \\u0939\\u0948\\u0964 \\u092c\\u0948\\u0915\\u0917\\u094d\\u0930\\u093e\\u0909\\u0902\\u0921 \\u091c\\u093f\\u0924\\u0928\\u093e \\u0938\\u0930\\u0932 \\u0914\\u0930 \\u0935\\u094d\\u092f\\u0915\\u094d\\u0924\\u093f \\u0915\\u093e \\u0905\\u0928\\u0941\\u092a\\u093e\\u0924 \\u091c\\u093f\\u0924\\u0928\\u093e \\u092c\\u0921\\u093c\\u093e \\u0939\\u094b\\u0917\\u093e, \\u092a\\u0930\\u093f\\u0923\\u093e\\u092e \\u0909\\u0924\\u0928\\u0947 \\u092c\\u0947\\u0939\\u0924\\u0930 \\u0939\\u094b\\u0902\\u0917\\u0947\\u0964\\n\\u0907\\u092e\\u0947\\u091c \\u0930\\u093f\\u092b\\u093e\\u0907\\u0928\\u0930 \\u092e\\u0949\\u0921\\u0932 \\u0915\\u092a\\u0921\\u093c\\u094b\\u0902 \\u0914\\u0930 \\u0924\\u094d\\u0935\\u091a\\u093e \\u0915\\u0940 \\u0921\\u093f\\u091f\\u0947\\u0932 \\u092c\\u0922\\u093c\\u093e \\u0938\\u0915\\u0924\\u093e \\u0939\\u0948, \\u0932\\u0947\\u0915\\u093f\\u0928 \\u0939\\u0930 \\u0915\\u0947\\u0938 \\u092e\\u0947\\u0902 \\u0909\\u092a\\u092f\\u0941\\u0915\\u094d\\u0924 \\u0928\\u0939\\u0940\\u0902 \\u0939\\u094b\\u0924\\u093e \\u0914\\u0930 \\u0915\\u092d\\u0940-\\u0915\\u092d\\u0940 \\u0909\\u0932\\u094d\\u091f\\u093e \\u0905\\u0938\\u0930 \\u092d\\u0940 \\u0926\\u0947 \\u0938\\u0915\\u0924\\u093e \\u0939\\u0948; \\u0938\\u093e\\u0935\\u0927\\u093e\\u0928\\u0940 \\u0938\\u0947 \\u0909\\u092a\\u092f\\u094b\\u0917 \\u0915\\u0930\\u0947\\u0902\\u0964\\n\\u0906\\u092a \\\"\\u092e\\u0947\\u0930\\u0947 \\u092e\\u0949\\u0921\\u0932\\\" \\u0914\\u0930 \\\"\\u092e\\u0947\\u0930\\u0940 \\u0905\\u0932\\u092e\\u093e\\u0930\\u0940\\\" \\u092e\\u0947\\u0902 \\u0905\\u0915\\u094d\\u0938\\u0930 \\u0909\\u092a\\u092f\\u094b\\u0917 \\u0939\\u094b\\u0928\\u0947 \\u0935\\u093e\\u0932\\u0940 \\u092e\\u0949\\u0921\\u0932/\\u0915\\u092a\\u0921\\u093c\\u094b\\u0902 \\u0915\\u0940 \\u0924\\u0938\\u094d\\u0935\\u0940\\u0930\\u0947\\u0902 \\u0938\\u0947\\u0935 \\u0915\\u0930 \\u0938\\u0915\\u0924\\u0947 \\u0939\\u0948\\u0902, \\u0924\\u093e\\u0915\\u093f \\u0939\\u0930 \\u092c\\u093e\\u0930 \\u0935\\u0939\\u0940 \\u0924\\u0938\\u094d\\u0935\\u0940\\u0930\\u0947\\u0902 \\u0926\\u094b\\u092c\\u093e\\u0930\\u093e \\u0905\\u092a\\u0932\\u094b\\u0921 \\u0928 \\u0915\\u0930\\u0928\\u0940 \\u092a\\u0921\\u093c\\u0947\\u0902\\u0964\\n\\u0905\\u0917\\u0930 \\u091f\\u094d\\u0930\\u093e\\u0908-\\u0911\\u0928 \\u0915\\u092a\\u0921\\u093c\\u093e \\u092c\\u0939\\u0941\\u0924 \\u091b\\u094b\\u091f\\u093e \\u0939\\u0948 \\u0914\\u0930 \\u092e\\u0949\\u0921\\u0932 \\u092b\\u094b\\u091f\\u094b \\u0915\\u0947 \\u0909\\u0928 \\u0939\\u093f\\u0938\\u094d\\u0938\\u094b\\u0902 \\u0915\\u094b \\u0915\\u0935\\u0930 \\u0928\\u0939\\u0940\\u0902 \\u0915\\u0930\\u0924\\u093e \\u091c\\u094b \\u092a\\u0939\\u0932\\u0947 \\u090f\\u0915\\u094d\\u0938\\u092a\\u094b\\u091c\\u093c \\u0928\\u0939\\u0940\\u0902 \\u0925\\u0947, \\u0924\\u094b \\u0909\\u0928 \\u0939\\u093f\\u0938\\u094d\\u0938\\u094b\\u0902 \\u0915\\u0940 \\u0924\\u094d\\u0935\\u091a\\u093e \\u0916\\u0930\\u093e\\u092c \\u091c\\u0928\\u0930\\u0947\\u091f \\u0939\\u094b \\u0938\\u0915\\u0924\\u0940 \\u0939\\u0948\\u0964 \\u092f\\u0939 \\u0938\\u094d\\u0925\\u093f\\u0924\\u093f \\u0906\\u092e\\u0924\\u094c\\u0930 \\u092a\\u0930 \\u092c\\u093f\\u0915\\u093f\\u0928\\u0940/\\u0907\\u0928\\u0930\\u0935\\u093f\\u092f\\u0930 \\u091f\\u094d\\u0930\\u093e\\u0908-\\u0911\\u0928 \\u092e\\u0947\\u0902 \\u0939\\u094b\\u0924\\u0940 \\u0939\\u0948; \\u0910\\u0938\\u0947 \\u092e\\u0947\\u0902 \\u092c\\u0947\\u0939\\u0924\\u0930 \\u092e\\u0949\\u0921\\u0932 \\u091a\\u0941\\u0928\\u0947\\u0902 \\u092f\\u093e \\u0909\\u0938 \\u0939\\u093f\\u0938\\u094d\\u0938\\u0947 \\u0915\\u0940 \\u0924\\u094d\\u0935\\u091a\\u093e \\u090f\\u0915\\u094d\\u0938\\u092a\\u094b\\u091c\\u093c \\u0915\\u0930\\u0928\\u0947 \\u0935\\u093e\\u0932\\u0940 \\u092e\\u0949\\u0921\\u0932 \\u092b\\u094b\\u091f\\u094b \\u0907\\u0938\\u094d\\u0924\\u0947\\u092e\\u093e\\u0932 \\u0915\\u0930\\u0947\\u0902\\u0964\\n\\u092e\\u0949\\u0921\\u0932 \\u092b\\u094b\\u091f\\u094b \\u092e\\u0947\\u0902 \\u0938\\u0902\\u0935\\u0947\\u0926\\u0928\\u0936\\u0940\\u0932 \\u0939\\u093f\\u0938\\u094d\\u0938\\u0947 \\u0928\\u0939\\u0940\\u0902 \\u0926\\u093f\\u0916\\u0928\\u0947 \\u091a\\u093e\\u0939\\u093f\\u090f\\u0964 \\u092c\\u0939\\u0941\\u0924 \\u0905\\u0927\\u093f\\u0915 \\u0916\\u0941\\u0932\\u0940 \\u092b\\u094b\\u091f\\u094b \\u0938\\u092e\\u0940\\u0915\\u094d\\u0937\\u093e \\u092e\\u0947\\u0902 \\u092b\\u0947\\u0932 \\u0939\\u094b\\u0915\\u0930 \\u091c\\u0928\\u0930\\u0947\\u0936\\u0928 \\u0935\\u093f\\u092b\\u0932 \\u0915\\u0930 \\u0938\\u0915\\u0924\\u0940 \\u0939\\u0948\\u0964\\n\\u0906\\u0917\\u0947 \\u0905\\u0932\\u0917-\\u0905\\u0932\\u0917 \\u092a\\u094b\\u091c\\u093c, \\u0907\\u0928\\u0930\\u0935\\u093f\\u092f\\u0930 \\u0914\\u0930 \\u092c\\u093f\\u0915\\u093f\\u0928\\u0940 \\u091c\\u0948\\u0938\\u0947 \\u092a\\u0930\\u093f\\u0926\\u0943\\u0936\\u094d\\u092f\\u094b\\u0902 \\u0915\\u0947 \\u0932\\u093f\\u090f \\u092c\\u0947\\u0939\\u0924\\u0930 \\u092e\\u0949\\u0921\\u0932 \\u0905\\u092a\\u0921\\u0947\\u091f \\u0915\\u093f\\u090f \\u091c\\u093e\\u090f\\u0902\\u0917\\u0947; \\u092a\\u094d\\u0932\\u0947\\u091f\\u092b\\u093c\\u0949\\u0930\\u094d\\u092e \\u0918\\u094b\\u0937\\u0923\\u093e\\u090f\\u0901 \\u0926\\u0947\\u0916\\u0924\\u0947 \\u0930\\u0939\\u0947\\u0902\\u0964\\n\\u092e\\u0949\\u0921\\u0932 \\u092b\\u094b\\u091f\\u094b \\u0938\\u093e\\u0907\\u0921-\\u0935\\u094d\\u092f\\u0942 \\u092e\\u0947\\u0902 \\u092d\\u0940 \\u0905\\u092a\\u0932\\u094b\\u0921 \\u0915\\u0940 \\u091c\\u093e \\u0938\\u0915\\u0924\\u0940 \\u0939\\u0948 (\\u092b\\u094d\\u0930\\u0902\\u091f-\\u0935\\u094d\\u092f\\u0942 \\u0906\\u092e\\u0924\\u094c\\u0930 \\u092a\\u0930 \\u092c\\u0947\\u0939\\u0924\\u0930 \\u0930\\u0939\\u0924\\u093e \\u0939\\u0948)\\u0964 \\u0905\\u0917\\u0930 \\u0906\\u092a\\u0915\\u0947 \\u092a\\u093e\\u0938 \\u092a\\u0942\\u0930\\u093e \\u090a\\u092a\\u0930\\u0940 \\u0914\\u0930 \\u0928\\u093f\\u091a\\u0932\\u093e \\u0915\\u092a\\u0921\\u093c\\u094b\\u0902 \\u0915\\u093e \\u0938\\u0947\\u091f \\u0928\\u0939\\u0940\\u0902 \\u0939\\u0948, \\u0924\\u094b \\u092a\\u0939\\u0928\\u0947 \\u0939\\u0941\\u090f \\u0935\\u094d\\u092f\\u0915\\u094d\\u0924\\u093f \\u0915\\u0940 \\u092b\\u094b\\u091f\\u094b \\u0938\\u0940\\u0927\\u0947 \\u0905\\u092a\\u0932\\u094b\\u0921 \\u0915\\u0930\\u0947\\u0902 (\\u0938\\u093f\\u0930\\u094d\\u092b \\\"\\u090a\\u092a\\u0930\\u0940 \\u091a\\u093f\\u0924\\u094d\\u0930\\\" \\u092e\\u0947\\u0902, \\\"\\u0928\\u093f\\u091a\\u0932\\u093e \\u091a\\u093f\\u0924\\u094d\\u0930\\\" \\u091c\\u0930\\u0942\\u0930\\u0940 \\u0928\\u0939\\u0940\\u0902)\\u0964 \\u092e\\u0949\\u0921\\u0932 \\u0909\\u0938 \\u0935\\u094d\\u092f\\u0915\\u094d\\u0924\\u093f \\u0915\\u0947 \\u090a\\u092a\\u0930\\u0940 \\u0914\\u0930 \\u0928\\u093f\\u091a\\u0932\\u0947 \\u0915\\u092a\\u0921\\u093c\\u094b\\u0902 \\u0915\\u094b \\u0905\\u092a\\u0928\\u0947-\\u0906\\u092a \\u092e\\u0949\\u0921\\u0932 \\u092b\\u094b\\u091f\\u094b \\u092a\\u0930 \\u091f\\u094d\\u0930\\u093e\\u0908 \\u0915\\u0930\\u0947\\u0917\\u093e\\u0964\",\"id\":\"Untuk versi web, disarankan menggunakan browser Chrome dan Safari.\\nFoto model boleh berupa pose samping, duduk, dan sebagainya, tetapi pose berdiri menghadap depan biasanya memberi hasil lebih baik.\\nJika tidak ada foto pakaian terpisah, Anda bisa langsung unggah foto orang yang sedang memakai pakaian pada slot \\\"Atasan\\\".\\nJika foto pakaian sulit didapat, gunakan \\\"Unggah Potong\\\" terlebih dahulu. Fitur \\\"Ekstraksi Pakaian\\\" akan ditambahkan nanti.\\nKualitas hasil bergantung pada model, kualitas gambar, dan latar belakang. Semakin sederhana latar belakang dan semakin besar proporsi orang di gambar, biasanya hasil semakin baik.\\nModel penyempurna gambar dapat meningkatkan detail pakaian dan kulit, tetapi tidak cocok untuk semua kasus dan kadang bisa memberi hasil sebaliknya. Gunakan dengan hati-hati.\\nAnda dapat menyimpan foto model dan pakaian yang sering dipakai di halaman \\\"Model Saya\\\" dan \\\"Lemari Saya\\\" agar tidak perlu unggah gambar yang sama setiap kali.\\nJika pakaian try-on terlalu kecil sehingga tidak menutupi area kulit yang sebelumnya tidak terekspos pada foto model, kulit di area tersebut bisa tergenerasi kurang baik. Kondisi ini umum saat mencoba bikini atau pakaian dalam; Anda bisa mengganti ke model yang lebih baik atau memakai foto model yang mengekspos area kulit terkait.\\nFoto model tidak boleh menampilkan area sensitif. Foto yang terlalu terbuka dapat gagal moderasi dan gagal generate.\\nKe depannya kami akan terus memperbarui model yang lebih baik untuk berbagai pose model serta skenario pakaian dalam/bikini. Pantau pengumuman platform.\\nFoto model bisa diunggah dari samping (foto depan biasanya lebih baik). Jika Anda tidak memiliki gambar atasan dan bawahan yang lengkap, Anda dapat langsung unggah foto orang yang sedang berpakaian (cukup ke slot \\\"Atasan\\\", tidak perlu slot bawahan). Model akan otomatis mencoba pakaian atas dan bawah orang tersebut pada foto model.\",\"vi\":\"\\u0110\\u1ed1i v\\u1edbi phi\\u00ean b\\u1ea3n web, n\\u00ean s\\u1eed d\\u1ee5ng Chrome v\\u00e0 Safari.\\nB\\u1ea1n c\\u00f3 th\\u1ec3 t\\u1ea3i l\\u00ean \\u1ea3nh ng\\u01b0\\u1eddi m\\u1eabu nh\\u00ecn t\\u1eeb b\\u00ean c\\u1ea1nh ho\\u1eb7c ng\\u01b0\\u1eddi m\\u1eabu ng\\u1ed3i, nh\\u01b0ng t\\u01b0 th\\u1ebf \\u0111\\u1ee9ng quay m\\u1eb7t v\\u1ec1 ph\\u00eda tr\\u01b0\\u1edbc th\\u01b0\\u1eddng mang l\\u1ea1i k\\u1ebft qu\\u1ea3 t\\u1ed1t h\\u01a1n.\\nN\\u1ebfu b\\u1ea1n kh\\u00f4ng c\\u00f3 h\\u00ecnh \\u1ea3nh qu\\u1ea7n \\u00e1o ri\\u00eang bi\\u1ec7t, b\\u1ea1n c\\u00f3 th\\u1ec3 t\\u1ea3i tr\\u1ef1c ti\\u1ebfp h\\u00ecnh \\u1ea3nh ng\\u01b0\\u1eddi m\\u1eb7c qu\\u1ea7n \\u00e1o l\\u00ean \\u00f4 \\\"\\u1ea2nh tr\\u00ean\\\".\\nN\\u1ebfu kh\\u00f3 l\\u1ea5y \\u0111\\u01b0\\u1ee3c h\\u00ecnh \\u1ea3nh h\\u00e0ng may m\\u1eb7c, tr\\u01b0\\u1edbc ti\\u00ean h\\u00e3y s\\u1eed d\\u1ee5ng \\\"C\\u1eaft x\\u00e9n t\\u1ea3i l\\u00ean\\\". T\\u00ednh n\\u0103ng \\\"Tr\\u00edch xu\\u1ea5t qu\\u1ea7n \\u00e1o\\\" s\\u1ebd \\u0111\\u01b0\\u1ee3c b\\u1ed5 sung sau.\\nCh\\u1ea5t l\\u01b0\\u1ee3ng th\\u1ebf h\\u1ec7 li\\u00ean quan \\u0111\\u1ebfn m\\u00f4 h\\u00ecnh, ch\\u1ea5t l\\u01b0\\u1ee3ng h\\u00ecnh \\u1ea3nh v\\u00e0 n\\u1ec1n. N\\u1ec1n \\u0111\\u01a1n gi\\u1ea3n h\\u01a1n v\\u00e0 t\\u1ef7 l\\u1ec7 ng\\u01b0\\u1eddi l\\u1edbn h\\u01a1n th\\u01b0\\u1eddng t\\u1ea1o ra ch\\u1ea5t l\\u01b0\\u1ee3ng t\\u1ed1t h\\u01a1n.\\nM\\u1eabu Image Refiner c\\u00f3 th\\u1ec3 c\\u1ea3i thi\\u1ec7n c\\u00e1c chi ti\\u1ebft tr\\u00ean qu\\u1ea7n \\u00e1o v\\u00e0 da nh\\u01b0ng kh\\u00f4ng ph\\u00f9 h\\u1ee3p v\\u1edbi m\\u1ecdi tr\\u01b0\\u1eddng h\\u1ee3p v\\u00e0 \\u0111\\u00f4i khi c\\u00f3 th\\u1ec3 khi\\u1ebfn k\\u1ebft qu\\u1ea3 tr\\u1edf n\\u00ean t\\u1ed3i t\\u1ec7 h\\u01a1n. S\\u1eed d\\u1ee5ng n\\u00f3 m\\u1ed9t c\\u00e1ch c\\u1ea9n th\\u1eadn.\\nB\\u1ea1n c\\u00f3 th\\u1ec3 t\\u1ea3i l\\u00ean v\\u00e0 l\\u01b0u c\\u00e1c h\\u00ecnh \\u1ea3nh m\\u1eabu v\\u00e0 trang ph\\u1ee5c \\u0111\\u01b0\\u1ee3c s\\u1eed d\\u1ee5ng th\\u01b0\\u1eddng xuy\\u00ean trong \\\"M\\u1eabu c\\u1ee7a t\\u00f4i\\\" v\\u00e0 \\\"T\\u1ee7 qu\\u1ea7n \\u00e1o c\\u1ee7a t\\u00f4i\\\" \\u0111\\u1ec3 kh\\u00f4ng c\\u1ea7n ph\\u1ea3i t\\u1ea3i l\\u00ean c\\u00f9ng m\\u1ed9t h\\u00ecnh \\u1ea3nh m\\u1ed7i l\\u1ea7n.\\nN\\u1ebfu trang ph\\u1ee5c th\\u1eed qu\\u00e1 nh\\u1ecf v\\u00e0 kh\\u00f4ng che \\u0111\\u01b0\\u1ee3c nh\\u1eefng v\\u00f9ng kh\\u00f4ng l\\u1ed9 ra tr\\u00ean \\u1ea3nh ng\\u01b0\\u1eddi m\\u1eabu g\\u1ed1c, da \\u1edf nh\\u1eefng v\\u00f9ng \\u0111\\u00f3 c\\u00f3 th\\u1ec3 hi\\u1ec3n th\\u1ecb k\\u00e9m. Trong tr\\u01b0\\u1eddng h\\u1ee3p n\\u00e0y (th\\u01b0\\u1eddng g\\u1eb7p khi th\\u1eed bikini/\\u0111\\u1ed3 l\\u00f3t), h\\u00e3y chuy\\u1ec3n sang ng\\u01b0\\u1eddi m\\u1eabu \\u0111\\u1eb9p h\\u01a1n ho\\u1eb7c s\\u1eed d\\u1ee5ng \\u1ea3nh ng\\u01b0\\u1eddi m\\u1eabu \\u0111\\u1ec3 l\\u1ed9 v\\u00f9ng da t\\u01b0\\u01a1ng \\u1ee9ng.\\n\\u1ea2nh ng\\u01b0\\u1eddi m\\u1eabu kh\\u00f4ng \\u0111\\u01b0\\u1ee3c \\u0111\\u1ec3 l\\u1ed9 nh\\u1eefng b\\u1ed9 ph\\u1eadn nh\\u1ea1y c\\u1ea3m tr\\u00ean c\\u01a1 th\\u1ec3. H\\u00ecnh \\u1ea3nh qu\\u00e1 h\\u1edf hang c\\u00f3 th\\u1ec3 kh\\u00f4ng \\u0111\\u01b0\\u1ee3c ki\\u1ec3m duy\\u1ec7t v\\u00e0 t\\u1ea1o ra.\\nCh\\u00fang t\\u00f4i s\\u1ebd ti\\u1ebfp t\\u1ee5c c\\u1eadp nh\\u1eadt c\\u00e1c m\\u1eabu t\\u1ed1t h\\u01a1n cho c\\u00e1c t\\u01b0 th\\u1ebf kh\\u00e1c nhau v\\u00e0 cho c\\u00e1c t\\u00ecnh hu\\u1ed1ng m\\u1eb7c \\u0111\\u1ed3 l\\u00f3t/bikini. H\\u00e3y l\\u00e0m theo th\\u00f4ng b\\u00e1o n\\u1ec1n t\\u1ea3ng.\\nB\\u1ea1n c\\u00f3 th\\u1ec3 t\\u1ea3i l\\u00ean \\u1ea3nh m\\u00f4 h\\u00ecnh xem t\\u1eeb b\\u00ean c\\u1ea1nh (\\u1ea3nh ch\\u00ednh di\\u1ec7n th\\u01b0\\u1eddng ho\\u1ea1t \\u0111\\u1ed9ng t\\u1ed1t h\\u01a1n). N\\u1ebfu ch\\u01b0a c\\u00f3 \\u0111\\u1ea7y \\u0111\\u1ee7 h\\u00ecnh \\u1ea3nh qu\\u1ea7n \\u00e1o tr\\u00ean v\\u00e0 d\\u01b0\\u1edbi, b\\u1ea1n c\\u00f3 th\\u1ec3 t\\u1ea3i tr\\u1ef1c ti\\u1ebfp h\\u00ecnh \\u1ea3nh ng\\u01b0\\u1eddi m\\u1eb7c qu\\u1ea7n \\u00e1o l\\u00ean (ch\\u1ec9 t\\u1ea3i l\\u00ean \\\"\\u1ea2nh tr\\u00ean\\\", kh\\u00f4ng c\\u1ea7n t\\u1ea3i \\u1ea3nh b\\u00ean d\\u01b0\\u1edbi). Ng\\u01b0\\u1eddi m\\u1eabu s\\u1ebd t\\u1ef1 \\u0111\\u1ed9ng th\\u1eed qu\\u1ea7n \\u00e1o tr\\u00ean v\\u00e0 d\\u01b0\\u1edbi c\\u1ee7a ng\\u01b0\\u1eddi \\u0111\\u00f3 v\\u00e0o \\u1ea3nh ng\\u01b0\\u1eddi m\\u1eabu.\",\"ar\":\"\\u0628\\u0627\\u0644\\u0646\\u0633\\u0628\\u0629 \\u0644\\u0625\\u0635\\u062f\\u0627\\u0631 \\u0627\\u0644\\u0648\\u064a\\u0628\\u060c \\u064a\\u0648\\u0635\\u0649 \\u0628\\u0627\\u0633\\u062a\\u062e\\u062f\\u0627\\u0645 Chrome \\u0648Safari.\\n\\u064a\\u0645\\u0643\\u0646\\u0643 \\u062a\\u062d\\u0645\\u064a\\u0644 \\u0635\\u0648\\u0631 \\u0644\\u0644\\u0639\\u0627\\u0631\\u0636\\u0629 \\u0630\\u0627\\u062a \\u0645\\u0646\\u0638\\u0631 \\u062c\\u0627\\u0646\\u0628\\u064a \\u0623\\u0648 \\u062c\\u0627\\u0644\\u0633\\u0629\\u060c \\u0644\\u0643\\u0646 \\u0648\\u0636\\u0639\\u064a\\u0627\\u062a \\u0627\\u0644\\u0648\\u0642\\u0648\\u0641 \\u0627\\u0644\\u0645\\u0648\\u0627\\u062c\\u0647\\u0629 \\u0644\\u0644\\u0623\\u0645\\u0627\\u0645 \\u0639\\u0627\\u062f\\u0629 \\u0645\\u0627 \\u062a\\u0624\\u062f\\u064a \\u0625\\u0644\\u0649 \\u0646\\u062a\\u0627\\u0626\\u062c \\u0623\\u0641\\u0636\\u0644.\\n\\u0625\\u0630\\u0627 \\u0644\\u0645 \\u064a\\u0643\\u0646 \\u0644\\u062f\\u064a\\u0643 \\u0635\\u0648\\u0631 \\u0645\\u0646\\u0641\\u0635\\u0644\\u0629 \\u0644\\u0644\\u0645\\u0644\\u0627\\u0628\\u0633\\u060c \\u0641\\u064a\\u0645\\u0643\\u0646\\u0643 \\u062a\\u062d\\u0645\\u064a\\u0644 \\u0635\\u0648\\u0631\\u0629 \\u0627\\u0644\\u0634\\u062e\\u0635 \\u0627\\u0644\\u0630\\u064a \\u064a\\u0631\\u062a\\u062f\\u064a \\u0645\\u0644\\u0627\\u0628\\u0633\\u0647 \\u0645\\u0628\\u0627\\u0634\\u0631\\u0629 \\u0641\\u064a \\u062e\\u0627\\u0646\\u0629 \\\"\\u0627\\u0644\\u0635\\u0648\\u0631\\u0629 \\u0627\\u0644\\u0639\\u0644\\u0648\\u064a\\u0629\\\".\\n\\u0625\\u0630\\u0627 \\u0643\\u0627\\u0646 \\u0645\\u0646 \\u0627\\u0644\\u0635\\u0639\\u0628 \\u0627\\u0644\\u062d\\u0635\\u0648\\u0644 \\u0639\\u0644\\u0649 \\u0635\\u0648\\u0631 \\u0627\\u0644\\u0645\\u0644\\u0627\\u0628\\u0633\\u060c \\u0627\\u0633\\u062a\\u062e\\u062f\\u0645 \\\"\\u062a\\u062d\\u0645\\u064a\\u0644 \\u0627\\u0644\\u0642\\u0637\\u0639\\\" \\u0623\\u0648\\u0644\\u0627\\u064b. \\u0633\\u064a\\u062a\\u0645 \\u0625\\u0636\\u0627\\u0641\\u0629 \\u0645\\u064a\\u0632\\u0629 \\\"\\u0627\\u0633\\u062a\\u062e\\u0631\\u0627\\u062c \\u0627\\u0644\\u0645\\u0644\\u0627\\u0628\\u0633\\\" \\u0644\\u0627\\u062d\\u0642\\u064b\\u0627.\\n\\u062a\\u0631\\u062a\\u0628\\u0637 \\u062c\\u0648\\u062f\\u0629 \\u0627\\u0644\\u062c\\u064a\\u0644 \\u0628\\u0627\\u0644\\u0646\\u0645\\u0648\\u0630\\u062c \\u0648\\u062c\\u0648\\u062f\\u0629 \\u0627\\u0644\\u0635\\u0648\\u0631\\u0629 \\u0648\\u0627\\u0644\\u062e\\u0644\\u0641\\u064a\\u0629. \\u0639\\u0627\\u062f\\u0629\\u064b \\u0645\\u0627 \\u062a\\u0646\\u062a\\u062c \\u0627\\u0644\\u062e\\u0644\\u0641\\u064a\\u0627\\u062a \\u0627\\u0644\\u0623\\u0628\\u0633\\u0637 \\u0648\\u0646\\u0633\\u0628\\u0629 \\u0627\\u0644\\u0623\\u0634\\u062e\\u0627\\u0635 \\u0627\\u0644\\u0623\\u0643\\u0628\\u0631 \\u062c\\u0648\\u062f\\u0629 \\u0623\\u0641\\u0636\\u0644.\\n\\u064a\\u0645\\u0643\\u0646 \\u0644\\u0646\\u0645\\u0648\\u0630\\u062c Image Refiner \\u062a\\u062d\\u0633\\u064a\\u0646 \\u062a\\u0641\\u0627\\u0635\\u064a\\u0644 \\u0627\\u0644\\u0645\\u0644\\u0627\\u0628\\u0633 \\u0648\\u0627\\u0644\\u062c\\u0644\\u062f\\u060c \\u0648\\u0644\\u0643\\u0646\\u0647 \\u063a\\u064a\\u0631 \\u0645\\u0646\\u0627\\u0633\\u0628 \\u0644\\u0643\\u0644 \\u062d\\u0627\\u0644\\u0629 \\u0648\\u0642\\u062f \\u064a\\u0624\\u062f\\u064a \\u0641\\u064a \\u0628\\u0639\\u0636 \\u0627\\u0644\\u0623\\u062d\\u064a\\u0627\\u0646 \\u0625\\u0644\\u0649 \\u062a\\u0641\\u0627\\u0642\\u0645 \\u0627\\u0644\\u0646\\u062a\\u0627\\u0626\\u062c. \\u0627\\u0633\\u062a\\u062e\\u062f\\u0645\\u0647 \\u0628\\u0639\\u0646\\u0627\\u064a\\u0629.\\n\\u064a\\u0645\\u0643\\u0646\\u0643 \\u062a\\u062d\\u0645\\u064a\\u0644 \\u0648\\u062d\\u0641\\u0638 \\u0635\\u0648\\u0631 \\u0627\\u0644\\u0639\\u0627\\u0631\\u0636\\u0627\\u062a \\u0648\\u0627\\u0644\\u0645\\u0644\\u0627\\u0628\\u0633 \\u0627\\u0644\\u0645\\u0633\\u062a\\u062e\\u062f\\u0645\\u0629 \\u0628\\u0634\\u0643\\u0644 \\u0645\\u062a\\u0643\\u0631\\u0631 \\u0641\\u064a \\\"\\u0639\\u0627\\u0631\\u0636\\u0627\\u062a\\u064a\\\" \\u0648\\\"\\u062e\\u0632\\u0627\\u0646\\u062a\\u064a\\\" \\u0644\\u0630\\u0644\\u0643 \\u0644\\u0627 \\u062a\\u062d\\u062a\\u0627\\u062c \\u0625\\u0644\\u0649 \\u062a\\u062d\\u0645\\u064a\\u0644 \\u0646\\u0641\\u0633 \\u0627\\u0644\\u0635\\u0648\\u0631 \\u0641\\u064a \\u0643\\u0644 \\u0645\\u0631\\u0629.\\n\\u0625\\u0630\\u0627 \\u0643\\u0627\\u0646\\u062a \\u0642\\u0637\\u0639\\u0629 \\u0627\\u0644\\u0645\\u0644\\u0627\\u0628\\u0633 \\u0627\\u0644\\u062a\\u064a \\u062a\\u0645\\u062a \\u062a\\u062c\\u0631\\u0628\\u062a\\u0647\\u0627 \\u0635\\u063a\\u064a\\u0631\\u0629 \\u062c\\u062f\\u064b\\u0627 \\u0648\\u0644\\u0627 \\u062a\\u063a\\u0637\\u064a \\u0627\\u0644\\u0645\\u0646\\u0627\\u0637\\u0642 \\u0627\\u0644\\u062a\\u064a \\u0644\\u0645 \\u064a\\u062a\\u0645 \\u0643\\u0634\\u0641\\u0647\\u0627 \\u0641\\u064a \\u0635\\u0648\\u0631\\u0629 \\u0627\\u0644\\u0646\\u0645\\u0648\\u0630\\u062c \\u0627\\u0644\\u0623\\u0635\\u0644\\u064a\\u0629\\u060c \\u0641\\u0642\\u062f \\u064a\\u0635\\u0628\\u062d \\u0627\\u0644\\u062c\\u0644\\u062f \\u0641\\u064a \\u062a\\u0644\\u0643 \\u0627\\u0644\\u0645\\u0646\\u0627\\u0637\\u0642 \\u0633\\u064a\\u0626\\u064b\\u0627. \\u0641\\u064a \\u0647\\u0630\\u0647 \\u0627\\u0644\\u062d\\u0627\\u0644\\u0629 (\\u0627\\u0644\\u0634\\u0627\\u0626\\u0639\\u0629 \\u0645\\u0639 \\u062a\\u062c\\u0627\\u0631\\u0628 \\u0627\\u0644\\u0628\\u064a\\u0643\\u064a\\u0646\\u064a/\\u0627\\u0644\\u0645\\u0644\\u0627\\u0628\\u0633 \\u0627\\u0644\\u062f\\u0627\\u062e\\u0644\\u064a\\u0629)\\u060c \\u0642\\u0645 \\u0628\\u0627\\u0644\\u062a\\u0628\\u062f\\u064a\\u0644 \\u0625\\u0644\\u0649 \\u0646\\u0645\\u0648\\u0630\\u062c \\u0623\\u0641\\u0636\\u0644 \\u0623\\u0648 \\u0627\\u0633\\u062a\\u062e\\u062f\\u0645 \\u0635\\u0648\\u0631\\u0629 \\u0646\\u0645\\u0648\\u0630\\u062c\\u064a\\u0629 \\u062a\\u0643\\u0634\\u0641 \\u0645\\u0646\\u0637\\u0642\\u0629 \\u0627\\u0644\\u062c\\u0644\\u062f \\u0627\\u0644\\u0645\\u0642\\u0627\\u0628\\u0644\\u0629.\\n\\u064a\\u062c\\u0628 \\u0623\\u0644\\u0627 \\u062a\\u0643\\u0634\\u0641 \\u0635\\u0648\\u0631 \\u0627\\u0644\\u0639\\u0627\\u0631\\u0636\\u0629 \\u0623\\u062c\\u0632\\u0627\\u0621 \\u0627\\u0644\\u062c\\u0633\\u0645 \\u0627\\u0644\\u062d\\u0633\\u0627\\u0633\\u0629. \\u0627\\u0644\\u0635\\u0648\\u0631 \\u0627\\u0644\\u0645\\u0641\\u0631\\u0637\\u0629 \\u0641\\u064a \\u0627\\u0644\\u0643\\u0634\\u0641 \\u0642\\u062f \\u062a\\u0641\\u0634\\u0644 \\u0641\\u064a \\u0627\\u0644\\u0627\\u0639\\u062a\\u062f\\u0627\\u0644 \\u0648\\u0627\\u0644\\u062a\\u0648\\u0644\\u064a\\u062f.\\n\\u0633\\u0646\\u0648\\u0627\\u0635\\u0644 \\u062a\\u062d\\u062f\\u064a\\u062b \\u0646\\u0645\\u0627\\u0630\\u062c \\u0623\\u0641\\u0636\\u0644 \\u0644\\u0645\\u062e\\u062a\\u0644\\u0641 \\u0627\\u0644\\u0623\\u0648\\u0636\\u0627\\u0639 \\u0648\\u0644\\u0633\\u064a\\u0646\\u0627\\u0631\\u064a\\u0648\\u0647\\u0627\\u062a \\u0627\\u0644\\u0645\\u0644\\u0627\\u0628\\u0633 \\u0627\\u0644\\u062f\\u0627\\u062e\\u0644\\u064a\\u0629/\\u0627\\u0644\\u0628\\u064a\\u0643\\u064a\\u0646\\u064a. \\u064a\\u0631\\u062c\\u0649 \\u0645\\u062a\\u0627\\u0628\\u0639\\u0629 \\u0625\\u0639\\u0644\\u0627\\u0646\\u0627\\u062a \\u0627\\u0644\\u0645\\u0646\\u0635\\u0629.\\n\\u064a\\u0645\\u0643\\u0646\\u0643 \\u062a\\u062d\\u0645\\u064a\\u0644 \\u0635\\u0648\\u0631 \\u0646\\u0645\\u0627\\u0630\\u062c \\u0627\\u0644\\u0639\\u0631\\u0636 \\u0627\\u0644\\u062c\\u0627\\u0646\\u0628\\u064a (\\u0639\\u0627\\u062f\\u0629\\u064b \\u0645\\u0627 \\u062a\\u0639\\u0645\\u0644 \\u0627\\u0644\\u0635\\u0648\\u0631 \\u0627\\u0644\\u0623\\u0645\\u0627\\u0645\\u064a\\u0629 \\u0628\\u0634\\u0643\\u0644 \\u0623\\u0641\\u0636\\u0644). \\u0625\\u0630\\u0627 \\u0644\\u0645 \\u064a\\u0643\\u0646 \\u0644\\u062f\\u064a\\u0643 \\u0635\\u0648\\u0631 \\u0643\\u0627\\u0645\\u0644\\u0629 \\u0644\\u0644\\u0645\\u0644\\u0627\\u0628\\u0633 \\u0627\\u0644\\u0639\\u0644\\u0648\\u064a\\u0629 \\u0648\\u0627\\u0644\\u0633\\u0641\\u0644\\u064a\\u0629\\u060c \\u0641\\u064a\\u0645\\u0643\\u0646\\u0643 \\u062a\\u062d\\u0645\\u064a\\u0644 \\u0635\\u0648\\u0631\\u0629 \\u0627\\u0644\\u0634\\u062e\\u0635 \\u0627\\u0644\\u0630\\u064a \\u064a\\u0631\\u062a\\u062f\\u064a \\u0645\\u0644\\u0627\\u0628\\u0633\\u0647 \\u0645\\u0628\\u0627\\u0634\\u0631\\u0629\\u064b (\\u0642\\u0645 \\u0628\\u062a\\u062d\\u0645\\u064a\\u0644\\u0647\\u0627 \\u0625\\u0644\\u0649 \\\"\\u0627\\u0644\\u0635\\u0648\\u0631\\u0629 \\u0627\\u0644\\u0639\\u0644\\u0648\\u064a\\u0629\\\" \\u0641\\u0642\\u0637\\u060c \\u0648\\u0644\\u0627 \\u062d\\u0627\\u062c\\u0629 \\u0644\\u062a\\u062d\\u0645\\u064a\\u0644 \\u0635\\u0648\\u0631\\u0629 \\u0633\\u0641\\u0644\\u064a\\u0629). \\u0633\\u064a\\u062d\\u0627\\u0648\\u0644 \\u0627\\u0644\\u0639\\u0627\\u0631\\u0636 \\u062a\\u0644\\u0642\\u0627\\u0626\\u064a\\u064b\\u0627 \\u0627\\u0631\\u062a\\u062f\\u0627\\u0621 \\u0627\\u0644\\u0645\\u0644\\u0627\\u0628\\u0633 \\u0627\\u0644\\u0639\\u0644\\u0648\\u064a\\u0629 \\u0648\\u0627\\u0644\\u0633\\u0641\\u0644\\u064a\\u0629 \\u0644\\u0630\\u0644\\u0643 \\u0627\\u0644\\u0634\\u062e\\u0635 \\u0639\\u0644\\u0649 \\u0635\\u0648\\u0631\\u0629 \\u0627\\u0644\\u0639\\u0627\\u0631\\u0636\\u0629.\",\"ja\":\"Web\\u7248\\u306e\\u5834\\u5408\\u306fChrome\\u3068Safari\\u3092\\u63a8\\u5968\\u3057\\u307e\\u3059\\u3002\\n\\u6a2a\\u304b\\u3089\\u898b\\u305f\\u5199\\u771f\\u3084\\u5ea7\\u3063\\u3066\\u3044\\u308b\\u30e2\\u30c7\\u30eb\\u306e\\u5199\\u771f\\u3092\\u30a2\\u30c3\\u30d7\\u30ed\\u30fc\\u30c9\\u3067\\u304d\\u307e\\u3059\\u304c\\u3001\\u901a\\u5e38\\u306f\\u6b63\\u9762\\u3092\\u5411\\u3044\\u305f\\u7acb\\u3061\\u30dd\\u30fc\\u30ba\\u306e\\u65b9\\u304c\\u826f\\u3044\\u7d50\\u679c\\u304c\\u5f97\\u3089\\u308c\\u307e\\u3059\\u3002\\n\\u5225\\u306e\\u8863\\u670d\\u306e\\u753b\\u50cf\\u304c\\u306a\\u3044\\u5834\\u5408\\u306f\\u3001\\u300c\\u4e0a\\u306e\\u753b\\u50cf\\u300d\\u30b9\\u30ed\\u30c3\\u30c8\\u306b\\u670d\\u3092\\u7740\\u305f\\u4eba\\u7269\\u306e\\u753b\\u50cf\\u3092\\u76f4\\u63a5\\u30a2\\u30c3\\u30d7\\u30ed\\u30fc\\u30c9\\u3067\\u304d\\u307e\\u3059\\u3002\\n\\u8863\\u670d\\u306e\\u753b\\u50cf\\u3092\\u5165\\u624b\\u3059\\u308b\\u306e\\u304c\\u96e3\\u3057\\u3044\\u5834\\u5408\\u306f\\u3001\\u307e\\u305a\\u300c\\u30af\\u30ed\\u30c3\\u30d7\\u30a2\\u30c3\\u30d7\\u30ed\\u30fc\\u30c9\\u300d\\u3092\\u4f7f\\u7528\\u3057\\u3066\\u304f\\u3060\\u3055\\u3044\\u3002 \\u300c\\u8863\\u670d\\u62bd\\u51fa\\u300d\\u6a5f\\u80fd\\u306f\\u4eca\\u5f8c\\u8ffd\\u52a0\\u3055\\u308c\\u308b\\u4e88\\u5b9a\\u3067\\u3059\\u3002\\n\\u751f\\u6210\\u54c1\\u8cea\\u306f\\u3001\\u30e2\\u30c7\\u30eb\\u3001\\u753b\\u50cf\\u54c1\\u8cea\\u3001\\u304a\\u3088\\u3073\\u80cc\\u666f\\u306b\\u95a2\\u9023\\u3057\\u307e\\u3059\\u3002\\u901a\\u5e38\\u3001\\u80cc\\u666f\\u304c\\u30b7\\u30f3\\u30d7\\u30eb\\u3067\\u4eba\\u7269\\u306e\\u6bd4\\u7387\\u304c\\u5927\\u304d\\u3044\\u307b\\u3069\\u3001\\u54c1\\u8cea\\u304c\\u9ad8\\u304f\\u306a\\u308a\\u307e\\u3059\\u3002\\nImage Refiner \\u30e2\\u30c7\\u30eb\\u306f\\u8863\\u670d\\u3084\\u808c\\u306e\\u8a73\\u7d30\\u3092\\u6539\\u5584\\u3067\\u304d\\u307e\\u3059\\u304c\\u3001\\u3059\\u3079\\u3066\\u306e\\u30b1\\u30fc\\u30b9\\u306b\\u9069\\u3057\\u3066\\u3044\\u308b\\u308f\\u3051\\u3067\\u306f\\u306a\\u304f\\u3001\\u5834\\u5408\\u306b\\u3088\\u3063\\u3066\\u306f\\u7d50\\u679c\\u304c\\u60aa\\u5316\\u3059\\u308b\\u53ef\\u80fd\\u6027\\u304c\\u3042\\u308a\\u307e\\u3059\\u3002\\u614e\\u91cd\\u306b\\u4f7f\\u7528\\u3057\\u3066\\u304f\\u3060\\u3055\\u3044\\u3002\\n\\u983b\\u7e41\\u306b\\u4f7f\\u7528\\u3059\\u308b\\u30e2\\u30c7\\u30eb\\u3084\\u8863\\u670d\\u306e\\u753b\\u50cf\\u3092\\u300c\\u30de\\u30a4 \\u30e2\\u30c7\\u30eb\\u300d\\u3068\\u300c\\u30de\\u30a4 \\u30af\\u30ed\\u30fc\\u30bc\\u30c3\\u30c8\\u300d\\u306b\\u30a2\\u30c3\\u30d7\\u30ed\\u30fc\\u30c9\\u3057\\u3066\\u4fdd\\u5b58\\u3067\\u304d\\u308b\\u305f\\u3081\\u3001\\u6bce\\u56de\\u540c\\u3058\\u753b\\u50cf\\u3092\\u30a2\\u30c3\\u30d7\\u30ed\\u30fc\\u30c9\\u3059\\u308b\\u5fc5\\u8981\\u306f\\u3042\\u308a\\u307e\\u305b\\u3093\\u3002\\n\\u8a66\\u7740\\u3057\\u305f\\u8863\\u670d\\u304c\\u5c0f\\u3055\\u3059\\u304e\\u3066\\u3001\\u5143\\u306e\\u30e2\\u30c7\\u30eb\\u306e\\u5199\\u771f\\u3067\\u306f\\u9732\\u51fa\\u3057\\u3066\\u3044\\u306a\\u3044\\u9818\\u57df\\u304c\\u8986\\u308f\\u308c\\u3066\\u3044\\u306a\\u3044\\u5834\\u5408\\u3001\\u305d\\u306e\\u9818\\u57df\\u306e\\u76ae\\u819a\\u306e\\u30ec\\u30f3\\u30c0\\u30ea\\u30f3\\u30b0\\u304c\\u4e0d\\u5341\\u5206\\u306b\\u306a\\u308b\\u53ef\\u80fd\\u6027\\u304c\\u3042\\u308a\\u307e\\u3059\\u3002\\u3053\\u306e\\u5834\\u5408\\uff08\\u30d3\\u30ad\\u30cb\\u3084\\u4e0b\\u7740\\u306e\\u8a66\\u7740\\u306b\\u3088\\u304f\\u3042\\u308a\\u307e\\u3059\\uff09\\u3001\\u3088\\u308a\\u512a\\u308c\\u305f\\u30e2\\u30c7\\u30eb\\u306b\\u5207\\u308a\\u66ff\\u3048\\u308b\\u304b\\u3001\\u5bfe\\u5fdc\\u3059\\u308b\\u808c\\u9818\\u57df\\u304c\\u9732\\u51fa\\u3057\\u3066\\u3044\\u308b\\u30e2\\u30c7\\u30eb\\u5199\\u771f\\u3092\\u4f7f\\u7528\\u3057\\u3066\\u304f\\u3060\\u3055\\u3044\\u3002\\n\\u30e2\\u30c7\\u30eb\\u306e\\u5199\\u771f\\u3067\\u306f\\u3001\\u654f\\u611f\\u306a\\u4f53\\u306e\\u90e8\\u5206\\u3092\\u9732\\u51fa\\u3057\\u3066\\u306f\\u306a\\u308a\\u307e\\u305b\\u3093\\u3002\\u904e\\u5ea6\\u306b\\u9732\\u51fa\\u3057\\u305f\\u753b\\u50cf\\u306f\\u30e2\\u30c7\\u30ec\\u30fc\\u30b7\\u30e7\\u30f3\\u3068\\u751f\\u6210\\u306b\\u5931\\u6557\\u3059\\u308b\\u53ef\\u80fd\\u6027\\u304c\\u3042\\u308a\\u307e\\u3059\\u3002\\n\\u3055\\u307e\\u3056\\u307e\\u306a\\u30dd\\u30fc\\u30ba\\u3084\\u4e0b\\u7740/\\u30d3\\u30ad\\u30cb\\u306e\\u30b7\\u30ca\\u30ea\\u30aa\\u5411\\u3051\\u306b\\u3001\\u3088\\u308a\\u512a\\u308c\\u305f\\u30e2\\u30c7\\u30eb\\u3092\\u66f4\\u65b0\\u3057\\u7d9a\\u3051\\u307e\\u3059\\u3002\\u30d7\\u30e9\\u30c3\\u30c8\\u30d5\\u30a9\\u30fc\\u30e0\\u306e\\u30a2\\u30ca\\u30a6\\u30f3\\u30b9\\u306b\\u5f93\\u3063\\u3066\\u304f\\u3060\\u3055\\u3044\\u3002\\n\\u6a2a\\u304b\\u3089\\u898b\\u305f\\u30e2\\u30c7\\u30eb\\u306e\\u5199\\u771f\\u3092\\u30a2\\u30c3\\u30d7\\u30ed\\u30fc\\u30c9\\u3067\\u304d\\u307e\\u3059 (\\u901a\\u5e38\\u306f\\u6b63\\u9762\\u304b\\u3089\\u898b\\u305f\\u5199\\u771f\\u306e\\u65b9\\u304c\\u9069\\u3057\\u3066\\u3044\\u307e\\u3059)\\u3002\\u4e0a\\u4e0b\\u306e\\u5b8c\\u5168\\u306a\\u8863\\u88c5\\u753b\\u50cf\\u304c\\u306a\\u3044\\u5834\\u5408\\u306f\\u3001\\u8863\\u88c5\\u3092\\u7740\\u305f\\u4eba\\u7269\\u306e\\u753b\\u50cf\\u3092\\u76f4\\u63a5\\u30a2\\u30c3\\u30d7\\u30ed\\u30fc\\u30c9\\u3067\\u304d\\u307e\\u3059\\uff08\\u300c\\u4e0a\\u306e\\u753b\\u50cf\\u300d\\u306e\\u307f\\u306b\\u30a2\\u30c3\\u30d7\\u30ed\\u30fc\\u30c9\\u3057\\u3001\\u4e0b\\u306e\\u753b\\u50cf\\u3092\\u30a2\\u30c3\\u30d7\\u30ed\\u30fc\\u30c9\\u3059\\u308b\\u5fc5\\u8981\\u306f\\u3042\\u308a\\u307e\\u305b\\u3093\\uff09\\u3002\\u30e2\\u30c7\\u30eb\\u5199\\u771f\\u306e\\u4e0a\\u3067\\u305d\\u306e\\u4eba\\u306e\\u4e0a\\u4e0b\\u3092\\u30e2\\u30c7\\u30eb\\u304c\\u81ea\\u52d5\\u3067\\u8a66\\u7740\\u3057\\u307e\\u3059\\u3002\",\"ko\":\"\\uc6f9 \\ubc84\\uc804\\uc758 \\uacbd\\uc6b0 Chrome, Safari\\ub97c \\uad8c\\uc7a5\\ud569\\ub2c8\\ub2e4.\\n\\uc606\\ubaa8\\uc2b5\\uc774\\ub098 \\uc549\\uc544 \\uc788\\ub294 \\ubaa8\\ub378 \\uc0ac\\uc9c4\\uc744 \\uc5c5\\ub85c\\ub4dc\\ud560 \\uc218 \\uc788\\uc9c0\\ub9cc \\uc77c\\ubc18\\uc801\\uc73c\\ub85c \\uc815\\uba74\\uc744 \\ubc14\\ub77c\\ubcf4\\uace0 \\uc11c \\uc788\\ub294 \\uc790\\uc138\\uac00 \\ub354 \\ub098\\uc740 \\uacb0\\uacfc\\ub97c \\uc0dd\\uc131\\ud569\\ub2c8\\ub2e4.\\n\\ubcc4\\ub3c4\\uc758 \\uc758\\uc0c1 \\uc774\\ubbf8\\uc9c0\\uac00 \\uc5c6\\uc744 \\uacbd\\uc6b0 '\\uc0c1\\ub2e8 \\uc774\\ubbf8\\uc9c0' \\uc2ac\\ub86f\\uc5d0 \\uc637\\uc744 \\uc785\\uc740 \\uc0ac\\ub78c \\uc774\\ubbf8\\uc9c0\\ub97c \\uc9c1\\uc811 \\uc5c5\\ub85c\\ub4dc\\ud560 \\uc218 \\uc788\\uc2b5\\ub2c8\\ub2e4.\\n\\uc758\\ub958 \\uc774\\ubbf8\\uc9c0\\ub97c \\uc5bb\\uae30 \\uc5b4\\ub824\\uc6b4 \\uacbd\\uc6b0 \\uba3c\\uc800 \\\"\\uc790\\ub974\\uae30 \\uc5c5\\ub85c\\ub4dc\\\"\\ub97c \\uc0ac\\uc6a9\\ud558\\uc138\\uc694. '\\uc758\\ub958 \\ucd94\\ucd9c' \\uae30\\ub2a5\\uc740 \\ucd94\\ud6c4 \\ucd94\\uac00\\ub420 \\uc608\\uc815\\uc785\\ub2c8\\ub2e4.\\n\\uc0dd\\uc131 \\ud488\\uc9c8\\uc740 \\ubaa8\\ub378, \\uc774\\ubbf8\\uc9c0 \\ud488\\uc9c8, \\ubc30\\uacbd\\uacfc \\uad00\\ub828\\uc774 \\uc788\\uc2b5\\ub2c8\\ub2e4. \\ubc30\\uacbd\\uc774 \\ub2e8\\uc21c\\ud558\\uace0 \\uc0ac\\ub78c \\ube44\\uc728\\uc774 \\ud074\\uc218\\ub85d \\uc77c\\ubc18\\uc801\\uc73c\\ub85c \\ud488\\uc9c8\\uc774 \\ub354 \\uc88b\\uc544\\uc9d1\\ub2c8\\ub2e4.\\nImage Refiner \\ubaa8\\ub378\\uc740 \\uc758\\ubcf5 \\ubc0f \\ud53c\\ubd80 \\uc138\\ubd80 \\uc0ac\\ud56d\\uc744 \\uac1c\\uc120\\ud560 \\uc218 \\uc788\\uc9c0\\ub9cc \\ubaa8\\ub4e0 \\uacbd\\uc6b0\\uc5d0 \\uc801\\ud569\\ud558\\uc9c0\\ub294 \\uc54a\\uc73c\\uba70 \\ub54c\\ub85c\\ub294 \\uacb0\\uacfc\\uac00 \\ub354 \\ub098\\ube60\\uc9c8 \\uc218 \\uc788\\uc2b5\\ub2c8\\ub2e4. \\uc8fc\\uc758\\ud574\\uc11c \\uc0ac\\uc6a9\\ud558\\uc138\\uc694.\\n\\uc790\\uc8fc \\uc0ac\\uc6a9\\ud558\\ub294 \\ubaa8\\ub378 \\ubc0f \\uc758\\uc0c1 \\uc774\\ubbf8\\uc9c0\\ub97c \\\"\\ub0b4 \\ubaa8\\ub378\\\", \\\"\\ub0b4 \\uc637\\uc7a5\\\"\\uc5d0 \\uc5c5\\ub85c\\ub4dc\\ud558\\uc5ec \\uc800\\uc7a5\\ud560 \\uc218 \\uc788\\uc5b4 \\ub9e4\\ubc88 \\uac19\\uc740 \\uc774\\ubbf8\\uc9c0\\ub97c \\uc5c5\\ub85c\\ub4dc\\ud560 \\ud544\\uc694\\uac00 \\uc5c6\\uc2b5\\ub2c8\\ub2e4.\\n\\uc2dc\\ucc29 \\uc758\\uc0c1\\uc774 \\ub108\\ubb34 \\uc791\\uc544\\uc11c \\uc6d0\\ubcf8 \\ubaa8\\ub378 \\uc0ac\\uc9c4\\uc5d0\\uc11c \\ub178\\ucd9c\\ub418\\uc9c0 \\uc54a\\uc740 \\ubd80\\ubd84\\uc744 \\ub36e\\uc9c0 \\ubabb\\ud558\\ub294 \\uacbd\\uc6b0 \\ud574\\ub2f9 \\ubd80\\ubd84\\uc758 \\ud53c\\ubd80\\uac00 \\uc81c\\ub300\\ub85c \\ud45c\\ud604\\ub418\\uc9c0 \\uc54a\\uc744 \\uc218 \\uc788\\uc2b5\\ub2c8\\ub2e4. \\uc774\\ub7f0 \\uacbd\\uc6b0(\\ube44\\ud0a4\\ub2c8/\\uc18d\\uc637 \\ucc29\\uc6a9 \\uc2dc \\ud754\\ud788 \\ubc1c\\uc0dd\\ud558\\ub294 \\ud604\\uc0c1), \\ub354 \\ub098\\uc740 \\ubaa8\\ub378\\ub85c \\uad50\\uccb4\\ud558\\uac70\\ub098 \\ud574\\ub2f9 \\ud53c\\ubd80 \\ubd80\\uc704\\ub97c \\ub178\\ucd9c\\ud558\\ub294 \\ubaa8\\ub378 \\uc0ac\\uc9c4\\uc744 \\uc0ac\\uc6a9\\ud558\\uc138\\uc694.\\n\\ubaa8\\ub378 \\uc0ac\\uc9c4\\uc740 \\ubbfc\\uac10\\ud55c \\uc2e0\\uccb4 \\ubd80\\uc704\\ub97c \\ub178\\ucd9c\\ud574\\uc11c\\ub294 \\uc548 \\ub429\\ub2c8\\ub2e4. \\uc9c0\\ub098\\uce58\\uac8c \\ub178\\ucd9c\\ub41c \\uc774\\ubbf8\\uc9c0\\ub294 \\uc870\\uc815 \\ubc0f \\uc0dd\\uc131\\uc5d0 \\uc2e4\\ud328\\ud560 \\uc218 \\uc788\\uc2b5\\ub2c8\\ub2e4.\\n\\ub2e4\\uc591\\ud55c \\ud3ec\\uc988\\uc640 \\uc18d\\uc637/\\ube44\\ud0a4\\ub2c8 \\uc2dc\\ub098\\ub9ac\\uc624\\uc5d0 \\ub9de\\uac8c \\ub354 \\ub098\\uc740 \\ubaa8\\ub378\\uc744 \\uacc4\\uc18d \\uc5c5\\ub370\\uc774\\ud2b8\\ud560 \\uc608\\uc815\\uc785\\ub2c8\\ub2e4. \\ud50c\\ub7ab\\ud3fc \\uacf5\\uc9c0\\ub97c \\ub530\\ub974\\uc138\\uc694.\\n\\uce21\\uba74 \\ubaa8\\ub378 \\uc0ac\\uc9c4\\uc744 \\uc5c5\\ub85c\\ub4dc\\ud560 \\uc218 \\uc788\\uc2b5\\ub2c8\\ub2e4. \\uc77c\\ubc18\\uc801\\uc73c\\ub85c \\uc815\\uba74 \\uc0ac\\uc9c4\\uc774 \\ub354 \\uc88b\\uc2b5\\ub2c8\\ub2e4. \\uc644\\uc804\\ud55c \\uc0c1\\uc758, \\ud558\\uc758 \\uc774\\ubbf8\\uc9c0\\uac00 \\uc5c6\\ub294 \\uacbd\\uc6b0, \\uc637\\uc744 \\uc785\\uc740 \\uc0ac\\ub78c \\uc774\\ubbf8\\uc9c0\\ub97c \\uc9c1\\uc811 \\uc5c5\\ub85c\\ub4dc\\ud560 \\uc218 \\uc788\\uc2b5\\ub2c8\\ub2e4('\\uc0c1\\ub2e8 \\uc774\\ubbf8\\uc9c0'\\uc5d0\\ub9cc \\uc5c5\\ub85c\\ub4dc, \\ud558\\uc758 \\uc774\\ubbf8\\uc9c0\\ub294 \\uc5c5\\ub85c\\ub4dc\\ud560 \\ud544\\uc694 \\uc5c6\\uc74c). \\ubaa8\\ub378\\uc740 \\ud574\\ub2f9 \\uc778\\ubb3c\\uc758 \\uc0c1\\uc758\\uc640 \\ud558\\uc758\\ub97c \\ubaa8\\ub378 \\uc0ac\\uc9c4\\uc5d0 \\uc790\\ub3d9\\uc73c\\ub85c \\uc785\\uc5b4\\ubcfc \\uac83\\uc785\\ub2c8\\ub2e4.\",\"ms\":\"Untuk versi web, Chrome dan Safari disyorkan.\\nAnda boleh memuat naik foto model pandangan sisi atau duduk, tetapi pose berdiri menghadap ke hadapan biasanya menjana hasil yang lebih baik.\\nJika anda tidak mempunyai imej pakaian yang berasingan, anda boleh terus memuat naik imej orang berpakaian dalam slot \\\"Imej Atas\\\".\\nJika imej pakaian sukar diperoleh, gunakan \\\"Pangkas Muat Naik\\\" dahulu. Ciri \\\"Ekstraksi Pakaian\\\" akan ditambahkan kemudian.\\nKualiti penjanaan berkaitan dengan model, kualiti imej dan latar belakang. Latar belakang yang lebih ringkas dan nisbah orang yang lebih besar biasanya menghasilkan kualiti yang lebih baik.\\nModel Penapis Imej boleh menambah baik butiran pakaian dan kulit, tetapi ia tidak sesuai untuk setiap kes dan kadangkala boleh memburukkan keputusan. Gunakan dengan berhati-hati.\\nAnda boleh memuat naik dan menyimpan imej model dan pakaian yang kerap digunakan dalam \\\"Model Saya\\\" dan \\\"Almari Saya\\\" supaya anda tidak perlu memuat naik imej yang sama setiap kali.\\nJika pakaian try-on terlalu kecil dan tidak menutupi kawasan yang tidak terdedah pada foto model asal, kulit di kawasan tersebut mungkin menjadi teruk. Dalam kes ini (biasa dengan percubaan bikini/seluar dalam), tukar kepada model yang lebih baik atau gunakan foto model yang mendedahkan kawasan kulit yang sepadan.\\nFoto model tidak boleh mendedahkan bahagian badan yang sensitif. Imej yang terlalu mendedahkan mungkin gagal penyederhanaan dan penjanaan.\\nKami akan terus mengemas kini model yang lebih baik untuk pose yang berbeza dan untuk senario seluar dalam/bikini. Sila ikuti pengumuman platform.\\nAnda boleh memuat naik foto model pandangan sisi (foto menghadap hadapan biasanya berfungsi lebih baik). Jika anda tidak mempunyai imej pakaian atas dan bawah yang lengkap, anda boleh memuat naik terus imej orang berpakaian (muat naik ke \\\"Imej Atas\\\" sahaja, tidak perlu memuat naik imej yang lebih rendah). Model secara automatik akan mencuba pakaian atas dan bawah orang itu pada foto model.\"}"
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
