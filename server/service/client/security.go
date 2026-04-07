package client

import (
	"context"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
)

type SecurityService struct{}

// ---------- 辅助：读取 sysConfig ----------

func (s *SecurityService) getConfigInt(key string, def int) int {
	var cfg client.SysConfig
	if err := global.GVA_DB.Where("config_key = ?", key).First(&cfg).Error; err == nil {
		v := def
		fmt.Sscanf(cfg.ConfigValue, "%d", &v)
		return v
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

// CheckCaptchaRateLimit 验证码请求频率限制（仅Redis支持，DB不做频控）
func (s *SecurityService) CheckCaptchaRateLimit(ip string) (allowed bool) {
	if global.GVA_REDIS == nil {
		return true // 验证码频率限制仅Redis支持
	}

	limit := s.getConfigInt("captcha_rate_limit", 10)

	key := fmt.Sprintf("captcha:rate:%s", ip)
	count, _ := global.GVA_REDIS.Get(context.Background(), key).Int()
	if count >= limit {
		return false
	}
	global.GVA_REDIS.Incr(context.Background(), key)
	global.GVA_REDIS.Expire(context.Background(), key, 60*time.Second)
	return true
}
