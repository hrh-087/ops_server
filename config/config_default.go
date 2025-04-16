package config

type DefaultConfig struct {
	GmUrl       string `mapstructure:"gm-url" json:"gm-url" yaml:"gm-url"`
	OnlineGmUrl string `mapstructure:"online-gm-url" json:"online-gm-url" yaml:"online-gm-url"`
}
