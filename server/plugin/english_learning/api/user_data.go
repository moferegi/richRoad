package api

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserDataApi struct{}

var userDataService = service.ServiceGroupApp.UserDataService

func normalizeUserDataPage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return page, pageSize
}

// SaveWordProgress 保存单词学习位置
// @Tags     EnglishUserData
// @Summary  保存单词学习位置
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.SaveWordProgressReq true "学习进度"
// @Success  200  {object} response.Response{msg=string} "保存成功"
// @Router   /englishLearning/userData/saveWordProgress [post]
func (a *UserDataApi) SaveWordProgress(c *gin.Context) {
	var req request.SaveWordProgressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}
	if claims.BaseClaims.ID == 0 {
		response.FailWithMessage("用户信息无效，请重新登录", c)
		return
	}
	if err = userDataService.SaveWordProgress(claims.BaseClaims.ID, req); err != nil {
		global.GVA_LOG.Error("保存学习进度失败", zap.Error(err))
		response.FailWithMessage("保存失败", c)
		return
	}
	response.OkWithMessage("保存成功", c)
}

// GetWordProgress 获取单词学习位置
// @Tags     EnglishUserData
// @Summary  获取单词学习位置
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=object,msg=string} "获取成功"
// @Router   /englishLearning/userData/getWordProgress [get]
func (a *UserDataApi) GetWordProgress(c *gin.Context) {
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}
	data, err := userDataService.GetWordProgress(claims.BaseClaims.ID)
	if err != nil {
		global.GVA_LOG.Error("获取学习进度失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(data, c)
}

// GetWatchProgress 获取视频观看位置
// @Tags     EnglishUserData
// @Summary  获取视频观看位置
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=object,msg=string} "获取成功"
// @Router   /englishLearning/userData/getWatchProgress [get]
func (a *UserDataApi) GetWatchProgress(c *gin.Context) {
	episodeIDStr := c.Query("episodeId")
	episodeID, err := strconv.ParseUint(episodeIDStr, 10, 64)
	if err != nil || episodeID == 0 {
		response.FailWithMessage("episodeId参数错误", c)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}
	data, err := userDataService.GetWatchProgress(claims.BaseClaims.ID, uint(episodeID))
	if err != nil {
		global.GVA_LOG.Error("获取视频进度失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(data, c)
}

// GetSeriesWatchProgressList 获取某剧集下的观看进度列表
// @Tags     EnglishUserData
// @Summary  获取某剧集下的观看进度列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=object,msg=string} "获取成功"
// @Router   /englishLearning/userData/getSeriesWatchProgressList [get]
func (a *UserDataApi) GetSeriesWatchProgressList(c *gin.Context) {
	seriesIDStr := c.Query("seriesId")
	seriesID, err := strconv.ParseUint(seriesIDStr, 10, 64)
	if err != nil || seriesID == 0 {
		response.FailWithMessage("seriesId参数错误", c)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}
	list, err := userDataService.GetSeriesWatchProgressList(claims.BaseClaims.ID, uint(seriesID))
	if err != nil {
		global.GVA_LOG.Error("获取剧集观看进度失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}

// GetWatchHistoryList 获取观看记录分页列表
// @Tags     EnglishUserData
// @Summary  获取当前用户观看记录分页列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.WatchHistorySearch true "分页参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/userData/getWatchHistoryList [get]
func (a *UserDataApi) GetWatchHistoryList(c *gin.Context) {
	var pageInfo request.WatchHistorySearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	page, pageSize := normalizeUserDataPage(pageInfo.Page, pageInfo.PageSize)

	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	list, total, err := userDataService.GetWatchHistoryList(claims.BaseClaims.ID, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取观看记录失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: page, PageSize: pageSize}, "获取成功", c)
}

// Collect 收藏
// @Tags     EnglishUserData
// @Summary  收藏单词/句子/视频
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.CollectionReq true "收藏参数"
// @Success  200  {object} response.Response{msg=string} "收藏成功"
// @Router   /englishLearning/userData/collect [post]
func (a *UserDataApi) Collect(c *gin.Context) {
	var req request.CollectionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}
	if err = userDataService.Collect(claims.BaseClaims.ID, req); err != nil {
		global.GVA_LOG.Error("收藏失败", zap.Error(err))
		response.FailWithMessage("收藏失败", c)
		return
	}
	response.OkWithMessage("收藏成功", c)
}

// Uncollect 取消收藏
// @Tags     EnglishUserData
// @Summary  取消收藏
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.CollectionReq true "收藏参数"
// @Success  200  {object} response.Response{msg=string} "取消成功"
// @Router   /englishLearning/userData/uncollect [delete]
func (a *UserDataApi) Uncollect(c *gin.Context) {
	var req request.CollectionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}
	if err = userDataService.Uncollect(claims.BaseClaims.ID, req); err != nil {
		global.GVA_LOG.Error("取消收藏失败", zap.Error(err))
		response.FailWithMessage("取消失败", c)
		return
	}
	response.OkWithMessage("取消成功", c)
}

// GetCollectionList 获取收藏列表
// @Tags     EnglishUserData
// @Summary  分页获取收藏列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.CollectionSearch true "分页参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/userData/getCollectionList [get]
func (a *UserDataApi) GetCollectionList(c *gin.Context) {
	var pageInfo request.CollectionSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}
	list, total, err := userDataService.GetCollectionList(claims.BaseClaims.ID, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取收藏列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize}, "获取成功", c)
}

// GetCollectionDetailList 获取收藏详情列表
// @Tags     EnglishUserData
// @Summary  分页获取收藏详情列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.CollectionSearch true "分页参数"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /englishLearning/userData/getCollectionDetailList [get]
func (a *UserDataApi) GetCollectionDetailList(c *gin.Context) {
	var pageInfo request.CollectionSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}
	list, total, err := userDataService.GetCollectionDetailList(claims.BaseClaims.ID, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取收藏详情列表失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize}, "获取成功", c)
}

// GetAsset 获取当前用户学习资产
// @Tags     EnglishUserData
// @Summary  获取当前用户学习资产
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=object,msg=string} "获取成功"
// @Router   /englishLearning/userData/getAsset [get]
func (a *UserDataApi) GetAsset(c *gin.Context) {
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}
	asset, err := userDataService.GetUserAsset(claims.BaseClaims.ID)
	if err != nil {
		global.GVA_LOG.Error("获取学习资产失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(asset, c)
}
