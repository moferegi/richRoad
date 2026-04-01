package shop

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type MarketingRewardRouter struct{}

// InitMarketingRewardRouter 初始化 营销奖励 路由信息
func (s *MarketingRewardRouter) InitMarketingRewardRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	marketingRewardRouter := Router.Group("marketingReward").Use(middleware.OperationRecord())
	marketingRewardRouterWithoutAuth := PublicRouter.Group("marketingReward")

	var marketingRewardApi = v1.ApiGroupApp.ShopApiGroup.MarketingRewardApi
	{
		marketingRewardRouter.POST("createMarketingReward", marketingRewardApi.CreateMarketingReward)       // 创建营销奖励
		marketingRewardRouter.DELETE("deleteMarketingReward", marketingRewardApi.DeleteMarketingReward)     // 删除营销奖励
		marketingRewardRouter.PUT("updateMarketingReward", marketingRewardApi.UpdateMarketingReward)        // 更新营销奖励
		marketingRewardRouter.GET("getMarketingRewardList", marketingRewardApi.GetMarketingRewardList)      // 获取营销奖励列表
	}
	{
		marketingRewardRouterWithoutAuth.GET("getMarketingRewardByType", marketingRewardApi.GetMarketingRewardByType) // 客户端根据类型获取奖励规则
	}
}
