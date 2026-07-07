package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
)

var i18nPayloadFieldKeys = map[string]struct{}{
	"name":                {},
	"title":               {},
	"content":             {},
	"countryName":         {},
	"seriesName":          {},
	"episodeName":         {},
	"description":         {},
	"explanation":         {},
	"translate":           {},
	"announcementContent": {},
}

func LocalizeI18nPayloadByContext(c *gin.Context, payload interface{}) interface{} {
	if payload == nil {
		return payload
	}
	if !shouldLocalizeForUniClient(c) {
		if c != nil {
			c.Header("X-I18n-Localized", "0")
		}
		return payload
	}
	lang := i18n.GetLang(c)
	if strings.TrimSpace(lang) == "" {
		lang = "zh"
	}
	if c != nil {
		c.Header("X-I18n-Localized", "1")
		c.Header("X-I18n-Lang", normalizeLocale(lang))
	}
	return localizeI18nAny(payload, lang)
}

func shouldLocalizeForUniClient(c *gin.Context) bool {
	if c == nil || c.Request == nil {
		return false
	}
	platform := strings.ToLower(strings.TrimSpace(c.GetHeader("X-Client-Platform")))
	if platform == "web" || platform == "h5" {
		return false
	}
	if platform == "uni" || platform == "uniapp" || platform == "uni-app" {
		return true
	}

	// 兜底：uni 端学习接口默认带 includeI18n=1，避免代理丢头导致不裁剪
	if strings.TrimSpace(c.Query("includeI18n")) == "1" {
		path := strings.ToLower(strings.TrimSpace(c.Request.URL.Path))
		if strings.Contains(path, "/englishlearning/") || strings.HasSuffix(path, "/popup/getactivepopups") {
			return true
		}
	}
	return false
}

func localizeI18nAny(value interface{}, lang string) interface{} {
	switch v := value.(type) {
	case nil, string, bool, float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, json.Number:
		return v
	case map[string]interface{}:
		if isI18nMap(v) {
			return pickI18nText(v, lang)
		}
		out := make(map[string]interface{}, len(v))
		for key, raw := range v {
			if _, ok := i18nPayloadFieldKeys[key]; ok {
				out[key] = localizeI18nFieldValue(raw, lang)
				continue
			}
			out[key] = localizeI18nAny(raw, lang)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(v))
		for i := range v {
			out[i] = localizeI18nAny(v[i], lang)
		}
		return out
	default:
		raw, err := json.Marshal(value)
		if err != nil {
			return value
		}
		var any interface{}
		if err = json.Unmarshal(raw, &any); err != nil {
			return value
		}
		switch any.(type) {
		case map[string]interface{}, []interface{}:
			return localizeI18nAny(any, lang)
		default:
			return any
		}
	}
}

func localizeI18nFieldValue(raw interface{}, lang string) interface{} {
	switch v := raw.(type) {
	case map[string]interface{}:
		if isI18nMap(v) {
			return pickI18nText(v, lang)
		}
		return localizeI18nAny(v, lang)
	case string:
		trimmed := strings.TrimSpace(v)
		if !strings.HasPrefix(trimmed, "{") {
			return v
		}
		obj := map[string]interface{}{}
		if err := json.Unmarshal([]byte(trimmed), &obj); err != nil {
			return v
		}
		if !isI18nMap(obj) {
			return v
		}
		return pickI18nText(obj, lang)
	default:
		return localizeI18nAny(v, lang)
	}
}

func isI18nMap(obj map[string]interface{}) bool {
	if len(obj) == 0 || len(obj) > 16 {
		return false
	}
	langKeyCount := 0
	for key := range obj {
		if !isLocaleKey(key) {
			return false
		}
		langKeyCount++
	}
	return langKeyCount > 0
}

func isLocaleKey(key string) bool {
	if key == "" {
		return false
	}
	normalized := normalizeLocale(key)
	parts := strings.Split(normalized, "-")
	if len(parts) == 0 || len(parts) > 2 {
		return false
	}
	if !isLocaleSegment(parts[0], 2, 3) {
		return false
	}
	if len(parts) == 2 && !isLocaleSegment(parts[1], 2, 4) {
		return false
	}
	return true
}

func isLocaleSegment(segment string, minLen int, maxLen int) bool {
	if len(segment) < minLen || len(segment) > maxLen {
		return false
	}
	for _, ch := range segment {
		if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') {
			return false
		}
	}
	return true
}

func pickI18nText(obj map[string]interface{}, lang string) string {
	normalizedLang := normalizeLocale(lang)
	if text := readMapString(obj, normalizedLang); text != "" {
		return text
	}
	if idx := strings.Index(normalizedLang, "-"); idx > 0 {
		if text := readMapString(obj, normalizedLang[:idx]); text != "" {
			return text
		}
	}
	if text := readMapString(obj, "zh-tw"); text != "" && normalizedLang == "zh-tw" {
		return text
	}
	if text := readMapString(obj, "en"); text != "" {
		return text
	}
	if text := readMapString(obj, "zh"); text != "" {
		return text
	}
	for _, raw := range obj {
		text := strings.TrimSpace(toString(raw))
		if text != "" {
			return text
		}
	}
	return ""
}

func readMapString(obj map[string]interface{}, key string) string {
	if key == "" {
		return ""
	}
	normalizedKey := normalizeLocale(key)
	for mapKey, mapValue := range obj {
		if normalizeLocale(mapKey) != normalizedKey {
			continue
		}
		return strings.TrimSpace(toString(mapValue))
	}
	return ""
}

func normalizeLocale(value string) string {
	locale := strings.ToLower(strings.TrimSpace(value))
	locale = strings.ReplaceAll(locale, "_", "-")
	return locale
}

func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case json.Number:
		return v.String()
	default:
		return fmt.Sprint(v)
	}
}
