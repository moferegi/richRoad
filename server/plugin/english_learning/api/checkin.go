package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/service"
	clientService "github.com/flipped-aurora/gin-vue-admin/server/service/client"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CheckinApi struct{}

// 获取单例Service
var checkinService = service.ServiceGroupApp.CheckinService

func normalizeCheckinPage(page int, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return page, pageSize
}

// DoCheckin
// @Tags     Checkin
// @Summary  执行每日打卡，按周期自动叠加积分
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{msg=string} "打卡成功，返回当次奖励"
// @Router   /englishLearning/checkin/do [post]
func (a *CheckinApi) DoCheckin(c *gin.Context) {
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	svc := clientService.SysConfigService{}
	basePoint := svc.GetConfigIntByKey("learning_checkin_base_point", 100)
	increment := svc.GetConfigIntByKey("learning_checkin_increment", 5)
	cycleDays := svc.GetConfigIntByKey("learning_checkin_cycle_days", 10)

	award, err := checkinService.PerformCheckin(claims.BaseClaims.ID, basePoint, increment, cycleDays)
	if err != nil {
		global.GVA_LOG.Error("打卡失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{"award": award}, "打卡成功", c)
}

// ExchangeTime
// @Tags     Checkin
// @Summary  积分兑换免费观看时长
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.ExchangePointsReq true "所需兑换消耗的积分"
// @Success  200  {object} response.Response{msg=string} "兑换成功"
// @Router   /englishLearning/checkin/exchange [post]
func (a *CheckinApi) ExchangeTime(c *gin.Context) {
	var req request.ExchangePointsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	svc := clientService.SysConfigService{}
	exchangeRate := svc.GetConfigIntByKey("points_exchange_rate", 100)
	if exchangeRate <= 0 {
		exchangeRate = 100
	}

	mins, err := checkinService.ExchangePoints(claims.BaseClaims.ID, req.Points, exchangeRate)
	if err != nil {
		global.GVA_LOG.Error("兑换失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{"minutes": mins}, "兑换成功", c)
}

// GetStats
// @Tags     Checkin
// @Summary  获取签到统计
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=object,msg=string} "获取成功"
// @Router   /englishLearning/checkin/getStats [get]
func (a *CheckinApi) GetStats(c *gin.Context) {
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	stats, checkedToday, err := checkinService.GetStats(claims.BaseClaims.ID)
	if err != nil {
		global.GVA_LOG.Error("获取签到统计失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	svc := clientService.SysConfigService{}
	dailyTarget := svc.GetConfigIntByKey("learning_daily_target", 10)
	basePoint := svc.GetConfigIntByKey("learning_checkin_base_point", 100)
	increment := svc.GetConfigIntByKey("learning_checkin_increment", 5)
	cycleDays := svc.GetConfigIntByKey("learning_checkin_cycle_days", 10)

	response.OkWithDetailed(gin.H{
		"continuousDays": stats.ContinuousDays,
		"totalDays":      stats.TotalDays,
		"wordsToday":     stats.WordsToday,
		"checkedToday":   checkedToday,
		"dailyTarget":    dailyTarget,
		"basePoint":      basePoint,
		"increment":      increment,
		"cycleDays":      cycleDays,
	}, "获取成功", c)
}

// GetPointRecordList
// @Tags     Checkin
// @Summary  获取英语学习积分流水列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.PointRecordSearch true "分页参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/checkin/getPointRecordList [get]
func (a *CheckinApi) GetPointRecordList(c *gin.Context) {
	var query request.PointRecordSearch
	if err := c.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	list, total, page, pageSize, err := checkinService.GetPointRecordList(claims.BaseClaims.ID, query)
	if err != nil {
		global.GVA_LOG.Error("获取英语学习积分流水失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	page, pageSize = normalizeCheckinPage(page, pageSize)
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, "获取成功", c)
}

// GetCheckinRecordList
// @Tags     Checkin
// @Summary  获取英语学习签到记录列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.CheckinRecordSearch true "分页参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/checkin/getCheckinRecordList [get]
func (a *CheckinApi) GetCheckinRecordList(c *gin.Context) {
	var query request.CheckinRecordSearch
	if err := c.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	list, total, page, pageSize, err := checkinService.GetCheckinRecordList(claims.BaseClaims.ID, query)
	if err != nil {
		global.GVA_LOG.Error("获取英语学习签到记录失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	page, pageSize = normalizeCheckinPage(page, pageSize)
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, "获取成功", c)
}
