package system

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
)

func TestGetSystemConfigRequiresSuperAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest("POST", "/system/getSystemConfig", nil)
	ctx.Set("claims", &request.CustomClaims{BaseClaims: request.BaseClaims{AuthorityId: 8881}})

	(&SystemApi{}).GetSystemConfig(ctx)

	if strings.Contains(recorder.Body.String(), "获取成功") || strings.Contains(recorder.Body.String(), "getSuccess") {
		t.Fatalf("expected 8881 to be denied full system config, got response: %s", recorder.Body.String())
	}
}

func TestSetSystemConfigRequiresSuperAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest("POST", "/system/setSystemConfig", strings.NewReader(`{}`))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set("claims", &request.CustomClaims{BaseClaims: request.BaseClaims{AuthorityId: 8881}})

	(&SystemApi{}).SetSystemConfig(ctx)

	if strings.Contains(recorder.Body.String(), "设置成功") || strings.Contains(recorder.Body.String(), "setSuccess") {
		t.Fatalf("expected 8881 to be denied setting system config, got response: %s", recorder.Body.String())
	}
}
