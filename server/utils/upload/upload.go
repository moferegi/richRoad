package upload

import (
	"mime/multipart"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// OSS 对象存储接口
// Author [SliverHorn](https://github.com/SliverHorn)
// Author [ccfish86](https://github.com/ccfish86)
type OSS interface {
	UploadFile(file *multipart.FileHeader) (string, string, error)
	DeleteFile(key string) error
}

// OSSWithFolder 支持文件夹上传的扩展接口
type OSSWithFolder interface {
	OSS
	UploadFileToFolder(file *multipart.FileHeader, folder string) (string, string, error)
}

// OSSWithFolderList 支持列举文件夹的扩展接口
type OSSWithFolderList interface {
	OSS
	ListFolders() ([]string, error)
}

// ListOSSFolders 列举当前 OSS 的文件夹（如 OSS 不支持则返回空切片）
func ListOSSFolders(oss OSS) ([]string, error) {
	if f, ok := oss.(OSSWithFolderList); ok {
		return f.ListFolders()
	}
	return []string{}, nil
}

// UploadFileToFolder 上传文件到指定文件夹（兼容方法）
// 如果 OSS 实现了 OSSWithFolder 接口则使用它，否则回退到普通上传
func UploadFileToFolder(oss OSS, file *multipart.FileHeader, folder string) (string, string, error) {
	if f, ok := oss.(OSSWithFolder); ok && folder != "" {
		return f.UploadFileToFolder(file, folder)
	}
	return oss.UploadFile(file)
}

// NewOss OSS的实例化方法
// Author [SliverHorn](https://github.com/SliverHorn)
// Author [ccfish86](https://github.com/ccfish86)
func NewOss() OSS {
	switch global.GVA_CONFIG.System.OssType {
	case "local":
		return &Local{}
	case "qiniu":
		return &Qiniu{}
	case "tencent-cos":
		return &TencentCOS{}
	case "aliyun-oss":
		return &AliyunOSS{}
	case "huawei-obs":
		return HuaWeiObs
	case "aws-s3":
		return &AwsS3{}
	case "cloudflare-r2":
		return &CloudflareR2{}
	case "minio":
		minioClient, err := GetMinio(global.GVA_CONFIG.Minio.Endpoint, global.GVA_CONFIG.Minio.AccessKeyId, global.GVA_CONFIG.Minio.AccessKeySecret, global.GVA_CONFIG.Minio.BucketName, global.GVA_CONFIG.Minio.UseSSL)
		if err != nil {
			global.GVA_LOG.Warn("你配置了使用minio，但是初始化失败，请检查minio可用性或安全配置: " + err.Error())
			panic("minio初始化失败") // 建议这样做，用户自己配置了minio，如果报错了还要把服务开起来，使用起来也很危险
		}
		return minioClient
	default:
		return &Local{}
	}
}
