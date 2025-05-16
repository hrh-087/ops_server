package config

type Prometheus struct {
	Addr              string `mapstructure:"addr" json:"addr" yaml:"addr"`
	SshPort           string `mapstructure:"ssh-port" json:"ssh-port" yaml:"ssh-port"`
	GameServerJsonDir string `mapstructure:"game-server-json-dir" json:"game-server-json-dir" yaml:"game-server-json-dir"`
	HostServerJsonDir string `mapstructure:"host-server-json-dir" json:"host-server-json-dir" yaml:"host-server-json-dir"`
	NodeExporterPort  string `mapstructure:"node-exporter-port" json:"node-exporter-port" yaml:"node-exporter-port"`
}
