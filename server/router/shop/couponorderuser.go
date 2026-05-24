package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CouponOrderUserRouter struct{}

// InitCouponOrderUserRouter 初始化 优惠券 路由信息
func (s *CouponOrderUserRouter) InitCouponOrderUserRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	couRouter := Router.Group("cou").Use(middleware.OperationRecord())
	couRouterWithoutRecord := Router.Group("cou")
	couRouterWithoutAuth := PublicRouter.Group("cou")
	{
		couRouter.POST("createCouponOrderUser", couApi.CreateCouponOrderUser)             // 新建优惠券
		couRouter.DELETE("deleteCouponOrderUser", couApi.DeleteCouponOrderUser)           // 删除优惠券
		couRouter.DELETE("deleteCouponOrderUserByIds", couApi.DeleteCouponOrderUserByIds) // 批量删除优惠券
		couRouter.PUT("updateCouponOrderUser", couApi.UpdateCouponOrderUser)              // 更新优惠券
	}
	{
		couRouterWithoutRecord.GET("findCouponOrderUser", couApi.FindCouponOrderUser)                   // 根据ID获取优惠券
		couRouterWithoutRecord.GET("getCouponOrderUserList", couApi.GetCouponOrderUserList)             // 获取优惠券列表
		couRouterWithoutRecord.GET("getCouponOrderUserDataSource", couApi.GetCouponOrderUserDataSource) // 获取优惠券数据源
	}
	{
		couRouterWithoutAuth.GET("getCouponOrderUserPublic", couApi.GetCouponOrderUserPublic) // 优惠券开放接口
	}
	{
		couRouter.POST("getAllClaimCoupon", couApi.GetAllClaimCoupon)         // 用户查询可领可用优惠券
		couRouter.POST("claimCouponByUser", couApi.ClaimCouponByUser)         // 用户领取优惠券
		couRouter.POST("adminIssueCouponToAll", couApi.AdminIssueCouponToAll) // 管理员向所有用户发放优惠券 (需要管理员权限中间件)
	}
}
