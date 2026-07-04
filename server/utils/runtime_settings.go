package utils

import (
	"os"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"gorm.io/gorm"
)

func GetSysConfigRawValue(sysConfigKey string) (string, bool) {
	key := strings.TrimSpace(sysConfigKey)
	if key == "" || global.GVA_DB == nil {
		return "", false
	}

	var cfg clientModel.SysConfig
	err := global.GVA_DB.
		Select("config_value").
		Where("config_key = ?", key).
		Take(&cfg).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", false
		}
		return "", false
	}

	return cfg.ConfigValue, true
}

func GetInt64Setting(sysConfigKey, envKey string, defaultValue int64) int64 {
	if value, ok := GetSysConfigRawValue(sysConfigKey); ok {
		if parsed, ok := parseInt64(strings.TrimSpace(value)); ok {
			return parsed
		}
	}

	if envValue := strings.TrimSpace(os.Getenv(strings.TrimSpace(envKey))); envValue != "" {
		if parsed, ok := parseInt64(envValue); ok {
			return parsed
		}
	}

	return defaultValue
}

func GetBoolSetting(sysConfigKey, envKey string, defaultValue bool) bool {
	if value, ok := GetSysConfigRawValue(sysConfigKey); ok {
		if parsed, ok := parseBool(strings.TrimSpace(value)); ok {
			return parsed
		}
	}

	if envValue := strings.TrimSpace(os.Getenv(strings.TrimSpace(envKey))); envValue != "" {
		if parsed, ok := parseBool(envValue); ok {
			return parsed
		}
	}

	return defaultValue
}

func parseInt64(raw string) (int64, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, false
	}

	parsed, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, false
	}

	return parsed, true
}

func parseBool(raw string) (bool, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return false, false
	}

	parsed, err := strconv.ParseBool(raw)
	if err != nil {
		return false, false
	}

	return parsed, true
}
