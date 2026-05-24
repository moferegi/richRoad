package initialize

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitSecurityData 注册IP封禁管理API、菜单、Casbin权限，并加载封禁IP缓存
func InitSecurityData() {
	if global.GVA_DB == nil {
		return
	}
	db := global.GVA_DB

	initSecurityApis(db)
	initSecurityMenu(db)
	initSecurityCasbin(db)
	cleanupCasbinManagementOvergrant(db)
	cleanupSystemSensitiveOvergrant(db)
	cleanupSensitiveIgnoreApis(db)
	ensureClientSelfProfilePermission(db)
	ensureSystemReloadPermission(db)
	loadBannedIPCache()
}

func cleanupSensitiveIgnoreApis(db *gorm.DB) {
	if db == nil || !db.Migrator().HasTable(&sysModel.SysIgnoreApi{}) {
		return
	}

	prefixes := []string{
		"/address/",
		"/api/",
		"/autoCode/",
		"/authority/",
		"/banner/",
		"/cart/",
		"/category/",
		"/authorityBtn/",
		"/casbin/",
		"/clientUser/",
		"/collect/",
		"/comment/",
		"/coupon/",
		"/cou/",
		"/cpr/",
		"/customer/",
		"/extDomain/",
		"/good/",
		"/goodPurchase/",
		"/kefu/",
		"/language/",
		"/marketingReward/",
		"/menu/",
		"/order/",
		"/phoneAreaCode/",
		"/presale/",
		"/promo/",
		"/qrcodePayment/",
		"/signIn/",
		"/sku/",
		"/skuSpec/",
		"/sysBannedIP/",
		"/sysConfig/",
		"/sysDictionary/",
		"/sysDictionaryDetail/",
		"/sysExportTemplate/",
		"/sysLoginLog/",
		"/sysOperationRecord/",
		"/sysParams/",
		"/sysVersion/",
		"/tag/",
		"/tryonCloth/",
		"/tryonModel/",
		"/tryonRechargeOrder/",
		"/tryonTask/",
		"/visitor/",
	}

	targets := []struct {
		Method string
		Path   string
	}{
		{Method: "POST", Path: "/init/initdb"},
		{Method: "POST", Path: "/init/checkdb"},
		{Method: "POST", Path: "/system/reloadSystem"},
		{Method: "GET", Path: "/info/getInfoDataSource"},
		{Method: "POST", Path: "/api/ignoreApi"},
		{Method: "POST", Path: "/api/enterSyncApi"},
		{Method: "POST", Path: "/api/createApi"},
		{Method: "POST", Path: "/api/deleteApi"},
		{Method: "DELETE", Path: "/api/deleteApisByIds"},
		{Method: "POST", Path: "/api/updateApi"},
		{Method: "POST", Path: "/user/admin_register"},
		{Method: "POST", Path: "/user/resetPassword"},
		{Method: "POST", Path: "/user/setUserAuthority"},
		{Method: "POST", Path: "/user/setUserAuthorities"},
		{Method: "DELETE", Path: "/user/deleteUser"},
		{Method: "PUT", Path: "/user/setUserInfo"},
		{Method: "POST", Path: "/autoCode/createTemp"},
		{Method: "POST", Path: "/autoCode/installPlugin"},
		{Method: "POST", Path: "/autoCode/removePlugin"},
		{Method: "POST", Path: "/autoCode/createPackage"},
		{Method: "POST", Path: "/autoCode/delPackage"},
		{Method: "POST", Path: "/fileUploadAndDownload/createScanUploadTicket"},
		{Method: "POST", Path: "/fileUploadAndDownload/deleteFile"},
		{Method: "POST", Path: "/fileUploadAndDownload/editFileName"},
		{Method: "POST", Path: "/fileUploadAndDownload/getFileList"},
		{Method: "POST", Path: "/fileUploadAndDownload/importURL"},
		{Method: "GET", Path: "/fileUploadAndDownload/listFolders"},
	}

	var affected int64
	for _, prefix := range prefixes {
		res := db.Where("path LIKE ?", prefix+"%").Delete(&sysModel.SysIgnoreApi{})
		if res.Error != nil {
			global.GVA_LOG.Warn("按前缀清理敏感忽略API失败", zap.String("prefix", prefix), zap.Error(res.Error))
			continue
		}
		affected += res.RowsAffected
	}
	for _, target := range targets {
		res := db.Where("method = ? AND path = ?", target.Method, target.Path).Delete(&sysModel.SysIgnoreApi{})
		if res.Error != nil {
			global.GVA_LOG.Warn("清理敏感忽略API失败", zap.String("path", target.Path), zap.String("method", target.Method), zap.Error(res.Error))
			continue
		}
		affected += res.RowsAffected
	}

	if affected > 0 {
		global.GVA_LOG.Warn("已清理敏感接口的忽略规则", zap.Int64("affected", affected))
	}
}

func ensureClientSelfProfilePermission(db *gorm.DB) {
	if db == nil {
		return
	}

	const path = "/clientUser/setClientUserInfo"
	const method = "POST"

	if db.Migrator().HasTable(&sysModel.SysApi{}) {
		var count int64
		db.Model(&sysModel.SysApi{}).Where("path = ? AND method = ?", path, method).Count(&count)
		if count == 0 {
			api := sysModel.SysApi{ApiGroup: "客户端用户", Method: method, Path: path, Description: "设置个人信息"}
			if err := db.Create(&api).Error; err != nil {
				global.GVA_LOG.Warn("补齐客户端个人信息API失败", zap.String("path", path), zap.Error(err))
			}
		}
	}

	if db.Migrator().HasTable("casbin_rule") {
		var count int64
		db.Table("casbin_rule").Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?", "p", "8080", path, method).Count(&count)
		if count == 0 {
			res := db.Exec("INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES (?, ?, ?, ?)", "p", "8080", path, method)
			if res.Error != nil {
				global.GVA_LOG.Warn("补齐客户端个人信息Casbin规则失败", zap.String("path", path), zap.Error(res.Error))
			}
		}
	}
}

func ensureSystemReloadPermission(db *gorm.DB) {
	if db == nil {
		return
	}

	const path = "/system/reloadSystem"
	const method = "POST"

	if db.Migrator().HasTable(&sysModel.SysApi{}) {
		var count int64
		db.Model(&sysModel.SysApi{}).Where("path = ? AND method = ?", path, method).Count(&count)
		if count == 0 {
			api := sysModel.SysApi{ApiGroup: "系统服务", Method: method, Path: path, Description: "重载系统"}
			if err := db.Create(&api).Error; err != nil {
				global.GVA_LOG.Warn("补齐系统重载API失败", zap.String("path", path), zap.Error(err))
			}
		}
	}

	if db.Migrator().HasTable("casbin_rule") {
		var count int64
		db.Table("casbin_rule").Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?", "p", "888", path, method).Count(&count)
		if count == 0 {
			res := db.Exec("INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES (?, ?, ?, ?)", "p", "888", path, method)
			if res.Error != nil {
				global.GVA_LOG.Warn("补齐系统重载Casbin规则失败", zap.String("path", path), zap.Error(res.Error))
			}
		}
	}

	cleanupLegacyManagedCasbinRules(db, "system_reload", []struct {
		Path   string
		Method string
	}{
		{Path: path, Method: method},
	}, []string{"888"})
}

func initSecurityApis(db *gorm.DB) {
	apis := []sysModel.SysApi{
		{ApiGroup: "安全管理", Method: "POST", Path: "/sysBannedIP/banIP", Description: "封禁IP"},
		{ApiGroup: "安全管理", Method: "POST", Path: "/sysBannedIP/unbanIP", Description: "解封IP"},
		{ApiGroup: "安全管理", Method: "GET", Path: "/sysBannedIP/getBannedIPList", Description: "获取封禁IP列表"},
		{ApiGroup: "安全管理", Method: "GET", Path: "/sysBannedIP/getAttackStats", Description: "获取IP攻击统计"},
	}
	for _, api := range apis {
		var count int64
		db.Model(&sysModel.SysApi{}).Where("path = ? AND method = ?", api.Path, api.Method).Count(&count)
		if count == 0 {
			if err := db.Create(&api).Error; err != nil {
				global.GVA_LOG.Error("注册安全管理API失败", zap.String("path", api.Path), zap.Error(err))
			}
		}
	}
}

func initSecurityMenu(db *gorm.DB) {
	// 幂等：已存在则跳过
	var count int64
	db.Model(&sysModel.SysBaseMenu{}).Where("name = ?", "ipBan").Count(&count)
	if count > 0 {
		return
	}

	// 找到 superAdmin 父菜单
	var parent sysModel.SysBaseMenu
	if err := db.Where("name = ?", "superAdmin").First(&parent).Error; err != nil {
		global.GVA_LOG.Error("未找到 superAdmin 父菜单，跳过 IP 封禁菜单初始化", zap.Error(err))
		return
	}

	menu := sysModel.SysBaseMenu{
		MenuLevel: 1,
		Hidden:    false,
		ParentId:  parent.ID,
		Path:      "ipBan",
		Name:      "ipBan",
		Component: "view/superAdmin/security/ipBan.vue",
		Sort:      20,
		Meta: sysModel.Meta{
			Title: "IP封禁管理",
			Icon:  "lock",
		},
	}
	if err := db.Create(&menu).Error; err != nil {
		global.GVA_LOG.Error("初始化 IP 封禁菜单失败", zap.Error(err))
		return
	}

	// 仅为超级管理员角色（888、8881）分配该菜单
	for _, authID := range []uint{888, 8881} {
		db.Exec("INSERT IGNORE INTO sys_authority_menus (sys_authority_authority_id, sys_base_menu_id) VALUES (?, ?)",
			fmt.Sprintf("%d", authID), menu.ID)
	}
	global.GVA_LOG.Info("IP 封禁管理菜单初始化成功")
}

func initSecurityCasbin(db *gorm.DB) {
	// 仅超级管理员和内部管理角色可访问IP封禁管理
	authIDs := []string{"888", "8881"}
	rules := []struct {
		Path   string
		Method string
	}{
		{"/sysBannedIP/banIP", "POST"},
		{"/sysBannedIP/unbanIP", "POST"},
		{"/sysBannedIP/getBannedIPList", "GET"},
		{"/sysBannedIP/getAttackStats", "GET"},
	}

	for _, authID := range authIDs {
		for _, rule := range rules {
			var count int64
			db.Table("casbin_rule").Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?",
				"p", authID, rule.Path, rule.Method).Count(&count)
			if count == 0 {
				db.Exec("INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES (?, ?, ?, ?)",
					"p", authID, rule.Path, rule.Method)
			}
		}
	}
}

func cleanupCasbinManagementOvergrant(db *gorm.DB) {
	if db == nil || !db.Migrator().HasTable("casbin_rule") {
		return
	}

	targets := []struct {
		Path   string
		Method string
	}{
		{Path: "/api/createApi", Method: "POST"},
		{Path: "/api/deleteApi", Method: "POST"},
		{Path: "/api/updateApi", Method: "POST"},
		{Path: "/api/deleteApisByIds", Method: "DELETE"},
		{Path: "/api/ignoreApi", Method: "POST"},
		{Path: "/api/freshCasbin", Method: "GET"},
		{Path: "/api/getApiList", Method: "POST"},
		{Path: "/api/getApiById", Method: "POST"},
		{Path: "/api/getAllApis", Method: "POST"},
		{Path: "/casbin/updateCasbin", Method: "POST"},
		{Path: "/casbin/getPolicyPathByAuthorityId", Method: "POST"},
	}

	var affected int64
	for _, target := range targets {
		res := db.Exec(
			"DELETE FROM casbin_rule WHERE ptype = ? AND v1 = ? AND v2 = ? AND v0 NOT IN ?",
			"p", target.Path, target.Method, []string{"888"},
		)
		if res.Error != nil {
			global.GVA_LOG.Warn("清理权限管理接口过授权失败", zap.String("path", target.Path), zap.String("method", target.Method), zap.Error(res.Error))
			continue
		}
		affected += res.RowsAffected
	}

	if affected > 0 {
		global.GVA_LOG.Warn("已清理权限管理接口过授权规则", zap.Int64("affected", affected))
	}
}

func cleanupSystemSensitiveOvergrant(db *gorm.DB) {
	if db == nil || !db.Migrator().HasTable("casbin_rule") {
		return
	}

	targets := []struct {
		Path   string
		Method string
	}{
		{Path: "/system/getSystemConfig", Method: "POST"},
		{Path: "/system/setSystemConfig", Method: "POST"},
		{Path: "/autoCode/getDB", Method: "GET"},
		{Path: "/autoCode/getTables", Method: "GET"},
		{Path: "/autoCode/getColumn", Method: "GET"},
		{Path: "/autoCode/preview", Method: "POST"},
		{Path: "/autoCode/createTemp", Method: "POST"},
		{Path: "/autoCode/addFunc", Method: "POST"},
		{Path: "/autoCode/mcp", Method: "POST"},
		{Path: "/autoCode/mcpList", Method: "POST"},
		{Path: "/autoCode/mcpTest", Method: "POST"},
		{Path: "/autoCode/getPackage", Method: "POST"},
		{Path: "/autoCode/delPackage", Method: "POST"},
		{Path: "/autoCode/createPackage", Method: "POST"},
		{Path: "/autoCode/getTemplates", Method: "GET"},
		{Path: "/autoCode/pubPlug", Method: "POST"},
		{Path: "/autoCode/installPlugin", Method: "POST"},
		{Path: "/autoCode/removePlugin", Method: "POST"},
		{Path: "/autoCode/getPluginList", Method: "GET"},
		{Path: "/autoCode/llmAuto", Method: "POST"},
		{Path: "/autoCode/initMenu", Method: "POST"},
		{Path: "/autoCode/initAPI", Method: "POST"},
		{Path: "/autoCode/initDictionary", Method: "POST"},
		{Path: "/autoCode/getMeta", Method: "POST"},
		{Path: "/autoCode/delSysHistory", Method: "POST"},
		{Path: "/autoCode/rollback", Method: "POST"},
		{Path: "/autoCode/getSysHistory", Method: "POST"},
		{Path: "/authority/createAuthority", Method: "POST"},
		{Path: "/authority/copyAuthority", Method: "POST"},
		{Path: "/authority/deleteAuthority", Method: "POST"},
		{Path: "/authority/updateAuthority", Method: "PUT"},
		{Path: "/authority/getAuthorityList", Method: "POST"},
		{Path: "/authority/setDataAuthority", Method: "POST"},
		{Path: "/menu/getMenuList", Method: "POST"},
		{Path: "/menu/addBaseMenu", Method: "POST"},
		{Path: "/menu/addMenuAuthority", Method: "POST"},
		{Path: "/menu/getMenuAuthority", Method: "POST"},
		{Path: "/menu/deleteBaseMenu", Method: "POST"},
		{Path: "/menu/updateBaseMenu", Method: "POST"},
		{Path: "/menu/getBaseMenuById", Method: "POST"},
		{Path: "/authorityBtn/setAuthorityBtn", Method: "POST"},
		{Path: "/authorityBtn/getAuthorityBtn", Method: "POST"},
		{Path: "/authorityBtn/canRemoveAuthorityBtn", Method: "POST"},
		{Path: "/sysDictionaryDetail/findSysDictionaryDetail", Method: "GET"},
		{Path: "/sysDictionaryDetail/updateSysDictionaryDetail", Method: "PUT"},
		{Path: "/sysDictionaryDetail/createSysDictionaryDetail", Method: "POST"},
		{Path: "/sysDictionaryDetail/getSysDictionaryDetailList", Method: "GET"},
		{Path: "/sysDictionaryDetail/deleteSysDictionaryDetail", Method: "DELETE"},
		{Path: "/sysDictionaryDetail/getDictionaryTreeList", Method: "GET"},
		{Path: "/sysDictionaryDetail/getDictionaryTreeListByType", Method: "GET"},
		{Path: "/sysDictionaryDetail/getDictionaryDetailsByParent", Method: "GET"},
		{Path: "/sysDictionaryDetail/getDictionaryPath", Method: "GET"},
		{Path: "/sysDictionary/findSysDictionary", Method: "GET"},
		{Path: "/sysDictionary/updateSysDictionary", Method: "PUT"},
		{Path: "/sysDictionary/getSysDictionaryList", Method: "GET"},
		{Path: "/sysDictionary/createSysDictionary", Method: "POST"},
		{Path: "/sysDictionary/deleteSysDictionary", Method: "DELETE"},
		{Path: "/sysDictionary/importSysDictionary", Method: "POST"},
		{Path: "/sysDictionary/exportSysDictionary", Method: "GET"},
		{Path: "/sysBannedIP/banIP", Method: "POST"},
		{Path: "/sysBannedIP/unbanIP", Method: "POST"},
		{Path: "/sysBannedIP/getBannedIPList", Method: "GET"},
		{Path: "/sysBannedIP/getAttackStats", Method: "GET"},
		{Path: "/user/admin_register", Method: "POST"},
		{Path: "/user/setUserAuthority", Method: "POST"},
		{Path: "/user/setUserAuthorities", Method: "POST"},
		{Path: "/user/resetPassword", Method: "POST"},
		{Path: "/user/deleteUser", Method: "DELETE"},
		{Path: "/user/setUserInfo", Method: "PUT"},
		{Path: "/user/getUserList", Method: "POST"},
		{Path: "/customer/customer", Method: "POST"},
		{Path: "/customer/customer", Method: "PUT"},
		{Path: "/customer/customer", Method: "DELETE"},
		{Path: "/customer/customer", Method: "GET"},
		{Path: "/customer/customerList", Method: "GET"},
	}

	var affected int64
	for _, target := range targets {
		res := db.Exec(
			"DELETE FROM casbin_rule WHERE ptype = ? AND v1 = ? AND v2 = ? AND v0 NOT IN ?",
			"p", target.Path, target.Method, []string{"888", "8881"},
		)
		if res.Error != nil {
			global.GVA_LOG.Warn("清理系统敏感接口过授权失败", zap.String("path", target.Path), zap.String("method", target.Method), zap.Error(res.Error))
			continue
		}
		affected += res.RowsAffected
	}

	if affected > 0 {
		global.GVA_LOG.Warn("已清理系统敏感接口过授权规则", zap.Int64("affected", affected))
	}
}

func loadBannedIPCache() {
	bannedIPService := service.ServiceGroupApp.SystemServiceGroup.BannedIPService
	bannedIPService.LoadBannedIPCache()
}
