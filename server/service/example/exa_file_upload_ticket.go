package example

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	scanUploadTicketTTLConfigKey = "security_scan_upload_ticket_ttl_seconds"
	scanUploadTicketTTLEnvKey    = "CS_SCAN_UPLOAD_TICKET_TTL_SECONDS"

	defaultScanUploadTicketTTLSeconds int64 = 180
	minScanUploadTicketTTLSeconds     int64 = 60
	maxScanUploadTicketTTLSeconds     int64 = 900

	scanUploadTicketRedisKeyPrefix = "scan_upload:ticket:"

	scanUploadTicketUploadRateLimitConfigKey = "security_scan_upload_ticket_upload_rate_limit_per_ip"
	scanUploadTicketUploadRateLimitEnvKey    = "CS_SCAN_UPLOAD_TICKET_UPLOAD_RATE_LIMIT_PER_IP"
	defaultScanUploadTicketUploadRateLimit   = int64(30)

	scanUploadTicketUploadRateWindowConfigKey = "security_scan_upload_ticket_upload_rate_limit_window_seconds"
	scanUploadTicketUploadRateWindowEnvKey    = "CS_SCAN_UPLOAD_TICKET_UPLOAD_RATE_LIMIT_WINDOW_SECONDS"
	defaultScanUploadTicketUploadRateWindow   = int64(60)

	scanUploadTicketUploadRateRedisKeyPrefix = "rate:scan_upload:ticket_upload:"
)

var (
	scanUploadTicketPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{20,128}$`)
	errScanUploadTicket     = errors.New("上传凭证无效或已过期")
	errScanUploadIssuer     = errors.New("用户未登录")

	scanUploadTicketConsumeScript = redis.NewScript(`
local value = redis.call('GET', KEYS[1])
if not value then
  return nil
end
redis.call('DEL', KEYS[1])
return value
`)

	localScanUploadTicketStore sync.Map

	localScanUploadTicketUploadRateLimiter = struct {
		mu       sync.Mutex
		counters map[string]localScanUploadTicketRateCounter
	}{
		counters: make(map[string]localScanUploadTicketRateCounter),
	}
)

type ScanUploadTicketPayload struct {
	ClassId        int    `json:"classId"`
	Folder         string `json:"folder"`
	UploadType     string `json:"uploadType"`
	UploadPosition string `json:"uploadPosition"`
	OperatorUserID uint   `json:"operatorUserId"`
	ExpiresAt      int64  `json:"expiresAt"`
}

type localScanUploadTicketEntry struct {
	Payload   ScanUploadTicketPayload
	ExpiresAt time.Time
}

type localScanUploadTicketRateCounter struct {
	Count     int64
	ExpiresAt time.Time
}

func loadScanUploadTicketTTL() time.Duration {
	seconds := utils.GetInt64Setting(scanUploadTicketTTLConfigKey, scanUploadTicketTTLEnvKey, defaultScanUploadTicketTTLSeconds)
	if seconds < minScanUploadTicketTTLSeconds {
		seconds = minScanUploadTicketTTLSeconds
	}
	if seconds > maxScanUploadTicketTTLSeconds {
		seconds = maxScanUploadTicketTTLSeconds
	}
	return time.Duration(seconds) * time.Second
}

func buildScanUploadTicketUploadRateIdentity(ip string, userAgent string) string {
	normalizedIP := strings.TrimSpace(ip)
	if normalizedIP == "" {
		normalizedIP = "unknown"
	}

	normalizedUA := strings.ToLower(strings.TrimSpace(userAgent))
	if normalizedUA == "" {
		normalizedUA = "unknown"
	} else {
		normalizedUA = trimRunes(normalizedUA, 128)
	}

	sum := sha1.Sum([]byte(normalizedIP + "|" + normalizedUA))
	return hex.EncodeToString(sum[:])
}

func checkScanUploadTicketUploadRateLimitLocal(identity string, limit int64, window time.Duration) (allowed bool, waitSeconds int) {
	if limit <= 0 {
		return true, 0
	}
	if window <= 0 {
		window = time.Duration(defaultScanUploadTicketUploadRateWindow) * time.Second
	}

	now := time.Now()
	localScanUploadTicketUploadRateLimiter.mu.Lock()
	defer localScanUploadTicketUploadRateLimiter.mu.Unlock()

	if len(localScanUploadTicketUploadRateLimiter.counters) > 10000 {
		for key, counter := range localScanUploadTicketUploadRateLimiter.counters {
			if now.After(counter.ExpiresAt) {
				delete(localScanUploadTicketUploadRateLimiter.counters, key)
			}
		}
	}

	counter, ok := localScanUploadTicketUploadRateLimiter.counters[identity]
	if !ok || now.After(counter.ExpiresAt) {
		localScanUploadTicketUploadRateLimiter.counters[identity] = localScanUploadTicketRateCounter{
			Count:     1,
			ExpiresAt: now.Add(window),
		}
		return true, 0
	}

	if counter.Count >= limit {
		remain := int(time.Until(counter.ExpiresAt).Seconds())
		if remain < 1 {
			remain = 1
		}
		return false, remain
	}

	counter.Count++
	localScanUploadTicketUploadRateLimiter.counters[identity] = counter
	return true, 0
}

func (e *FileUploadAndDownloadService) CheckScanUploadTicketUploadRateLimit(ip string, userAgent string) (allowed bool, waitSeconds int) {
	limit := utils.GetInt64Setting(scanUploadTicketUploadRateLimitConfigKey, scanUploadTicketUploadRateLimitEnvKey, defaultScanUploadTicketUploadRateLimit)
	if limit <= 0 {
		return true, 0
	}

	windowSeconds := utils.GetInt64Setting(scanUploadTicketUploadRateWindowConfigKey, scanUploadTicketUploadRateWindowEnvKey, defaultScanUploadTicketUploadRateWindow)
	if windowSeconds <= 0 {
		windowSeconds = defaultScanUploadTicketUploadRateWindow
	}
	window := time.Duration(windowSeconds) * time.Second
	identity := buildScanUploadTicketUploadRateIdentity(ip, userAgent)

	if global.GVA_REDIS == nil {
		return checkScanUploadTicketUploadRateLimitLocal(identity, limit, window)
	}

	key := scanUploadTicketUploadRateRedisKeyPrefix + identity
	count, err := global.GVA_REDIS.Incr(context.Background(), key).Result()
	if err != nil {
		return checkScanUploadTicketUploadRateLimitLocal(identity, limit, window)
	}
	if count == 1 {
		_ = global.GVA_REDIS.Expire(context.Background(), key, window).Err()
	}
	if count > limit {
		ttl, ttlErr := global.GVA_REDIS.TTL(context.Background(), key).Result()
		if ttlErr != nil || ttl <= 0 {
			return false, int(window.Seconds())
		}
		wait := int(ttl.Seconds())
		if wait < 1 {
			wait = 1
		}
		return false, wait
	}

	return true, 0
}

func generateScanUploadTicket() (string, error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func normalizeScanUploadTicket(raw string) (string, error) {
	ticket := strings.TrimSpace(raw)
	if !scanUploadTicketPattern.MatchString(ticket) {
		return "", errScanUploadTicket
	}
	return ticket, nil
}

func (e *FileUploadAndDownloadService) CreateScanUploadTicket(classId int, folder string, uploadType string, uploadPosition string, operatorUserID uint) (ticket string, expireAt time.Time, err error) {
	if operatorUserID == 0 {
		return "", time.Time{}, errScanUploadIssuer
	}
	if classId < 0 {
		classId = 0
	}

	ticket, err = generateScanUploadTicket()
	if err != nil {
		return "", time.Time{}, err
	}

	ttl := loadScanUploadTicketTTL()
	expireAt = time.Now().Add(ttl)
	payload := ScanUploadTicketPayload{
		ClassId:        classId,
		Folder:         strings.TrimSpace(folder),
		UploadType:     trimRunes(strings.TrimSpace(uploadType), 64),
		UploadPosition: trimRunes(strings.TrimSpace(uploadPosition), 64),
		OperatorUserID: operatorUserID,
		ExpiresAt:      expireAt.UnixMilli(),
	}

	encoded, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return "", time.Time{}, marshalErr
	}

	if global.GVA_REDIS != nil {
		key := scanUploadTicketRedisKeyPrefix + ticket
		if redisErr := global.GVA_REDIS.Set(context.Background(), key, encoded, ttl).Err(); redisErr == nil {
			return ticket, expireAt, nil
		} else {
			global.GVA_LOG.Warn("store scan upload ticket in redis failed, fallback to local cache", zap.Error(redisErr))
		}
	}

	localScanUploadTicketStore.Store(ticket, localScanUploadTicketEntry{Payload: payload, ExpiresAt: expireAt})
	return ticket, expireAt, nil
}

func (e *FileUploadAndDownloadService) ConsumeScanUploadTicket(rawTicket string) (ScanUploadTicketPayload, error) {
	var empty ScanUploadTicketPayload
	ticket, err := normalizeScanUploadTicket(rawTicket)
	if err != nil {
		return empty, err
	}

	if global.GVA_REDIS != nil {
		payload, consumeErr := consumeScanUploadTicketFromRedis(ticket)
		if consumeErr == nil {
			return payload, nil
		}
		if !errors.Is(consumeErr, redis.Nil) {
			global.GVA_LOG.Warn("consume scan upload ticket from redis failed, fallback to local cache", zap.Error(consumeErr))
		}
	}

	entryValue, ok := localScanUploadTicketStore.LoadAndDelete(ticket)
	if !ok {
		return empty, errScanUploadTicket
	}
	entry, ok := entryValue.(localScanUploadTicketEntry)
	if !ok || time.Now().After(entry.ExpiresAt) || entry.Payload.ExpiresAt <= 0 || time.Now().UnixMilli() > entry.Payload.ExpiresAt {
		return empty, errScanUploadTicket
	}
	return entry.Payload, nil
}

func consumeScanUploadTicketFromRedis(ticket string) (ScanUploadTicketPayload, error) {
	var payload ScanUploadTicketPayload
	key := scanUploadTicketRedisKeyPrefix + ticket
	value, err := scanUploadTicketConsumeScript.Run(context.Background(), global.GVA_REDIS, []string{key}).Text()
	if err != nil {
		return payload, err
	}
	if jsonErr := json.Unmarshal([]byte(value), &payload); jsonErr != nil {
		return payload, errScanUploadTicket
	}
	if payload.ExpiresAt <= 0 || time.Now().UnixMilli() > payload.ExpiresAt {
		return payload, errScanUploadTicket
	}
	return payload, nil
}
