package utils

import (
	"os"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
)

func GetSysConfigRawValue(key string) (string, bool) {
	key = strings.TrimSpace(key)
	if key == "" || global.GVA_DB == nil {
		return "", false
	}

	var cfg clientModel.SysConfig
	if err := global.GVA_DB.Select("config_value").Where("config_key = ?", key).First(&cfg).Error; err != nil {
		return "", false
	}

	return strings.TrimSpace(cfg.ConfigValue), true
}

func GetSysConfigBoolValue(key string) (bool, bool) {
	raw, ok := GetSysConfigRawValue(key)
	if !ok || raw == "" {
		return false, false
	}

	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "yes", "on":
		return true, true
	case "0", "false", "no", "off":
		return false, true
	default:
		return false, false
	}
}

func GetSysConfigInt64Value(key string) (int64, bool) {
	raw, ok := GetSysConfigRawValue(key)
	if !ok || raw == "" {
		return 0, false
	}

	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, false
	}
	return value, true
}

func GetBoolEnv(key string, defaultValue bool) bool {
	raw := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if raw == "" {
		return defaultValue
	}
	switch raw {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return defaultValue
	}
}

func GetInt64Env(key string, defaultValue int64) int64 {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return defaultValue
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func GetBoolSetting(sysConfigKey, envKey string, defaultValue bool) bool {
	if value, ok := GetSysConfigBoolValue(sysConfigKey); ok {
		return value
	}
	if strings.TrimSpace(envKey) != "" {
		return GetBoolEnv(envKey, defaultValue)
	}
	return defaultValue
}

func GetInt64Setting(sysConfigKey, envKey string, defaultValue int64) int64 {
	if value, ok := GetSysConfigInt64Value(sysConfigKey); ok {
		return value
	}
	if strings.TrimSpace(envKey) != "" {
		return GetInt64Env(envKey, defaultValue)
	}
	return defaultValue
}
