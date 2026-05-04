package initialize

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
	"go.uber.org/zap"
)

// Menu 初始化客服系统菜单，并将菜单分配给超级管理员角色
func Menu(ctx context.Context) {
	entities := []model.SysBaseMenu{
		{
			ParentId:  0,
			Path:      "customerService",
			Name:      "customerService",
			Hidden:    false,
			Component: "view/routerHolder.vue", // GVA 父菜单标准组件
			Sort:      90,
			Meta: model.Meta{
				Title: "客服系统",
				Icon:  "Service",
			},
		},
		{
			Path:      "workbench",
			Name:      "csWorkbench",
			Hidden:    false,
			Component: "plugin/customer_service/view/workbench.vue",
			Sort:      1,
			Meta: model.Meta{
				Title: "坐席工作台",
				Icon:  "ChatLineRound",
			},
		},
		{
			Path:      "history",
			Name:      "csHistory",
			Hidden:    false,
			Component: "plugin/customer_service/view/history.vue",
			Sort:      2,
			Meta: model.Meta{
				Title: "历史会话",
				Icon:  "Clock",
			},
		},
		{
			Path:      "agents",
			Name:      "csAgents",
			Hidden:    false,
			Component: "plugin/customer_service/view/agents.vue",
			Sort:      3,
			Meta: model.Meta{
				Title: "坐席管理",
				Icon:  "User",
			},
		},
		{
			Path:      "quickReplies",
			Name:      "csQuickReplies",
			Hidden:    false,
			Component: "plugin/customer_service/view/quickReplies.vue",
			Sort:      4,
			Meta: model.Meta{
				Title: "快捷回复",
				Icon:  "Lightning",
			},
		},
		{
			Path:      "blacklist",
			Name:      "csBlacklist",
			Hidden:    false,
			Component: "plugin/customer_service/view/blacklist.vue",
			Sort:      5,
			Meta: model.Meta{
				Title: "黑名单",
				Icon:  "CircleClose",
			},
		},
		{
			Path:      "config",
			Name:      "csConfig",
			Hidden:    false,
			Component: "plugin/customer_service/view/config.vue",
			Sort:      6,
			Meta: model.Meta{
				Title: "客服配置",
				Icon:  "Setting",
			},
		},
		{
			Path:      "changeLog",
			Name:      "csChangeLog",
			Hidden:    false,
			Component: "plugin/customer_service/view/changeLog.vue",
			Sort:      7,
			Meta: model.Meta{
				Title: "配置变更记录",
				Icon:  "Document",
			},
		},
	}
	utils.RegisterMenus(entities...)

	// 将所有已注册菜单分配给超级管理员（888）
	assignMenusToAdmin(entities)
}

// assignMenusToAdmin 将客服系统菜单分配给 888 超级管理员角色
func assignMenusToAdmin(entities []model.SysBaseMenu) {
	if global.GVA_DB == nil {
		return
	}
	// 按 name 查出实际菜单 ID
	names := make([]string, 0, len(entities))
	for _, e := range entities {
		names = append(names, e.Name)
	}
	var menus []model.SysBaseMenu
	if err := global.GVA_DB.Where("name IN ?", names).Find(&menus).Error; err != nil {
		zap.L().Error("客服菜单查询失败", zap.Error(err))
		return
	}
	for _, m := range menus {
		result := global.GVA_DB.Exec(
			"INSERT IGNORE INTO sys_authority_menus (sys_authority_authority_id, sys_base_menu_id) VALUES (?, ?)",
			fmt.Sprintf("%d", 888), m.ID,
		)
		if result.Error != nil {
			zap.L().Warn("分配客服菜单失败", zap.Error(result.Error), zap.Uint("menuID", m.ID))
		}
	}
}
