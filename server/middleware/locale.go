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
		i18n.LangEn: true,
		i18n.LangMn: true,
		i18n.LangVi: true,
		i18n.LangAr: true,
		i18n.LangJa: true,
		i18n.LangKo: true,
		i18n.LangMs: true,
	}
	return func(c *gin.Context) {
		lang := strings.TrimSpace(c.GetHeader("Accept-Language"))
		// 取第一个语言标签（如 "en-US,en;q=0.9" → "en"）
		if idx := strings.IndexAny(lang, ",-;"); idx > 0 {
			lang = lang[:idx]
		}
		lang = strings.ToLower(lang)
		if !supported[lang] {
			lang = i18n.LangZh
		}
		i18n.SetLang(c, lang)
		c.Next()
	}
}
