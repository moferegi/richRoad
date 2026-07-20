package client

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ExternalLinkDomain 外部链接域名管理（支持云存储配置）
type ExternalLinkDomain struct {
	global.GVA_MODEL
	Name         string  `json:"name" form:"name" gorm:"column:name;size:100;comment:域名名称;"`
	Domain       string  `json:"domain" form:"domain" gorm:"column:domain;size:500;comment:域名地址(如https://cdn.example.com);"`
	IsDefault    *bool   `json:"isDefault" form:"isDefault" gorm:"column:is_default;default:false;comment:是否默认;"`
	IsEnabled    *bool   `json:"isEnabled" form:"isEnabled" gorm:"column:is_enabled;default:true;comment:是否启用;"`
	Sort         int     `json:"sort" form:"sort" gorm:"column:sort;default:0;comment:排序(越大越前);"`
	Remark       string  `json:"remark" form:"remark" gorm:"column:remark;size:500;comment:备注;"`

	// 云存储通用配置
	CloudType       string  `json:"cloudType" form:"cloudType" gorm:"column:cloud_type;size:50;default:'';comment:云类型(qiniu/r2/b2/s3);"`
	AccessKey       string  `json:"accessKey" form:"accessKey" gorm:"column:access_key;size:500;comment:AccessKey/APIKey;"`
	SecretKey       string  `json:"secretKey" form:"secretKey" gorm:"column:secret_key;size:500;comment:SecretKey;"`
	Bucket          string  `json:"bucket" form:"bucket" gorm:"column:bucket;size:200;comment:存储桶名称;"`
	RegionEndpoint  string  `json:"regionEndpoint" form:"regionEndpoint" gorm:"column:region_endpoint;size:500;comment:区域/Endpoint;"`

	// 云存储专属字段
	Zone        string `json:"zone" form:"zone" gorm:"column:zone;size:50;comment:存储区域(七牛专用);"`
	BaseURL     string `json:"baseUrl" form:"baseUrl" gorm:"column:base_url;size:500;comment:CDN/自定义域名;"`
	AccountId   string `json:"accountId" form:"accountId" gorm:"column:account_id;size:200;comment:账户ID(R2专用);"`
	UseHTTPS    *bool  `json:"useHttps" form:"useHttps" gorm:"column:use_https;default:true;comment:是否使用HTTPS;"`
	UseCdnDomain *bool `json:"useCdnDomain" form:"useCdnDomain" gorm:"column:use_cdn_domain;default:false;comment:是否使用CDN域名(七牛专用);"`

	// 上传策略
	UploadScope    string `json:"uploadScope" form:"uploadScope" gorm:"column:upload_scope;size:50;default:'all';comment:上传用途(video/image/audio/all);"`
	DefaultUpload  *bool  `json:"defaultUpload" form:"defaultUpload" gorm:"column:default_upload;default:false;comment:是否默认上传云;"`
	SyncUpload     *bool  `json:"syncUpload" form:"syncUpload" gorm:"column:sync_upload;default:false;comment:是否同步上传至该云;"`
}

func (ExternalLinkDomain) TableName() string {
	return "client_external_link_domain"
}
