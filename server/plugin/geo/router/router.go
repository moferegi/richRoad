package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/geo/api"
	"github.com/gin-gonic/gin"
)

type GeoRouter struct {
}

func (s *GeoRouter) InitGeoRouter(Router *gin.RouterGroup) {
	plugApi := api.ApiGroupApp.GeoApi
	publicRouter := Router.Group("")
	authRouter := Router.Group("").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		publicRouter.GET("getGeos", plugApi.GetGeos)
		authRouter.GET("getGeo", plugApi.GetGeo)
		authRouter.PUT("editGeo", plugApi.EditGeo)
		authRouter.POST("createGeo", plugApi.CreateGeo)
		authRouter.DELETE("deleteGeo", plugApi.DeleteGeo)
	}
}
