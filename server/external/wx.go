package external

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/external/wx"
)

const jscode2session = "https://api.weixin.qq.com/sns/jscode2session"

func GetOpenID(code string) (*wx.Jscode2sessionRes, error) {
	var params RequestParams
	params.URL = jscode2session
	params.Method = "GET"
	params.Body = map[string]interface{}{
		"appid":      global.GVA_CONFIG.Wxpay.AppID,
		"secret":     global.GVA_CONFIG.Wxpay.Secret,
		"js_code":    code,
		"grant_type": "authorization_code",
	}
	res, err := HttpRequest[wx.Jscode2sessionRes](params)
	if err != nil {
		return nil, err
	}
	return res, nil
}
