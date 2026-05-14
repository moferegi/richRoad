package client

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
)

type SysConfigService struct{}

var deprecatedSysConfigKeys = []string{
	"shoe_models",
	"tryon_provider_mode",
	"tryon_provider_url",
	"tryon_provider_token",
}

type AliyunTryonQuotaItem struct {
	ModelKey                 string                 `json:"modelKey"`
	ModelUsage               string                 `json:"modelUsage"`
	Model                    string                 `json:"model"`
	Provider                 string                 `json:"provider"`
	FreeQuotaTotal           int                    `json:"freeQuotaTotal"`
	UsedSuccessCount         int64                  `json:"usedSuccessCount"`
	RemainingEstimate        int64                  `json:"remainingEstimate"`
	TokenQuotaList           []AliyunTokenQuotaItem `json:"tokenQuotaList"`
	TokenConfigCount         int                    `json:"tokenConfigCount"`
	TokenExhaustedCount      int                    `json:"tokenExhaustedCount"`
	RefinerEnabled           bool                   `json:"refinerEnabled"`
	RefinerModel             string                 `json:"refinerModel"`
	RefinerFreeQuotaTotal    int                    `json:"refinerFreeQuotaTotal"`
	RefinerUsedSuccessCount  int64                  `json:"refinerUsedSuccessCount"`
	RefinerRemainingEstimate int64                  `json:"refinerRemainingEstimate"`
	LastRefreshedAt          string                 `json:"lastRefreshedAt"`
	EstimateDescription      string                 `json:"estimateDescription"`
}

type AliyunTokenQuotaItem struct {
	TokenMasked           string `json:"tokenMasked"`
	TokenFingerprint      string `json:"tokenFingerprint"`
	TokenFingerprintShort string `json:"tokenFingerprintShort"`
	FreeQuotaTotal        int    `json:"freeQuotaTotal"`
	UsedSuccessCount      int64  `json:"usedSuccessCount"`
	RemainingEstimate     int64  `json:"remainingEstimate"`
	Exhausted             bool   `json:"exhausted"`
}

// GetSysConfigList 分页获取系统参数列表
func (s *SysConfigService) GetSysConfigList(info clientReq.SysConfigSearch) (list []client.SysConfig, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&client.SysConfig{})
	db = db.Where("config_key NOT IN ?", deprecatedSysConfigKeys)

	if info.ConfigGroup != "" {
		db = db.Where("config_group = ?", info.ConfigGroup)
	}
	if info.ConfigKey != "" {
		db = db.Where("config_key LIKE ?", "%"+info.ConfigKey+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("config_group, id").Limit(limit).Offset(offset).Find(&list).Error
	return
}

// UpdateSysConfig 更新参数值（只允许改value和remark）
func (s *SysConfigService) UpdateSysConfig(id uint, configValue string, remark string) error {
	var existed client.SysConfig
	if err := global.GVA_DB.Select("config_key").Where("id = ?", id).First(&existed).Error; err != nil {
		return err
	}

	err := global.GVA_DB.Model(&client.SysConfig{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"config_value": configValue,
			"remark":       remark,
		}).Error
	if err != nil {
		return err
	}

	if strings.EqualFold(strings.TrimSpace(existed.ConfigKey), "tryon_models") {
		resetTryonTokenQuotaExhaustedCache()
	}

	return nil
}

// GetConfigByKey 根据key获取配置值
func (s *SysConfigService) GetConfigByKey(key string) (string, error) {
	var config client.SysConfig
	err := global.GVA_DB.Where("config_key = ?", key).First(&config).Error
	if err != nil {
		return "", err
	}
	return config.ConfigValue, nil
}

// GetConfigIntByKey 获取int类型参数
func (s *SysConfigService) GetConfigIntByKey(key string, defaultVal int) int {
	val, err := s.GetConfigByKey(key)
	if err != nil || val == "" {
		return defaultVal
	}
	var result int
	_, _ = fmt.Sscanf(val, "%d", &result)
	if result <= 0 {
		return defaultVal
	}
	return result
}

// GetConfigByGroup 按分组获取所有配置
func (s *SysConfigService) GetConfigByGroup(group string) (list []client.SysConfig, err error) {
	err = global.GVA_DB.Where("config_group = ?", group).Order("id").Find(&list).Error
	return
}

// GetAnnouncementConfig 获取公告配置（便捷方法）
func (s *SysConfigService) GetAnnouncementConfig() (map[string]string, error) {
	list, err := s.GetConfigByGroup("announcement")
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, cfg := range list {
		result[cfg.ConfigKey] = cfg.ConfigValue
	}
	return result, nil
}

func (s *SysConfigService) GetAliyunTryonQuotaEstimate(modelKey string) (list []AliyunTryonQuotaItem, err error) {
	raw, err := s.GetConfigByKey("tryon_models")
	if err != nil {
		return nil, err
	}

	models := make([]tryonModelConfig, 0)
	if strings.TrimSpace(raw) != "" {
		if err = json.Unmarshal([]byte(raw), &models); err != nil {
			return nil, err
		}
	}

	modelKey = strings.TrimSpace(modelKey)
	nowText := time.Now().Format(time.RFC3339)
	for i := range models {
		item := &models[i]
		key := strings.TrimSpace(item.Key)
		if key == "" {
			continue
		}
		if modelKey != "" && !strings.EqualFold(key, modelKey) {
			continue
		}
		if !isAliyunQuotaSupportedModel(item) {
			continue
		}
		modelUsage := item.usageValue()

		tokenCandidates := resolveAliyunModelQuotaTokenCandidates(item, modelUsage)
		quotaCandidates := buildTokenQuotaCandidates(item, modelUsage, tokenCandidates)
		tokenQuotaList := make([]AliyunTokenQuotaItem, 0, len(quotaCandidates))

		usedSuccessCount := int64(0)
		remainingEstimate := int64(0)
		freeQuotaTotal := resolveDefaultModelQuotaTotal(item, modelUsage)
		tokenExhaustedCount := 0

		if len(quotaCandidates) > 0 {
			fingerprints := make([]string, 0, len(quotaCandidates))
			for _, candidate := range quotaCandidates {
				if strings.TrimSpace(candidate.Fingerprint) != "" {
					fingerprints = append(fingerprints, candidate.Fingerprint)
				}
			}

			usedCountByFingerprint, queryErr := queryModelTokenUsedSuccessCount(key, modelUsage, fingerprints)
			if queryErr != nil {
				return nil, queryErr
			}

			freeQuotaTotal = 0
			for _, candidate := range quotaCandidates {
				usedCount := usedCountByFingerprint[candidate.Fingerprint]
				remaining := int64(candidate.QuotaTotal) - usedCount
				if remaining < 0 {
					remaining = 0
				}

				if strings.TrimSpace(candidate.Fingerprint) != "" {
					if remaining <= 0 {
						markTokenQuotaExhausted(key, modelUsage, candidate.Fingerprint)
					} else {
						clearTokenQuotaExhausted(key, modelUsage, candidate.Fingerprint)
					}
				}

				if remaining <= 0 {
					tokenExhaustedCount++
				}

				freeQuotaTotal += candidate.QuotaTotal
				usedSuccessCount += usedCount
				remainingEstimate += remaining

				tokenQuotaList = append(tokenQuotaList, AliyunTokenQuotaItem{
					TokenMasked:           maskTokenForDisplay(candidate.Token),
					TokenFingerprint:      candidate.Fingerprint,
					TokenFingerprintShort: shortTokenFingerprint(candidate.Fingerprint),
					FreeQuotaTotal:        candidate.QuotaTotal,
					UsedSuccessCount:      usedCount,
					RemainingEstimate:     remaining,
					Exhausted:             remaining <= 0,
				})
			}
		} else {
			err = global.GVA_DB.Model(&client.ModelCallLog{}).
				Where("status = ? AND COALESCE(NULLIF(TRIM(model_key), ''), '') = ? AND COALESCE(NULLIF(TRIM(model_usage), ''), 'tryon') = ?", tryonTaskStatusSuccess, key, modelUsage).
				Count(&usedSuccessCount).Error
			if err != nil {
				return nil, err
			}

			remainingEstimate = int64(freeQuotaTotal) - usedSuccessCount
			if remainingEstimate < 0 {
				remainingEstimate = 0
			}
		}

		list = append(list, AliyunTryonQuotaItem{
			ModelKey:                 key,
			ModelUsage:               modelUsage,
			Model:                    strings.TrimSpace(item.Model),
			Provider:                 strings.TrimSpace(item.Provider),
			FreeQuotaTotal:           freeQuotaTotal,
			UsedSuccessCount:         usedSuccessCount,
			RemainingEstimate:        remainingEstimate,
			TokenQuotaList:           tokenQuotaList,
			TokenConfigCount:         len(tokenQuotaList),
			TokenExhaustedCount:      tokenExhaustedCount,
			RefinerEnabled:           modelUsage == "refiner",
			RefinerModel:             "",
			RefinerFreeQuotaTotal:    0,
			RefinerUsedSuccessCount:  0,
			RefinerRemainingEstimate: 0,
			LastRefreshedAt:          nowText,
			EstimateDescription:      "本地估算值（按 token 指纹统计成功调用并从各 token 原始额度扣减），官方免费额度请以百炼控制台为准",
		})
	}

	return list, nil
}

func resolveAliyunModelQuotaTokenCandidates(item *tryonModelConfig, modelUsage string) []string {
	if item == nil {
		return nil
	}

	modelUsage = normalizeTryonModelUsageValue(modelUsage)
	switch modelUsage {
	case "refiner":
		return item.refinerAuthTokenCandidates(strings.TrimSpace(item.authToken()))
	case "beautify":
		candidates := make([]string, 0, 4)
		candidates = appendUniqueTokenCandidate(candidates, item.BeautifyToken)
		for _, token := range item.authTokenCandidates() {
			candidates = appendUniqueTokenCandidate(candidates, token)
		}
		return candidates
	default:
		return item.authTokenCandidates()
	}
}

func isAliyunQuotaSupportedModel(item *tryonModelConfig) bool {
	if item == nil {
		return false
	}
	provider := strings.ToLower(strings.TrimSpace(item.Provider))
	if strings.Contains(provider, "aliyun") || strings.Contains(provider, "dashscope") {
		return true
	}
	model := strings.ToLower(strings.TrimSpace(item.Model))
	if strings.HasPrefix(model, "aitryon") {
		return true
	}
	if strings.Contains(model, "retouch") || strings.Contains(model, "facebody") {
		return true
	}
	key := strings.ToLower(strings.TrimSpace(item.Key))
	if strings.Contains(key, "aliyun") || strings.Contains(key, "aitryon") {
		return true
	}

	endpoint := strings.ToLower(strings.TrimSpace(item.endpointURL()))
	if strings.Contains(endpoint, "dashscope.aliyuncs.com") || strings.Contains(endpoint, "facebody") {
		return true
	}

	return false
}
