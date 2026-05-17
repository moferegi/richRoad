package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils" // 引入utils包
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CouponOrderUserApi struct{}

// CreateCouponOrderUser 创建优惠券
// @Tags CouponOrderUser
// @Summary 创建优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shop.CouponOrderUser true "创建优惠券"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /cou/createCouponOrderUser [post]
func (couApi *CouponOrderUserApi) CreateCouponOrderUser(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var cou shop.CouponOrderUser
	err := c.ShouldBindJSON(&cou)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = couService.CreateCouponOrderUser(ctx, &cou)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteCouponOrderUser 删除优惠券
// @Tags CouponOrderUser
// @Summary 删除优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shop.CouponOrderUser true "删除优惠券"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /cou/deleteCouponOrderUser [delete]
func (couApi *CouponOrderUserApi) DeleteCouponOrderUser(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := couService.DeleteCouponOrderUser(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteCouponOrderUserByIds 批量删除优惠券
// @Tags CouponOrderUser
// @Summary 批量删除优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /cou/deleteCouponOrderUserByIds [delete]
func (couApi *CouponOrderUserApi) DeleteCouponOrderUserByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := couService.DeleteCouponOrderUserByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateCouponOrderUser 更新优惠券
// @Tags CouponOrderUser
// @Summary 更新优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shop.CouponOrderUser true "更新优惠券"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /cou/updateCouponOrderUser [put]
func (couApi *CouponOrderUserApi) UpdateCouponOrderUser(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var cou shop.CouponOrderUser
	err := c.ShouldBindJSON(&cou)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = couService.UpdateCouponOrderUser(ctx, cou)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindCouponOrderUser 用id查询优惠券
// @Tags CouponOrderUser
// @Summary 用id查询优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询优惠券"
// @Success 200 {object} response.Response{data=shop.CouponOrderUser,msg=string} "查询成功"
// @Router /cou/findCouponOrderUser [get]
func (couApi *CouponOrderUserApi) FindCouponOrderUser(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	recou, err := couService.GetCouponOrderUser(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(recou, c)
}

// GetCouponOrderUserList 分页获取优惠券列表
// @Tags CouponOrderUser
// @Summary 分页获取优惠券列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query shopReq.CouponOrderUserSearch true "分页获取优惠券列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /cou/getCouponOrderUserList [get]
func (couApi *CouponOrderUserApi) GetCouponOrderUserList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo shopReq.CouponOrderUserSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := couService.GetCouponOrderUserInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetCouponOrderUserDataSource 获取CouponOrderUser的数据源
// @Tags CouponOrderUser
// @Summary 获取CouponOrderUser的数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /cou/getCouponOrderUserDataSource [get]
func (couApi *CouponOrderUserApi) GetCouponOrderUserDataSource(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口为获取数据源定义的数据
	dataSource, err := couService.GetCouponOrderUserDataSource(ctx)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(dataSource, c)
}

// GetCouponOrderUserPublic 不需要鉴权的优惠券接口
// @Tags CouponOrderUser
// @Summary 不需要鉴权的优惠券接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /cou/getCouponOrderUserPublic [get]
func (couApi *CouponOrderUserApi) GetCouponOrderUserPublic(c *gin.Context) {
	// ... existing code ...
}

// GetAllClaimCoupon 获取用户可领取的优惠券列表
// @Tags CouponOrderUser
// @Summary 获取用户可领取的优惠券列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shopReq.GetClaimCouponRequest true "商品ID列表"
// @Success 200 {object} response.Response{data=[]map[string]interface{},msg=string} "获取成功"
// @Router /cou/getAllClaimCoupon [post]
func (couApi *CouponOrderUserApi) GetAllClaimCoupon(c *gin.Context) {
	ctx := c.Request.Context()
	var req shopReq.GetClaimCouponRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("用户未登录", c)
		return
	}

	coupons, err := couService.GetAllClaimCoupon(ctx, userID, req.GoodIds)
	if err != nil {
		global.GVA_LOG.Error("获取优惠券列表失败!", zap.Error(err), zap.Uint("userID", userID))
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithData(i18n.LocalizeResponseData(c, coupons), c)
}

// ClaimCouponByUser 用户领取优惠券
// @Tags CouponOrderUser
// @Summary 用户领取优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shopReq.ClaimCouponRequest true "优惠券ID"
// @Success 200 {object} response.Response{msg=string} "领取成功"
// @Router /cou/claimCouponByUser [post]
func (couApi *CouponOrderUserApi) ClaimCouponByUser(c *gin.Context) {
	ctx := c.Request.Context()
	var req shopReq.ClaimCouponRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("用户未登录", c)
		return
	}

	couponNum, err := couService.ClaimCouponByUser(ctx, userID, req.CouponID)
	if err != nil {
		global.GVA_LOG.Error("用户领取优惠券失败!", zap.Error(err), zap.Uint("userID", userID), zap.Int("couponID", req.CouponID))
		response.FailWithMessage("领取失败: "+err.Error(), c)
		return
	}
	response.OkWithData(couponNum, c)
}

// AdminIssueCouponToAll 管理员向所有用户发放优惠券
// @Tags CouponOrderUser
// @Summary 管理员向所有用户发放优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shopReq.AdminIssueCouponRequest true "优惠券ID"
// @Success 200 {object} response.Response{msg=string} "发放成功"
// @Router /cou/adminIssueCouponToAll [post]
func (couApi *CouponOrderUserApi) AdminIssueCouponToAll(c *gin.Context) {
	ctx := c.Request.Context()
	var req shopReq.AdminIssueCouponRequest // 使用 AdminIssueCouponRequest，即使只用 CouponID，保持一致性或未来扩展
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = couService.IssueCouponToAllUsers(ctx, req.CouponID)
	if err != nil {
		global.GVA_LOG.Error("管理员发放优惠券失败!", zap.Error(err), zap.Int("couponID", req.CouponID))
		response.FailWithMessage("发放失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("发放成功", c)
}
