package shop

import (
	"encoding/json"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/shop"
	shopReq "github.com/flipped-aurora/gin-vue-admin/server/model/shop/request"
	"github.com/flipped-aurora/gin-vue-admin/server/passport/sf"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrderApi struct {
}

var orderService = service.ServiceGroupApp.ShopServiceGroup.OrderService

// CreateOrder 创建订单
// @Tags Order
// @Summary 创建订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Order true "创建订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /order/createOrder [post]
func (orderApi *OrderApi) CreateOrder(c *gin.Context) {
	var order shop.Order
	err := c.ShouldBindJSON(&order)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := orderService.CreateOrder(&order); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// PlaceOrder 下单
// @Tags Order
// @Summary 下单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Order true "下单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"下单成功"}"
// @Router /order/placeOrder [post]
func (orderApi *OrderApi) PlaceOrder(c *gin.Context) {
	var order shop.Order
	err := c.ShouldBindJSON(&order)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	order.UserID = utils.GetUserID(c)

	if orderID, orderNo, err := orderService.PlaceOrder(&order); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(i18n.LocalizeResponseData(c, gin.H{
			"orderID": orderID,
			"orderNo": orderNo,
		}), c)
	}
}

// PlaceOrderByCart 购物车下单
// @Tags Order
// @Summary 购物车下单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"购物车下单"}"
// @Router /order/placeOrderByCart [post]
func (orderApi *OrderApi) PlaceOrderByCart(c *gin.Context) {
	userID := utils.GetUserID(c)
	var req shopReq.PlaceOrderByCartRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if orderID, orderNo, err := orderService.PlaceOrderByCart(userID, req); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(i18n.LocalizeResponseData(c, gin.H{
			"orderID": orderID,
			"orderNo": orderNo,
		}), c)
	}
}

func (orderApi *OrderApi) ChangeOrderCoupon(c *gin.Context) {

	orderID := c.Query("orderID")
	couponNum := c.Query("couponNum")
	userID := utils.GetUserID(c)

	if err := orderService.ChangeOrderCoupon(userID, orderID, couponNum); err != nil {
		global.GVA_LOG.Error("变更失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.Ok(c)
	}
}

// ChangeOrderPoints 变更订单积分抵扣
// @Tags Order
// @Summary 变更订单积分抵扣
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param orderID query string true "订单ID"
// @Param usePoints query bool true "是否使用积分"
// @Success 200 {object} response.Response{msg=string} "变更成功"
// @Router /order/changeOrderPoints [post]
func (orderApi *OrderApi) ChangeOrderPoints(c *gin.Context) {
	orderID := c.Query("orderID")
	usePointsStr := c.Query("usePoints")
	userID := utils.GetUserID(c)

	usePoints := usePointsStr == "true"

	if err := orderService.ChangeOrderPoints(userID, orderID, usePoints); err != nil {
		global.GVA_LOG.Error("积分变更失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.Ok(c)
	}
}

// ApplyRefund 申请退款
// @Tags Order
// @Summary 申请退款
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shopReq.RefundApplyReq true "退款申请"
// @Success 200 {object} response.Response{msg=string} "申请成功"
// @Router /order/applyRefund [post]
func (orderApi *OrderApi) ApplyRefund(c *gin.Context) {
	userID := utils.GetUserID(c)
	var req shopReq.RefundApplyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := orderService.ApplyRefund(userID, req); err != nil {
		global.GVA_LOG.Error("退款申请失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("申请成功", c)
	}
}

// RefundOrder 后台退款处理
// @Tags Order
// @Summary 后台退款处理
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shopReq.RefundHandleReq true "退款处理"
// @Success 200 {object} response.Response{msg=string} "退款成功"
// @Router /order/refundOrder [post]
func (orderApi *OrderApi) RefundOrder(c *gin.Context) {
	var req shopReq.RefundHandleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	authorityID := utils.GetUserAuthorityId(c)
	if authorityID != 888 {
		response.FailWithMessage("无权限操作", c)
		return
	}
	if err := orderService.RefundOrder(req.OrderID, req.Remark); err != nil {
		global.GVA_LOG.Error("退款失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("退款成功", c)
	}
}

// UpdateOrderStatus 更新订单状态
// @Tags Order
// @Summary 更新订单状态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query string true "订单ID"
// @Param status query string true "订单状态"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新订单状态"}"
// @Router /order/updateOrderStatus [post]
func (orderApi *OrderApi) UpdateOrderStatus(c *gin.Context) {
	ID := c.Query("ID")
	status := c.Query("status")
	if status != "3" && status != "4" && status != "8" {
		response.FailWithMessage("状态错误", c)
		return
	}
	authorityID := utils.GetUserAuthorityId(c)
	if status == "8" && (authorityID == 888 || authorityID == 8881) {
		response.FailWithMessage("仅用户可提交付款确认", c)
		return
	}
	var err error
	if authorityID == 888 || authorityID == 8881 {
		err = orderService.UpdateOrderStatus(nil, ID, status)
	} else {
		userID := utils.GetUserID(c)
		err = orderService.UpdateOrderStatusForUser(userID, ID, status)
	}
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// ConfirmPayment 管理员确认收款
// @Tags Order
// @Summary 管理员确认收款（将待付款订单标记为已付款）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query string true "订单ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"确认收款成功"}"
// @Router /order/confirmPayment [post]
func (orderApi *OrderApi) ConfirmPayment(c *gin.Context) {
	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("订单ID不能为空", c)
		return
	}
	if err := orderService.ConfirmPayment(ID); err != nil {
		global.GVA_LOG.Error("确认收款失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("确认收款成功", c)
	}
}

// SelfOrderList 我的订单列表
// @Tags Order
// @Summary 我的订单列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"我的订单列表"}"
// @Router /order/selfOrderList [get]
func (orderApi *OrderApi) SelfOrderList(c *gin.Context) {
	userID := utils.GetUserID(c)
	var order shopReq.OrderSearch
	if err := c.ShouldBindQuery(&order); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	order.UserID = utils.Pointer(int(userID))
	if order.PageSize == 0 {
		order.PageSize = 10
	}
	if order.Page == 0 {
		order.Page = 1
	}
	if list, total, err := orderService.GetOrderInfoList(order); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		pageResult := response.PageResult{
			List:     list,
			Total:    total,
			Page:     order.Page,
			PageSize: order.PageSize,
		}
		response.OkWithData(i18n.LocalizeResponseData(c, pageResult), c)
	}
}

// SelfOrder 我的订单
// @Tags Order
// @Summary 我的订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query string true "订单ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"我的订单"}"
// @Router /order/selfOrder [get]
func (orderApi *OrderApi) SelfOrder(c *gin.Context) {
	userID := utils.GetUserID(c)
	ID := c.Query("ID")
	if order, err := orderService.GetOrder(ID, userID); err != nil {
		global.GVA_LOG.Error("查询失败："+err.Error(), zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
	} else {
		response.OkWithData(i18n.LocalizeResponseData(c, order), c)
	}
}

func (orderApi *OrderApi) SelfOrderComment(c *gin.Context) {
	userID := utils.GetUserID(c)
	ID := c.Query("ID")
	SKUID := c.Query("SKUID")
	goodID := c.Query("goodID")
	if order, err := orderService.SelfOrderComment(ID, goodID, SKUID, userID); err != nil {
		global.GVA_LOG.Error("查询失败："+err.Error(), zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
	} else {
		response.OkWithData(i18n.LocalizeResponseData(c, order), c)
	}
}

// DeleteOrder 删除订单
// @Tags Order
// @Summary 删除订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query string true "订单ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /order/deleteOrder [delete]
func (orderApi *OrderApi) DeleteOrder(c *gin.Context) {
	ID := c.Query("ID")
	if err := orderService.DeleteOrder(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteOrderByIds 批量删除订单
// @Tags Order
// @Summary 批量删除订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /order/deleteOrderByIds [delete]
func (orderApi *OrderApi) DeleteOrderByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	if err := orderService.DeleteOrderByIds(IDs); err != nil {
		global.GVA_LOG.Error(err.Error(), zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateOrder 更新订单
// @Tags Order
// @Summary 更新订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shop.Order true "更新订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /order/updateOrder [put]
func (orderApi *OrderApi) UpdateOrder(c *gin.Context) {
	var order shop.Order
	err := c.ShouldBindJSON(&order)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	authorityID := utils.GetUserAuthorityId(c)
	if authorityID == 888 || authorityID == 8881 {
		err = orderService.UpdateOrder(order)
	} else {
		userID := utils.GetUserID(c)
		err = orderService.UpdateOrderForUser(userID, order)
	}

	if err != nil {
		global.GVA_LOG.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindOrder 用id查询订单
// @Tags Order
// @Summary 用id查询订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shop.Order true "用id查询订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /order/findOrder [get]
func (orderApi *OrderApi) FindOrder(c *gin.Context) {
	ID := c.Query("ID")
	if reorder, err := orderService.GetOrder(ID, 0); err != nil {
		global.GVA_LOG.Error(err.Error(), zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(i18n.LocalizeResponseData(c, gin.H{"reorder": reorder}), c)
	}
}

// GetOrderList 分页获取订单列表
// @Tags Order
// @Summary 分页获取订单列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query shopReq.OrderSearch true "分页获取订单列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /order/getOrderList [get]
func (orderApi *OrderApi) GetOrderList(c *gin.Context) {
	var pageInfo shopReq.OrderSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := orderService.GetOrderInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(i18n.LocalizeResponseData(c, response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}), "获取成功", c)
	}
}

// GetOrderPublic 不需要鉴权的订单接口
// @Tags Order
// @Summary 不需要鉴权的订单接口
// @accept application/json
// @Produce application/json
// @Param data query shopReq.OrderSearch true "分页获取订单列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /order/getOrderList [get]
func (orderApi *OrderApi) GetOrderPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的订单接口信息",
	}, "获取成功", c)
}

func (orderApi *OrderApi) CheckRouters(c *gin.Context) {
	express := c.Query("express")
	if err, routers := sf.SfPassPort.SearchRouters(express); err != nil {
		global.GVA_LOG.Error(err.Error(), zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
	} else {
		var req map[string]interface{}
		_ = json.Unmarshal([]byte(routers), &req)
		response.OkWithData(req, c)
	}
}

// BatchUpdateOrderStatus 批量更新订单状态
// @Tags Order
// @Summary 批量更新订单状态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body shopReq.BatchUpdateOrderStatusReq true "批量更新订单状态"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量更新成功"}"
// @Router /order/batchUpdateOrderStatus [post]
func (orderApi *OrderApi) BatchUpdateOrderStatus(c *gin.Context) {
	var req shopReq.BatchUpdateOrderStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	authorityID := utils.GetUserAuthorityId(c)
	if authorityID != 888 {
		response.FailWithMessage("无权限操作", c)
		return
	}
	if err := orderService.BatchUpdateOrderStatus(req.IDs, req.Status); err != nil {
		global.GVA_LOG.Error("批量更新失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("批量更新成功", c)
	}
}
