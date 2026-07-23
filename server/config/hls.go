package config

type Hls struct {
	GlobalKey  string `mapstructure:"global-key" json:"global-key" yaml:"global-key"`       // AES-128 加密密钥(16字节hex)
	BackendURL string `mapstructure:"backend-url" json:"backend-url" yaml:"backend-url"` // 已废弃：前端动态检测后端地址，不再需要此字段
}
