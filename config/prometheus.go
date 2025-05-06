package config

type Prometheus struct {
	Addr              string `mapstructure:"addr" json:"addr" yaml:"addr"`
	GameServerJsonDir string `mapstructure:"game-server-json-dir" json:"game-server-json-dir" yaml:"game-server-json-dir"`
}
