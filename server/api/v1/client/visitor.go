package client

import (
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VisitorApi struct {
}

var visitorService = service.ServiceGroupApp.ClientServiceGroup.VisitorService

func normalizeVisitorText(raw string, maxLen int) string {
	text := strings.TrimSpace(raw)
	if maxLen <= 0 {
		return ""
	}
	runes := []rune(text)
	if len(runes) > maxLen {
		return string(runes[:maxLen])
	}
	return text
}

// Heartbeat 访客心跳上报
// @Tags Visitor
// @Summary 访客心跳上报（无需认证）
// @accept application/json
// @Produce application/json
// @Param data body clientReq.HeartbeatRequest true "心跳数据"
// @Success 200 {object} response.Response{msg=string} "上报成功"
// @Router /visitor/heartbeat [post]
func (visitorApi *VisitorApi) Heartbeat(c *gin.Context) {
	var req clientReq.HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	req.VisitorID = normalizeVisitorText(req.VisitorID, 64)
	if req.VisitorID == "" {
		failClientWithKey(c, "visitorIDRequired")
		return
	}
	req.PagePath = normalizeVisitorText(req.PagePath, 255)
	req.Referer = normalizeVisitorText(req.Referer, 512)
	req.Platform = normalizeVisitorText(req.Platform, 32)
	req.EventCategory = normalizeVisitorText(req.EventCategory, 64)
	req.EventAction = normalizeVisitorText(req.EventAction, 64)
	req.EventLabel = normalizeVisitorText(req.EventLabel, 128)
	req.EventExtra = normalizeVisitorText(req.EventExtra, 4096)
	req.Language = normalizeVisitorText(req.Language, 16)
	req.SessionID = normalizeVisitorText(req.SessionID, 64)

	// 获取真实IP
	ip := c.ClientIP()
	userAgent := normalizeVisitorText(c.GetHeader("User-Agent"), 512)
	if !visitorService.CheckHeartbeatRateLimit(ip) {
		failClientWithKey(c, "visitorHeartbeatTooFrequent")
		return
	}
	if !visitorService.ShouldPersistHeartbeat(req.VisitorID) {
		response.OkWithMessage(i18n.T(c, "ok"), c)
		return
	}

	// 公开接口不主动解析 token，只有经过鉴权中间件注入 claims 时才识别登录用户。
	var userID uint
	if claimsValue, exists := c.Get("claims"); exists {
		if claims, ok := claimsValue.(*systemReq.CustomClaims); ok && claims != nil {
			userID = claims.BaseClaims.ID
		}
	}

	log := &client.VisitorLog{
		VisitorID:     req.VisitorID,
		UserID:        userID,
		IP:            ip,
		UserAgent:     userAgent,
		Platform:      req.Platform,
		PagePath:      req.PagePath,
		Referer:       req.Referer,
		EventCategory: req.EventCategory,
		EventAction:   req.EventAction,
		EventLabel:    req.EventLabel,
		EventValue:    req.EventValue,
		EventExtra:    req.EventExtra,
		ScreenWidth:   req.ScreenWidth,
		ScreenHeight:  req.ScreenHeight,
		Language:      req.Language,
		SessionID:     req.SessionID,
	}

	if err := visitorService.Heartbeat(log); err != nil {
		global.GVA_LOG.Error("访客心跳上报失败!", zap.Error(err))
		failClientWithKey(c, "fail")
		return
	}
	response.OkWithMessage(i18n.T(c, "ok"), c)
}

// GetVisitorLogList 获取访客日志列表
// @Tags Visitor
// @Summary 获取访客日志列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query clientReq.VisitorLogSearch true "分页获取访客日志列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /visitor/getVisitorLogList [get]
func (visitorApi *VisitorApi) GetVisitorLogList(c *gin.Context) {
	if !isClientAdminAuthority(utils.GetUserAuthorityId(c)) {
		failClientWithKey(c, "noPermission")
		return
	}
	var pageInfo clientReq.VisitorLogSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	list, total, err := visitorService.GetVisitorLogList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		failClientWithKey(c, "getFail")
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, i18n.T(c, "getSuccess"), c)
}

// GetVisitorSummaryList 获取访客汇总列表
// @Tags Visitor
// @Summary 获取访客汇总统计列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query clientReq.VisitorSummarySearch true "分页获取访客汇总列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /visitor/getVisitorSummaryList [get]
func (visitorApi *VisitorApi) GetVisitorSummaryList(c *gin.Context) {
	if !isClientAdminAuthority(utils.GetUserAuthorityId(c)) {
		failClientWithKey(c, "noPermission")
		return
	}
	var pageInfo clientReq.VisitorSummarySearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	list, total, err := visitorService.GetVisitorSummaryList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		failClientWithKey(c, "getFail")
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, i18n.T(c, "getSuccess"), c)
}

// GetTodayStats 获取今日实时统计
// @Tags Visitor
// @Summary 获取今日实时访客统计
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取成功"
// @Router /visitor/getTodayStats [get]
func (visitorApi *VisitorApi) GetTodayStats(c *gin.Context) {
	if !isClientAdminAuthority(utils.GetUserAuthorityId(c)) {
		failClientWithKey(c, "noPermission")
		return
	}
	stats, err := visitorService.GetTodayStats()
	if err != nil {
		global.GVA_LOG.Error("获取今日统计失败!", zap.Error(err))
		failClientWithKey(c, "getFail")
		return
	}
	response.OkWithDetailed(stats, i18n.T(c, "getSuccess"), c)
}

// GetKefuGuideStats 获取客服引导漏斗统计
// @Tags Visitor
// @Summary 获取客服引导漏斗统计（确认率、复制率、跳转率）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query clientReq.KefuGuideStatsSearch false "客服引导统计查询条件"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取成功"
// @Router /visitor/getKefuGuideStats [get]
func (visitorApi *VisitorApi) GetKefuGuideStats(c *gin.Context) {
	if !isClientAdminAuthority(utils.GetUserAuthorityId(c)) {
		failClientWithKey(c, "noPermission")
		return
	}
	var query clientReq.KefuGuideStatsSearch
	if err := c.ShouldBindQuery(&query); err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}

	stats, err := visitorService.GetKefuGuideStats(query)
	if err != nil {
		global.GVA_LOG.Error("获取客服引导统计失败!", zap.Error(err))
		failClientWithKey(c, "getFail")
		return
	}
	response.OkWithDetailed(stats, i18n.T(c, "getSuccess"), c)
}

// AggregateDailySummary 手动触发日汇总聚合
// @Tags Visitor
// @Summary 手动触发日汇总聚合
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param date query string true "日期(YYYY-MM-DD)"
// @Success 200 {object} response.Response{msg=string} "聚合成功"
// @Router /visitor/aggregateDailySummary [post]
func (visitorApi *VisitorApi) AggregateDailySummary(c *gin.Context) {
	if !isClientAdminAuthority(utils.GetUserAuthorityId(c)) {
		failClientWithKey(c, "noPermission")
		return
	}
	date := c.Query("date")
	if date == "" {
		failClientWithKey(c, "dateRequired")
		return
	}
	if err := visitorService.AggregateDailySummary(date); err != nil {
		global.GVA_LOG.Error("聚合失败!", zap.Error(err))
		failClientWithKey(c, "fail")
		return
	}
	response.OkWithMessage(i18n.T(c, "success"), c)
}
