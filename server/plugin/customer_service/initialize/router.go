package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/router"
	"github.com/gin-gonic/gin"
)

// Router 注册客服系统路由
func Router(engine *gin.Engine) {
	// public: 无需权限（WebSocket 在此分组，自行验证 token）
	public := engine.Group(global.GVA_CONFIG.System.RouterPrefix).Group("")
	// private: 需要 JWT + Casbin
	private := engine.Group(global.GVA_CONFIG.System.RouterPrefix).Group("")
	private.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())

	router.Router.CustomerServiceRouter.InitCustomerServiceRouter(public, private)
}
