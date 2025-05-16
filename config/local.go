package config

type Local struct {
	Path      string `mapstructure:"path" json:"path" yaml:"path"`
	JsonDir   string `mapstructure:"json-dir" json:"json-dir" yaml:"json-dir"`
	StorePath string `mapstructure:"store-path" json:"store-path" yaml:"store-path"`
}
