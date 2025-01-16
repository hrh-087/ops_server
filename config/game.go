package config

type Game struct {
	GamePath           string `mapstructure:"game-path" json:"game-path" yaml:"game-path"`
	GameScriptPath     string `mapstructure:"game-script-path" json:"game-script-path" yaml:"game-script-path"`
	GameScriptAutoPath string `mapstructure:"game-script-auto-path" json:"game-script-auto-path" yaml:"game-script-auto-path"`
}
