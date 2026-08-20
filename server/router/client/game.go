package client

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type GameRouter struct{}

// InitGameRouter 初始化游戏路由
func (g *GameRouter) InitGameRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	gameRouter := Router.Group("game")
	gameAdminRouter := PublicRouter.Group("game/admin")

	var gameApi = v1.ApiGroupApp.ClientApiGroup.GameApi
	{
		// Uni 端接口
		gameRouter.GET("getGameList", gameApi.GetGameList)
		gameRouter.GET("getGameCategory", gameApi.GetGameCategory)
		gameRouter.GET("getDifficultyCategories", gameApi.GetDifficultyCategories)
		gameRouter.GET("getDifficultyCategory", gameApi.GetDifficultyCategory)
		gameRouter.GET("getLevelDetail", gameApi.GetLevelDetail)
		gameRouter.GET("getLevelList", gameApi.GetLevelListByCategory)
		gameRouter.POST("submitLevelResult", gameApi.SubmitLevelResult)
		gameRouter.GET("getUserProgress", gameApi.GetUserProgress)
		gameRouter.GET("getLeaderboard", gameApi.GetLeaderboard)

		// Web 管理端接口
		gameAdminRouter.POST("createCategory", gameApi.CreateCategory)
		gameAdminRouter.PUT("updateCategory", gameApi.UpdateCategory)
		gameAdminRouter.DELETE("deleteCategory", gameApi.DeleteCategory)
		gameAdminRouter.GET("getCategoryList", gameApi.GetCategoryList)
		gameAdminRouter.POST("createDifficultyCategory", gameApi.CreateDifficultyCategory)
		gameAdminRouter.PUT("updateDifficultyCategory", gameApi.UpdateDifficultyCategory)
		gameAdminRouter.DELETE("deleteDifficultyCategory", gameApi.DeleteDifficultyCategory)
		gameAdminRouter.GET("getDifficultyCategoryList", gameApi.GetDifficultyCategoryList)
		gameAdminRouter.POST("createLevel", gameApi.CreateLevel)
		gameAdminRouter.PUT("updateLevel", gameApi.UpdateLevel)
		gameAdminRouter.DELETE("deleteLevel", gameApi.DeleteLevel)
		gameAdminRouter.GET("getLevelList", gameApi.GetLevelList)
		gameAdminRouter.GET("getLeaderboard", gameApi.GetLeaderboardAdmin)
		gameAdminRouter.GET("getUserProgress", gameApi.GetUserProgressAdmin)
		gameAdminRouter.POST("setUserProgress", gameApi.SetUserProgress)
	}
}
