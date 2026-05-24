package initialize

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitInviteData 自动注册邀请系统和系统参数模块的API、菜单、casbin权限
func InitInviteData() {
	if global.GVA_DB == nil {
		return
	}
	db := global.GVA_DB

	initInviteApis(db)
	initInviteMenu(db)
	initInviteCasbin(db)
}

func initInviteApis(db *gorm.DB) {
	apis := []sysModel.SysApi{
		{ApiGroup: "系统参数", Method: "GET", Path: "/sysConfig/getSysConfigList", Description: "获取系统参数列表"},
		{ApiGroup: "系统参数", Method: "GET", Path: "/sysConfig/getAliyunTryonQuotaEstimate", Description: "获取阿里试衣模型额度估算"},
		{ApiGroup: "系统参数", Method: "PUT", Path: "/sysConfig/updateSysConfig", Description: "更新系统参数"},
		{ApiGroup: "客户端用户", Method: "GET", Path: "/clientUser/getSubordinates", Description: "获取下级用户列表"},
		{ApiGroup: "客户端用户", Method: "POST", Path: "/clientUser/adjustTryonPoint", Description: "后台调整用户试衣币"},
		{ApiGroup: "客户端用户", Method: "GET", Path: "/clientUser/getMyInviteInfo", Description: "获取我的邀请信息"},
		{ApiGroup: "客户端用户", Method: "GET", Path: "/clientUser/getMySubordinates", Description: "获取我的下级列表"},
	}
	for _, api := range apis {
		var count int64
		db.Model(&sysModel.SysApi{}).Where("path = ? AND method = ?", api.Path, api.Method).Count(&count)
		if count == 0 {
			if err := db.Create(&api).Error; err != nil {
				global.GVA_LOG.Error("初始化邀请系统API失败", zap.Error(err))
			}
		}
	}
}

func initInviteMenu(db *gorm.DB) {
	var count int64
	db.Model(&sysModel.SysBaseMenu{}).Where("name = ?", "sysConfig").Count(&count)
	if count > 0 {
		return
	}

	var parent sysModel.SysBaseMenu
	if err := db.Where("name = ?", "client").First(&parent).Error; err != nil {
		global.GVA_LOG.Error("未找到客户端父菜单，跳过系统参数菜单初始化", zap.Error(err))
		return
	}

	menu := sysModel.SysBaseMenu{
		MenuLevel: 1,
		Hidden:    false,
		ParentId:  parent.ID,
		Path:      "sysConfig",
		Name:      "sysConfig",
		Component: "view/client/sysConfig/sysConfig.vue",
		Sort:      12,
		Meta: sysModel.Meta{
			Title: "系统参数",
			Icon:  "setting",
		},
	}
	if err := db.Create(&menu).Error; err != nil {
		global.GVA_LOG.Error("初始化系统参数菜单失败", zap.Error(err))
		return
	}

	var authorities []sysModel.SysAuthority
	db.Find(&authorities)
	for _, auth := range authorities {
		db.Exec("INSERT IGNORE INTO sys_authority_menus (sys_authority_authority_id, sys_base_menu_id) VALUES (?, ?)",
			fmt.Sprintf("%d", auth.AuthorityId), menu.ID)
	}
	global.GVA_LOG.Info("系统参数菜单初始化成功")
}

func initInviteCasbin(db *gorm.DB) {
	var authorities []sysModel.SysAuthority
	db.Find(&authorities)

	allRolePaths := []struct {
		Path   string
		Method string
	}{
		{"/clientUser/getMyInviteInfo", "GET"},
		{"/clientUser/getMySubordinates", "GET"},
	}

	adminOnlyPaths := []struct {
		Path   string
		Method string
	}{
		{"/sysConfig/getSysConfigList", "GET"},
		{"/sysConfig/getAliyunTryonQuotaEstimate", "GET"},
		{"/sysConfig/updateSysConfig", "PUT"},
		{"/clientUser/getSubordinates", "GET"},
		{"/clientUser/adjustTryonPoint", "POST"},
	}

	for _, auth := range authorities {
		authId := fmt.Sprintf("%d", auth.AuthorityId)
		for _, p := range allRolePaths {
			var count int64
			db.Table("casbin_rule").Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?",
				"p", authId, p.Path, p.Method).Count(&count)
			if count == 0 {
				db.Exec("INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES (?, ?, ?, ?)",
					"p", authId, p.Path, p.Method)
			}
		}
	}

	for _, authId := range []string{"888", "8881"} {
		for _, p := range adminOnlyPaths {
			var count int64
			db.Table("casbin_rule").Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?",
				"p", authId, p.Path, p.Method).Count(&count)
			if count == 0 {
				db.Exec("INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES (?, ?, ?, ?)",
					"p", authId, p.Path, p.Method)
			}
		}
	}

	// 收敛历史脏数据：移除非管理员对系统参数更新接口的访问权限
	db.Exec(
		"DELETE FROM casbin_rule WHERE ptype = ? AND v1 = ? AND v2 = ? AND v0 NOT IN (?, ?)",
		"p", "/sysConfig/updateSysConfig", "PUT", "888", "8881",
	)
	// 收敛历史脏数据：移除非管理员对系统参数列表接口的访问权限
	db.Exec(
		"DELETE FROM casbin_rule WHERE ptype = ? AND v1 = ? AND v2 = ? AND v0 NOT IN (?, ?)",
		"p", "/sysConfig/getSysConfigList", "GET", "888", "8881",
	)
	// 收敛历史脏数据：移除非管理员对阿里模型额度估算接口的访问权限
	db.Exec(
		"DELETE FROM casbin_rule WHERE ptype = ? AND v1 = ? AND v2 = ? AND v0 NOT IN (?, ?)",
		"p", "/sysConfig/getAliyunTryonQuotaEstimate", "GET", "888", "8881",
	)
	db.Exec(
		"DELETE FROM casbin_rule WHERE ptype = ? AND v1 = ? AND v2 = ? AND v0 NOT IN (?, ?)",
		"p", "/clientUser/adjustTryonPoint", "POST", "888", "8881",
	)
	db.Exec(
		"DELETE FROM casbin_rule WHERE ptype = ? AND v1 = ? AND v2 = ? AND v0 NOT IN (?, ?)",
		"p", "/clientUser/getSubordinates", "GET", "888", "8881",
	)
}
