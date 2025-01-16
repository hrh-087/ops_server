package config

type Asynq struct {
	Addr          string `mapstructure:"addr" json:"addr" yaml:"addr"`                      // redis连接地址
	Port          string `mapstructure:"port" json:"port" yaml:"port"`                      // redis端口
	Db            int    `mapstructure:"db" json:"db" yaml:"db"`                            // redis数据库
	Pass          string `mapstructure:"password" json:"password" yaml:"password"`          // redis密码
	MaxRetryCount int    `mapstructure:"max-retry" json:"max-retry" yaml:"max-retry"`       // 重试次数
	Retention     int64  `mapstructure:"retention" json:"retention" yaml:"retention"`       // 缓存时间/小时
	Concurrency   int    `mapstructure:"concurrency" json:"concurrency" yaml:"concurrency"` // 并发数量
	Timeout       int64  `mapstructure:"timeout" json:"timeout" yaml:"timeout"`             // 任务执行超时时间/秒
}
