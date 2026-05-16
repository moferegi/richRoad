package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
)

type SysLanguageService struct{}

// CreateSysLanguage 创建语言
func (s *SysLanguageService) CreateSysLanguage(lang *client.SysLanguage) (err error) {
	err = global.GVA_DB.Create(lang).Error
	return err
}

// DeleteSysLanguage 删除语言
func (s *SysLanguageService) DeleteSysLanguage(ID string) (err error) {
	err = global.GVA_DB.Delete(&client.SysLanguage{}, "id = ?", ID).Error
	return err
}

// UpdateSysLanguage 更新语言
func (s *SysLanguageService) UpdateSysLanguage(lang client.SysLanguage) (err error) {
	// 如果设置为默认，先清除其他默认
	if lang.IsDefault != nil && *lang.IsDefault {
		global.GVA_DB.Model(&client.SysLanguage{}).Where("id != ?", lang.ID).Update("is_default", false)
	}
	err = global.GVA_DB.Model(&client.SysLanguage{}).Where("id = ?", lang.ID).Updates(&lang).Error
	return err
}

// GetSysLanguage 根据ID获取语言
func (s *SysLanguageService) GetSysLanguage(ID string) (lang client.SysLanguage, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&lang).Error
	return
}

// GetSysLanguageList 分页获取语言列表
func (s *SysLanguageService) GetSysLanguageList(info clientReq.SysLanguageSearch) (list []client.SysLanguage, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&client.SysLanguage{})

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Order("sort ASC").Find(&list).Error
	return
}

// GetEnabledLanguages 获取已启用的语言列表(客户端用)
func (s *SysLanguageService) GetEnabledLanguages() (list []client.SysLanguage, err error) {
	err = global.GVA_DB.Where("is_enabled = ?", true).Order("sort ASC").Find(&list).Error
	return
}

// TranslateI18n 基于 Google 免费接口执行文本翻译
func (s *SysLanguageService) TranslateI18n(req clientReq.TranslateI18nRequest) (map[string]string, error) {
	text := strings.TrimSpace(req.Text)
	if text == "" {
		return nil, errors.New("empty text")
	}

	sourceLang := normalizeTranslateLang(req.Source)
	if sourceLang == "" {
		sourceLang = "zh-CN"
	}

	httpClient := &http.Client{Timeout: 15 * time.Second}
	translations := make(map[string]string, len(req.Targets))
	seen := make(map[string]struct{}, len(req.Targets))

	for _, target := range req.Targets {
		rawTarget := strings.TrimSpace(target)
		if rawTarget == "" {
			continue
		}
		if _, ok := seen[rawTarget]; ok {
			continue
		}
		seen[rawTarget] = struct{}{}

		normalizedTarget := normalizeTranslateLang(rawTarget)
		if normalizedTarget == "" {
			continue
		}
		if strings.EqualFold(normalizedTarget, sourceLang) {
			translations[rawTarget] = text
			continue
		}

		translated, err := s.translateByGoogle(httpClient, text, sourceLang, normalizedTarget)
		if err != nil {
			return nil, err
		}
		translations[rawTarget] = translated
	}

	if len(translations) == 0 {
		return nil, errors.New("no available target languages")
	}

	return translations, nil
}

func (s *SysLanguageService) translateByGoogle(httpClient *http.Client, text, sourceLang, targetLang string) (string, error) {
	params := url.Values{}
	params.Set("client", "gtx")
	params.Set("sl", sourceLang)
	params.Set("tl", targetLang)
	params.Set("dt", "t")
	params.Set("q", text)

	endpoint := "https://translate.googleapis.com/translate_a/single?" + params.Encode()
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("translate request failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	result, err := parseGoogleTranslateResult(body)
	if err != nil {
		return "", err
	}

	return result, nil
}

func parseGoogleTranslateResult(body []byte) (string, error) {
	var payload []any
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}
	if len(payload) == 0 {
		return "", errors.New("empty translate payload")
	}

	segments, ok := payload[0].([]any)
	if !ok {
		return "", errors.New("invalid translate payload")
	}

	var builder strings.Builder
	for _, segment := range segments {
		part, ok := segment.([]any)
		if !ok || len(part) == 0 {
			continue
		}
		text, ok := part[0].(string)
		if !ok || strings.TrimSpace(text) == "" {
			continue
		}
		builder.WriteString(text)
	}

	translated := strings.TrimSpace(builder.String())
	if translated == "" {
		return "", errors.New("empty translated text")
	}

	return translated, nil
}

func normalizeTranslateLang(code string) string {
	normalized := strings.TrimSpace(strings.ReplaceAll(strings.ToLower(code), "_", "-"))
	switch normalized {
	case "", "auto":
		return ""
	case "zh", "zh-cn", "zh-hans":
		return "zh-CN"
	case "zh-tw", "zh-hk", "zh-hant":
		return "zh-TW"
	default:
		return normalized
	}
}
