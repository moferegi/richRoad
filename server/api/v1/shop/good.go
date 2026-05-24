package shop

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/datatypes"
)

type GoodApi struct {
}

var goodService = service.ServiceGroupApp.ShopServiceGroup.GoodService

type publicSkuResponse struct {
	ID                  uint           `json:"ID"`
	Name                string         `json:"name"`
	Picture             string         `json:"picture"`
	ExternalPicturePath string         `json:"externalPicturePath"`
	UpperImage          string         `json:"upperImage"`
	LowerImage          string         `json:"lowerImage"`
	Description         string         `json:"description"`
	Price               uint           `json:"price"`
	PriceI18n           string         `json:"priceI18n"`
	Inventory           uint           `json:"inventory"`
	Specs               datatypes.JSON `json:"specs"`
	Attrs               datatypes.JSON `json:"attrs"`
	GoodID              uint           `json:"goodID"`
	SaleNum             uint           `json:"saleNum"`
}

type publicGoodResponse struct {
	ID                  uint                `json:"ID"`
	Description         string              `json:"description"`
	Tags                datatypes.JSON      `json:"tags"`
	Specs               datatypes.JSON      `json:"specs"`
	Attrs               datatypes.JSON      `json:"attrs"`
	ImageUrl            string              `json:"imageUrl"`
	UpperImage          string              `json:"upperImage"`
	LowerImage          string              `json:"lowerImage"`
	Banner              datatypes.JSON      `json:"banner"`
	Price               *float64            `json:"price"`
	PriceI18n           string              `json:"priceI18n"`
	Rating              *float64            `json:"rating"`
	ReviewCount         *int                `json:"reviewCount"`
	SaleCount           *int                `json:"saleCount"`
	Title               string              `json:"title"`
	CategoryID          *int                `json:"categoryID"`
	Postage             *float64            `json:"postage"`
	Discount            *int                `json:"discount"`
	SKUS                []publicSkuResponse `json:"skus"`
	Detail              string              `json:"detail"`
	SaleNum             uint                `json:"saleNum"`
	ViewNum             int                 `json:"view_num"`
	ExternalImagePath   string              `json:"externalImagePath"`
	PointsEnabled       *bool               `json:"pointsEnabled"`
	PointsMaxUse        *int                `json:"pointsMaxUse"`
	PointsUseTimes      *int                `json:"pointsUseTimes"`
	IsPresale           *bool               `json:"isPresale"`
	PresaleQty          *int                `json:"presaleQty"`
	PresaleSold         *int                `json:"presaleSold"`
	PresaleStart        *time.Time          `json:"presaleStart"`
	PresaleEnd          *time.Time          `json:"presaleEnd"`
	PresaleEnabled      *bool               `json:"presaleEnabled"`
	PresalePopupEnabled *bool               `json:"presalePopupEnabled"`
	PresalePopupTitle   string              `json:"presalePopupTitle"`
	PresalePopupContent string              `json:"presalePopupContent"`
	CouponAvailable     bool                `json:"couponAvailable,omitempty"`
}

func toPublicSkuResponse(item shop.Sku) publicSkuResponse {
	return publicSkuResponse{
		ID:                  item.ID,
		Name:                item.Name,
		Picture:             item.Picture,
		ExternalPicturePath: item.ExternalPicturePath,
		UpperImage:          item.UpperImage,
		LowerImage:          item.LowerImage,
		Description:         item.Description,
		Price:               item.Price,
		PriceI18n:           item.PriceI18n,
		Inventory:           item.Inventory,
		Specs:               item.Specs,
		Attrs:               item.Attrs,
		GoodID:              item.GoodID,
		SaleNum:             item.SaleNum,
	}
}

func toPublicSkuResponses(list []shop.Sku) []publicSkuResponse {
	result := make([]publicSkuResponse, 0, len(list))
	for _, item := range list {
		result = append(result, toPublicSkuResponse(item))
	}
	return result
}

func toPublicGoodResponse(item shop.Good) publicGoodResponse {
	return publicGoodResponse{
		ID:                  item.ID,
		Description:         item.Description,
		Tags:                item.Tags,
		Specs:               item.Specs,
		Attrs:               item.Attrs,
		ImageUrl:            item.ImageUrl,
		UpperImage:          item.UpperImage,
		LowerImage:          item.LowerImage,
		Banner:              item.Banner,
		Price:               item.Price,
		PriceI18n:           item.PriceI18n,
		Rating:              item.Rating,
		ReviewCount:         item.ReviewCount,
		SaleCount:           item.SaleCount,
		Title:               item.Title,
		CategoryID:          item.CategoryID,
		Postage:             item.Postage,
		Discount:            item.Discount,
		SKUS:                toPublicSkuResponses(item.SKUS),
		Detail:              item.Detail,
		SaleNum:             item.SaleNum,
		ViewNum:             item.ViewNum,
		ExternalImagePath:   item.ExternalImagePath,
		PointsEnabled:       item.PointsEnabled,
		PointsMaxUse:        item.PointsMaxUse,
		PointsUseTimes:      item.PointsUseTimes,
		IsPresale:           item.IsPresale,
		PresaleQty:          item.PresaleQty,
		PresaleSold:         item.PresaleSold,
		PresaleStart:        item.PresaleStart,
		PresaleEnd:          item.PresaleEnd,
		PresaleEnabled:      item.PresaleEnabled,
		PresalePopupEnabled: item.PresalePopupEnabled,
		PresalePopupTitle:   item.PresalePopupTitle,
		PresalePopupContent: item.PresalePopupContent,
	}
}

func toPublicGoodResponses(list []shop.Good) []publicGoodResponse {
	result := make([]publicGoodResponse, 0, len(list))
	for _, item := range list {
		result = append(result, toPublicGoodResponse(item))
	}
	return result
}

// CreateGood 创建商品
// @Tags Good
// @Summary 创建商品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Good true "创建商品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /good/createGood [post]
func (goodApi *GoodApi) CreateGood(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	var good shop.Good
	err := c.ShouldBindJSON(&good)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := goodService.CreateGood(&good); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteGood 删除商品
// @Tags Good
// @Summary 删除商品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Good true "删除商品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /good/deleteGood [delete]
func (goodApi *GoodApi) DeleteGood(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	ID := c.Query("ID")
	if err := goodService.DeleteGood(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteGoodByIds 批量删除商品
// @Tags Good
// @Summary 批量删除商品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /good/deleteGoodByIds [delete]
func (goodApi *GoodApi) DeleteGoodByIds(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	IDs := c.QueryArray("IDs[]")
	if err := goodService.DeleteGoodByIds(IDs); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateGood 更新商品
// @Tags Good
// @Summary 更新商品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Good true "更新商品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /good/updateGood [put]
func (goodApi *GoodApi) UpdateGood(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	var good shop.Good
	err := c.ShouldBindJSON(&good)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := goodService.UpdateGood(good); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindGood 用id查询商品
// @Tags Good
// @Summary 用id查询商品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shop.Good true "用id查询商品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /good/findGood [get]
func (goodApi *GoodApi) FindGood(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	authorityID := utils.GetUserAuthorityId(c)
	if regood, err := goodService.GetGood(ID, userID, authorityID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(i18n.LocalizeResponseData(c, gin.H{"regood": regood}), c)
	}
}

// GetGoodHistory 查询用户的商品浏览历史
// @Tags Good
// @Summary 查询用户的商品浏览历史
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取历史记录成功"}"
// @Router /good/getGoodHistory [get]
func (goodApi *GoodApi) GetGoodHistory(c *gin.Context) {
	userID := utils.GetUserID(c)
	if goods, err := goodService.GetGoodHistory(userID); err != nil {
		global.GVA_LOG.Error("获取历史记录失败!", zap.Error(err))
		response.FailWithMessage("获取历史记录失败", c)
	} else {
		response.OkWithData(i18n.LocalizeResponseData(c, toPublicGoodResponses(goods)), c)
	}
}

// ClearGoodHistory 清空用户的商品浏览历史
// @Tags Good
// @Summary 清空用户的商品浏览历史
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{\"success\":true,\"data\":{},\"msg\":\"清空成功\"}"
// @Router /good/clearGoodHistory [delete]
func (goodApi *GoodApi) ClearGoodHistory(c *gin.Context) {
	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("用户未登录", c)
		return
	}

	if err := goodService.ClearGoodHistory(userID); err != nil {
		global.GVA_LOG.Error("清空历史记录失败!", zap.Error(err))
		response.FailWithMessage("清空历史记录失败", c)
	} else {
		response.OkWithMessage("清空成功", c)
	}
}

// GetGoodList 分页获取商品列表
// @Tags Good
// @Summary 分页获取商品列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shopReq.GoodSearch true "分页获取商品列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /good/getGoodList [get]
func (goodApi *GoodApi) GetGoodList(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	var pageInfo shopReq.GoodSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if list, total, err := goodService.GetGoodInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		pageResult := response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}
		response.OkWithDetailed(i18n.LocalizeResponseData(c, pageResult), "获取成功", c)
	}
}

// GetGoodPublic 不需要鉴权的商品接口
// @Tags Good
// @Summary 不需要鉴权的商品接口
// @accept application/json
// @Produce application/json
// @Param data query shopReq.GoodSearch true "分页获取商品列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /good/getGoodList [get]
func (goodApi *GoodApi) GetGoodPublic(c *gin.Context) {
	if ID := c.Query("ID"); ID != "" {
		if regood, err := goodService.GetPublicGood(ID); err != nil {
			global.GVA_LOG.Error("查询失败!", zap.Error(err))
			response.FailWithMessage("查询失败", c)
		} else {
			response.OkWithData(i18n.LocalizeResponseData(c, gin.H{"regood": toPublicGoodResponse(regood)}), c)
		}
		return
	}

	var pageInfo shopReq.GoodSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	active := true
	pageInfo.Status = &active
	pageInfo.ExcludeHiddenCategories = true

	if list, total, err := goodService.GetGoodInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		pageResult := response.PageResult{
			List:     toPublicGoodResponses(list),
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}
		response.OkWithDetailed(i18n.LocalizeResponseData(c, pageResult), "获取成功", c)
	}
}
