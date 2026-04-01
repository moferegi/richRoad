package initialize

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
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
	initNewModulesCasbin(db)
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
		// 国际区号
		{ApiGroup: "国际区号", Method: "POST", Path: "/phoneAreaCode/createPhoneAreaCode", Description: "创建国际区号"},
		{ApiGroup: "国际区号", Method: "DELETE", Path: "/phoneAreaCode/deletePhoneAreaCode", Description: "删除国际区号"},
		{ApiGroup: "国际区号", Method: "PUT", Path: "/phoneAreaCode/updatePhoneAreaCode", Description: "更新国际区号"},
		{ApiGroup: "国际区号", Method: "GET", Path: "/phoneAreaCode/getPhoneAreaCodeList", Description: "获取国际区号列表"},
		// 签到
		{ApiGroup: "签到管理", Method: "POST", Path: "/signIn/doSignIn", Description: "用户签到"},
		{ApiGroup: "签到管理", Method: "GET", Path: "/signIn/getSignInStatus", Description: "获取签到状态"},
		{ApiGroup: "签到管理", Method: "GET", Path: "/signIn/getSignInRecords", Description: "获取签到记录"},
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
			menuDef{"dashboard", "dashboard", "view/shop/dashboard/dashboard.vue", "数据看板", "data-analysis", shopParent.ID, 1},
			menuDef{"presale", "presale", "view/shop/presale/presale.vue", "预售管理", "clock", shopParent.ID, 24},
		)
	}

	if clientParentFound {
		menus = append(menus,
			menuDef{"language", "language", "view/client/language/language.vue", "语言管理", "edit", clientParent.ID, 13},
			menuDef{"phoneAreaCode", "phoneAreaCode", "view/client/phoneAreaCode/phoneAreaCode.vue", "国际区号", "phone", clientParent.ID, 14},
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

func initNewModulesCasbin(db *gorm.DB) {
	var authorities []sysModel.SysAuthority
	db.Find(&authorities)

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
		// 区号
		{"/phoneAreaCode/createPhoneAreaCode", "POST"},
		{"/phoneAreaCode/deletePhoneAreaCode", "DELETE"},
		{"/phoneAreaCode/updatePhoneAreaCode", "PUT"},
		{"/phoneAreaCode/getPhoneAreaCodeList", "GET"},
		// 签到
		{"/signIn/doSignIn", "POST"},
		{"/signIn/getSignInStatus", "GET"},
		{"/signIn/getSignInRecords", "GET"},
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
