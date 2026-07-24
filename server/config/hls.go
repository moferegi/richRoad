package config

type Hls struct {
	GlobalKey            string `mapstructure:"global-key" json:"global-key" yaml:"global-key"`                         // AES-128 加密密钥(16字节hex)
	HlsKeyRatePerMinute  int    `mapstructure:"hls-key-rate-per-minute" json:"hls-key-rate-per-minute" yaml:"hls-key-rate-per-minute"` // /hlsKey 每IP每分钟限流次数，0=不限
}
