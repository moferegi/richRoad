package middleware

import (
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
)

// Locale 从请求头 Accept-Language 中提取语言并存入 gin.Context
func Locale() gin.HandlerFunc {
	supported := map[string]bool{
		i18n.LangZh: true,
		"zh-tw":     true,
		i18n.LangEn: true,
		i18n.LangMn: true,
		"th":        true,
		"hi":        true,
		"id":        true,
		i18n.LangVi: true,
		i18n.LangAr: true,
		i18n.LangJa: true,
		i18n.LangKo: true,
		i18n.LangMs: true,
	}
	return func(c *gin.Context) {
		lang := strings.TrimSpace(c.GetHeader("Accept-Language"))
		// 取首个语言标签（如 "en-US,en;q=0.9" -> "en-us"，保留地区后缀）
		if idx := strings.Index(lang, ","); idx > 0 {
			lang = lang[:idx]
		}
		if idx := strings.Index(lang, ";"); idx > 0 {
			lang = lang[:idx]
		}
		lang = strings.ToLower(lang)
		lang = strings.ReplaceAll(lang, "_", "-")

		// 常见变体归一，保证 zh-TW 可正确透传为 zh-tw
		switch lang {
		case "zh-hant", "zh-hk", "zh-mo":
			lang = "zh-tw"
		case "zh-hans", "zh-cn", "zh-sg":
			lang = i18n.LangZh
		}
		if !supported[lang] {
			// 对 en-us / ja-jp 这类前缀做回退
			if idx := strings.Index(lang, "-"); idx > 0 {
				base := lang[:idx]
				if supported[base] {
					lang = base
				} else {
					lang = i18n.LangZh
				}
			} else {
				lang = i18n.LangZh
			}
		}
		i18n.SetLang(c, lang)
		c.Next()
	}
}
