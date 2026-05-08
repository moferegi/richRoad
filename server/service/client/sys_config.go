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

type AliyunTryonQuotaItem struct {
	ModelKey                 string `json:"modelKey"`
	Model                    string `json:"model"`
	Provider                 string `json:"provider"`
	FreeQuotaTotal           int    `json:"freeQuotaTotal"`
	UsedSuccessCount         int64  `json:"usedSuccessCount"`
	RemainingEstimate        int64  `json:"remainingEstimate"`
	RefinerEnabled           bool   `json:"refinerEnabled"`
	RefinerModel             string `json:"refinerModel"`
	RefinerFreeQuotaTotal    int    `json:"refinerFreeQuotaTotal"`
	RefinerUsedSuccessCount  int64  `json:"refinerUsedSuccessCount"`
	RefinerRemainingEstimate int64  `json:"refinerRemainingEstimate"`
	LastRefreshedAt          string `json:"lastRefreshedAt"`
	EstimateDescription      string `json:"estimateDescription"`
}

// GetSysConfigList 分页获取系统参数列表
func (s *SysConfigService) GetSysConfigList(info clientReq.SysConfigSearch) (list []client.SysConfig, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&client.SysConfig{})

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
	return global.GVA_DB.Model(&client.SysConfig{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"config_value": configValue,
			"remark":       remark,
		}).Error
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

		freeQuotaTotal := item.FreeQuotaTotal
		if freeQuotaTotal <= 0 {
			freeQuotaTotal = 400
		}

		var usedSuccessCount int64
		err = global.GVA_DB.Model(&client.TryonTask{}).
			Where("provider = ? AND status = ?", key, tryonTaskStatusSuccess).
			Count(&usedSuccessCount).Error
		if err != nil {
			return nil, err
		}

		remainingEstimate := int64(freeQuotaTotal) - usedSuccessCount
		if remainingEstimate < 0 {
			remainingEstimate = 0
		}

		supportsRefiner := item.supportsRefinerValue()
		refinerModel := ""
		refinerFreeQuotaTotal := 0
		var refinerUsedSuccessCount int64
		var refinerRemainingEstimate int64

		if supportsRefiner {
			refinerModel = item.refinerModelValue()
			refinerFreeQuotaTotal = item.RefinerFreeQuotaTotal
			if refinerFreeQuotaTotal <= 0 {
				if freeQuotaTotal > 0 {
					refinerFreeQuotaTotal = freeQuotaTotal
				} else {
					refinerFreeQuotaTotal = 400
				}
			}

			err = global.GVA_DB.Model(&client.TryonTask{}).
				Where("provider = ? AND enable_refiner = ? AND refiner_status = ?", key, true, tryonRefinerStatusSuccess).
				Count(&refinerUsedSuccessCount).Error
			if err != nil {
				return nil, err
			}

			refinerRemainingEstimate = int64(refinerFreeQuotaTotal) - refinerUsedSuccessCount
			if refinerRemainingEstimate < 0 {
				refinerRemainingEstimate = 0
			}
		}

		list = append(list, AliyunTryonQuotaItem{
			ModelKey:                 key,
			Model:                    strings.TrimSpace(item.Model),
			Provider:                 strings.TrimSpace(item.Provider),
			FreeQuotaTotal:           freeQuotaTotal,
			UsedSuccessCount:         usedSuccessCount,
			RemainingEstimate:        remainingEstimate,
			RefinerEnabled:           supportsRefiner,
			RefinerModel:             refinerModel,
			RefinerFreeQuotaTotal:    refinerFreeQuotaTotal,
			RefinerUsedSuccessCount:  refinerUsedSuccessCount,
			RefinerRemainingEstimate: refinerRemainingEstimate,
			LastRefreshedAt:          nowText,
			EstimateDescription:      "本地估算值（按当前系统成功任务数统计，含基础与精修接口），官方免费额度请以百炼控制台为准",
		})
	}

	return list, nil
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
	key := strings.ToLower(strings.TrimSpace(item.Key))
	return strings.Contains(key, "aliyun") || strings.Contains(key, "aitryon")
}
