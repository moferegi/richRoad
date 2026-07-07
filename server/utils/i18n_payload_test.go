package utils

import (
	"net/http/httptest"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
)

func TestLocalizeI18nPayloadByContext_UniHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/api/englishLearning/content/getVideoCategoryList", nil)
	req.Header.Set("X-Client-Platform", "uni")
	c.Request = req
	i18n.SetLang(c, "mn")

	payload := map[string]interface{}{
		"name": `{"zh":"中文分类","mn":"Монгол ангилал"}`,
	}
	got := LocalizeI18nPayloadByContext(c, payload).(map[string]interface{})
	if got["name"] != "Монгол ангилал" {
		t.Fatalf("expected mn text, got %v", got["name"])
	}
	if w.Header().Get("X-I18n-Localized") != "1" {
		t.Fatalf("expected localized header 1")
	}
}

func TestLocalizeI18nPayloadByContext_IncludeI18nFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/api/englishLearning/video/getSentenceList?includeI18n=1", nil)
	c.Request = req
	i18n.SetLang(c, "mn")

	payload := []interface{}{
		map[string]interface{}{
			"translate": map[string]interface{}{"zh": "中文", "mn": "Монгол"},
		},
	}
	got := LocalizeI18nPayloadByContext(c, payload).([]interface{})
	row := got[0].(map[string]interface{})
	if row["translate"] != "Монгол" {
		t.Fatalf("expected mn translate, got %v", row["translate"])
	}
}

func TestLocalizeI18nPayloadByContext_WebBypass(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/api/englishLearning/content/getVideoCategoryList?includeI18n=1", nil)
	req.Header.Set("X-Client-Platform", "web")
	c.Request = req
	i18n.SetLang(c, "mn")

	payload := map[string]interface{}{
		"name": `{"zh":"中文分类","mn":"Монгол ангилал"}`,
	}
	got := LocalizeI18nPayloadByContext(c, payload).(map[string]interface{})
	if got["name"] != `{"zh":"中文分类","mn":"Монгол ангилал"}` {
		t.Fatalf("expected raw json string for web, got %v", got["name"])
	}
	if w.Header().Get("X-I18n-Localized") != "0" {
		t.Fatalf("expected localized header 0")
	}
}
