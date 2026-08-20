package client

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GameApi struct{}

var gameService = service.ServiceGroupApp.ClientServiceGroup.GameService

// ==================== Uni 端接口 ====================

// GetGameCategory 获取单个游戏分类
// @Tags Game
// @Summary 获取单个游戏分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param gameID query int true "游戏ID"
// @Success 200 {object} response.Response{data=client.GameCategory,msg=string} "获取成功"
// @Router /game/getGameCategory [get]
func (api *GameApi) GetGameCategory(c *gin.Context) {
	gameIDStr := c.Query("gameID")
	gameID, err := strconv.ParseUint(gameIDStr, 10, 64)
	if err != nil || gameID == 0 {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	cat, err := gameService.GetGameCategoryByID(uint(gameID))
	if err != nil {
		global.GVA_LOG.Error("获取游戏分类失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(cat, i18n.T(c, "getSuccess"), c)
}

// GetDifficultyCategory 获取单个难度分类
// @Tags Game
// @Summary 获取单个难度分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param catID query int true "难度分类ID"
// @Success 200 {object} response.Response{data=client.GameDifficultyCategory,msg=string} "获取成功"
// @Router /game/getDifficultyCategory [get]
func (api *GameApi) GetDifficultyCategory(c *gin.Context) {
	catIDStr := c.Query("catID")
	catID, err := strconv.ParseUint(catIDStr, 10, 64)
	if err != nil || catID == 0 {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	cat, err := gameService.GetDifficultyCategoryByID(uint(catID))
	if err != nil {
		global.GVA_LOG.Error("获取难度分类失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(cat, i18n.T(c, "getSuccess"), c)
}

// GetGameList 获取游戏列表
// @Tags Game
// @Summary 获取游戏列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]client.GameCategory,msg=string} "获取成功"
// @Router /game/getGameList [get]
func (api *GameApi) GetGameList(c *gin.Context) {
	list, err := gameService.GetGameCategoryList()
	if err != nil {
		global.GVA_LOG.Error("获取游戏列表失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(list, i18n.T(c, "getSuccess"), c)
}

// SetUserProgress 管理端手动设置用户进度
// @Tags GameAdmin
// @Summary 管理端手动设置用户进度
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.SetUserProgressReq true "用户ID, 关卡ID, 状态"
// @Success 200 {object} response.Response{msg=string} "操作成功"
// @Router /game/admin/setUserProgress [post]
func (api *GameApi) SetUserProgress(c *gin.Context) {
	var req clientModel.SetUserProgressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if err := gameService.SetUserProgress(req.UserID, req.LevelID, req.Status); err != nil {
		global.GVA_LOG.Error("设置用户进度失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "updateFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
}

// GetDifficultyCategories 获取难度分类列表
// @Tags Game
// @Summary 获取难度分类列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param gameID query int true "游戏ID"
// @Success 200 {object} response.Response{data=[]client.GameDifficultyCategory,msg=string} "获取成功"
// @Router /game/getDifficultyCategories [get]
func (api *GameApi) GetDifficultyCategories(c *gin.Context) {
	gameIDStr := c.Query("gameID")
	gameID, err := strconv.ParseUint(gameIDStr, 10, 64)
	if err != nil || gameID == 0 {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	list, err := gameService.GetDifficultyCategoryList(uint(gameID))
	if err != nil {
		global.GVA_LOG.Error("获取难度分类失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(list, i18n.T(c, "getSuccess"), c)
}

// GetLevelDetail 获取关卡详情（含数字，用于开始游戏）
// @Tags Game
// @Summary 获取关卡详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param levelID query int true "关卡ID"
// @Success 200 {object} response.Response{data=client.GameLevel,msg=string} "获取成功"
// @Router /game/getLevelDetail [get]
func (api *GameApi) GetLevelDetail(c *gin.Context) {
	levelIDStr := c.Query("levelID")
	levelID, err := strconv.ParseUint(levelIDStr, 10, 64)
	if err != nil || levelID == 0 {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	level, err := gameService.GetLevelByID(uint(levelID))
	if err != nil {
		global.GVA_LOG.Error("获取关卡详情失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(level, i18n.T(c, "getSuccess"), c)
}

// GetLevelListByCategory 获取关卡列表（Uni端）
// @Tags Game
// @Summary 获取关卡列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param categoryID query int true "难度分类ID"
// @Success 200 {object} response.Response{data=[]client.GameLevel,msg=string} "获取成功"
// @Router /game/getLevelList [get]
func (api *GameApi) GetLevelListByCategory(c *gin.Context) {
	categoryIDStr := c.Query("categoryID")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil || categoryID == 0 {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	list, err := gameService.GetLevelList(uint(categoryID))
	if err != nil {
		global.GVA_LOG.Error("获取关卡列表失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(list, i18n.T(c, "getSuccess"), c)
}

// SubmitLevelResult 提交闯关结果
// @Tags Game
// @Summary 提交闯关结果
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param levelID query int true "关卡ID"
// @Success 200 {object} response.Response{msg=string} "闯关成功"
// @Router /game/submitLevelResult [post]
func (api *GameApi) SubmitLevelResult(c *gin.Context) {
	userID := utils.GetUserID(c)
	levelIDStr := c.Query("levelID")
	levelID, err := strconv.ParseUint(levelIDStr, 10, 64)
	if err != nil || levelID == 0 {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if err := gameService.SubmitLevelResult(userID, uint(levelID)); err != nil {
		global.GVA_LOG.Error("提交闯关结果失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "operationFailed"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "operationSuccess"), c)
}

// GetUserProgress 获取用户闯关进度
// @Tags Game
// @Summary 获取用户闯关进度
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param categoryID query int true "难度分类ID"
// @Success 200 {object} response.Response{data=[]client.GameUserProgress,msg=string} "获取成功"
// @Router /game/getUserProgress [get]
func (api *GameApi) GetUserProgress(c *gin.Context) {
	userID := utils.GetUserID(c)
	categoryIDStr := c.Query("categoryID")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil || categoryID == 0 {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	list, err := gameService.GetUserProgress(userID, uint(categoryID))
	if err != nil {
		global.GVA_LOG.Error("获取用户进度失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(list, i18n.T(c, "getSuccess"), c)
}

// GetLeaderboard 获取排行榜
// @Tags Game
// @Summary 获取排行榜
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param type query string false "类型: all=累计 daily=每日"
// @Param limit query int false "数量限制"
// @Success 200 {object} response.Response{data=[]object,msg=string} "获取成功"
// @Router /game/getLeaderboard [get]
func (api *GameApi) GetLeaderboard(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)
	boardType := c.DefaultQuery("type", "all")

	var list interface{}
	var err error
	if boardType == "daily" {
		list, err = gameService.GetDailyLeaderboard(limit)
	} else {
		list, err = gameService.GetAllTimeLeaderboard(limit)
	}
	if err != nil {
		global.GVA_LOG.Error("获取排行榜失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(list, i18n.T(c, "getSuccess"), c)
}

// ==================== Web 管理端接口 ====================

// CreateCategory 创建游戏大分类
// @Tags GameAdmin
// @Summary 创建游戏大分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.GameCategory true "游戏大分类"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /game/admin/createCategory [post]
func (api *GameApi) CreateCategory(c *gin.Context) {
	var cat clientModel.GameCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if err := gameService.CreateGameCategory(&cat); err != nil {
		global.GVA_LOG.Error("创建游戏分类失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "createFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "createSuccess"), c)
}

// UpdateCategory 更新游戏大分类
// @Tags GameAdmin
// @Summary 更新游戏大分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.GameCategory true "游戏大分类"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /game/admin/updateCategory [put]
func (api *GameApi) UpdateCategory(c *gin.Context) {
	var cat clientModel.GameCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if err := gameService.UpdateGameCategory(&cat); err != nil {
		global.GVA_LOG.Error("更新游戏分类失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "updateFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
}

// DeleteCategory 删除游戏大分类
// @Tags GameAdmin
// @Summary 删除游戏大分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query int true "分类ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /game/admin/deleteCategory [delete]
func (api *GameApi) DeleteCategory(c *gin.Context) {
	idStr := c.Query("ID")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if err := gameService.DeleteGameCategory(uint(id)); err != nil {
		global.GVA_LOG.Error("删除游戏分类失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "deleteFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "deleteSuccess"), c)
}

// GetCategoryList 获取游戏大分类列表（管理端）
// @Tags GameAdmin
// @Summary 获取游戏大分类列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /game/admin/getCategoryList [get]
func (api *GameApi) GetCategoryList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	list, total, err := gameService.GetGameCategoryListAdmin(page, pageSize)
	if err != nil {
		global.GVA_LOG.Error("获取游戏分类列表失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: page, PageSize: pageSize,
	}, i18n.T(c, "getSuccess"), c)
}

// CreateDifficultyCategory 创建难度分类
// @Tags GameAdmin
// @Summary 创建难度分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.GameDifficultyCategory true "难度分类"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /game/admin/createDifficultyCategory [post]
func (api *GameApi) CreateDifficultyCategory(c *gin.Context) {
	var cat clientModel.GameDifficultyCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if err := gameService.CreateDifficultyCategory(&cat); err != nil {
		global.GVA_LOG.Error("创建难度分类失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "createFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "createSuccess"), c)
}

// UpdateDifficultyCategory 更新难度分类
// @Tags GameAdmin
// @Summary 更新难度分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.GameDifficultyCategory true "难度分类"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /game/admin/updateDifficultyCategory [put]
func (api *GameApi) UpdateDifficultyCategory(c *gin.Context) {
	var cat clientModel.GameDifficultyCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if err := gameService.UpdateDifficultyCategory(&cat); err != nil {
		global.GVA_LOG.Error("更新难度分类失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "updateFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
}

// DeleteDifficultyCategory 删除难度分类
// @Tags GameAdmin
// @Summary 删除难度分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query int true "难度分类ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /game/admin/deleteDifficultyCategory [delete]
func (api *GameApi) DeleteDifficultyCategory(c *gin.Context) {
	idStr := c.Query("ID")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if err := gameService.DeleteDifficultyCategory(uint(id)); err != nil {
		global.GVA_LOG.Error("删除难度分类失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "deleteFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "deleteSuccess"), c)
}

// GetDifficultyCategoryList 获取难度分类列表（管理端）
// @Tags GameAdmin
// @Summary 获取难度分类列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param gameID query int true "游戏ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /game/admin/getDifficultyCategoryList [get]
func (api *GameApi) GetDifficultyCategoryList(c *gin.Context) {
	gameIDStr := c.Query("gameID")
	gameID, err := strconv.ParseUint(gameIDStr, 10, 64)
	if err != nil || gameID == 0 {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	list, total, err := gameService.GetDifficultyCategoryListAdmin(uint(gameID), page, pageSize)
	if err != nil {
		global.GVA_LOG.Error("获取难度分类列表失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: page, PageSize: pageSize,
	}, i18n.T(c, "getSuccess"), c)
}

// CreateLevel 创建关卡
// @Tags GameAdmin
// @Summary 创建关卡
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.GameLevel true "关卡"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /game/admin/createLevel [post]
func (api *GameApi) CreateLevel(c *gin.Context) {
	var level clientModel.GameLevel
	if err := c.ShouldBindJSON(&level); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if err := gameService.CreateLevel(&level); err != nil {
		global.GVA_LOG.Error("创建关卡失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "createFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "createSuccess"), c)
}

// UpdateLevel 更新关卡
// @Tags GameAdmin
// @Summary 更新关卡
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.GameLevel true "关卡"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /game/admin/updateLevel [put]
func (api *GameApi) UpdateLevel(c *gin.Context) {
	var level clientModel.GameLevel
	if err := c.ShouldBindJSON(&level); err != nil {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if err := gameService.UpdateLevel(&level); err != nil {
		global.GVA_LOG.Error("更新关卡失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "updateFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
}

// DeleteLevel 删除关卡
// @Tags GameAdmin
// @Summary 删除关卡
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query int true "关卡ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /game/admin/deleteLevel [delete]
func (api *GameApi) DeleteLevel(c *gin.Context) {
	idStr := c.Query("ID")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	if err := gameService.DeleteLevel(uint(id)); err != nil {
		global.GVA_LOG.Error("删除关卡失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "deleteFail"), c)
		return
	}
	response.OkWithMessage(i18n.T(c, "deleteSuccess"), c)
}

// GetLevelList 获取关卡列表（管理端）
// @Tags GameAdmin
// @Summary 获取关卡列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param categoryID query int true "难度分类ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /game/admin/getLevelList [get]
func (api *GameApi) GetLevelList(c *gin.Context) {
	categoryIDStr := c.Query("categoryID")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil || categoryID == 0 {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	list, total, err := gameService.GetLevelListAdmin(uint(categoryID), page, pageSize)
	if err != nil {
		global.GVA_LOG.Error("获取关卡列表失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: page, PageSize: pageSize,
	}, i18n.T(c, "getSuccess"), c)
}

// GetLeaderboardAdmin 管理端排行榜
// @Tags GameAdmin
// @Summary 管理端获取排行榜
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param type query string false "all/daily"
// @Param limit query int false "数量"
// @Success 200 {object} response.Response{data=[]object,msg=string} "获取成功"
// @Router /game/admin/getLeaderboard [get]
func (api *GameApi) GetLeaderboardAdmin(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)
	boardType := c.DefaultQuery("type", "all")

	var list interface{}
	var err error
	if boardType == "daily" {
		list, err = gameService.GetDailyLeaderboard(limit)
	} else {
		list, err = gameService.GetAllTimeLeaderboard(limit)
	}
	if err != nil {
		global.GVA_LOG.Error("获取排行榜失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(list, i18n.T(c, "getSuccess"), c)
}

// GetUserProgressAdmin 管理端查看用户进度
// @Tags GameAdmin
// @Summary 管理端查看用户进度
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param userID query int true "用户ID"
// @Param categoryID query int true "难度分类ID"
// @Success 200 {object} response.Response{data=[]client.GameUserProgress,msg=string} "获取成功"
// @Router /game/admin/getUserProgress [get]
func (api *GameApi) GetUserProgressAdmin(c *gin.Context) {
	userIDStr := c.Query("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	categoryIDStr := c.Query("categoryID")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil || categoryID == 0 {
		response.FailWithMessage(i18n.T(c, "invalidParams"), c)
		return
	}
	list, err := gameService.GetUserProgress(uint(userID), uint(categoryID))
	if err != nil {
		global.GVA_LOG.Error("获取用户进度失败!", zap.Error(err))
		response.FailWithMessage(i18n.T(c, "getFail"), c)
		return
	}
	response.OkWithDetailed(list, i18n.T(c, "getSuccess"), c)
}
