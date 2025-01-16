package config

type Mysql struct {
	GeneralDB `yaml:",inline" mapstructure:",squash"`
}

// Dsn 拼接连接字符串
func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}
