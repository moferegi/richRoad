package config

type Sf struct {
	PartnerID     string `mapstructure:"partner-id" json:"partner-id" yaml:"partner-id"`                // 顾客编码
	CheckCode     string `mapstructure:"check-code" json:"check-code" yaml:"check-code"`                // 顾客校验码
	SandCheckCode string `mapstructure:"sand-check-code" json:"sand-check-code" yaml:"sand-check-code"` // 沙箱顾客校验码
	IsSandbox     bool   `mapstructure:"is-sandbox" json:"is-sandbox" yaml:"is-sandbox"`                // 是否是沙箱环境
}
