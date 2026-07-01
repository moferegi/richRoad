package middleware

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	CtxEnglishTrialOnly   = "english_trial_only"
	CtxEnglishHasFullAuth = "english_has_full_auth"
	CtxEnglishUserID      = "english_user_id"
)

func setLearningAuthContext(c *gin.Context, userID uint, hasFullAuth bool) {
	c.Set(CtxEnglishUserID, userID)
	c.Set(CtxEnglishHasFullAuth, hasFullAuth)
	c.Set(CtxEnglishTrialOnly, !hasFullAuth)
}

// LearningAuth 英语学习体系全局鉴权中间件
// 该中间件负责拦截判断用户的资产情况，以决定是下发完整权限，还是打上 8% 试看标记
func LearningAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authzSvc := service.ServiceGroupApp.LearningAuthzService

		// 1. 获取 JWT 声明中的系统 userID
		claims, err := utils.GetClaims(c)
		if err != nil || claims.BaseClaims.ID == 0 {
			setLearningAuthContext(c, 0, false)
			c.Next()
			return
		}

		userID := claims.BaseClaims.ID

		hasFullAuth, _, resolveErr := authzSvc.ResolveGlobalFullAuth(userID)
		if resolveErr != nil {
			global.GVA_LOG.Warn("学习权限上下文构建失败，降级为试看态", zap.Error(resolveErr), zap.Uint("userID", userID))
			hasFullAuth = false
		}

		setLearningAuthContext(c, userID, hasFullAuth)

		c.Next()
	}
}
