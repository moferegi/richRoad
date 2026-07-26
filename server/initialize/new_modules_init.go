package initialize

import (
	"errors"
	"fmt"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitNewModulesData 自动注册新增模块（进货、收款码、弹窗、营销奖励、看板、预售、语言、区号、签到）的API、菜单、casbin权限
func InitNewModulesData() {
	if global.GVA_DB == nil {
		return
	}
	db := global.GVA_DB

	initNewModulesApis(db)
	initNewModulesMenus(db)
	initNewModulesMenuBtns(db)
	initNewModulesCasbin(db)
	initLanguageSeedData(db)
}

func initNewModulesApis(db *gorm.DB) {
	apis := []sysModel.SysApi{
		// 进货记录
		{ApiGroup: "进货管理", Method: "POST", Path: "/goodPurchase/createGoodPurchase", Description: "创建进货记录"},
		{ApiGroup: "进货管理", Method: "DELETE", Path: "/goodPurchase/deleteGoodPurchase", Description: "删除进货记录"},
		{ApiGroup: "进货管理", Method: "PUT", Path: "/goodPurchase/updateGoodPurchase", Description: "更新进货记录"},
		{ApiGroup: "进货管理", Method: "GET", Path: "/goodPurchase/getGoodPurchaseList", Description: "获取进货记录列表"},
		{ApiGroup: "进货管理", Method: "GET", Path: "/goodPurchase/getGoodPurchaseSummary", Description: "获取进货汇总"},
		// 收款码
		{ApiGroup: "收款码管理", Method: "POST", Path: "/qrcodePayment/createQrcodePayment", Description: "创建收款码"},
		{ApiGroup: "收款码管理", Method: "DELETE", Path: "/qrcodePayment/deleteQrcodePayment", Description: "删除收款码"},
		{ApiGroup: "收款码管理", Method: "PUT", Path: "/qrcodePayment/updateQrcodePayment", Description: "更新收款码"},
		{ApiGroup: "收款码管理", Method: "GET", Path: "/qrcodePayment/getQrcodePaymentList", Description: "获取收款码列表"},
		// 弹窗管理
		{ApiGroup: "弹窗管理", Method: "POST", Path: "/popup/createPopup", Description: "创建弹窗"},
		{ApiGroup: "弹窗管理", Method: "DELETE", Path: "/popup/deletePopup", Description: "删除弹窗"},
		{ApiGroup: "弹窗管理", Method: "PUT", Path: "/popup/updatePopup", Description: "更新弹窗"},
		{ApiGroup: "弹窗管理", Method: "GET", Path: "/popup/getPopupList", Description: "获取弹窗列表"},
		{ApiGroup: "弹窗管理", Method: "GET", Path: "/popup/getPopupPagePathOptions", Description: "获取弹窗页面路径选项"},
		// 营销奖励
		{ApiGroup: "营销奖励", Method: "POST", Path: "/marketingReward/createMarketingReward", Description: "创建营销奖励"},
		{ApiGroup: "营销奖励", Method: "DELETE", Path: "/marketingReward/deleteMarketingReward", Description: "删除营销奖励"},
		{ApiGroup: "营销奖励", Method: "PUT", Path: "/marketingReward/updateMarketingReward", Description: "更新营销奖励"},
		{ApiGroup: "营销奖励", Method: "GET", Path: "/marketingReward/getMarketingRewardList", Description: "获取营销奖励列表"},
		// 数据看板
		{ApiGroup: "数据看板", Method: "GET", Path: "/dashboard/getOverview", Description: "获取数据看板概览"},
		// 预售管理
		{ApiGroup: "预售管理", Method: "GET", Path: "/presale/getPresaleParticipants", Description: "获取预售参与者"},
		// 语言管理
		{ApiGroup: "语言管理", Method: "POST", Path: "/language/createLanguage", Description: "创建语言"},
		{ApiGroup: "语言管理", Method: "DELETE", Path: "/language/deleteLanguage", Description: "删除语言"},
		{ApiGroup: "语言管理", Method: "PUT", Path: "/language/updateLanguage", Description: "更新语言"},
		{ApiGroup: "语言管理", Method: "GET", Path: "/language/getLanguageList", Description: "获取语言列表"},
		{ApiGroup: "语言管理", Method: "POST", Path: "/language/translateI18n", Description: "翻译多语言文本"},
		// 国际区号
		{ApiGroup: "国际区号", Method: "POST", Path: "/phoneAreaCode/createPhoneAreaCode", Description: "创建国际区号"},
		{ApiGroup: "国际区号", Method: "DELETE", Path: "/phoneAreaCode/deletePhoneAreaCode", Description: "删除国际区号"},
		{ApiGroup: "国际区号", Method: "PUT", Path: "/phoneAreaCode/updatePhoneAreaCode", Description: "更新国际区号"},
		{ApiGroup: "国际区号", Method: "GET", Path: "/phoneAreaCode/getPhoneAreaCodeList", Description: "获取国际区号列表"},
		// 签到
		{ApiGroup: "签到管理", Method: "POST", Path: "/signIn/doSignIn", Description: "用户签到"},
		{ApiGroup: "签到管理", Method: "GET", Path: "/signIn/getSignInStatus", Description: "获取签到状态"},
		{ApiGroup: "签到管理", Method: "GET", Path: "/signIn/getSignInRecords", Description: "获取签到记录"},
		{ApiGroup: "签到管理", Method: "GET", Path: "/signIn/getSignInList", Description: "管理端获取签到列表"},
		{ApiGroup: "签到管理", Method: "DELETE", Path: "/signIn/deleteSignIn", Description: "管理端删除签到记录"},
		// 试衣任务
		{ApiGroup: "试衣任务", Method: "POST", Path: "/tryonTask/createTryonTask", Description: "创建试衣任务"},
		{ApiGroup: "试衣任务", Method: "POST", Path: "/tryonTask/applyTryonBeautify", Description: "执行试衣结果智能美肤"},
		{ApiGroup: "试衣任务", Method: "DELETE", Path: "/tryonTask/deleteTryonTask", Description: "删除试衣任务(管理端)"},
		{ApiGroup: "试衣任务", Method: "DELETE", Path: "/tryonTask/deleteTryonTaskByIds", Description: "批量删除试衣任务(管理端)"},
		{ApiGroup: "试衣任务", Method: "GET", Path: "/tryonTask/findTryonTask", Description: "根据ID查询试衣任务"},
		{ApiGroup: "试衣任务", Method: "GET", Path: "/tryonTask/getMyTryonTaskList", Description: "获取我的试衣任务列表"},
		{ApiGroup: "试衣任务", Method: "GET", Path: "/tryonTask/getTryonTaskList", Description: "获取试衣任务列表(管理端)"},
		{ApiGroup: "试衣任务", Method: "GET", Path: "/tryonTask/getTryonTaskStats", Description: "获取试衣任务统计(管理端)"},
		{ApiGroup: "试衣任务", Method: "GET", Path: "/tryonTask/getTryonTaskTrend", Description: "获取试衣任务趋势(管理端)"},
		{ApiGroup: "试衣币记录", Method: "GET", Path: "/cpr/getPointRecordStats", Description: "获取试衣币统计(管理端)"},
		{ApiGroup: "客户端用户", Method: "POST", Path: "/clientUser/adjustTryonPoint", Description: "后台调整用户试衣币"},
		// 试衣币充值订单
		{ApiGroup: "试衣币充值订单", Method: "POST", Path: "/tryonRechargeOrder/createTryonRechargeOrder", Description: "创建试衣币充值订单"},
		{ApiGroup: "试衣币充值订单", Method: "POST", Path: "/tryonRechargeOrder/updateTryonRechargeOrderPayMethod", Description: "更新试衣币充值订单支付方式"},
		{ApiGroup: "试衣币充值订单", Method: "POST", Path: "/tryonRechargeOrder/submitTryonRechargeOrderPayment", Description: "提交试衣币充值订单付款确认"},
		{ApiGroup: "试衣币充值订单", Method: "POST", Path: "/tryonRechargeOrder/cancelTryonRechargeOrder", Description: "取消试衣币充值订单"},
		{ApiGroup: "试衣币充值订单", Method: "POST", Path: "/tryonRechargeOrder/confirmTryonRechargeOrderPayment", Description: "确认试衣币充值订单支付(管理端)"},
		{ApiGroup: "试衣币充值订单", Method: "GET", Path: "/tryonRechargeOrder/selfTryonRechargeOrder", Description: "获取我的试衣币充值订单"},
		{ApiGroup: "试衣币充值订单", Method: "GET", Path: "/tryonRechargeOrder/findTryonRechargeOrder", Description: "获取试衣币充值订单(管理端)"},
		{ApiGroup: "试衣币充值订单", Method: "GET", Path: "/tryonRechargeOrder/getMyTryonRechargeOrderList", Description: "获取我的试衣币充值订单列表"},
		{ApiGroup: "试衣币充值订单", Method: "GET", Path: "/tryonRechargeOrder/getTryonRechargeOrderList", Description: "获取试衣币充值订单列表(管理端)"},
		// 我的模特
		{ApiGroup: "我的模特", Method: "POST", Path: "/tryonModel/createTryonModel", Description: "创建我的模特"},
		{ApiGroup: "我的模特", Method: "PUT", Path: "/tryonModel/updateTryonModel", Description: "重命名我的模特"},
		{ApiGroup: "我的模特", Method: "DELETE", Path: "/tryonModel/deleteTryonModel", Description: "删除我的模特"},
		{ApiGroup: "我的模特", Method: "DELETE", Path: "/tryonModel/deleteTryonModelByIds", Description: "批量删除我的模特"},
		{ApiGroup: "我的模特", Method: "GET", Path: "/tryonModel/getMyTryonModelList", Description: "获取我的模特列表"},
		{ApiGroup: "我的模特", Method: "GET", Path: "/tryonModel/getTryonModelList", Description: "获取模特列表(管理端)"},
		// 我的衣橱
		{ApiGroup: "我的衣橱", Method: "POST", Path: "/tryonCloth/createTryonCloth", Description: "创建我的衣橱"},
		{ApiGroup: "我的衣橱", Method: "PUT", Path: "/tryonCloth/updateTryonCloth", Description: "更新我的衣橱"},
		{ApiGroup: "我的衣橱", Method: "DELETE", Path: "/tryonCloth/deleteTryonCloth", Description: "删除我的衣橱"},
		{ApiGroup: "我的衣橱", Method: "DELETE", Path: "/tryonCloth/deleteTryonClothByIds", Description: "批量删除我的衣橱"},
		{ApiGroup: "我的衣橱", Method: "GET", Path: "/tryonCloth/getMyTryonClothList", Description: "获取我的衣橱列表"},
		{ApiGroup: "我的衣橱", Method: "GET", Path: "/tryonCloth/getTryonClothList", Description: "获取我的衣橱列表(管理端)"},
		// 外部链接域名
		{ApiGroup: "外部链接域名", Method: "POST", Path: "/extDomain/createExternalLinkDomain", Description: "创建外部链接域名"},
		{ApiGroup: "外部链接域名", Method: "DELETE", Path: "/extDomain/deleteExternalLinkDomain", Description: "删除外部链接域名"},
		{ApiGroup: "外部链接域名", Method: "PUT", Path: "/extDomain/updateExternalLinkDomain", Description: "更新外部链接域名"},
		{ApiGroup: "外部链接域名", Method: "GET", Path: "/extDomain/getExternalLinkDomainList", Description: "获取外部链接域名列表"},
		{ApiGroup: "外部链接域名", Method: "POST", Path: "/extDomain/setDefaultDomain", Description: "设置默认域名"},
		// 系统配置扩展
		{ApiGroup: "系统配置", Method: "GET", Path: "/sysConfig/getSysConfigByGroup", Description: "按分组获取配置"},
		{ApiGroup: "系统配置", Method: "GET", Path: "/sysConfig/getSysConfigByKey", Description: "按Key获取配置"},
		{ApiGroup: "系统配置", Method: "GET", Path: "/sysConfig/getPaymentConfig", Description: "获取支付方式配置"},
		{ApiGroup: "系统配置", Method: "GET", Path: "/sysConfig/getUniPreferredPayConfig", Description: "获取uni期望支付方式配置"},
		{ApiGroup: "系统配置", Method: "GET", Path: "/sysConfig/getTryonConfig", Description: "获取试衣配置"},
		{ApiGroup: "系统配置", Method: "GET", Path: "/sysConfig/getAliyunTryonQuotaEstimate", Description: "获取阿里试衣模型额度估算"},
		// 访客统计扩展
		{ApiGroup: "访客统计", Method: "GET", Path: "/visitor/getKefuGuideStats", Description: "获取客服引导漏斗统计"},
		// 数据库巡检
		{ApiGroup: "数据库巡检", Method: "GET", Path: "/dbInspector/getOverview", Description: "获取数据库巡检总览"},
		{ApiGroup: "数据库巡检", Method: "POST", Path: "/dbInspector/autoFix", Description: "自动修复数据库或Redis连接"},
		{ApiGroup: "数据库巡检", Method: "POST", Path: "/dbInspector/deleteRecordsByRange", Description: "按日期范围真删除记录并联动清理文件"},
	}
	for _, api := range apis {
		var count int64
		db.Model(&sysModel.SysApi{}).Where("path = ? AND method = ?", api.Path, api.Method).Count(&count)
		if count == 0 {
			if err := db.Create(&api).Error; err != nil {
				global.GVA_LOG.Error("初始化新模块API失败", zap.Error(err))
			}
		}
	}
}

func initNewModulesMenus(db *gorm.DB) {
	// 查找"商城管理"父菜单
	var shopParent sysModel.SysBaseMenu
	shopParentFound := db.Where("name = ?", "shop").First(&shopParent).Error == nil

	// 查找"客户端"父菜单
	var clientParent sysModel.SysBaseMenu
	clientParentFound := db.Where("name = ?", "client").First(&clientParent).Error == nil

	// 查找"系统工具"父菜单
	var systemToolsParent sysModel.SysBaseMenu
	systemToolsParentFound := db.Where("name = ?", "systemTools").First(&systemToolsParent).Error == nil

	// 查找或创建"英语端"父菜单
	var englishAppParent sysModel.SysBaseMenu
	if err := db.Where("name = ?", "englishApp").First(&englishAppParent).Error; err != nil {
		englishAppParent = sysModel.SysBaseMenu{
			MenuLevel: 0,
			Hidden:    false,
			ParentId:  0,
			Path:      "englishApp",
			Name:      "englishApp",
			Component: "view/index.vue",
			Sort:      2,
			Meta: sysModel.Meta{
				Title: "英语端",
				Icon:  "english",
			},
		}
		if err := db.Create(&englishAppParent).Error; err != nil {
			global.GVA_LOG.Error("创建英语端父菜单失败", zap.Error(err))
		} else {
			// 为新父菜单关联所有角色
			var authorities []sysModel.SysAuthority
			db.Find(&authorities)
			for _, auth := range authorities {
				db.Exec("INSERT IGNORE INTO sys_authority_menus (sys_authority_authority_id, sys_base_menu_id) VALUES (?, ?)",
					fmt.Sprintf("%d", auth.AuthorityId), englishAppParent.ID)
			}
			global.GVA_LOG.Info("英语端父菜单初始化成功")
		}
		// 重新查询获取 ID
		db.Where("name = ?", "englishApp").First(&englishAppParent)
	}

	// 迁移旧子菜单：将 englishLearningWord / englishLearningVideo 从 client 移到 englishApp
	if clientParentFound && englishAppParent.ID > 0 {
		db.Model(&sysModel.SysBaseMenu{}).Where("name = ? AND parent_id = ?", "englishLearningWord", clientParent.ID).Update("parent_id", englishAppParent.ID)
		db.Model(&sysModel.SysBaseMenu{}).Where("name = ? AND parent_id = ?", "englishLearningVideo", clientParent.ID).Update("parent_id", englishAppParent.ID)
	}

	type menuDef struct {
		name      string
		path      string
		component string
		title     string
		icon      string
		parentID  uint
		sort      int
	}

	var menus []menuDef

	if shopParentFound {
		menus = append(menus,
			menuDef{"goodPurchase", "goodPurchase", "view/shop/goodPurchase/goodPurchase.vue", "进货管理", "goods", shopParent.ID, 20},
			menuDef{"qrcodePayment", "qrcodePayment", "view/shop/qrcodePayment/qrcodePayment.vue", "收款码管理", "credit-card", shopParent.ID, 21},
			menuDef{"popup", "popup", "view/shop/popup/popup.vue", "弹窗管理", "message-box", shopParent.ID, 22},
			menuDef{"marketingReward", "marketingReward", "view/shop/marketingReward/marketingReward.vue", "营销奖励", "present", shopParent.ID, 23},
			menuDef{"shopDashboard", "shopDashboard", "view/shop/dashboard/dashboard.vue", "数据看板", "data-analysis", shopParent.ID, 1},
			menuDef{"presale", "presale", "view/shop/presale/presale.vue", "预售管理", "clock", shopParent.ID, 24},
			menuDef{"couponorderuser", "couponorderuser", "view/shop/couponorderuser/couponorderuser.vue", "领券记录", "ticket", shopParent.ID, 25},
			menuDef{"skuSpec", "skuSpec", "view/shop/skuSpec/skuSpec.vue", "规格字典", "notebook", shopParent.ID, 26},
		)
	}

	if clientParentFound {
		menus = append(menus,
			menuDef{"language", "language", "view/client/language/language.vue", "语言管理", "edit", clientParent.ID, 13},
			menuDef{"phoneAreaCode", "phoneAreaCode", "view/client/phoneAreaCode/phoneAreaCode.vue", "国际区号", "phone", clientParent.ID, 14},
			menuDef{"externalLinkDomain", "externalLinkDomain", "view/client/externalLinkDomain/externalLinkDomain.vue", "外部链接域名", "link", clientParent.ID, 15},
			menuDef{"signInManage", "signInManage", "view/client/signIn/signIn.vue", "签到管理", "calendar", clientParent.ID, 16},
			menuDef{"tryonTaskManage", "tryonTaskManage", "view/client/tryonTask/tryonTask.vue", "试衣任务", "camera", clientParent.ID, 17},
			menuDef{"tryonModelManage", "tryonModelManage", "view/client/tryonModel/tryonModel.vue", "我的模特管理", "avatar", clientParent.ID, 18},
			menuDef{"tryonClothManage", "tryonClothManage", "view/client/tryonCloth/tryonCloth.vue", "我的衣橱管理", "goods", clientParent.ID, 19},
			menuDef{"tryonPointRecord", "tryonPointRecord", "view/client/tryonPointRecord/tryonPointRecord.vue", "试衣币记录", "coin", clientParent.ID, 20},
			menuDef{"tryonRechargeOrder", "tryonRechargeOrder", "view/client/tryonRechargeOrder/tryonRechargeOrder.vue", "试衣币充值订单", "wallet", clientParent.ID, 21},
		)
		// 英语学习相关菜单已移至 englishApp 父菜单下，此处移除不再重复创建
	}

	if englishAppParent.ID > 0 {
		menus = append(menus,
			menuDef{"englishLearningWord", "englishLearningWord", "plugin/english_learning/view/word.vue", "英语单词管理", "reading", englishAppParent.ID, 1},
			menuDef{"englishLearningVideo", "englishLearningVideo", "plugin/english_learning/view/video.vue", "英语视频字幕", "video-play", englishAppParent.ID, 2},
			menuDef{"englishCheckinRecord", "englishCheckinRecord", "view/client/englishCheckinRecord/englishCheckinRecord.vue", "签到记录", "calendar", englishAppParent.ID, 3},
			menuDef{"englishCollections", "englishCollections", "view/client/englishCollections/englishCollections.vue", "英语收藏", "star-on", englishAppParent.ID, 4},
			menuDef{"englishErrorLog", "englishErrorLog", "view/client/englishErrorLog/englishErrorLog.vue", "错词本", "document-copy", englishAppParent.ID, 5},
			menuDef{"englishWatchHistory", "englishWatchHistory", "view/client/englishWatchHistory/englishWatchHistory.vue", "观看历史", "time", englishAppParent.ID, 6},
			menuDef{"englishPointHistory", "englishPointHistory", "view/client/englishPointHistory/englishPointHistory.vue", "英语积分", "coin", englishAppParent.ID, 7},
			menuDef{"englishFreeTimeHistory", "englishFreeTimeHistory", "view/client/englishFreeTimeHistory/englishFreeTimeHistory.vue", "时长明细", "timer", englishAppParent.ID, 8},
			menuDef{"englishVideoTags", "englishVideoTags", "view/client/videoTag/videoTag.vue", "视频标签", "price-tag", englishAppParent.ID, 9},
		)
	}

	if systemToolsParentFound {
		menus = append(menus,
			menuDef{"dbInspector", "dbInspector", "view/systemTools/dbInspector/index.vue", "数据库巡检与清理", "cpu", systemToolsParent.ID, 10},
		)
	}

	var authorities []sysModel.SysAuthority
	db.Find(&authorities)

	for _, m := range menus {
		var count int64
		db.Model(&sysModel.SysBaseMenu{}).Where("name = ?", m.name).Count(&count)
		if count > 0 {
			continue
		}
		menu := sysModel.SysBaseMenu{
			MenuLevel: 1,
			Hidden:    false,
			ParentId:  m.parentID,
			Path:      m.path,
			Name:      m.name,
			Component: m.component,
			Sort:      m.sort,
			Meta: sysModel.Meta{
				Title: m.title,
				Icon:  m.icon,
			},
		}
		if err := db.Create(&menu).Error; err != nil {
			global.GVA_LOG.Error("初始化菜单失败: "+m.name, zap.Error(err))
			continue
		}
		for _, auth := range authorities {
			db.Exec("INSERT IGNORE INTO sys_authority_menus (sys_authority_authority_id, sys_base_menu_id) VALUES (?, ?)",
				fmt.Sprintf("%d", auth.AuthorityId), menu.ID)
		}
		global.GVA_LOG.Info(m.title + "菜单初始化成功")
	}
}

func initNewModulesMenuBtns(db *gorm.DB) {
	type btnDef struct {
		Name string
		Desc string
	}

	menuBtnDefs := map[string][]btnDef{
		"popup": {
			{Name: "createPopup", Desc: "创建弹窗"},
			{Name: "deletePopup", Desc: "删除弹窗"},
			{Name: "updatePopup", Desc: "更新弹窗"},
			{Name: "getPopupList", Desc: "获取弹窗列表"},
			{Name: "getPopupPagePathOptions", Desc: "获取弹窗页面路径选项"},
		},
	}

	for menuName, defs := range menuBtnDefs {
		var menu sysModel.SysBaseMenu
		if err := db.Select("id", "name").Where("name = ?", menuName).First(&menu).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				global.GVA_LOG.Error("查询菜单失败: "+menuName, zap.Error(err))
			}
			continue
		}

		for _, btn := range defs {
			var count int64
			db.Model(&sysModel.SysBaseMenuBtn{}).
				Where("sys_base_menu_id = ? AND name = ?", menu.ID, btn.Name).
				Count(&count)
			if count > 0 {
				continue
			}

			row := sysModel.SysBaseMenuBtn{
				Name:          btn.Name,
				Desc:          btn.Desc,
				SysBaseMenuID: menu.ID,
			}
			if err := db.Create(&row).Error; err != nil {
				global.GVA_LOG.Error("初始化菜单按钮失败: "+menuName+"/"+btn.Name, zap.Error(err))
			}
		}
	}
}

func initNewModulesCasbin(db *gorm.DB) {
	var authorities []sysModel.SysAuthority
	db.Find(&authorities)

	seedAuthorities := make([]sysModel.SysAuthority, 0, len(authorities))
	for _, auth := range authorities {
		if auth.AuthorityId == 888 || auth.AuthorityId == 8881 || auth.AuthorityId == 8080 || auth.AuthorityId == 9528 {
			seedAuthorities = append(seedAuthorities, auth)
		}
	}

	paths := []struct {
		Path   string
		Method string
	}{
		// 进货管理
		{"/goodPurchase/createGoodPurchase", "POST"},
		{"/goodPurchase/deleteGoodPurchase", "DELETE"},
		{"/goodPurchase/updateGoodPurchase", "PUT"},
		{"/goodPurchase/getGoodPurchaseList", "GET"},
		{"/goodPurchase/getGoodPurchaseSummary", "GET"},
		// 收款码
		{"/qrcodePayment/createQrcodePayment", "POST"},
		{"/qrcodePayment/deleteQrcodePayment", "DELETE"},
		{"/qrcodePayment/updateQrcodePayment", "PUT"},
		{"/qrcodePayment/getQrcodePaymentList", "GET"},
		// 弹窗
		{"/popup/createPopup", "POST"},
		{"/popup/deletePopup", "DELETE"},
		{"/popup/updatePopup", "PUT"},
		{"/popup/getPopupList", "GET"},
		{"/popup/getPopupPagePathOptions", "GET"},
		// 营销奖励
		{"/marketingReward/createMarketingReward", "POST"},
		{"/marketingReward/deleteMarketingReward", "DELETE"},
		{"/marketingReward/updateMarketingReward", "PUT"},
		{"/marketingReward/getMarketingRewardList", "GET"},
		// 数据看板
		{"/dashboard/getOverview", "GET"},
		// 预售
		{"/presale/getPresaleParticipants", "GET"},
		// 语言
		{"/language/createLanguage", "POST"},
		{"/language/deleteLanguage", "DELETE"},
		{"/language/updateLanguage", "PUT"},
		{"/language/getLanguageList", "GET"},
		{"/language/translateI18n", "POST"},
		// 区号
		{"/phoneAreaCode/createPhoneAreaCode", "POST"},
		{"/phoneAreaCode/deletePhoneAreaCode", "DELETE"},
		{"/phoneAreaCode/updatePhoneAreaCode", "PUT"},
		{"/phoneAreaCode/getPhoneAreaCodeList", "GET"},
		// 签到
		{"/signIn/doSignIn", "POST"},
		{"/signIn/getSignInStatus", "GET"},
		{"/signIn/getSignInRecords", "GET"},
		{"/signIn/getSignInList", "GET"},
		{"/signIn/deleteSignIn", "DELETE"},
		// 试衣任务
		{"/tryonTask/createTryonTask", "POST"},
		{"/tryonTask/applyTryonBeautify", "POST"},
		{"/tryonTask/deleteTryonTask", "DELETE"},
		{"/tryonTask/deleteTryonTaskByIds", "DELETE"},
		{"/tryonTask/findTryonTask", "GET"},
		{"/tryonTask/getMyTryonTaskList", "GET"},
		{"/tryonTask/getTryonTaskList", "GET"},
		{"/tryonTask/getTryonTaskStats", "GET"},
		{"/tryonTask/getTryonTaskTrend", "GET"},
		{"/cpr/getPointRecordStats", "GET"},
		{"/clientUser/adjustTryonPoint", "POST"},
		// 试衣币充值订单
		{"/tryonRechargeOrder/createTryonRechargeOrder", "POST"},
		{"/tryonRechargeOrder/updateTryonRechargeOrderPayMethod", "POST"},
		{"/tryonRechargeOrder/submitTryonRechargeOrderPayment", "POST"},
		{"/tryonRechargeOrder/cancelTryonRechargeOrder", "POST"},
		{"/tryonRechargeOrder/confirmTryonRechargeOrderPayment", "POST"},
		{"/tryonRechargeOrder/selfTryonRechargeOrder", "GET"},
		{"/tryonRechargeOrder/findTryonRechargeOrder", "GET"},
		{"/tryonRechargeOrder/getMyTryonRechargeOrderList", "GET"},
		{"/tryonRechargeOrder/getTryonRechargeOrderList", "GET"},
		// 我的模特
		{"/tryonModel/createTryonModel", "POST"},
		{"/tryonModel/updateTryonModel", "PUT"},
		{"/tryonModel/deleteTryonModel", "DELETE"},
		{"/tryonModel/deleteTryonModelByIds", "DELETE"},
		{"/tryonModel/getMyTryonModelList", "GET"},
		{"/tryonModel/getTryonModelList", "GET"},
		// 我的衣橱
		{"/tryonCloth/createTryonCloth", "POST"},
		{"/tryonCloth/updateTryonCloth", "PUT"},
		{"/tryonCloth/deleteTryonCloth", "DELETE"},
		{"/tryonCloth/deleteTryonClothByIds", "DELETE"},
		{"/tryonCloth/getMyTryonClothList", "GET"},
		{"/tryonCloth/getTryonClothList", "GET"},
		// 外部链接域名
		{"/extDomain/createExternalLinkDomain", "POST"},
		{"/extDomain/deleteExternalLinkDomain", "DELETE"},
		{"/extDomain/updateExternalLinkDomain", "PUT"},
		{"/extDomain/getExternalLinkDomainList", "GET"},
		{"/extDomain/setDefaultDomain", "POST"},
		// 系统配置扩展
		{"/sysConfig/getSysConfigByGroup", "GET"},
		{"/sysConfig/getSysConfigByKey", "GET"},
		{"/sysConfig/getPaymentConfig", "GET"},
		{"/sysConfig/getUniPreferredPayConfig", "GET"},
		// 访客统计扩展
		{"/visitor/getKefuGuideStats", "GET"},
		// 数据库巡检
		{"/dbInspector/getOverview", "GET"},
		{"/dbInspector/autoFix", "POST"},
		{"/dbInspector/deleteRecordsByRange", "POST"},
		// 视频标签管理（englishVideoTags / englishLearningVideo 共用）
		{"/videoTag/getVideoTagList", "GET"},
		{"/videoTag/findVideoTag", "GET"},
		{"/videoTag/createVideoTag", "POST"},
		{"/videoTag/deleteVideoTag", "DELETE"},
		{"/videoTag/deleteVideoTagByIds", "DELETE"},
		{"/videoTag/updateVideoTag", "PUT"},
		// 英语学习管理端
		{"/englishLearning/admin/getCheckinRecordList", "GET"},
		{"/englishLearning/admin/getPointRecordList", "GET"},
		{"/englishLearning/admin/getFreeTimeRecordList", "GET"},
		{"/englishLearning/admin/getWatchHistoryList", "GET"},
		{"/englishLearning/admin/getCollectionList", "GET"},
		{"/englishLearning/admin/getWordErrorLogList", "GET"},
		{"/englishLearning/admin/getUserList", "GET"},
		{"/englishLearning/content/getVideoEpisodeListByTag", "GET"},
	}

	for _, auth := range seedAuthorities {
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

	cleanupLegacyManagedCasbinRules(db, "new_modules", paths, []string{"888", "8881", "8080", "9528"})

	adminOnlyPaths := []struct {
		Path   string
		Method string
	}{
		{"/goodPurchase/createGoodPurchase", "POST"},
		{"/goodPurchase/deleteGoodPurchase", "DELETE"},
		{"/goodPurchase/updateGoodPurchase", "PUT"},
		{"/goodPurchase/getGoodPurchaseList", "GET"},
		{"/goodPurchase/getGoodPurchaseSummary", "GET"},
		{"/qrcodePayment/createQrcodePayment", "POST"},
		{"/qrcodePayment/deleteQrcodePayment", "DELETE"},
		{"/qrcodePayment/updateQrcodePayment", "PUT"},
		{"/qrcodePayment/getQrcodePaymentList", "GET"},
		{"/popup/createPopup", "POST"},
		{"/popup/deletePopup", "DELETE"},
		{"/popup/updatePopup", "PUT"},
		{"/popup/getPopupList", "GET"},
		{"/popup/getPopupPagePathOptions", "GET"},
		{"/marketingReward/createMarketingReward", "POST"},
		{"/marketingReward/deleteMarketingReward", "DELETE"},
		{"/marketingReward/updateMarketingReward", "PUT"},
		{"/marketingReward/getMarketingRewardList", "GET"},
		{"/dashboard/getOverview", "GET"},
		{"/presale/getPresaleParticipants", "GET"},
		{"/language/createLanguage", "POST"},
		{"/language/deleteLanguage", "DELETE"},
		{"/language/updateLanguage", "PUT"},
		{"/phoneAreaCode/createPhoneAreaCode", "POST"},
		{"/phoneAreaCode/deletePhoneAreaCode", "DELETE"},
		{"/phoneAreaCode/updatePhoneAreaCode", "PUT"},
		{"/phoneAreaCode/getPhoneAreaCodeList", "GET"},
		{"/signIn/getSignInList", "GET"},
		{"/signIn/deleteSignIn", "DELETE"},
		{"/tryonTask/deleteTryonTask", "DELETE"},
		{"/tryonTask/deleteTryonTaskByIds", "DELETE"},
		{"/tryonTask/getTryonTaskList", "GET"},
		{"/tryonTask/getTryonTaskStats", "GET"},
		{"/tryonTask/getTryonTaskTrend", "GET"},
		{"/extDomain/createExternalLinkDomain", "POST"},
		{"/extDomain/deleteExternalLinkDomain", "DELETE"},
		{"/extDomain/updateExternalLinkDomain", "PUT"},
		{"/extDomain/getExternalLinkDomainList", "GET"},
		{"/extDomain/setDefaultDomain", "POST"},
		{"/sysConfig/getSysConfigList", "GET"},
		{"/sysConfig/getSysConfigByGroup", "GET"},
		{"/sysConfig/getSysConfigByKey", "GET"},
		{"/sysConfig/getAliyunTryonQuotaEstimate", "GET"},
		{"/sysConfig/updateSysConfig", "PUT"},
		{"/dbInspector/getOverview", "GET"},
		{"/dbInspector/autoFix", "POST"},
		{"/dbInspector/deleteRecordsByRange", "POST"},
		{"/clientUser/adjustTryonPoint", "POST"},
		{"/tryonRechargeOrder/confirmTryonRechargeOrderPayment", "POST"},
		{"/tryonRechargeOrder/findTryonRechargeOrder", "GET"},
		{"/tryonRechargeOrder/getTryonRechargeOrderList", "GET"},
		{"/visitor/getKefuGuideStats", "GET"},
	}
	for _, p := range adminOnlyPaths {
		db.Exec(
			"DELETE FROM casbin_rule WHERE ptype = ? AND v1 = ? AND v2 = ? AND v0 NOT IN (?, ?)",
			"p", p.Path, p.Method, "888", "8881",
		)
	}

	grantEnglishLearningOpsCasbin(db)
}

func cleanupLegacyManagedCasbinRules(db *gorm.DB, source string, paths []struct {
	Path   string
	Method string
}, keepAuthorities []string) {
	if db == nil {
		return
	}

	source = strings.TrimSpace(source)
	if source == "" {
		return
	}

	allowedRuleSet := make(map[string]struct{}, len(paths))
	for _, p := range paths {
		path := strings.TrimSpace(p.Path)
		method := strings.ToUpper(strings.TrimSpace(p.Method))
		if path == "" || method == "" {
			continue
		}
		allowedRuleSet[path+"#"+method] = struct{}{}
	}

	keepAuthoritySet := make(map[string]struct{}, len(keepAuthorities))
	for _, authority := range keepAuthorities {
		authority = strings.TrimSpace(authority)
		if authority == "" {
			continue
		}
		keepAuthoritySet[authority] = struct{}{}
	}

	type casbinRuleRow struct {
		V0 string `gorm:"column:v0"`
		V1 string `gorm:"column:v1"`
		V2 string `gorm:"column:v2"`
	}

	rules := make([]casbinRuleRow, 0)
	if err := db.Table("casbin_rule").
		Select("v0", "v1", "v2").
		Where("ptype = ? AND v3 = ?", "p", source).
		Scan(&rules).Error; err != nil {
		global.GVA_LOG.Warn("查询历史托管casbin规则失败", zap.Error(err), zap.String("source", source))
		return
	}

	if len(rules) == 0 {
		return
	}

	for _, rule := range rules {
		authority := strings.TrimSpace(rule.V0)
		ruleKey := strings.TrimSpace(rule.V1) + "#" + strings.ToUpper(strings.TrimSpace(rule.V2))

		if len(keepAuthoritySet) > 0 {
			if _, keep := keepAuthoritySet[authority]; !keep {
				if err := db.Table("casbin_rule").Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ? AND v3 = ?", "p", authority, rule.V1, rule.V2, source).Delete(nil).Error; err != nil {
					global.GVA_LOG.Warn("清理历史托管casbin规则失败", zap.Error(err), zap.String("authority", authority), zap.String("path", rule.V1), zap.String("method", rule.V2))
				}
				continue
			}
		}

		if _, allowed := allowedRuleSet[ruleKey]; allowed {
			continue
		}

		if err := db.Table("casbin_rule").Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ? AND v3 = ?", "p", authority, rule.V1, rule.V2, source).Delete(nil).Error; err != nil {
			global.GVA_LOG.Warn("清理过期托管casbin规则失败", zap.Error(err), zap.String("authority", authority), zap.String("path", rule.V1), zap.String("method", rule.V2))
		}
	}
}

func grantEnglishLearningOpsCasbin(db *gorm.DB) {
	if db == nil {
		return
	}

	type authorityMenuRow struct {
		AuthorityID uint `gorm:"column:authority_id"`
	}

	var authorityRows []authorityMenuRow
	if err := db.
		Table("sys_authority_menus sam").
		Select("DISTINCT sam.sys_authority_authority_id AS authority_id").
		Joins("JOIN sys_base_menus sbm ON sbm.id = sam.sys_base_menu_id").
		Where("sbm.name IN ?", []string{"englishLearningWord", "englishLearningVideo", "englishVideoTags"}).
		Scan(&authorityRows).Error; err != nil {
		global.GVA_LOG.Warn("查询英语学习菜单角色失败", zap.Error(err))
		return
	}

	if len(authorityRows) == 0 {
		return
	}

	dependentPaths := []struct {
		Path   string
		Method string
	}{
		{Path: "/language/getLanguageList", Method: "GET"},
		{Path: "/language/translateI18n", Method: "POST"},
		{Path: "/fileUploadAndDownload/upload", Method: "POST"},
		// 视频标签（englishVideoTags 管理页 / englishLearningVideo 剧集编辑共用）
		{Path: "/videoTag/getVideoTagList", Method: "GET"},
		{Path: "/videoTag/findVideoTag", Method: "GET"},
		{Path: "/videoTag/createVideoTag", Method: "POST"},
		{Path: "/videoTag/deleteVideoTag", Method: "DELETE"},
		{Path: "/videoTag/deleteVideoTagByIds", Method: "DELETE"},
		{Path: "/videoTag/updateVideoTag", Method: "PUT"},
	}

	for _, row := range authorityRows {
		if row.AuthorityID == 0 {
			continue
		}

		authorityID := fmt.Sprintf("%d", row.AuthorityID)
		for _, p := range dependentPaths {
			var count int64
			if err := db.Table("casbin_rule").Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?", "p", authorityID, p.Path, p.Method).Count(&count).Error; err != nil {
				global.GVA_LOG.Warn("查询英语学习依赖权限失败", zap.Error(err), zap.String("authorityId", authorityID), zap.String("path", p.Path), zap.String("method", p.Method))
				continue
			}
			if count > 0 {
				continue
			}

			if err := db.Exec("INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES (?, ?, ?, ?)", "p", authorityID, p.Path, p.Method).Error; err != nil {
				global.GVA_LOG.Warn("补齐英语学习依赖权限失败", zap.Error(err), zap.String("authorityId", authorityID), zap.String("path", p.Path), zap.String("method", p.Method))
			}
		}
	}
}

// initLanguageSeedData 初始化默认语言数据（幂等，缺失语言自动补齐）
func initLanguageSeedData(db *gorm.DB) {
	boolTrue := true
	boolFalse := false
	sort := func(n int) *int { return &n }

	langs := []clientModel.SysLanguage{
		{Code: "mn", Name: "蒙古国语", NativeName: "Монгол", IsEnabled: &boolTrue, IsDefault: &boolTrue, Sort: sort(0)},
		{Code: "zh", Name: "中文", NativeName: "中文", IsEnabled: &boolTrue, IsDefault: &boolFalse, Sort: sort(1)},
		{Code: "en", Name: "英语", NativeName: "English", IsEnabled: &boolTrue, IsDefault: &boolFalse, Sort: sort(2)},
		{Code: "th", Name: "泰文", NativeName: "ไทย", IsEnabled: &boolTrue, IsDefault: &boolFalse, Sort: sort(3)},
		{Code: "hi", Name: "印度语", NativeName: "हिन्दी", IsEnabled: &boolTrue, IsDefault: &boolFalse, Sort: sort(4)},
		{Code: "zh-TW", Name: "中文繁体", NativeName: "繁體中文", IsEnabled: &boolTrue, IsDefault: &boolFalse, Sort: sort(5)},
		{Code: "id", Name: "印度尼西亚语", NativeName: "Bahasa Indonesia", IsEnabled: &boolTrue, IsDefault: &boolFalse, Sort: sort(6)},
		{Code: "vi", Name: "越南语", NativeName: "Tiếng Việt", IsEnabled: &boolTrue, IsDefault: &boolFalse, Sort: sort(7)},
		{Code: "ar", Name: "阿拉伯语", NativeName: "العربية", IsEnabled: &boolTrue, IsDefault: &boolFalse, Sort: sort(8)},
		{Code: "ja", Name: "日语", NativeName: "日本語", IsEnabled: &boolTrue, IsDefault: &boolFalse, Sort: sort(9)},
		{Code: "ko", Name: "韩语", NativeName: "한국어", IsEnabled: &boolTrue, IsDefault: &boolFalse, Sort: sort(10)},
		{Code: "ms", Name: "马来语", NativeName: "Bahasa Melayu", IsEnabled: &boolTrue, IsDefault: &boolFalse, Sort: sort(11)},
	}

	existingRows := make([]clientModel.SysLanguage, 0)
	if err := db.Model(&clientModel.SysLanguage{}).Select("code").Find(&existingRows).Error; err != nil {
		global.GVA_LOG.Error("查询已存在语言失败", zap.Error(err))
		return
	}

	existingCodes := make(map[string]struct{}, len(existingRows))
	for _, row := range existingRows {
		code := strings.ToLower(strings.TrimSpace(row.Code))
		if code == "" {
			continue
		}
		existingCodes[code] = struct{}{}
	}

	insertedCount := 0

	for _, lang := range langs {
		code := strings.ToLower(strings.TrimSpace(lang.Code))
		if _, exists := existingCodes[code]; exists {
			continue
		}
		if err := db.Create(&lang).Error; err != nil {
			global.GVA_LOG.Error("初始化语言数据失败: "+lang.Code, zap.Error(err))
			continue
		}
		insertedCount++
		existingCodes[code] = struct{}{}
	}

	if insertedCount > 0 {
		global.GVA_LOG.Info(fmt.Sprintf("语言数据初始化完成，新增 %d 种语言", insertedCount))
	} else {
		global.GVA_LOG.Info("语言数据初始化完成，未发现缺失语言")
	}
}
