package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VisitorApi struct {
}

var visitorService = service.ServiceGroupApp.ClientServiceGroup.VisitorService

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
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if req.VisitorID == "" {
		response.FailWithMessage(i18n.T(c, "visitorIDRequired"), c)
		return
	}

	// 获取真实IP
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// 尝试获取登录用户ID（公开接口，可能未登录）
	userID := utils.GetUserID(c)

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
		response.FailWithMessage(i18n.T(c, "fail"), c)
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
	var pageInfo clientReq.VisitorLogSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	list, total, err := visitorService.GetVisitorLogList(pageInfo)
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
	var pageInfo clientReq.VisitorSummarySearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	list, total, err := visitorService.GetVisitorSummaryList(pageInfo)
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

// GetTodayStats 获取今日实时统计
// @Tags Visitor
// @Summary 获取今日实时访客统计
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "获取成功"
// @Router /visitor/getTodayStats [get]
func (visitorApi *VisitorApi) GetTodayStats(c *gin.Context) {
	stats, err := visitorService.GetTodayStats()
	if err != nil {
		global.GVA_LOG.Error("获取今日统计失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
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
	var query clientReq.KefuGuideStatsSearch
	if err := c.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}

	stats, err := visitorService.GetKefuGuideStats(query)
	if err != nil {
		global.GVA_LOG.Error("获取客服引导统计失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
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
	date := c.Query("date")
	if date == "" {
		response.FailWithMessage(i18n.T(c, "dateRequired"), c)
		return
	}
	if err := visitorService.AggregateDailySummary(date); err != nil {
		global.GVA_LOG.Error("聚合失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "fail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "success"), c)
}
