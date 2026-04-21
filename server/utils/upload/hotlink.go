package upload

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// SignURL 对文件路径生成防盗链签名URL
// 如果防盗链未启用，返回 cdnDomain + path 的明文URL
// 如果防盗链启用，返回带 sign 和 t 参数的签名URL（HMAC-SHA256）
// filePath 可以是路径（/xxx.mp4）或完整URL（https://cdn.example.com/xxx.mp4）
func SignURL(filePath string) string {
	cfg := global.GVA_CONFIG.Hotlink

	// 如果传入的是完整URL，提取其路径部分
	if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
		if parsed, err := url.Parse(filePath); err == nil {
			filePath = parsed.Path
		}
	}

	if cfg.CdnDomain == "" {
		// 没有配置CDN域名，使用当前OSS的baseURL
		return getOssBaseURL() + "/" + strings.TrimPrefix(filePath, "/")
	}

	domain := strings.TrimRight(cfg.CdnDomain, "/")
	path := "/" + strings.TrimPrefix(filePath, "/")

	if !cfg.Enabled || cfg.SignKey == "" {
		// 防盗链关闭，返回明文URL
		return domain + path
	}

	// 生成签名URL
	expire := cfg.ExpireSeconds
	if expire <= 0 {
		expire = 3600 // 默认1小时
	}
	deadline := time.Now().Unix() + expire
	deadlineHex := strconv.FormatInt(deadline, 16)

	// 签名算法: HMAC-SHA256(key, path + ":" + deadlineHex)
	mac := hmac.New(sha256.New, []byte(cfg.SignKey))
	mac.Write([]byte(path + ":" + deadlineHex))
	sign := hex.EncodeToString(mac.Sum(nil))

	return fmt.Sprintf("%s%s?sign=%s&t=%s", domain, path, sign, deadlineHex)
}

// VerifySignURL 验证签名URL是否有效
func VerifySignURL(rawURL string) bool {
	cfg := global.GVA_CONFIG.Hotlink
	if !cfg.Enabled || cfg.SignKey == "" {
		return true
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	sign := u.Query().Get("sign")
	t := u.Query().Get("t")
	if sign == "" || t == "" {
		return false
	}

	// 验证过期时间
	deadline, err := strconv.ParseInt(t, 16, 64)
	if err != nil {
		return false
	}
	if time.Now().Unix() > deadline {
		return false
	}

	// 验证签名: HMAC-SHA256(key, path + ":" + t)
	path := u.Path
	mac := hmac.New(sha256.New, []byte(cfg.SignKey))
	mac.Write([]byte(path + ":" + t))
	expectedSign := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(sign), []byte(expectedSign))
}

// getOssBaseURL 获取当前OSS的baseURL
func getOssBaseURL() string {
	switch global.GVA_CONFIG.System.OssType {
	case "qiniu":
		return global.GVA_CONFIG.Qiniu.ImgPath
	case "aws-s3":
		return global.GVA_CONFIG.AwsS3.BaseURL
	case "cloudflare-r2":
		return global.GVA_CONFIG.CloudflareR2.BaseURL
	case "aliyun-oss":
		return global.GVA_CONFIG.AliyunOSS.BucketUrl
	case "tencent-cos":
		return global.GVA_CONFIG.TencentCOS.BaseURL
	case "minio":
		return global.GVA_CONFIG.Minio.BucketUrl
	default:
		return ""
	}
}
