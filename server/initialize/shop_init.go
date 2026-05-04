package initialize

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitShopData 自动注册商城核心模块（商品、订单、购物车、优惠券等）的API、菜单、casbin权限
// 仅在数据缺失时插入，幂等执行
func InitShopData() {
	if global.GVA_DB == nil {
		return
	}
	db := global.GVA_DB

	initShopApis(db)
	initShopCasbin(db)
	fixKefuMenuComponent(db)
}

// fixKefuMenuComponent 修正 kefu 菜单指向正确的客服配置页组件
// 历史数据可能指向自动生成的 kefuService 页，需要修正为带开关的 kefu 配置页
func fixKefuMenuComponent(db *gorm.DB) {
	result := db.Model(&sysModel.SysBaseMenu{}).
		Where("name = ? AND component != ?", "kefu", "view/shop/kefu/kefu.vue").
		Update("component", "view/shop/kefu/kefu.vue")
	if result.Error != nil {
		global.GVA_LOG.Error("修正 kefu 菜单组件路径失败", zap.Error(result.Error))
	} else if result.RowsAffected > 0 {
		global.GVA_LOG.Info("kefu 菜单组件路径已修正为 view/shop/kefu/kefu.vue")
	}
}

func initShopApis(db *gorm.DB) {
	apis := []sysModel.SysApi{
		// 商品管理
		{ApiGroup: "商品管理", Method: "POST", Path: "/good/createGood", Description: "创建商品"},
		{ApiGroup: "商品管理", Method: "DELETE", Path: "/good/deleteGood", Description: "删除商品"},
		{ApiGroup: "商品管理", Method: "DELETE", Path: "/good/deleteGoodByIds", Description: "批量删除商品"},
		{ApiGroup: "商品管理", Method: "PUT", Path: "/good/updateGood", Description: "更新商品"},
		// 轮播图管理
		{ApiGroup: "轮播图管理", Method: "POST", Path: "/banner/createBanner", Description: "创建轮播图"},
		{ApiGroup: "轮播图管理", Method: "DELETE", Path: "/banner/deleteBanner", Description: "删除轮播图"},
		{ApiGroup: "轮播图管理", Method: "DELETE", Path: "/banner/deleteBannerByIds", Description: "批量删除轮播图"},
		{ApiGroup: "轮播图管理", Method: "PUT", Path: "/banner/updateBanner", Description: "更新轮播图"},
		{ApiGroup: "轮播图管理", Method: "GET", Path: "/banner/findBanner", Description: "查找轮播图"},
		// 分类管理
		{ApiGroup: "分类管理", Method: "POST", Path: "/category/createCategory", Description: "创建分类"},
		{ApiGroup: "分类管理", Method: "DELETE", Path: "/category/deleteCategory", Description: "删除分类"},
		{ApiGroup: "分类管理", Method: "DELETE", Path: "/category/deleteCategoryByIds", Description: "批量删除分类"},
		{ApiGroup: "分类管理", Method: "PUT", Path: "/category/updateCategory", Description: "更新分类"},
		{ApiGroup: "分类管理", Method: "GET", Path: "/category/findCategory", Description: "查找分类"},
		// SKU管理
		{ApiGroup: "SKU管理", Method: "POST", Path: "/sku/createSku", Description: "创建SKU"},
		{ApiGroup: "SKU管理", Method: "DELETE", Path: "/sku/deleteSku", Description: "删除SKU"},
		{ApiGroup: "SKU管理", Method: "DELETE", Path: "/sku/deleteSkuByIds", Description: "批量删除SKU"},
		{ApiGroup: "SKU管理", Method: "PUT", Path: "/sku/updateSku", Description: "更新SKU"},
		{ApiGroup: "SKU管理", Method: "GET", Path: "/sku/findSku", Description: "查找SKU"},
		{ApiGroup: "SKU管理", Method: "GET", Path: "/sku/getSkuList", Description: "获取SKU列表"},
		// 购物车管理
		{ApiGroup: "购物车管理", Method: "POST", Path: "/cart/createCart", Description: "创建购物车"},
		{ApiGroup: "购物车管理", Method: "DELETE", Path: "/cart/deleteCart", Description: "删除购物车"},
		{ApiGroup: "购物车管理", Method: "DELETE", Path: "/cart/deleteCartByIds", Description: "批量删除购物车"},
		{ApiGroup: "购物车管理", Method: "PUT", Path: "/cart/updateCart", Description: "更新购物车"},
		{ApiGroup: "购物车管理", Method: "POST", Path: "/cart/addCart", Description: "添加购物车"},
		{ApiGroup: "购物车管理", Method: "POST", Path: "/cart/cutCart", Description: "减少购物车数量"},
		{ApiGroup: "购物车管理", Method: "GET", Path: "/cart/getSelfCart", Description: "获取我的购物车"},
		{ApiGroup: "购物车管理", Method: "GET", Path: "/cart/clearCart", Description: "清空购物车"},
		{ApiGroup: "购物车管理", Method: "GET", Path: "/cart/findCart", Description: "查找购物车"},
		{ApiGroup: "购物车管理", Method: "GET", Path: "/cart/getCartList", Description: "获取购物车列表"},
		// 订单管理
		{ApiGroup: "订单管理", Method: "POST", Path: "/order/createOrder", Description: "创建订单"},
		{ApiGroup: "订单管理", Method: "DELETE", Path: "/order/deleteOrder", Description: "删除订单"},
		{ApiGroup: "订单管理", Method: "DELETE", Path: "/order/deleteOrderByIds", Description: "批量删除订单"},
		{ApiGroup: "订单管理", Method: "PUT", Path: "/order/updateOrder", Description: "更新订单"},
		{ApiGroup: "订单管理", Method: "POST", Path: "/order/placeOrder", Description: "下单"},
		{ApiGroup: "订单管理", Method: "POST", Path: "/order/placeOrderByCart", Description: "购物车下单"},
		{ApiGroup: "订单管理", Method: "POST", Path: "/order/changeOrderCoupon", Description: "切换订单优惠券"},
		{ApiGroup: "订单管理", Method: "POST", Path: "/order/changeOrderPoints", Description: "切换订单积分"},
		{ApiGroup: "订单管理", Method: "POST", Path: "/order/applyRefund", Description: "申请退款"},
		{ApiGroup: "订单管理", Method: "POST", Path: "/order/refundOrder", Description: "处理退款"},
		{ApiGroup: "订单管理", Method: "POST", Path: "/order/updateOrderStatus", Description: "更新订单状态"},
		{ApiGroup: "订单管理", Method: "POST", Path: "/order/confirmPayment", Description: "确认收款"},
		{ApiGroup: "订单管理", Method: "POST", Path: "/order/batchUpdateOrderStatus", Description: "批量更新订单状态"},
		{ApiGroup: "订单管理", Method: "GET", Path: "/order/selfOrder", Description: "我的订单详情"},
		{ApiGroup: "订单管理", Method: "GET", Path: "/order/selfOrderList", Description: "我的订单列表"},
		{ApiGroup: "订单管理", Method: "GET", Path: "/order/selfOrderComment", Description: "我的订单评论"},
		{ApiGroup: "订单管理", Method: "GET", Path: "/order/findOrder", Description: "查找订单"},
		{ApiGroup: "订单管理", Method: "GET", Path: "/order/getOrderList", Description: "获取订单列表"},
		// 评论管理
		{ApiGroup: "评论管理", Method: "POST", Path: "/comment/createComment", Description: "创建评论"},
		{ApiGroup: "评论管理", Method: "DELETE", Path: "/comment/deleteComment", Description: "删除评论"},
		{ApiGroup: "评论管理", Method: "DELETE", Path: "/comment/deleteCommentByIds", Description: "批量删除评论"},
		{ApiGroup: "评论管理", Method: "PUT", Path: "/comment/updateComment", Description: "更新评论"},
		// 标签管理
		{ApiGroup: "标签管理", Method: "POST", Path: "/tag/createTag", Description: "创建标签"},
		{ApiGroup: "标签管理", Method: "DELETE", Path: "/tag/deleteTag", Description: "删除标签"},
		{ApiGroup: "标签管理", Method: "DELETE", Path: "/tag/deleteTagByIds", Description: "批量删除标签"},
		{ApiGroup: "标签管理", Method: "PUT", Path: "/tag/updateTag", Description: "更新标签"},
		{ApiGroup: "标签管理", Method: "GET", Path: "/tag/findTag", Description: "查找标签"},
		{ApiGroup: "标签管理", Method: "GET", Path: "/tag/getTagList", Description: "获取标签列表"},
		// 优惠券管理
		{ApiGroup: "优惠券管理", Method: "POST", Path: "/Cou/createCoupon", Description: "创建优惠券"},
		{ApiGroup: "优惠券管理", Method: "DELETE", Path: "/Cou/deleteCoupon", Description: "删除优惠券"},
		{ApiGroup: "优惠券管理", Method: "DELETE", Path: "/Cou/deleteCouponByIds", Description: "批量删除优惠券"},
		{ApiGroup: "优惠券管理", Method: "PUT", Path: "/Cou/updateCoupon", Description: "更新优惠券"},
		{ApiGroup: "优惠券管理", Method: "GET", Path: "/Cou/findCoupon", Description: "查找优惠券"},
		{ApiGroup: "优惠券管理", Method: "GET", Path: "/Cou/getCouponList", Description: "获取优惠券列表"},
		// 用户优惠券管理
		{ApiGroup: "用户优惠券", Method: "POST", Path: "/cou/createCouponOrderUser", Description: "创建用户优惠券"},
		{ApiGroup: "用户优惠券", Method: "DELETE", Path: "/cou/deleteCouponOrderUser", Description: "删除用户优惠券"},
		{ApiGroup: "用户优惠券", Method: "DELETE", Path: "/cou/deleteCouponOrderUserByIds", Description: "批量删除用户优惠券"},
		{ApiGroup: "用户优惠券", Method: "PUT", Path: "/cou/updateCouponOrderUser", Description: "更新用户优惠券"},
		{ApiGroup: "用户优惠券", Method: "GET", Path: "/cou/findCouponOrderUser", Description: "查找用户优惠券"},
		{ApiGroup: "用户优惠券", Method: "GET", Path: "/cou/getCouponOrderUserList", Description: "获取用户优惠券列表"},
		{ApiGroup: "用户优惠券", Method: "POST", Path: "/cou/getAllClaimCoupon", Description: "获取所有可领优惠券"},
		{ApiGroup: "用户优惠券", Method: "POST", Path: "/cou/claimCouponByUser", Description: "用户领取优惠券"},
		{ApiGroup: "用户优惠券", Method: "POST", Path: "/cou/adminIssueCouponToAll", Description: "管理员发放优惠券"},
		// 促销管理
		{ApiGroup: "促销管理", Method: "POST", Path: "/promo/createPromotion", Description: "创建促销"},
		{ApiGroup: "促销管理", Method: "DELETE", Path: "/promo/deletePromotion", Description: "删除促销"},
		{ApiGroup: "促销管理", Method: "DELETE", Path: "/promo/deletePromotionByIds", Description: "批量删除促销"},
		{ApiGroup: "促销管理", Method: "PUT", Path: "/promo/updatePromotion", Description: "更新促销"},
		{ApiGroup: "促销管理", Method: "GET", Path: "/promo/findPromotion", Description: "查找促销"},
		{ApiGroup: "促销管理", Method: "GET", Path: "/promo/getPromotionList", Description: "获取促销列表"},
		// 客服管理
		{ApiGroup: "客服管理", Method: "POST", Path: "/kefu/createKefu", Description: "创建客服"},
		{ApiGroup: "客服管理", Method: "DELETE", Path: "/kefu/deleteKefu", Description: "删除客服"},
		{ApiGroup: "客服管理", Method: "DELETE", Path: "/kefu/deleteKefuByIds", Description: "批量删除客服"},
		{ApiGroup: "客服管理", Method: "PUT", Path: "/kefu/updateKefu", Description: "更新客服"},
		{ApiGroup: "客服管理", Method: "GET", Path: "/kefu/findKefu", Description: "查找客服"},
		{ApiGroup: "客服管理", Method: "GET", Path: "/kefu/getKefuList", Description: "获取客服列表"},
		// 浏览历史
		{ApiGroup: "浏览历史", Method: "GET", Path: "/good/getGoodHistory", Description: "获取浏览历史"},
		// 地址管理
		{ApiGroup: "地址管理", Method: "POST", Path: "/address/createAddress", Description: "创建地址"},
		{ApiGroup: "地址管理", Method: "DELETE", Path: "/address/deleteAddress", Description: "删除地址"},
		{ApiGroup: "地址管理", Method: "PUT", Path: "/address/updateAddress", Description: "更新地址"},
		{ApiGroup: "地址管理", Method: "GET", Path: "/address/findAddress", Description: "查找地址"},
		{ApiGroup: "地址管理", Method: "GET", Path: "/address/getAddressList", Description: "获取地址列表"},
		// 收藏管理
		{ApiGroup: "收藏管理", Method: "POST", Path: "/collect/createCollect", Description: "创建收藏"},
		{ApiGroup: "收藏管理", Method: "DELETE", Path: "/collect/deleteCollect", Description: "删除收藏"},
		{ApiGroup: "收藏管理", Method: "GET", Path: "/collect/getCollectList", Description: "获取收藏列表"},
		{ApiGroup: "收藏管理", Method: "GET", Path: "/collect/checkCollect", Description: "检查是否收藏"},
		// 积分记录
		{ApiGroup: "积分管理", Method: "POST", Path: "/cpr/createPointRecord", Description: "创建积分记录"},
		{ApiGroup: "积分管理", Method: "DELETE", Path: "/cpr/deletePointRecord", Description: "删除积分记录"},
		{ApiGroup: "积分管理", Method: "PUT", Path: "/cpr/updatePointRecord", Description: "更新积分记录"},
		{ApiGroup: "积分管理", Method: "GET", Path: "/cpr/findPointRecord", Description: "查找积分记录"},
		{ApiGroup: "积分管理", Method: "GET", Path: "/cpr/getPointRecordList", Description: "获取积分记录列表"},
		// 系统配置
		{ApiGroup: "系统配置", Method: "PUT", Path: "/sysConfig/updateSysConfig", Description: "更新系统配置"},
		{ApiGroup: "系统配置", Method: "GET", Path: "/sysConfig/getSysConfigList", Description: "获取系统配置列表"},
		{ApiGroup: "系统配置", Method: "GET", Path: "/sysConfig/getConfigByKey", Description: "按key获取配置"},
		{ApiGroup: "系统配置", Method: "GET", Path: "/sysConfig/getSysConfigByGroup", Description: "按分组获取配置"},
		{ApiGroup: "系统配置", Method: "GET", Path: "/sysConfig/getSysConfigByKey", Description: "按key获取单条配置"},
		// 外部链接域名
		{ApiGroup: "外部链接域名", Method: "POST", Path: "/extDomain/createExternalLinkDomain", Description: "创建外部链接域名"},
		{ApiGroup: "外部链接域名", Method: "DELETE", Path: "/extDomain/deleteExternalLinkDomain", Description: "删除外部链接域名"},
		{ApiGroup: "外部链接域名", Method: "PUT", Path: "/extDomain/updateExternalLinkDomain", Description: "更新外部链接域名"},
		{ApiGroup: "外部链接域名", Method: "GET", Path: "/extDomain/getExternalLinkDomainList", Description: "获取外部链接域名列表"},
		{ApiGroup: "外部链接域名", Method: "POST", Path: "/extDomain/setDefaultDomain", Description: "设置默认域名"},
		// 客户端用户
		{ApiGroup: "客户端用户", Method: "GET", Path: "/clientUser/getUserInfo", Description: "获取用户信息"},
		{ApiGroup: "客户端用户", Method: "GET", Path: "/clientUser/getUserList", Description: "获取用户列表"},
		{ApiGroup: "客户端用户", Method: "PUT", Path: "/clientUser/updateUser", Description: "更新用户"},
		{ApiGroup: "客户端用户", Method: "DELETE", Path: "/clientUser/deleteUser", Description: "删除用户"},
		{ApiGroup: "客户端用户", Method: "POST", Path: "/clientUser/changePassword", Description: "修改密码"},
		{ApiGroup: "客户端用户", Method: "POST", Path: "/clientUser/setPhoneVerified", Description: "验证后设置手机号"},
		// SKU规格字典
		{ApiGroup: "SKU规格字典", Method: "POST", Path: "/skuSpec/createSkuSpec", Description: "创建SKU规格字典"},
		{ApiGroup: "SKU规格字典", Method: "DELETE", Path: "/skuSpec/deleteSkuSpec", Description: "删除SKU规格字典"},
		{ApiGroup: "SKU规格字典", Method: "DELETE", Path: "/skuSpec/deleteSkuSpecByIds", Description: "批量删除SKU规格字典"},
		{ApiGroup: "SKU规格字典", Method: "PUT", Path: "/skuSpec/updateSkuSpec", Description: "更新SKU规格字典"},
		{ApiGroup: "SKU规格字典", Method: "GET", Path: "/skuSpec/findSkuSpec", Description: "查找SKU规格字典"},
		{ApiGroup: "SKU规格字典", Method: "GET", Path: "/skuSpec/getSkuSpecList", Description: "获取SKU规格字典列表"},
		{ApiGroup: "SKU规格字典", Method: "GET", Path: "/skuSpec/getAllSkuSpecs", Description: "获取所有SKU规格字典"},
	}
	for _, api := range apis {
		var count int64
		db.Model(&sysModel.SysApi{}).Where("path = ? AND method = ?", api.Path, api.Method).Count(&count)
		if count == 0 {
			if err := db.Create(&api).Error; err != nil {
				global.GVA_LOG.Error("初始化商城API失败", zap.Error(err))
			}
		}
	}
}

func initShopCasbin(db *gorm.DB) {
	var authorities []sysModel.SysAuthority
	db.Find(&authorities)

	paths := []struct {
		Path   string
		Method string
	}{
		// 商品
		{"/good/createGood", "POST"},
		{"/good/deleteGood", "DELETE"},
		{"/good/deleteGoodByIds", "DELETE"},
		{"/good/updateGood", "PUT"},
		// 轮播图
		{"/banner/createBanner", "POST"},
		{"/banner/deleteBanner", "DELETE"},
		{"/banner/deleteBannerByIds", "DELETE"},
		{"/banner/updateBanner", "PUT"},
		{"/banner/findBanner", "GET"},
		// 分类
		{"/category/createCategory", "POST"},
		{"/category/deleteCategory", "DELETE"},
		{"/category/deleteCategoryByIds", "DELETE"},
		{"/category/updateCategory", "PUT"},
		{"/category/findCategory", "GET"},
		// SKU
		{"/sku/createSku", "POST"},
		{"/sku/deleteSku", "DELETE"},
		{"/sku/deleteSkuByIds", "DELETE"},
		{"/sku/updateSku", "PUT"},
		{"/sku/findSku", "GET"},
		{"/sku/getSkuList", "GET"},
		// 购物车
		{"/cart/createCart", "POST"},
		{"/cart/deleteCart", "DELETE"},
		{"/cart/deleteCartByIds", "DELETE"},
		{"/cart/updateCart", "PUT"},
		{"/cart/addCart", "POST"},
		{"/cart/cutCart", "POST"},
		{"/cart/getSelfCart", "GET"},
		{"/cart/clearCart", "GET"},
		{"/cart/findCart", "GET"},
		{"/cart/getCartList", "GET"},
		// 订单
		{"/order/createOrder", "POST"},
		{"/order/deleteOrder", "DELETE"},
		{"/order/deleteOrderByIds", "DELETE"},
		{"/order/updateOrder", "PUT"},
		{"/order/placeOrder", "POST"},
		{"/order/placeOrderByCart", "POST"},
		{"/order/changeOrderCoupon", "POST"},
		{"/order/changeOrderPoints", "POST"},
		{"/order/applyRefund", "POST"},
		{"/order/refundOrder", "POST"},
		{"/order/updateOrderStatus", "POST"},
		{"/order/confirmPayment", "POST"},
		{"/order/batchUpdateOrderStatus", "POST"},
		{"/order/selfOrder", "GET"},
		{"/order/selfOrderList", "GET"},
		{"/order/selfOrderComment", "GET"},
		{"/order/findOrder", "GET"},
		{"/order/getOrderList", "GET"},
		// 评论
		{"/comment/createComment", "POST"},
		{"/comment/deleteComment", "DELETE"},
		{"/comment/deleteCommentByIds", "DELETE"},
		{"/comment/updateComment", "PUT"},
		// 标签
		{"/tag/createTag", "POST"},
		{"/tag/deleteTag", "DELETE"},
		{"/tag/deleteTagByIds", "DELETE"},
		{"/tag/updateTag", "PUT"},
		{"/tag/findTag", "GET"},
		{"/tag/getTagList", "GET"},
		// 优惠券
		{"/Cou/createCoupon", "POST"},
		{"/Cou/deleteCoupon", "DELETE"},
		{"/Cou/deleteCouponByIds", "DELETE"},
		{"/Cou/updateCoupon", "PUT"},
		{"/Cou/findCoupon", "GET"},
		{"/Cou/getCouponList", "GET"},
		// 用户优惠券
		{"/cou/createCouponOrderUser", "POST"},
		{"/cou/deleteCouponOrderUser", "DELETE"},
		{"/cou/deleteCouponOrderUserByIds", "DELETE"},
		{"/cou/updateCouponOrderUser", "PUT"},
		{"/cou/findCouponOrderUser", "GET"},
		{"/cou/getCouponOrderUserList", "GET"},
		{"/cou/getAllClaimCoupon", "POST"},
		{"/cou/claimCouponByUser", "POST"},
		{"/cou/adminIssueCouponToAll", "POST"},
		// 促销
		{"/promo/createPromotion", "POST"},
		{"/promo/deletePromotion", "DELETE"},
		{"/promo/deletePromotionByIds", "DELETE"},
		{"/promo/updatePromotion", "PUT"},
		{"/promo/findPromotion", "GET"},
		{"/promo/getPromotionList", "GET"},
		// 客服
		{"/kefu/createKefu", "POST"},
		{"/kefu/deleteKefu", "DELETE"},
		{"/kefu/deleteKefuByIds", "DELETE"},
		{"/kefu/updateKefu", "PUT"},
		{"/kefu/findKefu", "GET"},
		{"/kefu/getKefuList", "GET"},
		// 浏览历史
		{"/good/getGoodHistory", "GET"},
		// 地址
		{"/address/createAddress", "POST"},
		{"/address/deleteAddress", "DELETE"},
		{"/address/updateAddress", "PUT"},
		{"/address/findAddress", "GET"},
		{"/address/getAddressList", "GET"},
		// 收藏
		{"/collect/createCollect", "POST"},
		{"/collect/deleteCollect", "DELETE"},
		{"/collect/getCollectList", "GET"},
		{"/collect/checkCollect", "GET"},
		// 积分
		{"/cpr/createPointRecord", "POST"},
		{"/cpr/deletePointRecord", "DELETE"},
		{"/cpr/updatePointRecord", "PUT"},
		{"/cpr/findPointRecord", "GET"},
		{"/cpr/getPointRecordList", "GET"},
		// 系统配置
		{"/sysConfig/updateSysConfig", "PUT"},
		{"/sysConfig/getSysConfigList", "GET"},
		{"/sysConfig/getConfigByKey", "GET"},
		{"/sysConfig/getSysConfigByGroup", "GET"},
		{"/sysConfig/getSysConfigByKey", "GET"},
		// 外部链接域名
		{"/extDomain/createExternalLinkDomain", "POST"},
		{"/extDomain/deleteExternalLinkDomain", "DELETE"},
		{"/extDomain/updateExternalLinkDomain", "PUT"},
		{"/extDomain/getExternalLinkDomainList", "GET"},
		{"/extDomain/setDefaultDomain", "POST"},
		// 客户端用户
		{"/clientUser/getUserInfo", "GET"},
		{"/clientUser/getUserList", "GET"},
		{"/clientUser/updateUser", "PUT"},
		{"/clientUser/deleteUser", "DELETE"},
		{"/clientUser/changePassword", "POST"},
		{"/clientUser/setPhoneVerified", "POST"},
		// SKU规格字典
		{"/skuSpec/createSkuSpec", "POST"},
		{"/skuSpec/deleteSkuSpec", "DELETE"},
		{"/skuSpec/deleteSkuSpecByIds", "DELETE"},
		{"/skuSpec/updateSkuSpec", "PUT"},
		{"/skuSpec/findSkuSpec", "GET"},
		{"/skuSpec/getSkuSpecList", "GET"},
		{"/skuSpec/getAllSkuSpecs", "GET"},
	}

	for _, auth := range authorities {
		authId := fmt.Sprintf("%d", auth.AuthorityId)
		for _, p := range paths {
			var count int64
			db.Table("casbin_rule").Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?",
				"p", authId, p.Path, p.Method).Count(&count)
			if count == 0 {
				db.Exec("INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES (?, ?, ?, ?)",
					"p", authId, p.Path, p.Method)
			}
		}
	}
}
