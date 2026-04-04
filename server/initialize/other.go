package initialize

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"os"
	"strings"

	"github.com/songzhibin97/gkit/cache/local_cache"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
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
		// points 分组
		{ConfigKey: "points_exchange_rate", ConfigValue: "100", ConfigName: "积分兑换比率", ConfigGroup: "points", Remark: "多少积分兑换1货币单位，如100积分=1元"},
		// order 分组
		{ConfigKey: "order_close_minutes", ConfigValue: "20", ConfigName: "订单自动关闭时间", ConfigGroup: "order", Remark: "未支付订单自动关闭的分钟数"},
		{ConfigKey: "order_auto_close_minutes", ConfigValue: "20", ConfigName: "订单自动关闭时间(旧key)", ConfigGroup: "order", Remark: "兼容旧key，同order_close_minutes"},
		// display 分组
		{ConfigKey: "presale_home_count", ConfigValue: "4", ConfigName: "首页预售展示数", ConfigGroup: "display", Remark: "uni首页展示的预售商品数量"},
		{ConfigKey: "sign_in_enabled", ConfigValue: "true", ConfigName: "签到功能开关", ConfigGroup: "display", Remark: "uni端是否显示签到悬浮按钮(true/false)"},
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
			global.GVA_DB.Create(&cfg)
		}
	}

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
