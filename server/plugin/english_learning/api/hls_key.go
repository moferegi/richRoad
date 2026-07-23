package api

import (
	"encoding/hex"
	"net/http"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ServeHlsKey 提供HLS AES-128解密密钥
// 播放器加载m3u8后，根据 #EXT-X-KEY 的URI请求此端点获取密钥
// URI 格式: /api/hlsKey?token={JWT}
// JWT 由切片时生成，包含 episodeId + purpose="hls_key" + 30天过期
// @Tags     EnglishContent
// @Summary  获取HLS解密密钥
// @accept   application/json
// @Produce  application/octet-stream
// @Param    token query string true "HLS密钥访问token"
// @Success  200  {file} binary "16字节AES-128密钥"
// @Failure  400  {object} string "token无效或密钥未配置"
// @Router   /api/hlsKey [get]
func ServeHlsKey(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.String(http.StatusBadRequest, "missing token")
		return
	}

	// 验证token
	claims, err := videoSliceService.ValidateHlsKeyToken(tokenStr)
	if err != nil {
		global.GVA_LOG.Warn("HLS密钥token验证失败", zap.Error(err))
		c.String(http.StatusForbidden, "invalid token")
		return
	}

	global.GVA_LOG.Debug("HLS密钥请求通过", zap.Uint("episodeId", claims.EpisodeID))

	// 从配置读取 global-key
	globalKeyHex := global.GVA_CONFIG.Hls.GlobalKey
	if globalKeyHex == "" {
		global.GVA_LOG.Error("HLS global-key 未配置")
		c.String(http.StatusInternalServerError, "hls key not configured")
		return
	}

	keyBytes, err := hex.DecodeString(globalKeyHex)
	if err != nil || len(keyBytes) != 16 {
		global.GVA_LOG.Error("HLS global-key 格式错误", zap.Error(err))
		c.String(http.StatusInternalServerError, "hls key invalid")
		return
	}

	// 返回16字节原始密钥
	c.Data(http.StatusOK, "application/octet-stream", keyBytes)
}
