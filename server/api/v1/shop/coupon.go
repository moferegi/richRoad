package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CouponApi struct{}

// CreateCoupon 创建优惠券
// @Tags Coupon
// @Summary 创建优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shop.Coupon true "创建优惠券"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /Cou/createCoupon [post]
func (CouApi *CouponApi) CreateCoupon(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	// 创建业务用Context
	ctx := c.Request.Context()

	var Cou shop.Coupon
	err := c.ShouldBindJSON(&Cou)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	err = couponService.CreateCoupon(ctx, &Cou)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteCoupon 删除优惠券
// @Tags Coupon
// @Summary 删除优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shop.Coupon true "删除优惠券"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /Cou/deleteCoupon [delete]
func (CouApi *CouponApi) DeleteCoupon(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := couponService.DeleteCoupon(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteCouponByIds 批量删除优惠券
// @Tags Coupon
// @Summary 批量删除优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /Cou/deleteCouponByIds [delete]
func (CouApi *CouponApi) DeleteCouponByIds(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := couponService.DeleteCouponByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateCoupon 更新优惠券
// @Tags Coupon
// @Summary 更新优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shop.Coupon true "更新优惠券"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /Cou/updateCoupon [put]
func (CouApi *CouponApi) UpdateCoupon(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var Cou shop.Coupon
	err := c.ShouldBindJSON(&Cou)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	err = couponService.UpdateCoupon(ctx, Cou)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindCoupon 用id查询优惠券
// @Tags Coupon
// @Summary 用id查询优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询优惠券"
// @Success 200 {object} response.Response{data=shop.Coupon,msg=string} "查询成功"
// @Router /Cou/findCoupon [get]
func (CouApi *CouponApi) FindCoupon(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	reCou, err := couponService.GetCoupon(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(reCou, c)
}

// GetCouponList 分页获取优惠券列表
// @Tags Coupon
// @Summary 分页获取优惠券列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query shopReq.CouponSearch true "分页获取优惠券列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /Cou/getCouponList [get]
func (CouApi *CouponApi) GetCouponList(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo shopReq.CouponSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := couponService.GetCouponInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetCouponDataSource 获取Coupon的数据源
// @Tags Coupon
// @Summary 获取Coupon的数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /Cou/getCouponDataSource [get]
func (CouApi *CouponApi) GetCouponDataSource(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口为获取数据源定义的数据
	dataSource, err := couponService.GetCouponDataSource(ctx)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(dataSource, c)
}
