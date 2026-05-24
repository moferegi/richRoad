package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	shopModel "github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CollectApi struct {
}

var collectService = service.ServiceGroupApp.ClientServiceGroup.CollectService

type publicCollectGoodResponse struct {
	ID                uint     `json:"ID"`
	ImageUrl          string   `json:"imageUrl"`
	ExternalImagePath string   `json:"externalImagePath"`
	Title             string   `json:"title"`
	Price             *float64 `json:"price"`
	PriceI18n         string   `json:"priceI18n"`
	Discount          *int     `json:"discount"`
	SaleNum           uint     `json:"saleNum"`
}

func toPublicCollectGoodResponse(item shopModel.Good) publicCollectGoodResponse {
	return publicCollectGoodResponse{
		ID:                item.ID,
		ImageUrl:          item.ImageUrl,
		ExternalImagePath: item.ExternalImagePath,
		Title:             item.Title,
		Price:             item.Price,
		PriceI18n:         item.PriceI18n,
		Discount:          item.Discount,
		SaleNum:           item.SaleNum,
	}
}

func toPublicCollectGoodResponses(list []shopModel.Good) []publicCollectGoodResponse {
	result := make([]publicCollectGoodResponse, 0, len(list))
	for _, item := range list {
		result = append(result, toPublicCollectGoodResponse(item))
	}
	return result
}

func isCollectAdmin(authorityID uint) bool {
	return authorityID == 888 || authorityID == 8881
}

// CreateCollect 创建收藏
// @Tags Collect
// @Summary 创建收藏
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.Collect true "创建收藏"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /collect/createCollect [post]
func (collectApi *CollectApi) CreateCollect(c *gin.Context) {
	var collect client.Collect
	err := c.ShouldBindJSON(&collect)
	if err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	collect.UserID = utils.GetUserID(c)
	if err := collectService.CreateCollect(&collect); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		failClientWithErr(c, err)
	} else {
		response.OkWithMessage(i18n.T(c, "success"), c)
	}
}

// DeleteCollect 删除收藏
// @Tags Collect
// @Summary 删除收藏
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.Collect true "删除收藏"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /collect/deleteCollect [delete]
func (collectApi *CollectApi) DeleteCollect(c *gin.Context) {
	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	allowAll := isCollectAdmin(utils.GetUserAuthorityId(c))
	if err := collectService.DeleteCollect(ID, userID, allowAll); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		failClientWithKey(c, "deleteFail")
	} else {
		response.OkWithMessage(i18n.T(c, "deleteSuccess"), c)
	}
}

// DeleteCollectByIds 批量删除收藏
// @Tags Collect
// @Summary 批量删除收藏
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /collect/deleteCollectByIds [delete]
func (collectApi *CollectApi) DeleteCollectByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	allowAll := isCollectAdmin(utils.GetUserAuthorityId(c))
	if err := collectService.DeleteCollectByIds(IDs, userID, allowAll); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		failClientWithKey(c, "batchDeleteFail")
	} else {
		response.OkWithMessage(i18n.T(c, "batchDeleteSuccess"), c)
	}
}

// UpdateCollect 更新收藏
// @Tags Collect
// @Summary 更新收藏
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body client.Collect true "更新收藏"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /collect/updateCollect [put]
func (collectApi *CollectApi) UpdateCollect(c *gin.Context) {
	var collect client.Collect
	err := c.ShouldBindJSON(&collect)
	if err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}

	userID := utils.GetUserID(c)
	allowAll := isCollectAdmin(utils.GetUserAuthorityId(c))
	if err := collectService.UpdateCollect(collect, userID, allowAll); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		failClientWithKey(c, "updateFail")
	} else {
		response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
	}
}

// FindCollect 用id查询收藏
// @Tags Collect
// @Summary 用id查询收藏
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query client.Collect true "用id查询收藏"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /collect/findCollect [get]
func (collectApi *CollectApi) FindCollect(c *gin.Context) {
	ID := c.Query("goodID")
	userID := utils.GetUserID(c)
	if ok, err := collectService.GetCollect(userID, ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		failClientWithKey(c, "queryFail")
	} else {
		response.OkWithData(i18n.LocalizeResponseData(c, ok), c)
	}
}

// GetCollectList 分页获取收藏列表
// @Tags Collect
// @Summary 分页获取收藏列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query clientReq.CollectSearch true "分页获取收藏列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /collect/getCollectList [get]
func (collectApi *CollectApi) GetCollectList(c *gin.Context) {
	var pageInfo clientReq.CollectSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		failClientWithKey(c, "invalidParams")
		return
	}
	pageInfo.UserID = utils.GetUserID(c)
	if list, total, err := collectService.GetCollectInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		failClientWithKey(c, "getFail")
	} else {
		pageResult := response.PageResult{
			List:     toPublicCollectGoodResponses(list),
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}
		response.OkWithDetailed(i18n.LocalizeResponseData(c, pageResult), i18n.T(c, "getSuccess"), c)
	}
}

// GetCollectPublic 不需要鉴权的收藏接口
// @Tags Collect
// @Summary 不需要鉴权的收藏接口
// @accept application/json
// @Produce application/json
// @Param data query clientReq.CollectSearch true "分页获取收藏列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /collect/getCollectPublic [get]
func (collectApi *CollectApi) GetCollectPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的收藏接口信息",
	}, i18n.T(c, "getSuccess"), c)
}
