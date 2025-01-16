package initialize

import (
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/initialize/internal"
	"ops-server/plugin"
)

func GormMysql() *gorm.DB {
	m := global.OPS_CONFIG.Mysql

	if m.Dbname == "" {
		return nil
	}

	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         255,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), internal.Gorm.Config(m.Prefix, m.Singular)); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		// 注册gorm插件
		err = db.Use(plugin.ProjectFilterPlugin{})
		if err != nil {
			global.OPS_LOG.Error("加载gorm插件失败!", zap.Error(err))
			return nil
		}
		return db
	}
}
