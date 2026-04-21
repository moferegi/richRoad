package config

type Hotlink struct {
	Enabled       bool   `mapstructure:"enabled" json:"enabled" yaml:"enabled"`                      // 防盗链开关
	CdnDomain     string `mapstructure:"cdn-domain" json:"cdn-domain" yaml:"cdn-domain"`             // CDN访问域名(如 https://cdn.yourdomain.com)
	SignKey       string `mapstructure:"sign-key" json:"sign-key" yaml:"sign-key"`                   // 签名密钥
	ExpireSeconds int64  `mapstructure:"expire-seconds" json:"expire-seconds" yaml:"expire-seconds"` // 签名URL过期时间(秒)
}
