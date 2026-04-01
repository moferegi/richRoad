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
		{ConfigKey: "invite_reward_points", ConfigValue: "10", ConfigName: "邀请奖励积分", ConfigGroup: "invite", Remark: "用户邀请好友注册后获得的积分奖励"},
		{ConfigKey: "payment_qr_code", ConfigValue: "", ConfigName: "收款二维码", ConfigGroup: "payment", Remark: "客户下单后展示的收款二维码图片地址"},
		{ConfigKey: "payment_contact", ConfigValue: "", ConfigName: "支付联系客服", ConfigGroup: "payment", Remark: "支付页面联系客服的链接或信息"},
		{ConfigKey: "maintenance_enabled", ConfigValue: "false", ConfigName: "维护模式开关", ConfigGroup: "system", Remark: "是否开启全站维护模式(true/false)"},
		{ConfigKey: "maintenance_message", ConfigValue: "系统维护中，请稍后再试", ConfigName: "维护提示语", ConfigGroup: "system", Remark: "维护模式下的提示信息"},
		{ConfigKey: "maintenance_bg_image", ConfigValue: "", ConfigName: "维护背景图", ConfigGroup: "system", Remark: "维护页面的背景图片地址"},
		{ConfigKey: "points_exchange_rate", ConfigValue: "100", ConfigName: "积分兑换比率", ConfigGroup: "points", Remark: "多少积分兑换1元(分)"},
		{ConfigKey: "order_close_minutes", ConfigValue: "15", ConfigName: "订单自动关闭时间", ConfigGroup: "order", Remark: "未支付订单自动关闭的分钟数"},
		{ConfigKey: "currency_symbol", ConfigValue: "¥", ConfigName: "货币符号", ConfigGroup: "system", Remark: "前端展示的货币符号"},
		{ConfigKey: "currency_unit", ConfigValue: "CNY", ConfigName: "货币单位", ConfigGroup: "system", Remark: "货币单位编码"},
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
