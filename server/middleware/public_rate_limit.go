package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
	"github.com/gin-gonic/gin"
)

const (
	publicRateLimitEnabledConfigKey = "security_public_rate_limit_enabled"
	publicRateLimitEnabledEnvKey    = "CS_PUBLIC_RATE_LIMIT_ENABLED"

	publicRateLimitWindowConfigKey = "security_public_rate_limit_window_seconds"
	publicRateLimitWindowEnvKey    = "CS_PUBLIC_RATE_LIMIT_WINDOW_SECONDS"

	publicRateLimitMaxConfigKey = "security_public_rate_limit_max_requests"
	publicRateLimitMaxEnvKey    = "CS_PUBLIC_RATE_LIMIT_MAX_REQUESTS"

	defaultPublicRateLimitWindowSeconds int64 = 60
	defaultPublicRateLimitMaxRequests   int64 = 300
)

// PublicRateLimit 对公开接口做按 IP 的统一限流。
func PublicRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !utils.GetBoolSetting(publicRateLimitEnabledConfigKey, publicRateLimitEnabledEnvKey, true) {
			c.Next()
			return
		}

		if global.GVA_REDIS == nil {
			c.Next()
			return
		}

		windowSeconds := utils.GetInt64Setting(publicRateLimitWindowConfigKey, publicRateLimitWindowEnvKey, defaultPublicRateLimitWindowSeconds)
		if windowSeconds < 1 {
			windowSeconds = defaultPublicRateLimitWindowSeconds
		}
		maxRequests := utils.GetInt64Setting(publicRateLimitMaxConfigKey, publicRateLimitMaxEnvKey, defaultPublicRateLimitMaxRequests)
		if maxRequests < 1 {
			c.Next()
			return
		}

		ctx := context.Background()
		key := fmt.Sprintf("rate:public:ip:%s", c.ClientIP())
		count, err := global.GVA_REDIS.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}
		if count == 1 {
			_ = global.GVA_REDIS.Expire(ctx, key, time.Duration(windowSeconds)*time.Second).Err()
		}

		if count > maxRequests {
			msg := i18n.T(c, "requestTooFrequent")
			if ttl, ttlErr := global.GVA_REDIS.TTL(ctx, key).Result(); ttlErr == nil && ttl > 0 {
				msg = fmt.Sprintf("%s (%ds)", msg, int(ttl.Seconds()))
			}
			c.JSON(http.StatusOK, gin.H{"code": response.ERROR, "msg": msg})
			c.Abort()
			return
		}

		c.Next()
	}
}
