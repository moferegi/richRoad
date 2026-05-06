package initialize

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitVisitorData 自动注册访客统计模块的菜单、API、casbin权限
// 仅在数据缺失时插入，幂等执行
func InitVisitorData() {
	if global.GVA_DB == nil {
		return
	}
	db := global.GVA_DB

	initVisitorApis(db)
	initVisitorMenu(db)
	initVisitorCasbin(db)
}

func initVisitorApis(db *gorm.DB) {
	apis := []sysModel.SysApi{
		{ApiGroup: "访客统计", Method: "GET", Path: "/visitor/getVisitorLogList", Description: "获取访客日志列表"},
		{ApiGroup: "访客统计", Method: "GET", Path: "/visitor/getVisitorSummaryList", Description: "获取访客汇总列表"},
		{ApiGroup: "访客统计", Method: "GET", Path: "/visitor/getTodayStats", Description: "获取今日实时统计"},
		{ApiGroup: "访客统计", Method: "GET", Path: "/visitor/getKefuGuideStats", Description: "获取客服引导漏斗统计"},
		{ApiGroup: "访客统计", Method: "POST", Path: "/visitor/aggregateDailySummary", Description: "手动触发日汇总聚合"},
	}
	for _, api := range apis {
		var count int64
		db.Model(&sysModel.SysApi{}).Where("path = ? AND method = ?", api.Path, api.Method).Count(&count)
		if count == 0 {
			if err := db.Create(&api).Error; err != nil {
				global.GVA_LOG.Error("初始化访客API失败", zap.Error(err))
			}
		}
	}
}

func initVisitorMenu(db *gorm.DB) {
	// 检查是否已存在
	var count int64
	db.Model(&sysModel.SysBaseMenu{}).Where("name = ?", "visitor").Count(&count)
	if count > 0 {
		return
	}

	// 查找"客户端"父菜单
	var parent sysModel.SysBaseMenu
	if err := db.Where("name = ?", "client").First(&parent).Error; err != nil {
		global.GVA_LOG.Error("未找到客户端父菜单，跳过访客菜单初始化", zap.Error(err))
		return
	}

	menu := sysModel.SysBaseMenu{
		MenuLevel: 1,
		Hidden:    false,
		ParentId:  parent.ID,
		Path:      "visitor",
		Name:      "visitor",
		Component: "view/client/visitor/visitor.vue",
		Sort:      11,
		Meta: sysModel.Meta{
			Title: "访客统计",
			Icon:  "data-line",
		},
	}
	if err := db.Create(&menu).Error; err != nil {
		global.GVA_LOG.Error("初始化访客菜单失败", zap.Error(err))
		return
	}

	// 为所有角色分配该菜单权限
	var authorities []sysModel.SysAuthority
	db.Find(&authorities)
	for _, auth := range authorities {
		db.Exec("INSERT IGNORE INTO sys_authority_menus (sys_authority_authority_id, sys_base_menu_id) VALUES (?, ?)",
			fmt.Sprintf("%d", auth.AuthorityId), menu.ID)
	}
	global.GVA_LOG.Info("访客统计菜单初始化成功")
}

func initVisitorCasbin(db *gorm.DB) {
	// 获取所有角色
	var authorities []sysModel.SysAuthority
	db.Find(&authorities)

	paths := []struct {
		Path   string
		Method string
	}{
		{"/visitor/getVisitorLogList", "GET"},
		{"/visitor/getVisitorSummaryList", "GET"},
		{"/visitor/getTodayStats", "GET"},
		{"/visitor/getKefuGuideStats", "GET"},
		{"/visitor/aggregateDailySummary", "POST"},
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
