package client

import (
	"strings"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"go.uber.org/zap"
)

const tryonTokenQuotaExhaustedCacheTTL = 2 * time.Minute

var tryonTokenQuotaExhaustedCache = struct {
	mu    sync.RWMutex
	items map[string]time.Time
}{
	items: map[string]time.Time{},
}

type modelTokenUsageCountRow struct {
	TokenFingerprint string `json:"tokenFingerprint"`
	UsedSuccessCount int64  `json:"usedSuccessCount"`
}

type tokenQuotaCandidate struct {
	Token       string
	Fingerprint string
	QuotaTotal  int
}

func buildTokenQuotaCacheKey(modelKey string, modelUsage string, fingerprint string) string {
	return strings.ToLower(strings.TrimSpace(modelKey)) + ":" + strings.ToLower(strings.TrimSpace(modelUsage)) + ":" + strings.TrimSpace(fingerprint)
}

func isTokenQuotaExhaustedCached(modelKey string, modelUsage string, fingerprint string) bool {
	key := buildTokenQuotaCacheKey(modelKey, modelUsage, fingerprint)
	if key == "::" || strings.TrimSpace(fingerprint) == "" {
		return false
	}

	now := time.Now()
	tryonTokenQuotaExhaustedCache.mu.RLock()
	expireAt, ok := tryonTokenQuotaExhaustedCache.items[key]
	tryonTokenQuotaExhaustedCache.mu.RUnlock()
	if !ok {
		return false
	}
	if expireAt.After(now) {
		return true
	}

	tryonTokenQuotaExhaustedCache.mu.Lock()
	if storedExpireAt, exists := tryonTokenQuotaExhaustedCache.items[key]; exists && !storedExpireAt.After(now) {
		delete(tryonTokenQuotaExhaustedCache.items, key)
	}
	tryonTokenQuotaExhaustedCache.mu.Unlock()
	return false
}

func markTokenQuotaExhausted(modelKey string, modelUsage string, fingerprint string) {
	key := buildTokenQuotaCacheKey(modelKey, modelUsage, fingerprint)
	if key == "::" || strings.TrimSpace(fingerprint) == "" {
		return
	}

	tryonTokenQuotaExhaustedCache.mu.Lock()
	tryonTokenQuotaExhaustedCache.items[key] = time.Now().Add(tryonTokenQuotaExhaustedCacheTTL)
	tryonTokenQuotaExhaustedCache.mu.Unlock()
}

func clearTokenQuotaExhausted(modelKey string, modelUsage string, fingerprint string) {
	key := buildTokenQuotaCacheKey(modelKey, modelUsage, fingerprint)
	if key == "::" || strings.TrimSpace(fingerprint) == "" {
		return
	}

	tryonTokenQuotaExhaustedCache.mu.Lock()
	delete(tryonTokenQuotaExhaustedCache.items, key)
	tryonTokenQuotaExhaustedCache.mu.Unlock()
}

func resetTryonTokenQuotaExhaustedCache() {
	tryonTokenQuotaExhaustedCache.mu.Lock()
	tryonTokenQuotaExhaustedCache.items = map[string]time.Time{}
	tryonTokenQuotaExhaustedCache.mu.Unlock()
}

func resolveDefaultModelQuotaTotal(modelCfg *tryonModelConfig, usage string) int {
	if modelCfg == nil {
		return 0
	}
	usage = normalizeTryonModelUsageValue(usage)
	if usage == "refiner" && modelCfg.RefinerFreeQuotaTotal > 0 {
		return modelCfg.RefinerFreeQuotaTotal
	}
	if modelCfg.FreeQuotaTotal > 0 {
		return modelCfg.FreeQuotaTotal
	}
	return 400
}

func resolveModelTokenQuotaTotal(modelCfg *tryonModelConfig, usage string, token string) int {
	token = normalizeTokenValue(token)
	if token == "" {
		return 0
	}

	defaultTotal := resolveDefaultModelQuotaTotal(modelCfg, usage)
	if modelCfg == nil {
		return defaultTotal
	}

	fingerprint := buildTokenFingerprint(token)
	for _, item := range modelCfg.TokenQuotas {
		configuredToken := normalizeTokenValue(item.Token)
		configuredFingerprint := strings.TrimSpace(item.TokenFingerprint)
		if configuredToken == "" && configuredFingerprint == "" {
			continue
		}

		matched := false
		if configuredToken != "" && configuredToken == token {
			matched = true
		}
		if !matched && configuredFingerprint != "" && strings.EqualFold(configuredFingerprint, fingerprint) {
			matched = true
		}
		if !matched {
			continue
		}

		if item.FreeQuotaTotal < 0 {
			return 0
		}
		return item.FreeQuotaTotal
	}

	if defaultTotal < 0 {
		return 0
	}
	return defaultTotal
}

func queryModelTokenUsedSuccessCount(modelKey string, modelUsage string, fingerprints []string) (map[string]int64, error) {
	result := make(map[string]int64)
	if global.GVA_DB == nil {
		return result, nil
	}
	modelKey = strings.TrimSpace(modelKey)
	modelUsage = normalizeTryonModelUsageValue(modelUsage)
	if modelKey == "" || modelUsage == "" {
		return result, nil
	}

	normalizedFingerprints := make([]string, 0, len(fingerprints))
	fingerprintSet := map[string]struct{}{}
	for _, item := range fingerprints {
		fingerprint := strings.TrimSpace(item)
		if fingerprint == "" {
			continue
		}
		if _, exists := fingerprintSet[fingerprint]; exists {
			continue
		}
		fingerprintSet[fingerprint] = struct{}{}
		normalizedFingerprints = append(normalizedFingerprints, fingerprint)
	}
	if len(normalizedFingerprints) == 0 {
		return result, nil
	}

	rows := make([]modelTokenUsageCountRow, 0)
	err := global.GVA_DB.Model(&clientModel.ModelCallLog{}).
		Select("token_fingerprint, COUNT(1) as used_success_count").
		Where("status = ? AND COALESCE(NULLIF(TRIM(model_key), ''), '') = ? AND COALESCE(NULLIF(TRIM(model_usage), ''), 'tryon') = ? AND token_fingerprint IN ?", tryonTaskStatusSuccess, modelKey, modelUsage, normalizedFingerprints).
		Group("token_fingerprint").
		Scan(&rows).Error
	if err != nil {
		return result, err
	}

	for _, row := range rows {
		fingerprint := strings.TrimSpace(row.TokenFingerprint)
		if fingerprint == "" {
			continue
		}
		result[fingerprint] = row.UsedSuccessCount
	}
	return result, nil
}

func buildTokenQuotaCandidates(modelCfg *tryonModelConfig, modelUsage string, tokens []string) []tokenQuotaCandidate {
	result := make([]tokenQuotaCandidate, 0, len(tokens))
	seen := map[string]struct{}{}
	for _, rawToken := range tokens {
		token := normalizeTokenValue(rawToken)
		if token == "" {
			continue
		}
		if _, exists := seen[token]; exists {
			continue
		}
		seen[token] = struct{}{}
		result = append(result, tokenQuotaCandidate{
			Token:       token,
			Fingerprint: buildTokenFingerprint(token),
			QuotaTotal:  resolveModelTokenQuotaTotal(modelCfg, modelUsage, token),
		})
	}
	return result
}

func (s *TryonTaskService) filterTokenCandidatesByQuota(modelCfg *tryonModelConfig, modelUsage string, tokens []string) (available []string, allExhausted bool) {
	modelUsage = normalizeTryonModelUsageValue(modelUsage)
	if modelCfg == nil {
		return tokens, false
	}

	modelKey := strings.TrimSpace(modelCfg.Key)
	candidates := buildTokenQuotaCandidates(modelCfg, modelUsage, tokens)
	if len(candidates) == 0 {
		return nil, false
	}
	if modelKey == "" || modelUsage == "" {
		available = make([]string, 0, len(candidates))
		for _, item := range candidates {
			available = append(available, item.Token)
		}
		return available, false
	}

	needLookup := make([]string, 0, len(candidates))
	available = make([]string, 0, len(candidates))
	needLookupSet := map[string]struct{}{}
	for _, item := range candidates {
		if strings.TrimSpace(item.Fingerprint) == "" {
			available = append(available, item.Token)
			continue
		}
		if item.QuotaTotal <= 0 {
			markTokenQuotaExhausted(modelKey, modelUsage, item.Fingerprint)
			continue
		}
		if isTokenQuotaExhaustedCached(modelKey, modelUsage, item.Fingerprint) {
			continue
		}
		if _, exists := needLookupSet[item.Fingerprint]; !exists {
			needLookupSet[item.Fingerprint] = struct{}{}
			needLookup = append(needLookup, item.Fingerprint)
		}
	}

	usedCountByFingerprint := map[string]int64{}
	if len(needLookup) > 0 {
		var err error
		usedCountByFingerprint, err = queryModelTokenUsedSuccessCount(modelKey, modelUsage, needLookup)
		if err != nil {
			global.GVA_LOG.Warn("查询 token 额度使用次数失败，回退不过滤", zap.Error(err), zap.String("modelKey", modelKey), zap.String("modelUsage", modelUsage))
			available = available[:0]
			for _, item := range candidates {
				available = append(available, item.Token)
			}
			return available, false
		}
	}

	availableSet := map[string]struct{}{}
	for _, item := range available {
		availableSet[item] = struct{}{}
	}

	for _, item := range candidates {
		if _, exists := availableSet[item.Token]; exists {
			continue
		}
		if strings.TrimSpace(item.Fingerprint) == "" {
			available = append(available, item.Token)
			availableSet[item.Token] = struct{}{}
			continue
		}
		if item.QuotaTotal <= 0 {
			continue
		}
		if isTokenQuotaExhaustedCached(modelKey, modelUsage, item.Fingerprint) {
			continue
		}
		usedCount := usedCountByFingerprint[item.Fingerprint]
		remaining := int64(item.QuotaTotal) - usedCount
		if remaining <= 0 {
			markTokenQuotaExhausted(modelKey, modelUsage, item.Fingerprint)
			continue
		}
		clearTokenQuotaExhausted(modelKey, modelUsage, item.Fingerprint)
		available = append(available, item.Token)
		availableSet[item.Token] = struct{}{}
	}

	if len(available) == 0 {
		return available, true
	}
	return available, false
}
