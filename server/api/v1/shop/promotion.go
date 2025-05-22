package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PromotionApi struct{}

// CreatePromotion 创建促销信息
// @Tags Promotion
// @Summary 创建促销信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shop.Promotion true "创建促销信息"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /Promo/createPromotion [post]
func (PromoApi *PromotionApi) CreatePromotion(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var Promo shop.Promotion
	err := c.ShouldBindJSON(&Promo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = PromoService.CreatePromotion(ctx, &Promo)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeletePromotion 删除促销信息
// @Tags Promotion
// @Summary 删除促销信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shop.Promotion true "删除促销信息"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /Promo/deletePromotion [delete]
func (PromoApi *PromotionApi) DeletePromotion(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := PromoService.DeletePromotion(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeletePromotionByIds 批量删除促销信息
// @Tags Promotion
// @Summary 批量删除促销信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /Promo/deletePromotionByIds [delete]
func (PromoApi *PromotionApi) DeletePromotionByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := PromoService.DeletePromotionByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdatePromotion 更新促销信息
// @Tags Promotion
// @Summary 更新促销信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body shop.Promotion true "更新促销信息"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /Promo/updatePromotion [put]
func (PromoApi *PromotionApi) UpdatePromotion(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var Promo shop.Promotion
	err := c.ShouldBindJSON(&Promo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = PromoService.UpdatePromotion(ctx, Promo)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindPromotion 用id查询促销信息
// @Tags Promotion
// @Summary 用id查询促销信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询促销信息"
// @Success 200 {object} response.Response{data=shop.Promotion,msg=string} "查询成功"
// @Router /Promo/findPromotion [get]
func (PromoApi *PromotionApi) FindPromotion(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	rePromo, err := PromoService.GetPromotion(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(rePromo, c)
}

// GetPromotionList 分页获取促销信息列表
// @Tags Promotion
// @Summary 分页获取促销信息列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query shopReq.PromotionSearch true "分页获取促销信息列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /Promo/getPromotionList [get]
func (PromoApi *PromotionApi) GetPromotionList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo shopReq.PromotionSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := PromoService.GetPromotionInfoList(ctx, pageInfo)
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

// GetPromotionPublic 不需要鉴权的促销信息接口
// @Tags Promotion
// @Summary 不需要鉴权的促销信息接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /Promo/getPromotionPublic [get]
/*func (PromoApi *PromotionApi) GetPromotionPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    PromoService.GetPromotionPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的促销信息接口信息",
    }, "获取成功", c)
}
*/
// GetPromotionPublic 不需要鉴权的促销信息接口
// @Tags Promotion
// @Summary 获取第一个开启状态的促销信息
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=shop.Promotion,msg=string} "获取成功"
// @Router /shop/promotion/getPromotionPublic [get]
func (PromoApi *PromotionApi) GetPromotionPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 调用service层获取促销信息
	promotion, err := PromoService.GetPromotionPublic(ctx)
	if err != nil {
		global.GVA_LOG.Error("获取促销信息失败!", zap.Error(err))
		response.FailWithMessage("获取促销信息失败", c)
		return
	}

	response.OkWithDetailed(promotion, "获取成功", c)
}
