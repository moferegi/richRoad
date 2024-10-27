package initialize

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/email"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/geo"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/plugin"
	"github.com/gin-gonic/gin"
)

func PluginInit(group *gin.RouterGroup, Plugin ...plugin.Plugin) {
	for i := range Plugin {
		fmt.Println(Plugin[i].RouterPath(), "注册开始!")
		PluginGroup := group.Group(Plugin[i].RouterPath())
		Plugin[i].Register(PluginGroup)
		fmt.Println(Plugin[i].RouterPath(), "注册成功!")
	}
}

func bizPluginV1(group ...*gin.RouterGroup) {
	private := group[0]
	public := group[1]

	PluginInit(private, email.CreateEmailPlug(
		global.GVA_CONFIG.Email.To,
		global.GVA_CONFIG.Email.From,
		global.GVA_CONFIG.Email.Host,
		global.GVA_CONFIG.Email.Secret,
		global.GVA_CONFIG.Email.Nickname,
		global.GVA_CONFIG.Email.Port,
		global.GVA_CONFIG.Email.IsSSL,
	))

	PluginInit(public, wxpay.CreateWxpayPlug(
		global.GVA_CONFIG.Wxpay.MchID,
		global.GVA_CONFIG.Wxpay.AppID,
		global.GVA_CONFIG.Wxpay.MchCertificateSerialNumber,
		global.GVA_CONFIG.Wxpay.MchAPIv3Key,
		global.GVA_CONFIG.Wxpay.PemPath,
		global.GVA_CONFIG.Wxpay.NotifyUrl,
	))

	PluginInit(public, geo.CreateGeoPlug())

	holder(public, private)
}
