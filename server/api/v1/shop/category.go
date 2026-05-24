package shop

import (
	"strconv"

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

type CategoryApi struct {
}

var categoryService = service.ServiceGroupApp.ShopServiceGroup.CategoryService

type publicCategoryResponse struct {
	ID               uint                     `json:"ID"`
	ParentID         uint                     `json:"parentId"`
	Title            string                   `json:"title"`
	Desc             string                   `json:"desc"`
	Goods            []publicGoodResponse     `json:"goods,omitempty"`
	Children         []publicCategoryResponse `json:"children,omitempty"`
	Icons            string                   `json:"icons"`
	ExternalIconPath string                   `json:"externalIconPath"`
}

func toPublicCategoryResponse(item shop.Category) publicCategoryResponse {
	children := make([]publicCategoryResponse, 0, len(item.Children))
	for _, child := range item.Children {
		children = append(children, toPublicCategoryResponse(child))
	}

	return publicCategoryResponse{
		ID:               item.ID,
		ParentID:         item.ParentID,
		Title:            item.Title,
		Desc:             item.Desc,
		Goods:            toPublicGoodResponses(item.Goods),
		Children:         children,
		Icons:            item.Icons,
		ExternalIconPath: item.ExternalIconPath,
	}
}

func toPublicCategoryResponses(list []shop.Category) []publicCategoryResponse {
	result := make([]publicCategoryResponse, 0, len(list))
	for _, item := range list {
		result = append(result, toPublicCategoryResponse(item))
	}
	return result
}

// CreateCategory 创建商品分类
// @Tags Category
// @Summary 创建商品分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Category true "创建商品分类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /category/createCategory [post]
func (categoryApi *CategoryApi) CreateCategory(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	var category shop.Category
	err := c.ShouldBindJSON(&category)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := categoryService.CreateCategory(&category); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteCategory 删除商品分类
// @Tags Category
// @Summary 删除商品分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Category true "删除商品分类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /category/deleteCategory [delete]
func (categoryApi *CategoryApi) DeleteCategory(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	ID := c.Query("ID")
	if err := categoryService.DeleteCategory(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteCategoryByIds 批量删除商品分类
// @Tags Category
// @Summary 批量删除商品分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /category/deleteCategoryByIds [delete]
func (categoryApi *CategoryApi) DeleteCategoryByIds(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	IDs := c.QueryArray("IDs[]")
	if err := categoryService.DeleteCategoryByIds(IDs); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateCategory 更新商品分类
// @Tags Category
// @Summary 更新商品分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Category true "更新商品分类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /category/updateCategory [put]
func (categoryApi *CategoryApi) UpdateCategory(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	var category shop.Category
	err := c.ShouldBindJSON(&category)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := categoryService.UpdateCategory(category); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindCategory 用id查询商品分类
// @Tags Category
// @Summary 用id查询商品分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shop.Category true "用id查询商品分类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /category/findCategory [get]
func (categoryApi *CategoryApi) FindCategory(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	ID := c.Query("ID")
	if recategory, err := categoryService.GetCategory(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"recategory": recategory}, c)
	}
}

// GetCategoryList 分页获取商品分类列表
// @Tags Category
// @Summary 分页获取商品分类列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shopReq.CategorySearch true "分页获取商品分类列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /category/getCategoryList [get]
func (categoryApi *CategoryApi) GetCategoryList(c *gin.Context) {
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}
	var pageInfo shopReq.CategorySearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if list, total, err := categoryService.GetCategoryInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

func (categoryApi *CategoryApi) GetCategoryMobile(c *gin.Context) {
	var parentID int = 0
	qid := c.Query("parentID")
	if qid != "" {
		if id, err := strconv.Atoi(qid); err == nil {
			parentID = id
		}
	}
	if data, err := categoryService.GetCategoryMobile(parentID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(i18n.LocalizeResponseData(c, toPublicCategoryResponses(data)), "获取成功", c)
	}
}

func (categoryApi *CategoryApi) GetChildrenCategoryAndProduct(c *gin.Context) {
	var parentID int = 0
	qid := c.Query("parentID")
	if qid != "" {
		if id, err := strconv.Atoi(qid); err == nil {
			parentID = id
		}
	}
	if data, err := categoryService.GetChildrenCategoryAndProduct(parentID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(i18n.LocalizeResponseData(c, toPublicCategoryResponses(data)), "获取成功", c)
	}
}

// GetCategoryPublic 不需要鉴权的商品分类接口
// @Tags Category
// @Summary 不需要鉴权的商品分类接口
// @accept application/json
// @Produce application/json
// @Param data query shopReq.CategorySearch true "分页获取商品分类列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /category/getCategoryList [get]
func (categoryApi *CategoryApi) GetCategoryPublic(c *gin.Context) {
	var pageInfo shopReq.CategorySearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	showInUni := true
	pageInfo.ShowInUni = &showInUni

	if list, total, err := categoryService.GetCategoryInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     toPublicCategoryResponses(list),
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
