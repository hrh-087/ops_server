package config

type Ops struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port string `mapstructure:"port" json:"port" yaml:"port"`
	User string `mapstructure:"user" json:"user" yaml:"user"`
	Name string `mapstructure:"name" json:"name" yaml:"name"`
}
