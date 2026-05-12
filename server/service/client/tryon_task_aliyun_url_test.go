package client

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

func TestNormalizeAliyunRetouchImageURL_RewriteCustomDomainToStandardOSS(t *testing.T) {
	originalOssType := global.GVA_CONFIG.System.OssType
	originalEndpoint := global.GVA_CONFIG.AliyunOSS.Endpoint
	originalBucketName := global.GVA_CONFIG.AliyunOSS.BucketName
	originalBucketURL := global.GVA_CONFIG.AliyunOSS.BucketUrl
	originalCDNDomain := global.GVA_CONFIG.Hotlink.CdnDomain
	defer func() {
		global.GVA_CONFIG.System.OssType = originalOssType
		global.GVA_CONFIG.AliyunOSS.Endpoint = originalEndpoint
		global.GVA_CONFIG.AliyunOSS.BucketName = originalBucketName
		global.GVA_CONFIG.AliyunOSS.BucketUrl = originalBucketURL
		global.GVA_CONFIG.Hotlink.CdnDomain = originalCDNDomain
	}()

	global.GVA_CONFIG.System.OssType = "aliyun-oss"
	global.GVA_CONFIG.AliyunOSS.Endpoint = "oss-cn-shanghai.aliyuncs.com"
	global.GVA_CONFIG.AliyunOSS.BucketName = "beautify-bucket"
	global.GVA_CONFIG.AliyunOSS.BucketUrl = "https://cdn.example.com"
	global.GVA_CONFIG.Hotlink.CdnDomain = "https://cdn.example.com"

	s := &TryonTaskService{}
	got, err := s.normalizeAliyunRetouchImageURL("https://cdn.example.com/cloth-on/middle-transfer/demo.jpg?sign=abc&t=123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "https://beautify-bucket.oss-cn-shanghai.aliyuncs.com/cloth-on/middle-transfer/demo.jpg"
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestNormalizeAliyunRetouchImageURL_KeepStandardOSSURL(t *testing.T) {
	originalEndpoint := global.GVA_CONFIG.AliyunOSS.Endpoint
	originalBucketName := global.GVA_CONFIG.AliyunOSS.BucketName
	originalBucketURL := global.GVA_CONFIG.AliyunOSS.BucketUrl
	defer func() {
		global.GVA_CONFIG.AliyunOSS.Endpoint = originalEndpoint
		global.GVA_CONFIG.AliyunOSS.BucketName = originalBucketName
		global.GVA_CONFIG.AliyunOSS.BucketUrl = originalBucketURL
	}()

	global.GVA_CONFIG.AliyunOSS.Endpoint = "oss-cn-shanghai.aliyuncs.com"
	global.GVA_CONFIG.AliyunOSS.BucketName = "beautify-bucket"
	global.GVA_CONFIG.AliyunOSS.BucketUrl = "https://cdn.example.com"

	s := &TryonTaskService{}
	input := "https://beautify-bucket.oss-cn-shanghai.aliyuncs.com/cloth-on/middle-transfer/demo.jpg"
	got, err := s.normalizeAliyunRetouchImageURL(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != input {
		t.Fatalf("expected %s, got %s", input, got)
	}
}

func TestFormatAliyunModelError_CustomDomainHint(t *testing.T) {
	msg := "InvalidImage.URL: 图片链接非法。对于上海oss链接请使用标准的oss域名，目前暂不支持绑定cdn域名和自定义域名"
	got := formatAliyunModelError(msg)
	if got != "modelAliyunRetouchNeedStandardOSS" {
		t.Fatalf("expected modelAliyunRetouchNeedStandardOSS, got %s", got)
	}
}
