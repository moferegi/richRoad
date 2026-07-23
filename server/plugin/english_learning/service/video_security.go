package service

import (
	"errors"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/upload"
	"gorm.io/gorm"
)

// SignLearningVideoURL 将学习视频源地址转换为完整可访问URL。
// 1) 有防盗链CDN配置 → 签名防盗链URL
// 2) 无防盗链 → 从ExternalLinkDomain获取默认云域名拼接
func SignLearningVideoURL(rawVideoURL string) string {
	rawVideoURL = strings.TrimSpace(strings.ReplaceAll(rawVideoURL, "\\", "/"))
	if rawVideoURL == "" {
		return ""
	}

	// HLS (m3u8) 视频跳过 CDN 防盗链签名。
	// m3u8 内部的 .ts 分片路径是相对路径，播放器加载时会自动拼接为
	// m3u8 同域名的绝对路径。如果对 m3u8 做了 CDN 签名，.ts 请求
	// 不会携带签名 token，会导致分片加载失败（403）。
	if isHlsVideoURL(rawVideoURL) {
		return resolveWithDefaultDomain(rawVideoURL)
	}

	// 有防盗链CDN：走签名逻辑
	if strings.TrimSpace(global.GVA_CONFIG.Hotlink.CdnDomain) != "" {
		normalizedPath, ok := normalizeLearningVideoSignPath(rawVideoURL)
		if !ok {
			return rawVideoURL
		}
		return upload.SignURL(normalizedPath)
	}

	// 无防盗链：提取路径，拼接默认ExternalLinkDomain
	return resolveWithDefaultDomain(rawVideoURL)
}

var (
	cachedExtDomain   string
	cachedExtDomainAt time.Time
	cachedExtDomainMu sync.RWMutex
)

const extDomainCacheTTL = 30 * time.Second

// getCachedDefaultDomain 获取默认ExternalLinkDomain，缓存30秒自动刷新
func getCachedDefaultDomain() string {
	cachedExtDomainMu.RLock()
	if cachedExtDomain != "" && time.Since(cachedExtDomainAt) < extDomainCacheTTL {
		d := cachedExtDomain
		cachedExtDomainMu.RUnlock()
		return d
	}
	cachedExtDomainMu.RUnlock()

	cachedExtDomainMu.Lock()
	defer cachedExtDomainMu.Unlock()
	// double-check：可能其他goroutine已经刷新了
	if cachedExtDomain != "" && time.Since(cachedExtDomainAt) < extDomainCacheTTL {
		return cachedExtDomain
	}

	if global.GVA_DB == nil {
		return ""
	}

	var domain clientModel.ExternalLinkDomain
	err := global.GVA_DB.Where("is_default = ? AND is_enabled = ?", true, true).First(&domain).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return cachedExtDomain // 返回旧值，不更新为""
		}
		err = global.GVA_DB.Where("is_enabled = ?", true).Order("sort DESC, id ASC").First(&domain).Error
		if err != nil {
			return cachedExtDomain
		}
	}

	result := strings.TrimSpace(domain.Domain)
	if result == "" {
		result = strings.TrimSpace(domain.BaseURL)
	}
	cachedExtDomain = strings.TrimRight(result, "/")
	cachedExtDomainAt = time.Now()
	return cachedExtDomain
}

// resolveWithDefaultDomain 提取rawURL中的路径，拼接默认ExternalLinkDomain
func resolveWithDefaultDomain(rawURL string) string {
	filePath := rawURL
	if strings.HasPrefix(rawURL, "http://") || strings.HasPrefix(rawURL, "https://") {
		parsed, err := url.Parse(rawURL)
		if err == nil {
			filePath = parsed.Path
		}
	}
	if !strings.HasPrefix(filePath, "/") {
		filePath = "/" + filePath
	}

	domain := getCachedDefaultDomain()
	if domain == "" {
		return rawURL
	}
	return domain + filePath
}

func normalizeLearningVideoSignPath(raw string) (string, bool) {
	filePath := raw
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		parsed, err := url.Parse(raw)
		if err != nil {
			return "", false
		}
		filePath = parsed.Path
	}

	cleanPath := path.Clean("/" + strings.TrimPrefix(filePath, "/"))
	if cleanPath == "/" || cleanPath == "." {
		return "", false
	}

	if bucket := getCloudflareR2BucketName(); bucket != "" {
		prefix := "/file/" + bucket + "/"
		if strings.HasPrefix(cleanPath, prefix) {
			cleanPath = "/" + strings.TrimPrefix(cleanPath, prefix)
		}
	}

	if cleanPath == "/" || cleanPath == "." {
		return "", false
	}

	return cleanPath, true
}

// isHlsVideoURL 判断是否为 HLS 流媒体地址（.m3u8）
func isHlsVideoURL(rawURL string) bool {
	return strings.Contains(strings.ToLower(rawURL), ".m3u8")
}

func getCloudflareR2BucketName() string {
	baseURL := strings.TrimSpace(global.GVA_CONFIG.CloudflareR2.BaseURL)
	if baseURL == "" {
		return ""
	}

	if strings.HasPrefix(baseURL, "http://") || strings.HasPrefix(baseURL, "https://") {
		parsed, err := url.Parse(baseURL)
		if err == nil {
			trimmed := strings.Trim(parsed.Path, "/")
			if trimmed != "" {
				parts := strings.Split(trimmed, "/")
				return parts[len(parts)-1]
			}
		}
	}

	trimmed := strings.Trim(baseURL, "/")
	if trimmed == "" {
		return ""
	}
	parts := strings.Split(trimmed, "/")
	return parts[len(parts)-1]
}
