package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CouponRouter struct{}

// InitCouponRouter 初始化 优惠券 路由信息
func (s *CouponRouter) InitCouponRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	CouRouter := Router.Group("Cou").Use(middleware.OperationRecord())
	CouRouterWithoutRecord := Router.Group("Cou")
	{
		CouRouter.POST("createCoupon", CouApi.CreateCoupon)             // 新建优惠券
		CouRouter.DELETE("deleteCoupon", CouApi.DeleteCoupon)           // 删除优惠券
		CouRouter.DELETE("deleteCouponByIds", CouApi.DeleteCouponByIds) // 批量删除优惠券
		CouRouter.PUT("updateCoupon", CouApi.UpdateCoupon)              // 更新优惠券
	}
	{
		CouRouterWithoutRecord.GET("findCoupon", CouApi.FindCoupon)                   // 根据ID获取优惠券
		CouRouterWithoutRecord.GET("getCouponList", CouApi.GetCouponList)             // 获取优惠券列表
		CouRouterWithoutRecord.GET("getCouponDataSource", CouApi.GetCouponDataSource) // 获取优惠券数据源
	}
}
