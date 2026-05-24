package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PromotionApi struct{}

type publicPromotionResponse struct {
	PromotionImage string  `json:"promotionImage"`
	Title          *string `json:"title"`
	Description    *string `json:"description"`
}

func toPublicPromotionResponse(item shop.Promotion) publicPromotionResponse {
	return publicPromotionResponse{
		PromotionImage: item.PromotionImage,
		Title:          item.Title,
		Description:    item.Description,
	}
}

// CreatePromotion 创建促销信息
// @Tags Promotion
// @Summary 创建促销信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=publicPromotionResponse,msg=string} "获取成功"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /promo/createPromotion [post]
func (promoApi *PromotionApi) CreatePromotion(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	// 创建业务用Context
	ctx := c.Request.Context()

	var Promo shop.Promotion
	err := c.ShouldBindJSON(&Promo)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	err = promoService.CreatePromotion(ctx, &Promo)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
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
// @Router /promo/deletePromotion [delete]
func (promoApi *PromotionApi) DeletePromotion(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := promoService.DeletePromotion(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
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
// @Router /promo/deletePromotionByIds [delete]
func (promoApi *PromotionApi) DeletePromotionByIds(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := promoService.DeletePromotionByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
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
// @Router /promo/updatePromotion [put]
func (promoApi *PromotionApi) UpdatePromotion(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var Promo shop.Promotion
	err := c.ShouldBindJSON(&Promo)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	err = promoService.UpdatePromotion(ctx, Promo)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
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
// @Router /promo/findPromotion [get]
func (promoApi *PromotionApi) FindPromotion(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	rePromo, err := promoService.GetPromotion(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
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
// @Router /promo/getPromotionList [get]
func (promoApi *PromotionApi) GetPromotionList(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo shopReq.PromotionSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := promoService.GetPromotionInfoList(ctx, pageInfo)
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

// GetPromotionPublic 不需要鉴权的促销信息接口
// @Tags Promotion
// @Summary 不需要鉴权的促销信息接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /promo/getPromotionPublic [get]
/*func (promoApi *PromotionApi) GetPromotionPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    promoService.GetPromotionPublic(ctx)
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
// @Success 200 {object} response.Response{data=publicPromotionResponse,msg=string} "获取成功"
// @Router /shop/promotion/getPromotionPublic [get]
func (promoApi *PromotionApi) GetPromotionPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 调用service层获取促销信息
	promotion, err := promoService.GetPromotionPublic(ctx)
	if err != nil {
		global.GVA_LOG.Error("获取促销信息失败!", zap.Error(err))
		response.FailWithMessage("获取促销信息失败", c)
		return
	}

	response.OkWithDetailed(i18n.LocalizeResponseData(c, toPublicPromotionResponse(promotion)), "获取成功", c)
}
