package wxpay

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/wxpay/router"
	"github.com/gin-gonic/gin"
)

type WxpayPlugin struct{}

func CreateWxpayPlug(
	MchID string,
	AppId string,
	MchCertificateSerialNumber string,
	MchAPIv3Key string,
	MchAPIv2Key string,
	CertPath string,
	KeyPath string,
	NotifyUrl string,
) *WxpayPlugin {
	global.GlobalConfig.MchID = MchID
	global.GlobalConfig.MchCertificateSerialNumber = MchCertificateSerialNumber
	global.GlobalConfig.MchAPIv3Key = MchAPIv3Key
	global.GlobalConfig.MchAPIv2Key = MchAPIv2Key
	global.GlobalConfig.CertPath = CertPath
	global.GlobalConfig.KeyPath = KeyPath
	global.GlobalConfig.AppID = AppId
	global.GlobalConfig.NotifyUrl = NotifyUrl

	utils.RegisterMenus(
		system.SysBaseMenu{
			Path:      "wxpay",
			Name:      "wxpay",
			Hidden:    false,
			Component: "view/routerHolder.vue",
			Sort:      1000,
			Meta: system.Meta{
				Title: "支付示例",
				Icon:  "school",
			},
		},
		system.SysBaseMenu{
			Path:      "wxpayexample",
			Name:      "wxpayexample",
			Hidden:    false,
			Component: "plugin/wxpay/view/index.vue",
			Sort:      0,
			Meta: system.Meta{
				Title: "支付示例",
				Icon:  "school",
			},
		},
	)

	utils.RegisterApis(
		system.SysApi{
			Path:        "/wxpay/getPayCode",
			Description: "获取支付Code",
			ApiGroup:    "微信支付",
			Method:      "POST",
		},
		system.SysApi{
			Path:        "/wxpay/getOrderById",
			Description: "获取订单支付状态",
			ApiGroup:    "微信支付",
			Method:      "POST",
		},
	)

	return &WxpayPlugin{}
}

func (*WxpayPlugin) Register(group *gin.RouterGroup) {
	router.RouterGroupApp.InitWxpayRouter(group)
}

func (*WxpayPlugin) RouterPath() string {
	return "wxpay"
}
