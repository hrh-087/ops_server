package config

type Server struct {
	Mysql   Mysql         `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	System  System        `mapstructure:"system" json:"system" yaml:"system"`
	Zap     Zap           `mapstructure:"zap" json:"zap" yaml:"zap"`
	JWT     JWT           `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Asynq   Asynq         `mapstructure:"asynq" json:"asynq" yaml:"asynq"`
	Game    Game          `mapstructure:"game" json:"game" yaml:"game"`
	Local   Local         `mapstructyre:"local" json:"local" yaml:"local"`
	Ops     Ops           `mapstructure:"ops" json:"ops" yaml:"ops"`
	Redis   Redis         `mapstructure:"redis" json:"redis" yaml:"redis"`
	Default DefaultConfig `mapstructure:"default" json:"default" yaml:"default"`
}
