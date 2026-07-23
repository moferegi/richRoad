package client

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientReq "github.com/flipped-aurora/gin-vue-admin/server/model/client/request"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	qiniuStorage "github.com/qiniu/go-sdk/v7/storage"
	"go.uber.org/zap"
)

// CloudStorageService 云存储操作服务
type CloudStorageService struct{}

// getCloudConfigs 获取所有已配置云存储的域名
func (s *CloudStorageService) getCloudConfigs() ([]client.ExternalLinkDomain, error) {
	var domains []client.ExternalLinkDomain
	err := global.GVA_DB.Where("is_enabled = ? AND cloud_type != ''", true).
		Order("sort DESC, id ASC").Find(&domains).Error
	return domains, err
}

// getCloudConfigByID 获取指定ID的云配置
func (s *CloudStorageService) getCloudConfigByID(id uint) (client.ExternalLinkDomain, error) {
	var domain client.ExternalLinkDomain
	err := global.GVA_DB.Where("id = ?", id).First(&domain).Error
	return domain, err
}

// buildS3Client 从域名配置构建S3客户端
func (s *CloudStorageService) buildS3Client(domain client.ExternalLinkDomain) (*s3.S3, error) {
	cloudType := strings.TrimSpace(domain.CloudType)
	endpoint := strings.TrimSpace(domain.RegionEndpoint)
	accountId := strings.TrimSpace(domain.AccountId)

	// R2：如果没有提供完整 Endpoint，用 AccountId 自动构造
	if cloudType == "r2" && endpoint == "" && accountId != "" {
		endpoint = fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId)
	}

	if endpoint == "" {
		return nil, fmt.Errorf("云存储区域/Endpoint未配置")
	}

	accessKey := strings.TrimSpace(domain.AccessKey)
	secretKey := strings.TrimSpace(domain.SecretKey)
	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("云存储密钥未配置")
	}

	cfg := &aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Region:      aws.String("auto"),
		DisableSSL:  aws.Bool(false),
	}

	// 区分是完整 endpoint 还是需要拼接前缀
	endpointLower := strings.ToLower(endpoint)
	if strings.HasPrefix(endpointLower, "http://") || strings.HasPrefix(endpointLower, "https://") {
		cfg.Endpoint = aws.String(endpoint)
	} else {
		// 纯域名格式：仅加 https:// 前缀，不拼接 bucket
		// S3ForcePathStyle 已启用，bucket 通过请求路径传递，无需放在 hostname 中
		cfg.Endpoint = aws.String("https://" + endpoint)
	}

	// 启用路径寻址模式 (兼容大多数S3兼容存储)
	cfg.S3ForcePathStyle = aws.Bool(true)

	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, fmt.Errorf("创建S3会话失败: %w", err)
	}

	return s3.New(sess), nil
}

// pingQiniu 检测七牛云连接
func (s *CloudStorageService) pingQiniu(domain client.ExternalLinkDomain) (*clientReq.CloudPingResp, error) {
	accessKey := strings.TrimSpace(domain.AccessKey)
	secretKey := strings.TrimSpace(domain.SecretKey)
	bucket := strings.TrimSpace(domain.Bucket)

	if accessKey == "" || secretKey == "" {
		return &clientReq.CloudPingResp{Success: false, Message: "七牛云密钥未配置"}, nil
	}
	if bucket == "" {
		return &clientReq.CloudPingResp{Success: false, Message: "七牛云Bucket未配置"}, nil
	}

	mac := qbox.NewMac(accessKey, secretKey)
	cfg := qiniuConfigForDomain(domain)
	bucketManager := qiniuStorage.NewBucketManager(mac, cfg)

	start := time.Now()
	_, _, _, _, err := bucketManager.ListFiles(bucket, "", "", "", 1)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return &clientReq.CloudPingResp{
			Success: false,
			Latency: latency,
			Message: fmt.Sprintf("连接失败: %v", err),
		}, nil
	}

	return &clientReq.CloudPingResp{
		Success: true,
		Latency: latency,
		Message: "连接成功",
	}, nil
}

// listQiniuFiles 列出七牛云文件
func (s *CloudStorageService) listQiniuFiles(domain client.ExternalLinkDomain, req clientReq.CloudListFilesReq) (*clientReq.CloudListFilesResp, error) {
	accessKey := strings.TrimSpace(domain.AccessKey)
	secretKey := strings.TrimSpace(domain.SecretKey)
	bucket := strings.TrimSpace(domain.Bucket)

	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("七牛云密钥未配置")
	}

	mac := qbox.NewMac(accessKey, secretKey)
	cfg := qiniuConfigForDomain(domain)
	bucketManager := qiniuStorage.NewBucketManager(mac, cfg)

	maxKeys := req.MaxKeys
	if maxKeys <= 0 || maxKeys > 1000 {
		maxKeys = 200
	}

	prefix := strings.TrimSpace(req.Prefix)
	marker := strings.TrimSpace(req.Marker)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 七牛 ListFiles 返回 (entries, commonPrefixes, nextMarker, hasNext, err)
	// 但旧版 API 可能不同，直接用 goroutine + channel 包装
	type listResult struct {
		entries        []qiniuStorage.ListItem
		commonPrefixes []string
		nextMarker     string
		hasNext        bool
		err            error
	}

	ch := make(chan listResult, 1)
	go func() {
		entries, commonPrefixes, nextMarker, hasNext, err := bucketManager.ListFiles(bucket, prefix, "/", marker, maxKeys)
		ch <- listResult{entries, commonPrefixes, nextMarker, hasNext, err}
	}()

	select {
	case result := <-ch:
		if result.err != nil {
			return nil, fmt.Errorf("列出文件失败: %w", result.err)
		}

		resp := &clientReq.CloudListFilesResp{
			Prefix:      prefix,
			IsTruncated: result.hasNext,
			Marker:      result.nextMarker,
		}

		for _, cp := range result.commonPrefixes {
			count, size := s.getDirStatsQiniu(bucketManager, bucket, cp)
			resp.Files = append(resp.Files, clientReq.CloudFileItem{
				Key:          cp,
				Size:         0,
				IsDir:        true,
				DirFileCount: count,
				DirSize:      size,
			})
		}

		for _, entry := range result.entries {
			if prefix != "" && entry.Key == prefix {
				continue
			}
			lastMod := time.Unix(entry.PutTime/10000000, 0).Format("2006-01-02 15:04:05")
			resp.Files = append(resp.Files, clientReq.CloudFileItem{
				Key:          entry.Key,
				Size:         entry.Fsize,
				LastModified: lastMod,
				IsDir:        false,
			})
		}

		return resp, nil

	case <-ctx.Done():
		return nil, fmt.Errorf("列出文件超时")
	}
}

// qiniuConfigForDomain 基于域名配置构建七牛存储配置
func qiniuConfigForDomain(domain client.ExternalLinkDomain) *qiniuStorage.Config {
	cfg := &qiniuStorage.Config{}
	if domain.UseHTTPS != nil {
		cfg.UseHTTPS = *domain.UseHTTPS
	} else {
		cfg.UseHTTPS = true
	}
	if domain.UseCdnDomain != nil {
		cfg.UseCdnDomains = *domain.UseCdnDomain
	}

	switch domain.Zone {
	case "ZoneHuadong":
		cfg.Zone = &qiniuStorage.ZoneHuadong
	case "ZoneHuabei":
		cfg.Zone = &qiniuStorage.ZoneHuabei
	case "ZoneHuanan":
		cfg.Zone = &qiniuStorage.ZoneHuanan
	case "ZoneBeimei":
		cfg.Zone = &qiniuStorage.ZoneBeimei
	case "ZoneXinjiapo":
		cfg.Zone = &qiniuStorage.ZoneXinjiapo
	}
	return cfg
}

// PingCloud 检测云存储连接
func (s *CloudStorageService) PingCloud(req clientReq.CloudPingReq) (*clientReq.CloudPingResp, error) {
	domain, err := s.getCloudConfigByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("域名不存在: %w", err)
	}

	cloudType := strings.TrimSpace(domain.CloudType)
	if cloudType == "" {
		return &clientReq.CloudPingResp{Success: false, Message: "该域名未配置云存储类型"}, nil
	}

	// 七牛云使用独立 SDK
	if cloudType == "qiniu" {
		return s.pingQiniu(domain)
	}

	s3Client, err := s.buildS3Client(domain)
	if err != nil {
		return &clientReq.CloudPingResp{Success: false, Message: err.Error()}, nil
	}

	bucket := strings.TrimSpace(domain.Bucket)
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = s3Client.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(1),
	})

	latency := time.Since(start).Milliseconds()

	if err != nil {
		return &clientReq.CloudPingResp{
			Success: false,
			Latency: latency,
			Message: fmt.Sprintf("连接失败: %v", err),
		}, nil
	}

	return &clientReq.CloudPingResp{
		Success: true,
		Latency: latency,
		Message: "连接成功",
	}, nil
}

// ListCloudFiles 列出云存储文件
func (s *CloudStorageService) ListCloudFiles(req clientReq.CloudListFilesReq) (*clientReq.CloudListFilesResp, error) {
	domain, err := s.getCloudConfigByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("域名不存在: %w", err)
	}

	cloudType := strings.TrimSpace(domain.CloudType)

	// 七牛云使用独立 SDK
	if cloudType == "qiniu" {
		return s.listQiniuFiles(domain, req)
	}

	s3Client, err := s.buildS3Client(domain)
	if err != nil {
		return nil, err
	}

	bucket := strings.TrimSpace(domain.Bucket)
	maxKeys := int64(req.MaxKeys)
	if maxKeys <= 0 || maxKeys > 1000 {
		maxKeys = 100
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		MaxKeys:   aws.Int64(maxKeys),
		Delimiter: aws.String("/"), // 始终启用目录层级分隔，返回 CommonPrefixes
	}

	prefix := strings.TrimSpace(req.Prefix)
	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	marker := strings.TrimSpace(req.Marker)
	if marker != "" {
		input.StartAfter = aws.String(marker)
	}

	result, err := s3Client.ListObjectsV2WithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("列出文件失败: %w", err)
	}

	resp := &clientReq.CloudListFilesResp{
		Prefix:      prefix,
		IsTruncated: aws.BoolValue(result.IsTruncated),
	}

	// 公共前缀（目录）
	for _, cp := range result.CommonPrefixes {
		dirPrefix := aws.StringValue(cp.Prefix)
		count, size := s.getDirStatsS3(s3Client, bucket, dirPrefix)
		resp.Files = append(resp.Files, clientReq.CloudFileItem{
			Key:          dirPrefix,
			Size:         0,
			LastModified: "",
			IsDir:        true,
			DirFileCount: count,
			DirSize:      size,
		})
	}

	// 文件列表
	for _, obj := range result.Contents {
		key := aws.StringValue(obj.Key)
		if prefix != "" && key == prefix {
			continue // 跳过目录自身的条目
		}
		lastMod := ""
		if obj.LastModified != nil {
			lastMod = obj.LastModified.Format("2006-01-02 15:04:05")
		}
		resp.Files = append(resp.Files, clientReq.CloudFileItem{
			Key:          key,
			Size:         aws.Int64Value(obj.Size),
			LastModified: lastMod,
			IsDir:        false,
		})
	}

	// 下一轮标记
	if resp.IsTruncated {
		if result.NextContinuationToken != nil {
			resp.Marker = aws.StringValue(result.NextContinuationToken)
		} else if len(resp.Files) > 0 {
			resp.Marker = resp.Files[len(resp.Files)-1].Key
		}
	}

	return resp, nil
}

// CompareDirectories 多云目录比对（以指定云为参考基准）
func (s *CloudStorageService) CompareDirectories(req clientReq.CloudCompareReq) (*clientReq.CloudCompareResp, error) {
	// 获取参考云
	refDomain, err := s.getCloudConfigByID(req.ReferenceID)
	if err != nil {
		return nil, fmt.Errorf("参考云不存在: %w", err)
	}

	prefix := strings.TrimSpace(req.Prefix)

	// 获取参考云文件列表
	refFiles, err := s.listAllFilesForDomain(refDomain, prefix)
	if err != nil {
		return nil, fmt.Errorf("获取参考云文件列表失败: %w", err)
	}

	// 获取其他已启用云
	allDomains, err := s.getCloudConfigs()
	if err != nil {
		return nil, err
	}

	resp := &clientReq.CloudCompareResp{
		ReferenceName: refDomain.Name,
	}

	for _, target := range allDomains {
		if target.ID == req.ReferenceID {
			continue
		}

		targetResult := clientReq.CloudCompareTarget{
			ID:        target.ID,
			Name:      target.Name,
			CloudType: target.CloudType,
			TotalRef:  len(refFiles),
		}

		targetFiles, err := s.listAllFilesForDomain(target, prefix)
		if err != nil {
			targetResult.Error = err.Error()
			resp.Targets = append(resp.Targets, targetResult)
			continue
		}

		targetResult.TotalTarget = len(targetFiles)

		// 构建目标文件集合用于快速查找
		targetSet := make(map[string]int64, len(targetFiles))
		for k, v := range targetFiles {
			targetSet[k] = v
		}

		// 遍历参考文件，找出目标缺失的
		for key, size := range refFiles {
			if _, ok := targetSet[key]; !ok {
				targetResult.Missing = append(targetResult.Missing, clientReq.CloudCompareFile{
					Key:  key,
					Size: size,
				})
			}
		}

		// 遍历目标文件，找出参考没有的（多余）
		for key, size := range targetFiles {
			if _, ok := refFiles[key]; !ok {
				targetResult.Extra = append(targetResult.Extra, clientReq.CloudCompareFile{
					Key:  key,
					Size: size,
				})
			}
		}

		resp.Targets = append(resp.Targets, targetResult)
	}

	return resp, nil
}

const maxCompareFiles = 100000 // 比对单云最多取文件数，超出截断

// listAllFilesForDomain 获取指定云下某前缀的所有文件（用于比对，仅返回文件不含目录）
func (s *CloudStorageService) listAllFilesForDomain(domain client.ExternalLinkDomain, prefix string) (map[string]int64, error) {
	cloudType := strings.TrimSpace(domain.CloudType)

	// 七牛云
	if cloudType == "qiniu" {
		return s.listAllQiniuFiles(domain, prefix)
	}

	// S3 兼容云
	s3Client, err := s.buildS3Client(domain)
	if err != nil {
		return nil, err
	}

	bucket := strings.TrimSpace(domain.Bucket)
	files := make(map[string]int64)

	input := &s3.ListObjectsV2Input{Bucket: aws.String(bucket), MaxKeys: aws.Int64(1000)}
	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	err = s3Client.ListObjectsV2PagesWithContext(ctx, input,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			for _, obj := range page.Contents {
				key := aws.StringValue(obj.Key)
				if strings.HasSuffix(key, "/") {
					continue
				}
				files[key] = aws.Int64Value(obj.Size)
				if len(files) >= maxCompareFiles {
					return false // 达到上限，停止翻页
				}
			}
			return true
		})

	return files, err
}

// listAllQiniuFiles 获取七牛云某前缀的所有文件
func (s *CloudStorageService) listAllQiniuFiles(domain client.ExternalLinkDomain, prefix string) (map[string]int64, error) {
	accessKey := strings.TrimSpace(domain.AccessKey)
	secretKey := strings.TrimSpace(domain.SecretKey)
	bucket := strings.TrimSpace(domain.Bucket)

	mac := qbox.NewMac(accessKey, secretKey)
	cfg := qiniuConfigForDomain(domain)
	bucketManager := qiniuStorage.NewBucketManager(mac, cfg)

	files := make(map[string]int64)
	marker := ""

	for {
		entries, _, nextMarker, hasNext, err := bucketManager.ListFiles(bucket, prefix, "", marker, 1000)

		if err != nil {
			return nil, fmt.Errorf("列出七牛文件失败: %w", err)
		}

		for _, entry := range entries {
			if strings.HasSuffix(entry.Key, "/") {
				continue
			}
			files[entry.Key] = entry.Fsize
			if len(files) >= maxCompareFiles {
				return files, nil // 达到上限，停止翻页
			}
		}

		if !hasNext {
			break
		}
		marker = nextMarker
	}

	return files, nil
}

// ==================== 上传功能 ====================

// GetDefaultUploadDomain 获取默认上传云（优先 defaultUpload，其次 isDefault）
func (s *CloudStorageService) GetDefaultUploadDomain() (*client.ExternalLinkDomain, error) {
	var domain client.ExternalLinkDomain

	// 优先找显式标记了 default_upload 的
	err := global.GVA_DB.Where("default_upload = ? AND is_enabled = ? AND cloud_type != ''", true, true).
		Order("sort DESC, id ASC").First(&domain).Error
	if err == nil {
		return &domain, nil
	}

	// 退而求其次，找 is_default 且启用云的
	err = global.GVA_DB.Where("is_default = ? AND is_enabled = ? AND cloud_type != ''", true, true).
		Order("sort DESC, id ASC").First(&domain).Error
	if err == nil {
		return &domain, nil
	}

	return nil, nil // 没有配置云端，调用方回退到 config.yaml
}

// GetSyncUploadDomains 获取所有启用了同步上传的云（排除默认上传云，避免重复）
func (s *CloudStorageService) GetSyncUploadDomains() ([]client.ExternalLinkDomain, error) {
	var domains []client.ExternalLinkDomain

	// 获取默认上传云ID，排除它
	defaultDomain, _ := s.GetDefaultUploadDomain()
	var excludeID uint
	if defaultDomain != nil {
		excludeID = defaultDomain.ID
	}

	query := global.GVA_DB.Where("sync_upload = ? AND is_enabled = ? AND cloud_type != ''", true, true)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Order("sort DESC, id ASC").Find(&domains).Error
	return domains, err
}

// UploadFileToCloud 上传文件到指定云存储，返回 (相对路径, 相对路径, error)
// 返回值中 filePath 和 key 均为相对于 bucket 根的文件路径，不含域名前缀
func (s *CloudStorageService) UploadFileToCloud(domain client.ExternalLinkDomain, file *multipart.FileHeader, folder string) (filePath string, key string, err error) {
	cloudType := strings.TrimSpace(domain.CloudType)

	if cloudType == "qiniu" {
		return s.uploadToQiniu(domain, file, folder)
	}

	// S3 兼容云（b2 / s3 / r2）
	return s.uploadToS3(domain, file, folder)
}

// uploadToS3 上传到 S3 兼容云
func (s *CloudStorageService) uploadToS3(domain client.ExternalLinkDomain, file *multipart.FileHeader, folder string) (string, string, error) {
	s3Client, err := s.buildS3Client(domain)
	if err != nil {
		return "", "", err
	}

	// 构造文件 key：folder/timestamp_filename.ext
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)
	var keyParts []string
	if folder != "" {
		keyParts = append(keyParts, folder)
	}
	keyParts = append(keyParts, fileKey)
	fullKey := strings.Join(keyParts, "/")

	contentType := strings.TrimSpace(file.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 读取文件内容到内存（管理后台文件通常不超过 50MB）
	f, openErr := file.Open()
	if openErr != nil {
		global.GVA_LOG.Error("上传云存储打开文件失败", zap.Error(openErr))
		return "", "", errors.New("打开文件失败: " + openErr.Error())
	}
	defer f.Close()

	data, readErr := io.ReadAll(f)
	if readErr != nil {
		global.GVA_LOG.Error("上传云存储读取文件失败", zap.Error(readErr))
		return "", "", errors.New("读取文件失败: " + readErr.Error())
	}
	bodyReader := bytes.NewReader(data)

	// 使用 PutObject（兼容性更好，避免 s3manager 多部分上传与 R2 的兼容问题）
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	_, uploadErr := s3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(strings.TrimSpace(domain.Bucket)),
		Key:         aws.String(fullKey),
		Body:        bodyReader,
		ContentType: aws.String(contentType),
	})

	if uploadErr != nil {
		global.GVA_LOG.Error("上传云存储失败",
			zap.String("endpoint", strings.TrimSpace(domain.RegionEndpoint)),
			zap.String("bucket", strings.TrimSpace(domain.Bucket)),
			zap.String("key", fullKey),
			zap.Error(uploadErr))
		return "", "", fmt.Errorf("上传云存储失败: %w", uploadErr)
	}

	global.GVA_LOG.Info("上传云存储成功",
		zap.String("bucket", strings.TrimSpace(domain.Bucket)),
		zap.String("key", fullKey),
		zap.Int64("size", file.Size))

	// 返回相对路径（不含域名），由 UNI 端用 getDefaultDomain 拼接
	return fullKey, fullKey, nil
}

// uploadToQiniu 上传到七牛云
func (s *CloudStorageService) uploadToQiniu(domain client.ExternalLinkDomain, file *multipart.FileHeader, folder string) (string, string, error) {
	accessKey := strings.TrimSpace(domain.AccessKey)
	secretKey := strings.TrimSpace(domain.SecretKey)
	bucket := strings.TrimSpace(domain.Bucket)

	if accessKey == "" || secretKey == "" {
		return "", "", errors.New("七牛云密钥未配置")
	}
	if bucket == "" {
		return "", "", errors.New("七牛云Bucket未配置")
	}

	putPolicy := qiniuStorage.PutPolicy{Scope: bucket}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := qiniuConfigForDomain(domain)
	formUploader := qiniuStorage.NewFormUploader(cfg)
	ret := qiniuStorage.PutRet{}

	f, openErr := file.Open()
	if openErr != nil {
		global.GVA_LOG.Error("上传七牛云打开文件失败", zap.Error(openErr))
		return "", "", errors.New("打开文件失败: " + openErr.Error())
	}
	defer f.Close()

	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)
	if folder != "" {
		fileKey = folder + "/" + fileKey
	}

	putErr := formUploader.Put(context.Background(), &ret, upToken, fileKey, f, file.Size, &qiniuStorage.PutExtra{})
	if putErr != nil {
		global.GVA_LOG.Error("上传七牛云失败", zap.Error(putErr))
		return "", "", errors.New("上传七牛云失败: " + putErr.Error())
	}

	// 返回相对路径
	return ret.Key, ret.Key, nil
}

// UploadBytesToCloud 将字节数据直接上传到指定云存储
// key 为完整的对象键（含文件夹前缀），如 "hls/123/index.m3u8"
// 返回 uploadedKey（与 key 一致）和可能的错误
func (s *CloudStorageService) UploadBytesToCloud(domain client.ExternalLinkDomain, data []byte, key string, contentType string) (string, error) {
	cloudType := strings.TrimSpace(domain.CloudType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	if cloudType == "qiniu" {
		return s.uploadBytesToQiniu(domain, data, key, contentType)
	}

	return s.uploadBytesToS3(domain, data, key, contentType)
}

func (s *CloudStorageService) uploadBytesToS3(domain client.ExternalLinkDomain, data []byte, key string, contentType string) (string, error) {
	s3Client, err := s.buildS3Client(domain)
	if err != nil {
		return "", err
	}

	bodyReader := bytes.NewReader(data)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	_, uploadErr := s3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(strings.TrimSpace(domain.Bucket)),
		Key:         aws.String(key),
		Body:        bodyReader,
		ContentType: aws.String(contentType),
	})

	if uploadErr != nil {
		global.GVA_LOG.Error("上传字节到S3云存储失败",
			zap.String("key", key),
			zap.Error(uploadErr))
		return "", fmt.Errorf("上传S3失败: %w", uploadErr)
	}

	return key, nil
}

func (s *CloudStorageService) uploadBytesToQiniu(domain client.ExternalLinkDomain, data []byte, key string, contentType string) (string, error) {
	accessKey := strings.TrimSpace(domain.AccessKey)
	secretKey := strings.TrimSpace(domain.SecretKey)
	bucket := strings.TrimSpace(domain.Bucket)

	if accessKey == "" || secretKey == "" {
		return "", errors.New("七牛云密钥未配置")
	}
	if bucket == "" {
		return "", errors.New("七牛云Bucket未配置")
	}

	putPolicy := qiniuStorage.PutPolicy{Scope: bucket}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := qiniuConfigForDomain(domain)
	formUploader := qiniuStorage.NewFormUploader(cfg)
	ret := qiniuStorage.PutRet{}

	bodyReader := bytes.NewReader(data)
	putErr := formUploader.Put(context.Background(), &ret, upToken, key, bodyReader, int64(len(data)), &qiniuStorage.PutExtra{})
	if putErr != nil {
		global.GVA_LOG.Error("上传字节到七牛云失败", zap.Error(putErr))
		return "", errors.New("上传七牛云失败: " + putErr.Error())
	}

	return ret.Key, nil
}

// DeleteCloudFiles 批量删除云存储文件
func (s *CloudStorageService) DeleteCloudFiles(req clientReq.CloudDeleteFilesReq) (*clientReq.CloudDeleteFilesResp, error) {
	// 校验日期密码
	today := time.Now().Format("20060102")
	if req.Password != today {
		return nil, errors.New("密码错误：请输入当天日期（YYYYMMDD格式，如" + today + "）")
	}

	if len(req.Keys) == 0 {
		return nil, errors.New("未指定要删除的文件")
	}

	domain, err := s.getCloudConfigByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("域名不存在: %w", err)
	}

	cloudType := strings.TrimSpace(domain.CloudType)

	if cloudType == "qiniu" {
		return s.deleteQiniuFiles(domain, req.Keys)
	}

	return s.deleteS3Files(domain, req.Keys)
}

func (s *CloudStorageService) deleteS3Files(domain client.ExternalLinkDomain, keys []string) (*clientReq.CloudDeleteFilesResp, error) {
	s3Client, err := s.buildS3Client(domain)
	if err != nil {
		return nil, err
	}

	bucket := strings.TrimSpace(domain.Bucket)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp := &clientReq.CloudDeleteFilesResp{}
	var failedKeys []string

	for _, key := range keys {
		_, delErr := s3Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if delErr != nil {
			global.GVA_LOG.Warn("删除S3文件失败",
				zap.String("key", key),
				zap.Error(delErr))
			failedKeys = append(failedKeys, key)
		} else {
			resp.DeletedCount++
		}
	}

	resp.FailedKeys = failedKeys
	return resp, nil
}

func (s *CloudStorageService) deleteQiniuFiles(domain client.ExternalLinkDomain, keys []string) (*clientReq.CloudDeleteFilesResp, error) {
	accessKey := strings.TrimSpace(domain.AccessKey)
	secretKey := strings.TrimSpace(domain.SecretKey)
	bucket := strings.TrimSpace(domain.Bucket)

	mac := qbox.NewMac(accessKey, secretKey)
	cfg := qiniuConfigForDomain(domain)
	bucketManager := qiniuStorage.NewBucketManager(mac, cfg)

	resp := &clientReq.CloudDeleteFilesResp{}
	var failedKeys []string

	for _, key := range keys {
		delErr := bucketManager.Delete(bucket, key)
		if delErr != nil {
			global.GVA_LOG.Warn("删除七牛云文件失败",
				zap.String("key", key),
				zap.Error(delErr))
			failedKeys = append(failedKeys, key)
		} else {
			resp.DeletedCount++
		}
	}

	resp.FailedKeys = failedKeys
	return resp, nil
}

// UploadFileToCloudByID 根据域名ID上传文件到云存储
func (s *CloudStorageService) UploadFileToCloudByID(id uint, file *multipart.FileHeader, folder string) (string, string, error) {
	domain, err := s.getCloudConfigByID(id)
	if err != nil {
		return "", "", fmt.Errorf("域名不存在: %w", err)
	}
	return s.UploadFileToCloud(domain, file, folder)
}

// ==================== 目录大小统计 ====================

// getDirStatsS3 统计 S3 兼容云某目录下的文件数和总大小（近似，最多取1000个）
func (s *CloudStorageService) getDirStatsS3(s3Client *s3.S3, bucket, dirPrefix string) (fileCount, totalSize int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		Prefix:  aws.String(dirPrefix),
		MaxKeys: aws.Int64(1000),
	}

	result, err := s3Client.ListObjectsV2WithContext(ctx, input)
	if err != nil {
		global.GVA_LOG.Warn("获取目录统计失败", zap.String("prefix", dirPrefix), zap.Error(err))
		return 0, 0
	}

	for _, obj := range result.Contents {
		key := aws.StringValue(obj.Key)
		if key == dirPrefix {
			continue
		}
		fileCount++
		totalSize += aws.Int64Value(obj.Size)
	}
	return
}

// getDirStatsQiniu 统计七牛云某目录下的文件数和总大小
func (s *CloudStorageService) getDirStatsQiniu(bucketManager *qiniuStorage.BucketManager, bucket, dirPrefix string) (fileCount, totalSize int64) {
	entries, _, _, _, err := bucketManager.ListFiles(bucket, dirPrefix, "", "", 1000)
	if err != nil {
		global.GVA_LOG.Warn("获取七牛目录统计失败", zap.String("prefix", dirPrefix), zap.Error(err))
		return 0, 0
	}

	for _, entry := range entries {
		if entry.Key == dirPrefix {
			continue
		}
		fileCount++
		totalSize += entry.Fsize
	}
	return
}

// ==================== 全局搜索 ====================

// SearchCloudFiles 全局搜索云存储文件（按文件名关键词）
func (s *CloudStorageService) SearchCloudFiles(req clientReq.SearchCloudFilesReq) (*clientReq.SearchCloudFilesResp, error) {
	domain, err := s.getCloudConfigByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("域名不存在: %w", err)
	}

	cloudType := strings.TrimSpace(domain.CloudType)
	if cloudType == "qiniu" {
		return s.searchQiniuFiles(domain, req)
	}
	return s.searchS3Files(domain, req)
}

func (s *CloudStorageService) searchS3Files(domain client.ExternalLinkDomain, req clientReq.SearchCloudFilesReq) (*clientReq.SearchCloudFilesResp, error) {
	s3Client, err := s.buildS3Client(domain)
	if err != nil {
		return nil, err
	}

	bucket := strings.TrimSpace(domain.Bucket)
	keyword := strings.ToLower(strings.TrimSpace(req.Keyword))
	maxKeys := int64(req.MaxKeys)
	if maxKeys <= 0 || maxKeys > 500 {
		maxKeys = 100
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp := &clientReq.SearchCloudFilesResp{}
	var matched []clientReq.CloudFileItem

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(1000),
	}

	err = s3Client.ListObjectsV2PagesWithContext(ctx, input,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			for _, obj := range page.Contents {
				key := aws.StringValue(obj.Key)
				if strings.HasSuffix(key, "/") {
					continue
				}
				if !strings.Contains(strings.ToLower(key), keyword) {
					continue
				}
				lastMod := ""
				if obj.LastModified != nil {
					lastMod = obj.LastModified.Format("2006-01-02 15:04:05")
				}
				matched = append(matched, clientReq.CloudFileItem{
					Key:          key,
					Size:         aws.Int64Value(obj.Size),
					LastModified: lastMod,
					IsDir:        false,
				})
				if int64(len(matched)) >= maxKeys {
					return false
				}
			}
			return int64(len(matched)) < maxKeys
		})

	resp.Files = matched
	if int64(len(matched)) >= maxKeys {
		resp.IsTruncated = true
	}

	return resp, err
}

func (s *CloudStorageService) searchQiniuFiles(domain client.ExternalLinkDomain, req clientReq.SearchCloudFilesReq) (*clientReq.SearchCloudFilesResp, error) {
	accessKey := strings.TrimSpace(domain.AccessKey)
	secretKey := strings.TrimSpace(domain.SecretKey)
	bucket := strings.TrimSpace(domain.Bucket)

	mac := qbox.NewMac(accessKey, secretKey)
	cfg := qiniuConfigForDomain(domain)
	bucketManager := qiniuStorage.NewBucketManager(mac, cfg)

	keyword := strings.ToLower(strings.TrimSpace(req.Keyword))
	maxKeys := req.MaxKeys
	if maxKeys <= 0 || maxKeys > 500 {
		maxKeys = 100
	}

	resp := &clientReq.SearchCloudFilesResp{}
	var matched []clientReq.CloudFileItem
	marker := ""

	for len(matched) < maxKeys {
		entries, _, nextMarker, hasNext, err := bucketManager.ListFiles(bucket, "", "", marker, 1000)
		if err != nil {
			return nil, fmt.Errorf("搜索七牛文件失败: %w", err)
		}

		for _, entry := range entries {
			if strings.HasSuffix(entry.Key, "/") {
				continue
			}
			if !strings.Contains(strings.ToLower(entry.Key), keyword) {
				continue
			}
			lastMod := time.Unix(entry.PutTime/10000000, 0).Format("2006-01-02 15:04:05")
			matched = append(matched, clientReq.CloudFileItem{
				Key:          entry.Key,
				Size:         entry.Fsize,
				LastModified: lastMod,
				IsDir:        false,
			})
			if len(matched) >= maxKeys {
				break
			}
		}

		if !hasNext || len(matched) >= maxKeys {
			break
		}
		marker = nextMarker
	}

	resp.Files = matched
	if len(matched) >= maxKeys {
		resp.IsTruncated = true
	}

	return resp, nil
}

// ==================== 下载功能 ====================

// GetFileDownloadURL 获取文件签名下载链接
func (s *CloudStorageService) GetFileDownloadURL(req clientReq.DownloadCloudFileReq) (*clientReq.DownloadCloudFileResp, error) {
	domain, err := s.getCloudConfigByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("域名不存在: %w", err)
	}

	cloudType := strings.TrimSpace(domain.CloudType)
	key := req.Key

	if cloudType == "qiniu" {
		return s.getQiniuDownloadURL(domain, key)
	}
	return s.getS3DownloadURL(domain, key)
}

func (s *CloudStorageService) getS3DownloadURL(domain client.ExternalLinkDomain, key string) (*clientReq.DownloadCloudFileResp, error) {
	s3Client, err := s.buildS3Client(domain)
	if err != nil {
		return nil, err
	}

	bucket := strings.TrimSpace(domain.Bucket)

	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	// 签名有效期 10 分钟
	url, err := req.Presign(10 * time.Minute)
	if err != nil {
		return nil, fmt.Errorf("生成签名URL失败: %w", err)
	}

	return &clientReq.DownloadCloudFileResp{DownloadURL: url}, nil
}

func (s *CloudStorageService) getQiniuDownloadURL(domain client.ExternalLinkDomain, key string) (*clientReq.DownloadCloudFileResp, error) {
	accessKey := strings.TrimSpace(domain.AccessKey)
	secretKey := strings.TrimSpace(domain.SecretKey)

	mac := qbox.NewMac(accessKey, secretKey)

	// 七牛：拼接域名 + key，如配置了域名则用域名，否则用默认
	baseURL := strings.TrimSpace(domain.Domain)
	if baseURL == "" {
		baseURL = strings.TrimSpace(domain.BaseURL)
	}
	if baseURL == "" {
		return nil, fmt.Errorf("七牛云未配置可访问域名")
	}
	baseURL = strings.TrimRight(baseURL, "/")

	deadline := time.Now().Add(10 * time.Minute).Unix()
	url := qiniuStorage.MakePrivateURL(mac, baseURL, key, deadline)

	return &clientReq.DownloadCloudFileResp{DownloadURL: url}, nil
}

// ZipCloudFolder 将云存储目录打包为 zip 流
// 返回 zip 读取器、建议文件名、错误
func (s *CloudStorageService) ZipCloudFolder(domain client.ExternalLinkDomain, prefix string) (io.ReadCloser, string, error) {
	cloudType := strings.TrimSpace(domain.CloudType)

	// 获取目录下所有文件的 key 和内容读取方法
	type fileRef struct {
		key  string
		size int64
	}

	var files []fileRef

	if cloudType == "qiniu" {
		list, err := s.listAllQiniuFiles(domain, prefix)
		if err != nil {
			return nil, "", err
		}
		for k, sz := range list {
			files = append(files, fileRef{key: k, size: sz})
		}
	} else {
		s3Client, err := s.buildS3Client(domain)
		if err != nil {
			return nil, "", err
		}
		bucket := strings.TrimSpace(domain.Bucket)
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		input := &s3.ListObjectsV2Input{
			Bucket:  aws.String(bucket),
			Prefix:  aws.String(prefix),
			MaxKeys: aws.Int64(1000),
		}

		err = s3Client.ListObjectsV2PagesWithContext(ctx, input,
			func(page *s3.ListObjectsV2Output, lastPage bool) bool {
				for _, obj := range page.Contents {
					key := aws.StringValue(obj.Key)
					if strings.HasSuffix(key, "/") || key == prefix {
						continue
					}
					files = append(files, fileRef{key: key, size: aws.Int64Value(obj.Size)})
					if len(files) >= 1000 {
						return false
					}
				}
				return len(files) < 1000
			})
		if err != nil {
			return nil, "", fmt.Errorf("列出目录文件失败: %w", err)
		}
	}

	if len(files) == 0 {
		return nil, "", errors.New("目录为空")
	}

	// 确定 zip 文件名
	folderName := strings.TrimRight(prefix, "/")
	if idx := strings.LastIndex(folderName, "/"); idx >= 0 {
		folderName = folderName[idx+1:]
	}
	if folderName == "" {
		folderName = "root"
	}
	zipFileName := folderName + ".zip"

	// 创建管道
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		zw := zip.NewWriter(pw)
		defer zw.Close()

		for _, f := range files {
			var content io.Reader

			if cloudType == "qiniu" {
				// 七牛云需要通过公开URL或签名URL下载
				accessKey := strings.TrimSpace(domain.AccessKey)
				secretKey := strings.TrimSpace(domain.SecretKey)
				mac := qbox.NewMac(accessKey, secretKey)
				baseURL := strings.TrimSpace(domain.Domain)
				if baseURL == "" {
					baseURL = strings.TrimSpace(domain.BaseURL)
				}
				baseURL = strings.TrimRight(baseURL, "/")
				deadline := time.Now().Add(10 * time.Minute).Unix()
				url := qiniuStorage.MakePrivateURL(mac, baseURL, f.key, deadline)
				
				httpClient := &http.Client{Timeout: 30 * time.Second}
				resp, httpErr := httpClient.Get(url)
				if httpErr != nil {
					global.GVA_LOG.Warn("zip下载七牛文件失败", zap.String("key", f.key), zap.Error(httpErr))
					continue
				}
				defer resp.Body.Close()
				content = resp.Body
			} else {
				s3Client, err := s.buildS3Client(domain)
				if err != nil {
					global.GVA_LOG.Warn("zip连接S3失败", zap.Error(err))
					continue
				}
				bucket := strings.TrimSpace(domain.Bucket)
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				out, err := s3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
					Bucket: aws.String(bucket),
					Key:    aws.String(f.key),
				})
				cancel()
				if err != nil {
					global.GVA_LOG.Warn("zip下载S3文件失败", zap.String("key", f.key), zap.Error(err))
					continue
				}
				content = out.Body
				defer out.Body.Close()
			}

			// zip 内的相对路径（去掉前缀目录）
			innerPath := strings.TrimPrefix(f.key, prefix)

			w, err := zw.Create(innerPath)
			if err != nil {
				global.GVA_LOG.Warn("zip创建条目失败", zap.String("key", f.key), zap.Error(err))
				continue
			}

			_, copyErr := io.Copy(w, content)
			if copyErr != nil {
				global.GVA_LOG.Warn("zip写入文件失败", zap.String("key", f.key), zap.Error(copyErr))
			}
		}
	}()

	return pr, zipFileName, nil
}

// ZipCloudFolderByID 根据域名ID打包目录为 zip
func (s *CloudStorageService) ZipCloudFolderByID(id uint, prefix string) (io.ReadCloser, string, error) {
	domain, err := s.getCloudConfigByID(id)
	if err != nil {
		return nil, "", fmt.Errorf("域名不存在: %w", err)
	}
	return s.ZipCloudFolder(domain, prefix)
}
