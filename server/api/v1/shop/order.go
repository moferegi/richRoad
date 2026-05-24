package shop

import (
	"encoding/json"
	"strings"
	"time"

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
	"gorm.io/datatypes"
)

type OrderApi struct {
}

var orderService = service.ServiceGroupApp.ShopServiceGroup.OrderService

func isOrderAdmin(authorityID uint) bool {
	return authorityID == 888 || authorityID == 8881
}

func failWithErr(c *gin.Context, err error) {
	if err == nil {
		response.FailWithMessage(i18n.T(c, "fail"), c)
		return
	}
	key := strings.TrimSpace(err.Error())
	if key != "" && i18n.HasKey(key) {
		response.FailWithMessage(i18n.T(c, key), c)
		return
	}
	response.FailWithMessage(i18n.T(c, "fail"), c)
}

func failWithKey(c *gin.Context, key string) {
	key = strings.TrimSpace(key)
	if key == "" || !i18n.HasKey(key) {
		key = "fail"
	}
	response.FailWithMessage(i18n.T(c, key), c)
}

type publicOrderDetailResponse struct {
	ID        uint               `json:"ID"`
	CreatedAt time.Time          `json:"CreatedAt"`
	OrderID   uint               `json:"orderID"`
	GoodID    uint               `json:"goodID"`
	Good      publicGoodResponse `json:"good"`
	SKUID     uint               `json:"skuID"`
	SKU       publicSkuResponse  `json:"sku"`
	Quantity  uint               `json:"quantity"`
	Price     uint               `json:"price"`
	PriceI18n string             `json:"priceI18n"`
	IsComment bool               `json:"isComment"`
}

type publicOrderResponse struct {
	ID                       uint                        `json:"ID"`
	CreatedAt                time.Time                   `json:"CreatedAt"`
	CouponNum                string                      `json:"couponNum"`
	OriginPrice              uint                        `json:"originPrice"`
	Discount                 uint                        `json:"discount"`
	TotalPrice               uint                        `json:"totalPrice"`
	UsePoints                bool                        `json:"usePoints"`
	PointsUsed               uint                        `json:"pointsUsed"`
	Status                   string                      `json:"status"`
	RefundReason             string                      `json:"refundReason"`
	RefundImages             datatypes.JSON              `json:"refundImages"`
	RefundAppliedAt          *time.Time                  `json:"refundAppliedAt"`
	RefundRemark             string                      `json:"refundRemark"`
	RefundHandledAt          *time.Time                  `json:"refundHandledAt"`
	Detail                   []publicOrderDetailResponse `json:"detail"`
	Express                  string                      `json:"express"`
	Phone                    string                      `json:"phone"`
	Name                     string                      `json:"name"`
	Province                 string                      `json:"province"`
	City                     string                      `json:"city"`
	Area                     string                      `json:"area"`
	Street                   string                      `json:"street"`
	CloseTime                time.Time                   `json:"closeTime"`
	IsPresale                bool                        `json:"isPresale"`
	PayMethod                string                      `json:"payMethod"`
	SettlementCurrency       string                      `json:"settlementCurrency"`
	SettlementCurrencySymbol string                      `json:"settlementCurrencySymbol"`
	ExpireAt                 *time.Time                  `json:"expireAt"`
	PaidAt                   *time.Time                  `json:"paidAt"`
	ReceivedAt               *time.Time                  `json:"receivedAt"`
	CancelledAt              *time.Time                  `json:"cancelledAt"`
}

type publicOrderCommentResponse struct {
	ID                       uint                      `json:"ID"`
	CreatedAt                time.Time                 `json:"CreatedAt"`
	TotalPrice               uint                      `json:"totalPrice"`
	Status                   string                    `json:"status"`
	RefundReason             string                    `json:"refundReason"`
	RefundImages             datatypes.JSON            `json:"refundImages"`
	RefundAppliedAt          *time.Time                `json:"refundAppliedAt"`
	RefundRemark             string                    `json:"refundRemark"`
	RefundHandledAt          *time.Time                `json:"refundHandledAt"`
	Detail                   publicOrderDetailResponse `json:"detail"`
	Express                  string                    `json:"express"`
	Phone                    string                    `json:"phone"`
	Name                     string                    `json:"name"`
	SettlementCurrency       string                    `json:"settlementCurrency"`
	SettlementCurrencySymbol string                    `json:"settlementCurrencySymbol"`
	Province                 string                    `json:"province"`
	City                     string                    `json:"city"`
	Area                     string                    `json:"area"`
	Street                   string                    `json:"street"`
	CloseTime                time.Time                 `json:"closeTime"`
	Comment                  shop.Comment              `json:"comment"`
}

func toPublicOrderDetailResponse(item shop.OrderDetailRes) publicOrderDetailResponse {
	return publicOrderDetailResponse{
		ID:        item.ID,
		CreatedAt: item.CreatedAt,
		OrderID:   item.OrderID,
		GoodID:    item.GoodID,
		Good:      toPublicGoodResponse(item.Good),
		SKUID:     item.SKUID,
		SKU:       toPublicSkuResponse(item.SKU),
		Quantity:  item.Quantity,
		Price:     item.Price,
		PriceI18n: item.PriceI18n,
		IsComment: item.IsComment,
	}
}

func toPublicOrderDetailResponses(list []shop.OrderDetailRes) []publicOrderDetailResponse {
	result := make([]publicOrderDetailResponse, 0, len(list))
	for _, item := range list {
		result = append(result, toPublicOrderDetailResponse(item))
	}
	return result
}

func toPublicOrderResponse(item shop.OrderRes) publicOrderResponse {
	return publicOrderResponse{
		ID:                       item.ID,
		CreatedAt:                item.CreatedAt,
		CouponNum:                item.CouponNum,
		OriginPrice:              item.OriginPrice,
		Discount:                 item.Discount,
		TotalPrice:               item.TotalPrice,
		UsePoints:                item.UsePoints,
		PointsUsed:               item.PointsUsed,
		Status:                   item.Status,
		RefundReason:             item.RefundReason,
		RefundImages:             item.RefundImages,
		RefundAppliedAt:          item.RefundAppliedAt,
		RefundRemark:             item.RefundRemark,
		RefundHandledAt:          item.RefundHandledAt,
		Detail:                   toPublicOrderDetailResponses(item.Detail),
		Express:                  item.Express,
		Phone:                    item.Phone,
		Name:                     item.Name,
		Province:                 item.Province,
		City:                     item.City,
		Area:                     item.Area,
		Street:                   item.Street,
		CloseTime:                item.CloseTime,
		IsPresale:                item.IsPresale,
		PayMethod:                item.PayMethod,
		SettlementCurrency:       item.SettlementCurrency,
		SettlementCurrencySymbol: item.SettlementCurrencySymbol,
		ExpireAt:                 item.ExpireAt,
		PaidAt:                   item.PaidAt,
		ReceivedAt:               item.ReceivedAt,
		CancelledAt:              item.CancelledAt,
	}
}

func toPublicOrderResponses(list []shop.OrderRes) []publicOrderResponse {
	result := make([]publicOrderResponse, 0, len(list))
	for _, item := range list {
		result = append(result, toPublicOrderResponse(item))
	}
	return result
}

func toPublicOrderCommentResponse(item shop.OrderCommentRes) publicOrderCommentResponse {
	return publicOrderCommentResponse{
		ID:                       item.ID,
		CreatedAt:                item.CreatedAt,
		TotalPrice:               item.TotalPrice,
		Status:                   item.Status,
		RefundReason:             item.RefundReason,
		RefundImages:             item.RefundImages,
		RefundAppliedAt:          item.RefundAppliedAt,
		RefundRemark:             item.RefundRemark,
		RefundHandledAt:          item.RefundHandledAt,
		Detail:                   toPublicOrderDetailResponse(item.Detail),
		Express:                  item.Express,
		Phone:                    item.Phone,
		Name:                     item.Name,
		SettlementCurrency:       item.SettlementCurrency,
		SettlementCurrencySymbol: item.SettlementCurrencySymbol,
		Province:                 item.Province,
		City:                     item.City,
		Area:                     item.Area,
		Street:                   item.Street,
		CloseTime:                item.CloseTime,
		Comment:                  item.Comment,
	}
}

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
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}

	var order shop.Order
	err := c.ShouldBindJSON(&order)
	if err != nil {
		failWithKey(c, "invalidParams")
		return
	}

	if err := orderService.CreateOrder(&order); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		failWithErr(c, err)
	} else {
		response.OkWithMessage(i18n.T(c, "createSuccess"), c)
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
		failWithKey(c, "invalidParams")
		return
	}
	order.UserID = utils.GetUserID(c)

	if orderID, orderNo, err := orderService.PlaceOrder(&order); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		failWithErr(c, err)
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
		failWithKey(c, "invalidParams")
		return
	}
	if orderID, orderNo, err := orderService.PlaceOrderByCart(userID, req); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		failWithErr(c, err)
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
		failWithErr(c, err)
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
		failWithErr(c, err)
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
		failWithKey(c, "invalidParams")
		return
	}
	if err := orderService.ApplyRefund(userID, req); err != nil {
		global.GVA_LOG.Error("退款申请失败!", zap.Error(err))
		failWithErr(c, err)
	} else {
		response.OkWithMessage(i18n.T(c, "submitSuccess"), c)
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
		failWithKey(c, "invalidParams")
		return
	}
	authorityID := utils.GetUserAuthorityId(c)
	if !isOrderAdmin(authorityID) {
		failWithKey(c, "noPermission")
		return
	}
	if err := orderService.RefundOrder(req.OrderID, req.Remark); err != nil {
		global.GVA_LOG.Error("退款失败!", zap.Error(err))
		failWithErr(c, err)
	} else {
		response.OkWithMessage(i18n.T(c, "confirmSuccess"), c)
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
		failWithKey(c, "orderStatusInvalid")
		return
	}
	authorityID := utils.GetUserAuthorityId(c)
	if status == "8" && isOrderAdmin(authorityID) {
		failWithKey(c, "orderOnlyUserCanSubmitPendingConfirm")
		return
	}
	var err error
	if isOrderAdmin(authorityID) {
		err = orderService.UpdateOrderStatus(nil, ID, status)
	} else {
		userID := utils.GetUserID(c)
		err = orderService.UpdateOrderStatusForUser(userID, ID, status)
	}
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		failWithErr(c, err)
	} else {
		response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
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
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}

	ID := c.Query("ID")
	if ID == "" {
		failWithKey(c, "orderIDRequired")
		return
	}
	if err := orderService.ConfirmPayment(ID); err != nil {
		global.GVA_LOG.Error("确认收款失败!", zap.Error(err))
		failWithErr(c, err)
	} else {
		response.OkWithMessage(i18n.T(c, "confirmSuccess"), c)
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
		failWithKey(c, "invalidParams")
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
		failWithKey(c, "getFail")
	} else {
		pageResult := response.PageResult{
			List:     toPublicOrderResponses(list),
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
		failWithKey(c, "queryFail")
	} else {
		response.OkWithData(i18n.LocalizeResponseData(c, toPublicOrderResponse(order)), c)
	}
}

func (orderApi *OrderApi) SelfOrderComment(c *gin.Context) {
	userID := utils.GetUserID(c)
	ID := c.Query("ID")
	SKUID := c.Query("SKUID")
	goodID := c.Query("goodID")
	if order, err := orderService.SelfOrderComment(ID, goodID, SKUID, userID); err != nil {
		global.GVA_LOG.Error("查询失败："+err.Error(), zap.Error(err))
		failWithKey(c, "queryFail")
	} else {
		response.OkWithData(i18n.LocalizeResponseData(c, toPublicOrderCommentResponse(order)), c)
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
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}

	ID := c.Query("ID")
	if err := orderService.DeleteOrder(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		failWithErr(c, err)
	} else {
		response.OkWithMessage(i18n.T(c, "deleteSuccess"), c)
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
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}

	IDs := c.QueryArray("IDs[]")
	if err := orderService.DeleteOrderByIds(IDs); err != nil {
		global.GVA_LOG.Error(err.Error(), zap.Error(err))
		failWithKey(c, "batchDeleteFail")
	} else {
		response.OkWithMessage(i18n.T(c, "batchDeleteSuccess"), c)
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
		failWithKey(c, "invalidParams")
		return
	}
	authorityID := utils.GetUserAuthorityId(c)
	if isOrderAdmin(authorityID) {
		err = orderService.UpdateOrder(order)
	} else {
		userID := utils.GetUserID(c)
		err = orderService.UpdateOrderForUser(userID, order)
	}

	if err != nil {
		global.GVA_LOG.Error(err.Error(), zap.Error(err))
		failWithErr(c, err)
	} else {
		response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
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
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}

	ID := c.Query("ID")
	if reorder, err := orderService.GetOrder(ID, 0); err != nil {
		global.GVA_LOG.Error(err.Error(), zap.Error(err))
		failWithKey(c, "queryFail")
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
	if !isOrderAdmin(utils.GetUserAuthorityId(c)) {
		failWithKey(c, "noPermission")
		return
	}

	var pageInfo shopReq.OrderSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		failWithKey(c, "invalidParams")
		return
	}
	if list, total, err := orderService.GetOrderInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		failWithErr(c, err)
	} else {
		response.OkWithDetailed(i18n.LocalizeResponseData(c, response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}), i18n.T(c, "getSuccess"), c)
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
	express := strings.TrimSpace(c.Query("express"))
	if express == "" || len(express) > 128 {
		failWithKey(c, "paramError")
		return
	}
	authorityID := utils.GetUserAuthorityId(c)
	if !isOrderAdmin(authorityID) {
		allowed, err := orderService.UserOwnsExpress(utils.GetUserID(c), express)
		if err != nil {
			global.GVA_LOG.Error("物流归属校验失败", zap.Error(err))
			failWithKey(c, "queryFail")
			return
		}
		if !allowed {
			failWithKey(c, "noPermission")
			return
		}
	}
	if err, routers := sf.SfPassPort.SearchRouters(express); err != nil {
		global.GVA_LOG.Error(err.Error(), zap.Error(err))
		failWithKey(c, "queryFail")
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
		failWithKey(c, "invalidParams")
		return
	}
	authorityID := utils.GetUserAuthorityId(c)
	if authorityID != 888 {
		failWithKey(c, "noPermission")
		return
	}
	if err := orderService.BatchUpdateOrderStatus(req.IDs, req.Status); err != nil {
		global.GVA_LOG.Error("批量更新失败!", zap.Error(err))
		failWithErr(c, err)
	} else {
		response.OkWithMessage(i18n.T(c, "updateSuccess"), c)
	}
}
