package config

type Wxpay struct {
	MchID                      string `mapstructure:"mch-id" json:"mch-id" yaml:"mch-id"`
	AppID                      string `mapstructure:"app-id" json:"app-id" yaml:"app-id"`
	Secret                     string `mapstructure:"secret" json:"secret" yaml:"secret"`
	MchCertificateSerialNumber string `mapstructure:"mch-certificate-serial-number" json:"mch-certificate-serial-number" yaml:"mch-certificate-serial-number"`
	MchAPIv3Key                string `mapstructure:"mch-api-v3-key" json:"mch-api-v3-key" yaml:"mch-api-v3-key"`
	CertPath                   string `mapstructure:"cert-path" json:"cert-path" yaml:"cert-path"`
	KeyPath                    string `mapstructure:"key-path" json:"key-path" yaml:"key-path"`
	NotifyUrl                  string `mapstructure:"notify-url" json:"notify-url" yaml:"notify-url"`
}
