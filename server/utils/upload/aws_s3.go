package upload

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
)

type AwsS3 struct{}

const (
	awsS3UploadMaxRetries = 3
	awsS3RetryBaseDelay   = 250 * time.Millisecond
)

//@author: [WqyJh](https://github.com/WqyJh)
//@object: *AwsS3
//@function: UploadFile
//@description: Upload file to Aws S3 using aws-sdk-go. See https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html#s3-examples-bucket-ops-upload-file-to-bucket
//@param: file *multipart.FileHeader
//@return: string, string, error

func (*AwsS3) UploadFile(file *multipart.FileHeader) (string, string, error) {
	return new(AwsS3).UploadFileToFolder(file, "")
}

// UploadFileToFolder 上传文件到指定文件夹
func (*AwsS3) UploadFileToFolder(file *multipart.FileHeader, folder string) (string, string, error) {
	sess := newSession()
	uploader := s3manager.NewUploader(sess)

	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)
	// 构建路径: PathPrefix/folder/fileKey
	pathParts := []string{}
	if global.GVA_CONFIG.AwsS3.PathPrefix != "" {
		pathParts = append(pathParts, global.GVA_CONFIG.AwsS3.PathPrefix)
	}
	if folder != "" {
		pathParts = append(pathParts, folder)
	}
	pathParts = append(pathParts, fileKey)

	filename := strings.Join(pathParts, "/")

	contentType := strings.TrimSpace(file.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	var lastErr error
	for attempt := 1; attempt <= awsS3UploadMaxRetries; attempt++ {
		f, openError := file.Open()
		if openError != nil {
			global.GVA_LOG.Error("function file.Open() failed", zap.Any("err", openError.Error()))
			return "", "", errors.New("function file.Open() failed, err:" + openError.Error())
		}

		_, err := uploader.Upload(&s3manager.UploadInput{
			Bucket:      aws.String(global.GVA_CONFIG.AwsS3.Bucket),
			Key:         aws.String(filename),
			Body:        f,
			ContentType: aws.String(contentType),
		})
		_ = f.Close()

		if err == nil {
			return global.GVA_CONFIG.AwsS3.BaseURL + "/" + filename, filename, nil
		}

		lastErr = err
		if attempt >= awsS3UploadMaxRetries || !shouldRetryAwsS3Upload(err) {
			break
		}

		backoff := time.Duration(attempt) * awsS3RetryBaseDelay
		global.GVA_LOG.Warn("aws s3 upload retry",
			zap.Int("attempt", attempt),
			zap.String("bucket", global.GVA_CONFIG.AwsS3.Bucket),
			zap.String("key", filename),
			zap.Duration("backoff", backoff),
			zap.Error(err),
		)
		time.Sleep(backoff)
	}

	if lastErr == nil {
		lastErr = errors.New("unknown s3 upload error")
	}
	global.GVA_LOG.Error("function uploader.Upload() failed", zap.Any("err", lastErr.Error()))
	return "", "", lastErr
}

func shouldRetryAwsS3Upload(err error) bool {
	if err == nil {
		return false
	}

	message := strings.ToLower(strings.TrimSpace(err.Error()))
	if message == "" {
		return false
	}

	transientHints := []string{
		"eof",
		"timeout",
		"temporarily unavailable",
		"connection reset",
		"connection aborted",
		"broken pipe",
		"tls handshake timeout",
		"request canceled",
		"http2: server sent goaway",
		"503",
	}

	for _, hint := range transientHints {
		if strings.Contains(message, hint) {
			return true
		}
	}

	return false
}

//@author: [WqyJh](https://github.com/WqyJh)
//@object: *AwsS3
//@function: DeleteFile
//@description: Delete file from Aws S3 using aws-sdk-go. See https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html#s3-examples-bucket-ops-delete-bucket-item
//@param: file *multipart.FileHeader
//@return: string, string, error

func (*AwsS3) DeleteFile(key string) error {
	session := newSession()
	svc := s3.New(session)
	filename := global.GVA_CONFIG.AwsS3.PathPrefix + "/" + key
	bucket := global.GVA_CONFIG.AwsS3.Bucket

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		global.GVA_LOG.Error("function svc.DeleteObject() failed", zap.Any("err", err.Error()))
		return errors.New("function svc.DeleteObject() failed, err:" + err.Error())
	}

	_ = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	return nil
}

// ListFolders 列举 bucket 中所有层级的文件夹路径（基于 path-prefix），支持分页
func (*AwsS3) ListFolders() ([]string, error) {
	sess := newSession()
	svc := s3.New(sess)

	basePrefix := ""
	if global.GVA_CONFIG.AwsS3.PathPrefix != "" {
		basePrefix = strings.TrimRight(global.GVA_CONFIG.AwsS3.PathPrefix, "/") + "/"
	}

	folderSet := make(map[string]struct{})
	var continuationToken *string

	for {
		input := &s3.ListObjectsV2Input{
			Bucket:  aws.String(global.GVA_CONFIG.AwsS3.Bucket),
			Prefix:  aws.String(basePrefix),
			MaxKeys: aws.Int64(1000),
		}
		if continuationToken != nil {
			input.ContinuationToken = continuationToken
		}

		result, err := svc.ListObjectsV2(input)
		if err != nil {
			return nil, err
		}

		for _, obj := range result.Contents {
			if obj.Key == nil {
				continue
			}
			// 去掉 basePrefix，取相对路径
			relKey := strings.TrimPrefix(*obj.Key, basePrefix)
			parts := strings.Split(relKey, "/")
			// 提取所有上级目录路径（不含文件名本身）
			for i := 1; i < len(parts); i++ {
				folderSet[strings.Join(parts[:i], "/")] = struct{}{}
			}
		}

		if result.IsTruncated == nil || !*result.IsTruncated {
			break
		}
		continuationToken = result.NextContinuationToken
	}

	folders := make([]string, 0, len(folderSet))
	for f := range folderSet {
		folders = append(folders, f)
	}
	sort.Strings(folders)
	return folders, nil
}

// newSession Create S3 session
func newSession() *session.Session {
	httpClient := &http.Client{
		Timeout: 90 * time.Second,
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			ForceAttemptHTTP2:     false,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   32,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
		},
	}

	sess, _ := session.NewSession(&aws.Config{
		Region:           aws.String(global.GVA_CONFIG.AwsS3.Region),
		Endpoint:         aws.String(global.GVA_CONFIG.AwsS3.Endpoint), //minio在这里设置地址,可以兼容
		S3ForcePathStyle: aws.Bool(global.GVA_CONFIG.AwsS3.S3ForcePathStyle),
		DisableSSL:       aws.Bool(global.GVA_CONFIG.AwsS3.DisableSSL),
		MaxRetries:       aws.Int(2),
		HTTPClient:       httpClient,
		Credentials: credentials.NewStaticCredentials(
			global.GVA_CONFIG.AwsS3.SecretID,
			global.GVA_CONFIG.AwsS3.SecretKey,
			"",
		),
	})
	return sess
}
