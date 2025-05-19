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
	PluginInit(public, wxpay.CreateWxpayPlug("1625910538", "wxd7d68c1e9c2c1008", "2EE50AFA04D96FEFDB41816D2F18F32D142C31E2", "9d8609c724026e2d035791b073476288", "", "./resource/wx/apiclient_cert.pem", "./resource/wx/apiclient_key.pem", "https://gin-vue-admin.com/wxpay/payAction"))

	PluginInit(public, geo.CreateGeoPlug())

	holder(public, private)
}
