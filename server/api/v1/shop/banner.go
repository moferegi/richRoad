package shop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BannerApi struct {
}

var bannerService = service.ServiceGroupApp.ShopServiceGroup.BannerService

// CreateBanner 创建轮播图
// @Tags Banner
// @Summary 创建轮播图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Banner true "创建轮播图"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /banner/createBanner [post]
func (bannerApi *BannerApi) CreateBanner(c *gin.Context) {
	var banner shop.Banner
	err := c.ShouldBindJSON(&banner)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := bannerService.CreateBanner(&banner); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteBanner 删除轮播图
// @Tags Banner
// @Summary 删除轮播图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Banner true "删除轮播图"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /banner/deleteBanner [delete]
func (bannerApi *BannerApi) DeleteBanner(c *gin.Context) {
	ID := c.Query("ID")
	if err := bannerService.DeleteBanner(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteBannerByIds 批量删除轮播图
// @Tags Banner
// @Summary 批量删除轮播图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /banner/deleteBannerByIds [delete]
func (bannerApi *BannerApi) DeleteBannerByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	if err := bannerService.DeleteBannerByIds(IDs); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateBanner 更新轮播图
// @Tags Banner
// @Summary 更新轮播图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Banner true "更新轮播图"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /banner/updateBanner [put]
func (bannerApi *BannerApi) UpdateBanner(c *gin.Context) {
	var banner shop.Banner
	err := c.ShouldBindJSON(&banner)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := bannerService.UpdateBanner(banner); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindBanner 用id查询轮播图
// @Tags Banner
// @Summary 用id查询轮播图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shop.Banner true "用id查询轮播图"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /banner/findBanner [get]
func (bannerApi *BannerApi) FindBanner(c *gin.Context) {
	ID := c.Query("ID")
	if rebanner, err := bannerService.GetBanner(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rebanner": rebanner}, c)
	}
}

// GetBannerList 分页获取轮播图列表
// @Tags Banner
// @Summary 分页获取轮播图列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shopReq.BannerSearch true "分页获取轮播图列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /banner/getBannerList [get]
func (bannerApi *BannerApi) GetBannerList(c *gin.Context) {
	var pageInfo shopReq.BannerSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if list, total, err := bannerService.GetBannerInfoList(pageInfo); err != nil {
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

// GetBannerPublic 不需要鉴权的轮播图接口
// @Tags Banner
// @Summary 不需要鉴权的轮播图接口
// @accept application/json
// @Produce application/json
// @Param data query shopReq.BannerSearch true "分页获取轮播图列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /banner/getBannerList [get]
func (bannerApi *BannerApi) GetBannerPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	response.OkWithDetailed(i18n.LocalizeResponseData(c, gin.H{
		"info": "不需要鉴权的轮播图接口信息",
	}), "获取成功", c)
}
