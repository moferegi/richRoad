package client

import (
	"encoding/json"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
)

func loadTryonModelUsageMap() map[string]string {
	usageMap := make(map[string]string)
	sysConfigService := &SysConfigService{}
	raw, err := sysConfigService.GetConfigByKey("tryon_models")
	if err != nil || strings.TrimSpace(raw) == "" {
		return usageMap
	}

	list := make([]tryonModelConfig, 0)
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		global.GVA_LOG.Warn("解析试衣模型配置失败，模型用途映射回退", zap.Error(err))
		return usageMap
	}

	for i := range list {
		modelKey := strings.ToLower(strings.TrimSpace(list[i].Key))
		if modelKey == "" {
			continue
		}
		usageMap[modelKey] = normalizeTryonModelUsageValue(list[i].usageValue())
	}

	return usageMap
}

func loadDefaultModelKeyByUsageAndProvider(targetUsage string) map[string]string {
	result := make(map[string]string)
	targetUsage = normalizeTryonModelUsageValue(targetUsage)
	sysConfigService := &SysConfigService{}
	raw, err := sysConfigService.GetConfigByKey("tryon_models")
	if err != nil || strings.TrimSpace(raw) == "" {
		return result
	}

	list := make([]tryonModelConfig, 0)
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		global.GVA_LOG.Warn("解析试衣模型配置失败，默认模型映射回退", zap.Error(err))
		return result
	}

	for i := range list {
		modelKey := strings.TrimSpace(list[i].Key)
		if modelKey == "" {
			continue
		}
		if normalizeTryonModelUsageValue(list[i].usageValue()) != targetUsage {
			continue
		}
		provider := strings.ToLower(strings.TrimSpace(list[i].Provider))
		if provider == "" {
			provider = inferProviderFromModelKey(modelKey)
		}
		if provider == "" {
			continue
		}
		if _, exists := result[provider]; exists {
			continue
		}
		result[provider] = modelKey
	}

	return result
}

func normalizeRefinerStatsModelKey(modelKey string, provider string, usageMap map[string]string, defaultRefinerByProvider map[string]string) string {
	trimmedKey := strings.TrimSpace(modelKey)
	if trimmedKey == "" {
		providerLower := strings.ToLower(strings.TrimSpace(provider))
		if fallback := strings.TrimSpace(defaultRefinerByProvider[providerLower]); fallback != "" {
			return fallback
		}
		return "default"
	}

	usageValue, exists := usageMap[strings.ToLower(trimmedKey)]
	if exists && normalizeTryonModelUsageValue(usageValue) == "refiner" {
		return trimmedKey
	}
	if !exists {
		return trimmedKey
	}

	providerLower := strings.ToLower(strings.TrimSpace(provider))
	if fallback := strings.TrimSpace(defaultRefinerByProvider[providerLower]); fallback != "" {
		return fallback
	}

	return trimmedKey
}

func loadTryonModelProviderMap() map[string]string {
	providerMap := make(map[string]string)
	sysConfigService := &SysConfigService{}
	raw, err := sysConfigService.GetConfigByKey("tryon_models")
	if err != nil || strings.TrimSpace(raw) == "" {
		return providerMap
	}

	list := make([]tryonModelConfig, 0)
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		global.GVA_LOG.Warn("解析试衣模型配置失败，统计回退到任务模型标识", zap.Error(err))
		return providerMap
	}

	for i := range list {
		modelKey := strings.ToLower(strings.TrimSpace(list[i].Key))
		if modelKey == "" {
			continue
		}
		provider := strings.TrimSpace(list[i].Provider)
		if provider == "" {
			provider = inferProviderFromModelKey(modelKey)
		}
		providerMap[modelKey] = provider
	}

	return providerMap
}

func inferProviderFromModelKey(modelKey string) string {
	key := strings.ToLower(strings.TrimSpace(modelKey))
	if key == "" || key == "default" {
		return ""
	}
	if strings.Contains(key, "aliyun") || strings.Contains(key, "dashscope") || strings.Contains(key, "aitryon") {
		return "aliyun"
	}
	if strings.Contains(key, "gradio") || strings.Contains(key, "huggingface") || strings.Contains(key, "hf") {
		return "gradio"
	}
	return ""
}

func resolveTryonModelProvider(modelKey string, providerMap map[string]string) string {
	normalizedKey := strings.ToLower(strings.TrimSpace(modelKey))
	if normalizedKey == "" {
		return ""
	}
	if providerMap != nil {
		if provider := strings.TrimSpace(providerMap[normalizedKey]); provider != "" {
			return provider
		}
	}
	if inferred := inferProviderFromModelKey(normalizedKey); inferred != "" {
		return inferred
	}
	if normalizedKey == "default" {
		return ""
	}
	return strings.TrimSpace(modelKey)
}

func normalizeTryonModelUsageValue(value string) string {
	usage := strings.ToLower(strings.TrimSpace(value))
	switch usage {
	case "tryon", "refiner", "parsing", "beautify":
		return usage
	default:
		return "tryon"
	}
}
