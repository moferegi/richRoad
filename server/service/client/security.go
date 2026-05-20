package client

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
)

type SecurityService struct{}

const (
	captchaRateLimitConfigKey = "captcha_rate_limit"
	captchaRateLimitEnvKey    = "CS_CAPTCHA_RATE_LIMIT"
	defaultCaptchaRateLimit   = int64(10)

	captchaRateWindowConfigKey = "captcha_rate_limit_window_seconds"
	captchaRateWindowEnvKey    = "CS_CAPTCHA_RATE_LIMIT_WINDOW_SECONDS"
	defaultCaptchaRateWindow   = int64(60)

	loginIPRateLimitConfigKey = "security_login_ip_rate_limit_per_minute"
	loginIPRateLimitEnvKey    = "CS_LOGIN_IP_RATE_LIMIT_PER_MINUTE"
	defaultLoginIPRateLimit   = int64(30)

	loginIPRateWindowConfigKey = "security_login_ip_rate_limit_window_seconds"
	loginIPRateWindowEnvKey    = "CS_LOGIN_IP_RATE_LIMIT_WINDOW_SECONDS"
	defaultLoginIPRateWindow   = int64(60)
)

type localCaptchaRateCounter struct {
	Count     int64
	ExpiresAt time.Time
}

var localCaptchaRateLimiter = struct {
	mu       sync.Mutex
	counters map[string]localCaptchaRateCounter
}{
	counters: make(map[string]localCaptchaRateCounter),
}

var localLoginIPRateLimiter = struct {
	mu       sync.Mutex
	counters map[string]localCaptchaRateCounter
}{
	counters: make(map[string]localCaptchaRateCounter),
}

// ---------- 辅助：读取 sysConfig ----------

func (s *SecurityService) getConfigInt(key string, def int) int {
	if strings.TrimSpace(key) == "" || global.GVA_DB == nil {
		return def
	}

	var cfg client.SysConfig
	if err := global.GVA_DB.Where("config_key = ?", key).First(&cfg).Error; err == nil {
		v, parseErr := strconv.Atoi(strings.TrimSpace(cfg.ConfigValue))
		if parseErr == nil {
			return v
		}
	}
	return def
}

// ========== 注册 IP 限制 ==========

// CheckRegisterIPLimit 检查IP注册限制（优先Redis，回退DB）
func (s *SecurityService) CheckRegisterIPLimit(ip string) (allowed bool, err error) {
	limit := s.getConfigInt("register_ip_limit", 3)
	if limit <= 0 {
		return true, nil
	}

	if global.GVA_REDIS != nil {
		key := fmt.Sprintf("register:ip:%s", ip)
		count, _ := global.GVA_REDIS.Get(context.Background(), key).Int()
		return count < limit, nil
	}

	// DB 回退
	var record client.RegisterIPRecord
	if err := global.GVA_DB.Where("ip = ?", ip).First(&record).Error; err != nil {
		return true, nil // 无记录，允许
	}
	return record.RegisterCount < limit, nil
}

// IncrementRegisterIP 登记IP注册次数
func (s *SecurityService) IncrementRegisterIP(ip string) {
	if global.GVA_REDIS != nil {
		key := fmt.Sprintf("register:ip:%s", ip)
		global.GVA_REDIS.Incr(context.Background(), key)
		return
	}

	// DB 回退
	var record client.RegisterIPRecord
	if err := global.GVA_DB.Where("ip = ?", ip).First(&record).Error; err != nil {
		global.GVA_DB.Create(&client.RegisterIPRecord{IP: ip, RegisterCount: 1})
	} else {
		global.GVA_DB.Model(&record).Update("register_count", record.RegisterCount+1)
	}
}

// ========== 登录失败限制 ==========

func (s *SecurityService) checkLoginIPRateLimitLocal(ip string, limit int64, window time.Duration) (allowed bool, waitSeconds int) {
	if limit <= 0 {
		return true, 0
	}
	if window <= 0 {
		window = time.Duration(defaultLoginIPRateWindow) * time.Second
	}

	now := time.Now()
	localLoginIPRateLimiter.mu.Lock()
	defer localLoginIPRateLimiter.mu.Unlock()

	if len(localLoginIPRateLimiter.counters) > 10000 {
		for key, counter := range localLoginIPRateLimiter.counters {
			if now.After(counter.ExpiresAt) {
				delete(localLoginIPRateLimiter.counters, key)
			}
		}
	}

	counter, ok := localLoginIPRateLimiter.counters[ip]
	if !ok || now.After(counter.ExpiresAt) {
		localLoginIPRateLimiter.counters[ip] = localCaptchaRateCounter{
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
	localLoginIPRateLimiter.counters[ip] = counter
	return true, 0
}

// CheckLoginIPRateLimit 登录接口按IP限流（优先 Redis，失败时回退本地内存窗口计数）
func (s *SecurityService) CheckLoginIPRateLimit(ip string) (allowed bool, waitSeconds int) {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		ip = "unknown"
	}

	limit := utils.GetInt64Setting(loginIPRateLimitConfigKey, loginIPRateLimitEnvKey, defaultLoginIPRateLimit)
	if limit <= 0 {
		return true, 0
	}
	windowSeconds := utils.GetInt64Setting(loginIPRateWindowConfigKey, loginIPRateWindowEnvKey, defaultLoginIPRateWindow)
	if windowSeconds <= 0 {
		windowSeconds = defaultLoginIPRateWindow
	}
	window := time.Duration(windowSeconds) * time.Second

	if global.GVA_REDIS == nil {
		return s.checkLoginIPRateLimitLocal(ip, limit, window)
	}

	key := fmt.Sprintf("login:rate:ip:%s", ip)
	count, err := global.GVA_REDIS.Incr(context.Background(), key).Result()
	if err != nil {
		return s.checkLoginIPRateLimitLocal(ip, limit, window)
	}
	if count == 1 {
		_ = global.GVA_REDIS.Expire(context.Background(), key, window).Err()
	}
	if count > limit {
		ttl, ttlErr := global.GVA_REDIS.TTL(context.Background(), key).Result()
		if ttlErr != nil || ttl <= 0 {
			return false, int(window.Seconds())
		}
		return false, int(ttl.Seconds())
	}

	return true, 0
}

// CheckLoginFail 检查登录失败限制（优先Redis，回退DB）
func (s *SecurityService) CheckLoginFail(username string) (allowed bool, waitSeconds int, err error) {
	maxFail := s.getConfigInt("login_fail_max", 5)
	waitDuration := s.getConfigInt("login_fail_wait_seconds", 900)

	if global.GVA_REDIS != nil {
		key := fmt.Sprintf("login:fail:%s", username)
		count, _ := global.GVA_REDIS.Get(context.Background(), key).Int()
		if count >= maxFail {
			ttl, _ := global.GVA_REDIS.TTL(context.Background(), key).Result()
			return false, int(ttl.Seconds()), nil
		}
		return true, 0, nil
	}

	// DB 回退
	var record client.LoginFailRecord
	if err := global.GVA_DB.Where("username = ?", username).First(&record).Error; err != nil {
		return true, 0, nil // 无记录，允许
	}
	if record.FailCount >= maxFail && time.Now().Before(record.LockedUntil) {
		remaining := int(time.Until(record.LockedUntil).Seconds())
		if remaining < 0 {
			remaining = 0
		}
		return false, remaining, nil
	}
	// 锁已过期，重置
	if record.FailCount >= maxFail && !time.Now().Before(record.LockedUntil) {
		global.GVA_DB.Model(&record).Updates(map[string]interface{}{"fail_count": 0})
	}
	_ = waitDuration // used in RecordLoginFail
	return true, 0, nil
}

// RecordLoginFail 记录登录失败
func (s *SecurityService) RecordLoginFail(username string) {
	maxFail := s.getConfigInt("login_fail_max", 5)
	waitDuration := s.getConfigInt("login_fail_wait_seconds", 900)

	if global.GVA_REDIS != nil {
		key := fmt.Sprintf("login:fail:%s", username)
		global.GVA_REDIS.Incr(context.Background(), key)
		global.GVA_REDIS.Expire(context.Background(), key, time.Duration(waitDuration)*time.Second)
		return
	}

	// DB 回退
	var record client.LoginFailRecord
	if err := global.GVA_DB.Where("username = ?", username).First(&record).Error; err != nil {
		// 第一次失败
		newRecord := client.LoginFailRecord{
			Username:    username,
			FailCount:   1,
			LockedUntil: time.Now().Add(time.Duration(waitDuration) * time.Second),
		}
		if newRecord.FailCount >= maxFail {
			newRecord.LockedUntil = time.Now().Add(time.Duration(waitDuration) * time.Second)
		}
		global.GVA_DB.Create(&newRecord)
		return
	}
	// 如果锁已过期，重新开始计数
	if record.FailCount >= maxFail && !time.Now().Before(record.LockedUntil) {
		global.GVA_DB.Model(&record).Updates(map[string]interface{}{
			"fail_count":   1,
			"locked_until": time.Now().Add(time.Duration(waitDuration) * time.Second),
		})
		return
	}
	newCount := record.FailCount + 1
	updates := map[string]interface{}{"fail_count": newCount}
	if newCount >= maxFail {
		updates["locked_until"] = time.Now().Add(time.Duration(waitDuration) * time.Second)
	}
	global.GVA_DB.Model(&record).Updates(updates)
}

// ClearLoginFail 登录成功后清除失败计数
func (s *SecurityService) ClearLoginFail(username string) {
	if global.GVA_REDIS != nil {
		key := fmt.Sprintf("login:fail:%s", username)
		global.GVA_REDIS.Del(context.Background(), key)
		return
	}

	// DB 回退：直接删除记录
	global.GVA_DB.Where("username = ?", username).Delete(&client.LoginFailRecord{})
}

// ========== 验证码频率限制 ==========

func (s *SecurityService) checkCaptchaRateLimitLocal(ip string, limit int64, window time.Duration) bool {
	if limit <= 0 {
		return true
	}
	if window <= 0 {
		window = time.Duration(defaultCaptchaRateWindow) * time.Second
	}

	now := time.Now()
	localCaptchaRateLimiter.mu.Lock()
	defer localCaptchaRateLimiter.mu.Unlock()

	if len(localCaptchaRateLimiter.counters) > 10000 {
		for key, counter := range localCaptchaRateLimiter.counters {
			if now.After(counter.ExpiresAt) {
				delete(localCaptchaRateLimiter.counters, key)
			}
		}
	}

	counter, ok := localCaptchaRateLimiter.counters[ip]
	if !ok || now.After(counter.ExpiresAt) {
		localCaptchaRateLimiter.counters[ip] = localCaptchaRateCounter{
			Count:     1,
			ExpiresAt: now.Add(window),
		}
		return true
	}

	if counter.Count >= limit {
		return false
	}

	counter.Count++
	localCaptchaRateLimiter.counters[ip] = counter
	return true
}

// CheckCaptchaRateLimit 验证码请求频率限制（优先 Redis，失败时回退本地内存窗口计数）
func (s *SecurityService) CheckCaptchaRateLimit(ip string) (allowed bool) {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		ip = "unknown"
	}
	limit := utils.GetInt64Setting(captchaRateLimitConfigKey, captchaRateLimitEnvKey, defaultCaptchaRateLimit)
	if limit <= 0 {
		return true
	}
	windowSeconds := utils.GetInt64Setting(captchaRateWindowConfigKey, captchaRateWindowEnvKey, defaultCaptchaRateWindow)
	if windowSeconds <= 0 {
		windowSeconds = defaultCaptchaRateWindow
	}
	window := time.Duration(windowSeconds) * time.Second

	if global.GVA_REDIS == nil {
		return s.checkCaptchaRateLimitLocal(ip, limit, window)
	}

	key := fmt.Sprintf("captcha:rate:%s", ip)
	count, err := global.GVA_REDIS.Incr(context.Background(), key).Result()
	if err != nil {
		return s.checkCaptchaRateLimitLocal(ip, limit, window)
	}
	if count == 1 {
		_ = global.GVA_REDIS.Expire(context.Background(), key, window).Err()
	}

	return count <= limit
}
