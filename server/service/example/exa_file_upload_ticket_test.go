package example

import (
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

func resetScanUploadTicketLocalTestState() {
	localScanUploadTicketStore = sync.Map{}
	localScanUploadTicketUploadRateLimiter.mu.Lock()
	localScanUploadTicketUploadRateLimiter.counters = make(map[string]localScanUploadTicketRateCounter)
	localScanUploadTicketUploadRateLimiter.mu.Unlock()
}

func TestScanUploadTicketConsumeOnce(t *testing.T) {
	previousRedis := global.GVA_REDIS
	global.GVA_REDIS = nil
	defer func() {
		global.GVA_REDIS = previousRedis
	}()
	resetScanUploadTicketLocalTestState()

	service := &FileUploadAndDownloadService{}
	ticket, _, err := service.CreateScanUploadTicket(1, "tryon/demo", "cloth", "tryon", 9528)
	if err != nil {
		t.Fatalf("CreateScanUploadTicket returned error: %v", err)
	}

	payload, err := service.ConsumeScanUploadTicket(ticket)
	if err != nil {
		t.Fatalf("first consume should succeed, got error: %v", err)
	}
	if payload.OperatorUserID != 9528 {
		t.Fatalf("unexpected operator id: %d", payload.OperatorUserID)
	}

	if _, err = service.ConsumeScanUploadTicket(ticket); err == nil {
		t.Fatalf("second consume should fail for one-time ticket")
	}
}

func TestScanUploadTicketInvalidAndExpired(t *testing.T) {
	previousRedis := global.GVA_REDIS
	global.GVA_REDIS = nil
	defer func() {
		global.GVA_REDIS = previousRedis
	}()
	resetScanUploadTicketLocalTestState()

	service := &FileUploadAndDownloadService{}
	if _, err := service.ConsumeScanUploadTicket("invalid@@ticket"); err == nil {
		t.Fatalf("invalid ticket format should be rejected")
	}

	expiredTicket := strings.Repeat("a", 24)
	expiredAt := time.Now().Add(-time.Minute)
	localScanUploadTicketStore.Store(expiredTicket, localScanUploadTicketEntry{
		Payload:   ScanUploadTicketPayload{ExpiresAt: expiredAt.UnixMilli()},
		ExpiresAt: expiredAt,
	})
	if _, err := service.ConsumeScanUploadTicket(expiredTicket); err == nil {
		t.Fatalf("expired ticket should be rejected")
	}
}

func TestScanUploadTicketUploadRateLimitLocal(t *testing.T) {
	previousRedis := global.GVA_REDIS
	global.GVA_REDIS = nil
	defer func() {
		global.GVA_REDIS = previousRedis
	}()
	resetScanUploadTicketLocalTestState()

	previousLimit := os.Getenv(scanUploadTicketUploadRateLimitEnvKey)
	previousWindow := os.Getenv(scanUploadTicketUploadRateWindowEnvKey)
	defer func() {
		_ = os.Setenv(scanUploadTicketUploadRateLimitEnvKey, previousLimit)
		_ = os.Setenv(scanUploadTicketUploadRateWindowEnvKey, previousWindow)
	}()
	_ = os.Setenv(scanUploadTicketUploadRateLimitEnvKey, "1")
	_ = os.Setenv(scanUploadTicketUploadRateWindowEnvKey, "60")

	service := &FileUploadAndDownloadService{}
	allowed, wait := service.CheckScanUploadTicketUploadRateLimit("127.0.0.1", "UA-A")
	if !allowed || wait != 0 {
		t.Fatalf("first request should pass, allowed=%v wait=%d", allowed, wait)
	}

	allowed, wait = service.CheckScanUploadTicketUploadRateLimit("127.0.0.1", "UA-A")
	if allowed || wait < 1 {
		t.Fatalf("second request should be rate-limited, allowed=%v wait=%d", allowed, wait)
	}

	allowed, wait = service.CheckScanUploadTicketUploadRateLimit("127.0.0.1", "UA-B")
	if !allowed || wait != 0 {
		t.Fatalf("different user-agent should have independent window, allowed=%v wait=%d", allowed, wait)
	}
}
