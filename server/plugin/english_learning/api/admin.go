package api

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminApi struct{}

var adminService = service.ServiceGroupApp.AdminService

// GetAdminCheckinRecordList 管理端：签到记录列表
// @Tags     EnglishAdmin
// @Summary  管理端签到记录列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    userId query uint false "用户ID"
// @Param    startDate query string false "开始日期 YYYY-MM-DD"
// @Param    endDate query string false "结束日期 YYYY-MM-DD"
// @Param    page query int false "页码"
// @Param    pageSize query int false "每页数量"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/admin/getCheckinRecordList [get]
func (a *AdminApi) GetCheckinRecordList(c *gin.Context) {
	userId := parseOptionalUserId(c)
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	page, pageSize := parsePage(c)

	list, total, err := adminService.GetAdminCheckinRecordList(userId, startDate, endDate, page, pageSize)
	if err != nil {
		global.GVA_LOG.Error("管理端获取签到记录列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: page, PageSize: pageSize}, "获取成功", c)
}

// GetAdminPointRecordList 管理端：积分记录列表
// @Tags     EnglishAdmin
// @Summary  管理端积分记录列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    userId query uint false "用户ID"
// @Param    changeType query string false "增减类型(increase/decrease)"
// @Param    startDate query string false "开始日期 YYYY-MM-DD"
// @Param    endDate query string false "结束日期 YYYY-MM-DD"
// @Param    page query int false "页码"
// @Param    pageSize query int false "每页数量"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/admin/getPointRecordList [get]
func (a *AdminApi) GetPointRecordList(c *gin.Context) {
	userId := parseOptionalUserId(c)
	changeType := c.Query("changeType")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	page, pageSize := parsePage(c)

	list, total, err := adminService.GetAdminPointRecordList(userId, changeType, startDate, endDate, page, pageSize)
	if err != nil {
		global.GVA_LOG.Error("管理端获取积分记录列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: page, PageSize: pageSize}, "获取成功", c)
}

// GetAdminFreeTimeRecordList 管理端：时长明细列表
// @Tags     EnglishAdmin
// @Summary  管理端时长明细列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    userId query uint false "用户ID"
// @Param    changeType query string false "增减类型(increase/decrease)"
// @Param    startDate query string false "开始日期 YYYY-MM-DD"
// @Param    endDate query string false "结束日期 YYYY-MM-DD"
// @Param    page query int false "页码"
// @Param    pageSize query int false "每页数量"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/admin/getFreeTimeRecordList [get]
func (a *AdminApi) GetFreeTimeRecordList(c *gin.Context) {
	userId := parseOptionalUserId(c)
	changeType := c.Query("changeType")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	page, pageSize := parsePage(c)

	list, total, err := adminService.GetAdminFreeTimeRecordList(userId, changeType, startDate, endDate, page, pageSize)
	if err != nil {
		global.GVA_LOG.Error("管理端获取时长明细列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: page, PageSize: pageSize}, "获取成功", c)
}

// GetAdminWatchHistoryList 管理端：观看历史列表
// @Tags     EnglishAdmin
// @Summary  管理端观看历史列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    userId query uint false "用户ID"
// @Param    startDate query string false "开始日期 YYYY-MM-DD"
// @Param    endDate query string false "结束日期 YYYY-MM-DD"
// @Param    page query int false "页码"
// @Param    pageSize query int false "每页数量"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/admin/getWatchHistoryList [get]
func (a *AdminApi) GetWatchHistoryList(c *gin.Context) {
	userId := parseOptionalUserId(c)
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	page, pageSize := parsePage(c)

	list, total, err := adminService.GetAdminWatchHistoryList(userId, startDate, endDate, page, pageSize)
	if err != nil {
		global.GVA_LOG.Error("管理端获取观看历史列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: page, PageSize: pageSize}, "获取成功", c)
}

// GetAdminCollectionList 管理端：收藏列表
// @Tags     EnglishAdmin
// @Summary  管理端收藏列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    userId query uint false "用户ID"
// @Param    targetType query int false "收藏类型 1单词 2句子 3视频"
// @Param    page query int false "页码"
// @Param    pageSize query int false "每页数量"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/admin/getCollectionList [get]
func (a *AdminApi) GetCollectionList(c *gin.Context) {
	userId := parseOptionalUserId(c)
	targetType, _ := strconv.Atoi(c.DefaultQuery("targetType", "0"))
	page, pageSize := parsePage(c)

	list, total, err := adminService.GetAdminCollectionList(userId, targetType, page, pageSize)
	if err != nil {
		global.GVA_LOG.Error("管理端获取收藏列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: page, PageSize: pageSize}, "获取成功", c)
}

// GetAdminWordErrorLogList 管理端：错词本列表
// @Tags     EnglishAdmin
// @Summary  管理端错词本列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    userId query uint false "用户ID"
// @Param    startDate query string false "开始日期 YYYY-MM-DD"
// @Param    endDate query string false "结束日期 YYYY-MM-DD"
// @Param    page query int false "页码"
// @Param    pageSize query int false "每页数量"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/admin/getWordErrorLogList [get]
func (a *AdminApi) GetWordErrorLogList(c *gin.Context) {
	userId := parseOptionalUserId(c)
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	page, pageSize := parsePage(c)

	list, total, err := adminService.GetAdminWordErrorLogList(userId, startDate, endDate, page, pageSize)
	if err != nil {
		global.GVA_LOG.Error("管理端获取错词本列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: page, PageSize: pageSize}, "获取成功", c)
}

// GetAdminUserList 管理端：用户列表（用于下拉筛选）
// @Tags     EnglishAdmin
// @Summary  管理端用户搜索列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    keyword query string false "搜索关键词(用户名/昵称/手机号)"
// @Success  200  {object} response.Response{data=[]service.AdminUserItem,msg=string} "获取成功"
// @Router   /englishLearning/admin/getUserList [get]
func (a *AdminApi) GetUserList(c *gin.Context) {
	keyword := c.Query("keyword")

	list, err := adminService.GetAdminUserList(keyword)
	if err != nil {
		global.GVA_LOG.Error("管理端获取用户列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}

// ---- helpers ----

func parseOptionalUserId(c *gin.Context) uint {
	idStr := c.Query("userId")
	if idStr == "" {
		return 0
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0
	}
	return uint(id)
}

func parsePage(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	return page, pageSize
}
