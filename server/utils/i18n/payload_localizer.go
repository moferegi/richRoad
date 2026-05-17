package i18n

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
)

var localizedLangKeys = map[string]struct{}{
	LangZh:  {},
	LangEn:  {},
	LangMn:  {},
	LangVi:  {},
	LangAr:  {},
	LangJa:  {},
	LangKo:  {},
	LangMs:  {},
	"zh-tw": {},
	"th":    {},
	"hi":    {},
	"id":    {},
}

// ShouldLocalizeResponseData 返回是否对当前请求执行单语化响应。
func ShouldLocalizeResponseData(c *gin.Context) bool {
	if c == nil {
		return false
	}

	if truthyParam(c.Query("includeI18n")) {
		return false
	}
	if strings.EqualFold(strings.TrimSpace(c.Query("i18nMode")), "full") {
		return false
	}
	if strings.EqualFold(strings.TrimSpace(c.GetHeader("X-I18n-Mode")), "full") {
		return false
	}
	if truthyParam(c.Query("langOnly")) || strings.EqualFold(strings.TrimSpace(c.GetHeader("X-I18n-Mode")), "localized") {
		return true
	}

	return strings.EqualFold(strings.TrimSpace(c.GetHeader("X-Client-App")), "uni")
}

// LocalizeResponseData 将多语言对象按当前请求语言裁剪为单语言。
func LocalizeResponseData(c *gin.Context, data interface{}) interface{} {
	if data == nil || !ShouldLocalizeResponseData(c) {
		return data
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return data
	}

	var generic interface{}
	if err := json.Unmarshal(raw, &generic); err != nil {
		return data
	}

	lang := normalizeLang(GetLang(c))
	if lang == "" {
		lang = LangZh
	}

	return localizePayloadValue(generic, lang)
}

func localizePayloadValue(value interface{}, lang string) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		if isLanguageObject(typed) {
			return pickLocalizedValue(typed, lang)
		}
		out := make(map[string]interface{}, len(typed))
		for key, item := range typed {
			out[key] = localizePayloadValue(item, lang)
		}
		return out
	case []interface{}:
		out := make([]interface{}, 0, len(typed))
		for _, item := range typed {
			out = append(out, localizePayloadValue(item, lang))
		}
		return out
	case string:
		trimmed := strings.TrimSpace(typed)
		if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
			var maybeObj map[string]interface{}
			if err := json.Unmarshal([]byte(trimmed), &maybeObj); err == nil && isLanguageObject(maybeObj) {
				return pickLocalizedValue(maybeObj, lang)
			}
		}
		return typed
	default:
		return value
	}
}

func pickLocalizedValue(values map[string]interface{}, lang string) interface{} {
	if v, ok := findLangValue(values, lang); ok {
		return localizePayloadValue(v, lang)
	}

	baseLang := lang
	if idx := strings.Index(baseLang, "-"); idx > 0 {
		baseLang = baseLang[:idx]
		if v, ok := findLangValue(values, baseLang); ok {
			return localizePayloadValue(v, lang)
		}
	}

	fallbacks := []string{LangEn, LangZh, LangMn, LangVi, LangAr, LangJa, LangKo, LangMs, "zh-tw", "th", "hi", "id"}
	for _, key := range fallbacks {
		if v, ok := findLangValue(values, key); ok {
			return localizePayloadValue(v, lang)
		}
	}

	for _, v := range values {
		if v != nil {
			return localizePayloadValue(v, lang)
		}
	}

	return ""
}

func findLangValue(values map[string]interface{}, lang string) (interface{}, bool) {
	target := normalizeLang(lang)
	for key, value := range values {
		if normalizeLang(key) == target {
			return value, true
		}
	}
	return nil, false
}

func isLanguageObject(values map[string]interface{}) bool {
	if len(values) == 0 {
		return false
	}

	hasKnownLang := false
	for key := range values {
		normalized := normalizeLang(key)
		if _, ok := localizedLangKeys[normalized]; ok {
			hasKnownLang = true
			continue
		}
		if !looksLikeLangKey(normalized) {
			return false
		}
	}
	return hasKnownLang
}

func normalizeLang(lang string) string {
	normalized := strings.ToLower(strings.TrimSpace(lang))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	return normalized
}

func looksLikeLangKey(key string) bool {
	if len(key) == 2 {
		return isAlpha(key[0]) && isAlpha(key[1])
	}
	if len(key) == 5 && key[2] == '-' {
		return isAlpha(key[0]) && isAlpha(key[1]) && isAlpha(key[3]) && isAlpha(key[4])
	}
	return false
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func truthyParam(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
