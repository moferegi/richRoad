package system

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	commonResp "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
)

func newExportTemplateTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	i18n.SetLang(ctx, i18n.LangZh)
	return ctx, recorder
}

func decodeResponse(t *testing.T, recorder *httptest.ResponseRecorder) commonResp.Response {
	t.Helper()
	var resp commonResp.Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response failed: %v", err)
	}
	return resp
}

func TestFailExportTemplateWithErr_UnknownErrorUsesFallback(t *testing.T) {
	ctx, recorder := newExportTemplateTestContext()

	rawErr := errors.New("internal sql parser panic")
	failExportTemplateWithErr(ctx, rawErr, "getFail")

	resp := decodeResponse(t, recorder)
	if resp.Msg != i18n.T(ctx, "getFail") {
		t.Fatalf("expected fallback message %q, got %q", i18n.T(ctx, "getFail"), resp.Msg)
	}
	if resp.Msg == rawErr.Error() {
		t.Fatalf("unexpected raw error leakage: %q", resp.Msg)
	}
}

func TestFailExportTemplateWithErr_KnownKeyTranslated(t *testing.T) {
	ctx, recorder := newExportTemplateTestContext()

	failExportTemplateWithErr(ctx, errors.New("exportTokenRequired"), "fail")

	resp := decodeResponse(t, recorder)
	expected := i18n.T(ctx, "exportTokenRequired")
	if resp.Msg != expected {
		t.Fatalf("expected translated key message %q, got %q", expected, resp.Msg)
	}
}

func TestFailExportTemplateWithKey_InvalidKeyUsesFallback(t *testing.T) {
	ctx, recorder := newExportTemplateTestContext()

	failExportTemplateWithKey(ctx, "unknown_not_defined_key", "invalidParams")

	resp := decodeResponse(t, recorder)
	expected := i18n.T(ctx, "invalidParams")
	if resp.Msg != expected {
		t.Fatalf("expected fallback key message %q, got %q", expected, resp.Msg)
	}
}
