package initialize

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
)

// InitCasbin 为客服系统所有需鉴权的 API 注册 Casbin 规则
// 幂等执行：仅在规则不存在时插入
func InitCasbin() {
	if global.GVA_DB == nil {
		return
	}
	db := global.GVA_DB

	// 需要对所有角色开放的路径（客户端用户接口）
	clientPaths := []struct{ Path, Method string }{
		{"/cs/conversation/getOrCreate", "POST"},
		{"/cs/conversation/rate", "POST"},
		{"/cs/message/history", "GET"},
		{"/cs/message/upload", "POST"},
		{"/cs/message/revoke", "POST"},
		{"/cs/quickReply/all", "GET"},
	}

	// 坐席接口：所有后台角色均可使用（任何后台账号都能成为坐席）
	agentPaths := []struct{ Path, Method string }{
		{"/cs/agent/conversation/list", "GET"},
		{"/cs/agent/conversation/accept", "POST"},
		{"/cs/agent/conversation/close", "POST"},
		{"/cs/agent/conversation/transfer", "POST"},
		{"/cs/agent/message/send", "POST"},
		{"/cs/agent/message/upload", "POST"},
		{"/cs/agent/quickReply/list", "GET"},
	}

	// 仅超级管理员（888）可用的路径（管理后台接口）
	adminPaths := []struct{ Path, Method string }{
		{"/cs/admin/conversation/list", "GET"},
		{"/cs/admin/conversation/assign", "POST"},
		{"/cs/admin/agent/list", "GET"},
		{"/cs/admin/agent/create", "POST"},
		{"/cs/admin/agent/update", "PUT"},
		{"/cs/admin/agent/delete", "DELETE"},
		{"/cs/admin/quickReply/create", "POST"},
		{"/cs/admin/quickReply/update", "PUT"},
		{"/cs/admin/quickReply/delete", "DELETE"},
		{"/cs/admin/blacklist/list", "GET"},
		{"/cs/admin/blacklist/add", "POST"},
		{"/cs/admin/blacklist/remove", "DELETE"},
		{"/cs/admin/config/get", "GET"},
		{"/cs/admin/config/update", "PUT"},
	}

	// 获取所有角色
	var authorities []sysModel.SysAuthority
	db.Find(&authorities)

	inserted := 0

	for _, auth := range authorities {
		authId := fmt.Sprintf("%d", auth.AuthorityId)

		// 客户端路径：给所有角色注册
		for _, p := range clientPaths {
			var count int64
			db.Table("casbin_rule").Where(
				"ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?",
				"p", authId, p.Path, p.Method,
			).Count(&count)
			if count == 0 {
				if err := db.Exec(
					"INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES (?, ?, ?, ?)",
					"p", authId, p.Path, p.Method,
				).Error; err != nil {
					zap.L().Warn("客服 Casbin 规则插入失败", zap.Error(err))
				} else {
					inserted++
				}
			}
		}

		// 坐席路径：仅后台角色注册，排除客户端角色 8080
		if auth.AuthorityId != 8080 {
			for _, p := range agentPaths {
				var count int64
				db.Table("casbin_rule").Where(
					"ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?",
					"p", authId, p.Path, p.Method,
				).Count(&count)
				if count == 0 {
					if err := db.Exec(
						"INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES (?, ?, ?, ?)",
						"p", authId, p.Path, p.Method,
					).Error; err != nil {
						zap.L().Warn("客服坐席 Casbin 规则插入失败", zap.Error(err))
					} else {
						inserted++
					}
				}
			}
		}
	}

	// 管理员路径：只给 888 注册
	for _, p := range adminPaths {
		var count int64
		db.Table("casbin_rule").Where(
			"ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?",
			"p", "888", p.Path, p.Method,
		).Count(&count)
		if count == 0 {
			if err := db.Exec(
				"INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES (?, ?, ?, ?)",
				"p", "888", p.Path, p.Method,
			).Error; err != nil {
				zap.L().Warn("客服管理 Casbin 规则插入失败", zap.Error(err))
			} else {
				inserted++
			}
		}
	}

	// 安全收敛：移除客户端角色 8080 误持有的坐席接口权限
	for _, p := range agentPaths {
		_ = db.Exec(
			"DELETE FROM casbin_rule WHERE ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?",
			"p", "8080", p.Path, p.Method,
		).Error
	}

	if inserted > 0 {
		zap.L().Info(fmt.Sprintf("客服系统 Casbin 规则初始化完成，新增 %d 条", inserted))
	}

	// 刷新内存中的 enforcer policy，确保新插入的规则立即生效
	// （sync.Once 可能已创建 enforcer，需要 reload 才能感知新规则）
	if e := utils.GetCasbin(); e != nil {
		if err := e.LoadPolicy(); err != nil {
			zap.L().Warn("客服系统 Casbin 规则 LoadPolicy 失败", zap.Error(err))
		}
	}
}
