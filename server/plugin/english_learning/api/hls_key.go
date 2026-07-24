package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ServeHlsKey 提供HLS AES-128解密密钥
// 播放器加载m3u8后，根据 #EXT-X-KEY 的URI请求此端点获取密钥
// URI 格式: /hlsKey?nonce={32字符hex}
// nonce 切片时随机生成，密钥 = HMAC-SHA256(global-key, nonce)[:16]，永久有效
// @Tags     EnglishContent
// @Summary  获取HLS解密密钥
// @accept   application/json
// @Produce  application/octet-stream
// @Param    nonce query string true "切片时生成的随机nonce（32字符hex）"
// @Success  200  {file} binary "16字节AES-128密钥"
// @Failure  400  {object} string "nonce无效"
// @Failure  500  {object} string "密钥未配置或格式错误"
// @Router   /hlsKey [get]
func ServeHlsKey(c *gin.Context) {
	nonceHex := c.Query("nonce")
	if len(nonceHex) != 32 {
		c.String(http.StatusBadRequest, "invalid nonce")
		return
	}
	// 校验 nonce 是合法 hex
	if _, err := hex.DecodeString(nonceHex); err != nil {
		c.String(http.StatusBadRequest, "invalid nonce")
		return
	}

	// 从配置读取 global-key
	globalKeyHex := global.GVA_CONFIG.Hls.GlobalKey
	if globalKeyHex == "" {
		global.GVA_LOG.Error("HLS global-key 未配置")
		c.String(http.StatusInternalServerError, "hls key not configured")
		return
	}

	globalKeyBytes, err := hex.DecodeString(globalKeyHex)
	if err != nil || len(globalKeyBytes) != 16 {
		global.GVA_LOG.Error("HLS global-key 格式错误", zap.Error(err))
		c.String(http.StatusInternalServerError, "hls key invalid")
		return
	}

	// 每集独立密钥：HMAC-SHA256(global-key, nonce)，取前16字节
	mac := hmac.New(sha256.New, globalKeyBytes)
	mac.Write([]byte(nonceHex))
	perVideoKey := mac.Sum(nil)[:16]

	// 返回16字节原始密钥
	c.Data(http.StatusOK, "application/octet-stream", perVideoKey)
}
