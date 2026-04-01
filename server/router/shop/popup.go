package shop

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PopupRouter struct{}

// InitPopupRouter 初始化 弹窗管理 路由信息
func (s *PopupRouter) InitPopupRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	popupRouter := Router.Group("popup").Use(middleware.OperationRecord())
	popupRouterWithoutAuth := PublicRouter.Group("popup")

	var popupApi = v1.ApiGroupApp.ShopApiGroup.PopupApi
	{
		popupRouter.POST("createPopup", popupApi.CreatePopup)       // 创建弹窗
		popupRouter.DELETE("deletePopup", popupApi.DeletePopup)     // 删除弹窗
		popupRouter.PUT("updatePopup", popupApi.UpdatePopup)        // 更新弹窗
		popupRouter.GET("getPopupList", popupApi.GetPopupList)      // 获取弹窗列表
	}
	{
		popupRouterWithoutAuth.GET("getActivePopups", popupApi.GetActivePopups) // 客户端获取生效弹窗
	}
}
