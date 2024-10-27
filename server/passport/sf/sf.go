package sf

import (
	"encoding/json"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/request"
	"github.com/gofrs/uuid/v5"
	"io"
	"strconv"
	"sync"
	"time"
)

const SERVER_CODE = "EXP_RECE_SEARCH_ROUTES"                                // 路由查询
const SERVER_HOST = "https://bspgw.sf-express.com/std/service"              // 正式接口
const SERVER_HOST_SANDBOX = "https://sfapi-sbox.sf-express.com/std/service" // 沙箱接口

const OAUTH2_SERVER = "https://sfapi.sf-express.com/oauth2/accessToken"              // OAuth2
const OAUTH2_SERVER_SANDBOX = "https://sfapi-sbox.sf-express.com/oauth2/accessToken" // OAuth2 沙箱

type sfPassPort struct {
	accessToken string
	expiresAt   time.Time
	mu          sync.Mutex
}

type OAuthResponse struct {
	ApiResponseID string `json:"apiResponseID"`
	ApiResultCode string `json:"apiResultCode"`
	ApiErrorMsg   string `json:"apiErrorMsg"`
	AccessToken   string `json:"accessToken"`
	ExpiresIn     int    `json:"expiresIn"`
}

var SfPassPort = new(sfPassPort)

func (sf *sfPassPort) GetOAuthToken(partnerID, secret string) (string, error) {
	sf.mu.Lock()
	defer sf.mu.Unlock()
	// 如果accessToken存在且未过期，直接返回
	if sf.accessToken != "" && time.Now().Before(sf.expiresAt) {
		return sf.accessToken, nil
	}

	var req = make(map[string]string)
	req["partnerID"] = partnerID
	req["secret"] = secret
	req["grantType"] = "password"
	resp, err := request.PostForm(OAUTH2_SERVER_SANDBOX, req)
	if err != nil {
		return "", err
	}
	var oauthResp OAuthResponse

	resb, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(resb, &oauthResp)

	// 更新accessToken和expiresAt
	sf.accessToken = oauthResp.AccessToken
	sf.expiresAt = time.Now().Add(time.Duration(oauthResp.ExpiresIn) * time.Second)
	if oauthResp.ApiResultCode != "A1000" {
		return "", errors.New(oauthResp.ApiErrorMsg)
	}
	return oauthResp.AccessToken, nil
}

type MessageData struct {
	TrackingType   int      `json:"trackingType"`
	TrackingNumber []string `json:"trackingNumber"`
}

type SearchRoutersRes struct {
	ApiResponseID string `json:"apiResponseID"`
	ApiErrorMsg   string `json:"apiErrorMsg"`
	ApiResultCode string `json:"apiResultCode"`
	ApiResultData string `json:"apiResultData"`
	Succ          string `json:"succ"`
}

func (sf *sfPassPort) SearchRouters(trackingNumber string) (error, string) {
	partnerID := global.GVA_CONFIG.Sf.PartnerID
	secret := sf.getSecret()
	authCode, err := sf.GetOAuthToken(partnerID, secret)
	if err != nil {
		return err, ""
	}
	var routerRequest = make(map[string]string)

	routerRequest["partnerID"] = partnerID
	u, _ := uuid.NewV4()
	routerRequest["requestID"] = u.String()
	routerRequest["serviceCode"] = SERVER_CODE
	timestampStr := strconv.Itoa(int(time.Now().Unix()))
	routerRequest["timestamp"] = timestampStr
	routerRequest["accessToken"] = authCode

	var msgData MessageData
	msgData.TrackingType = 1
	msgData.TrackingNumber = []string{trackingNumber}
	msgDataB, err := json.Marshal(msgData)
	routerRequest["msgData"] = string(msgDataB)
	resp, err := request.PostForm(SERVER_HOST_SANDBOX, routerRequest)
	if err != nil {
		return err, ""
	}
	var searchRoutersRes SearchRoutersRes
	resb, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, ""
	}
	err = json.Unmarshal(resb, &searchRoutersRes)
	if err != nil {
		return err, ""
	}
	if searchRoutersRes.ApiResultCode != "A1000" {
		return errors.New(searchRoutersRes.ApiErrorMsg), ""
	}
	return nil, searchRoutersRes.ApiResultData
}

func (sf *sfPassPort) getSecret() string {
	if global.GVA_CONFIG.Sf.IsSandbox {
		return global.GVA_CONFIG.Sf.SandCheckCode
	}
	return global.GVA_CONFIG.Sf.CheckCode
}
