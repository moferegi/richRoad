package client

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SysConfigApi struct{}

var publicConfigKeyAllowlist = map[string]struct{}{
	"app_logo":                            {},
	"app_name":                            {},
	"currency_suffix":                     {},
	"currency_symbol":                     {},
	"order_logistics_enabled":             {},
	"order_refund_enabled":                {},
	"payment_qr_code":                     {},
	"payment_auto_enabled":                {},
	"payment_manual_qrcode_enabled":       {},
	"payment_manual_contact_enabled":      {},
	"payment_manual_methods":              {},
	"payment_uni_preferred_methods":       {},
	"payment_wechat_enabled":              {},
	"payment_alipay_enabled":              {},
	"payment_bank_cn_enabled":             {},
	"payment_bank_us_enabled":             {},
	"payment_bank_mn_enabled":             {},
	"payment_paypal_enabled":              {},
	"payment_tip_text":                    {},
	"payment_tip_text_color":              {},
	"payment_tip_text_size":               {},
	"points_exchange_rate":                {},
	"presale_home_count":                  {},
	"review_pic_enabled":                  {},
	"shop_kefu_enabled":                   {},
	"sign_in_enabled":                     {},
	"tryon_cost_points":                   {},
	"tryon_fail_refund_percent":           {},
	"tryon_append_parsing_failed_tip":     {},
	"tryon_append_refiner_failed_tip":     {},
	"tryon_parsing_failed_tip_text":       {},
	"tryon_refiner_failed_tip_text":       {},
	"tryon_guest_init_points":             {},
	"tryon_invite_register_reward_points": {},
	"tryon_models":                        {},
	"tryon_recharge_plans":                {},
	"tryon_register_reward_points":        {},
}

func isPublicConfigKeyAllowed(key string) bool {
	_, ok := publicConfigKeyAllowlist[key]
	return ok
}

func isSysConfigAdmin(authorityId uint) bool {
	return authorityId == 888 || authorityId == 8881
}

// GetSysConfigList 分页获取系统参数列表
// @Tags SysConfig
// @Summary 分页获取系统参数列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.SysConfigSearch true "分页获取系统参数列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /sysConfig/getSysConfigList [get]
func (s *SysConfigApi) GetSysConfigList(c *gin.Context) {
	if !isSysConfigAdmin(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return
	}

	var pageInfo request.SysConfigSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	list, total, err := sysConfigService.GetSysConfigList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, i18n.T(c, "getSuccess"), c)
}

// GetModelCallLogList 分页获取模型调用日志列表
// @Tags SysConfig
// @Summary 分页获取模型调用日志列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.ModelCallLogSearch true "分页获取模型调用日志列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /sysConfig/getModelCallLogList [get]
func (s *SysConfigApi) GetModelCallLogList(c *gin.Context) {
	if !isSysConfigAdmin(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return
	}

	var pageInfo request.ModelCallLogSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	list, total, err := sysConfigService.GetModelCallLogList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取模型调用日志失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, i18n.T(c, "getSuccess"), c)
}

// UpdateSysConfig 更新系统参数
// @Tags SysConfig
// @Summary 更新系统参数
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body updateSysConfigReq true "更新系统参数"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /sysConfig/updateSysConfig [put]
func (s *SysConfigApi) UpdateSysConfig(c *gin.Context) {
	authorityId := utils.GetUserAuthorityId(c)
	if !isSysConfigAdmin(authorityId) {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return
	}

	var req struct {
		ID          uint   `json:"id" binding:"required"`
		ConfigValue string `json:"configValue"`
		Remark      string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if err := sysConfigService.UpdateSysConfig(req.ID, req.ConfigValue, req.Remark); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "updateFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
}

// GetSysConfigByGroup 按分组获取系统参数
// @Tags SysConfig
// @Summary 按分组获取系统参数
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param configGroup query string true "配置分组"
// @Success 200 {object} response.Response{data=[]client.SysConfig,msg=string} "获取成功"
// @Router /sysConfig/getSysConfigByGroup [get]
func (s *SysConfigApi) GetSysConfigByGroup(c *gin.Context) {
	group := c.Query("configGroup")
	if group == "" {
		response.FailWithMessage(i18n.T(c, "configGroupRequired"), c)
		return
	}
	list, err := sysConfigService.GetConfigByGroup(group)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(list, i18n.T(c, "getSuccess"), c)
}

// GetSysConfigByKey 按key获取单条系统参数
// @Tags SysConfig
// @Summary 按key获取单条系统参数
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param configKey query string true "配置键名"
// @Success 200 {object} response.Response{data=string,msg=string} "获取成功"
// @Router /sysConfig/getSysConfigByKey [get]
func (s *SysConfigApi) GetSysConfigByKey(c *gin.Context) {
	key := c.Query("configKey")
	if key == "" {
		response.FailWithMessage(i18n.T(c, "configKeyRequired"), c)
		return
	}
	if !isPublicConfigKeyAllowed(key) {
		response.FailWithMessage(i18n.T(c, "configNotPublic"), c)
		return
	}
	val, err := sysConfigService.GetConfigByKey(key)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(val, i18n.T(c, "getSuccess"), c)
}

// GetLoginConfig 获取登录相关配置（公开接口）
// @Tags SysConfig
// @Summary 获取登录相关配置
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]string,msg=string} "获取成功"
// @Router /sysConfig/getLoginConfig [get]
func (s *SysConfigApi) GetLoginConfig(c *gin.Context) {
	phoneEnabled, _ := sysConfigService.GetConfigByKey("phone_login_enabled")
	usernameEnabled, _ := sysConfigService.GetConfigByKey("username_login_enabled")
	defaultMethod, _ := sysConfigService.GetConfigByKey("default_login_method")
	passwordChangeEnabled, _ := sysConfigService.GetConfigByKey("password_change_enabled")
	usernameRegex, _ := sysConfigService.GetConfigByKey("username_regex")
	passwordRegex, _ := sysConfigService.GetConfigByKey("password_regex")
	usernameRegexTip, _ := sysConfigService.GetConfigByKey("username_regex_tip")
	passwordRegexTip, _ := sysConfigService.GetConfigByKey("password_regex_tip")
	response.OkWithDetailed(map[string]string{
		"phone_login_enabled":     phoneEnabled,
		"username_login_enabled":  usernameEnabled,
		"default_login_method":    defaultMethod,
		"password_change_enabled": passwordChangeEnabled,
		"username_regex":          usernameRegex,
		"password_regex":          passwordRegex,
		"username_regex_tip":      usernameRegexTip,
		"password_regex_tip":      passwordRegexTip,
	}, i18n.T(c, "getSuccess"), c)
}

// GetTryonConfig 获取试衣配置（公开接口）
// @Tags SysConfig
// @Summary 获取试衣配置
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]string,msg=string} "获取成功"
// @Router /sysConfig/getTryonConfig [get]
func (s *SysConfigApi) GetTryonConfig(c *gin.Context) {
	guestInit, _ := sysConfigService.GetConfigByKey("tryon_guest_init_points")
	registerReward, _ := sysConfigService.GetConfigByKey("tryon_register_reward_points")
	inviteRegisterReward, _ := sysConfigService.GetConfigByKey("tryon_invite_register_reward_points")
	costPoints, _ := sysConfigService.GetConfigByKey("tryon_cost_points")
	failRefundPercent, _ := sysConfigService.GetConfigByKey("tryon_fail_refund_percent")
	appendParsingFailedTip, _ := sysConfigService.GetConfigByKey("tryon_append_parsing_failed_tip")
	parsingFailedTipText, _ := sysConfigService.GetConfigByKey("tryon_parsing_failed_tip_text")
	appendRefinerFailedTip, _ := sysConfigService.GetConfigByKey("tryon_append_refiner_failed_tip")
	refinerFailedTipText, _ := sysConfigService.GetConfigByKey("tryon_refiner_failed_tip_text")
	tryonModels, _ := sysConfigService.GetConfigByKey("tryon_models")

	if guestInit == "" {
		guestInit = "0"
	}
	if registerReward == "" {
		registerReward = "8"
	}
	if inviteRegisterReward == "" {
		inviteRegisterReward = "0"
	}
	if costPoints == "" {
		costPoints = "1"
	}
	if failRefundPercent == "" {
		failRefundPercent = "100"
	}
	if strings.TrimSpace(appendParsingFailedTip) == "" {
		appendParsingFailedTip = "false"
	}
	if strings.TrimSpace(parsingFailedTipText) == "" {
		parsingFailedTipText = `{"zh":"试衣成功，但分割增强失败，相关金币已退回","en":"Try-on succeeded, but parsing enhancement failed. Related coins have been refunded.","mn":"Туршилт амжилттай боловч segmentation enhancement амжилтгүй боллоо. Холбогдох зоос буцаан олгогдлоо."}`
	}
	if strings.TrimSpace(appendRefinerFailedTip) == "" {
		appendRefinerFailedTip = "true"
	}
	if strings.TrimSpace(refinerFailedTipText) == "" {
		refinerFailedTipText = `{"zh":"试衣成功，但精修失败，金币已退回","en":"Try-on succeeded, but refiner failed. Coins have been refunded.","mn":"Туршилт амжилттай боловч нарийвчлал амжилтгүй боллоо. Зоос буцаан олгогдсон."}`
	}
	if strings.TrimSpace(tryonModels) == "" {
		tryonModels = defaultTryonModelsConfig()
	}

	response.OkWithDetailed(map[string]string{
		"tryon_guest_init_points":             guestInit,
		"tryon_register_reward_points":        registerReward,
		"tryon_invite_register_reward_points": inviteRegisterReward,
		"tryon_cost_points":                   costPoints,
		"tryon_fail_refund_percent":           failRefundPercent,
		"tryon_append_parsing_failed_tip":     appendParsingFailedTip,
		"tryon_parsing_failed_tip_text":       parsingFailedTipText,
		"tryon_append_refiner_failed_tip":     appendRefinerFailedTip,
		"tryon_refiner_failed_tip_text":       refinerFailedTipText,
		"tryon_models":                        tryonModels,
	}, i18n.T(c, "getSuccess"), c)
}

// GetAliyunTryonQuotaEstimate 获取阿里试衣模型剩余额度估算（管理端）
// @Tags SysConfig
// @Summary 获取阿里试衣模型剩余额度估算（管理端）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param modelKey query string false "模型key，不传则返回全部阿里模型"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取成功"
// @Router /sysConfig/getAliyunTryonQuotaEstimate [get]
func (s *SysConfigApi) GetAliyunTryonQuotaEstimate(c *gin.Context) {
	if !isSysConfigAdmin(utils.GetUserAuthorityId(c)) {
		response.FailWithMessage(i18n.T(c, "noPermission"), c)
		return
	}

	var query request.AliyunTryonQuotaSearch
	if err := c.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	list, err := sysConfigService.GetAliyunTryonQuotaEstimate(query.ModelKey)
	if err != nil {
		global.GVA_LOG.Error("获取阿里模型额度估算失败", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}

	response.OkWithDetailed(gin.H{
		"list": list,
		"meta": gin.H{
			"source":               "local_success_count_estimate",
			"officialApiAvailable": false,
			"officialHint":         "阿里云百炼当前未提供可由 API-Key 直接查询免费额度余量的公开HTTP接口，请以控制台数据为准。",
		},
	}, i18n.T(c, "getSuccess"), c)
}

// GetAnnouncementConfig 获取公告配置
// @Tags SysConfig
// @Summary 获取公告配置
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]string,msg=string} "获取成功"
// @Router /sysConfig/getAnnouncementConfig [get]
func (s *SysConfigApi) GetAnnouncementConfig(c *gin.Context) {
	result, err := sysConfigService.GetAnnouncementConfig()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(result, i18n.T(c, "getSuccess"), c)
}

func parseBoolConfig(raw string, defaultVal bool) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return defaultVal
	}
}

func getConfigBoolOrDefault(key string, defaultVal bool) bool {
	val, err := sysConfigService.GetConfigByKey(key)
	if err != nil {
		return defaultVal
	}
	return parseBoolConfig(val, defaultVal)
}

func defaultPaymentMethodName(key string) map[string]string {
	switch key {
	case "qrcode":
		return map[string]string{"zh": "二维码支付", "en": "QR Payment", "mn": "QR төлбөр"}
	case "contact":
		return map[string]string{"zh": "联系客服", "en": "Contact Support", "mn": "Хэрэглэгчийн дэмжлэг"}
	case "wechat":
		return map[string]string{"zh": "微信支付", "en": "WeChat Pay", "mn": "WeChat Pay"}
	case "alipay":
		return map[string]string{"zh": "支付宝", "en": "Alipay", "mn": "Alipay"}
	case "bank_card_cn":
		return map[string]string{"zh": "银行卡(国内)", "en": "Bank Card (CN)", "mn": "Банкны карт (CN)"}
	case "bank_card_us":
		return map[string]string{"zh": "银行卡(美国)", "en": "Bank Card (US)", "mn": "Банкны карт (US)"}
	case "bank_card_mn":
		return map[string]string{"zh": "银行卡(蒙古)", "en": "Bank Card (MN)", "mn": "Банкны карт (MN)"}
	case "paypal":
		return map[string]string{"zh": "PayPal", "en": "PayPal", "mn": "PayPal"}
	default:
		fallback := strings.TrimSpace(key)
		if fallback == "" {
			fallback = "Payment"
		}
		return map[string]string{"zh": fallback, "en": fallback, "mn": fallback}
	}
}

func extractStringValue(value interface{}) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(typed)
	case json.Number:
		return strings.TrimSpace(typed.String())
	case float64:
		return strings.TrimSpace(strconv.FormatFloat(typed, 'f', -1, 64))
	case float32:
		return strings.TrimSpace(strconv.FormatFloat(float64(typed), 'f', -1, 32))
	case int:
		return strconv.Itoa(typed)
	case int32:
		return strconv.Itoa(int(typed))
	case int64:
		return strconv.FormatInt(typed, 10)
	case bool:
		if typed {
			return "true"
		}
		return "false"
	default:
		return strings.TrimSpace(fmt.Sprintf("%v", typed))
	}
}

func extractBoolValue(value interface{}, defaultVal bool) bool {
	if value == nil {
		return defaultVal
	}
	return parseBoolConfig(extractStringValue(value), defaultVal)
}

func extractIntValue(value interface{}, defaultVal int) int {
	text := extractStringValue(value)
	if text == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(text)
	if err != nil {
		return defaultVal
	}
	return n
}

func extractI18nMap(value interface{}) map[string]string {
	result := map[string]string{}
	appendFallback := func(text string) map[string]string {
		text = strings.TrimSpace(text)
		if text == "" {
			return result
		}
		result["zh"] = text
		result["en"] = text
		result["mn"] = text
		return result
	}

	switch typed := value.(type) {
	case nil:
		return result
	case string:
		text := strings.TrimSpace(typed)
		if text == "" {
			return result
		}
		if strings.HasPrefix(text, "{") {
			var parsed map[string]interface{}
			if err := json.Unmarshal([]byte(text), &parsed); err == nil {
				return extractI18nMap(parsed)
			}
		}
		return appendFallback(text)
	case map[string]interface{}:
		for k, v := range typed {
			key := strings.TrimSpace(k)
			if key == "" {
				continue
			}
			valueText := extractStringValue(v)
			if valueText == "" {
				continue
			}
			result[key] = valueText
		}
		return result
	case map[string]string:
		for k, v := range typed {
			key := strings.TrimSpace(k)
			valueText := strings.TrimSpace(v)
			if key == "" || valueText == "" {
				continue
			}
			result[key] = valueText
		}
		return result
	default:
		return appendFallback(extractStringValue(typed))
	}
}

func pickMethodLabel(name map[string]string, fallback string) string {
	for _, key := range []string{"zh", "en", "mn"} {
		if text := strings.TrimSpace(name[key]); text != "" {
			return text
		}
	}
	for _, text := range name {
		if text = strings.TrimSpace(text); text != "" {
			return text
		}
	}
	return strings.TrimSpace(fallback)
}

func methodSortValue(method gin.H) int {
	if method == nil {
		return 999
	}
	if val, ok := method["sort"]; ok {
		return extractIntValue(val, 999)
	}
	return 999
}

func methodKeyValue(method gin.H) string {
	if method == nil {
		return ""
	}
	return strings.TrimSpace(extractStringValue(method["key"]))
}

func appendLegacyPaymentMethod(methods []gin.H, key string, manual bool, sortValue int) []gin.H {
	name := defaultPaymentMethodName(key)
	label := pickMethodLabel(name, key)
	methods = append(methods, gin.H{
		"key":          key,
		"label":        label,
		"name":         name,
		"manual":       manual,
		"enabled":      true,
		"sort":         sortValue,
		"image":        "",
		"externalPath": "",
		"copyText":     map[string]string{},
	})
	return methods
}

func defaultUniPreferredPayMethodCopyText() map[string]string {
	return map[string]string{
		"zh": "您好，我的订单号是 {orderID}，期望使用 {payMethod} 支付，请协助提供收款方式并处理订单。",
		"en": "Hi, my order number is {orderID}. I expect to pay via {payMethod}. Please provide the receiving method and help process this order.",
		"mn": "Сайн байна уу, миний захиалгын дугаар {orderID}. Би {payMethod} аргаар төлөхийг хүсэж байна. Хүлээн авах мэдээлэл өгч, захиалгыг боловсруулж өгнө үү.",
	}
}

func defaultUniPreferredPayMethods() []gin.H {
	copyText := defaultUniPreferredPayMethodCopyText()
	methods := []gin.H{
		{
			"key":          "wechat",
			"label":        "微信支付",
			"name":         defaultPaymentMethodName("wechat"),
			"enabled":      true,
			"sort":         10,
			"image":        "cloth-on/up/wechat.png",
			"externalPath": "",
			"copyText":     copyText,
		},
		{
			"key":          "alipay",
			"label":        "支付宝",
			"name":         defaultPaymentMethodName("alipay"),
			"enabled":      true,
			"sort":         20,
			"image":        "cloth-on/up/alipay.png",
			"externalPath": "",
			"copyText":     copyText,
		},
		{
			"key":          "bank_card_cn",
			"label":        "银行卡(国内)",
			"name":         defaultPaymentMethodName("bank_card_cn"),
			"enabled":      true,
			"sort":         30,
			"image":        "cloth-on/up/bank-card-cn.png",
			"externalPath": "",
			"copyText":     copyText,
		},
	}
	return methods
}

func parseUniPreferredPayMethods(raw string) []gin.H {
	methods := make([]gin.H, 0, 8)
	if strings.TrimSpace(raw) == "" {
		return methods
	}

	var configuredMethods []map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &configuredMethods); err != nil {
		return methods
	}

	for _, item := range configuredMethods {
		key := strings.ToLower(strings.TrimSpace(extractStringValue(item["key"])))
		if key == "" || key == "qrcode" || key == "contact" {
			continue
		}
		if !extractBoolValue(item["enabled"], true) {
			continue
		}

		name := extractI18nMap(item["name"])
		if len(name) == 0 {
			name = defaultPaymentMethodName(key)
		}

		label := strings.TrimSpace(extractStringValue(item["label"]))
		if label == "" {
			label = pickMethodLabel(name, key)
		}

		copyText := extractI18nMap(item["copyText"])
		if len(copyText) == 0 {
			copyText = defaultUniPreferredPayMethodCopyText()
		}

		methods = append(methods, gin.H{
			"key":          key,
			"label":        label,
			"name":         name,
			"enabled":      true,
			"sort":         extractIntValue(item["sort"], 999),
			"image":        strings.TrimSpace(extractStringValue(item["image"])),
			"externalPath": strings.TrimSpace(extractStringValue(item["externalPath"])),
			"copyText":     copyText,
		})
	}

	sort.SliceStable(methods, func(i, j int) bool {
		iSort := methodSortValue(methods[i])
		jSort := methodSortValue(methods[j])
		if iSort == jSort {
			return methodKeyValue(methods[i]) < methodKeyValue(methods[j])
		}
		return iSort < jSort
	})

	return methods
}

func defaultTryonModelsConfig() string {
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
		"refinerModelKey": "aliyun_aitryon_refiner",
		"parsingModelKey": "aliyun_aitryon_parsing",
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
		"refinerModelKey": "aliyun_aitryon_refiner",
		"parsingModelKey": "aliyun_aitryon_parsing",
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
		"token": "",
		"supportsBeautify": true
	}
]`
}

// GetPaymentConfig 获取支付方式配置（公开接口）
// @Tags SysConfig
// @Summary 获取支付方式配置
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取成功"
// @Router /sysConfig/getPaymentConfig [get]
func (s *SysConfigApi) GetPaymentConfig(c *gin.Context) {
	autoEnabled := getConfigBoolOrDefault("payment_auto_enabled", false)
	manualQrcodeEnabled := getConfigBoolOrDefault("payment_manual_qrcode_enabled", true)
	manualContactEnabled := getConfigBoolOrDefault("payment_manual_contact_enabled", true)

	wechatEnabled := autoEnabled && getConfigBoolOrDefault("payment_wechat_enabled", false)
	alipayEnabled := autoEnabled && getConfigBoolOrDefault("payment_alipay_enabled", false)
	bankCNEnabled := autoEnabled && getConfigBoolOrDefault("payment_bank_cn_enabled", false)
	bankUSEnabled := autoEnabled && getConfigBoolOrDefault("payment_bank_us_enabled", false)
	bankMNEnabled := autoEnabled && getConfigBoolOrDefault("payment_bank_mn_enabled", false)
	paypalEnabled := autoEnabled && getConfigBoolOrDefault("payment_paypal_enabled", false)

	autoChannelEnabled := map[string]bool{
		"wechat":       wechatEnabled,
		"alipay":       alipayEnabled,
		"bank_card_cn": bankCNEnabled,
		"bank_card_us": bankUSEnabled,
		"bank_card_mn": bankMNEnabled,
		"paypal":       paypalEnabled,
	}

	methods := make([]gin.H, 0, 8)
	manualMethodsRaw, _ := sysConfigService.GetConfigByKey("payment_manual_methods")
	if strings.TrimSpace(manualMethodsRaw) != "" {
		var configuredMethods []map[string]interface{}
		if err := json.Unmarshal([]byte(manualMethodsRaw), &configuredMethods); err == nil {
			for _, item := range configuredMethods {
				key := strings.ToLower(strings.TrimSpace(extractStringValue(item["key"])))
				if key == "" {
					continue
				}

				manual := extractBoolValue(item["manual"], key == "qrcode" || key == "contact")
				enabled := extractBoolValue(item["enabled"], true)
				switch key {
				case "qrcode":
					enabled = enabled && manualQrcodeEnabled
				case "contact":
					enabled = enabled && manualContactEnabled
				default:
					if !manual {
						enabled = enabled && autoEnabled
						if legacyEnabled, ok := autoChannelEnabled[key]; ok {
							enabled = enabled && legacyEnabled
						}
					}
				}
				if !enabled {
					continue
				}

				name := extractI18nMap(item["name"])
				if len(name) == 0 {
					name = defaultPaymentMethodName(key)
				}
				label := strings.TrimSpace(extractStringValue(item["label"]))
				if label == "" {
					label = pickMethodLabel(name, key)
				}

				methods = append(methods, gin.H{
					"key":          key,
					"label":        label,
					"name":         name,
					"manual":       manual,
					"enabled":      true,
					"sort":         extractIntValue(item["sort"], 999),
					"image":        strings.TrimSpace(extractStringValue(item["image"])),
					"externalPath": strings.TrimSpace(extractStringValue(item["externalPath"])),
					"copyText":     extractI18nMap(item["copyText"]),
				})
			}
		}
	}

	if len(methods) == 0 {
		if manualQrcodeEnabled {
			methods = appendLegacyPaymentMethod(methods, "qrcode", true, 10)
		}
		if manualContactEnabled {
			methods = appendLegacyPaymentMethod(methods, "contact", true, 20)
		}
		if wechatEnabled {
			methods = appendLegacyPaymentMethod(methods, "wechat", false, 30)
		}
		if alipayEnabled {
			methods = appendLegacyPaymentMethod(methods, "alipay", false, 40)
		}
		if bankCNEnabled {
			methods = appendLegacyPaymentMethod(methods, "bank_card_cn", false, 50)
		}
		if bankUSEnabled {
			methods = appendLegacyPaymentMethod(methods, "bank_card_us", false, 60)
		}
		if bankMNEnabled {
			methods = appendLegacyPaymentMethod(methods, "bank_card_mn", false, 70)
		}
		if paypalEnabled {
			methods = appendLegacyPaymentMethod(methods, "paypal", false, 80)
		}
	}

	sort.SliceStable(methods, func(i, j int) bool {
		iSort := methodSortValue(methods[i])
		jSort := methodSortValue(methods[j])
		if iSort == jSort {
			return methodKeyValue(methods[i]) < methodKeyValue(methods[j])
		}
		return iSort < jSort
	})

	// 至少保留一个人工渠道，避免用户无法支付
	if len(methods) == 0 {
		methods = appendLegacyPaymentMethod(methods, "contact", true, 999)
		manualContactEnabled = true
	}

	response.OkWithDetailed(gin.H{
		"autoEnabled": autoEnabled,
		"manual": gin.H{
			"qrcode":  manualQrcodeEnabled,
			"contact": manualContactEnabled,
		},
		"channels": gin.H{
			"wechat":  wechatEnabled,
			"alipay":  alipayEnabled,
			"bank_cn": bankCNEnabled,
			"bank_us": bankUSEnabled,
			"bank_mn": bankMNEnabled,
			"paypal":  paypalEnabled,
		},
		"methods": methods,
	}, i18n.T(c, "getSuccess"), c)
}

// GetUniPreferredPayConfig 获取uni联系客服页期望支付方式配置（公开接口）
// @Tags SysConfig
// @Summary 获取uni期望支付方式配置
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取成功"
// @Router /sysConfig/getUniPreferredPayConfig [get]
func (s *SysConfigApi) GetUniPreferredPayConfig(c *gin.Context) {
	preferredMethodsRaw, _ := sysConfigService.GetConfigByKey("payment_uni_preferred_methods")
	preferredMethods := parseUniPreferredPayMethods(preferredMethodsRaw)
	if len(preferredMethods) == 0 {
		preferredMethods = defaultUniPreferredPayMethods()
	}

	response.OkWithDetailed(gin.H{
		"methods": preferredMethods,
	}, i18n.T(c, "getSuccess"), c)
}
