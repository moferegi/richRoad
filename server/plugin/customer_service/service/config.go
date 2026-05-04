package service

import (
	"sort"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/customer_service/model"
)

type ConfigService struct{}

const (
	defaultAutoCloseMinutes = 30
	defaultUploadMaxSizeMB  = 5
	maxUploadMaxSizeMB      = 50
	defaultUploadAllowExt   = "jpg,jpeg,png,webp,gif"
)

func normalizeUploadAllowExt(raw string) string {
	parts := strings.Split(raw, ",")
	seen := make(map[string]struct{}, len(parts))
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		ext := strings.ToLower(strings.TrimSpace(part))
		ext = strings.TrimPrefix(ext, ".")
		if ext == "" {
			continue
		}
		valid := true
		for _, ch := range ext {
			if (ch < 'a' || ch > 'z') && (ch < '0' || ch > '9') {
				valid = false
				break
			}
		}
		if !valid {
			continue
		}
		if _, ok := seen[ext]; ok {
			continue
		}
		seen[ext] = struct{}{}
		result = append(result, ext)
	}
	if len(result) == 0 {
		return defaultUploadAllowExt
	}
	sort.Strings(result)
	return strings.Join(result, ",")
}

func normalizeConfig(cfg model.CsConfig) model.CsConfig {
	if cfg.AutoCloseMinutes < 0 {
		cfg.AutoCloseMinutes = 0
	}
	if cfg.AutoCloseMinutes > 1440 {
		cfg.AutoCloseMinutes = 1440
	}

	if cfg.UploadMaxSizeMB <= 0 {
		cfg.UploadMaxSizeMB = defaultUploadMaxSizeMB
	}
	if cfg.UploadMaxSizeMB > maxUploadMaxSizeMB {
		cfg.UploadMaxSizeMB = maxUploadMaxSizeMB
	}

	cfg.UploadAllowExt = normalizeUploadAllowExt(cfg.UploadAllowExt)

	cfg.DefaultAvatarURL = strings.TrimSpace(cfg.DefaultAvatarURL)
	if len(cfg.DefaultAvatarURL) > 512 {
		cfg.DefaultAvatarURL = cfg.DefaultAvatarURL[:512]
	}

	return cfg
}

// GetConfig 获取客服全局配置（不存在时返回默认值）
func (s *ConfigService) GetConfig() (model.CsConfig, error) {
	var cfg model.CsConfig
	err := global.GVA_DB.First(&cfg).Error
	if err != nil {
		// 未初始化时返回默认配置
		return model.CsConfig{
			PlatEnabled:      false,
			AutoCloseMinutes: defaultAutoCloseMinutes,
			UploadMaxSizeMB:  defaultUploadMaxSizeMB,
			UploadAllowExt:   defaultUploadAllowExt,
			DefaultAvatarURL: "",
		}, nil
	}
	return normalizeConfig(cfg), nil
}

// UpdateConfig 更新客服全局配置
func (s *ConfigService) UpdateConfig(platEnabled bool, autoCloseMinutes int, uploadMaxSizeMB int, uploadAllowExt, defaultAvatarURL string) error {
	normalized := normalizeConfig(model.CsConfig{
		PlatEnabled:      platEnabled,
		AutoCloseMinutes: autoCloseMinutes,
		UploadMaxSizeMB:  uploadMaxSizeMB,
		UploadAllowExt:   uploadAllowExt,
		DefaultAvatarURL: defaultAvatarURL,
	})

	var cfg model.CsConfig
	result := global.GVA_DB.First(&cfg)
	if result.Error != nil {
		cfg = normalized
		return global.GVA_DB.Create(&cfg).Error
	}
	return global.GVA_DB.Model(&cfg).Updates(map[string]interface{}{
		"plat_enabled":       normalized.PlatEnabled,
		"auto_close_minutes": normalized.AutoCloseMinutes,
		"upload_max_size_mb": normalized.UploadMaxSizeMB,
		"upload_allow_ext":   normalized.UploadAllowExt,
		"default_avatar_url": normalized.DefaultAvatarURL,
	}).Error
}

// UpdatePlatEnabled 更新平台客服开关（兼容旧调用）
func (s *ConfigService) UpdatePlatEnabled(enabled bool) error {
	cfg, _ := s.GetConfig()
	return s.UpdateConfig(enabled, cfg.AutoCloseMinutes, cfg.UploadMaxSizeMB, cfg.UploadAllowExt, cfg.DefaultAvatarURL)
}

// StartAutoCloseWorker 启动后台定时器，检查空闲会话超时自动关闭
// 每分钟检查一次，关闭超过 AutoCloseMinutes 分钟无新消息的 active 会话
func (s *ConfigService) StartAutoCloseWorker() {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			s.runAutoClose()
		}
	}()
}

func (s *ConfigService) runAutoClose() {
	cfg, err := s.GetConfig()
	if err != nil || cfg.AutoCloseMinutes <= 0 {
		return
	}
	cutoff := time.Now().Add(-time.Duration(cfg.AutoCloseMinutes) * time.Minute)

	// 查找 active 且最近消息时间早于 cutoff 的会话
	var convIDs []uint
	global.GVA_DB.Raw(`
		SELECT c.id FROM cs_conversations c
		WHERE c.status = 'active'
		  AND (
		    (
		      SELECT MAX(m.created_at) FROM cs_messages m WHERE m.conversation_id = c.id
		    ) < ?
		    OR (
		      NOT EXISTS (SELECT 1 FROM cs_messages m2 WHERE m2.conversation_id = c.id)
		      AND c.created_at < ?
		    )
		  )
	`, cutoff, cutoff).Scan(&convIDs)

	for _, id := range convIDs {
		_ = Service.ConversationService.Close(id, model.ConvClosedBySystem)
	}
}
