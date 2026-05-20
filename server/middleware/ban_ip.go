package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

const bannedIPRedisHash = "sys:banned_ips"

// BanIPCheck 检查请求IP是否被封禁的中间件
func BanIPCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		// 优先查 Redis（O(1) 性能）
		if global.GVA_REDIS != nil {
			val, err := global.GVA_REDIS.HGet(context.Background(), bannedIPRedisHash, ip).Result()
			if err == nil {
				if val == "0" {
					// 永久封禁
					c.JSON(http.StatusForbidden, gin.H{
						"code": response.ERROR,
						"msg":  fmt.Sprintf("您的IP（%s）已被封禁，请联系管理员", ip),
					})
					c.Abort()
					return
				}
				expiry, parseErr := strconv.ParseInt(val, 10, 64)
				if parseErr == nil && time.Now().Unix() < expiry {
					// 临时封禁未到期
					remaining := time.Until(time.Unix(expiry, 0)).Minutes()
					c.JSON(http.StatusForbidden, gin.H{
						"code": response.ERROR,
						"msg":  fmt.Sprintf("您的IP（%s）已被临时封禁，剩余 %.0f 分钟", ip, remaining),
					})
					c.Abort()
					return
				}
				// 过期封禁，从 Redis 移除（DB 清理由 BannedIPService 定期处理）
				if parseErr == nil {
					global.GVA_REDIS.HDel(context.Background(), bannedIPRedisHash, ip)
				}
			}
			c.Next()
			return
		}

		// Redis 不可用时降级查 DB
		if global.GVA_DB != nil {
			var count int64
			now := time.Now()
			global.GVA_DB.Table("sys_banned_ips").
				Where("ip = ? AND deleted_at IS NULL AND (expired_at IS NULL OR expired_at > ?)", ip, now).
				Count(&count)
			if count > 0 {
				c.JSON(http.StatusForbidden, gin.H{
					"code": response.ERROR,
					"msg":  fmt.Sprintf("您的IP（%s）已被封禁，请联系管理员", ip),
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
