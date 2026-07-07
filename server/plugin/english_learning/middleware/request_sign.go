package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

const (
	signHeaderTimestamp = "X-Req-Ts"
	signHeaderNonce     = "X-Req-Nonce"
	signHeaderSignature = "X-Req-Sign"
	signSalt            = "uni-api-sign-v1"
	nonceTTLSeconds     = 300
)

func LearningRequestSignGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !shouldVerifySign(c) {
			c.Next()
			return
		}

		tsStr := strings.TrimSpace(c.GetHeader(signHeaderTimestamp))
		nonce := strings.TrimSpace(c.GetHeader(signHeaderNonce))
		signature := strings.TrimSpace(c.GetHeader(signHeaderSignature))
		if tsStr == "" || nonce == "" || signature == "" {
			response.FailWithMessage("签名参数缺失", c)
			c.Abort()
			return
		}

		ts, err := strconv.ParseInt(tsStr, 10, 64)
		if err != nil {
			response.FailWithMessage("签名时间戳无效", c)
			c.Abort()
			return
		}
		now := time.Now().Unix()
		if now-ts > nonceTTLSeconds || ts-now > 30 {
			response.FailWithMessage("签名已过期", c)
			c.Abort()
			return
		}

		token := strings.TrimSpace(c.GetHeader("x-token"))
		expected := buildReqSignature(c.Request.Method, c.FullPath(), c.Request.URL.RawQuery, tsStr, nonce, token)
		if !hmac.Equal([]byte(strings.ToLower(signature)), []byte(strings.ToLower(expected))) {
			response.FailWithMessage("签名校验失败", c)
			c.Abort()
			return
		}

		if isReplayNonce(nonce, token) {
			response.FailWithMessage("请求重复提交", c)
			c.Abort()
			return
		}

		c.Next()
	}
}

func shouldVerifySign(c *gin.Context) bool {
	if c == nil || c.Request == nil {
		return false
	}
	if !utils.GetBoolSetting("learning_api_sign_enabled", "LEARNING_API_SIGN_ENABLED", true) {
		return false
	}
	platform := strings.ToLower(strings.TrimSpace(c.GetHeader("X-Client-Platform")))
	if platform == "uni" || platform == "uniapp" || platform == "uni-app" {
		return true
	}
	return false
}

func buildReqSignature(method, path, rawQuery, ts, nonce, token string) string {
	material := strings.Join([]string{
		strings.ToUpper(strings.TrimSpace(method)),
		strings.TrimSpace(path),
		strings.TrimSpace(rawQuery),
		strings.TrimSpace(ts),
		strings.TrimSpace(nonce),
	}, "|")
	key := sha256.Sum256([]byte(strings.TrimSpace(token) + "|" + signSalt))
	h := hmac.New(sha256.New, key[:])
	_, _ = h.Write([]byte(material))
	return hex.EncodeToString(h.Sum(nil))
}

func isReplayNonce(nonce string, token string) bool {
	if global.GVA_REDIS == nil {
		return false
	}
	tokenHash := sha256.Sum256([]byte(strings.TrimSpace(token)))
	key := fmt.Sprintf("EL_SIGN_NONCE:%s:%s", hex.EncodeToString(tokenHash[:8]), nonce)
	ok, err := global.GVA_REDIS.SetNX(context.Background(), key, "1", time.Duration(nonceTTLSeconds)*time.Second).Result()
	if err != nil {
		return false
	}
	return !ok
}
