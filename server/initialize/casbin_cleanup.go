package initialize

import (
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	cleanupLegacyCasbinOverGrantConfigKey = "security_cleanup_legacy_casbin_overgrant"
	cleanupLegacyCasbinOverGrantEnvKey    = "CS_CLEANUP_LEGACY_CASBIN_OVERGRANT"
)

func cleanupLegacyManagedCasbinRules(db *gorm.DB, moduleName string, paths []struct {
	Path   string
	Method string
}, allowedRoleIDs []string) {
	if !utils.GetBoolSetting(cleanupLegacyCasbinOverGrantConfigKey, cleanupLegacyCasbinOverGrantEnvKey, false) || db == nil {
		return
	}
	if len(paths) == 0 || len(allowedRoleIDs) == 0 {
		return
	}

	seen := make(map[string]struct{}, len(paths))
	var affected int64
	for _, p := range paths {
		path := strings.TrimSpace(p.Path)
		method := strings.TrimSpace(strings.ToUpper(p.Method))
		if path == "" || method == "" {
			continue
		}
		key := method + " " + path
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		res := db.Exec(
			"DELETE FROM casbin_rule WHERE ptype = ? AND v1 = ? AND v2 = ? AND v0 NOT IN ?",
			"p", path, method, allowedRoleIDs,
		)
		if res.Error != nil {
			global.GVA_LOG.Warn("清理历史 Casbin 过授权失败",
				zap.String("module", moduleName),
				zap.String("path", path),
				zap.String("method", method),
				zap.Error(res.Error),
			)
			continue
		}
		affected += res.RowsAffected
	}

	if affected > 0 {
		global.GVA_LOG.Warn("已清理历史 Casbin 过授权规则",
			zap.String("module", moduleName),
			zap.Int64("affected", affected),
			zap.Strings("allowedRoleIDs", allowedRoleIDs),
		)
	}
}
