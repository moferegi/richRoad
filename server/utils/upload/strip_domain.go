package upload

import (
	"net/url"
	"strings"
)

// knownDomains 已知云存储域名前缀列表（入库前自动剥离）
var knownDomains = []string{
	// 七牛云
	"qiniucdn.com",
	"qiniu.com",
	"clouddn.com",
	// Cloudflare R2
	"r2.cloudflarestorage.com",
	"r2.dev",
	// Backblaze B2
	"backblazeb2.com",
	// 阿里云 OSS
	"aliyuncs.com",
	"oss-cn-",
	// 腾讯云 COS
	"myqcloud.com",
	// AWS S3
	"amazonaws.com",
	// MinIO (自定义域名通常不剥离，但以防万一)
	// 本地存储
	"localhost",
	"127.0.0.1",
}

// StripDomain 从完整 URL 中剥离域名，返回相对路径
// 输入："https://cdn.example.com/english-learn/pic/file.jpg"
// 输出："english-learn/pic/file.jpg"
// 若 URL 已经是相对路径或解析失败，原样返回
func StripDomain(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}
	// 已经是相对路径
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		if strings.HasPrefix(rawURL, "//") {
			rawURL = "https:" + rawURL
		} else {
			return rawURL
		}
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	path := strings.TrimPrefix(parsed.Path, "/")
	if path == "" {
		return ""
	}
	return path
}

// IsKnownStorageDomain 判断 URL 是否属于已知云存储域名
func IsKnownStorageDomain(rawURL string) bool {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return false
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	host := strings.ToLower(parsed.Host)
	for _, domain := range knownDomains {
		if strings.Contains(host, domain) {
			return true
		}
	}
	return false
}
