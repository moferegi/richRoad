package shop

import (
	"strconv"
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
)

type PresaleApi struct{}

var presaleService = service.ServiceGroupApp.ShopServiceGroup.PresaleService

type publicPresaleGoodResponse struct {
	ID                uint       `json:"ID"`
	Title             string     `json:"title"`
	ImageUrl          string     `json:"imageUrl"`
	ExternalImagePath string     `json:"externalImagePath"`
	Price             *float64   `json:"price"`
	PriceI18n         string     `json:"priceI18n"`
	PresaleQty        *int       `json:"presaleQty"`
	PresaleSold       *int       `json:"presaleSold"`
	PresaleStart      *time.Time `json:"presaleStart"`
	PresaleEnd        *time.Time `json:"presaleEnd"`
}

func toPublicPresaleGoodResponses(list []shop.Good) []publicPresaleGoodResponse {
	result := make([]publicPresaleGoodResponse, 0, len(list))
	for _, item := range list {
		result = append(result, publicPresaleGoodResponse{
			ID:                item.ID,
			Title:             item.Title,
			ImageUrl:          item.ImageUrl,
			ExternalImagePath: item.ExternalImagePath,
			Price:             item.Price,
			PriceI18n:         item.PriceI18n,
			PresaleQty:        item.PresaleQty,
			PresaleSold:       item.PresaleSold,
			PresaleStart:      item.PresaleStart,
			PresaleEnd:        item.PresaleEnd,
		})
	}
	return result
}

// GetPresaleGoodList 获取预售商品列表（客户端）
// @Tags Presale
// @Summary 获取预售商品列表
// @accept application/json
// @Produce application/json
// @Param data query shopReq.PresaleListRequest true "分页获取预售商品列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /presale/getPresaleGoodList [get]
func (api *PresaleApi) GetPresaleGoodList(c *gin.Context) {
	var pageInfo shopReq.PresaleListRequest
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if list, total, err := presaleService.GetPresaleGoodList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取预售列表失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		pageResult := response.PageResult{
			List:     toPublicPresaleGoodResponses(list),
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}
		response.OkWithDetailed(i18n.LocalizeResponseData(c, pageResult), "获取成功", c)
	}
}

// CheckPresaleAvailable 检查预售商品是否可购买
// @Tags Presale
// @Summary 检查预售商品是否可购买
// @accept application/json
// @Produce application/json
// @Param goodId query string true "商品ID"
// @Param quantity query int true "购买数量"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /presale/checkAvailable [get]
func (api *PresaleApi) CheckPresaleAvailable(c *gin.Context) {
	goodId := c.Query("goodId")
	gid, _ := strconv.ParseUint(goodId, 10, 64)
	available, msg, err := presaleService.CheckPresaleAvailable(uint(gid))
	if err != nil {
		global.GVA_LOG.Error("检查预售可用性失败!", zap.Error(err))
		response.FailWithMessage("检查失败", c)
		return
	}
	response.OkWithDetailed(i18n.LocalizeResponseData(c, gin.H{
		"available": available,
		"message":   msg,
	}), "获取成功", c)
}

// GetPresaleParticipants 获取预售参与者（管理端）
// @Tags Presale
// @Summary 获取预售参与者列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param goodId query string true "商品ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /presale/getPresaleParticipants [get]
func (api *PresaleApi) GetPresaleParticipants(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	goodId := c.Query("goodID")
	gid, _ := strconv.ParseUint(goodId, 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if list, total, err := presaleService.GetPresaleParticipants(uint(gid), page, pageSize); err != nil {
		global.GVA_LOG.Error("获取预售参与者失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		}, "获取成功", c)
	}
}
