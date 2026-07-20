package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

const (
	uniSignHeaderTs    = "X-Req-Ts"
	uniSignHeaderNonce = "X-Req-Nonce"
	uniSignHeaderSign  = "X-Req-Sign"
	uniSignSalt        = "uni-api-sign-v1"
	uniSignMaxAgeSecs  = 300 // 时间窗口 ±5 分钟
)

// nonceCache 内存级 nonce 去重，避免重放攻击
var (
	nonceCache   = make(map[string]int64)
	nonceCacheMu sync.RWMutex
)

func init() {
	// 每 60 秒清理过期 nonce
	go func() {
		for {
			time.Sleep(60 * time.Second)
			cleanExpiredNonces()
		}
	}()
}

func cleanExpiredNonces() {
	nonceCacheMu.Lock()
	defer nonceCacheMu.Unlock()
	now := time.Now().Unix()
	for k, ts := range nonceCache {
		if now-ts > uniSignMaxAgeSecs*2 {
			delete(nonceCache, k)
		}
	}
}

// UniSignVerify Uni 端请求签名校验中间件
// 通过 sysConfig 中 learning_api_sign_enabled 开关控制是否校验
func UniSignVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isUniSignEnabled() {
			c.Next()
			return
		}
		// 仅校验 X-Client-Platform=uni 的请求
		platform := strings.ToLower(strings.TrimSpace(c.GetHeader("X-Client-Platform")))
		if platform != "uni" && platform != "uniapp" && platform != "uni-app" {
			c.Next()
			return
		}
		// 公开接口（无需 token 的登录/注册等）跳过签名校验
		if isUniPublicPath(c.Request.URL.Path) {
			c.Next()
			return
		}
		if !verifyUniSign(c) {
			c.Abort()
			return
		}
		c.Next()
	}
}

func isUniSignEnabled() bool {
	return utils.GetBoolSetting("learning_api_sign_enabled", "LEARNING_API_SIGN_ENABLED", true)
}

// isUniPublicPath 公开路径无需签名校验
func isUniPublicPath(path string) bool {
	publicPaths := []string{
		"/api/base/captcha",
		"/api/clientUser/login",
		"/api/clientUser/register",
		"/api/clientUser/phoneLogin",
		"/api/clientUser/phoneRegister",
		"/api/sysConfig/getLoginConfig",
		"/api/sysConfig/getSysConfigByKey",
		"/api/phoneAreaCode/getEnabledPhoneAreaCodes",
		"/api/banner/getBannerList",
	}
	lower := strings.ToLower(strings.TrimSpace(path))
	for _, p := range publicPaths {
		if strings.HasPrefix(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

func verifyUniSign(c *gin.Context) bool {
	tsStr := strings.TrimSpace(c.GetHeader(uniSignHeaderTs))
	nonce := strings.TrimSpace(c.GetHeader(uniSignHeaderNonce))
	sign := strings.TrimSpace(c.GetHeader(uniSignHeaderSign))

	// 三个签名头缺一不可
	if tsStr == "" || nonce == "" || sign == "" {
		_ = c.AbortWithError(401, fmt.Errorf("missing sign headers"))
		return false
	}

	// 时间戳校验：防重放，不允许超过 ±5 分钟
	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		_ = c.AbortWithError(401, fmt.Errorf("invalid ts"))
		return false
	}
	now := time.Now().Unix()
	if absDiff(now, ts) > uniSignMaxAgeSecs {
		_ = c.AbortWithError(401, fmt.Errorf("ts expired"))
		return false
	}

	// nonce 去重：同一个 nonce 只能使用一次
	if exists := checkAndStoreNonce(nonce, now); exists {
		_ = c.AbortWithError(401, fmt.Errorf("nonce replayed"))
		return false
	}

	// 计算期望签名并比较
	token := utils.GetToken(c)
	if token == "" {
		return false
	}
	rawPath := strings.TrimSpace(c.Request.URL.Path)
	rawQuery := c.Request.URL.RawQuery
	method := strings.ToUpper(strings.TrimSpace(c.Request.Method))

	expectedSign := buildUniSign(method, rawPath, rawQuery, tsStr, nonce, token)
	if !hmac.Equal([]byte(sign), []byte(expectedSign)) {
		_ = c.AbortWithError(401, fmt.Errorf("sign mismatch"))
		return false
	}

	return true
}

func buildUniSign(method, path, query, ts, nonce, token string) string {
	material := fmt.Sprintf("%s|%s|%s|%s|%s",
		strings.ToUpper(strings.TrimSpace(method)),
		strings.TrimSpace(path),
		strings.TrimSpace(query),
		strings.TrimSpace(ts),
		strings.TrimSpace(nonce),
	)
	signKey := deriveUniSignKey(token)
	mac := hmac.New(sha256.New, signKey)
	mac.Write([]byte(material))
	return hex.EncodeToString(mac.Sum(nil))
}

func deriveUniSignKey(token string) []byte {
	sum := sha256.Sum256([]byte(strings.TrimSpace(token) + "|" + uniSignSalt))
	return sum[:]
}

func checkAndStoreNonce(nonce string, now int64) bool {
	nonceCacheMu.Lock()
	defer nonceCacheMu.Unlock()
	if _, exists := nonceCache[nonce]; exists {
		return true
	}
	nonceCache[nonce] = now
	return false
}

func absDiff(a, b int64) int64 {
	if a > b {
		return a - b
	}
	return b - a
}
