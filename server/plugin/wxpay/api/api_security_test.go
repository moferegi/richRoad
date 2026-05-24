package api

import (
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func newNotifyGuardContext() *gin.Context {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest("POST", "/wxpay/payAction", nil)
	return ctx
}

func setNotifyHeaders(ctx *gin.Context, ts int64, nonce string) {
	ctx.Request.Header.Set("Wechatpay-Timestamp", strconv.FormatInt(ts, 10))
	ctx.Request.Header.Set("Wechatpay-Nonce", nonce)
	ctx.Request.Header.Set("Wechatpay-Signature", "dummy-signature")
	ctx.Request.Header.Set("Wechatpay-Serial", "dummy-serial")
}

func TestValidateNotifyBasicReplayGuard_MissingHeaders(t *testing.T) {
	ctx := newNotifyGuardContext()
	if _, err := validateNotifyBasicReplayGuard(ctx); err == nil {
		t.Fatalf("expected error for missing headers")
	}
}

func TestValidateNotifyBasicReplayGuard_TimestampExpired(t *testing.T) {
	wxpayNotifyNonceCache = &notifyNonceCache{items: map[string]time.Time{}}
	ctx := newNotifyGuardContext()
	setNotifyHeaders(ctx, time.Now().Add(-20*time.Minute).Unix(), "nonce-expired")

	if _, err := validateNotifyBasicReplayGuard(ctx); err == nil {
		t.Fatalf("expected error for expired timestamp")
	}
}

func TestValidateNotifyBasicReplayGuard_ReplayNonce(t *testing.T) {
	wxpayNotifyNonceCache = &notifyNonceCache{items: map[string]time.Time{}}
	ctx1 := newNotifyGuardContext()
	setNotifyHeaders(ctx1, time.Now().Unix(), "nonce-replay")
	nonce, err := validateNotifyBasicReplayGuard(ctx1)
	if err != nil {
		t.Fatalf("unexpected first pass error: %v", err)
	}
	if err := markNotifyNonceUsed(nonce); err != nil {
		t.Fatalf("unexpected first nonce mark error: %v", err)
	}

	ctx2 := newNotifyGuardContext()
	setNotifyHeaders(ctx2, time.Now().Unix(), "nonce-replay")
	nonce, err = validateNotifyBasicReplayGuard(ctx2)
	if err != nil {
		t.Fatalf("unexpected second header validation error: %v", err)
	}
	if err := markNotifyNonceUsed(nonce); err == nil {
		t.Fatalf("expected replay nonce error")
	}
}
