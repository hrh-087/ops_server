package internal

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"ops-server/config"
	"ops-server/global"
	"os"
	"time"
)

type _gorm struct{}

var Gorm = new(_gorm)

func (g *_gorm) Config(prefix string, singular bool) *gorm.Config {
	var general config.GeneralDB

	switch global.OPS_CONFIG.System.DbType {
	case "mysql":
		general = global.OPS_CONFIG.Mysql.GeneralDB
	default:
		general = global.OPS_CONFIG.Mysql.GeneralDB
	}

	return &gorm.Config{
		Logger: logger.New(NewWriter(general, log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      general.LogLevel(),
			Colorful:      true,
		}),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
			SingularTable: singular,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}

}
