package middleware

import (
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
)

var supportedLocaleSet = map[string]struct{}{
	i18n.LangZh: {},
	i18n.LangEn: {},
	i18n.LangMn: {},
	i18n.LangVi: {},
	i18n.LangAr: {},
	i18n.LangJa: {},
	i18n.LangKo: {},
	i18n.LangMs: {},
	"zh-tw":     {},
	"th":        {},
	"hi":        {},
	"id":        {},
}

func normalizeLangTag(tag string) string {
	normalized := strings.ToLower(strings.TrimSpace(tag))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	return normalized
}

func resolveRequestLang(header string) string {
	if strings.TrimSpace(header) == "" {
		return i18n.LangZh
	}

	parts := strings.Split(header, ",")
	for _, part := range parts {
		candidate := strings.TrimSpace(part)
		if candidate == "" {
			continue
		}
		if semi := strings.Index(candidate, ";"); semi >= 0 {
			candidate = candidate[:semi]
		}

		normalized := normalizeLangTag(candidate)
		if normalized == "" {
			continue
		}

		if _, ok := supportedLocaleSet[normalized]; ok {
			return normalized
		}

		// 将常见繁体区域映射到 zh-tw，其余区域回退到基础语言。
		if normalized == "zh-hk" || normalized == "zh-mo" || normalized == "zh-hant" {
			return "zh-tw"
		}

		if idx := strings.Index(normalized, "-"); idx > 0 {
			base := normalized[:idx]
			if _, ok := supportedLocaleSet[base]; ok {
				return base
			}
		}
	}

	return i18n.LangZh
}

// Locale 从请求头 Accept-Language 中提取语言并存入 gin.Context
func Locale() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := resolveRequestLang(c.GetHeader("Accept-Language"))
		i18n.SetLang(c, lang)
		c.Next()
	}
}
