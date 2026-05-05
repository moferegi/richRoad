package client

import (
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SysConfigApi struct{}

var publicConfigKeyAllowlist = map[string]struct{}{
	"app_logo":                       {},
	"app_name":                       {},
	"currency_suffix":                {},
	"currency_symbol":                {},
	"order_logistics_enabled":        {},
	"order_refund_enabled":           {},
	"payment_qr_code":                {},
	"payment_auto_enabled":           {},
	"payment_manual_qrcode_enabled":  {},
	"payment_manual_contact_enabled": {},
	"payment_wechat_enabled":         {},
	"payment_alipay_enabled":         {},
	"payment_bank_cn_enabled":        {},
	"payment_bank_us_enabled":        {},
	"payment_bank_mn_enabled":        {},
	"payment_paypal_enabled":         {},
	"payment_tip_text":               {},
	"payment_tip_text_color":         {},
	"payment_tip_text_size":          {},
	"points_exchange_rate":           {},
	"presale_home_count":             {},
	"review_pic_enabled":             {},
	"shop_kefu_enabled":              {},
	"sign_in_enabled":                {},
	"tryon_cost_points":              {},
	"tryon_fail_refund_percent":      {},
	"tryon_guest_init_points":        {},
	"tryon_models":                   {},
	"tryon_register_reward_points":   {},
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
		response.FailWithMessage("无权限查看系统参数列表", c)
		return
	}

	var pageInfo request.SysConfigSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := sysConfigService.GetSysConfigList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
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
		response.FailWithMessage("无权限更新系统参数", c)
		return
	}

	var req struct {
		ID          uint   `json:"id" binding:"required"`
		ConfigValue string `json:"configValue"`
		Remark      string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysConfigService.UpdateSysConfig(req.ID, req.ConfigValue, req.Remark); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
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
		response.FailWithMessage("configGroup不能为空", c)
		return
	}
	list, err := sysConfigService.GetConfigByGroup(group)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
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
		response.FailWithMessage("configKey不能为空", c)
		return
	}
	if !isPublicConfigKeyAllowed(key) {
		response.FailWithMessage("该配置不对外开放", c)
		return
	}
	val, err := sysConfigService.GetConfigByKey(key)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(val, "获取成功", c)
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
	}, "获取成功", c)
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
	costPoints, _ := sysConfigService.GetConfigByKey("tryon_cost_points")
	failRefundPercent, _ := sysConfigService.GetConfigByKey("tryon_fail_refund_percent")
	tryonModels, _ := sysConfigService.GetConfigByKey("tryon_models")

	if guestInit == "" {
		guestInit = "3"
	}
	if registerReward == "" {
		registerReward = "8"
	}
	if costPoints == "" {
		costPoints = "1"
	}
	if failRefundPercent == "" {
		failRefundPercent = "100"
	}
	if strings.TrimSpace(tryonModels) == "" {
		tryonModels = defaultTryonModelsConfig()
	}

	response.OkWithDetailed(map[string]string{
		"tryon_guest_init_points":      guestInit,
		"tryon_register_reward_points": registerReward,
		"tryon_cost_points":            costPoints,
		"tryon_fail_refund_percent":    failRefundPercent,
		"tryon_models":                 tryonModels,
	}, "获取成功", c)
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
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
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

func defaultTryonModelsConfig() string {
	return `[{"key":"aliyun_aitryon","enabled":true,"scenes":["clothes","shoes"],"model":"aitryon","name":{"zh":"阿里AI试衣（基础）","en":"Aliyun AI Try-On (Basic)","mn":"Aliyun AI өмсгөл (Суурь)"},"desc":{"zh":"基础版试衣模型，速度更快，适合日常试衣。","en":"Basic try-on model with faster generation for everyday use.","mn":"Өдөр тутмын туршилтад тохирох, хурдан суурь загвар."},"cost":1,"provider":"aliyun","mode":"prod","url":"https://dashscope.aliyuncs.com/api/v1/services/aigc/image2image/image-synthesis","taskQueryUrl":"https://dashscope.aliyuncs.com/api/v1/tasks/{task_id}","token":"","resolution":-1,"restoreFace":true},{"key":"aliyun_aitryon_plus","enabled":true,"scenes":["clothes","shoes"],"model":"aitryon-plus","name":{"zh":"阿里AI试衣（Plus）","en":"Aliyun AI Try-On (Plus)","mn":"Aliyun AI өмсгөл (Plus)"},"desc":{"zh":"Plus版细节更好，适合高质量试衣图。","en":"Higher quality rendering with better texture and logo details.","mn":"Нэхмэл, логог илүү сайн сэргээдэг өндөр чанарын загвар."},"cost":1,"provider":"aliyun","mode":"prod","url":"https://dashscope.aliyuncs.com/api/v1/services/aigc/image2image/image-synthesis","taskQueryUrl":"https://dashscope.aliyuncs.com/api/v1/tasks/{task_id}","token":"","resolution":-1,"restoreFace":true},{"key":"aliyun_aitryon_parsing","enabled":true,"scenes":["takeoff"],"model":"aitryon-parsing-v1","name":{"zh":"阿里取衣分割","en":"Aliyun Takeoff Parsing","mn":"Aliyun хувцас салгах"},"desc":{"zh":"用于取衣区分割模特服饰并输出可用服饰图。","en":"Segments garment regions for takeoff area and outputs reusable garment images.","mn":"Загварын хувцсыг ялган авч, дахин ашиглах зургийг гаргана."},"cost":1,"provider":"aliyun","mode":"prod","url":"https://dashscope.aliyuncs.com/api/v1/services/vision/image-process/process","token":"","clothesType":["upper"]}]`
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

	methods := make([]gin.H, 0, 8)
	if manualQrcodeEnabled {
		methods = append(methods, gin.H{"key": "qrcode", "label": "二维码支付", "manual": true, "enabled": true})
	}
	if manualContactEnabled {
		methods = append(methods, gin.H{"key": "contact", "label": "联系客服", "manual": true, "enabled": true})
	}
	if wechatEnabled {
		methods = append(methods, gin.H{"key": "wechat", "label": "微信支付", "manual": false, "enabled": true})
	}
	if alipayEnabled {
		methods = append(methods, gin.H{"key": "alipay", "label": "支付宝", "manual": false, "enabled": true})
	}
	if bankCNEnabled {
		methods = append(methods, gin.H{"key": "bank_card_cn", "label": "银行卡(国内)", "manual": false, "enabled": true})
	}
	if bankUSEnabled {
		methods = append(methods, gin.H{"key": "bank_card_us", "label": "Bank Card (US)", "manual": false, "enabled": true})
	}
	if bankMNEnabled {
		methods = append(methods, gin.H{"key": "bank_card_mn", "label": "Банкны карт (MN)", "manual": false, "enabled": true})
	}
	if paypalEnabled {
		methods = append(methods, gin.H{"key": "paypal", "label": "PayPal", "manual": false, "enabled": true})
	}

	// 至少保留一个人工渠道，避免用户无法支付
	if len(methods) == 0 {
		methods = append(methods, gin.H{"key": "contact", "label": "联系客服", "manual": true, "enabled": true})
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
	}, "获取成功", c)
}
