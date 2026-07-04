package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
)

type SysLanguageService struct{}

const (
	translateTimeoutMSConfigKey = "language_translate_timeout_ms"
	translateTimeoutMSEnvKey    = "LANGUAGE_TRANSLATE_TIMEOUT_MS"
	defaultTranslateTimeoutMS   = int64(45000)
	minTranslateTimeoutMS       = int64(5000)
	maxTranslateTimeoutMS       = int64(180000)

	translateRetryConfigKey = "language_translate_retry_count"
	translateRetryEnvKey    = "LANGUAGE_TRANSLATE_RETRY_COUNT"
	defaultTranslateRetry   = int64(2)
	maxTranslateRetry       = int64(5)

	translateEndpointsConfigKey = "language_translate_endpoints"
	translateEndpointsEnvKey    = "LANGUAGE_TRANSLATE_ENDPOINTS"

	translateProviderURLConfigKey = "language_translate_provider_url"
	translateProviderURLEnvKey    = "LANGUAGE_TRANSLATE_PROVIDER_URL"

	translateProviderAPIKeyConfigKey = "language_translate_provider_api_key"
	translateProviderAPIKeyEnvKey    = "LANGUAGE_TRANSLATE_PROVIDER_API_KEY"

	translateProviderAPIKeyHeaderConfigKey = "language_translate_provider_api_key_header"
	translateProviderAPIKeyHeaderEnvKey    = "LANGUAGE_TRANSLATE_PROVIDER_API_KEY_HEADER"
)

var defaultTranslateEndpoints = []string{
	"https://translate.googleapis.com/translate_a/single",
	"https://translate.google.com/translate_a/single",
	"https://translate.google.com.hk/translate_a/single",
}

type translateHTTPError struct {
	StatusCode int
	Body       string
}

type translateProviderConfig struct {
	URL          string
	APIKey       string
	APIKeyHeader string
}

func (c translateProviderConfig) enabled() bool {
	return strings.TrimSpace(c.URL) != ""
}

func (e *translateHTTPError) Error() string {
	if e == nil {
		return "translate request failed"
	}
	if strings.TrimSpace(e.Body) == "" {
		return fmt.Sprintf("translate request failed: status=%d", e.StatusCode)
	}
	return fmt.Sprintf("translate request failed: status=%d body=%s", e.StatusCode, e.Body)
}

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

// TranslateI18n 优先使用 Google 网关翻译，失败时可回退外部翻译服务
func (s *SysLanguageService) TranslateI18n(req clientReq.TranslateI18nRequest) (map[string]string, error) {
	text := strings.TrimSpace(req.Text)
	if text == "" {
		return nil, errors.New("empty text")
	}

	sourceLang := normalizeTranslateLang(req.Source)
	if sourceLang == "" {
		sourceLang = "zh-CN"
	}

	httpClient := &http.Client{Timeout: loadTranslateTimeout()}
	endpoints := loadTranslateEndpoints()
	retryCount := loadTranslateRetryCount()
	providerConfig := loadTranslateProviderConfig()
	translations := make(map[string]string, len(req.Targets))
	seen := make(map[string]struct{}, len(req.Targets))
	failedTargets := make([]string, 0)
	var lastErr error

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

		translated, err := s.translateByGoogleWithFallback(httpClient, endpoints, retryCount, text, sourceLang, normalizedTarget)
		if err != nil {
			if providerConfig.enabled() {
				providerTranslated, providerErr := s.translateByExternalProvider(httpClient, providerConfig, text, sourceLang, normalizedTarget)
				if providerErr == nil {
					global.GVA_LOG.Warn("Google翻译失败，已回退到外部翻译服务",
						zap.String("source", sourceLang),
						zap.String("target", normalizedTarget),
						zap.Error(err),
					)
					translations[rawTarget] = providerTranslated
					continue
				}

				lastErr = fmt.Errorf("google=%v; external=%v", err, providerErr)
			} else {
				lastErr = err
			}

			failedTargets = append(failedTargets, rawTarget)
			continue
		}
		translations[rawTarget] = translated
	}

	if len(failedTargets) > 0 {
		global.GVA_LOG.Warn("部分目标语言翻译失败",
			zap.String("source", sourceLang),
			zap.Strings("failedTargets", failedTargets),
			zap.Error(lastErr),
		)
	}

	if len(translations) == 0 {
		if lastErr != nil {
			return nil, fmt.Errorf("translate service unavailable: %w", lastErr)
		}
		return nil, errors.New("no available target languages")
	}

	return translations, nil
}

func (s *SysLanguageService) translateByGoogleWithFallback(httpClient *http.Client, endpoints []string, retryCount int, text, sourceLang, targetLang string) (string, error) {
	if len(endpoints) == 0 {
		endpoints = defaultTranslateEndpoints
	}
	if retryCount < 1 {
		retryCount = 1
	}

	var lastErr error
	for _, endpoint := range endpoints {
		for attempt := 1; attempt <= retryCount; attempt++ {
			translated, err := s.translateByGoogle(httpClient, endpoint, text, sourceLang, targetLang)
			if err == nil {
				return translated, nil
			}

			lastErr = err
			if !isRetryableTranslateError(err) {
				break
			}
			if attempt < retryCount {
				time.Sleep(time.Duration(attempt) * 200 * time.Millisecond)
			}
		}
	}

	if lastErr == nil {
		lastErr = errors.New("translate endpoint unavailable")
	}
	return "", lastErr
}

func (s *SysLanguageService) translateByExternalProvider(httpClient *http.Client, cfg translateProviderConfig, text, sourceLang, targetLang string) (string, error) {
	if !cfg.enabled() {
		return "", errors.New("external translate provider not configured")
	}

	payload := map[string]any{
		"text":    text,
		"source":  sourceLang,
		"target":  targetLang,
		"targets": []string{targetLang},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, cfg.URL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "gin-vue-admin/translate-fallback")

	if strings.TrimSpace(cfg.APIKey) != "" {
		header := strings.TrimSpace(cfg.APIKeyHeader)
		if header == "" {
			header = "Authorization"
		}
		if strings.EqualFold(header, "Authorization") {
			req.Header.Set(header, "Bearer "+cfg.APIKey)
		} else {
			req.Header.Set(header, cfg.APIKey)
			if req.Header.Get("Authorization") == "" {
				req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
			}
		}
		req.Header.Set("X-API-Key", cfg.APIKey)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", &translateHTTPError{StatusCode: resp.StatusCode, Body: string(raw)}
	}

	translated, err := parseExternalTranslateResult(raw, targetLang)
	if err != nil {
		return "", err
	}

	return translated, nil
}

func (s *SysLanguageService) translateByGoogle(httpClient *http.Client, endpoint, text, sourceLang, targetLang string) (string, error) {
	endpointURL, err := url.Parse(strings.TrimSpace(endpoint))
	if err != nil {
		return "", err
	}

	params := endpointURL.Query()
	params.Set("client", "gtx")
	params.Set("sl", sourceLang)
	params.Set("tl", targetLang)
	params.Set("dt", "t")
	params.Set("q", text)
	endpointURL.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, endpointURL.String(), nil)
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
		return "", &translateHTTPError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	result, err := parseGoogleTranslateResult(body)
	if err != nil {
		return "", err
	}

	return result, nil
}

func isRetryableTranslateError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() || netErr.Temporary() {
			return true
		}
	}

	var httpErr *translateHTTPError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode == http.StatusTooManyRequests || httpErr.StatusCode >= http.StatusInternalServerError
	}

	errText := strings.ToLower(err.Error())
	if strings.Contains(errText, "timeout") || strings.Contains(errText, "tempor") {
		return true
	}

	return false
}

func loadTranslateTimeout() time.Duration {
	timeoutMs := utils.GetInt64Setting(translateTimeoutMSConfigKey, translateTimeoutMSEnvKey, defaultTranslateTimeoutMS)
	if timeoutMs < minTranslateTimeoutMS {
		timeoutMs = minTranslateTimeoutMS
	}
	if timeoutMs > maxTranslateTimeoutMS {
		timeoutMs = maxTranslateTimeoutMS
	}
	return time.Duration(timeoutMs) * time.Millisecond
}

func loadTranslateRetryCount() int {
	retry := utils.GetInt64Setting(translateRetryConfigKey, translateRetryEnvKey, defaultTranslateRetry)
	if retry < 1 {
		retry = 1
	}
	if retry > maxTranslateRetry {
		retry = maxTranslateRetry
	}
	return int(retry)
}

func loadTranslateEndpoints() []string {
	raw := getStringSetting(translateEndpointsConfigKey, translateEndpointsEnvKey)
	if strings.TrimSpace(raw) == "" {
		return append([]string(nil), defaultTranslateEndpoints...)
	}

	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == '\n' || r == ';'
	})

	result := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, item := range parts {
		candidate := strings.TrimSpace(item)
		if candidate == "" {
			continue
		}
		if _, err := url.ParseRequestURI(candidate); err != nil {
			continue
		}
		if _, exists := seen[candidate]; exists {
			continue
		}
		seen[candidate] = struct{}{}
		result = append(result, candidate)
	}

	if len(result) == 0 {
		return append([]string(nil), defaultTranslateEndpoints...)
	}

	return result
}

func loadTranslateProviderConfig() translateProviderConfig {
	providerURL := getStringSetting(translateProviderURLConfigKey, translateProviderURLEnvKey)
	if strings.TrimSpace(providerURL) == "" {
		// 兼容用户已接入的 TTS 网关：若该网关同时提供翻译接口，可直接复用。
		providerURL = getStringSetting("learning_tts_provider_url", "LEARNING_TTS_PROVIDER_URL")
	}

	apiKey := getStringSetting(translateProviderAPIKeyConfigKey, translateProviderAPIKeyEnvKey)
	if strings.TrimSpace(apiKey) == "" {
		apiKey = getStringSetting("learning_tts_api_key", "LEARNING_TTS_API_KEY")
	}

	apiKeyHeader := getStringSetting(translateProviderAPIKeyHeaderConfigKey, translateProviderAPIKeyHeaderEnvKey)
	if strings.TrimSpace(apiKeyHeader) == "" {
		apiKeyHeader = "Authorization"
	}

	return translateProviderConfig{
		URL:          strings.TrimSpace(providerURL),
		APIKey:       strings.TrimSpace(apiKey),
		APIKeyHeader: strings.TrimSpace(apiKeyHeader),
	}
}

func getStringSetting(sysConfigKey, envKey string) string {
	if value, ok := utils.GetSysConfigRawValue(sysConfigKey); ok {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}

	if strings.TrimSpace(envKey) == "" {
		return ""
	}

	return strings.TrimSpace(os.Getenv(envKey))
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

func parseExternalTranslateResult(body []byte, targetLang string) (string, error) {
	var payload any
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}

	translated := extractExternalTranslatedText(payload, targetLang, 0)
	if translated == "" {
		return "", errors.New("external provider missing translated text")
	}

	return translated, nil
}

func extractExternalTranslatedText(payload any, targetLang string, depth int) string {
	if depth > 5 || payload == nil {
		return ""
	}

	switch typed := payload.(type) {
	case string:
		return strings.TrimSpace(typed)
	case []any:
		for _, item := range typed {
			if text := extractExternalTranslatedText(item, targetLang, depth+1); text != "" {
				return text
			}
		}
	case map[string]any:
		if text := extractMapTextByLanguageKey(typed, targetLang); text != "" {
			return text
		}

		for _, key := range []string{"translation", "translatedText", "translated_text", "targetText", "target_text", "result"} {
			if value, ok := mapLookupCaseInsensitive(typed, key); ok {
				if text := strings.TrimSpace(anyToString(value)); text != "" {
					return text
				}
			}
		}

		if value, ok := mapLookupCaseInsensitive(typed, "translations"); ok {
			if text := extractExternalTranslatedText(value, targetLang, depth+1); text != "" {
				return text
			}
		}

		for _, key := range []string{"data", "payload", "response", "result"} {
			if value, ok := mapLookupCaseInsensitive(typed, key); ok {
				if text := extractExternalTranslatedText(value, targetLang, depth+1); text != "" {
					return text
				}
			}
		}
	}

	return ""
}

func extractMapTextByLanguageKey(data map[string]any, targetLang string) string {
	normalizedTarget := normalizeExternalLanguageKey(targetLang)
	if normalizedTarget == "" {
		return ""
	}

	for key, value := range data {
		if normalizeExternalLanguageKey(key) != normalizedTarget {
			continue
		}
		if text := strings.TrimSpace(anyToString(value)); text != "" {
			return text
		}
	}

	return ""
}

func normalizeExternalLanguageKey(raw string) string {
	normalized := strings.TrimSpace(strings.ToLower(raw))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	return normalized
}

func mapLookupCaseInsensitive(data map[string]any, key string) (any, bool) {
	for k, v := range data {
		if strings.EqualFold(strings.TrimSpace(k), strings.TrimSpace(key)) {
			return v, true
		}
	}
	return nil, false
}

func anyToString(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case fmt.Stringer:
		return typed.String()
	case json.Number:
		return typed.String()
	default:
		return strings.TrimSpace(fmt.Sprintf("%v", typed))
	}
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
