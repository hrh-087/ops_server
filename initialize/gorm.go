package initialize

import (
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/model/system"
	"os"
)

func Gorm() *gorm.DB {
	switch global.OPS_CONFIG.System.DbType {
	case "mysql":
		return GormMysql()
	default:
		return GormMysql()
	}
}

func RegisterTables() {
	db := global.OPS_DB

	err := db.AutoMigrate(
		system.SysUser{},         // 用户表
		system.SysAuthority{},    // 角色表
		system.SysAuthorityBtn{}, // 角色按钮表

		system.SysBaseMenu{},          // 基础菜单表
		system.SysBaseMenuParameter{}, // 基础菜单参数表
		system.SysBaseMenuBtn{},       // 基础按钮表

		system.SysApi{},       // api表
		system.SysIgnoreApi{}, // 忽略的api

		system.SysProject{},      // 项目表
		system.SysSshAuth{},      // ssh秘钥管理
		system.SysCloudProduce{}, // 云平台
		system.SysGamePlatform{}, // 游戏渠道

		system.SysAssetsServer{}, // 资产服务器
		system.SysAssetsServerPort{},
		system.SysAssetsRedis{},
		system.SysAssetsMongo{},
		system.SysAssetsMysql{},
		system.SysAssetsKafka{},
		system.SysAssetsLb{},
		system.SysAssetsListener{},

		// game
		system.SysGameServer{},
		system.SysGameType{},

		// job
		system.Job{},
		system.JobTask{},
		system.JobCommand{},
		system.GameUpdate{},

		system.SysOperationRecord{},
	)
	if err != nil {
		panic("register table failed")
	}

	err = db.AutoMigrate()

	if err != nil {
		os.Exit(0)
	}
}
