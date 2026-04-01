package client

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type LanguageRouter struct{}

// InitLanguageRouter 初始化 语言管理 路由信息
func (s *LanguageRouter) InitLanguageRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	languageRouter := Router.Group("language").Use(middleware.OperationRecord())
	languageRouterWithoutAuth := PublicRouter.Group("language")

	var languageApi = v1.ApiGroupApp.ClientApiGroup.LanguageApi
	{
		languageRouter.POST("createLanguage", languageApi.CreateLanguage)       // 创建语言
		languageRouter.DELETE("deleteLanguage", languageApi.DeleteLanguage)     // 删除语言
		languageRouter.PUT("updateLanguage", languageApi.UpdateLanguage)        // 更新语言
		languageRouter.GET("getLanguageList", languageApi.GetLanguageList)      // 获取语言列表（管理端）
	}
	{
		languageRouterWithoutAuth.GET("getEnabledLanguages", languageApi.GetEnabledLanguages) // 客户端获取启用的语言
	}
}
