package english_learning

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	systemMiddleware "github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/initialize"
	learningMiddleware "github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/middleware"
	englishRouter "github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/router"
	interfaces "github.com/flipped-aurora/gin-vue-admin/server/utils/plugin/v2"
	"github.com/gin-gonic/gin"
)

var _ interfaces.Plugin = (*plugin)(nil)

var Plugin = new(plugin)

type plugin struct{}

func init() {
	interfaces.Register(Plugin)
}

func (p *plugin) Register(engine *gin.Engine) {
	ctx := context.Background()
	initialize.Gorm(ctx)

	group := engine.Group(global.GVA_CONFIG.System.RouterPrefix).Group("englishLearning")
	group.Use(systemMiddleware.Locale(), systemMiddleware.JWTAuth(), learningMiddleware.LearningAuth())
	group.Use(learningMiddleware.LearningReadRateLimit())
	group.Use(learningMiddleware.LearningRequestSignGuard())
	group.Use(learningMiddleware.LearningResponseEncrypt())

	englishRouter.RouterGroupApp.InitUserLearningAssetRouter(group)
	englishRouter.RouterGroupApp.InitEnglishWordRouter(group)
	englishRouter.RouterGroupApp.InitVideoSubtitleRouter(group)
	englishRouter.RouterGroupApp.InitCheckinRouter(group)
	englishRouter.RouterGroupApp.InitContentRouter(group)
	englishRouter.RouterGroupApp.InitUserDataRouter(group)
}
