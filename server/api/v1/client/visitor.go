package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
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
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.VisitorID == "" {
		response.FailWithMessage("visitorId is required", c)
		return
	}

	// 获取真实IP
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// 尝试获取登录用户ID（公开接口，可能未登录）
	userID := utils.GetUserID(c)

	log := &client.VisitorLog{
		VisitorID:    req.VisitorID,
		UserID:       userID,
		IP:           ip,
		UserAgent:    userAgent,
		Platform:     req.Platform,
		PagePath:     req.PagePath,
		Referer:      req.Referer,
		ScreenWidth:  req.ScreenWidth,
		ScreenHeight: req.ScreenHeight,
		Language:     req.Language,
		SessionID:    req.SessionID,
	}

	if err := visitorService.Heartbeat(log); err != nil {
		global.GVA_LOG.Error("访客心跳上报失败!", zap.Error(err))
		response.FailWithMessage("上报失败", c)
		return
	}
	response.OkWithMessage("ok", c)
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
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := visitorService.GetVisitorLogList(pageInfo)
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
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := visitorService.GetVisitorSummaryList(pageInfo)
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
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(stats, "获取成功", c)
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
		response.FailWithMessage("date is required", c)
		return
	}
	if err := visitorService.AggregateDailySummary(date); err != nil {
		global.GVA_LOG.Error("聚合失败!", zap.Error(err))
		response.FailWithMessage("聚合失败", c)
		return
	}
	response.OkWithMessage("聚合成功", c)
}
