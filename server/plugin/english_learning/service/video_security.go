package service

import (
	"net/url"
	"path"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/upload"
)

// SignLearningVideoURL 将学习视频源地址转换为防盗链签名地址。
// 当未配置 CDN 域名时，保持原始地址，避免影响本地开发或外链回放。
func SignLearningVideoURL(rawVideoURL string) string {
	rawVideoURL = strings.TrimSpace(strings.ReplaceAll(rawVideoURL, "\\", "/"))
	if rawVideoURL == "" {
		return ""
	}

	if strings.TrimSpace(global.GVA_CONFIG.Hotlink.CdnDomain) == "" {
		return rawVideoURL
	}

	normalizedPath, ok := normalizeLearningVideoSignPath(rawVideoURL)
	if !ok {
		return rawVideoURL
	}

	return upload.SignURL(normalizedPath)
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
