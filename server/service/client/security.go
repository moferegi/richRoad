package client

import (
	"context"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
)

type SecurityService struct{}

// CheckRegisterIPLimit 检查IP注册限制
func (s *SecurityService) CheckRegisterIPLimit(ip string) (allowed bool, err error) {
	var config client.SysConfig
	err = global.GVA_DB.Where("config_key = ?", "register_ip_limit").First(&config).Error
	if err != nil {
		return true, nil // 配置不存在则不限制
	}

	limit := 3
	fmt.Sscanf(config.ConfigValue, "%d", &limit)
	if limit <= 0 {
		return true, nil
	}

	// 使用Redis计数
	if global.GVA_REDIS != nil {
		key := fmt.Sprintf("register:ip:%s", ip)
		count, _ := global.GVA_REDIS.Get(context.Background(), key).Int()
		return count < limit, nil
	}

	return true, nil
}

// IncrementRegisterIP 登记IP注册次数
func (s *SecurityService) IncrementRegisterIP(ip string) {
	if global.GVA_REDIS != nil {
		key := fmt.Sprintf("register:ip:%s", ip)
		global.GVA_REDIS.Incr(context.Background(), key)
		// 永不过期，如需按天重置可加TTL
	}
}

// CheckLoginFail 检查登录失败限制
func (s *SecurityService) CheckLoginFail(username string) (allowed bool, waitSeconds int, err error) {
	if global.GVA_REDIS == nil {
		return true, 0, nil
	}

	var maxFail int = 5
	var waitDuration int = 900 // 秒

	var config client.SysConfig
	if e := global.GVA_DB.Where("config_key = ?", "login_fail_max").First(&config).Error; e == nil {
		fmt.Sscanf(config.ConfigValue, "%d", &maxFail)
	}
	if e := global.GVA_DB.Where("config_key = ?", "login_fail_wait_seconds").First(&config).Error; e == nil {
		fmt.Sscanf(config.ConfigValue, "%d", &waitDuration)
	}

	key := fmt.Sprintf("login:fail:%s", username)
	count, _ := global.GVA_REDIS.Get(context.Background(), key).Int()

	if count >= maxFail {
		ttl, _ := global.GVA_REDIS.TTL(context.Background(), key).Result()
		return false, int(ttl.Seconds()), nil
	}
	return true, 0, nil
}

// RecordLoginFail 记录登录失败
func (s *SecurityService) RecordLoginFail(username string) {
	if global.GVA_REDIS == nil {
		return
	}

	var waitDuration int = 900
	var config client.SysConfig
	if e := global.GVA_DB.Where("config_key = ?", "login_fail_wait_seconds").First(&config).Error; e == nil {
		fmt.Sscanf(config.ConfigValue, "%d", &waitDuration)
	}

	key := fmt.Sprintf("login:fail:%s", username)
	global.GVA_REDIS.Incr(context.Background(), key)
	global.GVA_REDIS.Expire(context.Background(), key, time.Duration(waitDuration)*time.Second)
}

// ClearLoginFail 登录成功后清除失败计数
func (s *SecurityService) ClearLoginFail(username string) {
	if global.GVA_REDIS == nil {
		return
	}
	key := fmt.Sprintf("login:fail:%s", username)
	global.GVA_REDIS.Del(context.Background(), key)
}

// CheckCaptchaRateLimit 验证码请求频率限制
func (s *SecurityService) CheckCaptchaRateLimit(ip string) (allowed bool) {
	if global.GVA_REDIS == nil {
		return true
	}

	// 从 sysConfig 读取频率限制
	limit := 10
	var cfg client.SysConfig
	if err := global.GVA_DB.Where("config_key = ?", "captcha_rate_limit").First(&cfg).Error; err == nil {
		if v, e := fmt.Sscanf(cfg.ConfigValue, "%d", &limit); v == 0 || e != nil {
			limit = 10
		}
	}

	key := fmt.Sprintf("captcha:rate:%s", ip)
	count, _ := global.GVA_REDIS.Get(context.Background(), key).Int()
	if count >= limit {
		return false
	}
	global.GVA_REDIS.Incr(context.Background(), key)
	global.GVA_REDIS.Expire(context.Background(), key, 60*time.Second)
	return true
}
