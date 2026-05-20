package shop

import (
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

type CartApi struct {
}

var cartService = service.ServiceGroupApp.ShopServiceGroup.CartService

// CreateCart 创建购物车
// @Tags Cart
// @Summary 创建购物车
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Cart true "创建购物车"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /cart/createCart [post]
func (cartApi *CartApi) CreateCart(c *gin.Context) {
	var cart shop.Cart
	err := c.ShouldBindJSON(&cart)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := cartService.CreateCart(&cart); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteCart 删除购物车
// @Tags Cart
// @Summary 删除购物车
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Cart true "删除购物车"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /cart/deleteCart [delete]
func (cartApi *CartApi) DeleteCart(c *gin.Context) {
	ID := c.Query("ID")
	if err := cartService.DeleteCart(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteCartByIds 批量删除购物车
// @Tags Cart
// @Summary 批量删除购物车
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /cart/deleteCartByIds [delete]
func (cartApi *CartApi) DeleteCartByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	if err := cartService.DeleteCartByIds(IDs); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateCart 更新购物车
// @Tags Cart
// @Summary 更新购物车
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Cart true "更新购物车"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /cart/updateCart [put]
func (cartApi *CartApi) UpdateCart(c *gin.Context) {
	var cart shop.Cart
	err := c.ShouldBindJSON(&cart)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := cartService.UpdateCart(cart); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindCart 用id查询购物车
// @Tags Cart
// @Summary 用id查询购物车
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shop.Cart true "用id查询购物车"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /cart/findCart [get]
func (cartApi *CartApi) FindCart(c *gin.Context) {
	ID := c.Query("ID")
	if recart, err := cartService.GetCart(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"recart": recart}, c)
	}
}

// GetCartList 分页获取购物车列表
// @Tags Cart
// @Summary 分页获取购物车列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shopReq.CartSearch true "分页获取购物车列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cart/getCartList [get]
func (cartApi *CartApi) GetCartList(c *gin.Context) {
	var pageInfo shopReq.CartSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("获取用户失败", c)
		return
	}

	if list, total, err := cartService.GetCartInfoList(pageInfo); err != nil {
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

// GetSelfCart 获取用户自身购物车
// @Tags Cart
// @Summary 获取用户自身购物车
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cart/getSelfCart [get]
func (cartApi *CartApi) GetSelfCart(c *gin.Context) {
	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("获取用户失败", c)
		return
	}

	if list, err := cartService.GetSelfCart(userID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(i18n.LocalizeResponseData(c, list), "获取成功", c)
	}
}

// AddCart 添加购物车
// @Tags Cart
// @Summary 添加购物车
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Cart true "添加购物车"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cart/addCart [post]
func (cartApi *CartApi) AddCart(c *gin.Context) {
	var cart shopReq.CartCreate
	err := c.ShouldBindJSON(&cart)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	userID := utils.GetUserID(c)
	cart.UserID = userID
	if err := cartService.AddCart(&cart); err != nil {
		global.GVA_LOG.Error("购物车操作失败!", zap.Error(err))
		response.FailWithMessage("购物车操作失败", c)
	} else {
		response.OkWithMessage("购物车操作成功", c)
	}
}

// CutCart 减少购物车
// @Tags Cart
// @Summary 减少购物车
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Cart true "减少购物车"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cart/cutCart [post]
func (cartApi *CartApi) CutCart(c *gin.Context) {
	var cart shopReq.CartCreate
	err := c.ShouldBindJSON(&cart)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	userID := utils.GetUserID(c)
	cart.UserID = userID
	if err := cartService.CutCart(&cart); err != nil {
		global.GVA_LOG.Error("购物车操作失败!", zap.Error(err))
		response.FailWithMessage("购物车操作失败", c)
	} else {
		response.OkWithMessage("购物车操作成功", c)
	}
}

// ClearCart 清空购物车
// @Tags Cart
// @Summary 清空购物车
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cart/clearCart [get]
func (cartApi *CartApi) ClearCart(c *gin.Context) {
	userID := utils.GetUserID(c)
	if err := cartService.ClearCart(userID); err != nil {
		global.GVA_LOG.Error("购物车操作失败!", zap.Error(err))
		response.FailWithMessage("购物车操作失败", c)
	} else {
		response.OkWithMessage("购物车操作成功", c)
	}
}

// GetCartPublic 不需要鉴权的购物车接口
// @Tags Cart
// @Summary 不需要鉴权的购物车接口
// @accept application/json
// @Produce application/json
// @Param data query shopReq.CartSearch true "分页获取购物车列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cart/getCartList [get]
//func (cartApi *CartApi) GetCartPublic(c *gin.Context) {
//	// 此接口不需要鉴权
//	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
//	response.OkWithDetailed(gin.H{
//		"info": "不需要鉴权的购物车接口信息",
//	}, "获取成功", c)
//}
